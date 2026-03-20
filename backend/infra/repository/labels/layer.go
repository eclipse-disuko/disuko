// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package labels

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/label"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

const LabelCollectionName = "labels"

type ILabelRepository interface {
	base.IBaseRepositoryWithHardDelete[*label.Label]
	FindByNameAndType(requestSession *logy.RequestSession, name string, labelType label.LabelType) *label.Label
	FindAllByType(requestSession *logy.RequestSession, labelType label.LabelType) []*label.Label
	LoadFromDb(requestSession *logy.RequestSession) int
}
