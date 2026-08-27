// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package analytics

import (
	"sort"
	"strings"
	"time"

	"github.com/eclipse-disuko/disuko/domain/analytics"
	"github.com/eclipse-disuko/disuko/domain/approval"

	"github.com/eclipse-disuko/disuko/domain"
	da "github.com/eclipse-disuko/disuko/domain/analytics"
	prComponents "github.com/eclipse-disuko/disuko/domain/project/components"
	"github.com/eclipse-disuko/disuko/helper/exception"
	analyticsRepo "github.com/eclipse-disuko/disuko/infra/repository/analytics"
	analyticscomponents "github.com/eclipse-disuko/disuko/infra/repository/analyticscomponents"
	"github.com/eclipse-disuko/disuko/infra/repository/analyticslicenses"
	"github.com/eclipse-disuko/disuko/infra/repository/analyticsoccurrences"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/infra/repository/database"
	"github.com/eclipse-disuko/disuko/infra/service/locks"
	"github.com/eclipse-disuko/disuko/logy"
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
	sbomType := options.sbomType
	if sbomType == "" {
		sbomType = da.SbomTypeLatest
	}

	l, acquired := h.lockService.Acquire(locks.Options{
		Blocking: true,
		Key:      lockKey,
		Timeout:  lockTimeout,
	})
	if !acquired {
		return
	}
	defer h.lockService.Release(l)

	qc := database.New().SetMatcher(database.AndChain(
		database.AttributeMatcher(
			"SBomKey",
			database.EQ,
			options.spdxFile.Key,
		),
		database.AttributeMatcher(
			"SBomType",
			database.EQ,
			sbomType,
		),
	))
	existing := len(h.analyticsRepository.Query(options.rs, qc)) > 0
	if existing {
		return
	}

	qc = database.New().SetMatcher(database.AndChain(
		database.AttributeMatcher(
			"ProjectVersionKey",
			database.EQ,
			options.version.Key,
		),
		database.AttributeMatcher(
			"SBomType",
			database.EQ,
			sbomType,
		),
	))
	prev := h.analyticsRepository.Query(options.rs, qc)
	if len(prev) > 0 {
		h.handleSpdxDeleted(options.rs, prev[0].SBomKey, true, sbomType)
	}

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
			if sbomType == da.SbomTypeLatest || sbomType == da.SbomTypeLatestApproved {
				if o, found := occurrenceUpdates[occKey(l.OrigName, sbomType)]; found {
					o.Count++
				} else {
					occurrenceUpdates[occKey(l.OrigName, sbomType)] = &da.Occurrence{
						OrigName:          l.OrigName,
						ReferencedLicense: l.ReferencedLicense,
						Count:             1,
						SBomType:          sbomType,
					}
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
				SBomStatus:     approval.StateInfo(options.spdxFile.ApprovalInfo.Status),
				SBomLastUpdate: options.spdxFile.Updated,
				SBomType:       sbomType,
			}
			bulkSession.AddEnt(&a)
			exception.TryCatch(func() {
				if !h.analyticsComponentsRepository.ExistByName(options.rs, componentName) {
					logy.Infof(options.rs, "storing component %s", componentName)
					h.analyticsComponentsRepository.AddToIndex(options.rs, componentName)
					bulkSessionComps.AddEnt(
						&da.Component{
							RootEntity: domain.NewRootEntity(),
							Name:       componentName,
						},
					)
				}
			}, func(e exception.Exception) {
				logy.Infof(options.rs, "failed to store component %s", e.Error)
			})
			exception.TryCatch(func() {
				if !h.analyticsLicensesRepository.ExistByName(options.rs, licenseName) {
					logy.Infof(options.rs, "storing license %s", licenseName)
					h.analyticsLicensesRepository.AddToIndex(options.rs, licenseName)
					bulkSessionLic.AddEnt(
						&da.License{
							RootEntity: domain.NewRootEntity(),
							Name:       licenseName,
						},
					)
				}
			}, func(e exception.Exception) {
				logy.Infof(options.rs, "failed to store license %s", e.Error)
			})

		}
	}
	h.processOccurrences(options.rs, occurrenceUpdates)
}

func (h *DbHandler) handleSpdxDeleted(session *logy.RequestSession, key string, alreadyAcquired bool, sbomType da.SbomType) {
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

	matcher := database.AttributeMatcher(
		"SBomKey",
		database.EQ,
		key,
	)
	if sbomType != "" {
		matcher = database.AndChain(
			matcher,
			database.AttributeMatcher(
				"SBomType",
				database.EQ,
				sbomType,
			),
		)
	}
	qc := database.New().SetMatcher(matcher)
	existing := h.analyticsRepository.Query(session, qc)
	if len(existing) == 0 {
		return
	}

	bulkSession := h.analyticsRepository.StartSession(base.DeleteSession, 3000)
	defer bulkSession.EndSession()
	occurrenceUpdates := make(map[string]*analytics.Occurrence)
	for _, result := range existing {
		entrySbomType := result.SBomType
		if entrySbomType == "" {
			entrySbomType = da.SbomTypeLatest
		}
		if entrySbomType != da.SbomTypeLatest && entrySbomType != da.SbomTypeLatestApproved {
			bulkSession.AddEnt(result)
			continue
		}
		for _, l := range result.Licenses.List {
			if l.OrigName != result.EntryLicense {
				continue
			}
			if o, found := occurrenceUpdates[occKey(l.OrigName, entrySbomType)]; found {
				o.Count++
			} else {
				occurrenceUpdates[occKey(l.OrigName, entrySbomType)] = &da.Occurrence{
					OrigName:          l.OrigName,
					ReferencedLicense: l.ReferencedLicense,
					Count:             1,
					SBomType:          entrySbomType,
				}
			}
			break
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
		if i := occIndex(curr, u.OrigName, u.SBomType); i != -1 {
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
		i := occIndex(curr, u.OrigName, u.SBomType)
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

func occIndex(haystack []*da.Occurrence, needle string, sbomType da.SbomType) int {
	for i, o := range haystack {
		entrySbomType := o.SBomType
		if entrySbomType == "" {
			entrySbomType = da.SbomTypeLatest
		}
		if o.OrigName == needle && entrySbomType == sbomType {
			return i
		}
	}
	return -1
}

func occKey(origName string, sbomType da.SbomType) string {
	return string(sbomType) + "|" + origName
}

func (h *DbHandler) HandleSpdxDeleted(session *logy.RequestSession, key string) {
	h.handleSpdxDeleted(session, key, false, "")
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
	h.analyticsLicensesRepository.Reset()
	h.analyticsComponentsRepository.DatabaseConn().Truncate()
	h.analyticsComponentsRepository.Reset()
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
	sbomType := options.SbomType
	if sbomType == "" {
		sbomType = da.SbomTypeLatest
	}
	logy.Infof(options.Rs, "searching for component %s and license %s", options.Component, options.License)
	foundComponents := h.analyticsRepository.FindByNameAndProjectKeysAndLicense(
		options.Rs,
		options.Component,
		options.ProjectKeys,
		options.License,
		sbomType,
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
			SBomType:           c.SBomType,
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
	for _, e := range existing {
		e.ReferencedLicense = referencedName
		h.analyticsOccurrencesRepository.Update(session, e)
	}
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
	for _, e := range existing {
		e.ReferencedLicense = ""
		h.analyticsOccurrencesRepository.Update(session, e)
	}
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

func (h *DbHandler) Occurrences(session *logy.RequestSession, sbomType analytics.SbomType) []*analytics.Occurrence {
	all := h.analyticsOccurrencesRepository.FindAll(session, false)
	if sbomType == "" || sbomType == da.SbomTypeLatestAndApproved {
		return all
	}
	filtered := make([]*analytics.Occurrence, 0, len(all))
	for _, o := range all {
		entrySbomType := o.SBomType
		if entrySbomType == "" {
			entrySbomType = da.SbomTypeLatest
		}
		if entrySbomType == sbomType {
			filtered = append(filtered, o)
		}
	}
	return filtered
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
