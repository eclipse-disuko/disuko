// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package service

import (
	"errors"
	"io"
	"strings"
	"time"

	"mercedes-benz.ghe.com/foss/disuko/domain/project/components"
	sbomlist2 "mercedes-benz.ghe.com/foss/disuko/domain/project/sbomlist"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/sbomlist"

	"mercedes-benz.ghe.com/foss/disuko/helper/stopwatch"
	schema2 "mercedes-benz.ghe.com/foss/disuko/infra/repository/schema"

	"mercedes-benz.ghe.com/foss/disuko/conf"
	"mercedes-benz.ghe.com/foss/disuko/domain"
	"mercedes-benz.ghe.com/foss/disuko/domain/project"
	"mercedes-benz.ghe.com/foss/disuko/domain/schema"
	"mercedes-benz.ghe.com/foss/disuko/helper"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/s3Helper"
	projectService "mercedes-benz.ghe.com/foss/disuko/infra/service/project"
	"mercedes-benz.ghe.com/foss/disuko/infra/service/spdx"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

func ValidateSbom(requestSession *logy.RequestSession, content string, spdxFile *project.SpdxFileBase, currentProject *project.Project, schemaRepository schema2.ISchemaRepository) error {
	activeSchemas := schemaRepository.FindActiveSchemas(requestSession)
	if len(activeSchemas) == 0 {
		return errors.New("no active schemas found")
	}
	var correspondingSchema *schema.SpdxSchema
	for _, s := range activeSchemas {
		if s.MatchesProjectLabel(currentProject.SchemaLabel) {
			correspondingSchema = s
			break
		}
	}
	if correspondingSchema == nil {
		return errors.New("there was no matching schema")
	}
	schemaValid, err := spdxFile.ValidateSpdxContent(requestSession, content, correspondingSchema)
	if !schemaValid {
		return err
	}
	spdxFile.SchemaValid = true
	spdxFile.SchemaId = correspondingSchema.Key
	spdxFile.SchemaName = correspondingSchema.Name + "-" + correspondingSchema.Version
	return nil
}

func UploadSbom(requestSession *logy.RequestSession, currentProject *project.Project,
	versionKey string, origin string, uploader string,
	file io.ReadSeekCloser, fileName string, tag string,
	holder projectService.RepositoryHolder,
	service *spdx.Service,
) (*project.SpdxFileBase, string) {
	sWOverAll := stopwatch.StopWatch{}
	sWOverAll.Start()

	hashValue := ""
	defer func() {
		sWOverAll.Stop()
		logy.Infof(requestSession, "SBOM Upload file over all time: %s with hash: %s", sWOverAll.DiffTime, hashValue)
	}()

	t := time.Now()
	spdxFile := &project.SpdxFileBase{
		ChildEntity:  domain.NewChildEntity(),
		Type:         schema.JSON,
		ContentValid: true,
		Uploaded:     &t,
		Origin:       origin,
		Uploader:     uploader,
		Tag:          tag,
	}

	wrap := helper.NewHashedReadAllWrap(file)
	content, err := wrap.ReadAllAndRewind()
	if err != nil {
		exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.SbomReadInReqFile, err.Error()), err)
	}
	spdxString := string(content)
	if conf.Config.Server.SbomValidationEnabled {
		err := ValidateSbom(requestSession, spdxString, spdxFile, currentProject, holder.SchemaRepository)
		if err != nil {
			return nil, err.Error()
		}
	}

	// S3 cleanup after error
	errorFlagForDeleteFileOnS3 := false
	filename := currentProject.GetFilePathSbom(spdxFile.Key, versionKey)
	defer func() {
		if errorFlagForDeleteFileOnS3 {
			s3Helper.DeleteFile(requestSession, filename)
		}
	}()

	metadata := s3Helper.Metadata(requestSession, currentProject, versionKey, fileName, uploader)
	// save SBOM file
	sWFile := stopwatch.StopWatch{}
	sWFile.Start()
	exception.TryCatchAndThrow(func() {
		s3Helper.SaveFile(requestSession, filename, file, metadata)
	}, func(exception exception.Exception) exception.Exception {
		errorFlagForDeleteFileOnS3 = true
		return exception
	})
	sWFile.Stop()
	logy.Infof(requestSession, "SBOM Upload file time: %s", sWFile.DiffTime)

	// check if SBOM exist on S3 and perform hash check
	hashValue, err = wrap.GetHash()
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorUnexpectError))
	sWFile.Start()
	s3Helper.PerformFileHashCheck(requestSession, "SBOM", filename, hashValue)
	sWFile.Stop()
	logy.Infof(requestSession, "SBOM Upload file hash check time: %s", sWFile.DiffTime)
	spdxFile.Hash = hashValue

	spdxFile.ExtractMetaInfo(spdxString)
	ci := project.FileContent(spdxString).ExtractComponentInfo(requestSession)
	if len(ci) > conf.Config.Server.SBomLimits.MaxComponents {
		errorFlagForDeleteFileOnS3 = true
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.SbomMaxComponents, conf.Config.Server.SBomLimits.MaxComponents))
	}
	ci.EnrichComponentInfos(requestSession)
	var (
		licenseErrComps   []string
		copyrightErrComps []string
		purlErrComps      []string
	)
	for _, c := range ci {
		if len(c.LicensesDeclared.List) > conf.Config.Server.SBomLimits.MaxLicensesPerComponent ||
			len(c.LicensesConcluded.List) > conf.Config.Server.SBomLimits.MaxLicensesPerComponent {
			licenseErrComps = append(licenseErrComps, c.Name+" ("+c.Version+")")
		}
		if len(c.CopyrightText) > conf.Config.Server.SBomLimits.MaxCopyRightTextSize {
			copyrightErrComps = append(copyrightErrComps, c.Name+" ("+c.Version+")")
		}
		if len(c.PURL) > conf.Config.Server.SBomLimits.MaxPURLSize {
			purlErrComps = append(purlErrComps, c.Name+" ("+c.Version+")")
		}
	}
	var componentErrors []string
	if len(licenseErrComps) > 0 {
		componentErrors = append(componentErrors,
			message.GetI18N(message.SbomMaxLicensesPerComponent, conf.Config.Server.SBomLimits.MaxLicensesPerComponent).Text+
				" Components: "+strings.Join(licenseErrComps, ",<br> "))
	}
	if len(copyrightErrComps) > 0 {
		componentErrors = append(componentErrors,
			message.GetI18N(message.SbomMaxCopyrightTextSize, conf.Config.Server.SBomLimits.MaxCopyRightTextSize).Text+
				" Components: "+strings.Join(copyrightErrComps, ",<br> "))
	}
	if len(purlErrComps) > 0 {
		componentErrors = append(componentErrors,
			message.GetI18N(message.SbomMaxPUrlSite, conf.Config.Server.SBomLimits.MaxPURLSize).Text+
				" Components: "+strings.Join(purlErrComps, ",<br> "))
	}
	if len(componentErrors) > 0 {
		errorFlagForDeleteFileOnS3 = true
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.SbomComponentsError, strings.Join(componentErrors, "; ")))
	}

	exception.TryCatchAndThrow(func() {
		spdxStats(requestSession, ci)
		service.WriteComponentInfos(requestSession, spdxFile.Key, &components.ComponentInfoList{
			UsedRefsListHash:     "",
			UsedLicenseRulesHash: "",
			ComponentInfos:       ci,
		})
		addSbomUpload(holder.SBOMListRepository, requestSession, versionKey, spdxFile)
		currentProject = holder.ProjectRepository.FindByKey(requestSession, currentProject.Key, false)
		version := currentProject.GetVersion(versionKey)
		version.UpdateStatusWhenUploadIsDone()
		currentProject.UpdateStatusToActive()
		version.Updated = time.Now()
		holder.ProjectRepository.Update(requestSession, currentProject)
	}, func(exception exception.Exception) exception.Exception {
		errorFlagForDeleteFileOnS3 = true
		return exception
	})
	return spdxFile, ""
}

func spdxStats(requestSession *logy.RequestSession, ci components.ComponentInfos) {
	processedComps := 0
	processedLicenses := 0

	for _, c := range ci {
		processedComps++
		processedLicenses += len(c.GetLicensesEffective().List)
	}

	logy.Infof(requestSession, "SpdxStats before saving, processed comps: %d processed licenses %d", processedComps, processedLicenses)
}

func addSbomUpload(sbomListRepo sbomlist.ISbomListRepository, requestSession *logy.RequestSession, versionKey string, spdxFile *project.SpdxFileBase) {
	sbomList := sbomListRepo.FindByKey(requestSession, versionKey, false)
	if sbomList == nil {
		sbomList := &sbomlist2.SbomList{}
		sbomList.Key = versionKey
		sbomList.SpdxFileHistory = append(sbomList.SpdxFileHistory, spdxFile)
		sbomListRepo.Save(requestSession, sbomList)
	} else {
		sbomList.SpdxFileHistory = append(sbomList.SpdxFileHistory, spdxFile)
		sbomListRepo.Update(requestSession, sbomList)
	}
}
