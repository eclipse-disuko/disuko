// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"net/http"

	"github.com/go-chi/render"
	"mercedes-benz.ghe.com/foss/disuko/domain/reviewremarks"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
	"mercedes-benz.ghe.com/foss/disuko/helper/roles"
	"mercedes-benz.ghe.com/foss/disuko/helper/validation"
)

func (projectHandler *ProjectHandler) BulkSetReviewRemarkStatus(w http.ResponseWriter, r *http.Request) {
	currentProject, version, requestSession := projectHandler.retrieveProjectAndVersion2(r)
	if currentProject.IsDeprecated() {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.DeprecatedProjectError))
	}

	username, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if currentProject.GetMember(username) == nil && !rights.IsProjectAnalyst() && !rights.IsDomainAdmin() {
		if !roles.CanAccessVehicleProjectOperations(rights, projectHandler.ProjectLabelService.HasVehiclePlatformLabel(requestSession, currentProject)) {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.ReadVersionReviewRemarks))
		}
	}

	var req reviewremarks.BulkReviewRemarkStatusRequest
	validation.DecodeAndValidate(r, &req, false)

	_, status := reviewremarks.ParseStatus(req.Status)

	remarks := projectHandler.ReviewRemarksRepository.FindByKey(requestSession, version.Key, false)

	remarkMap := make(map[string]*reviewremarks.Remark)
	for _, r := range remarks.Remarks {
		remarkMap[r.Key] = r
	}

	for _, remarkKey := range req.RemarkKeys {
		remark, exists := remarkMap[remarkKey]
		if !exists {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.GetRemarkUnknownKeyError))
		}
		if remark.Status != reviewremarks.Open && remark.Status != reviewremarks.InProgress {
			exception.ThrowExceptionClientMessage3(message.GetI18N(message.UpdateRemarkInvalidStatusError))
		}
	}

	// Update all remarks
	switch status {
	case reviewremarks.Closed:
		for _, remarkKey := range req.RemarkKeys {
			remark := remarkMap[remarkKey]
			remark.Close(username, string(status))
		}
	case reviewremarks.Cancelled:
		for _, remarkKey := range req.RemarkKeys {
			remark := remarkMap[remarkKey]
			remark.Cancel(username, string(status))
		}
	}

	// Store all changes
	projectHandler.ReviewRemarksRepository.Update(requestSession, remarks)

	responseData := SuccessResponse{
		Success: true,
		Message: "bulk status update successful",
	}
	render.JSON(w, r, responseData)
}
