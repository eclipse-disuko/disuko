// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"time"

	"github.com/eclipse-disuko/disuko/domain/overallreview"

	"github.com/eclipse-disuko/disuko/domain/department"
	"github.com/eclipse-disuko/disuko/domain/obligation"
	"github.com/eclipse-disuko/disuko/domain/project/components"
	"github.com/eclipse-disuko/disuko/domain/user"
	"github.com/eclipse-disuko/disuko/logy"

	"github.com/eclipse-disuko/disuko/domain/license"
)

func (entity *DisclosureDocumentMeta) ToDto(department *department.Department, missing bool) DisclosureDocumentMetaDto {
	dto := DisclosureDocumentMetaDto{
		SupplierName:    entity.SupplierName,
		SupplierAddress: entity.SupplierAddress,
		SupplierNr:      entity.SupplierNr,
		DeptMissing:     missing,
	}
	if department != nil {
		dto.SupplierDept = department.ToDto()
	}
	return dto
}

func (entity *NoticeContactMeta) ToDto() NoticeContactMetaDto {
	return NoticeContactMetaDto{
		Address: entity.Address,
		Email:   entity.Email,
	}
}

func (entity *CustomerMeta) ToDto(department *department.Department, missing bool) CustomerMetaDto {
	dto := CustomerMetaDto{
		Address:     entity.Address,
		FRI:         entity.FRI,
		SRI:         entity.SRI,
		UserFRI:     nil,
		UserSRI:     nil,
		DeptMissing: missing,
	}
	if department != nil {
		dto.Dept = department.ToDto()
	}
	return dto
}

type NoticeContactMetaDto struct {
	Address string `json:"address" validate:"lte=300"`
	Email   string `json:"email" validate:"lte=80"`
}

type CustomerMetaDto struct {
	Dept        *department.DepartmentDto `json:"dept" validate:"OmitEmptyStruct"`
	DeptMissing bool                      `json:"deptMissing"`
	Address     string                    `json:"address" validate:"lte=300"`
	FRI         string                    `json:"fRI" validate:"NeFieldIfSet=SRI,RealInternalUser"`
	SRI         string                    `json:"sRI" validate:"RealInternalUser"`
	UserFRI     *user.UserDto             `json:"userFRI"`
	UserSRI     *user.UserDto             `json:"userSRI"`
}

type DisclosureDocumentMetaDto struct {
	SupplierName    string                    `json:"supplierName" validate:"lte=80"`
	SupplierAddress string                    `json:"supplierAddress" validate:"lte=300"`
	SupplierNr      string                    `json:"supplierNr" validate:"lte=25"`
	SupplierDept    *department.DepartmentDto `json:"supplierDept"`
	DeptMissing     bool                      `json:"deptMissing"`
}

type ProjectSettingsRequest struct {
	DocumentMeta      DisclosureDocumentMeta `json:"documentMeta" validate:"OmitEmptySubStructWith=SupplierDept SupplierExtraData.External=false"`
	CustomerMeta      CustomerMetaDto        `json:"customerMeta"`
	NoticeContactMeta NoticeContactMetaDto   `json:"noticeContactMeta"`
	SupplierExtraData SupplierExtraData      `json:"supplierExtraData"`
	NoFossProject     bool                   `json:"noFossProject"`
	CustomIds         []ProjectCustomIdDto   `json:"customIds"`
}

type ProjectSettingsDto struct {
	DocumentMeta      DisclosureDocumentMetaDto `json:"documentMeta"`
	CustomerMeta      CustomerMetaDto           `json:"customerMeta"`
	SupplierExtraData SupplierExtraDataDto      `json:"supplierExtraData"`
	NoticeContactMeta NoticeContactMetaDto      `json:"noticeContactMeta"`
	NoFossProject     bool                      `json:"noFossProject"`
}

type UserManagementDto struct {
	Users []ProjectMemberDto `json:"users"`
}

func (entity *UserManagementEntity) ToDto(requestSession *logy.RequestSession, userProvider user.IUserDtoProvider) UserManagementDto {
	users := make([]ProjectMemberDto, 0)
	for _, user := range entity.Users {
		users = append(users, user.ToDto(requestSession, userProvider))
	}
	return UserManagementDto{
		Users: users,
	}
}

type ProjectRoleDto struct {
	ProjectName   string   `json:"projectName"`
	ProjectKey    string   `json:"projectKey"`
	UserId        string   `json:"userId"`
	UserType      UserType `json:"userType"`
	IsResponsible bool     `json:"responsible"`
}

type SubscriptionsDto struct {
	Spdx          bool `json:"spdx"`
	OverallReview bool `json:"overallReview"`
}

type ProjectMemberDto struct {
	UserId        string        `json:"userId"`
	UserType      UserType      `json:"userType"`
	Created       time.Time     `json:"created,omitempty"`
	Comment       string        `json:"comment"`
	IsResponsible bool          `json:"responsible"`
	UserProfile   *user.UserDto `json:"userProfile"`
}

func (entity *Subscriptions) ToDto() SubscriptionsDto {
	return SubscriptionsDto{
		Spdx:          entity.Spdx,
		OverallReview: entity.OverallReview,
	}
}

func (dto *SubscriptionsDto) ToEntity() Subscriptions {
	return Subscriptions{
		Spdx:          dto.Spdx,
		OverallReview: dto.OverallReview,
	}
}

func (entity *ProjectMemberEntity) ToDto(requestSession *logy.RequestSession, userProvider user.IUserDtoProvider) ProjectMemberDto {
	userProfile := userProvider.FindByUserId(requestSession, entity.UserId)
	if userProfile == nil {
		userProfile = &user.User{}
	}
	return ProjectMemberDto{
		UserId:        entity.UserId,
		UserType:      entity.UserType,
		UserProfile:   userProfile.ToDto(),
		Comment:       entity.Comment,
		IsResponsible: entity.IsResponsible,
	}
}

func (entity *ProjectMemberEntity) ToProjectRoleDto(pr *Project) ProjectRoleDto {
	return ProjectRoleDto{
		ProjectName:   pr.Name,
		ProjectKey:    pr.Key,
		UserType:      entity.UserType,
		IsResponsible: entity.IsResponsible,
	}
}

type TokenDto struct {
	Key         string      `json:"_key"`
	Company     string      `json:"company"`
	Description string      `json:"description"`
	Expiry      string      `json:"expiry"`
	Created     time.Time   `json:"created,omitempty"`
	TokenSecret string      `json:"tokenSecret"`
	Status      TokenStatus `json:"status"`
}

func (entity *Token) ToDto() TokenDto {
	return TokenDto{
		Key:         entity.Key,
		Company:     entity.Company,
		Description: entity.Description,
		Expiry:      entity.Expiry,
		TokenSecret: "",
		Status:      entity.Status,
		Created:     entity.Created,
	}
}

func (entity *Token) ToDtoWithSecret() TokenDto {
	return TokenDto{
		Key:         entity.Key,
		Company:     entity.Company,
		Description: entity.Description,
		Expiry:      entity.Expiry,
		TokenSecret: entity.TokenSecret,
		Status:      entity.Status,
	}
}

type ProjectsResponse struct {
	Projects []ProjectSlimDto `json:"projects"`
	Count    int              `json:"count"`
}

type ProjectsChildren struct {
	List     []ProjectChildrenCombiDto `json:"list"`
	Projects []ProjectSlimDto          `json:"projects"`
}

type ProjectRequestDto struct {
	Name            string                  `json:"name" validate:"required,gte=3,lte=80"`
	SchemaLabel     string                  `json:"schemaLabel" validate:"required,gte=3,lte=50"`
	PolicyLabels    []string                `json:"policyLabels" validate:"dive,gte=3,lte=80"`
	ProjectLabels   []string                `json:"projectLabels" validate:"dive,gte=3,lte=80"`
	FreeLabels      []string                `json:"freeLabels" validate:"dive,gte=1,lte=20"`
	Children        []string                `json:"children" validate:"dive,gte=3,lte=80"`
	Description     string                  `json:"description" validate:"lte=10000"`
	Owner           string                  `json:"owner" validate:"lte=50"`
	IsGroup         bool                    `json:"isGroup"`
	IsNoFoss        bool                    `json:"isNoFoss"`
	CreationMode    CreationProjectMode     `json:"creationMode"`
	ProjectSettings *ProjectSettingsRequest `json:"projectSettings"`
	ApplicationMeta ApplicationMetaDto      `json:"applicationMeta"`
	Parent          string                  `json:"parent"`
	ParentName      string                  `json:"parentName"`
}

type SupplierExtraData struct {
	// TODO: move into something named "ApproverPresets"
	FRI      string `json:"fRI" validate:"NeFieldIfSet=SRI,RealUser"`
	SRI      string `json:"sRI" validate:"RealUser"`
	External bool   `json:"external"`
}

type SupplierExtraDataDto struct {
	// TODO: move into something named "ApproverPresets"
	FRI      string        `json:"fRI" validate:"NeFieldIfSet=SRI,RealInternalUser"`
	SRI      string        `json:"sRI" validate:"RealInternalUser"`
	UserFRI  *user.UserDto `json:"userFRI"`
	UserSRI  *user.UserDto `json:"userSRI"`
	External bool          `json:"external"`
}

type CreationProjectMode string

const (
	DEFAULT CreationProjectMode = "default"
)

type ProjectMemberRequestDto struct {
	TargetUser    string   `json:"targetUser" validate:"required,gt=1,lte=100,RealUser"`
	UserType      UserType `json:"userType" validate:"required,gt=1,lte=20,SupportedUserType"`
	IsResponsible bool     `json:"responsible"`
	Comment       string   `json:"comment" validate:"lte=100"`
}

type MultiProjectMemberRequestDto struct {
	TargetUser     string   `json:"targetUser" validate:"required,gt=1,lte=100,RealUser"`
	UserType       UserType `json:"userType" validate:"required,gt=1,lte=20,SupportedUserType"`
	IsResponsible  bool     `json:"responsible"`
	Comment        string   `json:"comment" validate:"lte=100"`
	TargetProjects []string `json:"targetProjects"`
}

func (d *MultiProjectMemberRequestDto) ToUserData() *ProjectMemberRequestDto {
	return &ProjectMemberRequestDto{
		TargetUser:    d.TargetUser,
		UserType:      d.UserType,
		IsResponsible: d.IsResponsible,
		Comment:       d.Comment,
	}
}

type Response struct {
	Id       string `json:"id"`
	GroupId  string `json:"groupId"`
	Name     string `json:"name"`
	TaskGuid string `json:"taskGuid"`
}

type ProjectSchemaResponse struct {
	Content string `json:"content"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ProjectPublicResponse struct {
	Name        string    `json:"name" example:"Dummy Project"`
	Uuid        string    `json:"uuid" example:"dummy-id-----6b9c-44a7-8e01-14e67ef4404a"`
	Created     time.Time `json:"created" example:"2022-02-17T09:33:47.675838556Z"`
	Updated     time.Time `json:"updated" example:"2023-04-14T09:41:28.077559111Z"`
	Schema      string    `json:"schema" example:"EnterpriseIt"`
	Description string    `json:"description" example:"This is a dummy project."`
	IsGroup     bool      `json:"isGroup" example:"true"`
} //	@name	Project

type VersionStatusPublicResponse struct {
	Name             string                                     `json:"name" example:"1.0"`
	Status           ProjectVersionStatusType                   `json:"status" example:"unreviewed"`
	OverallReview    *overallreview.OverallReviewPublicResponse `json:"overallReview,omitempty"`
	LastSbomUploaded *time.Time                                 `json:"lastSbomUploaded,omitempty" example:"2023-05-15T10:52:39.187559111Z"`
} //	@name	VersionStatusPublicResponse

type ProjectStatusPublicResponse struct {
	Status        ProjectStatusType             `json:"status" example:"active"`
	VersionStatus []VersionStatusPublicResponse `json:"versionStatus,omitempty"`
} //	@name	ProjectStatus

type VersionAllResponse struct {
	ProjectVersions []ProjectVersion `json:"projectVersions"`
	Count           int              `json:"count"`
}

type VersionRequestDto struct {
	Name        string `json:"name" validate:"required,gte=1,lte=80" example:"1.0"`
	Description string `json:"description" validate:"lte=1000" example:""`
} //	@name	VersionRequest

type ProjectSearchDto struct {
	Tag string `json:"tag" validate:"required,gte=1,lte=80" example:"1.0"`
} //	@name	ProjectSearch

type ProjectSearchResDto struct {
	ChannelId string `json:"channelId" example:"dummy-id-----6b9c-44a7-8e01-14e67ef4404a"`
	SbomId    string `json:"sbomId" example:"dummy-id-----6b9c-44a7-8e01-14e67ef4404a"`
} //	@name	ProjectSearchRes

type SPDXUploadResponse struct {
	DocIsValid              bool   `json:"docIsValid"`
	ValidationFailedMessage string `json:"validationFailedMessage" example:""`
	Hash                    string `json:"hash" example:"some-hash----b75f6c1ac76d1517bfef85"`
	FileUploaded            bool   `json:"fileUploaded"`
	Id                      string `json:"id" example:"SPDXRef-DOCUMENT"`
	SbomGuid                string `json:"sbomguid" example:"dummy-id-----6b9c-44a7-8e01-14e67ef4404a"`
} //	@name	SpdxUploadResponse

type SPDXSetTagRequestDto struct {
	Tag string `json:"tag" validate:"required,gte=1,lte=80" example:"1.0"`
}

type VersionCreationResponseMin struct {
	Name    string `json:"name" example:"1.0"`
	Uuid    string `json:"uuid" example:"817e18e3-c0c7-4552-b4ca-aac875aee990"`
	Success bool   `json:"success"`
	Message string `json:"message" example:"Resource created"`
}

type VersionPublicResponseMin struct {
	Name string `json:"name" example:"1.0"`
	Uuid string `json:"uuid" example:"817e18e3-c0c7-4552-b4ca-aac875aee990"`
} //	@name	VersionPublicResponseMin

type VersionPublicResponse struct {
	Name             string                                     `json:"name" example:"1.0"`
	Description      string                                     `json:"description" example:"version description"`
	Status           ProjectVersionStatusType                   `json:"status" example:"unreviewed"`
	OverallReview    *overallreview.OverallReviewPublicResponse `json:"overallReview,omitempty"`
	LastSbomUploaded *time.Time                                 `json:"lastSbomUploaded,omitempty" example:"2023-05-15T10:52:39.187559111Z"`
	Uuid             string                                     `json:"uuid" example:"817e18e3-c0c7-4552-b4ca-aac875aee990"`
} //	@name	VersionDetails

type ExternalSourcePublicResponseDto struct {
	URL      string    `json:"url" validate:"required,gt=1,lte=2000,url" example:"file:///SOME_PATH/PUBLIC/foss"`
	Comment  string    `json:"comment" validate:"lte=2000" example:"Test"`
	Created  time.Time `json:"created" example:"2023-06-02T20:14:11.358265366Z"`
	Origin   string    `json:"origin" example:"UI"`
	Uploader string    `json:"uploader" example:"Name of Uploader"`
} //	@name	ExternalSource

type VersionHistoryPublicResponse struct {
	Name    string    `json:"name" example:"SBOM Demonstration/v1"`
	Updated time.Time `json:"updated" example:"2023-05-26T06:20:36.051198856Z"`
	Valid   bool      `json:"valid" example:"false"`
	Id      string    `json:"id" example:"dummy-id----4672-ad7b-757530f2580b"`
} //	@name	VersionHistory

type SPDXMetaPublicResponse struct {
	Name     string    `json:"name" example:"SBOM Demonstration/v2"`
	Id       string    `json:"id" example:"SPDXRef-DOCUMENT"`
	Version  string    `json:"version" example:"SPDX-2.2"`
	Creators string    `json:"creators" example:"Tool"`
	Created  time.Time `json:"created" example:"2023-06-13T07:58:40.572647367Z"`
	Uploaded time.Time `json:"uploaded" example:"2023-06-13T07:58:40.572647367Z"`
	Status   bool      `json:"status"`
	IsRetain bool      `json:"isRetain" example:"false"`
	IsLocked bool      `json:"isLocked" example:"true"`
	Tag      string    `json:"tag" example:"release-2025"`
} //	@name	SpdxMetaData

type ScanRemarkStatus string

const (
	PROBLEM     = string("PROBLEM")
	WARNING     = string("WARNING")
	INFORMATION = string("INFORMATION")
)

type QualityScanRemarks struct {
	Status            ScanRemarkStatus                  `json:"status"`
	RemarkKey         string                            `json:"remarkKey"`
	Name              string                            `json:"name"`
	Version           string                            `json:"version"`
	Type              ComponentType                     `json:"type"`
	DescriptionKey    string                            `json:"descriptionKey"`
	SpdxId            string                            `json:"spdxId"`
	PolicyRuleStatus  []*components.PolicyRuleStatusDto `json:"policyRuleStatus"`
	UnmatchedLicenses []*components.UnmatchedLicenseDto `json:"unmatchedLicenses"`
}

type ComponentType string

const (
	PACKAGE ComponentType = "Package"
	FILE    ComponentType = "File"
	SNIPPET ComponentType = "Snippet"
	ROOT    ComponentType = "Root"
	PROJECT ComponentType = "Project"
)

type QualityLicenseRemarks2 struct {
	License     string                      `json:"license"`
	Warnings    bool                        `json:"warnings"`
	Alarms      bool                        `json:"alarms"`
	Obligations []*obligation.ObligationDto `json:"obligations"`
	Affected    []*AffectedComponent        `json:"affected"`
}

type AffectedComponent struct {
	SpdxId           string                         `json:"spdxid"`
	Name             string                         `json:"name"`
	Version          string                         `json:"version"`
	PolicyRuleStatus []*components.PolicyRuleStatus `json:"policyRuleStatus"`
}

type QualityLicenseRemarks struct {
	Status           string                         `json:"status"`
	Type             string                         `json:"type"`
	Remark           string                         `json:"remark"`
	Name             string                         `json:"name"`
	Version          string                         `json:"version"`
	License          string                         `json:"license"`
	Description      string                         `json:"description"`
	SpdxId           string                         `json:"spdxId"`
	PolicyRuleStatus []*components.PolicyRuleStatus `json:"policyRuleStatus"`
}

type QualityGeneralRemarks struct{}

type SpdxStatusInformation struct {
	Disclaimer     string                `json:"disclaimer" example:"Disclaimer"`
	ScanRemarks    string                `json:"scanRemarks" example:"Some scan remarks"`
	LicenseRemarks string                `json:"licenseRemarks" example:"Some license remarks"`
	GeneralRemarks string                `json:"generalRemarks" example:"Some general remarks"`
	Components     []SpdxStatusComponent `json:"components"`
} //	@name	SpdxStatusInformation

// SpdxStatusComponent is a merge of ComponentInfo and QualityScanRemarks and QualityLicenseRemarks
type SpdxStatusComponent struct {
	SpdxId           string                     `json:"spdxId" example:"dummy-spdx-id----component23456"`
	License          string                     `json:"license" example:"MIT"`
	Name             string                     `json:"name" example:"@some/component"`
	Version          string                     `json:"version" example:"7.19.6"`
	ScanRemarks      []SpdxStatusScanRemarks    `json:"scanRemarks"`
	LicenseRemarks   []SpdxStatusLicenseRemarks `json:"licenseRemarks"`
	PolicyRuleStatus []SpdxStatusPolicy         `json:"policyRuleStatus"`
	UsedAliases      []UsedAlias                `json:"usedAliases"`
	UsedDecision     *UsedDecision              `json:"usedDecision,omitempty"`
} //	@name	SpdxStatusComponent

// UsedAlias contains information about a alias which is used to identiy a component license.
type UsedAlias struct {
	Name           string `json:"name" example:"licenseRefMIT"`
	ReferencedName string `json:"referencedName" example:"MIT"`
} //	@name	UsedAlias

// UsedDecision contains information about a alias decision.
type UsedDecision struct {
	Expression  string `json:"expression" example:"MPL-2.0 OR Apache-2.0"`
	LicenseID   string `json:"licenseID" example:"Apache-2.0"`
	LicenseName string `json:"name" example:"Apache License 2.0"`
} //	@name	UsedDecision

type SpdxStatusScanRemarks struct {
	Status      ScanRemarkStatus `json:"status" example:"INFORMATION"`
	Remark      string           `json:"remark" example:"Display copyright notice"`
	Description string           `json:"description" example:"Description of scan remarks"`
} //	@name	SpdxStatusScanRemarks

type SpdxStatusLicenseRemarks struct {
	Status         string `json:"status" example:"INFORMATION"`
	Remark         string `json:"remark" example:"Display copyright notice"`
	Type           string `json:"type" example:"obligation"`
	LicenseMatched string `json:"licenseMatched" example:"MIT"`
	Description    string `json:"description" example:"Description of license remarks"`
} //	@name	SpdxStatusLicenceRemarks

// SpdxStatusPolicy is a public JSON representation for PolicyRuleStatus
type SpdxStatusPolicy struct {
	Name           string           `json:"name" example:"USE IT"`
	LicenseMatched string           `json:"licenseMatched" example:"MIT"`
	Type           license.ListType `json:"type" example:"allow"`
	Used           bool             `json:"used"`
	Description    string           `json:"description" example:"Description of policy"`
} //	@name	SpdxStatusPolicy

type NoticeComponent struct {
	Name        string `json:"name" example:"@some/component"`
	Version     string `json:"version" example:"7.19.6"`
	LicenseName string `json:"licenseName" example:"MIT License"`
	LicenseID   string `json:"licenseID" example:"MIT"`
	Copyright   string `json:"copyright" example:"Copyright"`
} //	@name	NoticeComponent

type NoticeLicense struct {
	Text      string `json:"text" example:"MIT License"`
	Name      string `json:"name" example:"MIT License"`
	LicenseID string `json:"id" example:"MIT"`
} //	@name	NoticeLicence

type NoticeJSONMeta struct {
	Title       string `json:"title" example:"Copyright notices and license information"`
	Description string `json:"description" example:"Description"`
} //	@name	NoticeMetaData

type NoticeFileJSON struct {
	Components []NoticeComponent `json:"components"`
	Licenses   []NoticeLicense   `json:"licenses"`
	Meta       NoticeJSONMeta    `json:"meta"`
} //	@name	NoticeFile

type SourceExternalDTO struct {
	Key        string     `json:"_key" validate:"lte=36" example:""`
	SourceType SourceType `json:"sourceType" validate:"lte=20"`
	URL        string     `json:"url" validate:"url" example:"file:///SOME_PATH/PUBLIC/foss"`
	Comment    string     `json:"comment" validate:"lte=2000" example:"Test"`
	Hash       string     `json:"hash" validate:"lte=50"`
	FileSize   int64      `json:"fileSize"`
	Created    time.Time  `json:"created" example:"2023-07-17T12:06:32.743823343Z"`
	Origin     string     `json:"origin" example:"API"`
	Uploader   string     `json:"uploader" example:"Uploader"`
} //	@name	SourceExternal

func (entity SupplierExtraData) ToDTO() SupplierExtraDataDto {
	return SupplierExtraDataDto{
		FRI:      entity.FRI,
		SRI:      entity.SRI,
		UserFRI:  nil,
		UserSRI:  nil,
		External: entity.External,
	}
}

func (entity *SourceExternal) ToDTO() SourceExternalDTO {
	return SourceExternalDTO{
		SourceType: ExternalLink,
		Key:        entity.Key,
		URL:        entity.URL,
		Comment:    entity.Comment,
		Origin:     entity.Origin,
		Uploader:   entity.Uploader,
		Created:    entity.Created,
	}
}

func (entity *SourceExternal) ToPublicDTO() ExternalSourcePublicResponseDto {
	return ExternalSourcePublicResponseDto{
		URL:      entity.URL,
		Comment:  entity.Comment,
		Created:  entity.Created,
		Origin:   entity.Origin,
		Uploader: entity.Uploader,
	}
}

type SbomStatusPublicResponseDto struct {
	Status string `json:"status" example:"red"`
}
