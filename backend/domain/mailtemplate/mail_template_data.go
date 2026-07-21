// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package mailtemplate

type ApprovalFinalizedVersionData struct {
	Num               int    `description:"1-based row number of this project/version in the finalized approval's version list"`
	ProjectName       string `description:"name of the project the version belongs to"`
	ProjectLink       string `description:"link to the project dashboard page"`
	VersionName       string `description:"name of the reviewed/approved version"`
	VersionLink       string `description:"link to the version dashboard page"`
	ReviewRemarksLink string `description:"link to the review remarks for this version's SBOM"`
}

type ApprovalFinalizedMailData struct {
	Username        string                         `description:"name of the mail recipient (creator, reviewer or approver depending on which mail is sent)"`
	Link            string                         `description:"link to the group's or project's approvals page"`
	State           string                         `description:"outcome of the approval/review in English, e.g. Approved, Declined, Aborted, OK, Not OK"`
	StateDE         string                         `description:"outcome of the approval/review in German, e.g. Genehmigt, Nicht genehmigt, Abgebrochen"`
	Requestor       string                         `description:"name of the user who requested the approval/review"`
	Comment         string                         `description:"comment left by the requestor when creating the approval"`
	Reviewer        string                         `description:"name of the user who performed the plausibility review"`
	DelegatedTo     string                         `description:"name of the user the review was delegated to, or '-' if not delegated"`
	ReviewerComment string                         `description:"comment left by the reviewer/approver on the plausibility review"`
	IsGroup         bool                           `description:"whether the approval/review is for a group of projects rather than a single project"`
	GroupName       string                         `description:"name of the group, set only when IsGroup is true"`
	GroupLink       string                         `description:"link to the group dashboard page, set only when IsGroup is true"`
	Versions        []ApprovalFinalizedVersionData `description:"list of project/version entries covered by this finalized approval, set only when IsGroup is true"`
}

type TaskApprovalVersionData struct {
	Num         int    `description:"1-based row number of this project/version in the approval task's version list"`
	ProjectName string `description:"name of the project the version belongs to"`
	ProjectLink string `description:"link to the project dashboard page"`
	VersionName string `description:"name of the version awaiting approval/review"`
	VersionLink string `description:"link to the version dashboard page"`
}

type TaskApprovalMailData struct {
	Username  string                    `description:"name of the mail recipient who has an open approval/review task"`
	Type      string                    `description:"kind of task requested in English, e.g. Review Request or Approval"`
	TypeDE    string                    `description:"kind of task requested in German, e.g. Prüfungsaufforderung or Freigabeaufforderung"`
	Link      string                    `description:"link to the task on the dashboard"`
	Requestor string                    `description:"name of the user who created the approval/review request"`
	Comment   string                    `description:"comment left by the requestor when creating the task"`
	IsGroup   bool                      `description:"whether the task covers a group of projects rather than a single project"`
	GroupName string                    `description:"name of the group, set only when IsGroup is true"`
	GroupLink string                    `description:"link to the group dashboard page, set only when IsGroup is true"`
	Versions  []TaskApprovalVersionData `description:"list of project/version entries covered by this task, set only when IsGroup is true"`
}

type ReviewCreatedMailData struct {
	Username         string `description:"name of the mail recipient subscribed to overall review notifications"`
	ProjectName      string `description:"name of the project the review was created for"`
	ReviewerFullName string `description:"name of the user who created the review"`
	Status           string `description:"review state in English"`
	StatusDE         string `description:"review state in German"`
	Comment          string `description:"comment left by the reviewer when creating the review"`
	ReviewsLink      string `description:"link to the review remarks page for the version"`
	ProjectLink      string `description:"link to the project dashboard page"`
	VersionName      string `description:"name of the version the review was created for"`
	VersionLink      string `description:"link to the version dashboard page"`
}

type SpdxUploadedMailData struct {
	Username    string `description:"name of the mail recipient subscribed to SPDX upload notifications"`
	ProjectName string `description:"name of the project the SBOM was uploaded to"`
	ProjectLink string `description:"link to the project dashboard page"`
	VersionName string `description:"name of the version the SBOM was uploaded to"`
	VersionLink string `description:"link to the version dashboard page"`
	SbomLink    string `description:"link to the uploaded SBOM component page"`
}

type ApprovalInactiveMailData struct {
	Username     string `description:"name of the mail recipient owning the inactive approval"`
	ProjectName  string `description:"name of the project whose approval has been inactive"`
	ProjectLink  string `description:"link to the project dashboard page"`
	DeletionDate string `description:"date on which the project will be deleted if the approval stays inactive"`
	InactiveDays int    `description:"number of days the approval has been inactive"`
}

type DummyDeletionMailData struct {
	Username    string `description:"name of the mail recipient owning the dummy project"`
	ProjectName string `description:"name of the dummy project scheduled for deletion"`
	Days        int    `description:"number of days remaining before the dummy project is deleted"`
}

var MailTemplateDataTypes = map[MailTemplateKey]any{
	MailTemplateKeyApprovalFinalized: ApprovalFinalizedMailData{},
	MailTemplateKeyApprovalInactive:  ApprovalInactiveMailData{},
	MailTemplateKeyDummyDeletion:     DummyDeletionMailData{},
	MailTemplateKeyReviewCreated:     ReviewCreatedMailData{},
	MailTemplateKeyReviewFinalized:   ApprovalFinalizedMailData{},
	MailTemplateKeySpdxUploaded:      SpdxUploadedMailData{},
	MailTemplateKeyTaskApproval:      TaskApprovalMailData{},
}
