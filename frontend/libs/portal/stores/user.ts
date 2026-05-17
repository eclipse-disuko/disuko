// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import INavItem from '@disclosure-portal/model/INavItem';
import ITile from '@disclosure-portal/model/ITile';
import SimpleProfileData from '@disclosure-portal/model/ProfileData';
import {Rights} from '@disclosure-portal/model/Rights';
import {UserDto} from '@disclosure-portal/model/Users';
import {useAppStore} from '@disclosure-portal/stores/app';
import {defineStore} from 'pinia';

export const createNavItemsGroup = function () {
  const rights = useUserStore().getRights;
  const items = [] as INavItem[];
  const adminItems = [] as INavItem[];
  items.push(
    {
      title: 'DASHBOARD',
      path: '/dashboard/home',
      iconName: 'mdi-view-dashboard',
      condition: true,
      active: false,
      tooltip: 'DASHBOARD',
    } as INavItem,
    {
      title: 'PROJECTS',
      path: '/dashboard/projects',
      iconName: 'mdi-list-box',
      condition: true,
      active: false,
      tooltip: 'PROJECTS',
    } as INavItem,
    {
      title: 'TASKS',
      path: '/dashboard/tasks',
      iconName: 'mdi-file-check',
      condition: true,
      active: false,
      tooltip: 'TASKS',
    } as INavItem,
    {
      title: 'LICENSES',
      path: '/dashboard/licenses',
      iconName: 'mdi-shield-check',
      condition: rights.allowLicense && rights.allowLicense.read,
      active: false,
      tooltip: 'LICENSES',
    } as INavItem,
    {
      title: 'POLICY_RULES',
      path: '/dashboard/policyrules',
      iconName: 'mdi-bank',
      condition: true,
      active: false,
      tooltip: 'POLICY_RULES',
    } as INavItem,
    {
      title: 'ANNOUNCEMENTS',
      path: '/dashboard/announcements',
      iconName: 'mdi-bullhorn',
      condition: true,
      active: false,
      tooltip: 'ANNOUNCEMENTS',
    } as INavItem,
  );

  if (rights.isAnyOfAdmin()) {
    adminItems.push(
      {
        title: 'AllProjects',
        path: '/dashboard/admin/projects',
        iconName: 'mdi-list-box-outline',
        condition: rights.hasAllProjectsReadonly(),
        active: false,
      } as INavItem,
      {
        title: 'LABELS',
        path: '/dashboard/admin/labels',
        iconName: 'mdi-label-outline',
        condition: rights.hasLabelAccess(),
        active: false,
      } as INavItem,
      {
        title: 'DB_TITLE_REVIEW_TEMPLATES',
        path: '/dashboard/admin/templates/review',
        iconName: 'mdi-text-box-outline',
        condition: rights.hasReviewTemplatesAcces(),
        active: false,
      } as INavItem,
      {
        title: 'Schemes',
        path: '/dashboard/admin/schemas',
        iconName: 'mdi-file-tree-outline',
        condition: rights.hasSchemaAccess(),
        active: false,
      } as INavItem,
      {
        title: 'CLASSIFICATIONS',
        path: '/dashboard/admin/classifications',
        iconName: 'mdi-grain',
        condition: rights.hasClassificationsAccess(),
        active: false,
      } as INavItem,
      {
        title: 'TOOLS',
        path: '/dashboard/admin/tools',
        iconName: 'mdi-hammer-screwdriver',
        condition: rights.hasToolsAccess() || rights.hasSampleDataAccess(),
        active: false,
      } as INavItem,
      {
        title: 'ADMIN_JOBS',
        path: '/dashboard/admin/jobs',
        iconName: 'mdi-window-shutter-cog',
        condition: rights.isApplicationAdmin(),
        active: false,
      } as INavItem,
      {
        title: 'USERS',
        path: '/dashboard/admin/users',
        iconName: 'mdi-account-multiple-outline',
        condition: rights.hasUsersAccess(),
        active: false,
      } as INavItem,
      {
        title: 'Analytics',
        path: '/dashboard/analytics/overview',
        iconName: 'mdi-chart-box-outline',
        condition: rights.isProjectAnalyst(),
        active: false,
      } as INavItem,
      {
        title: 'CHECKLISTS',
        path: '/dashboard/admin/checklist',
        iconName: 'mdi-format-list-checks',
        condition: rights.isDomainAdmin() || rights.isFOSSOffice(),
        active: false,
      } as INavItem,
      {
        title: 'NEWSBOX',
        path: '/dashboard/admin/newsbox',
        iconName: 'mdi-newspaper-variant',
        condition: rights.isDomainAdmin() || rights.isApplicationAdmin(),
        active: false,
      } as INavItem,
      {
        title: 'Internationalization',
        path: '/dashboard/admin/i18n',
        iconName: 'mdi-translate',
        condition: rights.isApplicationAdmin(),
        active: false,
      } as INavItem,
      {
        title: 'TAB_ADMIN_USER_MANAGEMENT',
        path: '/dashboard/admin/userManagement',
        iconName: 'mdi-account-cog-outline',
        condition: rights.isDomainAdmin(),
        active: false,
      } as INavItem,
    );
  }

  // TODO: Using a store outside of a component, composable or store
  useAppStore().setNavItemGroup(items, adminItems);

  const res: ITile[] = [];
  const addTile = (title: string, url: string, icon: string, cnt?: number) => {
    res.push({
      cnt: cnt || -1,
      title,
      url,
      icon,
      expand: false,
      expandGroup: false,
    } as ITile);
  };
  addTile('PROJECTS', '/dashboard/projects', 'mdi-list-box-outline', 0);
  addTile('TASKS', '/dashboard/tasks', '');
  if (rights.allowLicense.read) {
    addTile('Licenses', '/dashboard/licenses', 'mdi-shield-check-outline', 0);
  }
  addTile('Policies', '/dashboard/policyrules', 'mdi-bank-outline', 0);
  // if (rights.isDomainAdmin() || rights.isProjectAnalyst()) {
  //   addTile('ANALYTICS', '/dashboard/analytics/overview', 'mdi-chart-box-outline');
  // }
  if (rights.isAnyOfAdmin()) {
    addTile('ADMIN', '/dashboard/admin', 'mdi-cog-outline');
  }
  // TODO: Using a store outside of a component, composable or store
  useAppStore().setTiles(res);
};
export const useUserStore = defineStore('user', {
  state: () => ({
    simpleProfileData: {
      rights: {} as Rights,
      profile: {} as UserDto,
      allowed: true,
    } as SimpleProfileData,
  }),
  actions: {
    setSimpleProfileData(simpleProfileData: SimpleProfileData) {
      Object.assign(this.simpleProfileData, simpleProfileData);
      this.simpleProfileData.rights = new Rights();
      Object.assign(this.simpleProfileData.rights, simpleProfileData.rights);
      this.simpleProfileData.allowed = true;
    },
    clear() {
      this.simpleProfileData.allowed = false;
    },
    updateTermsOfUse(termsOfUse: boolean, termsOfUseDate: string) {
      this.simpleProfileData.profile.termsOfUse = termsOfUse;
      this.simpleProfileData.profile.termsOfUseDate = termsOfUseDate;
    },
  },
  getters: {
    getRights(): Rights {
      return this.simpleProfileData.rights;
    },
    getProfile(): UserDto {
      return this.simpleProfileData.profile;
    },
  },
});
