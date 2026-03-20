// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package spdx_license

import (
	"github.com/eclipse-disuko/disuko/domain/license"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
)

type spdxLicensesRepositoryStruct struct {
	base.BaseRepositoryWithHardDelete[*license.License]
}

func NewSpdxLicenseRepository(requestSession *logy.RequestSession) ISpdxLicensesRepository {
	return &spdxLicensesRepositoryStruct{
		BaseRepositoryWithHardDelete: base.CreateRepositoryWithHardDelete[*license.License](
			requestSession,
			SpdxLicensesCollectionName,
			func() *license.License {
				return &license.License{}
			},
			nil,
			"",
			nil,
			nil),
	}
}
