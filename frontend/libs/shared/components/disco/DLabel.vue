// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script lang="ts">
import {PropType, defineComponent} from 'vue';

export default defineComponent({
  props: {
    iconName: {
      type: String,
      required: false,
    },
    labelName: {
      type: String,
      required: false,
    },
    cssClazzes: {
      type: String,
      required: false,
    },
    parentProps: {
      type: Object as PropType<Record<string, any>>,
      required: false,
    },
    closable: {
      type: Boolean,
      required: false,
      default: false,
    },
  },
});
</script>

<template>
  <v-chip
    :class="cssClazzes && cssClazzes.length > 0 ? cssClazzes : 'px-2 py-2 mr-1'"
    variant="outlined"
    size="small"
    color="labels"
    label
    :prepend-icon="iconName"
    v-bind="parentProps"
    :closable="closable"
  >
    {{ labelName }}
  </v-chip>
</template>
