// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package licenseremarks

import (
	"sort"
	"strings"

	"github.com/eclipse-disuko/disuko/domain/license"
	"github.com/eclipse-disuko/disuko/domain/obligation"
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/domain/project/components"
	licenseRepo "github.com/eclipse-disuko/disuko/infra/repository/license"
	obligationRepo "github.com/eclipse-disuko/disuko/infra/repository/obligation"
	"github.com/eclipse-disuko/disuko/logy"
)

type licenseremarks struct {
	containsWarnings bool
	containsAlarms   bool
	license          *license.License
	obligations      map[string]*obligation.Obligation
	affected         map[string]components.ComponentResult
}

func CreateQualityLicenseRemarks(requestSession *logy.RequestSession, licenseRepository licenseRepo.ILicensesRepository, obligationRepository obligationRepo.IObligationRepository, evalRes *components.EvaluationResult) []*project.QualityLicenseRemarks2 {
	if evalRes == nil {
		return nil
	}
	result := make(map[string]*licenseremarks, 0)
	for _, compRes := range evalRes.Results {
		for _, license := range compRes.Component.GetLicensesEffective().List {
			if !license.Known {
				continue
			}
			if _, found := result[license.ReferencedLicense]; !found {
				result[license.ReferencedLicense] = &licenseremarks{
					license:     nil,
					obligations: make(map[string]*obligation.Obligation),
					affected:    map[string]components.ComponentResult{},
				}
			}
			result[license.ReferencedLicense].affected[compRes.Component.SpdxId] = compRes
		}
	}

	licenseIds := make([]string, 0)
	for id := range result {
		licenseIds = append(licenseIds, id)
	}
	obligationsMap := make(map[string]bool, 0)
	obligationsKeys := make([]string, 0)
	licensesSlice := licenseRepository.FindByIds(requestSession, licenseIds)
	for _, license := range licensesSlice {
		result[license.LicenseId].license = license
		for _, key := range license.Meta.ObligationsKeyList {
			if _, found := obligationsMap[key]; !found {
				obligationsMap[key] = true
				obligationsKeys = append(obligationsKeys, key)
			}
		}
	}
	obligationsSlice := obligationRepository.FindByKeys(requestSession, obligationsKeys, false)
	obligations := make(map[string]*obligation.Obligation)
	for _, obligation := range obligationsSlice {
		obligations[obligation.Key] = obligation
	}

	for _, lremarks := range result {
		for _, obligationKey := range lremarks.license.Meta.ObligationsKeyList {
			obl, found := obligations[obligationKey]
			if !found {
				continue
			}
			lremarks.obligations[obligationKey] = obl
			if strings.ToLower(string(obl.WarnLevel)) == obligation.Warning {
				lremarks.containsWarnings = true
			}
			if strings.ToLower(string(obl.WarnLevel)) == obligation.Alarm {
				lremarks.containsAlarms = true
			}
		}
	}

	return toDto(result)
}

func toDto(in map[string]*licenseremarks) []*project.QualityLicenseRemarks2 {
	res := make([]*project.QualityLicenseRemarks2, 0)
	for _, lremarks := range in {
		if len(lremarks.obligations) == 0 {
			continue
		}
		l := &project.QualityLicenseRemarks2{
			License:     lremarks.license.LicenseId,
			Warnings:    lremarks.containsWarnings,
			Alarms:      lremarks.containsAlarms,
			Obligations: []*obligation.ObligationDto{},
			Affected:    []*project.AffectedComponent{},
		}
		for _, obl := range lremarks.obligations {
			l.Obligations = append(l.Obligations, obligation.ToDto(obl))
		}
		for _, affected := range lremarks.affected {
			l.Affected = append(l.Affected, &project.AffectedComponent{
				SpdxId:           affected.Component.SpdxId,
				Name:             affected.Component.Name,
				Version:          affected.Component.Version,
				PolicyRuleStatus: affected.Status,
			})
		}
		sort.Slice(l.Affected, func(i, j int) bool {
			return strings.ToLower(l.Affected[i].Name) < strings.ToLower(l.Affected[j].Name)
		})
		res = append(res, l)
	}
	getWeightFn := func(lr *project.QualityLicenseRemarks2) int {
		if lr.Alarms {
			return 2
		} else if lr.Warnings {
			return 1
		}
		return 0
	}
	sort.Slice(res, func(i, j int) bool {
		return getWeightFn(res[i]) > getWeightFn(res[j])
	})
	return res
}
