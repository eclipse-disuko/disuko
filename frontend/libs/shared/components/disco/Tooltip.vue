<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {TOOLTIP_OPEN_DELAY_IN_MS} from '@shared/utils/constant';
import {computed} from 'vue';
import DExternalLink from './DExternalLink.vue';

interface Props {
  activator?: string;
  location?: 'top' | 'bottom' | 'left' | 'right';
  disabled?: boolean;
  text?: string;
  contentClass?: string;
  asParent?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  activator: 'parent',
  location: 'bottom',
  disabled: undefined,
  text: undefined,
  contentClass: 'dpTooltip',
  asParent: false,
});

interface TextSegment {
  type: 'text' | 'url';
  content: string;
}

const textSegments = computed<TextSegment[]>(() => {
  if (!props.text) return [];

  const urlRegex = /(https?:\/\/[^\s]+)/g;
  const segments: TextSegment[] = [];
  let lastIndex = 0;
  let match;

  while ((match = urlRegex.exec(props.text)) !== null) {
    // Add text before the URL
    if (match.index > lastIndex) {
      segments.push({
        type: 'text',
        content: props.text.substring(lastIndex, match.index),
      });
    }

    // Add the URL
    segments.push({
      type: 'url',
      content: match[0],
    });

    lastIndex = match.index + match[0].length;
  }

  // Add remaining text after the last URL
  if (lastIndex < props.text.length) {
    segments.push({
      type: 'text',
      content: props.text.substring(lastIndex),
    });
  }

  // If no URLs were found, return the whole text as a single segment
  if (segments.length === 0) {
    segments.push({
      type: 'text',
      content: props.text,
    });
  }

  return segments;
});

const hasLinks = computed(() => {
  return textSegments.value.some((segment) => segment.type === 'url');
});
</script>

<template>
  <v-tooltip
    v-if="asParent"
    :location="location"
    :disabled="disabled"
    :content-class="contentClass"
    :open-delay="TOOLTIP_OPEN_DELAY_IN_MS"
    :close-delay="hasLinks ? 500 : undefined"
    :interactive="hasLinks"
    v-bind="$attrs">
    <template #activator="{props: tooltipProps}">
      <span v-bind="tooltipProps">
        <slot></slot>
      </span>
    </template>
    <div>
      <span v-if="hasLinks">
        <template v-for="(segment, index) in textSegments" :key="index">
          <DExternalLink
            v-if="segment.type === 'url'"
            :url="segment.content"
            :text="segment.content"
            :tooltip="false" />
          <span v-else>{{ segment.content }}</span>
        </template>
      </span>
      <span v-else-if="text">{{ text }}</span>
    </div>
  </v-tooltip>

  <v-tooltip
    v-else
    :activator="activator"
    :location="location"
    :disabled="disabled"
    :text="hasLinks ? undefined : text"
    :content-class="contentClass"
    :open-delay="TOOLTIP_OPEN_DELAY_IN_MS"
    v-bind="$attrs">
    <template v-if="text">
      <span v-if="hasLinks">
        <template v-for="(segment, index) in textSegments" :key="index">
          <DExternalLink
            v-if="segment.type === 'url'"
            :url="segment.content"
            :text="segment.content"
            :tooltip="false" />
          <span v-else>{{ segment.content }}</span>
        </template>
      </span>
      <span v-else>{{ text }}</span>
    </template>
    <slot v-else></slot>
  </v-tooltip>
</template>
