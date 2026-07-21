<script setup lang="ts">
import {Version} from '@cli/models/Project';
import {projectService} from '@cli/services/projectService';
import {useAppStore} from '@cli/stores/app';
import useSnackbar from '@shared/composables/useSnackbar';
import TableActionButtons, {TableActionButtonsProps} from '@shared/components/TableActionButtons.vue';
import {useClipboard} from '@shared/utils/clipboard';
import {computed, ref, ref as vueRef} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRouter} from 'vue-router';
import {getOverallReviewTranslationKey} from '@disclosure-portal/utils/Table';

const maxVersions = 10;
const appStore = useAppStore();
const router = useRouter();

const statusClasses = (status: string) => ({
  [`pvStatus${!status ? 'new' : status}`]: true,
  pStatusFilter: true,
});

const {t} = useI18n();
const {info: snackbar} = useSnackbar();
const {copyToClipboard} = useClipboard();
const search = ref('');
const selectedFilterStatus = ref<string[]>([]);
const statusFilterOpened = ref(false);
const updatedSort = [{key: 'lastSbomUploaded', order: 'desc' as const}];

const headers = computed(() => [
  {title: t('COL_ACTIONS'), key: 'actions', align: 'center' as const, width: 50},
  {title: t('COL_STATUS'), key: 'status', align: 'center' as const, width: 80},
  {title: t('COL_CHANNEL'), key: 'name', align: 'start' as const, width: 130},
  {title: t('COL_DESCRIPTION'), key: 'description', align: 'start' as const, width: 130},
  {title: t('COL_LAST_SBOM_UPLOADED'), key: 'lastSbomUploaded', align: 'start' as const, width: 160},
]);

const currentProject = computed(() => appStore.getCurrentProject());
const versions = computed(() => currentProject.value?.versions ?? []);

const possibleStatus = computed(() => {
  const statuses = Array.from(new Set(versions.value.map((v: Version) => v.status || 'new')));
  return statuses.map((s) => ({
    text: t(getOverallReviewTranslationKey(s)),
    value: s,
    classes: statusClasses(s),
  }));
});

const filteredList = computed(() => {
  return versions.value.filter((v: Version) => {
    return selectedFilterStatus.value.length === 0 || selectedFilterStatus.value.includes(v.status || 'new');
  });
});

const versionDialogFormRef = vueRef();

const openAddDialog = () => {
  versionDialogFormRef.value?.open({
    projectID: currentProject.value?.uuid,
    version: undefined,
  });
};

const maxVersionsReached = computed((): boolean => {
  return versions.value.length >= maxVersions;
});

const confirmDelete = async (version: Version) => {
  try {
    const projectUuid = currentProject.value?.uuid;
    if (projectUuid) {
      await projectService.deleteProjectVersion(projectUuid, version.name);
      await appStore.refetchCurrentProject(projectUuid);
      snackbar(t('VERSION_DELETED_SUCCESSFULLY'));
    }
  } catch (error: any) {
    const errorMessage = error?.response?.data?.message || error?.message || t('ERROR_DELETING_VERSION');
    snackbar(errorMessage);
  }
};

const getStrWithMaxLength = (max: number, str: string) => {
  return str.length > max ? str.slice(0, max) + '...' : str;
};

const getReferenceInfoForClipboard = (item: Version): string => {
  const projectName = currentProject.value?.name ?? '';
  const projectId = currentProject.value?.uuid ?? '';
  const versionName = item.name ?? '';
  const description = item.description ?? '';
  return [
    `Project Name: ${projectName}`,
    `Project Identifier: ${projectId}`,
    `Project Version: ${versionName}`,
    `Version description: ${description}`,
  ].join('\n');
};

const getActionButtons = (_item: Version): TableActionButtonsProps['buttons'] => {
  return [
    {
      icon: 'mdi-content-copy',
      hint: t('TT_COPY_REFERENCE_INFO'),
      event: 'copy',
      show: true,
    },
    {
      icon: 'mdi-delete',
      hint: t('TT_delete_version'),
      event: 'delete',
      show: true,
    },
  ];
};

const deleteTitle = computed(
  () =>
    (item: Version): string =>
      `${t('Delete')} "${item.name}"`,
);

const showSbomStatus = async (item: Version) => {
  const project = currentProject.value;
  const version = item?.name || item?.version;
  if (!project || !version) {
    console.error('Missing required data:', {project, version, item});
    return;
  }

  const success = await appStore.navigateToSbom(project, version, router);
  if (!success) {
    snackbar(t('Please select a parent token first'));
  }
};
</script>

<template>
  <TableLayout has-tab has-title data-versions="versions">
    <template #description v-if="$slots.default">
      <slot></slot>
    </template>
    <template #buttons>
      <DCActionButton :disabled="maxVersionsReached" :text="t('BTN_ADD')" icon="mdi-plus" @click="openAddDialog" />
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
      <div ref="tableVersions" class="fill-height">
        <v-data-table
          fixed-header
          :sort-by="updatedSort"
          :search="search"
          :headers="headers"
          :items="filteredList"
          :footer-props="{'items-per-page-options': [10, 50, 100, -1]}"
          class="striped-table fill-height"
          density="compact"
          hover
          @click:row="(_event: Event, {item}: {item: Version}) => showSbomStatus(item)">
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
                  <Stack direction="row" justify="end" class="ma-1 mr-2">
                    <DIconButton icon="mdi-close" @clicked="statusFilterOpened = false" color="default" />
                  </Stack>
                  <v-select
                    v-model="selectedFilterStatus"
                    :items="possibleStatus"
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
                          <span :class="item.value.classes"> {{ title }}</span>
                        </template>
                      </v-list-item>
                    </template>
                    <template v-slot:selection="{item, index}">
                      <div v-if="index === 0" class="d-flex align-center">
                        <span :class="item.value.classes">{{ !item.value ? 'new' : item.title }}</span>
                      </div>
                      <span v-if="index === 1" class="pAdditionalFilter">
                        + {{ selectedFilterStatus.length - 1 }} others
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
          <template v-slot:[`item.lastSbomUploaded`]="{item}">
            <DDateCellWithTooltip :value="item.lastSbomUploaded" />
          </template>
          <template v-slot:[`item.description`]="{item}">
            <v-tooltip open-delay="300" bottom max-width="480" content-class="dpTooltip">
              <template v-slot:activator="{props}">
                <span v-bind="props"> {{ getStrWithMaxLength(50, '' + item.description) }}</span>
              </template>
              {{ '' + item.description }}
            </v-tooltip>
          </template>
          <template v-slot:[`item.actions`]="{item}">
            <DeleteConfirmationDialog
              :title="deleteTitle(item)"
              :message="t('MSG_DELETE_CHANNEL_CONFIRMATION')"
              @confirmed="() => confirmDelete(item)">
              <template #default="{showDialog}">
                <TableActionButtons
                  variant="compact"
                  :buttons="getActionButtons(item)"
                  @copy="copyToClipboard(getReferenceInfoForClipboard(item))"
                  @delete="showDialog" />
              </template>
            </DeleteConfirmationDialog>
          </template>
        </v-data-table>
      </div>
    </template>
  </TableLayout>
  <!-- Version Dialog Form -->
  <CliVersionDialogForm ref="versionDialogFormRef" />
</template>
