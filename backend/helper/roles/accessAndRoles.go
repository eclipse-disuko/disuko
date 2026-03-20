// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package roles

import (
	"net/http"
	"strings"

	"mercedes-benz.ghe.com/foss/disuko/conf"
	"mercedes-benz.ghe.com/foss/disuko/helper"

	"mercedes-benz.ghe.com/foss/disuko/domain/oauth"
	"mercedes-benz.ghe.com/foss/disuko/domain/project"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/jwt"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/labels"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

func GetUsernameFromRequest(requestSession *logy.RequestSession, r *http.Request) string {
	token := jwt.ExtractTokenMetadata(requestSession, r)
	return token.Username
}

func GetAccessAndRolesRightsFromRequest(requestSession *logy.RequestSession, r *http.Request) (string, *oauth.AccessAndRolesRights) {
	token := jwt.ExtractTokenMetadata(requestSession, r)
	checkDisabledUserAccess(token)
	rights := GetAccessAndRolesRightsFromClaim(*token)
	return token.Username, rights
}

func checkDisabledUserAccess(token *jwt.TokenData) {
	if !token.IsEnabled {
		exception.ThrowExceptionSendUserDisabledResponse()
	}
}

func GetAccessAndRolesRightsFromClaim(userInfo jwt.TokenData) *oauth.AccessAndRolesRights {
	// performance issue logy.Infof(requestSession, "accessAndRole::GetAccessAndRolesRightsFromClaim %v", userInfo.Hashed(requestSession))
	rightSet := oauth.AccessAndRolesRights{}
	rightSet.SetDefault()

	rightSet.IsInternal = userInfo.IsInternalEmployee || helper.Contains(userInfo.Username, conf.Config.InternalUsersAllowList)

	if rightSet.IsInternal {
		rightSet.SetForInternal()
	} else {
		rightSet.SetForNonInternal()
	}

	if strings.Contains(userInfo.Email, "@extaccount.com") {
		rightSet.SetForExternal()
	}

	groups := strings.Split(userInfo.Groups, jwt.GROUPS_TOKEN)

	rightSet.Groups = make([]string, 0)
	for _, group := range groups {
		if group == FossOfficeUser {
			rightSet.SetForFOSSOffice()
			rightSet.Groups = append(rightSet.Groups, oauth.UserFOSSOffice)
		}
		if group == LicenseAdmin {
			rightSet.SetForLicenseAdmin()
			rightSet.Groups = append(rightSet.Groups, oauth.UserLicenseManager)
		}
		if group == PolicyAdmin {
			rightSet.SetForPolicyAdmin()
			rightSet.Groups = append(rightSet.Groups, oauth.UserPolicyManager)
		}
		if group == ProjectAnalyst {
			rightSet.SetForProjectAnalyst()
			rightSet.Groups = append(rightSet.Groups, oauth.UserProjectAnalyst)
		}
		if group == ApplicationAdmin {
			rightSet.SetForApplicationAdmin()
			rightSet.Groups = append(rightSet.Groups, oauth.UserApplicationAdmin)
		}
		if group == DomainAdmin {
			rightSet.SetForDomainAdmin()
			rightSet.Groups = append(rightSet.Groups, oauth.UserDomainAdmin)
		}

	}

	// performance issue logy.Infof(requestSession, "current rights %v", rightSet.ToString())
	return &rightSet
}

func GetAndCheckProjectRights(requestSession *logy.RequestSession,
	request *http.Request, prj *project.Project, softly bool,
) (string, *oauth.AccessAndRolesRights) {
	var tokenData *jwt.TokenData = nil
	if softly {
		exception.TryCatchAndThrow(func() {
			tokenData = jwt.ExtractTokenMetadata(requestSession, request)
		}, func(exc exception.Exception) exception.Exception {
			exc.HttpErrorCode = http.StatusExpectationFailed
			return exc
		})
	} else {
		tokenData = jwt.ExtractTokenMetadata(requestSession, request)
	}

	checkDisabledUserAccess(tokenData)

	rights := GetAccessAndRolesRightsFromClaim(*tokenData)
	// Project scope rights
	setAccessRights(tokenData.Username, prj, rights)

	return tokenData.Username, rights
}

func FilterProjectsWithoutAccess(requestSession *logy.RequestSession, request *http.Request, projects *[]*project.Project, access project.UserType) error {
	tokenData := jwt.ExtractTokenMetadata(requestSession, request)
	rights := GetAccessAndRolesRightsFromClaim(*tokenData)
	if rights.AllowProject.Read {
		return nil
	}

	newList := make([]*project.Project, 0)
	for _, prj := range *projects {
		userFound := false
		for _, value := range prj.UserManagement.Users {
			if value.UserId != tokenData.Username {
				continue
			}
			if access != "" && value.UserType != access {
				break
			}
			userFound = true
			break
		}
		if userFound {
			newList = append(newList, prj)
		}
	}
	*projects = newList

	return nil
}

func GetProjectAccess(tokenData *jwt.TokenData, prj *project.Project) (*oauth.AccessAndRolesRights, error) {
	rights := GetAccessAndRolesRightsFromClaim(*tokenData)
	// Project scope rights
	setAccessRights(tokenData.Username, prj, rights)
	return rights, nil
}

// setAccessRights Sets adjusted for owner/supplier/viewer rights. Access to the specific operation must be checked at the caller place
func setAccessRights(userName string, prj *project.Project, rights *oauth.AccessAndRolesRights) {
	var userType project.UserType
	for _, projectUser := range prj.UserManagement.Users {
		if projectUser.UserId == userName {
			userType = projectUser.UserType
			break
		}
	}
	if len(userType) == 0 {
		if rights.AllowProject.Read || rights.AllowProjectGroup.Read || rights.AllowAllProjectUserManagement.Read || rights.AllowAllProjectTokenManagement.Read {
			// user not listed in project members but read access granted - do not overwrite, read access granted by Project Analyst only to view other projects or by Domain Admin to access User/Token Management over all projects independent on project role
			return
		} else {
			exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorAAR, message.ErrorKeyProjectAccess), "user has no access")
		}
	}
	switch userType {
	case project.OWNER:
		rights.Groups = append(rights.Groups, string(project.OWNER))
		rights.SetOwnerRights()
	case project.SUPPLIER:
		rights.Groups = append(rights.Groups, string(project.SUPPLIER))
		rights.SetSupplierRights()
	case project.VIEWER:
		rights.Groups = append(rights.Groups, string(project.VIEWER))
		rights.SetViewerRights()
	}
	// All things are fine: access rights are set, no error occured
}
func CanAccessVehicleProjectOperations(rights *oauth.AccessAndRolesRights, hasVehiclePlatformLabel bool) bool {
	return rights.IsFossOffice() && hasVehiclePlatformLabel
}

func CheckProjectTypeAccess(requestSession *logy.RequestSession, rights *oauth.AccessAndRolesRights, pr *project.Project, labelRepo labels.ILabelRepository, oauthAccessLevel oauth.AccessLevel) bool {
	for _, policyLabelKey := range pr.PolicyLabels {
		labelObj := labelRepo.FindByKey(requestSession, policyLabelKey, true)
		if labelObj == nil {
			continue
		}
		return !rights.HasProjectTypeAccess(labelObj.Name, oauthAccessLevel)
	}
	return false
}
