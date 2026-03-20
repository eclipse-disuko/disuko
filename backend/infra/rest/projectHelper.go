// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"strconv"
	"time"

	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"

	"github.com/eclipse-disuko/disuko/conf"
	project2 "github.com/eclipse-disuko/disuko/infra/repository/project"
	"github.com/eclipse-disuko/disuko/logy"
)

func TryNewFileUpload(requestSession *logy.RequestSession, projectKey string, projectRepo project2.IProjectRepository) {
	// clean uploads older the one-hour
	project := projectRepo.FindByKey(requestSession, projectKey, false)
	newListOfUploadTimes := make([]time.Time, 0)
	beforeOneHour := time.Now().Add(-time.Hour)
	for _, uploadTime := range project.LastFileUploads {
		if uploadTime.After(beforeOneHour) {
			newListOfUploadTimes = append(newListOfUploadTimes, uploadTime)
		}
	}

	// count uploads of the last hour
	if len(newListOfUploadTimes) >= conf.Config.Server.MaxUploadPerHourPerProject {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorFileMaxUploadPerHourReached, project.Name, project.Key, strconv.Itoa(len(newListOfUploadTimes)), strconv.Itoa(conf.Config.Server.MaxUploadPerHourPerProject)))
	}

	// add new upload time
	project.LastFileUploads = append(newListOfUploadTimes, time.Now())
	projectRepo.Update(requestSession, project)
}
