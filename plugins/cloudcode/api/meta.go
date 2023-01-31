package api

import (
	"github.com/pocketbase/pocketbase"
	lua "github.com/yuin/gopher-lua"
)

func getMetaModule(l *lua.LState) *lua.LTable {
	m := l.NewTable()
	l.RawSet(m, lua.LString("version"), lua.LString(pocketbase.Version))
	return m
}
