// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<template>
  <Markdown class="markdown" html :source="markdownText" />
</template>

<script setup lang="ts">
import MarkdownIt from 'markdown-it';
import {onMounted, ref} from 'vue';
import Markdown from 'vue3-markdown-it';

const props = defineProps({
  text: {
    type: String,
    required: true,
  },
});

const md = new MarkdownIt({html: true});
const markdownText = ref('');

onMounted(() => {
  // add _target="blank" to all links
  // see https://github.com/markdown-it/markdown-it/blob/master/docs/architecture.md#renderer
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
      // @ts-expect-error
      tokens[idx].attrs[aIndex][1] = '_blank';
    }
    return defaultRender(tokens, idx, options, env, self);
  };

  markdownText.value = md.render(props.text);
});
</script>
