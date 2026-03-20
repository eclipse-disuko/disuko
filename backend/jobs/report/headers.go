// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package report

const (
	colGroup Col = iota
	colGuid
	colName
	colLink
	colStatus
	colGroupName
	colGroupId
	colNonFoss
	colIsDummy
	colUpdated
	colCreated
	colApplicationName
	colApplicationId
	colApplicationSecondaryId
	colSchemaLabel
	colPolicyLabels
	colTags
	colProjectLabels
	colSubscribers
	colOwnerCompanyName
	colOwnerCompanyId
	colOwnerDepartmentId
	colOwnerDepartmentTitle
	colOwnerDepartmentAbbreviation
	colSupplierCompanyName
	colSupplierCompanyId
	colSupplierDepartmentId
	colSupplierDepartmentTitle
	colSupplierDepartmentAbbreviation
	colSupplierDepartmentExternal
	colProjectResponsibleUserid
	colProjectResponsibleEmail
	colProjectResponsibleFullName
	colProjectSboms
	colLastUpload
	colProjectTokens
	colActiveTokens
	colApprovedApprovalUpdated
	colApprovedApprovalTotal
	colApprovedApprovalDenied
	colApprovedApprovalUnasserted
	colApprovedApprovalWarned
	colApprovedApprovalQuestioned
	colApprovedApprovalAllowed
	colApprovedApprovalLink
	colLatestApprovalStatus
	colLatestApprovalStatusDetails
	colLatestApprovalUpdated
	colLatestApprovalTotal
	colLatestApprovalDenied
	colLatestApprovalUnasserted
	colLatestApprovalWarned
	colLatestApprovalQuestioned
	colLatestApprovalAllowed
	colLatestApprovalWeakCopyLeft
	colLatestApprovalStrongCopyLeft
	colLatestApprovalNetworkCopyLeft
	colLatestApprovalAndLicenseExp
	colLatestApprovalOrLicenseExp
	colLatestApprovalWithLicenseExp
	colLatestApprovalMixedLicenseExp
	colLatestApprovalMassiveAndExp
	colLatestApprovalMassiveOrExp
	colLatestApprovalSourceCodeReference
	colLatestApprovalKeepSourceCode
	colLatestApprovalGNU_CCSObligation
	colLatestApprovalNoFoss
	colLatestExternalApprovalStatus
	colLatestExternalApprovalUpdated
	colLatestExternalApprovalTotal
	colLatestExternalApprovalDenied
	colLatestExternalApprovalUnasserted
	colLatestExternalApprovalWarned
	colLatestExternalApprovalQuestioned
	colLatestExternalApprovalAllowed
	colLatestExternalApprovalLink
	colLatestExternalApprovalWeakCopyLeft
	colLatestExternalApprovalStrongCopyLeft
	colLatestExternalApprovalNetworkCopyLeft
	colLatestExternalApprovalAndLicenseExp
	colLatestExternalApprovalOrLicenseExp
	colLatestExternalApprovalWithLicenseExp
	colLatestExternalApprovalMixedLicenseExp
	colLatestExternalApprovalMassiveAndExp
	colLatestExternalApprovalMassiveOrExp
	colLatestExternalApprovalSourceCodeReference
	colLatestExternalApprovalKeepSourceCode
	colLatestExternalApprovalGNU_CCSObligation
	colLatestExternalApprovalNoFoss
	colLatestSbomTotal
	colLatestSbomDenied
	colLatestSbomUnasserted
	colLatestSbomWarned
	colLatestSbomQuestioned
	colLatestSbomAllowed
	colLatestSbomAndLicenseExp
	colLatestSbomOrLicenseExp
	colLatestSbomWithLicenseExp
	colLatestSbomMixedLicenseExp
	colLatestSbomMassiveAndExp
	colLatestSbomMassiveOrExp
	colLatestSbomSourceCodeReference
	colLatestSbomWeakCopyLeft
	colLatestSbomStrongCopyLeft
	colLatestSbomNetworkCopyLeft
	colLatestSbomKeepSourceCode
	colLatestSbomGNU_CCSObligation
	colLatestSbomNoFoss
	colNumberOfCodeReference
	colActiveLicenseDecisionRules
	colInactiveLicenseDecisionRules
	colActivePolicyDecisionRules
	colInactivePolicyDecisionRules
	colActiveDeniedPolicyDecision
	colInactiveDeniedPolicyDecision
	colManuallyLockedSBOM
	colTotalLockedSBOM
	colLatestStatusReviewDate
	colLatestStatusReviewStatus
	colLatestE2ReviewDate
	colLatestE2ReviewStatus
	colLatestE2ReviewComment
)

var colHeaders = []string{
	colGroup:                                     "group",
	colGuid:                                      "guid",
	colName:                                      "name",
	colLink:                                      "link",
	colStatus:                                    "status",
	colGroupName:                                 "group name",
	colGroupId:                                   "group id",
	colNonFoss:                                   "non foss",
	colIsDummy:                                   "dummy",
	colUpdated:                                   "updated",
	colCreated:                                   "created",
	colApplicationName:                           "application name",
	colApplicationId:                             "application id",
	colApplicationSecondaryId:                    "secondary application id",
	colSchemaLabel:                               "schema label",
	colPolicyLabels:                              "policy labels",
	colTags:                                      "tags",
	colSubscribers:                               "subscribers number",
	colProjectLabels:                             "project labels",
	colOwnerCompanyName:                          "owner company name",
	colOwnerCompanyId:                            "owner company id",
	colOwnerDepartmentId:                         "owner department id",
	colOwnerDepartmentTitle:                      "owner department title",
	colOwnerDepartmentAbbreviation:               "owner department abbreviation",
	colSupplierCompanyName:                       "supplier company name",
	colSupplierCompanyId:                         "supplier company id",
	colSupplierDepartmentId:                      "supplier department id",
	colSupplierDepartmentTitle:                   "supplier department title",
	colSupplierDepartmentAbbreviation:            "supplier department abbreviation",
	colSupplierDepartmentExternal:                "supplier department external",
	colProjectResponsibleUserid:                  "project responsible userid",
	colProjectResponsibleEmail:                   "project responsible email",
	colProjectResponsibleFullName:                "project responsible full name",
	colProjectSboms:                              "project sboms",
	colLastUpload:                                "last upload",
	colProjectTokens:                             "project tokens",
	colActiveTokens:                              "active tokens",
	colApprovedApprovalUpdated:                   "approved approval updated",
	colApprovedApprovalTotal:                     "approved approval total",
	colApprovedApprovalDenied:                    "approved approval denied",
	colApprovedApprovalUnasserted:                "approved approval unasserted",
	colApprovedApprovalWarned:                    "approved approval warned",
	colApprovedApprovalQuestioned:                "approved approval questioned",
	colApprovedApprovalAllowed:                   "approved approval allowed",
	colApprovedApprovalLink:                      "approved approval link",
	colLatestApprovalStatus:                      "latest approval status",
	colLatestApprovalStatusDetails:               "latest approval status details",
	colLatestApprovalUpdated:                     "latest approval updated",
	colLatestApprovalTotal:                       "latest approval total",
	colLatestApprovalDenied:                      "latest approval denied",
	colLatestApprovalUnasserted:                  "latest approval unasserted",
	colLatestApprovalWarned:                      "latest approval warned",
	colLatestApprovalQuestioned:                  "latest approval questioned",
	colLatestApprovalAllowed:                     "latest approval allowed",
	colLatestApprovalWeakCopyLeft:                "latest approval weak copyleft",
	colLatestApprovalStrongCopyLeft:              "latest approval strong copyleft",
	colLatestApprovalNetworkCopyLeft:             "latest approval network copyleft",
	colLatestApprovalAndLicenseExp:               "latest approval AND license expression",
	colLatestApprovalOrLicenseExp:                "latest approval OR license expression",
	colLatestApprovalWithLicenseExp:              "latest approval WITH license expression",
	colLatestApprovalMixedLicenseExp:             "latest approval mixed AND-OR license expression",
	colLatestApprovalMassiveAndExp:               "latest approval massive AND license expression",
	colLatestApprovalMassiveOrExp:                "latest approval massive OR license expression",
	colLatestApprovalSourceCodeReference:         "latest approval channel code references",
	colLatestApprovalKeepSourceCode:              "latest approval obligation source available",
	colLatestApprovalGNU_CCSObligation:           "latest approval obligation CCS",
	colLatestApprovalNoFoss:                      "latest approval warning non-foss",
	colLatestExternalApprovalStatus:              "latest external approval status",
	colLatestExternalApprovalUpdated:             "latest external approval updated",
	colLatestExternalApprovalTotal:               "latest external approval total",
	colLatestExternalApprovalDenied:              "latest external approval denied",
	colLatestExternalApprovalUnasserted:          "latest external approval unasserted",
	colLatestExternalApprovalWarned:              "latest external approval warned",
	colLatestExternalApprovalQuestioned:          "latest external approval questioned",
	colLatestExternalApprovalAllowed:             "latest external approval allowed",
	colLatestExternalApprovalLink:                "latest external approval link",
	colLatestExternalApprovalWeakCopyLeft:        "latest external approval weak copyleft",
	colLatestExternalApprovalStrongCopyLeft:      "latest external approval strong copyleft",
	colLatestExternalApprovalNetworkCopyLeft:     "latest external approval network copyleft",
	colLatestExternalApprovalAndLicenseExp:       "latest external approval AND license expression",
	colLatestExternalApprovalOrLicenseExp:        "latest external approval OR license expression",
	colLatestExternalApprovalWithLicenseExp:      "latest external approval WITH license expression",
	colLatestExternalApprovalMixedLicenseExp:     "latest external approval mixed AND-OR license expression",
	colLatestExternalApprovalMassiveAndExp:       "latest external approval massive AND license expression",
	colLatestExternalApprovalMassiveOrExp:        "latest external approval massive OR license expression",
	colLatestExternalApprovalSourceCodeReference: "latest external approval channel code references",
	colLatestExternalApprovalKeepSourceCode:      "latest external approval obligation source available",
	colLatestExternalApprovalGNU_CCSObligation:   "latest external approval obligation CCS",
	colLatestExternalApprovalNoFoss:              "latest external approval warning non-foss",
	colLatestSbomTotal:                           "latest sbom total",
	colLatestSbomDenied:                          "latest sbom denied",
	colLatestSbomUnasserted:                      "latest sbom unasserted",
	colLatestSbomWarned:                          "latest sbom warned",
	colLatestSbomQuestioned:                      "latest sbom questioned",
	colLatestSbomAllowed:                         "latest sbom allowed",
	colLatestSbomAndLicenseExp:                   "latest sbom AND license expression",
	colLatestSbomOrLicenseExp:                    "latest sbom OR license expression",
	colLatestSbomWithLicenseExp:                  "latest sbom WITH license expression",
	colLatestSbomMixedLicenseExp:                 "latest sbom mixed AND-OR license expression",
	colLatestSbomMassiveAndExp:                   "latest sbom massive AND license expression",
	colLatestSbomMassiveOrExp:                    "latest sbom massive OR license expression",
	colLatestSbomSourceCodeReference:             "latest sbom channel code references",
	colLatestSbomWeakCopyLeft:                    "latest sbom weak copyleft",
	colLatestSbomStrongCopyLeft:                  "latest sbom strong copyleft",
	colLatestSbomNetworkCopyLeft:                 "latest sbom network copyleft",
	colLatestSbomKeepSourceCode:                  "latest sbom obligation source available",
	colLatestSbomGNU_CCSObligation:               "latest sbom obligation CCS",
	colLatestSbomNoFoss:                          "latest sbom warning non-foss",
	colNumberOfCodeReference:                     "number of code references",
	colActiveLicenseDecisionRules:                "active license decision rules",
	colInactiveLicenseDecisionRules:              "inactive license decision rules",
	colActivePolicyDecisionRules:                 "active policy decision rules",
	colInactivePolicyDecisionRules:               "inactive policy decision rules",
	colActiveDeniedPolicyDecision:                "active denied policy decision",
	colInactiveDeniedPolicyDecision:              "inactive denied policy decision",
	colManuallyLockedSBOM:                        "manually locked sbom",
	colTotalLockedSBOM:                           "total locked sbom",
	colLatestStatusReviewDate:                    "latest status review date",
	colLatestStatusReviewStatus:                  "latest status review status",
	colLatestE2ReviewDate:                        "latest Management review date",
	colLatestE2ReviewStatus:                      "latest Management review status",
	colLatestE2ReviewComment:                     "latest Management review comment",
}
