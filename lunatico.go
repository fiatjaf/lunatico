package lunatico

import (
	"fmt"
	"reflect"
)

type LuaState interface {
	SetGlobal(name string)
	GetGlobal(name string)
	Pop(n int)
	GetTop() int
	Type(index int) int
	ToBoolean(index int) bool
	ToString(index int) string
	ToNumber(index int) float64
	ToInteger(index int) int
	ObjLen(index int) uint
	Next(index int) bool
	CreateTable(narr int, nrec int)
	RawSet(index int)
	RawSeti(index int, n int)
	PushGoFunction(f func(LuaState) int)
	RaiseError(msg string)
	ArgError(narg int, extramsg string) int
	PushString(str string)
	PushNumber(n float64)
	PushBoolean(b bool)
	PushNil()
}

// From lua.h.
const (
	LUA_TNONE          = (-1)
	LUA_TNIL           = 0
	LUA_TBOOLEAN       = 1
	LUA_TLIGHTUSERDATA = 2
	LUA_TNUMBER        = 3
	LUA_TSTRING        = 4
	LUA_TTABLE         = 5
	LUA_TFUNCTION      = 6
	LUA_TUSERDATA      = 7
	LUA_TTHREAD        = 8
	LUA_NUMTAGS        = 9
)

// utils
func SetGlobals(L LuaState, globals map[string]interface{}) {
	for k, v := range globals {
		PushAny(L, v)
		L.SetGlobal(k)
	}
}

func GetGlobals(L LuaState, names ...string) map[string]interface{} {
	globals := make(map[string]interface{})
	for _, name := range names {
		L.GetGlobal(name)
		v := ReadAny(L, -1)
		globals[name] = v
		L.Pop(1)
	}
	return globals
}

func GetFullStack(L LuaState) []interface{} {
	tip := L.GetTop()
	values := make([]interface{}, tip)
	for i := 1; i <= tip; i++ {
		v := ReadAny(L, i)
		values[i-1] = v
	}
	return values
}

// read stuff
func ReadAny(L LuaState, pos int) interface{} {
	switch L.Type(pos) {
	case LUA_TNIL:
		return nil
	case LUA_TNUMBER:
		return L.ToNumber(pos)
	case LUA_TBOOLEAN:
		return L.ToBoolean(pos)
	case LUA_TSTRING:
		return L.ToString(pos)
	case LUA_TTABLE:
		return ReadTable(L, pos)
	case LUA_TFUNCTION:
		return nil
	}
	return nil
}

func ReadString(L LuaState, pos int) (v string) {
	switch L.Type(pos) {
	case LUA_TNUMBER:
		return fmt.Sprint(L.ToNumber(pos))
	case LUA_TBOOLEAN:
		return fmt.Sprint(L.ToBoolean(pos))
	case LUA_TSTRING:
		return L.ToString(pos)
	}
	return ""
}

func ReadTable(L LuaState, pos int) interface{} {
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

	for L.Next(pos) {
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
func PushMap(L LuaState, m map[string]interface{}) {
	L.CreateTable(0, len(m))
	for k, v := range m {
		PushAny(L, k)
		PushAny(L, v)
		L.RawSet(-3)
	}
}

func PushSlice(L LuaState, s []interface{}) {
	L.CreateTable(len(s), 0)
	for i, v := range s {
		PushAny(L, v)
		L.RawSeti(-2, i+1)
	}
}

func PushAny(L LuaState, ival interface{}) {
	if ival == nil {
		L.PushNil()
		return
	}

	rv := reflect.ValueOf(ival)
	switch rv.Kind() {
	case reflect.Func:
		L.PushGoFunction(func(L LuaState) int {
			fnType := rv.Type()

			fnArgs := fnType.NumIn()           // includes a potential variadic argument
			givenArgs := L.GetTop()            // args passed to function
			variadic := rv.Type().IsVariadic() // means the last argument is ...

			var numArgs int
			if variadic {
				// when variadic we can ignore the last argument
				// or accept many of it
				if givenArgs+1 >= fnArgs {
					numArgs = givenArgs
				} else {
					numArgs = fnArgs
				}
			} else {
				// function is limited to the number of fnArgs
				numArgs = fnArgs
			}

			// when it's less there's nothing we can do
			if numArgs > givenArgs {
				L.RaiseError(fmt.Sprintf("got %d arguments, needed %d", numArgs, givenArgs))
			}

			args := make([]reflect.Value, numArgs)
			for i := 0; i < numArgs; i++ {
				arg := ReadAny(L, i+1)
				av := reflect.ValueOf(arg)

				var requiredType reflect.Type
				if i >= fnArgs-1 && variadic {
					requiredType = fnType.In(fnArgs - 1).Elem()
				} else {
					requiredType = fnType.In(i)
				}

				if av.IsValid() {
					if !av.Type().ConvertibleTo(requiredType) {
						L.ArgError(i+1, fmt.Sprintf("wrong argument type: got %s, wanted %s",
							av.Kind().String(), requiredType.Kind().String()))
					}

					av = av.Convert(requiredType)
				} else {
					av = reflect.New(requiredType).Elem()

					if !av.Type().AssignableTo(requiredType) {
						L.ArgError(i+1, fmt.Sprintf("wrong argument type: got '%s', wanted %s",
							arg, requiredType.Kind().String()))
					}
				}

				args[i] = av
			}

			defer func() {
				// recover from panics during function run
				if err := recover(); err != nil {
					L.RaiseError(fmt.Sprintf("function panic: %s", err))
				}
			}()
			returned := rv.Call(args)

			for _, ret := range returned {
				PushAny(L, ret.Interface())
			}

			return len(returned)
		})
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
	case reflect.Ptr, reflect.Struct:
		// if it has an Error() or String() method, call these instead of pushing nil.
		method, ok := rv.Type().MethodByName("Error")
		if ok {
			goto callmethod
		}
		method, ok = rv.Type().MethodByName("String")
		if ok {
			goto callmethod
		}

		goto justpushnil
	callmethod:
		if method.Type.NumIn() == 1 /* 1 because the struct itself is an argument */ &&
			method.Type.NumOut() == 1 &&
			method.Type.Out(0).Kind() == reflect.String {

			res := method.Func.Call([]reflect.Value{rv})
			L.PushString(res[0].String())
			break
		}
	justpushnil:
		L.PushNil()
	default:
		L.PushNil()
	}
}
