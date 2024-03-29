package lunatico

import (
	"errors"
	"testing"

	"github.com/aarzilli/golua/lua"
)

func TestBasic(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.OpenLibs()

	SetGlobals(L, map[string]interface{}{
		"fromgo": 10,
	})

	err := L.DoString(`
      togo = fromgo * 3
      emptyobject = {}
      emptyarray = {__emptyarray = 'yuhu'}
    `)
	if err != nil {
		t.Errorf("Execution error: %s", err)
	}

	values := GetGlobals(L, "togo", "emptyobject", "emptyarray")
	if v, ok := values["togo"].(float64); !ok {
		t.Errorf("togo is not a number")
	} else if v != 30 {
		t.Errorf("Got wrong value %f, wanted 30", v)
	}
	if _, ok := values["emptyobject"].(map[string]interface{}); !ok {
		t.Errorf("emptyobject is not an object")
	}
	if _, ok := values["emptyarray"].([]interface{}); !ok {
		t.Errorf("emptyarray is not an array")
	}
}

func TestFunctions(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.OpenLibs()

	SetGlobals(L, map[string]interface{}{
		"multiply": func(v int, times int) int { return v * times },
		"sum": func(xx ...interface{}) (res uint64) {
			for _, x := range xx {
				switch v := x.(type) {
				case float64:
					res += uint64(v)
				case map[string]interface{}:
					for _, item := range v {
						if n, ok := item.(float64); ok {
							res += uint64(n)
						}
					}
				case []interface{}:
					for _, item := range v {
						if n, ok := item.(float64); ok {
							res += uint64(n)
						}
					}
				}
			}
			return res
		},
		"multi": func() (string, string, string) { return "a", "b", "c" },
		"check_one": func(one int) (int, error) {
			if one == 1 {
				return 1, nil
			} else {
				return 0, errors.New("not one")
			}
		},
		"returns_nil": func() map[string]interface{} { return nil },
		"returns_map": func() map[string]int { return map[string]int{"a": 2} },
	})

	err := L.DoString(`
      v1 = multiply(12, 12)
      v2 = sum({7, 8, 9})
      v3 = sum({a=7, b=8, c=9})
      v4 = sum(3, 4, {b=8}, {9})
      f, s, t = multi()
      v5 = {f, s, t}
      one, err = check_one(0)
      v6 = {one, err}
      one, err = check_one(1)
      v7 = {one, err}
      v8 = function () return 'x' end
      v9 = not returns_nil()
      v10 = returns_nil()
      v11 = returns_map()
    `)
	if err != nil {
		t.Errorf("Execution error: %s", err)
	}

	values := GetGlobals(L,
		"v1", "v2", "v3", "v4", "v5", "v6", "v7", "v8", "v9", "v10",
		"v11")

	if v, ok := values["v1"].(float64); !ok {
		t.Errorf("v1 is not a number")
	} else if v != 144 {
		t.Errorf("got wrong value %f, wanted 144", v)
	}

	if v, ok := values["v2"].(float64); !ok {
		t.Errorf("v2 is not a number")
	} else if v != 24 {
		t.Errorf("got wrong value %f, wanted 24", v)
	}

	if v, ok := values["v3"].(float64); !ok {
		t.Errorf("v3 is not a number")
	} else if v != 24 {
		t.Errorf("got wrong value %f, wanted 24", v)
	}

	if v, ok := values["v4"].(float64); !ok {
		t.Errorf("v4 is not a number")
	} else if v != 24 {
		t.Errorf("got wrong value %f, wanted 24", v)
	}

	if v, ok := values["v5"].([]interface{}); !ok {
		t.Errorf("v5 is not an array")
	} else {
		if v[0].(string) != "a" {
			t.Errorf("v5.1 is %s, wanted %s", v[0].(string), "a")
		}
		if v[1].(string) != "b" {
			t.Errorf("v5.2 is %s, wanted %s", v[0].(string), "b")
		}
		if v[2].(string) != "c" {
			t.Errorf("v5.3 is %s, wanted %s", v[0].(string), "c")
		}
	}

	if v, ok := values["v6"].([]interface{}); !ok {
		t.Errorf("v6 is not an array")
	} else {
		if v[0].(float64) != 0 {
			t.Errorf("v6.one should be zero, got %v", v[0].(float64))
		}
		if v[1].(string) == "" {
			t.Errorf("v6.err should be a non-empty string, got %s", v[1].(string))
		}
	}

	if v, ok := values["v7"].([]interface{}); !ok {
		t.Errorf("v7 is not an array")
	} else {
		if v[0].(float64) != 1 {
			t.Errorf("v7.one should be 1, got %v", v[0].(float64))
		}
		if len(v) > 1 {
			t.Errorf("v7.err should be nothing, got %s", v[1])
		}
	}

	if _, ok := values["v8"].(*LuaFunction); !ok {
		t.Errorf("v8 is not a *LuaFunction")
	}

	if values["v9"] != true {
		t.Errorf("v9 should be true, but it's %v", values["v9"])
	}

	if values["v10"] != nil {
		t.Errorf("v10 should be nil, but it's %v", values["v10"])
	}

	if m, ok := values["v11"].(map[string]interface{}); !ok {
		t.Errorf("v11 should be a map, but it's %v", values["v11"])
	} else if m["a"] != float64(2) {
		t.Errorf("v11[\"a\"] should be a 2, but it's %v", m["a"])
	}
}

func TestSomeValues(t *testing.T) {
	L := lua.NewState()
	defer L.Close()
	L.OpenLibs()

	SetGlobals(L, map[string]interface{}{
		"x": map[string]interface{}{
			"k": []float64{1.123, 8, 999999999999},
			"l": uint(77),
		},
		"y": []string{"y"},
	})

	err := L.DoString("lalala = 121211")
	if err != nil {
		t.Errorf("Execution error: %s", err)
	}

	values := GetGlobals(L, "x", "y", "z")
	x := values["x"].(map[string]interface{})
	if x["l"] != float64(77) {
		t.Errorf("%v != %v", x["l"], 77)
	}
	y := values["y"].([]interface{})
	if y[0].(string) != "y" {
		t.Errorf("%v != %v", y[0], "y")
	}
}
