<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {ConfirmationType, IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import {OverallReview, OverallReviewState, VersionSlim} from '@disclosure-portal/model/VersionDetails';
import versionService from '@disclosure-portal/services/version';
import {useAppStore} from '@disclosure-portal/stores/app';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useSbomStore} from '@disclosure-portal/stores/sbom.store';
import {formatDateAndTime, getOverallReviewTranslationKey} from '@disclosure-portal/utils/Table';
import {getStrWithMaxLength, openUrl} from '@disclosure-portal/utils/View';
import {TableActionButtonsProps} from '@shared/components/TableActionButtons.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import {DataTableHeader, DataTableHeaderFilterItems, DataTableItem} from '@shared/types/table';
import {useClipboard} from '@shared/utils/clipboard';
import {TOOLTIP_OPEN_DELAY_IN_MS} from '@shared/utils/constant';
import dayjs from 'dayjs';
import _ from 'lodash';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRouter} from 'vue-router';
import {SortItem} from 'vuetify/lib/components/VDataTable/composables/sort';
import {getIconColor, getVersionStateIcon} from '@disclosure-portal/utils/Table';

const {t} = useI18n();
const {info} = useSnackbar();
const {copyToClipboard} = useClipboard();
const router = useRouter();
const appStore = useAppStore();
const projectStore = useProjectStore();
const sbomStore = useSbomStore();
const currentProject = computed(() => projectStore.currentProject!);
const labelTools = computed(() => appStore.getLabelsTools);

const search = ref('');
const selectedFilterStatus = ref<string[]>([]);
const confirmConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);
const versionDialog = ref();
const confirmVisible = ref(false);

const maxVersions = 10;

const showActions = computed(
  () =>
    currentProject.value?.allowProjectRead ||
    currentProject.value?.allowProjectEdit ||
    currentProject.value?.allowProjectDelete,
);

const headers = computed((): DataTableHeader[] => [
  ...(showActions.value
    ? [
        {
          title: t('COL_ACTIONS'),
          key: 'actions',
          align: 'center',
          width: 100,
          sortable: false,
        } as DataTableHeader,
      ]
    : []),
  {
    title: t('COL_STATUS'),
    sortable: true,
    align: 'center',
    width: 110,
    value: 'status',
  },
  {
    title: t('COL_VERSION'),
    key: 'name',
    width: 130,
    align: 'start',
  },
  {
    title: t('COL_DESCRIPTION'),
    key: 'description',
    align: 'start',
    width: 300,
    sortable: false,
  },
  {
    title: t('COL_CREATED'),
    key: 'created',
    align: 'start',
    width: 120,
  },
  {
    title: t('COL_UPDATED'),
    key: 'updated',
    align: 'start',
    width: 120,
  },
]);

const versions = computed((): VersionSlim[] => {
  return Object.values(currentProject.value?.versions ?? []);
});
const filteredList = computed((): VersionSlim[] => {
  return Array.isArray(versions.value) ? versions.value.filter(filterOnStatus) : [];
});
const possibleStatuses = computed((): DataTableHeaderFilterItems[] => {
  if (!versions.value) {
    return [];
  }
  return _.chain(versions.value)
    .uniqBy((item: VersionSlim) => {
      if (item.status == OverallReviewState.UNREVIEWED.toLowerCase()) {
        return item.status;
      }
      const recentOverallReview = getRecentOverallReview(item.overallReviews);
      return recentOverallReview === null || !recentOverallReview.state ? 'new' : item.status;
    })
    .map((item: VersionSlim) => {
      if (item.status == OverallReviewState.UNREVIEWED.toLowerCase()) {
        return {
          text: t(getOverallReviewTranslationKey(item.status)),
          value: item.status,
          iconColor: getIconColor(item.status),
          icon: getVersionStateIcon(item.status),
        } as DataTableHeaderFilterItems;
      }
      const recentOverallReview = getRecentOverallReview(item.overallReviews);
      const status = !recentOverallReview || !recentOverallReview.state ? 'new' : item.status;
      return {
        text: t(getOverallReviewTranslationKey(status)),
        value: status,
        iconColor: getIconColor(status),
        icon: getVersionStateIcon(status),
      } as DataTableHeaderFilterItems;
    })
    .value();
});
const maxVersionsReached = computed((): boolean => {
  return versions.value.length >= maxVersions;
});
const filterOnStatus = (item: VersionSlim) => {
  return (
    selectedFilterStatus.value.length === 0 || selectedFilterStatus.value.includes(!item.status ? 'new' : item.status)
  );
};
const getReferenceInfoForClipboard = (item: VersionSlim): string => {
  const schemaLabel = currentProject.value?.schemaLabel ?? '';
  const schemaLabelName = labelTools.value.schemaLabelsMap[schemaLabel]
    ? labelTools.value.schemaLabelsMap[schemaLabel].name
    : 'UNKNOWN_LABEL';
  const policyLabelNames = currentProject.value?.policyLabels
    .map((l: string) =>
      labelTools.value.policyLabelsMap[l] ? labelTools.value.policyLabelsMap[l].name : 'UNKNOWN_LABEL',
    )
    .join(', ');
  return `Disclosure Portal Project Version Reference

Project Name: ${currentProject.value?.name}
Project Identifier: ${currentProject.value?._key}
Project Schema Label: ${schemaLabelName}
Project Policy Labels: ${policyLabelNames}
Project Version: ${item.name}
Version Identifier: ${item._key}
Reference Timestamp: ${formatDateAndTime(dayjs().toISOString())} (server time)
Deliveries Link: ${window.location.origin}/#/dashboard/projects/${encodeURIComponent(currentProject.value?._key)}/versions/${encodeURIComponent(item._key)}/component`;
};
const getRecentOverallReview = (overallReviews: OverallReview[]): OverallReview | null => {
  if (!overallReviews || overallReviews.length === 0) {
    return null;
  }
  return overallReviews.reduce((recent, current) => {
    return dayjs(current.updated).isAfter(dayjs(recent.updated)) ? current : recent;
  }, overallReviews[0]);
};

const updatedSort: SortItem[] = [{key: 'updated', order: 'desc'}];

const openVersion = (event: Event, item: DataTableItem<VersionSlim>) => {
  const url = `/dashboard/projects/${encodeURIComponent(currentProject.value._key)}/versions/${encodeURIComponent(item.item._key)}`;
  openUrl(url, router);
};

const doDelete = async (config: IConfirmationDialogConfig) => {
  if (config.okButtonIsDisabled) return;
  await versionService.deleteVersion(currentProject.value._key, config.key);
  info(t('DIALOG_version_delete_success'));
  await sbomStore.fetchAllSBOMsFlat(true);
  await projectStore.fetchProjectByKey(currentProject.value._key);
};
const showConfirm = async (item: VersionSlim) => {
  await versionService.getApprovalOrReviewUsage(currentProject.value?._key, item._key).then((r) => {
    const isInUse = r.data.success;
    if (isInUse) {
      confirmConfig.value = {
        type: ConfirmationType.NOT_SET,
        title: 'DLG_WARNING_TITLE',
        key: '',
        name: '',
        description: 'VERSION_IN_APPROVAL',
        okButton: 'Btn_delete',
        okButtonIsDisabled: true,
      };
    } else {
      confirmConfig.value = {
        type: ConfirmationType.DELETE,
        key: item._key,
        name: item.name,
        description: 'DLG_CONFIRMATION_DESCRIPTION',
        okButton: 'Btn_delete',
      };
    }
    confirmVisible.value = true;
  });
};

const copyVersionToClipboard = (item: VersionSlim) => {
  const content = getReferenceInfoForClipboard(item);
  copyToClipboard(content);
};

const editVersion = (item: VersionSlim) => {
  versionDialog.value?.open({
    version: item,
  });
};

const actionButtons = computed((): TableActionButtonsProps['buttons'] => {
  return [
    {
      icon: 'mdi-content-copy',
      hint: t('TT_COPY_REFERENCE_INFO'),
      event: 'copy',
      show: currentProject.value?.allowProjectRead,
    },
    {
      icon: 'mdi-pencil',
      hint: t('TT_edit_version'),
      event: 'edit',
      show: currentProject.value?.allowProjectEdit,
    },
    {
      icon: 'mdi-delete',
      hint: t('TT_delete_version'),
      event: 'delete',
      show: currentProject.value?.allowProjectDelete,
    },
  ];
});
</script>

<template>
  <TableLayout has-tab has-title>
    <template #description v-if="$slots.default">
      <slot></slot>
    </template>
    <template #buttons>
      <DCActionButton
        :disabled="maxVersionsReached"
        :text="t('BTN_ADD')"
        icon="mdi-plus"
        :hint="maxVersionsReached ? t('TT_no_new_version') : t('TT_new_version')"
        @clicked="versionDialog?.open({})"
        v-if="
          currentProject && currentProject.accessRights && currentProject.accessRights.allowProjectVersion.create
        " />
      <v-spacer></v-spacer>
      <DSearchField v-model="search" />
    </template>
    <template #table>
      <div ref="tableVersions" class="fill-height">
        <v-data-table
          fixed-header
          :sort-by="updatedSort"
          :search="search"
          :headers="headers"
          :items="filteredList"
          class="striped-table fill-height"
          density="compact"
          :footer-props="{
            'items-per-page-options': [10, 50, 100, -1],
          }"
          @click:row="openVersion">
          <template #[`header.status`]="{column, getSortIcon, toggleSort}">
            <GridFilterHeader :column="column" :getSortIcon="getSortIcon" :toggleSort="toggleSort">
              <template #filter>
                <GridHeaderFilterIcon
                  v-model="selectedFilterStatus"
                  :column="column"
                  :label="t('COL_STATUS')"
                  :allItems="possibleStatuses">
                </GridHeaderFilterIcon>
              </template>
            </GridFilterHeader>
          </template>
          <template #[`item.status`]="{item}">
            <DVersionStateWithTooltip :version="item"></DVersionStateWithTooltip>
          </template>
          <template #[`item.created`]="{item}">
            <DDateCellWithTooltip :value="item.created" />
          </template>
          <template #[`item.updated`]="{item}">
            <DDateCellWithTooltip :value="item.updated" />
          </template>
          <template #[`item.description`]="{item}">
            <v-tooltip :open-delay="TOOLTIP_OPEN_DELAY_IN_MS" bottom max-width="480" content-class="dpTooltip">
              <template #activator="{props}">
                <span v-bind="props"> {{ getStrWithMaxLength(120, '' + item.description) }}</span>
              </template>
              {{ '' + item.description }}
            </v-tooltip>
          </template>
          <template #[`item.actions`]="{item}">
            <TableActionButtons
              variant="compact"
              :buttons="actionButtons"
              @copy="copyVersionToClipboard(item)"
              @edit="editVersion(item)"
              @delete="showConfirm(item)" />
          </template>
        </v-data-table>
      </div>
    </template>
  </TableLayout>

  <VersionDialogForm ref="versionDialog"></VersionDialogForm>
  <ConfirmationDialog
    v-model:showDialog="confirmVisible"
    :config="confirmConfig"
    @confirm="doDelete"></ConfirmationDialog>
</template>
