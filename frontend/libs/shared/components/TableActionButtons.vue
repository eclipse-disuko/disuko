// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import {useTableActionSlider} from '@shared/composables/useTableActionSlider';
import {computed} from 'vue';

interface Button {
  icon: string;
  event: string;
  hint?: string;
  disabled?: boolean;
  show?: boolean;
  color?: string;
}

export interface TableActionButtonsProps {
  buttons: Button[];
  variant?: 'normal' | 'minimal' | 'compact' | 'slider';
}

const props = withDefaults(defineProps<TableActionButtonsProps>(), {
  variant: 'normal',
});

const emit = defineEmits<{
  slideOut: [value: number];
  slideIn: [value: number];
  [key: string]: [value?: number];
}>();

const shownButtons = computed(() => props.buttons.filter((button) => button.show ?? true));
const outsideButtons = computed(() => shownButtons.value.slice(0, 1));
const remainingButtons = computed(() => shownButtons.value.slice(1));

const {sliderWidth, expandedMaxWidth, setupTableActionSlider, stopSlideInTimerAndSlideOut, startSlideInTimer} =
  useTableActionSlider();

if (props.variant === 'slider') {
  setupTableActionSlider(
    props.buttons.length,
    () => emit('slideOut', sliderWidth.value),
    () => emit('slideIn', sliderWidth.value),
  );
}
</script>

<template>
  <div
    class="h-[100%] flex items-center"
    :class="{'justify-center': variant !== 'slider', 'justify-start': variant === 'slider'}">
    <!-- Minimal Variant: All buttons in an extra menu -->
    <template v-if="variant === 'minimal'">
      <ExtraMenu>
        <div v-for="button in shownButtons" :key="button.icon">
          <DIconButton
            :icon="button.icon"
            :hint="button.hint"
            :color="button.color"
            :disabled="button.disabled"
            @clicked="emit(button.event)" />
        </div>
      </ExtraMenu>
    </template>

    <!-- Normal Variant: All buttons displayed -->
    <template v-else-if="variant === 'normal'">
      <div v-for="button in buttons" :key="button.icon" class="size-10">
        <DIconButton
          v-if="button.show ?? true"
          :icon="button.icon"
          :hint="button.hint"
          :color="button.color"
          :disabled="button.disabled"
          @clicked="emit(button.event)" />
      </div>
    </template>

    <!-- Normal Variant: All buttons displayed -->
    <template v-else-if="variant === 'slider'">
      <div
        class="h-full flex items-center justify-start pr-5"
        @click.stop
        @mouseenter="stopSlideInTimerAndSlideOut"
        @mouseleave="startSlideInTimer">
        <div class="flex h-[40px]" :style="{width: expandedMaxWidth + 'px'}">
          <template v-for="button in buttons" :key="button.icon">
            <div
              v-if="button?.show ?? true"
              class="d-inline size-10"
              @click.stop="!button?.disabled ? emit(button.event) : null">
              <v-btn
                plain
                size="small"
                variant="text"
                density="default"
                :icon="button.icon"
                :color="button.color || 'primary'"
                :disabled="Boolean(button?.disabled) || false" />
              <Tooltip v-if="button.hint && !button?.disabled" location="bottom" :text="button.hint" />
            </div>
          </template>
        </div>
      </div>
    </template>

    <!-- Compact Variant: When there are 2 buttons, show them without menu -->
    <template v-else-if="variant === 'compact' && shownButtons.length <= 2">
      <div v-for="button in shownButtons" :key="button.icon" class="size-10">
        <DIconButton
          :icon="button.icon"
          :hint="button.hint"
          :color="button.color"
          :disabled="button.disabled"
          @clicked="emit(button.event)"></DIconButton>
      </div>
    </template>

    <!-- Compact Variant: First button displayed, rest in extra menu -->
    <template v-else-if="variant === 'compact' && shownButtons.length > 2">
      <div v-for="button in outsideButtons" :key="button.icon" class="size-10">
        <DIconButton
          :icon="button.icon"
          :hint="button.hint"
          :color="button.color"
          :disabled="button.disabled"
          @clicked="emit(button.event)" />
      </div>

      <div v-if="remainingButtons.length > 0" class="size-10">
        <ExtraMenu>
          <div v-for="button in remainingButtons" :key="button.icon">
            <DIconButton
              :icon="button.icon"
              :hint="button.hint"
              :color="button.color"
              :disabled="button.disabled"
              @clicked="emit(button.event)" />
          </div>
        </ExtraMenu>
      </div>
    </template>
  </div>
</template>

<style lang="scss">
.action-slider-table > .v-table > .v-table__wrapper > table {
  > thead > tr > th:nth-child(2) {
    padding-right: 0 !important;
    padding-left: 0 !important;
    transition: width ease-in-out 0.2s;
  }
  > tbody > tr > td:first-child {
    padding-right: 0 !important;
    padding-left: 0 !important;
  }
  > tbody > tr > td:nth-child(2) {
    padding-right: 0 !important;
    padding-left: 0 !important;
  }
}
</style>
