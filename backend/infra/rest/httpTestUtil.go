// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package rest

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strings"
	"testing"

	"github.com/eclipse-disuko/disuko/domain/user"
	"github.com/eclipse-disuko/disuko/helper/jwt"
	"github.com/eclipse-disuko/disuko/helper/validation"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
)

var requestSessionTest = &logy.RequestSession{ReqID: "TEST"}

func RecordHTTPRequest(t *testing.T, method string, path string, body []byte, handler http.HandlerFunc, params map[string]string) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, path, bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	ctx := chi.NewRouteContext()
	if params != nil {
		for param, value := range params {
			ctx.URLParams.Add(param, value)
		}
	}
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	req.RemoteAddr = "localhost"
	userData := jwt.CreateUserData(requestSessionTest, &user.User{TermsOfUse: false}, "DEMO", "demoMailAddress", strings.Split("FOSSDP.domain_admin", jwt.GROUPS_TOKEN), jwt.GROUP_TYPE_DAIMLER, true, req)
	tokenDetails := jwt.CreateToken(userData)

	req.Header.Set("Content-Type", string(validation.ContentTypeJson))
	req.Header.Set("Authorization", "BEARER "+tokenDetails.AccessToken)
	//req.Header.Set(middleware.RequestIDHeader, requestSessionTest.ReqID)

	recorder := httptest.NewRecorder()

	httpServer := http.HandlerFunc(handler)

	httpServer.ServeHTTP(recorder, req.WithContext(context.WithValue(req.Context(), middleware.RequestIDKey, requestSessionTest.ReqID)))
	return recorder
}

func CheckHTTPCode(t *testing.T, recorder *httptest.ResponseRecorder, expectedCode int) bool {
	return assert.True(t, recorder.Code == expectedCode, "handler returned unexpected code %d, should be %d", recorder.Code, expectedCode)
}

func RecordHTTPMultipartRequest(t *testing.T, method string, path string, postData string, handler http.HandlerFunc, params map[string]string) *httptest.ResponseRecorder {
	reader := strings.NewReader(postData)

	b := bytes.Buffer{} // buffer to write the request payload into
	writer := multipart.NewWriter(&b)
	part, _ := CreateFormFile(writer, "file", "somefile.json")
	_, err := io.Copy(part, reader)
	if err != nil {
		t.Fatal(err)
	}
	err = writer.Close()
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest(method, path, bytes.NewReader(b.Bytes()))
	if err != nil {
		t.Fatal(err)
	}

	ctx := chi.NewRouteContext()
	if params != nil {
		for param, value := range params {
			ctx.URLParams.Add(param, value)
		}
	}

	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	req.RemoteAddr = "localhost"
	userData := jwt.CreateUserData(requestSessionTest, &user.User{TermsOfUse: false}, "DEMO", "demoMailAddress", strings.Split("FOSSDP.domain_admin",
		jwt.GROUPS_TOKEN), jwt.GROUP_TYPE_DAIMLER, true, req)
	tokenDetails := jwt.CreateToken(userData)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "BEARER "+tokenDetails.AccessToken)

	recorder := httptest.NewRecorder()

	httpServer := http.HandlerFunc(handler)
	httpServer.ServeHTTP(recorder, req.WithContext(context.WithValue(req.Context(), middleware.RequestIDKey, requestSessionTest.ReqID)))
	return recorder
}

func CreateFormFile(mw *multipart.Writer, fieldname, filename string) (io.Writer, error) {
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition",
		fmt.Sprintf(`form-data; name="%s"; filename="%s"`,
			escapeQuotes(fieldname), escapeQuotes(filename)))
	h.Set("Content-Type", string(validation.ContentTypeJson))
	return mw.CreatePart(h)
}

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}
