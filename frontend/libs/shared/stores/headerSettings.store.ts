// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {defineStore} from 'pinia';
import type {HeaderSettings} from '@shared/types/table';

const STORE_KEY = 'gridHeaderSettings';

export const useHeaderSettingsStore = defineStore(STORE_KEY, {
  state: () => ({}) as Record<string, HeaderSettings>,
});
