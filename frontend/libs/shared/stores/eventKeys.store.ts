// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {defineStore} from 'pinia';
import {reactive, toRefs, watch} from 'vue';
import {onKeyDown, onKeyUp, useWindowFocus} from '@vueuse/core';

export const useEventKeysStore = defineStore('eventKeysStore', () => {
  const state = reactive({
    controlIsPressed: false,
    shiftIsPressed: false,
    eventListenerDown: false,
    eventListenerUp: false,
    windowFocus: false,
    eventListenerDownStop: (): void => void 0,
    eventListenerUpStop: (): void => void 0,
  });

  const windowFocus = useWindowFocus();

  const releaseKeys = () => {
    state.controlIsPressed = false;
    state.shiftIsPressed = false;
  };

  watch([() => state.shiftIsPressed, () => state.controlIsPressed, windowFocus], () => {
    if (
      !windowFocus.value &&
      state.eventListenerDown &&
      state.eventListenerUp &&
      (state.shiftIsPressed || state.controlIsPressed)
    ) {
      releaseKeys();
    }
  });

  const initEventKeyStore = () => {
    if (!state.eventListenerDown) {
      state.eventListenerDownStop = onKeyDown(true, (e) => {
        if (e?.shiftKey) {
          state.shiftIsPressed = true;
        }
        if (e?.ctrlKey) {
          state.controlIsPressed = true;
        }
      });
      state.eventListenerDown = true;
    }

    if (!state.eventListenerUp) {
      state.eventListenerUpStop = onKeyUp(true, (e) => {
        if (!e?.shiftKey) {
          state.shiftIsPressed = false;
        }
        if (!e?.ctrlKey) {
          state.controlIsPressed = false;
        }
      });
      state.eventListenerUp = true;
    }
  };

  return {
    ...toRefs(state),
    initEventKeyStore,
    releaseKeys,
  };
});
