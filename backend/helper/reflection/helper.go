// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package reflection

import (
	"reflect"
)

func IsSlice(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Slice
}

func IsMap(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Map
}

func IsArary(v interface{}) bool {
	return reflect.TypeOf(v).Kind() == reflect.Array
}

func ToPointer[TYPE interface{}](value TYPE) *TYPE {
	return &value
}
