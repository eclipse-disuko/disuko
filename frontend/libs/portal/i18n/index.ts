// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {createI18n} from 'vue-i18n';

import sharedDE from '@shared/i18n/locales/de.json';
import sharedEN from '@shared/i18n/locales/en.json';
import de from './locales/de.json';
import en from './locales/en.json';

const storedLang = localStorage.getItem('appLanguage');
const supportedLocales = ['en', 'de'] as const;
const normalizedStoredLang = supportedLocales.includes(storedLang as 'en' | 'de') ? storedLang : null;
const browserLang = navigator.language.toLowerCase().startsWith('de') ? 'de' : 'en';

const i18n = createI18n({
  legacy: false,
  locale: normalizedStoredLang || browserLang,
  fallbackLocale: 'en',
  messages: {
    en: {
      ...sharedEN,
      ...en,
    },
    de: {
      ...sharedDE,
      ...de,
    },
  },
});

export default i18n;
