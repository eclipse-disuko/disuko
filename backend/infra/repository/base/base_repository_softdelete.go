// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package base

import (
	"mercedes-benz.ghe.com/foss/disuko/domain"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type IBaseRepositoryWithSoftDelete[ENTITY domain.IRootEntity] interface {
	IBaseRepository[ENTITY]
	IBaseRepositoryWithHardDelete[ENTITY]
	FindByKeyWithDeleted(requestSession *logy.RequestSession, key string, optimized bool) ENTITY
	FindAllWithDeleted(requestSession *logy.RequestSession, optimized bool) []ENTITY
	FindAllKeysWithDeleted(requestSession *logy.RequestSession) []string
	CountAllWithDeleted(requestSession *logy.RequestSession) int
	ExistByKeyWithDeleted(requestSession *logy.RequestSession, key string) bool
	DeleteHard(requestSession *logy.RequestSession, key string)
}
type BaseRepositoryWithSoftDelete[ENTITY domain.IRootEntity] struct {
	BaseRepository[ENTITY]
	BaseRepositoryWithHardDelete[ENTITY]
}

func CreateRepositoryWithSoftDelete[ENTITY domain.IRootEntity](
	requestSession *logy.RequestSession,
	collectionName string,
	entityCreator func() ENTITY,
	preDelete func(*logy.RequestSession, ENTITY),
	optimizedUnsetAttributes []string,
	indexes [][]string) BaseRepositoryWithSoftDelete[ENTITY] {

	baseRepo := BaseRepository[ENTITY]{
		CollectionName:          collectionName,
		EntityCreator:           entityCreator,
		PreDelete:               preDelete,
		OptimizeUnsetAttributes: optimizedUnsetAttributes,
		IsSoftDelete:            true,
	}
	baseRepo.Database = NewDatabase()
	baseRepo.Database.Init(requestSession, collectionName, append(defaultIndexes, indexes...))

	repo := BaseRepositoryWithSoftDelete[ENTITY]{
		BaseRepository: baseRepo,
		BaseRepositoryWithHardDelete: BaseRepositoryWithHardDelete[ENTITY]{
			BaseRepository: baseRepo,
		},
	}
	return repo
}

func (repo *BaseRepositoryWithSoftDelete[ENTITY]) FindByKeyWithDeleted(requestSession *logy.RequestSession, key string, optimized bool) ENTITY {
	return repo.findByKey(requestSession, key, optimized, true)
}

func (repo *BaseRepositoryWithSoftDelete[ENTITY]) FindAllWithDeleted(requestSession *logy.RequestSession, optimized bool) []ENTITY {
	return repo.findAll(requestSession, optimized, true)
}

func (repo *BaseRepositoryWithSoftDelete[ENTITY]) FindAllKeysWithDeleted(requestSession *logy.RequestSession) []string {
	return repo.findAllKeys(requestSession, true)
}

func (repo *BaseRepositoryWithSoftDelete[ENTITY]) CountAllWithDeleted(requestSession *logy.RequestSession) int {
	return repo.countAll(requestSession, true)
}

func (repo *BaseRepositoryWithHardDelete[ENTITY]) ExistByKeyWithDeleted(requestSession *logy.RequestSession, key string) bool {
	return repo.existByKey(requestSession, key, true)
}

func (repo *BaseRepositoryWithSoftDelete[ENTITY]) DeleteHard(requestSession *logy.RequestSession, key string) {
	entity := repo.findByKey(requestSession, key, false, false)
	if repo.PreDelete != nil {
		repo.PreDelete(requestSession, entity)
	}
	repo.Database.Delete(key)
}
