<script setup lang="ts">
import SbomComponentDetailDialog from '@cli/components/dialogs/SbomComponentDetailDialog.vue';
import {comparePolicyStatus, PolicyState} from '@cli/models/PolicyRule';
import {SpdxStatusComponent, SpdxStatusInformation} from '@cli/models/Sbom';
import {projectService} from '@cli/services/projectService';
import {
  getIconColorForPolicyType,
  getIconForPolicyType,
  policyStateToTranslationKey,
} from '@disclosure-portal/utils/View';
import TableLayout from '@shared/layouts/TableLayout.vue';
import {TOOLTIP_OPEN_DELAY_IN_MS} from '@shared/utils/constant';
import {computed, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRoute} from 'vue-router';
import type {DataTableHeader, DataTableSortItem} from 'vuetify';

const route = useRoute();
const projectUuid = computed(() => route.params.id as string);
const version = computed(() => route.params.version as string);
const sbomId = computed(() => route.params.spdx as string);

const {t} = useI18n();
const search = ref('');
const sortItems = ref<DataTableSortItem[]>([{key: 'prStatus', order: 'asc'}]);

const spdxStatus = ref<SpdxStatusInformation | null>(null);
const loading = ref(false);
const error = ref<string | null>(null);

const getSbomStatus = async () => {
  if (!projectUuid.value || !version.value || !sbomId.value) {
    return;
  }

  loading.value = true;
  error.value = null;
  try {
    spdxStatus.value = await projectService.getSbomStatus(projectUuid.value, version.value, sbomId.value);
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Failed to get SBOM status';
    console.error('Error getting SBOM status:', e);
  } finally {
    loading.value = false;
  }
};

const getEffectivePrStatus = (item: SpdxStatusComponent): string => {
  if (item.prStatus === '') {
    return PolicyState.NOT_SET;
  }

  return item.prStatus || item.policyRuleStatus?.[0]?.type || PolicyState.NOT_SET;
};

watch(
  [projectUuid, version, sbomId],
  () => {
    if (version.value && sbomId.value) {
      getSbomStatus();
    }
  },
  {immediate: true},
);

// Component headers
const componentHeaders = computed<DataTableHeader[]>(() => [
  {
    key: 'prStatus',
    value: 'prStatus',
    title: t('COL_SPDX_POLICY'),
    align: 'center',
    width: 100,
    sortable: true,
    class: 'text-uppercase text-medium-emphasis text-subtitle-2 font-weight-bold pl-0',
    mustSort: true,
    sort: comparePolicyStatus,
  },
  {
    key: 'usedDecision',
    title: t('COL_DECISION'),
    align: 'center',
    width: 120,
    sortable: false,
    class: 'text-uppercase text-medium-emphasis text-subtitle-2 font-weight-bold',
  },
  {
    key: 'name',
    title: t('COL_COMPONENT_NAME'),
    align: 'start',
    width: 150,
    sortable: true,
    class: 'text-uppercase text-medium-emphasis text-subtitle-2 font-weight-bold',
  },
  {
    key: 'version',
    title: t('COL_VERSION'),
    align: 'start',
    width: 100,
    sortable: true,
    class: 'text-uppercase text-medium-emphasis text-subtitle-2 font-weight-bold',
  },
  {
    key: 'license',
    title: t('COL_COMPONENT_LICENSE'),
    align: 'start',
    width: 100,
    sortable: true,
    class: 'text-uppercase text-medium-emphasis text-subtitle-2 font-weight-bold',
  },
  {
    key: 'spdxId',
    title: t('COL_COMPONENT_ID'),
    align: 'start',
    width: 200,
    sortable: true,
    class: 'text-uppercase text-medium-emphasis text-subtitle-2 font-weight-bold',
  },
]);

// Computed components list
const components = computed(() => spdxStatus.value?.components || []);

// Dialog ref
const remarksDialog = ref<InstanceType<typeof SbomComponentDetailDialog> | null>(null);

const onRowClick = (_event: Event, dataItem: unknown) => {
  const row = (dataItem as {item?: SpdxStatusComponent})?.item;
  if (remarksDialog.value && row) {
    remarksDialog.value.open(row);
  }
};
</script>

<template>
  <div class="h-full w-full">
    <TableLayout has-tab has-title>
      <template #buttons>
        <v-spacer></v-spacer>
        <v-text-field
          autocomplete="off"
          v-model="search"
          append-inner-icon="mdi-magnify"
          :label="t('labelSearch')"
          clearable
          density="compact"
          variant="outlined"
          hide-details
          style="max-width: 400px"></v-text-field>
      </template>
      <template #table>
        <div class="fill-height">
          <v-progress-linear v-if="loading" indeterminate></v-progress-linear>
          <v-alert v-if="error" type="error" class="mb-4">{{ error }}</v-alert>
          <v-data-table
            density="compact"
            fixed-header
            :headers="componentHeaders"
            :items="components"
            :search="search"
            :sort-by="sortItems"
            :items-per-page="25"
            :footer-props="{
              'items-per-page-options': [10, 25, 50, 100, -1],
            }"
            hover
            @click:row="onRowClick"
            class="striped-table fill-height cursor-pointer">
            <template v-slot:item.prStatus="{item}">
              <v-tooltip :open-delay="TOOLTIP_OPEN_DELAY_IN_MS" location="bottom" content-class="dpTooltip">
                <template v-slot:activator="{props}">
                  <v-icon small v-bind="props" :color="getIconColorForPolicyType(getEffectivePrStatus(item))">
                    {{ getIconForPolicyType(getEffectivePrStatus(item)) }}
                  </v-icon>
                </template>
                <span v-if="item.policyRuleStatus?.length">
                  <div v-for="(prStatus, index) in item.policyRuleStatus" :key="index">
                    {{ prStatus.name }} ({{ prStatus.licenseMatched }})
                    <div v-if="prStatus.description" class="text-caption">
                      {{ prStatus.description }}
                    </div>
                  </div>
                </span>
                <span v-else>
                  {{ t(policyStateToTranslationKey(PolicyState.NOT_SET)) }}
                </span>
              </v-tooltip>
            </template>
            <template v-slot:item.usedDecision="{item}">
              <span v-if="item?.usedDecision">
                <v-tooltip :open-delay="TOOLTIP_OPEN_DELAY_IN_MS" location="bottom" content-class="dpTooltip">
                  <template v-slot:activator="{props}">
                    <v-icon small v-bind="props">mdi-information-outline</v-icon>
                  </template>
                  <span>
                    <div class="text-subtitle-1">{{ t('TT_LICENSE_DECISION') }}</div>
                    <div class="d-text d-secondary-text">{{ t('TT_FOR_THE_EXPRESSION') }}</div>
                    <div class="d-text d-secondary-text">{{ item.usedDecision.expression }}</div>
                    <div class="d-text d-secondary-text">{{ t('TT_A_DECISION_WAS_MADE') }}</div>
                    <div class="text-subtitle-1">{{ item.usedDecision.name }} ({{ item.usedDecision.licenseID }})</div>
                  </span>
                </v-tooltip>
              </span>
            </template>
          </v-data-table>
        </div>
      </template>
    </TableLayout>
  </div>
  <SbomComponentDetailDialog ref="remarksDialog" />
</template>
