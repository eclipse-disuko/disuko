// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package licenserefresh

import (
	"io"
	"net/http"
	"strings"
	"time"

	"mercedes-benz.ghe.com/foss/disuko/helper/client_utils"
	"mercedes-benz.ghe.com/foss/disuko/scheduler"

	"github.com/google/go-cmp/cmp"
	"github.com/tidwall/gjson"
	"mercedes-benz.ghe.com/foss/disuko/conf"
	"mercedes-benz.ghe.com/foss/disuko/domain"
	"mercedes-benz.ghe.com/foss/disuko/domain/audit"
	"mercedes-benz.ghe.com/foss/disuko/domain/job"
	"mercedes-benz.ghe.com/foss/disuko/domain/license"
	auditHelper "mercedes-benz.ghe.com/foss/disuko/helper/audit"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
	"mercedes-benz.ghe.com/foss/disuko/helper/validation"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	licenseRepo "mercedes-benz.ghe.com/foss/disuko/infra/repository/license"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/obligation"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/spdx_license"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type Job struct {
	licenseRepository      licenseRepo.ILicensesRepository
	obligationRepository   obligation.IObligationRepository
	spdxLicensesRepository spdx_license.ISpdxLicensesRepository
	retriever              IRetrieveJsonBytes
}

func Init(lr licenseRepo.ILicensesRepository, or obligation.IObligationRepository, sr spdx_license.ISpdxLicensesRepository, retriever IRetrieveJsonBytes) *Job {
	if retriever == nil {
		retriever = &retrieveJsonBytes{}
	}
	return &Job{
		licenseRepository:      lr,
		obligationRepository:   or,
		spdxLicensesRepository: sr,
		retriever:              retriever,
	}
}

// TODO: refactor
func (j *Job) Execute(rs *logy.RequestSession, info job.Job) scheduler.ExecutionResult {
	var log job.Log
	log.AddEntry(job.Info, "started")
	s := j.spdxLicensesRepository.StartSession(base.DeleteSession, 2000)
	allEntities := j.spdxLicensesRepository.FindAll(rs, false)
	for _, key := range allEntities {
		s.AddEnt(key)
	}
	s.EndSession()

	url := conf.Config.SPDXLicense.LicensesInfoPath
	licenseFile, err := j.retriever.RetrieveJsonBytesFromUrl(url)
	if err != nil {
		log.AddEntry(job.Error, "Cannot retrieve main licenses file, error=%s", err.Error())
		return scheduler.ExecutionResult{
			Success: false,
			Log:     log,
		}
	}
	licenseJson := gjson.GetBytes(licenseFile, "licenses.#.licenseId")
	licenseMetaListJson := gjson.GetBytes(licenseFile, "licenses").Array()
	licenseIdsArray := licenseJson.Array()
	countNewLicenses := 0
	countUnchangedLicenses := 0
	countUpdatedLicenses := 0
	countDiffLicenses := 0
	countAllLoadedLicenses := len(licenseIdsArray)
	updatedLicenseIds := make([]string, 0)
	diffLicenseIds := make([]string, 0)
	newLicenseIds := make([]string, 0)
	errorLicenseIds := make([]string, 0)
	var customRes interface{}
	for _, licenseMetaJson := range licenseMetaListJson {
		hasError := false
		licenseId := licenseMetaJson.Get("licenseId").String()
		detailsUrl := conf.Config.SPDXLicense.LicenseDetailsDirPath + licenseId + ".json"
		licenseDetailsFile, err := j.retriever.RetrieveJsonBytesFromUrl(detailsUrl)
		if err != nil {
			log.AddEntry(job.Error, "Cannot retrieve details license file for license=%s, error=%s", licenseId, err.Error())
			return scheduler.ExecutionResult{
				Success: false,
				Log:     log,
			}
		}
		licenseUrl := licenseMetaJson.Get("reference").String()

		newLicense := &license.License{
			RootEntity:            domain.NewRootEntity(),
			IsDeprecatedLicenseId: gjson.GetBytes(licenseDetailsFile, "isDeprecatedLicenseId").Bool(),
			Name:                  strings.TrimSpace(gjson.GetBytes(licenseDetailsFile, "name").String()),
			LicenseId:             strings.TrimSpace(gjson.GetBytes(licenseDetailsFile, "licenseId").String()),
			Text:                  gjson.GetBytes(licenseDetailsFile, "licenseText").String(),
			Source:                license.PublicLicenseDb,
			Meta: license.MetaData{
				Family:             license.NotDeclared,
				ApprovalState:      license.NotSet,
				ReviewState:        license.NotReviewed,
				ReviewDateStr:      "",
				ObligationsKeyList: make([]string, 0),
				LicenseUrl:         licenseUrl,
				SourceUrl:          detailsUrl,
			},
		}

		if len(newLicense.Name) < 1 {
			errorLicenseIds = append(errorLicenseIds, newLicense.LicenseId+" has empty name - "+detailsUrl)
			hasError = true
		}
		if len(newLicense.LicenseId) < 1 {
			errorLicenseIds = append(errorLicenseIds, newLicense.Name+" has empty LicenseId - "+detailsUrl)
			hasError = true
		}
		if !validation.IsSpdxIdentifier(newLicense.LicenseId) {
			errorLicenseIds = append(errorLicenseIds, newLicense.Name+" ("+newLicense.LicenseId+") has invalid LicenseId - "+detailsUrl)
			hasError = true
		}

		if !hasError {
			existingLicense := j.licenseRepository.FindById(rs, newLicense.LicenseId)
			licenseExist := err == nil && existingLicense != nil

			if licenseExist {
				var oldLicenseAudit = existingLicense.ToAudit(rs, j.obligationRepository)

				logMessage := ""
				diff := false
				updated := false
				licenseDiff := ""
				licenseUpdated := ""
				//Name
				if len(existingLicense.Name) > 0 {
					if existingLicense.Name != newLicense.Name {
						diff = true
						logMessage = "[license import] license (" + newLicense.LicenseId + ") has changed: name"
						licenseDiff += existingLicense.LicenseId + " [Name"
					}
				} else {
					if len(newLicense.Name) > 0 {
						updated = true
						existingLicense.Name = newLicense.Name
						logMessage = "[license import] license (" + newLicense.LicenseId + ") was updated: name"
						licenseUpdated += getAIfEmpty(licenseUpdated, existingLicense.LicenseId+" [Name", ", Name")
					}
				}

				//Text
				if len(existingLicense.Text) > 0 {
					if existingLicense.Text != newLicense.Text {
						diff = true
						logMessage = "[license import] license (" + newLicense.LicenseId + ") has changed: text"
						licenseDiff += getAIfEmpty(licenseDiff, existingLicense.LicenseId+" [Text", ", Text")
					}
				} else {
					if len(newLicense.Text) > 0 {
						updated = true
						existingLicense.Text = newLicense.Text
						logMessage = "[license import] license (" + newLicense.LicenseId + ") was updated: text"
						licenseUpdated += getAIfEmpty(licenseUpdated, existingLicense.LicenseId+" [Text", ", Text")
					}
				}

				//LicenseUrl
				if len(existingLicense.Meta.LicenseUrl) > 0 {
					if existingLicense.Meta.LicenseUrl != newLicense.Meta.LicenseUrl {
						diff = true
						logMessage = "[license import] license (" + newLicense.LicenseId + ") has changed: licenseUrl"
						licenseDiff += getAIfEmpty(licenseDiff, existingLicense.LicenseId+" [LicenseUrl", ", LicenseUrl")
					}
				} else {
					if len(newLicense.Meta.LicenseUrl) > 0 {
						updated = true
						existingLicense.Meta.LicenseUrl = newLicense.Meta.LicenseUrl
						logMessage = "[license import] license (" + newLicense.LicenseId + ") was updated: licenseUrl"
						licenseUpdated += getAIfEmpty(licenseUpdated, existingLicense.LicenseId+" [LicenseUrl", ", LicenseUrl")
					}
				}

				//SourceUrl
				if len(existingLicense.Meta.SourceUrl) > 0 {
					if existingLicense.Meta.SourceUrl != newLicense.Meta.SourceUrl {
						diff = true
						logMessage = "[license import] license (" + newLicense.LicenseId + ") has changed: sourceUrl"
						licenseDiff += getAIfEmpty(licenseDiff, existingLicense.LicenseId+" [SourceUrl", ", SourceUrl")
					}
				} else {
					if len(newLicense.Meta.SourceUrl) > 0 {
						updated = true
						existingLicense.Meta.SourceUrl = newLicense.Meta.SourceUrl
						logMessage = "[license import] license (" + newLicense.LicenseId + ") was updated: sourceUrl"
						licenseUpdated += getAIfEmpty(licenseUpdated, existingLicense.LicenseId+" [SourceUrl", ", SourceUrl")
					}
				}

				// spdx deprecation status
				if existingLicense.IsDeprecatedLicenseId {
					if !newLicense.IsDeprecatedLicenseId {
						diff = true
						logMessage = "[license import] license (" + newLicense.LicenseId + ") has changed: IsDeprecatedLicenseId"
						licenseDiff += getAIfEmpty(licenseDiff, existingLicense.LicenseId+" [IsDeprecatedLicenseId", ", IsDeprecatedLicenseId")
					}
				} else {
					if newLicense.IsDeprecatedLicenseId {
						updated = true
						existingLicense.IsDeprecatedLicenseId = newLicense.IsDeprecatedLicenseId
						logMessage = "[license import] license (" + newLicense.LicenseId + ") was updated: IsDeprecatedLicenseId"
						licenseUpdated += getAIfEmpty(licenseUpdated, existingLicense.LicenseId+" [IsDeprecatedLicenseId", ", IsDeprecatedLicenseId")
					}
				}

				if len(logMessage) > 0 {
					logy.Warnw(rs, logMessage)
				}

				if updated {
					countUpdatedLicenses++
					licenseUpdated += getAIfEmpty(licenseUpdated, "", "]")
					updatedLicenseIds = append(updatedLicenseIds, licenseUpdated)
					existingLicense.Updated = time.Now()
					actualizedLicenseAudit := existingLicense.ToAudit(rs, j.obligationRepository)
					if !cmp.Equal(actualizedLicenseAudit.Text, oldLicenseAudit.Text) {
						oldLicenseAudit.Text = ""
					}
					auditHelper.CreateAndAddAuditEntry(&existingLicense.Container, "SYSTEM", message.LicenseUpdated, audit.DiffWithReporter, actualizedLicenseAudit, oldLicenseAudit)
					j.licenseRepository.Update(rs, existingLicense)
					logy.Warnw(rs, "[license import] license ("+newLicense.LicenseId+") updated")
				} else if diff {
					countDiffLicenses++
					licenseDiff += getAIfEmpty(licenseDiff, "", "]") + " - " + detailsUrl
					diffLicenseIds = append(diffLicenseIds, licenseDiff)
					j.spdxLicensesRepository.Save(rs, newLicense)
				} else {
					countUnchangedLicenses++
				}
			} else {
				countNewLicenses++
				newLicenseIds = append(newLicenseIds, newLicense.LicenseId)

				licenseAudit := newLicense.ToAudit(rs, j.obligationRepository)
				auditHelper.CreateAndAddAuditEntry(&newLicense.Container, "SYSTEM", message.LicenseCreated, audit.DiffWithReporter, licenseAudit, &license.LicenseAudit{})
				j.licenseRepository.Save(rs, newLicense)
				logy.Warnw(rs, "[license import] new license ("+newLicense.LicenseId+") imported")
			}
		}

		countHandledLicenses := countNewLicenses + countUnchangedLicenses + countDiffLicenses + countUpdatedLicenses + len(errorLicenseIds)
		res := struct {
			Added       int      `json:"added"`
			Unchanged   int      `json:"unchanged"`
			Changed     int      `json:"changed"`
			Differences int      `json:"differences"`
			Errors      int      `json:"errors"`
			Handled     int      `json:"handled"`
			Total       int      `json:"total"`
			AddedLics   []string `json:"addedLics"`
			UpdatedLics []string `json:"updatedLics"`
			DiffLics    []string `json:"diffLics"`
			ErrorLics   []string `json:"errorLics"`
		}{
			countNewLicenses,
			countUnchangedLicenses,
			countUpdatedLicenses,
			countDiffLicenses,
			len(errorLicenseIds),
			countHandledLicenses,
			countAllLoadedLicenses,
			newLicenseIds,
			updatedLicenseIds,
			diffLicenseIds,
			errorLicenseIds,
		}
		customRes = res
	}
	log.AddEntry(job.Info, "finished")
	return scheduler.ExecutionResult{
		Success:   true,
		Log:       log,
		CustomRes: customRes,
	}
}

func getAIfEmpty(source string, a string, b string) string {
	if len(source) > 0 {
		return b
	}
	return a
}

type IRetrieveJsonBytes interface {
	RetrieveJsonBytesFromUrl(url string) ([]byte, error)
}

type retrieveJsonBytes struct {
}

func (retrieveJsonBytes *retrieveJsonBytes) RetrieveJsonBytesFromUrl(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Transport: client_utils.GetTransport(true)}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	file, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	print(file)
	return file, nil
}
