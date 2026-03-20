// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package migration

import (
	"github.com/eclipse-disuko/disuko/domain"
	"github.com/eclipse-disuko/disuko/helper/hash"
	"github.com/eclipse-disuko/disuko/logy"
)

type Migration struct {
	domain.RootEntity `bson:"inline"`
	Name              string
}

func New(session *logy.RequestSession, name string) *Migration {
	return &Migration{
		RootEntity: domain.NewRootEntityWithKey(hash.Hash(session, name)),
		Name:       name,
	}
}
