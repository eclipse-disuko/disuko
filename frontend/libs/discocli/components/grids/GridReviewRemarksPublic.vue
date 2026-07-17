<script setup lang="ts">
import CliReviewRemarkDialog from '@cli/components/dialogs/CliReviewRemarkDialog.vue';
import {compareRRLevel, type RemarkWithVersion} from '@cli/models/ReviewRemark';
import {projectService} from '@cli/services/projectService';
import {useAppStore} from '@cli/stores/app';
import type {ReviewRemarkLevel} from '@disclosure-portal/model/Quality';
import {getIconColorReviewRemarkLevel, getIconReviewRemarkLevel} from '@disclosure-portal/utils/View';
import Tooltip from '@shared/components/disco/Tooltip.vue';
import TableActionButtons, {TableActionButtonsProps} from '@shared/components/TableActionButtons.vue';
import {useTableFilter} from '@shared/composables/useTableFilter';
import {useClipboard} from '@shared/utils/clipboard';
import {computed, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRoute} from 'vue-router';
import type {DataTableHeader} from 'vuetify';

const {t} = useI18n();
const route = useRoute();
const appStore = useAppStore();
const COMPONENTS_TOOLTIP_THRESHOLD = 30;
// Get version from URL params if it exists
const urlVersion = computed(() => {
  const routeVersion = route.params.version;
  return typeof routeVersion === 'string' ? routeVersion : null;
});
const projectUuid = computed(() => route.params.id as string | null);

const remarks = ref<RemarkWithVersion[]>([]);

const search = ref('');
const loading = ref(true);
const sortItems = ref([{key: 'level', order: 'desc' as const}]);
const {copyToClipboard} = useClipboard();

const currentProject = computed(() => appStore.getCurrentProject());
const versions = computed(() => currentProject.value?.versions || []);

const headers = computed<DataTableHeader[]>(() => [
  {key: 'actions', title: t('COL_ACTIONS'), align: 'center', width: 120, sortable: false, class: 'tableHeaderCell'},
  {key: 'version', title: t('COL_CHANNEL'), align: 'center', width: 140, sortable: true, class: 'tableHeaderCell'},
  {key: 'comment', title: t('COL_COMMENT'), align: 'center', width: 120, sortable: false, class: 'tableHeaderCell'},
  {
    key: 'level',
    title: t('COL_LEVEL'),
    align: 'center',
    width: 120,
    sortable: true,
    class: 'tableHeaderCell',
    sort: compareRRLevel,
  },
  {key: 'status', title: t('COL_STATUS'), align: 'start', width: 120, sortable: true, class: 'tableHeaderCell'},
  {key: 'title', title: t('COL_REVIEW_REMARK'), align: 'start', width: 200, sortable: false, class: 'tableHeaderCell'},
  {
    key: 'components',
    title: t('COL_COMPONENTS'),
    align: 'start',
    width: 120,
    sortable: false,
    class: 'tableHeaderCell',
  },
  {
    key: 'sbomName',
    title: t('COL_SBOM_REFERENCE'),
    align: 'start',
    width: 200,
    sortable: false,
    class: 'tableHeaderCell',
  },
  {key: 'licenses', title: t('COL_LICENSES'), align: 'start', width: 125, sortable: false, class: 'tableHeaderCell'},
  {key: 'author', title: t('COL_CREATOR'), align: 'start', width: 100, sortable: true, class: 'tableHeaderCell'},
  {key: 'origin', title: t('COL_ORIGIN'), align: 'center', width: 100, sortable: false, class: 'tableHeaderCell'},
  {key: 'closed', title: t('COL_CLOSED'), align: 'start', width: 100, sortable: true, class: 'tableHeaderCell'},
  {key: 'created', title: t('COL_CREATED'), align: 'start', width: 100, sortable: true, class: 'tableHeaderCell'},
  {key: 'updated', title: t('COL_UPDATED'), align: 'start', width: 100, sortable: true, class: 'tableHeaderCell'},
]);

const selectedFilterStatus = ref(['OPEN', 'IN_PROGRESS']);
const statusFilterMenu = ref(false);
const possibleStatus = computed(() => {
  const allStatuses = Array.from(new Set(remarks.value.map((r) => r.status).filter(Boolean)));
  return allStatuses.map((s) => ({
    value: s,
    text: t('REMARK_STATUS_' + s),
  }));
});

const selectedFilterVersion = ref<string[]>([]);
const versionFilterMenu = ref(false);
const possibleVersions = computed(() => {
  const allVersions = Array.from(new Set(remarks.value.map((r) => r.version).filter(Boolean)));
  return allVersions.map((v) => ({
    value: v,
    text: v,
  }));
});

const matchesSelectedStatus = (remark: RemarkWithVersion) =>
  !selectedFilterStatus.value.length || (remark.status && selectedFilterStatus.value.includes(remark.status));

const matchesSelectedVersion = (remark: RemarkWithVersion) =>
  !selectedFilterVersion.value.length || (remark.version && selectedFilterVersion.value.includes(remark.version));

const withFormattedComponents = (remark: RemarkWithVersion) => {
  const componentsText = formatComponents(remark.components);
  const licensesText = formatLicenses(remark.licenses);
  return {
    ...remark,
    _componentsText: componentsText,
    _componentsNeedsTooltip: componentsText.length > COMPONENTS_TOOLTIP_THRESHOLD,
    _licensesText: licensesText,
  };
};

const customFilter = useTableFilter([
  'title',
  'description',
  '_componentsText',
  '_licensesText',
  'sbomName',
  'author',
  'status',
]);

const filteredList = computed(() => {
  return remarks.value
    .filter((remark) => matchesSelectedStatus(remark) && matchesSelectedVersion(remark))
    .map(withFormattedComponents);
});

type ComponentObj = {componentName?: string; componentVersion?: string};
type LicenseObj = {licenseName?: string; licenseId?: string};

function formatComponents(components: string[] | ComponentObj[] | string | null | undefined): string {
  if (!components) return '';
  if (Array.isArray(components)) {
    return components
      .map((comp) =>
        typeof comp === 'object' && comp !== null
          ? `${comp.componentName || ''}${comp.componentVersion ? ` (${comp.componentVersion})` : ''}`.trim()
          : comp,
      )
      .filter(Boolean)
      .join('; ');
  }
  return String(components);
}

function formatLicenses(licenses: string[] | LicenseObj[] | string | null | undefined): string {
  if (!licenses) return '-';
  if (Array.isArray(licenses)) {
    return licenses
      .map((lic) =>
        typeof lic === 'object' && lic !== null
          ? `${lic.licenseName || ''}${lic.licenseId ? ` (${lic.licenseId})` : ''}`.trim()
          : lic,
      )
      .filter(Boolean)
      .join('; ');
  }
  return String(licenses);
}

const getRemarkInfoForClipboard = (item: RemarkWithVersion): string => {
  const lines: string[] = [];
  lines.push(`${t('COL_CHANNEL')}: ${item.version || '-'}`);
  lines.push(`${t('COL_LEVEL')}: ${item.level || '-'}`);
  lines.push(`${t('COL_STATUS')}: ${item.status || '-'}`);
  lines.push(`${t('COL_COMPONENTS')}: ${formatComponents(item.components)}`);
  lines.push(`${t('COL_SBOM_NAME')}: ${item.sbomName || '-'}`);
  lines.push(`${t('COL_LICENSES')}: ${formatLicenses(item.licenses)}`);
  lines.push(`${t('COL_REVIEW_REMARK')}: ${item.title || '-'}`);
  lines.push(`${t('COL_ORIGIN')}: ${item.origin || '-'}`);
  lines.push(`${t('COL_AUTHOR')}: ${item.author || '-'}`);
  lines.push(`${t('COL_CLOSED')}: ${item.closed || '-'}`);
  lines.push(`${t('COL_CREATED')}: ${item.created || '-'}`);
  lines.push(`${t('COL_UPDATED')}: ${item.updated || '-'}`);
  return lines.join('\n');
};

const getActionButtons = (_item: RemarkWithVersion): TableActionButtonsProps['buttons'] => {
  return [
    {
      icon: 'mdi-comment-plus-outline',
      hint: t('BTN_COMMENT'),
      event: 'comment',
      show: true,
    },

    {
      icon: 'mdi-content-copy',
      hint: t('TT_COPY_REFERENCE_INFO'),
      event: 'copy',
      show: true,
    },
  ];
};

const fetchRemarks = async () => {
  if (!projectUuid.value) {
    return;
  }

  loading.value = true;
  remarks.value = [];

  try {
    if (urlVersion.value) {
      const versionRemarks = await projectService.getVersionReviewRemarks(projectUuid.value, urlVersion.value);
      if (versionRemarks) {
        remarks.value = versionRemarks.map((remark) => ({
          ...remark,
          version: urlVersion.value!,
        }));
      }
    } else if (versions.value.length > 0) {
      const results = await Promise.all(
        versions.value.map(async (version) => {
          if (!projectUuid.value) return [];
          const versionRemarks = await projectService.getVersionReviewRemarks(projectUuid.value, version.name);
          return (
            versionRemarks?.map((remark) => ({
              ...remark,
              version: version.name,
            })) || []
          );
        }),
      );
      remarks.value = results.flat();
    }
  } catch (error) {
    console.error('Error fetching remarks:', error);
    remarks.value = [];
  } finally {
    loading.value = false;
  }
};

onMounted(() => {
  fetchRemarks();
});

const reviewRemarkDialog = ref<InstanceType<typeof CliReviewRemarkDialog>>();

function openReviewRemarkDialog(item?: RemarkWithVersion) {
  reviewRemarkDialog.value?.open(item ? item : undefined);
}
</script>

<template>
  <TableLayout has-tab has-title>
    <template #description v-if="$slots.default">
      <slot></slot>
    </template>
    <template v-else #description>
      <div v-if="!urlVersion" class="text-body-2 line-clamp-1 max-w-[50%]">
        {{ t('TABLE_DESC_PROJECT_ALL_REVIEW_REMARKS') }}
      </div>
    </template>
    <template #buttons>
      <v-spacer></v-spacer>
      <v-text-field
        v-model="search"
        append-inner-icon="mdi-magnify"
        :label="t('labelSearch')"
        variant="outlined"
        clearable
        density="compact"
        hide-details
        style="max-width: 400px" />
    </template>

    <template #table>
      <div class="fill-height">
        <v-data-table
          fixed-header
          :loading="loading"
          :sort-by="sortItems"
          :search="search"
          :custom-filter="customFilter"
          :headers="headers"
          :items="filteredList"
          :items-per-page="100"
          :footer-props="{'items-per-page-options': [10, 25, 50, 100, -1]}"
          must-sort
          class="striped-table fill-height"
          density="compact"
          @click:row="(_, {item}) => openReviewRemarkDialog(item)">
          <template v-slot:[`item.comment`]="{item}">
            <span @click.stop="openReviewRemarkDialog(item)" class="inline-block cursor-pointer p-2">
              <v-badge
                v-if="item && item.events && item.events.length > 0"
                :content="String(item.events.length)"
                color="mbti"
                class="dpTextBlack"
                overlap
                bottom>
                <span style="display: inline-flex; min-width: 24px; justify-content: center">
                  <v-icon aria-label="Show comments"> mdi-comment-text-outline </v-icon>
                </span>
              </v-badge>
              <v-icon v-else aria-label="Show comments"> mdi-comment-outline </v-icon>
            </span>
          </template>
          <template v-slot:[`item.version`]="{item}">
            <span class="text-medium-emphasis">{{ item.version }}</span>
          </template>
          <template v-slot:[`item.level`]="{item}">
            <v-tooltip :open-delay="300" bottom content-class="dpTooltip">
              <template v-slot:activator="{}">
                <v-icon :color="getIconColorReviewRemarkLevel(item.level as ReviewRemarkLevel)">
                  {{ getIconReviewRemarkLevel(item.level as ReviewRemarkLevel) }}
                </v-icon>
              </template>
              <span>{{ t('REMARK_LEVEL_' + (item.level || '')) }}</span>
            </v-tooltip>
          </template>
          <template v-slot:[`item.status`]="{item}">
            <span>{{ t('REMARK_STATUS_' + (item.status || '')) }}</span>
          </template>
          <template v-slot:[`item.components`]="{item}">
            <span v-if="!item._componentsNeedsTooltip">{{ item._componentsText }}</span>
            <Tooltip v-else as-parent :text="item._componentsText">
              <span class="block max-w-[180px] truncate">
                {{ item._componentsText }}
              </span>
            </Tooltip>
          </template>
          <template v-slot:[`item.sbomName`]="{item}">
            <span>{{ item.sbomName }}</span>
          </template>
          <template v-slot:[`item.licenses`]="{item}">
            <span>
              <template v-if="Array.isArray(item.licenses)">
                <template v-for="(lic, idx) in item.licenses" :key="idx">
                  <span v-if="lic && typeof lic === 'object' && 'licenseName' in lic && 'licenseId' in lic">
                    {{ lic.licenseName }}<template v-if="lic.licenseId"> ({{ lic.licenseId }})</template
                    ><span v-if="idx < item.licenses.length - 1">; </span>
                  </span>
                </template>
              </template>
              <template v-else>{{ item.licenses || '' }}</template>
            </span>
          </template>
          <template v-slot:[`item.description`]="{item}">
            <span>{{ item.description || '' }}</span>
          </template>
          <template v-slot:[`item.author`]="{item}">
            <span>{{ item.author }}</span>
          </template>
          <template v-slot:[`item.origin`]="{item}">
            <span>{{ item.origin }}</span>
          </template>
          <template v-slot:[`item.closed`]="{item}">
            <DDateCellWithTooltip :value="item.closed || undefined" />
          </template>
          <template v-slot:[`item.created`]="{item}">
            <DDateCellWithTooltip :value="item.created" />
          </template>
          <template v-slot:[`item.updated`]="{item}">
            <DDateCellWithTooltip :value="item.updated" />
          </template>
          <template v-slot:[`item.actions`]="{item}">
            <TableActionButtons
              variant="compact"
              :buttons="getActionButtons(item)"
              @comment="openReviewRemarkDialog(item)"
              @copy="copyToClipboard(getRemarkInfoForClipboard(item))" />
          </template>
          <template v-slot:header.version="{column, getSortIcon, toggleSort}">
            <div class="v-data-table-header__content">
              <span>{{ column.title }}</span>
              <v-menu :close-on-content-click="false" v-model="versionFilterMenu">
                <template v-slot:activator="{props}">
                  <DIconButton
                    :parentProps="props"
                    icon="mdi-filter-variant"
                    :hint="t('TT_SHOW_FILTER')"
                    :color="selectedFilterVersion.length > 0 ? 'primary' : 'default'" />
                </template>
                <div class="bg-background" style="width: 220px">
                  <v-row class="d-flex ma-1 mr-2 justify-end">
                    <DCloseButton @click="versionFilterMenu = false" size="x-small" />
                  </v-row>
                  <v-select
                    v-model="selectedFilterVersion"
                    :items="possibleVersions"
                    :label="t('Lbl_filter_channel')"
                    clearable
                    multiple
                    item-title="text"
                    item-value="value"
                    variant="outlined"
                    density="compact"
                    hide-details />
                </div>
              </v-menu>
              <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
            </div>
          </template>
          <template v-slot:header.status="{column, getSortIcon, toggleSort}">
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
                <div class="bg-background" style="width: 220px">
                  <v-row class="d-flex ma-1 mr-2 justify-end">
                    <DCloseButton @click="statusFilterMenu = false" size="x-small" />
                  </v-row>
                  <v-select
                    v-model="selectedFilterStatus"
                    :items="possibleStatus"
                    :label="t('Lbl_filter_status')"
                    clearable
                    multiple
                    item-title="text"
                    item-value="value"
                    variant="outlined"
                    density="compact"
                    hide-details />
                </div>
              </v-menu>
              <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
            </div>
          </template>
        </v-data-table>
      </div>
    </template>
  </TableLayout>
  <CliReviewRemarkDialog ref="reviewRemarkDialog" @refresh="fetchRemarks" />
</template>
