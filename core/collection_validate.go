package core

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core/validators"
	"github.com/pocketbase/pocketbase/tools/dbutils"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/search"
	"github.com/pocketbase/pocketbase/tools/types"
)

var collectionNameRegex = regexp.MustCompile(`^\w+$`)

func onCollectionValidate(e *CollectionEvent) error {
	var original *Collection
	if !e.Collection.IsNew() {
		original = &Collection{}
		if err := e.App.ModelQuery(original).Model(e.Collection.LastSavedPK(), original); err != nil {
			return fmt.Errorf("failed to fetch old collection state: %w", err)
		}
	}

	validator := newCollectionValidator(
		e.Context,
		e.App,
		e.Collection,
		original,
	)

	if err := validator.run(); err != nil {
		return err
	}

	return e.Next()
}

func newCollectionValidator(ctx context.Context, app App, new, original *Collection) *collectionValidator {
	validator := &collectionValidator{
		ctx:      ctx,
		app:      app,
		new:      new,
		original: original,
	}

	// load old/original collection
	if validator.original == nil {
		validator.original = NewCollection(validator.new.Type, "")
	}

	return validator
}

type collectionValidator struct {
	original *Collection
	new      *Collection
	app      App
	ctx      context.Context
}

type optionsValidator interface {
	validate(cv *collectionValidator) error
}

func (validator *collectionValidator) run() error {
	if validator.original.IsNew() {
		validator.new.updateGeneratedIdIfExists(validator.app)
	}

	// generate fields from the query (overwriting any explicit user defined fields)
	if validator.new.IsView() {
		validator.new.Fields, _ = validator.app.CreateViewFields(validator.new.ViewQuery)
	}

	// validate base fields
	baseErr := validation.ValidateStruct(validator.new,
		validation.Field(
			&validator.new.Id,
			validation.Required,
			validation.When(
				validator.original.IsNew(),
				validation.Length(1, 100),
				validation.Match(DefaultIdRegex),
				validation.By(validators.UniqueId(validator.app.DB(), validator.new.TableName())),
			).Else(
				validation.By(validators.Equal(validator.original.Id)),
			),
		),
		validation.Field(
			&validator.new.System,
			validation.By(validator.ensureNoSystemFlagChange),
		),
		validation.Field(
			&validator.new.Type,
			validation.Required,
			validation.In(
				CollectionTypeBase,
				CollectionTypeAuth,
				CollectionTypeView,
			),
			validation.By(validator.ensureNoTypeChange),
		),
		validation.Field(
			&validator.new.Name,
			validation.Required,
			validation.Length(1, 255),
			validation.By(checkForVia),
			validation.Match(collectionNameRegex),
			validation.By(validator.ensureNoSystemNameChange),
			validation.By(validator.checkUniqueName),
		),
		validation.Field(
			&validator.new.Fields,
			validation.By(validator.checkFieldDuplicates),
			validation.By(validator.checkMinFields),
			validation.When(
				!validator.new.IsView(),
				validation.By(validator.ensureNoSystemFieldsChange),
				validation.By(validator.ensureNoFieldsTypeChange),
			),
			validation.When(validator.new.IsAuth(), validation.By(validator.checkReservedAuthKeys)),
			validation.By(validator.checkFieldValidators),
		),
		validation.Field(
			&validator.new.ListRule,
			validation.By(validator.checkRule),
			validation.By(validator.ensureNoSystemRuleChange(validator.original.ListRule)),
		),
		validation.Field(
			&validator.new.ViewRule,
			validation.By(validator.checkRule),
			validation.By(validator.ensureNoSystemRuleChange(validator.original.ViewRule)),
		),
		validation.Field(
			&validator.new.CreateRule,
			validation.When(validator.new.IsView(), validation.Nil),
			validation.By(validator.checkRule),
			validation.By(validator.ensureNoSystemRuleChange(validator.original.CreateRule)),
		),
		validation.Field(
			&validator.new.UpdateRule,
			validation.When(validator.new.IsView(), validation.Nil),
			validation.By(validator.checkRule),
			validation.By(validator.ensureNoSystemRuleChange(validator.original.UpdateRule)),
		),
		validation.Field(
			&validator.new.DeleteRule,
			validation.When(validator.new.IsView(), validation.Nil),
			validation.By(validator.checkRule),
			validation.By(validator.ensureNoSystemRuleChange(validator.original.DeleteRule)),
		),
		validation.Field(&validator.new.Indexes, validation.By(validator.checkIndexes)),
	)

	optionsErr := validator.validateOptions()

	return validators.JoinValidationErrors(baseErr, optionsErr)
}

func (validator *collectionValidator) checkUniqueName(value any) error {
	v, _ := value.(string)

	// ensure unique collection name
	if !validator.app.IsCollectionNameUnique(v, validator.original.Id) {
		return validation.NewError("validation_collection_name_exists", "Collection name must be unique (case insensitive).")
	}

	// ensure that the collection name doesn't collide with the id of any collection
	dummyCollection := &Collection{}
	if validator.app.ModelQuery(dummyCollection).Model(v, dummyCollection) == nil {
		return validation.NewError("validation_collection_name_id_duplicate", "The name must not match an existing collection id.")
	}

	// ensure that there is no existing internal table with the provided name
	if validator.original.Name != v && // has changed
		validator.app.IsCollectionNameUnique(v) && // is not a collection (in case it was presaved)
		validator.app.HasTable(v) {
		return validation.NewError("validation_collection_name_invalid", "The name shouldn't match with an existing internal table.")
	}

	return nil
}

func (validator *collectionValidator) ensureNoSystemNameChange(value any) error {
	v, _ := value.(string)

	if !validator.original.IsNew() && validator.original.System && v != validator.original.Name {
		return validation.NewError("validation_collection_system_name_change", "System collection name cannot be changed.")
	}

	return nil
}

func (validator *collectionValidator) ensureNoSystemFlagChange(value any) error {
	v, _ := value.(bool)

	if !validator.original.IsNew() && v != validator.original.System {
		return validation.NewError("validation_collection_system_flag_change", "System collection state cannot be changed.")
	}

	return nil
}

func (validator *collectionValidator) ensureNoTypeChange(value any) error {
	v, _ := value.(string)

	if !validator.original.IsNew() && v != validator.original.Type {
		return validation.NewError("validation_collection_type_change", "Collection type cannot be changed.")
	}

	return nil
}

func (validator *collectionValidator) ensureNoFieldsTypeChange(value any) error {
	v, ok := value.(FieldsList)
	if !ok {
		return validators.ErrUnsupportedValueType
	}

	errs := validation.Errors{}

	for i, field := range v {
		oldField := validator.original.Fields.GetById(field.GetId())

		if oldField != nil && oldField.Type() != field.Type() {
			errs[strconv.Itoa(i)] = validation.NewError(
				"validation_field_type_change",
				"Field type cannot be changed.",
			)
		}
	}
	if len(errs) > 0 {
		return errs
	}

	return nil
}

func (validator *collectionValidator) checkFieldDuplicates(value any) error {
	fields, ok := value.(FieldsList)
	if !ok {
		return validators.ErrUnsupportedValueType
	}

	totalFields := len(fields)
	ids := make([]string, 0, totalFields)
	names := make([]string, 0, totalFields)

	for i, field := range fields {
		if list.ExistInSlice(field.GetId(), ids) {
			return validation.Errors{
				strconv.Itoa(i): validation.Errors{
					"id": validation.NewError(
						"validation_duplicated_field_id",
						fmt.Sprintf("Duplicated or invalid field id %q", field.GetId()),
					),
				},
			}
		}

		// field names are used as db columns and should be case insensitive
		nameLower := strings.ToLower(field.GetName())

		if list.ExistInSlice(nameLower, names) {
			return validation.Errors{
				strconv.Itoa(i): validation.Errors{
					"name": validation.NewError(
						"validation_duplicated_field_name",
						"Duplicated or invalid field name {{.fieldName}}",
					).SetParams(map[string]any{
						"fieldName": field.GetName(),
					}),
				},
			}
		}

		ids = append(ids, field.GetId())
		names = append(names, nameLower)
	}

	return nil
}

func (validator *collectionValidator) checkFieldValidators(value any) error {
	fields, ok := value.(FieldsList)
	if !ok {
		return validators.ErrUnsupportedValueType
	}

	errs := validation.Errors{}

	for i, field := range fields {
		if err := field.ValidateSettings(validator.ctx, validator.app, validator.new); err != nil {
			errs[strconv.Itoa(i)] = err
		}
	}

	if len(errs) > 0 {
		return errs
	}

	return nil
}

func (cv *collectionValidator) checkViewQuery(value any) error {
	v, _ := value.(string)
	if v == "" {
		return nil // nothing to check
	}

	if _, err := cv.app.CreateViewFields(v); err != nil {
		return validation.NewError(
			"validation_invalid_view_query",
			fmt.Sprintf("Invalid query - %s", err.Error()),
		)
	}

	return nil
}

var reservedAuthKeys = []string{"passwordConfirm", "oldPassword"}

func (cv *collectionValidator) checkReservedAuthKeys(value any) error {
	fields, ok := value.(FieldsList)
	if !ok {
		return validators.ErrUnsupportedValueType
	}

	if !cv.new.IsAuth() {
		return nil // not an auth collection
	}

	errs := validation.Errors{}
	for i, field := range fields {
		if list.ExistInSlice(field.GetName(), reservedAuthKeys) {
			errs[strconv.Itoa(i)] = validation.Errors{
				"name": validation.NewError(
					"validation_reserved_field_name",
					"The field name is reserved and cannot be used.",
				),
			}
		}
	}
	if len(errs) > 0 {
		return errs
	}

	return nil
}

func (cv *collectionValidator) checkMinFields(value any) error {
	fields, ok := value.(FieldsList)
	if !ok {
		return validators.ErrUnsupportedValueType
	}

	if len(fields) == 0 {
		return validation.ErrRequired
	}

	// all collections must have an "id" PK field
	idField, _ := fields.GetByName(FieldNameId).(*TextField)
	if idField == nil || !idField.PrimaryKey {
		return validation.NewError("validation_missing_primary_key", `Missing or invalid "id" PK field.`)
	}

	switch cv.new.Type {
	case CollectionTypeAuth:
		passwordField, _ := fields.GetByName(FieldNamePassword).(*PasswordField)
		if passwordField == nil {
			return validation.NewError("validation_missing_password_field", `System "password" field is required.`)
		}
		if !passwordField.Hidden || !passwordField.System {
			return validation.Errors{FieldNamePassword: ErrMustBeSystemAndHidden}
		}

		tokenKeyField, _ := fields.GetByName(FieldNameTokenKey).(*TextField)
		if tokenKeyField == nil {
			return validation.NewError("validation_missing_tokenKey_field", `System "tokenKey" field is required.`)
		}
		if !tokenKeyField.Hidden || !tokenKeyField.System {
			return validation.Errors{FieldNameTokenKey: ErrMustBeSystemAndHidden}
		}

		emailField, _ := fields.GetByName(FieldNameEmail).(*EmailField)
		if emailField == nil {
			return validation.NewError("validation_missing_email_field", `System "email" field is required.`)
		}
		if !emailField.System {
			return validation.Errors{FieldNameEmail: ErrMustBeSystem}
		}

		emailVisibilityField, _ := fields.GetByName(FieldNameEmailVisibility).(*BoolField)
		if emailVisibilityField == nil {
			return validation.NewError("validation_missing_emailVisibility_field", `System "emailVisibility" field is required.`)
		}
		if !emailVisibilityField.System {
			return validation.Errors{FieldNameEmailVisibility: ErrMustBeSystem}
		}

		verifiedField, _ := fields.GetByName(FieldNameVerified).(*BoolField)
		if verifiedField == nil {
			return validation.NewError("validation_missing_verified_field", `System "verified" field is required.`)
		}
		if !verifiedField.System {
			return validation.Errors{FieldNameVerified: ErrMustBeSystem}
		}

		return nil
	}

	return nil
}

func (validator *collectionValidator) ensureNoSystemFieldsChange(value any) error {
	fields, ok := value.(FieldsList)
	if !ok {
		return validators.ErrUnsupportedValueType
	}

	if validator.original.IsNew() {
		return nil // not an update
	}

	for _, oldField := range validator.original.Fields {
		if !oldField.GetSystem() {
			continue
		}

		newField := fields.GetById(oldField.GetId())

		if newField == nil || oldField.GetName() != newField.GetName() {
			return validation.NewError("validation_system_field_change", "System fields cannot be deleted or renamed.")
		}
	}

	return nil
}

func (cv *collectionValidator) checkFieldsForUniqueIndex(value any) error {
	names, ok := value.([]string)
	if !ok {
		return validators.ErrUnsupportedValueType
	}

	if len(names) == 0 {
		return nil // nothing to check
	}

	for _, name := range names {
		field := cv.new.Fields.GetByName(name)
		if field == nil {
			return validation.NewError("validation_missing_field", "Invalid or missing field {{.fieldName}}").
				SetParams(map[string]any{"fieldName": name})
		}

		if _, ok := dbutils.FindSingleColumnUniqueIndex(cv.new.Indexes, name); !ok {
			return validation.NewError("validation_missing_unique_constraint", "The field {{.fieldName}} doesn't have a UNIQUE constraint.").
				SetParams(map[string]any{"fieldName": name})
		}
	}

	return nil
}

// note: value could be either *string or string
func (validator *collectionValidator) checkRule(value any) error {
	var vStr string

	v, ok := value.(*string)
	if ok {
		if v != nil {
			vStr = *v
		}
	} else {
		vStr, ok = value.(string)
	}
	if !ok {
		return validators.ErrUnsupportedValueType
	}

	if vStr == "" {
		return nil // nothing to check
	}

	r := NewRecordFieldResolver(validator.app, validator.new, nil, true)
	_, err := search.FilterData(vStr).BuildExpr(r)
	if err != nil {
		return validation.NewError("validation_invalid_rule", "Invalid rule. Raw error: "+err.Error())
	}

	return nil
}

func (validator *collectionValidator) ensureNoSystemRuleChange(oldRule *string) validation.RuleFunc {
	return func(value any) error {
		if validator.original.IsNew() || !validator.original.System {
			return nil // not an update of a system collection
		}

		rule, ok := value.(*string)
		if !ok {
			return validators.ErrUnsupportedValueType
		}

		if (rule == nil && oldRule == nil) ||
			(rule != nil && oldRule != nil && *rule == *oldRule) {
			return nil
		}

		return validation.NewError("validation_collection_system_rule_change", "System collection API rule cannot be changed.")
	}
}

func (cv *collectionValidator) checkIndexes(value any) error {
	indexes, _ := value.(types.JSONArray[string])

	if cv.new.IsView() && len(indexes) > 0 {
		return validation.NewError(
			"validation_indexes_not_supported",
			"View collections don't support indexes.",
		)
	}

	duplicatedNames := make(map[string]struct{}, len(indexes))
	duplicatedDefinitions := make(map[string]struct{}, len(indexes))

	for i, rawIndex := range indexes {
		parsed := dbutils.ParseIndex(rawIndex)

		// always set a table name because it is ignored anyway in order to keep it in sync with the collection name
		parsed.TableName = "validator"

		if !parsed.IsValid() {
			return validation.Errors{
				strconv.Itoa(i): validation.NewError(
					"validation_invalid_index_expression",
					"Invalid CREATE INDEX expression.",
				),
			}
		}

		if _, isDuplicated := duplicatedNames[strings.ToLower(parsed.IndexName)]; isDuplicated {
			return validation.Errors{
				strconv.Itoa(i): validation.NewError(
					"validation_duplicated_index_name",
					"The index name already exists.",
				),
			}
		}
		duplicatedNames[strings.ToLower(parsed.IndexName)] = struct{}{}

		// ensure that the index name is not used in another collection
		var usedTblName string
		_ = cv.app.DB().Select("tbl_name").
			From("sqlite_master").
			AndWhere(dbx.HashExp{"type": "index"}).
			AndWhere(dbx.NewExp("LOWER([[tbl_name]])!=LOWER({:oldName})", dbx.Params{"oldName": cv.original.Name})).
			AndWhere(dbx.NewExp("LOWER([[tbl_name]])!=LOWER({:newName})", dbx.Params{"newName": cv.new.Name})).
			AndWhere(dbx.NewExp("LOWER([[name]])=LOWER({:indexName})", dbx.Params{"indexName": parsed.IndexName})).
			Limit(1).
			Row(&usedTblName)
		if usedTblName != "" {
			return validation.Errors{
				strconv.Itoa(i): validation.NewError(
					"validation_existing_index_name",
					"The index name is already used in {{.usedTableName}} collection.",
				).SetParams(map[string]any{"usedTableName": usedTblName}),
			}
		}

		// reset non-important identifiers
		parsed.SchemaName = "validator"
		parsed.IndexName = "validator"
		parsedDef := parsed.Build()

		if _, isDuplicated := duplicatedDefinitions[parsedDef]; isDuplicated {
			return validation.Errors{
				strconv.Itoa(i): validation.NewError(
					"validation_duplicated_index_definition",
					"The index definition already exists.",
				),
			}
		}
		duplicatedDefinitions[parsedDef] = struct{}{}

		// note: we don't check the index table name because it is always
		// overwritten by the SyncRecordTableSchema to allow
		// easier partial modifications (eg. changing only the collection name).
		// if !strings.EqualFold(parsed.TableName, form.Name) {
		// 	return validation.Errors{
		// 		strconv.Itoa(i): validation.NewError(
		// 			"validation_invalid_index_table",
		// 			fmt.Sprintf("The index table must be the same as the collection name."),
		// 		),
		// 	}
		// }
	}

	// ensure that unique indexes on system fields are not changed or removed
	if !cv.original.IsNew() {
	OLD_INDEXES_LOOP:
		for _, oldIndex := range cv.original.Indexes {
			oldParsed := dbutils.ParseIndex(oldIndex)
			if !oldParsed.Unique {
				continue
			}

			// reset collate and sort since they are not important for the unique constraint
			for i := range oldParsed.Columns {
				oldParsed.Columns[i].Collate = ""
				oldParsed.Columns[i].Sort = ""
			}

			oldParsedStr := oldParsed.Build()

			for _, column := range oldParsed.Columns {
				for _, f := range cv.original.Fields {
					if !f.GetSystem() || !strings.EqualFold(column.Name, f.GetName()) {
						continue
					}

					var hasMatch bool
					for _, newIndex := range cv.new.Indexes {
						newParsed := dbutils.ParseIndex(newIndex)

						// exclude the non-important identifiers from the check
						newParsed.SchemaName = oldParsed.SchemaName
						newParsed.IndexName = oldParsed.IndexName
						newParsed.TableName = oldParsed.TableName

						// exclude partial constraints
						newParsed.Where = oldParsed.Where

						// reset collate and sort
						for i := range newParsed.Columns {
							newParsed.Columns[i].Collate = ""
							newParsed.Columns[i].Sort = ""
						}

						if oldParsedStr == newParsed.Build() {
							hasMatch = true
							break
						}
					}

					if !hasMatch {
						return validation.NewError(
							"validation_invalid_unique_system_field_index",
							"Unique index definition on system fields ({{.fieldName}}) is invalid or missing.",
						).SetParams(map[string]any{"fieldName": f.GetName()})
					}

					continue OLD_INDEXES_LOOP
				}
			}
		}
	}

	// check for required indexes
	//
	// note: this is in case the indexes were removed manually when creating/importing new auth collections
	// and technically it is not necessary because on app.Save() the missing indexes will be reinserted by the system collection hook
	if cv.new.IsAuth() {
		requiredNames := []string{FieldNameTokenKey, FieldNameEmail}
		for _, name := range requiredNames {
			if _, ok := dbutils.FindSingleColumnUniqueIndex(indexes, name); !ok {
				return validation.NewError(
					"validation_missing_required_unique_index",
					`Missing required unique index for field "{{.fieldName}}".`,
				).SetParams(map[string]any{"fieldName": name})
			}
		}
	}

	return nil
}

func (validator *collectionValidator) validateOptions() error {
	switch validator.new.Type {
	case CollectionTypeAuth:
		return validator.new.collectionAuthOptions.validate(validator)
	case CollectionTypeView:
		return validator.new.collectionViewOptions.validate(validator)
	}

	return nil
}
