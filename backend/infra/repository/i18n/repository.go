// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
// SPDX-License-Identifier: Apache-2.0

package i18n

import (
	"time"

	"github.com/eclipse-disuko/disuko/domain/i18n"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/logy"
)

type i18nRepositoryStruct struct {
	base.BaseRepositoryWithSoftDelete[*i18n.I18nLocale]
}

// NewI18nRepository creates a new instance of the i18n repository.
func NewI18nRepository(requestSession *logy.RequestSession) II18nRepository {
	return &i18nRepositoryStruct{
		BaseRepositoryWithSoftDelete: base.CreateRepositoryWithSoftDelete[*i18n.I18nLocale](
			requestSession,
			I18nLocaleCollectionName,
			func() *i18n.I18nLocale {
				return &i18n.I18nLocale{}
			},
			nil,
			nil,
			nil),
	}
}

// FindByLocaleCode retrieves a locale by its language code.
func (repository *i18nRepositoryStruct) FindByLocaleCode(requestSession *logy.RequestSession, localeCode string, deleted bool) *i18n.I18nLocale {
	return repository.FindByKey(requestSession, localeCode, deleted)
}

// SetTranslation updates a translation entry for a locale, creating the locale if it does not exist.
func (repository *i18nRepositoryStruct) SetTranslation(requestSession *logy.RequestSession, localeCode string, key string, value string, description string, updatedBy string) {
	locale := repository.FindByKey(requestSession, localeCode, false)
	alreadyExists := locale != nil

	if locale == nil {
		locale = i18n.NewI18nLocale(localeCode)
	}

	entry := locale.GetEntry(key)
	if entry == nil {
		entry = i18n.NewI18nEntry(key, value, description)
		entry.CreatedBy = updatedBy
	} else {
		entry.Value = value
		entry.Description = description
		entry.UpdatedAt = time.Now()
	}
	entry.UpdatedBy = updatedBy
	locale.SetEntry(entry)

	if alreadyExists {
		repository.Update(requestSession, locale)
	} else {
		repository.Save(requestSession, locale)
	}
}

// SetLocaleMetadata updates locale metadata and creates locale if missing.
// When creating a new locale, it copies all translation keys from the default locale.
func (repository *i18nRepositoryStruct) SetLocaleMetadata(requestSession *logy.RequestSession, localeCode string, displayName string, nativeName string, isDefault bool, scope string) {
	locale := repository.FindByKey(requestSession, localeCode, false)
	alreadyExists := locale != nil

	if locale == nil {
		locale = i18n.NewI18nLocale(localeCode)
		// Before saving a new locale, copy keys from default locale
		copyKeysFromDefaultLocaleToObject(repository, requestSession, locale, localeCode)
	}

	if isDefault {
		allLocales := repository.FindAll(requestSession, false)
		for _, current := range allLocales {
			if current != nil && current.LocaleCode != localeCode && current.IsDefault {
				current.IsDefault = false
				repository.Update(requestSession, current)
			}
		}
	}

	locale.DisplayName = displayName
	locale.NativeName = nativeName
	locale.IsDefault = isDefault
	locale.Scope = scope

	if alreadyExists {
		repository.Update(requestSession, locale)
	} else {
		repository.Save(requestSession, locale)
	}
}

// copyKeysFromDefaultLocaleToObject copies all translation keys from default locale to a locale object before it's saved.
func copyKeysFromDefaultLocaleToObject(repository *i18nRepositoryStruct, requestSession *logy.RequestSession, newLocale *i18n.I18nLocale, newLocaleCode string) {
	// Find the default locale
	var defaultLocale *i18n.I18nLocale
	allLocales := repository.FindAll(requestSession, false)
	for _, locale := range allLocales {
		if locale != nil && locale.IsDefault && locale.LocaleCode != newLocaleCode {
			defaultLocale = locale
			break
		}
	}

	// If no default locale exists, try to find "en" specifically as base
	if defaultLocale == nil && newLocaleCode != "en" {
		defaultLocale = repository.FindByLocaleCode(requestSession, "en", false)
	}

	// If still no locale found but there's at least one other locale, use the first one
	if defaultLocale == nil && len(allLocales) > 0 {
		for _, locale := range allLocales {
			if locale != nil && locale.LocaleCode != newLocaleCode {
				defaultLocale = locale
				break
			}
		}
	}

	// Copy all entries from the default/base locale to the new locale object
	if defaultLocale != nil {
		if defaultLocale.Entries == nil || len(defaultLocale.Entries) == 0 {
			// No entries to copy, exit silently
			return
		}

		for key, entry := range defaultLocale.Entries {
			if entry != nil {
				// Create new entry with the same key and value
				newEntry := i18n.NewI18nEntry(key, entry.Value, entry.Description)
				newEntry.CreatedBy = entry.CreatedBy
				newLocale.SetEntry(newEntry)
			}
		}
	}
}

// GetTranslation retrieves a translation value for a locale and key.
func (repository *i18nRepositoryStruct) GetTranslation(requestSession *logy.RequestSession, localeCode string, key string) (string, bool) {
	locale := repository.FindByKey(requestSession, localeCode, false)
	if locale == nil {
		return "", false
	}

	entry := locale.GetEntry(key)
	if entry == nil {
		return "", false
	}

	return entry.Value, true
}

// FindAllEntries returns all translation entries for a locale.
func (repository *i18nRepositoryStruct) FindAllEntries(requestSession *logy.RequestSession, localeCode string) map[string]*i18n.I18nEntry {
	locale := repository.FindByKey(requestSession, localeCode, false)
	if locale == nil {
		return make(map[string]*i18n.I18nEntry)
	}

	if locale.Entries == nil {
		return make(map[string]*i18n.I18nEntry)
	}

	return locale.Entries
}

// DeleteTranslation removes a translation entry from a locale.
func (repository *i18nRepositoryStruct) DeleteTranslation(requestSession *logy.RequestSession, localeCode string, key string) {
	locale := repository.FindByKey(requestSession, localeCode, false)
	if locale == nil {
		return
	}

	locale.RemoveEntry(key)
	repository.Update(requestSession, locale)
}

// GetLocaleCount returns the total number of active locales.
func (repository *i18nRepositoryStruct) GetLocaleCount(requestSession *logy.RequestSession) int {
	locales := repository.FindAll(requestSession, false)
	if locales == nil {
		return 0
	}
	return len(locales)
}

// GetEntryCountForLocale returns the number of translation entries for a locale.
func (repository *i18nRepositoryStruct) GetEntryCountForLocale(requestSession *logy.RequestSession, localeCode string) int {
	locale := repository.FindByKey(requestSession, localeCode, false)
	if locale == nil {
		return 0
	}
	return locale.GetEntryCount()
}

// DeleteLocale permanently removes a locale. Returns false if the locale is the default.
func (repository *i18nRepositoryStruct) DeleteLocale(requestSession *logy.RequestSession, localeCode string) bool {
	locale := repository.FindByKey(requestSession, localeCode, false)
	if locale == nil {
		return true
	}
	if locale.IsDefault {
		return false
	}
	repository.DeleteHard(requestSession, localeCode)
	return true
}
