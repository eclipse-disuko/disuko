// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

/**
 * plugins/index.ts
 *
 * Automatically included in `./src/main.ts`
 */

import i18n from '@disclosure-portal/i18n';
import config from '@shared/utils/config';
import {createPinia} from 'pinia';
import type {App} from 'vue';
import router from '../router';
import {ConfigSymbol} from '../types/symbols';
import vuetify from './vuetify';

export function createDiscoPinia() {
  return createPinia();
}

export function registerPlugins(app: App) {
  app.use(createDiscoPinia());
  app.provide(ConfigSymbol, config);
  app.config.globalProperties.$config = config;
  app.use(vuetify);
  app.use(i18n);
  app.use(router);
}
