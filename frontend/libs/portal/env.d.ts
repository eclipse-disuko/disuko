// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import 'vue-router';

export {};

/* eslint-disable */
declare module '*.vue' {
  import type {DefineComponent} from 'vue';
  const component: DefineComponent<{}, {}, any>;
  export default component;
}

declare module 'vue-router' {
  interface RouteMeta {
    // is optional
    helpText?: Record<'de' | 'en', string>;
  }
}
