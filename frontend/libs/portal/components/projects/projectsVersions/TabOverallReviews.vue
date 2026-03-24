<script setup lang="ts">
import useDimensions from '@disclosure-portal/composables/useDimensions';
import {PolicyLabels} from '@disclosure-portal/constants/policyLabels';
import {OverallReview, OverallReviewState} from '@disclosure-portal/model/VersionDetails';
import {useLabelStore} from '@disclosure-portal/stores/label.store';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useSbomStore} from '@disclosure-portal/stores/sbom.store';
import eventBus from '@disclosure-portal/utils/eventbus';
import {RightsUtils} from '@disclosure-portal/utils/Rights';
import {formatDateAndTime, getIconColor, getVersionStateIcon} from '@disclosure-portal/utils/Table';
import {DataTableHeader, SortItem} from '@shared/types/table';
import {computed, nextTick, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const labelStore = useLabelStore();
const projectStore = useProjectStore();
const sbomStore = useSbomStore();
const {calculateHeight} = useDimensions();

const sortBy: SortItem[] = [{key: 'created', order: 'desc'}];

const tableHeight = ref(110);
const tableGrid = ref<HTMLElement | null>(null);
const overallReviewDialog = ref();
const overallAuditDialog = ref();

const currentProject = computed(() => projectStore.currentProject!);
const spdxHistory = computed(() => sbomStore.getChannelSpdxs);
const selectedSpdx = computed(() => sbomStore.getSelectedSpdx);
const version = computed(() => sbomStore.getCurrentVersion);
const hasVehiclePlatformChildren = computed(() => projectStore.hasVehiclePlatformChildren);

const isVehiclePlatform = computed(() => {
  if (!currentProject.value) return false;
  for (const lbl of currentProject.value.policyLabels) {
    if (labelStore.getLabelByKey(lbl)?.name === PolicyLabels.VEHICLE_PLATFORM) {
      return true;
    }
  }
  return false;
});

const canAddAudit = computed(() => {
  const hasPermission = RightsUtils.isFOSSOffice();
  return spdxHistory.value.length > 0 && hasPermission && !currentProject.value.isDeprecated;
});

const headers = computed<DataTableHeader[]>(() => [
  {
    title: t('COL_STATUS'),
    align: 'center',
    class: 'tableHeaderCell',
    value: 'state',
    sortable: true,
    sort: compareStatus,
    width: 100,
  },
  {
    title: t('COL_SPDX_FILENAME'),
    align: 'start',
    class: 'tableHeaderCell',
    value: 'sbom',
    width: 250,
  },
  {
    title: t('COL_COMMENT'),
    align: 'start',
    class: 'tableHeaderCell',
    value: 'comment',
    width: 200,
    sortable: true,
  },
  {
    title: t('COL_CREATOR'),
    align: 'start',
    class: 'tableHeaderCell',
    value: 'creator',
    width: 200,
    sortable: true,
  },
  {
    title: t('COL_CREATED'),
    align: 'start',
    class: 'tableHeaderCell',
    value: 'created',
    width: 150,
    sortable: true,
  },
]);
const items = computed<OverallReview[]>(() =>
  !version.value.overallReviews || version.value.overallReviews.length === 0 ? [] : version.value.overallReviews,
);
const enumToLowerCase = (overallReviewState: OverallReviewState): string =>
  overallReviewState ? OverallReviewState[overallReviewState].toLowerCase() : '';

const updateTableHeight = async () => {
  await nextTick();
  if (tableGrid.value) {
    tableHeight.value = calculateHeight(tableGrid.value, false, true, ['divIntroText']);
  }
};

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
  overallReviewDialog.value.open(
    currentProject.value._key,
    version.value._key,
    spdxHistory.value,
    selectedSpdx.value,
    currentProject.value.approvablespdx.spdxkey,
  );
};

const showOverallAuditDialog = async () => {
  overallAuditDialog.value.open(
    currentProject.value._key,
    version.value._key,
    spdxHistory.value,
    selectedSpdx.value,
    currentProject.value.approvablespdx.spdxkey,
  );
};

onMounted(async () => {
  eventBus.on('window-resize', updateTableHeight);
  await updateTableHeight();
});
</script>

<template>
  <v-container fluid v-if="currentProject">
    <v-row justify="space-between" class="measure-height">
      <v-col xs="12" sm="12" md="6" lg="5" class="pa-1">
        <DCActionButton
          v-if="
            spdxHistory.length > 0 &&
            currentProject.accessRights.allowProjectVersion.read &&
            !currentProject.isDeprecated
          "
          large
          class="mx-2"
          :text="t('BTN_ADD')"
          icon="mdi-plus"
          :hint="t('TT_overall_review')"
          @click="showOverallReviewDialog"></DCActionButton>
        <DCActionButton
          v-if="canAddAudit"
          large
          class="mx-2"
          :text="t('BTN_ADD_AUDIT')"
          icon="mdi-plus"
          :hint="t('TT_overall_audit')"
          @click="showOverallAuditDialog"></DCActionButton>
      </v-col>
      <v-col cols="12" class="pb-0">
        <div ref="tableGrid">
          <v-data-table
            density="compact"
            class="striped-table custom-data-table"
            :headers="headers"
            fixed-header
            item-key="updated"
            :sort-by="sortBy"
            :height="tableHeight"
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
        </div>
      </v-col>
    </v-row>
    <OverallReviewDialog ref="overallReviewDialog" visible></OverallReviewDialog>
    <OverallAuditDialog ref="overallAuditDialog" visible></OverallAuditDialog>
  </v-container>
</template>
