<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script lang="ts" setup>
import {useAttrs} from 'vue';
import {useI18n} from 'vue-i18n';

defineOptions({
  inheritAttrs: false,
});

type SelectValue = string | number | boolean;

interface Props {
  modelValue: SelectValue[];
  items: unknown[];
  label: string;
  itemTitle?: string;
  itemValue?: string;
  clearable?: boolean;
  othersLabel?: string;
  color?: string;
  listProps?: Record<string, unknown>;
}

const props = withDefaults(defineProps<Props>(), {
  itemTitle: 'text',
  itemValue: 'value',
  clearable: true,
  othersLabel: 'OTHERS',
  color: 'inputActiveBorderColor',
  listProps: () => ({class: 'striped-filter-dd py-0'}),
});

const emit = defineEmits<{
  'update:modelValue': [value: SelectValue[]];
}>();

const {t} = useI18n();
const attrs = useAttrs();

const multiSelectClass =
  'd-multi-select w-full min-w-0 ' +
  '[&_.v-field]:min-w-0 [&_.v-field__input]:flex-nowrap [&_.v-field__input]:overflow-hidden ' +
  '[&_.v-select__selection]:overflow-hidden ' +
  '[&_.v-select__selection-text]:overflow-hidden [&_.v-select__selection-text]:text-ellipsis [&_.v-select__selection-text]:whitespace-nowrap';

const maxSelectionTitleLength = 12;

const getSelectionTitle = (title: unknown): string => {
  const text = String(title ?? '');
  if (text.length <= maxSelectionTitleLength) {
    return text;
  }
  return `${text.slice(0, maxSelectionTitleLength)}...`;
};
</script>

<template>
  <v-select
    :class="multiSelectClass"
    :model-value="modelValue"
    :items="items"
    :label="label"
    :item-title="itemTitle"
    :item-value="itemValue"
    variant="outlined"
    density="compact"
    hide-details
    :clearable="clearable"
    :color="color"
    :list-props="listProps"
    multiple
    transition="scale-transition"
    persistent-clear
    v-bind="attrs"
    @update:modelValue="emit('update:modelValue', $event as SelectValue[])">
    <template v-slot:selection="{item, index}">
      <span v-if="index === 0" class="pFilterEntry inline-block max-w-full truncate align-middle">
        {{ getSelectionTitle(item.title) }}
      </span>
      <span v-if="index === 1" class="pAdditionalFilter whitespace-nowrap">
        (+{{ modelValue.length - 1 }} {{ t(othersLabel, modelValue.length) }})
      </span>
    </template>
  </v-select>
</template>
