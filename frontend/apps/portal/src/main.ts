/**
 * main.ts
 *
 * Bootstraps Vuetify and other plugins then mounts the App`
 */

import '@mdi/font/css/materialdesignicons.min.css';
import {setup} from '@shared/utils/config';
import dayjs from 'dayjs';
import dayjsPluginUTC from 'dayjs-plugin-utc';
import {createApp} from 'vue';
import App from './App.vue';
import {registerPlugins} from './plugins';
import('./styles/themes/default/markdown.scss');

dayjs.extend(dayjsPluginUTC);
setup().then(() => {
  const app = createApp(App);
  registerPlugins(app);
  app.mount('#app');
});
