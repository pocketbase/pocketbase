package forms

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
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
		validation.Field(
			&form.Sql,
			validation.Required,
			validation.By(form.checkSqlValid),
		),
	)
}

func (form *ViewUpsert) checkUniqueName(value any) error {
	v, _ := value.(string)

	if !form.config.Dao.IsViewNameUnique(v, form.view.Id) {
		return validation.NewError("validation_view_name_exists", "View name must be unique (case insensitive).")
	}
	if (form.view.IsNew() || !strings.EqualFold(v, form.view.Name)) && form.config.Dao.HasView(v) {
		return validation.NewError("validation_view_name_table_exists", "The View name must be also unique view name.")
	}

	return nil
}

func (form *ViewUpsert) checkSqlValid(value any) error {
	query, _ := value.(string)
	dbPath := filepath.Join(form.config.App.DataDir(), "data.db")
	params := "mode=ro"
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?%s", dbPath, params))
	defer db.Close()
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	res, err := tx.Exec(query)
	if err != nil {
		return validation.NewError("validation_sql_invalid", err.Error())
	}
	tx.Rollback()
	n, _ := res.RowsAffected()
	if n > 0 {
		return validation.NewError("validation_sql_invalid", "SQL should not affect data")
	}
	if err != nil {
		return validation.NewError("validation_sql_invalid", err.Error())
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

	if form.view.IsNew() {
		// custom insertion id can be set only on create
		if form.Id != "" {
			form.view.MarkAsNew()
			form.view.SetId(form.Id)
		}
	}

	// system collections cannot be renamed
	if form.view.IsNew() {
		form.view.Name = form.Name
	}

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
	r := resolvers.NewRecordFieldResolver(dao, dummy, nil)

	_, err := search.FilterData(*listRule).BuildExpr(r)
	if err != nil {
		return validation.NewError("validation_collection_rule", "Invalid filter rule.")
	}

	return nil
}
