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

const {sliderWidth, setupTableActionSlider, stopSlideInTimerAndSlideOut, startSlideInTimer} = useTableActionSlider();

if (props.variant === 'slider') {
  setupTableActionSlider(
    () => emit('slideOut', sliderWidth.value),
    () => emit('slideIn', sliderWidth.value),
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
            @clicked="emit(button.event)"></DIconButton>
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
          @clicked="emit(button.event)"></DIconButton>
      </div>
    </template>

    <!-- Normal Variant: All buttons displayed -->
    <template v-else-if="variant === 'slider'">
      <div class="flex justify-start pl-8">
        <v-btn
          plain
          size="small"
          variant="text"
          icon
          color="primary"
          class="size-10"
          @click.stop
          @mouseenter="stopSlideInTimerAndSlideOut"
          @mouseleave="startSlideInTimer">
          <v-icon>mdi-dots-horizontal</v-icon>
          <Tooltip location="bottom" text="Actions"></Tooltip>
        </v-btn>
        <div
          v-for="button in buttons"
          :key="button.icon"
          class="size-10"
          @mouseenter="stopSlideInTimerAndSlideOut"
          @mouseleave="startSlideInTimer">
          <DIconButton
            v-if="button.show ?? true"
            :icon="button.icon"
            :hint="button.hint"
            :color="button.color"
            :disabled="button.disabled"
            @clicked="emit(button.event)"></DIconButton>
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
          @clicked="emit(button.event)"></DIconButton>
      </div>

      <div v-if="remainingButtons.length > 0" class="size-10">
        <ExtraMenu>
          <div v-for="button in remainingButtons" :key="button.icon">
            <DIconButton
              :icon="button.icon"
              :hint="button.hint"
              :color="button.color"
              :disabled="button.disabled"
              @clicked="emit(button.event)"></DIconButton>
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
