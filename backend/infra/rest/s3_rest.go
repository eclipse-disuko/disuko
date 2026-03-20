// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"fmt"
	"html"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/minio/minio-go/v7"
	"mercedes-benz.ghe.com/foss/disuko/helper"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
	"mercedes-benz.ghe.com/foss/disuko/helper/roles"
	"mercedes-benz.ghe.com/foss/disuko/helper/s3Helper"
	"mercedes-benz.ghe.com/foss/disuko/helper/stopwatch"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type S3 struct {
	LastLoadTestResult string
}

func (s3 *S3) S3TriggerLoadTestHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.AllowS3Tests.Create && rights.AllowS3Tests.Read && rights.AllowS3Tests.Update && rights.AllowS3Tests.Delete) {
		exception.ThrowExceptionSendDeniedResponse()
	}

	exception.RunAsyncAndLogException(requestSession, func() {
		s3.runLoadTest(requestSession, w, r)
	})

	render.Status(r, http.StatusOK)
}

func (s3 *S3) S3TriggerLoadTestHandlerStatus(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.AllowS3Tests.Create && rights.AllowS3Tests.Read && rights.AllowS3Tests.Update && rights.AllowS3Tests.Delete) {
		exception.ThrowExceptionSendDeniedResponse()
	}

	render.JSON(w, r, s3.LastLoadTestResult)
}

func (s3 *S3) runLoadTest(requestSession *logy.RequestSession, w http.ResponseWriter, r *http.Request) {
	s3.LastLoadTestResult = ""
	sizeInMbStr := html.EscapeString(chi.URLParam(r, "sizeInMb"))
	sizeInMb, err := strconv.Atoi(sizeInMbStr)
	if err != nil {
		logy.Errorf(requestSession, "%v", err)
		exception.ThrowExceptionServerMessage(message.GetI18N(message.S3Error, err.Error()), "")
	}

	logs := &strings.Builder{}
	//upload test file
	fileName := "TEMP_" + strconv.FormatInt(time.Now().Unix(), 10) + ".dat"
	logs.WriteString("------ UPLOAD FILE: " + fileName + " with " + sizeInMbStr + " MB ------\n")
	uploadTime := stopwatch.StopWatch{}
	uploadTime.Start()
	testFile := s3Helper.UploadTempFileToS3(requestSession, sizeInMb, fileName)
	uploadTime.Stop()

	logs.WriteString("------ DOWNLOAD FILE: " + fileName + " ------\n")
	downloadTime := stopwatch.StopWatch{}
	downloadTime.Start()
	fileReader := s3Helper.S3Client.ReadFile(requestSession, fileName)
	fileContent := s3Helper.ConvertToStringAndClose(fileReader)
	downloadTime.Stop()

	logs.WriteString("------ DELETE FILE: " + fileName + " ------\n")
	s3Helper.S3Client.DeleteFile(requestSession, fileName)

	logs.WriteString("------ FINISH ------\n")
	logs.WriteString("Test File Size Upload expected (MB): " + sizeInMbStr + "\n")
	logs.WriteString("Test File Size Upload is (MB): " + fmt.Sprint(testFile.Step) + "\n")
	receivedFileSizeInMB := helper.ByteToMB(int64(len(fileContent)))
	logs.WriteString("Test File Size Download is (MB): " + fmt.Sprint(receivedFileSizeInMB) + "\n")
	logs.WriteString("Test File Name: " + fileName + "\n")
	logs.WriteString("Upload Time: " + fmt.Sprint(uploadTime.DiffTime) + "\n")
	logs.WriteString("Upload Rate: " + uploadTime.GetTransferRate(float64(testFile.Step)) + "\n")
	logs.WriteString("Download Time: " + fmt.Sprint(downloadTime.DiffTime) + "\n")
	logs.WriteString("Download Rate: " + downloadTime.GetTransferRate(receivedFileSizeInMB) + "\n")

	s3.LastLoadTestResult = logs.String()
}

func (s3 *S3) S3GetAllHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.AllowS3Tests.Create && rights.AllowS3Tests.Read && rights.AllowS3Tests.Update && rights.AllowS3Tests.Delete) {
		exception.ThrowExceptionSendDeniedResponse()
	}
	folder := chi.URLParam(r, "folder")
	folder = s3.unEscape(requestSession, w, folder)

	files := s3Helper.S3Client.ListObjects(requestSession, folder)

	result := make([]minio.ObjectInfo, 10)
	for file := range files {
		result = append(result, file)
	}
	render.JSON(w, r, result)
}

func (s3 *S3) S3GetMetadataFileHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.AllowS3Tests.Create && rights.AllowS3Tests.Read && rights.AllowS3Tests.Update && rights.AllowS3Tests.Delete) {
		exception.ThrowExceptionSendDeniedResponse()
	}
	fileName := chi.URLParam(r, "filename")
	fileName = s3.unEscape(requestSession, w, fileName)

	metaData := s3Helper.S3Client.ReadFileMetaData(requestSession, fileName)

	render.JSON(w, r, metaData)
}

func (s3 *S3) S3GetTextFileHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.AllowS3Tests.Create && rights.AllowS3Tests.Read && rights.AllowS3Tests.Update && rights.AllowS3Tests.Delete) {
		exception.ThrowExceptionSendDeniedResponse()
	}
	fileName := chi.URLParam(r, "filename")
	fileName = s3.unEscape(requestSession, w, fileName)

	downloadBytesReader := s3Helper.S3Client.ReadFile(requestSession, fileName)

	downloadBytes := s3Helper.ConvertToStringAndClose(downloadBytesReader)

	render.JSON(w, r, string(downloadBytes[:]))
}

func (s3 *S3) S3DeleteObjectHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.AllowS3Tests.Create && rights.AllowS3Tests.Read && rights.AllowS3Tests.Update && rights.AllowS3Tests.Delete) {
		exception.ThrowExceptionSendDeniedResponse()
	}
	fileName := chi.URLParam(r, "filename")
	fileName = s3.unEscape(requestSession, w, fileName)

	s3Helper.S3Client.DeleteFile(requestSession, fileName)
	render.JSON(w, r, "Object: "+fileName+" deleted")
}

func (s3 *S3) S3StoreFileHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.AllowS3Tests.Create && rights.AllowS3Tests.Read && rights.AllowS3Tests.Update && rights.AllowS3Tests.Delete) {
		exception.ThrowExceptionSendDeniedResponse()
	}
	fileName := chi.URLParam(r, "filename")
	fileName = s3.unEscape(requestSession, w, fileName)

	metadata := make(map[string]string)
	metadata["src"] = "server/infra/rest/s3_rest.go -> S3StoreFileHandler(..)"
	metadata["ReqID"] = requestSession.ReqID
	s3Helper.S3Client.UploadObject(requestSession, fileName, r.Body, metadata)
	render.JSON(w, r, fmt.Sprintf("File: %s uploaded to S3.", fileName))
}

// ### Helper functions ###
func (s3 *S3) unEscape(requestSession *logy.RequestSession, w http.ResponseWriter, fileName string) string {
	fileName, err := url.PathUnescape(fileName)
	if err != nil {
		logy.Errorf(requestSession, "%v", err)
		exception.ThrowExceptionServerMessage(message.GetI18N(message.HttpPathUnescape, err.Error()), "")
	}
	return fileName
}
