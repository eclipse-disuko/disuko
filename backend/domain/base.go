// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package domain

import (
	"time"

	"github.com/google/uuid"
)

type IRootEntity interface {
	IChildEntity
	GetRef() string
	SetRef(ref string)
}

type IChildEntity interface {
	IBase
	IsOptimized() bool
	SetOptimized(optimized bool)
}

type ISoftDelete interface {
	IsDeleted() bool
	SetDeleted(deleted bool)
}

type RootEntity struct {
	ChildEntity `bson:",inline"`
	Rev         string `json:"_rev"`
}

func NewRootEntity() RootEntity {
	return SetRootEntity(uuid.NewString())
}

func NewRootEntityWithKey(key string) RootEntity {
	return SetRootEntity(key)
}

func SetRootEntity(uuid string) RootEntity {
	return RootEntity{ChildEntity: SetChildEntity(uuid)}
}

func NewChildEntity() ChildEntity {
	return SetChildEntity(uuid.NewString())
}

func SetChildEntity(uuid string) ChildEntity {
	return ChildEntity{Key: uuid, Created: time.Now(), Updated: time.Now()}
}

type ChildEntity struct {
	Key       string `bson:"_id" json:"_key" validate:"lte=36"`
	Created   time.Time
	Updated   time.Time
	Optimized bool `bson:"-" json:"-"`
}

type BaseState struct {
	IsRunning bool      `default:"false" json:"isRunning"`
	ReqID     string    `json:"reqID"`
	StartTime time.Time `json:"startTime"`
	EndTime   time.Time `json:"endTime"`
}

func (entity *ChildEntity) IsOptimized() bool {
	return entity.Optimized
}

func (entity *ChildEntity) SetOptimized(optimized bool) {
	entity.Optimized = optimized
}

func (entity *ChildEntity) GetKey() string {
	return entity.Key
}

func (entity *ChildEntity) SetKey(key string) {
	entity.Key = key
}

func (entity *ChildEntity) GetCreated() time.Time {
	return entity.Created
}

func (entity *ChildEntity) SetCreated(created time.Time) {
	entity.Created = created
}

func (entity *ChildEntity) GetUpdated() time.Time {
	return entity.Updated
}

func (entity *ChildEntity) SetUpdated(updated time.Time) {
	entity.Updated = updated
}

func (entity *RootEntity) GetRef() string {
	return entity.Rev
}

func (entity *RootEntity) SetRef(ref string) {
	entity.Rev = ref
}

type SoftDelete struct {
	Deleted bool
}

func (entity *SoftDelete) IsDeleted() bool {
	return entity.Deleted
}

func (entity *SoftDelete) SetDeleted(deleted bool) {
	entity.Deleted = deleted
}

func IsSoftDelete(entity any) bool {
	switch entity.(type) {
	case ISoftDelete:
		return true
	default:
		return false
	}
}

func ToSoftDelete(entity any) ISoftDelete {
	// The way of go for type checking! Do you know a better way?
	switch entity.(type) {
	case ISoftDelete:
		return entity.(ISoftDelete)
	default:
		return nil
	}
}
