// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	i18nDomain "github.com/eclipse-disuko/disuko/domain/i18n"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/helper/roles"
	"github.com/eclipse-disuko/disuko/helper/validation"
	i18nRepo "github.com/eclipse-disuko/disuko/infra/repository/i18n"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type I18nHandler struct {
	I18nRepository i18nRepo.II18nRepository
}

// protectedLocales lists locale codes that cannot be deleted via the API.
var protectedLocales = map[string]struct{}{
	"en": {},
	"de": {},
}

func parseLocaleImportJSON(fileName string, payload []byte) (map[string]string, *i18nDomain.I18nImportIssueDto) {
	var result map[string]string
	if err := json.Unmarshal(payload, &result); err != nil {
		return nil, &i18nDomain.I18nImportIssueDto{
			FileName: fileName,
			Code:     "INVALID_JSON",
			Message:  "File contains invalid JSON or non-string values",
		}
	}
	return result, nil
}

func ensureI18nWriteAccess(requestSession *logy.RequestSession, r *http.Request) {
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.IsApplicationAdmin() && !rights.IsDomainAdmin() {
		exception.ThrowExceptionSendDeniedResponse()
	}
}

func isValidLocaleCode(code string) bool {
	if len(code) < 2 || len(code) > 35 {
		return false
	}
	for _, c := range code {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '-' || c == '_') {
			return false
		}
	}
	return true
}

func (handler *I18nHandler) GetLocale(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	requestedLocale := strings.TrimSpace(chi.URLParam(r, "locale"))
	if !isValidLocaleCode(requestedLocale) {
		exception.ThrowExceptionBadRequestResponse()
	}

	locale, fallbackUsed := handler.I18nRepository.FindByLocaleCodeOrDefault(requestSession, requestedLocale)
	if locale == nil {
		exception.ThrowExceptionClient404Message(message.GetI18N(message.ErrorDbNotFound), "i18n locale not found: "+requestedLocale)
	}

	render.JSON(w, r, locale.ToDTO(fallbackUsed))
}

func (handler *I18nHandler) ExportLocaleJSON(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	localeCode := strings.TrimSpace(chi.URLParam(r, "locale"))
	if !isValidLocaleCode(localeCode) {
		exception.ThrowExceptionBadRequestResponse()
	}

	locale := handler.I18nRepository.FindByLocaleCode(requestSession, localeCode, false)
	if locale == nil {
		exception.ThrowExceptionClient404Message(message.GetI18N(message.ErrorDbNotFound), "i18n locale not found: "+localeCode)
	}

	entries := make(map[string]string)
	for key, entry := range locale.Entries {
		entries[key] = entry.Value
	}

	body, err := json.Marshal(entries)
	if err != nil {
		exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorUnexpectError), err)
	}

	filename := fmt.Sprintf("locale.%s.json", strings.ToLower(localeCode))
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(body)
}

func (handler *I18nHandler) ImportLocaleJSON(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	ensureI18nWriteAccess(requestSession, r)

	localeCode := strings.TrimSpace(chi.URLParam(r, "locale"))
	if !isValidLocaleCode(localeCode) {
		exception.ThrowExceptionBadRequestResponse()
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, i18nDomain.I18nImportResponseDto{
			Success:          false,
			ValidationPassed: false,
			Locale:           localeCode,
			Errors: []i18nDomain.I18nImportIssueDto{{
				Code:    "INVALID_MULTIPART",
				Message: "Upload must be multipart/form-data",
			}},
		})
		return
	}

	var fileHeader *multipart.FileHeader
	if files, ok := r.MultipartForm.File["file"]; ok && len(files) > 0 {
		fileHeader = files[0]
	}

	if fileHeader == nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, i18nDomain.I18nImportResponseDto{
			Success:          false,
			ValidationPassed: false,
			Locale:           localeCode,
			Errors: []i18nDomain.I18nImportIssueDto{{
				Code:    "NO_FILE",
				Message: "No JSON file uploaded",
			}},
		})
		return
	}

	fileReader, err := fileHeader.Open()
	if err != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, i18nDomain.I18nImportResponseDto{
			Success:          false,
			ValidationPassed: false,
			Locale:           localeCode,
			Errors: []i18nDomain.I18nImportIssueDto{{
				FileName: fileHeader.Filename,
				Code:     "FILE_READ_ERROR",
				Message:  "Unable to open uploaded file",
			}},
		})
		return
	}
	content, readErr := io.ReadAll(fileReader)
	_ = fileReader.Close()
	if readErr != nil {
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, i18nDomain.I18nImportResponseDto{
			Success:          false,
			ValidationPassed: false,
			Locale:           localeCode,
			Errors: []i18nDomain.I18nImportIssueDto{{
				FileName: fileHeader.Filename,
				Code:     "FILE_READ_ERROR",
				Message:  "Unable to read uploaded file",
			}},
		})
		return
	}

	parsedEntries, parseIssue := parseLocaleImportJSON(fileHeader.Filename, content)
	if parseIssue != nil {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, i18nDomain.I18nImportResponseDto{
			Success:          false,
			ValidationPassed: false,
			Locale:           localeCode,
			FilesProcessed:   1,
			Errors:           []i18nDomain.I18nImportIssueDto{*parseIssue},
		})
		return
	}

	if len(parsedEntries) == 0 {
		render.JSON(w, r, i18nDomain.I18nImportResponseDto{
			Success:          true,
			ValidationPassed: true,
			Locale:           localeCode,
			FilesProcessed:   1,
			TotalKeysParsed:  0,
			Appended:         0,
			Updated:          0,
			Unchanged:        0,
		})
		return
	}

	existingLocale := handler.I18nRepository.FindByLocaleCode(requestSession, localeCode, false)

	currentUser := roles.GetUsernameFromRequest(requestSession, r)
	if strings.TrimSpace(currentUser) == "" {
		currentUser = "SYSTEM"
	}

	appended, updated, unchanged := 0, 0, 0

	for key, newValue := range parsedEntries {
		if existingLocale != nil {
			if existingEntry := existingLocale.GetEntry(key); existingEntry != nil {
				if existingEntry.Value == newValue {
					unchanged++
					continue
				}
				updated++
			} else {
				appended++
			}
		} else {
			appended++
		}
		handler.I18nRepository.SetTranslation(requestSession, localeCode, key, newValue, "Imported from JSON", currentUser)
	}

	render.JSON(w, r, i18nDomain.I18nImportResponseDto{
		Success:          true,
		ValidationPassed: true,
		Locale:           localeCode,
		FilesProcessed:   1,
		TotalKeysParsed:  len(parsedEntries),
		Appended:         appended,
		Updated:          updated,
		Unchanged:        unchanged,
	})
}

func (handler *I18nHandler) GetTranslationByKey(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	requestedLocale := strings.TrimSpace(chi.URLParam(r, "locale"))
	key := strings.TrimSpace(chi.URLParam(r, "key"))
	if !isValidLocaleCode(requestedLocale) || key == "" {
		exception.ThrowExceptionBadRequestResponse()
	}

	value, actualLocale, found := handler.I18nRepository.GetTranslationWithFallback(requestSession, requestedLocale, key)
	if found {
		render.JSON(w, r, i18nDomain.I18nTranslationResponseDto{
			LocaleCode:   actualLocale,
			RequestedKey: key,
			Value:        value,
			FallbackUsed: actualLocale != requestedLocale,
		})
		return
	}

	// Last-resort frontend-safe fallback: return key itself.
	render.JSON(w, r, i18nDomain.I18nTranslationResponseDto{
		LocaleCode:   requestedLocale,
		RequestedKey: key,
		Value:        key,
		FallbackUsed: true,
	})
}

func (handler *I18nHandler) GetLocales(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	all := handler.I18nRepository.FindAll(requestSession, false)
	result := make([]i18nDomain.I18nLocaleListResponseDto, 0, len(all))
	for _, locale := range all {
		result = append(result, locale.ToListDTO())
	}
	render.JSON(w, r, result)
}

func (handler *I18nHandler) UpsertLocaleMetadata(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	ensureI18nWriteAccess(requestSession, r)

	localeCode := strings.TrimSpace(chi.URLParam(r, "locale"))
	if !isValidLocaleCode(localeCode) {
		exception.ThrowExceptionBadRequestResponse()
	}

	var req i18nDomain.I18nLocaleUpsertRequestDto
	validation.DecodeAndValidate(r, &req, false)
	handler.I18nRepository.UpsertLocaleMetadata(requestSession, localeCode, strings.TrimSpace(req.DisplayName), strings.TrimSpace(req.NativeName), req.IsDefault, strings.TrimSpace(req.Scope))

	render.JSON(w, r, SuccessResponse{Success: true})
}

func (handler *I18nHandler) UpsertTranslationByKey(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	ensureI18nWriteAccess(requestSession, r)

	localeCode := strings.TrimSpace(chi.URLParam(r, "locale"))
	key := strings.TrimSpace(chi.URLParam(r, "key"))
	if !isValidLocaleCode(localeCode) || key == "" {
		exception.ThrowExceptionBadRequestResponse()
	}

	var req i18nDomain.I18nTranslationUpsertRequestDto
	validation.DecodeAndValidate(r, &req, false)
	currentUser := roles.GetUsernameFromRequest(requestSession, r)
	handler.I18nRepository.SetTranslation(requestSession, localeCode, key, req.Value, req.Description, currentUser)

	render.JSON(w, r, SuccessResponse{Success: true})
}

func (handler *I18nHandler) DeleteTranslationByKey(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	ensureI18nWriteAccess(requestSession, r)

	localeCode := strings.TrimSpace(chi.URLParam(r, "locale"))
	key := strings.TrimSpace(chi.URLParam(r, "key"))
	if !isValidLocaleCode(localeCode) || key == "" {
		exception.ThrowExceptionBadRequestResponse()
	}

	handler.I18nRepository.DeleteTranslation(requestSession, localeCode, key)
	render.JSON(w, r, SuccessResponse{Success: true})
}

func (handler *I18nHandler) DeleteLocale(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	ensureI18nWriteAccess(requestSession, r)

	localeCode := strings.TrimSpace(chi.URLParam(r, "locale"))
	if !isValidLocaleCode(localeCode) {
		exception.ThrowExceptionBadRequestResponse()
	}

	if _, isProtected := protectedLocales[strings.ToLower(localeCode)]; isProtected {
		exception.ThrowExceptionClientWithHttpCode(message.BadRequest, "Locale cannot be deleted", localeCode, http.StatusBadRequest)
	}

	if ok := handler.I18nRepository.DeleteLocale(requestSession, localeCode); !ok {
		exception.ThrowExceptionClientWithHttpCode(message.BadRequest, "Default locale cannot be deleted", localeCode, http.StatusBadRequest)
	}

	render.JSON(w, r, SuccessResponse{Success: true})
}
