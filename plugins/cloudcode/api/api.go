package api

import (
	"github.com/pocketbase/pocketbase/core"
	lua "github.com/yuin/gopher-lua"
)

// _app is a reference to the core.App that was used to initialize the cloud code system.
// This is used by some modules (via getApp()) to interact with the main pocketbase application.
var _app *core.App = nil

func getApp() core.App {
	return *_app
}

func Bind(a *core.App, l *lua.LState) {
	_app = a

	root := l.NewTable()
	//l.RawSet(root, lua.LString("db"), getDbModule(l))
	l.RawSet(root, lua.LString("events"), getEventsModule(l))
	l.RawSet(root, lua.LString("meta"), getMetaModule(l))
	l.RawSet(root, lua.LString("mail"), getMailModule(l))

	l.SetGlobal("pb", root)
}
