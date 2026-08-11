// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {Group, Rights} from '@shared/user/models/Rights';
import {useUserStore} from '@shared/user/stores/user.store';
import {useAppStore} from '@disclosure-portal/stores/app';
import {useNewsboxStore} from '@disclosure-portal/stores/newsbox.store';
import {useWizardStore} from '@disclosure-portal/stores/wizard.store';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {nextTick} from 'vue';
import {describe, expect, it, vi} from 'vitest';
import Home from '../Home.vue';

const {getDashboardCountsMock} = vi.hoisted(() => ({
  getDashboardCountsMock: vi.fn(),
}));

vi.mock('@shared/user/services/profile.service', () => ({
  default: {
    getDashboardCounts: getDashboardCountsMock,
  },
}));

const childStubs = {
  Stack: {template: '<div class="stack"><slot /></div>'},
  HomeTile: {
    template: '<button type="button" class="home-tile" @click="$emit(\'click\')">{{ title }}</button>',
    props: ['type', 'title', 'icon', 'description', 'url', 'showBadge', 'badgeContent'],
    emits: ['click'],
  },
  NewsboxDialog: true,
};

const createWrapper = (rightsGroups: Group[] = [Group.UserApplicationAdmin]) =>
  mountView(Home, {
    childStubs,
    beforePiniaMount: (pinia) => {
      useUserStore(pinia).setSimpleProfileData({
        rights: {groups: rightsGroups, allowProject: {create: true}} as unknown as Rights,
        profile: {forename: 'Jane'} as never,
        allowed: true,
      });
    },
  });

describe('Home', () => {
  it('sets the dashboard breadcrumb and loads dashboard counts on mount', async () => {
    getDashboardCountsMock.mockResolvedValue({activeJobCount: 3, hasNewNewsboxItem: true});

    const {pinia} = createWrapper();
    await nextTick();
    await nextTick();

    const breadcrumbs = useBreadcrumbsStore(pinia);
    expect(breadcrumbs.setCurrentBreadcrumbs).toHaveBeenCalledWith([breadcrumbs.dashboard]);

    expect(getDashboardCountsMock).toHaveBeenCalledTimes(1);
    const appStore = useAppStore(pinia);
    expect(appStore.updateTileCounts).toHaveBeenCalledWith({activeJobCount: 3, hasNewNewsboxItem: true});
    const newsboxStore = useNewsboxStore(pinia);
    expect(newsboxStore.hasNewNewsboxItem).toBe(true);
  });

  it('opens the new-project wizard when the create-project tile is clicked', async () => {
    getDashboardCountsMock.mockResolvedValue({activeJobCount: 0, hasNewNewsboxItem: false});

    const {wrapper, pinia} = createWrapper();
    await nextTick();
    await nextTick();

    const wizardStore = useWizardStore(pinia);
    vi.mocked(wizardStore.openWizard).mockImplementation(() => Promise.resolve());

    await wrapper.find('.home-tile').trigger('click');

    expect(wizardStore.openWizard).toHaveBeenCalledTimes(1);
  });

  it('opens the newsbox when the "did you know" tile is clicked', async () => {
    getDashboardCountsMock.mockResolvedValue({activeJobCount: 0, hasNewNewsboxItem: false});

    const {wrapper, pinia} = createWrapper();
    await nextTick();
    await nextTick();

    const tiles = wrapper.findAll('.home-tile');
    await tiles[tiles.length - 1].trigger('click');

    const newsboxStore = useNewsboxStore(pinia);
    expect(newsboxStore.showNewsbox).toBe(true);
  });
});
