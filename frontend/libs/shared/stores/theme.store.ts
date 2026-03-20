// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useStorage} from '@vueuse/core';
import {defineStore} from 'pinia';
import {computed, watch} from 'vue';
import {useTheme} from 'vuetify';

const storeKey = 'disco-theme';

export enum ThemeColor {
  light = 'light',
  dark = 'dark',
}

type ThemeColors = keyof typeof ThemeColor;

export const useThemeStore = defineStore(storeKey, () => {
  const theme = useTheme();
  const store = useStorage(storeKey, theme.global.name.value);

  const current = computed(() =>
    Object.keys(ThemeColor).includes(ThemeColor[store.value as ThemeColors])
      ? ThemeColor[store.value as ThemeColors]
      : ThemeColor.light,
  );

  const toggle = () => {
    store.value = store.value === ThemeColor.light ? ThemeColor.dark : ThemeColor.light;
  };

  const updateHtmlClass = (color: string) => {
    const html = document.documentElement;
    html.classList.remove(ThemeColor.light, ThemeColor.dark);
    if (Object.keys(ThemeColor).includes(ThemeColor[color as ThemeColors])) {
      html.classList.add(color);
    }
  };

  updateHtmlClass(store.value);

  watch(store, (color: string) => {
    if (Object.keys(ThemeColor).includes(ThemeColor[color as ThemeColors])) {
      theme.global.name.value = color;
      updateHtmlClass(color);
    }
  });

  return {
    current,
    toggle,
  };
});
