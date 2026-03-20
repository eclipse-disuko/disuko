// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package newsbox

import (
	"github.com/eclipse-disuko/disuko/domain/newsbox"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
)

const collName = "newsbox"

type IRepo interface {
	base.IBaseRepositoryWithHardDelete[*newsbox.Item]
}
