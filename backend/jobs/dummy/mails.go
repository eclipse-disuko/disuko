// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package dummy

import (
	"time"

	"mercedes-benz.ghe.com/foss/disuko/domain/job"
	"mercedes-benz.ghe.com/foss/disuko/domain/label"
	"mercedes-benz.ghe.com/foss/disuko/domain/project"
	"mercedes-benz.ghe.com/foss/disuko/domain/user"
	"mercedes-benz.ghe.com/foss/disuko/helper/mail"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/database"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/labels"
	projectRepo "mercedes-benz.ghe.com/foss/disuko/infra/repository/project"
	userRepo "mercedes-benz.ghe.com/foss/disuko/infra/repository/user"
	"mercedes-benz.ghe.com/foss/disuko/logy"
	"mercedes-benz.ghe.com/foss/disuko/scheduler"
)

type MailJob struct {
	projectRepository projectRepo.IProjectRepository
	labelRepository   labels.ILabelRepository
	mailClient        mail.Client
	userRepository    userRepo.IUsersRepository
}

func InitMailJob(
	projectRepository projectRepo.IProjectRepository,
	labelRepository labels.ILabelRepository,
	mailClient mail.Client,
	userRepository userRepo.IUsersRepository,
) *MailJob {
	return &MailJob{
		projectRepository: projectRepository,
		labelRepository:   labelRepository,
		mailClient:        mailClient,
		userRepository:    userRepository,
	}
}

func (j *MailJob) Execute(rs *logy.RequestSession, info job.Job) scheduler.ExecutionResult {
	var log job.Log
	log.AddEntry(job.Info, "started")

	// Fetch projekt label "dummy"
	dummyLabel := j.labelRepository.FindByNameAndType(rs, label.DUMMY, label.PROJECT)
	if dummyLabel == nil {
		log.AddEntry(job.Error, "label 'dummy' not found")
		return scheduler.ExecutionResult{
			Success: false,
			Log:     log,
		}
	}

	cutoff := time.Now().UTC().AddDate(0, 0, -12).Format(time.RFC3339Nano)
	qc := database.New().SetMatcher(
		database.AndChain(
			database.AttributeMatcher(
				"Deleted",
				database.EQ,
				false,
			),
			database.AttributeMatcher(
				"Created",
				database.LT,
				cutoff,
			),
			database.ArrayElemMatcher(
				"ProjectLabels",
				database.EQ,
				dummyLabel.Key,
			),
		),
	)
	projects := j.projectRepository.Query(rs, qc)
	log.AddEntry(job.Info, "found %d dummy projects for possible notification", len(projects))

	for _, prj := range projects {
		resp := prj.ProjectResponsible()
		if resp == nil {
			continue
		}
		respUser := j.userRepository.FindByUserId(rs, resp.UserId)
		if respUser == nil || (!respUser.Deprovisioned.IsZero() && time.Now().After(respUser.Deprovisioned)) {
			continue
		}

		del := prj.Created.UTC().AddDate(0, 3, 0)
		daysUntil := int(time.Until(del).Hours() / 24)
		if daysUntil == 30 {
			j.sendMail(prj, respUser, &log, 30)
		}
		if daysUntil == 14 {
			j.sendMail(prj, respUser, &log, 14)
		}
	}

	return scheduler.ExecutionResult{
		Success: true,
		Log:     log,
	}
}

func (j *MailJob) sendMail(prj *project.Project, resp *user.User, log *job.Log, days int) {
	mailData := struct {
		Username    string
		ProjectName string
		Days        int
	}{}
	mailData.Username = resp.Forename + " " + resp.Lastname
	mailData.ProjectName = prj.Name
	mailData.Days = days
	templ := "dummyDeletion"
	err := j.mailClient.Send(resp.Email, templ, mailData)
	if err != nil {
		log.AddEntry(job.Error, "Failed to send the email to %s", resp.Email)
	} else {
		log.AddEntry(job.Info, "Email sent successfully to %s", resp.Email)
	}
}
