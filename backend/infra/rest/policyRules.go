// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"fmt"
	"net/http"
	"time"

	changeloglist2 "github.com/eclipse-disuko/disuko/domain/changeloglist"
	"github.com/eclipse-disuko/disuko/infra/repository/changeloglist"
	"github.com/eclipse-disuko/disuko/infra/repository/database"
	"github.com/eclipse-disuko/disuko/infra/repository/user"
	"github.com/eclipse-disuko/disuko/observermngmt"

	"github.com/eclipse-disuko/disuko/infra/repository/sbomlist"

	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/infra/repository/labels"
	"github.com/eclipse-disuko/disuko/infra/repository/policyrules"
	project3 "github.com/eclipse-disuko/disuko/infra/service/policy"

	"github.com/eclipse-disuko/disuko/domain/audit"
	"github.com/eclipse-disuko/disuko/domain/license"
	"github.com/eclipse-disuko/disuko/helper"
	auditHelper "github.com/eclipse-disuko/disuko/helper/audit"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/helper/roles"
	"github.com/eclipse-disuko/disuko/helper/validation"
	license2 "github.com/eclipse-disuko/disuko/infra/repository/license"
	project2 "github.com/eclipse-disuko/disuko/infra/repository/project"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/go-cmp/cmp"
)

type PolicyRulesHandler struct {
	PolicyRulesRepository   policyrules.IPolicyRulesRepository
	ProjectRepository       project2.IProjectRepository
	LicenseRepository       license2.ILicensesRepository
	LabelRepository         labels.ILabelRepository
	PolicyRulesService      project3.Service
	SbomListRepository      sbomlist.ISbomListRepository
	ChangeLogListRepository changeloglist.IChangeLogListRepository
	UserRepository          user.IUsersRepository
}

func (policyRulesHandler *PolicyRulesHandler) PolicyRulesGetByIdHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)

	policyRule := retrievePolicyRule(requestSession, policyRulesHandler, r)
	if policyRule == nil {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorDbNotFound))
	}

	render.JSON(w, r, policyRule.ToDto())
}

func (policyRulesHandler *PolicyRulesHandler) PolicyRulesChangeLogGetAllHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)

	id := chi.URLParam(r, "id")
	changeLogList := policyRulesHandler.ChangeLogListRepository.FindByKey(requestSession, id, false)
	if changeLogList != nil {
		render.JSON(w, r, changeloglist2.ToDtos(changeLogList.ChangeLogs))
	}
}

func (policyRulesHandler *PolicyRulesHandler) PolicyRulesTrailGetAllHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)

	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.AllowPolicy.Create && rights.AllowPolicy.Update && rights.AllowPolicy.Delete) {
		exception.ThrowExceptionSendDeniedResponse()
	}

	policyRule := retrievePolicyRule(requestSession, policyRulesHandler, r)

	auditTrail := make([]audit.AuditDto, 0)
	for _, item := range policyRule.GetAuditTrail() {
		auditTrail = append(auditTrail, item.ToDto())
	}
	render.JSON(w, r, auditTrail)
}

func (policyRulesHandler *PolicyRulesHandler) PolicyRulesGetHandler(w http.ResponseWriter, r *http.Request) {
	policyRulesHandler.HandlePolicyRulesGet(w, r)
}

func (policyRulesHandler *PolicyRulesHandler) HandlePolicyRulesGet(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	lists := policyRulesHandler.PolicyRulesRepository.FindAll(requestSession, false)

	render.JSON(w, r, license.ToPolicyRuleDtoList(lists))
}

func (policyRulesHandler *PolicyRulesHandler) PolicyRulesUpdateHandler(w http.ResponseWriter, r *http.Request) {
	policyRulesHandler.PolicyRulesUpdateOrCreateHandler(w, r, false)
}

func (policyRulesHandler *PolicyRulesHandler) PolicyRulesCopyHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	username, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if (!rights.AllowPolicy.Read) || (!rights.AllowPolicy.Create) {
		exception.ThrowExceptionSendDeniedResponse()
	}
	rule := retrievePolicyRule(requestSession, policyRulesHandler, r)
	if rule == nil {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorDbNotFound, policyrules.PolicyRulesCollectionName))
	}
	if rule.Deprecated {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.PolicyDeprecated))
	}
	copied := *rule
	copied.Name += " - Copy"
	copied.RootEntity = domain.NewRootEntity()
	copied.AuditTrail = nil
	copied.Active = false

	audit := copied.ToAudit(requestSession, policyRulesHandler.LabelRepository)
	auditOriginal := rule.ToAudit(requestSession, policyRulesHandler.LabelRepository)

	auditHelper.CreateAndAddAuditEntry(&copied.Container, username, message.PolicyRulesCopied, cmp.Diff, audit, auditOriginal)

	policyRulesHandler.PolicyRulesRepository.Save(requestSession, &copied)
	w.WriteHeader(200)
}

func (policyRulesHandler *PolicyRulesHandler) DeprecateHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	username, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.IsPolicyManager() {
		exception.ThrowExceptionSendDeniedResponse()
	}
	rule := retrievePolicyRule(requestSession, policyRulesHandler, r)
	if rule == nil {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorDbNotFound, policyrules.PolicyRulesCollectionName))
	}
	if rule.Deprecated {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.PolicyDeprecated))
	}
	beforeAudit := rule.ToAudit(requestSession, policyRulesHandler.LabelRepository)
	rule.Deprecated = true
	rule.DeprecatedDate = time.Now()
	rule.Active = false
	afterAudit := rule.ToAudit(requestSession, policyRulesHandler.LabelRepository)
	auditHelper.CreateAndAddAuditEntry(&rule.Container, username, message.PolicyRulesUpdated, cmp.Diff, afterAudit, beforeAudit)
	policyRulesHandler.PolicyRulesRepository.Update(requestSession, rule)
	render.JSON(w, r, rule.ToDto())
}

func (policyRulesHandler *PolicyRulesHandler) PolicyRulesCreateHandler(w http.ResponseWriter, r *http.Request) {
	policyRulesHandler.PolicyRulesUpdateOrCreateHandler(w, r, true)
}

func (policyRulesHandler *PolicyRulesHandler) PolicyRulesUpdateOrCreateHandler(w http.ResponseWriter, r *http.Request, isNew bool) {
	requestSession := logy.GetRequestSession(r)
	username, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if (!rights.AllowPolicy.Update && !isNew) || (!rights.AllowPolicy.Create && isNew) {
		exception.ThrowExceptionSendDeniedResponse()
	}
	var ruleRequestDto license.PolicyRuleRequestDto
	validation.DecodeAndValidate(r, &ruleRequestDto, false)
	if ruleRequestDto.Calculated {
		componentsAllow, componentsWarn, componentsDeny := policyRulesHandler.PolicyRulesService.CalculatePolicyRuleComponents(requestSession, ruleRequestDto.CalculatedConfig)
		ruleRequestDto.ComponentsAllow = componentsAllow
		ruleRequestDto.ComponentsWarn = componentsWarn
		ruleRequestDto.ComponentsDeny = componentsDeny
	} else {
		qc := database.New().SetMatcher(
			database.AttributeMatcher(
				"Deleted",
				database.EQ,
				false,
			),
		).SetKeep([]string{"licenseId"})
		allLics := policyRulesHandler.LicenseRepository.Query(requestSession, qc)

		known := make(map[string]bool)
		for _, k := range allLics {
			known[k.LicenseId] = true
		}

		ruleRequestDto.ComponentsAllow = policyRulesHandler.removeUnknownLicenses(requestSession, ruleRequestDto.ComponentsAllow, known)
		ruleRequestDto.ComponentsWarn = policyRulesHandler.removeUnknownLicenses(requestSession, ruleRequestDto.ComponentsWarn, known)
		ruleRequestDto.ComponentsDeny = policyRulesHandler.removeUnknownLicenses(requestSession, ruleRequestDto.ComponentsDeny, known)
	}

	policyRulesLoadedByName := policyRulesHandler.PolicyRulesRepository.FindByName(requestSession, ruleRequestDto.Name)
	var policyRuleEntity *license.PolicyRules
	var oldPolicyRuleEntityAudit *license.PolicyRulesAudit
	if isNew {
		if policyRulesLoadedByName != nil {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.CreatePolicyListAlreadyExists))
		}

		policyRuleEntity = &license.PolicyRules{
			RootEntity: domain.NewRootEntity(),
		}
	} else {
		policyRuleEntity = retrievePolicyRule(requestSession, policyRulesHandler, r)
		if policyRuleEntity == nil {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorDbNotFound, policyrules.PolicyRulesCollectionName))
		}
		if policyRuleEntity.Deprecated {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.PolicyDeprecated))
		}
		if policyRulesLoadedByName != nil && policyRulesLoadedByName.Name != "" && policyRulesLoadedByName.Key != policyRuleEntity.Key {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.UpdatePolicyListAlreadyExists))
		}
		oldPolicyRuleEntityAudit = policyRuleEntity.ToAudit(requestSession, policyRulesHandler.LabelRepository)
	}

	policyRuleEntity.Description = ruleRequestDto.Description
	policyRuleEntity.ComponentsAllow = ruleRequestDto.ComponentsAllow
	policyRuleEntity.ComponentsWarn = ruleRequestDto.ComponentsWarn
	policyRuleEntity.ComponentsDeny = ruleRequestDto.ComponentsDeny
	policyRuleEntity.LabelSets = cleanUpLabelSets(requestSession, policyRulesHandler.LabelRepository, ruleRequestDto.LabelSets)
	policyRuleEntity.Updated = time.Now()
	policyRuleEntity.Name = ruleRequestDto.Name
	policyRuleEntity.Auxiliary = ruleRequestDto.Auxiliary
	policyRuleEntity.Active = ruleRequestDto.Active
	policyRuleEntity.ApplyToAll = ruleRequestDto.ApplyToAll
	policyRuleEntity.Calculated = ruleRequestDto.Calculated
	policyRuleEntity.CalculatedConfig = ruleRequestDto.CalculatedConfig

	policyRulesAudit := policyRuleEntity.ToAudit(requestSession, policyRulesHandler.LabelRepository)
	if isNew {
		auditHelper.CreateAndAddAuditEntry(&policyRuleEntity.Container, username, message.PolicyRulesCreated, cmp.Diff, policyRulesAudit, license.PolicyRulesAudit{})
		policyRulesHandler.PolicyRulesRepository.Save(requestSession, policyRuleEntity)

		observermngmt.FireEvent(observermngmt.DatabaseEntryAddedOrDeleted, observermngmt.DatabaseSizeChange{
			RequestSession: requestSession,
			CollectionName: policyrules.PolicyRulesCollectionName,
			Rights:         rights,
			Username:       username,
		})
	} else {
		auditHelper.CreateAndAddAuditEntry(&policyRuleEntity.Container, username, message.PolicyRulesUpdated, cmp.Diff, policyRulesAudit, oldPolicyRuleEntityAudit)
		policyRulesHandler.PolicyRulesRepository.Update(requestSession, policyRuleEntity)
	}

	render.JSON(w, r, policyRuleEntity.ToDto())
}

func (policyRulesHandler *PolicyRulesHandler) removeUnknownLicenses(requestSession *logy.RequestSession, licenses []string, known map[string]bool) []string {
	var licensesForRemove []string
	for _, license := range licenses {
		if _, ok := known[license]; !ok {
			logy.Errorw(requestSession, "Unknown license found in policy rule: "+license+", remove from list!")
			licensesForRemove = append(licensesForRemove, license)
		}
	}

	for _, license := range licensesForRemove {
		licenses = helper.RemoveStrFromSlice(license, licenses)
	}

	return licenses
}

func cleanUpLabelSets(requestSession *logy.RequestSession, labelRepository labels.ILabelRepository, labelSets [][]string) [][]string {
	labelSetsCleaned := helper.RemoveDuplicates(labelSets)
	// validate labels
	for _, labelSet := range labelSetsCleaned {
		CheckIfLabelsExistOrThrowException(requestSession, labelRepository, labelSet)
	}
	return labelSetsCleaned
}

func (policyRulesHandler *PolicyRulesHandler) PolicyRulesDeleteHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowPolicy.Delete {
		exception.ThrowExceptionSendDeniedResponse()
	}
	id := chi.URLParam(r, "id")

	rule := retrievePolicyRule(requestSession, policyRulesHandler, r)
	if rule == nil {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorDbNotFound, policyrules.PolicyRulesCollectionName))
	}
	if rule.Deprecated {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.PolicyDeprecated))
	}

	policyRulesHandler.PolicyRulesRepository.Delete(requestSession, id)

	w.WriteHeader(200)

	fmt.Fprintf(w, "Successfully deleted allow/deny list with id %s \n", id)
}

func (handler *PolicyRulesHandler) CreateCSVHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.AllowPolicy.Create && rights.AllowPolicy.Update && rights.AllowPolicy.Delete) {
		exception.ThrowExceptionSendDeniedResponse()
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	project3.CreateCSV(&w, requestSession, handler.PolicyRulesRepository, handler.LicenseRepository)
}

func (handler *PolicyRulesHandler) CreateRuleCSVHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	w.Header().Set("Content-Type", "application/octet-stream")
	id := chi.URLParam(r, "id")
	project3.CreateRuleCSV(&w, requestSession, handler.PolicyRulesRepository, handler.LicenseRepository, id)
}

// PolicyRulesGetExternHandler godoc
//
//	@Summary	Get policy rules of project
//	@Id			getProjectPolicyRules
//	@Produce	json
//	@Param		uuid	path		string									true	"Project UUID e.g.: dummy-id---xxx-4413-yyy-24f060311111"
//	@Success	200		{object}	[]license.PolicyRulePublicResponseDto	"Policy Rules"
//	@Failure	404		{object}	exception.HttpError404					"NotFound Error"
//	@Failure	401		{object}	exception.HttpError						"Unauthorized Error"
//	@Router		/v1/projects/{uuid}/policyrules [get]
//	@security	Bearer
func (policyRulesHandler *PolicyRulesHandler) PolicyRulesGetExternHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	currentProject, _ := policyRulesHandler.retrieveProjectFromPublicRequest(requestSession, r, false)

	responseData := policyRulesHandler.PolicyRulesService.CollectPolicyRulesForProject(requestSession, currentProject, nil)

	render.JSON(w, r, responseData)
}

func retrievePolicyRule(requestSession *logy.RequestSession,
	policyRulesHandler *PolicyRulesHandler, r *http.Request,
) *license.PolicyRules {
	id := chi.URLParam(r, "id")
	return policyRulesHandler.PolicyRulesRepository.FindByKey(requestSession, id, false)
}
