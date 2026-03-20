// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"slices"

	"github.com/eclipse-disuko/disuko/connector/application"
	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/audit"
	"github.com/eclipse-disuko/disuko/domain/department"
	"github.com/eclipse-disuko/disuko/domain/label"
	"github.com/eclipse-disuko/disuko/domain/project"
	auditHelper "github.com/eclipse-disuko/disuko/helper/audit"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	auditloglistRepo "github.com/eclipse-disuko/disuko/infra/repository/auditloglist"
	departmentRepo "github.com/eclipse-disuko/disuko/infra/repository/department"
	labelsRepo "github.com/eclipse-disuko/disuko/infra/repository/labels"
	projectRepo "github.com/eclipse-disuko/disuko/infra/repository/project"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/google/go-cmp/cmp"
	"github.com/jinzhu/copier"
)

type WizardService struct {
	LabelRepository        labelsRepo.ILabelRepository
	ProjectRepository      projectRepo.IProjectRepository
	DepartmentRepository   departmentRepo.IDepartmentRepository
	AuditLogListRepository auditloglistRepo.IAuditLogListRepository
	ApplicationConnector   *application.Connector
}

var platformLabelMap map[project.TargetPlatform]string = map[project.TargetPlatform]string{
	project.PlatformEnterprise: label.ENTERPRISE_PLATFORM,
	project.PlatformMobile:     label.MOBILE_PLATFORM,
	project.PlatformVehicle:    label.VEHICLE_PLATFORM,
	project.PlatformOther:      label.OTHER_PLATFORM,
}

var architectureLabelMap map[project.Architecture]string = map[project.Architecture]string{
	project.ArchitectureBackend:         label.BACKEND_LAYER,
	project.ArchitectureFrontend:        label.FRONTEND_LAYER,
	project.ArchitectureVehicleOffboard: label.OFFBOARD,
	project.ArchitectureVehicleOnboard:  label.ONBOARD,
}

var usersLabelMap map[project.TargetUsers]string = map[project.TargetUsers]string{
	project.TargetUsersCustomer:        label.CUSTOMER_USERS,
	project.TargetUsersCompany:         label.COMPANY_USERS,
	project.TargetUsersBusinessPartner: label.BP_USERS,
}

var distributionLabelMap map[project.DistributionTarget]string = map[project.DistributionTarget]string{
	project.DistributionTargetsBusinessPartner: label.BP_TARGET,
	project.DistributionTargetsCompany:         label.COMPANY_TARGET,
}

var developmentLabelMap map[project.Development]string = map[project.Development]string{
	project.DevelopmentsExternal: label.EXTERNAL_DEVELOP,
	project.DevelopmentsInhouse:  label.INHOUSE_DEVELOP,
	project.DevelopmentsInternal: label.INTERNAL_DEVELOP,
}

type combination struct {
	platform []project.TargetPlatform
	arch     []project.Architecture
	users    []project.TargetUsers
	dist     []project.DistributionTarget
	dev      []project.Development
}

var allowedCombinations = []combination{
	{
		[]project.TargetPlatform{project.PlatformEnterprise},
		[]project.Architecture{project.ArchitectureFrontend, project.ArchitectureBackend},
		nil,
		nil,
		nil,
	},
	{
		[]project.TargetPlatform{project.PlatformMobile},
		[]project.Architecture{project.ArchitectureFrontend, project.ArchitectureBackend},
		nil,
		nil,
		nil,
	},
	{
		[]project.TargetPlatform{project.PlatformVehicle},
		[]project.Architecture{project.ArchitectureVehicleOffboard, project.ArchitectureVehicleOnboard},
		[]project.TargetUsers{project.TargetUsersEmpty},
		[]project.DistributionTarget{project.DistributionTargetsEmpty},
		nil,
	},
	{
		[]project.TargetPlatform{project.PlatformOther},
		[]project.Architecture{project.ArchitectureEmpty},
		[]project.TargetUsers{project.TargetUsersEmpty},
		[]project.DistributionTarget{project.DistributionTargetsEmpty},
		nil,
	},
}

func (s *WizardService) CreateGroup(rs *logy.RequestSession, req *project.WizardGroupDto, requester string) project.ProjectDto {
	if !s.validateCombinationOfGroup(req) {
		exception.ThrowExceptionBadRequestResponse()
	}
	var (
		parent     *project.Project
		parentKey  string
		parentName string
	)
	if req.ParentKey != "" {
		parent = s.ProjectRepository.FindByKey(rs, req.ParentKey, false)
		if parent == nil {
			exception.ThrowExceptionBadRequestResponse()
		} else {
			parentKey = parent.Key
			parentName = parent.Name
		}
	}

	sl := s.LabelRepository.FindByNameAndType(rs, label.COMMON_STANDARD, label.SCHEMA)
	if sl == nil {
		exception.ThrowExceptionBadRequestResponse()
	}

	var (
		docDept  *department.Department
		custDept *department.Department
	)

	newPr := project.Project{
		RootEntity:    domain.NewRootEntity(),
		Name:          req.Name,
		Description:   req.Description,
		Versions:      map[string]*project.ProjectVersion{},
		SchemaLabel:   sl.Key,
		PolicyLabels:  s.policyLabelsForGroups(rs, req),
		ProjectLabels: s.projectLabelsForGroups(rs, req),
		FreeLabels:    []string{},
		Children:      nil,
		UserManagement: project.UserManagementEntity{
			Users: []*project.ProjectMemberEntity{
				{
					ChildEntity:   domain.NewChildEntity(),
					UserId:        requester,
					UserType:      project.OWNER,
					Comment:       "",
					IsResponsible: true,
				},
			},
		},
		Status:          project.Ready,
		IsGroup:         true,
		Parent:          parentKey,
		ParentName:      parentName,
		ApplicationMeta: project.ApplicationMeta{},

		CustomerMeta:      s.customerMeta(rs, &req.Settings, &custDept),
		DocumentMeta:      s.documentMeta(rs, req.Development, &req.Settings, &docDept),
		NoticeContactMeta: s.noticeContactMeta(&req.Settings),
		SupplierExtraData: s.supplierExtraData(req.Development),

		IsNoFoss: req.Settings.NoFossProject,
	}
	s.ProjectRepository.Save(rs, &newPr)

	projectAuditEntries := make([]*audit.Audit, 0)
	projectAuditEntries = append(projectAuditEntries, auditHelper.CreateAuditEntry(requester, message.GroupCreated, cmp.Diff, newPr, project.Project{}))
	s.AuditLogListRepository.CreateAuditEntriesByKey(rs, newPr.Key, projectAuditEntries)

	return newPr.ToDto(docDept, false, custDept, false, req.GetIsDummy())
}

func (s *WizardService) PreviewGroup(rs *logy.RequestSession, req *project.WizardGroupDto) project.WizardGroupDto {
	if !s.validateCombinationOfGroup(req) {
		exception.ThrowExceptionBadRequestResponse()
	}
	if req.ParentKey != "" {
		p := s.ProjectRepository.FindByKey(rs, req.ParentKey, true)
		if p == nil {
			exception.ThrowExceptionBadRequestResponse()
		}
	}

	labels := append(s.policyLabelsForGroups(rs, req), s.projectLabelsForGroups(rs, req)...)
	s.appendLabelKey(rs, label.COMMON_STANDARD, label.SCHEMA, &labels)

	res := project.WizardGroupDto{
		Name:               req.Name,
		Description:        req.Description,
		Architecture:       req.Architecture,
		DistributionTarget: req.DistributionTarget,
		Development:        req.Development,
		IsDummy:            req.IsDummy,
		IsGroup:            req.IsGroup,
		ParentKey:          req.ParentKey,
		Settings:           req.Settings,
		Labels:             labels,
	}

	if req.Development == project.DevelopmentsInhouse {
		res.Settings.DocumentMeta.SupplierDept = res.Settings.CustomerMeta.Dept
	}

	return res
}

func (s *WizardService) Preview(rs *logy.RequestSession, req *project.WizardProjectDto) project.WizardProjectDto {
	if !s.validateCombination(req) {
		exception.ThrowExceptionBadRequestResponse()
	}
	if req.ParentKey != "" {
		p := s.ProjectRepository.FindByKey(rs, req.ParentKey, true)
		if p == nil {
			exception.ThrowExceptionBadRequestResponse()
		}
	}

	app := s.appMeta(rs, &req.ApplicationMeta)

	labels := append(s.policyLabels(rs, req), s.projectLabels(rs, req)...)
	s.appendLabelKey(rs, label.COMMON_STANDARD, label.SCHEMA, &labels)

	res := project.WizardProjectDto{
		Name:               req.Name,
		Description:        req.Description,
		ApplicationMeta:    app.ToDto(),
		TargetPlatform:     req.TargetPlatform,
		Architecture:       req.Architecture,
		DistributionTarget: req.DistributionTarget,
		TargetUsers:        req.TargetUsers,
		Development:        req.Development,
		IsDummy:            req.IsDummy,
		IsGroup:            req.IsGroup,
		ParentKey:          req.ParentKey,
		Settings:           req.Settings,
		Labels:             labels,
	}

	if req.Development == project.DevelopmentsInhouse {
		res.Settings.DocumentMeta.SupplierDept = res.Settings.CustomerMeta.Dept
	}

	return res
}

func (s *WizardService) Create(rs *logy.RequestSession, req *project.WizardProjectDto, requester string) project.ProjectDto {
	if !s.validateCombination(req) {
		exception.ThrowExceptionBadRequestResponse()
	}
	var (
		parent     *project.Project
		parentKey  string
		parentName string
	)
	if req.ParentKey != "" {
		if req.GetIsGroup() {
			exception.ThrowExceptionBadRequestResponse()
		}
		parent = s.ProjectRepository.FindByKey(rs, req.ParentKey, false)
		if parent == nil {
			exception.ThrowExceptionBadRequestResponse()
		}
		parentKey = parent.Key
		parentName = parent.Name
	}

	sl := s.LabelRepository.FindByNameAndType(rs, label.COMMON_STANDARD, label.SCHEMA)
	if sl == nil {
		exception.ThrowExceptionBadRequestResponse()
	}

	app := s.appMeta(rs, &req.ApplicationMeta)

	var (
		docDept  *department.Department
		custDept *department.Department
		main     = project.ProjectVersion{
			ChildEntity: domain.NewChildEntity(),
			Name:        message.DefaultBranchMainName,
			Description: message.DefaultBranchMainDescription,
			Status:      "new",
		}
		dev = project.ProjectVersion{
			ChildEntity: domain.NewChildEntity(),
			Name:        message.DefaultBranchDevName,
			Description: message.DefaultBranchDevDescription,
			Status:      "new",
		}
	)

	newPr := project.Project{
		RootEntity:  domain.NewRootEntity(),
		Name:        req.Name,
		Description: req.Description,
		Versions: map[string]*project.ProjectVersion{
			main.Key: &main,
			dev.Key:  &dev,
		},
		SchemaLabel:   sl.Key,
		PolicyLabels:  s.policyLabels(rs, req),
		ProjectLabels: s.projectLabels(rs, req),
		FreeLabels:    []string{},
		UserManagement: project.UserManagementEntity{
			Users: []*project.ProjectMemberEntity{
				{
					ChildEntity:   domain.NewChildEntity(),
					UserId:        requester,
					UserType:      project.OWNER,
					Comment:       "",
					IsResponsible: true,
				},
			},
		},
		Status:          project.Ready,
		IsGroup:         req.GetIsGroup(),
		Parent:          parentKey,
		ParentName:      parentName,
		ApplicationMeta: app,

		CustomerMeta:      s.customerMeta(rs, &req.Settings, &custDept),
		DocumentMeta:      s.documentMeta(rs, req.Development, &req.Settings, &docDept),
		NoticeContactMeta: s.noticeContactMeta(&req.Settings),
		SupplierExtraData: s.supplierExtraData(req.Development),

		IsNoFoss: req.Settings.NoFossProject,
	}
	s.ProjectRepository.Save(rs, &newPr)

	if parent != nil {
		parent.Children = append(parent.Children, newPr.Key)
		s.ProjectRepository.Update(rs, parent)
	}

	auditMsg := message.ProjectCreated
	if req.GetIsGroup() {
		auditMsg = message.GroupCreated
	}
	projectAuditEntries := make([]*audit.Audit, 0)
	projectAuditEntries = append(projectAuditEntries, auditHelper.CreateAuditEntry(requester, auditMsg, cmp.Diff, newPr, project.Project{}))
	s.AuditLogListRepository.CreateAuditEntriesByKey(rs, newPr.Key, projectAuditEntries)

	return newPr.ToDto(docDept, false, custDept, false, req.GetIsDummy())
}

func (s *WizardService) Update(rs *logy.RequestSession, req *project.WizardProjectDto, existing *project.Project, requester string) project.ProjectDto {
	if !s.validateCombination(req) {
		exception.ThrowExceptionBadRequestResponse()
	}
	before := project.Project{}
	copier.Copy(&before, existing)

	existing.Name = req.Name
	existing.Description = req.Description

	existing.PolicyLabels = s.policyLabels(rs, req)
	existing.ProjectLabels = s.projectLabels(rs, req)

	existing.ApplicationMeta = s.appMeta(rs, &req.ApplicationMeta)

	existing.IsNoFoss = req.Settings.NoFossProject

	var (
		docDept  *department.Department
		custDept *department.Department
	)
	existing.CustomerMeta = s.updatedCustomerMeta(rs, req, existing.CustomerMeta, &custDept)
	existing.DocumentMeta = s.documentMeta(rs, req.Development, &req.Settings, &docDept)
	existing.NoticeContactMeta = s.noticeContactMeta(&req.Settings)
	existing.SupplierExtraData = s.updatedSupplierExtraData(req, existing.SupplierExtraData)

	s.AuditLogListRepository.CreateAuditEntryByKey(rs, existing.Key, requester, message.ProjectUpdated, cmp.Diff, existing, &before)
	s.ProjectRepository.Update(rs, existing)

	return existing.ToDto(docDept, false, custDept, false, req.GetIsDummy())
}

func (s *WizardService) Get(rs *logy.RequestSession, pr *project.Project) project.WizardAttributesDto {
	policyLables := s.LabelRepository.FindAllByType(rs, label.POLICY)
	projectLabels := s.LabelRepository.FindAllByType(rs, label.PROJECT)
	var res project.WizardAttributesDto
OUTER_PL:
	for _, key := range pr.PolicyLabels {
		label := findLabel(key, policyLables)
		if label == nil {
			logy.Warnf(rs, "label with key %s not found", key)
			continue
		}
		for p, lt := range platformLabelMap {
			if lt == label.Name {
				res.TargetPlatform = p
				continue OUTER_PL
			}
		}
		for a, lt := range architectureLabelMap {
			if lt == label.Name {
				res.Architecture = a
				continue OUTER_PL
			}
		}
		for t, lt := range usersLabelMap {
			if lt == label.Name {
				res.TargetUsers = t
				continue OUTER_PL
			}
		}
		for d, lt := range distributionLabelMap {
			if lt == label.Name {
				res.DistributionTarget = d
				continue OUTER_PL
			}
		}
	}
OUTER_PR:
	for _, key := range pr.ProjectLabels {
		label := findLabel(key, projectLabels)
		if label == nil {
			logy.Warnf(rs, "label with key %s not found", key)
			continue
		}
		for d, lt := range developmentLabelMap {
			if lt == label.Name {
				res.Development = d
				break OUTER_PR
			}
		}
	}
	return res
}

func findLabel(key string, labels []*label.Label) *label.Label {
	for _, l := range labels {
		if l.Key == key {
			return l
		}
	}
	return nil
}

func (s *WizardService) validateCombinationOfGroup(req *project.WizardGroupDto) bool {
	for _, allowed := range allowedCombinations {
		if allowed.arch != nil && !slices.Contains(allowed.arch, req.Architecture) {
			continue
		}
		if allowed.dist != nil && !slices.Contains(allowed.dist, req.DistributionTarget) {
			continue
		}
		if allowed.dev != nil && !slices.Contains(allowed.dev, req.Development) {
			continue
		}
		return true
	}
	return false
}

func (s *WizardService) validateCombination(req *project.WizardProjectDto) bool {
	for _, allowed := range allowedCombinations {
		if allowed.platform != nil && !slices.Contains(allowed.platform, req.TargetPlatform) {
			continue
		}
		if allowed.arch != nil && !slices.Contains(allowed.arch, req.Architecture) {
			continue
		}
		if allowed.users != nil && !slices.Contains(allowed.users, req.TargetUsers) {
			continue
		}
		if allowed.dist != nil && !slices.Contains(allowed.dist, req.DistributionTarget) {
			continue
		}
		if allowed.dev != nil && !slices.Contains(allowed.dev, req.Development) {
			continue
		}
		return true
	}
	return false
}

func (s *WizardService) policyLabels(rs *logy.RequestSession, req *project.WizardProjectDto) []string {
	var res []string
	res = make([]string, 0)
	s.appendLabelKey(rs, platformLabelMap[req.TargetPlatform], label.POLICY, &res)
	s.appendLabelKey(rs, architectureLabelMap[req.Architecture], label.POLICY, &res)
	s.appendLabelKey(rs, distributionLabelMap[req.DistributionTarget], label.POLICY, &res)
	s.appendLabelKey(rs, usersLabelMap[req.TargetUsers], label.POLICY, &res)
	return res
}

func (s *WizardService) projectLabels(rs *logy.RequestSession, req *project.WizardProjectDto) []string {
	var res []string
	res = make([]string, 0)
	s.appendLabelKey(rs, developmentLabelMap[req.Development], label.PROJECT, &res)
	if req.GetIsDummy() {
		s.appendLabelKey(rs, label.DUMMY, label.PROJECT, &res)
	}
	return res
}

func (s *WizardService) policyLabelsForGroups(rs *logy.RequestSession, req *project.WizardGroupDto) []string {
	var res []string
	res = make([]string, 0)
	s.appendLabelKey(rs, architectureLabelMap[req.Architecture], label.POLICY, &res)
	s.appendLabelKey(rs, distributionLabelMap[req.DistributionTarget], label.POLICY, &res)
	return res
}

func (s *WizardService) projectLabelsForGroups(rs *logy.RequestSession, req *project.WizardGroupDto) []string {
	var res []string
	res = make([]string, 0)
	s.appendLabelKey(rs, developmentLabelMap[req.Development], label.PROJECT, &res)
	if req.GetIsDummy() {
		// does not allow dummy groups for now
		// s.appendLabelKey(rs, label.DUMMY, label.PROJECT, &res)
	}
	return res
}

func (s *WizardService) appendLabelKey(rs *logy.RequestSession, name string, labelType label.LabelType, labels *[]string) {
	if name == "" {
		return
	}
	l := s.LabelRepository.FindByNameAndType(rs, name, labelType)
	if l == nil {
		exception.ThrowExceptionBadRequestResponse()
	}
	*labels = append(*labels, l.Key)
}

func (s *WizardService) appMeta(rs *logy.RequestSession, applicationMeta *project.ApplicationMetaDto) project.ApplicationMeta {
	var res project.ApplicationMeta
	// TODO: on mb-disuko app is needed
	if applicationMeta.Id != "" && s.ApplicationConnector != nil {
		a := s.ApplicationConnector.GetApplication(rs, applicationMeta.Id)
		res = project.ApplicationMeta{
			Id:           a.Id,
			SecondaryId:  a.SecondaryId,
			Name:         a.Name,
			ExternalLink: a.Link,
		}
	}
	return res
}

func (s *WizardService) customerMeta(rs *logy.RequestSession, req *project.ProjectSettingsDto, dept **department.Department) project.CustomerMeta {
	if req.CustomerMeta.Dept == nil {
		exception.ThrowExceptionBadRequestResponse()
	}
	if d := s.DepartmentRepository.FindByKey(rs, req.CustomerMeta.Dept.DeptId, false); d != nil {
		*dept = d
	} else {
		exception.ThrowExceptionBadRequestResponse()
	}
	res := project.CustomerMeta{
		ChildEntity: domain.NewChildEntity(),
		Address:     req.CustomerMeta.Address,
		DeptId:      req.CustomerMeta.Dept.DeptId,
	}
	return res
}

func (s *WizardService) updatedCustomerMeta(rs *logy.RequestSession, req *project.WizardProjectDto, existing project.CustomerMeta, dept **department.Department) project.CustomerMeta {
	m := s.customerMeta(rs, &req.Settings, dept)
	m.FRI = existing.FRI
	m.SRI = existing.SRI
	return m
}

func (s *WizardService) documentMeta(rs *logy.RequestSession, development project.Development, settings *project.ProjectSettingsDto, dept **department.Department) project.DisclosureDocumentMeta {
	switch development {
	case project.DevelopmentsExternal:
		return project.DisclosureDocumentMeta{
			ChildEntity:     domain.NewChildEntity(),
			SupplierName:    settings.DocumentMeta.SupplierName,
			SupplierAddress: settings.DocumentMeta.SupplierAddress,
			SupplierNr:      settings.DocumentMeta.SupplierNr,
		}
	case project.DevelopmentsInhouse, project.DevelopmentsInternal:
		if settings.DocumentMeta.SupplierDept == nil {
			exception.ThrowExceptionBadRequestResponse()
		}
		if d := s.DepartmentRepository.FindByKey(rs, settings.DocumentMeta.SupplierDept.DeptId, false); d != nil {
			*dept = d
		} else {
			exception.ThrowExceptionBadRequestResponse()
		}
		return project.DisclosureDocumentMeta{
			ChildEntity:     domain.NewChildEntity(),
			SupplierAddress: settings.DocumentMeta.SupplierAddress,
			SupplierDeptId:  settings.DocumentMeta.SupplierDept.DeptId,
		}
	}
	return project.DisclosureDocumentMeta{}
}

func (s *WizardService) noticeContactMeta(settings *project.ProjectSettingsDto) project.NoticeContactMeta {
	return project.NoticeContactMeta{
		ChildEntity: domain.NewChildEntity(),
		Address:     settings.NoticeContactMeta.Address,
	}
}

func (s *WizardService) supplierExtraData(development project.Development) project.SupplierExtraData {
	return project.SupplierExtraData{
		External: development == project.DevelopmentsExternal,
	}
}

func (s *WizardService) updatedSupplierExtraData(req *project.WizardProjectDto, existing project.SupplierExtraData) project.SupplierExtraData {
	d := s.supplierExtraData(req.Development)
	d.FRI = existing.FRI
	d.SRI = existing.SRI
	return d
}
