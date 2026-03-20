// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package internalToken

import (
	"time"

	"github.com/eclipse-disuko/disuko/domain/internalToken"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
)

const collectionName = "basicauths"

type IRepo interface {
	base.IBaseRepositoryWithHardDelete[*internalToken.InternalToken]
	CheckAuth(requestSession *logy.RequestSession, key, token string) *internalToken.InternalToken
	Add(requestSession *logy.RequestSession, name, description string, expiry time.Time, caps []internalToken.Capability) (*internalToken.InternalToken, string)
	Revoke(requestSession *logy.RequestSession, key string)
	Renew(requestSession *logy.RequestSession, key string) (*internalToken.InternalToken, string)
}
