// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
)

const ProjectCollectionName = "projects"

type IProjectRepository interface {
	base.IBaseRepositoryWithSoftDelete[*project.Project]
	FindByKeyWithDeleted(requestSession *logy.RequestSession, key string, optimized bool) *project.Project
	FindAllForUser(requestSession *logy.RequestSession, userId string) []*project.Project

	FindRecentByUpdatedForUser(requestSession *logy.RequestSession, userId string, limit int) []*project.Project
	ExistsBySchemaLabel(requestSession *logy.RequestSession, label string) bool
	ExistsByPolicyLabel(requestSession *logy.RequestSession, label string) bool
	CountForUser(requestSession *logy.RequestSession, userId string) int
	CountAllGroups(requestSession *logy.RequestSession) int
	CountGroupsForUser(requestSession *logy.RequestSession, userId string) int
}
