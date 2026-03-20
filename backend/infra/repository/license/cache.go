// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package license

import (
	"errors"
	"sync"

	"github.com/eclipse-disuko/disuko/domain/license"
	"github.com/eclipse-disuko/disuko/logy"
	"golang.org/x/net/context"
)

type Cache struct {
	mu             sync.Mutex
	cache          map[string]*license.License
	repo           ILicensesRepository
	requestSession *logy.RequestSession
	ctx            context.Context
}

func NewLicenseCache(requestSession *logy.RequestSession, repo ILicensesRepository) *Cache {
	return &Cache{
		mu:             sync.Mutex{},
		cache:          make(map[string]*license.License),
		repo:           repo,
		requestSession: requestSession,
		ctx:            context.Background(),
	}
}

func (c *Cache) Get(licenseId string) (*license.License, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if lic, ok := c.cache[licenseId]; ok {
		return lic, nil
	}

	lic := c.repo.FindById(c.requestSession, licenseId)
	if lic == nil {
		return nil, errors.New("license not found")
	}
	c.cache[licenseId] = lic
	return lic, nil
}
