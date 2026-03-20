// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package count

import (
	"net/http"
	"time"

	rt "github.com/eclipse-disuko/disuko/infra/repository/reviewtemplates"

	"github.com/eclipse-disuko/disuko/domain/statistic"
	"github.com/eclipse-disuko/disuko/helper/roles"
	"github.com/eclipse-disuko/disuko/infra/repository/database"
	"github.com/eclipse-disuko/disuko/infra/repository/labels"
	"github.com/eclipse-disuko/disuko/infra/repository/license"
	newsboxRepo "github.com/eclipse-disuko/disuko/infra/repository/newsbox"
	"github.com/eclipse-disuko/disuko/infra/repository/obligation"
	"github.com/eclipse-disuko/disuko/infra/repository/policyrules"
	project2 "github.com/eclipse-disuko/disuko/infra/repository/project"
	"github.com/eclipse-disuko/disuko/infra/repository/schema"

	"github.com/eclipse-disuko/disuko/domain/user"
	userRepo "github.com/eclipse-disuko/disuko/infra/repository/user"
	"github.com/eclipse-disuko/disuko/logy"
)

func GetDashboardCounts(adminRequest bool, requestSession *logy.RequestSession, projectRepository project2.IProjectRepository,
	licenseRepository license.ILicensesRepository, policyRulesRepository policyrules.IPolicyRulesRepository, labelRepository labels.ILabelRepository, schemaRepository schema.ISchemaRepository, obligationRepository obligation.IObligationRepository, userRepository userRepo.IUsersRepository, reviewTemplateRepository rt.IReviewTemplateRepository, newsboxRepository newsboxRepo.IRepo, w http.ResponseWriter, r *http.Request) *statistic.DashboardCounts {
	counts := statistic.DashboardCounts{
		ProjectCount:        getProjectCount(adminRequest, requestSession, projectRepository, w, r),
		LicenseCount:        getLicenseCount(requestSession, licenseRepository, w, r),
		PolicyRuleCount:     getPolicyRuleCount(requestSession, policyRulesRepository, w, r),
		LabelCount:          getLabelCount(requestSession, labelRepository, w, r),
		SchemaCount:         getSchemaCount(requestSession, schemaRepository, w, r),
		ObligationCount:     getObligationCount(requestSession, obligationRepository, w, r),
		UserCount:           getUserCount(requestSession, userRepository, w, r),
		DisclosureCount:     getDisclosureCount(adminRequest, requestSession, projectRepository, w, r),
		ReviewTemplateCount: getReviewTemplateCount(requestSession, reviewTemplateRepository, w, r),
		ActiveJobCount:      getActiveJobCount(requestSession, userRepository, r),
		HasNewNewsboxItem:   getNewsboxState(requestSession, userRepository, newsboxRepository, w, r),
	}
	return &counts
}

func getObligationCount(requestSession *logy.RequestSession, repository obligation.IObligationRepository, w http.ResponseWriter, r *http.Request) int {
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowObligation.Read {
		return 0
	}

	qbRes := repository.CountAll(requestSession)
	return qbRes
}

func getSchemaCount(requestSession *logy.RequestSession, schemaRepository schema.ISchemaRepository, w http.ResponseWriter, r *http.Request) int {
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowSchema.Read {
		return 0
	}

	return schemaRepository.CountAll(requestSession)
}

func getLabelCount(requestSession *logy.RequestSession, labelRepository labels.ILabelRepository, w http.ResponseWriter, r *http.Request) int {
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowPolicy.Read {
		return 0
	}
	return labelRepository.CountAll(requestSession)
}

func getPolicyRuleCount(requestSession *logy.RequestSession, policyRulesRepository policyrules.IPolicyRulesRepository, w http.ResponseWriter, r *http.Request) int {
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowPolicy.Read {
		return 0
	}

	qbRes := policyRulesRepository.CountAll(requestSession)
	return qbRes
}

func getLicenseCount(requestSession *logy.RequestSession, licenseRepository license.ILicensesRepository, w http.ResponseWriter, r *http.Request) int {
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowLicense.Read {
		return 0
	}
	return licenseRepository.CountAll(requestSession)
}

func getProjectCount(adminRequest bool, requestSession *logy.RequestSession, projectRepository project2.IProjectRepository, w http.ResponseWriter, r *http.Request) int {
	if adminRequest {
		qbRes := projectRepository.CountAll(requestSession)
		return qbRes
		// not needed for admin view
		// err = roles.FilterProjectsWithoutAccess(requestSession, w, r, &projects)
	} else {
		username := roles.GetUsernameFromRequest(requestSession, r)
		return projectRepository.CountForUser(requestSession, username)
	}

}

func getDisclosureCount(adminRequest bool, requestSession *logy.RequestSession, projectRepository project2.IProjectRepository, w http.ResponseWriter, r *http.Request) int {
	if adminRequest {
		qbRes := projectRepository.CountAllGroups(requestSession)
		return qbRes
		// not needed for admin view
		// err = roles.FilterProjectsWithoutAccess(requestSession, w, r, &projects)
	} else {
		username := roles.GetUsernameFromRequest(requestSession, r)
		return projectRepository.CountGroupsForUser(requestSession, username)
	}
}

func getUserCount(requestSession *logy.RequestSession, userRepository userRepo.IUsersRepository, w http.ResponseWriter, r *http.Request) int {
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.AllowUsers.Create && rights.AllowUsers.Read && rights.AllowUsers.Update && rights.AllowUsers.Delete) {
		return 0
	}

	qbRes := userRepository.CountAll(requestSession)
	return qbRes
}

func getReviewTemplateCount(requestSession *logy.RequestSession, repository rt.IReviewTemplateRepository, w http.ResponseWriter, r *http.Request) int {
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowObligation.Read && !rights.AllowLicense.Read && !rights.AllowPolicy.Read {
		return 0
	}
	return repository.CountAll(requestSession)
}

func getActiveJobCount(requestSession *logy.RequestSession, userRepository userRepo.IUsersRepository, r *http.Request) int {
	name, _ := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	u := userRepository.FindByUserId(requestSession, name)
	if u == nil {
		return 0
	}
	res := 0

	for _, t := range u.Tasks {
		if t.Status == user.TaskActive {
			res++
		}
	}
	return res
}

func getNewsboxState(requestSession *logy.RequestSession, userRepository userRepo.IUsersRepository, newsboxRepository newsboxRepo.IRepo, w http.ResponseWriter, r *http.Request) bool {
	user, _ := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)

	currentUserProfile := userRepository.FindByUserId(requestSession, user)
	userLastSeenId := currentUserProfile.NewsboxLastSeenId

	qc := &database.QueryConfig{}
	qc.SetSort(database.SortConfig{
		database.SortAttribute{
			Name:  "Created",
			Order: database.ASC,
		},
	})
	all := newsboxRepository.Query(requestSession, qc)

	foundLastSeen := false
	hasNonExpiredItems := false
	hasNewItemsAfterLastSeen := false
	for _, item := range all {
		if item == nil {
			continue
		}

		if item.Key == userLastSeenId {
			foundLastSeen = true
			continue
		}

		isNonExpired := item.Expiry.IsZero() || item.Expiry.After(time.Now())
		if !isNonExpired {
			continue
		}
		hasNonExpiredItems = true

		if foundLastSeen {
			hasNewItemsAfterLastSeen = true
			break
		}
	}
	// If user has never seen any items, return true only if there are non-expired items
	if userLastSeenId == "" {
		return hasNonExpiredItems
	}

	// If last seen item was found, return true only if there are new items after it
	if foundLastSeen {
		return hasNewItemsAfterLastSeen
	}

	// If last seen item was not found but there are non-expired items, assume there are new items
	return hasNonExpiredItems
}
