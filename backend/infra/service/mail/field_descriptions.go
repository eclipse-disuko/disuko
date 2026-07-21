// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package mail

import (
	"reflect"
)

func FieldDescriptions(v interface{}) map[string]string {
	result := make(map[string]string)
	collectFieldDescriptions(reflect.TypeOf(v), "", result)
	return result
}

func collectFieldDescriptions(t reflect.Type, prefix string, result map[string]string) {
	if t == nil {
		return
	}

	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}

		path := prefix + field.Name
		fieldType := field.Type
		isSlice := fieldType.Kind() == reflect.Slice

		if isSlice {
			path += "[]"
			fieldType = fieldType.Elem()
		}

		for fieldType.Kind() == reflect.Pointer {
			fieldType = fieldType.Elem()
		}

		if fieldType.Kind() == reflect.Struct {
			collectFieldDescriptions(fieldType, path+".", result)
			continue
		}

		result[path] = field.Tag.Get("description")
	}
}
