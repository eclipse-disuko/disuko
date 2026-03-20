<script setup lang="ts">
import ExtraMenu from '@shared/components/disco/ExtraMenu.vue';
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
  variant?: 'normal' | 'minimal' | 'compact';
}
const props = withDefaults(defineProps<TableActionButtonsProps>(), {
  variant: 'normal',
});
const emit = defineEmits<{
  [key: string]: [];
}>();

const shownButtons = computed(() => props.buttons.filter((button) => button.show ?? true));
const outsideButtons = computed(() => shownButtons.value.slice(0, 1));
const remainingButtons = computed(() => shownButtons.value.slice(1));
</script>

<template>
  <div class="flex justify-center items-center">
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
