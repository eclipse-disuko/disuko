// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package report

type Report struct {
	CustomIDNames []string
	Projects      []Project
}

type Project struct {
	Group                                     string   `column:"group"`
	Guid                                      string   `column:"guid"`
	Name                                      string   `column:"name"`
	Link                                      string   `column:"link"`
	Status                                    string   `column:"status"`
	GroupName                                 string   `column:"group name"`
	GroupId                                   string   `column:"group id"`
	NonFoss                                   string   `column:"non foss"`
	IsDummy                                   string   `column:"dummy"`
	Updated                                   string   `column:"updated"`
	Created                                   string   `column:"created"`
	ApplicationName                           string   `column:"application name"`
	ApplicationId                             string   `column:"application id"`
	ApplicationSecondaryId                    string   `column:"secondary application id"`
	SchemaLabel                               string   `column:"schema label"`
	PolicyLabels                              string   `column:"policy labels"`
	Tags                                      string   `column:"tags"`
	ProjectLabels                             string   `column:"project labels"`
	Subscribers                               string   `column:"subscribers number"`
	OwnerCompanyName                          string   `column:"owner company name"`
	OwnerCompanyId                            string   `column:"owner company id"`
	OwnerDepartmentId                         string   `column:"owner department id"`
	OwnerDepartmentTitle                      string   `column:"owner department title"`
	OwnerDepartmentAbbreviation               string   `column:"owner department abbreviation"`
	SupplierCompanyName                       string   `column:"supplier company name"`
	SupplierCompanyId                         string   `column:"supplier company id"`
	SupplierDepartmentId                      string   `column:"supplier department id"`
	SupplierDepartmentTitle                   string   `column:"supplier department title"`
	SupplierDepartmentAbbreviation            string   `column:"supplier department abbreviation"`
	SupplierDepartmentExternal                string   `column:"supplier department external"`
	ProjectResponsibleUserid                  string   `column:"project responsible userid"`
	ProjectResponsibleEmail                   string   `column:"project responsible email"`
	ProjectResponsibleFullName                string   `column:"project responsible full name"`
	ProjectSboms                              string   `column:"project sboms"`
	LastUpload                                string   `column:"last upload"`
	ProjectTokens                             string   `column:"project tokens"`
	ActiveTokens                              string   `column:"active tokens"`
	ApprovedApprovalUpdated                   string   `column:"approved approval updated"`
	ApprovedApprovalTotal                     string   `column:"approved approval total"`
	ApprovedApprovalDenied                    string   `column:"approved approval denied"`
	ApprovedApprovalUnasserted                string   `column:"approved approval unasserted"`
	ApprovedApprovalWarned                    string   `column:"approved approval warned"`
	ApprovedApprovalQuestioned                string   `column:"approved approval questioned"`
	ApprovedApprovalAllowed                   string   `column:"approved approval allowed"`
	ApprovedApprovalLink                      string   `column:"approved approval link"`
	LatestApprovalStatus                      string   `column:"latest approval status"`
	LatestApprovalStatusDetails               string   `column:"latest approval status details"`
	LatestApprovalUpdated                     string   `column:"latest approval updated"`
	LatestApprovalTotal                       string   `column:"latest approval total"`
	LatestApprovalDenied                      string   `column:"latest approval denied"`
	LatestApprovalUnasserted                  string   `column:"latest approval unasserted"`
	LatestApprovalWarned                      string   `column:"latest approval warned"`
	LatestApprovalQuestioned                  string   `column:"latest approval questioned"`
	LatestApprovalAllowed                     string   `column:"latest approval allowed"`
	LatestApprovalWeakCopyLeft                string   `column:"latest approval weak copyleft"`
	LatestApprovalStrongCopyLeft              string   `column:"latest approval strong copyleft"`
	LatestApprovalNetworkCopyLeft             string   `column:"latest approval network copyleft"`
	LatestApprovalAndLicenseExp               string   `column:"latest approval AND license expression"`
	LatestApprovalOrLicenseExp                string   `column:"latest approval OR license expression"`
	LatestApprovalWithLicenseExp              string   `column:"latest approval WITH license expression"`
	LatestApprovalMixedLicenseExp             string   `column:"latest approval mixed AND-OR license expression"`
	LatestApprovalMassiveAndExp               string   `column:"latest approval massive AND license expression"`
	LatestApprovalMassiveOrExp                string   `column:"latest approval massive OR license expression"`
	LatestApprovalSourceCodeReference         string   `column:"latest approval channel code references"`
	LatestApprovalKeepSourceCode              string   `column:"latest approval obligation source available"`
	LatestApprovalGNU_CCSObligation           string   `column:"latest approval obligation CCS"`
	LatestApprovalNoFoss                      string   `column:"latest approval warning non-foss"`
	LatestExternalApprovalStatus              string   `column:"latest external approval status"`
	LatestExternalApprovalUpdated             string   `column:"latest external approval updated"`
	LatestExternalApprovalTotal               string   `column:"latest external approval total"`
	LatestExternalApprovalDenied              string   `column:"latest external approval denied"`
	LatestExternalApprovalUnasserted          string   `column:"latest external approval unasserted"`
	LatestExternalApprovalWarned              string   `column:"latest external approval warned"`
	LatestExternalApprovalQuestioned          string   `column:"latest external approval questioned"`
	LatestExternalApprovalAllowed             string   `column:"latest external approval allowed"`
	LatestExternalApprovalLink                string   `column:"latest external approval link"`
	LatestExternalApprovalWeakCopyLeft        string   `column:"latest external approval weak copyleft"`
	LatestExternalApprovalStrongCopyLeft      string   `column:"latest external approval strong copyleft"`
	LatestExternalApprovalNetworkCopyLeft     string   `column:"latest external approval network copyleft"`
	LatestExternalApprovalAndLicenseExp       string   `column:"latest external approval AND license expression"`
	LatestExternalApprovalOrLicenseExp        string   `column:"latest external approval OR license expression"`
	LatestExternalApprovalWithLicenseExp      string   `column:"latest external approval WITH license expression"`
	LatestExternalApprovalMixedLicenseExp     string   `column:"latest external approval mixed AND-OR license expression"`
	LatestExternalApprovalMassiveAndExp       string   `column:"latest external approval massive AND license expression"`
	LatestExternalApprovalMassiveOrExp        string   `column:"latest external approval massive OR license expression"`
	LatestExternalApprovalSourceCodeReference string   `column:"latest external approval channel code references"`
	LatestExternalApprovalKeepSourceCode      string   `column:"latest external approval obligation source available"`
	LatestExternalApprovalGNU_CCSObligation   string   `column:"latest external approval obligation CCS"`
	LatestExternalApprovalNoFoss              string   `column:"latest external approval warning non-foss"`
	LatestSbomTotal                           string   `column:"latest sbom total"`
	LatestSbomDenied                          string   `column:"latest sbom denied"`
	LatestSbomUnasserted                      string   `column:"latest sbom unasserted"`
	LatestSbomWarned                          string   `column:"latest sbom warned"`
	LatestSbomQuestioned                      string   `column:"latest sbom questioned"`
	LatestSbomAllowed                         string   `column:"latest sbom allowed"`
	LatestSbomAndLicenseExp                   string   `column:"latest sbom AND license expression"`
	LatestSbomOrLicenseExp                    string   `column:"latest sbom OR license expression"`
	LatestSbomWithLicenseExp                  string   `column:"latest sbom WITH license expression"`
	LatestSbomMixedLicenseExp                 string   `column:"latest sbom mixed AND-OR license expression"`
	LatestSbomMassiveAndExp                   string   `column:"latest sbom massive AND license expression"`
	LatestSbomMassiveOrExp                    string   `column:"latest sbom massive OR license expression"`
	LatestSbomSourceCodeReference             string   `column:"latest sbom channel code references"`
	LatestSbomWeakCopyLeft                    string   `column:"latest sbom weak copyleft"`
	LatestSbomStrongCopyLeft                  string   `column:"latest sbom strong copyleft"`
	LatestSbomNetworkCopyLeft                 string   `column:"latest sbom network copyleft"`
	LatestSbomKeepSourceCode                  string   `column:"latest sbom obligation source available"`
	LatestSbomGNU_CCSObligation               string   `column:"latest sbom obligation CCS"`
	LatestSbomNoFoss                          string   `column:"latest sbom warning non-foss"`
	NumberOfCodeReference                     string   `column:"number of code references"`
	ActiveLicenseDecisionRules                string   `column:"active license decision rules"`
	InactiveLicenseDecisionRules              string   `column:"inactive license decision rules"`
	ActivePolicyDecisionRules                 string   `column:"active policy decision rules"`
	InactivePolicyDecisionRules               string   `column:"inactive policy decision rules"`
	ActiveDeniedPolicyDecision                string   `column:"active denied policy decision"`
	InactiveDeniedPolicyDecision              string   `column:"inactive denied policy decision"`
	ManuallyLockedSBOM                        string   `column:"manually locked sbom"`
	TotalLockedSBOM                           string   `column:"total locked sbom"`
	LatestStatusReviewDate                    string   `column:"latest status review date"`
	LatestStatusReviewStatus                  string   `column:"latest status review status"`
	LatestE2ReviewDate                        string   `column:"latest Management review date"`
	LatestE2ReviewStatus                      string   `column:"latest Management review status"`
	LatestE2ReviewComment                     string   `column:"latest Management review comment"`
	CustomIDs                                 []string `column:"customids"`
}
