import {createI18n} from 'vue-i18n';

import sharedDE from '@shared/i18n/locales/de.json';
import sharedEN from '@shared/i18n/locales/en.json';
import de from './locales/de.json';
import en from './locales/en.json';

const i18n = createI18n({
  legacy: false,
  locale: localStorage.getItem('appLanguage') || navigator.language.includes('en') ? 'en' : navigator.language || 'en',
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
