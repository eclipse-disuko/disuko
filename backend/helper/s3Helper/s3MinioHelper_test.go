// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package s3Helper

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"testing"

	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/stretchr/testify/assert"
)

const runS3Tests = false

var requestSessionTest = &logy.RequestSession{ReqID: "TEST"}

func PrepareS3Test(t *testing.T) *MinioS3Client {
	conf.Config.S3.AwsEndPoint = ""
	conf.Config.S3.AwsAccessKeyId = ""
	conf.Config.S3.BucketName = ""
	conf.Config.S3.IsEnabled = true

	client := CreateOrGetMinioClient(requestSessionTest)
	return client
}

func Test_S3WithMinio(t *testing.T) {
	logy.Infof(requestSessionTest, "START Test_S3WithMinio")
	if !runS3Tests {
		return
	}
	client := PrepareS3Test(t)

	logy.Infof(requestSessionTest, "Upload file")
	fileName := "test/test.json"
	testString := "TEST... ich bin ein Test"
	client.UploadObject(requestSessionTest, fileName, bytes.NewReader([]byte(testString)), make(map[string]string))

	logy.Infof(requestSessionTest, "Read meta data")
	meta := client.ReadFileMetaData(requestSessionTest, fileName)
	MustNotNil(t, meta)
	logy.Infof(requestSessionTest, "meta.Key: "+meta.Key)
	logy.Infof(requestSessionTest, "meta.Size: "+fmt.Sprint(meta.Size))
	logy.Infof(requestSessionTest, "meta.ContentType: "+meta.ContentType)
	logy.Infof(requestSessionTest, "meta.LastModified: "+fmt.Sprint(meta.LastModified))
	assert.Greater(t, meta.Size, int64(0))

	logy.Infof(requestSessionTest, "load file")
	fileReader := client.ReadFile(requestSessionTest, fileName)
	loadedString, err := io.ReadAll(fileReader)
	MustNil(t, err)
	logy.Infof(requestSessionTest, "loadedString: "+string(loadedString))
	MustNotNil(t, loadedString)
	assert.Equal(t, testString, string(loadedString))

	logy.Infof(requestSessionTest, "Delete file")
	client.DeleteFile(requestSessionTest, fileName)

	logy.Infof(requestSessionTest, "Check: Can not load")
	fileReader = client.ReadFile(requestSessionTest, fileName)
	_, err = io.ReadAll(fileReader)
	MustNotNil(t, err)

	logy.Infof(requestSessionTest, "END Test_S3WithMinio")
}

func Test_S3WithMinio_ListAll(t *testing.T) {
	logy.Infof(requestSessionTest, "START Test_S3WithMinio_ListAll")
	if !runS3Tests {
		return
	}
	client := PrepareS3Test(t)
	listObjects := client.ListObjects(requestSessionTest, "/")
	for object := range listObjects {
		if object.Err != nil {
			log.Fatal(object.Err)
			return
		}
		logy.Infof(requestSessionTest, "%v", object)
	}
	logy.Infof(requestSessionTest, "END Test_S3WithMinio_ListAll")
}
