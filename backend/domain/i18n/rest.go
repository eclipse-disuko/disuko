// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
// SPDX-License-Identifier: Apache-2.0

package i18n

type I18nLocaleResponseDto struct {
	LocaleCode   string            `json:"localeCode"`
	DisplayName  string            `json:"displayName"`
	NativeName   string            `json:"nativeName"`
	IsDefault    bool              `json:"isDefault"`
	Scope        string            `json:"scope"`
	EntryCount   int               `json:"entryCount"`
	Entries      map[string]string `json:"entries"`
	FallbackUsed bool              `json:"fallbackUsed"`
}

type I18nLocaleListResponseDto struct {
	LocaleCode  string `json:"localeCode"`
	DisplayName string `json:"displayName"`
	NativeName  string `json:"nativeName"`
	IsDefault   bool   `json:"isDefault"`
	Scope       string `json:"scope"`
	EntryCount  int    `json:"entryCount"`
}

type I18nTranslationResponseDto struct {
	LocaleCode   string `json:"localeCode"`
	RequestedKey string `json:"requestedKey"`
	Value        string `json:"value"`
	FallbackUsed bool   `json:"fallbackUsed"`
}

type I18nTranslationUpsertRequestDto struct {
	Value       string `json:"value"`
	Description string `json:"description"`
}

type I18nLocaleUpsertRequestDto struct {
	DisplayName string `json:"displayName"`
	NativeName  string `json:"nativeName"`
	IsDefault   bool   `json:"isDefault"`
	Scope       string `json:"scope"`
}

type I18nMigrationLocaleResultDto struct {
	LocaleCode  string `json:"localeCode"`
	SourceCount int    `json:"sourceCount"`
	TargetCount int    `json:"targetCount"`
	Upserted    int    `json:"upserted"`
	Verified    bool   `json:"verified"`
}

type I18nMigrationResponseDto struct {
	DryRun        bool                           `json:"dryRun"`
	Scope         string                         `json:"scope"`
	IncludeShared bool                           `json:"includeShared"`
	Locales       []I18nMigrationLocaleResultDto `json:"locales"`
	Success       bool                           `json:"success"`
}

type I18nImportIssueDto struct {
	FileName string `json:"fileName"`
	Key      string `json:"key,omitempty"`
	Code     string `json:"code"`
	Message  string `json:"message"`
}

type I18nImportResponseDto struct {
	Success          bool                 `json:"success"`
	ValidationPassed bool                 `json:"validationPassed"`
	Locale           string               `json:"locale"`
	FilesProcessed   int                  `json:"filesProcessed"`
	TotalKeysParsed  int                  `json:"totalKeysParsed"`
	Appended         int                  `json:"appended"`
	Updated          int                  `json:"updated"`
	Unchanged        int                  `json:"unchanged"`
	Errors           []I18nImportIssueDto `json:"errors,omitempty"`
}
