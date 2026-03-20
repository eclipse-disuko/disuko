// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package test

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	dpconfig2 "github.com/eclipse-disuko/disuko/infra/repository/dpconfig"
	"github.com/eclipse-disuko/disuko/infra/repository/labels"
	"github.com/eclipse-disuko/disuko/infra/repository/policyrules"
	"github.com/eclipse-disuko/disuko/infra/repository/sbomlist"

	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/dpconfig"
	"github.com/eclipse-disuko/disuko/domain/label"
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/roles"
	"github.com/eclipse-disuko/disuko/helper/s3Helper"
	"github.com/eclipse-disuko/disuko/helper/stopwatch"
	"github.com/eclipse-disuko/disuko/infra/repository/license"
	"github.com/eclipse-disuko/disuko/infra/repository/obligation"
	project2 "github.com/eclipse-disuko/disuko/infra/repository/project"
	"github.com/eclipse-disuko/disuko/infra/repository/schema"
	"github.com/eclipse-disuko/disuko/infra/service"
	projectService "github.com/eclipse-disuko/disuko/infra/service/project"
	"github.com/eclipse-disuko/disuko/infra/service/spdx"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/go-chi/render"
)

type SampleDataHandler struct {
	DpConfigRepo          *dpconfig2.DBConfigRepository
	ProjectRepository     project2.IProjectRepository
	LicensesRepository    license.ILicensesRepository
	PolicyRulesRepository policyrules.IPolicyRulesRepository
	ObligationRepository  obligation.IObligationRepository
	SchemaRepository      schema.ISchemaRepository
	LabelRepository       labels.ILabelRepository
	SbomListRepository    sbomlist.ISbomListRepository
	SpdxService           *spdx.Service
}

func (handler *SampleDataHandler) GetStateCreateSampleDataHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)

	state, _ := handler.getStateAndCheckRights(w, r, requestSession)
	if state == nil {
		state = &dpconfig.SampleDataCreationState{}
	}
	render.JSON(w, r, state)
}

func (handler *SampleDataHandler) StopStateCreateSampleDataHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	state, _ := handler.getStateAndCheckRights(w, r, requestSession)
	if state != nil {
		state.EndTime = time.Now()
		state.IsRunning = false
		handler.saveState(requestSession, state)
	}
	w.WriteHeader(200)
	w.Write([]byte("STOPPED"))
}

func (handler *SampleDataHandler) StartCreateSampleDataHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	state, userName := handler.getStateAndCheckRights(w, r, requestSession)
	if state == nil || !state.IsRunning {
		state = &dpconfig.SampleDataCreationState{
			RootEntity: domain.NewRootEntity(),
			BaseState: domain.BaseState{
				StartTime: time.Now(),
				IsRunning: true,
				ReqID:     requestSession.ReqID,
			},
		}
	} else {
		w.WriteHeader(200)
		w.Write([]byte("Already running"))
		return
	}

	cntSample, err := strconv.Atoi(r.URL.Query().Get("cnt"))
	if err != nil {
		cntSample = 1
	}
	state.TargetCount = cntSample

	withFileUpload, err := strconv.ParseBool(r.URL.Query().Get("fileUpload"))
	if err != nil {
		withFileUpload = false
	}
	state.WithFileUpload = withFileUpload
	handler.saveState(requestSession, state)

	exception.RunAsyncAndLogException(requestSession, func() {
		handler.startCreateSampleData(requestSession, userName, state, withFileUpload)
	})

	w.WriteHeader(200)
	w.Write([]byte("STARTED"))
}

func (handler *SampleDataHandler) startCreateSampleData(requestSession *logy.RequestSession,
	userName string, state *dpconfig.SampleDataCreationState, withSbomFileUpload bool,
) {
	defer exception.CatchExceptionWithCustom(exception.NewExceptionHandler2(requestSession),
		func(exception2 exception.Exception) {
			state.HasErrors = true
			state.IsRunning = false
			state.LastError = exception2.ErrorMessage + ", raw= " + exception2.ErrorRaw
			exception.TryCatchAndLog(requestSession, func() {
				handler.DpConfigRepo.SampleDataCreationState.Save(requestSession, state)
			})
		})

	logy.Infof(requestSession, "StartCreateSampleDataHandler - START")

	sWOverAll := stopwatch.StopWatch{}
	sWOverAll.Start()
	sWPart := stopwatch.StopWatch{}
	sWPart.Start()
	allProjects := handler.ProjectRepository.FindAll(requestSession, true)
	sWPart.Stop()
	logy.Infof(requestSession, "FindAll all projects: %s", sWPart.DiffTime)

	sWPart.Start()
	for _, prj := range allProjects {
		if len(prj.FreeLabels) > 0 {
			inStr := strings.Join(prj.FreeLabels, ",")
			if strings.Contains(inStr, "auto generated") {
				handler.ProjectRepository.Delete(requestSession, prj.Key)
			}
		}
	}
	sWPart.Stop()
	logy.Infof(requestSession, "Delete all (auto generated) projects: %s", sWPart.DiffTime)

	for i := 0; i < state.TargetCount; i++ {
		sWPart.Start()
		state = handler.getState(requestSession)
		state.CreatedCount++
		handler.saveState(requestSession, state)
		sWPart.Stop()
		logy.Infof(requestSession, "Get/Save state: %s", sWPart.DiffTime)
		if !state.IsRunning {
			break
		}

		logy.Infof(requestSession, "Create (auto generated) project index: "+strconv.Itoa(i))

		sWPart.Start()
		newProject := handler.createProject(i, requestSession, userName)
		handler.setSchema(requestSession, newProject)
		handler.setPolicyLabels(requestSession, newProject)
		versionKey := handler.createVersions(requestSession, newProject, state)
		handler.updateProject(requestSession, newProject)
		sWPart.Stop()
		logy.Infof(requestSession, "Create/save project: %s", sWPart.DiffTime)

		if withSbomFileUpload {
			sWPart.Start()
			handler.uploadSbom(requestSession, userName, conf.Config.Server.BasePath+"/resources/backendSampleSPDX.json", newProject, versionKey)
			sWPart.Stop()
			logy.Infof(requestSession, "Upload sbom: %s", sWPart.DiffTime)
		}
	}
	sWOverAll.Stop()
	logy.Infof(requestSession, "Create sample projects overall time: %s", sWOverAll.DiffTime)
	logy.Infof(requestSession, "StartCreateSampleDataHandler - END")

	if state.IsRunning {
		state.IsRunning = false
		state.EndTime = time.Now()
		handler.saveState(requestSession, state)
	}
}

func (handler *SampleDataHandler) createProject(i int, requestSession *logy.RequestSession, owner string) *project.Project {
	newProject := project.CreateNewProject(project.ProjectRequestDto{
		Name:         "New Project " + strconv.Itoa(i),
		Owner:        owner,
		CreationMode: project.DEFAULT,
	})
	newProject.PolicyLabels = make([]string, 0)
	newProject.FreeLabels = make([]string, 0)
	newProject.FreeLabels = append(newProject.FreeLabels, "auto generated")

	handler.ProjectRepository.Save(requestSession, newProject)
	return newProject
}

func (handler *SampleDataHandler) setSchema(requestSession *logy.RequestSession, newProject *project.Project) {
	schemaLabelsActive := handler.LabelRepository.FindByNameAndType(requestSession, "common standard", label.SCHEMA)
	newProject.SchemaLabel = schemaLabelsActive.Key
}

func (handler *SampleDataHandler) setPolicyLabels(requestSession *logy.RequestSession, newProject *project.Project) {
	// enterprise platform, mobile platform, other platform
	policyLabel := handler.LabelRepository.FindByNameAndType(requestSession, "enterprise platform", label.POLICY)
	newProject.PolicyLabels = append(newProject.PolicyLabels, policyLabel.Key)

	// architecture: frontend layer, backend layer;
	policyLabel = handler.LabelRepository.FindByNameAndType(requestSession, "backend layer", label.POLICY)
	newProject.PolicyLabels = append(newProject.PolicyLabels, policyLabel.Key)

	// target user: entity users, group users, external users
	policyLabel = handler.LabelRepository.FindByNameAndType(requestSession, "entity users", label.POLICY)
	newProject.PolicyLabels = append(newProject.PolicyLabels, policyLabel.Key)

	// distribution target: entity target, external target
	policyLabel = handler.LabelRepository.FindByNameAndType(requestSession, "entity target", label.POLICY)
	newProject.PolicyLabels = append(newProject.PolicyLabels, policyLabel.Key)
}

func (handler *SampleDataHandler) uploadSbom(requestSession *logy.RequestSession,
	userName string, fileName string, newProject *project.Project, versionKey string,
) {
	sWFile := stopwatch.StopWatch{}
	sWFile.Start()
	fileReader := s3Helper.ReadFileFromLocalFileSystem(fileName)
	sWFile.Stop()
	logy.Infof(requestSession, "read file: %s", sWFile.DiffTime)

	sWFile.Start()

	holder := projectService.RepositoryHolder{
		ProjectRepository:  handler.ProjectRepository,
		LicenseRepository:  handler.LicensesRepository,
		SBOMListRepository: handler.SbomListRepository,
		SchemaRepository:   handler.SchemaRepository,
	}
	service.UploadSbom(requestSession, newProject,
		versionKey, project.OriginServer,
		userName,
		fileReader,
		fileName,
		"",
		holder,
		handler.SpdxService,
	)
	sWFile.Stop()
	logy.Infof(requestSession, "Upload file: %s", sWFile.DiffTime)
}

func (handler *SampleDataHandler) createVersions(requestSession *logy.RequestSession, project *project.Project, state *dpconfig.SampleDataCreationState) string {
	project.CreateNewProjectVersionIfNameNotUsed("1.0", "my version 1.0")
	project.CreateNewProjectVersionIfNameNotUsed("1.8", "my version 1.8")
	project.CreateNewProjectVersionIfNameNotUsed("1.9", "my version 1.9")
	versionName := "2.0"
	versionKey := project.CreateNewProjectVersionIfNameNotUsed(versionName, "my version 2.0")
	return versionKey
}

func (handler *SampleDataHandler) updateProject(requestSession *logy.RequestSession, project *project.Project) {
	handler.ProjectRepository.Update(requestSession, project)
}

func (handler *SampleDataHandler) saveState(requestSession *logy.RequestSession, state *dpconfig.SampleDataCreationState) {
	handler.DpConfigRepo.SampleDataCreationState.Save(requestSession, state)
}

func (handler *SampleDataHandler) getStateAndCheckRights(w http.ResponseWriter, r *http.Request, requestSession *logy.RequestSession) (*dpconfig.SampleDataCreationState, string) {
	userName, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.AllowSampleData.Create && rights.AllowSampleData.Read && rights.AllowSampleData.Update && rights.AllowSampleData.Delete) {
		exception.ThrowExceptionSendDeniedResponse()
	}

	state := handler.getState(requestSession)
	return state, userName
}

func (handler *SampleDataHandler) getState(requestSession *logy.RequestSession) *dpconfig.SampleDataCreationState {
	return handler.DpConfigRepo.SampleDataCreationState.Get(requestSession)
}
