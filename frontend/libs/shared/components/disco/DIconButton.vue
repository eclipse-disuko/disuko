<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script lang="ts" setup>
import {Density} from 'vuetify/lib/composables/density';

interface Props {
  icon: string;
  hint?: string;
  color?: string;
  disabled?: boolean;
  density?: Density;
  location?: string;
  parentProps?: Record<string, any>;
}

withDefaults(defineProps<Props>(), {
  color: 'primary',
  disabled: false,
  density: 'default',
  location: 'bottom',
});

const emit = defineEmits(['clicked']);
</script>

<template>
  <div class="d-inline">
    <v-btn
      v-bind="parentProps"
      plain
      size="small"
      variant="text"
      :density="density"
      :icon="icon"
      :color="color"
      :disabled="disabled"
      @click.stop="emit('clicked')"></v-btn>
    <Tooltip v-if="hint" :location="location" :text="hint"></Tooltip>
  </div>
</template>
