// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
// SPDX-License-Identifier: Apache-2.0

package i18n

import (
	"github.com/eclipse-disuko/disuko/domain"
)

type I18nLocale struct {
	domain.RootEntity `bson:",inline"`
	domain.SoftDelete `bson:",inline"`

	DisplayName string
	NativeName  string
	IsDefault   bool
	Entries     map[string]*I18nEntry
	Scope       string
}

type I18nEntry struct {
	domain.ChildEntity `bson:",inline"`
	Value       string
	Description string
	CreatedBy   string
	UpdatedBy   string
}

func NewI18nLocale(localeCode string) *I18nLocale {
	return &I18nLocale{
		RootEntity: domain.NewRootEntityWithKey(localeCode),
		Entries:    make(map[string]*I18nEntry),
	}
}

func NewI18nEntry(key, value, description string) *I18nEntry {
	return &I18nEntry{
		ChildEntity: domain.SetChildEntity(key),
		Value:       value,
		Description: description,
	}
}

func (locale *I18nLocale) SetEntry(entry *I18nEntry) {
	if locale.Entries == nil {
		locale.Entries = make(map[string]*I18nEntry)
	}
	locale.Entries[entry.Key] = entry
}

func (locale *I18nLocale) GetEntry(key string) *I18nEntry {
	if locale.Entries == nil {
		return nil
	}
	return locale.Entries[key]
}

func (locale *I18nLocale) RemoveEntry(key string) {
	if locale.Entries == nil {
		return
	}
	delete(locale.Entries, key)
}

func (locale *I18nLocale) ToListDTO() I18nLocaleListResponseDto {
	return I18nLocaleListResponseDto{
		LocaleCode:  locale.Key,
		DisplayName: locale.DisplayName,
		NativeName:  locale.NativeName,
		IsDefault:   locale.IsDefault,
		Scope:       locale.Scope,
		EntryCount:  len(locale.Entries),
	}
}

func (locale *I18nLocale) ToDTO(fallbackUsed bool) I18nLocaleResponseDto {
	entries := make(map[string]string, len(locale.Entries))
	for key, value := range locale.Entries {
		entries[key] = value.Value
	}
	return I18nLocaleResponseDto{
		LocaleCode:   locale.Key,
		DisplayName:  locale.DisplayName,
		NativeName:   locale.NativeName,
		IsDefault:    locale.IsDefault,
		Scope:        locale.Scope,
		EntryCount:   len(entries),
		Entries:      entries,
		FallbackUsed: fallbackUsed,
	}
}
