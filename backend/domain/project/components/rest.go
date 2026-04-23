// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package components

import (
	"strings"

	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/domain/license"
	"github.com/eclipse-disuko/disuko/domain/licenserules"
	obligation2 "github.com/eclipse-disuko/disuko/domain/obligation"
	"github.com/eclipse-disuko/disuko/domain/policydecisions"
	"github.com/eclipse-disuko/disuko/helper/message"
	licRepo "github.com/eclipse-disuko/disuko/infra/repository/license"
	"github.com/eclipse-disuko/disuko/infra/repository/obligation"
	"github.com/eclipse-disuko/disuko/logy"
)

type PolicyRuleStatusDto struct {
	Key                         string           `json:"key"`
	Name                        string           `json:"name"`
	LicenseMatched              string           `json:"licenseMatched"`
	Type                        license.ListType `json:"type"`
	Used                        bool             `json:"used"`
	Description                 string           `json:"description"`
	IsDecisionMade              bool             `json:"isDecisionMade"`
	CanMakeWarnedDecision       bool             `json:"canMakeWarnedDecision"`
	CanMakeDeniedDecision       bool             `json:"canMakeDeniedDecision"`
	DeniedDecisionDeniedReason  string           `json:"deniedDecisionDeniedReason"`
	LicenseRecommendationWeight *float64         `json:"licenseRecommendationWeight"`
}
type UnmatchedLicenseDto struct {
	OrigName       string `json:"orig"`
	ReferencedName string `json:"referenced"`
	Known          bool   `json:"known"`
}

type ComponentInfoDto struct {
	SpdxId             string                 `json:"spdxId"`
	Name               string                 `json:"name"`
	Version            string                 `json:"version"`
	LicenseEffective   string                 `json:"licenseEffective"`
	License            string                 `json:"license"`
	LicenseDeclared    string                 `json:"licenseDeclared"`
	LicenseComments    string                 `json:"licenseComments"`
	WorstFamily        string                 `json:"worstFamily"`
	CopyrightText      string                 `json:"copyrightText"`
	Description        string                 `json:"description"`
	DownloadLocation   string                 `json:"downloadLocation"`
	Type               ComponentType          `json:"type"`
	Modified           bool                   `json:"modified"`
	Questioned         bool                   `json:"questioned"`
	Unasserted         bool                   `json:"unasserted"`
	PolicyRuleStatus   []*PolicyRuleStatusDto `json:"policyRuleStatus"`
	UnmatchedLicenses  []*UnmatchedLicenseDto `json:"unmatchedLicenses"`
	LicenseApplied     LicenseAppliedType     `json:"licenseApplied"`
	PURL               string                 `json:"purl"`
	PrStatus           string                 `json:"prStatus"`
	UsedPolicyRule     string                 `json:"usedPolicyRule"`
	CanChooseLicense   bool                   `json:"canChooseLicense"`
	ChoiceDeniedReason string                 `json:"choiceDeniedReason"`

	LicenseRuleApplied *licenserules.LicenseRuleSlimDto `json:"licenseRuleApplied"`

	PolicyDecisionsApplied     []*policydecisions.PolicyDecisionSlimDto `json:"policyDecisionsApplied"`
	PolicyDecisionDeniedReason string                                   `json:"policyDecisionDeniedReason"`

	LicenseRecommended *string `json:"licenseRecommended"`
}

type ComponentInfoSlimDto struct {
	SpdxId            string `json:"spdxId"`
	Name              string `json:"name"`
	Version           string `json:"version"`
	LicenseExpression string `json:"licenseExpression"`
}

func (ci ComponentInfo) ToComponentInfoSlimDto() *ComponentInfoSlimDto {
	return &ComponentInfoSlimDto{
		SpdxId:            ci.SpdxId,
		Name:              ci.Name,
		Version:           ci.Version,
		LicenseExpression: ci.EffectiveLicensesString(),
	}
}

type ComponentsInfoResponse struct {
	ComponentInfo                  []ComponentInfoDto `json:"componentInfo"`
	ComponentStats                 ComponentStats     `json:"componentStats"`
	BulkPolicyDecisionDeniedReason string             `json:"bulkPolicyDecisionDeniedReason"`
}

func (entity *ComponentResult) ToComponentInfoDto(
	isResponsible bool,
	policyDecisionDeniedReason string,
	isAllowDeniedPolicyDecision bool,
	licensesRepository licRepo.ILicensesRepository,
	rs *logy.RequestSession,
	cache *recommendationCache,
) *ComponentInfoDto {
	status, rule := entity.GetUsedPolicyRule()

	var (
		deniedReason string
		canChoose    = entity.Component.LicenseRuleApplied == nil && entity.Component.GetLicensesEffective().Op == OR
	)
	if canChoose {
		if !isResponsible {
			deniedReason = message.ChoiceDeniedResp
		} else if len(entity.Component.GetLicensesEffective().List) > 4 {
			deniedReason = message.ChoiceDeniedMassive
		}
	}

	if entity.Component.Version == "" {
		if policyDecisionDeniedReason == "" {
			policyDecisionDeniedReason = message.DecisionDeniedComponentVersionNotSet
		}
		if deniedReason == "" {
			deniedReason = message.DecisionDeniedComponentVersionNotSet
		}
	}

	var licenseRecommended *string
	policyStatusDtos := ToPolicyStatusDto(entity.Status, isAllowDeniedPolicyDecision)
	if !conf.IsProdEnv() && deniedReason == "" && canChoose {
		licenseRecommended = recommendLicense(policyStatusDtos, licensesRepository, rs, cache)
	}

	return &ComponentInfoDto{
		SpdxId:                     entity.Component.SpdxId,
		Name:                       entity.Component.Name,
		Version:                    entity.Component.Version,
		LicenseEffective:           entity.Component.EffectiveLicensesString(),
		License:                    entity.Component.License,
		LicenseDeclared:            entity.Component.LicenseDeclared,
		LicenseComments:            entity.Component.LicenseComments,
		WorstFamily:                string(entity.Component.WorstFamily()),
		CopyrightText:              entity.Component.CopyrightText,
		Description:                entity.Component.Description,
		DownloadLocation:           entity.Component.DownloadLocation,
		Type:                       entity.Component.Type,
		LicenseApplied:             entity.Component.GetLicenseAppliedType(),
		Modified:                   entity.Component.Modified,
		Questioned:                 entity.Questioned,
		Unasserted:                 entity.Unasserted,
		PolicyRuleStatus:           policyStatusDtos,
		UnmatchedLicenses:          ToUnmatchedDto(entity.Unmatched),
		PrStatus:                   status,
		UsedPolicyRule:             rule,
		PURL:                       entity.Component.PURL,
		CanChooseLicense:           canChoose,
		ChoiceDeniedReason:         deniedReason,
		LicenseRuleApplied:         entity.Component.LicenseRuleApplied.ToSlimDto(),
		PolicyDecisionsApplied:     policydecisions.ToSlimDtos(entity.Component.PolicyDecisionsApplied),
		PolicyDecisionDeniedReason: policyDecisionDeniedReason,
		LicenseRecommended:         licenseRecommended,
	}
}

type recommendationCache struct {
	licenses    map[string]*license.License
	obligations map[string]*obligation2.Obligation
}

func recommendLicense(
	policyStatusDtos []*PolicyRuleStatusDto,
	licensesRepository licRepo.ILicensesRepository,
	rs *logy.RequestSession,
	cache *recommendationCache,
) *string {
	if len(policyStatusDtos) == 0 {
		return nil
	}

	var allows []*PolicyRuleStatusDto
	var warns []*PolicyRuleStatusDto

	for _, ps := range policyStatusDtos {
		if ps == nil {
			continue
		}

		switch ps.Type {
		case license.ALLOW:
			allows = append(allows, ps)
		case license.WARN:
			warns = append(warns, ps)
		}
	}

	if len(allows) == 1 {
		return &allows[0].LicenseMatched
	}
	if len(allows) > 1 {
		return recommendByClassificationWeight(allows, licensesRepository, rs, cache)
	}

	if len(warns) == 1 {
		return &warns[0].LicenseMatched
	}
	if len(warns) > 1 {
		return recommendByClassificationWeight(warns, licensesRepository, rs, cache)
	}

	return nil
}

func recommendByClassificationWeight(
	policyStatusDtos []*PolicyRuleStatusDto,
	licensesRepository licRepo.ILicensesRepository,
	rs *logy.RequestSession,
	cache *recommendationCache,
) *string {
	if licensesRepository == nil || rs == nil || cache == nil {
		return nil
	}

	licensesWeightMap := calculateLicenseWeights(policyStatusDtos, licensesRepository, rs, cache)
	applyLicenseWeights(policyStatusDtos, licensesWeightMap)

	return findRecommendedLicense(licensesWeightMap)
}

func findRecommendedLicense(licensesWeightMap map[string]float64) *string {
	var recommendedLicense string
	var bestScore float64
	found := false
	bestScoreCount := 0

	for licenseId, score := range licensesWeightMap {
		if !found || score < bestScore {
			recommendedLicense = licenseId
			bestScore = score
			found = true
			bestScoreCount = 1
		} else if score == bestScore {
			bestScoreCount++
		}
	}

	if !found || bestScoreCount > 1 {
		return nil
	}

	return &recommendedLicense
}

func calculateLicenseWeights(policyStatusDtos []*PolicyRuleStatusDto, licensesRepository licRepo.ILicensesRepository, rs *logy.RequestSession, cache *recommendationCache) map[string]float64 {
	licensesWeightMap := make(map[string]float64)
	processedLicenses := make(map[string]struct{})

	for _, ps := range policyStatusDtos {
		licenseId := ps.LicenseMatched
		if _, alreadyProcessed := processedLicenses[licenseId]; alreadyProcessed {
			continue
		}
		processedLicenses[licenseId] = struct{}{}

		l := getLicenseFromCacheOrRepo(cache, licensesRepository, rs, licenseId)
		if l == nil {
			continue
		}

		for _, obligationKey := range l.Meta.ObligationsKeyList {
			o := cache.obligations[obligationKey]
			if o == nil {
				continue
			}

			switch strings.ToLower(string(o.WarnLevel)) {
			case obligation2.Information:
				licensesWeightMap[licenseId] += obligation2.InfoWeight
			case obligation2.Warning:
				licensesWeightMap[licenseId] += obligation2.WarnWeight
			case obligation2.Alarm:
				licensesWeightMap[licenseId] += obligation2.AlarmWeight
			}
		}

		if _, exists := licensesWeightMap[licenseId]; !exists {
			licensesWeightMap[licenseId] = 0
		}
	}
	return licensesWeightMap
}

func applyLicenseWeights(
	policyStatusDtos []*PolicyRuleStatusDto,
	licensesWeightMap map[string]float64,
) {
	for _, ps := range policyStatusDtos {
		ps.LicenseRecommendationWeight = new(licensesWeightMap[ps.LicenseMatched])
	}
}

func getLicenseFromCacheOrRepo(cache *recommendationCache, licensesRepository licRepo.ILicensesRepository, rs *logy.RequestSession, licenseId string) *license.License {
	if l, ok := cache.licenses[licenseId]; ok {
		return l
	}
	l := licensesRepository.FindById(rs, licenseId)
	if l != nil {
		cache.licenses[licenseId] = l
	}
	return l
}

func (entity *EvaluationResult) ToComponentInfoDtos(isResponsible bool,
	policyDecisionDeniedReason string,
	isAllowDeniedPolicyDecision bool,
	obligationProvider obligation.IObligationRepository,
	licensesRepository licRepo.ILicensesRepository,
	rs *logy.RequestSession,
) []ComponentInfoDto {
	dtos := make([]ComponentInfoDto, 0)

	cache := &recommendationCache{
		licenses:    make(map[string]*license.License),
		obligations: make(map[string]*obligation2.Obligation),
	}

	if obligationProvider != nil && rs != nil {
		allObligations := obligationProvider.FindAll(rs, false)
		for _, o := range allObligations {
			if o == nil || o.Key == "" {
				continue
			}
			cache.obligations[o.Key] = o
		}
	}

	for _, compRes := range entity.Results {
		dtos = append(dtos, *compRes.ToComponentInfoDto(isResponsible, policyDecisionDeniedReason, isAllowDeniedPolicyDecision, licensesRepository, rs, cache))
	}
	return dtos
}

func ToPolicyStatusDto(status []*PolicyRuleStatus, isAllowDeniedPolicyDecision bool) []*PolicyRuleStatusDto {
	dtos := make([]*PolicyRuleStatusDto, 0)
	for _, s := range status {
		dtos = append(dtos, &PolicyRuleStatusDto{
			Key:                        s.Key,
			Name:                       s.Name,
			LicenseMatched:             s.LicenseMatched,
			Type:                       s.Type,
			Used:                       s.Used,
			Description:                s.Description,
			IsDecisionMade:             s.IsDecisionMade,
			CanMakeWarnedDecision:      s.CanMakeWarnedDecision,
			CanMakeDeniedDecision:      s.CanMakeDeniedDecision && isAllowDeniedPolicyDecision,
			DeniedDecisionDeniedReason: s.DeniedDecisionDeniedReason,
		})
	}
	return dtos
}

func ToUnmatchedDto(unmatched []*UnmatchedLicense) []*UnmatchedLicenseDto {
	dtos := make([]*UnmatchedLicenseDto, 0)
	for _, u := range unmatched {
		dtos = append(dtos, &UnmatchedLicenseDto{
			OrigName:       u.OrigName,
			ReferencedName: u.ReferencedName,
			Known:          u.Known,
		})
	}
	return dtos
}
