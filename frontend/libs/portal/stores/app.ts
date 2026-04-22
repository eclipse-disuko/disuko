// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import INavItem, {INavItemGroup} from '@disclosure-portal/model/INavItem';
import ITile from '@disclosure-portal/model/ITile';
import {GetDashboardCounts} from '@disclosure-portal/services/admin';
import sessionService from '@disclosure-portal/services/session';
import {LabelsTools} from '@disclosure-portal/utils/Labels';
import {useStorage} from '@vueuse/core';
import {defineStore} from 'pinia';
import {computed, reactive, toRefs, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRoute} from 'vue-router';

const SUPPORTED_LANGUAGES = ['en', 'de'] as const;

const normalizeSupportedLanguage = (language: string): (typeof SUPPORTED_LANGUAGES)[number] => {
  const normalized = language.trim().toLowerCase();
  return normalized === 'de' ? 'de' : 'en';
};

function resolveInitialAppLanguage(): string {
  const stored = localStorage.getItem('appLanguage');
  if (stored && stored.trim()) {
    const normalized = normalizeSupportedLanguage(stored);
    localStorage.setItem('appLanguage', normalized);
    return normalized;
  }
  const locale = useI18n().locale.value;
  const initial = locale && locale.trim() ? normalizeSupportedLanguage(locale) : 'en';
  localStorage.setItem('appLanguage', initial);
  return initial;
}

export const useAppStore = defineStore('app', () => {
  // State as reactive object with type
  const state = reactive({
    appLanguage: resolveInitialAppLanguage(),
    publishedLanguages: [...SUPPORTED_LANGUAGES] as string[],
    publishedLanguageLabels: {} as Record<string, {displayName?: string; nativeName?: string}>,
    LabelsTools: new LabelsTools(),
    tiles: [] as ITile[],
    alternateRender: false,
    navItemGroup: {
      items: [],
      adminItem: {
        title: '',
        path: '',
        iconName: '',
        condition: false,
        active: false,
        tooltip: '',
        subItems: [],
      } as INavItem,
    } as INavItemGroup,
    tokenRefresherIsRunning: false,
    notificationMessage: '',
    dummyDesignMode: false,
    shouldReloadApprovals: false,
  });

  const notificationClosed = useStorage('disco-notification-closed', false, sessionStorage);

  // Actions
  const fetchLabelsTools = async () => {
    try {
      await state.LabelsTools.loadLabels();
    } catch (error) {
      console.error(error);
    }
  };

  const checkIfTokenMustRefresh = async () => {
    try {
      await sessionService.getRefreshAccessToken();
    } catch (error) {
      console.error(error);
    }
    setTimeout(() => checkIfTokenMustRefresh(), 1000 * 60 * 2);
  };

  const setNotification = (msg: string) => {
    state.notificationMessage = msg;
  };

  const setNavItemGroup = (items: INavItem[], adminItems: INavItem[]) => {
    state.navItemGroup.items = items;
    if (adminItems.length > 0) {
      Object.assign(state.navItemGroup.adminItem, {
        title: 'ADMIN_DASHBOARD',
        path: '/dashboard/admin',
        iconName: 'mdi-account-cog',
        condition: true,
        active: false,
        tooltip: 'ADMIN_DASHBOARD',
        subItems: [] as INavItem[],
      });
      state.navItemGroup.adminItem.subItems = adminItems;
    } else {
      Object.assign(state.navItemGroup.adminItem, {
        title: '',
        path: '',
        iconName: '',
        condition: false,
        active: false,
        tooltip: '',
        subItems: [],
      } as INavItem);
    }

    setNavItemActive(route.path);
  };
  const setNavItemActive = (currentPath: string) => {
    if (!state.navItemGroup) return;
    state.navItemGroup.items.forEach((navItem) => {
      navItem.active = currentPath.includes(navItem.path);
    });
    if (state.navItemGroup && state.navItemGroup.adminItem.subItems) {
      state.navItemGroup.adminItem.subItems.forEach((navItem) => {
        navItem.active = currentPath.includes(navItem.path);
      });
      const oneOfAdminSubItemsActive = state.navItemGroup.adminItem.subItems.some((item) => item.active);
      state.navItemGroup.adminItem.active =
        currentPath.includes(state.navItemGroup.adminItem.path) ||
        (state.navItemGroup.adminItem.subItems && oneOfAdminSubItemsActive);
    }
  };
  const route = useRoute();
  const setTiles = (tiles: ITile[]) => {
    state.tiles = [];
    state.tiles.push(...tiles);
  };

  const startTokenRefresher = () => {
    if (state.tokenRefresherIsRunning) {
      return;
    }
    state.tokenRefresherIsRunning = true;
    checkIfTokenMustRefresh().then((r) => {
      console.log('checkIfTokenMustRefresh', r);
    });
  };

  const setLanguage = (language: string) => {
    const normalized = normalizeSupportedLanguage(language);
    state.appLanguage = normalized;
    localStorage.setItem('appLanguage', state.appLanguage);
  };

  const setPublishedLanguages = (languages: Array<string | {code: string; displayName?: string; nativeName?: string}>) => {
    const normalizedObjects = (languages || [])
      .map((item) => {
        if (typeof item === 'string') {
          const code = item.trim().toLowerCase();
          return code ? {code} : null;
        }
        const code = (item.code || '').trim().toLowerCase();
        if (!code) return null;
        return {
          code,
          displayName: item.displayName?.trim(),
          nativeName: item.nativeName?.trim(),
        };
      })
      .filter((item): item is {code: string; displayName?: string; nativeName?: string} => item !== null)
      .filter((item) => SUPPORTED_LANGUAGES.includes(item.code as (typeof SUPPORTED_LANGUAGES)[number]));

    const uniqueByCode = Array.from(new Map(normalizedObjects.map((item) => [item.code, item])).values());
    state.publishedLanguages = [...SUPPORTED_LANGUAGES];
    state.publishedLanguageLabels = uniqueByCode.reduce(
      (acc, item) => {
        acc[item.code] = {displayName: item.displayName, nativeName: item.nativeName};
        return acc;
      },
      {} as Record<string, {displayName?: string; nativeName?: string}>,
    );
    setLanguage(state.appLanguage);
  };

  const toggleLanguage = () => {
    setLanguage(state.appLanguage === 'de' ? 'en' : 'de');
  };

  const setDummyDesignMode = (isDummy: boolean) => {
    state.dummyDesignMode = isDummy;
  };

  const unsetDummyDesignMode = () => {
    state.dummyDesignMode = false;
  };

  watch(
    () => route.path,
    () => {
      setNavItemActive(route.path);
    },
    {immediate: true},
  );

  const setShouldReloadApprovals = (value: boolean) => {
    state.shouldReloadApprovals = value;
  };
  const updateTileCounts = (counts: GetDashboardCounts) => {
    for (const tile of state.tiles) {
      if (tile.url === '/dashboard/tasks') tile.cnt = counts.activeJobCount;
      if (tile.url === '/dashboard/projects') tile.cnt = counts.projectCount;
      if (tile.url === '/dashboard/licenses') tile.cnt = counts.licenseCount;
      if (tile.url === '/dashboard/policyrules') tile.cnt = counts.policyRuleCount;
    }
  };

  // Getters
  const getLabelsTools = computed(() => state.LabelsTools);
  const getAppLanguage = computed(() => state.appLanguage);
  const getPublishedLanguages = computed(() => state.publishedLanguages);
  const getPublishedLanguageLabels = computed(() => state.publishedLanguageLabels);

  return {
    // State
    ...toRefs(state),
    notificationClosed,

    // Actions
    updateTileCounts,
    fetchLabelsTools,
    checkIfTokenMustRefresh,
    setNotification,
    setNavItemGroup,
    setTiles,
    startTokenRefresher,
    toggleLanguage,
    setLanguage,
    setPublishedLanguages,
    setDummyDesignMode,
    unsetDummyDesignMode,
    setShouldReloadApprovals,

    // Getters
    getLabelsTools,
    getAppLanguage,
    getPublishedLanguages,
    getPublishedLanguageLabels,
  };
});
