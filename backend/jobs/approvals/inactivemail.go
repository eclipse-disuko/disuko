// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package approvals

import (
	"time"

	"github.com/eclipse-disuko/disuko/domain/approval"
	"github.com/eclipse-disuko/disuko/domain/job"
	"github.com/eclipse-disuko/disuko/helper/mail"
	"github.com/eclipse-disuko/disuko/infra/repository/approvallist"
	userRepo "github.com/eclipse-disuko/disuko/infra/repository/user"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/eclipse-disuko/disuko/scheduler"
)

const (
	inactiveMailTemplate = "approvalInactiveMail"
	sendMailOnDay        = 18
)

type inactiveMailData struct {
	Username     string
	ProjectName  string
	DeletionDate string
	InactiveDays int
}

type InactiveMail struct {
	approvalListRepo approvallist.IApprovalListRepository
	userRepo         userRepo.IUsersRepository
	mailClient       mail.Client
}

func InitInactiveMail(approvalListRepo approvallist.IApprovalListRepository, userRepo userRepo.IUsersRepository, mailClient mail.Client) *InactiveMail {
	return &InactiveMail{
		approvalListRepo: approvalListRepo,
		userRepo:         userRepo,
		mailClient:       mailClient,
	}
}

func (j *InactiveMail) Execute(rs *logy.RequestSession, info job.Job) scheduler.ExecutionResult {
	var log job.Log
	log.AddEntry(job.Info, "started")

	approvalLists := j.approvalListRepo.FindAll(rs, false)
	for _, list := range approvalLists {
		for _, appr := range list.Approvals {
			if !isOngoing(&appr) {
				continue
			}
			inactive := time.Since(appr.Updated)
			if inactive < sendMailOnDay*24*time.Hour || inactive >= (sendMailOnDay+1)*24*time.Hour {
				continue
			}
			deletionDate := appr.Updated.Add(abortOnDay * 24 * time.Hour).Format("2006-01-02")
			log.AddEntry(job.Info, "approval %s (project %s, type %s) ongoing for %s", appr.Key, list.Key, appr.Type, inactive.Round(time.Second))
			j.notifyRecipients(rs, &appr, list.Key, deletionDate, &log)
		}
	}

	log.AddEntry(job.Info, "finished")
	return scheduler.ExecutionResult{
		Success: true,
		Log:     log,
	}
}

func (j *InactiveMail) notifyRecipients(rs *logy.RequestSession, appr *approval.Approval, projectKey string, deletionDate string, log *job.Log) {
	seen := make(map[string]bool)

	recipients := []string{appr.Creator}
	switch appr.Type {
	case approval.TypeInternal:
		for i := 0; i < 4; i++ {
			recipients = append(recipients, appr.Internal.GetApproverName(approval.Approver(i)))
		}
	case approval.TypePlausibility:
		recipients = append(recipients, appr.Plausibility.Approver)
	}

	for _, userId := range recipients {
		if userId == "" || seen[userId] {
			continue
		}
		seen[userId] = true
		u := j.userRepo.FindByUserId(rs, userId)
		if u == nil || u.Email == "" || !u.Deprovisioned.IsZero() {
			continue
		}
		data := inactiveMailData{
			Username:     u.Forename + " " + u.Lastname,
			ProjectName:  projectKey,
			DeletionDate: deletionDate,
			InactiveDays: sendMailOnDay,
		}
		err := j.mailClient.Send(u.Email, inactiveMailTemplate, data)
		if err != nil {
			log.AddEntry(job.Error, "failed to send mail to %s: %s", u.Email, err)
			logy.Errorf(rs, "failed to send inactive mail to %s: %s", u.Email, err)
		} else {
			log.AddEntry(job.Info, "mail sent to %s for approval %s", u.Email, appr.Key)
		}
	}
}

func isOngoing(a *approval.Approval) bool {
	switch a.Type {
	case approval.TypeInternal:
		return a.Internal.IsActive()
	case approval.TypePlausibility:
		return a.Plausibility.IsActive()
	case approval.TypeExternal:
		return a.External.State == approval.Pending
	}
	return false
}
