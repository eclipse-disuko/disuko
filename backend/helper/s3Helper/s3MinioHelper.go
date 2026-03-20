// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package s3Helper

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioS3Client struct {
	minioClient *minio.Client
}

var minioS3Client *MinioS3Client = nil

func CreateOrGetMinioClient(requestSession *logy.RequestSession) *MinioS3Client {
	checkIfS3IsEnabledOrThrowException()

	if minioS3Client != nil && minioS3Client.minioClient.IsOnline() {
		return minioS3Client
	}
	endpoint := strings.ToLower(conf.Config.S3.AwsEndPoint)
	endpointSplitted := strings.Split(endpoint, "://")
	endpoint = endpointSplitted[1]
	useSSL := endpointSplitted[0] == "https"
	accessKeyID := conf.Config.S3.AwsAccessKeyId
	secretAccessKey := conf.Config.S3.AwsSecretAccessKey

	opt := &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	}

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, opt)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorS3ClientInit))

	logS3(requestSession, "Create Minio S3 Client")
	minioS3Client = &MinioS3Client{minioClient: minioClient}
	return minioS3Client
}

func checkIfS3IsEnabledOrThrowException() {
	if !conf.Config.S3.IsEnabled {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorS3Disabled), "")
	}
}

func (client *MinioS3Client) ListObjects(requestSession *logy.RequestSession,
	folder string) <-chan minio.ObjectInfo {
	checkIfS3IsEnabledOrThrowException()

	logS3(requestSession, "List S3 (Minio)")
	opts := minio.ListObjectsOptions{
		UseV1:     true,
		Recursive: true,
		Prefix:    folder,
	}

	// List all objects from a bucket-name with a matching prefix.
	return client.minioClient.ListObjects(context.Background(),
		conf.Config.S3.BucketName, opts)
}

func (client *MinioS3Client) DeleteFile(requestSession *logy.RequestSession, fileName string) {
	checkIfS3IsEnabledOrThrowException()

	logS3(requestSession, "Delete file on S3 (Minio), filename: "+fileName)
	err := client.minioClient.RemoveObject(context.Background(),
		conf.Config.S3.BucketName, fileName, minio.RemoveObjectOptions{})
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorS3Delete, fileName))
}

func (client *MinioS3Client) CopyFile(requestSession *logy.RequestSession, from, to string) {
	logS3(requestSession, fmt.Sprintf("Copy file on S3 (Minio), from %s to %s", from, to))

	src := minio.CopySrcOptions{
		Bucket: conf.Config.S3.BucketName,
		Object: from,
	}
	dst := minio.CopyDestOptions{
		Bucket: conf.Config.S3.BucketName,
		Object: to,
	}
	_, err := client.minioClient.CopyObject(context.Background(), dst, src)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorS3Copy, from, to))
}

func (client *MinioS3Client) ReadFileMetaData(requestSession *logy.RequestSession, fileName string) minio.ObjectInfo {
	checkIfS3IsEnabledOrThrowException()

	logS3(requestSession, "Read file meta on S3 (Minio), filename: "+fileName)
	objectReader, err := client.minioClient.StatObject(context.Background(),
		conf.Config.S3.BucketName, fileName, minio.StatObjectOptions{})
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorS3ReadMetaData, fileName))
	return objectReader
}
func (client *MinioS3Client) Exist(requestSession *logy.RequestSession, fileName string) bool {
	checkIfS3IsEnabledOrThrowException()

	logS3(requestSession, "Check if a exist file on S3 (Minio), filename: "+fileName)
	_, err := client.minioClient.StatObject(context.Background(),
		conf.Config.S3.BucketName, fileName, minio.StatObjectOptions{})
	if err != nil {
		resp := minio.ToErrorResponse(err)
		if resp.Code == "NoSuchKey" || resp.Code == "NotFound" {
			return false
		}
		exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorS3ReadMetaData, fileName))
	}
	return true
}

func (client *MinioS3Client) ReadFile(requestSession *logy.RequestSession, fileName string) io.ReadCloser {
	checkIfS3IsEnabledOrThrowException()

	logS3(requestSession, "Read file on S3 (Minio), filename: "+fileName)
	objectReader, err := client.minioClient.GetObject(context.Background(),
		conf.Config.S3.BucketName, fileName, minio.GetObjectOptions{})
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorS3Read, fileName))
	return objectReader
}

func (client *MinioS3Client) UploadObject(requestSession *logy.RequestSession, fileName string, fileContent io.Reader,
	metadata map[string]string) int64 {
	checkIfS3IsEnabledOrThrowException()

	logS3(requestSession, "Write file on S3 (Minio), filename: "+fileName)
	objectInfo, err := client.minioClient.PutObject(context.Background(),
		conf.Config.S3.BucketName, fileName, fileContent, -1, minio.PutObjectOptions{
			UserMetadata: metadata,
		})
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorS3Put, fileName))
	logy.Infof(requestSession, "[S3] PutObject responded with ETag: %s", objectInfo.ETag)
	return objectInfo.Size
}

const logPrefix = "[S3] "

func logS3(requestSession *logy.RequestSession, message string) {
	logy.Infof(requestSession, logPrefix+message)
}
