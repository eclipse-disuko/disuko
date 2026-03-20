// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package userstats

import (
	"encoding/json"
	"strings"

	"mercedes-benz.ghe.com/foss/disuko/domain"
	"mercedes-benz.ghe.com/foss/disuko/domain/audit"
	"mercedes-benz.ghe.com/foss/disuko/domain/job"
	"mercedes-benz.ghe.com/foss/disuko/domain/oauth"
	domainUser "mercedes-benz.ghe.com/foss/disuko/domain/user"
	"mercedes-benz.ghe.com/foss/disuko/domain/userstats"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/labels"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/license"
	newsboxRepo "mercedes-benz.ghe.com/foss/disuko/infra/repository/newsbox"

	"mercedes-benz.ghe.com/foss/disuko/infra/repository/obligation"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/policyrules"
	projectRepo "mercedes-benz.ghe.com/foss/disuko/infra/repository/project"
	rt "mercedes-benz.ghe.com/foss/disuko/infra/repository/reviewtemplates"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/schema"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/user"
	userstatsRepo "mercedes-benz.ghe.com/foss/disuko/infra/repository/userstats"
	"mercedes-benz.ghe.com/foss/disuko/logy"
	"mercedes-benz.ghe.com/foss/disuko/scheduler"
)

type Config struct {
	Username       string
	Rights         *oauth.AccessAndRolesRights
	AdminRequest   bool
	CollectionName string
}

type Job struct {
	projectRepo              projectRepo.IProjectRepository
	licenseRepo              license.ILicensesRepository
	policyrulesRepo          policyrules.IPolicyRulesRepository
	labelRepo                labels.ILabelRepository
	schemaRepository         schema.ISchemaRepository
	obligationRepository     obligation.IObligationRepository
	userRepository           user.IUsersRepository
	reviewTemplateRepository rt.IReviewTemplateRepository
	userstatsRepository      userstatsRepo.IUserStatsRepository
	newsboxRepository        newsboxRepo.IRepo
}

func Init(projectRepo projectRepo.IProjectRepository, licenseRepo license.ILicensesRepository, labelRepo labels.ILabelRepository, policyrulesRepo policyrules.IPolicyRulesRepository, schemaRepository schema.ISchemaRepository, obligationRepository obligation.IObligationRepository, userRepository user.IUsersRepository, reviewTemplateRepository rt.IReviewTemplateRepository, userstatsRepository userstatsRepo.IUserStatsRepository, newsbox newsboxRepo.IRepo) *Job {
	return &Job{
		projectRepo:              projectRepo,
		licenseRepo:              licenseRepo,
		policyrulesRepo:          policyrulesRepo,
		labelRepo:                labelRepo,
		schemaRepository:         schemaRepository,
		obligationRepository:     obligationRepository,
		userRepository:           userRepository,
		reviewTemplateRepository: reviewTemplateRepository,
		userstatsRepository:      userstatsRepository,
		newsboxRepository:        newsbox,
	}
}

func (j *Job) Execute(rs *logy.RequestSession, info job.Job) scheduler.ExecutionResult {
	var log job.Log
	log.AddEntry(job.Info, "started")
	var config Config

	if err := json.Unmarshal([]byte(info.Config), &config); err != nil {
		log.AddEntry(job.Error, "decoding config failed: %s", err)
		return scheduler.ExecutionResult{
			Success: false,
			Log:     log,
		}
	}

	var entry *userstats.UserStatus
	entryExists := true
	entry = j.userstatsRepository.FindByUserId(rs, config.Username)
	if entry == nil {
		entry = &userstats.UserStatus{
			RootEntity: domain.NewRootEntity(),
			Container:  audit.Container{},
			User:       strings.ToUpper(config.Username),
		}
		entryExists = false
	}
	isUser := !config.AdminRequest
	if isUser {
		entry.ActiveJobCount = j.getActiveJobCount(rs, config)   // user dependent
		entry.DisclosureCount = j.getDisclosureCount(rs, config) // user dependent
		entry.ProjectCount = j.getProjectCount(rs, config)       // user dependent
		if entryExists {
			j.userstatsRepository.Update(rs, entry)
		} else {
			j.userstatsRepository.Save(rs, entry)
		}
		config.AdminRequest = true // overwrite to get other counts for the admin stats
	}

	if config.AdminRequest {
		entry = j.userstatsRepository.FindByUserId(rs, userstats.SystemStatsUser)
		entryExists = entry != nil
		if entry == nil {
			entry = &userstats.UserStatus{
				RootEntity:          domain.NewRootEntity(),
				Container:           audit.Container{},
				User:                strings.ToUpper(userstats.SystemStatsUser),
				LicenseCount:        -1,
				SchemaCount:         -1,
				ReviewTemplateCount: -1,
				ObligationCount:     -1,
				LabelCount:          -1,
				PolicyRuleCount:     -1,
				UserCount:           -1,
			}
		}
		entry.ActiveJobCount = -1
		entry.DisclosureCount = -1
		entry.ProjectCount = -1

		if config.CollectionName == license.LicensesCollectionName || entry.LicenseCount == -1 {
			entry.LicenseCount = j.getLicenseCount(rs, config)
		}
		if config.CollectionName == schema.SpdxSchemaCollectionName || entry.SchemaCount == -1 {
			entry.SchemaCount = j.getSchemaCount(rs, config)
		}
		if config.CollectionName == rt.ReviewTemplateCollectionName || entry.ReviewTemplateCount == -1 {
			entry.ReviewTemplateCount = j.getReviewTemplateCount(rs, config)
		}
		if config.CollectionName == obligation.ObligationCollectionName || entry.ObligationCount == -1 {
			entry.ObligationCount = j.getObligationCount(rs, config)
		}
		if config.CollectionName == labels.LabelCollectionName || entry.LabelCount == -1 {
			entry.LabelCount = j.getLabelCount(rs, config)
		}
		if config.CollectionName == policyrules.PolicyRulesCollectionName || entry.PolicyRuleCount == -1 {
			entry.PolicyRuleCount = j.getPolicyRuleCount(rs, config)
		}
		if config.CollectionName == user.UsersCollectionName || entry.UserCount == -1 {
			entry.UserCount = j.getUserCount(rs, config)
		}
		if config.CollectionName == projectRepo.ProjectCollectionName || entry.ProjectCount == -1 {
			entry.ProjectCount = j.getProjectCount(rs, config)
		}
		if entryExists {
			j.userstatsRepository.Update(rs, entry)
		} else {
			j.userstatsRepository.Save(rs, entry)
		}
	}

	log.AddEntry(job.Info, "finished")
	return scheduler.ExecutionResult{
		Success: true,
		Log:     log,
	}
}

func (j *Job) getObligationCount(rs *logy.RequestSession, config Config) int {

	if !config.Rights.AllowObligation.Read {
		return 0
	}

	qbRes := j.obligationRepository.CountAll(rs)
	return qbRes
}

func (j *Job) getSchemaCount(rs *logy.RequestSession, config Config) int {
	if !config.Rights.AllowSchema.Read {
		return 0
	}

	return j.schemaRepository.CountAll(rs)
}

func (j *Job) getLabelCount(rs *logy.RequestSession, config Config) int {
	if !config.Rights.AllowPolicy.Read {
		return 0
	}
	return j.labelRepo.CountAll(rs)
}

func (j *Job) getPolicyRuleCount(rs *logy.RequestSession, config Config) int {
	if !config.Rights.AllowPolicy.Read {
		return 0
	}

	qbRes := j.policyrulesRepo.CountAll(rs)
	return qbRes
}

func (j *Job) getLicenseCount(rs *logy.RequestSession, config Config) int {
	if !config.Rights.AllowLicense.Read {
		return 0
	}
	return j.licenseRepo.CountAll(rs)
}

func (j *Job) getProjectCount(rs *logy.RequestSession, config Config) int {
	if config.AdminRequest {
		qbRes := j.projectRepo.CountAll(rs)
		return qbRes
		// not needed for admin view
		// err = roles.FilterProjectsWithoutAccess(requestSession, w, r, &projects)
	} else {

		return j.projectRepo.CountForUser(rs, config.Username)
	}

}

func (j *Job) getDisclosureCount(rs *logy.RequestSession, config Config) int {
	if config.AdminRequest {
		qbRes := j.projectRepo.CountAllGroups(rs)
		return qbRes
		// not needed for admin view
		// err = roles.FilterProjectsWithoutAccess(requestSession, w, r, &projects)
	} else {
		return j.projectRepo.CountGroupsForUser(rs, config.Username)
	}
}

func (j *Job) getUserCount(rs *logy.RequestSession, config Config) int {
	if !(config.Rights.AllowUsers.Create && config.Rights.AllowUsers.Read && config.Rights.AllowUsers.Update && config.Rights.AllowUsers.Delete) {
		return 0
	}

	qbRes := j.userRepository.CountAll(rs)
	return qbRes
}

func (j *Job) getReviewTemplateCount(rs *logy.RequestSession, config Config) int {
	if !config.Rights.AllowObligation.Read && !config.Rights.AllowLicense.Read && !config.Rights.AllowPolicy.Read {
		return 0
	}
	return j.reviewTemplateRepository.CountAll(rs)
}

func (j *Job) getActiveJobCount(rs *logy.RequestSession, config Config) int {

	u := j.userRepository.FindByUserId(rs, config.Username)
	if u == nil {
		return 0
	}
	res := 0

	for _, t := range u.Tasks {
		if t.Status == domainUser.TaskActive {
			res++
		}
	}
	return res
}
