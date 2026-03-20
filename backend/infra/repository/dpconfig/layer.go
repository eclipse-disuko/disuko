// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package dpconfig

import (
	"github.com/eclipse-disuko/disuko/domain/dpconfig"
	"github.com/eclipse-disuko/disuko/domain/integrity"
	"github.com/eclipse-disuko/disuko/domain/notification"
	"github.com/eclipse-disuko/disuko/logy"
)

const DPConfigCollectionName = "dpconfig"

type DBConfigRepository struct {
	IntegrityCheck          iBaseConfigRepository[*integrity.DbIntegrityResult]
	SampleDataCreationState iBaseConfigRepository[*dpconfig.SampleDataCreationState]
	DpConfig                iBaseConfigRepository[*dpconfig.DPConfig]
	Notification            iBaseConfigRepository[*notification.Notification]
}

func NewDbConfigRepository(requestSession *logy.RequestSession) *DBConfigRepository {
	return &DBConfigRepository{
		IntegrityCheck:          newIntegrityCheckRepository(requestSession),
		SampleDataCreationState: newSampleDataCreationRepository(requestSession),
		DpConfig:                newDPConfigRepository(requestSession),
		Notification:            newNotificationRepository(requestSession),
	}
}
