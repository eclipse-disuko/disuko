// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package sort

import "sort"

type ExtractorFunc[T any, P any] func(T) P
type LessThanFunc[P any] func(a, b P) bool
type Implementor[T any, P any] struct {
	items         []T
	extractorFunc ExtractorFunc[T, P]
	lessThanFunc  LessThanFunc[P]
}

func (impl Implementor[T, P]) Len() int { return len(impl.items) }
func (impl Implementor[T, P]) Less(i, j int) bool {
	return impl.lessThanFunc(impl.extractorFunc(impl.items[i]), impl.extractorFunc(impl.items[j]))
}
func (impl Implementor[T, P]) Swap(i, j int) {
	impl.items[i], impl.items[j] = impl.items[j], impl.items[i]
}

func StringLessThan(s1, s2 string) bool { return s1 < s2 }
func Int64LessThan(i1, i2 int64) bool   { return i1 < i2 }
func BoolLessThan(b1, b2 bool) bool     { return !b1 && b2 }

func Sort[T any, P any](items []T, extractorFunc ExtractorFunc[T, P], lessThanFunc LessThanFunc[P], asc bool) {
	if asc {
		sort.Sort(Implementor[T, P]{items: items, extractorFunc: extractorFunc, lessThanFunc: lessThanFunc})
	} else {
		sort.Sort(sort.Reverse(Implementor[T, P]{items: items, extractorFunc: extractorFunc, lessThanFunc: lessThanFunc}))
	}
}
