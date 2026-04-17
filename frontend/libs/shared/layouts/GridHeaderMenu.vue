<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {DataTableHeaderFilterItems} from '@shared/types/table';
import {ref} from 'vue';
import {useI18n} from 'vue-i18n';

interface Props {
  allItems: DataTableHeaderFilterItems[];
  selectedItems: string[];
  showReset: boolean;
  cardTitle: string;
  resetHint: string;
  selectLabel: string;
}

defineProps<Props>();

const emit = defineEmits<{
  reset: [];
  update: [value: string[]];
}>();

const {t} = useI18n();

const showMenu = ref(false);
</script>

<template>
  <v-menu offset-y :close-on-content-click="false" v-model="showMenu">
    <template v-slot:activator="{props}">
      <slot name="activator" :props="props" />
    </template>
    <v-card class="w-[400px]">
      <Stack class="gap-1">
        <Stack direction="row" justify="between" align="center" class="ma-1 mr-2 ml-4">
          <span class="align-self-center text-base">{{ cardTitle }}</span>
          <DCloseButton @click="showMenu = false" />
        </Stack>
        <Stack
          direction="row"
          justify="between"
          align="center"
          class="mb-4 ml-4"
          :class="{
            'mr-6': !showReset,
          }">
          <v-select
            :model-value="selectedItems"
            :items="allItems"
            :label="selectLabel"
            :clearable="false"
            multiple
            menu
            hide-details
            variant="outlined"
            density="compact"
            transition="scale-transition"
            @update:modelValue="(value: string[]) => emit('update', value)">
            <template v-slot:item="{props, item}">
              <v-list-item v-bind="props" :title="undefined" :disabled="item.raw?.disabled" class="px-2 py-0">
                <template v-slot:prepend="{isSelected}">
                  <v-checkbox hide-details :model-value="isSelected" />
                </template>
                <v-icon v-if="item.raw?.icon" small :color="item.raw.iconColor">{{ item.raw.icon }}</v-icon>
                <span
                  :style="{color: item.raw?.textColor || 'inherit'}"
                  class="text-sm"
                  :class="{'ml-1': item.raw?.icon, 'font-bold': item.raw?.textBold}"
                >{{ item.raw?.text || item.raw.value }}</span
                >
                <v-chip
                  v-if="item.raw?.chip"
                  :color="item.raw?.chipColor || 'default'"
                  label
                  size="x-small"
                  class="ml-1">
                  {{ item.raw.chip }}
                </v-chip>
              </v-list-item>
            </template>
            <template v-slot:selection="{item, index}">
              <div v-if="index === 0" class="d-flex align-center">
                <v-icon v-if="item.raw?.icon" small :color="item.raw.iconColor">{{ item.raw.icon }}</v-icon>
                <span
                  v-if="index === 0"
                  :style="{color: item.raw?.textColor || 'inherit'}"
                  class="text-sm"
                  :class="{'ml-1': item.raw?.icon, 'font-bold': item.raw?.textBold}"
                >{{ item.raw?.text || item.raw.value }}</span
                >
                <v-chip
                  v-if="item.raw?.chip"
                  :color="item.raw?.chipColor || 'default'"
                  label
                  size="x-small"
                  class="ml-1">
                  {{ item.raw.chip }}
                </v-chip>
              </div>
              <span v-if="index === 1" class="text-sm">
                (+{{ selectedItems.length - 1 }} {{ t('OTHERS', selectedItems.length) }})
              </span>
            </template>
          </v-select>
          <DIconButton
            v-if="showReset"
            color="primary"
            variant="text"
            :hint="resetHint"
            icon="mdi-autorenew"
            class="mr-3"
            @clicked="emit('reset')" />
        </Stack>
      </Stack>
    </v-card>
  </v-menu>
</template>
