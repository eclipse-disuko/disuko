<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
/**
 * Dumb step component, only shows a step button with header and error text.
 * Accepts props to display and emit index on click.
 * */
import {StepType} from '@disclosure-portal/model/NewWizard';
import {computed} from 'vue';
import {useI18n} from 'vue-i18n';
import {useLanguageStore} from '@shared/stores/language.store';
import {storeToRefs} from 'pinia';

export interface StepProps {
  step: StepType;
  currentIndex: number;
}

const props = defineProps<StepProps>();
const emit = defineEmits(['click']);

const languageStore = useLanguageStore();
const {appLanguage} = storeToRefs(languageStore);
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
  <div class="relative flex max-w-[110px] min-w-[85px] flex-col items-center space-y-1">
    <v-btn
      :class="btnClasses"
      :variant="!isActive ? 'outlined' : undefined"
      :disabled="!isClickable"
      flat
      @click="emit('click')">
      <template v-if="step.isCompleted"><v-icon size="18">mdi-check</v-icon></template>
      <template v-else>
        <span v-if="isActive" class="font-bold text-white">{{ step.index + 1 }}</span>
        <span v-else>{{ step.index + 1 }}</span>
      </template>
    </v-btn>
    <div
      class="text-body-2 w-full px-1 text-center break-words hyphens-auto"
      :class="{'opacity-50': !isActive}"
      :lang="appLanguage">
      {{ t(step.i18nKey) }}
    </div>
    <small
      class="block min-h-[14px] w-full text-center text-[11px] leading-[1.1] break-words text-yellow-700 dark:text-yellow-500">
      {{ step.errorText || ' ' }}
    </small>
  </div>
</template>
