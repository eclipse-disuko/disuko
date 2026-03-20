<template>
  <div class="pa-0 expanding-container">
    <v-row class="shrink pb-2">
      <v-spacer></v-spacer>
      <v-col cols="12" xs="12" sm="8" md="4" lg="2" class="pa-3 pb-0">
        <v-text-field
          autocomplete="off"
          variant="outlined"
          v-model="search"
          append-inner-icon="mdi-magnify"
          :label="t('labelSearch')"
          clearable
          density="compact"
          single-line
          hide-details></v-text-field>
      </v-col>
    </v-row>
    <v-row>
      <v-col cols="12" xs="12" v-if="allObligations" class="pa-3">
        <v-data-table
          class="striped-table"
          :headers="headers"
          :search="search"
          fixed-header
          height="445"
          density="compact"
          return-object
          show-select
          select-strategy="all"
          v-model="selectedObligations"
          :sort-by="sortBy"
          items-per-page="-1"
          :hide-default-footer="true"
          :items="allObligations"
          :item-class="getCssClassForTableRow">
          <template v-slot:bottom>
            <v-row>
              <v-col class="d-flex justify-end paddingRightItems mr-7 mb-4 fontColorItems">
                <span class="fontColorItems font-weight-light">
                  {{ t('TABLE_ITEMS') }}
                  <span class="font-weight-light fontColorItems"> {{ allObligations.length }}</span>
                </span>
              </v-col>
            </v-row>
          </template>
          <template v-slot:item.created="{item}">
            <DDateCellWithTooltip :value="item.created" />
          </template>
          <template v-slot:item.type="{item}">
            {{ getTextOfType(item.type) }}
          </template>
          <template v-slot:item.name="{item}">
            {{ viewTools.getNameForLanguage(item) }}
          </template>
          <template v-slot:item.description="{item}">
            <v-tooltip :open-delay="TOOLTIP_OPEN_DELAY_IN_MS" location="bottom">
              <template v-slot:activator="{props}">
                <span v-bind="props">{{ viewTools.getDescriptionForLanguage(item, true) }}</span>
                <span v-if="item.description.length > 120">...</span>
              </template>
              <span v-if="item.description">{{ viewTools.getDescriptionForLanguage(item) }}</span>
            </v-tooltip>
          </template>
          <template v-slot:item.warnLevel="{item}">
            <span>
              <v-tooltip :open-delay="TOOLTIP_OPEN_DELAY_IN_MS" location="bottom">
                <template v-slot:activator="{props}">
                  <v-icon v-bind="props" :color="getIconColorOfLevel(item.warnLevel)">
                    {{ getIconOfLevel(item.warnLevel) }}
                  </v-icon>
                </template>
                <span>{{ getTextOfLevel(item.warnLevel) }}</span>
              </v-tooltip>
            </span>
          </template>
        </v-data-table>
      </v-col>
    </v-row>
  </div>
</template>

<script lang="ts" setup>
import {compareLevel} from '@disclosure-portal/model/Quality';
import {getCssClassForTableRow} from '@disclosure-portal/utils/Table';
import useViewTools, {getIconColorOfLevel, getIconOfLevel} from '@disclosure-portal/utils/View';
import DDateCellWithTooltip from '@shared/components/disco/DDateCellWithTooltip.vue';
import {DataTableHeader, SortItem} from '@shared/types/table';
import {TOOLTIP_OPEN_DELAY_IN_MS} from '@shared/utils/constant';
import {onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';

import {useView} from '@disclosure-portal/composables/useView';
import {IObligation} from '@disclosure-portal/model/IObligation';
import AdminService from '@disclosure-portal/services/admin';
import {computed} from 'vue';

const {t} = useI18n();
const viewTools = useViewTools();
const {getTextOfLevel, getTextOfType} = useView();

const emits = defineEmits(['update:obligations']);
const props = defineProps<{
  obligations: IObligation[];
}>();

const selectedObligations = computed({
  get: () => props.obligations,
  set: (value: IObligation[]) => {
    emits('update:obligations', value);
  },
});

const allObligations = ref<IObligation[]>([]);
onMounted(async () => {
  const response = (await AdminService.getAllObligations()).data;
  allObligations.value = response.items;
});

const headers = ref<DataTableHeader[]>([
  {
    title: t('COL_TYPE'),
    align: 'start',
    filterable: true,
    class: 'tableHeaderCell',
    value: 'type',
  },
  {
    title: t('COL_WARN_LEVEL'),
    align: 'start',
    filterable: true,
    class: 'tableHeaderCell',
    width: 80,
    value: 'warnLevel',
    sort: compareLevel,
  },
  {
    title: t('COL_SHORT_NAME'),
    align: 'start',
    filterable: true,
    class: 'tableHeaderCell',
    value: 'name',
  },
  {
    title: t('COL_DESCRIPTION'),
    align: 'start',
    class: 'tableHeaderCell',
    value: 'description',
  },
  {
    title: t('COL_CREATED'),
    align: 'center',
    class: 'tableHeaderCell',
    value: 'created',
  },
]);
const search = ref('');

// Assuming SortItem is supposed to have properties 'key' and 'order',
// update the object structure accordingly.
const sortBy: SortItem[] = [{key: 'type', order: 'asc'}];
</script>
