package daos

import (
	"fmt"
	"strings"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/models/schema"
)

func (dao *Dao) ViewQuery() *dbx.SelectQuery {
	return dao.ModelQuery(&models.View{})
}

func (dao *Dao) RecordFromViewQuery(view *models.View) *dbx.SelectQuery {
	viewName := view.Name
	selectCols := fmt.Sprintf("%s.*", dao.DB().QuoteSimpleColumnName(viewName))

	return dao.DB().Select(selectCols).From(viewName)
}

// FindViewByIdOrName finds the view by its name or id.
func (dao *Dao) FindViewByIdOrName(nameOrId string) (*models.View, error) {
	model := &models.View{}

	err := dao.ViewQuery().
		AndWhere(dbx.Or(
			dbx.HashExp{"id": nameOrId},
			dbx.HashExp{"name": nameOrId},
		)).
		Limit(1).
		One(model)
	if err != nil {
		return nil, err
	}

	return model, nil
}

// IsViewNameUnique checks that there is no existing view
// with the provided name (case insensitive!).
func (dao *Dao) IsViewNameUnique(name string, excludeId string) bool {
	if name == "" {
		return false
	}

	var exists bool
	err := dao.ViewQuery().
		Select("count(*)").
		AndWhere(dbx.Not(dbx.HashExp{"id": excludeId})).
		AndWhere(dbx.NewExp("LOWER([[name]])={:name}", dbx.Params{"name": strings.ToLower(name)})).
		Limit(1).
		Row(&exists)

	return err == nil && !exists
}

// DeleteView deletes the provided View model.
func (dao *Dao) DeleteView(view *models.View) error {
	return dao.RunInTransaction(func(txDao *Dao) error {
		sql := "DROP VIEW " + dao.DB().QuoteSimpleTableName(view.Name)
		_, err := dao.DB().NewQuery(sql).Execute()
		if err != nil {
			return err
		}
		return txDao.Delete(view)
	})
}

func (dao *Dao) CreateOrReplaceView(view *models.View) error {
	// drop if exists
	sql := "DROP VIEW IF EXISTS " + dao.DB().QuoteSimpleTableName(view.Name)
	_, err := dao.DB().NewQuery(sql).Execute()
	if err != nil {
		return err
	}
	// create
	sql = fmt.Sprintf("CREATE VIEW %s AS %s", dao.DB().QuoteSimpleTableName(view.Name), view.Sql)
	_, err = dao.DB().NewQuery(sql).Execute()

	return err
}

func (dao *Dao) SaveView(view *models.View) (*models.View, error) {
	var err error
	dao.RunInTransaction(func(txDao *Dao) error {
		err = txDao.CreateOrReplaceView(view)
		return err
	})
	if err != nil {
		return nil, err
	}
	Schema, err := dao.GetViewSchema(view.Name)
	if err != nil {
		return nil, err
	}
	view.Schema = Schema
	err = dao.Save(view)
	if err != nil {
		return nil, err
	}
	return view, nil
}

func (dao *Dao) GetViewSchema(viewName string) (s schema.Schema, err error) {
	rows := []dbx.NullStringMap{}
	s = schema.Schema{}
	// this query will return the view columns underlying type ('TEXT','BOOLEAN','REAL','JSON')
	sql := fmt.Sprintf(`PRAGMA table_info(%s)`, dao.db.QuoteSimpleTableName(viewName))
	err = dao.DB().NewQuery(sql).All(&rows)
	if err != nil {
		return
	}
	for _, row := range rows {
		// we cant exactly know if the text type is (text, date, url, ...etc), because the underlying type is the same.
		// we will do an extra check on the field when it's retrieved
		fieldType := ""
		switch strings.ToLower(row["type"].String) {
		case "boolean":
			fieldType = schema.FieldTypeBool
		case "json":
			fieldType = schema.FieldTypeJson
		case "real":
			fieldType = schema.FieldTypeNumber
		default:
			fieldType = schema.FieldTypeText
		}
		newField := schema.SchemaField{
			Name: row["name"].String,
			Type: fieldType,
		}
		s.AddField(&newField)
	}
	return
}

func (dao *Dao) HasView(viewName string) bool {
	var exists bool

	err := dao.DB().Select("count(*)").
		From("sqlite_schema").
		AndWhere(dbx.HashExp{"type": "view"}).
		AndWhere(dbx.NewExp("LOWER([[name]])=LOWER({:viewName})", dbx.Params{"viewName": viewName})).
		Limit(1).Row(&exists)

	return err == nil && exists
}
