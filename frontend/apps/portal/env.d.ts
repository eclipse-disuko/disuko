// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import 'vue-router';

/* eslint-disable */
declare module '*.vue' {
  import type {DefineComponent} from 'vue';
  const component: DefineComponent<{}, {}, any>;
  export default component;
}

declare module '*.md';

declare module 'vue-router' {
  interface RouteMeta {
    // is optional
    helpText?: Record<Lang, string>;
    title?: Record<Lang, string>;
  }
}
