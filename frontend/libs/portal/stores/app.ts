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
import {useRoute} from 'vue-router';

export const useAppStore = defineStore('app', () => {
  // State as reactive object with type
  const state = reactive({
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
    setDummyDesignMode,
    unsetDummyDesignMode,
    setShouldReloadApprovals,

    // Getters
    getLabelsTools,
  };
});
