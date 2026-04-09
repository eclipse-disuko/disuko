<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script lang="ts" setup>
import {UserApproval} from '@disclosure-portal/model/Users';
import AdminService from '@disclosure-portal/services/admin';
import Stack from '@shared/layouts/Stack.vue';
import TableLayout from '@shared/layouts/TableLayout.vue';
import {ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRoute} from 'vue-router';

type ApprovalTab = {
  id: string;
  buttonIcon: string;
  buttonText: string;
};

const {t} = useI18n();
const route = useRoute();

const loading = ref(false);
const approvals = ref<UserApproval[]>([]);

loading.value = true;
AdminService.getUserApprovals(route.params.uuid as string).then((data) => {
  approvals.value = data;
  loading.value = false;
});

const tabs = ref<ApprovalTab[]>([
  {
    id: 'ongoingApprovals',
    buttonIcon: 'mdi-progress-clock',
    buttonText: 'TAB_ONGOING_APPROVALS',
  },
  {
    id: 'finalizedApprovals',
    buttonIcon: 'mdi-check-circle-outline',
    buttonText: 'TAB_FINALIZED_APPROVALS',
  },
]);

const selectedTabId = ref<string>('ongoingApprovals');
</script>

<template>
  <TableLayout has-tab has-title gap="0">
    <template #buttons>
      <Stack direction="row" class="mb-4">
        <v-btn
          v-for="tab in tabs"
          size="small"
          :key="tab.id"
          @click="selectedTabId = tab.id"
          :variant="selectedTabId === tab.id ? 'tonal' : 'text'"
          :class="{active: selectedTabId === tab.id}"
          class="text-none card-border"
          min-width="130px">
          <v-icon color="primary" class="pr-2">{{ tab.buttonIcon }}</v-icon>
          {{ t(tab.buttonText) }}
        </v-btn>
      </Stack>
      <v-spacer></v-spacer>
    </template>
    <template #table>
      <div class="h-full" v-if="selectedTabId === 'ongoingApprovals'">
        <GridOngoingUserApprovals :approvals="approvals.filter((a) => a.isActive)" :loading="loading" />
      </div>
      <div class="h-full" v-if="selectedTabId === 'finalizedApprovals'">
        <GridFinalizedUserApprovals :approvals="approvals.filter((a) => !a.isActive)" :loading="loading" />
      </div>
    </template>
  </TableLayout>
</template>
