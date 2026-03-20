// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package project

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"

	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/helper/hash"

	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/domain/project/approvable"
	"github.com/eclipse-disuko/disuko/domain/project/pdocument"
	"github.com/eclipse-disuko/disuko/domain/schema"
	"github.com/google/uuid"
)

const (
	OriginApi    = "API"
	OriginUi     = "UI"
	OriginServer = "SERVER"
)

type Project struct {
	domain.RootEntity `bson:",inline"`
	domain.SoftDelete `bson:",inline"`
	Name              string
	ApplicationId     *string
	Versions          map[string]*ProjectVersion
	Description       string
	SchemaLabel       string
	PolicyLabels      []string
	ProjectLabels     []string
	FreeLabels        []string

	CorrespondingSchema *schema.SpdxSchema
	UserManagement      UserManagementEntity
	Token               []*Token

	Status            ProjectStatusType `json:"status"`
	LastFileUploads   []time.Time
	DocumentMeta      DisclosureDocumentMeta
	CustomerMeta      CustomerMeta
	NoticeContactMeta NoticeContactMeta
	IsGroup           bool
	Children          []string // children guids
	Parent            string   // parent guid
	ParentName        string   // parent name
	Documents         []*pdocument.PDocument
	SupplierExtraData SupplierExtraData
	ApprovableSPDX    approvable.ApprovableSPDX
	IsNoFoss          bool
	ApplicationMeta   ApplicationMeta
	/*
	   This means, that the project was not complete loaded!
	   Do not save this project otherwise, we will delete the not loaded data in the database.
	*/
	Optimized       bool `bson:"-" json:"-"`
	CustomIds       []ProjectCustomId
	HasApproval     bool
	HasChildren     bool
	HasSBOMToRetain bool
}

func (p *Project) GetDocumentByFileNameWithIndex(fileName string, index int) *pdocument.PDocument {
	for _, document := range p.Documents {
		// nil safety
		if document == nil {
			continue
		}
		dbDocFileName := pdocument.GetFileNameWithIndex(document.Type, document.ApprovalId, pdocument.LangStrToTag(document.Lang), int(*document.VersionIndex))
		if dbDocFileName == fileName && *document.VersionIndex == pdocument.DocumentVersion(index) {
			return document
		}
	}
	return nil
}

func (p *Project) GetFilePathBaseProject() string {
	return RemoveDoubleSlash(strings.Join([]string{
		conf.Config.Server.GetUploadPath(),
		strconv.Itoa(p.Created.Year()),
		fmt.Sprintf("%02d", int(p.Created.Month())),
		p.Key,
	}, "/"))
}

func RemoveDoubleSlash(source string) string {
	return strings.Replace(source, "//", "/", -1)
}

func (p *Project) GetFilePathBaseProjectAndVersion(versionKey string) string {
	return RemoveDoubleSlash(strings.Join([]string{
		p.GetFilePathBaseProject(),
		"versions",
		versionKey,
	}, "/"))
}

func (p *Project) GetFilePathSbom(spdxKey string, versionKey string) string {
	return RemoveDoubleSlash(strings.Join([]string{
		p.GetFilePathBaseProjectAndVersion(versionKey),
		"sbom",
		spdxKey,
	}, "/"))
}

func (p *Project) GetFilePathDocumentForProject(docFileName string) string {
	return RemoveDoubleSlash(strings.Join([]string{
		p.GetFilePathBaseProject(),
		"documents",
		docFileName,
	}, "/"))
}

type Department struct {
	domain.ChildEntity `bson:",inline"`
	DeptId             string
	ParentDeptId       string
	ValidFrom          *time.Time
	DescriptionEnglish string
	OrgAbbreviation    string
	Skz                string
	CompanyCode        string
	CompanyName        string
}

type DisclosureDocumentMeta struct {
	domain.ChildEntity `bson:",inline"`
	SupplierName       string
	SupplierAddress    string
	SupplierNr         string
	SupplierDeptId     string
	SupplierDept       *Department
}

type NoticeContactMeta struct {
	domain.ChildEntity `bson:",inline"`
	Address            string
	Email              string
}

type CustomerMeta struct {
	domain.ChildEntity `bson:",inline"`

	// TODO: Delete after migration
	Dept *Department

	DeptId string

	Address string
	// TODO: move into something named "ApproverPresets"
	FRI string `json:"fRI"`
	SRI string `json:"sRI"`
}

type UserManagementEntity struct {
	domain.ChildEntity `bson:",inline"`
	Users              []*ProjectMemberEntity
}

type ProjectMemberEntity struct {
	domain.ChildEntity `bson:",inline"`
	UserId             string
	UserType           UserType
	Comment            string
	IsResponsible      bool
	Subscriptions      Subscriptions
}

type Subscriptions struct {
	Spdx          bool
	OverallReview bool
}

type UserType string

const (
	OWNER               UserType = "Owner"
	SUPPLIER            UserType = "Supplier"
	VIEWER              UserType = "Viewer"
	PROJECT_RESPONSIBLE UserType = "ProjectResponsible"
)

type ProjectStatusType string

const (
	Ready      = "ready"
	Active     = "active"
	Deprecated = "deprecated"
	Approved   = "approved"
	Archived   = "archived"
)

type ProjectVersionStatusType string

const (
	PV_New                    = "new"
	PV_Unreviewed             = "unreviewed"
	PV_Rejected               = "rejected"
	PV_Approved               = "approved"
	PV_Freezed                = "freezed"
	PV_Acceptable             = "acceptable"
	PV_AcceptableAfterChanges = "acceptable_after_changes"
	PV_NotAcceptable          = "not_acceptable"
	PV_Audited                = "audited"
)

func ParseStatusType(state string) (valid bool, result ProjectVersionStatusType) {
	switch state {
	case PV_New:
		valid, result = true, PV_New
	case PV_Unreviewed:
		valid, result = true, PV_Unreviewed
	case PV_Rejected:
		valid, result = true, PV_Rejected
	case PV_Approved:
		valid, result = true, PV_Approved
	case PV_Freezed:
		valid, result = true, PV_Freezed
	case PV_Acceptable:
		valid, result = true, PV_Acceptable
	case PV_AcceptableAfterChanges:
		valid, result = true, PV_AcceptableAfterChanges
	case PV_Audited:
		valid, result = true, PV_Audited
	case PV_NotAcceptable:
		valid, result = true, PV_NotAcceptable
	}
	return
}

func CreateNewProject(projectData ProjectRequestDto) *Project {
	var ownersArray []*ProjectMemberEntity
	ownersArray = append(ownersArray, &ProjectMemberEntity{
		ChildEntity:   domain.NewChildEntity(),
		UserId:        projectData.Owner,
		UserType:      OWNER,
		IsResponsible: true,
	})

	newProject := &Project{
		RootEntity: domain.NewRootEntity(),
		Versions:   map[string]*ProjectVersion{},
		Name:       projectData.Name,
		UserManagement: UserManagementEntity{
			Users: ownersArray,
		},
		Status: Ready,
	}
	if projectData.ProjectSettings == nil {
		return newProject
	}
	newProject.DocumentMeta = projectData.ProjectSettings.DocumentMeta
	if projectData.ProjectSettings.DocumentMeta.SupplierDept != nil {
		newProject.DocumentMeta.SupplierDeptId = projectData.ProjectSettings.DocumentMeta.SupplierDept.DeptId
	}
	newProject.CustomerMeta = CustomerMeta{
		ChildEntity: domain.ChildEntity{},
		Address:     projectData.ProjectSettings.CustomerMeta.Address,
	}
	if projectData.ProjectSettings.CustomerMeta.Dept != nil {
		newProject.CustomerMeta.DeptId = projectData.ProjectSettings.CustomerMeta.Dept.DeptId
	}
	newProject.IsNoFoss = projectData.ProjectSettings.NoFossProject

	if projectData.Parent != "" {
		newProject.Parent = projectData.Parent
		newProject.ParentName = projectData.ParentName
	}
	return newProject
}

func (project *Project) GetVersion(versionKey string) *ProjectVersion {
	projectVersion := project.Versions[versionKey]
	version := projectVersion
	if version == nil {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorVersionMissing), "")
	}
	if version.Deleted {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorVersionDeleted), "")
	}
	return version
}

func (entity *Project) AddDocument(docs ...*pdocument.PDocument) {
	entity.Documents = append(entity.Documents, docs...)
}

func (project *Project) GenerateAndAddToken(tokenData Token) Token {
	tokenSecret := uuid.New().String()
	tokenForDatabase := Token{
		ChildEntity: domain.NewChildEntity(),
		Company:     tokenData.Company,
		Description: tokenData.Description,
		Expiry:      tokenData.Expiry,
		TokenSecret: hash.GetSha256Hash([]byte(tokenSecret)),
		Status:      ACTIVE,
	}
	tokenForWeb := tokenForDatabase
	tokenForWeb.TokenSecret = tokenSecret
	project.Token = append(project.Token, &tokenForDatabase)
	project.Updated = time.Now()

	return tokenForWeb
}

func (project *Project) RevokeToken(key string) {
	tokenFound := false
	for k, v := range project.Token {
		if v.Key == key {
			project.Token[k].Status = REVOKED
			tokenFound = true
		}
	}
	if !tokenFound {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorTokenNotFoundForKey, key), "")
	}
}

func (project *Project) RenewToken(key string) *Token {
	for _, oldToken := range project.Token {
		if oldToken.Key == key {
			oldToken.Status = REVOKED
			newToken := project.GenerateAndAddToken(*oldToken)
			return &newToken
		}
	}
	exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorTokenNotFoundForKey, key), "")
	return nil
}

func (project *Project) AddUser(userData ProjectMemberRequestDto) {
	project.CheckIfUserAlreadyExists(userData.TargetUser)

	entity := &ProjectMemberEntity{
		ChildEntity:   domain.NewChildEntity(),
		UserId:        userData.TargetUser,
		UserType:      userData.UserType,
		Comment:       userData.Comment,
		IsResponsible: userData.IsResponsible,
	}
	project.UserManagement.Users = append(project.UserManagement.Users, entity)
	project.Updated = time.Now()
}

func (project *Project) ProjectResponsible() *ProjectMemberEntity {
	for _, u := range project.UserManagement.Users {
		if !u.IsResponsible {
			continue
		}
		return u
	}
	return nil
}

func (project *Project) IsResponsible(id string) bool {
	for _, u := range project.UserManagement.Users {
		if u.UserId != id {
			continue
		}
		if u.IsResponsible {
			return true
		}
	}
	return false
}

func (project *Project) UUID() string {
	return project.Key
}

func (project *Project) UpdateProjectMember(userId string, userData ProjectMemberRequestDto) {
	if userData.TargetUser != userId {
		project.CheckIfUserAlreadyExists(userData.TargetUser)
	}

	var targetUser *ProjectMemberEntity
	ownersCount := 0
	for _, u := range project.UserManagement.Users {
		if u.UserType == OWNER {
			ownersCount++
		}
		if u.UserId != userId {
			continue
		} else {
			targetUser = u
		}
	}

	if targetUser == nil {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorUserNotFound, userData.TargetUser), "")
	}

	if targetUser.UserType == OWNER &&
		ownersCount == 1 && (userData.UserType != OWNER ||
		userData.TargetUser != userId) {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorProjectLastOwnerCanNotChanged, userData.TargetUser), "")
	}

	targetUser.UserId = userData.TargetUser
	targetUser.UserType = userData.UserType
	targetUser.Comment = userData.Comment
	targetUser.IsResponsible = userData.IsResponsible

	project.Updated = time.Now()
}

func (project *Project) DeleteUser(userId string) {
	userFound := false
	ownersCount := 0
	cleanedUsers := make([]*ProjectMemberEntity, 0)
	for _, value := range project.UserManagement.Users {
		if value.UserType == OWNER {
			ownersCount++
		}
		if value.UserId == userId {
			userFound = true
			if value.UserType == OWNER {
				ownersCount--
			}
		} else {
			cleanedUsers = append(cleanedUsers, value)
		}
	}

	if !userFound {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorDbNotFound), "user "+userId+" not found in DB")
	}

	// check if this is an attempt to delete a last owner, what is not allowed
	if ownersCount == 0 {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorProjectDeleteLastMember), " "+userId+" not found in DB")
	}

	project.UserManagement.Users = cleanedUsers
	project.Updated = time.Now()
}

func (project *Project) UpdateProjectData(newData ProjectRequestDto, connectorSet bool) {
	project.Name = newData.Name
	project.SchemaLabel = newData.SchemaLabel
	project.PolicyLabels = newData.PolicyLabels
	project.ProjectLabels = newData.ProjectLabels
	project.FreeLabels = newData.FreeLabels
	project.Description = newData.Description
	project.IsGroup = newData.IsGroup
	if newData.ProjectSettings != nil {
		project.SupplierExtraData = newData.ProjectSettings.SupplierExtraData
	}
	project.IsNoFoss = newData.IsNoFoss
	project.ApplicationMeta = newData.ApplicationMeta.ToEntity()
	project.Updated = time.Now()
}

// Contains tells whether a contains x.
func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func (project *Project) PrepareUpdateChild(newData ProjectRequestDto, isOwnerOrDummyCb func(string) bool) ([]string, []string, []string) {
	newChildren := make([]string, 0)
	added := make([]string, 0)
	for _, n := range newData.Children {
		if Contains(project.Children, n) {
			newChildren = append(newChildren, n)
			continue
		}
		if !isOwnerOrDummyCb(n) {
			continue
		}
		newChildren = append(newChildren, n)
		added = append(added, n)
	}
	removed := make([]string, 0)
	for _, o := range project.Children {
		if Contains(newData.Children, o) {
			continue
		}
		removed = append(removed, o)
	}

	project.Children = newChildren
	return added, removed, newChildren
}

func (project *Project) SetDocumentMeta(documentMeta DisclosureDocumentMetaDto) {
	project.DocumentMeta.SupplierName = documentMeta.SupplierName
	project.DocumentMeta.SupplierNr = documentMeta.SupplierNr
	project.DocumentMeta.SupplierAddress = documentMeta.SupplierAddress
	project.DocumentMeta.Updated = time.Now()
	project.Updated = time.Now()
}

func (project *Project) SetProjectSettings(settings *ProjectSettingsRequest) {
	project.DocumentMeta.SupplierName = settings.DocumentMeta.SupplierName
	project.DocumentMeta.SupplierNr = settings.DocumentMeta.SupplierNr
	project.DocumentMeta.SupplierAddress = settings.DocumentMeta.SupplierAddress
	if settings.DocumentMeta.SupplierDept != nil {
		project.DocumentMeta.SupplierDeptId = settings.DocumentMeta.SupplierDept.DeptId
	} else {
		project.DocumentMeta.SupplierDeptId = ""
	}
	project.DocumentMeta.Updated = time.Now()

	project.CustomerMeta.Address = settings.CustomerMeta.Address
	if settings.CustomerMeta.Dept != nil {
		project.CustomerMeta.DeptId = settings.CustomerMeta.Dept.DeptId
	} else {
		project.CustomerMeta.DeptId = ""
	}
	project.CustomerMeta.FRI = settings.CustomerMeta.FRI
	project.CustomerMeta.SRI = settings.CustomerMeta.SRI
	project.CustomerMeta.Updated = time.Now()

	project.NoticeContactMeta.Address = settings.NoticeContactMeta.Address
	project.NoticeContactMeta.Email = settings.NoticeContactMeta.Email
	project.NoticeContactMeta.Updated = time.Now()

	project.SupplierExtraData.FRI = settings.SupplierExtraData.FRI
	project.SupplierExtraData.SRI = settings.SupplierExtraData.SRI
	project.SupplierExtraData.External = settings.SupplierExtraData.External
	project.Updated = time.Now()
	project.IsNoFoss = settings.NoFossProject

	project.CustomIds = []ProjectCustomId{}
	for _, c := range settings.CustomIds {
		project.CustomIds = append(project.CustomIds, ProjectCustomId{
			ChildEntity: domain.NewChildEntity(),
			TechnicalId: c.TechnicalId,
			Value:       c.Value,
		})
	}
}

func (project *Project) FindVersionByName(name string) *ProjectVersion {
	for _, version := range project.Versions {
		if version.Deleted {
			continue
		}
		if version.Name == name {
			return version
		}
	}
	return nil
}

func (project *Project) CreateNewProjectVersionIfNameNotUsed(name string, description string) string {
	if project.Versions == nil {
		project.Versions = make(map[string]*ProjectVersion)
	}
	if project.FindVersionByName(name) != nil {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorVersionAlreadyExist, name), "")
	}
	projectVersion := ProjectVersion{
		ChildEntity: domain.NewChildEntity(),
		Name:        name,
		Description: description,
		Updated:     time.Now(),
		Created:     time.Now(),
		Status:      PV_New,
	}
	project.Versions[projectVersion.Key] = &projectVersion
	project.Updated = time.Now()

	return projectVersion.Key
}

func (project *Project) DeleteVersion(versionKey string) {
	project.Updated = time.Now()
	project.Versions[versionKey].Deleted = true
}

func (project *Project) UpdateVersion(version *ProjectVersion, newName string, description string) {
	version.Name = newName
	version.Description = description
	version.Updated = time.Now()
	project.Updated = time.Now()
}

func (project *Project) GetVersions() []ProjectVersion {
	this := *project
	projectVersions := make([]ProjectVersion, 0)
	for _, v := range this.Versions {
		if v.Deleted {
			continue
		}
		projectVersions = append(projectVersions, *v)
	}
	return projectVersions
}

func (project *Project) GetStatus() ProjectStatusType {
	if len(project.Status) == 0 {
		return Ready
	}
	return project.Status
}

func (project *Project) ExpireTokens() bool {
	statusChanged := false
	for i := 0; i < len(project.Token); i++ {
		token := project.Token[i]
		if token.Status == ACTIVE && token.IsExpired() {
			project.Token[i].Status = EXPIRED
			statusChanged = true
		}
	}
	return statusChanged
}

func (project *Project) GetActiveTokenByKey(key string) *Token {
	for _, value := range project.Token {
		if value.Key == key {
			if value.IsExpired() {
				continue
			}
			if value.Status == ACTIVE {
				return value
			}
		}
	}
	return nil
}

func (project *Project) GetActiveToken(token string) *Token {
	tokenHash := hash.GetSha256Hash([]byte(token))
	for _, value := range project.Token {
		if value.TokenSecret == tokenHash {
			if value.IsExpired() {
				value.Status = EXPIRED
			}
			if value.Status == ACTIVE {
				return value
			}
		}
	}
	return nil
}

func (project *Project) FindCorrespondingSchema(activeSchemas []*schema.SpdxSchema) *schema.SpdxSchema {
	if len(activeSchemas) == 0 {
		return nil
	}
	for _, activeSchema := range activeSchemas {
		if activeSchema.MatchesProjectLabel(project.SchemaLabel) {
			return activeSchema
		}
	}
	return nil
}

func (project *Project) CheckIfUserAlreadyExists(targetUser string) {
	for _, value := range project.UserManagement.Users {
		if value.UserId == targetUser {
			exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorProjectMemberAlreadyExist, targetUser), "")
		}
	}
}

func (project *Project) CheckIfUserAlreadyExistsSoft(targetUser string) bool {
	for _, value := range project.UserManagement.Users {
		if value.UserId == targetUser {
			return true
		}
	}
	return false
}

func (project *Project) GetMember(userId string) *ProjectMemberEntity {
	for _, value := range project.UserManagement.Users {
		if value.UserId == userId {
			return value
		}
	}
	return nil
}

func (project *Project) GetToken(tokenKey string) (*Token, error) {
	for _, token := range project.Token {
		if token.Key == tokenKey {
			return token, nil
		}
	}
	return nil, errors.New("no token found for key " + tokenKey)
}

func (project *Project) UpdateStatusToActive() {
	if len(project.Status) == 0 || project.Status == Ready {
		project.Status = Active
	}
}

func (project *Project) IsDeprecated() bool {
	return project.Status == Deprecated
}

func (project *Project) DeprecateProject() {
	if project.IsDeprecated() {
		return
	}

	project.Status = Deprecated

	project.RevokeAllTokens()
	project.SetAllUsersToViewerAndUnsubscribe()
}

func (project *Project) RevokeAllTokens() {
	for _, t := range project.Token {
		t.Status = REVOKED
	}
}

func (project *Project) SetAllUsersToViewerAndUnsubscribe() {
	for _, u := range project.UserManagement.Users {
		if u.IsResponsible {
			u.IsResponsible = false
		}
		if u.Subscriptions.Spdx {
			u.Subscriptions.Spdx = false
		}
		if u.Subscriptions.OverallReview {
			u.Subscriptions.OverallReview = false
		}
		u.UserType = VIEWER
	}
}

func (project *Project) HasParent() bool {
	return project.Parent != ""
}

func (c *CustomerMeta) Diff(d CustomerMeta) bool {
	if c.SRI != d.SRI || c.FRI != d.FRI || c.Address != d.Address {
		return true
	}
	return (c.DeptId != d.DeptId)
}

func (d *DisclosureDocumentMeta) Diff(di DisclosureDocumentMeta) bool {
	return d.SupplierName != di.SupplierName || d.SupplierAddress != di.SupplierAddress ||
		d.SupplierDeptId != di.SupplierDeptId || d.SupplierNr != di.SupplierNr
}

func (entity *Project) GetDocuments() []*pdocument.PDocument {
	if entity.Documents == nil {
		entity.Documents = make([]*pdocument.PDocument, 0)
	}
	return entity.Documents
}

type ApplicationMeta struct {
	Id           string
	SecondaryId  string
	Name         string
	ExternalLink string
}

func (entity ApplicationMeta) ToDto() ApplicationMetaDto {
	return ApplicationMetaDto{
		Id:           entity.Id,
		SecondaryId:  entity.SecondaryId,
		Name:         entity.Name,
		ExternalLink: entity.ExternalLink,
	}
}
