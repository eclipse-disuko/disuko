// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package roles

import (
	"net/http"
	"testing"

	"github.com/eclipse-disuko/disuko/domain/user"

	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/helper/jwt"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/helper/test"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/stretchr/testify/assert"
)

var requestSessionTest = &logy.RequestSession{ReqID: "TEST"}

func Test_isExternalUser(t *testing.T) {

	//create token
	tokenData := jwt.TokenData{
		Username:           "testUser",
		Groups:             "testGroup1,TestGroup2",
		Email:              "email@test.de",
		RemoteAddress:      "localhost",
		IsInternalEmployee: false,
	}
	token := jwt.CreateToken(tokenData)
	assert.NotNil(t, token)
	rights := GetAccessAndRolesRightsFromClaim(tokenData)

	assert.Equal(t, rights.AllowProject.Create, false)
	assert.Equal(t, rights.AllowProject.Read, false)
	assert.Equal(t, rights.AllowProject.Update, false)
	assert.Equal(t, rights.AllowProject.Delete, false)
}

func Test_isInternalUser(t *testing.T) {

	//create token
	tokenData := jwt.TokenData{
		Username:           "testUser",
		Groups:             "testGroup1,TestGroup2",
		Email:              "email@test.de",
		RemoteAddress:      "localhost",
		IsInternalEmployee: true,
	}
	token := jwt.CreateToken(tokenData)
	assert.NotNil(t, token)
	rights := GetAccessAndRolesRightsFromClaim(tokenData)

	assert.Equal(t, rights.AllowProject.Create, true)
	assert.Equal(t, rights.AllowProject.Read, false)
	assert.Equal(t, rights.AllowProject.Update, false)
	assert.Equal(t, rights.AllowProject.Delete, false)
}

func Test_trimPortFromRemoteAddress(t *testing.T) {
	assert.NotEqual(t, "localhost:8080", jwt.TrimPortFromRemoteAddress("localhost:8080"))

	assert.Equal(t, "localhost", jwt.TrimPortFromRemoteAddress("localhost"))
	assert.Equal(t, "localhost", jwt.TrimPortFromRemoteAddress("localhost:80"))
	assert.Equal(t, "localhost", jwt.TrimPortFromRemoteAddress("localhost:"))

	assert.Equal(t, "192.18.23.4", jwt.TrimPortFromRemoteAddress("192.18.23.4:80"))

	assert.Equal(t, "::", jwt.TrimPortFromRemoteAddress("[::]:8560"))
	assert.Equal(t, "2a01:598:9995:eea6:c186:989b:acc9:18ae",
		jwt.TrimPortFromRemoteAddress("[2a01:598:9995:eea6:c186:989b:acc9:18ae]:8560"))
	assert.Equal(t, "2a01:598:9995:eea6:c186:989b:acc9:18ae",
		jwt.TrimPortFromRemoteAddress("2a01:598:9995:eea6:c186:989b:acc9:18ae"))
}

func Test_CreateAndVerifyToken_Success(t *testing.T) {

	//create token
	tokenData := jwt.TokenData{
		Username:      "testUser",
		Groups:        "testGroup1,TestGroup2",
		Email:         "email@test.de",
		RemoteAddress: "localhost",
	}
	token := jwt.CreateToken(tokenData)
	assert.NotNil(t, token)

	//create http request
	httpRequest := &http.Request{
		Header: http.Header{},
	}
	httpRequest.Header.Add("Authorization", "Bearer "+token.AccessToken)
	httpRequest.RemoteAddr = "localhost"

	//verify token
	metadata := jwt.ExtractTokenMetadata(requestSessionTest, httpRequest)
	assert.NotNil(t, metadata)
	assert.Equal(t, &tokenData, metadata)
}

func Test_CreateAndVerifyToken_FailedSigning(t *testing.T) {

	//create token
	tokenData := jwt.TokenData{
		Username:      "testUser",
		Groups:        "testGroup1,TestGroup2",
		Email:         "email@test.de",
		RemoteAddress: "localhost",
	}
	token := jwt.CreateToken(tokenData)
	assert.NotNil(t, token)

	//create http request
	httpRequest := &http.Request{
		Header: http.Header{},
	}
	httpRequest.Header.Add("Authorization", "Bearer "+token.AccessToken)
	httpRequest.RemoteAddr = "localhost"

	//verify token
	//change secrect
	orgSecrect := conf.Config.Auth.AccessSecret
	conf.Config.Auth.AccessSecret = "adsdadwqe"
	test.ExpectException(t, message.ErrorAAR, func() {
		jwt.ExtractTokenMetadata(requestSessionTest, httpRequest)
	})

	//reset secret
	conf.Config.Auth.AccessSecret = orgSecrect
}

func Test_CreateAndVerifyToken_FailedWrongRemoteAddress(t *testing.T) {

	//create token
	tokenData := jwt.TokenData{
		Username:      "testUser",
		Groups:        "testGroup1,TestGroup2",
		Email:         "email@test.de",
		RemoteAddress: "localhost",
	}
	token := jwt.CreateToken(tokenData)
	assert.NotNil(t, token)

	//create http request
	httpRequest := &http.Request{
		Header: http.Header{},
	}
	httpRequest.Header.Add("Authorization", "Bearer "+token.AccessToken)
	httpRequest.RemoteAddr = "localhost2"

	//verify token
	orgPreventTokenHiJacking := conf.Config.OAuth2.PreventTokenHiJacking
	conf.Config.OAuth2.PreventTokenHiJacking = true
	test.ExpectException(t, message.ErrorAAR, func() {
		jwt.ExtractTokenMetadata(requestSessionTest, httpRequest)
	})
	conf.Config.OAuth2.PreventTokenHiJacking = orgPreventTokenHiJacking
}

func Test_CreateAndVerifyToken_SuccessWrongRemoteAddress(t *testing.T) {

	//create token
	tokenData := jwt.TokenData{
		Username:      "testUser",
		Groups:        "testGroup1,TestGroup2",
		Email:         "email@test.de",
		RemoteAddress: "localhost",
	}
	token := jwt.CreateToken(tokenData)
	assert.NotNil(t, token)

	//create http request
	httpRequest := &http.Request{
		Header: http.Header{},
	}
	httpRequest.Header.Add("Authorization", "Bearer "+token.AccessToken)
	httpRequest.RemoteAddr = "localhost2"

	//verify token
	orgPreventTokenHiJacking := conf.Config.OAuth2.PreventTokenHiJacking
	conf.Config.OAuth2.PreventTokenHiJacking = false
	metadata := jwt.ExtractTokenMetadata(requestSessionTest, httpRequest)
	conf.Config.OAuth2.PreventTokenHiJacking = orgPreventTokenHiJacking
	assert.NotNil(t, metadata)
	assert.Equal(t, &tokenData, metadata)
}

func TestTokenData_Hashed(t1 *testing.T) {
	tests := []struct {
		name   string
		fields jwt.TokenData
		want   jwt.TokenData
	}{
		{
			name: "hash token data", fields: jwt.TokenData{
				Username:           "john",
				Groups:             "some group",
				GroupType:          "type",
				Email:              "lala@lala.lala",
				RemoteAddress:      "192.168.42.1",
				IsInternalEmployee: false,
			}, want: jwt.TokenData{
				Username:           "afd22b4dcd3cfe92bac7aafa3e783a2eee3e5a52a9a36c7d1c7e37d15c76eb65",
				Groups:             "some group",
				GroupType:          "type",
				Email:              "dd4525cd67a7f34bed7fe246c0a657d6b9420714230e9edba9ccaa6ab7023aca",
				RemoteAddress:      "684bddbe6c365e91929bdf7f2a82e9101c98cc540fa580ce6ca8888ade68f1e1",
				IsInternalEmployee: false,
			},
		},
		{
			name: "do not hash empty fields", fields: jwt.TokenData{
				Username:           "john",
				Groups:             "some group",
				GroupType:          "type",
				Email:              "",
				RemoteAddress:      "",
				IsInternalEmployee: false,
			}, want: jwt.TokenData{
				Username:           "afd22b4dcd3cfe92bac7aafa3e783a2eee3e5a52a9a36c7d1c7e37d15c76eb65",
				Groups:             "some group",
				GroupType:          "type",
				Email:              "",
				RemoteAddress:      "",
				IsInternalEmployee: false,
			},
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := jwt.TokenData{
				Username:           tt.fields.Username,
				Groups:             tt.fields.Groups,
				GroupType:          tt.fields.GroupType,
				Email:              tt.fields.Email,
				RemoteAddress:      tt.fields.RemoteAddress,
				IsInternalEmployee: tt.fields.IsInternalEmployee,
			}
			assert.Equalf(t1, tt.want, t.Hashed(requestSessionTest), "Hashed()")
		})
	}
}

func TestCreateUserData(t *testing.T) {
	// Mock the request session, user, and http request
	requestSession := &logy.RequestSession{}
	testUser := &user.User{
		TermsOfUse: false,
		Active:     true,
	}
	username := "testuser"
	email := "test@example.com"
	groups := []string{"group1", "group2"}
	groupType := "testGroupType"
	isInternalEmployee := true
	request := &http.Request{
		RemoteAddr: "127.0.0.1:8080",
		Header: http.Header{
			"Cache-Control": {"max-age=604800"},
		},
	}

	// Set up the expected token data
	expectedTokenData := jwt.TokenData{
		Username:           username,
		TermsOfUse:         testUser.TermsOfUse,
		IsEnabled:          testUser.Active,
		Email:              email,
		Groups:             "group1;group2", // Assuming GROUPS_TOKEN is ","
		GroupType:          groupType,
		IsInternalEmployee: isInternalEmployee,
		RemoteAddress:      "127.0.0.1", // Assuming TrimPortFromRemoteAddress works correctly
	}

	//conf.Config.OAuth2.DebugLog = true

	// Call the function under test
	tokenData := jwt.CreateUserData(requestSession, testUser, username, email, groups, groupType, isInternalEmployee, request)

	// Assert that the returned token data matches the expected token data
	assert.Equal(t, expectedTokenData, tokenData, "The token data should match the expected values")
}
