/**
 * plugins/index.ts
 *
 * Automatically included in `./src/main.ts`
 */

import i18n from '@cli/i18n';
import {createPinia} from 'pinia';
import type {App} from 'vue';
import router from '../router';
import vuetify from './vuetify';

export function registerPlugins(app: App) {
  app.use(createPinia());
  app.use(vuetify);
  app.use(i18n);
  app.use(router);
}
