<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {ConfirmationType, IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import {IDefaultSelectItem} from '@disclosure-portal/model/IObligation';
import {OverallReview, OverallReviewState, VersionSlim} from '@disclosure-portal/model/VersionDetails';
import versionService from '@disclosure-portal/services/version';
import {useAppStore} from '@disclosure-portal/stores/app';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useSbomStore} from '@disclosure-portal/stores/sbom.store';
import {formatDateAndTime, getOverallReviewTranslationKey} from '@disclosure-portal/utils/Table';
import {getStrWithMaxLength, openUrl} from '@disclosure-portal/utils/View';
import TableActionButtons, {TableActionButtonsProps} from '@shared/components/TableActionButtons.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import {DataTableHeader, DataTableItem} from '@shared/types/table';
import {useClipboard} from '@shared/utils/clipboard';
import {TOOLTIP_OPEN_DELAY_IN_MS} from '@shared/utils/constant';
import dayjs from 'dayjs';
import _ from 'lodash';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRouter} from 'vue-router';
import {SortItem} from 'vuetify/lib/components/VDataTable/composables/sort';

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
const statusFilterOpened = ref(false);
const confirmConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);
const versionDialog = ref();
const confirmVisible = ref(false);

const maxVersions = 10;

const headers = computed((): DataTableHeader[] => {
  const res: DataTableHeader[] = [
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
  ];
  if (
    currentProject.value?.allowProjectRead ||
    currentProject.value?.allowProjectEdit ||
    currentProject.value?.allowProjectDelete
  ) {
    res.unshift({
      title: t('COL_ACTIONS'),
      key: 'actions',
      align: 'center',
      width: 100,
      sortable: false,
    });
  }
  return res;
});
const versions = computed((): VersionSlim[] => {
  return Object.values(currentProject.value?.versions ?? []);
});
const filteredList = computed((): VersionSlim[] => {
  return Array.isArray(versions.value) ? versions.value.filter(filterOnStatus) : [];
});
const possibleStatuses = computed((): IDefaultSelectItem[] => {
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
        } as IDefaultSelectItem;
      }
      const recentOverallReview = getRecentOverallReview(item.overallReviews);
      const status = !recentOverallReview || !recentOverallReview.state ? 'new' : item.status;
      return {
        text: t(getOverallReviewTranslationKey(status)),
        value: status,
      } as IDefaultSelectItem;
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

const getActionButtons = (): TableActionButtonsProps['buttons'] => {
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
};
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
      <v-text-field
        autocomplete="off"
        style="max-width: 500px"
        model="search"
        append-inner-icon="mdi-magnify"
        :label="t('labelSearch')"
        variant="outlined"
        clearable
        density="compact"
        hide-details />
    </template>
    <template #table>
      <div ref="tableVersions" class="fill-height">
        <v-data-table
          fixed-header
          :sort-by="updatedSort"
          :search="search"
          :headers="headers"
          :items="filteredList"
          @click:row="openVersion"
          :footer-props="{
            'items-per-page-options': [10, 50, 100, -1],
          }"
          class="striped-table fill-height"
          density="compact">
          <template v-slot:[`header.status`]="{column, getSortIcon, toggleSort}">
            <div class="v-data-table-header__content">
              <span>{{ column.title }}</span>
              <v-menu :close-on-content-click="false" v-model="statusFilterOpened">
                <template v-slot:activator="{props}">
                  <DIconButton
                    :parentProps="props"
                    icon="mdi-filter-variant"
                    :hint="t('TT_SHOW_FILTER')"
                    :color="selectedFilterStatus.length > 0 ? 'primary' : 'default'" />
                </template>
                <div class="bg-background" style="width: 280px">
                  <v-row class="d-flex ma-1 mr-2 justify-end">
                    <DIconButton icon="mdi-close" @clicked="statusFilterOpened = false" color="default" />
                  </v-row>
                  <v-select
                    v-model="selectedFilterStatus"
                    :items="possibleStatuses"
                    class="pa-2 mx-2 pb-4"
                    :label="t('Lbl_filter_status')"
                    clearable
                    multiple
                    item-title="text"
                    item-value="value"
                    variant="outlined"
                    density="compact"
                    menu
                    transition="scale-transition"
                    persistent-clear
                    :list-props="{class: 'striped-filter-dd py-0'}">
                    <template v-slot:item="{item, props}">
                      <v-list-item v-bind="props" class="px-2 py-0">
                        <template v-slot:prepend="{isSelected}">
                          <v-checkbox hide-details :model-value="isSelected" />
                        </template>
                        <template v-slot:title="{title}">
                          <span :class="'pvStatus' + (!item.value ? 'new' : item.value) + ' pStatusFilter'">
                            {{ title }}</span
                          >
                        </template>
                      </v-list-item>
                    </template>
                    <template v-slot:selection="{item, index}">
                      <div v-if="index === 0" class="d-flex align-center">
                        <span :class="'pvStatus' + (!item.value ? 'new' : item.value) + ' pStatusFilter'">{{
                          !item.value ? 'new' : item.title
                        }}</span>
                      </div>
                      <span v-if="index === 1" class="pAdditionalFilter">
                        +{{ selectedFilterStatus.length - 1 }} others
                      </span>
                    </template>
                  </v-select>
                </div>
              </v-menu>
              <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
            </div>
          </template>
          <template v-slot:[`item.status`]="{item}">
            <DVersionStateWithTooltip :version="item"></DVersionStateWithTooltip>
          </template>
          <template v-slot:[`item.created`]="{item}">
            <DDateCellWithTooltip :value="item.created" />
          </template>
          <template v-slot:[`item.updated`]="{item}">
            <DDateCellWithTooltip :value="item.updated" />
          </template>
          <template v-slot:[`item.description`]="{item}">
            <v-tooltip :open-delay="TOOLTIP_OPEN_DELAY_IN_MS" bottom max-width="480" content-class="dpTooltip">
              <template v-slot:activator="{props}">
                <span v-bind="props"> {{ getStrWithMaxLength(120, '' + item.description) }}</span>
              </template>
              {{ '' + item.description }}
            </v-tooltip>
          </template>
          <template v-slot:[`item.actions`]="{item}">
            <TableActionButtons
              variant="compact"
              :buttons="getActionButtons(item)"
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
