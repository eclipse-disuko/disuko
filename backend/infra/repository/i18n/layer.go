// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
// SPDX-License-Identifier: Apache-2.0

package i18n

import (
	"github.com/eclipse-disuko/disuko/domain/i18n"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
)

const I18nLocaleCollectionName = "i18nLocale"

// II18nRepository defines the interface for i18n persistence operations.
type II18nRepository interface {
	// Base CRUD operations with soft delete support
	base.IBaseRepositoryWithSoftDelete[*i18n.I18nLocale]

	// FindByLocaleCode retrieves a locale by its language code (en, de, etc.)
	FindByLocaleCode(requestSession *logy.RequestSession, localeCode string, deleted bool) *i18n.I18nLocale

	// SetTranslation updates a translation entry for a locale.
	SetTranslation(requestSession *logy.RequestSession, localeCode string, key string, value string, description string, updatedBy string)

	// SetLocaleMetadata updates locale metadata and creates locale if missing.
	SetLocaleMetadata(requestSession *logy.RequestSession, localeCode string, displayName string, nativeName string, isDefault bool, scope string)

	// DeleteLocale permanently removes a locale (hard delete). Returns false if locale is default.
	DeleteLocale(requestSession *logy.RequestSession, localeCode string) bool

	// GetTranslation retrieves a translation value for a locale and key.
	GetTranslation(requestSession *logy.RequestSession, localeCode string, key string) (string, bool)

	// FindAllEntries returns all translation entries for a locale.
	FindAllEntries(requestSession *logy.RequestSession, localeCode string) map[string]*i18n.I18nEntry

	// DeleteTranslation removes a translation entry from a locale.
	DeleteTranslation(requestSession *logy.RequestSession, localeCode string, key string)

	// GetLocaleCount returns the total number of active locales.
	GetLocaleCount(requestSession *logy.RequestSession) int

	// GetEntryCountForLocale returns the number of translation entries for a locale.
	GetEntryCountForLocale(requestSession *logy.RequestSession, localeCode string) int
}
