package core

import (
	"errors"
	"fmt"
	"sync"

	"github.com/pocketbase/dbx"
)

// RunInTransaction wraps fn into a transaction for the regular app database.
//
// It is safe to nest RunInTransaction calls as long as you use the callback's txApp.
func (app *BaseApp) RunInTransaction(fn func(txApp App) error) error {
	return app.runInTransaction(app.NonconcurrentDB(), fn, false)
}

// AuxRunInTransaction wraps fn into a transaction for the auxiliary app database.
//
// It is safe to nest RunInTransaction calls as long as you use the callback's txApp.
func (app *BaseApp) AuxRunInTransaction(fn func(txApp App) error) error {
	return app.runInTransaction(app.AuxNonconcurrentDB(), fn, true)
}

func (app *BaseApp) runInTransaction(db dbx.Builder, fn func(txApp App) error, isForAuxDB bool) error {
	switch txOrDB := db.(type) {
	case *dbx.Tx:
		// run as part of the already existing transaction
		return fn(app)
	case *dbx.DB:
		var txApp *BaseApp
		txErr := txOrDB.Transactional(func(tx *dbx.Tx) error {
			txApp = app.createTxApp(tx, isForAuxDB)
			return fn(txApp)
		})

		// execute all after event calls on transaction complete
		if txApp != nil && txApp.txInfo != nil {
			afterFuncErr := txApp.txInfo.runAfterFuncs(txErr)
			if afterFuncErr != nil {
				return errors.Join(txErr, afterFuncErr)
			}
		}

		return txErr
	default:
		return errors.New("failed to start transaction (unknown db type)")
	}
}

// createTxApp shallow clones the current app and assigns a new tx state.
func (app *BaseApp) createTxApp(tx *dbx.Tx, isForAuxDB bool) *BaseApp {
	clone := *app

	if isForAuxDB {
		clone.auxConcurrentDB = tx
		clone.auxNonconcurrentDB = tx
	} else {
		clone.concurrentDB = tx
		clone.nonconcurrentDB = tx
	}

	clone.txInfo = &txAppInfo{
		parent:     app,
		isForAuxDB: isForAuxDB,
	}

	return &clone
}

type txAppInfo struct {
	parent     *BaseApp
	afterFuncs []func(txErr error) error
	mu         sync.Mutex
	isForAuxDB bool
}

func (a *txAppInfo) onAfterFunc(fn func(txErr error) error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.afterFuncs = append(a.afterFuncs, fn)
}

// note: can be called only once because txAppInfo is cleared
func (a *txAppInfo) runAfterFuncs(txErr error) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	var errs []error

	for _, call := range a.afterFuncs {
		if err := call(txErr); err != nil {
			errs = append(errs, err)
		}
	}

	a.afterFuncs = nil

	if len(errs) > 0 {
		return fmt.Errorf("transaction afterFunc errors: %w", errors.Join(errs...))
	}

	return nil
}
