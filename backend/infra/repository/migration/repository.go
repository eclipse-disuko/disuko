// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package migration

import (
	"github.com/eclipse-disuko/disuko/domain/migration"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
)

type labelsRepositoryStruct struct {
	base.BaseRepositoryWithHardDelete[*migration.Migration]
}

func NewMigrationRepository(requestSession *logy.RequestSession) IMigrationRepository {
	return &labelsRepositoryStruct{
		BaseRepositoryWithHardDelete: base.CreateRepositoryWithHardDelete[*migration.Migration](
			requestSession,
			MigrationCollectionName,
			func() *migration.Migration {
				return &migration.Migration{}
			},
			nil,
			"",
			nil,
			nil),
	}
}
