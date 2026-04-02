<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<template>
  <v-btn
    :variant="variant"
    :disabled="disabled"
    :size="size"
    :block="block"
    :color="color ? color : 'primary'"
    :loading="loading"
    :class="isDialogButton ? '' : 'text-none'"
    class="pointer-events-auto"
    @click="$emit('clicked')">
    <div class="d-inline discoActionBtnHover">
      <v-icon :color="color ? color : 'primary'" v-if="icon">{{ icon }}</v-icon>
      <span
        class="font-weight-bold px-1"
        v-if="text"
        :style="color && isDisabledTextColor ? 'color:' + color + ' !important;' : ''">
        {{ text }}
      </span>
    </div>

    <tooltip v-if="hint" :text="hint"></tooltip>
  </v-btn>
</template>

<script setup lang="ts">
interface Props {
  icon?: string;
  text?: string;
  isDisabledTextColor?: boolean;
  isDialogButton?: boolean;
  block?: boolean;
  hint?: string;
  variant?: 'flat' | 'text' | 'elevated' | 'tonal' | 'outlined' | 'plain';
  color?: string;
  size?: string;
  disabled?: boolean;
  loading?: boolean;
}

withDefaults(defineProps<Props>(), {
  isDisabledTextColor: true,
  isDialogButton: false,
  variant: 'tonal',
  color: 'primary',
});
</script>
