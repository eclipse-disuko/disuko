// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package jwt

import (
	"fmt"
	"net/http"
	"net/netip"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"mercedes-benz.ghe.com/foss/disuko/conf"
	"mercedes-benz.ghe.com/foss/disuko/domain/user"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/hash"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

const JK_USERNAME string = "username"
const JK_GROUPS string = "groups"
const JK_GROUP_TYPE string = "grouptype"
const JK_EMAIL string = "email"
const JK_EXPIRES string = "exp"
const JK_AUTHORIZED string = "authorized"
const JK_REMOTE_ADDRESS string = "remote_address"
const JK_IS_INTERNAL string = "isInternal"
const JK_IS_ENABLED string = "isEnabled"

const GROUP_TYPE_DAIMLER string = "0"
const GROUP_TYPE_TRUCKS string = "1"
const GROUP_TYPE_SHARED string = "S"

const GROUPS_TOKEN string = ";"

type TokenData struct {
	Username           string `json:"username"`
	Groups             string `json:"groups"`
	GroupType          string `json:"grouptype"`
	Email              string `json:"email"`
	RemoteAddress      string `json:"remoteaddress"`
	IsInternalEmployee bool   `json:"isInternalEmployee"`
	TermsOfUse         bool   `json:"hasTermsOfUse"`
	IsEnabled          bool   `json:"isEnabled"`
}

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AtExpires    time.Time
	RtExpires    time.Time
}

func (t TokenData) Hashed(requestSession *logy.RequestSession) TokenData {
	return TokenData{
		Username:           hash.Hash(requestSession, t.Username),
		Groups:             t.Groups,
		GroupType:          t.GroupType,
		Email:              hash.Hash(requestSession, t.Email),
		RemoteAddress:      hash.Hash(requestSession, t.RemoteAddress),
		IsInternalEmployee: false,
	}
}

func CreateTechnicalUserdata(request *http.Request) TokenData {
	remoteAddress := TrimPortFromRemoteAddress(request.RemoteAddr)

	userData := TokenData{
		Username:           "Tech user",
		Email:              "noanswer@disclosureportal.de",
		Groups:             "FOSSDP.domain_admin",
		GroupType:          GROUP_TYPE_DAIMLER,
		IsInternalEmployee: true,
		RemoteAddress:      remoteAddress,
		IsEnabled:          !conf.IsProdEnv(),
		TermsOfUse:         true,
	}
	return userData
}
func CreateUserData(requestSession *logy.RequestSession, user *user.User, username string, email string, groups []string, groupType string, isInternalEmployee bool, request *http.Request) TokenData {
	if conf.Config.OAuth2.DebugLog {
		logy.Infof(requestSession, "")
		logy.Infof(requestSession, "################################################")
		logy.Infof(requestSession, "DEBUG: print HTTP Header for token request - start")
		allowList := []string{
			"Sec-Fetch-Dest",
			"Referer",
			"Priority",
			"Sec-Ch-Ua-Mobile",
			"Sec-Ch-Ua-Platform",
			"Cache-Control",
			"Sec-Fetch-Site",
			"Sec-Fetch-Mode",
			"Sec-Ch-Ua",
			"Sec-Fetch-User",
			"Upgrade-Insecure-Requests",
			"User-Agent",
			"Accept",
			"Accept-Encoding",
			"Accept-Language",
			"Cookie",
		}
		for _, element := range allowList {
			value, ok := request.Header[element]
			if !ok {
				continue
			}
			logy.Infof(requestSession, element, "=", value)
		}
		logy.Infof(requestSession, "DEBUG: print HTTP Header for token request - end")
		logy.Infof(requestSession, "################################################")
		logy.Infof(requestSession, "")
	}

	//remove port from remote address, because this must not be the same
	remoteAddress := TrimPortFromRemoteAddress(request.RemoteAddr)

	userData := TokenData{
		Username:           username,
		TermsOfUse:         user.TermsOfUse,
		IsEnabled:          user.Active,
		Email:              email,
		Groups:             strings.Join(groups, GROUPS_TOKEN),
		GroupType:          groupType,
		IsInternalEmployee: isInternalEmployee,
		RemoteAddress:      remoteAddress,
	}
	return userData
}

func TrimPortFromRemoteAddress(remoteAddress string) string {
	splitted := strings.Split(remoteAddress, ":")
	if len(splitted) <= 2 {
		return splitted[0]
	}
	if remoteAddress[0] == '[' {
		ap, err := netip.ParseAddrPort(remoteAddress)
		if err != nil {
			return remoteAddress
		}
		return ap.Addr().String()
	}
	return remoteAddress
}

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestSession := logy.GetRequestSession(r)
		ExtractAndVerifyAccessToken(requestSession, r)
		next.ServeHTTP(w, r)
	})
}

func CreateToken(userData TokenData) *TokenDetails {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * time.Duration(conf.Config.Auth.AccessTokenExpiresInMinutes))

	td.RtExpires = time.Now().Add(time.Hour * time.Duration(conf.Config.Auth.RefreshTokenExpiresInHours))

	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims[JK_AUTHORIZED] = true
	atClaims[JK_USERNAME] = userData.Username
	atClaims[JK_EMAIL] = userData.Email
	atClaims[JK_GROUPS] = userData.Groups
	atClaims[JK_REMOTE_ADDRESS] = userData.RemoteAddress
	atClaims[JK_EXPIRES] = td.AtExpires.Unix()
	atClaims[JK_GROUP_TYPE] = userData.GroupType
	atClaims[JK_IS_INTERNAL] = userData.IsInternalEmployee
	atClaims[JK_IS_ENABLED] = userData.IsEnabled
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	var err error
	td.AccessToken, err = at.SignedString([]byte(conf.Config.Auth.AccessSecret))
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorTokenCreate))

	//Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims[JK_USERNAME] = userData.Username
	rtClaims[JK_EMAIL] = userData.Email
	rtClaims[JK_GROUPS] = userData.Groups
	rtClaims[JK_REMOTE_ADDRESS] = userData.RemoteAddress
	rtClaims[JK_IS_ENABLED] = userData.IsEnabled
	rtClaims[JK_EXPIRES] = td.RtExpires.Unix()
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(conf.Config.Auth.RefreshSecret))
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorTokenCreate))

	return td
}

func CreateAccessTokenStrFromClaims(sourceClaims jwt.MapClaims) (token string, expiry time.Time) {
	expires := time.Now().Add(time.Minute * time.Duration(conf.Config.Auth.AccessTokenExpiresInMinutes))

	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims[JK_AUTHORIZED] = true
	atClaims[JK_USERNAME] = ExtractClaimStr(sourceClaims, JK_USERNAME)
	atClaims[JK_EMAIL] = ExtractClaimStr(sourceClaims, JK_EMAIL)
	atClaims[JK_GROUPS] = ExtractClaimStr(sourceClaims, JK_GROUPS)
	atClaims[JK_REMOTE_ADDRESS] = ExtractClaimStr(sourceClaims, JK_REMOTE_ADDRESS)
	atClaims[JK_EXPIRES] = expires.Unix()
	atClaims[JK_GROUP_TYPE] = ExtractClaimStr(sourceClaims, JK_GROUP_TYPE)
	atClaims[JK_IS_INTERNAL] = ExtractClaimBool(sourceClaims, JK_IS_INTERNAL)
	atClaims[JK_IS_ENABLED] = ExtractClaimBool(sourceClaims, JK_IS_ENABLED)
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	tokenStr, err := at.SignedString([]byte(conf.Config.Auth.AccessSecret))
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorTokenCreate))

	return tokenStr, expires
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")

	cookie, err := r.Cookie("oauth.a")
	if err == nil {
		bearToken = "Bearer " + cookie.Value
	}

	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func ExtractTokenMetadata(requestSession *logy.RequestSession, r *http.Request) *TokenData {
	token := ExtractAndVerifyAccessToken(requestSession, r)
	claims := ExtractClaims(token)

	username := ExtractClaimStr(claims, JK_USERNAME)
	email := ExtractClaimStr(claims, JK_EMAIL)
	groups := ExtractClaimStr(claims, JK_GROUPS)
	groupType := ExtractClaimStr(claims, JK_GROUP_TYPE)
	remoteAddress := ExtractClaimStr(claims, JK_REMOTE_ADDRESS)
	isInternal := ExtractClaimBool(claims, JK_IS_INTERNAL)
	isEnabled := ExtractClaimBool(claims, JK_IS_ENABLED)

	return &TokenData{
		Username:           username,
		Email:              email,
		Groups:             groups,
		GroupType:          groupType,
		RemoteAddress:      remoteAddress,
		IsInternalEmployee: isInternal,
		IsEnabled:          isEnabled,
	}
}

func ExtractTokenMetadataExternal(requestSession *logy.RequestSession, raw string) *TokenData {
	token := ParseToken(raw, conf.Config.Auth.AccessSecret)
	claims := ExtractClaims(token)

	username := ExtractClaimStr(claims, JK_USERNAME)
	email := ExtractClaimStr(claims, JK_EMAIL)
	groups := ExtractClaimStr(claims, JK_GROUPS)
	groupType := ExtractClaimStr(claims, JK_GROUP_TYPE)
	remoteAddress := ExtractClaimStr(claims, JK_REMOTE_ADDRESS)
	isInternal := ExtractClaimBool(claims, JK_IS_INTERNAL)
	isEnabled := ExtractClaimBool(claims, JK_IS_ENABLED)

	return &TokenData{
		Username:           username,
		Email:              email,
		Groups:             groups,
		GroupType:          groupType,
		RemoteAddress:      remoteAddress,
		IsInternalEmployee: isInternal,
		IsEnabled:          isEnabled,
	}
}

func ExtractClaimStr(claims jwt.MapClaims, key string) string {
	value, ok := claims[key].(string)
	if !ok {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.ErrorAAR, "Invalid token"), "missing "+key+" in token (string value)")
	}
	return value
}

func ExtractClaimBool(claims jwt.MapClaims, key string) bool {
	value, ok := claims[key].(bool)
	if !ok {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.ErrorAAR, "Invalid token"), "missing "+key+" in token (bool value)")
	}
	return value
}

func ExtractClaimInt64(claims jwt.MapClaims, key string) int64 {
	value, ok := claims[key].(float64)
	if !ok {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.ErrorAAR, "Invalid token"), "missing "+key+" in token (int64 value)")
	}
	return int64(value)
}

func ExtractAndVerifyAccessToken(requestSession *logy.RequestSession, r *http.Request) *jwt.Token {
	tokenString := ExtractToken(r)
	return ParseAndVerifyToken(requestSession, r, tokenString, conf.Config.Auth.AccessSecret)
}

func ParseAndVerifyToken(requestSession *logy.RequestSession, r *http.Request, tokenString string, secret string) *jwt.Token {
	token := ParseToken(tokenString, secret)
	VerifyRemoteAddress(requestSession, r, token)
	return token
}

func VerifyRemoteAddress(requestSession *logy.RequestSession, r *http.Request, token *jwt.Token) {
	claims := ExtractClaims(token)

	//flag is needed, if a http request from outside should use the bearer
	if conf.Config.OAuth2.PreventTokenHiJacking {
		remoteAddress, ok := claims[JK_REMOTE_ADDRESS].(string)
		if !ok || len(strings.TrimSpace(remoteAddress)) < 1 {
			exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.ErrorAAR, "Invalid token"), "missing or empty remote address in token")
		}
		if conf.Config.OAuth2.DebugLog {
			logy.Infof(requestSession, "remote address in token: %s", remoteAddress)
		}
		//remove port from remote address, because this must not be the same
		currentRemoteAddress := TrimPortFromRemoteAddress(r.RemoteAddr)
		if currentRemoteAddress != remoteAddress {
			exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.ErrorAAR, "Invalid token"), "remote address in token does not match")
		}
	}
}

func ExtractClaims(token *jwt.Token) jwt.MapClaims {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.ErrorAAR, "Invalid token"), "Invalid claims in token")
	}
	return claims
}

func ParseToken(tokenString string, secret string) *jwt.Token {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.ErrorAAR, "Invalid token"), err.Error())
	}

	if !token.Valid {
		exception.ThrowExceptionSendDeniedResponseRaw(message.GetI18N(message.ErrorAAR, "Invalid token"), "Invalid token")
	}
	return token
}
