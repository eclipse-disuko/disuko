// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package migration

import (
	"github.com/eclipse-disuko/disuko/domain/migration"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
)

const MigrationCollectionName = "migrations"

type IMigrationRepository interface {
	base.IBaseRepositoryWithHardDelete[*migration.Migration]
}
