// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"io"
	"math"
	"net/http"
	"net/url"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	sorthelper "github.com/eclipse-disuko/disuko/helper/sort"
	"go.uber.org/zap/zapcore"

	"github.com/eclipse-disuko/disuko/domain/internalToken"
	obligation2 "github.com/eclipse-disuko/disuko/domain/obligation"

	"github.com/eclipse-disuko/disuko/domain/search"
	"github.com/eclipse-disuko/disuko/helper/filter"
	"github.com/eclipse-disuko/disuko/infra/repository/database"
	"github.com/eclipse-disuko/disuko/observermngmt"

	"github.com/google/go-cmp/cmp"
	levenshtein "github.com/ka-weihe/fast-levenshtein"
	"golang.org/x/text/unicode/norm"

	"github.com/eclipse-disuko/disuko/domain/audit"
	auditHelper "github.com/eclipse-disuko/disuko/helper/audit"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/helper/validation"
	"github.com/eclipse-disuko/disuko/infra/repository/jobs"
	"github.com/eclipse-disuko/disuko/infra/repository/obligation"
	"github.com/eclipse-disuko/disuko/infra/repository/policyrules"
	"github.com/eclipse-disuko/disuko/infra/repository/spdx_license"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"

	"github.com/eclipse-disuko/disuko/domain/license"
	"github.com/eclipse-disuko/disuko/helper"
	"github.com/eclipse-disuko/disuko/helper/roles"
	license2 "github.com/eclipse-disuko/disuko/infra/repository/license"
	"github.com/eclipse-disuko/disuko/logy"
)

type LicensesHandler struct {
	LicenseRepository     license2.ILicensesRepository
	JobRepository         jobs.IJobsRepository
	ObligationRepository  obligation.IObligationRepository
	PolicyRulesRepository policyrules.IPolicyRulesRepository
	SpdxLicenseRepository spdx_license.ISpdxLicensesRepository
}

func (licensesHandler *LicensesHandler) SearchHandler(w http.ResponseWriter, r *http.Request) {
	requestVersion := r.Header.Get("X-Client-Version")

	if requestVersion == "2.0" {
		// Or we check the request by
		// var body map[string]interface{}
		// json.NewDecoder(r.Body).Decode(&body)
		//  if _, exists := body["newVue3SpecificField"]; exists
		var searchOptionsVue3 search.RequestSearchOptionsNew
		validation.DecodeAndValidate(r, &searchOptionsVue3, false)
		licensesHandler.searchHandler(w, r, &searchOptionsVue3)
	} else {
		var searchOptions search.RequestSearchOptions
		validation.DecodeAndValidate(r, &searchOptions, false)
		licensesHandler.searchHandler(w, r, &searchOptions)
	}
}

func (licensesHandler *LicensesHandler) searchHandler(w http.ResponseWriter, r *http.Request, searchOptions search.SortableOptions) {
	requestSession := logy.GetRequestSession(r)

	qc := database.New().SetMatcher(
		database.AttributeMatcher(
			"Deleted",
			database.EQ,
			false,
		),
	).SetKeep([]string{
		"Updated",
		"Created",
		"meta",
		"licenseId",
		"name",
		"source",
		"Aliases",
	})

	licenses := licensesHandler.LicenseRepository.Query(requestSession, qc)

	dtos := make([]*license.LicenseSlimDto, 0)

	possibleCharts := make(map[string]int)
	possibleSources := make(map[string]int)
	possibleFamilies := make(map[string]int)
	possibleApproval := map[string]int{
		license.NotSet:     0,
		license.Pending:    0,
		license.Check:      0,
		license.Assigning:  0,
		license.Approved:   0,
		license.Forbidden:  0,
		license.Deprecated: 0,
	}
	possibleType := make(map[string]int)

	extractors := map[string]func(*license.LicenseSlimDto) string{
		"isLicenseChart": func(item *license.LicenseSlimDto) string { return strconv.FormatBool(item.Meta.IsLicenseChart) },
		"source":         func(item *license.LicenseSlimDto) string { return string(item.Source) },
		"family":         func(item *license.LicenseSlimDto) string { return item.Meta.Family.Value() },
		"approvalState":  func(item *license.LicenseSlimDto) string { return item.Meta.ApprovalState.Value() },
		"licenseType":    func(item *license.LicenseSlimDto) string { return item.Meta.LicenseType.Value() },
	}

	allClassifications := licensesHandler.ObligationRepository.FindAll(requestSession, false)
	classificationCountMap := make(map[string]license.ClassificationWithCount)
	classificationMap := make(map[string]*obligation2.ObligationDto)

	for _, classification := range allClassifications {
		classificationMap[classification.GetKey()] = classification.ToDto()
	}
	for _, l := range licenses {
		arrayExtractors := map[string]func(*license.LicenseSlimDto) []string{
			"classifications": func(item *license.LicenseSlimDto) []string {
				if len(l.Meta.ObligationsKeyList) == 0 {
					return []string{}
				}
				var names []string

				classifications := getClassificationsByIds(l.Meta.ObligationsKeyList, classificationMap)
				for _, classification := range classifications {
					names = append(names, classification.Name)
					names = append(names, classification.NameDe)
				}
				return names
			},
		}
		dto := l.ToSlimDto()
		var classificationDtos []obligation2.ObligationSlimDto

		for _, c := range l.Meta.ObligationsKeyList {
			classificationDto := getClassificationById(c, classificationMap)
			if classificationDto != nil {
				classificationDtos = append(classificationDtos, *classificationDto.ToSlimDto())
			}
			if entry, exists := classificationCountMap[c]; exists {
				entry.Count++
				classificationCountMap[c] = entry
			} else {
				classificationCountMap[c] = license.ClassificationWithCount{
					Classification: classificationDto,
					Count:          1,
				}
			}
		}

		if len(l.Meta.ObligationsKeyList) == 0 {
			if entry, exists := classificationCountMap[""]; exists {
				entry.Count++
				classificationCountMap[""] = entry
			} else {
				classificationCountMap[""] = license.ClassificationWithCount{
					Count: 1,
				}
			}
		}

		dto.Meta.Classifications = classificationDtos
		dto.Meta.PrevalentClassificationLevel = findPrevalentLevel(classificationDtos)

		possibleCharts[strconv.FormatBool(dto.Meta.IsLicenseChart)]++
		possibleSources[string(dto.Source)]++
		possibleFamilies[string(dto.Meta.Family)]++
		possibleApproval[string(dto.Meta.ApprovalState)]++
		possibleType[string(dto.Meta.LicenseType)]++

		if filter.MatchesCriteria(dto, searchOptions, extractors, arrayExtractors) {
			dtos = append(dtos, dto)
		}
	}

	possibleClassifications := make([]license.ClassificationWithCount, 0, len(classificationCountMap))
	for _, classificationWithCount := range classificationCountMap {
		possibleClassifications = append(possibleClassifications, classificationWithCount)
	}

	sort.Slice(possibleClassifications, func(i, j int) bool {
		return possibleClassifications[i].Count > possibleClassifications[j].Count
	})

	result := license.LicensesResponse{
		Licenses: dtos,
		Count:    len(dtos),
		Meta: license.PossibleFilterValues{
			PossibleCharts:          possibleCharts,
			PossibleSources:         possibleSources,
			PossibleFamilies:        possibleFamilies,
			PossibleApproval:        possibleApproval,
			PossibleType:            possibleType,
			PossibleClassifications: possibleClassifications,
		},
	}

	if searchOptions.ShouldOrder() {
		asc := searchOptions.IsSortAsc()
		key := searchOptions.GetSortKey()
		if key == "name" {
			sorthelper.Sort(result.Licenses, func(dto *license.LicenseSlimDto) string { return dto.Name }, sorthelper.StringLessThan, asc)
		} else if key == "licenseId" {
			sorthelper.Sort(result.Licenses, func(dto *license.LicenseSlimDto) string { return dto.LicenseId }, sorthelper.StringLessThan, asc)
		} else if key == "meta.approvalState" {
			sorthelper.Sort(result.Licenses, func(dto *license.LicenseSlimDto) string { return dto.Meta.ApprovalState.Value() }, sorthelper.StringLessThan, asc)
		} else if key == "meta.family" {
			sorthelper.Sort(result.Licenses, func(dto *license.LicenseSlimDto) string { return dto.Meta.Family.Value() }, sorthelper.StringLessThan, asc)
		} else if key == "meta.licenseType" {
			sorthelper.Sort(result.Licenses, func(dto *license.LicenseSlimDto) string { return dto.Meta.LicenseType.Value() }, sorthelper.StringLessThan, asc)
		} else if key == "source" {
			sorthelper.Sort(result.Licenses, func(dto *license.LicenseSlimDto) string { return string(dto.Source) }, sorthelper.StringLessThan, asc)
		} else if key == "updated" {
			sorthelper.Sort(result.Licenses, func(dto *license.LicenseSlimDto) int64 { return dto.Updated.Unix() }, sorthelper.Int64LessThan, asc)
		} else if key == "meta.classifications" {
			sorthelper.Sort(result.Licenses, func(dto *license.LicenseSlimDto) int64 {
				return obligation2.GetLevelWeight(dto.Meta.PrevalentClassificationLevel)
			}, sorthelper.Int64LessThan, asc)
		} else if key == "meta.isLicenseChart" {
			sorthelper.Sort(result.Licenses, func(dto *license.LicenseSlimDto) bool {
				return dto.Meta.IsLicenseChart
			}, sorthelper.BoolLessThan, asc)
		}
	}

	if searchOptions.HasPaginationActive() && len(result.Licenses) > 0 {
		lowIndex := (searchOptions.GetPage() - 1) * searchOptions.GetItemsPerPage()
		highIndex := lowIndex + searchOptions.GetItemsPerPage()
		if highIndex > len(result.Licenses) {
			highIndex = len(result.Licenses)
		}
		if lowIndex > highIndex {
			lowIndex = 0 // reset page number
		}
		result.Licenses = result.Licenses[lowIndex:highIndex]
	}

	render.JSON(w, r, result)
}

func findPrevalentLevel(classificationDtos []obligation2.ObligationSlimDto) obligation2.WarnLevel {
	prevalentLevel := obligation2.Information
	if len(classificationDtos) == 0 {
		return obligation2.WarnLevel(prevalentLevel)
	}
	for _, dto := range classificationDtos {
		switch strings.ToUpper(string(dto.WarnLevel)) {
		case "ALARM":
			return obligation2.Alarm
		case "WARNING":
			prevalentLevel = obligation2.Warning
		}
	}
	return obligation2.WarnLevel(prevalentLevel)
}

func getClassificationsByIds(idList []string, classificationMap map[string]*obligation2.ObligationDto) []obligation2.ObligationDto {
	result := make([]obligation2.ObligationDto, 0)
	for _, key := range idList {
		classification := getClassificationById(key, classificationMap)
		result = append(result, *classification)
	}
	return result
}

func getClassificationById(id string, classificationMap map[string]*obligation2.ObligationDto) *obligation2.ObligationDto {
	classification, exists := classificationMap[id]
	if exists {
		return classification
	}
	return nil
}

func (licensesHandler *LicensesHandler) LicensesGetAllHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	licenses := licensesHandler.LicenseRepository.FindAll(requestSession, true)
	allClassifications := licensesHandler.ObligationRepository.FindAll(requestSession, false)

	licenseDtos := make([]*license.LicenseSlimDto, 0)
	classificationMap := make(map[string]*obligation2.ObligationDto)
	for _, l := range licenses {
		for _, classification := range allClassifications {
			classificationMap[classification.GetKey()] = classification.ToDto()
		}
		dto := l.ToSlimDto()
		var classificationDtos []obligation2.ObligationSlimDto
		for _, c := range l.Meta.ObligationsKeyList {
			classificationDto := getClassificationById(c, classificationMap)
			if classificationDto != nil {
				classificationDtos = append(classificationDtos, *classificationDto.ToSlimDto())
			}
		}
		dto.Meta.Classifications = classificationDtos
		dto.Meta.PrevalentClassificationLevel = findPrevalentLevel(classificationDtos)

		licenseDtos = append(licenseDtos, dto)
	}
	render.JSON(w, r, licenseDtos)
}

func (licensesHandler *LicensesHandler) GetSpdxLicensesCount(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)

	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowTools.Update {
		exception.ThrowExceptionSendDeniedResponse()
	}

	qbCount := licensesHandler.SpdxLicenseRepository.CountAll(requestSession)
	render.JSON(w, r, CountResponse{
		Count: qbCount,
	})
}

func (licensesHandler *LicensesHandler) GetLicensesDiffs(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)

	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowTools.Update {
		exception.ThrowExceptionSendDeniedResponse()
	}

	licenseDiffs := make([]license.LicenseDiffDto, 0)
	newSpdxLicenses := licensesHandler.SpdxLicenseRepository.FindAll(requestSession, false)
	for _, newSpdxLicense := range newSpdxLicenses {
		licenseId := newSpdxLicense.LicenseId
		oldLicense := licensesHandler.LicenseRepository.FindById(requestSession, licenseId)
		oldLicenseDto := oldLicense.ToDto(requestSession, licensesHandler.ObligationRepository)
		newLicenseDto := newSpdxLicense.ToDto(requestSession, nil)
		licenseDiffs = append(licenseDiffs, license.LicenseDiffDto{
			LicenseId:  licenseId,
			OldLicense: oldLicenseDto,
			NewLicense: newLicenseDto,
		})
	}

	render.JSON(w, r, licenseDiffs)
}

func (licensesHandler *LicensesHandler) GetAllPolicyRulesAssignmentsForThisLicenceHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)

	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowLicense.Read && !rights.AllowPolicy.Update {
		exception.ThrowExceptionSendDeniedResponse()
	}
	licenceId := chi.URLParam(r, "licenceId")
	policyRules := licensesHandler.PolicyRulesRepository.FindAll(requestSession, false)

	var policyRulesForLicense license.PolicyRulesForLicenseDto
	policyRulesForLicense.Id = licenceId
	policyRulesAssignments := make([]license.PolicyRulesAssignmentDto, 0)

	for _, rule := range policyRules {
		status := license.StatusActive
		if !rule.Active {
			status = license.StatusInactive
		}
		if rule.Deprecated {
			status = license.StatusDeprecated
		}
		policyRulesAssignment := license.PolicyRulesAssignmentDto{
			Status:      status,
			Key:         rule.Key,
			Name:        rule.Name,
			Description: rule.Description,
			Type:        license.NOT_SET,
		}
		if helper.Contains(licenceId, rule.ComponentsAllow) {
			policyRulesAssignment.Type = license.ALLOW
		}
		if helper.Contains(licenceId, rule.ComponentsWarn) {
			policyRulesAssignment.Type = license.WARN
		}
		if helper.Contains(licenceId, rule.ComponentsDeny) {
			policyRulesAssignment.Type = license.DENY
		}
		policyRulesAssignments = append(policyRulesAssignments, policyRulesAssignment)
	}
	policyRulesForLicense.PolicyRulesAssignments = policyRulesAssignments

	render.JSON(w, r, policyRulesForLicense)
}

func (licensesHandler *LicensesHandler) UpdateAllPolicyRulesAssignmentsForThisLicenceHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	username, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowPolicy.Update {
		exception.ThrowExceptionSendDeniedResponse()
	}
	licenceId := chi.URLParam(r, "licenceId")
	var policyRulesForLicense license.PolicyRulesForLicenseDto
	validation.DecodeAndValidate(r, &policyRulesForLicense, false)

	keys := make([]string, 0)
	changedRulesMap := make(map[string]license.PolicyRulesAssignmentDto)
	for _, rule := range policyRulesForLicense.PolicyRulesAssignments {
		keys = append(keys, rule.Key)
		changedRulesMap[rule.Key] = rule
	}

	var keyAttributes []database.MatchGroup
	for _, key := range keys {
		keyAttributes = append(keyAttributes,
			database.AttributeMatcher(
				licensesHandler.JobRepository.DatabaseConn().GetKeyAttribute(),
				database.EQ,
				key,
			),
		)
	}

	qc := database.New().SetMatcher(
		database.AndChain(
			database.OrChain(
				keyAttributes...,
			),
			database.AttributeMatcher(
				"Deleted",
				database.EQ,
				false,
			),
			database.AttributeMatcher(
				"Deprecated",
				database.EQ,
				false,
			),
		),
	)

	policyRules := licensesHandler.PolicyRulesRepository.Query(requestSession, qc)

	for _, ruleEntity := range policyRules {
		oldRulesAudit := ruleEntity.ToAudit(requestSession, nil)
		// Do not check for contains otherwise all slices are iterating twice
		if changedRulesMap[ruleEntity.Key].Type == license.NOT_SET {
			ruleEntity.ComponentsAllow = helper.RemoveStrFromSlice(licenceId, ruleEntity.ComponentsAllow)
			ruleEntity.ComponentsWarn = helper.RemoveStrFromSlice(licenceId, ruleEntity.ComponentsWarn)
			ruleEntity.ComponentsDeny = helper.RemoveStrFromSlice(licenceId, ruleEntity.ComponentsDeny)
		}
		if changedRulesMap[ruleEntity.Key].Type == license.ALLOW {
			ruleEntity.ComponentsAllow = append(ruleEntity.ComponentsAllow, licenceId)
			ruleEntity.ComponentsWarn = helper.RemoveStrFromSlice(licenceId, ruleEntity.ComponentsWarn)
			ruleEntity.ComponentsDeny = helper.RemoveStrFromSlice(licenceId, ruleEntity.ComponentsDeny)
		}
		if changedRulesMap[ruleEntity.Key].Type == license.WARN {
			ruleEntity.ComponentsAllow = helper.RemoveStrFromSlice(licenceId, ruleEntity.ComponentsAllow)
			ruleEntity.ComponentsWarn = append(ruleEntity.ComponentsWarn, licenceId)
			ruleEntity.ComponentsDeny = helper.RemoveStrFromSlice(licenceId, ruleEntity.ComponentsDeny)
		}
		if changedRulesMap[ruleEntity.Key].Type == license.DENY {
			ruleEntity.ComponentsAllow = helper.RemoveStrFromSlice(licenceId, ruleEntity.ComponentsAllow)
			ruleEntity.ComponentsWarn = helper.RemoveStrFromSlice(licenceId, ruleEntity.ComponentsWarn)
			ruleEntity.ComponentsDeny = append(ruleEntity.ComponentsDeny, licenceId)
		}
		ruleEntity.Updated = time.Now()
		rulesAudit := ruleEntity.ToAudit(requestSession, nil)
		auditHelper.CreateAndAddAuditEntry(&ruleEntity.Container, username, message.PolicyRulesUpdated, cmp.Diff, rulesAudit, oldRulesAudit)
	}
	licensesHandler.PolicyRulesRepository.UpdateList(requestSession, policyRules)

	responseData := SuccessResponse{
		Success: true,
		Message: "Rules updated",
	}
	render.JSON(w, r, responseData)
}

func (licensesHandler *LicensesHandler) GetCountOfPolicyRuleUsingThisLicenceHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)

	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowLicense.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}
	licenceId := chi.URLParam(r, "licenceId")
	policyRules := licensesHandler.PolicyRulesRepository.FindAll(requestSession, false)
	count := 0
	for _, rule := range policyRules {
		if helper.Contains(licenceId, rule.ComponentsAllow) {
			count++
			continue
		}
		if helper.Contains(licenceId, rule.ComponentsWarn) {
			count++
			continue
		}
		if helper.Contains(licenceId, rule.ComponentsDeny) {
			count++
			continue
		}
	}
	render.JSON(w, r, CountResponse{Count: count})
}

func (licensesHandler *LicensesHandler) GetCountOfLicencesUsingThisObligationHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)

	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowLicense.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}
	obligationKey := chi.URLParam(r, "obligationKey")
	// TODO optimize
	count := licensesHandler.LicenseRepository.CountByObligationKey(requestSession, obligationKey)
	render.JSON(w, r, CountResponse{Count: count})
}

func (licensesHandler *LicensesHandler) LicensesGetHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	licenseIds := chi.URLParam(r, "ids")
	licenseIds, err := url.QueryUnescape(licenseIds)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ErrorKeyRequestParamNotValid, "ids"), zapcore.InfoLevel)
	licenseIdsAsArray := strings.Split(licenseIds, " ")

	entities := make([]license.License, 0)
	for _, licenseId := range licenseIdsAsArray {
		if strings.ToLower(licenseId) != "and" && strings.ToLower(licenseId) != "or" {
			existingLicense := licensesHandler.LicenseRepository.FindById(requestSession, licenseId)
			if existingLicense != nil {
				entities = append(entities, *existingLicense)
			}
		}
	}

	// An empty list of the entities is not an error in this case. Empty list is handled by appropriate info in Frontend.
	// Furthermore there are Common Licenses extracted from Component Details if Component License ID matches the Common License ID
	// Convert to DTO, if an entity found
	licenseDtos := make([]license.LicenseDto, 0)
	if len(entities) > 0 {
		for _, entity := range entities {
			dto := entity.ToDto(requestSession, licensesHandler.ObligationRepository)
			licenseDtos = append(licenseDtos, *dto)
		}
	}

	// convert to json
	render.JSON(w, r, licenseDtos)
}

func (licensesHandler *LicensesHandler) LicenseGetHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	licenseId := licensesHandler.extractLicenceId(r)

	l := licensesHandler.LicenseRepository.FindById(requestSession, licenseId)
	if l == nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.LicenseDataMissing), "")
	}

	licenseDto := l.ToDto(requestSession, licensesHandler.ObligationRepository)

	render.JSON(w, r, licenseDto)
}

func (licensesHandler *LicensesHandler) extractLicenceId(r *http.Request) string {
	licenseIdEncoded := chi.URLParam(r, "id")
	licenseId, err := url.QueryUnescape(licenseIdEncoded)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.ErrorKeyRequestParamNotValid, "id"), zapcore.InfoLevel)
	return licenseId
}

func (licensesHandler *LicensesHandler) LicenseTrailGetAllHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)

	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.AllowLicense.Create && rights.AllowLicense.Update && rights.AllowLicense.Delete) {
		exception.ThrowExceptionSendDeniedResponse()
	}

	licenseId := licensesHandler.extractLicenceId(r)

	l := licensesHandler.LicenseRepository.FindById(requestSession, licenseId)
	if l == nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.DataMissing), "")
	}

	auditTrail := make([]audit.AuditDto, 0)
	for _, item := range l.GetAuditTrail() {
		auditTrail = append(auditTrail, item.ToDto())
	}
	render.JSON(w, r, auditTrail)
}

func (licensesHandler *LicensesHandler) LicenseHeadHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	licenseId := licensesHandler.extractLicenceId(r)

	l := licensesHandler.LicenseRepository.FindByIdCaseInsensitive(requestSession, licenseId)
	if l == nil {
		render.JSON(w, r, FoundResponse{
			Found: false,
		})
		return
	}
	render.JSON(w, r, FoundResponse{
		Found: true,
	})
}

func (licensesHandler *LicensesHandler) AliasHeadHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	licenseId := chi.URLParam(r, "alias")

	currentRefs := licensesHandler.LicenseRepository.GetLicenseRefs(requestSession)
	_, found := currentRefs[strings.TrimSpace(strings.ToLower(licenseId))]
	if !found {
		render.JSON(w, r, FoundResponse{
			Found: false,
		})
		return
	}
	render.JSON(w, r, FoundResponse{
		Found: true,
	})
}

func (licensesHandler *LicensesHandler) LicenseNameHeadHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	licenseId := licensesHandler.extractLicenceId(r)

	l := licensesHandler.LicenseRepository.FindByName(requestSession, licenseId)
	if l == nil {
		render.JSON(w, r, FoundResponse{
			Found: false,
		})
		return
	}
	render.JSON(w, r, FoundResponse{
		Found: true,
	})
}

func (licensesHandler *LicensesHandler) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	username, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowLicense.Delete {
		exception.ThrowExceptionSendDeniedResponse()
	}
	licenseId := licensesHandler.extractLicenceId(r)

	existingLicense := licensesHandler.LicenseRepository.FindById(requestSession, licenseId)
	if existingLicense == nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.LicenseUpdate, licenseId), "")
	}

	oldLicenseAudit := existingLicense.ToAudit(requestSession, licensesHandler.ObligationRepository)

	if existingLicense.Source == license.PublicLicenseDb {
		exception.ThrowExceptionSendDeniedResponse()
	}

	// delete license from policy rule
	policyRules := licensesHandler.PolicyRulesRepository.FindAll(requestSession, false)
	// TODO let this for future use!
	const userHasRightsToDeleteEvenIsUsed = false
	for _, rule := range policyRules {
		wasChanged := false
		if helper.Contains(existingLicense.LicenseId, rule.ComponentsAllow) {
			rule.ComponentsAllow = helper.RemoveStrFromSlice(existingLicense.LicenseId, rule.ComponentsAllow)
			wasChanged = true
		}
		if helper.Contains(existingLicense.LicenseId, rule.ComponentsWarn) {
			rule.ComponentsWarn = helper.RemoveStrFromSlice(existingLicense.LicenseId, rule.ComponentsWarn)
			wasChanged = true
		}
		if helper.Contains(existingLicense.LicenseId, rule.ComponentsDeny) {
			rule.ComponentsDeny = helper.RemoveStrFromSlice(existingLicense.LicenseId, rule.ComponentsDeny)
			wasChanged = true
		}

		if wasChanged {
			if userHasRightsToDeleteEvenIsUsed {
				// TODO On Exception, what should we do? Normally we should rollback the hole transaction, but we have no transaction.
				licensesHandler.PolicyRulesRepository.Update(requestSession, rule)
			} else {
				exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorInUse, "license"))
			}
		}
	}

	existingLicenseAudit := existingLicense.ToAudit(requestSession, licensesHandler.ObligationRepository)

	// License Text changed? -> Take the new Text completely into diff
	if !cmp.Equal(existingLicenseAudit.Text, oldLicenseAudit.Text) {
		oldLicenseAudit.Text = ""
	}

	auditHelper.CreateAndAddAuditEntry(&existingLicense.Container, username, message.LicenseDeleted, audit.DiffWithReporter, existingLicenseAudit, oldLicenseAudit)

	// update first to save auditLog
	licensesHandler.LicenseRepository.Update(requestSession, existingLicense)
	// then delete license
	licensesHandler.LicenseRepository.Delete(requestSession, existingLicense.Key)
	observermngmt.FireEvent(observermngmt.LicenseDeleted, observermngmt.LicenseData{
		RequestSession: requestSession,
		Id:             existingLicense.LicenseId,
	})

	w.WriteHeader(http.StatusOK)
}

func (licensesHandler *LicensesHandler) DeleteSpdxHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	key := chi.URLParam(r, "key")
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowLicense.Delete {
		exception.ThrowExceptionSendDeniedResponse()
	}

	licensesHandler.SpdxLicenseRepository.Delete(requestSession, key)

	w.WriteHeader(http.StatusOK)
}

func (licensesHandler *LicensesHandler) UpdateAcceptedChangesHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	licensesHandler.handleUpdate(requestSession, w, r, true, message.SpdxDatabaseRefreshed)
}

func (licensesHandler *LicensesHandler) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	licensesHandler.handleUpdate(requestSession, w, r, false, message.LicenseUpdated)
}

func (licensesHandler *LicensesHandler) handleUpdate(requestSession *logy.RequestSession, w http.ResponseWriter, r *http.Request, acceptAllChanges bool, auditTitle string) {
	key := chi.URLParam(r, "key")
	username, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowLicense.Update {
		exception.ThrowExceptionSendDeniedResponse()
	}

	var licenseData license.LicenseDto
	validation.DecodeAndValidate(r, &licenseData, false)

	if !validation.IsSpdxIdentifier(licenseData.LicenseId) {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorValidationNotValidSpdxIdentifier))
	}

	existingLicenseById := licensesHandler.LicenseRepository.FindByIdCaseInsensitive(requestSession, licenseData.LicenseId)
	if existingLicenseById != nil && (existingLicenseById.Key != licenseData.Key) {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.DataExistsLicenseId, licenseData.LicenseId), "already exists")
	}

	existingLicense := licensesHandler.LicenseRepository.FindByKey(requestSession, key, false)
	if existingLicense == nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.LicenseUpdate, key), "")
	}

	oldLicenseAudit := existingLicense.ToAudit(requestSession, licensesHandler.ObligationRepository)

	for i := range licenseData.Aliases {
		licenseData.Aliases[i].LicenseId = strings.TrimSpace(licenseData.Aliases[i].LicenseId)
		if !validation.IsSpdxAlias(licenseData.Aliases[i].LicenseId) {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorValidationNotValidSpdxIdentifier, licenseData.Aliases[i].LicenseId))
		}
	}

	aliases := license.AliasesToEntity(licenseData.Aliases)
	insertKeys(aliases)
	currentRefs := licensesHandler.LicenseRepository.GetLicenseRefs(requestSession)
	if !checkForDuplicatesUpdate(existingLicense.LicenseId, licenseData.LicenseId, aliases, currentRefs) {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.IdConflict), "")
	}
	del, new := aliasDiff(existingLicense.Aliases, aliases)
	for _, d := range del {
		observermngmt.FireEvent(observermngmt.LicenseAliasDeleted, observermngmt.AliasData{
			RequestSession: requestSession,
			Id:             existingLicense.LicenseId,
			Alias:          d.LicenseId,
		})
	}
	for _, n := range new {
		observermngmt.FireEvent(observermngmt.LicenseAliasAdded, observermngmt.AliasData{
			RequestSession: requestSession,
			Id:             existingLicense.LicenseId,
			Alias:          n.LicenseId,
		})
	}
	existingLicense.Aliases = aliases
	existingLicense.Updated = time.Now()

	existingLicense.Meta = *licenseData.Meta.ToEntity()
	existingLicense.Text = licenseData.Text
	if existingLicense.Source == license.CUSTOM {
		existingLicense.LicenseId = licenseData.LicenseId
		existingLicense.Name = licenseData.Name
	}

	if acceptAllChanges {
		existingLicense.Name = licenseData.Name
		existingLicense.IsDeprecatedLicenseId = licenseData.IsDeprecatedLicenseId
	}

	existingLicenseAudit := existingLicense.ToAudit(requestSession, licensesHandler.ObligationRepository)

	// License Text changed? -> Take the new Text completely into diff
	if !cmp.Equal(existingLicenseAudit.Text, oldLicenseAudit.Text) {
		oldLicenseAudit.Text = ""
	}

	auditHelper.CreateAndAddAuditEntry(&existingLicense.Container, username, auditTitle, audit.DiffWithReporter, existingLicenseAudit, oldLicenseAudit)
	licensesHandler.LicenseRepository.Update(requestSession, existingLicense)
	render.JSON(w, r, existingLicense.LicenseId)
}

func aliasDiff(prev []license.Alias, updated []license.Alias) (del []license.Alias, new []license.Alias) {
	for _, p := range prev {
		if aliasInArray(p.LicenseId, updated) {
			continue
		}
		del = append(del, p)
	}
	for _, u := range updated {
		if aliasInArray(u.LicenseId, prev) {
			continue
		}
		new = append(new, u)
	}
	return
}

func aliasInArray(needle string, haystack []license.Alias) bool {
	for _, a := range haystack {
		if a.LicenseId != needle {
			continue
		}
		return true
	}
	return false
}

func AliasesToEntity(aliasDto []license.AliasDto) {
	panic("unimplemented")
}

func insertKeys(requested []license.Alias) {
	for i := 0; i < len(requested); i++ {
		if strings.HasPrefix(requested[i].Key, "added") {
			requested[i].Key = uuid.NewString()
		}
	}
}

func checkForDuplicatesCreate(lId string, requested []license.Alias, refList license.LicenseRefs) bool {
	if _, found := refList[lId]; found {
		return false
	}
	for _, search := range requested {
		if _, found := refList[strings.ToLower(search.LicenseId)]; found {
			return false
		}
		if strings.EqualFold(search.LicenseId, lId) {
			return false
		}
		occurences := 0
		for _, a := range requested {
			if a.LicenseId == search.LicenseId {
				occurences++
				if occurences > 1 {
					return false
				}
			}
		}
	}
	return true
}

func checkForDuplicatesUpdate(oldId, newId string, requested []license.Alias, refList license.LicenseRefs) bool {
	if _, found := refList[strings.ToLower(newId)]; found && !strings.EqualFold(oldId, newId) {
		return false
	}
	for _, search := range requested {
		if duplicateFor, found := refList[strings.ToLower(search.LicenseId)]; found && !strings.EqualFold(oldId, duplicateFor.ID) {
			return false
		}
		if strings.EqualFold(search.LicenseId, newId) {
			return false
		}
		occurences := 0
		for _, a := range requested {
			if a.LicenseId == search.LicenseId {
				occurences++
				if occurences > 1 {
					return false
				}
			}
		}
	}
	return true
}

func (licensesHandler *LicensesHandler) LicensePostHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	username, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowLicense.Create {
		exception.ThrowExceptionSendDeniedResponse()
	}

	var licenseData license.LicenseDto
	validation.DecodeAndValidate(r, &licenseData, false)

	if !validation.IsSpdxIdentifier(licenseData.LicenseId) {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorValidationNotValidSpdxIdentifier, licenseData.LicenseId))
	}

	existingLicenseById := licensesHandler.LicenseRepository.FindByIdCaseInsensitive(requestSession, licenseData.LicenseId)
	if existingLicenseById != nil {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.DataExistsLicenseId, licenseData.LicenseId), "already exists")
	}

	for i := range licenseData.Aliases {
		licenseData.Aliases[i].LicenseId = strings.TrimSpace(licenseData.Aliases[i].LicenseId)
		if !validation.IsSpdxAlias(licenseData.Aliases[i].LicenseId) {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorValidationNotValidSpdxIdentifier, licenseData.Aliases[i].LicenseId))
		}
	}

	aliases := license.AliasesToEntity(licenseData.Aliases)
	insertKeys(aliases)
	currentRefs := licensesHandler.LicenseRepository.GetLicenseRefs(requestSession)
	if !checkForDuplicatesCreate(licenseData.LicenseId, aliases, currentRefs) {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.IdConflict))
	}

	licenseData.Key = uuid.New().String()
	licenseData.Updated = time.Now()
	licenseData.Source = license.CUSTOM

	l := licenseData.ToEntity()
	licenseAudit := l.ToAudit(requestSession, licensesHandler.ObligationRepository)
	l.Aliases = aliases
	auditHelper.CreateAndAddAuditEntry(&l.Container, username, message.LicenseCreated, audit.DiffWithReporter, licenseAudit, &license.LicenseAudit{})
	licensesHandler.LicenseRepository.Save(requestSession, l)

	observermngmt.FireEvent(observermngmt.DatabaseEntryAddedOrDeleted, observermngmt.DatabaseSizeChange{
		RequestSession: requestSession,
		CollectionName: license2.LicensesCollectionName,
		Rights:         rights,
		Username:       username,
	})
	observermngmt.FireEvent(observermngmt.LicenseAdded, observermngmt.LicenseData{
		RequestSession: requestSession,
		Id:             l.LicenseId,
	})
	for _, a := range l.Aliases {
		observermngmt.FireEvent(observermngmt.LicenseAdded, observermngmt.AliasData{
			RequestSession: requestSession,
			Id:             l.LicenseId,
			Alias:          a.LicenseId,
		})
	}

	render.JSON(w, r, licenseData)
}

func (licensesHandler *LicensesHandler) LicenseCompareHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	licenseText := simplifyText(string(bodyBytes))

	type licenseWithSimilarityHolder struct {
		License    license.License
		Similarity float64
	}
	// Find all licenses in DB
	allLicenses := licensesHandler.LicenseRepository.FindAll(requestSession, false)
	licenseWithSimilarities := make([]licenseWithSimilarityHolder, 0)

	// Calc similarity with Levenshtein and add to list
	for _, l := range allLicenses {
		simplifiedCompareText := simplifyText(l.Text)

		distance := levenshtein.Distance(licenseText, simplifiedCompareText)
		lenA, lenB := utf8.RuneCountInString(licenseText), utf8.RuneCountInString(simplifiedCompareText)
		similarity := 1 - float64(distance)/math.Max(float64(lenA), float64(lenB))

		licenseWithSimilarities = append(licenseWithSimilarities, licenseWithSimilarityHolder{License: *l, Similarity: similarity})
	}

	// sort DESC by similarity
	sort.SliceStable(licenseWithSimilarities, func(i, j int) bool {
		return licenseWithSimilarities[i].Similarity > licenseWithSimilarities[j].Similarity
	})

	// New list with top 10 items
	topNo := 10
	if len(licenseWithSimilarities) < topNo {
		topNo = len(licenseWithSimilarities)
	}
	topItems := licenseWithSimilarities[:topNo]

	dtoList := make([]license.LicenseWithSimilarityDto, 0)
	for _, item := range topItems {
		dto := item.License.ToDto(requestSession, nil)
		dtoList = append(dtoList, license.LicenseWithSimilarityDto{License: *dto, Similarity: item.Similarity})
	}
	render.JSON(w, r, dtoList)
}

func (licensesHandler *LicensesHandler) CustomLicenses(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)

	auth := extractInternalToken(r.Context())
	if !slices.Contains(auth.Capabilities, internalToken.CustomLicenses) {
		exception.ThrowExceptionSendDeniedResponse()
	}

	raw := r.URL.Query().Get("since")
	var since time.Time
	if raw != "" {
		var err error
		since, err = time.Parse(time.RFC3339, raw)
		if err != nil {
			exception.ThrowExceptionBadRequestResponse()
		}
	}

	matcher := database.AndChain(
		database.AttributeMatcher(
			"Deleted",
			database.EQ,
			false,
		),
		database.OrChain(
			database.AttributeMatcher(
				"meta.approvalState",
				database.EQ,
				license.Approved,
			),
			database.AttributeMatcher(
				"meta.approvalState",
				database.EQ,
				license.Deprecated,
			),
		),
		database.OrChain(
			database.AttributeMatcher(
				"meta.licenseType",
				database.EQ,
				license.OpenSource,
			),
			database.AttributeMatcher(
				"meta.licenseType",
				database.EQ,
				license.PublicDomain,
			),
		),
		database.AttributeMatcher(
			"source",
			database.EQ,
			license.CUSTOM,
		),
		database.AttributeMatcher(
			"licenseId",
			database.LIKE,
			"LicenseRef-MB-%",
		),
	)

	var res license.CustomLicensesDto
	qc := database.New().SetMatcher(matcher).SetKeep([]string{"licenseId"})
	lics := licensesHandler.LicenseRepository.Query(requestSession, qc)
	for _, lic := range lics {
		res.AllIDs = append(res.AllIDs, lic.LicenseId)
	}

	if !since.IsZero() {
		matcher = database.AndChain(
			matcher,
			database.OrChain(
				database.AttributeMatcher(
					"Created",
					database.GT,
					since.String(),
				),
				database.AttributeMatcher(
					"Updated",
					database.GT,
					since.String(),
				),
			),
		)
	}
	qc = database.New().SetMatcher(matcher).SetKeep([]string{"licenseId", "text"})
	lics = licensesHandler.LicenseRepository.Query(requestSession, qc)
	res.Texts = make(map[string]string)
	for _, lic := range lics {
		res.Texts[lic.LicenseId] = lic.Text
	}

	render.JSON(w, r, res)
}

func (licensesHandler *LicensesHandler) LookupHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	var body license.LookupRequestDto
	validation.DecodeAndValidate(r, &body, false)

	var conds []database.MatchGroup
	for _, id := range body.Ids {
		conds = append(conds, database.AttributeMatcher("licenseId", database.EQ, id))
	}
	qc := database.New().SetMatcher(
		database.AndChain(
			database.AttributeMatcher("Deleted", database.EQ, false),
			database.OrChain(conds...),
		),
	).SetKeep([]string{
		"licenseId",
		"name",
	})

	lics := licensesHandler.LicenseRepository.Query(requestSession, qc)
	var res license.LookupResponseDto
	res.Items = license.ToSlimDtos(lics)
	render.JSON(w, r, res)
}

func simplifyText(input string) string {
	re := regexp.MustCompile(`[[:punct:]]`)
	withoutPunctuation := re.ReplaceAllString(input, "")
	withoutLineBreaks := strings.ReplaceAll(withoutPunctuation, "\n", "")
	withoutSpaces := strings.ReplaceAll(withoutLineBreaks, " ", "")
	lowerCaseText := strings.ToLower(withoutSpaces)
	simplified := norm.NFC.String(lowerCaseText)

	return simplified
}
