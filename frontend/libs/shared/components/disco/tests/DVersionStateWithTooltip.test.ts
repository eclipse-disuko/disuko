// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {VersionSlim, VersionSlimDto} from '@disclosure-portal/model/VersionDetails';
import {config, mount} from '@vue/test-utils';
import dayjs from 'dayjs';
import {afterAll, beforeAll, describe, expect, it} from 'vitest';
import DVersionStateWithTooltip from '../DVersionStateWithTooltip.vue';

const VTooltipStub = {
  props: ['location', 'maxWidth', 'contentClass'],
  template:
    '<div class="v-tooltip"><slot name="activator" :props="{}" /><div class="tooltip-content"><slot /></div></div>',
};

const VIconStub = {
  props: ['color', 'size'],
  template: '<i class="v-icon" :color="color"><slot /></i>',
};

describe('DVersionStateWithTooltip', () => {
  const originalTMock = config.global.mocks?.$t;

  beforeAll(() => {
    if (config.global.mocks && '$t' in config.global.mocks) {
      delete config.global.mocks.$t;
    }
  });

  afterAll(() => {
    if (originalTMock) {
      config.global.mocks = {
        ...config.global.mocks,
        $t: originalTMock,
      };
    }
  });

  const makeVersion = (overrides: Partial<VersionSlimDto> = {}) =>
    new VersionSlim({...new VersionSlimDto(), status: 'acceptable', ...overrides});

  const createWrapper = (props: Record<string, unknown>) => {
    return mount(DVersionStateWithTooltip, {
      props,
      global: {stubs: {'v-tooltip': VTooltipStub, 'v-icon': VIconStub}},
    });
  };

  it('renders the icon with the color matching the version status', () => {
    const wrapper = createWrapper({version: makeVersion({status: 'acceptable'})});

    const icon = wrapper.find('.v-icon');
    expect(icon.attributes('color')).toBe('versionApproved');
    expect(icon.text()).toBe('mdi-comment-check');
  });

  it('renders the translated overall review status in the tooltip', () => {
    const wrapper = createWrapper({version: makeVersion({status: 'not_acceptable'})});

    expect(wrapper.find('.tooltip-content').text()).toContain('Problem');
  });

  it('does not render recent overall review details when there are none', () => {
    const wrapper = createWrapper({version: makeVersion({overallReviews: []})});

    expect(wrapper.find('.tooltip-content').text()).not.toContain('For SBOM:');
  });

  it('shows the most recently updated overall review comment and sbom name', () => {
    const wrapper = createWrapper({
      version: makeVersion({
        overallReviews: [
          {
            created: '',
            updated: '2024-01-01T00:00:00Z',
            state: 'ACCEPTABLE' as never,
            comment: 'Older comment',
            sbomId: '',
            sbomName: 'old.spdx',
            sbomUploaded: '2024-01-01T00:00:00Z',
            creator: 'jdoe',
            creatorFullName: 'Jane Doe',
          },
          {
            created: '',
            updated: '2024-06-01T00:00:00Z',
            state: 'ACCEPTABLE' as never,
            comment: 'Latest comment',
            sbomId: '',
            sbomName: 'new.spdx',
            sbomUploaded: dayjs().format(),
            creator: 'asmith',
            creatorFullName: 'Ann Smith',
          },
        ],
      }),
    });

    const text = wrapper.find('.tooltip-content').text();
    expect(text).toContain('Latest comment');
    expect(text).toContain('new.spdx');
    expect(text).not.toContain('Older comment');
  });

  it('shows the sbom outdated warning when the most recent sbom is older than 30 days', () => {
    const wrapper = createWrapper({
      version: makeVersion({
        overallReviews: [
          {
            created: '',
            updated: '2024-01-01T00:00:00Z',
            state: 'ACCEPTABLE' as never,
            comment: 'Stale review',
            sbomId: '',
            sbomName: 'old.spdx',
            sbomUploaded: dayjs().subtract(40, 'day').format(),
            creator: 'jdoe',
            creatorFullName: 'Jane Doe',
          },
        ],
      }),
    });

    expect(wrapper.find('.tooltip-content').text()).toContain('SBOM is older than 30 days and might be outdated');
  });

  it('does not show the sbom outdated warning when the sbom is recent', () => {
    const wrapper = createWrapper({
      version: makeVersion({
        overallReviews: [
          {
            created: '',
            updated: '2024-01-01T00:00:00Z',
            state: 'ACCEPTABLE' as never,
            comment: 'Fresh review',
            sbomId: '',
            sbomName: 'new.spdx',
            sbomUploaded: dayjs().format(),
            creator: 'jdoe',
            creatorFullName: 'Jane Doe',
          },
        ],
      }),
    });

    expect(wrapper.find('.tooltip-content').text()).not.toContain('SBOM is older than 30 days and might be outdated');
  });

  it('applies the ml-7 class to the icon when isGroup is true', () => {
    const wrapper = createWrapper({version: makeVersion(), isGroup: true});

    expect(wrapper.find('.v-icon').classes()).toContain('ml-7');
  });
});
