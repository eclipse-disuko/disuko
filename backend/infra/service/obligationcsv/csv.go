// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package obligationcsv

import (
	"encoding/csv"
	"net/http"
	"strings"

	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/infra/repository/database"
	"github.com/eclipse-disuko/disuko/infra/repository/obligation"

	license2 "github.com/eclipse-disuko/disuko/domain/license"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/infra/repository/license"
	"github.com/eclipse-disuko/disuko/logy"
)

func CreateCSV(w *http.ResponseWriter, requestSession *logy.RequestSession, obligationsRepository obligation.IObligationRepository, licenseRepository license.ILicensesRepository) {
	var csvHeader = []string{
		"License name",
		"License ID",
		"License Family",
		"License Type",
		"Evaluation State",
		"Risk Checked",
		"Review State",
		"Review Date",
		"Legal Comments",
		"License Chart Text",
		"License Chart Status",
		"Source",
		"SPDX Status",
		"Created",
		"Updated",
		"LicenseURL",
		"SourceURL",
		"Aliases",
	}

	classifcations := obligationsRepository.FindAllSortedByName(requestSession)
	for _, c := range classifcations {
		csvHeader = append(csvHeader, c.Name+" ("+string(c.Type)+"/"+string(c.WarnLevel)+")")
	}

	csvWriter := csv.NewWriter(*w)
	csvWriter.Comma = ';'
	defer csvWriter.Flush()

	if err := csvWriter.Write(csvHeader); err != nil {
		exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCsvGeneration, "licenses and classifications", "header"), err)
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

	var licenses []*license2.License
	licenses = qbRes

	for _, l := range licenses {
		aliases := ""
		for _, a := range l.Aliases {
			aliases += a.LicenseId + ", "
		}
		aliases = strings.TrimSuffix(aliases, ", ")

		chartText := ""
		if l.Meta.Evaluation != "" {
			chartText = "available"
		}
		chartStatus := "no"
		if l.Meta.IsLicenseChart {
			chartStatus = "yes"
		}
		comments := ""
		if l.Meta.LegalComments != "" {
			comments = "available"
		}
		status := "ok"
		if l.IsDeprecatedLicenseId {
			status = "deprecated"
		}

		riskCheckedStr := ""
		audits := l.AuditTrail
		for _, audit := range audits {
			if strings.Contains(audit.MetaJSON, "-: check") {
				riskCheckedStr = audit.Created.Format("02.01.2006")
			}
		}

		var row = []string{
			l.Name,
			l.LicenseId,
			string(l.Meta.Family),
			string(l.Meta.LicenseType),
			string(l.Meta.ApprovalState),
			riskCheckedStr,
			string(l.Meta.ReviewState),
			l.Meta.ReviewDateStr,
			comments,
			chartText,
			chartStatus,
			string(l.Source),
			status,
			l.Created.Format("02.01.2006"),
			l.Updated.Format("02.01.2006"),
			l.Meta.LicenseUrl,
			l.Meta.SourceUrl,
			aliases,
		}
		for _, c := range classifcations {
			if sliceContains(c.Key, l.Meta.ObligationsKeyList) {
				row = append(row, "x")
			} else {
				row = append(row, "")
			}
		}
		if err := csvWriter.Write(row); err != nil {
			exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCsvGeneration, "licenses and classifications", "data"), err)
		}
	}
}

func sliceContains(needle string, haystack []string) bool {
	for _, tmp := range haystack {
		if tmp == needle {
			return true
		}
	}
	return false
}
