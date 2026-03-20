// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package integrity

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/minio/minio-go/v7"
	"mercedes-benz.ghe.com/foss/disuko/conf"
	"mercedes-benz.ghe.com/foss/disuko/domain/integrity"
	"mercedes-benz.ghe.com/foss/disuko/domain/project"
	"mercedes-benz.ghe.com/foss/disuko/domain/project/pdocument"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
	"mercedes-benz.ghe.com/foss/disuko/helper/s3Helper"
	project2 "mercedes-benz.ghe.com/foss/disuko/infra/repository/project"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/sbomlist"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

const (
	versions_anchor  = "versions"
	documents_anchor = "documents"
	reports_anchor   = "reports"
)

func checkIfAllFilesFromS3ExistInTheDb(requestSession *logy.RequestSession, projectRepository project2.IProjectRepository,
	sbomListRepo sbomlist.ISbomListRepository,
	fixIt bool,
	state *integrity.DbIntegrityResult) {
	logy.Infof(requestSession, "Start analyse files on S3")
	filesOnS3 := s3Helper.ListObjects(requestSession, conf.Config.Server.GetUploadPath())
	for file := range filesOnS3 {
		if len(file.Key) < 1 {
			//ignore ghost files, sometime happens on S3 Mock
			logy.Errorf(requestSession, "Found file ghost! ")
			continue
		}
		state.CountUploadsOnS3++

		errorMsg := ""
		var exp *exception.Exception

		fileNameOnS3 := extractFilename(file.Key)
		deleteFile := false
		if len(fileNameOnS3) < 1 {
			errorMsg = "no file uuid/name found in path: " + file.Key
		} else {
			if strings.Contains(file.Key, versions_anchor) {
				errorMsg, exp, deleteFile = checkFoundFileOnS3OnVersionLevel(requestSession, projectRepository, file, fileNameOnS3, sbomListRepo, fixIt)
			} else if strings.Contains(file.Key, reports_anchor) {
				continue
			} else {
				errorMsg, exp, deleteFile = checkFoundDocumentOnS3OnProjectLevel(requestSession, projectRepository, file, fileNameOnS3, fixIt)
			}
			if deleteFile && fixIt {
				s3Helper.DeleteFile(requestSession, file.Key)
				errorMsg += " - file will be deleted"
			}
		}

		addMissingDBFileOnError(errorMsg, file, fixIt, exp, state, deleteFile)

		if (state.CountUploadsOnS3 % 100) == 0 {
			logy.Infof(requestSession, "analysed files on S3: "+strconv.Itoa(state.CountUploadsOnS3))
		}
	}

	logy.Infof(requestSession, "End analyse files on S3: "+strconv.Itoa(state.CountUploadsOnS3))

}

func addMissingDBFileOnError(errorMsg string, file minio.ObjectInfo, fixIt bool,
	exp *exception.Exception, state *integrity.DbIntegrityResult, fileDeleted bool) {
	if len(errorMsg) > 0 {
		missing := &integrity.MissingDBFile{
			FilePath:    file.Key,
			Upload:      file.LastModified,
			Message:     errorMsg,
			MetaData:    file.UserMetadata,
			Fixed:       fixIt,
			FileDeleted: fileDeleted,
		}
		if exp != nil {
			missing.ErrorMessage = exp.ErrorMessage
			missing.ErrorCode = exp.ErrorCode
			missing.ErrorRaw = exp.ErrorRaw
		}
		state.MissingFileOnDB = append(state.MissingFileOnDB, missing)
	}
}

func checkFoundFileOnS3OnVersionLevel(requestSession *logy.RequestSession,
	projectRepository project2.IProjectRepository,
	file minio.ObjectInfo,
	fileNameOnS3 string,
	sbomListRepo sbomlist.ISbomListRepository,
	fixIt bool) (errorMsg string, exp *exception.Exception, deleteFile bool) {

	//SBOM or DOCUMENT on version level
	projectUuid := extractProjectFromS3Path(file.Key, versions_anchor)
	if len(projectUuid) < 1 {
		return "no project uuid found in path: " + file.Key, nil, true
	}

	versionUuid := extractVersionFromS3Path(file.Key)
	if len(versionUuid) < 1 {
		return "no version uuid found in path: " + file.Key, nil, true
	}

	hash := ""
	deleteFile = false
	//SBOM
	sbomList := sbomListRepo.FindByKeyWithDeleted(requestSession, versionUuid, false)

	var spdxFile *project.SpdxFileBase
	if sbomList != nil {
		for _, file := range sbomList.SpdxFileHistory {
			if file.Key == fileNameOnS3 {
				spdxFile = file
				break
			}
		}
	}

	if spdxFile == nil {
		return "no sbom found for path: " + file.Key, nil, true
	}

	errorMsg, exp, hash = checkIfFileIntact(requestSession, file.Key, spdxFile.Hash, fixIt)
	//fix hash in db
	if fixIt {
		if len(hash) > 0 {
			//fix hash in db
			spdxFile.Hash = hash
			sbomListRepo.UpdateWithoutTimestamp(requestSession, sbomList)
			errorMsg += " - Adding/Replace HASH in DB Entry"
		} else if len(errorMsg) > 0 {
			deleteFile = true
			errorMsg += " - file will be deleted"
		}
	}
	return errorMsg, exp, deleteFile
}

func checkFoundDocumentOnS3OnProjectLevel(requestSession *logy.RequestSession,
	projectRepository project2.IProjectRepository,
	file minio.ObjectInfo, fileNameOnS3 string, fixIt bool) (errorMsg string, exp *exception.Exception, deleteFile bool) {
	//DOCUMENT on project level
	projectUuid := extractProjectFromS3Path(file.Key, documents_anchor)
	if len(projectUuid) < 1 {
		return "no project uuid found in path: " + file.Key, nil, true
	}

	projectEntity := projectRepository.FindByKeyWithDeleted(requestSession, projectUuid, false)
	if projectEntity == nil {
		return "no project found for uuid " + projectUuid + " in path: " + file.Key, nil, true
	}

	index := extractVersionIndex(fileNameOnS3)
	foundDocument := projectEntity.GetDocumentByFileNameWithIndex(fileNameOnS3, index)
	if foundDocument == nil {
		return "no document ref found in db for path: " + file.Key, nil, true
	}

	errorMsg, exp, hash := checkIfFileIntact(requestSession, file.Key, foundDocument.Hash, true)

	deleteFile = false
	if fixIt {
		if len(hash) > 0 {
			// update hash in db
			foundDocument.Hash = hash
			projectRepository.UpdateWithoutTimestamp(requestSession, projectEntity)
			errorMsg += " - Adding/Replace HASH in DB Entry"
		} else if len(errorMsg) > 0 {
			deleteFile = true
			errorMsg += " - file will be deleted"
		}

	}
	return errorMsg, exp, deleteFile
}

func checkIfFileIntact(requestSession *logy.RequestSession, filePath string, hash string, getRemoteHash bool) (errorMsg string, exp *exception.Exception, remoteFileHash string) {
	remoteFileHash = ""
	dbHasHash := len(hash) == 0
	if dbHasHash {
		errorMsg = "no HASH in DB"
	}
	loadHash := dbHasHash
	exception.TryCatch(func() {
		s3Helper.ReadFileFully(requestSession, filePath, hash)
	}, func(exception exception.Exception) {
		if getRemoteHash && exception.ErrorCode == message.ErrorFileHashCheck {
			loadHash = true
		}
		errorMsg += " could not load found document for path: " + filePath
		exp = &exception
	})
	if loadHash && (exp == nil || exp.ErrorCode == message.ErrorFileHashCheck) {
		remoteFileHash = s3Helper.GetFileHash(requestSession, filePath)
	}
	errorMsg = strings.Trim(errorMsg, "")
	return errorMsg, exp, remoteFileHash
}

func extractProjectFromS3Path(path string, anchor string) string {
	re := regexp.MustCompile(`.*[\\/](.*?)[\\/]` + anchor + `[\\/].*`)
	match := re.FindStringSubmatch(path)

	if len(match) < 2 {
		return ""
	}

	return match[1]
}

func extractVersionFromS3Path(path string) string {
	re := regexp.MustCompile(`.*[\\/]versions[\\/](.*?)[\\/].*`)
	match := re.FindStringSubmatch(path)

	if len(match) < 2 {
		return ""
	}

	return match[1]
}

func extractFilename(path string) string {
	re := regexp.MustCompile(`.*[\\/](.*)`)
	match := re.FindStringSubmatch(path)

	if len(match) < 2 {
		return ""
	}

	return match[1]
}

func extractVersionIndex(fileName string) int {
	re := regexp.MustCompile(`^([a-f0-9]{8}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{4}-[a-f0-9]{12})_([a-z]+)_?([0-4]?)-?(en|de)?\.(pdf|json)$`)
	match := re.FindStringSubmatch(fileName)
	if len(match) < 4 {
		return int(pdocument.NONE_VERSION)
	}

	index, err := strconv.Atoi(match[3])
	if err != nil {
		index = int(pdocument.NONE_VERSION)
	}
	return index
}
