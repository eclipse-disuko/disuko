// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package dummy

import (
	"fmt"
	"slices"
	"time"

	"golang.org/x/text/language"
	"mercedes-benz.ghe.com/foss/disuko/domain/approval"
	"mercedes-benz.ghe.com/foss/disuko/domain/job"
	"mercedes-benz.ghe.com/foss/disuko/domain/label"
	"mercedes-benz.ghe.com/foss/disuko/domain/project/pdocument"
	sbomlist2 "mercedes-benz.ghe.com/foss/disuko/domain/project/sbomlist"
	reviewremarks2 "mercedes-benz.ghe.com/foss/disuko/domain/reviewremarks"
	user2 "mercedes-benz.ghe.com/foss/disuko/domain/user"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/s3Helper"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/approvallist"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/auditloglist"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/base"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/database"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/labels"
	projectRepo "mercedes-benz.ghe.com/foss/disuko/infra/repository/project"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/reviewremarks"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/sbomlist"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/user"
	"mercedes-benz.ghe.com/foss/disuko/infra/service/cache"
	"mercedes-benz.ghe.com/foss/disuko/logy"
	"mercedes-benz.ghe.com/foss/disuko/scheduler"
)

type DeletionJob struct {
	projectRepository       projectRepo.IProjectRepository
	labelRepository         labels.ILabelRepository
	sbomListRepository      sbomlist.ISbomListRepository
	reviewRemarksRepository reviewremarks.IReviewRemarksRepository
	approvalListRepository  approvallist.IApprovalListRepository
	userRepository          user.IUsersRepository
	auditLogListRepository  auditloglist.IAuditLogListRepository
}

func InitDeletionJob(
	projectRepository projectRepo.IProjectRepository,
	labelRepository labels.ILabelRepository,
	sbomListRepository sbomlist.ISbomListRepository,
	reviewRemarksRepository reviewremarks.IReviewRemarksRepository,
	approvalListRepository approvallist.IApprovalListRepository,
	userRepository user.IUsersRepository,
	auditLogListRepository auditloglist.IAuditLogListRepository,
) *DeletionJob {
	return &DeletionJob{
		projectRepository:       projectRepository,
		labelRepository:         labelRepository,
		sbomListRepository:      sbomListRepository,
		reviewRemarksRepository: reviewRemarksRepository,
		approvalListRepository:  approvalListRepository,
		userRepository:          userRepository,
		auditLogListRepository:  auditLogListRepository,
	}
}

func (j *DeletionJob) Execute(rs *logy.RequestSession, info job.Job) scheduler.ExecutionResult {
	var log job.Log
	log.AddEntry(job.Info, "started")

	// Fetch projekt label "dummy"
	dummyLabel := j.labelRepository.FindByNameAndType(rs, label.DUMMY, label.PROJECT)
	if dummyLabel == nil {
		log.AddEntry(job.Error, "label 'dummy' not found")
		return scheduler.ExecutionResult{
			Success: false,
			Log:     log,
		}
	}

	// Search for all dummy projects older than 3 months
	cutoff := time.Now().UTC().AddDate(0, -3, 0).Format(time.RFC3339Nano)
	qc := database.New().SetMatcher(
		database.AndChain(
			database.AttributeMatcher(
				"Deleted",
				database.EQ,
				false,
			),
			database.AttributeMatcher(
				"Created",
				database.LT,
				cutoff,
			),
			database.ArrayElemMatcher(
				"ProjectLabels",
				database.EQ,
				dummyLabel.Key,
			),
		),
	)
	projects := j.projectRepository.Query(rs, qc)
	log.AddEntry(job.Info, "found %d dummy projects for deletion", len(projects))

	sbomLists := make([]*sbomlist2.SbomList, 0)
	reviewRemarkLists := make([]*reviewremarks2.ReviewRemarks, 0)

	approvalLists := make([]*approval.ApprovalList, 0)
	userTasksMap := make(map[string][]string)

	for _, prj := range projects {
		// Delete each SBOM and Cache for each channel of each project
		for _, v := range prj.GetVersions() {
			sbomList := j.sbomListRepository.FindByKey(rs, v.Key, false)
			reviewRemarkList := j.reviewRemarksRepository.FindByKey(rs, v.Key, false)

			if sbomList == nil && reviewRemarkList == nil {
				continue
			}

			if sbomList != nil {
				for _, sbom := range sbomList.SpdxFileHistory {
					sbomFile := prj.GetFilePathSbom(sbom.Key, v.Key)
					exception.TryCatchAndLog(rs, func() {
						s3Helper.DeleteFile(rs, sbomFile)
					})
					cacheFilePath := fmt.Sprintf(cache.CachePath, sbom.Key)
					exception.TryCatchAndLog(rs, func() {
						s3Helper.DeleteFile(rs, cacheFilePath)
					})
				}
				sbomLists = append(sbomLists, sbomList)
			}

			if reviewRemarkList != nil {
				reviewRemarkLists = append(reviewRemarkLists, reviewRemarkList)
			}
		}

		// Delete Approvals from corresponding Approval and User Tasks
		approvalList := j.approvalListRepository.FindByKey(rs, prj.Key, false)
		if approvalList != nil {
			for _, appr := range approvalList.Approvals {
				taskUuid := appr.Key
				creator := appr.Creator
				// Only Plausi and Internal Approval produce tasks
				if appr.Type == approval.TypePlausibility {
					addTaskForUser(userTasksMap, creator, taskUuid)
					addTaskForUser(userTasksMap, appr.Plausibility.Approver, taskUuid)
				} else if appr.Type == approval.TypeInternal {
					addTaskForUser(userTasksMap, creator, taskUuid)
					for _, approver := range appr.Internal.Approver {
						addTaskForUser(userTasksMap, approver, taskUuid)
					}
				}
			}
			approvalLists = append(approvalLists, approvalList)
		}

		// Delete document files of each Project
		// Only Internal and External Approvals produce document files
		// Process prj.Documents to resolve each belonging file and delete them
		// Take care about document's VersionIndex, Type, Language to collect the all files
		if prj.Documents != nil {
			for _, doc := range prj.Documents {
				var langTag *language.Tag
				if doc.Lang != "" {
					if t, err := language.Parse(doc.Lang); err == nil {
						langTag = &t
					}
				}
				versionIndex := doc.VersionIndex
				targetFileName := pdocument.GetFileNameWithIndex(doc.Type, doc.ApprovalId, langTag, int(*versionIndex))
				completeFileNameInS3 := prj.GetFilePathDocumentForProject(targetFileName)
				exception.TryCatchAndLog(rs, func() {
					s3Helper.DeleteFile(rs, completeFileNameInS3)
				})
			}
		}

		// Delete stil remaining project related files, if any
		projectPath := prj.GetFilePathBaseProject()
		filesCount := s3Helper.CountFiles(rs, projectPath).CntFiles
		projectFiles := s3Helper.ListObjects(rs, projectPath)
		if filesCount > 0 {
			// Log as an error for the first time to be notified in Grafana Dashboard
			msgFmt := "Found %d still remaining files for dummy project %s(%s) after deletion, deleting them now. Enhance deletion process to avoid possible data inconsistency"
			logy.Errorf(rs, msgFmt, filesCount, prj.Name, prj.Key)
			log.AddEntry(job.Error, msgFmt, filesCount, prj.Name, prj.Key)
			for file := range projectFiles {
				if len(file.Key) < 1 {
					// ignore ghost files, sometime happens on S3 Mock
					logy.Errorf(rs, "Found file ghost! ")
					continue
				}

				exception.TryCatchAndLog(rs, func() {
					s3Helper.DeleteFile(rs, file.Key)
				})
			}
		}

		j.projectRepository.DeleteHard(rs, prj.Key)
		log.AddEntry(job.Info, "dummy project '%s' (%s) deleted", prj.Name, prj.Key)
		j.auditLogListRepository.DeleteHard(rs, prj.Key)
		log.AddEntry(job.Info, "audit log for dummy project '%s' (%s) deleted", prj.Name, prj.Key)
	}

	// Perform deletion of things aggregated over all projects
	if len(sbomLists) > 0 {
		// Bulk deletion (hard) of all SBOM Lists for each version of each project
		s := j.sbomListRepository.StartSession(base.DeleteSession, 100)
		for _, sbomList := range sbomLists {
			s.AddEnt(sbomList)
		}
		s.EndSession()
	}
	if len(reviewRemarkLists) > 0 {
		// Bulk deletion (hard) of all Review Remarks Lists for each version of each project
		s := j.reviewRemarksRepository.StartSession(base.DeleteSession, 100)
		for _, reviewRemarks := range reviewRemarkLists {
			s.AddEnt(reviewRemarks)
		}
		s.EndSession()
	}
	// Deletion (hard) of Approval List of each project
	for _, approvalList := range approvalLists {
		j.approvalListRepository.DeleteHard(rs, approvalList.Key)
	}
	// Delete all tasks for users related to each project
	if len(userTasksMap) > 0 {
		users := make([]*user2.User, 0)
		for userId, taskUuids := range userTasksMap {
			usr := j.userRepository.FindByUserId(rs, userId)
			newUserTasks := make([]user2.Task, 0)
			for _, userTask := range usr.Tasks {
				if !slices.Contains(taskUuids, userTask.TargetGuid) {
					newUserTasks = append(newUserTasks, userTask)
				}
			}
			usr.Tasks = newUserTasks
			users = append(users, usr)
		}
		j.userRepository.UpdateList(rs, users)
	}

	log.AddEntry(job.Info, "finished")
	return scheduler.ExecutionResult{
		Success: true,
		Log:     log,
	}
}

func addTaskForUser(userTasksMap map[string][]string, user string, taskUuid string) {
	if user == "" {
		return
	}
	userTasksMap[user] = append(userTasksMap[user], taskUuid)
}
