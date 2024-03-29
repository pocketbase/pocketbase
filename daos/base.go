// Package daos handles common PocketBase DB model manipulations.
//
// Think of daos as DB repository and service layer in one.
package daos

import (
	"errors"
	"fmt"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
)

// New creates a new Dao instance with the provided db builder
// (for both async and sync db operations).
func New(db dbx.Builder) *Dao {
	return NewMultiDB(db, db)
}

// NewMultiDB creates a new Dao instance with the provided dedicated
// async and sync db builders.
func NewMultiDB(concurrentDB, nonconcurrentDB dbx.Builder) *Dao {
	return &Dao{
		concurrentDB:      concurrentDB,
		nonconcurrentDB:   nonconcurrentDB,
		MaxLockRetries:    8,
		ModelQueryTimeout: 30 * time.Second,
	}
}

// Dao handles various db operations.
//
// You can think of Dao as a repository and service layer in one.
type Dao struct {
	// in a transaction both refer to the same *dbx.TX instance
	concurrentDB    dbx.Builder
	nonconcurrentDB dbx.Builder

	// MaxLockRetries specifies the default max "database is locked" auto retry attempts.
	MaxLockRetries int

	// ModelQueryTimeout is the default max duration of a running ModelQuery().
	//
	// This field has no effect if an explicit query context is already specified.
	ModelQueryTimeout time.Duration

	// write hooks
	BeforeCreateFunc func(eventDao *Dao, m models.Model, action func() error) error
	AfterCreateFunc  func(eventDao *Dao, m models.Model) error
	BeforeUpdateFunc func(eventDao *Dao, m models.Model, action func() error) error
	AfterUpdateFunc  func(eventDao *Dao, m models.Model) error
	BeforeDeleteFunc func(eventDao *Dao, m models.Model, action func() error) error
	AfterDeleteFunc  func(eventDao *Dao, m models.Model) error
}

// DB returns the default dao db builder (*dbx.DB or *dbx.TX).
//
// Currently the default db builder is dao.concurrentDB but that may change in the future.
func (dao *Dao) DB() dbx.Builder {
	return dao.ConcurrentDB()
}

// ConcurrentDB returns the dao concurrent (aka. multiple open connections)
// db builder (*dbx.DB or *dbx.TX).
//
// In a transaction the concurrentDB and nonconcurrentDB refer to the same *dbx.TX instance.
func (dao *Dao) ConcurrentDB() dbx.Builder {
	return dao.concurrentDB
}

// NonconcurrentDB returns the dao nonconcurrent (aka. single open connection)
// db builder (*dbx.DB or *dbx.TX).
//
// In a transaction the concurrentDB and nonconcurrentDB refer to the same *dbx.TX instance.
func (dao *Dao) NonconcurrentDB() dbx.Builder {
	return dao.nonconcurrentDB
}

// Clone returns a new Dao with the same configuration options as the current one.
func (dao *Dao) Clone() *Dao {
	clone := *dao

	return &clone
}

// WithoutHooks returns a new Dao with the same configuration options
// as the current one, but without create/update/delete hooks.
func (dao *Dao) WithoutHooks() *Dao {
	clone := dao.Clone()

	clone.BeforeCreateFunc = nil
	clone.AfterCreateFunc = nil
	clone.BeforeUpdateFunc = nil
	clone.AfterUpdateFunc = nil
	clone.BeforeDeleteFunc = nil
	clone.AfterDeleteFunc = nil

	return clone
}

// ModelQuery creates a new preconfigured select query with preset
// SELECT, FROM and other common fields based on the provided model.
func (dao *Dao) ModelQuery(m models.Model) *dbx.SelectQuery {
	tableName := m.TableName()

	return dao.DB().
		Select("{{" + tableName + "}}.*").
		From(tableName).
		WithBuildHook(func(query *dbx.Query) {
			query.WithExecHook(execLockRetry(dao.ModelQueryTimeout, dao.MaxLockRetries))
		})
}

// FindById finds a single db record with the specified id and
// scans the result into m.
func (dao *Dao) FindById(m models.Model, id string) error {
	return dao.ModelQuery(m).Where(dbx.HashExp{"id": id}).Limit(1).One(m)
}

type afterCallGroup struct {
	Model    models.Model
	EventDao *Dao
	Action   string
}

// RunInTransaction wraps fn into a transaction.
//
// It is safe to nest RunInTransaction calls as long as you use the txDao.
func (dao *Dao) RunInTransaction(fn func(txDao *Dao) error) error {
	switch txOrDB := dao.NonconcurrentDB().(type) {
	case *dbx.Tx:
		// nested transactions are not supported by default
		// so execute the function within the current transaction
		// ---
		// create a new dao with the same hooks to avoid semaphore deadlock when nesting
		txDao := New(txOrDB)
		txDao.MaxLockRetries = dao.MaxLockRetries
		txDao.ModelQueryTimeout = dao.ModelQueryTimeout
		txDao.BeforeCreateFunc = dao.BeforeCreateFunc
		txDao.BeforeUpdateFunc = dao.BeforeUpdateFunc
		txDao.BeforeDeleteFunc = dao.BeforeDeleteFunc
		txDao.AfterCreateFunc = dao.AfterCreateFunc
		txDao.AfterUpdateFunc = dao.AfterUpdateFunc
		txDao.AfterDeleteFunc = dao.AfterDeleteFunc

		return fn(txDao)
	case *dbx.DB:
		afterCalls := []afterCallGroup{}

		txError := txOrDB.Transactional(func(tx *dbx.Tx) error {
			txDao := New(tx)

			if dao.BeforeCreateFunc != nil {
				txDao.BeforeCreateFunc = func(eventDao *Dao, m models.Model, action func() error) error {
					return dao.BeforeCreateFunc(eventDao, m, action)
				}
			}
			if dao.BeforeUpdateFunc != nil {
				txDao.BeforeUpdateFunc = func(eventDao *Dao, m models.Model, action func() error) error {
					return dao.BeforeUpdateFunc(eventDao, m, action)
				}
			}
			if dao.BeforeDeleteFunc != nil {
				txDao.BeforeDeleteFunc = func(eventDao *Dao, m models.Model, action func() error) error {
					return dao.BeforeDeleteFunc(eventDao, m, action)
				}
			}

			if dao.AfterCreateFunc != nil {
				txDao.AfterCreateFunc = func(eventDao *Dao, m models.Model) error {
					afterCalls = append(afterCalls, afterCallGroup{m, eventDao, "create"})
					return nil
				}
			}
			if dao.AfterUpdateFunc != nil {
				txDao.AfterUpdateFunc = func(eventDao *Dao, m models.Model) error {
					afterCalls = append(afterCalls, afterCallGroup{m, eventDao, "update"})
					return nil
				}
			}
			if dao.AfterDeleteFunc != nil {
				txDao.AfterDeleteFunc = func(eventDao *Dao, m models.Model) error {
					afterCalls = append(afterCalls, afterCallGroup{m, eventDao, "delete"})
					return nil
				}
			}

			return fn(txDao)
		})
		if txError != nil {
			return txError
		}

		// execute after event calls on successful transaction
		// (note: using the non-transaction dao to allow following queries in the after hooks)
		var errs []error
		for _, call := range afterCalls {
			var err error
			switch call.Action {
			case "create":
				err = dao.AfterCreateFunc(dao, call.Model)
			case "update":
				err = dao.AfterUpdateFunc(dao, call.Model)
			case "delete":
				err = dao.AfterDeleteFunc(dao, call.Model)
			}

			if err != nil {
				errs = append(errs, err)
			}
		}
		if len(errs) > 0 {
			return fmt.Errorf("after transaction errors: %w", errors.Join(errs...))
		}

		return nil
	}

	return errors.New("failed to start transaction (unknown dao.NonconcurrentDB() instance)")
}

// Delete deletes the provided model.
func (dao *Dao) Delete(m models.Model) error {
	if !m.HasId() {
		return errors.New("ID is not set")
	}

	return dao.lockRetry(func(retryDao *Dao) error {
		action := func() error {
			if err := retryDao.NonconcurrentDB().Model(m).Delete(); err != nil {
				return err
			}

			if retryDao.AfterDeleteFunc != nil {
				retryDao.AfterDeleteFunc(retryDao, m)
			}

			return nil
		}

		if retryDao.BeforeDeleteFunc != nil {
			return retryDao.BeforeDeleteFunc(retryDao, m, action)
		}

		return action()
	})
}

// Save persists the provided model in the database.
//
// If m.IsNew() is true, the method will perform a create, otherwise an update.
// To explicitly mark a model for update you can use m.MarkAsNotNew().
func (dao *Dao) Save(m models.Model) error {
	if m.IsNew() {
		return dao.lockRetry(func(retryDao *Dao) error {
			return retryDao.create(m)
		})
	}

	return dao.lockRetry(func(retryDao *Dao) error {
		return retryDao.update(m)
	})
}

func (dao *Dao) update(m models.Model) error {
	if !m.HasId() {
		return errors.New("ID is not set")
	}

	if m.GetCreated().IsZero() {
		m.RefreshCreated()
	}

	m.RefreshUpdated()

	action := func() error {
		if v, ok := any(m).(models.ColumnValueMapper); ok {
			dataMap := v.ColumnValueMap()

			_, err := dao.NonconcurrentDB().Update(
				m.TableName(),
				dataMap,
				dbx.HashExp{"id": m.GetId()},
			).Execute()

			if err != nil {
				return err
			}
		} else if err := dao.NonconcurrentDB().Model(m).Update(); err != nil {
			return err
		}

		if dao.AfterUpdateFunc != nil {
			return dao.AfterUpdateFunc(dao, m)
		}

		return nil
	}

	if dao.BeforeUpdateFunc != nil {
		return dao.BeforeUpdateFunc(dao, m, action)
	}

	return action()
}

func (dao *Dao) create(m models.Model) error {
	if !m.HasId() {
		// auto generate id
		m.RefreshId()
	}

	// mark the model as "new" since the model now always has an ID
	m.MarkAsNew()

	if m.GetCreated().IsZero() {
		m.RefreshCreated()
	}

	if m.GetUpdated().IsZero() {
		m.RefreshUpdated()
	}

	action := func() error {
		if v, ok := any(m).(models.ColumnValueMapper); ok {
			dataMap := v.ColumnValueMap()
			if _, ok := dataMap["id"]; !ok {
				dataMap["id"] = m.GetId()
			}

			_, err := dao.NonconcurrentDB().Insert(m.TableName(), dataMap).Execute()
			if err != nil {
				return err
			}
		} else if err := dao.NonconcurrentDB().Model(m).Insert(); err != nil {
			return err
		}

		// clears the "new" model flag
		m.MarkAsNotNew()

		if dao.AfterCreateFunc != nil {
			return dao.AfterCreateFunc(dao, m)
		}

		return nil
	}

	if dao.BeforeCreateFunc != nil {
		return dao.BeforeCreateFunc(dao, m, action)
	}

	return action()
}

func (dao *Dao) lockRetry(op func(retryDao *Dao) error) error {
	retryDao := dao

	return baseLockRetry(func(attempt int) error {
		if attempt == 2 {
			// assign new Dao without the before hooks to avoid triggering
			// the already fired before callbacks multiple times
			retryDao = NewMultiDB(dao.concurrentDB, dao.nonconcurrentDB)
			retryDao.AfterCreateFunc = dao.AfterCreateFunc
			retryDao.AfterUpdateFunc = dao.AfterUpdateFunc
			retryDao.AfterDeleteFunc = dao.AfterDeleteFunc
		}

		return op(retryDao)
	}, dao.MaxLockRetries)
}
