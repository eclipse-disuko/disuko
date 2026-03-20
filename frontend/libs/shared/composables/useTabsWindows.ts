// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {CONTROL_IS_PRESSED, SHIFT_IS_PRESSED} from '@disclosure-portal/keyState';
import {MaybeRef} from '@vueuse/core';
import {computed, isRef, onMounted, ref, toRef, watch} from 'vue';
import {useRoute, useRouter} from 'vue-router';

export const useTabsWindows = (
  baseUrlUnknown: MaybeRef<string>,
  urlPartList: string[],
  maybeDefaultUrl?: string,
  urlSuffix?: MaybeRef<string>,
) => {
  const router = useRouter();
  const route = useRoute();

  const baseUrl = isRef(baseUrlUnknown) ? baseUrlUnknown : toRef(baseUrlUnknown);
  const suffix = urlSuffix ? (isRef(urlSuffix) ? urlSuffix : toRef(urlSuffix)) : toRef('');

  const changeUrlForTab = async (tabUrl: string) => {
    if (CONTROL_IS_PRESSED || SHIFT_IS_PRESSED) {
      // open in new browser tab or window
      return;
    }
    // Handle refs to tabs if necessary
    const finalUrl = suffix.value ? `${tabUrl}/${suffix.value}` : tabUrl;
    return router.push(finalUrl);
  };

  const currentTabUrlMatch = (currentTab: string) => isTabFromUrl.value && tabFromUrl.value === currentTab;

  const isTabFromUrl = computed(() => Boolean(route.params?.tab));

  const tabFromUrl = computed(() => {
    if (isTabFromUrl.value) {
      return Array.isArray(route.params.tab) ? route.params.tab[0] : route.params.tab;
    } else {
      return '';
    }
  });

  const tabUrl = computed(() => {
    const tabs: Record<string, string> = {};
    urlPartList.forEach((url) => {
      tabs[url] = `${baseUrl.value}/${url}`;
    });
    return tabs;
  });

  const defaultUrl = maybeDefaultUrl || urlPartList[0];

  const selectedTab = ref(defaultUrl);

  onMounted(async () => {
    if (!tabFromUrl.value) {
      await changeUrlForTab(tabUrl.value[selectedTab.value]);
    } else if (!currentTabUrlMatch(selectedTab.value)) {
      selectedTab.value = tabFromUrl.value;
    }
  });

  watch(selectedTab, async (newTab) => {
    if (!currentTabUrlMatch(newTab)) {
      await changeUrlForTab(tabUrl.value[newTab]);
    }
  });

  watch(tabFromUrl, async (newTab) => {
    if (newTab !== selectedTab.value) {
      selectedTab.value = newTab || defaultUrl;
    }
  });

  return {
    tabUrl,
    selectedTab,
    isTabFromUrl,
    tabFromUrl,
    currentTabUrlMatch,
    changeUrlForTab,
  };
};
