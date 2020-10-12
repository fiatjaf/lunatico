package azure

import (
	"testing"

	"github.com/Azure/golua/lua"
	"github.com/Azure/golua/std"
	"github.com/fiatjaf/lunatico"
)

func TestBasic(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	std.Open(L)
	l := LuaState{State: L}
	lunatico.RunTestBasic(t, l)
}

func TestFunctions(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	std.Open(L)
	l := LuaState{State: L}
	lunatico.RunTestFunctions(t, l)
}

func TestSomeValues(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	std.Open(L)
	l := LuaState{State: L}
	lunatico.RunTestSomeValues(t, l)
}
