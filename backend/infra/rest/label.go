// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"mercedes-benz.ghe.com/foss/disuko/domain"
	"mercedes-benz.ghe.com/foss/disuko/domain/job"
	"mercedes-benz.ghe.com/foss/disuko/domain/label"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
	"mercedes-benz.ghe.com/foss/disuko/helper/roles"
	"mercedes-benz.ghe.com/foss/disuko/helper/validation"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/labels"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/policyrules"
	project2 "mercedes-benz.ghe.com/foss/disuko/infra/repository/project"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/schema"
	"mercedes-benz.ghe.com/foss/disuko/infra/service/labelcsv"
	"mercedes-benz.ghe.com/foss/disuko/logy"
	"mercedes-benz.ghe.com/foss/disuko/observermngmt"
	"mercedes-benz.ghe.com/foss/disuko/scheduler"
)

type LabelHandler struct {
	LabelRepository   labels.ILabelRepository
	ProjectRepository project2.IProjectRepository
	SchemaRepository  schema.ISchemaRepository
	PolicyRepository  policyrules.IPolicyRulesRepository
	Scheduler         *scheduler.Scheduler
}

func (labelHandler *LabelHandler) GetLabels(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowLabel.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}
	labels := labelHandler.LabelRepository.FindAll(requestSession, false)
	queryType := r.URL.Query().Get("type")

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	resultLabels := make([]label.LabelResponseDto, 0)
	for _, currentLabel := range labels {
		var labelType string
		if currentLabel.Type == label.SCHEMA {
			labelType = "SCHEMA"
		} else if currentLabel.Type == label.PROJECT {
			labelType = "PROJECT"
		} else {
			labelType = "POLICY"
		}
		if len(queryType) == 0 || labelType == queryType {
			resultLabels = append(resultLabels, *currentLabel.ToDto())
		}
	}
	render.JSON(w, r, resultLabels)
}

func (labelHandler *LabelHandler) UpdateLabel(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowLabel.Update {
		exception.ThrowExceptionSendDeniedResponse()
	}
	labelRequest := extractLabelRequestBody(r)

	currentLabel := retrieveLabelForKey(requestSession, labelHandler.LabelRepository, w, r)

	updatedLabelName := strings.TrimSpace(strings.ToLower(labelRequest.Name))

	exisitingLabel := labelHandler.LabelRepository.FindByNameAndType(requestSession, updatedLabelName, currentLabel.Type)
	if exisitingLabel != nil && exisitingLabel.Key != currentLabel.Key {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorLabelAlreadyExist))
	}

	labelType := label.ConvertToLabelType(labelRequest.Type)

	if currentLabel.Type != labelType {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorLabelTypeCanNotChanged))
	}
	currentLabel.Update(updatedLabelName, labelRequest.Description)

	labelHandler.LabelRepository.Update(requestSession, currentLabel)

	err := labelHandler.Scheduler.ExecuteJobManual(requestSession, job.LabelLoadDb)
	if err != nil {
		exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorStartingJob), err)
	}

	w.WriteHeader(200)
}

func (labelHandler *LabelHandler) CreateLabel(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	username, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowLabel.Create {
		exception.ThrowExceptionSendDeniedResponse()
	}
	labelRequest := extractLabelRequestBody(r)
	labelType := label.ConvertToLabelType(labelRequest.Type)

	newLabelName := strings.TrimSpace(strings.ToLower(labelRequest.Name))

	exisitingLabel := labelHandler.LabelRepository.FindByNameAndType(requestSession, newLabelName, labelType)
	if exisitingLabel != nil {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.ErrorLabelAlreadyExist))
	}

	newLabel := &label.Label{
		RootEntity:  domain.NewRootEntity(),
		Name:        newLabelName,
		Description: labelRequest.Description,
		Type:        labelType,
	}

	labelHandler.LabelRepository.Save(requestSession, newLabel)

	err := labelHandler.Scheduler.ExecuteJobManual(requestSession, job.LabelLoadDb)
	if err != nil {
		exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorStartingJob), err)
	}

	observermngmt.FireEvent(observermngmt.DatabaseEntryAddedOrDeleted, observermngmt.DatabaseSizeChange{
		RequestSession: requestSession,
		CollectionName: labels.LabelCollectionName,
		Rights:         rights,
		Username:       username,
	})

	w.WriteHeader(200)
}

func (labelHandler *LabelHandler) DeleteLabel(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowLabel.Delete {
		exception.ThrowExceptionSendDeniedResponse()
	}
	currentLabel := retrieveLabelForKey(requestSession, labelHandler.LabelRepository, w, r)

	labelUsed := false
	if currentLabel.Type == label.SCHEMA {
		labelUsed = labelHandler.checkIfSchemaLabelIsUsed(requestSession, currentLabel)
	} else {
		labelUsed = labelHandler.checkIfPolicyLabelIsUsed(requestSession, currentLabel)
	}
	if labelUsed {
		exception.ThrowExceptionClientWithHttpCode(message.ErrorLabelUsed, message.GetI18N(message.ErrorLabelUsed).Text, "", exception.HTTP_CODE_SHOW_NO_REQUEST_ID)
	}

	labelHandler.LabelRepository.Delete(requestSession, currentLabel.Key)

	err := labelHandler.Scheduler.ExecuteJobManual(requestSession, job.LabelLoadDb)
	if err != nil {
		exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorStartingJob), err)
	}

	// Wait for job to finish
	time.Sleep(time.Second)

	w.WriteHeader(200)

	fmt.Fprintf(w, "Successfully deleted label with id %s \n", currentLabel.Key)
}

func (labelHandler *LabelHandler) CreateCSVHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !(rights.AllowLabel.Create && rights.AllowLabel.Update && rights.AllowLabel.Delete) {
		exception.ThrowExceptionSendDeniedResponse()
	}
	labelcsv.CreateCSV(&w, requestSession, labelHandler.PolicyRepository, labelHandler.LabelRepository)
}

func (labelHandler *LabelHandler) checkIfSchemaLabelIsUsed(requestSession *logy.RequestSession, label *label.Label) bool {
	exists := labelHandler.ProjectRepository.ExistsBySchemaLabel(requestSession, label.Key)
	if exists {
		return true
	}

	return labelHandler.SchemaRepository.ExistsByLabel(requestSession, label.Key)
}

func (labelHandler *LabelHandler) checkIfPolicyLabelIsUsed(requestSession *logy.RequestSession, label *label.Label) bool {
	exists := labelHandler.ProjectRepository.ExistsByPolicyLabel(requestSession, label.Key)
	if exists {
		return true
	}
	return labelHandler.PolicyRepository.ExistsByLabel(requestSession, label.Key)
}

func retrieveLabelForKey(requestSession *logy.RequestSession, repo labels.ILabelRepository, w http.ResponseWriter, r *http.Request) *label.Label {
	labelKey := chi.URLParam(r, "id")

	currentLabel := repo.FindByKey(requestSession, labelKey, false)
	if currentLabel == nil {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorDbNotFound), "Key is missing for label with id "+labelKey)
	}
	return currentLabel
}

func extractLabelRequestBody(r *http.Request) label.LabelRequestDto {
	var labelRequest label.LabelRequestDto
	validation.DecodeAndValidate(r, &labelRequest, false)
	return labelRequest
}

func CheckIfLabelsExistOrThrowException(requestSession *logy.RequestSession, labelRepository labels.ILabelRepository, keys []string) {
	for _, key := range keys {
		CheckIfLabelExistOrThrowException(requestSession, labelRepository, key)
	}
}

func CheckIfLabelExistOrThrowException(requestSession *logy.RequestSession, labelRepository labels.ILabelRepository, key string) {
	if !labelRepository.ExistByKey(requestSession, key) {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorLabelNotExist), "label.Key="+key)
	}
}
