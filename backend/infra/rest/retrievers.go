// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/eclipse-disuko/disuko/domain/project"
	sbomlist2 "github.com/eclipse-disuko/disuko/domain/project/sbomlist"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	project2 "github.com/eclipse-disuko/disuko/infra/repository/project"
	"github.com/eclipse-disuko/disuko/infra/repository/sbomlist"
	userRepo "github.com/eclipse-disuko/disuko/infra/repository/user"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/go-chi/chi/v5"
)

func retrieveProjectAndVersionFromPublicRequest(rs *logy.RequestSession, prRepo project2.IProjectRepository, userRepo userRepo.IUsersRepository, r *http.Request) (*project.Project, *project.ProjectVersion, string) {
	currentProject, origin := retrieveProjectFromPublicRequest(rs, prRepo, userRepo, r, true, true)

	versionEscaped := chi.URLParam(r, "version")
	versionName, err := url.QueryUnescape(versionEscaped)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ParamVersionWrong))

	versionNameLen := len(versionName)
	if versionNameLen <= 0 {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ParamVersionEmpty), "")
	}
	if versionNameLen > 80 {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ParamVersionToLong), "")
	}
	version := currentProject.FindVersionByName(versionName)
	if version == nil {
		exception.ThrowExceptionClient404Message3(message.GetI18N(message.FindVersion))
	}
	if version.Deleted {
		exception.ThrowExceptionClient404Message3(message.GetI18N(message.ErrorVersionDeleted))
	}
	return currentProject, version, origin
}

func retrieveProjectFromPublicRequest(rs *logy.RequestSession, prRepo project2.IProjectRepository, userRepo userRepo.IUsersRepository, r *http.Request, withVersions bool, denyDeprecated bool) (*project.Project, string) {
	prID := extractProjectKeyFromRequest(r)
	pr := prRepo.FindByKey(rs, prID, !withVersions)
	if pr == nil {
		exception.ThrowExceptionClient404Message(message.GetI18N(message.ErrorDbRead, project2.ProjectCollectionName), "project not found: "+prID)
	}
	if denyDeprecated && pr.IsDeprecated() {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.DeprecatedProjectError), "")
	}

	expired := pr.ExpireTokens()
	if expired {
		newToken := pr.Token
		if !withVersions {
			full := prRepo.FindByKey(rs, pr.Key, false)
			full.Token = newToken
			prRepo.Update(rs, full)
		} else {
			prRepo.Update(rs, pr)
		}
	}

	accessCookie, err := r.Cookie("access")
	if err == nil {
		return pr, projectAccessAuth(rs, prRepo, pr, accessCookie)
	}
	authHeader := r.Header.Get("Authorization")

	s := strings.Split(authHeader, " ")
	if len(s) != 2 {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.DiscoTokenUnauthorized, "Invalid disco token"), "Malformed token provided")
	}
	if s[0] == DiscoBearer {
		return pr, projectTokenAuth(rs, prRepo, pr, assertTokenUUID(authHeader, DiscoBearer)).Origin()
	}
	if s[0] == Bearer {
		return pr, patAuth(rs, pr, prRepo, userRepo, s[1])
	}
	exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.DiscoTokenUnauthorized, "Invalid disco token"), "Malformed token provided")
	return nil, ""
}

func retrieveProject2(repo project2.IProjectRepository, r *http.Request, withVersions bool) (*project.Project, *logy.RequestSession) {
	requestSession := logy.GetRequestSession(r)
	projectUUID := extractProjectKeyFromRequest(r)

	currentProject := repo.FindByKey(requestSession, projectUUID, !withVersions)
	if currentProject == nil {
		exception.ThrowExceptionClient404Message(message.GetI18N(message.ErrorDbRead, project2.ProjectCollectionName), "project not found: "+projectUUID)
	}

	return currentProject, requestSession
}

func retrieveProjectAndVersion2(repo project2.IProjectRepository, r *http.Request) (*project.Project, *project.ProjectVersion, *logy.RequestSession) {
	projectUUID := extractProjectKeyFromRequest(r)
	versionKey := extractVersionKeyFromRequest(r)

	pr, requestSession := retrieveProject2(repo, r, true)
	version, ok := pr.Versions[versionKey]
	if !ok {
		exception.ThrowExceptionClient404Message(message.GetI18N(message.ErrorDbRead, project2.ProjectCollectionName), "version not found: "+versionKey+" in project: "+projectUUID)
	}

	version = pr.Versions[versionKey]
	if version.Deleted {
		exception.ThrowExceptionClient404Message(message.GetI18N(message.ErrorVersionMissing, project2.ProjectCollectionName), "version not found: "+versionKey+" in project: "+projectUUID)
	}

	return pr, version, requestSession
}

func retrieveSbomListAndLatestFile(requestSession *logy.RequestSession, repo sbomlist.ISbomListRepository, key string) (*sbomlist2.SbomList, *project.SpdxFileBase) {
	sbomList := repo.FindByKey(requestSession, key, false)
	if sbomList == nil || len(sbomList.SpdxFileHistory) == 0 {
		return nil, nil
	}
	return sbomList, sbomList.SpdxFileHistory[len(sbomList.SpdxFileHistory)-1]
}

func retrieveSbomListAndFile(requestSession *logy.RequestSession, repo sbomlist.ISbomListRepository, key, fileKey string) (*sbomlist2.SbomList, *project.SpdxFileBase) {
	sbomList := repo.FindByKey(requestSession, key, false)
	if sbomList == nil || len(sbomList.SpdxFileHistory) == 0 {
		return nil, nil
	}
	var res *project.SpdxFileBase
	for _, spdx := range sbomList.SpdxFileHistory {
		if spdx.Key == fileKey {
			return sbomList, spdx
		}
	}
	return sbomList, res
}

func (s *SPDXHandler) retrieveSbomListAndFile(requestSession *logy.RequestSession, key, fileKey string) (*sbomlist2.SbomList, *project.SpdxFileBase) {
	return retrieveSbomListAndFile(requestSession, s.SbomListRepository, key, fileKey)
}

func (s *SPDXHandler) retrieveSbomListAndLatestFile(requestSession *logy.RequestSession, key string) (*sbomlist2.SbomList, *project.SpdxFileBase) {
	return retrieveSbomListAndLatestFile(requestSession, s.SbomListRepository, key)
}

func (p *ProjectHandler) RetrieveSbomListAndFile(requestSession *logy.RequestSession, key, fileKey string) (*sbomlist2.SbomList, *project.SpdxFileBase) {
	return retrieveSbomListAndFile(requestSession, p.SbomListRepository, key, fileKey)
}

func (p *ProjectHandler) retrieveSbomListAndLatestFile(requestSession *logy.RequestSession, key string) (*sbomlist2.SbomList, *project.SpdxFileBase) {
	return retrieveSbomListAndLatestFile(requestSession, p.SbomListRepository, key)
}

func (p *ProjectHandler) retrieveProjectAndVersionFromPublicRequest(rs *logy.RequestSession, r *http.Request) (*project.Project, *project.ProjectVersion, string) {
	return retrieveProjectAndVersionFromPublicRequest(rs, p.ProjectRepository, p.UserRepository, r)
}

func (h *ProjectHandler) retrieveProjectFromPublicRequest(rs *logy.RequestSession, r *http.Request, withVersions bool) (*project.Project, string) {
	return retrieveProjectFromPublicRequest(rs, h.ProjectRepository, h.UserRepository, r, withVersions, true)
}

func (h *PolicyRulesHandler) retrieveProjectFromPublicRequest(rs *logy.RequestSession, r *http.Request, withVersions bool) (*project.Project, string) {
	return retrieveProjectFromPublicRequest(rs, h.ProjectRepository, h.UserRepository, r, withVersions, true)
}

func (s *SPDXHandler) retrieveProjectAndVersionFromPublicRequest(rs *logy.RequestSession, r *http.Request) (*project.Project, *project.ProjectVersion, string) {
	return retrieveProjectAndVersionFromPublicRequest(rs, s.ProjectRepository, s.UserRepository, r)
}

func (p *ProjectHandler) retrieveProjectAndVersion2(r *http.Request) (*project.Project, *project.ProjectVersion, *logy.RequestSession) {
	return retrieveProjectAndVersion2(p.ProjectRepository, r)
}

func (s *SPDXHandler) retrieveProjectAndVersion2(r *http.Request) (*project.Project, *project.ProjectVersion, *logy.RequestSession) {
	return retrieveProjectAndVersion2(s.ProjectRepository, r)
}

func (s *SPDXHandler) retrieveProject2(r *http.Request, withVersions bool) (*project.Project, *logy.RequestSession) {
	return retrieveProject2(s.ProjectRepository, r, withVersions)
}

func (p *ProjectHandler) retrieveProject2(r *http.Request, withVersions bool) (*project.Project, *logy.RequestSession) {
	return retrieveProject2(p.ProjectRepository, r, withVersions)
}
