<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import Tooltip from '@shared/components/disco/Tooltip.vue';
import {useForm} from '@shared/composables/useForm';
import {computed} from 'vue';

const props = withDefaults(
  defineProps<{
    modelValue: string | string[];
    readonly?: boolean;
    label: string;
    required?: boolean;
    items: {name: string; value: any}[] | string[];
    multiple?: boolean;
    simpleList?: boolean;
    help?: string;
  }>(),
  {
    readonly: false,
    required: false,
    multiple: false,
    simpleList: false,
  },
);

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void;
}>();

const {rules} = useForm();
const value = computed({
  get: () => props.modelValue,
  set: (val: string) => emit('update:modelValue', val),
});
</script>

<template>
  <v-select
    class="group"
    v-if="!simpleList"
    :multiple="multiple ?? false"
    :readonly="readonly ?? false"
    :variant="readonly ? 'solo' : 'outlined'"
    v-model="value"
    :label="label"
    :rules="required ? [rules.required] : []"
    :class="{required: required && !readonly}"
    :items="items"
    item-title="name"
    item-value="value"
    v-bind="$attrs"
    :clearable="!readonly"
    persistent-placeholder
    hide-details>
    <template v-if="help" #append-inner>
      <Tooltip :text="help" as-parent>
        <v-icon
          icon="mdi-help-circle-outline"
          class="cursor-help text-gray-400 opacity-0 group-focus-within:opacity-100 group-hover:opacity-100 transition-opacity duration-250"></v-icon>
      </Tooltip>
    </template>
  </v-select>
  <v-select
    v-else
    class="group"
    :multiple="multiple ?? false"
    :readonly="readonly ?? false"
    :variant="readonly ? 'solo' : 'outlined'"
    v-model="value"
    :label="label"
    :rules="required ? [rules.required] : []"
    :class="{required: required && !readonly}"
    :items="items"
    v-bind="$attrs"
    persistent-placeholder
    hide-details>
    <template v-if="help" #append-inner>
      <Tooltip :text="help" as-parent>
        <v-icon
          icon="mdi-help-circle-outline"
          class="cursor-help text-gray-400 opacity-0 group-focus-within:opacity-100 group-hover:opacity-100 transition-opacity duration-250"></v-icon>
      </Tooltip>
    </template>
  </v-select>
</template>
