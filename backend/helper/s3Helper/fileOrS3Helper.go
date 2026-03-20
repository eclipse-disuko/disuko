// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package s3Helper

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"mercedes-benz.ghe.com/foss/disuko/helper"

	"mercedes-benz.ghe.com/foss/disuko/conf"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/hash"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
	"mercedes-benz.ghe.com/foss/disuko/helper/stopwatch"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

var S3Client *MinioS3Client

func SetUpS3Client(requestSession *logy.RequestSession) {
	if !conf.Config.S3.IsEnabled {
		return
	}
	S3Client = CreateOrGetMinioClient(requestSession)
}

func CreateFolderIfNotExist(requestSession *logy.RequestSession, folder string, panicOnError bool) {
	if conf.Config.S3.IsEnabled {
		// nothing to do, S3 has no folder itself, folders part of the filename
		return
	}

	CreateFolderIfNotExistOnLocalFileSystem(requestSession, folder, panicOnError)
}

func CreateFolderIfNotExistOnLocalFileSystem(requestSession *logy.RequestSession, folder string, panicOnError bool) {
	// use normal local file system
	filePath := folder
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		err = os.MkdirAll(filePath, 0750)
		if err != nil {
			if panicOnError {
				exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorCreateFolder, folder))
				logy.Fatalf(requestSession, "failed to create folder %s, %s", folder, err)
			}
		}
	}
}

func ListObjects(requestSession *logy.RequestSession, folder string) <-chan minio.ObjectInfo {
	if conf.Config.S3.IsEnabled {
		return S3Client.ListObjects(requestSession, folder)
	} else {
		return ListObjectsOnLocalFileSystem(requestSession, folder)
	}
}

type CountMeta struct {
	CntPDF   int
	CntJson  int
	CntFiles int
	CntSBOM  int
}

func CountFiles(requestSession *logy.RequestSession, folder string) CountMeta {
	result := CountMeta{
		CntPDF:   0,
		CntJson:  0,
		CntFiles: 0,
		CntSBOM:  0,
	}
	folder = strings.Trim(folder, "/")
	files := ListObjects(requestSession, folder)
	for file := range files {
		if len(file.Key) < 1 {
			// ignore ghost files, sometime happens on S3 Mock
			logy.Errorf(requestSession, "Found file ghost! ")
			continue
		}
		if strings.Contains(file.Key, ".pdf") {
			if strings.HasSuffix(file.Key, "-de.pdf") {
				// count only the base disclosure files
				result.CntPDF++
			}
		} else if strings.Contains(file.Key, ".json") {
			result.CntJson++
		} else {
			if strings.Contains(file.Key, "/sbom/") {
				result.CntSBOM++
			}
		}
		result.CntFiles++
	}
	return result
}

func ListObjectsOnLocalFileSystem(requestSession *logy.RequestSession, dir string) <-chan minio.ObjectInfo {
	out := make(chan minio.ObjectInfo)

	exception.RunAsyncAndLogException(requestSession, func() {
		defer close(out)

		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorUnexpectError))

			if !info.IsDir() {
				objInfo := minio.ObjectInfo{Key: path}
				out <- objInfo
			}

			return nil
		})

		exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorUnexpectError))
	})

	return out
}

func SaveObjectToFile(requestSession *logy.RequestSession, fileName string, content interface{}, metadata map[string]string) int64 {
	jsonBytes, err := json.MarshalIndent(content, "", "  ")
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorJsonMarshalling))
	return SaveFile(requestSession, fileName, bytes.NewReader(jsonBytes), metadata)
}

func SaveObjectToFileAndGetHash(requestSession *logy.RequestSession, fileName string, content interface{}, metadata map[string]string) string {
	jsonBytes, err := json.MarshalIndent(content, "", "  ")
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorJsonMarshalling))
	wrap := helper.NewSha256ReaderWrapper(bytes.NewReader(jsonBytes))
	SaveFile(requestSession, fileName, wrap, metadata)
	return wrap.GetHash()
}

func SaveObjectToLocalFileAndGetHash(requestSession *logy.RequestSession, fileName string, content interface{}, metadata map[string]string) string {
	jsonBytes, err := json.MarshalIndent(content, "", "  ")
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorJsonMarshalling))
	wrap := helper.NewSha256ReaderWrapper(bytes.NewReader(jsonBytes))
	SaveFileToLocalFileSystem(requestSession, fileName, bytes.NewReader(jsonBytes))
	return wrap.GetHash()
}

func SaveFile(requestSession *logy.RequestSession, fileName string, content io.Reader, metadata map[string]string) int64 {
	return SaveFileAndGetLogs(requestSession, fileName, content, metadata)
}

func SaveFileAndGetHash(requestSession *logy.RequestSession, fileName string, content io.Reader, metadata map[string]string) string {
	wrap := helper.NewSha256ReaderWrapper(content)
	SaveFileAndGetLogs(requestSession, fileName, wrap, metadata)
	return wrap.GetHash()
}

func SaveFileAndGetLogs(requestSession *logy.RequestSession, fileName string,
	content io.Reader, metadata map[string]string,
) int64 {
	sWFile := stopwatch.StopWatch{}
	sWFile.Start()
	existFile := ExistFile(requestSession, fileName)
	sWFile.Stop()
	logy.Infof(requestSession, "ExistFile time: %s", sWFile.DiffTime)
	if existFile {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorFileExistAlready, fileName), "")
	}

	if conf.Config.S3.IsEnabled {
		return S3Client.UploadObject(requestSession, fileName, content, metadata)
	}
	return SaveFileToLocalFileSystem(requestSession, fileName, content)
}

func Copy(requestSession *logy.RequestSession, from, to string) {
	if conf.Config.S3.IsEnabled {
		S3Client.CopyFile(requestSession, from, to)
	} else {
		file := ReadFile(requestSession, from)
		defer file.Close()
		SaveFileToLocalFileSystem(requestSession, to, file)
	}
}

func CopyLocalFilesystem(from, to string) error {
	source, err := os.Open(from)
	if err != nil {
		return fmt.Errorf("opening source file: %w", err)
	}
	defer source.Close()

	destination, err := os.Create(to)
	if err != nil {
		return fmt.Errorf("creating destination file: %w", err)
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	if err != nil {
		return fmt.Errorf("copy: %w", err)
	}
	return nil
}

func SaveObjectToLocalFileSystem(requestSession *logy.RequestSession, fileName string, content interface{}) {
	jsonBytes, err := json.MarshalIndent(content, "", "  ")
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorJsonMarshalling))
	SaveFileToLocalFileSystem(requestSession, fileName, bytes.NewReader(jsonBytes))
}

func SaveFileToLocalFileSystem(requestSession *logy.RequestSession, fileName string, content io.Reader) int64 {
	dir := filepath.Dir(fileName)
	CreateFolderIfNotExist(requestSession, dir, false)

	f, err := os.Create(fileName)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorCreateFile, fileName))
	defer f.Close()

	writtenBytes, err := io.Copy(f, content)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorFileWrite, fileName))

	return writtenBytes
}

func ExistFile(requestSession *logy.RequestSession, filePath string) bool {
	if conf.Config.S3.IsEnabled {
		return S3Client.Exist(requestSession, filePath)
	}
	return ExistsOnLocalFileSystem(filePath)
}

func ExistsOnLocalFileSystem(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func ReadFileFromLocalFileSystem(filePath string) io.ReadSeekCloser {
	file, err := os.Open(filePath)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorOsOpenFile, filePath))
	return file
}

func DeleteFile(requestSession *logy.RequestSession, filePath string) {
	if conf.Config.S3.IsEnabled {
		S3Client.DeleteFile(requestSession, filePath)
	}
	DeleteFileFromLocalFileSystem(requestSession, filePath)
}

func DeleteFileFromLocalFileSystem(requestSession *logy.RequestSession, filePath string) {
	if conf.Config.S3.IsEnabled {
		S3Client.DeleteFile(requestSession, filePath)
	} else {
		_, err := os.Stat(filePath)
		if err == nil {
			err = os.Remove(filePath)
			exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorFileDelete, filePath))
		}
	}
}

func ReadFile(requestSession *logy.RequestSession, filePath string) io.ReadCloser {
	if conf.Config.S3.IsEnabled {
		return S3Client.ReadFile(requestSession, filePath)
	}
	return ReadFileFromLocalFileSystem(filePath)
}

func ReadFileFully(requestSession *logy.RequestSession, filePath string, sha256Hash string) []byte {
	file := ReadFile(requestSession, filePath)
	result := ReadAllAndClose(file)

	if len(sha256Hash) > 0 { // for old data without hash value
		s3Hash := hash.GetSha256Hash(result)
		if s3Hash != sha256Hash {
			exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorFileHashCheck, filePath, s3Hash, sha256Hash), "")
		}
	}
	return result
}

func GetFileHash(requestSession *logy.RequestSession, filePath string) string {
	file := ReadFile(requestSession, filePath)
	result := ReadAllAndClose(file)
	return hash.GetSha256Hash(result)
}

func ReadTextFile(requestSession *logy.RequestSession, filePath string, sha256Hash string) *string {
	readedBytes := ReadFileFully(requestSession, filePath, sha256Hash)
	resultString := string(readedBytes)
	return &resultString
}

func ConvertToStringAndClose(source io.ReadCloser) string {
	content := ReadAllAndClose(source)
	return string(content)
}

func ReadAllAndClose(source io.ReadCloser) []byte {
	defer source.Close()
	content, err := io.ReadAll(source)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorReadAllAndClose))
	return content
}

func PerformDownload(requestSession *logy.RequestSession, w *http.ResponseWriter, filePath string, hash string) {
	// FILE SEND
	(*w).WriteHeader(http.StatusOK)
	fileContent := ReadFileFully(requestSession, filePath, hash)
	_, err := (*w).Write(fileContent)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.WritingContent))
}

func PerformFileHashCheck(requestSession *logy.RequestSession, fileType string, filePath string, sha256Hash string) {
	// HASH CHECK
	if len(sha256Hash) == 0 {
		// for old data without hash value
		return
	}
	readerForHash := ReadFile(requestSession, filePath)
	defer readerForHash.Close()

	sha256hash := sha256.New()
	_, err := io.Copy(sha256hash, readerForHash)
	exception.HandleErrorClientMessage(err, message.GetI18N(message.WritingContent))

	if hash.EncodeSha256ToString(sha256hash) != sha256Hash {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ReadFileCorrupted, fileType), "")
	}
}
