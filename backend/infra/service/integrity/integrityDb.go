// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package integrity

import (
	"strconv"
	"strings"
	"time"

	"mercedes-benz.ghe.com/foss/disuko/infra/repository/database"

	"mercedes-benz.ghe.com/foss/disuko/domain/integrity"
	"mercedes-benz.ghe.com/foss/disuko/domain/project"
	"mercedes-benz.ghe.com/foss/disuko/domain/project/pdocument"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	project2 "mercedes-benz.ghe.com/foss/disuko/infra/repository/project"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/sbomlist"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

func checkIfAllFilesFromDbExistOnS3(requestSession *logy.RequestSession, projectRepository project2.IProjectRepository,
	sbomListRepo sbomlist.ISbomListRepository,
	fixIt bool,
	state *integrity.DbIntegrityResult,
) {
	projectCount := projectRepository.CountAllWithDeleted(requestSession)
	logy.Infof(requestSession, "Start analyse projects, project count: "+strconv.Itoa(projectCount))

	offset := 0
	limit := 100
	for {
		errorMsg := ""
		var exp *exception.Exception

		qc := database.New().SetLimit(offset, limit)
		qbRawRes := projectRepository.Query(requestSession, qc)
		var qbRes []string
		for _, e := range qbRawRes {
			qbRes = append(qbRes, e.GetKey())
		}
		for _, projectKey := range qbRes {
			state.CountProjects++
			projectEntity := projectRepository.FindByKeyWithDeleted(requestSession, projectKey, false)
			for _, versionEntity := range projectEntity.Versions {
				// SBOMs
				sbomListEntity := sbomListRepo.FindByKeyWithDeleted(requestSession, versionEntity.Key, false)
				if sbomListEntity == nil {
					continue
				}
				for position, spdxFile := range sbomListEntity.SpdxFileHistory {
					state.CountFilesOnDB++
					sbomS3Path := projectEntity.GetFilePathSbom(spdxFile.Key, versionEntity.Key)
					remoteHash := ""
					errorMsg, exp, remoteHash = checkIfFileIntact(requestSession, sbomS3Path, spdxFile.Hash, fixIt)
					if fixIt {
						if len(remoteHash) > 0 {
							spdxFile.Hash = remoteHash
							errorMsg += " - Adding/Replace HASH in DB Entry"
						} else if len(errorMsg) > 0 {
							sbomListEntity.SpdxFileHistory = RemoveElementFromSlice(sbomListEntity.SpdxFileHistory, position)
							errorMsg += " - Remove file DB Entry"
						}
						sbomListRepo.UpdateWithoutTimestamp(requestSession, sbomListEntity)
					}

					if len(errorMsg) > 0 {
						addMissingS3File(errorMsg, versionEntity, projectEntity, sbomS3Path, spdxFile.MetaInfo.Name,
							spdxFile.Uploaded, spdxFile.Key, exp, state, fixIt)
					}
				}

			}

			// Documents -> Project level
			for position, document := range projectEntity.Documents {
				state.CountFilesOnDB++
				fileName := pdocument.GetFileNameWithIndex(document.Type, document.ApprovalId, pdocument.LangStrToTag(document.Lang), int(*document.VersionIndex))
				fileS3Path := projectEntity.GetFilePathDocumentForProject(fileName)

				remoteHash := ""
				errorMsg, exp, remoteHash = checkIfFileIntact(requestSession, fileS3Path, document.Hash, fixIt)

				if fixIt {
					if len(remoteHash) > 0 {
						document.Hash = remoteHash
						errorMsg += " - Adding/Replace HASH in DB Entry"
					} else if len(errorMsg) > 0 {
						document.Hash = remoteHash
						projectEntity.Documents = RemoveElementFromSlice(projectEntity.Documents, position)
						errorMsg += " - Remove file DB Entry"
					}

					projectRepository.UpdateWithoutTimestamp(requestSession, projectEntity)
				}

				if len(errorMsg) > 0 {
					addMissingS3File(errorMsg, nil, projectEntity, fileS3Path, fileName,
						&document.Created, document.Key, exp, state, fixIt)
				}

			}
		}

		logy.Infof(requestSession, "analysed projects: "+strconv.Itoa(state.CountProjects)+" / "+strconv.Itoa(projectCount))

		if len(qbRes) < limit {
			break
		}
		offset += limit
	}

	if state.CountProjects != projectCount {
		addError2(requestSession, state, "Different analysed project count: "+strconv.Itoa(state.CountProjects)+" / "+strconv.Itoa(projectCount))
	}

	logy.Infof(requestSession, "End analyse projects: "+strconv.Itoa(state.CountProjects)+" / "+strconv.Itoa(projectCount))
}

func addMissingS3File(errorMsg string, versionEntity *project.ProjectVersion, projectEntity *project.Project,
	sbomS3Path string, orgFileName string, uploaded *time.Time, fileKeyInDb string, exp *exception.Exception,
	state *integrity.DbIntegrityResult, fixed bool,
) {
	errorMsg = strings.Trim(errorMsg, "")
	versionName := ""
	versionUuid := ""
	versionIsDeleted := false
	if versionEntity != nil {
		versionName = versionEntity.Name
		versionUuid = versionEntity.Key
		versionIsDeleted = versionEntity.Deleted
	}
	missingFileOnS3 := &integrity.MissingS3File{
		ProjectUuid:      projectEntity.UUID(),
		ProjectName:      projectEntity.Name,
		ProjectIsDeleted: projectEntity.Deleted,
		VersionName:      versionName,
		VersionUuid:      versionUuid,
		VersionIsDeleted: versionIsDeleted,
		S3FileName:       sbomS3Path,
		OrgFileName:      orgFileName,
		Upload:           uploaded,
		Message:          errorMsg,
		Fixed:            fixed,
		DbEntityKey:      fileKeyInDb,
	}
	if exp != nil {
		missingFileOnS3.ErrorMessage = exp.ErrorMessage
		missingFileOnS3.ErrorCode = exp.ErrorCode
		missingFileOnS3.ErrorRaw = exp.ErrorRaw
	}
	state.MissingFileOnS3 = append(state.MissingFileOnS3, missingFileOnS3)
}
