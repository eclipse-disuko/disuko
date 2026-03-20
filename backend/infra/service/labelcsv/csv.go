// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package labelcsv

import (
	"encoding/csv"
	"net/http"

	"mercedes-benz.ghe.com/foss/disuko/domain/label"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/labels"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/policyrules"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

func CreateCSV(w *http.ResponseWriter, requestSession *logy.RequestSession, policyRepository policyrules.IPolicyRulesRepository, labelRepository labels.ILabelRepository) {
	var csvHeader = []string{
		"Policy name",
		"Created",
		"Updated",
		"Description",
	}

	allLabels := labelRepository.FindAll(requestSession, false)
	for _, l := range allLabels {
		if l.Type == label.POLICY {
			csvHeader = append(csvHeader, l.Name+" (Created: "+l.Created.Format("02.01.2006")+" / Description: "+l.Description+")")
		}
	}
	csvWriter := csv.NewWriter(*w)
	csvWriter.Comma = ';'
	defer csvWriter.Flush()

	if csvErr := csvWriter.Write(csvHeader); csvErr != nil {
		exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCsvGeneration, "policies and labels", "header"), csvErr)
	}

	rules := policyRepository.FindAll(requestSession, false)
	for _, r := range rules {
		for _, labelsSet := range r.LabelSets {
			var row = []string{
				r.Name,
				r.Created.Format("02.01.2006 15:04"),
				r.Updated.Format("02.01.2006 15:04"),
				r.Description,
			}
			for _, l := range allLabels {
				if l.Type == label.POLICY {
					if sliceContains(l.Key, labelsSet) {
						row = append(row, "x")
					} else {
						row = append(row, "")
					}
				}
			}
			if csvErr := csvWriter.Write(row); csvErr != nil {
				exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCsvGeneration, "policies and labels", "data"), csvErr)
			}
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
