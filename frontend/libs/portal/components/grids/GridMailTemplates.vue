<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import MailTemplateDialog from '@disclosure-portal/components/dialog/MailTemplateDialog.vue';
import {MailTemplate} from '@disclosure-portal/model/MailTemplate';
import mailTemplatesService from '@disclosure-portal/services/mailtemplates.service';
import {DataTableHeader, SortItem} from '@shared/types/table';
import {useTableActionSlider} from '@shared/composables/useTableActionSlider';
import {onMounted, ref, computed} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const {sliderWidth} = useTableActionSlider();

const dialog = ref<InstanceType<typeof MailTemplateDialog>>();
const items = ref<MailTemplate[]>([]);
const loaded = ref(false);
const sortBy: SortItem[] = [{key: '_key', order: 'asc'}];

const headers = computed((): DataTableHeader[] => [
  {title: t('COL_ACTIONS'), align: 'start', value: 'actions', width: sliderWidth.value, sortable: false},
  {title: t('MAIL_TEMPLATE_KEY'), align: 'start', value: '_key', sortable: true},
  {title: t('MAIL_TEMPLATE_BCC'), align: 'start', value: 'bcc', sortable: false},
  {title: t('MAIL_TEMPLATE_CC'), align: 'start', value: 'cc', sortable: false},
]);

const reload = async () => {
  loaded.value = false;
  const res = await mailTemplatesService.getAll();
  items.value = res.data;
  loaded.value = true;
};

onMounted(reload);
</script>

<template>
  <TableLayout>
    <template #buttons>
      <h1 class="text-h5">{{ t('MAIL_TEMPLATES') }}</h1>
    </template>
    <template #table>
      <div class="fill-height action-slider-table">
        <v-data-table
          density="compact"
          class="striped-table fill-height"
          :loading="!loaded"
          item-key="_key"
          :items="items"
          :headers="headers"
          :items-per-page="50"
          fixed-header
          :sort-by="sortBy">
          <template #[`item.actions`]="{item}">
            <TableActionButtons
              variant="slider"
              :buttons="[{icon: 'mdi-pencil', event: 'edit'}]"
              @edit="dialog?.open(item)" />
          </template>
        </v-data-table>
      </div>
    </template>
  </TableLayout>
  <MailTemplateDialog ref="dialog" @reload="reload" />
</template>
