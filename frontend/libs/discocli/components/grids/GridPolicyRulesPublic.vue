<template>
  <TableLayout has-tab has-title>
    <template #description> </template>

    <template #buttons>
      <v-spacer></v-spacer>
      <v-text-field
        v-model="search"
        append-inner-icon="mdi-magnify"
        :label="t('labelSearch')"
        variant="outlined"
        clearable
        density="compact"
        hide-details
        class="max-w-[400px]" />
    </template>

    <template #table>
      <div ref="tablePolicyRules" class="fill-height">
        <v-data-table
          fixed-header
          :loading="loading"
          :headers="headers"
          :items="policyRules"
          :search="search"
          :sort-by="sortItems"
          :footer-props="{'items-per-page-options': [10, 50, 100, -1]}"
          density="compact"
          class="striped-table fill-height">
          <template #[`item.actions`]="{item}">
            <TableActionButtons
              variant="compact"
              :buttons="getActionButtons(item as unknown as PolicyRule)"
              @copy="copyToClipboard(getPolicyRuleInfoForClipboard(item as unknown as PolicyRule))"
              @showLicenses="showLicenses(item as unknown as PolicyRule)" />
          </template>

          <template v-slot:no-data>
            <v-alert type="info" variant="tonal" :text="t('NO_POLICY_RULES')" class="ma-2" />
          </template>
        </v-data-table>
      </div>
    </template>
  </TableLayout>
  <PolicyLicensesDialog
    :show="showLicensesDialog"
    :policy-rule="selectedRule"
    @update:show="showLicensesDialog = $event"
    @close="selectedRule = null" />
</template>

<script setup lang="ts">
import PolicyLicensesDialog from '@cli/components/dialogs/PolicyLicensesDialog.vue';
import type {PolicyRule} from '@cli/models/PolicyRule';
import {projectService} from '@cli/services/projectService';
import {useAppStore} from '@cli/stores/app';
import {formatDateTime} from '@disclosure-portal/utils/View';
import TableActionButtons, {TableActionButtonsProps} from '@shared/components/TableActionButtons.vue';
import TableLayout from '@shared/layouts/TableLayout.vue';
import {useClipboard} from '@shared/utils/clipboard';
import {computed, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import type {DataTableHeader} from 'vuetify';

const {t} = useI18n();
const appStore = useAppStore();
const currentProject = computed(() => appStore.getCurrentProject());
const rules = ref<PolicyRule[]>([]);
const search = ref('');
const loading = ref(false);
const sortItems = ref([{key: 'name', order: false}]);
const showLicensesDialog = ref(false);
const selectedRule = ref<PolicyRule | null>(null);
const {copyToClipboard} = useClipboard();

const showLicenses = (item: PolicyRule) => {
  try {
    selectedRule.value = item;
    showLicensesDialog.value = true;
  } catch (error) {
    console.error('Error in showLicenses:', error);
  }
};

const fetchPolicyRules = async () => {
  const projectUuid = currentProject.value?.uuid;
  if (!projectUuid) return;

  loading.value = true;
  try {
    const fetchedRules = await projectService.getProjectPolicyRules(projectUuid);
    rules.value = (fetchedRules || []).map((rule: unknown) => {
      const r = rule as Record<string, unknown>;

      // Map according to the actual API response structure
      const mappedRule: PolicyRule = {
        id: String(r.Key || r.key || r.id || ''),
        key: String(r.Key || r.key || r.id || ''),
        name: String(r.Name || r.name || ''),
        type: String(r.Type || r.type || ''),
        description: String(r.Description || r.description || ''),
        created: String(r.Created || r.created || r.createdAt || ''),
        updated: String(r.Updated || r.updated || r.updatedAt || ''),
        licenses: ((r.Licenses || r.licenses || []) as Array<Record<string, unknown>>).map((license) => ({
          key: String(license.Key || license.key || ''),
          identifier: String(license.Identifier || license.identifier || ''),
          name: String(license.Name || license.name || ''),
          policyType: String(r.Type || r.type || ''), // Add policy type to each license
        })),
      };
      return mappedRule;
    });
  } catch (error) {
    console.error('Error fetching policy rules:', error);
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchPolicyRules();
});

const headers = computed<DataTableHeader[]>(() => [
  {
    key: 'actions',
    title: t('COL_ACTIONS'),
    align: 'center',
    sortable: false,
    width: 120,
    class: 'tableHeaderCell',
  },
  {
    key: 'name',
    title: t('COL_POLICY_NAME'),
    align: 'start',
    sortable: true,
    width: 150,
    class: 'tableHeaderCell',
  },
  {
    key: 'description',
    title: t('COL_DESCRIPTION'),
    align: 'start',
    sortable: true,
    width: 450,
    class: 'tableHeaderCell',
  },
]);

const policyRules = computed(() => {
  // Group rules by key
  const groupedRules = rules.value.reduce(
    (acc, rule) => {
      if (!acc[rule.key]) {
        // First rule with this key
        acc[rule.key] = {
          ...rule,
          licenses: [...(rule.licenses || [])],
          policyTypes: new Set([rule.type]),
        };
      } else {
        // Merge licenses without duplicates
        const existingLicenseKeys = new Set(acc[rule.key].licenses.map((l) => l.key));
        const newLicenses = rule.licenses?.filter((license) => !existingLicenseKeys.has(license.key)) || [];
        acc[rule.key].licenses.push(...newLicenses);
        acc[rule.key].policyTypes.add(rule.type);

        // Combine names if different
        if (rule.name !== acc[rule.key].name) {
          acc[rule.key].name = `${acc[rule.key].name}, ${rule.name}`;
        }

        // Combine descriptions if different
        if (rule.description && rule.description !== acc[rule.key].description) {
          acc[rule.key].description = `${acc[rule.key].description || ''}\n${rule.description}`.trim();
        }
      }
      return acc;
    },
    {} as Record<string, PolicyRule & {policyTypes: Set<string>}>,
  );

  // Convert grouped rules back to array
  return Object.values(groupedRules).map((rule) => ({
    ...rule,
    type: Array.from(rule.policyTypes).join(', '), // Keep track of all policy types
  }));
});

const getPolicyRuleInfoForClipboard = (item: PolicyRule): string => {
  return `Policy Rule Reference
  Project Name: ${currentProject.value?.name}
  Project Identifier: ${currentProject.value?.uuid}
  Rule Name: ${item.name}
  Rule Key: ${item.key}
  Type: ${item.type}
  Description: ${item.description || '-'}
  Reference Timestamp: ${formatDateTime(new Date().toISOString())} (UTC)`;
};

const getActionButtons = (_item: PolicyRule): TableActionButtonsProps['buttons'] => {
  return [
    {
      icon: 'mdi-content-copy',
      hint: t('TT_COPY_REFERENCE_INFO'),
      event: 'copy',
      show: true,
    },
    {
      icon: 'mdi-open-in-new',
      hint: t('POLICY_RULE'),
      event: 'showLicenses',
      show: true,
    },
  ];
};
</script>
