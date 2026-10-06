<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<template>
  <TableLayout has-title has-tab>
    <template #buttons>
      <v-spacer></v-spacer>
      <DSearchField v-model="search" />
    </template>
    <template #table>
      <PolicyRulesTable class="fill-height" v-model="items" :loading="!dataLoaded"></PolicyRulesTable>
    </template>
  </TableLayout>
</template>

<script setup lang="ts">
import {PolicyRulesAssignmentsDto} from '@disclosure-portal/model/PolicyRule';
import {default as LicenseService} from '@disclosure-portal/services/license';
import {onBeforeMount, ref} from 'vue';

const props = defineProps<{
  licenseId: string;
  disabled?: boolean;
}>();

const search = ref('');
const items = ref<PolicyRulesAssignmentsDto[]>([]);
const dataLoaded = ref(false);

const reload = async () => {
  dataLoaded.value = false;
  const response = await LicenseService.getPolicyRuleAssignmentsForThisLicence(props.licenseId);
  items.value = response.policyRulesAssignments;
  dataLoaded.value = true;
};

onBeforeMount(async () => {
  await reload();
});
</script>
