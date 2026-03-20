// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package department

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/department"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

const DepartmentsCollectionName = "departments"

type IDepartmentRepository interface {
	base.IBaseRepositoryWithHardDelete[*department.Department]
	FindBySearchStr(requestSession *logy.RequestSession, searchStr string) []*department.Department
	SaveDepartments(requestSession *logy.RequestSession, deps []*department.Department)
	LoadFromDb(requestSession *logy.RequestSession) int
	GetByDeptId(requestSession *logy.RequestSession, id string) *department.Department
}
