// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {defineStore} from 'pinia';
import {getStoreLanguage, isLanguage} from '@shared/utils/language';
import {useLocale} from 'vuetify/framework';

const resolveInitialAppLanguage = (locale: string): Lang => {
  const storeLanguage = getStoreLanguage();

  if (storeLanguage) {
    return storeLanguage;
  }

  if (isLanguage(locale)) {
    return locale as Lang;
  }

  return 'en';
};

export const useLanguageStore = defineStore('language', () => {
  const {locale} = useI18n();
  const vuetifyLocale = useLocale();

  const appLanguage = ref<Lang>(resolveInitialAppLanguage(locale.value));

  const setLanguage = (language: Lang) => {
    appLanguage.value = language;
  };

  const toggleLanguage = () => {
    setLanguage(appLanguage.value === 'en' ? 'de' : 'en');
  };

  const shareLanguage = () => {
    locale.value = appLanguage.value;
    vuetifyLocale.current.value = appLanguage.value;
    localStorage.setItem('appLanguage', appLanguage.value);
  };

  const initializeLanguage = shareLanguage;

  watch(appLanguage, () => {
    shareLanguage();
  });

  return {
    appLanguage,
    setLanguage,
    toggleLanguage,
    initializeLanguage,
  };
});
