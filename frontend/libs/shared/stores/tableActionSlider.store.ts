// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {defineStore} from 'pinia';
import {reactive, toRefs} from 'vue';

export const useTableActionSliderStore = defineStore('tableActionSlider', () => {
  const sliderState = reactive({
    baseWidth: 0,
    slideInTimeout: null as ReturnType<typeof setTimeout> | null,
    sliderWidth: 0,
  });

  return {
    ...toRefs(sliderState),
  };
});
