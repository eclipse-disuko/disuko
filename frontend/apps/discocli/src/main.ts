/**
 * main.ts
 *
 * Bootstraps Vuetify and other plugins then mounts the App`
 */

import '@mdi/font/css/materialdesignicons.min.css';
import './styles/layers.css';
import './styles/tailwind.css';
import {registerPlugins} from './plugins';

import dayjs from 'dayjs';
import dayjsPluginUTC from 'dayjs/plugin/utc';
import {createApp} from 'vue';
import App from './App.vue';
import('./styles/themes/default/markdown.scss');

dayjs.extend(dayjsPluginUTC);

const app = createApp(App);

registerPlugins(app);

app.mount('#app');
