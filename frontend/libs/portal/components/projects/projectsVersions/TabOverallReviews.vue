<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {OverallReview, OverallReviewState} from '@disclosure-portal/model/VersionDetails';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useSbomStore} from '@disclosure-portal/stores/sbom.store';
import {RightsUtils} from '@disclosure-portal/utils/Rights';
import {formatDateAndTime, getIconColor, getVersionStateIcon} from '@disclosure-portal/utils/Table';
import {DataTableHeader, SortItem} from '@shared/types/table';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const projectStore = useProjectStore();
const sbomStore = useSbomStore();

const sortBy: SortItem[] = [{key: 'created', order: 'desc'}];

const overallReviewDialog = ref();
const overallAuditDialog = ref();

const currentProject = computed(() => projectStore.currentProject!);
const spdxHistory = computed(() => sbomStore.getChannelSpdxs);

const canAddAudit = computed(() => {
  const hasPermission = RightsUtils.isFOSSOffice();
  return spdxHistory.value.length > 0 && hasPermission && !currentProject.value.isDeprecated;
});

const headers = computed<DataTableHeader[]>(() => [
  {
    title: t('COL_STATUS'),
    align: 'center',
    value: 'state',
    sortable: true,
    sort: compareStatus,
    width: 100,
  },
  {
    title: t('COL_SPDX_FILENAME'),
    align: 'start',
    value: 'sbom',
    width: 250,
  },
  {
    title: t('COL_COMMENT'),
    align: 'start',
    value: 'comment',
    width: 200,
    sortable: true,
  },
  {
    title: t('COL_CREATOR'),
    align: 'start',
    value: 'creator',
    width: 200,
    sortable: true,
  },
  {
    title: t('COL_CREATED'),
    align: 'start',
    value: 'created',
    width: 150,
    sortable: true,
  },
]);
const items = computed<OverallReview[]>(
  () => currentProject.value.versions[sbomStore.currentVersionKey]?.overallReviews ?? [],
);
const enumToLowerCase = (overallReviewState: OverallReviewState): string =>
  overallReviewState ? OverallReviewState[overallReviewState].toLowerCase() : '';

const compareStatus = (a: OverallReviewState, b: OverallReviewState): number => {
  const levelWeight: Map<OverallReviewState, number> = new Map<OverallReviewState, number>([
    [OverallReviewState.AUDITED, 0],
    [OverallReviewState.ACCEPTABLE, 1],
    [OverallReviewState.ACCEPTABLE_AFTER_CHANGES, 2],
    [OverallReviewState.UNREVIEWED, 3],
    [OverallReviewState.NOT_ACCEPTABLE, 4],
  ]);
  return levelWeight.get(a)! - levelWeight.get(b)!;
};

const showOverallReviewDialog = async () => {
  overallReviewDialog.value.open();
};

const showOverallAuditDialog = async () => {
  overallAuditDialog.value.open();
};
</script>

<template>
  <TableLayout has-tab has-title>
    <template #buttons>
      <DCActionButton
        v-if="
          spdxHistory.length > 0 && currentProject.accessRights.allowProjectVersion.read && !currentProject.isDeprecated
        "
        large
        :text="t('BTN_ADD')"
        icon="mdi-plus"
        :hint="t('TT_overall_review')"
        @click="showOverallReviewDialog"></DCActionButton>
      <DCActionButton
        v-if="canAddAudit"
        large
        :text="t('BTN_ADD_AUDIT')"
        icon="mdi-plus"
        :hint="t('TT_overall_audit')"
        @click="showOverallAuditDialog"></DCActionButton>
    </template>
    <template #table>
      <v-data-table
        density="compact"
        class="striped-table custom-data-table h-full"
        :headers="headers"
        fixed-header
        item-key="updated"
        :sort-by="sortBy"
        sort-desc
        :items="items"
        :footer-props="{'items-per-page-options': [10, 50, 100, -1]}">
        <template v-slot:[`item.state`]="{item}">
          <v-icon :color="getIconColor(enumToLowerCase(item.state))" small>
            {{ getVersionStateIcon(enumToLowerCase(item.state)) }}
          </v-icon>
        </template>
        <template v-slot:[`item.comment`]="{item}">
          <Truncated>{{ item.comment }}</Truncated>
        </template>
        <template v-slot:[`item.sbom`]="{item}">
          {{ `${formatDateAndTime(item.sbomUploaded)} - ${item.sbomName}` }}
        </template>
        <template v-slot:[`item.creator`]="{item}">
          {{ `${item.creatorFullName} (${item.creator})` }}
        </template>
        <template v-slot:[`item.created`]="{item}">
          <DDateCellWithTooltip :value="item.created"></DDateCellWithTooltip>
        </template>
      </v-data-table>
    </template>
  </TableLayout>
  <OverallReviewDialog ref="overallReviewDialog" visible></OverallReviewDialog>
  <OverallAuditDialog ref="overallAuditDialog" visible></OverallAuditDialog>
</template>
