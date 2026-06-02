// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
// SPDX-License-Identifier: Apache-2.0

package i18n

import (
	"time"

	"github.com/eclipse-disuko/disuko/domain/i18n"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/infra/repository/database"
	"github.com/eclipse-disuko/disuko/logy"
)

type i18nRepositoryStruct struct {
	base.BaseRepositoryWithSoftDelete[*i18n.I18nLocale]
}

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
			[][]string{{"Scope"}}),
	}
}

func (repository *i18nRepositoryStruct) FindByLocaleCode(requestSession *logy.RequestSession, localeCode string, deleted bool) *i18n.I18nLocale {
	return repository.FindByKey(requestSession, localeCode, deleted)
}

func (repository *i18nRepositoryStruct) FindDefault(requestSession *logy.RequestSession) *i18n.I18nLocale {
	qc := database.New().SetMatcher(database.AttributeMatcher("IsDefault", database.EQ, true))
	results := repository.Query(requestSession, qc)
	if len(results) == 0 {
		return nil
	}
	return results[0]
}

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
		entry.Updated = time.Now()
	}
	entry.UpdatedBy = updatedBy
	locale.SetEntry(entry)

	if alreadyExists {
		repository.Update(requestSession, locale)
	} else {
		repository.Save(requestSession, locale)
	}
}

func (repository *i18nRepositoryStruct) UpsertLocaleMetadata(requestSession *logy.RequestSession, localeCode string, displayName string, nativeName string, isDefault bool, scope string) {
	locale := repository.FindByKey(requestSession, localeCode, false)
	alreadyExists := locale != nil

	if locale == nil {
		locale = i18n.NewI18nLocale(localeCode)
		repository.copyKeysFromDefaultLocale(requestSession, locale, localeCode)
	}

	if isDefault {
		if prev := repository.FindDefault(requestSession); prev != nil && prev.Key != localeCode {
			prev.IsDefault = false
			repository.Update(requestSession, prev)
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

func (repository *i18nRepositoryStruct) copyKeysFromDefaultLocale(requestSession *logy.RequestSession, newLocale *i18n.I18nLocale, newLocaleCode string) {
	defaultLocale := repository.FindDefault(requestSession)
	if (defaultLocale == nil || defaultLocale.Key == newLocaleCode) && newLocaleCode != "en" {
		defaultLocale = repository.FindByLocaleCode(requestSession, "en", false)
	}
	if defaultLocale == nil || defaultLocale.Key == newLocaleCode {
		return
	}
	for key, entry := range defaultLocale.Entries {
		newEntry := i18n.NewI18nEntry(key, entry.Value, entry.Description)
		newEntry.CreatedBy = entry.CreatedBy
		newLocale.SetEntry(newEntry)
	}
}

func (repository *i18nRepositoryStruct) FindByLocaleCodeOrDefault(requestSession *logy.RequestSession, localeCode string) (*i18n.I18nLocale, bool) {
	if locale := repository.FindByLocaleCode(requestSession, localeCode, false); locale != nil {
		return locale, false
	}
	return repository.FindDefault(requestSession), true
}

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

func (repository *i18nRepositoryStruct) GetTranslationWithFallback(requestSession *logy.RequestSession, localeCode string, key string) (string, string, bool) {
	if value, ok := repository.GetTranslation(requestSession, localeCode, key); ok {
		return value, localeCode, true
	}
	defaultLocale := repository.FindDefault(requestSession)
	if defaultLocale != nil && defaultLocale.Key != localeCode {
		if value, ok := repository.GetTranslation(requestSession, defaultLocale.Key, key); ok {
			return value, defaultLocale.Key, true
		}
	}
	return "", "", false
}

func (repository *i18nRepositoryStruct) DeleteTranslation(requestSession *logy.RequestSession, localeCode string, key string) {
	locale := repository.FindByKey(requestSession, localeCode, false)
	if locale == nil {
		return
	}

	locale.RemoveEntry(key)
	repository.Update(requestSession, locale)
}

func (repository *i18nRepositoryStruct) GetLocaleCount(requestSession *logy.RequestSession) int {
	return len(repository.FindAllKeys(requestSession))
}

// Returns false if the locale is set as default.
func (repository *i18nRepositoryStruct) DeleteLocale(requestSession *logy.RequestSession, localeCode string) bool {
	locale := repository.FindByKey(requestSession, localeCode, false)
	if locale == nil {
		return true
	}
	if locale.IsDefault {
		return false
	}
	repository.Delete(requestSession, localeCode)
	return true
}
