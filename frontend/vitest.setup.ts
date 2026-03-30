// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import en from '@disclosure-portal/i18n/locales/en.json';
import {config} from '@vue/test-utils';
import {vi} from 'vitest';
import {createI18n} from 'vue-i18n';

global.ResizeObserver = vi.fn().mockImplementation(() => ({
  observe: vi.fn(),
  unobserve: vi.fn(),
  disconnect: vi.fn(),
}));

const i18n = createI18n({
  legacy: false,
  locale: 'en',
  messages: {
    en: en,
  },
});

// Mock vue-router
vi.mock('vue-router', () => ({
  useRoute: vi.fn(() => ({
    path: '/',
    name: 'home',
    params: {},
    query: {},
    hash: '',
    fullPath: '/',
    matched: [],
    meta: {},
    redirectedFrom: undefined,
  })),
  useRouter: vi.fn(() => ({
    push: vi.fn(),
    replace: vi.fn(),
    go: vi.fn(),
    back: vi.fn(),
    forward: vi.fn(),
  })),
}));

config.global.plugins = [i18n];
config.global.mocks = {
  $t: (key: string) => key,
};
