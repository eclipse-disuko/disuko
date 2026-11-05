// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package startup

import (
	"testing"
	"time"

	"github.com/eclipse-disuko/disuko/domain"
	rr "github.com/eclipse-disuko/disuko/domain/reviewremarks"
	i18nRepo "github.com/eclipse-disuko/disuko/infra/repository/i18n"
	remarkRepo "github.com/eclipse-disuko/disuko/infra/repository/reviewremarks"
	templateRepo "github.com/eclipse-disuko/disuko/infra/repository/reviewtemplates"
	"github.com/eclipse-disuko/disuko/logy"
)

type informationRemarkRepository struct {
	remarkRepo.IReviewRemarksRepository
	items   []*rr.ReviewRemarks
	updates int
}

func (repo *informationRemarkRepository) FindAllWithDeleted(*logy.RequestSession, bool) []*rr.ReviewRemarks {
	return repo.items
}

func (repo *informationRemarkRepository) UpdateWithoutTimestamp(*logy.RequestSession, *rr.ReviewRemarks) {
	repo.updates++
}

type informationTemplateRepository struct {
	templateRepo.IReviewTemplateRepository
	items   []*rr.ReviewTemplate
	updates int
}

func (repo *informationTemplateRepository) FindAllWithDeleted(*logy.RequestSession, bool) []*rr.ReviewTemplate {
	return repo.items
}

func (repo *informationTemplateRepository) UpdateWithoutTimestamp(*logy.RequestSession, *rr.ReviewTemplate) {
	repo.updates++
}

type remarkLevelTranslationRepository struct {
	i18nRepo.II18nRepository
	values  map[string]string
	updates int
}

func (repo *remarkLevelTranslationRepository) GetTranslation(_ *logy.RequestSession, locale, key string) (string, bool) {
	value, found := repo.values[locale+"/"+key]
	return value, found
}

func (repo *remarkLevelTranslationRepository) SetTranslation(_ *logy.RequestSession, locale, key, value, description, updatedBy string) {
	repo.values[locale+"/"+key] = value
	repo.updates++
}

func TestMigrateReviewRemarkLevelTranslations(t *testing.T) {
	repo := &remarkLevelTranslationRepository{values: map[string]string{
		"en/REMARK_LEVEL_GREEN":  "Information",
		"de/REMARK_LEVEL_GREEN":  "Information",
		"en/RR_CHART_ACCEPTABLE": "Information",
		"de/RR_CHART_ACCEPTABLE": "Custom label",
	}}
	handler := StartUpHandler{I18nRepository: repo}
	session := logy.NewRequestSession()
	handler.migrateReviewRemarkLevelTranslations(session)
	handler.migrateReviewRemarkLevelTranslations(session)
	for key, want := range map[string]string{
		"en/REMARK_LEVEL_GREEN":  "Closed",
		"de/REMARK_LEVEL_GREEN":  "Geschlossen",
		"en/RR_CHART_ACCEPTABLE": "Closed",
		"de/RR_CHART_ACCEPTABLE": "Custom label",
	} {
		if repo.values[key] != want {
			t.Errorf("translation %s = %q, want %q", key, repo.values[key], want)
		}
	}
	if repo.updates != 3 {
		t.Errorf("updates = %d, want 3", repo.updates)
	}
}

func TestMigrateInformationReviewRemarkLevels(t *testing.T) {
	for _, status := range []rr.Status{rr.Open, rr.InProgress, rr.Closed, rr.Cancelled, rr.UnsetStatus} {
		for _, level := range []rr.Level{rr.Green, rr.Blue, rr.Yellow, rr.Red, rr.Gray, rr.UnsetLevel} {
			t.Run(string(status)+"/"+string(level), func(t *testing.T) {
				entity := &rr.ReviewRemarks{RootEntity: domain.NewRootEntity(), SoftDelete: domain.SoftDelete{Deleted: true}}
				entity.Remarks = []*rr.Remark{nil, {Level: level, Status: status}}
				before := entity.Updated
				repo := &informationRemarkRepository{items: []*rr.ReviewRemarks{entity}}
				handler := StartUpHandler{ReviewRemarkRepository: repo}
				session := logy.NewRequestSession()
				handler.migrateInformationReviewRemarkLevels(session)
				wantLevel, wantUpdates := level, 0
				if level == rr.Green && (status == rr.Open || status == rr.InProgress) {
					wantLevel, wantUpdates = rr.Blue, 1
				}
				if entity.Remarks[1].Level != wantLevel || entity.Remarks[1].Status != status {
					t.Fatalf("got level %q, status %q; want %q, %q", entity.Remarks[1].Level, entity.Remarks[1].Status, wantLevel, status)
				}
				handler.migrateInformationReviewRemarkLevels(session)
				if repo.updates != wantUpdates || !entity.Updated.Equal(before) {
					t.Fatalf("updates = %d, want %d; timestamp preserved = %t", repo.updates, wantUpdates, entity.Updated.Equal(before))
				}
			})
		}
	}
}

func TestMigrateInformationReviewTemplateLevels(t *testing.T) {
	for _, level := range []rr.Level{rr.Green, rr.Blue, rr.Yellow, rr.Red, rr.Gray, rr.UnsetLevel} {
		t.Run(string(level), func(t *testing.T) {
			entity := &rr.ReviewTemplate{RootEntity: domain.NewRootEntity(), SoftDelete: domain.SoftDelete{Deleted: true}, Level: level}
			entity.Updated = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
			before := entity.Updated
			repo := &informationTemplateRepository{items: []*rr.ReviewTemplate{entity}}
			handler := StartUpHandler{ReviewTemplateRepository: repo}
			session := logy.NewRequestSession()
			handler.migrateInformationReviewTemplateLevels(session)
			wantLevel, wantUpdates := level, 0
			if level == rr.Green {
				wantLevel, wantUpdates = rr.Blue, 1
			}
			if entity.Level != wantLevel {
				t.Fatalf("level = %q, want %q", entity.Level, wantLevel)
			}
			handler.migrateInformationReviewTemplateLevels(session)
			if repo.updates != wantUpdates || !entity.Updated.Equal(before) {
				t.Fatalf("updates = %d, want %d; timestamp preserved = %t", repo.updates, wantUpdates, entity.Updated.Equal(before))
			}
		})
	}
}
