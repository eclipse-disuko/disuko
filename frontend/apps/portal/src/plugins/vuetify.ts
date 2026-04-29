// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

/**
 * plugins/vuetify.ts
 *
 * Framework documentation: https://vuetifyjs.com`
 */

// Styles
import 'vuetify/styles';

// Composables
import {createVuetify} from 'vuetify';
import {mdi} from 'vuetify/iconsets/mdi';
import dark from './dark';
import light from './light';
import {de, en} from 'vuetify/locale';
import {getDefaultTheme} from '@shared/utils/theme';
import {getStoreLanguage} from '@shared/utils/language';

// https://vuetifyjs.com/en/introduction/why-vuetify/#feature-guides
export default createVuetify({
  locale: {
    locale: getStoreLanguage() || 'en',
    fallback: 'en',
    messages: {en, de},
  },
  theme: {
    defaultTheme: getDefaultTheme(),
    themes: {
      dark,
      light,
    },
  },

  icons: {
    defaultSet: 'mdi',
    sets: {
      mdi,
    },
  },
});
