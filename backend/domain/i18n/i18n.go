// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
// SPDX-License-Identifier: Apache-2.0

package i18n

import (
	"time"

	"github.com/eclipse-disuko/disuko/domain"
)

// I18nLocale represents a language locale with all translation entries for that language.
type I18nLocale struct {
	domain.RootEntity `bson:",inline"`
	domain.SoftDelete `bson:",inline"`

	LocaleCode  string
	DisplayName string
	NativeName  string
	IsDefault   bool
	Entries     map[string]*I18nEntry
	Scope       string
}

// I18nEntry represents a single translation entry within a locale.
type I18nEntry struct {
	Key         string
	Value       string
	Description string
	CreatedBy   string
	CreatedAt   time.Time
	UpdatedBy   string
	UpdatedAt   time.Time
}

// NewI18nLocale creates a new I18nLocale entity.
func NewI18nLocale(localeCode string) *I18nLocale {
	return &I18nLocale{
		RootEntity: domain.NewRootEntityWithKey(localeCode),
		LocaleCode: localeCode,
		Entries:    make(map[string]*I18nEntry),
	}
}

// NewI18nEntry creates a new translation entry.
func NewI18nEntry(key, value, description string) *I18nEntry {
	return &I18nEntry{
		Key:         key,
		Value:       value,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// SetEntry updates or adds a translation entry.
func (locale *I18nLocale) SetEntry(entry *I18nEntry) {
	if locale.Entries == nil {
		locale.Entries = make(map[string]*I18nEntry)
	}
	locale.Entries[entry.Key] = entry
}

// GetEntry retrieves a translation entry by key.
func (locale *I18nLocale) GetEntry(key string) *I18nEntry {
	if locale.Entries == nil {
		return nil
	}
	return locale.Entries[key]
}

// RemoveEntry removes a translation entry by key.
func (locale *I18nLocale) RemoveEntry(key string) {
	if locale.Entries == nil {
		return
	}
	delete(locale.Entries, key)
}

// GetEntryCount returns the number of translation entries.
func (locale *I18nLocale) GetEntryCount() int {
	if locale.Entries == nil {
		return 0
	}
	return len(locale.Entries)
}
