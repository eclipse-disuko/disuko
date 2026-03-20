// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package license

import (
	"strings"

	"github.com/eclipse-disuko/disuko/domain/license"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/infra/repository/database"
	"github.com/eclipse-disuko/disuko/logy"
)

type licensesRepositoryStruct struct {
	base.BaseRepositoryWithSoftDelete[*license.License]
}

var deleteCheckFunc = func(requestSession *logy.RequestSession, entity *license.License) {
	if entity.Source != license.CUSTOM {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorLicenseCanOnlyDeleteCustom, (*entity).GetKey()))
	}
}

var createEmptyEntityFunc = func() *license.License {
	return &license.License{}
}

func NewLicenseRepository(requestSession *logy.RequestSession) ILicensesRepository {
	repo := &licensesRepositoryStruct{
		BaseRepositoryWithSoftDelete: base.CreateRepositoryWithSoftDelete[*license.License](
			requestSession,
			LicensesCollectionName,
			createEmptyEntityFunc,
			deleteCheckFunc,
			[]string{
				"text",
				"AuditTrail",
			},
			[][]string{
				{"licenseId"},
				{"Deleted", "licenseId"},
				{"name"},
			},
		),
	}
	repo.PreUpdate = func(requestSession *logy.RequestSession, l *license.License) {
		loaded := repo.FindById(requestSession, l.LicenseId)
		if loaded == nil {
			return
		}
		if loaded.Key == l.Key {
			return
		}
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorLicenseIdInUse, l.LicenseId), "")
	}
	repo.PreSave = func(requestSession *logy.RequestSession, l *license.License) {
		loaded := repo.FindById(requestSession, l.LicenseId)
		if loaded == nil {
			return
		}
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorLicenseIdInUse, l.LicenseId), "")
	}
	return repo
}

func (lr *licensesRepositoryStruct) FindByName(requestSession *logy.RequestSession, name string) *license.License {
	qc := database.New().SetMatcher(
		database.AndChain(
			database.AttributeMatcher(
				"name",
				database.EQ,
				name,
			),
			database.AttributeMatcher(
				"Deleted",
				database.EQ,
				false,
			),
		),
	).SetLimit(0, 1)

	qbLs := lr.Query(requestSession, qc)
	var qbL *license.License
	if qbLs != nil {
		qbL = qbLs[0]
	}
	return qbL
}

func (lr *licensesRepositoryStruct) FindById(requestSession *logy.RequestSession, id string) *license.License {
	qc := database.New().SetMatcher(
		database.AndChain(
			database.AttributeMatcher(
				"licenseId",
				database.EQ,
				id,
			),
			database.AttributeMatcher(
				"Deleted",
				database.EQ,
				false,
			),
		),
	)
	qbLs := lr.Query(requestSession, qc)
	var qbL *license.License
	if qbLs != nil {
		qbL = qbLs[0]
	}
	return qbL
}

func (lr *licensesRepositoryStruct) FindByIdOptimized(requestSession *logy.RequestSession, id string) *license.License {
	qc := database.New().SetMatcher(
		database.AndChain(
			database.AttributeMatcher(
				"licenseId",
				database.EQ,
				id,
			),
			database.AttributeMatcher(
				"Deleted",
				database.EQ,
				false,
			),
		),
	).SetUnset(lr.BaseRepository.OptimizeUnsetAttributes)

	qbLs := lr.Query(requestSession, qc)
	var qbL *license.License
	if qbLs != nil {
		qbL = qbLs[0]
	}
	return qbL
}

func (lr *licensesRepositoryStruct) FindByIdCaseInsensitive(requestSession *logy.RequestSession, id string) *license.License {
	qc := database.New().SetMatcher(
		database.AndChain(
			database.AttributeMatcher(
				"licenseId",
				database.EQI,
				id,
			),
			database.AttributeMatcher(
				"Deleted",
				database.EQ,
				false,
			),
		),
	)
	qbLs := lr.Query(requestSession, qc)
	var qbL *license.License
	if qbLs != nil {
		qbL = qbLs[0]
	}
	return qbL
}

func (lr *licensesRepositoryStruct) FindByObligationKey(requestSession *logy.RequestSession, key string) []*license.License {
	qc := database.New().SetMatcher(
		database.AndChain(
			database.ArrayElemMatcher(
				"meta.obligationsKeyList",
				database.EQ,
				key,
			),
			database.AttributeMatcher(
				"Deleted",
				database.EQ,
				false,
			),
		),
	)

	qbL := lr.Query(requestSession, qc)
	return qbL
}

func (lr *licensesRepositoryStruct) CountByObligationKey(requestSession *logy.RequestSession, key string) int {
	qc := database.New().SetMatcher(
		database.AndChain(
			database.ArrayElemMatcher(
				"meta.obligationsKeyList",
				database.EQ,
				key,
			),
			database.AttributeMatcher(
				"Deleted",
				database.EQ,
				false,
			),
		),
	).SetKeep(
		[]string{
			lr.Database.GetKeyAttribute(),
		},
	)
	qbRes := len(lr.Query(requestSession, qc))
	return qbRes
}

// GetLicenseRefs gathers every license id and family + every alias for a license id and merges them into an object.
// In this object they key is the reference (aliased name or original id) and the value is the actual license id.
func (lr *licensesRepositoryStruct) GetLicenseRefs(requestSession *logy.RequestSession) license.LicenseRefs {
	qc := database.New().SetMatcher(
		database.AttributeMatcher(
			"Deleted",
			database.EQ,
			false,
		),
	).SetKeep([]string{"licenseId", "Aliases", "meta.family", "meta.approvalState"})
	licenses := lr.Query(requestSession, qc)
	res := license.LicenseRefs{}
	for _, l := range licenses {
		res[strings.ToLower(l.LicenseId)] = license.LicenseRef{
			ID:            l.LicenseId,
			Family:        l.Meta.Family,
			ApprovalState: l.Meta.ApprovalState,
		}
		for _, alias := range l.Aliases {
			res[strings.ToLower(alias.LicenseId)] = license.LicenseRef{
				ID:            l.LicenseId,
				Family:        l.Meta.Family,
				ApprovalState: l.Meta.ApprovalState,
			}
		}
	}

	return res
}

func (lr *licensesRepositoryStruct) FindByIds(requestSession *logy.RequestSession, ids []string) []*license.License {
	var licenses []*license.License
	for _, id := range ids {
		qc := database.New().SetMatcher(
			database.AndChain(
				database.AttributeMatcher(
					"licenseId",
					database.EQ,
					id,
				),
				database.AttributeMatcher(
					"Deleted",
					database.EQ,
					false,
				),
			),
		).SetKeep([]string{"licenseId", "meta", "name"})
		qbLs := lr.Query(requestSession, qc)
		if qbLs != nil {
			licenses = append(licenses, qbLs[0])
		}
	}
	return licenses
}
