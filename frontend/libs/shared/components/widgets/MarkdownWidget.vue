<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import MarkdownIt from 'markdown-it';
import {onMounted, ref, watch} from 'vue';

interface Props {
  id: string;
  text: string;
}
const props = defineProps<Props>();
const markdownText = ref('');
const md = new MarkdownIt({html: true});

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

const renderMarkdown = (text: string) => {
  markdownText.value = md.render(text);
};

onMounted(() => {
  renderMarkdown(props.text);
});

watch(
  () => props.text,
  (newValue) => {
    if (newValue) {
      renderMarkdown(newValue);
    }
  },
);
</script>

<template>
  <div class="markdown" v-html="markdownText"></div>
  <slot name="default"> </slot>
</template>
