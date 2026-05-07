<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {computed, nextTick, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';

defineOptions({inheritAttrs: false});

interface Props {
  disabled?: boolean;
}
defineProps<Props>();
const model = defineModel<string>({default: ''});

const {t} = useI18n();

const isExpanded = ref(model.value !== '');
const inputRef = ref<HTMLInputElement | null>(null);

const isFieldEmpty = computed(() => {
  return !model.value || model.value.trim() === '';
});

watch(model, (val) => {
  if (val !== '') {
    isExpanded.value = true;
  }
});

async function expand() {
  isExpanded.value = true;
  await nextTick();
  setTimeout(() => {
    inputRef.value?.focus();
  }, 150);
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    if (isFieldEmpty.value) {
      isExpanded.value = false;
    }
  }
}

function onBlur() {
  if (isFieldEmpty.value) {
    isExpanded.value = false;
  }
}
</script>

<template>
  <div class="align-center inline-flex">
    <Transition name="d-search-expand" mode="out-in">
      <v-text-field
        v-if="isExpanded"
        ref="inputRef"
        v-model="model"
        v-bind="$attrs"
        autocomplete="off"
        :width="400"
        append-inner-icon="mdi-magnify"
        variant="outlined"
        density="compact"
        :label="t('labelSearch')"
        :disabled="disabled"
        single-line
        hide-details
        clearable
        @keydown="onKeydown"
        @blur="onBlur" />
      <v-btn
        v-else
        variant="tonal"
        color="primary"
        :disabled="disabled"
        prepend-icon="mdi-magnify"
        class="text-none h-10"
        @click="expand">
        <span class="font-bold">{{ t('labelSearch') }}</span>
      </v-btn>
    </Transition>
  </div>
</template>

<style scoped>
.d-search-expand-enter-active {
  transition:
    opacity 0.1s ease-in,
    max-width 0.15s ease-in;
  overflow: hidden;
}

.d-search-expand-leave-active {
  transition:
    opacity 0.1s ease-out,
    max-width 0.1s ease-out;
  overflow: hidden;
}

.d-search-expand-enter-from,
.d-search-expand-leave-to {
  opacity: 0;
  max-width: 0;
}

.d-search-expand-enter-to,
.d-search-expand-leave-from {
  opacity: 1;
  max-width: 400px;
}
</style>
