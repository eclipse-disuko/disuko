// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import {ConfirmationType, IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import ConfirmationDialog from '@disclosure-portal/components/dialog/ConfirmationDialog.vue';
import NewPolicyRuleDialog from '@disclosure-portal/components/dialog/NewPolicyRuleDialog.vue';
import {IDefaultSelectItem} from '@disclosure-portal/model/ISelectItem';
import Label from '@disclosure-portal/model/Label';
import PolicyRule from '@disclosure-portal/model/PolicyRule';
import AdminService from '@disclosure-portal/services/admin';
import policyRuleService from '@disclosure-portal/services/policyrules';
import {downloadFile} from '@disclosure-portal/utils/download';
import {RightsUtils} from '@disclosure-portal/utils/Rights';
import {formatDateAndTime, getCssClassForTableRow} from '@disclosure-portal/utils/Table';
import {openUrl} from '@disclosure-portal/utils/url';
import {getStrWithMaxLength} from '@disclosure-portal/utils/View';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import DDateCellWithTooltip from '@shared/components/disco/DDateCellWithTooltip.vue';
import DIconButton from '@shared/components/disco/DIconButton.vue';
import {TableActionButtonsProps} from '@shared/components/TableActionButtons.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import TableLayout from '@shared/layouts/TableLayout.vue';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {DataTableHeader, DataTableItem, SortItem} from '@shared/types/table';
import dayjs from 'dayjs';
import {computed, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRouter} from 'vue-router';

const {t} = useI18n();
const breadcrumbs = useBreadcrumbsStore();
const {info} = useSnackbar();
const router = useRouter();

const confirmConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);
const confirmVisible = ref(false);
const confirmDeprConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);
const confirmDeprVisible = ref(false);
const confirmCopyConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);
const confirmCopyVisible = ref(false);
const possibleStatus = ref<IDefaultSelectItem[]>([
  {
    text: t('PR_STATUS_ACTIVE'),
    value: 'active',
  },
  {
    text: t('PR_STATUS_INACTIVE'),
    value: 'inactive',
  },
  {
    text: t('PR_STATUS_DEPRECATED'),
    value: 'deprecated',
  },
]);
const menu = ref(false);
const selectedFilterStatus = ref<string[]>(['active', 'inactive']);
const items = ref<PolicyRule[]>([]);
const search = ref('');
const isPolicyManager = ref(false);
const policyLabels = ref<Label[]>([]);
const sortItems = ref<SortItem[]>([{key: 'Name', order: 'asc'}]);
const policyRuleDialogRef = ref();
const currentPolicyRuleForAction = ref<PolicyRule | null>(null);

const initBreadcrumbs = () => {
  breadcrumbs.setCurrentBreadcrumbs([
    {title: t('BC_Dashboard'), disabled: false, href: '/dashboard/home'},
    {title: t('POLICY_RULES'), disabled: false, href: '/dashboard/policyrules'},
  ]);
};

const reload = async () => {
  items.value = (await policyRuleService.getAllPolicyRules()).data;
  await reloadLabels();
};

const reloadLabels = async () => {
  policyLabels.value = (await AdminService.getPolicyLabels()).data;
};

const showDeletionConfirmationDialog = (item: PolicyRule) => {
  confirmConfig.value = {
    type: ConfirmationType.DELETE,
    key: item._key,
    name: item.Name,
    okButtonIsDisabled: false,
    okButton: 'Btn_delete',
    description: 'DLG_CONFIRMATION_DESCRIPTION',
  } as IConfirmationDialogConfig;
  confirmVisible.value = true;
};

const onClickRow = (event: Event, table: DataTableItem<PolicyRule>) => {
  openUrl('/dashboard/policyrules/' + table.item._key, router);
};

const showDeprecationConfirmationDialog = async (pr: PolicyRule) => {
  confirmDeprConfig.value = {
    type: ConfirmationType.DEPRECATE,
    key: pr._key,
    name: pr.Name,
    description: 'DLG_PR_DEPRECATION_CONFIRMATION_DESCRIPTION',
    emphasiseText: 'PR_DEPRECATION_UNREVERTABLE',
    emphasiseConfirmationText: 'PR_DEPRECATION_UNREVERTABLE_CONFIRM',
    okButton: 'BTN_DEPRECATE',
  };
  confirmDeprVisible.value = true;
};

const doDeletePolicyRule = async (config: IConfirmationDialogConfig) => {
  await AdminService.deletePolicyRule(config.key);
  info(t('DIALOG_policy_rule_delete_success'));
  await reload();
};

const doDeprecate = async (config: IConfirmationDialogConfig) => {
  await AdminService.deprecatePolicyRule(config.key);
  await reload();
  useSnackbar().info(t('DIALOG_pr_deprecate_success'));
};

const doCopy = async (config: IConfirmationDialogConfig) => {
  await AdminService.copyPolicyRule(config.key);
  await reload();
  useSnackbar().info(t('DIALOG_pr_copy_success'));
};

const downloadCsv = async () => {
  downloadFile(
    `licenses_and_policies_${dayjs(new Date()).format('YYYY-MM-DD_hh_mm_ss')}.csv`,
    AdminService.downloadLPcsv(),
    true,
  );
};

const filteredList = computed<PolicyRule[]>(() => {
  if (!Array.isArray(items.value)) {
    return [];
  }
  return items.value.filter((pr) => selectedFilterStatus.value.some((s) => s == pr.Status));
});

const downloadSingleCsv = async (id: string) => {
  downloadFile(
    `policy_${id}_${dayjs(new Date()).format('YYYY-MM-DD_hh_mm_ss')}.csv`,
    policyRuleService.downloadSingleLPcsv(id),
    true,
  );
};

const showCopyConfirmationDialog = (item: PolicyRule) => {
  confirmCopyConfig.value = {
    type: ConfirmationType.NOT_SET,
    key: item._key,
    name: item.Name,
    okButtonIsDisabled: false,
    okButton: 'BTN_COPY',
    description: 'DLG_CONFIRMATION_COPY_PR',
  } as IConfirmationDialogConfig;
  confirmCopyVisible.value = true;
};

const editPolicyRule = (item: PolicyRule) => {
  currentPolicyRuleForAction.value = item;
  policyRuleDialogRef.value?.showDialog();
};

const onPolicyRuleDialogClosed = async () => {
  currentPolicyRuleForAction.value = null;
  await reload();
};

const getActionButtons = (item: PolicyRule): TableActionButtonsProps['buttons'] => {
  const canManage = isPolicyManager.value && !item.Deprecated;

  return [
    {
      icon: 'mdi-pencil',
      hint: t('TT_edit_rule'),
      event: 'edit',
      show: canManage,
    },
    {
      icon: 'mdi-delete',
      hint: t('TT_delete_rule'),
      event: 'delete',
      show: canManage,
    },
    {
      icon: 'mdi-download',
      hint: t('TT_download_single_policy_csv'),
      event: 'download',
      show: true,
    },
    {
      icon: 'mdi-plus-circle-multiple-outline',
      hint: t('BTN_COPY'),
      event: 'copy',
      show: canManage,
    },
    {
      icon: 'mdi-archive-outline',
      hint: t('TT_deprecate_pr'),
      event: 'deprecate',
      show: canManage,
    },
  ];
};

const customFilterTable = (value: string, search: string) => {
  if (value != null && value) {
    const dateTime = formatDateAndTime(value);
    if (dateTime && dateTime !== 'Invalid date') {
      return dateTime.indexOf(search) > -1;
    }
    return value.toLowerCase().indexOf(search.toLowerCase()) > -1;
  }
  return false;
};
onMounted(async () => {
  isPolicyManager.value = RightsUtils.isPolicyManager();
  initBreadcrumbs();
  await reload();
});

const headers = computed(() => {
  return [
    {
      title: t('COL_ACTIONS'),
      align: 'center',
      filterable: true,
      class: 'tableHeaderCell',
      value: 'actions',
      width: 80,
      sortable: false,
    },
    {
      title: t('COL_STATUS'),
      align: 'start',
      class: 'tableHeaderCell',
      value: 'Status',
      width: 110,
      sortable: true,
    },
    {
      title: t('COL_NAME'),
      align: 'start',
      filterable: true,
      class: 'tableHeaderCell',
      value: 'Name',
      width: 300,
      sortable: true,
    },
    {
      title: t('DESCRIPTION'),
      align: 'start',
      filterable: false,
      class: 'tableHeaderCell',
      value: 'Description',
      width: 350,
      sortable: false,
    },
    {
      title: t('TOTAL'),
      align: 'center',
      filterable: false,
      class: 'tableHeaderCell',
      value: 'Total',
      width: 75,
      sortable: false,
    },
    {
      title: t('ALLOWED'),
      align: 'center',
      filterable: false,
      class: 'tableHeaderCell',
      value: 'Allowed',
      width: 75,
      sortable: false,
    },
    {
      title: t('WARNED'),
      align: 'center',
      filterable: false,
      class: 'tableHeaderCell',
      value: 'Warned',
      width: 75,
      sortable: false,
    },
    {
      title: t('DENIED'),
      align: 'center',
      filterable: false,
      class: 'tableHeaderCell',
      value: 'Denied',
      width: 75,
      sortable: false,
    },
    {
      title: t('CREATED'),
      align: 'center',
      filterable: true,
      class: 'tableHeaderCell',
      value: 'created',
      width: 108,
      sortable: true,
    },
    {
      title: t('UPDATED'),
      align: 'center',
      filterable: true,
      class: 'tableHeaderCell',
      value: 'updated',
      width: 108,
      sortable: true,
    },
  ] as DataTableHeader[];
});
</script>

<template>
  <TableLayout data-testid="policyrules">
    <template #buttons>
      <h1 class="text-h5">{{ t('POLICY_RULES') }}</h1>
      <NewPolicyRuleDialog v-slot="{showDialog}" :policy-labels="policyLabels" @reload="reload">
        <DCActionButton
          :text="t('BTN_ADD')"
          icon="mdi-plus"
          :hint="t('TT_add_rule')"
          @click="showDialog"
          v-if="isPolicyManager" />
      </NewPolicyRuleDialog>
      <v-spacer></v-spacer>
      <DCActionButton
        icon="mdi-download"
        :text="t('BTN_DOWNLOAD')"
        :hint="t('TT_download_policy_csv')"
        @click="downloadCsv"
        v-if="isPolicyManager" />
      <v-text-field
        autocomplete="off"
        :max-width="500"
        density="compact"
        variant="outlined"
        v-model="search"
        append-inner-icon="mdi-magnify"
        :label="t('labelSearch')"
        clearable
        single-line
        hide-details></v-text-field>
    </template>
    <template #table>
      <div ref="tableGridPolicyRules" class="fill-height">
        <v-data-table
          density="compact"
          class="striped-table fill-height"
          :headers="headers"
          fixed-header
          @click:row="onClickRow"
          :custom-filter="customFilterTable"
          :items="filteredList"
          :item-class="getCssClassForTableRow"
          :items-per-page="-1"
          :search="search"
          :sort-by="sortItems">
          <template v-slot:[`header.Status`]="{column, getSortIcon, toggleSort}">
            <div class="v-data-table-header__content">
              <span>{{ column.title }}</span>
              <v-menu offset-y :close-on-content-click="false" v-model="menu">
                <template v-slot:activator="{props}">
                  <DIconButton
                    :parentProps="props"
                    icon="mdi-filter-variant"
                    :hint="t('TT_SHOW_FILTER')"
                    :color="selectedFilterStatus && selectedFilterStatus.length > 0 ? 'primary' : 'secondary'" />
                </template>
                <div style="width: 320px" class="bg-background">
                  <v-card>
                    <v-row class="d-flex justify-end ma-1 mr-2">
                      <DCloseButton @click="menu = false" />
                    </v-row>

                    <v-select
                      variant="outlined"
                      v-model="selectedFilterStatus"
                      density="compact"
                      class="mx-2 pa-2"
                      autofocus
                      clearable
                      :items="possibleStatus"
                      :label="t('lbl_filter_on_status')"
                      hide-details
                      item-title="text"
                      item-value="value"
                      multiple
                      menu
                      transition="scale-transition"
                      persistent-clear
                      :list-props="{class: 'striped-filter-dd py-0'}">
                      <template v-slot:item="{props, item}">
                        <v-list-item v-bind="props" :title="undefined" class="py-0 px-2">
                          <template v-slot:prepend="{isSelected}">
                            <v-checkbox hide-details :model-value="isSelected" />
                          </template>
                          <span :class="'prStatus' + item.value + ' prStatusFilter'">{{ item.title }}</span>
                        </v-list-item>
                      </template>
                      <template v-slot:selection="{item, index}">
                        <div v-if="index === 0" class="d-flex align-center">
                          <span :class="'prStatus' + item.value + ' prStatusFilter'">{{ item.title }}</span>
                        </div>
                        <span v-if="index === 1" class="pAdditionalFilter">
                          +{{ selectedFilterStatus.length - 1 }} others
                        </span>
                      </template>
                    </v-select>
                  </v-card>
                </div>
              </v-menu>
              <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
            </div>
          </template>
          <template v-slot:[`item.Status`]="{item}">
            <span :class="'prStatus' + item.Status">{{ t('PR_STATUS_' + item.Status.toUpperCase()) }}</span>
          </template>
          <template v-slot:item.Allowed="{item}">
            {{ item.ComponentsAllow.length }}
          </template>
          <template v-slot:item.Warned="{item}">
            {{ item.ComponentsWarn.length }}
          </template>
          <template v-slot:item.Denied="{item}">
            {{ item.ComponentsDeny.length }}
          </template>
          <template v-slot:item.Total="{item}">
            {{ item.ComponentsAllow.length + item.ComponentsDeny.length + item.ComponentsWarn.length }}
          </template>
          <template v-slot:item.Description="{item}">
            <v-tooltip :text="item.Description" width="320" location="bottom" content-class="dpTooltip">
              <template v-slot:activator="{props}"
                ><span v-bind="props">{{ getStrWithMaxLength(95, item.Description) }}</span>
              </template></v-tooltip
            >
          </template>
          <template v-slot:item.updated="{item}">
            <DDateCellWithTooltip :value="item.updated" />
          </template>
          <template v-slot:item.created="{item}">
            <DDateCellWithTooltip :value="item.created" />
          </template>
          <template v-slot:item.actions="{item}">
            <TableActionButtons
              variant="compact"
              :buttons="getActionButtons(item)"
              @edit="editPolicyRule(item)"
              @delete="showDeletionConfirmationDialog(item as PolicyRule)"
              @download="downloadSingleCsv(item._key)"
              @copy="showCopyConfirmationDialog(item as PolicyRule)"
              @deprecate="showDeprecationConfirmationDialog(item as PolicyRule)" />
          </template>
        </v-data-table>
      </div>
    </template>
  </TableLayout>

  <NewPolicyRuleDialog
    v-if="currentPolicyRuleForAction"
    ref="policyRuleDialogRef"
    :policy-labels="policyLabels"
    :policy-rule="currentPolicyRuleForAction"
    @reload="onPolicyRuleDialogClosed" />
  <ConfirmationDialog v-model:showDialog="confirmVisible" :config="confirmConfig" @confirm="doDeletePolicyRule" />
  <ConfirmationDialog v-model:showDialog="confirmDeprVisible" :config="confirmDeprConfig" @confirm="doDeprecate" />
  <ConfirmationDialog v-model:showDialog="confirmCopyVisible" :config="confirmCopyConfig" @confirm="doCopy" />
</template>
