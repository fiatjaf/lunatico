package lunatico

import (
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
    `)
	if err != nil {
		t.Errorf("Execution error: %s", err)
	}

	values := GetGlobals(L, "togo")
	if v, ok := values["togo"].(float64); !ok {
		t.Errorf("togo is not a number")
	} else if v != 30 {
		t.Errorf("Got wrong value %f, wanted 30", v)
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
	})

	err := L.DoString(`
      v1 = multiply(12, 12)
      v2 = sum({7, 8, 9})
      v3 = sum({a=7, b=8, c=9})
      v4 = sum(3, 4, {b=8}, {9})
    `)
	if err != nil {
		t.Errorf("Execution error: %s", err)
	}

	values := GetGlobals(L, "v1", "v2", "v3", "v4")

	if v, ok := values["v1"].(float64); !ok {
		t.Errorf("v1 is not a number")
	} else if v != 144 {
		t.Errorf("Got wrong value %f, wanted 144", v)
	}

	if v, ok := values["v2"].(float64); !ok {
		t.Errorf("v2 is not a number")
	} else if v != 24 {
		t.Errorf("Got wrong value %f, wanted 24", v)
	}

	if v, ok := values["v3"].(float64); !ok {
		t.Errorf("v3 is not a number")
	} else if v != 24 {
		t.Errorf("Got wrong value %f, wanted 24", v)
	}

	if v, ok := values["v4"].(float64); !ok {
		t.Errorf("v4 is not a number")
	} else if v != 24 {
		t.Errorf("Got wrong value %f, wanted 24", v)
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
