// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {defineStore} from 'pinia';
import {reactive, toRefs} from 'vue';

interface TableActionSliderState {
  slideInTimeout: ReturnType<typeof setTimeout> | null;
  sliderWidth: number;
  baseWidth: number;
}

export const useTableActionSliderStore = defineStore('tableActionSlider', () => {
  const baseWidth = 100;

  const sliderInt: TableActionSliderState = {
    slideInTimeout: null,
    sliderWidth: baseWidth,
    baseWidth,
  };

  const sliderState = reactive(sliderInt);

  return {
    ...toRefs(sliderState),
  };
});
