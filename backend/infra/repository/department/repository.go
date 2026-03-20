// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package department

import (
	"regexp"
	"sync"

	"mercedes-benz.ghe.com/foss/disuko/domain/department"
	"mercedes-benz.ghe.com/foss/disuko/helper"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type departmentRepository struct {
	base.BaseRepositoryWithHardDelete[*department.Department]
	mem   map[string]*department.Department
	mutex sync.RWMutex
}

func NewDepartmentRepository(requestSession *logy.RequestSession) IDepartmentRepository {
	innerRepo := &departmentRepository{
		BaseRepositoryWithHardDelete: base.CreateRepositoryWithHardDelete(
			requestSession,
			DepartmentsCollectionName,
			func() *department.Department {
				return &department.Department{}
			},
			nil,
			"",
			nil,
			nil),
	}
	innerRepo.LoadFromDb(requestSession)
	return innerRepo
}

var wordPattern = regexp.MustCompile(`\S+`)

func (repo *departmentRepository) FindBySearchStr(requestSession *logy.RequestSession, searchStr string) []*department.Department {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()
	matchingDeparments := make(map[string]*department.Department)
	tokens := wordPattern.FindAllString(searchStr, -1)
	if len(tokens) == 0 {
		return nil
	}
OuterDep:
	for _, dep := range repo.mem {
		for _, token := range tokens {
			if !helper.StringContainsI(dep.DescriptionEnglish, token) &&
				!helper.StringContainsI(dep.OrgAbbreviation, token) &&
				!helper.StringContainsI(dep.Skz, token) &&
				!helper.StringContainsI(dep.CompanyCode, token) &&
				!helper.StringContainsI(dep.CompanyName, token) {
				continue OuterDep
			}
		}
		matchingDeparments[dep.Key] = dep
	}
	res := make([]*department.Department, 0, len(matchingDeparments))
	for _, dep := range matchingDeparments {
		res = append(res, dep)
	}

	return res
}

func (repo *departmentRepository) GetByDeptId(requestSession *logy.RequestSession, id string) *department.Department {
	repo.mutex.RLock()
	defer repo.mutex.RUnlock()
	dep, ok := repo.mem[id]
	if !ok {
		return nil
	}
	return dep
}

func (repo *departmentRepository) LoadFromDb(requestSession *logy.RequestSession) int {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	repo.mem = make(map[string]*department.Department)
	all := repo.FindAll(requestSession, false)
	for _, dep := range all {
		repo.mem[dep.Key] = dep
	}
	return len(all)
}

func (repo *departmentRepository) SaveDepartments(requestSession *logy.RequestSession, deps []*department.Department) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	s := repo.StartSession(base.DeleteSession, 2000)
	allEntities := repo.FindAll(requestSession, false)
	for _, key := range allEntities {
		s.AddEnt(key)
	}
	s.EndSession()
	repo.mem = nil

	s = repo.StartSession(base.UpdateSession, 2000)
	repo.mem = make(map[string]*department.Department)
	for _, dep := range deps {
		s.AddEnt(dep)
		repo.mem[dep.Key] = dep
	}
	s.EndSession()
}
