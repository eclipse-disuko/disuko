// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/eclipse-disuko/disuko/domain/audit"
	"github.com/eclipse-disuko/disuko/domain/license"
	policydecisions2 "github.com/eclipse-disuko/disuko/domain/policydecisions"
	"github.com/eclipse-disuko/disuko/domain/project"
	audit2 "github.com/eclipse-disuko/disuko/helper/audit"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/helper/roles"
	"github.com/eclipse-disuko/disuko/helper/validation"
	project2 "github.com/eclipse-disuko/disuko/infra/repository/project"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/go-cmp/cmp"
	"github.com/jinzhu/copier"
)

func (h *ProjectHandler) WizardGroupPreviewHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowProject.Create {
		exception.ThrowExceptionSendDeniedResponse()
	}
	req := extractWizardGroupBody(r)

	render.JSON(w, r, h.WizardService.PreviewGroup(requestSession, req))
}

func (h *ProjectHandler) WizardPreviewHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowProject.Create {
		exception.ThrowExceptionSendDeniedResponse()
	}
	req := extractWizardBody(r)

	render.JSON(w, r, h.WizardService.Preview(requestSession, req))
}

func (h *ProjectHandler) WizardCreateGroupHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	user, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowProject.Create {
		exception.ThrowExceptionSendDeniedResponse()
	}
	req := extractWizardGroupBody(r)

	render.JSON(w, r, h.WizardService.CreateGroup(requestSession, req, user))
}

func (h *ProjectHandler) WizardCreateHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	user, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowProject.Create {
		exception.ThrowExceptionSendDeniedResponse()
	}
	req := extractWizardBody(r)

	render.JSON(w, r, h.WizardService.Create(requestSession, req, user))
}

func (h *ProjectHandler) WizardUpdateHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)

	projectUUIDEscaped := chi.URLParam(r, "id")
	projectUUID, _ := url.QueryUnescape(projectUUIDEscaped)
	err := validation.CheckUuid(projectUUID)
	if err != nil {
		exception.ThrowExceptionBadRequestResponse()
	}

	currentProject := h.ProjectRepository.FindByKey(requestSession, projectUUID, false)
	if currentProject == nil {
		exception.ThrowExceptionClient404Message(message.GetI18N(message.ErrorDbRead, project2.ProjectCollectionName), "project not found: "+projectUUID)
	}
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	username, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowProject.Update {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.UpdateProject))
	}

	req := extractWizardBody(r)

	policyLabelsBefore := currentProject.PolicyLabels

	updatedProjectDto := h.WizardService.Update(requestSession, req, currentProject, username)

	policyLabelsAfter := updatedProjectDto.PolicyLabels

	if !equalStringSlices(policyLabelsBefore, policyLabelsAfter) {
		h.cancelDeniedDecisions(requestSession, projectUUID, username)
	}

	render.JSON(w, r, updatedProjectDto)
}

func (h *ProjectHandler) cancelDeniedDecisions(requestSession *logy.RequestSession, projectUUID string, username string) {
	policyDecisions := h.PolicyDecisionsRepository.FindByKey(requestSession, projectUUID, false)
	if policyDecisions != nil {
		auditEntries := make([]*audit.Audit, 0)
		var changed = false
		for _, pd := range policyDecisions.Decisions {
			if pd.PolicyEvaluated == string(license.DENY) && pd.Active {
				oldPd := policydecisions2.PolicyDecision{}
				copier.Copy(&oldPd, pd)

				pd.Updated = time.Now()
				pd.Active = false

				auditEntries = append(auditEntries, audit2.CreateAuditEntry(username, message.PolicyDecisionUpdated, cmp.Diff, pd, &oldPd))
				changed = true
			}
		}
		if changed {
			h.AuditLogListRepository.CreateAuditEntriesByKey(requestSession, projectUUID, auditEntries)
			h.PolicyDecisionsRepository.Update(requestSession, policyDecisions)
		}
	}
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	slices.Sort(a)
	slices.Sort(b)

	return slices.Equal(a, b)
}

func (h *ProjectHandler) WizardGetHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)

	projectUUIDEscaped := chi.URLParam(r, "id")
	projectUUID, _ := url.QueryUnescape(projectUUIDEscaped)
	err := validation.CheckUuid(projectUUID)
	if err != nil {
		exception.ThrowExceptionBadRequestResponse()
	}

	currentProject := h.ProjectRepository.FindByKey(requestSession, projectUUID, false)
	if currentProject == nil {
		exception.ThrowExceptionClient404Message(message.GetI18N(message.ErrorDbRead, project2.ProjectCollectionName), "project not found: "+projectUUID)
	}
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	_, rights := roles.GetAndCheckProjectRights(requestSession, r, currentProject, false)
	if !rights.AllowProject.Update {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.UpdateProject))
	}

	wizardAttributesDto := h.WizardService.Get(requestSession, currentProject)

	policyDecisions := h.PolicyDecisionsRepository.FindByKey(requestSession, projectUUID, false)
	wizardAttributesDto.HasDeniedDecisions = hasActiveDeniedDecision(policyDecisions)

	render.JSON(w, r, wizardAttributesDto)
}

func extractWizardBody(r *http.Request) *project.WizardProjectDto {
	var res project.WizardProjectDto
	validation.DecodeAndValidate(r, &res, false)
	return &res
}

func extractWizardGroupBody(r *http.Request) *project.WizardGroupDto {
	var res project.WizardGroupDto
	validation.DecodeAndValidate(r, &res, false)
	return &res
}
