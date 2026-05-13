// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package calculatedpolicyrules

import (
	"github.com/eclipse-disuko/disuko/domain/job"
	"github.com/eclipse-disuko/disuko/helper"
	auditHelper "github.com/eclipse-disuko/disuko/helper/audit"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/infra/repository/policyrules"
	"github.com/eclipse-disuko/disuko/infra/service/policy"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/eclipse-disuko/disuko/scheduler"
	"github.com/google/go-cmp/cmp"
)

type Job struct {
	policyRulesRepository policyrules.IPolicyRulesRepository
	policyService         *policy.Service
}

func Init(policyRulesRepo policyrules.IPolicyRulesRepository, policyService *policy.Service) *Job {
	return &Job{
		policyRulesRepository: policyRulesRepo,
		policyService:         policyService,
	}
}

func (j *Job) Execute(rs *logy.RequestSession, info job.Job) scheduler.ExecutionResult {
	var log job.Log
	log.AddEntry(job.Info, "started")

	allRules := j.policyRulesRepository.FindAll(rs, false)
	updatedCount := 0

	for _, rule := range allRules {
		if !rule.Calculated {
			continue
		}

		allow, warn, deny := j.policyService.CalculatePolicyRuleComponents(rs, rule.CalculatedConfig)

		changed := !helper.EqualsStringSlicesIgnoreOrder(rule.ComponentsAllow, allow) ||
			!helper.EqualsStringSlicesIgnoreOrder(rule.ComponentsWarn, warn) ||
			!helper.EqualsStringSlicesIgnoreOrder(rule.ComponentsDeny, deny)

		if !changed {
			continue
		}

		log.AddEntry(job.Info, "updated rule '%s': allow %d->%d, warn %d->%d, deny %d->%d",
			rule.Name, len(rule.ComponentsAllow), len(allow), len(rule.ComponentsWarn), len(warn), len(rule.ComponentsDeny), len(deny))

		oldAudit := rule.ToAudit(rs, nil)

		rule.ComponentsAllow = allow
		rule.ComponentsWarn = warn
		rule.ComponentsDeny = deny
		newAudit := rule.ToAudit(rs, nil)
		auditHelper.CreateAndAddAuditEntry(&rule.Container, "SYSTEM", message.PolicyRulesUpdated, cmp.Diff, newAudit, oldAudit)
		j.policyRulesRepository.Update(rs, rule)
		updatedCount++
	}

	log.AddEntry(job.Info, "finished: %d rules updated, %d skipped", updatedCount, len(allRules)-updatedCount)
	return scheduler.ExecutionResult{
		Success: true,
		Log:     log,
	}
}
