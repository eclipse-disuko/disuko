<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import MarkdownIt from 'markdown-it';
import {onMounted, ref, watch} from 'vue';

const props = defineProps({
  text: {
    type: String,
    required: true,
  },
});

const md = new MarkdownIt({html: true});
const markdownText = ref('');

const defaultRender =
  md.renderer.rules.link_open ||
  function (tokens, idx, options, env, self) {
    return self.renderToken(tokens, idx, options);
  };

md.renderer.rules.link_open = function (tokens, idx, options, env, self) {
  const aIndex = tokens[idx].attrIndex('target');
  if (aIndex < 0) {
    tokens[idx].attrPush(['target', '_blank']);
  } else {
    tokens[idx].attrs![aIndex][1] = '_blank';
  }
  return defaultRender(tokens, idx, options, env, self);
};

onMounted(() => {
  markdownText.value = md.render(props.text);
});

watch(
  () => props.text,
  (newValue) => {
    if (newValue) {
      markdownText.value = md.render(newValue);
    }
  },
);
</script>

<template>
  <div>
    <div class="markdown" v-html="markdownText"></div>
    <slot />
  </div>
</template>
