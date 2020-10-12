package azure

import (
	"github.com/Azure/golua/lua"
	"github.com/fiatjaf/lunatico"
)

// Wrapper type.

type LuaState struct {
	*lua.State
}

var azure2lua = [9]int{
	lua.NoneType:     lunatico.LUA_TNONE,
	lua.NilType:      lunatico.LUA_TNIL,
	lua.BoolType:     lunatico.LUA_TBOOLEAN,
	lua.NumberType:   lunatico.LUA_TNUMBER,
	lua.StringType:   lunatico.LUA_TSTRING,
	lua.FuncType:     lunatico.LUA_TFUNCTION,
	lua.UserDataType: lunatico.LUA_TUSERDATA,
	lua.ThreadType:   lunatico.LUA_TTHREAD,
	lua.TableType:    lunatico.LUA_TTABLE,
}

func (L LuaState) Type(index int) int {
	return azure2lua[L.State.TypeAt(index)]
}

func (L LuaState) PushGoFunction(f func(lunatico.LuaState) int) {
	wrapper := func(L *lua.State) int {
		l := LuaState{State: L}
		return f(l)
	}
	L.State.PushClosure(wrapper, 0)
}

func (L LuaState) CreateTable(narr int, nrec int) {
	L.State.NewTableSize(narr, nrec)
}

func (L LuaState) GetGlobal(name string) {
	L.State.GetGlobal(name)
}

func (L LuaState) GetTop() int {
	return L.State.Top()
}

func (L LuaState) ObjLen(index int) uint {
	length := L.State.RawLen(index)

	// Cut "holes" from the end of the table.
	for length > 0 {
		typ := L.State.RawGetIndex(index, length)
		L.State.Pop()
		if typ != lua.NoneType && typ != lua.NilType {
			break
		}
		length--
	}

	return uint(length)
}

func (L LuaState) Pop(n int) {
	L.State.PopN(n)
}

func (L LuaState) PushString(str string) {
	L.State.Push(str)
}

func (L LuaState) PushNumber(n float64) {
	L.State.Push(n)
}

func (L LuaState) PushBoolean(b bool) {
	L.State.Push(b)
}

func (L LuaState) PushNil() {
	L.State.Push(nil)
}

func (L LuaState) RaiseError(msg string) {
	L.State.Errorf("%s", msg)
}

func (L LuaState) RawSeti(index int, n int) {
	L.State.RawSetIndex(index, n)
}

func (L LuaState) ToBoolean(index int) bool {
	return L.State.ToBool(index)
}

func (L LuaState) ToInteger(index int) int {
	return int(L.State.ToInt(index))
}

func (L LuaState) DoString(str string) error {
	return L.State.ExecText(str)
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
