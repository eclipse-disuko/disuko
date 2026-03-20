// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"mercedes-benz.ghe.com/foss/disuko/conf"
	"mercedes-benz.ghe.com/foss/disuko/helper/s3Helper"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

func (s *Server) setupS3(requestSession *logy.RequestSession) {
	s3Helper.SetUpS3Client(requestSession)
	setUpFiles(requestSession)
}

func setUpFiles(requestSession *logy.RequestSession) {
	// ignore errors. will panic anyway.
	s3Helper.CreateFolderIfNotExist(requestSession, conf.Config.Server.Uploadpath, true)
}
