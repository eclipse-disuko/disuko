<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<template>
  <v-tooltip
    v-if="hint"
    :open-delay="TOOLTIP_OPEN_DELAY_IN_MS"
    :text="hint"
    location="bottom"
    content-class="dpTooltip">
    <template v-slot:activator="{props}">
      <v-list-item v-bind="props" @click="$emit('click')" class="px-4">
        <v-icon style="width: 30px" :color="color ? color : 'primary'" v-if="icon">{{ icon }}</v-icon>
        <span class="d-subtitle px-3">{{ text }}</span>
      </v-list-item>
    </template>
  </v-tooltip>
</template>
<script lang="ts">
import {TOOLTIP_OPEN_DELAY_IN_MS} from '@shared/utils/constant';
import {defineComponent} from 'vue';

export default defineComponent({
  emits: ['click'],
  props: {
    icon: {
      type: String,
      required: false,
    },
    text: {
      type: String,
      required: false,
    },
    isDisabledTextColor: {
      type: Boolean,
      default: true,
      required: false,
    },
    block: {
      type: Boolean,
      required: false,
    },
    hint: {
      type: String,
      required: false,
    },
    color: {
      type: String,
      required: false,
    },
    disabled: {
      type: Boolean,
      required: false,
    },
    loading: {
      type: Boolean,
      required: false,
    },
  },
  setup() {
    return {
      TOOLTIP_OPEN_DELAY_IN_MS,
    };
  },
});
</script>
