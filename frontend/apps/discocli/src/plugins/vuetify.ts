/**
 * plugins/vuetify.ts
 *
 * Framework documentation: https://vuetifyjs.com`
 */

// Styles
import '../styles/layers.css';
import 'vuetify/styles';

// Composables
import DayJsAdapter from '@date-io/dayjs';
import deDatesLang from 'dayjs/locale/de';
import {createVuetify} from 'vuetify';
import {mdi} from 'vuetify/iconsets/mdi';
import {VDateInput} from 'vuetify/labs/VDateInput';
import {de, en} from 'vuetify/locale';
import dark from './dark';
import light from './light';
import {getStoreLanguage} from '@shared/utils/language';
import {getDefaultTheme} from '@shared/utils/theme';

// https://vuetifyjs.com/en/introduction/why-vuetify/#feature-guides
export default createVuetify({
  locale: {
    locale: getStoreLanguage() || 'en',
    fallback: 'en',
    messages: {en, de},
  },
  components: {
    VDateInput,
  },
  date: {
    adapter: DayJsAdapter,
    locale: {
      en: deDatesLang,
      de: deDatesLang,
    },
  },
  display: {
    thresholds: {
      xs: 0,
      sm: 640,
      md: 768,
      lg: 1024,
      xl: 1280,
    },
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
