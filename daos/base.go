// Package daos handles common PocketBase DB model manipulations.
//
// Think of daos as DB repository and service layer in one.
package daos

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/models"
	"golang.org/x/sync/semaphore"
)

const DefaultMaxFailRetries = 5

// New creates a new Dao instance with the provided db builder.
func New(db dbx.Builder) *Dao {
	return &Dao{
		db: db,
	}
}

// Dao handles various db operations.
// Think of Dao as a repository and service layer in one.
type Dao struct {
	db  dbx.Builder
	sem *semaphore.Weighted
	mux sync.RWMutex

	BeforeCreateFunc func(eventDao *Dao, m models.Model) error
	AfterCreateFunc  func(eventDao *Dao, m models.Model)
	BeforeUpdateFunc func(eventDao *Dao, m models.Model) error
	AfterUpdateFunc  func(eventDao *Dao, m models.Model)
	BeforeDeleteFunc func(eventDao *Dao, m models.Model) error
	AfterDeleteFunc  func(eventDao *Dao, m models.Model)
}

// DB returns the internal db builder (*dbx.DB or *dbx.TX).
func (dao *Dao) DB() dbx.Builder {
	return dao.db
}

// Block acquires a lock and blocks all other go routines that uses
// the Dao instance until dao.Continue() is called, effectively making
// the concurrent requests to perform synchronous db operations.
//
// This method should be used only as a last resort and as a workaround
// for the SQLITE_BUSY error when mixing read&write in a transaction.
//
// Example:
//
// 	func someLongRunningTransaction() error {
// 		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 		defer cancel()
// 		if err := app.Dao().Block(ctx); err != nil {
// 			return err
// 		}
// 		defer app.Dao().Continue()
//
//  	return app.Dao().RunInTransaction(func (txDao *daos.Dao) error {
//  	    // some long running read&write transaction...
//  	})
// 	}
func (dao *Dao) Block(ctx context.Context) error {
	if dao.sem == nil {
		dao.mux.Lock()
		dao.sem = semaphore.NewWeighted(1)
		dao.mux.Unlock()
	}

	return dao.sem.Acquire(ctx, 1)
}

// Continue releases the previously acquired Block() lock.
func (dao *Dao) Continue() {
	if dao.sem == nil {
		return
	}

	dao.sem.Release(1)
}

// ModelQuery creates a new query with preset Select and From fields
// based on the provided model argument.
func (dao *Dao) ModelQuery(m models.Model) *dbx.SelectQuery {
	tableName := m.TableName()
	return dao.db.Select("{{" + tableName + "}}.*").From(tableName)
}

// FindById finds a single db record with the specified id and
// scans the result into m.
func (dao *Dao) FindById(m models.Model, id string) error {
	return dao.ModelQuery(m).Where(dbx.HashExp{"id": id}).Limit(1).One(m)
}

type afterCallGroup struct {
	Action   string
	EventDao *Dao
	Model    models.Model
}

// RunInTransaction wraps fn into a transaction.
//
// It is safe to nest RunInTransaction calls.
func (dao *Dao) RunInTransaction(fn func(txDao *Dao) error) error {
	switch txOrDB := dao.db.(type) {
	case *dbx.Tx:
		// nested transactions are not supported by default
		// so execute the function within the current transaction

		// create a new dao with the same hooks to avoid semaphore deadlock when nesting
		txDao := New(txOrDB)
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
				txDao.BeforeCreateFunc = func(eventDao *Dao, m models.Model) error {
					return dao.BeforeCreateFunc(eventDao, m)
				}
			}
			if dao.BeforeUpdateFunc != nil {
				txDao.BeforeUpdateFunc = func(eventDao *Dao, m models.Model) error {
					return dao.BeforeUpdateFunc(eventDao, m)
				}
			}
			if dao.BeforeDeleteFunc != nil {
				txDao.BeforeDeleteFunc = func(eventDao *Dao, m models.Model) error {
					return dao.BeforeDeleteFunc(eventDao, m)
				}
			}

			if dao.AfterCreateFunc != nil {
				txDao.AfterCreateFunc = func(eventDao *Dao, m models.Model) {
					afterCalls = append(afterCalls, afterCallGroup{"create", eventDao, m})
				}
			}
			if dao.AfterUpdateFunc != nil {
				txDao.AfterUpdateFunc = func(eventDao *Dao, m models.Model) {
					afterCalls = append(afterCalls, afterCallGroup{"update", eventDao, m})
				}
			}
			if dao.AfterDeleteFunc != nil {
				txDao.AfterDeleteFunc = func(eventDao *Dao, m models.Model) {
					afterCalls = append(afterCalls, afterCallGroup{"delete", eventDao, m})
				}
			}

			return fn(txDao)
		})

		if txError == nil {
			// execute after event calls on successful transaction
			for _, call := range afterCalls {
				switch call.Action {
				case "create":
					dao.AfterCreateFunc(call.EventDao, call.Model)
				case "update":
					dao.AfterUpdateFunc(call.EventDao, call.Model)
				case "delete":
					dao.AfterDeleteFunc(call.EventDao, call.Model)
				}
			}
		}

		return txError
	}

	return errors.New("Failed to start transaction (unknown dao.db)")
}

// Delete deletes the provided model.
func (dao *Dao) Delete(m models.Model) error {
	if !m.HasId() {
		return errors.New("ID is not set")
	}

	return dao.failRetry(func(retryDao *Dao) error {
		if retryDao.BeforeDeleteFunc != nil {
			if err := retryDao.BeforeDeleteFunc(retryDao, m); err != nil {
				return err
			}
		}

		if err := retryDao.db.Model(m).Delete(); err != nil {
			return err
		}

		if retryDao.AfterDeleteFunc != nil {
			retryDao.AfterDeleteFunc(retryDao, m)
		}

		return nil
	}, DefaultMaxFailRetries)
}

// Save upserts (update or create if primary key is not set) the provided model.
func (dao *Dao) Save(m models.Model) error {
	if m.IsNew() {
		return dao.failRetry(func(retryDao *Dao) error {
			return retryDao.create(m)
		}, DefaultMaxFailRetries)
	}

	return dao.failRetry(func(retryDao *Dao) error {
		return retryDao.update(m)
	}, DefaultMaxFailRetries)
}

func (dao *Dao) update(m models.Model) error {
	if !m.HasId() {
		return errors.New("ID is not set")
	}

	if m.GetCreated().IsZero() {
		m.RefreshCreated()
	}

	m.RefreshUpdated()

	if dao.BeforeUpdateFunc != nil {
		if err := dao.BeforeUpdateFunc(dao, m); err != nil {
			return err
		}
	}

	if v, ok := any(m).(models.ColumnValueMapper); ok {
		dataMap := v.ColumnValueMap()

		_, err := dao.db.Update(
			m.TableName(),
			dataMap,
			dbx.HashExp{"id": m.GetId()},
		).Execute()

		if err != nil {
			return err
		}
	} else {
		if err := dao.db.Model(m).Update(); err != nil {
			return err
		}
	}

	if dao.AfterUpdateFunc != nil {
		dao.AfterUpdateFunc(dao, m)
	}

	return nil
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

	if dao.BeforeCreateFunc != nil {
		if err := dao.BeforeCreateFunc(dao, m); err != nil {
			return err
		}
	}

	if v, ok := any(m).(models.ColumnValueMapper); ok {
		dataMap := v.ColumnValueMap()
		if _, ok := dataMap["id"]; !ok {
			dataMap["id"] = m.GetId()
		}

		_, err := dao.db.Insert(m.TableName(), dataMap).Execute()
		if err != nil {
			return err
		}
	} else {
		if err := dao.db.Model(m).Insert(); err != nil {
			return err
		}
	}

	// clears the "new" model flag
	m.MarkAsNotNew()

	if dao.AfterCreateFunc != nil {
		dao.AfterCreateFunc(dao, m)
	}

	return nil
}

func (dao *Dao) failRetry(op func(retryDao *Dao) error, maxRetries int) error {
	retryDao := dao
	attempts := 1

Retry:
	if attempts == 2 {
		// assign new Dao without the before hooks to avoid triggering
		// the already fired before callbacks multiple times
		retryDao = &Dao{
			db:              dao.db,
			AfterCreateFunc: dao.AfterCreateFunc,
			AfterUpdateFunc: dao.AfterUpdateFunc,
			AfterDeleteFunc: dao.AfterDeleteFunc,
		}
	}

	// execute
	err := op(retryDao)

	if err != nil &&
		attempts < maxRetries &&
		// note: we are checking the err message so that we can handle both the cgo and noncgo errors
		strings.Contains(err.Error(), "database is locked") {
		// wait and retry
		time.Sleep(time.Duration(200*attempts) * time.Millisecond)
		attempts++
		goto Retry
	}

	return err
}
