// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
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

func parseLocaleImportJSON(fileName string, payload []byte) (map[string]string, []i18nDomain.I18nImportIssueDto) {
	issues := make([]i18nDomain.I18nImportIssueDto, 0)
	result := make(map[string]string)

	decoder := json.NewDecoder(bytes.NewReader(payload))
	token, err := decoder.Token()
	if err != nil {
		issues = append(issues, i18nDomain.I18nImportIssueDto{
			FileName: fileName,
			Code:     "INVALID_JSON",
			Message:  "File contains invalid JSON",
		})
		return nil, issues
	}

	delim, ok := token.(json.Delim)
	if !ok || delim != '{' {
		issues = append(issues, i18nDomain.I18nImportIssueDto{
			FileName: fileName,
			Code:     "INVALID_ROOT",
			Message:  "JSON root must be an object",
		})
		return nil, issues
	}

	seenKeys := make(map[string]struct{})
	for decoder.More() {
		keyToken, err := decoder.Token()
		if err != nil {
			issues = append(issues, i18nDomain.I18nImportIssueDto{
				FileName: fileName,
				Code:     "INVALID_JSON",
				Message:  "Failed to parse JSON key",
			})
			return nil, issues
		}

		key, ok := keyToken.(string)
		if !ok {
			issues = append(issues, i18nDomain.I18nImportIssueDto{
				FileName: fileName,
				Code:     "INVALID_KEY",
				Message:  "JSON object keys must be strings",
			})
			return nil, issues
		}

		if _, exists := seenKeys[key]; exists {
			issues = append(issues, i18nDomain.I18nImportIssueDto{
				FileName: fileName,
				Key:      key,
				Code:     "DUPLICATE_KEY",
				Message:  "Duplicate key found in JSON file",
			})
		}
		seenKeys[key] = struct{}{}

		var rawValue json.RawMessage
		if err := decoder.Decode(&rawValue); err != nil {
			issues = append(issues, i18nDomain.I18nImportIssueDto{
				FileName: fileName,
				Key:      key,
				Code:     "INVALID_VALUE",
				Message:  "Failed to parse JSON value",
			})
			continue
		}

		var textValue string
		if err := json.Unmarshal(rawValue, &textValue); err != nil {
			issues = append(issues, i18nDomain.I18nImportIssueDto{
				FileName: fileName,
				Key:      key,
				Code:     "UNSUPPORTED_VALUE_TYPE",
				Message:  "Only string values are supported",
			})
			continue
		}

		result[key] = textValue
	}

	endToken, err := decoder.Token()
	if err != nil {
		issues = append(issues, i18nDomain.I18nImportIssueDto{
			FileName: fileName,
			Code:     "INVALID_JSON",
			Message:  "Failed to close JSON object",
		})
		return nil, issues
	}

	endDelim, ok := endToken.(json.Delim)
	if !ok || endDelim != '}' {
		issues = append(issues, i18nDomain.I18nImportIssueDto{
			FileName: fileName,
			Code:     "INVALID_JSON",
			Message:  "Invalid JSON object ending",
		})
	}

	if decoder.More() {
		issues = append(issues, i18nDomain.I18nImportIssueDto{
			FileName: fileName,
			Code:     "INVALID_JSON",
			Message:  "Unexpected trailing JSON tokens",
		})
	}

	if len(issues) > 0 {
		return nil, issues
	}

	return result, nil
}

func marshalSortedLocaleEntries(entries map[string]string) ([]byte, error) {
	keys := make([]string, 0, len(entries))
	for key := range entries {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	buffer := bytes.NewBufferString("{")
	for index, key := range keys {
		if index > 0 {
			buffer.WriteString(",")
		}

		encodedKey, err := json.Marshal(key)
		if err != nil {
			return nil, err
		}
		encodedValue, err := json.Marshal(entries[key])
		if err != nil {
			return nil, err
		}

		buffer.Write(encodedKey)
		buffer.WriteString(":")
		buffer.Write(encodedValue)
	}
	buffer.WriteString("}")

	return buffer.Bytes(), nil
}

func flattenI18nEntries(prefix string, input map[string]interface{}, out map[string]string) {
	for key, value := range input {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		switch typed := value.(type) {
		case string:
			out[fullKey] = typed
		case map[string]interface{}:
			flattenI18nEntries(fullKey, typed, out)
		default:
			out[fullKey] = fmt.Sprint(typed)
		}
	}
}

func readLocaleJSON(filePath string) (map[string]string, error) {
	raw, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	parsed := make(map[string]interface{})
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, err
	}

	result := make(map[string]string)
	flattenI18nEntries("", parsed, result)
	return result, nil
}

func localeNames(localeCode string) (string, string) {
	switch strings.ToLower(localeCode) {
	case "en":
		return "English", "English"
	case "de":
		return "German", "Deutsch"
	default:
		return strings.ToUpper(localeCode), strings.ToUpper(localeCode)
	}
}

func ensureI18nWriteAccess(requestSession *logy.RequestSession, r *http.Request) {
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.IsApplicationAdmin() && !rights.IsDomainAdmin() {
		exception.ThrowExceptionSendDeniedResponse()
	}
}

func (handler *I18nHandler) findDefaultLocale(requestSession *logy.RequestSession) (string, bool) {
	allLocales := handler.I18nRepository.FindAll(requestSession, false)
	for _, locale := range allLocales {
		if locale != nil && locale.IsDefault {
			return locale.LocaleCode, true
		}
	}
	return "", false
}

func (handler *I18nHandler) GetLocale(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	requestedLocale := strings.TrimSpace(chi.URLParam(r, "locale"))
	if requestedLocale == "" {
		exception.ThrowExceptionBadRequestResponse()
	}

	locale := handler.I18nRepository.FindByLocaleCode(requestSession, requestedLocale, false)
	fallbackUsed := false
	if locale == nil {
		defaultLocaleCode, ok := handler.findDefaultLocale(requestSession)
		if !ok {
			exception.ThrowExceptionClient404Message(message.GetI18N(message.ErrorDbNotFound), "i18n locale not found: "+requestedLocale)
		}
		locale = handler.I18nRepository.FindByLocaleCode(requestSession, defaultLocaleCode, false)
		fallbackUsed = true
	}
	if locale == nil {
		exception.ThrowExceptionClient404Message(message.GetI18N(message.ErrorDbNotFound), "default i18n locale not found")
	}

	entries := make(map[string]string)
	for key, value := range locale.Entries {
		if value != nil {
			entries[key] = value.Value
		}
	}

	render.JSON(w, r, i18nDomain.I18nLocaleResponseDto{
		LocaleCode:   locale.LocaleCode,
		DisplayName:  locale.DisplayName,
		NativeName:   locale.NativeName,
		IsDefault:    locale.IsDefault,
		Scope:        locale.Scope,
		EntryCount:   len(entries),
		Entries:      entries,
		FallbackUsed: fallbackUsed,
	})
}

func (handler *I18nHandler) ExportLocaleJSON(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	localeCode := strings.TrimSpace(chi.URLParam(r, "locale"))
	if localeCode == "" {
		exception.ThrowExceptionBadRequestResponse()
	}

	locale := handler.I18nRepository.FindByLocaleCode(requestSession, localeCode, false)
	if locale == nil {
		exception.ThrowExceptionClient404Message(message.GetI18N(message.ErrorDbNotFound), "i18n locale not found: "+localeCode)
	}

	entries := make(map[string]string)
	for key, value := range locale.Entries {
		if value != nil {
			entries[key] = value.Value
		}
	}

	body, err := marshalSortedLocaleEntries(entries)
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
	if localeCode == "" {
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

	fileHeaders := make([]*multipart.FileHeader, 0)
	if r.MultipartForm != nil {
		if files, ok := r.MultipartForm.File["files"]; ok {
			fileHeaders = append(fileHeaders, files...)
		}
		if files, ok := r.MultipartForm.File["file"]; ok {
			fileHeaders = append(fileHeaders, files...)
		}
	}

	if len(fileHeaders) == 0 {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, i18nDomain.I18nImportResponseDto{
			Success:          false,
			ValidationPassed: false,
			Locale:           localeCode,
			Errors: []i18nDomain.I18nImportIssueDto{{
				Code:    "NO_FILES",
				Message: "No JSON files uploaded",
			}},
		})
		return
	}

	issues := make([]i18nDomain.I18nImportIssueDto, 0)
	mergedEntries := make(map[string]string)
	keyToFile := make(map[string]string)
	totalKeysParsed := 0

	for _, fileHeader := range fileHeaders {
		fileReader, err := fileHeader.Open()
		if err != nil {
			issues = append(issues, i18nDomain.I18nImportIssueDto{
				FileName: fileHeader.Filename,
				Code:     "FILE_READ_ERROR",
				Message:  "Unable to open uploaded file",
			})
			continue
		}

		content, readErr := io.ReadAll(fileReader)
		_ = fileReader.Close()
		if readErr != nil {
			issues = append(issues, i18nDomain.I18nImportIssueDto{
				FileName: fileHeader.Filename,
				Code:     "FILE_READ_ERROR",
				Message:  "Unable to read uploaded file",
			})
			continue
		}

		parsedEntries, parseIssues := parseLocaleImportJSON(fileHeader.Filename, content)
		if len(parseIssues) > 0 {
			issues = append(issues, parseIssues...)
			continue
		}

		totalKeysParsed += len(parsedEntries)
		for key, value := range parsedEntries {
			if sourceFile, exists := keyToFile[key]; exists {
				issues = append(issues, i18nDomain.I18nImportIssueDto{
					FileName: fileHeader.Filename,
					Key:      key,
					Code:     "DUPLICATE_KEY_ACROSS_FILES",
					Message:  "Key is duplicated across uploaded files (also in " + sourceFile + ")",
				})
				continue
			}

			keyToFile[key] = fileHeader.Filename
			mergedEntries[key] = value
		}
	}

	if len(issues) > 0 {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, i18nDomain.I18nImportResponseDto{
			Success:          false,
			ValidationPassed: false,
			Locale:           localeCode,
			FilesProcessed:   len(fileHeaders),
			TotalKeysParsed:  totalKeysParsed,
			Errors:           issues,
		})
		return
	}

	if len(mergedEntries) == 0 {
		render.JSON(w, r, i18nDomain.I18nImportResponseDto{
			Success:          true,
			ValidationPassed: true,
			Locale:           localeCode,
			FilesProcessed:   len(fileHeaders),
			TotalKeysParsed:  0,
			Appended:         0,
			Updated:          0,
			Unchanged:        0,
		})
		return
	}

	existingLocale := handler.I18nRepository.FindByLocaleCode(requestSession, localeCode, false)
	existingEntries := make(map[string]string)
	if existingLocale != nil {
		for key, value := range existingLocale.Entries {
			if value != nil {
				existingEntries[key] = value.Value
			}
		}
	}

	currentUser := roles.GetUsernameFromRequest(requestSession, r)
	if strings.TrimSpace(currentUser) == "" {
		currentUser = "SYSTEM"
	}

	appended := 0
	updated := 0
	unchanged := 0

	keys := make([]string, 0, len(mergedEntries))
	for key := range mergedEntries {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		newValue := mergedEntries[key]
		if oldValue, exists := existingEntries[key]; exists {
			if oldValue == newValue {
				unchanged++
				continue
			}
			updated++
		} else {
			appended++
		}

		handler.I18nRepository.SetTranslation(requestSession, localeCode, key, newValue, "Imported from JSON", currentUser)
	}

	render.JSON(w, r, i18nDomain.I18nImportResponseDto{
		Success:          true,
		ValidationPassed: true,
		Locale:           localeCode,
		FilesProcessed:   len(fileHeaders),
		TotalKeysParsed:  totalKeysParsed,
		Appended:         appended,
		Updated:          updated,
		Unchanged:        unchanged,
	})
}

func (handler *I18nHandler) GetTranslationByKey(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	requestedLocale := strings.TrimSpace(chi.URLParam(r, "locale"))
	key := strings.TrimSpace(chi.URLParam(r, "key"))
	if requestedLocale == "" || key == "" {
		exception.ThrowExceptionBadRequestResponse()
	}

	if value, ok := handler.I18nRepository.GetTranslation(requestSession, requestedLocale, key); ok {
		render.JSON(w, r, i18nDomain.I18nTranslationResponseDto{
			LocaleCode:   requestedLocale,
			RequestedKey: key,
			Value:        value,
			FallbackUsed: false,
		})
		return
	}

	defaultLocaleCode, foundDefault := handler.findDefaultLocale(requestSession)
	if foundDefault {
		if value, ok := handler.I18nRepository.GetTranslation(requestSession, defaultLocaleCode, key); ok {
			render.JSON(w, r, i18nDomain.I18nTranslationResponseDto{
				LocaleCode:   defaultLocaleCode,
				RequestedKey: key,
				Value:        value,
				FallbackUsed: true,
			})
			return
		}
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
		if locale == nil {
			continue
		}
		entryCount := 0
		if locale.Entries != nil {
			entryCount = len(locale.Entries)
		}
		result = append(result, i18nDomain.I18nLocaleListResponseDto{
			LocaleCode:  locale.LocaleCode,
			DisplayName: locale.DisplayName,
			NativeName:  locale.NativeName,
			IsDefault:   locale.IsDefault,
			Scope:       locale.Scope,
			EntryCount:  entryCount,
		})
	}
	render.JSON(w, r, result)
}

func (handler *I18nHandler) UpsertLocaleMetadata(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	ensureI18nWriteAccess(requestSession, r)

	localeCode := strings.TrimSpace(chi.URLParam(r, "locale"))
	if localeCode == "" {
		exception.ThrowExceptionBadRequestResponse()
	}

	var req i18nDomain.I18nLocaleUpsertRequestDto
	validation.DecodeAndValidate(r, &req, false)
	handler.I18nRepository.SetLocaleMetadata(requestSession, localeCode, strings.TrimSpace(req.DisplayName), strings.TrimSpace(req.NativeName), req.IsDefault, strings.TrimSpace(req.Scope))

	render.JSON(w, r, SuccessResponse{Success: true})
}

func (handler *I18nHandler) UpsertTranslationByKey(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	ensureI18nWriteAccess(requestSession, r)

	localeCode := strings.TrimSpace(chi.URLParam(r, "locale"))
	key := strings.TrimSpace(chi.URLParam(r, "key"))
	if localeCode == "" || key == "" {
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
	if localeCode == "" || key == "" {
		exception.ThrowExceptionBadRequestResponse()
	}

	handler.I18nRepository.DeleteTranslation(requestSession, localeCode, key)
	render.JSON(w, r, SuccessResponse{Success: true})
}

func (handler *I18nHandler) DeleteLocale(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	ensureI18nWriteAccess(requestSession, r)

	localeCode := strings.TrimSpace(chi.URLParam(r, "locale"))
	if localeCode == "" {
		exception.ThrowExceptionBadRequestResponse()
	}

	normalizedLocaleCode := strings.ToLower(localeCode)
	if normalizedLocaleCode == "en" || normalizedLocaleCode == "de" {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"message": "Locale cannot be deleted"})
		return
	}

	ok := handler.I18nRepository.DeleteLocale(requestSession, localeCode)
	if !ok {
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, map[string]string{"message": "Default locale cannot be deleted"})
		return
	}

	render.JSON(w, r, SuccessResponse{Success: true})
}

func (handler *I18nHandler) MigrateFromJSON(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	ensureI18nWriteAccess(requestSession, r)

	username := roles.GetUsernameFromRequest(requestSession, r)
	if username == "" {
		username = "SYSTEM"
	}

	dryRun := true
	if rawDryRun := strings.TrimSpace(r.URL.Query().Get("dryRun")); rawDryRun != "" {
		parsed, err := strconv.ParseBool(rawDryRun)
		if err != nil {
			exception.ThrowExceptionBadRequestResponse()
		}
		dryRun = parsed
	}

	includeShared := true
	if rawIncludeShared := strings.TrimSpace(r.URL.Query().Get("includeShared")); rawIncludeShared != "" {
		parsed, err := strconv.ParseBool(rawIncludeShared)
		if err != nil {
			exception.ThrowExceptionBadRequestResponse()
		}
		includeShared = parsed
	}

	scope := strings.TrimSpace(r.URL.Query().Get("scope"))
	if scope == "" {
		scope = "portal"
	}

	baseDir := filepath.Clean("../frontend/libs")
	sourceFolders := []string{filepath.Join(baseDir, scope, "i18n", "locales")}
	if includeShared {
		sourceFolders = append(sourceFolders, filepath.Join(baseDir, "shared", "i18n", "locales"))
	}

	entriesByLocale := make(map[string]map[string]string)
	for _, folder := range sourceFolders {
		matches, err := filepath.Glob(filepath.Join(folder, "*.json"))
		if err != nil {
			exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorUnexpectError), err)
		}
		for _, filePath := range matches {
			localeCode := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
			parsedEntries, err := readLocaleJSON(filePath)
			if err != nil {
				exception.ThrowExceptionServerMessageWithError(message.GetI18N(message.ErrorUnexpectError), err)
			}
			if _, ok := entriesByLocale[localeCode]; !ok {
				entriesByLocale[localeCode] = make(map[string]string)
			}
			for key, value := range parsedEntries {
				entriesByLocale[localeCode][key] = value
			}
		}
	}

	if len(entriesByLocale) == 0 {
		exception.ThrowExceptionClient404Message(message.GetI18N(message.ErrorDbNotFound), "no locale JSON files found for migration")
	}

	results := make([]i18nDomain.I18nMigrationLocaleResultDto, 0, len(entriesByLocale))
	for localeCode, entries := range entriesByLocale {
		displayName, nativeName := localeNames(localeCode)
		sourceCount := len(entries)
		upserted := 0

		if !dryRun {
			handler.I18nRepository.SetLocaleMetadata(requestSession, localeCode, displayName, nativeName, strings.EqualFold(localeCode, "en"), scope)
			for key, value := range entries {
				handler.I18nRepository.SetTranslation(requestSession, localeCode, key, value, "Migrated from JSON", username)
				upserted++
			}
		}

		targetCount := handler.I18nRepository.GetEntryCountForLocale(requestSession, localeCode)
		if dryRun {
			upserted = sourceCount
		}

		verified := targetCount >= sourceCount
		if dryRun {
			verified = true
		}

		results = append(results, i18nDomain.I18nMigrationLocaleResultDto{
			LocaleCode:  localeCode,
			SourceCount: sourceCount,
			TargetCount: targetCount,
			Upserted:    upserted,
			Verified:    verified,
		})
	}

	response := i18nDomain.I18nMigrationResponseDto{
		DryRun:        dryRun,
		Scope:         scope,
		IncludeShared: includeShared,
		Locales:       results,
		Success:       true,
	}
	render.JSON(w, r, response)
}
