// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package base

import (
	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/logy"
)

type IBaseRepositoryWithHardDelete[ENTITY domain.IRootEntity] interface {
	IBaseRepository[ENTITY]
	FindByKey(requestSession *logy.RequestSession, key string, optimized bool) ENTITY
	FindByKeys(requestSession *logy.RequestSession, keys []string, optimized bool) []ENTITY
	FindAll(requestSession *logy.RequestSession, optimized bool) []ENTITY
	FindAllKeys(requestSession *logy.RequestSession) []string
	CountAll(requestSession *logy.RequestSession) int
	ExistByKey(requestSession *logy.RequestSession, key string) bool
}
type BaseRepositoryWithHardDelete[ENTITY domain.IRootEntity] struct {
	BaseRepository[ENTITY]
}

func CreateRepositoryWithHardDelete[ENTITY domain.IRootEntity](
	requestSession *logy.RequestSession,
	collectionName string,
	entityCreator func() ENTITY,
	preDelete func(*logy.RequestSession, ENTITY),
	optimizeUnsetQuery string,
	optimizeUnsetAttributes []string,
	indexes [][]string) BaseRepositoryWithHardDelete[ENTITY] {

	repo := BaseRepositoryWithHardDelete[ENTITY]{
		BaseRepository: BaseRepository[ENTITY]{
			CollectionName:          collectionName,
			EntityCreator:           entityCreator,
			PreDelete:               preDelete,
			OptimizeUnsetAttributes: optimizeUnsetAttributes,
		},
	}

	repo.Database = NewDatabase()
	repo.Database.Init(requestSession, collectionName, append(defaultIndexes, indexes...))
	return repo
}

func (repo *BaseRepositoryWithHardDelete[ENTITY]) FindByKey(requestSession *logy.RequestSession, key string, optimized bool) ENTITY {
	return repo.findByKey(requestSession, key, optimized, false)
}
func (repo *BaseRepositoryWithHardDelete[ENTITY]) FindByKeys(requestSession *logy.RequestSession, keys []string, optimized bool) []ENTITY {
	return repo.findByKeys(requestSession, keys, optimized, false)
}

func (repo *BaseRepositoryWithHardDelete[ENTITY]) FindByKeyWithDeleted(requestSession *logy.RequestSession, key string, optimized bool) ENTITY {
	return repo.findByKey(requestSession, key, optimized, true)
}

func (repo *BaseRepositoryWithHardDelete[ENTITY]) FindAll(requestSession *logy.RequestSession, optimized bool) []ENTITY {
	return repo.findAll(requestSession, optimized, false)
}

func (repo *BaseRepositoryWithHardDelete[ENTITY]) FindAllKeys(requestSession *logy.RequestSession) []string {
	return repo.findAllKeys(requestSession, false)
}

func (repo *BaseRepositoryWithHardDelete[ENTITY]) CountAll(requestSession *logy.RequestSession) int {
	return repo.countAll(requestSession, false)
}

func (repo *BaseRepositoryWithHardDelete[ENTITY]) ExistByKey(requestSession *logy.RequestSession, key string) bool {
	return repo.existByKey(requestSession, key, false)
}
