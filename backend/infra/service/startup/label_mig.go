// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package startup

import (
	"github.com/eclipse-disuko/disuko/domain/label"
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/google/go-cmp/cmp"
)

type labelsLookup map[string]*label.Label

var replacements map[string]string = map[string]string{
	label.EXTERNAL_USERS:  label.CUSTOMER_USERS,
	label.GROUP_USERS:     label.COMPANY_USERS,
	label.ENTITY_USERS:    label.COMPANY_USERS,
	label.ENTITY_TARGET:   label.COMPANY_TARGET,
	label.EXTERNAL_TARGET: label.BP_TARGET,
}

func (startUpHandler *StartUpHandler) migrateLabels(requestSession *logy.RequestSession) {
	logy.Infof(requestSession, "migrateLabels - START")

	exception.TryCatchAndLog(requestSession, func() {
		labelsByKey := make(labelsLookup)
		labelsByName := make(labelsLookup)

		all := startUpHandler.LabelRepository.FindAll(requestSession, false)
		for _, l := range all {
			labelsByKey[l.Key] = l
			labelsByName[l.Name] = l
		}

		projects := startUpHandler.ProjectRepository.FindAllKeysWithDeleted(requestSession)
		logy.Infof(requestSession, "migrateLabels - Found %d projects to process", len(projects))
		failedCount := 0
		for i, projectID := range projects {
			if i%1000 == 0 {
				logy.Infof(requestSession, "migrateLabels - progress: %d/%d", i, len(projects))
			}
			exception.TryCatch(func() {
				pr := startUpHandler.ProjectRepository.FindByKeyWithDeleted(requestSession, projectID, false)
				if pr == nil {
					logy.Warnf(requestSession, "migrateLabels - Project not found: %s", projectID)
					failedCount++
					return
				}

				before := project.PolicyLabelsToAudit(pr.PolicyLabels, labelsByKey)

				replaced := startUpHandler.replaceLabels(requestSession, pr, labelsByKey, labelsByName)
				if !replaced {
					logy.Infof(requestSession, "migrateLabels - nothing replaced in %s (%s)", pr.Name, pr.Key)
					return
				}

				after := project.PolicyLabelsToAudit(pr.PolicyLabels, labelsByKey)
				startUpHandler.AuditLogListRepository.CreateAuditEntryByKey(requestSession, pr.Key, "SYSTEM", message.PolicyLabelsUpdated, cmp.Diff, after, before)

				// TODO: tbd with timestamp?
				startUpHandler.ProjectRepository.UpdateWithoutTimestamp(requestSession, pr)
				logy.Infof(requestSession, "migrateLabels - Updated policy labels for: %s (%s)", pr.Name, pr.Key)
			}, func(e exception.Exception) {
				logy.Infof(requestSession, "migrateLabels - failed on %s", projectID)
				exception.LogException(requestSession, e)
				failedCount++
			})

		}

		logy.Infof(requestSession, "migrateLabels - done. Processed: %d, Failed: %d",
			len(projects), failedCount)
	})

	logy.Infof(requestSession, "migrateLabels - END")
}

func (startUpHandler *StartUpHandler) replaceLabels(requestSession *logy.RequestSession, pr *project.Project, labelsByKey labelsLookup, labelsByName labelsLookup) bool {
	var replaced bool
	for i, lk := range pr.PolicyLabels {
		label, ok := labelsByKey[lk]
		if !ok {
			continue
		}
		repl, ok := replacements[label.Name]
		if !ok {
			continue
		}
		replLabel, ok := labelsByName[repl]
		if !ok {
			logy.Warnf(requestSession, "could not find replacement label %s", repl)
			continue
		}
		logy.Infof(requestSession, "replacing label %s (%s) -> %s (%s)", label.Name, label.Key, replLabel.Name, replLabel.Key)
		pr.PolicyLabels[i] = replLabel.Key
		replaced = true
	}
	return replaced
}
