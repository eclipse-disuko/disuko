<script setup lang="ts">
import TableLayout from '@shared/layouts/TableLayout.vue';
import License from '@disclosure-portal/model/License';
import {computed} from 'vue';
import {DataTableHeader} from '@shared/types/table';
import {useI18n} from 'vue-i18n';

const props = defineProps<{
  license: License;
}>();

const {t} = useI18n();

const headers = computed((): DataTableHeader[] => {
  return [
    {
      title: t('COL_NAME'),
      align: 'start',
      width: '180',
      value: 'licenseId',
    },
    {
      title: t('COL_DESCRIPTION'),
      align: 'start',
      value: 'description',
    },
  ];
});
</script>

<template>
  <TableLayout has-tab has-title>
    <template #table>
      <v-data-table
        v-if="props.license?.aliases"
        :headers="headers"
        fixed-header
        density="compact"
        class="striped-table fill-height"
        item-key="_key"
        :items-per-page="-1"
        :footer-props="{
          'items-per-page-options': [10, 50, 100, -1],
        }"
        :items="props.license?.aliases">
        <template v-slot:item.licenseId="{item}">
          <v-text-field
            autocomplete="off"
            class="my-1 pt-5"
            v-if="item._key === ''"
            density="compact"
            solo
            variant="outlined"
            v-model="item.licenseId" />
          <span v-else>{{ item.licenseId }}</span>
        </template>
        <template #[`item.description`]="{item}">
          <v-text-field
            autocomplete="off"
            class="my-1 pt-5"
            v-if="item._key === ''"
            density="compact"
            variant="outlined"
            v-model="item.description" />
          <span v-else>{{ item.description }}</span>
        </template>
      </v-data-table>
    </template>
  </TableLayout>
</template>
