// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package integrity

import (
	"time"

	"github.com/minio/minio-go/v7"
	"mercedes-benz.ghe.com/foss/disuko/conf"
	"mercedes-benz.ghe.com/foss/disuko/domain/integrity"
	"mercedes-benz.ghe.com/foss/disuko/domain/project"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/s3Helper"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/approvallist"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/dpconfig"
	project2 "mercedes-benz.ghe.com/foss/disuko/infra/repository/project"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/sbomlist"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type AnalyseCache struct {
	UploadsInS3 map[string]minio.ObjectInfo

	FilesInDb map[string]DbFileMeta
}

type DbFileMeta struct {
	Project     *project.Project
	Version     *project.ProjectVersion
	Type        string
	S3FileName  string
	OrgFileName string
	Upload      *time.Time
	Hash        string
	DbEntityKey string
}

func AnalyseDataIntegrity(requestSession *logy.RequestSession, projectRepository project2.IProjectRepository,
	dpConfigRepository *dpconfig.DBConfigRepository,
	sbomListRepo sbomlist.ISbomListRepository,
	approvalListRepo approvallist.IApprovalListRepository,
	fixIt bool) {
	exception.TryCatchAndLog(requestSession, func() {
		logy.Infof(requestSession, "AnalyseFiles - START")
		state := LoadDbIntegrityResult(requestSession, dpConfigRepository)
		exception.TryCatchAndThrow(func() {
			if fixIt && !conf.Config.Server.EnableFixDataIntegrity {
				addError2(requestSession, state, "Could not fix DataIntegrity, because is disabled!")
				fixIt = false
			}

			runnerCount := 3
			runnerChan := make(chan int, runnerCount)

			exception.RunAsyncAndLogExceptionAndInform(requestSession, func() {
				defer markedAsDone(runnerChan)
				//DB
				checkIfAllFilesFromDbExistOnS3(requestSession, projectRepository,
					sbomListRepo,
					fixIt,
					state)
			}, func(exception exception.Exception) exception.Exception {
				addError(requestSession, state, &exception)
				return exception
			})

			exception.RunAsyncAndLogExceptionAndInform(requestSession, func() {
				defer markedAsDone(runnerChan)
				//S3
				checkIfAllFilesFromS3ExistInTheDb(requestSession, projectRepository,
					sbomListRepo,
					fixIt,
					state)
			}, func(exception exception.Exception) exception.Exception {
				addError(requestSession, state, &exception)
				return exception
			})

			exception.RunAsyncAndLogExceptionAndInform(requestSession, func() {
				defer markedAsDone(runnerChan)
				//DB cross-references
				checkIfEachApprovalHasDocumentsOnProject(requestSession,
					projectRepository,
					approvalListRepo,
					fixIt,
					state)
			}, func(exception exception.Exception) exception.Exception {
				addError(requestSession, state, &exception)
				return exception
			})

			waitUntilAllDone(requestSession, runnerChan, runnerCount)

			//load metadata
			if !checkIsRunning(requestSession, dpConfigRepository) {
				return
			}
			loadMetadata(requestSession, state)

			state.ErrorsCount = len(state.Errors)
			state.MissingFileOnS3Count = len(state.MissingFileOnS3)
			state.MissingFileOnDBCount = len(state.MissingFileOnDB)
			state.MissingDocRefsOnProjectCount = len(state.MissingDocRefsOnProject)

			state.IsRunning = false
			state.EndTime = time.Now()
			SaveDbIntegrityResult(requestSession, state, dpConfigRepository)
			logy.Infof(requestSession, "AnalyseFiles - END")
		}, func(exception exception.Exception) exception.Exception {
			state.IsRunning = false
			state.EndTime = time.Now()
			addError(requestSession, state, &exception)
			SaveDbIntegrityResult(requestSession, state, dpConfigRepository)
			logy.Infof(requestSession, "AnalyseFiles - FAILED")
			return exception
		})
	})
}

func LoadDbIntegrityResult(requestSession *logy.RequestSession, repository *dpconfig.DBConfigRepository) *integrity.DbIntegrityResult {
	return repository.IntegrityCheck.Get(requestSession)
}

func SaveDbIntegrityResult(requestSession *logy.RequestSession, state *integrity.DbIntegrityResult, repository *dpconfig.DBConfigRepository) {
	repository.IntegrityCheck.Save(requestSession, state)
}

func checkIsRunning(requestSession *logy.RequestSession, dpConfigRepository *dpconfig.DBConfigRepository) bool {
	return LoadDbIntegrityResult(requestSession, dpConfigRepository).IsRunning
}

func loadMetadata(requestSession *logy.RequestSession, dbIntegrityResult *integrity.DbIntegrityResult) {
	runnerCount := 1
	doneLoadFileMetadataCount := make(chan int, runnerCount)

	exception.RunAsyncAndLogExceptionAndInform(requestSession, func() {
		loadMetadata2(requestSession, dbIntegrityResult, dbIntegrityResult.MissingFileOnDB, doneLoadFileMetadataCount)
	}, func(exception exception.Exception) exception.Exception {
		addError(requestSession, dbIntegrityResult, &exception)
		return exception
	})

	waitUntilAllDone(requestSession, doneLoadFileMetadataCount, runnerCount)

}

func loadMetadata2(requestSession *logy.RequestSession,
	dbIntegrityResult *integrity.DbIntegrityResult, files []*integrity.MissingDBFile,
	done chan int) {
	defer markedAsDone(done)
	if !conf.Config.S3.IsEnabled {
		return
	}
	for _, file := range files {
		if file.FileDeleted {
			continue
		}
		exception.TryCatch(func() {
			objectInfo := s3Helper.S3Client.ReadFileMetaData(requestSession, file.FilePath)
			file.MetaData = objectInfo.UserMetadata
		}, func(exception exception.Exception) {
			addError(requestSession, dbIntegrityResult, &exception)
		})
	}
}

func waitUntilAllDone(requestSession *logy.RequestSession, done chan int, max int) {
	logy.Infof(requestSession, "waitUntilAllDone - START max: %x", max)
	for i := 0; i < max; {
		i += <-done
		logy.Infof(requestSession, "waitUntilAllDone - i:%x max: %x", i, max)
	}
	close(done)
	logy.Infof(requestSession, "waitUntilAllDone - END")
}

func markedAsDone(done chan int) {
	done <- 1
}

func RemoveElementFromSlice[T any](slice []T, position int) []T {
	if position < 0 || position >= len(slice) {
		return slice
	}
	if position == len(slice)-1 {
		return slice[:position]
	}
	return append(slice[:position], slice[position+1:]...)
}

func addError(requestSession *logy.RequestSession, result *integrity.DbIntegrityResult, exc *exception.Exception) {
	errStr := ""
	if exc != nil {
		errStr = exc.ToString()
	}
	addError2(requestSession, result, errStr)
}

func addError2(requestSession *logy.RequestSession, result *integrity.DbIntegrityResult, message string) {
	logy.Errorf(requestSession, "Error: %s", message)
	result.Errors = append(result.Errors, message)
}
