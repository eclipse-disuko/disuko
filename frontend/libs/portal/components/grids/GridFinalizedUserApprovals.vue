<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script lang="ts" setup>
import {UserApproval} from '@disclosure-portal/model/Users';
import {DataTableHeader} from '@shared/types/table';
import {computed} from 'vue';
import {useI18n} from 'vue-i18n';

const props = defineProps<{
  approvals: UserApproval[];
  loading: boolean;
}>();

const {t} = useI18n();

const items = computed(() => props.approvals.filter((a) => !a.isActive));

const headers = computed<DataTableHeader[]>(() => [
  {
    title: t('COL_APPROVAL_TYPE'),
    align: 'center',
    value: 'approvalType',
    width: 140,
    sortable: true,
  },
  {
    title: t('COL_IS_CREATOR'),
    align: 'center',
    value: 'isCreator',
    width: 120,
  },
  {
    title: t('COL_IS_APPROVER'),
    align: 'center',
    value: 'isApprover',
    width: 120,
  },
  {
    title: t('COL_PROJECT_NAME'),
    align: 'start',
    value: 'projectName',
    sortable: true,
  },
]);
</script>

<template>
  <v-data-table
    :loading="loading"
    :items="items"
    :headers="headers"
    fixed-header
    class="striped-table custom-data-table h-full"
    :items-per-page="50">
    <template #item.approvalType="{item}">
      {{ t('COL_APPROVAL_TITLE_TYPE_' + item.approvalType) }}
    </template>
    <template #item.isCreator="{item}">
      <v-icon icon="mdi-check" :color="item.isCreator ? 'primary' : 'tableBorderColor'"></v-icon>
    </template>
    <template #item.isApprover="{item}">
      <v-icon icon="mdi-check" :color="item.isApprover ? 'primary' : 'tableBorderColor'"></v-icon>
    </template>
  </v-data-table>
</template>
