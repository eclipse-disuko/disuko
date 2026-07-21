<script setup lang="ts">
import {Project, Version} from '@cli/models/Project';
import {useAppStore} from '@cli/stores/app';
import DDateCellWithTooltip from '@shared/components/disco/DDateCellWithTooltip.vue';
import DVersionStateWithTooltip from '@shared/components/disco/DVersionStateWithTooltip.vue';
import {useTableFilter} from '@shared/composables/useTableFilter';
import TableLayout from '@shared/layouts/TableLayout.vue';
import {computed, defineProps, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRouter} from 'vue-router';
import type {DataTableHeader} from 'vuetify/lib/components/VDataTable/types';
const appStore = useAppStore();
const search = ref('');

const props = defineProps<{children: Project[]}>();
const {t} = useI18n();
const router = useRouter();
const currentProject = computed(() => appStore.getCurrentProject());

const childrenList = computed(() => {
  if (!Array.isArray(props.children)) return [];
  return props.children.flatMap((project) => {
    const versions = project.versions || [];
    return versions.map((version) => ({
      project,
      version,
      projectKey: project.uuid,
      projectName: project.name,
      description: version.description,
      status: version.status,
      lastSbomUploaded: version.lastSbomUploaded,
    }));
  });
});

const sortItems = [{key: 'version.lastSbomUploaded', order: 'desc' as const}];
const groupBy = [{key: 'projectKey'}];

const customFilter = useTableFilter(['projectName', 'version.name', 'description']);

const showSbomStatus = async (
  item: {project: Project; version: Version} | {raw?: {project: Project; version: Version}},
) => {
  const normalized = (item as any)?.raw ?? item;
  const project = normalized.project;
  const version = normalized.version?.name;

  if (!project?.uuid || !version) {
    console.error('Missing required data:', {projectId: project?.uuid, version, item});
    return;
  }

  await appStore.navigateToSbom(project, version, router);
};

const openProject = (project: Project) => {
  const parent = currentProject.value;
  if (parent) {
    appStore.setParentProject({
      uuid: parent.uuid,
      name: parent.name,
    });
  }
  router.push(`/projects/${project.uuid}/overview`);
};

const headers = computed<DataTableHeader[]>(() => [
  {key: 'data-table-group', title: t('COL_PROJECT_NAME'), width: 140},
  {key: 'projectKey', title: ' ', align: 'start', width: '60'},
  {key: 'status', title: t('COL_STATUS'), align: 'center', width: 120},
  {key: 'version.name', title: t('COL_CHANNEL'), align: 'start', width: 140},
  {key: 'description', title: t('COL_DESCRIPTION'), align: 'start', width: 140},
  {key: 'version.lastSbomUploaded', title: t('COL_LAST_SBOM_UPLOADED'), align: 'start', width: 220},
]);
</script>

<template>
  <TableLayout has-tab has-title>
    <template #description v-if="$slots.default">
      <slot></slot>
    </template>

    <template #buttons>
      <span class="d-headline-2 pr-2">{{ t('RELATED_PROJECTS') }}</span>
      <v-spacer></v-spacer>
      <v-text-field
        autocomplete="off"
        class="max-w-[400px]"
        v-model="search"
        :label="t('labelSearch')"
        append-inner-icon="mdi-magnify"
        variant="outlined"
        clearable
        density="compact"
        hide-details />
    </template>
    <template #table>
      <div class="fill-height">
        <v-data-table
          density="compact"
          class="striped-table fill-height"
          fixed-header
          :sort-by="sortItems"
          :group-by="groupBy"
          :headers="headers"
          :items="childrenList"
          :search="search"
          :custom-filter="customFilter"
          :items-per-page="-1"
          :no-data-text="t('NO_DATA_AVAILABLE')"
          hover
          @click:row="(_event: Event, {item}: {item: any}) => showSbomStatus(item)">
          <template v-slot:group-header="{item, toggleGroup, isGroupOpen}">
            <tr class="cursor-pointer" @click.stop="openProject(item.items[0].raw.project)">
              <td :colspan="headers.length">
                <v-icon class="mr-2" color="primary" @click.stop="toggleGroup(item)">
                  {{ isGroupOpen(item) ? 'mdi-chevron-up' : 'mdi-chevron-down' }}
                </v-icon>
                <span>
                  {{ item.items[0].raw.project.name }}
                </span>
              </td>
            </tr>
          </template>

          <template v-slot:[`item.projectKey`]>
            <span>&nbsp;</span>
          </template>

          <template v-slot:[`item.status`]="{item}">
            <DVersionStateWithTooltip :version="item.version" />
          </template>

          <template v-slot:[`item.version.name`]="{item}">
            <span class="cursor-pointer" @click.stop="showSbomStatus(item)">{{ item.version.name }}</span>
          </template>

          <template v-slot:[`item.description`]="{item}">
            <v-tooltip bottom max-width="480">
              <template #activator="{props}">
                <span v-bind="props">
                  {{
                    item.description
                      ? item.description.length > 50
                        ? item.description.slice(0, 50) + '...'
                        : item.description
                      : ''
                  }}
                </span>
              </template>
              {{ item.description || '' }}
            </v-tooltip>
          </template>

          <template v-slot:[`item.version.lastSbomUploaded`]="{item}">
            <DDateCellWithTooltip :value="item.version.lastSbomUploaded" />
          </template>
        </v-data-table>
      </div>
    </template>
  </TableLayout>
</template>
