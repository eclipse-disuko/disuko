// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export interface I18nLocaleListItem {
  localeCode: string;
  displayName: string;
  nativeName: string;
  isDefault: boolean;
  scope: string;
  entryCount: number;
}

export interface I18nLocaleResponse {
  localeCode: string;
  displayName: string;
  nativeName: string;
  isDefault: boolean;
  scope: string;
  entryCount: number;
  entries: Record<string, string>;
  fallbackUsed: boolean;
}

export interface I18nImportIssue {
  fileName: string;
  key?: string;
  code: string;
  message: string;
}

export interface I18nImportResponse {
  success: boolean;
  validationPassed: boolean;
  locale: string;
  filesProcessed: number;
  totalKeysParsed: number;
  appended: number;
  updated: number;
  unchanged: number;
  errors?: I18nImportIssue[];
}
