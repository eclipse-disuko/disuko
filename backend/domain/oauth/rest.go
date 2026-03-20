// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package oauth

import (
	"fmt"
	"time"

	"github.com/eclipse-disuko/disuko/domain/label"
	"github.com/eclipse-disuko/disuko/domain/user"
	"github.com/eclipse-disuko/disuko/helper"
)

const (
	UserInternal         = "Internal"
	UserNonInternal      = "NonInternal"
	UserLicenseManager   = "LicenseManager"
	UserPolicyManager    = "PolicyManager"
	UserProjectAnalyst   = "ProjectAnalyst"
	UserDomainAdmin      = "DomainAdmin"
	UserApplicationAdmin = "ApplicationAdmin"
	UserFOSSOffice       = "FOSSOffice"
)

type AccessLevel string

const (
	AccessLevelRead   AccessLevel = "Read"
	AccessLevelCreate AccessLevel = "Create"
	AccessLevelUpdate AccessLevel = "Update"
	AccessLevelDelete AccessLevel = "Delete"
)

type OAuthTokenResponse struct {
	Rights   *AccessAndRolesRights `json:"rights"`
	UserGuid string                `json:"userGuid"`
	Profile  *user.UserDto         `json:"profile"`
}

type RefreshTokenResponseDto struct {
	AccessToken string    `json:"accessToken"`
	Expiry      time.Time `json:"expiry,omitempty"`
}

type RefreshTokenRequestDto struct {
	RefreshToken string `json:"refreshToken" validate:"required,gte=50,lte=1000"`
}

type AccessResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

// CRUDRights describes classical rights to perform CRUD operation on an object.
type CRUDRights struct {
	Create bool `json:"create"`
	Read   bool `json:"read"`
	Update bool `json:"update"`
	Delete bool `json:"delete"`
}

// ActionRights describes right flags to perform specific actions: Uploading, downloading.
type ActionRights struct {
	Upload   bool `json:"upload"`
	Download bool `json:"download"`
	Delete   bool `json:"delete"`
}
type ProjectTypeRights struct {
	VehiclePlatform    CRUDRights `json:"vehiclePlatform"`
	EnterprisePlatform CRUDRights `json:"enterprisePlatform"`
	MobilePlatform     CRUDRights `json:"mobilePlatform"`
	OtherPlatform      CRUDRights `json:"otherPlatform"`
}

func (ptr *ProjectTypeRights) setAll(value bool) {
	ptr.VehiclePlatform.SetAll(value)
	ptr.EnterprisePlatform.SetAll(value)
	ptr.MobilePlatform.SetAll(value)
	ptr.OtherPlatform.SetAll(value)
}

type CRUDRightsAssigned struct {
	CRUDRights
	ReadWhenAssigned   bool `json:"readWhenAssigned"`
	UpdateWhenAssigned bool `json:"updateWhenAssigned"`
	DeleteWhenAssigned bool `json:"deleteWhenAssigned"`
}

type AccessAndRolesRights struct {
	AllowSchema                    CRUDRights        `json:"allowSchema"`
	AllowLabel                     CRUDRights        `json:"allowLabel"`
	AllowPolicy                    CRUDRights        `json:"allowPolicy"`
	AllowProjectPolicy             CRUDRights        `json:"allowProjectPolicy"`
	AllowLicense                   CRUDRights        `json:"allowLicense"`
	AllowProject                   CRUDRights        `json:"allowProject"`
	AllowProjectVersion            CRUDRights        `json:"allowProjectVersion"`
	AllowAllProjectUserManagement  CRUDRights        `json:"allowAllProjectUserManagement"`
	AllowAllProjectTokenManagement CRUDRights        `json:"allowAllProjectTokenManagement"`
	AllowProjectUserManagement     CRUDRights        `json:"allowProjectUserManagement"`
	AllowProjectTokenManagement    CRUDRights        `json:"allowProjectTokenManagement"`
	AllowProjectAudit              CRUDRights        `json:"allowProjectAudit"`
	AllowSBOMAction                ActionRights      `json:"allowSBOMAction"`
	AllowCCSAction                 ActionRights      `json:"allowCCSAction"`
	AllowObligation                CRUDRights        `json:"allowObligation"`
	AllowProjectGroup              CRUDRights        `json:"allowProjectGroup"`
	IsInternal                     bool              `json:"isInternal"`
	AllowDisclosureDocument        ActionRights      `json:"allowDisclosureDocument"`
	AllowRequestApproval           CRUDRights        `json:"allowRequestApproval"`
	AllowRequestPlausi             CRUDRights        `json:"allowRequestPlausi"`
	AllowTools                     CRUDRights        `json:"allowTools"`
	AllowStyleguide                CRUDRights        `json:"allowStyleguide"`
	AllowS3Tests                   CRUDRights        `json:"allowS3Tests"`
	AllowSampleData                CRUDRights        `json:"allowSampleData"`
	AllowUsers                     CRUDRights        `json:"allowUsers"`
	Groups                         []string          `json:"groups"`
	AllowTask                      CRUDRights        `json:"allowTask"`
	AllowAnnouncement              CRUDRights        `json:"allowAnnouncement"`
	AllowChecklist                 CRUDRights        `json:"allowChecklist"`
	AllowReviewTemplates           CRUDRights        `json:"allowReviewTemplates"`
	AllowNewsbox                   CRUDRights        `json:"allowNewsbox"`
	AllowInternalToken             CRUDRights        `json:"allowInternalToken"`
	AllowCustomIDs                 CRUDRights        `json:"allowCustomIDs"`
	AllowExecuteChecklist          bool              `json:"allowExecuteChecklist"`
	AllowProjectType               ProjectTypeRights `json:"allowProjectType"`
	AllowLicenseDecision           CRUDRights        `json:"allowLicenseDecision"`
	AllowPolicyDecision            CRUDRights        `json:"allowPolicyDecision"`
	AllowFeatureFlag               CRUDRights        `json:"allowFeatureFlag"`
}

type AccessRightsDto struct {
	AllowProject                   CRUDRights `json:"allowProject"`
	AllowProjectGroup              CRUDRights `json:"allowProjectGroup"`
	AllowRequestApproval           CRUDRights `json:"allowRequestApproval"`
	AllowRequestPlausi             CRUDRights `json:"allowRequestPlausi"`
	AllowAllProjectUserManagement  CRUDRights `json:"allowAllProjectUserManagement"`
	AllowAllProjectTokenManagement CRUDRights `json:"allowAllProjectTokenManagement"`
	AllowProjectAudit              CRUDRights `json:"allowProjectAudit"`
	AllowSchema                    CRUDRights `json:"allowSchema"`
	AllowLabel                     CRUDRights `json:"allowLabel"`
	AllowPolicy                    CRUDRights `json:"allowPolicy"`
	AllowLicense                   CRUDRights `json:"allowLicense"`
	AllowObligation                CRUDRights `json:"allowObligation"`
	AllowTools                     CRUDRights `json:"allowTools"`
	AllowStyleguide                CRUDRights `json:"allowStyleguide"`
	AllowS3Tests                   CRUDRights `json:"allowS3Tests"`
	AllowSampleData                CRUDRights `json:"allowSampleData"`
	AllowUsers                     CRUDRights `json:"allowUsers"`
	AllowTask                      CRUDRights `json:"allowTask"`
	AllowAnnouncement              CRUDRights `json:"allowAnnouncement"`
	AllowChecklist                 CRUDRights `json:"allowChecklist"`
	AllowReviewTemplates           CRUDRights `json:"allowReviewTemplates"`
	AllowNewsbox                   CRUDRights `json:"allowNewsbox"`
	AllowInternalToken             CRUDRights `json:"allowInternalToken"`
	AllowCustomIDs                 CRUDRights `json:"allowCustomIDs"`
	AllowFeatureFlag               CRUDRights `json:"allowFeatureFlag"`
}

type ProjectAccessRightsDto struct {
	AllowProject                CRUDRights   `json:"allowProject"`
	AllowProjectGroup           CRUDRights   `json:"allowProjectGroup"`
	AllowProjectVersion         CRUDRights   `json:"allowProjectVersion"`
	AllowProjectPolicy          CRUDRights   `json:"allowProjectPolicy"`
	AllowProjectUserManagement  CRUDRights   `json:"allowProjectUserManagement"`
	AllowProjectTokenManagement CRUDRights   `json:"allowProjectTokenManagement"`
	AllowSBOMAction             ActionRights `json:"allowSBOMAction"`
	AllowCCSAction              ActionRights `json:"allowCCSAction"`
	AllowDisclosureDocument     ActionRights `json:"allowDisclosureDocument"`
	AllowLicenseDecision        CRUDRights   `json:"allowLicenseDecision"`
	AllowPolicyDecision         CRUDRights   `json:"allowPolicyDecision"`
	AllowRequestApproval        CRUDRights   `json:"allowRequestApproval"`
}

func (aarr *AccessAndRolesRights) ToAccessRightsDto() *AccessRightsDto {
	return &AccessRightsDto{
		AllowProject:                   aarr.AllowProject,
		AllowProjectGroup:              aarr.AllowProjectGroup,
		AllowRequestApproval:           aarr.AllowRequestApproval,
		AllowRequestPlausi:             aarr.AllowRequestPlausi,
		AllowAllProjectUserManagement:  aarr.AllowAllProjectUserManagement,
		AllowAllProjectTokenManagement: aarr.AllowAllProjectTokenManagement,
		AllowProjectAudit:              aarr.AllowProjectAudit,
		AllowSchema:                    aarr.AllowSchema,
		AllowLabel:                     aarr.AllowLabel,
		AllowPolicy:                    aarr.AllowPolicy,
		AllowLicense:                   aarr.AllowLicense,
		AllowObligation:                aarr.AllowObligation,
		AllowTools:                     aarr.AllowTools,
		AllowStyleguide:                aarr.AllowStyleguide,
		AllowS3Tests:                   aarr.AllowS3Tests,
		AllowSampleData:                aarr.AllowSampleData,
		AllowUsers:                     aarr.AllowUsers,
		AllowTask:                      aarr.AllowTask,
		AllowAnnouncement:              aarr.AllowAnnouncement,
		AllowChecklist:                 aarr.AllowChecklist,
		AllowReviewTemplates:           aarr.AllowReviewTemplates,
		AllowNewsbox:                   aarr.AllowNewsbox,
		AllowInternalToken:             aarr.AllowInternalToken,
		AllowCustomIDs:                 aarr.AllowCustomIDs,
		AllowFeatureFlag:               aarr.AllowFeatureFlag,
	}
}

func (aarr *AccessAndRolesRights) ToProjectAccessRightsDto() *ProjectAccessRightsDto {
	return &ProjectAccessRightsDto{
		AllowProject:                aarr.AllowProject,
		AllowProjectGroup:           aarr.AllowProjectGroup,
		AllowProjectVersion:         aarr.AllowProjectVersion,
		AllowProjectPolicy:          aarr.AllowProjectPolicy,
		AllowProjectUserManagement:  aarr.AllowProjectUserManagement,
		AllowProjectTokenManagement: aarr.AllowProjectTokenManagement,
		AllowSBOMAction:             aarr.AllowSBOMAction,
		AllowCCSAction:              aarr.AllowCCSAction,
		AllowDisclosureDocument:     aarr.AllowDisclosureDocument,
		AllowLicenseDecision:        aarr.AllowLicenseDecision,
		AllowPolicyDecision:         aarr.AllowPolicyDecision,
		AllowRequestApproval:        aarr.AllowRequestApproval,
	}
}

func (aarr *AccessAndRolesRights) IsLicenseManager() bool {
	return aarr.hasRole(UserLicenseManager)
}

func (aarr *AccessAndRolesRights) IsPolicyManager() bool {
	return aarr.hasRole(UserPolicyManager)
}

func (aarr *AccessAndRolesRights) IsProjectAnalyst() bool {
	return aarr.hasRole(UserProjectAnalyst)
}

func (aarr *AccessAndRolesRights) IsDomainAdmin() bool {
	return aarr.hasRole(UserDomainAdmin)
}

func (aarr *AccessAndRolesRights) IsApplicationAdmin() bool {
	return aarr.hasRole(UserApplicationAdmin)
}

func (aarr *AccessAndRolesRights) hasRole(role string) bool {
	return helper.Contains(role, aarr.Groups)
}

func (aarr *AccessAndRolesRights) IsFossOffice() bool {
	return aarr.hasRole(UserFOSSOffice)
}

func (aarr *AccessAndRolesRights) HasProjectTypeAccess(projectLabel string, accessLevel AccessLevel) bool {
	var rights *CRUDRights

	switch projectLabel {
	case label.VEHICLE_PLATFORM:
		rights = &aarr.AllowProjectType.VehiclePlatform
	case label.ENTERPRISE_PLATFORM:
		rights = &aarr.AllowProjectType.EnterprisePlatform
	case label.MOBILE_PLATFORM:
		rights = &aarr.AllowProjectType.MobilePlatform
	case label.OTHER_PLATFORM:
		rights = &aarr.AllowProjectType.OtherPlatform
	default:
		return false
	}

	switch accessLevel {
	case AccessLevelCreate:
		return rights.Create
	case AccessLevelRead:
		return rights.Read
	case AccessLevelUpdate:
		return rights.Update
	case AccessLevelDelete:
		return rights.Delete
	}

	return false
}

func (r *CRUDRights) SetAll(value bool) CRUDRights {
	r.Create = value
	r.Read = value
	r.Update = value
	r.Delete = value
	return *r
}

func (r *ActionRights) SetAll(value bool) {
	r.Upload = value
	r.Download = value
	r.Delete = value
}

func (r *ActionRights) SetReadOnly() {
	r.Download = true
	r.Upload = false
	r.Delete = false
}

func (r *CRUDRightsAssigned) SetAll(value bool) {
	r.Create = value
	r.Read = value
	r.Update = value
	r.Delete = value
	r.DeleteWhenAssigned = value
	r.UpdateWhenAssigned = value
	r.ReadWhenAssigned = value
}

func (r *CRUDRights) SetReadOnly() CRUDRights {
	r.Create = false
	r.Read = true
	r.Update = false
	r.Delete = false
	return *r
}

func (r *CRUDRights) ToString() interface{} {
	str := ""
	if r.Create {
		str += "Create "
	}
	if r.Read {
		str += "Read "
	}
	if r.Update {
		str += "Update "
	}
	if r.Delete {
		str += "Delete "
	}
	return str
}

func (r *ActionRights) ToString() interface{} {
	str := ""
	if r.Upload {
		str += "Upload "
	}
	if r.Download {
		str += "Download "
	}
	return str
}

func (r *AccessAndRolesRights) setAllFalse() {
	r.AllowSchema.SetAll(false)
	r.AllowLabel.SetAll(false)
	r.AllowPolicy.SetAll(false)
	r.AllowProjectPolicy.SetAll(false)
	r.AllowLicense.SetAll(false)
	r.AllowProject.SetAll(false)
	r.AllowRequestApproval.SetAll(false)
	r.AllowProjectGroup.SetAll(false)
	r.AllowProjectVersion.SetAll(false)
	r.AllowAllProjectUserManagement.SetAll(false)
	r.AllowAllProjectTokenManagement.SetAll(false)
	r.AllowProjectAudit.SetAll(false)
	r.AllowProjectUserManagement.SetAll(false)
	r.AllowProjectTokenManagement.SetAll(false)
	r.AllowSBOMAction.SetAll(false)
	r.AllowCCSAction.SetAll(false)
	r.AllowObligation.SetAll(false)
	r.AllowTools.SetAll(false)
	r.AllowStyleguide.SetAll(false)
	r.AllowS3Tests.SetAll(false)
	r.AllowSampleData.SetAll(false)
	r.AllowUsers.SetAll(false)
	r.AllowTask.SetAll(false)
	r.AllowAnnouncement.SetAll(false)
	r.AllowProjectType.setAll(false)
}

func (r *AccessAndRolesRights) setEachDefault() {
	r.AllowSchema.SetReadOnly()
	r.AllowLabel.SetReadOnly()
	r.AllowPolicy.SetReadOnly()
	r.AllowProjectPolicy.SetAll(false)
	r.AllowLicense.SetReadOnly()
	r.AllowProject.SetAll(false)
	r.AllowRequestApproval.SetAll(false)
	r.AllowProjectGroup.SetAll(false)
	r.AllowProjectVersion.SetAll(false)
	r.AllowAllProjectUserManagement.SetAll(false)
	r.AllowAllProjectTokenManagement.SetAll(false)
	r.AllowProjectAudit.SetAll(false)
	r.AllowProjectUserManagement.SetAll(false)
	r.AllowProjectTokenManagement.SetAll(false)
	r.AllowSBOMAction.SetAll(false)
	r.AllowCCSAction.SetAll(false)
	r.AllowObligation.SetReadOnly()
	r.AllowTools.SetAll(false)
	r.AllowStyleguide.SetAll(false)
	r.AllowS3Tests.SetAll(false)
	r.AllowSampleData.SetAll(false)
	r.AllowUsers.SetAll(false)
	r.AllowDisclosureDocument.SetAll(false)
	r.AllowTask.SetReadOnly()
	r.AllowAnnouncement.SetReadOnly()
	r.AllowChecklist.SetAll(false)
	r.AllowNewsbox.SetAll(false)
	r.AllowInternalToken.SetAll(false)
	r.AllowFeatureFlag.SetAll(false)
	r.AllowCustomIDs.SetAll(false)
	r.AllowReviewTemplates.SetAll(false)
}

func (r *AccessAndRolesRights) SetDefault() {
	r.setEachDefault()
}

func (r *AccessAndRolesRights) SetForInternal() {
	r.AllowProject.Create = true
	r.AllowProjectGroup.Create = true
	r.AllowProjectType.VehiclePlatform.Create = true
	r.AllowProjectType.MobilePlatform.Create = true
	r.AllowProjectType.EnterprisePlatform.Create = true
	r.AllowProjectType.OtherPlatform.Create = true
}

func (r *AccessAndRolesRights) SetForNonInternal() {
	r.AllowSchema.SetAll(false)
	r.AllowLicense.SetAll(false)
	r.AllowObligation.SetAll(false)
	r.AllowPolicy.SetReadOnly()
	r.AllowProject.SetAll(false)
	r.AllowProjectGroup.SetAll(false)
	r.AllowProjectType.setAll(false)
}

func (r *AccessAndRolesRights) SetForLicenseAdmin() {
	r.AllowObligation.SetAll(true)
	r.AllowLicense.SetAll(true)
	r.AllowReviewTemplates.SetAll(true)
}

func (r *AccessAndRolesRights) SetForPolicyAdmin() {
	r.AllowPolicy.SetAll(true)
	r.AllowReviewTemplates.SetAll(true)
}

func (r *AccessAndRolesRights) SetForFOSSOffice() {
	r.AllowProject.Read = true
	r.AllowProjectVersion.Read = true
	r.AllowCCSAction.SetReadOnly()
	r.AllowProjectVersion.SetReadOnly()
	r.AllowProjectPolicy.SetReadOnly()
	r.AllowProjectGroup.Read = true
	r.AllowReviewTemplates.SetAll(true)
	r.AllowChecklist.SetAll(true)
	r.AllowExecuteChecklist = true
	r.AllowProjectType.VehiclePlatform.SetReadOnly()
}

func (r *AccessAndRolesRights) SetForProjectAnalyst() {
	r.AllowProject.Read = true
	r.AllowProjectGroup.Read = true
	r.AllowProjectPolicy.Read = true
	r.AllowProjectVersion.Read = true
	r.AllowAllProjectUserManagement.Read = true
	r.AllowCCSAction.SetReadOnly()
	r.AllowSBOMAction.SetReadOnly()
	r.AllowDisclosureDocument.Download = true
	r.AllowExecuteChecklist = true

	r.AllowProjectType.VehiclePlatform.SetReadOnly()
	r.AllowProjectType.MobilePlatform.SetReadOnly()
	r.AllowProjectType.EnterprisePlatform.SetReadOnly()
	r.AllowProjectType.OtherPlatform.SetReadOnly()
}

func (r *AccessAndRolesRights) SetForApplicationAdmin() {
	r.AllowLabel.SetAll(true)
	r.AllowSchema.SetAll(true)
	r.AllowTools.SetAll(true)
	r.AllowSampleData.SetAll(true)
	r.AllowStyleguide.SetAll(true)
	r.AllowNewsbox.SetAll(true)
	r.AllowInternalToken.SetAll(true)
	r.AllowCustomIDs.SetAll(true)
	r.AllowFeatureFlag.SetAll(true)
}

func (r *AccessAndRolesRights) SetForDomainAdmin() {
	// Just having a role "domain admin" does not grant full access on project.
	// But in case "domain admin" is an owner of a project - will be granted in further checks for owner/supplier/viewer roles

	r.AllowProject.Read = true
	r.AllowRequestApproval.Create = true
	r.AllowRequestPlausi.Create = true
	r.AllowProjectGroup.Read = true
	r.AllowProjectPolicy.Read = true

	r.AllowProjectType.VehiclePlatform.SetReadOnly()
	r.AllowProjectType.MobilePlatform.SetReadOnly()
	r.AllowProjectType.EnterprisePlatform.SetReadOnly()
	r.AllowProjectType.OtherPlatform.SetReadOnly()

	r.AllowProjectVersion.SetReadOnly()
	r.AllowAllProjectUserManagement.SetAll(true)
	r.AllowAllProjectTokenManagement.SetAll(true)
	r.AllowProjectAudit.Create = true
	r.AllowProjectAudit.Read = true

	r.AllowSBOMAction.SetReadOnly()
	r.AllowCCSAction.SetReadOnly()

	r.AllowS3Tests.SetAll(true)
	r.AllowUsers.SetAll(true)
	r.AllowDisclosureDocument.Download = true

	r.AllowChecklist.SetAll(true)
	r.AllowReviewTemplates.SetAll(true)
	r.AllowNewsbox.SetAll(true)

	r.AllowExecuteChecklist = true
}

func (r *AccessAndRolesRights) SetForExternal() {
	r.AllowSchema.SetAll(false)
	r.AllowLicense.SetAll(false)
	r.AllowObligation.SetAll(false)
	r.AllowPolicy.SetAll(false)
	r.AllowProject.SetAll(false)
	r.AllowProjectGroup.SetAll(false)
	r.AllowProjectType.setAll(false)
}

func (r *AccessAndRolesRights) SetOwnerRights() {
	r.AllowProjectPolicy.SetReadOnly()
	r.AllowProject.Read = true
	r.AllowProject.Update = true
	r.AllowProject.Delete = true
	r.AllowRequestPlausi.Read = true
	r.AllowProjectVersion.SetAll(true)
	r.AllowProjectUserManagement.SetAll(true)
	r.AllowProjectTokenManagement.SetAll(true)
	r.AllowProjectGroup.Read = true
	r.AllowProjectGroup.Update = true
	r.AllowProjectGroup.Delete = true
	r.AllowSBOMAction.SetAll(true)
	r.AllowCCSAction.SetAll(true)
	r.AllowDisclosureDocument.Download = true
	r.AllowLicenseDecision.Read = true
	r.AllowPolicyDecision.Read = true
	r.AllowRequestApproval.SetAll(true)
	r.AllowExecuteChecklist = true
}

func (r *AccessAndRolesRights) SetProjectResponsibleRights() {
	r.AllowProjectPolicy.SetReadOnly()
	r.AllowProject.SetAll(true)
	r.AllowRequestPlausi.SetAll(true)
	r.AllowRequestApproval.SetAll(true)

	r.AllowProjectType.VehiclePlatform.SetAll(true)
	r.AllowProjectType.MobilePlatform.SetAll(true)
	r.AllowProjectType.EnterprisePlatform.SetAll(true)
	r.AllowProjectType.OtherPlatform.SetAll(true)
	r.AllowProjectVersion.SetAll(true)
	r.AllowProjectUserManagement.SetAll(true)
	r.AllowProjectTokenManagement.SetAll(true)
	r.AllowProjectGroup.SetAll(true)
	r.AllowSBOMAction.SetAll(true)
	r.AllowCCSAction.SetAll(true)
	r.AllowDisclosureDocument.Download = true
	r.AllowLicenseDecision.SetAll(true)
	r.AllowPolicyDecision.SetAll(true)
}

func (r *AccessAndRolesRights) SetSupplierRights() {
	r.AllowProjectPolicy.SetReadOnly()
	r.AllowProject.SetReadOnly()
	r.AllowProjectVersion.Create = true
	r.AllowProjectVersion.Read = true
	r.AllowProjectVersion.Update = true
	r.AllowProjectVersion.Delete = false
	r.AllowProjectUserManagement.SetReadOnly()
	r.AllowProjectTokenManagement.SetAll(true)
	r.AllowProjectGroup.SetReadOnly()
	r.AllowSBOMAction.SetAll(true)
	r.AllowCCSAction.SetAll(true)
	r.AllowDisclosureDocument.Download = true
	r.AllowExecuteChecklist = true

	r.AllowProjectType.VehiclePlatform.SetReadOnly()
	r.AllowProjectType.MobilePlatform.SetReadOnly()
	r.AllowProjectType.EnterprisePlatform.SetReadOnly()
	r.AllowProjectType.OtherPlatform.SetReadOnly()
}

func (r *AccessAndRolesRights) SetViewerRights() {
	r.AllowProjectPolicy.SetReadOnly()
	r.AllowProject.SetReadOnly()
	r.AllowProjectVersion.SetReadOnly()
	r.AllowProjectUserManagement.SetReadOnly()
	r.AllowProjectTokenManagement.SetAll(false)
	r.AllowProjectGroup.SetReadOnly()
	r.AllowSBOMAction.SetReadOnly()
	r.AllowCCSAction.SetReadOnly()
	r.AllowDisclosureDocument.Download = true
	r.AllowExecuteChecklist = true

	r.AllowProjectType.VehiclePlatform.SetReadOnly()
	r.AllowProjectType.MobilePlatform.SetReadOnly()
	r.AllowProjectType.EnterprisePlatform.SetReadOnly()
	r.AllowProjectType.OtherPlatform.SetReadOnly()
}

func (r *AccessAndRolesRights) SetPublicApiRights() {
	r.setAllFalse()
	r.AllowProject.SetReadOnly()
	r.AllowProjectGroup.SetReadOnly()
	r.AllowProjectVersion.Create = true
	r.AllowProjectVersion.Read = true
	r.AllowSBOMAction.Upload = true
	r.AllowCCSAction.Upload = true

	r.AllowProjectType.VehiclePlatform.SetReadOnly()
	r.AllowProjectType.MobilePlatform.SetReadOnly()
	r.AllowProjectType.EnterprisePlatform.SetReadOnly()
	r.AllowProjectType.OtherPlatform.SetReadOnly()
}

func (r *AccessAndRolesRights) ToString() interface{} {
	return fmt.Sprintf("Schema(%s) Label(%s) Policy(%s) License(%s) Project(%s) ProvectVersion(%s) ProjectUserManagement(%s) ProjectTokenManagement(%s) SBOMAction(%s) CCSAction(%s)",
		r.AllowSchema.ToString(), r.AllowLabel.ToString(), r.AllowPolicy.ToString(), r.AllowLicense.ToString(), r.AllowProject.ToString(),
		r.AllowProjectVersion.ToString(), r.AllowProjectUserManagement.ToString(), r.AllowProjectTokenManagement.ToString(), r.AllowSBOMAction.ToString(), r.AllowCCSAction.ToString())
}
