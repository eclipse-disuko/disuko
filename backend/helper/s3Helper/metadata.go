// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package s3Helper

import (
	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/logy"
)

func MetadataForApplication(requestSession *logy.RequestSession, fileName, uploader string) map[string]string {
	metadata := make(map[string]string)
	metadata["ReqID"] = requestSession.ReqID
	metadata["ProjectName"] = "DISCO"
	metadata["ProjectUUID"] = "SYSTEM"
	metadata["FileName"] = fileName
	metadata["ServerIp"] = conf.Config.Server.LocalIp
	metadata["Uploader"] = uploader
	return metadata
}

func Metadata(requestSession *logy.RequestSession, currentProject *project.Project, versionKey, fileName, uploader string) map[string]string {
	metadata := make(map[string]string)
	metadata["ReqID"] = requestSession.ReqID
	metadata["ProjectName"] = currentProject.Name
	metadata["ProjectUUID"] = currentProject.UUID()
	if len(versionKey) > 0 {
		metadata["VersionName"] = currentProject.Versions[versionKey].Name
		metadata["VersionUUID"] = versionKey
	}
	metadata["FileName"] = fileName
	metadata["ServerIp"] = conf.Config.Server.LocalIp
	metadata["Uploader"] = uploader
	return metadata
}
