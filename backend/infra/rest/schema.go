// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/schema"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/helper/roles"
	"github.com/eclipse-disuko/disuko/helper/validation"
	"github.com/eclipse-disuko/disuko/infra/repository/labels"
	schemaRepository "github.com/eclipse-disuko/disuko/infra/repository/schema"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/eclipse-disuko/disuko/observermngmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/xeipuuv/gojsonschema"
)

type SchemaHandler struct {
	SchemaRepository schemaRepository.ISchemaRepository
	LabelRepository  labels.ILabelRepository
}

func (schemaHandler *SchemaHandler) SchemaActivateHandler(w http.ResponseWriter, r *http.Request) {
	schemaToActivate, requestSession := retrieveSchema(schemaHandler.SchemaRepository, w, r)

	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowSchema.Update {
		exception.ThrowExceptionSendDeniedResponse()
	}

	checkIfSchemaJSONIsParsable(schemaToActivate.Content)

	activeSchemas := schemaHandler.SchemaRepository.FindActiveSchemas(requestSession)
	schemaToDeactivate := &schema.SpdxSchema{}
	for _, v := range activeSchemas {
		if (len(schemaToActivate.Label) > 0 && schemaToActivate.Label == v.Label) || (len(schemaToActivate.Label) == 0 && len(v.Label) == 0) {
			schemaToDeactivate = v
			break
		}
	}

	schemaToActivate.ActivateSchema()

	var schemasToSave []*schema.SpdxSchema
	if schemaHandler.schemaFound(schemaToDeactivate) {
		schemaToDeactivate.DeactivateSchema()
		schemasToSave = []*schema.SpdxSchema{schemaToActivate, schemaToDeactivate}
	} else {
		schemasToSave = []*schema.SpdxSchema{schemaToActivate}
	}

	schemaHandler.SchemaRepository.UpdateList(requestSession, schemasToSave)
	result := schema.StatusResponseDto{
		Success: true,
		Message: "Schema successfully activated",
	}
	render.JSON(w, r, result)
}

func checkIfSchemaJSONIsParsable(schemaContent string) {
	// or use gojsonschema.NewStringLoader() when the content is in string
	schemaLoader := gojsonschema.NewStringLoader("{}")
	spdxLoader := gojsonschema.NewStringLoader(schemaContent)
	_, err := gojsonschema.Validate(schemaLoader, spdxLoader)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ParseSchema))
}

func (schemaHandler *SchemaHandler) schemaFound(schema *schema.SpdxSchema) bool {
	return len(schema.Key) > 0
}

func (schemaHandler *SchemaHandler) SchemaGetHandler(w http.ResponseWriter, r *http.Request) {
	schemaObject, _ := retrieveSchema(schemaHandler.SchemaRepository, w, r)
	render.JSON(w, r, schemaObject.ToDto())
}

func (schemaHandler *SchemaHandler) SchemaGetAllHandler(w http.ResponseWriter, r *http.Request) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowSchema.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}

	schemas := schemaHandler.SchemaRepository.FindAll(requestSession, true)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	render.JSON(w, r, schema.ToSpdxSchemaDtoList(schemas))
}

func (schemaHandler *SchemaHandler) SchemaUploadHandler(w http.ResponseWriter, r *http.Request) {
	validation.CheckExpectedContentType(r, validation.ContentTypeFormData)

	requestSession := logy.GetRequestSession(r)
	user, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowSchema.Create {
		exception.ThrowExceptionSendDeniedResponse()
	}

	mr, err := r.MultipartReader()
	exception.HandleErrorClientMessage(err, message.GetI18N(message.MultipartReader))

	schemaRequest := schema.SchemaRequestDto{}
	var contents string
	for {
		part, err := mr.NextPart()

		if err == io.EOF {
			break
		}
		exception.HandleErrorClientMessage(err, message.GetI18N(message.PartReader))
		if part.FormName() == "file" {
			var buf bytes.Buffer
			_, err = io.Copy(&buf, part)
			exception.HandleErrorClientMessage(err, message.GetI18N(message.CopyTempFile))
			contents = buf.String()
			validation.CheckExpectedContentType2(part.Header, []validation.ContentType{
				validation.ContentTypeJson,
				validation.ContentTypeOctets,
			})
		}

		if part.FormName() == "schema" {
			validation.DecodePartAndValidate(part, &schemaRequest, false)
		}
	}

	if len(contents) == 0 {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.NoSchemaProvided))
	}

	currentSchema := schemaHandler.SchemaRepository.FindSpdxSchemaByNameAndVersion(requestSession, strings.TrimSpace(schemaRequest.Name), strings.TrimSpace(schemaRequest.Version))

	if currentSchema != nil {
		exception.ThrowExceptionClientMessage3(message.GetI18N(message.FoundExistingSchema))
	}

	_, err = gojsonschema.NewStringLoader(contents).LoadJSON()

	exception.HandleErrorClientMessage(err, message.GetI18N(message.ValidateSchema))

	// validate label
	CheckIfLabelExistOrThrowException(requestSession, schemaHandler.LabelRepository, schemaRequest.Label)

	newSchema := &schema.SpdxSchema{
		RootEntity:  domain.NewRootEntity(),
		Content:     contents,
		Type:        schema.JSON,
		Active:      false,
		Name:        strings.TrimSpace(schemaRequest.Name),
		Version:     strings.TrimSpace(schemaRequest.Version),
		Description: strings.TrimSpace(schemaRequest.Description),
		Label:       strings.TrimSpace(schemaRequest.Label),
	}

	schemaHandler.SchemaRepository.Save(requestSession, newSchema)

	observermngmt.FireEvent(observermngmt.DatabaseEntryAddedOrDeleted, observermngmt.DatabaseSizeChange{
		RequestSession: requestSession,
		CollectionName: schemaRepository.SpdxSchemaCollectionName,
		Rights:         rights,
		Username:       user,
	})

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	_, err = w.Write([]byte(newSchema.Key))
	exception.HandleErrorClientMessage(err, message.GetI18N(message.WritingContent))
	w.WriteHeader(201)
}

func (schemaHandler *SchemaHandler) SchemaDownloadHandler(w http.ResponseWriter, r *http.Request) {
	schemaObject, _ := retrieveSchema(schemaHandler.SchemaRepository, w, r)

	_, err := w.Write([]byte(schemaObject.Content))
	exception.HandleErrorClientMessage(err, message.GetI18N(message.WritingContent))
	w.Header().Set("Content-Type", schemaObject.ContentTypeAsString())
	w.WriteHeader(200)
}

func retrieveSchema(schemaRepo schemaRepository.ISchemaRepository, w http.ResponseWriter, r *http.Request) (*schema.SpdxSchema, *logy.RequestSession) {
	requestSession := logy.GetRequestSession(r)
	_, rights := roles.GetAccessAndRolesRightsFromRequest(requestSession, r)
	if !rights.AllowSchema.Read {
		exception.ThrowExceptionSendDeniedResponse()
	}

	schemaId := chi.URLParam(r, "id")

	schemaObject := schemaRepo.FindByKey(requestSession, schemaId, false)

	if schemaObject == nil {
		exception.ThrowExceptionClient404Message3(message.GetI18N(message.SchemaNotFound, schemaId))
	}
	return schemaObject, requestSession
}
