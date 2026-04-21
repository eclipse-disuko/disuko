<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {ConfirmationType, IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import Icons from '@disclosure-portal/constants/icons';
import {PolicyLabels} from '@disclosure-portal/constants/policyLabels';
import {Decision} from '@disclosure-portal/model/Decision';
import projectService from '@disclosure-portal/services/projects';
import {useLabelStore} from '@disclosure-portal/stores/label.store';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useUserStore} from '@disclosure-portal/stores/user';
import {RightsUtils} from '@disclosure-portal/utils/Rights';
import {escapeHtml} from '@disclosure-portal/utils/Validation';
import {formatDateTime} from '@disclosure-portal/utils/View';
import {DataTableHeader, DataTableHeaderFilterItems, SortItem} from '@shared/types/table';
import {computed, onMounted, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {useHeaderSettings} from '@shared/composables/useHeaderSettings';

const {t} = useI18n();
const projectStore = useProjectStore();
const userStore = useUserStore();
const labelStore = useLabelStore();

const search = ref('');
const sortBy: SortItem[] = [{key: 'updated', order: 'desc'}];
const tableHeight = ref(0);
const loading = ref(false);
const items = ref<Decision[]>([]);
const cancelVisible = ref(false);
const confirmCancelConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);

const currentProject = computed(() => projectStore.currentProject!);

const isVehicleProject = computed(() =>
  currentProject.value.policyLabels.some(
    (lbl) => labelStore.getLabelByKey(lbl)?.name === PolicyLabels.VEHICLE_PLATFORM,
  ),
);

const possibleStatus = computed((): DataTableHeaderFilterItems[] => {
  return [...new Set(items.value.map((item) => (item.active ? 'active' : 'cancelled')))].map((item) => ({
    value: item,
    text: item === 'active' ? t('TT_LICENSE_RULE_ACTIVE') : t('TT_LICENSE_RULE_CANCELLED'),
    iconColor: item === 'active' ? 'success' : 'warning',
    icon: Icons.CIRCLE_FILLED,
  }));
});

const possibleType = computed((): DataTableHeaderFilterItems[] =>
  [...new Set(items.value.map((item) => item.type))].map((item) => ({value: item})),
);

const selectedStatus = ref<string[]>([]);
const selectedType = ref<string[]>([]);

const filteredList = computed(() => {
  return items.value.filter((item: Decision) => filterOnStatus(item) && filterOnType(item));
});

const filterOnStatus = (item: Decision): boolean => {
  if (!selectedStatus.value.length) {
    return true;
  }
  return selectedStatus.value.includes(item.active ? 'active' : 'cancelled');
};

const filterOnType = (item: Decision): boolean => {
  if (!selectedType.value.length) {
    return true;
  }
  return selectedType.value.includes(item.type);
};

const headers: DataTableHeader[] = [
  {
    title: 'COL_ACTIONS',
    align: 'center',
    width: 120,
    minWidth: 120,
    maxWidth: 130,
    value: 'actions',
    sortable: false,
  },
  {
    title: 'COL_STATUS',
    align: 'center',
    value: 'active',
    sortable: true,
    width: 150,
    minWidth: 150,
    maxWidth: 160,
  },
  {
    title: 'COL_TYPE',
    align: 'start',
    value: 'type',
    sortable: true,
    width: 120,
    minWidth: 120,
    maxWidth: 130,
  },
  {
    title: 'COL_SPDX_FILENAME',
    align: 'start',
    value: 'sbomName',
    sortable: true,
    width: 350,
    minWidth: 350,
    maxWidth: 350,
  },
  {
    title: 'COL_COMPONENT_NAME',
    align: 'start',
    value: 'componentName',
    sortable: true,
    width: 250,
    minWidth: 250,
  },
  {
    title: 'COL_COMPONENT_VERSION',
    align: 'start',
    value: 'componentVersion',
    sortable: true,
    width: 100,
    maxWidth: 200,
  },
  {
    title: 'COL_LICENSE_EXPRESSION',
    align: 'start',
    value: 'licenseExpression',
    sortable: true,
    width: 300,
    minWidth: 300,
  },
  {
    title: 'COL_DECISION',
    align: 'start',
    value: 'licenseDecisionId',
    sortable: true,
    width: 300,
    minWidth: 300,
  },
  {
    title: 'COL_COMMENT',
    align: 'start',
    value: 'comment',
    sortable: true,
    width: 120,
  },
  {
    title: 'COL_CREATED',
    align: 'start',
    value: 'created',
    width: 108,
    minWidth: 108,
    maxWidth: 130,
    sortable: true,
  },
  {
    title: 'COL_UPDATED',
    align: 'start',
    value: 'updated',
    width: 108,
    minWidth: 108,
    maxWidth: 130,
    sortable: true,
  },
  {
    title: 'COL_CREATOR',
    align: 'start',
    value: 'creator',
    width: 120,
    sortable: true,
  },
];

const tableName = 'LicenseRules';
const headerSettings = useHeaderSettings({tableName, headers});
const {filteredHeaders} = headerSettings;

const icons = Icons;

const reload = async (): Promise<void> => {
  loading.value = true;
  items.value = (await projectService.getDecisions(currentProject.value._key)).data;
  loading.value = false;
};

const formatText = (text: string): string => {
  text = escapeHtml(text);
  if (text.includes(' AND ') || text.includes(' OR ')) {
    return text
      .replace(/ AND /g, ' <strong class="db-highlight">AND</strong> ')
      .replace(/ OR /g, ' <strong class="db-highlight">OR</strong> ');
  }
  return text;
};

const isCancelActionVisible = (item: Decision) => {
  if (currentProject.value.isDeprecated) return false;

  if (!item.active) return false;

  const isProjectResponsible = userStore.getProfile.user === currentProject.value.responsible;
  const isFossOffice = RightsUtils.isFOSSOffice();
  const isDomainAdmin = RightsUtils.isDomainAdmin();
  const isVehicle = isVehicleProject.value;

  if (item.policyEvaluated === 'warn' || item.type === 'license') {
    return isVehicle ? isFossOffice : isProjectResponsible;
  }

  if (item.policyEvaluated === 'deny') {
    return isVehicle ? isFossOffice : isDomainAdmin;
  }

  return false;
};

const cancelLicenseRule = (item: Decision) => {
  confirmCancelConfig.value = {
    key: item.key,
    name: '',
    type: ConfirmationType.CONFIRM,
    description: t('DLG_CONFIRMATION_DESCRIPTION_CANCEL_LICENSE_RULE', {
      name: item.componentName,
      decision: item.licenseDecisionName,
    }),
    okButton: 'Btn_confirm',
  };
};

const doCancelLicenseRule = async (config: IConfirmationDialogConfig) => {
  const decision = items.value.find((d) => d.key === config.key);
  if (decision?.type === 'policy') {
    await projectService.cancelPolicyDecision(currentProject.value._key, config.key);
  } else {
    await projectService.cancelLicenseRule(currentProject.value._key, config.key);
  }

  await reload();
};

watch(
  () => currentProject.value,
  () => reload(),
  {deep: true},
);

onMounted(async () => {
  await reload();
});
</script>

<template>
  <TableLayout has-title has-tab>
    <template #buttons>
      <v-spacer></v-spacer>
      <DSearchField v-model="search" />
    </template>
    <template #table>
      <v-data-table
        density="compact"
        class="striped-table custom-data-table h-full"
        :headers="filteredHeaders"
        fixed-header
        item-key="key"
        :sort-by="sortBy"
        :height="tableHeight"
        sort-desc
        :items="filteredList"
        :search="search"
        items-per-page="50"
        :footer-props="{'items-per-page-options': [10, 50, 100, -1]}">
        <template v-slot:[`header.actions`]="{column}">
          <HeaderSettings :column="column" :grid-name="tableName" />
          <span class="ml-6">{{ column.title }}</span>
        </template>
        <template v-slot:[`header.active`]="{column, toggleSort, getSortIcon}">
          <span class="mr-1">{{ column.title }}</span>
          <GridHeaderFilterIcon
            v-model="selectedStatus"
            :column="column"
            :label="t('STATUS')"
            :allItems="possibleStatus">
          </GridHeaderFilterIcon>
          <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
        </template>
        <template v-slot:[`header.type`]="{column, toggleSort, getSortIcon}">
          <span class="mr-1">{{ column.title }}</span>
          <GridHeaderFilterIcon v-model="selectedType" :column="column" :label="t('TYPE')" :allItems="possibleType" />
          <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
        </template>
        <template v-slot:[`item.actions`]="{item}">
          <DIconButton
            icon="mdi-cancel"
            :hint="t('TT_cancel_license_rule')"
            @clicked="cancelLicenseRule(item)"
            v-if="isCancelActionVisible(item)" />
        </template>
        <template v-slot:[`item.active`]="{item}">
          <v-icon size="x-small" :color="item.active ? 'success' : 'warning'">{{ icons.CIRCLE_FILLED }}</v-icon>
          <Tooltip>
            <span>{{ item.active ? t('TT_LICENSE_RULE_ACTIVE') : t('TT_LICENSE_RULE_CANCELLED') }}</span>
          </Tooltip>
        </template>
        <template v-slot:[`item.sbomName`]="{item}">
          <span>{{ formatDateTime(item.sbomUploaded) }}&nbsp;-&nbsp;{{ item.sbomName }}</span>
          <br />
          <span class="font-weight-bold">UUID: </span>
          <span>{{ item.sbomId }}</span>
        </template>
        <template v-slot:[`item.licenseExpression`]="{item}">
          <Truncated><span v-html="formatText(item.licenseExpression)"></span></Truncated>
        </template>
        <template v-slot:[`item.licenseDecisionId`]="{item}">
          <Truncated v-if="item.type === 'license'">{{ item.licenseDecisionId }}</Truncated>
          <PolicyDecisionItem v-else :decision="item" :previewMode="false" :license-id="item.licenseMatchedId" />
        </template>
        <template v-slot:[`item.created`]="{item}">
          <DDateCellWithTooltip :value="item.created" />
        </template>
        <template v-slot:[`item.updated`]="{item}">
          <DDateCellWithTooltip :value="item.updated" />
        </template>
      </v-data-table>
    </template>
  </TableLayout>

  <ConfirmationDialog v-model:showDialog="cancelVisible" :config="confirmCancelConfig" @confirm="doCancelLicenseRule" />
</template>
