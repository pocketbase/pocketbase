package core

import (
	"github.com/pocketbase/pocketbase/tools/hook"
)

// LambdaFunctionEvent represents a lambda function change event.
type LambdaFunctionEvent struct {
	hook.Event
	App            App
	LambdaFunction *LambdaFunction
}

// LambdaFunctionErrorEvent represents a lambda function error event.
type LambdaFunctionErrorEvent struct {
	hook.Event
	App          App
	LambdaFunction *LambdaFunction
	Error        error
}

// LambdaFunctionExecuteEvent represents a lambda function execution event.
type LambdaFunctionExecuteEvent struct {
	hook.Event
	App     App
	Context *LambdaFunctionContext
	Result  *LambdaFunctionResult
}

const systemHookIdLambdaFunction = "__pbLambdaFunctionSystemHook__"

func (app *BaseApp) registerLambdaFunctionHooks() {
	// Hook into model events for LambdaFunction
	app.OnModelValidate().Bind(&hook.Handler[*ModelEvent]{
		Id: systemHookIdLambdaFunction,
		Func: func(me *ModelEvent) error {
			if ef, ok := me.Model.(*LambdaFunction); ok {
				event := &LambdaFunctionEvent{
					Event:        me.Event,
					App:          me.App,
					LambdaFunction: ef,
				}
				return app.OnLambdaFunctionValidate().Trigger(event, func(e *LambdaFunctionEvent) error {
					me.Model = e.LambdaFunction
					return me.Next()
				})
			}
			return me.Next()
		},
		Priority: -99,
	})

	app.OnModelCreate().Bind(&hook.Handler[*ModelEvent]{
		Id: systemHookIdLambdaFunction,
		Func: func(me *ModelEvent) error {
			if ef, ok := me.Model.(*LambdaFunction); ok {
				event := &LambdaFunctionEvent{
					Event:        me.Event,
					App:          me.App,
					LambdaFunction: ef,
				}
				return app.OnLambdaFunctionCreate().Trigger(event, func(e *LambdaFunctionEvent) error {
					me.Model = e.LambdaFunction
					return me.Next()
				})
			}
			return me.Next()
		},
		Priority: -99,
	})

	app.OnModelUpdate().Bind(&hook.Handler[*ModelEvent]{
		Id: systemHookIdLambdaFunction,
		Func: func(me *ModelEvent) error {
			if ef, ok := me.Model.(*LambdaFunction); ok {
				event := &LambdaFunctionEvent{
					Event:        me.Event,
					App:          me.App,
					LambdaFunction: ef,
				}
				return app.OnLambdaFunctionUpdate().Trigger(event, func(e *LambdaFunctionEvent) error {
					me.Model = e.LambdaFunction
					return me.Next()
				})
			}
			return me.Next()
		},
		Priority: -99,
	})

	app.OnModelDelete().Bind(&hook.Handler[*ModelEvent]{
		Id: systemHookIdLambdaFunction,
		Func: func(me *ModelEvent) error {
			if ef, ok := me.Model.(*LambdaFunction); ok {
				event := &LambdaFunctionEvent{
					Event:        me.Event,
					App:          me.App,
					LambdaFunction: ef,
				}
				return app.OnLambdaFunctionDelete().Trigger(event, func(e *LambdaFunctionEvent) error {
					me.Model = e.LambdaFunction
					return me.Next()
				})
			}
			return me.Next()
		},
		Priority: -99,
	})

	// Success hooks
	app.OnModelAfterCreateSuccess().Bind(&hook.Handler[*ModelEvent]{
		Id: systemHookIdLambdaFunction,
		Func: func(me *ModelEvent) error {
			if ef, ok := me.Model.(*LambdaFunction); ok {
				event := &LambdaFunctionEvent{
					Event:        me.Event,
					App:          me.App,
					LambdaFunction: ef,
				}
				return app.OnLambdaFunctionAfterCreateSuccess().Trigger(event, func(e *LambdaFunctionEvent) error {
					return me.Next()
				})
			}
			return me.Next()
		},
		Priority: -99,
	})

	app.OnModelAfterUpdateSuccess().Bind(&hook.Handler[*ModelEvent]{
		Id: systemHookIdLambdaFunction,
		Func: func(me *ModelEvent) error {
			if ef, ok := me.Model.(*LambdaFunction); ok {
				event := &LambdaFunctionEvent{
					Event:        me.Event,
					App:          me.App,
					LambdaFunction: ef,
				}
				return app.OnLambdaFunctionAfterUpdateSuccess().Trigger(event, func(e *LambdaFunctionEvent) error {
					return me.Next()
				})
			}
			return me.Next()
		},
		Priority: -99,
	})

	app.OnModelAfterDeleteSuccess().Bind(&hook.Handler[*ModelEvent]{
		Id: systemHookIdLambdaFunction,
		Func: func(me *ModelEvent) error {
			if ef, ok := me.Model.(*LambdaFunction); ok {
				event := &LambdaFunctionEvent{
					Event:        me.Event,
					App:          me.App,
					LambdaFunction: ef,
				}
				return app.OnLambdaFunctionAfterDeleteSuccess().Trigger(event, func(e *LambdaFunctionEvent) error {
					return me.Next()
				})
			}
			return me.Next()
		},
		Priority: -99,
	})

	// Error hooks
	app.OnModelAfterCreateError().Bind(&hook.Handler[*ModelErrorEvent]{
		Id: systemHookIdLambdaFunction,
		Func: func(me *ModelErrorEvent) error {
			if ef, ok := me.Model.(*LambdaFunction); ok {
				event := &LambdaFunctionErrorEvent{
					Event:        me.Event,
					App:          me.App,
					LambdaFunction: ef,
					Error:        me.Error,
				}
				return app.OnLambdaFunctionAfterCreateError().Trigger(event, func(e *LambdaFunctionErrorEvent) error {
					return me.Next()
				})
			}
			return me.Next()
		},
		Priority: -99,
	})

	app.OnModelAfterUpdateError().Bind(&hook.Handler[*ModelErrorEvent]{
		Id: systemHookIdLambdaFunction,
		Func: func(me *ModelErrorEvent) error {
			if ef, ok := me.Model.(*LambdaFunction); ok {
				event := &LambdaFunctionErrorEvent{
					Event:        me.Event,
					App:          me.App,
					LambdaFunction: ef,
					Error:        me.Error,
				}
				return app.OnLambdaFunctionAfterUpdateError().Trigger(event, func(e *LambdaFunctionErrorEvent) error {
					return me.Next()
				})
			}
			return me.Next()
		},
		Priority: -99,
	})

	app.OnModelAfterDeleteError().Bind(&hook.Handler[*ModelErrorEvent]{
		Id: systemHookIdLambdaFunction,
		Func: func(me *ModelErrorEvent) error {
			if ef, ok := me.Model.(*LambdaFunction); ok {
				event := &LambdaFunctionErrorEvent{
					Event:        me.Event,
					App:          me.App,
					LambdaFunction: ef,
					Error:        me.Error,
				}
				return app.OnLambdaFunctionAfterDeleteError().Trigger(event, func(e *LambdaFunctionErrorEvent) error {
					return me.Next()
				})
			}
			return me.Next()
		},
		Priority: -99,
	})
}

// OnLambdaFunctionValidate returns the OnLambdaFunctionValidate hook.
func (app *BaseApp) OnLambdaFunctionValidate() *hook.Hook[*LambdaFunctionEvent] {
	return app.onLambdaFunctionValidate
}

// OnLambdaFunctionCreate returns the OnLambdaFunctionCreate hook.
func (app *BaseApp) OnLambdaFunctionCreate() *hook.Hook[*LambdaFunctionEvent] {
	return app.onLambdaFunctionCreate
}

// OnLambdaFunctionUpdate returns the OnLambdaFunctionUpdate hook.
func (app *BaseApp) OnLambdaFunctionUpdate() *hook.Hook[*LambdaFunctionEvent] {
	return app.onLambdaFunctionUpdate
}

// OnLambdaFunctionDelete returns the OnLambdaFunctionDelete hook.
func (app *BaseApp) OnLambdaFunctionDelete() *hook.Hook[*LambdaFunctionEvent] {
	return app.onLambdaFunctionDelete
}

// OnLambdaFunctionAfterCreateSuccess returns the OnLambdaFunctionAfterCreateSuccess hook.
func (app *BaseApp) OnLambdaFunctionAfterCreateSuccess() *hook.Hook[*LambdaFunctionEvent] {
	return app.onLambdaFunctionAfterCreateSuccess
}

// OnLambdaFunctionAfterUpdateSuccess returns the OnLambdaFunctionAfterUpdateSuccess hook.
func (app *BaseApp) OnLambdaFunctionAfterUpdateSuccess() *hook.Hook[*LambdaFunctionEvent] {
	return app.onLambdaFunctionAfterUpdateSuccess
}

// OnLambdaFunctionAfterDeleteSuccess returns the OnLambdaFunctionAfterDeleteSuccess hook.
func (app *BaseApp) OnLambdaFunctionAfterDeleteSuccess() *hook.Hook[*LambdaFunctionEvent] {
	return app.onLambdaFunctionAfterDeleteSuccess
}

// OnLambdaFunctionAfterCreateError returns the OnLambdaFunctionAfterCreateError hook.
func (app *BaseApp) OnLambdaFunctionAfterCreateError() *hook.Hook[*LambdaFunctionErrorEvent] {
	return app.onLambdaFunctionAfterCreateError
}

// OnLambdaFunctionAfterUpdateError returns the OnLambdaFunctionAfterUpdateError hook.
func (app *BaseApp) OnLambdaFunctionAfterUpdateError() *hook.Hook[*LambdaFunctionErrorEvent] {
	return app.onLambdaFunctionAfterUpdateError
}

// OnLambdaFunctionAfterDeleteError returns the OnLambdaFunctionAfterDeleteError hook.
func (app *BaseApp) OnLambdaFunctionAfterDeleteError() *hook.Hook[*LambdaFunctionErrorEvent] {
	return app.onLambdaFunctionAfterDeleteError
}

// OnLambdaFunctionBeforeExecute returns the OnLambdaFunctionBeforeExecute hook.
func (app *BaseApp) OnLambdaFunctionBeforeExecute() *hook.Hook[*LambdaFunctionExecuteEvent] {
	return app.onLambdaFunctionBeforeExecute
}

// OnLambdaFunctionAfterExecute returns the OnLambdaFunctionAfterExecute hook.
func (app *BaseApp) OnLambdaFunctionAfterExecute() *hook.Hook[*LambdaFunctionExecuteEvent] {
	return app.onLambdaFunctionAfterExecute
}