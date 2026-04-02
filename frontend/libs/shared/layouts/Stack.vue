<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {computed, useAttrs} from 'vue';

/**
 * Stack Component is used to display a list of items in a flexible layout,
 * whether vertical or horizontal. By default, it stacks items vertically and has gap-3.
 * ```[item] [item] [item]```
 * or
 * ```
 * [item]
 * [item]
 * [item]
 * ```
 *
 * @example Basic usage
 * <Stack>
 *   <div>Item 1</div>
 *   <div>Item 2</div>
 * </Stack>
 *
 * @example Row with different alignment and custom gap
 * <Stack direction="row" justify="between" align="center" class="gap-8">
 *   <div>Left aligned</div>
 *   <div>Right aligned</div>
 * </Stack>
 *
 * @example With inherited classes
 * <Stack class="bg-gray-100 p-4 rounded-md">
 *   <div>Content in a gray box with padding</div>
 * </Stack>
 */
type Props = {
  direction?: 'row' | 'column' | 'row-reverse' | 'column-reverse';
  justify?: 'start' | 'end' | 'center' | 'between' | 'around' | 'evenly';
  align?: 'start' | 'end' | 'center' | 'baseline' | 'stretch';
  wrap?: boolean;
  tag?: string;
};
const props = withDefaults(defineProps<Props>(), {
  direction: 'column',
  wrap: false,
  tag: 'div',
});

const attrs = useAttrs();

// Static mappings for Tailwind classes to ensure they are discoverable
const directionClasses = {
  row: 'flex-row',
  column: 'flex-col',
  'row-reverse': 'flex-row-reverse',
  'column-reverse': 'flex-col-reverse',
};

const justifyClasses = {
  start: 'justify-start',
  end: 'justify-end',
  center: 'justify-center',
  between: 'justify-between',
  around: 'justify-around',
  evenly: 'justify-evenly',
};

const alignClasses = {
  start: 'items-start',
  end: 'items-end',
  center: 'items-center',
  baseline: 'items-baseline',
  stretch: 'items-stretch',
};

const stackClasses = computed(() => {
  const classes = ['flex', directionClasses[props.direction]];

  if (props.justify) {
    classes.push(justifyClasses[props.justify]);
  }
  if (props.align) {
    classes.push(alignClasses[props.align]);
  }

  if (props.wrap) {
    classes.push('flex-wrap');
  }

  const externalClasses = attrs.class as string | undefined;
  const hasExternalGap = externalClasses && /(^|\s)gap-\S+/.test(externalClasses);
  if (!hasExternalGap) {
    classes.push('gap-3');
  }

  return classes;
});
</script>

<template>
  <component :is="tag" :class="[stackClasses]">
    <slot></slot>
  </component>
</template>
