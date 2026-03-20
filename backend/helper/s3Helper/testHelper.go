// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package s3Helper

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

func UploadTempFileToS3(requestSession *logy.RequestSession, sizeInMb int, filePath string) *TempFile {
	Create1MbByteArray()

	tempFile := &TempFile{
		SizeInMb: sizeInMb,
	}
	metadata := make(map[string]string)
	metadata["desc"] = "Test file (server/helper/s3Helper/testHelper.go:13); sizeInMb=sizeInMb"
	S3Client.UploadObject(requestSession, filePath, tempFile, metadata)

	return tempFile
}

type TempFile struct {
	SizeInMb    int
	Step        int
	Template1Mb []byte
}

func (a *TempFile) Read(p []byte) (int, error) {
	if a.Step == a.SizeInMb {
		return 0, io.EOF
	}
	a.Step++

	if a.Template1Mb == nil {
		a.Template1Mb = Create1MbByteArray()
	}

	copy(p, a.Template1Mb)

	return len(a.Template1Mb), nil
}

func Create1MbByteArray() []byte {
	length := 1024
	content1KByte := make([]byte, length)
	seededRand := big.NewInt(int64(length))
	for i := 0; i < len(content1KByte); i++ {
		number, err := rand.Int(rand.Reader, seededRand)
		if err != nil {
			return nil
		}
		//content1KByte[i] = byte(rand.Int())
		content1KByte[i] = byte(number.Int64())
	}

	content1MByte := make([]byte, 0)
	for i := 0; i < 1024; i++ {
		content1MByte = append(content1MByte, content1KByte...)
	}

	return content1MByte
}

func MustNoError(t *testing.T, data interface{}, err error) {
	MustNotNil(t, data)
	MustNil(t, err)
}

func MustNotNil(t *testing.T, data interface{}) {
	assert.NotNil(t, data)
	if data == nil {
		panic("MustNoNil!")
	}
}
func MustTrue(t *testing.T, check bool) {
	assert.Equal(t, check, true)
	if !check {
		panic("MustTrue!")
	}
}

func MustFalse(t *testing.T, check bool) {
	assert.Equal(t, check, false)
	if check {
		panic("MustFalse!")
	}
}

func MustNil(t *testing.T, data interface{}) {
	assert.Nil(t, data)
	if data != nil {
		panic("MustNil, but was: " + fmt.Sprint(data))
	}
}
