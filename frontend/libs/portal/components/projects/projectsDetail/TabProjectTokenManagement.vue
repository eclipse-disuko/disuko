<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {ConfirmationType, IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import {Token} from '@disclosure-portal/model/Project';
import projectService from '@disclosure-portal/services/projects';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {getCssClassForTableRow} from '@disclosure-portal/utils/Table';
import {TableActionButtonsProps} from '@shared/components/TableActionButtons.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import TableLayout from '@shared/layouts/TableLayout.vue';
import {DataTableHeader, DataTableHeaderFilterItems, SortItem} from '@shared/types/table';
import _ from 'lodash';
import {computed, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const projectStore = useProjectStore();
const {info} = useSnackbar();

const tokenGrid = ref<HTMLElement | null>(null);
const tokens = ref<Token[]>([]);
const dataAreLoaded = ref(false);
const selectedFilterStatus = ref<string[]>([]);
const create = ref(false);
const update = ref(false);
const del = ref(false);
const confirmationDialogConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);
const confirmationDialogVisible = ref(false);
const tokenIssuedDialogVisible = ref(false);
const sortBy: SortItem[] = [{key: 'created', order: 'desc'}];
const tokenRef = ref<Token>({} as Token);
const renewed = ref(false);

const tokenHeaders = computed((): DataTableHeader[] => [
  {
    title: t('COL_ACTIONS'),
    align: 'center',
    value: 'actions',
    width: 120,
    sortable: false,
  },
  {
    title: t('TOKEN_NAME'),
    align: 'start',
    sortable: true,
    value: 'company',
    width: 160,
  },
  {
    title: t('COL_DESCRIPTION'),
    align: 'start',
    sortable: true,
    width: 180,
    value: 'description',
  },
  {
    title: t('COL_STATUS'),
    align: 'start',
    sortable: true,
    value: 'status',
    width: 160,
  },
  {
    title: t('COL_CREATED'),
    align: 'start',
    sortable: true,
    value: 'created',
    width: 160,
  },
  {
    title: t('COL_EXPIRY'),
    align: 'start',
    sortable: true,
    value: 'expiry',
    width: 160,
  },
]);

const projectModel = computed(() => projectStore.currentProject!);
const filteredList = computed(() => _.filter(tokens.value, filterOnStatus));

function filterOnStatus(item: Token): boolean {
  return selectedFilterStatus.value.length === 0 || _.includes(selectedFilterStatus.value, item.status);
}

const possibleStatus = computed((): DataTableHeaderFilterItems[] => {
  if (!tokens.value) {
    return [];
  }

  const uniqueStatuses = [...new Set(tokens.value.map(({status}) => status))];

  return uniqueStatuses.map((value: string) => {
    return {
      value,
    } as DataTableHeaderFilterItems;
  });
});

onMounted(async () => {
  await reload();
});

async function reload() {
  dataAreLoaded.value = false;
  create.value =
    projectModel.value.accessRights.allowProjectTokenManagement.update ||
    projectModel.value.accessRights.allowAllProjectTokenManagement.update;
  update.value =
    projectModel.value.accessRights.allowProjectTokenManagement.update ||
    projectModel.value.accessRights.allowAllProjectTokenManagement.update;
  del.value =
    projectModel.value.accessRights.allowProjectTokenManagement.delete ||
    projectModel.value.accessRights.allowAllProjectTokenManagement.delete;
  if (projectModel.value._key) {
    projectService.getTokens(projectModel.value._key).then((response) => {
      tokens.value = response;
      dataAreLoaded.value = true;
    });
  }
}

function renew(item: Token) {
  confirmationDialogConfig.value = {
    type: ConfirmationType.RENEW,
    key: item._key,
    name: item.description,
    description: 'DLG_CONFIRMATION_DESCRIPTION_RENEW',
    okButton: 'Btn_renew',
    okButtonIsDisabled: false,
  };
  confirmationDialogVisible.value = true;
}

function revoke(item: Token) {
  confirmationDialogConfig.value = {
    type: ConfirmationType.REVOKE,
    key: item._key,
    name: item.description,
    description: 'DLG_CONFIRMATION_DESCRIPTION_REVOKE',
    okButton: 'Btn_revoke',
    okButtonIsDisabled: false,
  };
  confirmationDialogVisible.value = true;
}

async function doRenewOrRevoke(config: IConfirmationDialogConfig) {
  if (config.type === ConfirmationType.REVOKE) {
    await projectService.revokeProjectToken(projectModel.value._key, config.key);
    info(t('DIALOG_token_revoke_success'));
    await reload();
    return;
  }
  if (config.type === ConfirmationType.RENEW) {
    tokenRef.value = (await projectService.renewProjectToken(projectModel.value._key, config.key)).data;
    renewed.value = true;
    tokenIssuedDialogVisible.value = true;
    await reload();
    return;
  }
}

async function onCreated(token: Token) {
  tokenRef.value = token;
  renewed.value = false;
  tokenIssuedDialogVisible.value = true;
  await reload();
}

const getActionButtons = (item: Token): TableActionButtonsProps['buttons'] => {
  const isActive = item?.status === 'active';
  const isNotDeprecated = !projectModel.value.isDeprecated;

  return [
    {
      icon: 'mdi-replay',
      hint: t('TT_renew_token'),
      event: 'renew',
      show: isActive && update.value && isNotDeprecated,
    },
    {
      icon: 'mdi-close',
      hint: t('TT_revoke_token'),
      event: 'revoke',
      show: isActive && del.value && isNotDeprecated,
    },
  ];
};
</script>

<template>
  <TableLayout has-title has-tab>
    <template #buttons>
      <CreateTokenDialog v-slot="{showDialog}" @onCreated="onCreated">
        <DCActionButton
          :text="t('BTN_ADD')"
          icon="mdi-plus"
          :hint="t('TT_new_token')"
          @click="showDialog"
          v-if="projectModel && create && !projectModel.isDeprecated" />
      </CreateTokenDialog>
    </template>
    <template #table>
      <div ref="tokenGrid" class="fill-height">
        <v-data-table
          density="compact"
          :loading="!dataAreLoaded"
          fixed-header
          :items-per-page="15"
          :sort-by="sortBy"
          sort-desc
          class="striped-table custom-data-table fill-height"
          :headers="tokenHeaders"
          :items="filteredList"
          :item-class="getCssClassForTableRow">
          <template #[`header.status`]="{column, getSortIcon, toggleSort}">
            <GridFilterHeader :column="column" :getSortIcon="getSortIcon" :toggleSort="toggleSort">
              <template #filter>
                <GridHeaderFilterIcon
                  v-model="selectedFilterStatus"
                  :column="column"
                  :label="t('COL_STATUS')"
                  :allItems="possibleStatus">
                </GridHeaderFilterIcon>
              </template>
            </GridFilterHeader>
          </template>
          <template #[`item.created`]="{item}">
            <DDateCellWithTooltip :value="item.created" />
          </template>
          <template #[`item.expiry`]="{item}">
            <DDateCellWithTooltip :value="item.expiry" />
          </template>
          <template #[`item.actions`]="{item}">
            <TableActionButtons
              variant="compact"
              :buttons="getActionButtons(item)"
              @renew="renew(item)"
              @revoke="revoke(item)" />
          </template>
        </v-data-table>
      </div>
    </template>
  </TableLayout>

  <TokenIssuedDialog
    v-model:showDialog="tokenIssuedDialogVisible"
    :token="tokenRef"
    :renewed="renewed"></TokenIssuedDialog>
  <ConfirmationDialog
    v-model:showDialog="confirmationDialogVisible"
    :config="confirmationDialogConfig"
    @confirm="doRenewOrRevoke" />
</template>
