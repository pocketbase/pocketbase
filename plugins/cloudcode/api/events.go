package api

import (
	lua "github.com/yuin/gopher-lua"
	"reflect"
)

func getEventsModule(l *lua.LState) *lua.LTable {
	m := l.SetFuncs(l.NewTable(), buildEventsExports())
	return m
}

// this function was a PITA to write. somebody smarter than me may want to refactor it a bit.
func buildEventsExports() map[string]lua.LGFunction {
	exports := make(map[string]lua.LGFunction)

	app := reflect.TypeOf(getApp())
	reflectedApp := reflect.ValueOf(getApp())

	// Export each On* method for Lua code. In Lua, the method name is lower camel cased.
	for i := 0; i < app.NumMethod(); i++ {
		method := app.Method(i)
		if method.Name[0:2] == "On" {
			hook := reflectedApp.MethodByName(method.Name).Call([]reflect.Value{})
			m := hook[0].MethodByName("Add")

			handlerType := m.Type().In(0)

			exports["o"+method.Name[1:]] = func(l *lua.LState) int {
				f := l.ToFunction(1)

				handler := reflect.MakeFunc(handlerType, func(args []reflect.Value) []reflect.Value {
					// TODO: pass event object to lua call somehow.
					err := l.CallByParam(lua.P{
						Fn:      f,
						NRet:    0,
						Protect: true,
					})

					return []reflect.Value{reflect.ValueOf(&err).Elem()}
				})

				m.Call([]reflect.Value{handler})

				return 0
			}
		}
	}

	return exports
}
