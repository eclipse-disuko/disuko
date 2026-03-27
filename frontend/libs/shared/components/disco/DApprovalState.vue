// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import {ApproverRoles, InternalApproval} from '@disclosure-portal/model/Approval';
import {formatDateAndTime} from '@disclosure-portal/utils/Table';
import {computed} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();

const props = defineProps<{
  internal: InternalApproval;
}>();
const roles = computed(
  () =>
    new Map<ApproverRoles, string>([
      [ApproverRoles.Supplier1, t('APPROVER_ROLE_DEVELOPER_FRI')],
      [ApproverRoles.Supplier2, t('APPROVER_ROLE_DEVELOPER_SRI')],
      [ApproverRoles.Customer1, t('APPROVER_ROLE_OWNER_FRI')],
      [ApproverRoles.Customer2, t('APPROVER_ROLE_OWNER_SRI')],
    ]),
);
</script>

<template>
  <v-table class="striped-table">
    <thead>
      <tr>
        <th width="100" class="pa-2 font-weight-bold">
          {{ t('COL_APPROVAL_HISTORY_ROLE') }}
        </th>
        <th width="100" class="pa-2 font-weight-bold">
          {{ t('COL_APPROVAL_HISTORY_STATE') }}
        </th>
        <th width="200" class="pa-2 font-weight-bold">
          {{ t('COL_APPROVAL_HISTORY_APPROVER') }}
        </th>
        <th width="250" class="pa-2 font-weight-bold">
          {{ t('COL_COMMENT') }}
        </th>
        <th width="150" class="pa-2 font-weight-bold">
          {{ t('COL_UPDATED') }}
        </th>
      </tr>
    </thead>
    <tbody>
      <template v-for="(approver, index) in props.internal.approver" :key="index">
        <tr v-if="approver">
          <td class="pa-2">
            <span class="userState">{{ roles.get(index) }}</span>
          </td>
          <td class="pa-2">
            <span class="userState" v-if="props.internal.states[index].state">
              {{ t('COL_APPROVAL_STATUS_' + props.internal.states[index].state) }}
            </span>
          </td>
          <td class="pa-2">
            <span class="userTitle">{{ props.internal.approverFullName[index] }} ({{ approver }})</span>
          </td>
          <td class="pa-2">
            <span class="userComment">{{ props.internal.comments[index] }}</span>
          </td>
          <td class="pa-2">
            <span class="userUpdated" v-if="props.internal.states[index].updated">
              {{ formatDateAndTime(props.internal.states[index].updated) }}
            </span>
            <span class="userUpdated" v-else></span>
          </td>
        </tr>
      </template>
    </tbody>
  </v-table>
</template>
