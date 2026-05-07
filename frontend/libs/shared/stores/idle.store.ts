// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import IdleInfo from '@shared/types/IdleInfo';
import eventBus from '@shared/utils/eventbus';
import {defineStore} from 'pinia';
import {nextTick, onMounted, onUnmounted, reactive, toRefs} from 'vue';

const defaultShowIdle = false;
const defaultIdleMessage = '';
const defaultProgressUnit = '';
const defaultProgress = -1;

/**
 * Idle store to manage the idle state of the application.
 *
 * Usage:
 * ``` ts
 * const idle = useIdleStore();
 * idle.show('Loading data...', 0);
 *
 * // if updating progress
 *
 * idle.update('Loading data...', 50);
 *
 * // later then
 *
 * idle.hide();
 * ```
 */
export const useIdleStore = defineStore('idle', () => {
  const state = reactive({
    showIdle: defaultShowIdle,
    idleMessage: defaultIdleMessage,
    progressUnit: defaultProgressUnit,
    progress: defaultProgress,
  });

  const onIDLE = ({idle}: {idle: IdleInfo}) => {
    if (idle) {
      state.showIdle = idle.show;
      state.idleMessage = idle.message;
      state.progressUnit = idle.progressUnit;
      state.progress = idle.progress;
    }
  };

  const hide = () => {
    setTimeout(() => {
      state.showIdle = defaultShowIdle;
      state.idleMessage = defaultIdleMessage;
      state.progressUnit = defaultProgressUnit;
      state.progress = defaultProgress;
    });
  };

  const update = (message?: string, progress?: number, progressUnit?: string) => {
    nextTick(() => {
      state.idleMessage = message ?? '';
      state.progressUnit = progressUnit ?? '';
      state.progress = progress ?? -1;
    });
  };

  const show = (message?: string, progress?: number, progressUnit?: string) => {
    nextTick(() => {
      state.showIdle = true;
      state.idleMessage = message ?? '';
      state.progressUnit = progressUnit ?? '';
      state.progress = progress ?? -1;
    });
  };

  onMounted(() => {
    eventBus.on('on-idle', onIDLE);
  });

  onUnmounted(() => {
    eventBus.off('on-idle', onIDLE);
  });

  return {
    ...toRefs(state),
    show,
    hide,
    update,
  };
});
