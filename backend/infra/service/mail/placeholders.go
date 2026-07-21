// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package mail

import (
	"fmt"
	"reflect"
)

func FillWithPlaceholders(v any) any {
	val := reflect.New(reflect.TypeOf(v)).Elem()
	fillPlaceholder(val.Type().Name(), val)
	return val.Interface()
}

func fillPlaceholder(name string, v reflect.Value) {
	switch v.Kind() {
	case reflect.String:
		v.SetString(fmt.Sprintf("[%s Placeholder]", name))
	case reflect.Bool:
		v.SetBool(true)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.SetInt(1)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v.SetUint(1)
	case reflect.Float32, reflect.Float64:
		v.SetFloat(1)
	case reflect.Pointer:
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		fillPlaceholder(name, v.Elem())
	case reflect.Slice:
		elem := reflect.New(v.Type().Elem()).Elem()
		fillPlaceholder(name, elem)
		slice := reflect.MakeSlice(v.Type(), 1, 1)
		slice.Index(0).Set(elem)
		v.Set(slice)
	case reflect.Struct:
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if !field.IsExported() {
				continue
			}
			fillPlaceholder(field.Name, v.Field(i))
		}
	}
}
