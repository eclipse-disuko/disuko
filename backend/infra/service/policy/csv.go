// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package policy

import (
	"encoding/csv"
	"net/http"
	"strings"

	"golang.org/x/text/encoding/unicode"
	lic "mercedes-benz.ghe.com/foss/disuko/domain/license"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/database"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/license"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/policyrules"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

func CreateCSV(w *http.ResponseWriter, requestSession *logy.RequestSession, policyRepository policyrules.IPolicyRulesRepository, licenseRepository license.ILicensesRepository) {
	var csvHeader = []string{
		"License name",
		"License ID",
		"License Family",
		"License Type",
		"Evaluation State",
		"Review State",
		"Review Date",
		"Source",
		"Created",
		"Updated",
		"LicenseURL",
		"Aliases",
	}

	policies := policyRepository.FindAll(requestSession, false)
	for _, p := range policies {
		csvHeader = append(csvHeader, p.Name+" (Created: "+p.Created.Format("02.01.2006 15:04")+" / Updated: "+p.Updated.Format("02.01.2006 15:04")+" / Description: "+p.Description+")")
	}

	convWriter := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewEncoder().Writer(*w)
	csvWriter := csv.NewWriter(convWriter)
	csvWriter.Comma = ';'

	if err := csvWriter.Write(csvHeader); err != nil {
		exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCsvGeneration, "licenses and policies", "header"), err)
	}

	qc := database.New().SetMatcher(
		database.AttributeMatcher(
			"Deleted",
			database.EQ,
			false,
		),
	).SetSort(database.SortConfig{
		database.SortAttribute{
			Name:  "name",
			Order: database.ASC,
		},
	})
	qbRes := licenseRepository.Query(requestSession, qc)

	var licenses = qbRes
	for _, l := range licenses {

		var row = []string{
			l.Name,
			l.LicenseId,
			string(l.Meta.Family),
			string(l.Meta.LicenseType),
			string(l.Meta.ApprovalState),
			string(l.Meta.ReviewState),
			l.Meta.ReviewDateStr,
			string(l.Source),
			l.Created.Format("02.01.2006"),
			l.Updated.Format("02.01.2006"),
			l.Meta.LicenseUrl,
			joinLicenseAliases(l.Aliases),
		}
		for _, p := range policies {
			if sliceContains(l.LicenseId, p.ComponentsAllow) {
				row = append(row, "allowed")
			} else if sliceContains(l.LicenseId, p.ComponentsDeny) {
				row = append(row, "denied")
			} else if sliceContains(l.LicenseId, p.ComponentsWarn) {
				row = append(row, "warned")
			} else {
				row = append(row, " ")
			}
		}
		if err := csvWriter.Write(row); err != nil {
			exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCsvGeneration, "licenses and policies", "data"), err)
		}
	}
	csvWriter.Flush()
}

func CreateRuleCSV(w *http.ResponseWriter, requestSession *logy.RequestSession, policyRepository policyrules.IPolicyRulesRepository, licenseRepository license.ILicensesRepository, id string) {
	var csvHeader = []string{
		"License name",
		"License ID",
		"License URL",
		"Policy Rule",
		"License Status",
		"License Aliases",
	}

	convWriter := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewEncoder().Writer(*w)
	csvWriter := csv.NewWriter(convWriter)
	csvWriter.Comma = ';'

	if err := csvWriter.Write(csvHeader); err != nil {
		exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCsvGeneration, "licenses and policies", "header"), err)
	}
	pol := policyRepository.FindByKey(requestSession, id, false)
	if pol == nil {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorDbNotFound), id+" not found in DB")
	}

	for _, id := range pol.ComponentsAllow {
		lic := licenseRepository.FindById(requestSession, id)
		if lic == nil {
			continue
		}
		var row = []string{
			lic.Name,
			lic.LicenseId,
			lic.Meta.LicenseUrl,
			pol.Name,
			"allowed",
			joinLicenseAliases(lic.Aliases),
		}
		if err := csvWriter.Write(row); err != nil {
			exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCsvGeneration, "licenses and policies", "data"), err)
		}
	}

	for _, id := range pol.ComponentsDeny {
		lic := licenseRepository.FindById(requestSession, id)
		if lic == nil {
			continue
		}

		var row = []string{
			lic.Name,
			lic.LicenseId,
			lic.Meta.LicenseUrl,
			pol.Name,
			"denied",
			joinLicenseAliases(lic.Aliases),
		}
		if err := csvWriter.Write(row); err != nil {
			exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCsvGeneration, "licenses and policies", "data"), err)
		}
	}

	for _, id := range pol.ComponentsWarn {
		lic := licenseRepository.FindById(requestSession, id)
		if lic == nil {
			continue
		}
		var row = []string{
			lic.Name,
			lic.LicenseId,
			lic.Meta.LicenseUrl,
			pol.Name,
			"warned",
		}
		if err := csvWriter.Write(row); err != nil {
			exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCsvGeneration, "licenses and policies", "data"), err)
		}
	}
	csvWriter.Flush()

}

func sliceContains(needle string, haystack []string) bool {
	for _, tmp := range haystack {
		if tmp == needle {
			return true
		}
	}
	return false
}

// joinLicenseAliases returns a comma-separated string of non-empty LicenseIds from the given aliases slice.
func joinLicenseAliases(aliases []lic.Alias) string {
	var aliasList []string
	for _, a := range aliases {
		if a.LicenseId != "" {
			aliasList = append(aliasList, a.LicenseId)
		}
	}
	return strings.Join(aliasList, ", ")
}
