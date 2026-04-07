<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {useTableActionSlider} from '@shared/composables/useTableActionSlider';
import {computed} from 'vue';
import {useI18n} from 'vue-i18n';

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
  slideToggle: [value: number];
  [key: string]: [value?: number];
}>();

const {t} = useI18n();

const shownButtons = computed(() => props.buttons.filter((button) => button.show ?? true));
const outsideButtons = computed(() => shownButtons.value.slice(0, 1));
const remainingButtons = computed(() => shownButtons.value.slice(1));

const {sliderWidth, baseWidth, setupTableActionSlider, stopSlideInTimerAndSlideOut, startSlideInTimer} = useTableActionSlider();

if (props.variant === 'slider') {
  setupTableActionSlider(
    () => emit('slideToggle', sliderWidth.value),
    props.buttons.length,
  );
}
</script>

<template>
  <div
    class="flex items-center"
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
        class="h-[100%] flex justify-start align-center pl-8 pr-5"
        @click.stop
        @mouseenter="stopSlideInTimerAndSlideOut"
        @mouseleave="startSlideInTimer">
        <v-btn
          v-if="shownButtons.length >= 2"
          plain
          size="small"
          variant="text"
          icon
          color="primary"
          class="size-10"
          @click.stop>
          <v-icon>mdi-dots-horizontal</v-icon>
          <Tooltip location="bottom" :text="t('OPEN_ACTIONS')" />
        </v-btn>
        <div
          v-else-if="(buttons[0]?.show ?? true) && !(buttons[0]?.disabled ?? false)"
          class="d-inline size-10"
          @click.stop="emit(buttons[0].event)">
          <v-btn
            plain
            size="small"
            variant="text"
            density="default"
            :icon="buttons[0].icon"
            :color="buttons[0].color || 'primary'" />
          <Tooltip v-if="buttons[0].hint" location="bottom" :text="buttons[0].hint" />
        </div>
        <template v-if="shownButtons.length >= 2">
          <template v-for="button in buttons" :key="button.icon">
            <div :style="{opacity: sliderWidth !== baseWidth ? 1 : 0}" class="transition-[opacity] ease-in-out duration-200">
              <div
                v-if="(button?.show ?? true) && !(button?.disabled ?? false)"
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
            </div>
          </template>
        </template>
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
  > thead > tr > th:first-child {
    transition: width ease-in-out 0.2s;
  }
  > tbody > tr > td:first-child {
    padding-right: 0 !important;
  }
}
</style>
