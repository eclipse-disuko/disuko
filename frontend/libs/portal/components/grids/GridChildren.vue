<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {IDefaultSelectItem} from '@disclosure-portal/model/IObligation';
import {Project} from '@disclosure-portal/model/Project';
import {ProjectChildrenCombiDto, ProjectSlimDto} from '@disclosure-portal/model/ProjectsResponse';
import {useDialogStore} from '@disclosure-portal/stores/dialog.store';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useWizardStore} from '@disclosure-portal/stores/wizard.store';
import {getStrWithMaxLength} from '@disclosure-portal/utils/View';
import useSnackbar from '@shared/composables/useSnackbar';
import {DataTableHeader, DataTableItem} from '@shared/types/table';
import {TOOLTIP_OPEN_DELAY_IN_MS} from '@shared/utils/constant';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRouter} from 'vue-router';

const emit = defineEmits(['openSettings']);

const {t} = useI18n();
const router = useRouter();
const projectStore = useProjectStore();
const dialogStore = useDialogStore();
const wizardStore = useWizardStore();

const currentProject = computed((): Project => projectStore.currentProject!);

const search = ref('');
const statusFilterMenu = ref(false);
const selectedFilterStatus = ref<string[]>([]);
const addChildrenProjectDialog = ref();
const errorDialog = ref();

const headers = computed<DataTableHeader[]>(() => {
  return [
    {key: 'data-table-group', title: t('COL_PROJECT_NAME'), width: 200},
    {title: t('COL_PROJECT_STATUS'), key: 'project.status', align: 'center', width: 200},
    {title: ' ', key: 'projectKey', align: 'start', width: '60'},
    {title: t('COL_STATUS'), key: 'status', align: 'center', width: 120},
    {title: t('COL_VERSION'), key: 'version.name', align: 'start', width: 140},
    {title: t('COL_DESCRIPTION'), key: 'description', align: 'start', width: 260},
    {title: t('COL_CREATED'), key: 'version.created', align: 'start', width: 160},
    {title: t('COL_UPDATED'), key: 'version.updated', align: 'start', width: 160},
  ];
});

const possibleStatus = computed(() => {
  const list = currentProject.value?.projectChildren?.list;
  if (!list?.length) return [];

  return [...new Set(list.map((item: ProjectChildrenCombiDto) => item.project.status))].map(
    (status: string) =>
      ({
        value: status,
        text: status === 'deprecated' ? t('PROJECT_DEPRECATED') : status,
      }) as IDefaultSelectItem,
  );
});

const filterOnStatus = (item: ProjectChildrenCombiDto): boolean => {
  if (!selectedFilterStatus.value.length) {
    return true;
  }
  return selectedFilterStatus.value.includes(item.project.status);
};

const filteredList = computed(() => {
  const items = currentProject.value?.projectChildren?.list;
  if (!items) return [];
  return items.filter(filterOnStatus);
});

const showCreateProjectDialog = async () => {
  if (!projectStore.areMandatoryProjectSettingsSet) {
    errorDialog.value?.open();
  } else {
    await wizardStore.openWizard({parentProject: currentProject.value});
  }
};

const openSettingsDialog = () => {
  dialogStore.isSettingsDialogOpen = true;
};

const openProject = (item: ProjectSlimDto) => {
  if (item.isDeleted) {
    return;
  }
  router.push({
    path: `${item.isGroup ? '/dashboard/groups/' : '/dashboard/projects/'}${encodeURIComponent(item._key)}`,
  });
};

const openVersion = (item: ProjectChildrenCombiDto) => {
  const version = item.version;
  if (version.isDeleted || item.project.isDeleted || !item.hasProjectReadAccess) {
    useSnackbar().info(t('SNACK_MISSING_RIGHTS'));
    return;
  }
  router.push({
    path: `/dashboard/projects/${encodeURIComponent(version.parentKey)}/versions/${encodeURIComponent(version._key)}`,
  });
};

const showAddChildrenProjectDialog = () => {
  if (!projectStore.areMandatoryProjectSettingsSet) {
    errorDialog.value?.open();
  } else {
    addChildrenProjectDialog.value?.open();
  }
};

const customFilter = (value: string, query: string, item?: any) => {
  if (!query || !item) return true;
  const searchQuery = query.toLowerCase();
  const rawItem = item.raw as ProjectChildrenCombiDto;
  const projectName = rawItem.project?.name?.toLowerCase() || '';

  return projectName.includes(searchQuery);
};

const getStatusClass = computed(() => (status?: string) => {
  const statusValue = status?.toLowerCase() || 'new';
  return {
    [`pStatus${statusValue}`]: true,
  };
});
</script>

<template>
  <TableLayout has-title has-tab>
    <template v-if="$slots.default" #description>
      <slot></slot>
    </template>
    <template #buttons>
      <span class="d-headline-2">{{ t('RELATED_PROJECTS') }}</span>

      <DCActionButton
        v-if="currentProject?.allowGroupCreate"
        large
        icon="mdi-plus"
        :hint="t('BTN_ADD_CHILDREN')"
        :text="t('BTN_ADD')"
        @click="showCreateProjectDialog" />
      <DCActionButton
        large
        icon="mdi-pencil"
        :hint="t('TT_children')"
        :text="t('BTN_EDIT')"
        @click.stop="showAddChildrenProjectDialog"
        v-if="currentProject?.allowGroupCreate" />
      <v-spacer></v-spacer>
      <v-text-field
        autocomplete="off"
        :max-width="500"
        v-model="search"
        :label="t('labelSearch')"
        append-inner-icon="mdi-magnify"
        variant="outlined"
        clearable
        density="compact"
        hide-details />
    </template>
    <template #table>
      <v-data-table
        v-if="currentProject?.projectChildren?.list"
        density="compact"
        class="striped-table fill-height"
        fixed-header
        :sort-by="[{key: 'version.updated', order: 'desc'}]"
        :group-by="[{key: 'projectKey'}]"
        :search="search"
        :custom-filter="customFilter"
        :headers="headers"
        :items="filteredList"
        @click:row="(event: Event, dataItem: DataTableItem<ProjectChildrenCombiDto>) => openVersion(dataItem.item)"
        :items-per-page="-1">
        <template v-slot:[`item.data-table-expand`]="{}"> x </template>
        <template v-slot:group-header="{item, toggleGroup, isGroupOpen}">
          <tr
            :class="{'cursor-pointer': !item.items[0].raw.project.isDeleted && item.items[0].raw.hasProjectReadAccess}"
            @click.stop="openProject(item.items[0].raw.project)">
            <td>
              <v-icon class="mr-2" color="primary" @click.stop="toggleGroup(item)">
                {{ isGroupOpen(item) ? 'mdi-chevron-up' : 'mdi-chevron-down' }}
              </v-icon>
              <span v-if="item.items[0].raw.project.isDeleted" class="disabledText">
                {{ item.items[0].raw.project.name }}<span class="deleted">&nbsp;{{ t('PROJECT_DELETED') }}</span>
              </span>
              <span v-else-if="!item.items[0].raw.hasProjectReadAccess" class="disabledText">
                {{ item.items[0].raw.project.name }}
                <span v-if="item.items[0].raw.project.status === 'deprecated'" class="disabledText">
                  &nbsp;[{{ t('PROJECT_DEPRECATED') }}]
                </span>
                <span class="deleted">&nbsp;{{ t('INSUFFICIENT_PERMISSIONS') }}</span>
              </span>
              <span v-else>
                {{ item.items[0].raw.project.name }}
                <span v-if="item.items[0].raw.project.status === 'deprecated'" class="disabledText">
                  &nbsp;[{{ t('PROJECT_DEPRECATED') }}]
                </span>
              </span>
            </td>
            <td class="text-center">
              <span :class="getStatusClass(item.items[0].raw.project.status)">
                {{
                  item.items[0].raw.project.status === 'deprecated'
                    ? t('PROJECT_DEPRECATED')
                    : item.items[0].raw.project.status
                }}
              </span>
            </td>
            <td :colspan="headers.length - 2"></td>
          </tr>
        </template>
        <template v-slot:[`item.projectKey`]="{}">
          <!-- blank cause in expand row-->
          <span>&nbsp;</span>
        </template>
        <template v-slot:[`item.status`]="{item}">
          <DVersionStateWithTooltip v-if="item.version" :version="item.version" :isGroup="true" />
        </template>
        <template v-slot:[`header.project.status`]="{column}">
          <div class="v-data-table-header__content">
            <span>{{ column.title }}</span>
            <v-menu :close-on-content-click="false" v-model="statusFilterMenu">
              <template v-slot:activator="{props}">
                <DIconButton
                  :parentProps="props"
                  icon="mdi-filter-variant"
                  :hint="t('TT_SHOW_FILTER')"
                  :color="selectedFilterStatus.length > 0 ? 'primary' : 'default'" />
              </template>
              <div class="bg-background" style="width: 280px">
                <v-row class="d-flex justify-end ma-1 mr-2">
                  <DIconButton icon="mdi-close" @clicked="statusFilterMenu = false" color="default" />
                </v-row>
                <v-select
                  v-model="selectedFilterStatus"
                  :items="possibleStatus"
                  class="mx-2 pa-2 pb-4"
                  :label="t('FILTER_BY_STATUS')"
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
                    <v-list-item v-bind="props" class="py-0 px-2">
                      <template v-slot:prepend="{isSelected}">
                        <v-checkbox hide-details :model-value="isSelected" />
                      </template>
                      <template v-slot:title>
                        <span :class="getStatusClass(item.raw.value)" class="pFilterEntry">{{ item.raw.text }}</span>
                      </template>
                    </v-list-item>
                  </template>
                  <template v-slot:selection="{item, index}">
                    <div v-if="index === 0" class="d-flex align-center">
                      <span class="pFilterEntry">{{ item.raw.text }}</span>
                    </div>
                    <span v-if="index === 1" class="pAdditionalFilter">
                      +{{ selectedFilterStatus.length - 1 }} others
                    </span>
                  </template>
                </v-select>
              </div>
            </v-menu>
          </div>
        </template>
        <template v-slot:[`item.project.status`]> </template>
        <template v-slot:[`item.version.updated`]="{item}">
          <DDateCellWithTooltip :value="item.version.updated" v-if="item.version" />
        </template>
        <template v-slot:[`item.version.created`]="{item}">
          <DDateCellWithTooltip :value="item.version.created" v-if="item.version" />
        </template>
        <template v-slot:[`item.description`]="{item}">
          <v-tooltip
            :open-delay="TOOLTIP_OPEN_DELAY_IN_MS"
            bottom
            max-width="480"
            v-if="item.version"
            content-class="dpTooltip">
            <template v-slot:activator="{props}">
              <span v-bind="props"> {{ getStrWithMaxLength(50, '' + item.version.description) }}</span>
            </template>
            {{ item.version.description }}
          </v-tooltip>
        </template>
      </v-data-table>
    </template>
  </TableLayout>

  <AddChildrenDialog ref="addChildrenProjectDialog"></AddChildrenDialog>
  <AddChildrenErrorDialog ref="errorDialog" @open-settings="openSettingsDialog"></AddChildrenErrorDialog>
</template>
