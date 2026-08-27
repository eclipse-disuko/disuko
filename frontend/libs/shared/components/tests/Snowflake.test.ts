// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {mount} from '@vue/test-utils';
import {describe, expect, it} from 'vitest';
import Snowflake from '../Snowflake.vue';

describe('Snowflake', () => {
  it('renders without throwing', () => {
    const wrapper = mount(Snowflake);

    expect(wrapper.find('.snowflake-container').exists()).toBe(true);
  });

  it('renders 15 snowflake elements', () => {
    const wrapper = mount(Snowflake);

    expect(wrapper.findAll('.snowflake')).toHaveLength(15);
  });

  it('assigns increasing delay and horizontal offset per snowflake', () => {
    const wrapper = mount(Snowflake);

    const snowflakes = wrapper.findAll('.snowflake');
    expect(snowflakes[0].attributes('style')).toContain('--delay: 0.5s');
    expect(snowflakes[0].attributes('style')).toContain('--x: 6.5%');
    expect(snowflakes[14].attributes('style')).toContain('--delay: 7.5s');
    expect(snowflakes[14].attributes('style')).toContain('--x: 97.5%');
  });
});
