// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useTableActionSliderStore} from '@shared/stores/tableActionSlider.store';
import {storeToRefs} from 'pinia';
import {computed, ref, watch} from 'vue';

export const useTableActionSlider = () => {
  // Need the store to have one timeout across all sliders
  const tableActionSliderStore = useTableActionSliderStore();
  const {slideInTimeout, baseWidth, sliderWidth} = storeToRefs(tableActionSliderStore);

  const buttonWidth = 40;
  const spaceAfter = 20;

  const buttonsLength = ref<number>(1);
  const slideInTimer = ref<number>(0);

  const slideOutAction = ref<() => unknown>(() => {});
  const slideInAction = ref<() => unknown>(() => {});

  const setupTableActionSlider = (
    newButtonsLength: number,
    newSlideOutAction?: () => unknown,
    newSlideInAction?: () => unknown,
    newBaseWidth?: number,
  ) => {
    if (newBaseWidth && newBaseWidth > baseWidth.value) {
      baseWidth.value = newBaseWidth;
    }

    if (newButtonsLength && newButtonsLength > buttonsLength.value && buttonsLength.value < newButtonsLength) {
      buttonsLength.value = newButtonsLength;
    }

    if (newSlideOutAction) {
      slideOutAction.value = newSlideOutAction;
    }

    if (newSlideInAction) {
      slideInAction.value = newSlideInAction;
    }
  };

  const expandedMaxWidth = computed(() => buttonsLength.value * buttonWidth + spaceAfter);

  const startSlideInTimer = () => {
    slideInTimer.value = 300;

    if (slideInTimeout.value) {
      clearTimeout(slideInTimeout.value);
      slideInTimeout.value = null;
    }

    slideInTimeout.value = setTimeout(() => {
      slideInTimer.value = 0;
    }, slideInTimer.value);
  };

  const slideOut = () => {
    sliderWidth.value = expandedMaxWidth.value;
    slideOutAction.value();
  };

  const slideIn = () => {
    sliderWidth.value = baseWidth.value;
    slideInAction.value();
  };

  const stopSlideInTimerAndSlideOut = () => {
    if (slideInTimeout.value) {
      clearTimeout(slideInTimeout.value);
      slideInTimeout.value = null;
    }

    slideOut();
  };

  watch(slideInTimer, () => {
    if (slideInTimer.value === 0) {
      slideIn();
    }
  });

  return {
    baseWidth,
    sliderWidth,
    expandedMaxWidth,
    slideOut,
    slideIn,
    setupTableActionSlider,
    startSlideInTimer,
    stopSlideInTimerAndSlideOut,
  };
};
