// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package ProjectLabelsService

import (
	"slices"

	label2 "mercedes-benz.ghe.com/foss/disuko/domain/label"
	"mercedes-benz.ghe.com/foss/disuko/domain/project"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/labels"
	projectRepo "mercedes-benz.ghe.com/foss/disuko/infra/repository/project"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type ProjectLabelService struct {
	ProjectRepo projectRepo.IProjectRepository
	LabelRepo   labels.ILabelRepository
}

func (s *ProjectLabelService) CheckProjectLabel(requestSession *logy.RequestSession, pr *project.Project, labelName string, labelType label2.LabelType) bool {
	platformLabel := s.LabelRepo.FindByNameAndType(requestSession, labelName, labelType)
	if platformLabel == nil {
		return false
	}

	labelKey := platformLabel.GetKey()
	switch labelType {
	case label2.POLICY:
		return slices.Contains(pr.PolicyLabels, labelKey)
	case label2.PROJECT:
		return slices.Contains(pr.FreeLabels, labelKey)
	case label2.SCHEMA:
		return pr.SchemaLabel == labelKey
	default:
		return false
	}
}

func (s *ProjectLabelService) HasLabelInGroupOrProject(requestSession *logy.RequestSession, pr *project.Project, labelName string, labelType label2.LabelType) bool {
	if !pr.IsGroup {
		return s.CheckProjectLabel(requestSession, pr, labelName, labelType)
	}
	for _, cId := range pr.Children {
		c := s.ProjectRepo.FindByKey(requestSession, cId, true)
		if c == nil {
			continue
		}
		if s.CheckProjectLabel(requestSession, c, labelName, labelType) {
			return true
		}
	}
	return false
}

func (s *ProjectLabelService) AllChildrenHaveLabel(requestSession *logy.RequestSession, pr *project.Project, labelName string, labelType label2.LabelType) bool {
	if !pr.IsGroup || len(pr.Children) == 0 {
		return false
	}
	for _, cId := range pr.Children {
		c := s.ProjectRepo.FindByKeyWithDeleted(requestSession, cId, true)
		if c.Deleted {
			continue
		}
		if !s.CheckProjectLabel(requestSession, c, labelName, labelType) {
			return false
		}
	}
	return true
}

func (s *ProjectLabelService) HasVehiclePlatformLabel(requestSession *logy.RequestSession, pr *project.Project) bool {
	return s.HasLabelInGroupOrProject(requestSession, pr, label2.VEHICLE_PLATFORM, label2.POLICY)
}

func (s *ProjectLabelService) OnlyVehicleChildren(requestSession *logy.RequestSession, pr *project.Project) bool {
	return s.AllChildrenHaveLabel(requestSession, pr, label2.VEHICLE_PLATFORM, label2.POLICY)
}
