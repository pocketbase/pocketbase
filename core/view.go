package core

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/tools/dbutils"
	"github.com/pocketbase/pocketbase/tools/inflector"
	"github.com/pocketbase/pocketbase/tools/security"
	"github.com/pocketbase/pocketbase/tools/tokenizer"
)

// DeleteView drops the specified view name.
//
// This method is a no-op if a view with the provided name doesn't exist.
//
// NB! Be aware that this method is vulnerable to SQL injection and the
// "dangerousViewName" argument must come only from trusted input!
func (app *BaseApp) DeleteView(dangerousViewName string) error {
	_, err := app.DB().NewQuery(fmt.Sprintf(
		"DROP VIEW IF EXISTS {{%s}}",
		dangerousViewName,
	)).Execute()

	return err
}

// SaveView creates (or updates already existing) persistent SQL view.
//
// NB! Be aware that this method is vulnerable to SQL injection and
// its arguments must come only from trusted input!
func (app *BaseApp) SaveView(dangerousViewName string, dangerousSelectQuery string) error {
	return app.RunInTransaction(func(txApp App) error {
		// delete old view (if exists)
		if err := txApp.DeleteView(dangerousViewName); err != nil {
			return err
		}

		dangerousSelectQuery = strings.Trim(strings.TrimSpace(dangerousSelectQuery), ";")

		// try to loosely detect multiple inline statements
		tk := tokenizer.NewFromString(dangerousSelectQuery)
		tk.Separators(';')
		if queryParts, _ := tk.ScanAll(); len(queryParts) > 1 {
			return errors.New("multiple statements are not supported")
		}

		// (re)create the view
		//
		// note: the query is wrapped in a secondary SELECT as a rudimentary
		// measure to discourage multiple inline sql statements execution
		viewQuery := fmt.Sprintf("CREATE VIEW {{%s}} AS SELECT * FROM (%s)", dangerousViewName, dangerousSelectQuery)
		if _, err := txApp.DB().NewQuery(viewQuery).Execute(); err != nil {
			return err
		}

		// fetch the view table info to ensure that the view was created
		// because missing tables or columns won't return an error
		if _, err := txApp.TableInfo(dangerousViewName); err != nil {
			// manually cleanup previously created view in case the func
			// is called in a nested transaction and the error is discarded
			txApp.DeleteView(dangerousViewName)

			return err
		}

		return nil
	})
}

// CreateViewFields creates a new FieldsList from the provided select query.
//
// There are some caveats:
// - The select query must have an "id" column.
// - Wildcard ("*") columns are not supported to avoid accidentally leaking sensitive data.
//
// NB! Be aware that this method is vulnerable to SQL injection and the
// "dangerousSelectQuery" argument must come only from trusted input!
func (app *BaseApp) CreateViewFields(dangerousSelectQuery string) (FieldsList, error) {
	result := NewFieldsList()

	suggestedFields, err := parseQueryToFields(app, dangerousSelectQuery)
	if err != nil {
		return result, err
	}

	// note wrap in a transaction in case the dangerousSelectQuery contains
	// multiple statements allowing us to rollback on any error
	txErr := app.RunInTransaction(func(txApp App) error {
		info, err := getQueryTableInfo(txApp, dangerousSelectQuery)
		if err != nil {
			return err
		}

		var hasId bool

		for _, row := range info {
			if row.Name == FieldNameId {
				hasId = true
			}

			var field Field

			if f, ok := suggestedFields[row.Name]; ok {
				field = f.field
			} else {
				field = defaultViewField(row.Name)
			}

			result.Add(field)
		}

		if !hasId {
			return errors.New("missing required id column (you can use `(ROW_NUMBER() OVER()) as id` if you don't have one)")
		}

		return nil
	})

	return result, txErr
}

// FindRecordByViewFile returns the original Record of the provided view collection file.
func (app *BaseApp) FindRecordByViewFile(viewCollectionModelOrIdentifier any, fileFieldName string, filename string) (*Record, error) {
	view, err := getCollectionByModelOrIdentifier(app, viewCollectionModelOrIdentifier)
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

		queryFields, err := parseQueryToFields(app, view.ViewQuery)
		if err != nil {
			return nil, err
		}

		for _, item := range queryFields {
			if item.collection == nil ||
				item.original == nil ||
				item.field.GetName() != fileFieldName {
				continue
			}

			if item.collection.IsView() {
				view = item.collection
				fileFieldName = item.original.GetName()
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

	cleanFieldName := inflector.Columnify(qf.original.GetName())

	record := &Record{}

	query := app.RecordQuery(qf.collection).Limit(1)

	if opt, ok := qf.original.(MultiValuer); !ok || !opt.IsMultiple() {
		query.AndWhere(dbx.HashExp{cleanFieldName: filename})
	} else {
		query.InnerJoin(
			fmt.Sprintf(`%s as {{_je_file}}`, dbutils.JSONEach(cleanFieldName)),
			dbx.HashExp{"_je_file.value": filename},
		)
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
	field Field

	// collection refers to the original field's collection model.
	// It could be nil if the found query field is not from a collection
	collection *Collection

	// original is the original found collection field.
	// It could be nil if the found query field is not from a collection
	original Field
}

func defaultViewField(name string) Field {
	return &JSONField{
		Name:    name,
		MaxSize: 1, // unused for views
	}
}

var castRegex = regexp.MustCompile(`(?is)^cast\s*\(.*\s+as\s+(\w+)\s*\)$`)

func parseQueryToFields(app App, selectQuery string) (map[string]*queryField, error) {
	p := new(identifiersParser)
	if err := p.parse(selectQuery); err != nil {
		return nil, err
	}

	collections, err := findCollectionsByIdentifiers(app, p.tables)
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

		// pk (always assume text field for now)
		if col.alias == FieldNameId {
			result[col.alias] = &queryField{
				field: &TextField{
					Name:       col.alias,
					System:     true,
					Required:   true,
					PrimaryKey: true,
					Pattern:    `^[a-z0-9]+$`,
				},
			}
			continue
		}

		// numeric aggregations
		if strings.HasPrefix(colLower, "count(") || strings.HasPrefix(colLower, "total(") {
			result[col.alias] = &queryField{
				field: &NumberField{
					Name: col.alias,
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
					field: &NumberField{
						Name: col.alias,
					},
				}
				continue
			case "text":
				result[col.alias] = &queryField{
					field: &TextField{
						Name: col.alias,
					},
				}
				continue
			case "boolean", "bool":
				result[col.alias] = &queryField{
					field: &BoolField{
						Name: col.alias,
					},
				}
				continue
			}
		}

		parts := strings.Split(col.original, ".")

		var fieldName string
		var collection *Collection

		if len(parts) == 2 {
			fieldName = parts[1]
			collection = collections[parts[0]]
		} else {
			fieldName = parts[0]
			collection = collections[mainTable.alias]
		}

		// fallback to the default field
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
		var field Field
		for _, f := range collection.Fields {
			if strings.EqualFold(f.GetName(), fieldName) {
				field = f
				break
			}
		}

		// fallback to the default field
		if field == nil {
			result[col.alias] = &queryField{
				field:      defaultViewField(col.alias),
				collection: collection,
			}
			continue
		}

		// convert to relation since it is an id reference
		if strings.EqualFold(fieldName, FieldNameId) {
			result[col.alias] = &queryField{
				field: &RelationField{
					Name:         col.alias,
					MaxSelect:    1,
					CollectionId: collection.Id,
				},
				collection: collection,
			}
			continue
		}

		// we fetch a brand new collection object to avoid using reflection
		// or having a dedicated Clone method for each field type
		tempCollection, err := app.FindCollectionByNameOrId(collection.Id)
		if err != nil {
			return nil, err
		}

		clone := tempCollection.Fields.GetById(field.GetId())
		if clone == nil {
			return nil, fmt.Errorf("missing expected field %q (%q) in collection %q", field.GetName(), field.GetId(), tempCollection.Name)
		}
		// set new random id to prevent duplications if the same field is aliased multiple times
		clone.SetId("_clone_" + security.PseudorandomString(4))
		clone.SetName(col.alias)

		result[col.alias] = &queryField{
			original:   field,
			field:      clone,
			collection: collection,
		}
	}

	return result, nil
}

func findCollectionsByIdentifiers(app App, tables []identifier) (map[string]*Collection, error) {
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

	result := make(map[string]*Collection, len(names))
	collections := make([]*Collection, 0, len(names))

	err := app.CollectionQuery().
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

func getQueryTableInfo(app App, selectQuery string) ([]*TableInfoRow, error) {
	tempView := "_temp_" + security.PseudorandomString(6)

	var info []*TableInfoRow

	txErr := app.RunInTransaction(func(txApp App) error {
		// create a temp view with the provided query
		err := txApp.SaveView(tempView, selectQuery)
		if err != nil {
			return err
		}

		// extract the generated view table info
		info, err = txApp.TableInfo(tempView)

		return errors.Join(err, txApp.DeleteView(tempView))
	})

	if txErr != nil {
		return nil, txErr
	}

	return info, nil
}

// -------------------------------------------------------------------
// Raw query identifiers parser
// -------------------------------------------------------------------

var (
	joinReplaceRegex     = regexp.MustCompile(`(?im)\s+(full\s+outer\s+join|left\s+outer\s+join|right\s+outer\s+join|full\s+join|cross\s+join|inner\s+join|outer\s+join|left\s+join|right\s+join|join)\s+?`)
	discardReplaceRegex  = regexp.MustCompile(`(?im)\s+(where|group\s+by|having|order|limit|with)\s+?`)
	commentsReplaceRegex = regexp.MustCompile(`(?m)(\/\*[\s\S]*?\*\/)|(--.+$)`)
)

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
	str = commentsReplaceRegex.ReplaceAllString(str, " ")
	str = joinReplaceRegex.ReplaceAllString(str, " __pb_join__ ")
	str = discardReplaceRegex.ReplaceAllString(str, " __pb_discard__ ")

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
		case "__pb_join__":
			skip = false

			// the previous part was also a join
			if partType == "join" {
				joinParts.WriteString(",")
			}

			partType = "join"
			activeBuilder = &joinParts
		case "__pb_discard__":
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
