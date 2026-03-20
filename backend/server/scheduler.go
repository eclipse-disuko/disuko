// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"context"

	"mercedes-benz.ghe.com/foss/disuko/jobs/labels"
	"mercedes-benz.ghe.com/foss/disuko/jobs/userstats"

	"mercedes-benz.ghe.com/foss/disuko/domain/job"
	"mercedes-benz.ghe.com/foss/disuko/jobs/analytics"
	"mercedes-benz.ghe.com/foss/disuko/jobs/announcement"
	"mercedes-benz.ghe.com/foss/disuko/jobs/departments"
	"mercedes-benz.ghe.com/foss/disuko/jobs/deprovisioning"
	"mercedes-benz.ghe.com/foss/disuko/jobs/dummy"
	"mercedes-benz.ghe.com/foss/disuko/jobs/fossdd"
	"mercedes-benz.ghe.com/foss/disuko/jobs/licenserefresh"
	"mercedes-benz.ghe.com/foss/disuko/jobs/notification"
	"mercedes-benz.ghe.com/foss/disuko/jobs/policyrule_changelog"
	"mercedes-benz.ghe.com/foss/disuko/jobs/report"
	"mercedes-benz.ghe.com/foss/disuko/jobs/termsofuse"
	"mercedes-benz.ghe.com/foss/disuko/logy"

	"mercedes-benz.ghe.com/foss/disuko/scheduler"
)

func (s *Server) setupScheduling(ctx context.Context, rs *logy.RequestSession) {
	s.scheduler = scheduler.Init(rs, s.repos.job, s.services.lock)

	depRefresh := departments.Init(s.repos.department, s.connectors.department, s.scheduler)
	s.scheduler.AddJobCb(job.DepartmentRefresh, depRefresh)

	depLoad := departments.InitLoadDb(s.repos.department)
	s.scheduler.AddJobCb(job.DepartmentLoadDb, depLoad)

	licenseRefresh := licenserefresh.Init(s.repos.licenses, s.repos.obligation, s.repos.spdxLicense, nil)
	s.scheduler.AddJobCb(job.LicenseRefresh, licenseRefresh)

	licenseAnnounce := announcement.Init(s.repos.announcements, s.repos.licenses)
	s.scheduler.AddJobCb(job.LicenseAnnouncements, licenseAnnounce)

	policyRuleChangeLog := policyrule_changelog.Init(s.repos.changeLogList, s.repos.changeLogs, s.repos.policyRules, s.repos.licenses)
	s.scheduler.AddJobCb(job.PolicyRuleChangeLogs, policyRuleChangeLog)

	terms := termsofuse.Init(s.repos.user)
	s.scheduler.AddJobCb(job.TermOfUseUpdate, terms)

	ana := analytics.Init(s.services.analytics)
	s.scheduler.AddJobCb(job.AnalyticsRebuild, ana)

	not := notification.Init(s.repos.dpConfig)
	s.scheduler.AddJobCb(job.Notification, not)

	dep := deprovisioning.Init(s.repos.user, s.connectors.userrole)
	s.scheduler.AddJobCb(job.Deprovisioning, dep)

	fossdd := fossdd.Init(
		s.repos.approvalList,
		s.repos.project,
		s.repos.policyRules,
		s.repos.user,
		&s.handlers.project,
		&s.services.fossdd,
		&s.services.projectLabelService,
	)
	s.scheduler.AddJobCb(job.FOSSDDGen, fossdd)
	userStatsOTJob := userstats.Init(
		s.repos.project, s.repos.licenses, s.repos.label,
		s.repos.policyRules, s.repos.schema,
		s.repos.obligation, s.repos.user,
		s.repos.reviewTemplate,
		s.repos.userstats,
		s.repos.newsbox,
	)
	s.scheduler.AddJobCb(job.CalculateUserStats, userStatsOTJob)

	rep := report.Init(
		s.repos.project,
		s.repos.user,
		s.repos.department,
		s.repos.label,
		s.repos.sbomList,
		s.repos.approvalList,
		s.repos.obligation,
		s.repos.policyRules,
		s.repos.licenseRules,
		s.services.spdx,
		s.repos.customid,
		&s.services.projectLabelService,
		s.repos.policyDecisions,
	)
	s.scheduler.AddJobCb(job.Report, rep)

	dummyDeletion := dummy.InitDeletionJob(
		s.repos.project,
		s.repos.label,
		s.repos.sbomList,
		s.repos.reviewRemarks,
		s.repos.approvalList,
		s.repos.user,
		s.repos.auditLogList,
	)
	s.scheduler.AddJobCb(job.DummyProjectDeletion, dummyDeletion)

	dummyMail := dummy.InitMailJob(
		s.repos.project,
		s.repos.label,
		s.mailClient,
		s.repos.user,
	)
	s.scheduler.AddJobCb(job.DummyMail, dummyMail)

	label := labels.Init(
		s.repos.label,
	)
	s.scheduler.AddJobCb(job.LabelLoadDb, label)

	go s.scheduler.Start(ctx)
	s.handlers.job.Scheduler = s.scheduler // todo ensure scheduler is set also found in observer.go
	s.handlers.project.Scheduler = s.scheduler
	s.handlers.count.Scheduler = s.scheduler
	s.handlers.label.Scheduler = s.scheduler
	s.services.export.Scheduler = s.scheduler
}
