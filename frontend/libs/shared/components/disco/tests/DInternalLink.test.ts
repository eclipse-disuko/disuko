// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {config, mount} from '@vue/test-utils';
import {afterAll, beforeAll, describe, expect, it} from 'vitest';
import DInternalLink from '../DInternalLink.vue';

describe('DInternalLink', () => {
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

  const createWrapper = (props: Record<string, unknown>) => {
    return mount(DInternalLink, {props});
  };

  it('renders without throwing', () => {
    const wrapper = createWrapper({url: '/foo', text: 'Foo'});

    expect(wrapper.exists()).toBe(true);
  });

  it('renders an anchor with the given href and text when url is set', () => {
    const wrapper = createWrapper({url: '/projects/1', text: 'Project 1'});
    const link = wrapper.find('a');

    expect(link.exists()).toBe(true);
    expect(link.attributes('href')).toBe('/projects/1');
    expect(link.text()).toBe('Project 1');
  });

  it('renders no anchor when url is empty', () => {
    const wrapper = createWrapper({url: '', text: 'Project 1'});

    expect(wrapper.find('a').exists()).toBe(false);
  });
});
