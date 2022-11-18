package forms

import (
	"fmt"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/resolvers"
	"github.com/pocketbase/pocketbase/tools/search"
	_ "modernc.org/sqlite"
)

// CollectionUpsert specifies a [models.Collection] upsert (create/update) form.
type ViewUpsert struct {
	config ViewUpsertConfig
	view   *models.View

	Id       string  `form:"id" json:"id"`
	Sql      string  `form:"sql" json:"sql"`
	Name     string  `form:"name" json:"name"`
	ListRule *string `form:"listRule" json:"listRule"`
}

// CollectionUpsertConfig is the [CollectionUpsert] factory initializer config.
//
// NB! App is a required struct member.
type ViewUpsertConfig struct {
	App core.App
	Dao *daos.Dao
}

// NewCollectionUpsert creates a new [CollectionUpsert] form with initializer
// config created from the provided [core.App] and [models.Collection] instances
// (for create you could pass a pointer to an empty Collection - `&models.Collection{}`).
//
// If you want to submit the form as part of another transaction, use
// [NewCollectionUpsertWithConfig] with explicitly set Dao.
func NewViewUpsert(app core.App, view *models.View) *ViewUpsert {
	return NewViewUpsertWithConfig(ViewUpsertConfig{
		App: app,
	}, view)
}

// NewCollectionUpsertWithConfig creates a new [CollectionUpsert] form
// with the provided config and [models.Collection] instance or panics on invalid configuration
// (for create you could pass a pointer to an empty Collection - `&models.Collection{}`).
func NewViewUpsertWithConfig(config ViewUpsertConfig, view *models.View) *ViewUpsert {
	form := &ViewUpsert{
		config: config,
		view:   view,
	}

	if form.config.App == nil || form.view == nil {
		panic("Invalid initializer config or nil upsert model.")
	}

	if form.config.Dao == nil {
		form.config.Dao = form.config.App.Dao()
	}

	// load defaults
	form.Id = form.view.Id
	form.Name = form.view.Name
	form.ListRule = form.view.ListRule
	form.Sql = form.view.Sql

	return form
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *ViewUpsert) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(
			&form.Sql,
			validation.Required,
			validation.By(form.checkSqlValid),
		),
		validation.Field(
			&form.Name,
			validation.Required,
			validation.Length(1, 255),
			validation.Match(collectionNameRegex),
			validation.By(form.checkUniqueName),
		),
		validation.Field(
			&form.Id,
			validation.When(
				form.view.IsNew(),
				validation.Length(models.DefaultIdLength, models.DefaultIdLength),
				validation.Match(idRegex),
			).Else(validation.In(form.view.Id)),
		),
	)
}

func (form *ViewUpsert) checkUniqueName(value any) error {
	v, _ := value.(string)

	if !form.config.Dao.IsViewNameUnique(v, form.view.Id) {
		return validation.NewError("validation_view_name_exists", "View name must be unique (case insensitive).")
	}

	if form.config.Dao.HasTable(form.view.Name) {
		return validation.NewError("validation_view_name_same_as_collection", fmt.Sprintf("View name '%s'must not be same as collection (case insensitive).", form.view.Name))
	}

	return nil
}

func (form *ViewUpsert) checkSqlValid(value any) error {
	// the query provided should not affect any data
	// testing before and after the query to make sure the data do not change
	query, _ := value.(string)
	query = strings.Split(query, ";")[0]

	db := form.config.Dao.DB().(*dbx.DB)
	tx, err := db.Begin()
	var totalChangesBefore int
	defer tx.Rollback()
	// testing changes before
	// total_changes() function return changes done by insert, update, delete

	// doing the test in transaction even if the data changes it can be reverted
	db.NewQuery("select total_changes()").Row(&totalChangesBefore)
	if err != nil {
		return err
	}
	_, err = tx.NewQuery(query).Rows()
	if err != nil {
		return validation.NewError("validation_sql_invalid", err.Error())
	}
	var totalChangesAfter int
	// testing changes after
	db.NewQuery("select total_changes()").Row(&totalChangesAfter)

	fmt.Println("before", totalChangesBefore, "after", totalChangesAfter)

	if totalChangesAfter > totalChangesBefore {
		return validation.NewError("validation_sql_invalid", "SQL should not affect data")
	}
	return nil
}

// Submit validates the form and upserts the form's Collection model.
//
// On success the related record table schema will be auto updated.
//
// You can optionally provide a list of InterceptorFunc to further
// modify the form behavior before persisting it.
func (form *ViewUpsert) Submit(interceptors ...InterceptorFunc) error {
	if err := form.Validate(); err != nil {
		return err
	}
	form.Sql = strings.Split(form.Sql, ";")[0]

	if form.view.IsNew() {
		// custom insertion id can be set only on create
		if form.Id != "" {
			form.view.MarkAsNew()
			form.view.SetId(form.Id)
		}
	}

	form.view.Name = form.Name
	form.view.ListRule = form.ListRule
	form.view.Sql = form.Sql

	return runInterceptors(func() error {
		var err error
		form.config.Dao.RunInTransaction(func(txDao *daos.Dao) error {
			var view *models.View
			view, err = form.config.Dao.SaveView(form.view)
			if err != nil {
				return err
			}
			// ListRule cannot be checked before the view is created
			// because the schema is dynamically generated
			err = checkViewListRule(view, form.config.Dao)
			return err
		})
		return err
	}, interceptors...)
}

func checkViewListRule(view *models.View, dao *daos.Dao) error {
	listRule := view.ListRule
	if listRule == nil || *listRule == "" {
		return nil // nothing to check
	}

	dummy := &models.Collection{Schema: view.Schema}
	r := resolvers.NewRecordFieldResolver(dao, dummy, nil, false)

	_, err := search.FilterData(*listRule).BuildExpr(r)
	if err != nil {
		return validation.NewError("validation_collection_rule", "Invalid filter rule.")
	}

	return nil
}
