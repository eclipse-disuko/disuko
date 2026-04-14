// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"net/http"

	"github.com/eclipse-disuko/disuko/infra/rest"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"

	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/jwt"
	"github.com/eclipse-disuko/disuko/helper/message"
)

type RouteExtender func(r chi.Router)

func (s *Server) setupRoutes(extenders ...RouteExtender) {
	s.r.Group(func(r chi.Router) {
		r.Use(jwt.Authenticator)
		r.Post("/api/v1/refreshToken", s.handlers.auth.HandleRefreshToken) // test missing
		for _, extend := range extenders {
			extend(r)
		}
		r.Route("/api/v1", func(r chi.Router) {
			r.Get("/counts/dashboard", s.handlers.count.GetDashboardCountsHandler)
			r.Route("/analyse", func(r chi.Router) {
				r.Get("/files/start", s.handlers.analyseFiles.AnalyseFilesHandlerStart)
				r.Get("/files/stop", s.handlers.analyseFiles.AnalyseFilesHandlerStop)
				r.Get("/files/status", s.handlers.analyseFiles.GetResultHandler)
			})
			r.Get("/departments/find/{searchStr}", s.handlers.department.Find)
			r.Get("/customids", s.handlers.customid.List)
			r.Route("/disclosures", func(r chi.Router) {
				r.Use(AbortOnProd)
				r.Get("/", s.handlers.project.GetAllDisclosures)
			})
			r.Route("/projects", func(r chi.Router) {
				r.Post("/", s.handlers.project.ProjectPostHandler)
				r.Post("/{uuid}/clone", s.handlers.project.CloneProjectPostHandler)
				r.Get("/", s.handlers.project.ProjectGetAllHandler)
				r.Post("/search", s.handlers.project.ProjectSearchHandler)
				r.Get("/recent", s.handlers.project.ProjectRecentHandler)
				r.Get("/{uuid}/checklists", s.handlers.project.ProjectFindApplicableChecklists)
				r.Get("/{uuid}/possibleChildren", s.handlers.project.ProjectGetPossibleChildrenHandler)
				r.Get("/{uuid}/allSbomLists", s.handlers.project.ProjectGetAllSbomlists)
				r.Get("/{uuid}/allSBOM", s.handlers.project.ProjectGetAllSbom)
				r.Put("/{uuid}", s.handlers.project.ProjectUpdateHandler)
				r.Put("/{uuid}/deprecate", s.handlers.project.ProjectDeprecateHandler)
				r.Get("/{uuid}/jobs/{key}", s.handlers.project.JobGetOnetimeStatus) // test missing
				r.Route("/{uuid}/users", func(r chi.Router) {
					r.Get("/", s.handlers.project.ProjectMemberGetAllHandler)
					r.Post("/", s.handlers.project.ProjectMemberAddHandler)
					r.Route("/{userId}", func(r chi.Router) {
						r.Put("/", s.handlers.project.ProjectMemberUpdateHandler)
						r.Delete("/", s.handlers.project.ProjectMemberDeleteHandler)
						r.Get("/pendingApprovalOrReviewUsage", s.handlers.project.ProjectMemberGetUsageInPendingApprovalOrRequest)
					})
					r.Get("/profile/search/{searchFragment}", s.handlers.project.ProjectUserGetAllBySearchFragmentHandler)
				})
				r.Get("/{uuid}/schema", s.handlers.project.ProjectSchemaGetHandler) // test missing
				r.Get("/{uuid}", s.handlers.project.ProjectGetHandler)
				r.Route("/{uuid}/settings", func(r chi.Router) {
					r.Get("/", s.handlers.project.ProjectGetSettingsHandler)
					r.Put("/", s.handlers.project.ProjectUpdateSettingsHandler)
				})
				r.Route("/{uuid}/children", func(r chi.Router) {
					r.Get("/", s.handlers.project.ProjectGetChildrenHandler)
					r.Route("/users", func(r chi.Router) {
						r.Get("/", s.handlers.project.ProjectChildrenMemberGetAllHandler)
						r.Post("/", s.handlers.project.ProjectChildrenMemberAddHandler)
					})
				})
				r.Get("/{uuid}/approvalOrReviewUsage", s.handlers.project.ProjectGetUsageInApprovalOrReviewRequest)
				r.Delete("/{uuid}", s.handlers.project.ProjectDeleteHandler)
				r.Put("/{uuid}/approvableSPDX", s.handlers.project.ProjectUpdateTaskApprovableSPDX)
				r.Get("/{uuid}/approvableinfo", s.handlers.project.ProjectGetApprovableInfo)
				r.Route("/{uuid}/approval", func(r chi.Router) {
					r.Post("/create/{approvalType}", s.handlers.project.ProjectCreateApproval)
					r.Route("/{appId}", func(r chi.Router) {
						r.Put("/", s.handlers.project.ProjectUpdateApproval)
						r.Get("/", s.handlers.project.ProjectGetApproval)
						r.Post("/fillCustomer", s.handlers.project.ProjectFillCustomer)
						r.Get("/approver/{approver}", s.handlers.project.GetApproverUser)
					})
					r.Get("/list", s.handlers.project.ProjectGetApprovalList)
					r.Get("/vehiclechildren", s.handlers.project.ProjectCheckVehicleChildren)
					r.Get("/vehiclechildrenonly", s.handlers.project.GroupOnlyVehicleChildren)
				})

				r.Route("/{uuid}/policyrules", func(r chi.Router) {
					r.Get("/", s.handlers.project.ProjectGetPolicyRulesHandler)
					r.Get("/{id}", s.handlers.project.ProjectGetPolicyRulesByIdHandler)
				})

				r.Route("/{uuid}/tokens", func(r chi.Router) {
					r.Get("/", s.handlers.project.ProjectTokenGetAllHandler) // test missing
					r.Post("/", s.handlers.project.ProjectTokenAddHandler)
					r.Put("/{token}", s.handlers.project.ProjectTokenRenewHandler)     // test missing
					r.Delete("/{token}", s.handlers.project.ProjectTokenRevokeHandler) // test missing
				})
				r.Route("/{uuid}/audit", func(r chi.Router) {
					r.Get("/", s.handlers.project.ProjectTrailGetAllHandler) // test missing
				})
				r.Route("/{uuid}/documents", func(r chi.Router) {
					r.Get("/", s.handlers.project.ProjectDocumentsGetAllHandler) // test missing
					r.Get("/downloadTask/{taskId}/{fileType}/{lang}", s.handlers.project.DownloadDocumentByTaskHandler)
					r.Get("/downloadTask/{taskId}/{fileType}/{lang}/{docVersion}", s.handlers.project.DownloadDocumentByTaskHandler)
					r.Get("/downloadTask/{taskId}/{fileType}/", s.handlers.project.DownloadDocumentByTaskHandler)
				})
				r.Route("/{uuid}/spdxCompare", func(r chi.Router) {
					r.Get("/{versionOld}/{spdxOld}/{versionNew}/{spdxNew}", s.handlers.spdx.SPDXCompareHandler)
				})
				r.Route("/{uuid}/subscriptions", func(r chi.Router) {
					r.Put("/", s.handlers.project.SetSubscribedHandler)
				})
				r.Route("/{uuid}/templates/review", func(r chi.Router) {
					r.Get("/", s.handlers.project.GetReviewTemplates)    // test missing
					r.Get("/{id}", s.handlers.project.GetReviewTemplate) // test missing
				})
				r.Route("/{uuid}/decisions", func(r chi.Router) {
					r.Get("/", s.handlers.project.GetDecisions) // test missing
				})
				r.Route("/{uuid}/licenserules", func(r chi.Router) {
					r.Put("/{licenseRuleId}/cancel", s.handlers.project.CancelLicenseRule) // test missing
				})
				r.Route("/{uuid}/policyDecisions", func(r chi.Router) {
					r.Put("/{policyDecisionId}/cancel", s.handlers.project.CancelPolicyDecision) // test missing
				})
				r.Route("/{uuid}/versions", func(r chi.Router) {
					r.Get("/", s.handlers.project.ProjectVersionGetAllHandler)
					r.Post("/", s.handlers.project.ProjectVersionCreateHandler)
					r.Route("/{version}", func(r chi.Router) {
						r.Get("/", s.handlers.project.ProjectVersionGetHandler)
						r.Get("/audit", s.handlers.project.ProjectVersionTrailGetAllHandler)
						r.Get("/approvalOrReviewUsage", s.handlers.project.ProjectVersionGetUsageInApprovalOrReviewRequest)
						r.Get("/stats", s.handlers.project.GetGeneralVersionStatsHandler)
						r.Delete("/", s.handlers.project.ProjectVersionDeleteHandler)
						r.Put("/", s.handlers.project.ProjectVersionUpdateHandler)
						r.Post("/createLicenseRule", s.handlers.project.CreateLicenseRule)                // test missing
						r.Post("/createPolicyDecision", s.handlers.project.CreatePolicyDecision)          // test missing
						r.Post("/createBulkPolicyDecision", s.handlers.project.CreateBulkPolicyDecisions) // test missing

						r.Route("/externalsources", func(r chi.Router) {
							r.Get("/", s.handlers.project.GetAllExternalSourcesHandler)
							r.Post("/", s.handlers.project.ExternalSourceCreateHandler)
							r.Route("/{sourceId}", func(r chi.Router) {
								r.Delete("/", s.handlers.project.ExternalSourceDeleteHandler)
								r.Put("/", s.handlers.project.ExternalSourceUpdateHandler)
							})
						})

						r.Get("/components/{sbomUuid}", s.handlers.project.ProjectVersionComponentsForSbom) // test missing
						r.Post("/checklists/{sbomUuid}", s.handlers.project.ExecuteChecklistsHandler)       // test missing
						r.Get("/components/{sbomUuid}/stats", s.handlers.project.GetSBOMStatsHandler)
						r.Get("/components/{sbomUuid}/{searchFragment}", s.handlers.project.ProjectVersionComponentsBySearch) // test missing
						r.Get("/components/{sbomUuid}/licenses", s.handlers.project.SbomAllLicensesGetHandler)                // test missing
						r.Get("/component/{sbomUuid}/{spdxId}", s.handlers.project.ComponentDetailsForSbomGetHandler)
						r.Get("/component/{sbomUuid}/{spdxId}/reviewremarks", s.handlers.project.ComponentReviewRemarksGetHandler) // test missing
						r.Get("/component/{sbomUuid}/{spdxId}/licenses", s.handlers.project.ComponentLicensesGetHandler)           // test missing,
						r.Route("/notice/{sbomUuid}", func(r chi.Router) {
							r.Get("/text", s.handlers.spdx.ExportNoticeFileForSbomAsTextHandler) // test missing
							r.Get("/html", s.handlers.spdx.ExportNoticeFileForSbomAsHTMLHandler) // test missing
							r.Get("/json", s.handlers.spdx.ExportNoticeFileForSbomAsJSONHandler) // test missing
						})
						r.Route("/quality", func(r chi.Router) {
							r.Get("/scanremarks/{sbomUuid}", s.handlers.project.ProjectVersionScanRemarksForSbom) // test missing
							r.Get("/scanremarks/{sbomUuid}/download", s.handlers.project.DownloadScanRemarksForSbomCsvHandler)
							r.Get("/licenseremarks/{sbomUuid}", s.handlers.project.ProjectVersionLicenseRemarksForSbom)              // test missing
							r.Get("/licenseremarks/{sbomUuid}/download", s.handlers.project.DownloadLicenseRemarksForSbomCsvHandler) // test missing
							r.Route("/reviewremarks", func(r chi.Router) {
								r.Get("/", s.handlers.project.GetReviewRemarks)
								r.Post("/", s.handlers.project.CreateReviewRemark)
								r.Get("/download", s.handlers.project.DownloadReviewRemarksHandler)
								r.Post("/bulk-status", s.handlers.project.BulkSetReviewRemarkStatus)
								r.Put("/{remarkId}", s.handlers.project.EditReviewRemark)
								r.Post("/{remarkId}/comments", s.handlers.project.CommentReviewRemark)
								r.Put("/{remarkId}/status", s.handlers.project.SetReviewRemarkStatus)
							})
						})
						r.Get("/sbomhistory", s.handlers.project.ProjectVersionSPDXHistory) // test missing
						r.Route("/spdx", func(r chi.Router) {
							r.Post("/", s.handlers.spdx.SPDXUploadFileHandler)
							r.Route("/{spdxFileKey}", func(r chi.Router) {
								r.Get("/", s.handlers.spdx.DownloadSPDXHistoryFileHandler) // test missing
								r.Delete("/", s.handlers.spdx.SpdxDeleteFileHandler)       // test missing
								r.Put("/tag", s.handlers.spdx.SpdxTagUpdateHandler)
								r.Put("/toggleLock", s.handlers.spdx.SpdxToggleLockHandler)
							})
						})
						r.Route("/overallreview", func(r chi.Router) {
							r.Post("/", s.handlers.project.CreateOverallReview)
						})
					})
				})
				r.Route("/group-wizard", func(r chi.Router) {
					r.Post("/preview", s.handlers.project.WizardGroupPreviewHandler)
					r.Post("/", s.handlers.project.WizardCreateGroupHandler)
				})
				r.Route("/wizard", func(r chi.Router) {
					r.Post("/preview", s.handlers.project.WizardPreviewHandler)
					r.Post("/", s.handlers.project.WizardCreateHandler)
					r.Put("/{id}", s.handlers.project.WizardUpdateHandler)
					r.Get("/{id}", s.handlers.project.WizardGetHandler)
				})
			})
			r.Get("/capabilities", s.handlers.cap.GetCapabilities) // test missing
			r.Route("/profile", func(r chi.Router) {
				// for later use
				r.Get("/", s.handlers.user.GetProfileData)                                     // test missing
				r.Put("/{uuid}", s.handlers.user.UpdateHandlerForUser)                         // test missing
				r.Get("/search/{searchFragment}", s.handlers.user.Get5BySearchFragmentHandler) // test missing
				r.Route("/tasks", func(r chi.Router) {
					r.Get("/", s.handlers.user.GetTaskList)                   // test missing
					r.Get("/csv", s.handlers.user.GetTaskListCsv)             // test missing
					r.Get("/{taskId}", s.handlers.user.GetTask)               // test missing
					r.Put("/{taskId}/delegate", s.handlers.user.DelegateTask) // delegate task endpoint
				})
				r.Get("/projectroles", s.handlers.user.GetProjectRoles) // test missing
			})
			r.Route("/newsbox/items", func(r chi.Router) {
				r.Get("/", s.handlers.user.GetNewsBoxItems)      // test missing
				r.Put("/{uuid}", s.handlers.user.UpdateLastSeen) // test missing
			})
			r.Route("/admin", func(r chi.Router) {
				r.Get("/system/profile", s.handlers.statistic.GetSystemProfileStats) // test missing
				r.Post("/mail/send", s.handlers.mail.SendMail)
				r.Get("/counts/dashboard", s.handlers.count.GetDashboardCountsForAdminHandler)
				r.Route("/jobs", func(r chi.Router) {
					r.Get("/", s.handlers.job.JobGetAllHandler)
					r.Put("/{jobType}", s.handlers.job.JobTriggerRun)                    // test missing
					r.Put("/onetime/{key}", s.handlers.job.JobRerunOnetime)              // test missing
					r.Put("/{jobType}/config", s.handlers.job.SetConfig)                 // test missing
					r.Get("/latest/{jobType}", s.handlers.job.JobGetLatestByTypeHandler) // test missing
				})
				r.Route("/users", func(r chi.Router) {
					r.Get("/", s.handlers.user.GetAllHandler) // test missing
					// for later use
					r.Put("/", s.handlers.user.UpdateHandlerForAdmin) // test missing
					r.Get("/termsOfUseCurrentVersion", s.handlers.user.GetTermsOfUseCurrentVersionHandler)
					r.Post("/search", s.handlers.user.SearchHandlerForAdmin) // test missing
					r.Get("/delete-personal-data", s.handlers.user.DeletePersonalDataDryRunHandler)
					r.Get("/get-personal-details/{username}", s.handlers.user.GetPersonalDetailsHandler)
					r.Delete("/delete-personal-data/{entity}/{id}", s.handlers.user.DeletePersonalDataByEntityIdHandler)
					r.Delete("/delete-personal-data/{entity}", s.handlers.user.DeletePersonalDataByEntityHandler)
					r.Route("/deletion-audit", func(r chi.Router) {
						r.Get("/operation/{operationId}", s.handlers.user.GetDeletionAuditByOperationHandler)
						r.Get("/admin/{adminUser}", s.handlers.user.GetDeletionAuditByAdminHandler)
						r.Get("/user/{username}", s.handlers.user.GetDeletionAuditTrailHandler)
					})
					r.Route("/{uuid}", func(r chi.Router) {
						r.Get("/", s.handlers.user.GetByUuidHandler)                                           // test missing
						r.Get("/audit", s.handlers.user.GetAuditTrailHandler)                                  // test missing
						r.Put("/roles", s.handlers.user.UpdateUserRolesHandlerForAdmin)                        // test missing
						r.Get("/tokens", s.handlers.user.GetNewTokensHandlerForAdmin)                          // test missing
						r.Get("/tokensNonInternal", s.handlers.user.GetNewTokensForNonInternalHandlerForAdmin) // test missing
						r.Put("/active", s.handlers.user.EnableDisableHandlerForAdmin)                         // test missing
						r.Get("/projectroles", s.handlers.user.GetProjectRolesForAdmin)                        // test missing
						r.Get("/tasks", s.handlers.user.GetTaskListForAdmin)                                   // test missing
					})
					r.Get("/mails/{userId}", s.handlers.user.GetUserMailByIdHandler)
				})
				r.Route("/obligations", func(r chi.Router) {
					r.Get("/", s.handlers.obligation.GetAllHandler)                  // test missing
					r.Post("/", s.handlers.obligation.PostHandler)                   // test missing
					r.Put("/{id}", s.handlers.obligation.UpdateHandler)              // test missing
					r.Get("/{id}", s.handlers.obligation.GetByIdHandler)             // test missing
					r.Get("/{id}/audit", s.handlers.obligation.GetAuditTrailHandler) // test missing
					r.Delete("/{id}", s.handlers.obligation.DeleteHandler)           // test missing
					r.Get("/csv", s.handlers.obligation.CreateCSVHandler)            // test missing
				})
				r.Route("/checklist", func(r chi.Router) {
					r.Post("/", s.handlers.checklist.PostHandler)         // test missing
					r.Get("/", s.handlers.checklist.GetAllHandler)        // test missing
					r.Put("/{id}", s.handlers.checklist.UpdateHandler)    // test missing
					r.Delete("/{id}", s.handlers.checklist.DeleteHandler) // test missing
					r.Route("/{id}/items", func(r chi.Router) {
						r.Post("/", s.handlers.checklist.CreateItemHandler)
						r.Put("/{itemId}", s.handlers.checklist.UpdateItemHandler)
						r.Delete("/{itemId}", s.handlers.checklist.DeleteItemHandler)
					})
				})
				r.Route("/newsbox/items", func(r chi.Router) {
					r.Post("/", s.handlers.newsbox.PostHandler)         // test missing
					r.Get("/", s.handlers.newsbox.GetAllHandler)        // test missing
					r.Put("/{id}", s.handlers.newsbox.UpdateHandler)    // test missing
					r.Delete("/{id}", s.handlers.newsbox.DeleteHandler) // test missing
				})
				r.Route("/schemas", func(r chi.Router) {
					r.Get("/", s.handlers.schema.SchemaGetAllHandler)
					r.Post("/", s.handlers.schema.SchemaUploadHandler)
					r.Get("/{id}", s.handlers.schema.SchemaGetHandler)
					r.Get("/{id}/download", s.handlers.schema.SchemaDownloadHandler) // test missing
					r.Post("/{id}/activation", s.handlers.schema.SchemaActivateHandler)
					r.Get("/knowledgebase/export", s.handlers.export.ExportSchemaKnowledgeBase)
					r.Route("/knowledgebase", func(r chi.Router) {
						r.Use(AbortOnProd)
						r.Post("/import", s.handlers.export.ImportSchemaKnowledgeBase)
					})
				})
				r.Route("/policyrules", func(r chi.Router) {
					r.Post("/", s.handlers.policyRules.PolicyRulesCreateHandler)               // test missing
					r.Put("/{id}", s.handlers.policyRules.PolicyRulesUpdateHandler)            // test missing
					r.Put("/{id}/copy", s.handlers.policyRules.PolicyRulesCopyHandler)         // test missing
					r.Put("/{id}/deprecate", s.handlers.policyRules.DeprecateHandler)          // test missing
					r.Get("/{id}/audit", s.handlers.policyRules.PolicyRulesTrailGetAllHandler) // test missing
					r.Delete("/{id}", s.handlers.policyRules.PolicyRulesDeleteHandler)         // test missing
					r.Get("/csv", s.handlers.policyRules.CreateCSVHandler)
					r.Get("/{id}/changelog", s.handlers.policyRules.PolicyRulesChangeLogGetAllHandler) // test missing
				})

				r.Route("/licenses", func(r chi.Router) {
					r.Route("/spdx", func(r chi.Router) {
						r.Get("/count", s.handlers.licenses.GetSpdxLicensesCount) // test missing
						r.Get("/diffs", s.handlers.licenses.GetLicensesDiffs)     // test missing
						r.Delete("/{key}", s.handlers.licenses.DeleteSpdxHandler)
					})
					r.Put("/{key}", s.handlers.licenses.UpdateAcceptedChangesHandler) // test missing
					r.Get("/knowledgebase/export", s.handlers.export.ExportLicenseKnowledgeBase)
					r.Route("/knowledgebase", func(r chi.Router) {
						r.Use(AbortOnProd)
						r.Post("/import", s.handlers.export.ImportLicenseKnowledgeBase)
					})
				})
				r.Route("/labels", func(r chi.Router) {
					r.Get("/", s.handlers.label.GetLabels)
					r.Post("/", s.handlers.label.CreateLabel)        // test missing
					r.Put("/{id}", s.handlers.label.UpdateLabel)     // test missing
					r.Delete("/{id}", s.handlers.label.DeleteLabel)  // test missing
					r.Get("/csv", s.handlers.label.CreateCSVHandler) // test missing
				})
				r.Route("/templates/review", func(r chi.Router) {
					r.Get("/", s.handlers.template.GetReviewTemplates)
					r.Get("/{id}", s.handlers.template.GetReviewTemplate)
					r.Post("/", s.handlers.template.CreateReviewTemplate)       // test missing
					r.Put("/{id}", s.handlers.template.UpdateReviewTemplate)    // test missing
					r.Delete("/{id}", s.handlers.template.DeleteReviewTemplate) // test missing
					r.Get("/csv", s.handlers.template.CreateCSVHandler)         // test missing
				})
				r.Route("/utils/sampledata", func(r chi.Router) {
					r.Use(AbortOnProd)
					r.Post("/", s.handlers.sampleData.StartCreateSampleDataHandler)       // test missing
					r.Get("/", s.handlers.sampleData.GetStateCreateSampleDataHandler)     // test missing
					r.Delete("/", s.handlers.sampleData.StopStateCreateSampleDataHandler) // test missing
				})
				r.Get("/utils/stats", s.handlers.statistic.GetSystemStats)    // test missing
				r.Put("/utils/stats", s.handlers.statistic.UpdateSystemStats) // test missing

				r.Get("/utils/rightsProject", s.handlers.accessRights.ProjectAccessRightsGetAllHandler) // test missing
				r.Get("/utils/rights", s.handlers.accessRights.AccessRightsGetAllHandler)               // test missing

				r.Route("/notification", func(r chi.Router) {
					r.Get("/", s.handlers.notification.NotificationGetHandler)  // test missing
					r.Post("/", s.handlers.notification.NotificationSetHandler) // test missing
				})
				r.Route("/internaltoken", func(r chi.Router) {
					r.Get("/", s.handlers.basicauth.List)
					r.Post("/", s.handlers.basicauth.Create)
					r.Put("/{uuid}", s.handlers.basicauth.Renew)
					r.Delete("/{uuid}", s.handlers.basicauth.Revoke)
				})
				r.Route("/customid", func(r chi.Router) {
					r.Get("/", s.handlers.customid.List)
					r.Post("/", s.handlers.customid.Create)
					r.Put("/{uuid}", s.handlers.customid.Update)
					r.Delete("/{uuid}", s.handlers.customid.Delete)
					r.Get("/{uuid}/usage", s.handlers.customid.Usage)
				})
			})
			r.Route("/api", func(r chi.Router) {
				r.Route("/applications", func(r chi.Router) {
					r.Get("/search", s.handlers.application.SearchHandler) // test missing
				})
			})
			r.Route("/analytics", func(r chi.Router) {
				r.Post("/init", s.handlers.project.ReinitialiseAnalytics)               // test missing
				r.Post("/search", s.handlers.project.ProjectComponentSearchHandler)     // test missing
				r.Post("/components/search", s.handlers.project.ComponentSearchHandler) // test missing
				r.Post("/licenses/search", s.handlers.project.LicensesSearchHandler)    // test missing
				r.Get("/occurrences", s.handlers.analytics.LicenseOccurrences)          // test missing
				r.Get("/stats", s.handlers.analytics.Statistic)                         // test missing
				r.Get("/report", s.handlers.analytics.Report)                           // test missing
			})
			r.Route("/licenses", func(r chi.Router) {
				r.Post("/", s.handlers.licenses.LicensePostHandler) // test missing
				r.Post("/search", s.handlers.licenses.SearchHandler)
				r.Post("/lookup", s.handlers.licenses.LookupHandler)

				r.Put("/{key}", s.handlers.licenses.UpdateHandler) // test missing
				r.Delete("/{id}", s.handlers.licenses.DeleteHandler)
				r.Get("/", s.handlers.licenses.LicensesGetAllHandler) // test missing
				r.Get("/{id}", s.handlers.licenses.LicenseGetHandler)
				r.Get("/{id}/audit", s.handlers.licenses.LicenseTrailGetAllHandler)
				r.Get("/exists/{id}", s.handlers.licenses.LicenseHeadHandler)
				r.Get("/aliases/{alias}", s.handlers.licenses.AliasHeadHandler)
				r.Get("/name/{id}", s.handlers.licenses.LicenseNameHeadHandler)
				r.Get("/list/{ids}", s.handlers.licenses.LicensesGetHandler)       // test missing
				r.Post("/text/compare", s.handlers.licenses.LicenseCompareHandler) // test missing
				r.Get("/obligation/{obligationKey}/usagecount",
					s.handlers.licenses.GetCountOfLicencesUsingThisObligationHandler) // test missing
				r.Route("/policyrules/{licenceId}", func(r chi.Router) {
					r.Get("/", s.handlers.licenses.GetAllPolicyRulesAssignmentsForThisLicenceHandler)     // test missing
					r.Put("/", s.handlers.licenses.UpdateAllPolicyRulesAssignmentsForThisLicenceHandler)  // test missing
					r.Get("/usagecount", s.handlers.licenses.GetCountOfPolicyRuleUsingThisLicenceHandler) // test missing
				})
			})
			r.Route("/policyrules", func(r chi.Router) {
				r.Get("/", s.handlers.policyRules.PolicyRulesGetHandler)         // test missing
				r.Get("/{id}", s.handlers.policyRules.PolicyRulesGetByIdHandler) // test missing
				r.Get("/{id}/csv", s.handlers.policyRules.CreateRuleCSVHandler)  // test missing
			})
			r.Route("/announcements", func(r chi.Router) {
				r.Get("/", s.handlers.announcements.AnnouncementsGetAllHandler)
			})
			r.Route("/filtersets", func(r chi.Router) { // test missing
				r.Get("/tables/{tablename}", s.handlers.filterSet.FilterSetsGetByTableHandler)
				r.Get("/{id}", s.handlers.filterSet.FilterSetGetHandler)
				r.Post("/", s.handlers.filterSet.FilterSetPostHandler)
				r.Delete("/{id}", s.handlers.filterSet.FilterSetDeleteHandler)
				r.Put("/{id}", s.handlers.filterSet.FilterSetUpdateHandler)
			})
		})

		// if conf.Config.S3.IsRestApiEnabled {
		// 	apiRoute.Route("/s3", func(r chi.Router) {
		// 		r.Get("/loadtest/{sizeInMb}", s.handlers.s3.S3TriggerLoadTest)
		// 		r.Get("/loadtest/status", s.handlers.s3.S3TriggerLoadTestStatus)
		// 		r.Get("/", s.handlers.s3.S3GetAll)
		// 		r.Get("/{folder}", s.handlers.s3.S3GetAll)
		// 		r.Delete("/{filename}", s.handlers.s3.S3DeleteObject)
		// 		r.Post("/{filename}", s.handlers.s3.S3StoreFile)
		// 		r.Get("/text/{filename}", s.handlers.s3.S3GetTextFile)
		// 		r.Get("/meta/{filename}", s.handlers.s3.S3GetMetadataFile)
		// 	})
		// }
	})

	s.r.Group(func(r chi.Router) {
		r.Use(s.basicauthMW.Authenticator)
		r.Route("/api/internal", func(r chi.Router) {
			r.Get("/report", s.handlers.analytics.InternalReport)
			r.Get("/customlicenses", s.handlers.licenses.CustomLicenses)
		})
	})

	//	@title			FOSS Disclosure Portal
	//	@version		1.0.0
	//	@description	FOSS Disclosure Portal automates and digitizes the process for disclosure of the Free and Open Source Software components, which are included in products and applications. It aims at a more efficient, transparent and digital software supply chain, enabling software suppliers to deliver information on used open source via a technical interface in a standardized exchange format as Software Bill of Materials (SBOM).
	//	@description
	//	@description	SPDX-FileCopyrightText: 2023 Mercedes-Benz Tech Innovation GmbH
	//	@description	SPDX-License-Identifier: MIT
	//	@termsOfService	https://mb4.me/FOSS_Disclosure_Portal_ToU

	//	@schemes	https
	//	@BasePath	/disco/v1

	//	@securityDefinitions.apiKey	Bearer
	//	@in							header
	//	@name						Authorization
	//	@description				API key to authorize requests. DISCO xxxx.xxxx.xxxx.xxx

	//	@externalDocs.description	OpenAPI
	//	@externalDocs.url			https://swagger.io/resources/open-api/
	s.r.Group(func(r chi.Router) {
		r.Route("/api/public/auth", func(r chi.Router) {
			r.Post("/login", s.handlers.publicAuth.Login)
			r.Get("/refresh", s.handlers.publicAuth.Refresh)
			r.Get("/logout", s.handlers.publicAuth.Logout)
			r.Get("/info", s.handlers.publicAuth.Info)
		})
		r.Route("/api/public/v1", func(r chi.Router) {
			r.Get("/groups/{uuid}/children", s.handlers.project.ProjectGetChildrenExternHandler)
			r.Route("/projects/{uuid}", func(r chi.Router) {
				r.Get("/", s.handlers.project.ProjectGetExternHandler)
				r.Get("/status", s.handlers.project.ProjectStatusExternHandler)
				r.Get("/policyrules", s.handlers.policyRules.PolicyRulesGetExternHandler)
				r.Get("/schema", s.handlers.project.ProjectSchemaExternHandler)
				r.Get("/versions", s.handlers.project.ProjectVersionGetListExternHandler)
				r.Post("/sbomcheck", s.handlers.project.ProjectSPDXExternCheckOnDemand)
				r.Post("/versions", s.handlers.project.ProjectVersionExternCreateHandler)
				r.Post("/search", s.handlers.project.ProjectSbomSearchHandler)
				r.Route("/versions/{version}", func(r chi.Router) {
					r.Get("/", s.handlers.project.ProjectVersionGetExternHandler)
					r.Delete("/", s.handlers.project.ProjectVersionDeleteExternHandler)
					r.Get("/ccs", s.handlers.project.CCSGetListExternHandler)
					r.Get("/reviewremarks", s.handlers.project.ProjectVersionReviewRemarksExtern)
					r.Post("/reviewremarks/{reviewRemarkUuid}", s.handlers.project.ProjectVersionReviewRemarksCommentExtern)
					r.Post("/ccs", s.handlers.project.CCSCreateExternHandler)
					r.Post("/sboms", s.handlers.spdx.SPDXUploadFileExternHandler)
					// deprecated but still left for backward compatibility
					r.Post("/sbom", s.handlers.spdx.SPDXUploadFileExternHandler)
					r.Get("/sboms", s.handlers.project.ProjectVersionSPDXHistoryExtern)
					r.Route("/sboms/{sbomUuid}", func(r chi.Router) {
						r.Get("/", s.handlers.project.ProjectVersionSPDXMetaByIDExtern)
						r.Put("/tag", s.handlers.spdx.PublicSpdxTagUpdateHandler)
						r.Put("/lock", s.handlers.spdx.PublicSpdxLockHandler)
						r.Put("/unlock", s.handlers.spdx.PublicSpdxUnlockHandler)
						r.Get("/check", s.handlers.project.ProjectVersionSPDXExternCheck)
						r.Get("/status", s.handlers.spdx.PublicSpdxStatusHandler)
						r.Route("/notice", func(r chi.Router) {
							r.Get("/text", s.handlers.spdx.ExportTextNoticeExtern)
							r.Get("/html", s.handlers.spdx.ExportHTMLNoticeExtern)
							r.Get("/json", s.handlers.spdx.ExportJSONNoticeExtern)
						})
					})
				})
			})
			r.Get("/jwtTest", s.handlers.project.JwtTest)
		})

		// V2 endpoints
		r.Route("/api/public/v2", func(r chi.Router) {
			r.Get("/projects/{uuid}/versions", s.handlers.project.ProjectVersionGetListExternHandlerV2)
			r.Get("/projects/{uuid}/versions/{version}/reviewremarks", s.handlers.project.ProjectVersionReviewRemarksExternV2)
		})
	})

	// application token needed
	s.r.Group(func(r chi.Router) {
		r.Use(EnsureAppToken)
		r.Get("/api/v1/utils/stats/trigger", s.handlers.statistic.TriggerUpdateSystemStatsHandler) // test missing
	})

	// k8s health & shutdown handler
	s.r.Group(func(r chi.Router) {
		r.Get("/healthz", func(writer http.ResponseWriter, request *http.Request) {
			if isShuttingDown.Load() {
				writer.WriteHeader(http.StatusServiceUnavailable)
			} else {
				writer.WriteHeader(http.StatusOK)
			}
		})
		r.Get("/shutdown", func(writer http.ResponseWriter, request *http.Request) {
			logy.Infof(logy.GetRequestSession(request), "Received shutdown signal, stop serving requests")
			isShuttingDown.Store(true)
			writer.WriteHeader(http.StatusOK)
		})
	})

	// no user or token needed routes
	s.r.Group(func(r chi.Router) {
		r.Get("/api/v1/oauth/login", s.handlers.auth.HandleRedirectToIAM)        // redirect from client to here
		r.Get("/api/v1/oauth/logout", s.handlers.auth.HandleRedirectToIAMLogout) // redirect from client to here
		r.Get("/api/v1/login", s.handlers.auth.HandleRequestTokenFromCode)       // redirect from iam to here -> redirects to client than
	})

	s.r.Get("/maintenance", maintenanceHandler)
}

func EnsureAppToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		applicationToken := jwt.ExtractToken(r)
		if len(applicationToken) != len(conf.Config.Server.ApplicationToken) || applicationToken != conf.Config.Server.ApplicationToken {
			exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.ErrorAAR, "Invalid application token"), "Invalid application token")
		}

		// all fine
		next.ServeHTTP(w, r)
	})
}

func AbortOnProd(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		env := conf.Config.Server.Env
		prodEnv := conf.Config.Server.ProdEnv
		if conf.IsProdEnv() {
			exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorRunNotOnProd), "env="+env+" prodEnv="+prodEnv)
		}
		next.ServeHTTP(w, r)
	})
}

func maintenanceHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusServiceUnavailable)
	render.JSON(w, r, rest.SuccessResponse{
		Success: false,
		Message: "temporarily down for maintenance",
	})
}
