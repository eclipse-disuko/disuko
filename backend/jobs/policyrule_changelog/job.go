// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package policyrule_changelog

import (
	"regexp"
	"sort"
	"time"

	"github.com/eclipse-disuko/disuko/domain/audit"
	changeloglist2 "github.com/eclipse-disuko/disuko/domain/changeloglist"
	"github.com/eclipse-disuko/disuko/domain/job"
	"github.com/eclipse-disuko/disuko/domain/license"
	"github.com/eclipse-disuko/disuko/infra/repository/changeloglist"
	"github.com/eclipse-disuko/disuko/infra/repository/changelogs"
	licenseRepo "github.com/eclipse-disuko/disuko/infra/repository/license"
	"github.com/eclipse-disuko/disuko/infra/repository/policyrules"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/eclipse-disuko/disuko/scheduler"
)

type Job struct {
	changeLogListRepository changeloglist.IChangeLogListRepository
	changeLogsRepository    changelogs.IChangeLogsRepository
	policyRuleRepository    policyrules.IPolicyRulesRepository
	licenseRepository       licenseRepo.ILicensesRepository
}

func Init(cllr changeloglist.IChangeLogListRepository, clr changelogs.IChangeLogsRepository, prr policyrules.IPolicyRulesRepository, lr licenseRepo.ILicensesRepository) *Job {
	return &Job{
		changeLogListRepository: cllr,
		changeLogsRepository:    clr,
		policyRuleRepository:    prr,
		licenseRepository:       lr,
	}
}

const (
	Added   = "Added"
	Removed = "Removed"
)

func (j *Job) Execute(rs *logy.RequestSession, info job.Job) scheduler.ExecutionResult {
	var log job.Log

	log.AddEntry(job.Info, "started")

	var (
		policyRules = j.policyRuleRepository.FindAll(rs, false)
		licenses    = j.licenseRepository.FindAll(rs, false)
		failed      = 0
		ok          = 0
		result      = changeloglist2.ChangeLogGenerationResultDto{
			ReqID:           rs.ReqID,
			ChangeLogsCount: 0,
			Errors:          make([]string, 0),
		}
	)
	licensesMap := make(map[string]string)
	for _, l := range licenses {
		licensesMap[l.LicenseId] = l.Name
	}

	allChangeLoglistsToSave := make([]*changeloglist2.ChangeLogList, 0)
	allChangeLoglistsToUpdate := make([]*changeloglist2.ChangeLogList, 0)
	allChangeLogs := make([]*changeloglist2.ChangeLog, 0)
	for _, pr := range policyRules {
		policyRuleChangeLogs := runChecks(pr, licensesMap, info.Updated)
		if len(policyRuleChangeLogs) == 0 {
			continue
		}

		changeLogs := make([]*changeloglist2.ChangeLog, 0, len(policyRuleChangeLogs))
		for _, policyRuleChangeLog := range policyRuleChangeLogs {
			cl, err := policyRuleChangeLog.ToChangeLog()
			if err == nil {
				changeLogs = append(changeLogs, cl)
				ok++
			} else {
				log.AddEntry(job.Error, "could not marshall change log content %s", err.Error())
				result.Errors = append(result.Errors, "could not marshall change log content"+err.Error())
				failed++
			}
		}
		log.AddEntry(job.Info, "saved change logs for policy rule %s", pr.Name)
		for k, val := range policyRuleChangeLogs {
			log.AddEntry(job.Info, "key / value = %d / %v", k, val)
		}

		changeLogList := j.changeLogListRepository.FindByKey(rs, pr.Key, false)
		if changeLogList == nil {
			changeLogList = &changeloglist2.ChangeLogList{}
			changeLogList.Key = pr.Key
			changeLogList.ChangeLogs = append(changeLogList.ChangeLogs, changeLogs...)
			allChangeLoglistsToSave = append(allChangeLoglistsToSave, changeLogList)
		} else {
			changeLogList.ChangeLogs = append(changeLogList.ChangeLogs, changeLogs...)
			allChangeLoglistsToUpdate = append(allChangeLoglistsToUpdate, changeLogList)
		}
		allChangeLogs = append(allChangeLogs, changeLogs...)
	}

	if len(allChangeLoglistsToSave) > 0 {
		j.changeLogListRepository.SaveList(rs, allChangeLoglistsToSave, false)
	}
	if len(allChangeLoglistsToUpdate) > 0 {
		j.changeLogListRepository.UpdateList(rs, allChangeLoglistsToUpdate)
	}
	if len(allChangeLogs) > 0 {
		j.changeLogsRepository.SaveList(rs, allChangeLogs, false)
	}

	result.ChangeLogsCount = ok
	log.AddEntry(job.Info, "finished: ok %d / failed %d", ok, failed)
	return scheduler.ExecutionResult{
		Success:   !(failed > 0),
		Log:       log,
		CustomRes: result,
	}
}

func runChecks(pr *license.PolicyRules, licensesMap map[string]string, lastRun time.Time) []*changeloglist2.PolicyRuleChangeLog {
	allPolicyRuleChangeLogs := make([]*changeloglist2.PolicyRuleChangeLog, 0)
	for _, auditTrail := range pr.AuditTrail {
		if auditTrail.Created.Before(lastRun) {
			// ignore older auditTrails
			continue
		}
		currentAuditChangeLogs := make([]*changeloglist2.PolicyRuleChangeLog, 0)
		for _, check := range policyRuleChangeLogChecks {
			changeLogs := check(pr, licensesMap, auditTrail)
			if len(changeLogs) > 0 {
				currentAuditChangeLogs = append(currentAuditChangeLogs, changeLogs...)
			}
		}
		reordered := reorder(currentAuditChangeLogs)
		allPolicyRuleChangeLogs = append(allPolicyRuleChangeLogs, reordered...)
	}
	return allPolicyRuleChangeLogs
}

func reorder(entries []*changeloglist2.PolicyRuleChangeLog) []*changeloglist2.PolicyRuleChangeLog {
	// Sort entries by LicenseId first, and then by Change, prioritizing "Added" over "Removed"
	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].LicenseId == entries[j].LicenseId {
			// "Added" should come before "Removed" for the same LicenseId
			return entries[i].Change == Added && entries[j].Change == Removed
		}
		return entries[i].LicenseId < entries[j].LicenseId
	})
	return entries
}

func createNewPolicyRuleChangeLog(pr *license.PolicyRules, when time.Time, licenseName, licenseId, policyStatus string, change string) *changeloglist2.PolicyRuleChangeLog {
	cl := changeloglist2.NewPolicyRuleChangeLog()
	cl.When = when
	cl.RefKey = pr.Key
	cl.RefName = pr.Name

	cl.LicenseName = licenseName
	cl.LicenseId = licenseId
	cl.PolicyStatus = policyStatus
	cl.Change = change
	return cl
}

var policyRuleChangeLogChecks = []func(pr *license.PolicyRules, licensesMap map[string]string, auditLog *audit.Audit) []*changeloglist2.PolicyRuleChangeLog{
	checkSingleChanges,
	checkMultipleChanges,
}

func checkSingleChanges(pr *license.PolicyRules, licensesMap map[string]string, auditLog *audit.Audit) []*changeloglist2.PolicyRuleChangeLog {
	pattern := `([+-])[\s\xA0]*\tComponents(Allow|Deny|Warn):[\s\xA0]*\[\]string\{\"([-.\w\d]+)\"\}`

	r, _ := regexp.Compile(pattern)
	matches := r.FindAllStringSubmatch(auditLog.MetaJSON, -1)

	var changeLogs []*changeloglist2.PolicyRuleChangeLog

	for _, match := range matches {
		if len(match) < 4 {
			continue
		}

		sign := match[1]
		policyStatus := match[2]
		licenseId := match[3]

		if licenseId == "" {
			continue
		}

		change := signToChange(sign)

		changeLogs = append(changeLogs, createNewPolicyRuleChangeLog(pr, auditLog.Created, licensesMap[licenseId], licenseId, policyStatus, change))
	}
	return changeLogs
}

func checkMultipleChanges(pr *license.PolicyRules, licensesMap map[string]string, auditLog *audit.Audit) []*changeloglist2.PolicyRuleChangeLog {
	pattern := `\tComponents(Allow|Deny|Warn):[\s\xA0]*\[\]string\{\n([^}]*)`

	r, _ := regexp.Compile(pattern)
	matches := r.FindAllStringSubmatch(auditLog.MetaJSON, -1)

	var changeLogs []*changeloglist2.PolicyRuleChangeLog

	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		policyStatus := match[1]

		innerPattern := `([+-])[\s\xA0]*\t\t\"([^"]+)\"`
		rInner, _ := regexp.Compile(innerPattern)
		values := rInner.FindAllStringSubmatch(match[2], -1)
		var innerChangeLogs []*changeloglist2.PolicyRuleChangeLog
		for _, value := range values {
			if len(value) < 3 {
				continue
			}

			sign := value[1]
			licenseId := value[2]

			change := signToChange(sign)

			// Here can be same License with '+' and '-', meaning just changing the position in the array.
			// Such entries must not be taken into account, it is not a change
			if !findAndRemoveConflicting(&innerChangeLogs, licenseId, change) {
				innerChangeLogs = append(innerChangeLogs, createNewPolicyRuleChangeLog(pr, auditLog.Created, licensesMap[licenseId], licenseId, policyStatus, change))
			}
		}
		changeLogs = append(changeLogs, innerChangeLogs...)
	}
	return changeLogs
}

func signToChange(sign string) string {
	var change string
	if sign == "+" {
		change = Added
	} else if sign == "-" {
		change = Removed
	}
	return change
}

func findAndRemoveConflicting(changeLogs *[]*changeloglist2.PolicyRuleChangeLog, licenseId string, change string) bool {
	for i, changeLog := range *changeLogs {
		if changeLog.LicenseId == licenseId && changeLog.Change != change {
			// Remove the entry, which is not a change and just position changing in the array
			*changeLogs = append((*changeLogs)[:i], (*changeLogs)[i+1:]...)
			return true
		}
	}
	return false
}
