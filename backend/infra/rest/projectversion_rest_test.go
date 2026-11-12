// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"testing"

	"github.com/eclipse-disuko/disuko/domain/project/components"
	"github.com/eclipse-disuko/disuko/domain/reviewremarks"
	remarkRepo "github.com/eclipse-disuko/disuko/infra/repository/reviewremarks"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/stretchr/testify/assert"
)

type remarkStatsRepository struct {
	remarkRepo.IReviewRemarksRepository
	remarks *reviewremarks.ReviewRemarks
}

func (repo *remarkStatsRepository) FindByKey(*logy.RequestSession, string, bool) *reviewremarks.ReviewRemarks {
	return repo.remarks
}

func TestCalculateReviewRemarkStatsSeparatesInformationAndClosed(t *testing.T) {
	repo := &remarkStatsRepository{remarks: &reviewremarks.ReviewRemarks{Remarks: []*reviewremarks.Remark{
		{Level: reviewremarks.Blue, Status: reviewremarks.Open},
		{Level: reviewremarks.Blue, Status: reviewremarks.InProgress},
		{Level: reviewremarks.Green, Status: reviewremarks.Closed},
		{Level: reviewremarks.Yellow, Status: reviewremarks.Open},
		{Level: reviewremarks.Red, Status: reviewremarks.InProgress},
		{Level: reviewremarks.Gray, Status: reviewremarks.Cancelled},
	}}}
	handler := ProjectHandler{ReviewRemarksRepository: repo}
	stats := handler.calculateReviewRemarkStats(logy.NewRequestSession(), "version")
	assert.Equal(t, components.ReviewRemarkStats{
		Total: 6, Information: 2, Acceptable: 1, AcceptableAfterChanges: 1, NotAcceptable: 1, Cancelled: 1,
	}, stats)
	repo.remarks = nil
	assert.Equal(t, components.ReviewRemarkStats{}, handler.calculateReviewRemarkStats(logy.NewRequestSession(), "version"))
}

func TestShouldValidateFilename(t *testing.T) {
	//given
	type test struct {
		file  string
		valid bool
	}
	tests := []test{
		{file: "test.zip", valid: true},
		{file: "test.gz", valid: true},
		{file: "test@#$%^&()_.gz", valid: true},
		{file: "a/a.zip", valid: false},
		{file: "a<.zip", valid: false},
		{file: "a>.zip", valid: false},
		{file: "a'.zip", valid: false},
		{file: "a\".zip", valid: false},
		{file: "a\\.zip", valid: false},
		{file: "a*.zip", valid: false},
		{file: "a?.zip", valid: false},
		{file: "a|.zip", valid: false},
	}

	//when
	for _, tc := range tests {
		got := isFilenameValid(tc.file)
		//then
		assert.Equal(t, tc.valid, got, "misfunctioning validation of filename for "+tc.file)
	}

}
