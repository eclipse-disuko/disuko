// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {defineStore} from 'pinia';
import {computed, reactive, toRefs} from 'vue';

export const useNotificationStore = defineStore('notification', () => {
  const state = reactive({
    message: '',
    dismissedText: '',
    closed: false,
  });

  const visible = computed(() => !state.closed && Boolean(state.message));

  return {
    ...toRefs(state),
    visible,
  };
});
