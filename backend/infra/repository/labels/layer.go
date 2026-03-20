// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package labels

import (
	"github.com/eclipse-disuko/disuko/domain/label"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
)

const LabelCollectionName = "labels"

type ILabelRepository interface {
	base.IBaseRepositoryWithHardDelete[*label.Label]
	FindByNameAndType(requestSession *logy.RequestSession, name string, labelType label.LabelType) *label.Label
	FindAllByType(requestSession *logy.RequestSession, labelType label.LabelType) []*label.Label
	LoadFromDb(requestSession *logy.RequestSession) int
}
