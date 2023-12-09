package daos

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
	"github.com/pocketbase/pocketbase/tools/inflector"
	"github.com/pocketbase/pocketbase/tools/list"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/pocketbase/pocketbase/tools/tokenizer"
	"github.com/pocketbase/pocketbase/tools/types"
)

// DeleteView drops the specified view name.
//
// This method is a no-op if a view with the provided name doesn't exist.
//
// Be aware that this method is vulnerable to SQL injection and the
// "name" argument must come only from trusted input!
func (dao *Dao) DeleteView(name string) error {
	_, err := dao.DB().NewQuery(fmt.Sprintf(
		"DROP VIEW IF EXISTS {{%s}}",
		name,
	)).Execute()

	return err
}

// SaveView creates (or updates already existing) persistent SQL view.
//
// Be aware that this method is vulnerable to SQL injection and the
// "selectQuery" argument must come only from trusted input!
func (dao *Dao) SaveView(name string, selectQuery string) error {
	return dao.RunInTransaction(func(txDao *Dao) error {
		// delete old view (if exists)
		if err := txDao.DeleteView(name); err != nil {
			return err
		}

		selectQuery = strings.Trim(strings.TrimSpace(selectQuery), ";")

		// try to eagerly detect multiple inline statements
		tk := tokenizer.NewFromString(selectQuery)
		tk.Separators(';')
		if queryParts, _ := tk.ScanAll(); len(queryParts) > 1 {
			return errors.New("multiple statements are not supported")
		}

		// (re)create the view
		//
		// note: the query is wrapped in a secondary SELECT as a rudimentary
		// measure to discourage multiple inline sql statements execution.
		viewQuery := fmt.Sprintf("CREATE VIEW {{%s}} AS SELECT * FROM (%s)", name, selectQuery)
		if _, err := txDao.DB().NewQuery(viewQuery).Execute(); err != nil {
			return err
		}

		// fetch the view table info to ensure that the view was created
		// because missing tables or columns won't return an error
		if _, err := txDao.TableInfo(name); err != nil {
			// manually cleanup previously created view in case the func
			// is called in a nested transaction and the error is discarded
			txDao.DeleteView(name)

			return err
		}

		return nil
	})
}

// CreateViewSchema creates a new view schema from the provided select query.
//
// There are some caveats:
// - The select query must have an "id" column.
// - Wildcard ("*") columns are not supported to avoid accidentally leaking sensitive data.
func (dao *Dao) CreateViewSchema(selectQuery string) (schema.Schema, error) {
	result := schema.NewSchema()

	suggestedFields, err := dao.parseQueryToFields(selectQuery)
	if err != nil {
		return result, err
	}

	// note wrap in a transaction in case the selectQuery contains
	// multiple statements allowing us to rollback on any error
	txErr := dao.RunInTransaction(func(txDao *Dao) error {
		tempView := "_temp_" + security.PseudorandomString(5)
		// create a temp view with the provided query
		if err := txDao.SaveView(tempView, selectQuery); err != nil {
			return err
		}
		defer txDao.DeleteView(tempView)

		// extract the generated view table info
		info, err := txDao.TableInfo(tempView)
		if err != nil {
			return err
		}

		var hasId bool

		for _, row := range info {
			if row.Name == schema.FieldNameId {
				hasId = true
			}

			if list.ExistInSlice(row.Name, schema.BaseModelFieldNames()) {
				continue // skip base model fields since they are not part of the schema
			}

			var field *schema.SchemaField

			if f, ok := suggestedFields[row.Name]; ok {
				field = f.field
			} else {
				field = defaultViewField(row.Name)
			}

			result.AddField(field)
		}

		if !hasId {
			return errors.New("missing required id column (you can use `(ROW_NUMBER() OVER()) as id` if you don't have one)")
		}

		return nil
	})

	return result, txErr
}

// FindRecordByViewFile returns the original models.Record of the
// provided view collection file.
func (dao *Dao) FindRecordByViewFile(
	viewCollectionNameOrId string,
	fileFieldName string,
	filename string,
) (*models.Record, error) {
	view, err := dao.FindCollectionByNameOrId(viewCollectionNameOrId)
	if err != nil {
		return nil, err
	}

	if !view.IsView() {
		return nil, errors.New("not a view collection")
	}

	var findFirstNonViewQueryFileField func(int) (*queryField, error)
	findFirstNonViewQueryFileField = func(level int) (*queryField, error) {
		// check the level depth to prevent infinite circular recursion
		// (the limit is arbitrary and may change in the future)
		if level > 5 {
			return nil, errors.New("reached the max recursion level of view collection file field queries")
		}

		queryFields, err := dao.parseQueryToFields(view.ViewOptions().Query)
		if err != nil {
			return nil, err
		}

		for _, item := range queryFields {
			if item.collection == nil ||
				item.original == nil ||
				item.field.Name != fileFieldName {
				continue
			}

			if item.collection.IsView() {
				view = item.collection
				fileFieldName = item.original.Name
				return findFirstNonViewQueryFileField(level + 1)
			}

			return item, nil
		}

		return nil, errors.New("no query file field found")
	}

	qf, err := findFirstNonViewQueryFileField(1)
	if err != nil {
		return nil, err
	}

	cleanFieldName := inflector.Columnify(qf.original.Name)

	record := &models.Record{}

	query := dao.RecordQuery(qf.collection).Limit(1)

	if opt, ok := qf.original.Options.(schema.MultiValuer); !ok || !opt.IsMultiple() {
		query.AndWhere(dbx.HashExp{cleanFieldName: filename})
	} else {
		query.InnerJoin(fmt.Sprintf(
			`json_each(CASE WHEN json_valid([[%s]]) THEN [[%s]] ELSE json_array([[%s]]) END) as {{_je_file}}`,
			cleanFieldName, cleanFieldName, cleanFieldName,
		), dbx.HashExp{"_je_file.value": filename})
	}

	if err := query.One(record); err != nil {
		return nil, err
	}

	return record, nil
}

// -------------------------------------------------------------------
// Raw query to schema helpers
// -------------------------------------------------------------------

type queryField struct {
	// field is the final resolved field.
	field *schema.SchemaField

	// collection refers to the original field's collection model.
	// It could be nil if the found query field is not from a collection schema.
	collection *models.Collection

	// original is the original found collection field.
	// It could be nil if the found query field is not from a collection schema.
	original *schema.SchemaField
}

func defaultViewField(name string) *schema.SchemaField {
	return &schema.SchemaField{
		Name: name,
		Type: schema.FieldTypeJson,
		Options: &schema.JsonOptions{
			MaxSize: 1, // the size doesn't matter in this case
		},
	}
}

var castRegex = regexp.MustCompile(`(?i)^cast\s*\(.*\s+as\s+(\w+)\s*\)$`)

func (dao *Dao) parseQueryToFields(selectQuery string) (map[string]*queryField, error) {
	p := new(identifiersParser)
	if err := p.parse(selectQuery); err != nil {
		return nil, err
	}

	collections, err := dao.findCollectionsByIdentifiers(p.tables)
	if err != nil {
		return nil, err
	}

	result := make(map[string]*queryField, len(p.columns))

	var mainTable identifier

	if len(p.tables) > 0 {
		mainTable = p.tables[0]
	}

	for _, col := range p.columns {
		colLower := strings.ToLower(col.original)

		// numeric aggregations
		if strings.HasPrefix(colLower, "count(") || strings.HasPrefix(colLower, "total(") {
			result[col.alias] = &queryField{
				field: &schema.SchemaField{
					Name: col.alias,
					Type: schema.FieldTypeNumber,
				},
			}
			continue
		}

		castMatch := castRegex.FindStringSubmatch(colLower)

		// numeric casts
		if len(castMatch) == 2 {
			switch castMatch[1] {
			case "real", "integer", "int", "decimal", "numeric":
				result[col.alias] = &queryField{
					field: &schema.SchemaField{
						Name: col.alias,
						Type: schema.FieldTypeNumber,
					},
				}
				continue
			case "text":
				result[col.alias] = &queryField{
					field: &schema.SchemaField{
						Name: col.alias,
						Type: schema.FieldTypeText,
					},
				}
				continue
			case "boolean", "bool":
				result[col.alias] = &queryField{
					field: &schema.SchemaField{
						Name: col.alias,
						Type: schema.FieldTypeBool,
					},
				}
				continue
			}
		}

		parts := strings.Split(col.original, ".")

		var fieldName string
		var collection *models.Collection

		if len(parts) == 2 {
			fieldName = parts[1]
			collection = collections[parts[0]]
		} else {
			fieldName = parts[0]
			collection = collections[mainTable.alias]
		}

		// fallback to the default field if the found column is not from a collection schema
		if collection == nil {
			result[col.alias] = &queryField{
				field: defaultViewField(col.alias),
			}
			continue
		}

		if fieldName == "*" {
			return nil, errors.New("dynamic column names are not supported")
		}

		// find the first field by name (case insensitive)
		var field *schema.SchemaField
		for _, f := range collection.Schema.Fields() {
			if strings.EqualFold(f.Name, fieldName) {
				field = f
				break
			}
		}

		if field != nil {
			clone := *field
			clone.Id = "" // unset to prevent duplications if the same field is aliased multiple times
			clone.Name = col.alias
			result[col.alias] = &queryField{
				field:      &clone,
				collection: collection,
				original:   field,
			}
			continue
		}

		if fieldName == schema.FieldNameId {
			// convert to relation since it is a direct id reference
			result[col.alias] = &queryField{
				field: &schema.SchemaField{
					Name: col.alias,
					Type: schema.FieldTypeRelation,
					Options: &schema.RelationOptions{
						MaxSelect:    types.Pointer(1),
						CollectionId: collection.Id,
					},
				},
				collection: collection,
			}
		} else if fieldName == schema.FieldNameCreated || fieldName == schema.FieldNameUpdated {
			result[col.alias] = &queryField{
				field: &schema.SchemaField{
					Name: col.alias,
					Type: schema.FieldTypeDate,
				},
				collection: collection,
			}
		} else if fieldName == schema.FieldNameUsername && collection.IsAuth() {
			result[col.alias] = &queryField{
				field: &schema.SchemaField{
					Name: col.alias,
					Type: schema.FieldTypeText,
				},
				collection: collection,
			}
		} else if fieldName == schema.FieldNameEmail && collection.IsAuth() {
			result[col.alias] = &queryField{
				field: &schema.SchemaField{
					Name: col.alias,
					Type: schema.FieldTypeEmail,
				},
				collection: collection,
			}
		} else if (fieldName == schema.FieldNameVerified || fieldName == schema.FieldNameEmailVisibility) && collection.IsAuth() {
			result[col.alias] = &queryField{
				field: &schema.SchemaField{
					Name: col.alias,
					Type: schema.FieldTypeBool,
				},
				collection: collection,
			}
		} else {
			result[col.alias] = &queryField{
				field:      defaultViewField(col.alias),
				collection: collection,
			}
		}
	}

	return result, nil
}

func (dao *Dao) findCollectionsByIdentifiers(tables []identifier) (map[string]*models.Collection, error) {
	names := make([]any, 0, len(tables))

	for _, table := range tables {
		if strings.Contains(table.alias, "(") {
			continue // skip expressions
		}
		names = append(names, table.original)
	}

	if len(names) == 0 {
		return nil, nil
	}

	result := make(map[string]*models.Collection, len(names))
	collections := make([]*models.Collection, 0, len(names))

	err := dao.CollectionQuery().
		AndWhere(dbx.In("name", names...)).
		All(&collections)
	if err != nil {
		return nil, err
	}

	for _, table := range tables {
		for _, collection := range collections {
			if collection.Name == table.original {
				result[table.alias] = collection
			}
		}
	}

	return result, nil
}

// -------------------------------------------------------------------
// Raw query identifiers parser
// -------------------------------------------------------------------

var joinReplaceRegex = regexp.MustCompile(`(?im)\s+(inner join|outer join|left join|right join|join)\s+?`)
var discardReplaceRegex = regexp.MustCompile(`(?im)\s+(where|group by|having|order|limit|with)\s+?`)
var commentsReplaceRegex = regexp.MustCompile(`(?m)(\/\*[\s\S]+\*\/)|(--.+$)`)

type identifier struct {
	original string
	alias    string
}

type identifiersParser struct {
	columns []identifier
	tables  []identifier
}

func (p *identifiersParser) parse(selectQuery string) error {
	str := strings.Trim(strings.TrimSpace(selectQuery), ";")
	str = joinReplaceRegex.ReplaceAllString(str, " _join_ ")
	str = discardReplaceRegex.ReplaceAllString(str, " _discard_ ")
	str = commentsReplaceRegex.ReplaceAllString(str, "")

	tk := tokenizer.NewFromString(str)
	tk.Separators(',', ' ', '\n', '\t')
	tk.KeepSeparator(true)

	var skip bool
	var partType string
	var activeBuilder *strings.Builder
	var selectParts strings.Builder
	var fromParts strings.Builder
	var joinParts strings.Builder

	for {
		token, err := tk.Scan()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}

		trimmed := strings.ToLower(strings.TrimSpace(token))

		switch trimmed {
		case "select":
			skip = false
			partType = "select"
			activeBuilder = &selectParts
		case "distinct":
			continue // ignore as it is not important for the identifiers parsing
		case "from":
			skip = false
			partType = "from"
			activeBuilder = &fromParts
		case "_join_":
			skip = false

			// the previous part was also a join
			if partType == "join" {
				joinParts.WriteString(",")
			}

			partType = "join"
			activeBuilder = &joinParts
		case "_discard_":
			// skip following tokens
			skip = true
		default:
			isJoin := partType == "join"

			if isJoin && trimmed == "on" {
				skip = true
			}

			if !skip && activeBuilder != nil {
				activeBuilder.WriteString(" ")
				activeBuilder.WriteString(token)
			}
		}
	}

	selects, err := extractIdentifiers(selectParts.String())
	if err != nil {
		return err
	}

	froms, err := extractIdentifiers(fromParts.String())
	if err != nil {
		return err
	}

	joins, err := extractIdentifiers(joinParts.String())
	if err != nil {
		return err
	}

	p.columns = selects
	p.tables = froms
	p.tables = append(p.tables, joins...)

	return nil
}

func extractIdentifiers(rawExpression string) ([]identifier, error) {
	rawTk := tokenizer.NewFromString(rawExpression)
	rawTk.Separators(',')

	rawIdentifiers, err := rawTk.ScanAll()
	if err != nil {
		return nil, err
	}

	result := make([]identifier, 0, len(rawIdentifiers))

	for _, rawIdentifier := range rawIdentifiers {
		tk := tokenizer.NewFromString(rawIdentifier)
		tk.Separators(' ', '\n', '\t')

		parts, err := tk.ScanAll()
		if err != nil {
			return nil, err
		}

		resolved, err := identifierFromParts(parts)
		if err != nil {
			return nil, err
		}

		result = append(result, resolved)
	}

	return result, nil
}

func identifierFromParts(parts []string) (identifier, error) {
	var result identifier

	switch len(parts) {
	case 3:
		if !strings.EqualFold(parts[1], "as") {
			return result, fmt.Errorf(`invalid identifier part - expected "as", got %v`, parts[1])
		}

		result.original = parts[0]
		result.alias = parts[2]
	case 2:
		result.original = parts[0]
		result.alias = parts[1]
	case 1:
		subParts := strings.Split(parts[0], ".")
		result.original = parts[0]
		result.alias = subParts[len(subParts)-1]
	default:
		return result, fmt.Errorf(`invalid identifier parts %v`, parts)
	}

	result.original = trimRawIdentifier(result.original)

	// we trim the single quote even though it is not a valid column quote character
	// because SQLite allows it if the context expects an identifier and not string literal
	// (https://www.sqlite.org/lang_keywords.html)
	result.alias = trimRawIdentifier(result.alias, "'")

	return result, nil
}

func trimRawIdentifier(rawIdentifier string, extraTrimChars ...string) string {
	trimChars := "`\"[];"
	if len(extraTrimChars) > 0 {
		trimChars += strings.Join(extraTrimChars, "")
	}

	parts := strings.Split(rawIdentifier, ".")

	for i := range parts {
		parts[i] = strings.Trim(parts[i], trimChars)
	}

	return strings.Join(parts, ".")
}
