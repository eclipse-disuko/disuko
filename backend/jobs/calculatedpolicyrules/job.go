// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package calculatedpolicyrules

import (
	"github.com/eclipse-disuko/disuko/domain/job"
	"github.com/eclipse-disuko/disuko/infra/repository/policyrules"
	"github.com/eclipse-disuko/disuko/infra/service/policy"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/eclipse-disuko/disuko/scheduler"
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
	skippedCount := 0

	for _, rule := range allRules {
		if rule == nil || !rule.Calculated || rule.Deleted {
			skippedCount++
			continue
		}

		allow, warn, deny := j.policyService.CalculatePolicyRuleComponents(rs, rule.CalculatedConfig)

		if hasChanged(rule.ComponentsAllow, allow) || hasChanged(rule.ComponentsWarn, warn) || hasChanged(rule.ComponentsDeny, deny) {
			oldAllow := len(rule.ComponentsAllow)
			oldWarn := len(rule.ComponentsWarn)
			oldDeny := len(rule.ComponentsDeny)

			rule.ComponentsAllow = allow
			rule.ComponentsWarn = warn
			rule.ComponentsDeny = deny

			j.policyRulesRepository.Update(rs, rule)

			updatedCount++
			log.AddEntry(job.Info, "updated rule '%s': allow %d->%d, warn %d->%d, deny %d->%d",
				rule.Name, oldAllow, len(allow), oldWarn, len(warn), oldDeny, len(deny))
		} else {
			skippedCount++
		}
	}

	log.AddEntry(job.Info, "finished: %d rules updated, %d skipped", updatedCount, skippedCount)
	return scheduler.ExecutionResult{
		Success: true,
		Log:     log,
	}
}

func hasChanged(old, new []string) bool {
	if len(old) != len(new) {
		return true
	}
	oldMap := make(map[string]bool, len(old))
	for _, v := range old {
		oldMap[v] = true
	}
	for _, v := range new {
		if !oldMap[v] {
			return true
		}
	}
	return false
}
