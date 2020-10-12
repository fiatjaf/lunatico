package aarzilli

import (
	"testing"

	"github.com/aarzilli/golua/lua"
	"github.com/fiatjaf/lunatico"
)

func TestBasic(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.OpenLibs()
	l := LuaState{State: L}
	lunatico.RunTestBasic(t, l)
}

func TestFunctions(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.OpenLibs()
	l := LuaState{State: L}
	lunatico.RunTestFunctions(t, l)
}

func TestSomeValues(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.OpenLibs()
	l := LuaState{State: L}
	lunatico.RunTestSomeValues(t, l)
}
