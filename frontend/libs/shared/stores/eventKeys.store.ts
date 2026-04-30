// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {defineStore} from 'pinia';
import {onMounted, onUnmounted, reactive, toRefs} from 'vue';
import {onKeyDown} from '@vueuse/core';

export const useEventKeysStore = defineStore('eventKeysStore', () => {
  const state = reactive({
    controlIsPressed: false,
    shiftIsPressed: false,
    eventListener: false,
    eventListenerStop: (): void => void 0,
  });

  const releaseKeys = () => {
    state.controlIsPressed = false;
    state.shiftIsPressed = false;
  };

  onMounted(() => {
    if (!state.eventListener) {
      state.eventListenerStop = onKeyDown(
        (e) => e?.shiftKey || e?.ctrlKey,
        (e) => {
          state.shiftIsPressed = Boolean(e?.shiftKey);
          state.controlIsPressed = Boolean(e?.ctrlKey);
        },
      );
      state.eventListener = true;
    }
  });

  onUnmounted(() => {
    if (state.eventListener) {
      state.eventListenerStop();
      state.eventListener = false;
      releaseKeys();
    }
  });

  return {
    ...toRefs(state),
    releaseKeys,
  };
});
