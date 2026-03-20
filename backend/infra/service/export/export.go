// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package export

import (
	"github.com/eclipse-disuko/disuko/domain/export"
	"github.com/eclipse-disuko/disuko/domain/job"
	"github.com/eclipse-disuko/disuko/domain/label"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/infra/repository/labels"
	"github.com/eclipse-disuko/disuko/infra/repository/license"
	"github.com/eclipse-disuko/disuko/infra/repository/obligation"
	"github.com/eclipse-disuko/disuko/infra/repository/policyrules"
	"github.com/eclipse-disuko/disuko/infra/repository/schema"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/eclipse-disuko/disuko/scheduler"
)

type Service struct {
	LicensesRepository    license.ILicensesRepository
	PolicyRulesRepository policyrules.IPolicyRulesRepository
	ObligationRepository  obligation.IObligationRepository
	LabelRepository       labels.ILabelRepository
	SchemaRepository      schema.ISchemaRepository
	Scheduler             *scheduler.Scheduler
}

func (exportService *Service) ExportLicenseKnowledgeBase(requestSession *logy.RequestSession) *export.ExportLicenseKnowledgeBaseDto {
	result := &export.ExportLicenseKnowledgeBaseDto{
		PolicyLabels: exportService.LabelRepository.FindAllByType(requestSession, label.POLICY),
		PolicyRules:  exportService.PolicyRulesRepository.FindAll(requestSession, false),
		Obligations:  exportService.ObligationRepository.FindAll(requestSession, false),
		Licenses:     exportService.LicensesRepository.FindAll(requestSession, false),
	}
	for _, license := range result.Licenses {
		license.Rev = ""
	}
	for _, policyLabel := range result.PolicyLabels {
		policyLabel.Rev = ""
	}
	for _, policyRule := range result.PolicyRules {
		policyRule.Rev = ""
	}
	for _, obligation := range result.Obligations {
		obligation.Rev = ""
	}
	return result
}

func (exportService *Service) ImportLicenseKnowledgeBase(requestSession *logy.RequestSession, data *export.ExportLicenseKnowledgeBaseDto) {
	// handle labels
	for _, policyLabelFromFile := range data.PolicyLabels {
		policyLabelFromDb := exportService.LabelRepository.FindByNameAndType(requestSession, policyLabelFromFile.Name, policyLabelFromFile.Type)
		if policyLabelFromDb == nil {
			exportService.LabelRepository.Save(requestSession, policyLabelFromFile)
		} else {
			for _, rule := range data.PolicyRules {
				newLabelSets := make([][]string, 0)
				for _, labelSet := range rule.LabelSets {
					newLabelSet := make([]string, 0)
					for _, label := range labelSet {
						if label == policyLabelFromFile.Key {
							label = policyLabelFromDb.Key
						}
						newLabelSet = append(newLabelSet, label)
					}
					newLabelSets = append(newLabelSets, newLabelSet)
				}
				rule.LabelSets = newLabelSets
			}
		}
	}

	err := exportService.Scheduler.ExecuteJobManual(requestSession, job.LabelLoadDb)
	if err != nil {
		exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorStartingJob), err)
	}

	// delete old data
	exportService.deleteAll(requestSession, license.LicensesCollectionName)
	exportService.deleteAll(requestSession, policyrules.PolicyRulesCollectionName)
	exportService.deleteAll(requestSession, obligation.ObligationCollectionName)
	// add new data
	exportService.ObligationRepository.SaveList(requestSession, data.Obligations, true)
	exportService.LicensesRepository.SaveList(requestSession, data.Licenses, true)
	exportService.PolicyRulesRepository.SaveList(requestSession, data.PolicyRules, true)
}

func (exportService *Service) ExportSchemaKnowledgeBase(requestSession *logy.RequestSession) *export.ExportSchemaKnowledgeBaseDto {
	result := &export.ExportSchemaKnowledgeBaseDto{
		SchemaLabels: exportService.LabelRepository.FindAllByType(requestSession, label.SCHEMA),
		Schemas:      exportService.SchemaRepository.FindAll(requestSession, false),
	}
	for _, schemaLabel := range result.SchemaLabels {
		schemaLabel.Rev = ""
	}
	for _, spdxSchema := range result.Schemas {
		spdxSchema.Rev = ""
	}
	return result
}

func (exportService *Service) ImportSchemaKnowledgeBase(requestSession *logy.RequestSession, data *export.ExportSchemaKnowledgeBaseDto) {
	// handle labels
	for _, schemaLabelFromFile := range data.SchemaLabels {
		schemaLabelFromDb := exportService.LabelRepository.FindByNameAndType(requestSession, schemaLabelFromFile.Name, schemaLabelFromFile.Type)
		if schemaLabelFromDb == nil {
			exportService.LabelRepository.Save(requestSession, schemaLabelFromFile)
		} else {
			for _, spdxSchema := range data.Schemas {
				if spdxSchema.Label == schemaLabelFromFile.Key {
					spdxSchema.Label = schemaLabelFromDb.Key
				}
			}
		}
	}

	err := exportService.Scheduler.ExecuteJobManual(requestSession, job.LabelLoadDb)
	if err != nil {
		exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorStartingJob), err)
	}

	// delete old data
	exportService.deleteAll(requestSession, schema.SpdxSchemaCollectionName)

	// add new data
	exportService.SchemaRepository.SaveList(requestSession, data.Schemas, true)
}

func (exportService *Service) deleteAll(requestSession *logy.RequestSession, collectionName string) {
	db := base.NewDatabase()
	db.Init(requestSession, collectionName, nil)
	db.Truncate()
}
