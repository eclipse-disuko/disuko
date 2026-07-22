// Plugins
import basicSsl from '@vitejs/plugin-basic-ssl';
import Vue from '@vitejs/plugin-vue';
import Components from 'unplugin-vue-components/vite';
import Markdown from 'unplugin-vue-markdown/vite';
import Vuetify, {transformAssetUrls} from 'vite-plugin-vuetify';
// Utilities
import VueI18nPlugin from '@intlify/unplugin-vue-i18n/vite';
import tailwindcss from '@tailwindcss/vite';
import fs from 'node:fs';
import path from 'node:path';
import {visualizer} from 'rollup-plugin-visualizer';
import {type CommonServerOptions, defineConfig} from 'vite';
import vueDevTools from 'vite-plugin-vue-devtools';
import {apiProxySettings} from '../../libs/shared/utils/proxySettings.js';
import packageInfo from '../../package.json';

export default defineConfig(({mode}) => {
  process.env.VITE_BUILD_DATE = new Date().toLocaleDateString('de-DE');
  process.env.VITE_VERSION = packageInfo.version || 'n/a';
  process.env.VITE_COMMIT = process.env.VITE_COMMIT || 'Local Build';
  process.env.VITE_BRANCH = process.env.VITE_BRANCH || 'n/a';

  const environment = (mode || 'development') as keyof typeof apiProxySettings;
  const certPath = path.resolve(__dirname, '../../../backend/server.crt');
  const keyPath = path.resolve(__dirname, '../../../backend/server.key');
  const hasCerts = fs.existsSync(certPath) && fs.existsSync(keyPath);

  const developmentServerConfig: CommonServerOptions = {
    port: 3000,
    host: 'disco.local',
    proxy: apiProxySettings[environment],
    cors: true,
    ...(hasCerts && {
      https: {
        cert: fs.readFileSync(certPath),
        key: fs.readFileSync(keyPath),
      },
    }),
  };

  const isDevelopment = mode === 'development';
  return {
    plugins: [
      Vue({
        template: {transformAssetUrls},
        include: [/\.vue$/, /\.md$/],
      }),
      Markdown({
        /* options */
      }),
      Vuetify(),
      Components({
        dts: true,
        types: [
          {
            from: 'vue-router',
            names: ['RouterLink', 'RouterView'],
          },
        ],
        dirs: [
          path.resolve(__dirname, './../../libs/portal/components'),
          path.resolve(__dirname, './../../libs/portal/layouts'),
          path.resolve(__dirname, './../../libs/shared'),
          path.resolve(__dirname, './../../libs/discocli'),
        ],
        deep: true,
      }),
      VueI18nPlugin({
        include: [
          path.resolve(__dirname, './../../libs/discocli/i18n/locales/**'),
          path.resolve(__dirname, './../../libs/shared/i18n/locales/**'),
        ],
        strictMessage: false,
        runtimeOnly: false,
      }),
      isDevelopment && vueDevTools({launchEditor: 'idea'}),
      !hasCerts && basicSsl(),
      visualizer({
        filename: `bundle-cli-ui.html`,
        open: true,
        template: 'sunburst',
      }),
      tailwindcss(),
    ].filter(Boolean),
    base: '/supplierportal/',
    resolve: {
      alias: {
        '@disclosure-portal': `${path.resolve(__dirname, './../../libs/portal')}/`, // Alias für Hauptprojekt
        '@shared': `${path.resolve(__dirname, './../../libs/shared')}/`, // Alias für Hauptprojekt
        '@cli': path.resolve(__dirname, './../../libs/discocli'),
      },
      extensions: ['.js', '.json', '.jsx', '.mjs', '.ts', '.tsx', '.vue', ''],
    },
    server: developmentServerConfig,
    assetsInclude: ['**/*.md'],
    build: {
      rollupOptions: {
        output: {
          assetFileNames: 'assets/[name].[ext]',
          chunkFileNames: (assetInfo: {name: string}) => {
            const name =
              assetInfo.name.endsWith('.vue_vue_type_style_index_0_lang') ||
              assetInfo.name.endsWith('.vue_vue_type_script_setup_true_lang')
                ? assetInfo.name.split('.')[0]
                : assetInfo.name;
            return `assets/${name}-[hash].js`;
          },
          entryFileNames: '[name].js',
        },
      },
      sourcemap: false, //TODO: Set to false in production
      target: 'esnext',
      cssCodeSplit: false,
    },
  };
});
