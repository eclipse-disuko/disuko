<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {nextTick, ref, watch} from 'vue';
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

watch(model, (val) => {
  if (val !== '') {
    isExpanded.value = true;
  }
});

async function expand() {
  isExpanded.value = true;
  await nextTick();
  inputRef.value?.focus();
}

function collapse() {
  if (model.value === '') {
    isExpanded.value = false;
  }
}

function focus() {
  if (!isExpanded.value) {
    expand();
  } else {
    inputRef.value?.focus();
  }
}

function blur() {
  inputRef.value?.blur();
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    if (model.value === '') {
      isExpanded.value = false;
    }
  }
}

function onBlur() {
  if (model.value === '') {
    isExpanded.value = false;
  }
}

defineExpose({focus, blur, expand, collapse});
</script>

<template>
  <div class="d-search-field">
    <Transition name="d-search-expand">
      <v-text-field
        v-if="isExpanded"
        ref="inputRef"
        v-model="model"
        v-bind="$attrs"
        autocomplete="off"
        :width="500"
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
        {{ t('labelSearch') }}
      </v-btn>
    </Transition>
  </div>
</template>

<style scoped>
.d-search-field {
  display: inline-flex;
  align-items: center;
}

.d-search-expand-enter-active,
.d-search-expand-leave-active {
  transition:
    opacity 0.2s ease,
    max-width 0.25s ease;
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
  max-width: 500px;
}
</style>
