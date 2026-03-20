// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/eclipse-disuko/disuko/domain/department"
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/infra/repository/user"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap/zapcore"
)

var UserRepository user.IUsersRepository

func Validate(data interface{}, handleErrorAsServerException bool, isDraft *bool) {
	// https://github.com/go-playground/validator
	validate := validator.New(validator.WithRequiredStructEnabled())
	validate.RegisterValidation("NeFieldIfSet", validateNeFieldIfSet)
	validate.RegisterValidation("RealUser", ValidateRealUser)
	validate.RegisterValidation("RealInternalUser", ValidateRealInternalUser)
	validate.RegisterValidation("SupportedUserType", ValidateSupportedUserType)
	validate.RegisterValidation("OmitEmptyStruct", ValidateOmitemptyStruct)
	validate.RegisterValidation("OmitEmptySubStructWith", ValidateOmitEmptySubStructWith)
	validate.RegisterAlias("RequiredIfNotDraft", "required_unless=IsDraft true")
	if isDraft != nil {
		validate.RegisterValidation("GtIfNotDraft", validateGtIfNotDraft(validate, *isDraft))
	}

	err := validate.Struct(data)
	if err != nil {
		if handleErrorAsServerException {
			exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorJsonValidatingForm), err.Error()+fmt.Sprintf(" %v", data), zapcore.InfoLevel)
		} else {
			exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorJsonValidatingForm), err.Error()+fmt.Sprintf(" %v", data), zapcore.InfoLevel)
		}
	}
}

func ValidateList[RESULT interface{}](data *[]RESULT, handleErrorAsServerException bool) {
	for _, member := range *data {
		Validate(member, handleErrorAsServerException, nil)
	}
}

func DecodeAndValidate(r *http.Request, targetDto interface{}, handleErrorAsServerException bool) {
	CheckExpectedContentType(r, ContentTypeJson)
	DecodePartAndValidate(r.Body, targetDto, handleErrorAsServerException)
}

func DecodePartAndValidate(source io.ReadCloser, targetDto interface{}, handleErrorAsServerException bool) {
	DecodePart(source, targetDto)
	Validate(targetDto, handleErrorAsServerException, nil)
}

func DecodePartAndValidateList[RESULT interface{}](source io.ReadCloser, targetsDto *[]RESULT, handleErrorAsServerException bool) {
	DecodePart(source, targetsDto)
	ValidateList(targetsDto, handleErrorAsServerException)
}

func DecodePart(source io.ReadCloser, targetDto any) {
	decoder := json.NewDecoder(source)
	err := decoder.Decode(targetDto)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorJsonDecodingInput))
}

func validateNeFieldIfSet(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	referencedByParamFieldValue := reflect.Indirect(fl.Parent()).FieldByName(fl.Param()).String()
	return len(value) == 0 && len(referencedByParamFieldValue) == 0 || value != referencedByParamFieldValue
}

func ValidateRealInternalUser(fl validator.FieldLevel) bool {
	internal := true
	return validateUser(fl.Field().String(), &internal)
}

func ValidateRealUser(fl validator.FieldLevel) bool {
	return validateUser(fl.Field().String(), nil)
}

func validateUser(name string, internal *bool) bool {
	if len(name) == 0 {
		return true
	}
	requestSession := &logy.RequestSession{ReqID: "VALIDATION-" + uuid.NewString()}
	if user := UserRepository.FindByUserId(requestSession, name); user != nil && len(user.User) > 0 {
		if internal == nil {
			return true
		}
		return user.IsInternal == *internal
	}
	logy.Errorf(requestSession, "User Validation - could not find user %s in database", name)
	return false
}

func ValidateSupportedUserType(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(project.UserType)
	return value == project.OWNER || value == project.SUPPLIER || value == project.VIEWER
}

func ValidateOmitemptyStruct(fl validator.FieldLevel) bool {
	value := fl.Field().Interface().(department.DepartmentDto)
	return len(value.DeptId) > 0
}

func ValidateOmitEmptySubStructWith(fl validator.FieldLevel) bool {
	parts := strings.SplitN(fl.Param(), " ", 2)
	pathsWithValue := strings.SplitN(parts[1], "=", 2)
	paths := strings.Split(pathsWithValue[0], ".")
	var valueOfNode bool
	boolNodeValue := reflect.Indirect(fl.Parent())
	for _, path := range paths {
		boolNodeValue = reflect.Indirect(boolNodeValue.FieldByName(path))
	}
	boolType := reflect.TypeOf(true)
	boolNodeType := boolNodeValue.Type()
	if boolNodeType != boolType {
		return false
	}
	valueOfNode = boolNodeValue.Bool()
	requiredValue, err := strconv.ParseBool(pathsWithValue[1])
	if err != nil {
		return false
	}

	if valueOfNode == requiredValue {
		subStructValue := reflect.Indirect(reflect.Indirect(fl.Field()).FieldByName(parts[0]))
		if !subStructValue.IsValid() {
			return false
		}
		valueOfSubStruct := subStructValue.Interface().(project.Department)
		return len(valueOfSubStruct.DeptId) > 0
	} else {
		return true
	}
}

func validateGtIfNotDraft(validate *validator.Validate, isDraft bool) func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		if isDraft {
			return true
		}
		return validate.Var(fl.Field().Interface(), "gt="+fl.Param()) == nil
	}
}

func validateRequiredIfNotDraft(validate *validator.Validate, isDraft bool) func(fl validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		if isDraft {
			return true
		}
		return validate.Var(fl.Field().Interface(), "required") == nil
	}
}
