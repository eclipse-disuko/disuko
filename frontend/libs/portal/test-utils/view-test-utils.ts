// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {vuetifyStubs} from '@disclosure-portal/test-utils/vuetify-stubs';
import {createTestingPinia} from '@pinia/testing';
import {Component, ComponentPropsOptions} from 'vue';
import {config, mount, MountingOptions, VueWrapper} from '@vue/test-utils';
import {vi} from 'vitest';

// vitest.setup.ts installs both a real vue-i18n plugin and a global `$t` mock. Views call
// useI18n() themselves, and vue-i18n's own mixin then collides with the mocked `$t` binding
// ("Cannot mutate <script setup> binding $t"), so the global mock must not be present for them.
if (config.global.mocks && '$t' in config.global.mocks) {
  delete config.global.mocks.$t;
}

interface MountViewOptions<T> {
  childStubs?: Record<string, unknown>;
  piniaOptions?: Parameters<typeof createTestingPinia>[0];
  mountOptions?: MountingOptions<T>;
  /** Runs against the testing pinia instance right after creation, before mount — for setting up
   * store state (e.g. via a store action) that must be in place before the component's first render. */
  beforePiniaMount?: (pinia: ReturnType<typeof createTestingPinia>) => void;
}

export const mountView = <T extends Component>(component: T, options: MountViewOptions<ComponentPropsOptions> = {}) => {
  const {childStubs = {}, piniaOptions, mountOptions, beforePiniaMount} = options;

  const pinia = createTestingPinia({
    createSpy: vi.fn,
    stubActions: false,
    ...piniaOptions,
  });

  beforePiniaMount?.(pinia);

  const wrapper: VueWrapper = mount(component, {
    ...mountOptions,
    global: {
      plugins: [pinia, ...(mountOptions?.global?.plugins ?? [])],
      stubs: {
        ...vuetifyStubs,
        ...childStubs,
        ...mountOptions?.global?.stubs,
      },
      ...mountOptions?.global,
    },
  } as MountingOptions<T>);

  return {wrapper, pinia};
};
