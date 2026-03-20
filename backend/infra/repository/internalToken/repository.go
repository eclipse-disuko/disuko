// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package internalToken

import (
	"crypto/rand"
	"math/big"
	"time"

	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/domain/internalToken"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Repo struct {
	base.BaseRepositoryWithHardDelete[*internalToken.InternalToken]
}

func NewRepo(requestSession *logy.RequestSession) IRepo {
	innerRepo := &Repo{
		BaseRepositoryWithHardDelete: base.CreateRepositoryWithHardDelete(
			requestSession,
			collectionName,
			func() *internalToken.InternalToken {
				return &internalToken.InternalToken{}
			},
			nil,
			"",
			nil,
			nil),
	}
	return innerRepo
}

func (r *Repo) CheckAuth(requestSession *logy.RequestSession, key, token string) *internalToken.InternalToken {
	t := r.FindByKey(requestSession, key, false)
	if t == nil {
		return nil
	}
	if time.Now().After(t.Expiry) || t.Revoked {
		return nil
	}
	spToken := conf.Config.BasicAuth.Pepper + token + t.Salt
	if bcrypt.CompareHashAndPassword([]byte(t.Token), []byte(spToken)) != nil {
		return nil
	}

	return t
}

func (r *Repo) Add(requestSession *logy.RequestSession, name, description string, expiry time.Time, caps []internalToken.Capability) (*internalToken.InternalToken, string) {
	salt := randomString(5)
	t := uuid.NewString()
	spPassword := conf.Config.BasicAuth.Pepper + t + salt
	crypted, err := bcrypt.GenerateFromPassword([]byte(spPassword), bcrypt.DefaultCost)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorEncrypting))

	n := internalToken.InternalToken{
		RootEntity:   domain.NewRootEntity(),
		Name:         name,
		Token:        string(crypted),
		Salt:         salt,
		Description:  description,
		Expiry:       expiry,
		Capabilities: caps,
	}
	r.Save(requestSession, &n)

	return &n, t
}

func (r *Repo) Revoke(requestSession *logy.RequestSession, key string) {
	t := r.FindByKey(requestSession, key, false)
	if t == nil {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorDbExist, key), "")
	}
	if time.Now().After(t.Expiry) {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.TokenExpired), "")
	}
	if t.Revoked {
		exception.ThrowExceptionBadRequestResponse()
	}
	t.Revoked = true
	r.Update(requestSession, t)
}

func (r *Repo) Renew(requestSession *logy.RequestSession, key string) (*internalToken.InternalToken, string) {
	t := r.FindByKey(requestSession, key, false)
	if t == nil {
		exception.ThrowExceptionClientMessage(message.GetI18N(message.ErrorDbExist, key), "")
	}
	if time.Now().After(t.Expiry) {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.TokenExpired), "")
	}
	if t.Revoked {
		exception.ThrowExceptionBadRequestResponse()
	}

	salt := randomString(5)
	sec := uuid.NewString()
	spPassword := conf.Config.BasicAuth.Pepper + sec + salt
	crypted, err := bcrypt.GenerateFromPassword([]byte(spPassword), bcrypt.DefaultCost)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorEncrypting))

	t.Salt = salt
	t.Token = string(crypted)

	r.Update(requestSession, t)

	return t, sec
}

func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	seededRand := big.NewInt(int64(len(letters)))
	s := make([]rune, n)
	for i := range s {
		number, err := rand.Int(rand.Reader, seededRand)
		if err != nil {
			return ""
		}
		s[i] = letters[number.Int64()]
	}
	return string(s)
}
