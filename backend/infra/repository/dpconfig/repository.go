// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package dpconfig

import (
	"reflect"

	"mercedes-benz.ghe.com/foss/disuko/domain"
	"mercedes-benz.ghe.com/foss/disuko/domain/dpconfig"
	"mercedes-benz.ghe.com/foss/disuko/domain/integrity"
	"mercedes-benz.ghe.com/foss/disuko/domain/notification"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type iBaseConfigRepository[ENTITY domain.IRootEntity] interface {
	Get(requestSession *logy.RequestSession) ENTITY
	Save(requestSession *logy.RequestSession, data ENTITY)
}
type baseConfigRepositoryStruct[ENTITY domain.IRootEntity] struct {
	base.BaseRepositoryWithHardDelete[ENTITY]
	StaticKey string
}

func newBaseConfigRepository[ENTITY domain.IRootEntity](requestSession *logy.RequestSession, staticKey string, entityCreator func() ENTITY) baseConfigRepositoryStruct[ENTITY] {
	return baseConfigRepositoryStruct[ENTITY]{
		BaseRepositoryWithHardDelete: base.CreateRepositoryWithHardDelete[ENTITY](
			requestSession,
			DPConfigCollectionName,
			entityCreator,
			nil,
			"",
			nil,
			nil),
		StaticKey: staticKey,
	}
}

func (lr *baseConfigRepositoryStruct[ENTITY]) Save(requestSession *logy.RequestSession, data ENTITY) {
	data.SetKey(lr.StaticKey)
	old := lr.FindByKey(requestSession, data.GetKey(), false)
	if reflect.ValueOf(old).IsNil() {
		lr.BaseRepositoryWithHardDelete.Save(requestSession, data)
	} else {
		data.SetRef(old.GetRef())
		lr.Update(requestSession, data)
	}
}

func (lr *baseConfigRepositoryStruct[ENTITY]) Get(requestSession *logy.RequestSession) ENTITY {
	data := lr.FindByKey(requestSession, lr.StaticKey, false)
	if reflect.ValueOf(data).IsNil() {
		data = lr.EntityCreator()
	}
	return data
}

const DP_DB_S3_CHECK_STATE = "DP_DB_S3_CHECK_STATE"

type integrityCheckRepositoryStruct struct {
	baseConfigRepositoryStruct[*integrity.DbIntegrityResult]
}

func newIntegrityCheckRepository(requestSession *logy.RequestSession) iBaseConfigRepository[*integrity.DbIntegrityResult] {
	return &integrityCheckRepositoryStruct{
		baseConfigRepositoryStruct: newBaseConfigRepository(requestSession, DP_DB_S3_CHECK_STATE, func() *integrity.DbIntegrityResult {
			return &integrity.DbIntegrityResult{}
		}),
	}
}

const DP_SAMPLE_DATA_STATE = "DP_SAMPLE_DATA_STATE"

type sampleDataCreationStateRepositoryStruct struct {
	baseConfigRepositoryStruct[*dpconfig.SampleDataCreationState]
}

func newSampleDataCreationRepository(requestSession *logy.RequestSession) iBaseConfigRepository[*dpconfig.SampleDataCreationState] {
	return &sampleDataCreationStateRepositoryStruct{
		baseConfigRepositoryStruct: newBaseConfigRepository(requestSession, DP_SAMPLE_DATA_STATE, func() *dpconfig.SampleDataCreationState {
			return &dpconfig.SampleDataCreationState{}
		}),
	}
}

const DP_CONFIG_KEY = "DP_CONFIG"

type dpConfigRepositoryStruct struct {
	baseConfigRepositoryStruct[*dpconfig.DPConfig]
}

func newDPConfigRepository(requestSession *logy.RequestSession) iBaseConfigRepository[*dpconfig.DPConfig] {
	return &dpConfigRepositoryStruct{
		baseConfigRepositoryStruct: newBaseConfigRepository(requestSession, DP_CONFIG_KEY, func() *dpconfig.DPConfig {
			return &dpconfig.DPConfig{}
		}),
	}
}

const DP_NOTIFICATION = "DP_NOTIFICATION"

type notificationStruct struct {
	baseConfigRepositoryStruct[*notification.Notification]
}

func newNotificationRepository(requestSession *logy.RequestSession) iBaseConfigRepository[*notification.Notification] {
	return &notificationStruct{
		baseConfigRepositoryStruct: newBaseConfigRepository(requestSession, DP_NOTIFICATION, func() *notification.Notification {
			return &notification.Notification{}
		}),
	}
}
