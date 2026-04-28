// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export function isLanguage(locale: string | null): Lang | '' {
  if (locale === 'en' || locale === 'de') {
    return locale;
  }
  return '';
}

export function getStoreLanguage(): Lang | '' {
  const language = localStorage.getItem('appLanguage');
  return isLanguage(language);
}
