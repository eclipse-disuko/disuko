// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rights

import (
	"net/http"

	"github.com/eclipse-disuko/disuko/domain/oauth"
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/roles"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/go-chi/render"
)

const PublicApi = "PublicApi"

type AccessRightsHandler struct{}

func (handler *AccessRightsHandler) ProjectAccessRightsGetAllHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	logy.Infof(requestSession, "GetAccessRights - START")
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowTools.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}

	defaultInitialAccessRights := oauth.AccessAndRolesRights{}
	defaultInitialAccessRights.SetDefault()
	// defaultInitialAccessRights.SetForDomainAdmin()

	roles := []string{string(project.PROJECT_RESPONSIBLE), string(project.OWNER), string(project.SUPPLIER), string(project.VIEWER), PublicApi}
	response := make(map[string]*oauth.ProjectAccessRightsDto)
	for _, role := range roles {
		accessAndRolesRights := applyForProjectRole(defaultInitialAccessRights, role)
		if accessAndRolesRights != nil {
			response[role] = accessAndRolesRights.ToProjectAccessRightsDto()
		}
	}

	render.JSON(w, r, response)
}

func (handler *AccessRightsHandler) AccessRightsGetAllHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	logy.Infof(requestSession, "GetAccessRights - START")
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)

	if !rights.AllowTools.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}

	defaultInitialAccessRights := oauth.AccessAndRolesRights{}
	defaultInitialAccessRights.SetDefault()

	roles := []string{oauth.UserInternal, oauth.UserNonInternal, oauth.UserLicenseManager, oauth.UserPolicyManager, oauth.UserProjectAnalyst, oauth.UserDomainAdmin, oauth.UserApplicationAdmin, oauth.UserFOSSOffice}
	response := make(map[string]*oauth.AccessRightsDto)
	for _, role := range roles {
		accessAndRolesRights := applyForRole(defaultInitialAccessRights, role)
		if accessAndRolesRights != nil {
			response[role] = accessAndRolesRights.ToAccessRightsDto()
		}
	}

	render.JSON(w, r, response)
}

func applyForProjectRole(rights oauth.AccessAndRolesRights, role string) *oauth.AccessAndRolesRights {
	switch role {
	case string(project.OWNER):
		rights.SetOwnerRights()
		return &rights
	case string(project.SUPPLIER):
		rights.SetSupplierRights()
		return &rights
	case string(project.VIEWER):
		rights.SetViewerRights()
		return &rights
	case string(project.PROJECT_RESPONSIBLE):
		rights.SetProjectResponsibleRights()
		return &rights
	case PublicApi:
		rights.SetPublicApiRights()
		return &rights
	}
	return nil
}

// applyForRole applies Access Rights for each role. "*Admin"-Roles are Internals too, therefore set the rights for Internal and then for the corresponding role
func applyForRole(rights oauth.AccessAndRolesRights, role string) *oauth.AccessAndRolesRights {
	switch role {
	case oauth.UserInternal:
		rights.SetForInternal()
		return &rights
	case oauth.UserNonInternal:
		rights.SetForNonInternal()
		return &rights
	case oauth.UserLicenseManager:
		rights.SetForInternal()
		rights.SetForLicenseAdmin()
		return &rights
	case oauth.UserPolicyManager:
		rights.SetForInternal()
		rights.SetForPolicyAdmin()
		return &rights
	case oauth.UserProjectAnalyst:
		rights.SetForInternal()
		rights.SetForProjectAnalyst()
		return &rights
	case oauth.UserDomainAdmin:
		rights.SetForInternal()
		rights.SetForDomainAdmin()
		return &rights
	case oauth.UserApplicationAdmin:
		rights.SetForInternal()
		rights.SetForApplicationAdmin()
		return &rights
	case oauth.UserFOSSOffice:
		rights.SetForInternal()
		rights.SetForFOSSOffice()
		return &rights
	}
	return nil
}
