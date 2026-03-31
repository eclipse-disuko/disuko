// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import {useSlots} from 'vue';
import {InternalDataTableHeader} from 'vuetify/lib/components/VDataTable/types';
import {IconValue} from 'vuetify/lib/composables/icons';

interface Props {
  column: InternalDataTableHeader;
  getSortIcon?: (column: InternalDataTableHeader) => IconValue;
  toggleSort?: (column: InternalDataTableHeader) => void;
}

const {column, getSortIcon, toggleSort} = defineProps<Props>();

const slots = useSlots();
</script>

<template>
  <div
    class="flex items-center"
    :class="{
      'justify-start': column.align === 'start',
      'justify-center': column.align === 'center',
      'justify-end': column.align === 'end',
    }">
    <slot v-if="slots.settings" name="settings" />
    <span v-if="column.title" :class="{'ml-6': slots.settings}" class="mr-1">{{ column.title }}</span>
    <slot v-if="slots.filter" name="filter" />
    <v-icon
      v-if="getSortIcon && toggleSort"
      class="v-data-table-header__sort-icon"
      :icon="getSortIcon(column)"
      @click="toggleSort(column)" />
  </div>
</template>
