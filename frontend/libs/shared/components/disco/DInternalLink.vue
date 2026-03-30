// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script lang="ts">
import Icons from '@disclosure-portal/constants/icons';
import {releaseKeys} from '@disclosure-portal/keyState';
import {defineComponent} from 'vue';

export default defineComponent({
  name: 'DInternalLink',
  props: {
    url: {
      type: String,
      required: true,
    },
    text: {
      type: String,
      required: true,
    },
  },
  setup(props) {
    const icons = Icons;

    const openURL = () => {
      if (props.url) {
        window.open(props.url, '_blank');
        releaseKeys();
      }
    };

    return {
      icons,
      openURL,
      releaseKeys,
    };
  },
});
</script>

<template>
  <span>
    <a color="primary" class="d-subtitle-2" :href="url" target="_blank" v-if="url" @click="releaseKeys">
      {{ text }}
    </a>
  </span>
</template>
