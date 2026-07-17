<script setup lang="ts">
import {SBOM, SpdxUploadResponse} from '@cli/models/Sbom';
import {projectService} from '@cli/services/projectService';
import {useAppStore} from '@cli/stores/app';
import {NameKeyIdentifier} from '@disclosure-portal/model/ProjectsResponse';
import {formatDateTime} from '@disclosure-portal/utils/View';
import {isAxiosError} from 'axios';
import DCopyClipboardButton from '@shared/components/disco/DCopyClipboardButton.vue';
import DDateCellWithTooltip from '@shared/components/disco/DDateCellWithTooltip.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import DialogLayout from '@shared/layouts/DialogLayout.vue';
import TableLayout from '@shared/layouts/TableLayout.vue';
import TableActionButtons, {TableActionButtonsProps} from '@shared/components/TableActionButtons.vue';
import {computed, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRouter, useRoute} from 'vue-router';
import type {DataTableHeader} from 'vuetify';
import DSpdxTagDialog from '../dialogs/CliSpdxTagDialog.vue';
import {useClipboard} from '@shared/utils/clipboard';
import {useTableFilter} from '@shared/composables/useTableFilter';
const selectedVersion = ref<NameKeyIdentifier | null>(null);

const appStore = useAppStore();
const currentProject = computed(() => appStore.getCurrentProject());
const {t} = useI18n();
const search = ref('');
const sortItems = ref([{key: 'details.uploaded', order: 'desc' as const}]);

const router = useRouter();
const route = useRoute();

// Get version from URL params if it exists
const urlVersion = computed(() => {
  const routeVersion = route.params.version;
  return typeof routeVersion === 'string' ? routeVersion : null;
});

const projectUuid = computed(() => route.params.id as string | null);

watch(
  urlVersion,
  (newVersion) => {
    if (newVersion) {
      selectedVersion.value = {name: newVersion, key: newVersion};
    }
  },
  {immediate: true},
);

const showSbomStatus = (item: SBOM) => {
  const projectId = currentProject.value?.uuid;
  const version = urlVersion.value || item.version;
  if (projectId && version) {
    const encodedProjectId = encodeURIComponent(projectId);
    const encodedVersion = encodeURIComponent(version);
    const spdxParam = item.id ? `/${encodeURIComponent(item.id)}` : '';
    const path = `/projects/${encodedProjectId}/versions/${encodedVersion}/components${spdxParam}`;
    router.push(path);
  } else {
    console.error('Missing required data:', {projectId, version, item});
  }
};

const toggleLock = async (item: SBOM) => {
  if (!projectUuid.value) return;
  try {
    if (item.details?.isLocked) {
      await projectService.unlockSbom(projectUuid.value, item.version, item.id);
    } else {
      await projectService.lockSbom(projectUuid.value, item.version, item.id);
    }
    await appStore.refetchCurrentProject(projectUuid.value);
  } catch (error) {
    console.error('Error toggling lock:', error);
    if (isAxiosError(error) && error.response?.data?.message) {
      errorSnackbar(error.response.data.message);
    } else {
      errorSnackbar(t('Error toggling SBOM lock'));
    }
  }
};

const selectedFile = ref<File | null>(null);
const uploading = ref(false);
const fileInputRef = ref<HTMLInputElement | null>(null);

const handleFileSelect = async (event: Event) => {
  const input = event.target as HTMLInputElement;
  if (input.files && input.files.length > 0) {
    selectedFile.value = input.files[0];
    await uploadSPDXFile();
  }
};

const clearFileSelection = () => {
  selectedFile.value = null;
  if (fileInputRef.value) {
    fileInputRef.value.value = '';
  }
};

const triggerFileInput = () => {
  if (!urlVersion.value && !selectedVersion.value?.name) {
    errorSnackbar(t('Please select a channel first'));
    return;
  }
  fileInputRef.value?.click();
};
const {info: snackbar, error: errorSnackbar} = useSnackbar();
const {copyToClipboard} = useClipboard();

// Error / validation dialog state
const showErrorDialog = ref(false);
const errorDialogMessage = ref(''); // backend message or generated summary
interface ValidationIssue {
  path: string;
  message: string;
}
const errorDialogIssues = ref<ValidationIssue[]>([]); // detailed validation issues parsed
const errorDialogReqID = ref<string>('');
const errorDialogCode = ref<string>('');

const uploadSPDXFile = async () => {
  const versionName = urlVersion.value || selectedVersion.value?.name;
  if (!versionName) {
    errorSnackbar(t('Please select a channel first'));
    return;
  }

  if (!selectedFile.value) {
    errorSnackbar(t('Please select a file to upload'));
    return;
  }

  uploading.value = true;
  try {
    if (!projectUuid.value) throw new Error('No project selected');
    const result = (await projectService.uploadSBOM(
      projectUuid.value,
      versionName,
      selectedFile.value,
    )) as SpdxUploadResponse | null;
    if (result && result.docIsValid) {
      snackbar(t('SBOM uploaded successfully'));
      await appStore.refetchCurrentProject(projectUuid.value);
      clearFileSelection();
    } else if (result) {
      errorDialogReqID.value = result.reqID || '';
      errorDialogCode.value = result.code || '';
      if (result.validationFailedMessage) {
        const raw = result.validationFailedMessage.trim();
        const lines = raw
          .split('\n')
          .map((l) => l.trim())
          .filter((l) => l.length > 0);
        errorDialogIssues.value = lines.map((l) => {
          const cleaned = l.replace(/^-\s*/, '').trim();
          const idx = cleaned.indexOf(':');
          if (idx !== -1) {
            const path = cleaned.substring(0, idx).trim();
            const msg = cleaned.substring(idx + 1).trim();
            return {path, message: msg};
          }
          return {path: cleaned, message: ''};
        });
        const rawMessage = result.message || '';
        errorDialogMessage.value =
          rawMessage ||
          t('SBOM_VALIDATION_DETAILS', {
            count: errorDialogIssues.value.length,
          });
        showErrorDialog.value = true;
      } else {
        errorDialogMessage.value = result.message || t('Failed to upload SBOM');
        errorDialogIssues.value = [];
        showErrorDialog.value = true;
      }
    }
  } catch (error) {
    console.error('Error uploading SBOM:', error);
    const errUnknown = error as unknown;
    const message =
      typeof errUnknown === 'object' && errUnknown && 'message' in errUnknown
        ? (errUnknown as {message?: string}).message || ''
        : '';
    errorDialogMessage.value = message || t('Error uploading SBOM');
    errorDialogIssues.value = [];
    showErrorDialog.value = true;
  } finally {
    uploading.value = false;
  }
};

const getSbomInfoForClipboard = (item: SBOM): string => {
  return `SBOM Reference
  Project Name: ${currentProject.value?.name}
  Project Identifier: ${currentProject.value?.uuid}
  SBOM Name: ${item.details?.name}
  SBOM Identifier: ${item.id}
  Version: ${item.version}
  Upload Date: ${formatDateTime(item.details?.uploaded)} (UTC)
  Reference Timestamp: ${formatDateTime(new Date().toISOString())} (UTC)`;
};

// Filter SBOMs optionally by URL version (list populated in projectService.getProject)
const filteredList = computed<SBOM[]>(() => {
  const sboms = (currentProject.value?.sboms || []) as SBOM[];
  if (urlVersion.value) return sboms.filter((s) => s.version === urlVersion.value);
  return sboms;
});

const customFilter = useTableFilter([
  'id',
  'name',
  'version',
  'details.tag',
  'details.name',
  'details.version',
  'details.creators',
  'details.uploaded',
  'updated',
]);

const headers: DataTableHeader[] = [
  {
    key: 'actions',
    title: t('COL_ACTIONS'),
    align: 'center',
    width: 160,
    sortable: false,
  },
  {
    key: 'name',
    title: t('COL_SBOM_NAME'),
    align: 'start',
    width: 380,
    sortable: true,
  },
  {
    key: 'version',
    title: t('COL_CHAENNEL'),
    align: 'start',
    width: 220,
    sortable: true,
  },
  {
    key: 'tag',
    title: t('COL_SBOM_TAG'),
    align: 'start',
    width: 220,
    sortable: true,
  },
  {
    key: 'details.version',
    title: t('COL_SBOM_FORMAT'),
    align: 'start',
    width: 220,
    sortable: true,
  },
  {
    key: 'details.creators',
    title: t('COL_CREATOR'),
    align: 'start',
    width: 220,
    sortable: true,
  },
  {
    key: 'details.uploaded',
    title: t('COL_LAST_SBOM_UPLOADED'),
    align: 'start',
    width: 220,
    sortable: true,
  },
];

watch(
  () => currentProject.value?.versions,
  (newVersions) => {
    if (Array.isArray(newVersions) && newVersions.length > 0) {
      const firstVersion = newVersions[0];
      if (firstVersion && !selectedVersion.value?.name) {
        selectedVersion.value = {
          name: firstVersion.name,
          key: firstVersion.name,
        };
      }
    } else {
      selectedVersion.value = null;
    }
  },
  {immediate: true},
);

function getStrWithMaxLength(max: number, str: string) {
  return str.length > max ? str.slice(0, max) + '...' : str;
}

const getActionButtons = (item: SBOM): TableActionButtonsProps['buttons'] => {
  return [
    {
      icon: item.details?.isLocked ? 'mdi-lock-outline' : 'mdi-lock-open-variant-outline',
      hint: item.details?.isLocked ? t('TT_unlock_sbom') : t('TT_lock_sbom'),
      event: 'toggleLock',
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
</script>

<template>
  <div>
    <TableLayout has-tab has-title>
      <template #buttons>
        <div class="d-flex align-center gap-2">
          <input ref="fileInputRef" type="file" accept=".json,.spdx" style="display: none" @change="handleFileSelect" />
          <DCActionButton
            :text="selectedFile ? selectedFile.name : t('BTN_UPLOAD')"
            icon="mdi-upload"
            :loading="uploading"
            :disabled="!selectedVersion?.name || uploading"
            :hint="selectedVersion?.name ? t('BTN_UPLOAD') : t('SBOM_UPLOAD_DISABLED')"
            @clicked="triggerFileInput"
            v-if="currentProject">
            <template v-if="selectedFile" #prepend>
              <v-icon icon="mdi-close" size="small" class="mr-2" @click.stop="clearFileSelection" />
            </template>
          </DCActionButton>
        </div>
        <v-select
          v-model="selectedVersion"
          density="compact"
          variant="outlined"
          item-title="name"
          return-object
          :items="currentProject.versions || []"
          :label="t('LBL_UPLOAD_CHANNEL')"
          v-if="currentProject"
          hide-details
          :hint="!currentProject.versions?.length ? t('NO_VERSIONS_AVAILABLE') : undefined"
          persistent-hint
          no-data-text="No versions available"
          class="max-w-[500px]"
          :disabled="Boolean(urlVersion)" />
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
          style="max-width: 400px"></v-text-field>
      </template>
      <template #table>
        <div ref="tableSbomDeliveries" class="fill-height">
          <v-data-table
            density="compact"
            fixed-header
            :sort-by="sortItems"
            :search="search"
            :custom-filter="customFilter"
            :headers="headers"
            :items="filteredList"
            :footer-props="{
              'items-per-page-options': [10, 50, 100, -1],
            }"
            hover
            selectable
            class="striped-table fill-height cursor-pointer"
            @click:row="(_event: Event, {item}: {item: SBOM}) => showSbomStatus(item)">
            <!-- Actions Column -->
            <template v-slot:[`item.actions`]="{item}">
              <TableActionButtons
                variant="compact"
                :buttons="getActionButtons(item)"
                @toggleLock="toggleLock(item)"
                @copy="copyToClipboard(getSbomInfoForClipboard(item))" />
            </template>

            <!-- Name Column -->
            <template v-slot:[`item.name`]="{item}">
              <DDateCellWithTooltip :value="item.details?.created" /> -&nbsp;{{ item.details?.name }}
              <div>
                <span class="font-weight-bold">UUID: </span>
                <span>{{ item.id }}</span>
              </div>
              <span v-if="item.details?.isRetain" class="font-weight-bold text-[rgb(var(--v-theme-success))]">
                {{ t('SBOM_MARKED_FOR_RETENTION') }}
              </span>
            </template>

            <!-- Version Column -->
            <template v-slot:[`item.version`]="{item}">
              {{ item.version }}
            </template>

            <!-- Tag Column -->
            <template v-slot:[`item.tag`]="{item}">
              <DSpdxTagDialog
                :presetTag="item.details?.tag"
                :versionID="item.version"
                :spdxID="item.id"
                :spdxName="item.version"
                :isCliApp="true"
                v-slot="{showDialog}">
                <v-chip
                  color="labelBackgroundColor"
                  class="mr-1 mb-1 px-2 py-2"
                  label
                  :link="!item.details?.isLocked"
                  @click.stop="showDialog">
                  <v-icon class="pr-2" small color="primary" left>mdi-label</v-icon>
                  <span v-if="!item.details?.tag" class="letterSpacing text-medium-emphasis">{{
                    t('SBOM_TAG_UNKNOWN')
                  }}</span>
                  <span v-else class="letterSpacing">{{ item.details?.tag }}</span>
                </v-chip>
              </DSpdxTagDialog>
            </template>

            <!-- Format Column -->
            <template v-slot:[`item.details.version`]="{item}">
              {{ item.details?.version }}
            </template>

            <!-- Creator Column -->
            <template v-slot:[`item.details.creators`]="{item}">
              <v-tooltip open-delay="300" bottom max-width="480" content-class="dpTooltip">
                <template v-slot:activator="{props}">
                  <span v-bind="props">{{ getStrWithMaxLength(50, '' + item.details?.creators) }}</span>
                </template>
                {{ '' + item.details?.creators }}
              </v-tooltip>
            </template>

            <!-- Upload Date Column -->
            <template v-slot:[`item.details.uploaded`]="{item}">
              <DDateCellWithTooltip :value="item.details?.uploaded" />
              <br />
              <span class="text-caption" v-if="item.updated !== item.details?.uploaded">
                {{ t('Updated') }}:
                <DDateCellWithTooltip :value="item.updated" />
              </span>
            </template>
          </v-data-table>
        </div>
      </template>
    </TableLayout>
    <v-dialog v-model="showErrorDialog" width="640" scrollable>
      <DialogLayout
        :config="{
          icon: 'mdi-alert',
          iconColor: 'warning',
          title: errorDialogCode ? errorDialogCode.toUpperCase() : t('upload_error_title'),
          primaryButton: {text: t('BTN_OK')},
        }"
        @primary-action="showErrorDialog = false"
        @close="showErrorDialog = false">
        <div class="error-message" v-html="errorDialogMessage"></div>
        <div v-if="errorDialogReqID" class="reqid-row d-flex align-center mt-4">
          <div class="reqid-text flex-grow-1">
            {{ t('DESC_ERROR') }}
            <span class="font-weight-bold">{{ errorDialogReqID }}</span>
          </div>
          <DCopyClipboardButton
            class="ml-2"
            size="x-small"
            :content="errorDialogReqID"
            :hint="t('TT_COPY_TO_CLIPBOARD_REQID')" />
        </div>
        <div v-if="errorDialogIssues.length" class="details-list mt-4">
          <div class="text-caption mb-2">
            {{ t('SBOM_VALIDATION_DETAILS', {count: errorDialogIssues.length}) }}
          </div>
          <table class="validation-table">
            <thead>
              <tr>
                <th class="path-col">Path</th>
                <th>Issue</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(issue, i) in errorDialogIssues" :key="i">
                <td class="path-col">
                  <code>{{ issue.path }}</code>
                </td>
                <td>{{ issue.message }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </DialogLayout>
    </v-dialog>
  </div>
</template>

<style scoped>
.details-list {
  max-height: 360px;
  overflow-y: auto;
}

.error-message {
  white-space: pre-line;
  line-height: 1.4;
}

.validation-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 13px;
}

.validation-table thead tr {
  background: var(--v-theme-surface-variant);
}

.validation-table th,
.validation-table td {
  text-align: left;
  vertical-align: top;
  padding: 4px 8px;
  border-bottom: 1px solid var(--v-theme-outline-variant);
}

.validation-table .path-col {
  width: 38%;
  font-family: monospace;
  word-break: break-all;
}
</style>
