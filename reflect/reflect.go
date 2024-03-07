package reflect

import (
	"fmt"
	"reflect"
)

func ReflectType(x interface{}) {
	v := reflect.TypeOf(x)
	fmt.Printf("type:%v\n", v)
}

func ReflectValue(x interface{}) {
	v := reflect.ValueOf(x)
	k := v.Kind()
	switch k {
	case reflect.Int64:
		fmt.Printf("type is int64, value is %d\n", int64(v.Int()))
	case reflect.Float32:
		fmt.Printf("type is float32, value is %f\n", float32(v.Float()))
	case reflect.Float64:
		fmt.Printf("type is float64, value is %f\n", float64(v.Float()))
	}
}
