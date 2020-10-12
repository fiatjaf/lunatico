package aarzilli

import (
	"github.com/aarzilli/golua/lua"
	"github.com/fiatjaf/lunatico"
)

// Wrapper type.

type LuaState struct {
	*lua.State
}

func (L LuaState) Type(index int) int {
	return int(L.State.Type(index))
}

func (L LuaState) PushGoFunction(f func(lunatico.LuaState) int) {
	L.State.PushGoFunction(func(L *lua.State) int {
		l := LuaState{State: L}
		return f(l)
	})
}

func (L LuaState) Next(index int) bool {
	return L.State.Next(index) != 0
}

// Wrapper functions to use as drop-in replacement for old lunatico.

func SetGlobals(L *lua.State, globals map[string]interface{}) {
	l := LuaState{State: L}
	lunatico.SetGlobals(l, globals)
}

func GetGlobals(L *lua.State, names ...string) map[string]interface{} {
	l := LuaState{State: L}
	return lunatico.GetGlobals(l, names...)
}

func GetFullStack(L *lua.State) []interface{} {
	l := LuaState{State: L}
	return lunatico.GetFullStack(l)
}

func ReadAny(L *lua.State, pos int) interface{} {
	l := LuaState{State: L}
	return lunatico.ReadAny(l, pos)
}

func ReadString(L *lua.State, pos int) (v string) {
	l := LuaState{State: L}
	return lunatico.ReadString(l, pos)
}

func ReadTable(L *lua.State, pos int) interface{} {
	l := LuaState{State: L}
	return lunatico.ReadTable(l, pos)
}

func PushMap(L *lua.State, m map[string]interface{}) {
	l := LuaState{State: L}
	lunatico.PushMap(l, m)
}

func PushSlice(L *lua.State, s []interface{}) {
	l := LuaState{State: L}
	lunatico.PushSlice(l, s)
}

func PushAny(L *lua.State, ival interface{}) {
	l := LuaState{State: L}
	lunatico.PushAny(l, ival)
}
