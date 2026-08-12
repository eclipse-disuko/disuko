// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import config from '@shared/utils/config';
import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {nextTick} from 'vue';
import {beforeEach, describe, expect, it} from 'vitest';
import FeatureFlags from '../FeatureFlags.vue';

// Note: FeatureFlags.vue re-initializes `featureFlags` from sessionStorage inside onMounted, which
// replaces the array the reactive-ref points to *after* the first render. The DOM must be given a
// chance to re-render against that second array (an extra `nextTick()`) before interacting with
// switches, otherwise events fire against stale (pre-onMounted) elements.
//
// It also mutates the shared `@shared/utils/config` singleton in place, and re-derives its notion of
// "original" values from that same live singleton whenever sessionStorage's cache is empty. Left
// mutated between tests, this makes an override look like a fresh "original" value on the next mount
// (see genuine bug note in the final report). Reset the booleans this view manages before each test.
const EXCLUDED_KEYS = new Set(['SERVER_URL', 'PUBLIC_API_ENDPOINT', 'OAUTH']);
const resetConfigBooleans = () => {
  Object.entries(config).forEach(([key, value]) => {
    if (typeof value === 'boolean' && !EXCLUDED_KEYS.has(key)) {
      (config as unknown as Record<string, boolean>)[key] = false;
    }
  });
};

describe('FeatureFlags', () => {
  beforeEach(() => {
    sessionStorage.clear();
    resetConfigBooleans();
  });

  it('sets the expected breadcrumbs on mount', () => {
    const {pinia} = mountView(FeatureFlags);

    const breadcrumbs = useBreadcrumbsStore(pinia);
    expect(breadcrumbs.setCurrentBreadcrumbs).toHaveBeenCalledWith([
      ...breadcrumbs.dashboardCrumbs,
      {title: 'Feature Flags'},
    ]);
  });

  it('renders without throwing and shows one card per non-excluded boolean config flag', () => {
    const {wrapper} = mountView(FeatureFlags);

    expect(wrapper.exists()).toBe(true);
    expect(wrapper.findAll('.v-card').length).toBeGreaterThan(0);
  });

  it('overrides a flag when its switch is toggled and persists it to sessionStorage', async () => {
    const {wrapper} = mountView(FeatureFlags);
    await nextTick();

    const switches = wrapper.findAll('input[type="checkbox"]');
    expect(switches.length).toBeGreaterThan(0);

    await switches[0].setValue(true);
    await nextTick();

    const stored = JSON.parse(sessionStorage.getItem('featureFlagOverrides') || '{}');
    expect(stored).toEqual({isProd: true});
    expect(wrapper.findAll('button').some((btn) => btn.text() === 'Reset All Flags')).toBe(true);
  });

  it('resets all overridden flags when "Reset All" is clicked', async () => {
    const {wrapper} = mountView(FeatureFlags);
    await nextTick();

    const switches = wrapper.findAll('input[type="checkbox"]');
    await switches[0].setValue(true);
    await nextTick();

    const resetAllButton = wrapper.findAll('button').find((btn) => btn.text() === 'Reset All Flags');
    expect(resetAllButton).toBeTruthy();
    await resetAllButton!.trigger('click');
    await nextTick();

    expect(sessionStorage.getItem('featureFlagOverrides')).toBeNull();
    expect(wrapper.findAll('button').some((btn) => btn.text() === 'Reset All Flags')).toBe(false);
  });
});
