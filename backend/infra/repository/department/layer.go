// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package department

import (
	"github.com/eclipse-disuko/disuko/domain/department"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
)

const DepartmentsCollectionName = "departments"

type IDepartmentRepository interface {
	base.IBaseRepositoryWithHardDelete[*department.Department]
	FindBySearchStr(requestSession *logy.RequestSession, searchStr string) []*department.Department
	SaveDepartments(requestSession *logy.RequestSession, deps []*department.Department)
	LoadFromDb(requestSession *logy.RequestSession) int
	GetByDeptId(requestSession *logy.RequestSession, id string) *department.Department
}
