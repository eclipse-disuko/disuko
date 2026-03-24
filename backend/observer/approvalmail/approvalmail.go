// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package approvalmail

import (
	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/domain/approval"
	"github.com/eclipse-disuko/disuko/domain/user"
	"github.com/eclipse-disuko/disuko/helper/mail"
	"github.com/eclipse-disuko/disuko/infra/repository/project"
	userRepo "github.com/eclipse-disuko/disuko/infra/repository/user"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/eclipse-disuko/disuko/observermngmt"
)

type ApprovalMail struct {
	mailClient  mail.Client
	userRepo    userRepo.IUsersRepository
	projectRepo project.IProjectRepository
}

func Init(mailClient mail.Client, userRepo userRepo.IUsersRepository, projectRepo project.IProjectRepository) *ApprovalMail {
	return &ApprovalMail{
		mailClient:  mailClient,
		userRepo:    userRepo,
		projectRepo: projectRepo,
	}
}

func (a *ApprovalMail) RegisterHandlers() {
	observermngmt.RegisterHandler(observermngmt.ApprovalTaskCreated, a.OnApprovalTask)
	observermngmt.RegisterHandler(observermngmt.ApprovalFinalized, a.OnApprovalFinalized)
}

func (a *ApprovalMail) OnApprovalFinalized(eventId observermngmt.EventId, arg interface{}) {
	data, ok := arg.(observermngmt.ApprovalData)
	if !ok {
		return
	}
	a.sendFinalizedMail(data)
}

func (a *ApprovalMail) OnApprovalTask(eventId observermngmt.EventId, arg interface{}) {
	data, ok := arg.(observermngmt.ApprovalTaskData)
	if !ok {
		return
	}
	a.sendTaskMail(data)
}

func (a *ApprovalMail) sendFinalizedMail(data observermngmt.ApprovalData) {
	type mailData struct {
		Username        string
		Link            string
		State           string
		StateDE         string
		Requestor       string
		Comment         string
		Reviewer        string
		DelegatedTo     string
		ReviewerComment string
		IsGroup         bool
		GroupName       string
		GroupLink       string
		Versions        []struct {
			Num               int
			ProjectName       string
			ProjectLink       string
			VersionName       string
			VersionLink       string
			ReviewRemarksLink string
		}
	}
	md := mailData{}
	creatorUser := a.userRepo.FindByUserId(data.RequestSession, data.Approval.Creator)
	md.Requestor = creatorUser.Forename + " " + creatorUser.Lastname
	pr := a.projectRepo.FindByKey(data.RequestSession, data.Approval.ProjectGuid, false)
	md.IsGroup = pr.IsGroup
	if pr.IsGroup {
		md.GroupName = pr.Name
		md.GroupLink = conf.Config.Server.DisukoHost + "/#/dashboard/groups/" + data.Approval.ProjectGuid
		md.Link = md.GroupLink + "/approvals"
		for i, a := range data.Approval.Info.Projects {
			prLink := conf.Config.Server.DisukoHost + "/#/dashboard/projects/" + a.ProjectKey
			versionLink := prLink + "/versions/" + a.ApprovableSPDX.VersionKey
			reviewRemarksLink := versionLink + "/sbomQuality/" + a.ApprovableSPDX.SpdxKey + "/reviewRemarks"
			md.Versions = append(md.Versions, struct {
				Num               int
				ProjectName       string
				ProjectLink       string
				VersionName       string
				VersionLink       string
				ReviewRemarksLink string
			}{
				Num:               i + 1,
				ProjectName:       a.ProjectName,
				ProjectLink:       prLink,
				VersionName:       a.ApprovableSPDX.VersionName,
				VersionLink:       versionLink,
				ReviewRemarksLink: reviewRemarksLink,
			})
		}
	} else {
		approvable := data.Approval.Info.Projects[0]
		prLink := conf.Config.Server.DisukoHost + "/#/dashboard/projects/" + data.Approval.ProjectGuid
		md.Link = prLink + "/approvals"
		versionLink := prLink + "/versions/" + approvable.ApprovableSPDX.VersionKey
		reviewRemarksLink := versionLink + "/sbomQuality/" + approvable.ApprovableSPDX.SpdxKey + "/reviewRemarks"
		md.Versions = append(md.Versions, struct {
			Num               int
			ProjectName       string
			ProjectLink       string
			VersionName       string
			VersionLink       string
			ReviewRemarksLink string
		}{
			ProjectName:       pr.Name,
			ProjectLink:       prLink,
			VersionName:       approvable.ApprovableSPDX.VersionName,
			VersionLink:       versionLink,
			ReviewRemarksLink: reviewRemarksLink,
		})

	}
	md.Comment = data.Approval.Comment

	switch data.Approval.Type {
	case approval.TypeInternal:
		if data.Approval.Internal.CustomerDone() {
			md.State = "Approved"
			md.StateDE = "Genehmigt"
		} else if !data.Approval.Internal.Aborted {
			md.State = "Declined"
			md.StateDE = "Nicht genehmigt"
		} else {
			md.State = "Aborted"
			md.StateDE = "Abgebrochen"
		}
		md.Username = creatorUser.Forename + " " + creatorUser.Lastname
		a.sendMail(data.RequestSession, creatorUser, creatorUser.Email, "approvalFinalized", md)
		users := []string{
			data.Approval.Internal.GetApproverName(approval.Supplier1),
			data.Approval.Internal.GetApproverName(approval.Supplier2),
			data.Approval.Internal.GetApproverName(approval.Customer1),
			data.Approval.Internal.GetApproverName(approval.Customer2),
		}
		for _, u := range users {
			if u == "" || u == data.Approval.Creator {
				continue
			}
			currentUser := a.userRepo.FindByUserId(data.RequestSession, u)
			if currentUser == nil || currentUser.Email == "" {
				continue
			}
			md.Username = currentUser.Forename + " " + currentUser.Lastname
			a.sendMail(data.RequestSession, currentUser, currentUser.Email, "approvalFinalized", md)
		}
	case approval.TypePlausibility:
		if data.Approval.Plausibility.State.State == approval.Approved {
			md.State = "OK"
			md.StateDE = "OK"
		} else if data.Approval.Plausibility.State.State == approval.Declined {
			md.State = "Not OK"
			md.StateDE = "Nicht OK"
		} else {
			md.State = "Aborted"
			md.StateDE = "Abgebrochen"
		}
		md.DelegatedTo = "-"
		if data.DelegatedTo != "" {
			delegationUser := a.userRepo.FindByUserId(data.RequestSession, data.DelegatedTo)
			if delegationUser == nil || delegationUser.Email == "" {
				break
			}
			md.DelegatedTo = delegationUser.Forename + " " + delegationUser.Lastname
		}
		approverUser := a.userRepo.FindByUserId(data.RequestSession, data.Approval.Plausibility.Approver)
		if approverUser == nil || approverUser.Email == "" {
			break
		}
		md.Reviewer = approverUser.Forename + " " + approverUser.Lastname
		md.ReviewerComment = data.Approval.Plausibility.ApproveComment

		md.Username = creatorUser.Forename + " " + creatorUser.Lastname
		a.sendMail(data.RequestSession, creatorUser, creatorUser.Email, "reviewFinalized", md)
		if approverUser.User == creatorUser.User {
			return
		}
		md.Username = md.Reviewer
		a.sendMail(data.RequestSession, approverUser, approverUser.Email, "reviewFinalized", md)
	}
}

func (a *ApprovalMail) sendTaskMail(data observermngmt.ApprovalTaskData) {
	mailData := struct {
		Username  string
		Type      string
		TypeDE    string
		Link      string
		Requestor string
		Comment   string
		IsGroup   bool
		GroupName string
		GroupLink string
		Versions  []struct {
			Num         int
			ProjectName string
			ProjectLink string
			VersionName string
			VersionLink string
		}
	}{}

	mailData.Link = conf.Config.Server.DisukoHost + "/#/dashboard/tasks/" + data.TaskId

	targetUser := a.userRepo.FindByUserId(data.RequestSession, data.TargetUser)
	if targetUser.Email == "" {
		return
	}
	mailData.Username = targetUser.Forename + " " + targetUser.Lastname
	creatorUser := a.userRepo.FindByUserId(data.RequestSession, data.Creator)
	mailData.Requestor = creatorUser.Forename + " " + creatorUser.Lastname

	pr := a.projectRepo.FindByKey(data.RequestSession, data.ProjectId, false)
	mailData.IsGroup = pr.IsGroup
	if pr.IsGroup {
		mailData.GroupName = pr.Name
		mailData.GroupLink = conf.Config.Server.DisukoHost + "/#/dashboard/groups/" + data.ProjectId
		for i, a := range data.Approvables {
			prLink := conf.Config.Server.DisukoHost + "/#/dashboard/projects/" + a.ProjectKey
			versionLink := prLink + "/versions/" + a.ApprovableSPDX.VersionKey
			mailData.Versions = append(mailData.Versions, struct {
				Num         int
				ProjectName string
				ProjectLink string
				VersionName string
				VersionLink string
			}{
				Num:         i + 1,
				ProjectName: a.ProjectName,
				ProjectLink: prLink,
				VersionName: a.ApprovableSPDX.VersionName,
				VersionLink: versionLink,
			})
		}
	} else {
		approvable := data.Approvables[0]
		prLink := conf.Config.Server.DisukoHost + "/#/dashboard/projects/" + data.ProjectId
		versionLink := prLink + "/versions/" + approvable.ApprovableSPDX.VersionKey
		mailData.Versions = append(mailData.Versions, struct {
			Num         int
			ProjectName string
			ProjectLink string
			VersionName string
			VersionLink string
		}{
			ProjectName: pr.Name,
			ProjectLink: prLink,
			VersionName: approvable.ApprovableSPDX.VersionName,
			VersionLink: versionLink,
		})

	}
	mailData.Comment = data.Comment
	mailData.Type = "Review Request"
	mailData.TypeDE = "Prüfungsaufforderung"
	if data.Type == approval.TypeInternal {
		mailData.Type = "Approval"
		mailData.TypeDE = "Freigabeaufforderung"
	}
	a.sendMail(data.RequestSession, targetUser, targetUser.Email, "taskApproval", mailData)
}

func (a *ApprovalMail) sendMail(rs *logy.RequestSession, user *user.User, mail string, templ string, data interface{}) {
	if !user.Deprovisioned.IsZero() {
		return
	}
	go func() {
		defer func() {
			if err := recover(); err != nil {
				logy.Errorf(rs, "Could not send email %v", err)
			}
		}()
		err := a.mailClient.Send(mail, templ, data)
		if err != nil {
			logy.Errorf(rs, "Failed to send the email: %v", err)
		} else {
			logy.Infof(rs, "Email sent successfully!")
		}
	}()
}
