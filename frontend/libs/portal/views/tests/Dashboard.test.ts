// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {Rights} from '@shared/user/models/Rights';
import {useUserStore} from '@shared/user/stores/user.store';
import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {reactive, ref} from 'vue';
import {flushPromises} from '@vue/test-utils';
import {beforeEach, describe, expect, it, vi} from 'vitest';
import Dashboard from '../Dashboard.vue';

const {themeToggleMock, toggleLanguageMock, createNavItemsGroupMock, logoutMock} = vi.hoisted(() => ({
  themeToggleMock: vi.fn(),
  toggleLanguageMock: vi.fn(),
  createNavItemsGroupMock: vi.fn(),
  logoutMock: vi.fn(),
}));

vi.mock('@shared/stores/theme.store', () => ({
  ThemeColor: {light: 'light', dark: 'dark'},
  useThemeStore: () => ({current: 'light', toggle: themeToggleMock}),
}));

vi.mock('@shared/stores/language.store', () => ({
  // Dashboard reads this via storeToRefs(), which pinia only turns into a live ref for
  // properties that are themselves refs on the raw target object - a plain string is skipped.
  useLanguageStore: () =>
    reactive({
      appLanguage: ref('en'),
      toggleLanguage: toggleLanguageMock,
    }),
}));

vi.mock('@disclosure-portal/stores/navigation', () => ({
  createNavItemsGroup: createNavItemsGroupMock,
}));

vi.mock('@disclosure-portal/utils/logout', () => ({
  logout: logoutMock,
}));

const childStubs = {
  NavItem: true,
  SubNavItem: true,
  DBreadcrumb: true,
  HelpDialog: true,
  Snowflake: true,
  DCloseButton: true,
  DSnackbar: true,
  'router-view': true,
  ProviderPrivacyDialog: {template: '<div><slot :showDialog="() => {}" /></div>'},
  ReleaseNotesDialog: {template: '<div><slot :showDialog="() => {}" /></div>'},
};

const createWrapper = () =>
  mountView(Dashboard, {
    childStubs,
    beforePiniaMount: (p) => {
      useUserStore(p).setSimpleProfileData({
        rights: {} as unknown as Rights,
        profile: {user: 'jane.doe'} as never,
        allowed: true,
      });
    },
  });

describe('Dashboard', () => {
  beforeEach(() => {
    themeToggleMock.mockReset();
    toggleLanguageMock.mockReset();
    createNavItemsGroupMock.mockReset();
    logoutMock.mockReset();
  });

  it('renders without throwing and shows the current user name', async () => {
    const {wrapper} = createWrapper();
    await flushPromises();

    expect(wrapper.exists()).toBe(true);
    expect(wrapper.text()).toContain('jane.doe');
  });

  it('toggles the app language and rebuilds the nav items', async () => {
    const {wrapper} = createWrapper();
    await flushPromises();

    await (wrapper.vm as unknown as {toggleLang: () => void}).toggleLang();

    expect(toggleLanguageMock).toHaveBeenCalledTimes(1);
    expect(createNavItemsGroupMock).toHaveBeenCalledTimes(1);
  });

  it('toggles the color theme via the account menu entry', async () => {
    const {wrapper} = createWrapper();
    await flushPromises();

    // The account-menu list lives inside a stubbed v-menu's default slot, so it's always present
    // in the DOM regardless of open state - find the "switch theme" entry and click it directly.
    const themeListItem = wrapper.findAll('.v-list-item').find((item) => item.text().includes('Dark Mode'));
    expect(themeListItem).toBeTruthy();
    await themeListItem!.trigger('click');

    expect(themeToggleMock).toHaveBeenCalledTimes(1);
  });

  it('clears the user store and logs out when logout is clicked', async () => {
    const {wrapper, pinia} = createWrapper();
    await flushPromises();

    (wrapper.vm as unknown as {logoutUser: () => void}).logoutUser();

    expect(logoutMock).toHaveBeenCalledTimes(1);
    expect(useUserStore(pinia).simpleProfileData.allowed).toBe(false);
  });
});
