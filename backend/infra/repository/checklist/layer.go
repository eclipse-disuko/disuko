// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package checklist

import (
	"github.com/eclipse-disuko/disuko/domain/checklist"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
)

const collName = "checklists"

type IChecklistRepository interface {
	base.IBaseRepositoryWithHardDelete[*checklist.Checklist]
}
