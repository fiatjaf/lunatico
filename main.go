package main

import (
	"fmt"
	"log"
	"reflect"

	"github.com/aarzilli/golua/lua"
)

func main() {
	L := lua.NewState()
	L.OpenLibs()

	var value interface{} = map[string]interface{}{
		"x": map[string]interface{}{
			"y": 2,
		},
	}

	SetGlobals(L, map[string]interface{}{
		"value":        value,
		"getvalueback": func(v interface{}) { value = v },
	})

	log.Print(L.DoString(`
value.x.y = 3
local z = value.x.z or {}
value.x.z = z
z.m = 'alskndalsdn'
local m = value.m or {}
value.m = m
m.z = {1,2,3}
getvalueback(value)
    `))

	fmt.Println(value)
}

// utils
func SetGlobals(L *lua.State, globals map[string]interface{}) {
	for k, v := range globals {
		PushAny(L, v)
		L.SetGlobal(k)
	}
}

func GetFullStack(L *lua.State) []interface{} {
	tip := L.GetTop()
	values := make([]interface{}, tip)
	for i := 1; i <= tip; i++ {
		v := ReadAny(L, i)
		values[i-1] = v
	}
	return values
}

// read stuff
func ReadAny(L *lua.State, pos int) interface{} {
	switch L.Type(pos) {
	case lua.LUA_TNIL:
		return nil
	case lua.LUA_TNUMBER:
		return L.ToNumber(pos)
	case lua.LUA_TBOOLEAN:
		return L.ToBoolean(pos)
	case lua.LUA_TSTRING:
		return L.ToString(pos)
	case lua.LUA_TTABLE:
		return ReadTable(L, pos)
	case lua.LUA_TFUNCTION:
		return nil
	}
	return nil
}

func ReadString(L *lua.State, pos int) (v string) {
	switch L.Type(pos) {
	case lua.LUA_TNUMBER:
		return fmt.Sprint(L.ToNumber(pos))
	case lua.LUA_TBOOLEAN:
		return fmt.Sprint(L.ToBoolean(pos))
	case lua.LUA_TSTRING:
		return L.ToString(pos)
	}
	return ""
}

func ReadTable(L *lua.State, pos int) interface{} {
	if pos < 0 {
		pos = L.GetTop() + 1 + pos
	}

	var object = make(map[string]interface{})
	var slice []interface{}

	isArray := true
	size := L.ObjLen(pos)
	if size == 0 {
		isArray = false
	} else {
		slice = make([]interface{}, size)
	}

	L.PushNil()

	for L.Next(pos) != 0 {
		val := ReadAny(L, -1)
		L.Pop(1)

		// array
		if isArray {
			if index := L.ToInteger(-1); index != 0 && index <= int(size) {
				slice[index-1] = val
			} else {
				isArray = false
			}
		}

		// object
		key := ReadString(L, -1)
		object[key] = val
	}

	if isArray {
		return slice
	} else {
		return object
	}
}

// push stuff
func PushMap(L *lua.State, m map[string]interface{}) {
	L.CreateTable(0, len(m))
	for k, v := range m {
		PushAny(L, k)
		PushAny(L, v)
		L.RawSet(-3)
	}
}

func PushSlice(L *lua.State, s []interface{}) {
	L.CreateTable(len(s), 0)
	for i, v := range s {
		PushAny(L, v)
		L.RawSeti(-2, i+1)
	}
}

func PushAny(L *lua.State, ival interface{}) {
	rv := reflect.ValueOf(ival)
	switch rv.Kind() {
	case reflect.String:
		L.PushString(rv.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		L.PushNumber(float64(rv.Int()))
	case reflect.Uint, reflect.Uintptr, reflect.Uint8,
		reflect.Uint16, reflect.Uint32, reflect.Uint64:
		L.PushNumber(float64(rv.Uint()))
	case reflect.Float32, reflect.Float64:
		L.PushNumber(rv.Float())
	case reflect.Bool:
		L.PushBoolean(rv.Bool())
	case reflect.Slice:
		size := rv.Len()
		slice := make([]interface{}, size)
		for i := 0; i < size; i++ {
			slice[i] = rv.Index(i).Interface()
		}
		PushSlice(L, slice)
	case reflect.Map:
		m := make(map[string]interface{}, rv.Len())
		for _, key := range rv.MapKeys() {
			m[fmt.Sprint(key)] = rv.MapIndex(key).Interface()
		}
		PushMap(L, m)
	case reflect.Func:
		L.PushGoFunction(func(L *lua.State) int {
			var args []reflect.Value
			for i := 1; !L.IsNone(i); i++ {
				arg := ReadAny(L, i)
				args = append(args, reflect.ValueOf(arg))
			}

			returned := rv.Call(args)

			for _, ret := range returned {
				PushAny(L, ret.Interface())
			}

			return len(returned)
		})
	default:
		L.PushNil()
	}
}
