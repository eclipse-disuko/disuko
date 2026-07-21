<script setup lang="ts">
import AddExternalSourceDialog from '@cli/components/dialogs/AddExternalSourceDialog.vue';
import {ExtendedExternalSource, ExternalSourceCodeReference} from '@cli/models/ExternalSourceCode';
import {projectService} from '@cli/services/projectService';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import DDateCellWithTooltip from '@shared/components/disco/DDateCellWithTooltip.vue';
import DExternalLink from '@shared/components/disco/DExternalLink.vue';
import Tooltip from '@shared/components/disco/Tooltip.vue';
import TableActionButtons, {TableActionButtonsProps} from '@shared/components/TableActionButtons.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import TableLayout from '@shared/layouts/TableLayout.vue';
import {useClipboard} from '@shared/utils/clipboard';
import {computed, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRoute} from 'vue-router';
import type {DataTableHeader} from 'vuetify';

const {t} = useI18n();
const route = useRoute();
const snackbar = useSnackbar();
const {copyToClipboard} = useClipboard();

// Get version from URL params
const urlVersion = computed(() => {
  const routeVersion = route.params.version;
  return typeof routeVersion === 'string' ? routeVersion : null;
});
const projectUuid = computed(() => route.params.id as string | null);

const canLoad = computed(() => !!projectUuid.value && !!urlVersion.value);

const externalSources = ref<ExternalSourceCodeReference[]>([]);
const loading = ref(false);
const search = ref('');
const sortItems = ref([{key: 'created', order: 'desc' as const}]);
const dlgAddSource = ref<InstanceType<typeof AddExternalSourceDialog> | null>(null);

const headers: DataTableHeader[] = [
  {key: 'actions', title: t('COL_ACTIONS'), align: 'center', width: 180, sortable: false},
  {key: 'url', title: t('COL_URL'), align: 'start', width: 350, sortable: true},
  {key: 'comment', title: t('COL_DESCRIPTION'), align: 'start', width: 300, sortable: true},
  {key: 'origin', title: t('COL_ORIGIN'), align: 'start', width: 120, sortable: true},
  {key: 'uploader', title: t('COL_UPLOADER'), align: 'start', width: 120, sortable: true},
  {key: 'created', title: t('COL_CREATED'), align: 'start', width: 120, sortable: true},
];

const getActionButtons = (_item: ExtendedExternalSource): TableActionButtonsProps['buttons'] => {
  return [
    {
      icon: 'mdi-content-copy',
      hint: t('TT_COPY_REFERENCE_INFO'),
      event: 'copy',
      show: true,
    },
  ];
};

// Helper functions
const getStrWithMaxLength = (max: number, str: string): string => {
  return str && str.length > max ? str.slice(0, max) + '...' : str;
};

const getReferenceInfoForClipboard = (item: ExternalSourceCodeReference): string => {
  return [
    `URL: ${item.url}`,
    `Description: ${item.comment || '-'}`,
    `Origin: ${item.origin || '-'}`,
    `Uploader: ${item.uploader || '-'}`,
    `Created: ${item.created || '-'}`,
  ].join('\n');
};

// Load data
const loadExternalSources = async () => {
  if (!canLoad.value) return;

  loading.value = true;
  try {
    if (projectUuid.value && urlVersion.value) {
      const result = await projectService.getExternalSourceCodeReferences(projectUuid.value, urlVersion.value);
      externalSources.value = result || [];
    }
  } catch (error) {
    console.error('Error loading external source code references:', error);
    externalSources.value = [];
  } finally {
    loading.value = false;
  }
};

// Load data when component mounts
onMounted(async () => {
  await loadExternalSources();
});

// Handle dialog submit event
const handleDialogSubmit = async (data: {url: string; comment?: string}) => {
  if (!projectUuid.value || !urlVersion.value) {
    console.warn('Missing projectUuid or urlVersion');
    return;
  }

  try {
    await projectService.createExternalSourceCodeReference(projectUuid.value, urlVersion.value, data);
    snackbar.info(t('DIALOG_source_code_add_success'));
    dlgAddSource.value?.close();
    await loadExternalSources();
  } catch (error) {
    console.error('Error adding external source:', error);
    snackbar.info(t('DIALOG_source_code_add_error'));
  }
};

const extendedExternalSources = computed<ExtendedExternalSource[]>(() => {
  return externalSources.value.map((item) => ({
    ...item,
    _displayUrl: item.url.startsWith('file://') ? item.url : getStrWithMaxLength(45, item.url),
    _isFileUrl: item.url.startsWith('file://'),
    _displayComment: getStrWithMaxLength(40, item.comment || '-'),
    _displayOrigin: getStrWithMaxLength(30, item.origin || '-'),
    _displayUploader: getStrWithMaxLength(30, item.uploader || '-'),
    _referenceInfo: getReferenceInfoForClipboard(item),
    _createdString: '' + (item.created || ''),
  }));
});
</script>

<template>
  <TableLayout has-tab has-title>
    <template #description>
      <div class="text-body-2 line-clamp-1 max-w-[50%] cursor-help">
        {{ t('SOURCE_CODE_DISCLAIMER_TEXT') }}
      </div>
      <Tooltip :text="t('SOURCE_CODE_DISCLAIMER_TEXT')" />
    </template>
    <template #buttons>
      <span class="text-h6">{{ t('TITLE_SOURCE_CODE') }}</span>
      <DCActionButton
        large
        :text="t('BTN_ADD')"
        class="text-none"
        icon="mdi-plus"
        :hint="t('TT_ADD_SOURCE')"
        @click="dlgAddSource?.open()" />
      <v-spacer></v-spacer>
      <v-text-field
        autocomplete="off"
        v-model="search"
        append-inner-icon="mdi-magnify"
        :label="t('labelSearch')"
        clearable
        density="compact"
        variant="outlined"
        hide-details
        style="max-width: 400px" />
    </template>
    <template #table>
      <div class="fill-height">
        <v-data-table
          density="compact"
          class="striped-table fill-height"
          :headers="headers"
          fixed-header
          item-key="url"
          :sort-by="sortItems"
          :search="search"
          :items="extendedExternalSources"
          :loading="loading"
          :footer-props="{
            'items-per-page-options': [10, 50, 100, -1],
          }">
          <!-- URL Column -->
          <template v-slot:item.url="{item}">
            <span v-if="!item._isFileUrl">
              <DExternalLink :text="item._displayUrl" :url="item.url"></DExternalLink>
              <Tooltip>{{ t('OPEN_URL_EXTERN') }} {{ item.url }}</Tooltip>
            </span>
            <span v-else>{{ item.url }}</span>
          </template>

          <!-- Description/Comment Column -->
          <template v-slot:item.comment="{item}">
            <span class="text-medium-emphasis">
              {{ item._displayComment }}
            </span>
            <Tooltip :text="item.comment || '-'" />
          </template>

          <!-- Origin Column -->
          <template v-slot:item.origin="{item}">
            <span>{{ item._displayOrigin }}</span>
            <Tooltip :text="item.origin || '-'" />
          </template>

          <!-- Uploader Column -->
          <template v-slot:item.uploader="{item}">
            <span>{{ item._displayUploader }}</span>
            <Tooltip :text="item.uploader || '-'" />
          </template>

          <!-- Created Column -->
          <template v-slot:item.created="{item}">
            <DDateCellWithTooltip :value="item._createdString" />
          </template>

          <!-- Actions Column -->
          <template v-slot:item.actions="{item}">
            <TableActionButtons
              variant="compact"
              :buttons="getActionButtons(item)"
              @copy="copyToClipboard(item._referenceInfo)" />
          </template>
        </v-data-table>
      </div>
    </template>
  </TableLayout>

  <!-- Add External Source Dialog -->
  <AddExternalSourceDialog ref="dlgAddSource" @submit="handleDialogSubmit" />
</template>
