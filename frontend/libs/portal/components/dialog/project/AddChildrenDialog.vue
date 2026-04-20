<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import Icons from '@disclosure-portal/constants/icons';
import ProjectPostRequest from '@disclosure-portal/model/ProjectPostRequest';
import {ProjectSlim} from '@disclosure-portal/model/ProjectsResponse';
import {useIdleStore} from '@disclosure-portal/stores/idle.store';
import {useLabelStore} from '@disclosure-portal/stores/label.store';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {formatDateAndTime, getCssClassForTableRow} from '@disclosure-portal/utils/Table';
import DLabel from '@shared/components/disco/DLabel.vue';
import {DataTableHeader} from '@shared/types/table';
import {computed, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {InternalItem} from 'vuetify/framework';

const {t} = useI18n();
const projectStore = useProjectStore();
const labelStore = useLabelStore();
const idleStore = useIdleStore();

const show = ref(false);
const search = ref('');

const icons = Icons;
const headers = ref<DataTableHeader[]>([
  {title: t('COL_STATUS'), sortable: true, filterable: true, class: 'tableHeaderCell', value: 'status', width: '80'},
  {title: t('COL_NAME'), align: 'start', class: 'tableHeaderCell', value: 'name', width: 240},
  {title: t('COL_LABELS'), align: 'start', class: 'tableHeaderCell', value: 'labels', width: 340, sortable: false},
]);

const selectedInGrid = ref<ProjectSlim[]>([]);
const keepSelected = ref<string[]>([]);
const allSelected = ref(false);

const projectsForSelection = computed(() =>
  (projectStore.projects || [])
    .filter((project: ProjectSlim) => project._key !== projectStore.currentProject?._key)
    .map((project) => ({
      ...project,
      _policyLabels: project.policyLabels.map((labelKey) => labelStore.getLabelByKey(labelKey)),
      _schemaLabel: labelStore.getLabelByKey(project.schemaLabel),
    })),
);
const indeterminate = computed(
  () => selectedInGrid.value.length > 0 && selectedInGrid.value.length < projectsForSelection.value.length,
);

const customFilterTable = (rawCellValue: string, searchTextRaw: string, internalItem?: InternalItem<ProjectSlim>) => {
  const row = internalItem!.raw;

  if (rawCellValue) {
    const dateTime = formatDateAndTime(rawCellValue);

    if (dateTime !== 'Invalid date') {
      return dateTime.indexOf(searchTextRaw) > -1;
    }

    const searchText = searchTextRaw.toLowerCase();

    const cellValueFound = rawCellValue.toLowerCase().includes(searchText);

    if (cellValueFound) {
      return cellValueFound;
    }

    const freeLabelFound = row.freeLabels.some((label) => label.toLowerCase().includes(searchText));

    if (freeLabelFound) {
      return freeLabelFound;
    }

    const schemaLabel = labelStore.getLabelByKey(row.schemaLabel);
    const schemaLabelFound = (schemaLabel && schemaLabel.name?.toLowerCase().includes(searchText)) || false;

    if (schemaLabelFound) {
      return schemaLabelFound;
    }

    return row.policyLabels.some((labelKey) => {
      const label = labelStore.getLabelByKey(labelKey);
      return label && label.name?.toLowerCase().includes(searchText);
    });
  }

  return false;
};

const updateAllSelectedState = () => {
  const selectable = projectsForSelection.value.filter((it) => !it.isInGroupApproval);
  const selected = selectedInGrid.value.filter((it) => !it.isInGroupApproval).length;
  allSelected.value = selectable.length > 0 && selected === selectable.length;
};

const evaluateHeaderCheckboxState = () => {
  const total = projectsForSelection.value.length;
  const selected = selectedInGrid.value.length;
  allSelected.value = selected === total && total > 0;
};

const isSelected = (p: ProjectSlim) => selectedInGrid.value.some((s) => s && s._key === p._key);

const toggleSelection = (p: ProjectSlim, value: boolean) => {
  if (p.isInGroupApproval && isSelected(p)) return;
  if (value) {
    if (!isSelected(p)) selectedInGrid.value.push(p);
  } else {
    const index = selectedInGrid.value.findIndex((s) => s && s._key === p._key);
    if (index >= 0) selectedInGrid.value.splice(index, 1);
  }
  updateAllSelectedState();
};

const toggleAllSelection = (value: boolean) => {
  allSelected.value = value;
  selectedInGrid.value = value
    ? projectsForSelection.value.slice()
    : projectsForSelection.value.filter((it) => it.isInGroupApproval);
};

const preselectItems = () => {
  const allItems = projectsForSelection.value;
  const preSelected = [...(projectStore.currentProject!.children || [])];
  const idxArray = preSelected.map((entry) => allItems.findIndex((e) => e._key === entry));
  selectedInGrid.value = idxArray.map((idx) => allItems[idx]).filter(Boolean) as ProjectSlim[];
  keepSelected.value = preSelected.filter((key) => !allItems.some((e) => e._key === key));
  evaluateHeaderCheckboxState();
};

const getSelected = () => {
  selectedInGrid.value = selectedInGrid.value.filter(Boolean);
  const keys = selectedInGrid.value.map((e) => e._key);
  return Array.from(new Set([...keys, ...keepSelected.value]));
};

const open = async () => {
  show.value = true;
  await projectStore.fetchProjectPossibleChildren(projectStore.currentProject!._key);
  preselectItems();
};

const save = async () => {
  idleStore.show();

  const projectPostRequest = new ProjectPostRequest();
  projectPostRequest.fillWithProjectModel(projectStore.currentProject!);
  projectPostRequest.projectSettings = null;
  projectPostRequest.children = getSelected();
  try {
    await projectStore.updateProject(projectPostRequest);
  } finally {
    idleStore.hide();
  }
  close();
};

const close = () => {
  show.value = false;
};

defineExpose({open, close});

watch(selectedInGrid, evaluateHeaderCheckboxState);
</script>

<template>
  <v-dialog v-model="show" scrollable persistent width="1200">
    <v-card class="pa-8 dDialog">
      <v-card-title>
        <Stack direction="row">
          <span class="text-h5 d-headline">{{ t('PROJECT_GROUP_DEFINITION') }}</span>
          <span class="grow"></span>
          <DCloseButton @click="close" />
        </Stack>
      </v-card-title>

      <v-card-text>
        <Stack>
          <Stack direction="row">
            <div class="grow"></div>
            <DSearchField v-model="search" />
          </Stack>

          <v-data-table
            :headers="headers"
            show-select
            fixed-header
            hide-default-footer
            density="compact"
            :items-per-page="-1"
            :footer-props="{'items-per-page-options': [10, 50, 100, -1]}"
            :sort-by="[{key: 'updated', order: 'desc'}]"
            :items="projectsForSelection"
            height="380"
            class="striped-table custom-data-table"
            item-key="_key"
            item-value="_key"
            v-model="selectedInGrid"
            :custom-filter="customFilterTable"
            :item-class="getCssClassForTableRow"
            :search="search">
            <template v-slot:[`header.data-table-select`]>
              <v-checkbox-btn
                :model-value="allSelected"
                :indeterminate="indeterminate"
                @update:model-value="toggleAllSelection" />
            </template>
            <template v-slot:[`item.data-table-select`]="{item: row}">
              <v-checkbox-btn
                :model-value="isSelected(row)"
                :disabled="row.isInGroupApproval"
                @update:modelValue="(value: boolean) => toggleSelection(row, value)" />
              <Tooltip v-if="row.isInGroupApproval">
                <span>{{ t('CHILD_PROJECT_UNDER_APPROVE') }}</span>
              </Tooltip>
            </template>
            <template v-slot:[`item.status`]="{item: row}">
              <span :class="'pStatus' + (!row.status ? 'new' : row.status)">{{
                !row.status ? 'new' : row.status
              }}</span>
            </template>
            <template v-slot:[`item.labels`]="{item: row}">
              <Stack class="py-2">
                <ProjectLabel :label="row._schemaLabel"></ProjectLabel>
                <span v-for="(l, i) in row.freeLabels" :key="'b' + i">
                  <DLabel :labelName="l" :iconName="icons.TAG" />
                  <Tooltip>
                    <span>{{ t('TT_free_label') }}</span>
                  </Tooltip>
                </span>

                <div class="m-[-4px]">
                  <ProjectLabel
                    v-for="(label, i) in row._policyLabels"
                    :key="'a' + i"
                    :label="label"
                    class="m-1 inline-block"></ProjectLabel>
                </div>
              </Stack>
            </template>
          </v-data-table>
        </Stack>
      </v-card-text>

      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn @click="close" plain class="secondary mr-8" color="primary" size="small">{{ t('BTN_CANCEL') }}</v-btn>
        <v-btn @click="save" color="primary" size="small" variant="flat">{{ t('Btn_save') }}</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
