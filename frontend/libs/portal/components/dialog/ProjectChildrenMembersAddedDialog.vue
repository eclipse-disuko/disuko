<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {ProjectChildrenMemberSuccessResponse} from '@disclosure-portal/model/Project';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import DCloseButton from '@shared/components/disco/DCloseButton.vue';
import {DataTableHeader} from '@shared/types/table';
import {ref} from 'vue';
import {useI18n} from 'vue-i18n';

const showDialog = defineModel<boolean>('showDialog');

defineProps<{
  targetUser: string;
  items: ProjectChildrenMemberSuccessResponse[];
}>();

const {t} = useI18n();
const headers = ref<DataTableHeader[]>([
  {
    title: t('COL_STATUS'),
    align: 'center',
    class: 'tableHeaderCell',
    value: 'success',
    width: 60,
  },
  {
    title: t('PROJECT'),
    align: 'start',
    class: 'tableHeaderCell',
    value: 'project',
  },
  {
    title: t('COL_MESSAGE'),
    align: 'start',
    class: 'tableHeaderCell',
    value: 'message',
  },
]);

const close = () => {
  showDialog.value = false;
};
</script>

<template>
  <v-form ref="form">
    <v-dialog v-model="showDialog" content-class="large" scrollable width="800" persistent>
      <v-card class="pa-8 dDialog">
        <v-card-title>
          <v-row>
            <v-col cols="10">
              <span class="text-h5">{{ t('PROJECT_CHILDREN_MEMBER_ADDED', {userId: targetUser}) }}</span>
            </v-col>
            <v-col cols="2" align="right">
              <DCloseButton @click="close" />
            </v-col>
          </v-row>
        </v-card-title>
        <v-card-text>
          <v-data-table
            density="compact"
            class="striped-table"
            :headers="headers"
            :items="items"
            item-key="projectKey"
            :footer-props="{
              'items-per-page-options': [10, 50, 100, -1],
            }">
            <template v-slot:item.success="{item}">
              <v-icon v-if="item.success" color="primary" small>mdi-check</v-icon>
              <v-icon v-else class="greyCheck" small>mdi-check</v-icon>
            </template>
            <template v-slot:item.project="{item}"> {{ item.projectName }} ({{ item.projectKey }}) </template>
            <template v-slot:item.message="{item}">
              {{ t('ADD_MULTI_USERS_' + item.message) }}
            </template>
          </v-data-table>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <DCActionButton isDialogButton size="small" variant="flat" @click="close" :text="t('BTN_OK')" />
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-form>
</template>
