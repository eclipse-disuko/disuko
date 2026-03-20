// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package spdx_license

import (
	"mercedes-benz.ghe.com/foss/disuko/domain/license"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/logy"
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
