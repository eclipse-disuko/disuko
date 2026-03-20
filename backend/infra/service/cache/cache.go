// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package cache

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/helper/s3Helper"
	"github.com/eclipse-disuko/disuko/infra/service/locks"
	"github.com/eclipse-disuko/disuko/logy"
)

const CachePath = "cache/%s.json"

type CacheService struct {
	RequestSession *logy.RequestSession
	WithLock       bool
	LockService    *locks.Service
}

func (s *CacheService) GetCacheEntry(key string, v any) bool {
	filePath := fmt.Sprintf(CachePath, key)
	if !s3Helper.ExistFile(s.RequestSession, filePath) {
		return false
	}
	if s.WithLock {
		l, acquired := s.LockService.Acquire(locks.Options{
			Key:      key,
			Blocking: true,
			Timeout:  time.Second * 10,
		})
		if !acquired {
			exception.ThrowExceptionServerMessage(message.GetI18N(message.ResourceInUse), "")
		}
		defer s.LockService.Release(l)
	}
	text := s3Helper.ReadTextFile(s.RequestSession, filePath, "")
	err := json.Unmarshal([]byte(*text), v)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.UnmarshallingCache))
	return true
}

func (s *CacheService) SaveCacheEntry(key string, data interface{}) {
	if s.WithLock {
		l, acquired := s.LockService.Acquire(locks.Options{
			Key:      key,
			Blocking: true,
			Timeout:  time.Second * 10,
		})
		if !acquired {
			exception.ThrowExceptionServerMessage(message.GetI18N(message.ResourceInUse), "")
		}
		defer s.LockService.Release(l)
	}
	filePath := fmt.Sprintf(CachePath, key)
	if s3Helper.ExistFile(s.RequestSession, filePath) {
		s3Helper.DeleteFile(s.RequestSession, filePath)
	}
	s3Helper.SaveObjectToFile(
		s.RequestSession,
		filePath,
		data,
		nil,
	)
}
