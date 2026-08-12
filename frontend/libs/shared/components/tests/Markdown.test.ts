// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {config, mount} from '@vue/test-utils';
import {afterAll, beforeAll, describe, expect, it} from 'vitest';
import {nextTick} from 'vue';
import Markdown from '../Markdown.vue';

describe('Markdown', () => {
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

  it('renders the given markdown text as HTML', async () => {
    const wrapper = mount(Markdown, {props: {text: '# Hello world'}});
    await nextTick();

    expect(wrapper.find('.markdown h1').text()).toBe('Hello world');
  });

  it('renders bold and italic markdown syntax', async () => {
    const wrapper = mount(Markdown, {props: {text: '**bold** and *italic*'}});
    await nextTick();

    expect(wrapper.find('.markdown strong').text()).toBe('bold');
    expect(wrapper.find('.markdown em').text()).toBe('italic');
  });

  it('forces links to open in a new tab', async () => {
    const wrapper = mount(Markdown, {props: {text: '[link](https://example.com)'}});
    await nextTick();

    const link = wrapper.find('.markdown a');
    expect(link.attributes('href')).toBe('https://example.com');
    expect(link.attributes('target')).toBe('_blank');
  });

  it('re-renders the markdown when the text prop changes', async () => {
    const wrapper = mount(Markdown, {props: {text: '# First'}});
    await nextTick();
    expect(wrapper.find('.markdown h1').text()).toBe('First');

    await wrapper.setProps({text: '# Second'});

    expect(wrapper.find('.markdown h1').text()).toBe('Second');
  });

  it('keeps the previous rendered content when the text prop becomes empty', async () => {
    const wrapper = mount(Markdown, {props: {text: '# First'}});

    await wrapper.setProps({text: ''});

    expect(wrapper.find('.markdown h1').text()).toBe('First');
  });

  it('renders the default slot content alongside the markdown', () => {
    const wrapper = mount(Markdown, {
      props: {text: '# Hello'},
      slots: {default: '<div class="extra">extra content</div>'},
    });

    expect(wrapper.find('.extra').exists()).toBe(true);
    expect(wrapper.find('.extra').text()).toBe('extra content');
  });
});
