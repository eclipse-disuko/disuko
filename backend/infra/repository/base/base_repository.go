// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package base

import (
	"errors"
	"time"

	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/infra/repository/database"
	"github.com/eclipse-disuko/disuko/infra/repository/database/couch"
	"github.com/eclipse-disuko/disuko/infra/repository/database/mongo"
	"github.com/eclipse-disuko/disuko/logy"
)

var defaultIndexes = [][]string{
	{"Created"},
	{"Updated"},
	{"Deleted"},
	{"_id", "Deleted"},
}

type IBaseRepository[ENTITY domain.IRootEntity] interface {
	Query(requestSession *logy.RequestSession, qc *database.QueryConfig) []ENTITY
	Save(requestSession *logy.RequestSession, entity ENTITY)
	SaveList(requestSession *logy.RequestSession, entities []ENTITY, silent bool)
	Update(requestSession *logy.RequestSession, entity ENTITY)
	UpdateWithoutTimestamp(requestSession *logy.RequestSession, entity ENTITY)
	UpdateList(requestSession *logy.RequestSession, entities []ENTITY)
	UpdateListSilent(requestSession *logy.RequestSession, entities []ENTITY)
	Delete(requestSession *logy.RequestSession, key string)
	DatabaseConn() IDatabase
	StartSession(st SessionType, flushThreshold int) *BulkSession[ENTITY]
}

type IBaseRepositoryForMigration[ENTITY domain.IRootEntity] interface {
	UpdateForMigration(requestSession *logy.RequestSession, entity ENTITY)
	UpdateListForMigration(requestSession *logy.RequestSession, entities []ENTITY)
}

type IDatabase interface {
	Init(rs *logy.RequestSession, tableName string, indexes [][]string)
	QueryQB(qc *database.QueryConfig, createResult func() interface{}) []interface{}
	Save(doc interface{}) (string, string)
	SaveBulk(docs []interface{}) []database.RevKeyHolder
	DeleteBulk(docs []interface{})
	Update(key string, oldRev string, doc interface{}) string
	Delete(key string)
	GetKeyAttribute() string
	Truncate()
	DropDatabase()
	CreateDatabase()
}

type BaseRepository[ENTITY domain.IRootEntity] struct {
	CollectionName string
	EntityCreator  func() ENTITY
	PreDelete      func(*logy.RequestSession, ENTITY)
	PreSave        func(*logy.RequestSession, ENTITY)
	PreUpdate      func(*logy.RequestSession, ENTITY)

	Database                IDatabase
	OptimizeUnsetAttributes []string
	IsSoftDelete            bool
}

func NewDatabase() IDatabase {
	if conf.Config.Database.Type == conf.DatabaseCouchDB {
		return &couch.DatabaseCouch{}
	} else if conf.Config.Database.Type == conf.DatabaseMongoDB {
		return &mongo.Database{}
	}

	return nil
}

type SessionType int

const (
	UpdateSession SessionType = iota
	DeleteSession
)

type BulkSession[ENTITY domain.IRootEntity] struct {
	st             SessionType
	flushThreshold int
	repo           *BaseRepository[ENTITY]
	Cache          []interface{}
}

func (repo *BaseRepository[ENTITY]) StartSession(st SessionType, flushThreshold int) *BulkSession[ENTITY] {
	session := BulkSession[ENTITY]{
		st:             st,
		flushThreshold: flushThreshold,
		repo:           repo,
	}
	session.Cache = make([]interface{}, 0)
	return &session
}

func (s *BulkSession[ENTITY]) AddEnt(entity ENTITY) {
	if s.st == UpdateSession {
		entity.SetUpdated(time.Now())
	}
	s.Cache = append(s.Cache, entity)
	if len(s.Cache) >= s.flushThreshold {
		if s.st == UpdateSession {
			s.repo.Database.SaveBulk(s.Cache)
		} else {
			s.repo.Database.DeleteBulk(s.Cache)
		}
		s.Cache = nil
	}
}

func (s *BulkSession[ENTITY]) EndSession() {
	if len(s.Cache) == 0 {
		return
	}
	if s.st == UpdateSession {
		s.repo.Database.SaveBulk(s.Cache)
	} else {
		s.repo.Database.DeleteBulk(s.Cache)
	}
	s.Cache = nil
}

func (repo *BaseRepository[ENTITY]) Query(requestSession *logy.RequestSession, qc *database.QueryConfig) []ENTITY {
	if requestSession != nil {
		requestSession.QueryCount++
	}
	start := time.Now()
	res := repo.Database.QueryQB(qc, func() interface{} {
		buf := repo.EntityCreator()
		return &buf
	})
	if requestSession != nil {
		requestSession.QueryTime += time.Since(start)
	}
	tRes := UnmarshalX[ENTITY](res)
	for _, e := range tRes {
		e.SetOptimized(qc.DoesModify())
	}
	return tRes
}

func (repo *BaseRepository[ENTITY]) Save(requestSession *logy.RequestSession, entity ENTITY) {
	if repo.PreSave != nil {
		repo.PreSave(requestSession, entity)
	}
	entity.SetCreated(time.Now())
	entity.SetUpdated(entity.GetCreated())

	if requestSession != nil {
		requestSession.QueryCount++
	}
	start := time.Now()
	key, rev := repo.Database.Save(entity)
	if requestSession != nil {
		requestSession.QueryTime += time.Since(start)
	}

	entity.SetKey(key)
	entity.SetRef(rev)
}

func (repo *BaseRepository[ENTITY]) SaveList(requestSession *logy.RequestSession, entities []ENTITY, silent bool) {
	bulkArray := make([]interface{}, 0)
	for _, entity := range entities {
		if silent == false {
			entity.SetUpdated(time.Now())
		}
		bulkArray = append(bulkArray, entity)
	}
	if requestSession != nil {
		requestSession.QueryCount++
	}
	start := time.Now()
	res := repo.Database.SaveBulk(bulkArray)
	if requestSession != nil {
		requestSession.QueryTime += time.Since(start)
	}
	if len(res) == len(entities) {
		for i, entity := range entities {
			entity.SetKey(res[i].Key)
			entity.SetRef(res[i].Rev)
		}
	}
}

func (repo *BaseRepository[ENTITY]) Update(requestSession *logy.RequestSession, entity ENTITY) {
	if repo.PreUpdate != nil {
		repo.PreUpdate(requestSession, entity)
	}
	repo.update(requestSession, entity, true)
}

func (repo *BaseRepository[ENTITY]) UpdateWithoutTimestamp(requestSession *logy.RequestSession, entity ENTITY) {
	if repo.PreUpdate != nil {
		repo.PreUpdate(requestSession, entity)
	}
	repo.update(requestSession, entity, false)
}

func (repo *BaseRepository[ENTITY]) UpdateList(requestSession *logy.RequestSession, entities []ENTITY) {
	for _, entity := range entities {
		repo.Update(requestSession, entity)
	}
}

func (repo *BaseRepository[ENTITY]) UpdateListSilent(requestSession *logy.RequestSession, entities []ENTITY) {
	for _, entity := range entities {
		repo.update(requestSession, entity, false)
	}
}

func (repo *BaseRepository[ENTITY]) UpdateForMigration(requestSession *logy.RequestSession, entity ENTITY) {
	repo.update(requestSession, entity, false)
}

func (repo *BaseRepository[ENTITY]) UpdateListForMigration(requestSession *logy.RequestSession, entities []ENTITY) {
	for _, entity := range entities {
		repo.UpdateForMigration(requestSession, entity)
	}
}

func (repo *BaseRepository[ENTITY]) update(requestSession *logy.RequestSession, entity ENTITY, setUpdateTimestamp bool) {
	if entity.IsOptimized() {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorDbTryUpdateOptimized, entity.GetKey()), "")
	}

	if setUpdateTimestamp {
		entity.SetUpdated(time.Now())
	}

	if requestSession != nil {
		requestSession.QueryCount++
	}
	start := time.Now()
	rev := repo.Database.Update(entity.GetKey(), entity.GetRef(), entity)
	if requestSession != nil {
		requestSession.QueryTime += time.Since(start)
	}
	entity.SetRef(rev)
}

func (repo *BaseRepository[ENTITY]) existByKey(requestSession *logy.RequestSession, key string, withDeleted bool) bool {
	qc := database.New().SetMatcher(
		database.AttributeMatcher(
			repo.Database.GetKeyAttribute(),
			database.EQ,
			key,
		),
	).SetKeep(
		[]string{
			repo.Database.GetKeyAttribute(),
		},
	)
	if repo.IsSoftDelete && !withDeleted {
		qc.SetMatcher(
			database.AndChain(
				*qc.Matcher,
				database.AttributeMatcher(
					"Deleted",
					database.EQ,
					false,
				),
			),
		)
	}
	qbRes := len(repo.Query(requestSession, qc)) > 0

	return qbRes
}

func (repo *BaseRepository[ENTITY]) findByKey(requestSession *logy.RequestSession, key string, optimized bool, withDeleted bool) ENTITY {
	qc := database.New()

	qc.SetMatcher(database.AttributeMatcher(
		repo.Database.GetKeyAttribute(),
		database.EQ,
		key,
	))
	if repo.IsSoftDelete && !withDeleted {
		qc.SetMatcher(database.AndChain(
			*qc.Matcher,
			database.AttributeMatcher(
				"Deleted",
				database.EQ,
				false,
			),
		))
	}
	if optimized && len(repo.OptimizeUnsetAttributes) > 0 {
		qc.SetUnset(repo.OptimizeUnsetAttributes)
	}
	qbList := repo.Query(requestSession, qc)
	var qbRes ENTITY
	if len(qbList) > 0 {
		qbRes = qbList[0]
	}

	return qbRes
}

func (repo *BaseRepository[ENTITY]) findByKeys(requestSession *logy.RequestSession, keys []string, optimized bool, withDeleted bool) []ENTITY {
	qc := database.New()

	res := qc.Matcher
	if len(keys) > 0 {
		res = &database.MatchGroup{}
		for _, key := range keys {
			g := database.AttributeMatcher(
				repo.Database.GetKeyAttribute(),
				database.EQ,
				key,
			)
			res.Operator = database.OR
			res.Chain = append(res.Chain, g)
		}
		qc.SetMatcher(*res)
	}
	if repo.IsSoftDelete && !withDeleted {
		if res == nil {
			qc.SetMatcher(database.AndChain(
				database.AttributeMatcher(
					"Deleted",
					database.EQ,
					false,
				),
			))
		} else {
			qc.SetMatcher(database.AndChain(
				*res,
				database.AttributeMatcher(
					"Deleted",
					database.EQ,
					false,
				),
			))
		}
	}

	if optimized && len(repo.OptimizeUnsetAttributes) > 0 {
		qc.SetUnset(repo.OptimizeUnsetAttributes)
	}
	qbList := repo.Query(requestSession, qc)

	return qbList
}

func (repo *BaseRepository[ENTITY]) Delete(requestSession *logy.RequestSession, key string) {
	entity := repo.findByKey(requestSession, key, false, false)
	if repo.PreDelete != nil {
		repo.PreDelete(requestSession, entity)
	}

	if repo.IsSoftDelete {
		softDelete := domain.ToSoftDelete(entity)
		softDelete.SetDeleted(true)
		repo.Update(requestSession, entity)
	} else {
		// hard delete
		if requestSession != nil {
			requestSession.QueryCount++
		}
		start := time.Now()
		repo.Database.Delete(key)
		if requestSession != nil {
			requestSession.QueryTime += time.Since(start)
		}
	}
}

func (repo *BaseRepository[ENTITY]) countAll(requestSession *logy.RequestSession, withDeleted bool) int {
	qc := database.New().SetKeep([]string{
		repo.Database.GetKeyAttribute(),
	})
	if repo.IsSoftDelete && !withDeleted {
		qc.SetMatcher(database.AttributeMatcher(
			"Deleted",
			database.EQ,
			false,
		))
	}
	return len(repo.Query(requestSession, qc))
}

func (repo *BaseRepository[ENTITY]) findAll(requestSession *logy.RequestSession, optimized bool, withDeleted bool) []ENTITY {
	qc := database.New()
	if repo.IsSoftDelete && !withDeleted {
		qc.SetMatcher(database.AttributeMatcher(
			"Deleted",
			database.EQ,
			false,
		))
	}
	if optimized && len(repo.OptimizeUnsetAttributes) > 0 {
		qc.SetUnset(repo.OptimizeUnsetAttributes)
	}

	qbRes := repo.Query(requestSession, qc)
	return qbRes
}

func (repo *BaseRepository[ENTITY]) findAllKeys(requestSession *logy.RequestSession, withDeleted bool) []string {
	qc := database.New().SetKeep([]string{
		repo.Database.GetKeyAttribute(),
	})

	if repo.IsSoftDelete && !withDeleted {
		qc.SetMatcher(database.AttributeMatcher(
			"Deleted",
			database.EQ,
			false,
		))
	}

	qbList := repo.Query(requestSession, qc)
	var qbRes []string
	for _, e := range qbList {
		qbRes = append(qbRes, e.GetKey())
	}
	return qbRes
}

func (repo *BaseRepository[ENTITY]) DatabaseConn() IDatabase {
	return repo.Database
}

func UnmarshalX[X interface{}](data []interface{}) []X {
	var entities []X

	for _, d := range data {
		entity, ok := d.(*X)
		if !ok {
			exception.HandleErrorServerMessage(errors.New("unexpected database result"), message.GetI18N(message.ErrorDbRead))
		}
		entities = append(entities, *entity)
	}
	return entities
}
