package main

import (
	"encoding/json"
	"fmt"
)

func main() {

	var j = `[{"x": 1, "y": 2.0}, {"x": 3, "y": 4.0}]`
	var c interface{}

	json.Unmarshal([]byte(j), &c)

	switch vt := c.(type) {
	case map[string]interface{}:
		fmt.Println("map")
	case []interface{}:
		fmt.Println("arr")
	default:
		fmt.Println(vt)
	}

	x, ok := c.([]interface{})

	fmt.Println("convert ok ", ok)

	for i, v := range x {
		fmt.Println(i, v)
		y := v.(map[string]interface{})
		switch y["x"].(type) {
		case int:
			fmt.Println("int", y["x"].(int))
		case float32:
			fmt.Println("float32", y["x"].(float32))
		case float64:
			fmt.Println("float64", y["x"].(float64))
		default:
			fmt.Println(y["x"])
		}

		switch y["y"].(type) {
		case int:
			fmt.Println("int", y["x"].(int))
		case float32:
			fmt.Println("float32", y["x"].(float32))
		case float64:
			fmt.Println("float64", y["x"].(float64))
		}
	}

	//	var j2 = []string{`[1, 2, 3, 4]`, `["1", "2", "3"]`,
	//	`[true, false, true]`, `[null, null, null]`}
	//
	//	for _, v := range j2 {
	//		json.Unmarshal([]byte(v), &c)
	//		fmt.Println(c)
	//
	//
	//		fmt.Println(f)
	//		switch reflect.TypeOf(c).Elem().Kind() {
	//		case reflect.Float64:
	//			fmt.Println("[]float64")
	//		case reflect.String:
	//			fmt.Println("[]string")
	//		case reflect.Bool:
	//			fmt.Println("[]bool")
	//		case reflect.Invalid:
	//			fmt.Println("[]nil")
	//		default:
	//			fmt.Println(reflect.TypeOf(c).Elem().Kind())
	//			fmt.Println(reflect.TypeOf(c).Elem().Name())
	//		}
	//	}

	var j3 = `{"x": null}`

	e := json.Unmarshal([]byte(j3), &c)
	fmt.Println(e)

	fmt.Println(c)
	v := c.(map[string]interface{})
	switch v["x"].(type) {
	case nil:
		fmt.Println("type is nil")
	default:
		fmt.Println(v["x"])
	}
}
