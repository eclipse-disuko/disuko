// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package analytics

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"mercedes-benz.ghe.com/foss/disuko/domain/analytics"
	"mercedes-benz.ghe.com/foss/disuko/domain/user/approval"

	"mercedes-benz.ghe.com/foss/disuko/domain"
	da "mercedes-benz.ghe.com/foss/disuko/domain/analytics"
	prComponents "mercedes-benz.ghe.com/foss/disuko/domain/project/components"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	analyticsRepo "mercedes-benz.ghe.com/foss/disuko/infra/repository/analytics"
	analyticscomponents "mercedes-benz.ghe.com/foss/disuko/infra/repository/analyticscomponents"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/analyticslicenses"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/analyticsoccurrences"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/database"
	"mercedes-benz.ghe.com/foss/disuko/infra/service/locks"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

const (
	lockKey     = "analock"
	lockTimeout = time.Hour
)

type DbHandler struct {
	analyticsRepository            analyticsRepo.IAnalyticsRepository
	analyticsComponentsRepository  analyticscomponents.IComponentsRepository
	analyticsLicensesRepository    analyticslicenses.ILicensesRepository
	analyticsOccurrencesRepository analyticsoccurrences.IOccurrencesRepository
	lockService                    *locks.Service
}

func InitDbHandler(
	analyticsRepo analyticsRepo.IAnalyticsRepository,
	analyticsCompRepo analyticscomponents.IComponentsRepository,
	analyticsLicRepo analyticslicenses.ILicensesRepository,
	analyticOccRepo analyticsoccurrences.IOccurrencesRepository,
	lockService *locks.Service,
) *DbHandler {
	return &DbHandler{
		analyticsRepository:            analyticsRepo,
		analyticsComponentsRepository:  analyticsCompRepo,
		analyticsLicensesRepository:    analyticsLicRepo,
		analyticsOccurrencesRepository: analyticOccRepo,
		lockService:                    lockService,
	}
}

func (h *DbHandler) HandleSpdxAdded(options SpdxAddedOptions) {
	l, acquired := h.lockService.Acquire(locks.Options{
		Blocking: true,
		Key:      lockKey,
		Timeout:  lockTimeout,
	})
	if !acquired {
		return
	}
	defer h.lockService.Release(l)

	qc := database.New().SetMatcher(database.AttributeMatcher(
		"SBomKey",
		database.EQ,
		options.spdxFile.Key,
	))
	existing := len(h.analyticsRepository.Query(options.rs, qc)) > 0
	if existing {
		return
	}

	qc = database.New().SetMatcher(database.AttributeMatcher(
		"ProjectVersionKey",
		database.EQ,
		options.version.Key,
	))
	prev := h.analyticsRepository.Query(options.rs, qc)
	if len(prev) > 0 {
		h.handleSpdxDeleted(options.rs, prev[0].SBomKey, true)
	}

	licenses := make(map[string]bool)
	components := make(map[string]bool)

	bulkSession := h.analyticsRepository.StartSession(base.UpdateSession, 3000)
	defer bulkSession.EndSession()
	bulkSessionComps := h.analyticsComponentsRepository.StartSession(base.UpdateSession, 3000)
	defer bulkSessionComps.EndSession()
	bulkSessionLic := h.analyticsLicensesRepository.StartSession(base.UpdateSession, 3000)
	defer bulkSessionLic.EndSession()
	occurrenceUpdates := make(map[string]*analytics.Occurrence)
	for _, result := range options.evalRes.Results {
		if result.Component.Type == prComponents.FILE || result.Component.Type == prComponents.SNIPPET {
			logy.Infof(options.rs, "ignoring component of type %s: %s", result.Component.Type, result.Component.Name)
			continue
		}
		for _, l := range result.Component.GetLicensesEffective().List {
			if o, found := occurrenceUpdates[l.OrigName]; found {
				o.Count++
			} else {
				occurrenceUpdates[l.OrigName] = &da.Occurrence{
					OrigName:          l.OrigName,
					ReferencedLicense: l.ReferencedLicense,
					Count:             1,
				}
			}
			licenseName := l.OrigName
			componentName := strings.ToLower(result.Component.Name)
			deptId := ""
			if options.project.CustomerMeta.DeptId != "" {
				deptId = options.project.CustomerMeta.DeptId
			}
			if options.parent != nil && options.parent.CustomerMeta.DeptId != "" {
				deptId = options.parent.CustomerMeta.DeptId
			}
			var responsibleName string
			if res := options.project.ProjectResponsible(); res != nil {
				responsibleName = res.UserId
			}
			a := da.Analytics{
				RootEntity:  domain.NewRootEntity(),
				ProjectKey:  options.project.Key,
				ProjectName: options.project.Name,
				Responsible: responsibleName,

				ProjectVersionKey:  options.version.Key,
				ProjectVersionName: options.version.Name,

				OwnerDeptId: deptId,

				ComponentName:    componentName,
				ComponentVersion: result.Component.Version,

				LicenseConcluded: result.Component.License,
				LicenseDeclared:  result.Component.LicenseDeclared,
				Licenses:         result.Component.GetLicensesEffective(),
				EntryLicense:     licenseName,

				SBomKey:        options.spdxFile.Key,
				SBomName:       options.spdxFile.MetaInfo.Name,
				SBomStatus:     approval.ApprovalStatus(options.spdxFile.ApprovalInfo.Status),
				SBomLastUpdate: options.spdxFile.Updated,
			}
			bulkSession.AddEnt(&a)
			exception.TryCatch(func() {
				componentVersion := strings.ToLower(result.Component.Version)
				key := fmt.Sprintf("%s-%s", componentName, componentVersion)
				if _, found := components[key]; !found {
					byNameAndVersion := h.analyticsComponentsRepository.FindByNameAndVersion(options.rs, componentName, componentVersion)
					if len(byNameAndVersion) == 0 {
						h.analyticsComponentsRepository.AddToIndex(options.rs, componentName)
						bulkSessionComps.AddEnt(
							&da.Component{
								RootEntity: domain.NewRootEntity(),
								Name:       componentName,
								Version:    componentVersion,
							},
						)
					}
					components[key] = true
				}
			}, func(e exception.Exception) {
				logy.Infof(options.rs, "failed to store component %s", e.Error)
			})
			exception.TryCatch(func() {
				if _, found := licenses[licenseName]; !found {
					byName := h.analyticsLicensesRepository.FindByName(options.rs, licenseName)
					if len(byName) == 0 {
						logy.Infof(options.rs, "storing license %s", licenseName)
						h.analyticsLicensesRepository.AddToIndex(options.rs, licenseName)
						bulkSessionLic.AddEnt(
							&da.License{
								RootEntity: domain.NewRootEntity(),
								Name:       licenseName,
							},
						)
					}
					licenses[licenseName] = true
				}
			}, func(e exception.Exception) {
				logy.Infof(options.rs, "failed to store license %s", e.Error)
			})

		}
	}
	h.processOccurrences(options.rs, occurrenceUpdates)
}

func (h *DbHandler) handleSpdxDeleted(session *logy.RequestSession, key string, alreadyAcquired bool) {
	if !alreadyAcquired {
		l, acquired := h.lockService.Acquire(locks.Options{
			Blocking: true,
			Key:      lockKey,
			Timeout:  lockTimeout,
		})
		if !acquired {
			return
		}
		defer h.lockService.Release(l)
	}

	qc := database.New().SetMatcher(database.AttributeMatcher(
		"SBomKey",
		database.EQ,
		key,
	))
	existing := h.analyticsRepository.Query(session, qc)
	if len(existing) == 0 {
		return
	}

	bulkSession := h.analyticsRepository.StartSession(base.DeleteSession, 3000)
	defer bulkSession.EndSession()
	occurrenceUpdates := make(map[string]*analytics.Occurrence)
	for _, result := range existing {
		for _, l := range result.Licenses.List {
			if o, found := occurrenceUpdates[l.OrigName]; found {
				o.Count++
			} else {
				occurrenceUpdates[l.OrigName] = &da.Occurrence{
					OrigName:          l.OrigName,
					ReferencedLicense: l.ReferencedLicense,
					Count:             1,
				}
			}
		}

		bulkSession.AddEnt(result)
	}
	h.processOccurrencesDeletion(session, occurrenceUpdates)
}

func (h *DbHandler) processOccurrences(session *logy.RequestSession, updates map[string]*analytics.Occurrence) {
	curr := h.analyticsOccurrencesRepository.FindAll(session, false)
	bulkSession := h.analyticsOccurrencesRepository.StartSession(base.UpdateSession, 100)
	defer bulkSession.EndSession()
	for _, u := range updates {
		var new *analytics.Occurrence
		if i := occIndex(curr, u.OrigName); i != -1 {
			new = curr[i]
			new.ReferencedLicense = u.ReferencedLicense
			new.Count += u.Count
			h.analyticsOccurrencesRepository.Update(session, new)
		} else {
			new = u
			new.RootEntity = domain.NewRootEntity()
			bulkSession.AddEnt(new)
		}
	}
}

func (h *DbHandler) processOccurrencesDeletion(session *logy.RequestSession, updates map[string]*analytics.Occurrence) {
	curr := h.analyticsOccurrencesRepository.FindAll(session, false)
	bulkSession := h.analyticsOccurrencesRepository.StartSession(base.DeleteSession, 100)
	defer bulkSession.EndSession()
	for _, u := range updates {
		i := occIndex(curr, u.OrigName)
		if i == -1 {
			continue
		}
		if curr[i].Count-u.Count <= 0 {
			bulkSession.AddEnt(curr[i])
		}
		new := curr[i]
		new.Count -= u.Count
		h.analyticsOccurrencesRepository.Update(session, new)
	}
}

func occIndex(haystack []*da.Occurrence, needle string) int {
	for i, o := range haystack {
		if o.OrigName == needle {
			return i
		}
	}
	return -1
}

func (h *DbHandler) HandleSpdxDeleted(session *logy.RequestSession, key string) {
	h.handleSpdxDeleted(session, key, false)
}

func (h *DbHandler) Reset() {
	l, acquired := h.lockService.Acquire(locks.Options{
		Blocking: true,
		Key:      lockKey,
		Timeout:  lockTimeout,
	})
	if !acquired {
		return
	}
	defer h.lockService.Release(l)

	h.analyticsRepository.DatabaseConn().Truncate()
	h.analyticsLicensesRepository.DatabaseConn().Truncate()
	h.analyticsComponentsRepository.DatabaseConn().Truncate()
	h.analyticsOccurrencesRepository.DatabaseConn().Truncate()
}

func (h *DbHandler) ResetWithStatus(statusChannel chan string) {
	l, acquired := h.lockService.Acquire(locks.Options{
		Blocking: true,
		Key:      lockKey,
		Timeout:  lockTimeout,
	})
	if !acquired {
		return
	}
	defer h.lockService.Release(l)

	statusChannel <- "deleting analytics repo..."
	h.analyticsRepository.DatabaseConn().Truncate()
	statusChannel <- "deleting analytics license repo..."
	h.analyticsLicensesRepository.DatabaseConn().Truncate()
	statusChannel <- "deleting analytics components repo..."
	h.analyticsComponentsRepository.DatabaseConn().Truncate()
	statusChannel <- "deleting analytics occurrences repo..."
	h.analyticsOccurrencesRepository.DatabaseConn().Truncate()
}

func (h *DbHandler) HandleSearch(options SearchOptions) analytics.ResponseAnalyticsSearch {
	logy.Infof(options.Rs, "searching for component %s and license %s", options.Component, options.License)
	foundComponents := h.analyticsRepository.FindByNameAndProjectKeysAndLicense(options.Rs,
		options.Component,
		options.ProjectKeys,
		options.License,
		options.Offset,
		options.Limit,
		options.SortCol,
		options.Asc,
	)
	items := make([]analytics.SearchResponseItem, 0)
	for _, c := range foundComponents {
		responseItem := analytics.SearchResponseItem{
			Name:               c.ProjectName,
			Key:                c.ProjectKey,
			ComponentName:      c.ComponentName,
			ComponentVersion:   c.ComponentVersion,
			ProjectVersionName: c.ProjectVersionName,
			ProjectVersionKey:  c.ProjectVersionKey,
			Responsible:        c.Responsible,
			LicenseConcluded:   c.LicenseConcluded,
			LicenseDeclared:    c.LicenseDeclared,
			EntryLicense:       c.EntryLicense,
			SBomName:           c.SBomName,
			SBomStatus:         c.SBomStatus,
			LastUpdate:         c.SBomLastUpdate,
			OwnerDeptId:        c.OwnerDeptId,
		}
		items = append(items, responseItem)
	}
	count := len(foundComponents)
	searchResponse := analytics.ResponseAnalyticsSearch{
		Success: true,
		Items:   items,
		Count:   100000,
	}
	logy.Infof(options.Rs, "found %d components", count)
	return searchResponse
}

func (h *DbHandler) HandleComponentSearch(session *logy.RequestSession, component string, exact bool) analytics.ResponseComponentsSearch {
	componentNames := h.analyticsComponentsRepository.SearchByName(session, strings.ToLower(component), exact)
	sort.Strings(componentNames)
	searchResponse := analytics.ResponseComponentsSearch{
		Components: componentNames,
	}

	logy.Infof(session, "found %d components", len(searchResponse.Components))
	return searchResponse
}

func (h *DbHandler) HandleLicenseIdAdded(session *logy.RequestSession, origName, referencedName string) {
	l, acquired := h.lockService.Acquire(locks.Options{
		Blocking: true,
		Key:      lockKey,
		Timeout:  lockTimeout,
	})
	if !acquired {
		return
	}
	defer h.lockService.Release(l)
	qc := database.New().SetMatcher(database.AttributeMatcher(
		"OrigName",
		database.EQ,
		origName,
	))
	existing := h.analyticsOccurrencesRepository.Query(session, qc)
	if len(existing) == 0 {
		return
	}
	existing[0].ReferencedLicense = referencedName
	h.analyticsOccurrencesRepository.Update(session, existing[0])
}

func (h *DbHandler) HandleLicenseIdDeleted(session *logy.RequestSession, id string) {
	l, acquired := h.lockService.Acquire(locks.Options{
		Blocking: true,
		Key:      lockKey,
		Timeout:  lockTimeout,
	})
	if !acquired {
		return
	}
	defer h.lockService.Release(l)
	qc := database.New().SetMatcher(database.AttributeMatcher(
		"ReferencedLicense",
		database.EQ,
		id,
	))
	existing := h.analyticsOccurrencesRepository.Query(session, qc)
	if len(existing) == 0 {
		return
	}
	existing[0].ReferencedLicense = ""
	h.analyticsOccurrencesRepository.Update(session, existing[0])
}

func (h *DbHandler) HandleLicenseSearch(session *logy.RequestSession, license string, exact bool) da.ResponseLicensesSearch {
	logy.Infof(session, "searching for license %s", license)
	licenseNames := h.analyticsLicensesRepository.SearchLicenceByName(session, strings.ToLower(license), exact)
	sort.Strings(licenseNames)
	searchResponse := analytics.ResponseLicensesSearch{
		Licenses: licenseNames,
	}

	logy.Infof(session, "found %d licenses", len(searchResponse.Licenses))
	return searchResponse
}

func (h *DbHandler) Occurrences(session *logy.RequestSession) []*analytics.Occurrence {
	return h.analyticsOccurrencesRepository.FindAll(session, false)
}

func (h *DbHandler) HandleCompanyChanged(session *logy.RequestSession, prKey string, companyId string) {
	l, acquired := h.lockService.Acquire(locks.Options{
		Blocking: true,
		Key:      lockKey,
		Timeout:  lockTimeout,
	})
	if !acquired {
		return
	}
	defer h.lockService.Release(l)
	qc := database.New().SetMatcher(database.AttributeMatcher(
		"ProjectKey",
		database.EQ,
		prKey,
	))
	as := h.analyticsRepository.Query(session, qc)
	for _, a := range as {
		a.OwnerDeptId = companyId
		h.analyticsRepository.Update(session, a)
	}
}

func (h *DbHandler) HandleResponsibleChanged(session *logy.RequestSession, prKey string, responsible string) {
	l, acquired := h.lockService.Acquire(locks.Options{
		Blocking: true,
		Key:      lockKey,
		Timeout:  lockTimeout,
	})
	if !acquired {
		return
	}
	defer h.lockService.Release(l)
	qc := database.New().SetMatcher(database.AttributeMatcher(
		"ProjectKey",
		database.EQ,
		prKey,
	))
	as := h.analyticsRepository.Query(session, qc)
	for _, a := range as {
		a.Responsible = responsible
		h.analyticsRepository.Update(session, a)
	}
}
