<script setup lang="ts">
/**
 * Dumb step component, only shows a step button with header and error text.
 * Accepts props to display and emit index on click.
 * */
import {StepType} from '@disclosure-portal/model/NewWizard';
import {useAppStore} from '@disclosure-portal/stores/app';
import {computed} from 'vue';
import {useI18n} from 'vue-i18n';

export interface StepProps {
  step: StepType;
  currentIndex: number;
}

const props = defineProps<StepProps>();
const emit = defineEmits(['click']);

const {appLanguage} = useAppStore();
const {t} = useI18n();

const isActive = computed(() => props.currentIndex === props.step.index);

const btnClasses = computed(() => {
  const base =
    'aspect-square size-10 !min-w-0 rounded-full flex items-center justify-center transition-colors duration-200';

  // Step is active
  if (isActive.value) {
    return `${base} bg-[rgb(var(--v-theme-primary))] text-white shadow-[0_0_0_2px_rgba(var(--v-theme-primary),0.25)]`;
  }

  // Step is completed
  if (props.step.isCompleted) {
    return `${base} border border-[rgba(var(--v-theme-primary),0.6)] text-[rgb(var(--v-theme-primary))]`;
  }

  // Step has been visited but not completed
  if (props.step.seen && !props.step.isCompleted) {
    return `${base} border border-gray-700 text-gray-700 dark:border-gray-100 dark:text-gray-100`;
  }

  // not seen and not completed - inactive, default state
  return `${base} bg-[rgba(var(--v-theme-secondary),0,8)] border border-[rgba(var(--v-theme-secondary),0.6)] pointer-events-none opacity-60`;
});

const isClickable = computed(() => {
  return isActive.value || props.step.isCompleted || props.step.seen;
});
</script>

<template>
  <div class="flex flex-col items-center space-y-1 relative min-w-[85px] max-w-[110px]">
    <v-btn
      :class="btnClasses"
      :variant="!isActive ? 'outlined' : undefined"
      :disabled="!isClickable"
      flat
      @click="emit('click')">
      <template v-if="props.step.isCompleted"><v-icon size="18">mdi-check</v-icon></template>
      <template v-else>
        <span v-if="isActive" class="text-white font-bold">{{ props.step.index + 1 }}</span>
        <span v-else>{{ props.step.index + 1 }}</span>
      </template>
    </v-btn>
    <div
      class="text-body-2 break-words hyphens-auto w-full text-center px-1"
      :class="{'opacity-50': !isActive}"
      :lang="appLanguage">
      {{ t(props.step.i18nKey) }}
    </div>
    <small
      class="block text-[11px] leading-[1.1] min-h-[14px] text-center break-words w-full text-yellow-700 dark:text-yellow-500">
      {{ props.step.errorText || ' ' }}
    </small>
  </div>
</template>
