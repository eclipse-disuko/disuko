// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package template

import (
	"encoding/csv"
	"net/http"

	"mercedes-benz.ghe.com/foss/disuko/domain/reviewremarks"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/database"
	rt "mercedes-benz.ghe.com/foss/disuko/infra/repository/reviewtemplates"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

func CreateCSV(w *http.ResponseWriter, requestSession *logy.RequestSession, templateRepository rt.IReviewTemplateRepository) {
	var csvHeader = []string{
		"Title",
		"Level",
		"Description",
		"Created",
		"Updated",
	}

	csvWriter := csv.NewWriter(*w)
	csvWriter.Comma = ';'
	defer csvWriter.Flush()

	if err := csvWriter.Write(csvHeader); err != nil {
		exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCsvGeneration, "review remark templates", "header"), err)
	}

	qc := database.New().SetMatcher(
		database.AttributeMatcher(
			"Deleted",
			database.EQ,
			false,
		),
	)
	qbRes := templateRepository.Query(requestSession, qc)

	var reviewTemplates []*reviewremarks.ReviewTemplate
	reviewTemplates = qbRes
	for _, template := range reviewTemplates {
		var lvlInWords string
		switch template.Level {
		case reviewremarks.Green:
			lvlInWords = "information"
		case reviewremarks.Yellow:
			lvlInWords = "investigation"
		case reviewremarks.Red:
			lvlInWords = "problem"
		}

		var row = []string{
			template.Title,
			lvlInWords,
			template.Description,
			template.Created.Format("02.01.2006"),
			template.Updated.Format("02.01.2006"),
		}
		if err := csvWriter.Write(row); err != nil {
			exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorCsvGeneration, "review remark templates", "data"), err)
		}
	}
}
