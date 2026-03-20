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
