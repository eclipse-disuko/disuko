// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

// For global usage of dialogs
import {defineStore} from 'pinia';
import {reactive, toRefs} from 'vue';

export const useDialogStore = defineStore('dialogStore', () => {
  const state = reactive({
    isSettingsDialogOpen: false,
  });
  return {...toRefs(state)};
});
