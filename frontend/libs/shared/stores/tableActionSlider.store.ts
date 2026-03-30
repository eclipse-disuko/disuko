// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {defineStore} from 'pinia';
import {ref} from 'vue';

export const useTableActionSliderStore = defineStore('tableActionSlider', () => {
  const slideInTimeout = ref<ReturnType<typeof setTimeout> | null>(null);

  return {
    slideInTimeout,
  };
});
