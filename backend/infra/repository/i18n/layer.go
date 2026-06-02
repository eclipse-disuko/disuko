// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
// SPDX-License-Identifier: Apache-2.0

package i18n

import (
	"github.com/eclipse-disuko/disuko/domain/i18n"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
)

const I18nLocaleCollectionName = "i18nLocale"

type II18nRepository interface {
	base.IBaseRepositoryWithSoftDelete[*i18n.I18nLocale]

	FindByLocaleCode(requestSession *logy.RequestSession, localeCode string, deleted bool) *i18n.I18nLocale
	FindDefault(requestSession *logy.RequestSession) *i18n.I18nLocale
	FindByLocaleCodeOrDefault(requestSession *logy.RequestSession, localeCode string) (locale *i18n.I18nLocale, fallbackUsed bool)
	SetTranslation(requestSession *logy.RequestSession, localeCode string, key string, value string, description string, updatedBy string)
	UpsertLocaleMetadata(requestSession *logy.RequestSession, localeCode string, displayName string, nativeName string, isDefault bool, scope string)
	// Returns false if locale is set as default.
	DeleteLocale(requestSession *logy.RequestSession, localeCode string) bool
	GetTranslation(requestSession *logy.RequestSession, localeCode string, key string) (string, bool)
	GetTranslationWithFallback(requestSession *logy.RequestSession, localeCode string, key string) (value string, actualLocale string, found bool)
	DeleteTranslation(requestSession *logy.RequestSession, localeCode string, key string)
	GetLocaleCount(requestSession *logy.RequestSession) int
}
