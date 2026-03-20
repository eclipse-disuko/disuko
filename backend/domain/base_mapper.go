// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package domain

import "time"

type IBase interface {
	GetRef() string
	SetRef(ref string)

	GetKey() string
	SetKey(key string)

	GetCreated() time.Time
	SetCreated(created time.Time)

	GetUpdated() time.Time
	SetUpdated(updated time.Time)
}

func SetBaseValues(source IBase, target IBase) {
	target.SetRef(source.GetRef())
	target.SetKey(source.GetKey())
	target.SetCreated(source.GetCreated())
	target.SetUpdated(source.GetUpdated())
}

func MapTo[SOURCE interface{}, TARGET interface{}](sourceList []*SOURCE, mapFunc func(source *SOURCE) *TARGET) []*TARGET {
	result := make([]*TARGET, 0)
	for _, source := range sourceList {
		target := mapFunc(source)
		result = append(result, target)
	}
	return result
}

func MapToLimit[SOURCE interface{}, TARGET interface{}](sourceList []*SOURCE, mapFunc func(source *SOURCE) *TARGET, limit int) []*TARGET {
	result := make([]*TARGET, 0)
	count := len(sourceList)
	if limit > 0 && limit < count {
		count = limit
	}
	for i := 0; i < count; i++ {
		target := mapFunc(sourceList[i])
		result = append(result, target)
	}
	return result
}

type ConvertableEntity[T any] interface {
	ToDto() T
}

func ToDtos[S ConvertableEntity[T], T any](l []S) []T {
	res := make([]T, 0, len(l))
	for _, s := range l {
		res = append(res, s.ToDto())
	}
	return res
}

type ConvertableDto[T any] interface {
	ToEntity() T
}

func ToEntities[S ConvertableDto[T], T any](l []S) []T {
	res := make([]T, 0, len(l))
	for _, s := range l {
		res = append(res, s.ToEntity())
	}
	return res
}
