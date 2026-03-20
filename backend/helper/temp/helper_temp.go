// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package temp

import (
	"crypto/rand"
	"encoding/hex"
	"os"
	"runtime"
	"strings"

	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/logy"
)

type TempHelper struct {
	Path           string
	RequestSession *logy.RequestSession
}

func generateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (tp *TempHelper) CreateRandomFolder() string {
	randomFolderName, err := generateRandomString(8)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorCreateFolder))
	targetPath := "/tmp/" + strings.ReplaceAll(randomFolderName, "/", "-")

	if runtime.GOOS == "windows" {
		targetPath = "C:/" + targetPath
	}
	tp.Path = targetPath

	err = os.Mkdir(targetPath, 0750)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.CouldNotWrite))
	return targetPath
}

func (tp *TempHelper) CreateFolder() string {
	// Prepare temporary folder includes remove all after work
	targetPath := "/tmp/" + strings.ReplaceAll(tp.RequestSession.ReqID, "/", "-")

	if runtime.GOOS == "windows" {
		targetPath = "C:/" + targetPath
	}
	tp.Path = targetPath

	err := os.Mkdir(targetPath, 0750)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.CouldNotWrite))
	return targetPath
}

func (tp *TempHelper) CreateSubFolder(name string) string {
	target := tp.Path + string(os.PathSeparator) + name
	err := os.Mkdir(target, 0755)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorCreateFolder))
	return target
}

func (tp *TempHelper) RemoveAll() {
	err := os.RemoveAll(tp.Path)
	if err != nil {
		if tp.RequestSession != nil {
			logy.Errorf(tp.RequestSession, "Cannot delete dir %s", tp.Path)
		}
	}
}

func (tp *TempHelper) WriteFile(fileName string, bytes []byte) string {
	completeFileName := tp.GetCompleteFileName(fileName)
	err := os.WriteFile(completeFileName, bytes, os.FileMode(int(0600)))
	if err != nil {
		if tp.RequestSession != nil {
			logy.Errorf(tp.RequestSession, "Cannot write file %s/%s", tp.Path, fileName)
		}
	}
	return completeFileName
}

func (tp *TempHelper) GetCompleteFileName(fileName string) string {
	return tp.Path + string(os.PathSeparator) + fileName
}
