<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script lang="ts" setup>
import {IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import ErrorDialogConfig from '@disclosure-portal/model/ErrorDialogConfig';
import IdleInfo from '@disclosure-portal/model/IdleInfo';
import {
  JOB_STATUS_FAILURE,
  JOB_STATUS_IN_PROGRESS,
  JOB_STATUS_SUCCESS,
  JOB_TYPE_DEPARTMENT,
  JOB_TYPE_LICENSE,
  JobDto,
  RefreshRes,
} from '@disclosure-portal/model/Job';
import AdminService from '@disclosure-portal/services/admin';
import {downloadFile} from '@disclosure-portal/utils/download';
import eventBus from '@disclosure-portal/utils/eventbus';
import useSnackbar from '@shared/composables/useSnackbar';
import config from '@shared/utils/config';
import dayjs from 'dayjs';
import {computed, nextTick, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import JsonViewer from 'vue-json-viewer';

const {t} = useI18n();
const {info: snack} = useSnackbar();

// **Reaktive Variablen**
const uploadURL = ref('');
const refreshDepartmentsJob = ref<JobDto>({} as JobDto);
const isRefreshingDepartments = ref(false);
const refreshLicensesJob = ref<JobDto>({} as JobDto);
const isReloadingLicenses = ref(false);
const spdxLicensesPresent = ref(false);

const dlgConfirmationDialog = ref();
const errorDialog = ref();
const upload = ref();
const diffDlg = ref();
const showConfirmationDialog = ref(false);

// **Computed Properties**
const lrStatusText = computed(() => {
  if (refreshLicensesJob.value.status === JOB_STATUS_SUCCESS) {
    const res = refreshLicensesJob.value.customRes as RefreshRes;
    isReloadingLicenses.value = false;
    return (
      'Licenses added: ' +
      res.added +
      '\n' +
      'Licenses unchanged: ' +
      res.unchanged +
      '\n' +
      'Licenses changed: ' +
      res.changed +
      '\n' +
      'Licenses differ: ' +
      res.differences +
      '\n' +
      'Licenses errors: ' +
      res.errors +
      '\n' +
      'Total number of licenses: ' +
      res.handled +
      '/' +
      res.total +
      '\n' +
      licList('Added IDs', res.addedLics) +
      licList('Updated IDs', res.updatedLics) +
      licList('Diff IDs', res.diffLics) +
      licList('Error IDs', res.errorLics)
    );
  } else if (refreshLicensesJob.value.status === JOB_STATUS_IN_PROGRESS) {
    isReloadingLicenses.value = true;
    return 'In progress...';
  } else {
    return 'Last run failed';
  }
});

// **Methoden und Funktionen**
function licList(title: string, list: string[]): string {
  if (list.length === 0) {
    return '';
  }
  return title + ': \n' + list.map((val: string) => '   * ' + val + '\n').join('');
}

async function open() {
  await getLicenseJobStatus();
  await getDepartmentJobStatus();
}

const isProd = computed(() => {
  return config.isProd;
});

const exportSchemaKnowledgeBase = async (): Promise<void> => {
  downloadFile(
    'SchemaKnowledgeBase_' + import.meta.env.MODE + '_' + dayjs().format('YYYY-MM-DD_hh_mm_ss') + '.json',
    AdminService.exportSchemaKnowledgeBase(),
    true,
  );
};

const exportLicenseKnowledgeBase = async (): Promise<void> => {
  downloadFile(
    'LicenseKnowledgeBase_' + import.meta.env.MODE + '_' + dayjs().format('YYYY-MM-DD_hh_mm_ss') + '.json',
    AdminService.exportLicenseKnowledgeBase(),
    true,
  );
};

const uploadImportLicenseKnowledgeBase = () => {
  showConfirmationDialog.value = true;
  dlgConfirmationDialog.value?.openWithoutDetails(
    'DLG_CONFIRMATION_DESCRIPTION_IMPORT_LICENSE_KNOWLEDGE_BASE',
    'BTN_import',
  );
};

const doUploadImportLicenseKnowledgeBase = () => {
  uploadURL.value = config.SERVER_URL + '/api/v1/admin/licenses/knowledgebase/import';
  upload.value?.uploadClick();
};

const uploadImportSchemaKnowledgeBase = () => {
  showConfirmationDialog.value = true;
  dlgConfirmationDialog.value?.openWithoutDetails(
    'DLG_CONFIRMATION_DESCRIPTION_IMPORT_SCHEMA_KNOWLEDGE_BASE',
    'BTN_import',
  );
};

const doUploadImportSchemaKnowledgeBase = () => {
  uploadURL.value = config.SERVER_URL + '/api/v1/admin/schemas/knowledgebase/import';
  upload.value!.uploadClick();
};

const uploadProgress = (file: File, progress: number) => {
  const validProgress = Math.max(0, Math.min(100, Math.round(progress || 0)));
  const idleInfo = new IdleInfo(true);
  const uploadPhaseProgress = Math.round(validProgress * 0.9);
  idleInfo.message = t('PROGRESS_UPLOADING') + ' (' + file.name + ')';
  idleInfo.progress = uploadPhaseProgress;
  idleInfo.progressUnit = '%';
  eventBus.emit('on-idle', {idle: idleInfo});
};

const setProcessingProgress = (fileName: string) => {
  const idleInfo = new IdleInfo(true);
  idleInfo.message = t('PROGRESS_PROCESSING') + ' (' + fileName + ')';
  idleInfo.progress = 95; // Show 95% during processing
  idleInfo.progressUnit = '%';
  eventBus.emit('on-idle', {idle: idleInfo});
};

const fileUploaded = (_file: File, response: any) => {
  setProcessingProgress(_file.name);
  setTimeout(() => {
    eventBus.emit('on-idle', {idle: new IdleInfo(false)});

    if (response.success) {
      snack(t(response.message));
      showConfirmationDialog.value = false;
    } else {
      const d = new ErrorDialogConfig();
      d.description = t('upload_error_message');
      d.title = '' + t('VALIDATE_SCHEMA');
      d.copyDesc = true;
      d.description += response.message + ' ' + response.raw;
      d.reqId = response.reqID;
      eventBus.emit('on-error', {error: d});
      showConfirmationDialog.value = false;
    }
  }, 500); // Brief delay to show processing state
};

const refreshLicenses = () => {
  dlgConfirmationDialog.value?.openWithoutDetails('DLG_CONFIRMATION_DESCRIPTION_REFRESH_LICENSES', 'BTN_refresh');
};

const openDiffDetails = () => {
  diffDlg.value?.open();
};

const doRefreshLicenses = async (): Promise<void> => {
  isReloadingLicenses.value = true;
  spdxLicensesPresent.value = false;
  await AdminService.startJob(JOB_TYPE_LICENSE);
  await getLicenseJobStatus(true);
};

const refreshDepartments = () => {
  dlgConfirmationDialog.value?.openWithoutDetails('DLG_CONFIRMATION_DESCRIPTION_REFRESH_DEPARTMENTS', 'BTN_refresh');
};

const doRefreshDepartments = async (): Promise<void> => {
  isRefreshingDepartments.value = true;
  await AdminService.startJob(JOB_TYPE_DEPARTMENT);
  await getDepartmentJobStatus(true);
};

const onConfirm = async (config: IConfirmationDialogConfig) => {
  if (config.description === 'DLG_CONFIRMATION_DESCRIPTION_REFRESH_LICENSES') {
    await doRefreshLicenses();
    showConfirmationDialog.value = false;
  } else if (config.description === 'DLG_CONFIRMATION_DESCRIPTION_REFRESH_DEPARTMENTS') {
    await doRefreshDepartments();
    showConfirmationDialog.value = false;
  } else if (config.description === 'DLG_CONFIRMATION_DESCRIPTION_IMPORT_LICENSE_KNOWLEDGE_BASE') {
    doUploadImportLicenseKnowledgeBase();
  } else if (config.description === 'DLG_CONFIRMATION_DESCRIPTION_IMPORT_SCHEMA_KNOWLEDGE_BASE') {
    doUploadImportSchemaKnowledgeBase();
  }
};

const getLicenseJobStatus = async (justStarted = false): Promise<void> => {
  await refreshSpdxLicensesDiffPresent();
  if (!justStarted) {
    // latest refreshLicensesJob
    const serverResponse = await AdminService.getJobLatest(JOB_TYPE_LICENSE);
    refreshLicensesJob.value = serverResponse.data;
    return;
  }
  // active refreshLicensesJob
  const refreshStatus = setInterval(async () => {
    const serverResponse = await AdminService.getJobLatest(JOB_TYPE_LICENSE);
    refreshLicensesJob.value = serverResponse.data;
    if (refreshLicensesJob.value.status === JOB_STATUS_SUCCESS) {
      clearInterval(refreshStatus);
      await refreshSpdxLicensesDiffPresent();
      isReloadingLicenses.value = false;
    }
  }, 1000);
};

const getDepartmentJobStatus = async (justStarted = false): Promise<void> => {
  if (!justStarted) {
    // latest refreshDepartmentsJob
    refreshDepartmentsJob.value = (await AdminService.getJobLatest(JOB_TYPE_DEPARTMENT)).data;
    return;
  }
  // active refreshDepartmentsJob
  const refreshStatus = setInterval(async () => {
    const jobData = (await AdminService.getJobLatest(JOB_TYPE_DEPARTMENT)).data;
    refreshDepartmentsJob.value = jobData;
    if (jobData.status === JOB_STATUS_SUCCESS || jobData.status === JOB_STATUS_FAILURE) {
      clearInterval(refreshStatus);
      isRefreshingDepartments.value = false;
    }
  }, 1000);
};

const refreshSpdxLicensesDiffPresent = async () => {
  spdxLicensesPresent.value = (await AdminService.getSpdxLicensesCount()).count > 0;
};

onMounted(() => {
  nextTick(open);
});
</script>
<template>
  <v-card class="pa-4 mb-3" style="height: auto; border: none !important">
    <!-- Schema Knowledge Base -->
    <v-row>
      <v-col cols="12" xs="12">
        <v-card class="pa-3">
          <v-row>
            <v-col cols="12" xs="12">
              <h2 class="pb-3">Schema Knowledge Base</h2>
              <span class="caption">(Schemas, SchemaLabels)</span>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="3">
              <DCActionButton
                @click="exportSchemaKnowledgeBase"
                large
                class="pa-3"
                :hint="t('TT_ExportSchemaKnowledgeBase')"
                :text="t('BTN_ExportSchemaKnowledgeBase')" />
            </v-col>
            <v-col cols="5" v-if="!isProd">
              <DCActionButton
                @click="uploadImportSchemaKnowledgeBase"
                large
                class="pa-3"
                :hint="t('TT_ImportSchemaKnowledgeBase')"
                :text="t('BTN_ImportSchemaKnowledgeBase')" />
              <v-row>
                <v-col cols="12" class="caption">
                  <div class="pa-3">{{ t('IMPORT_WILL_DO_FOLLOWING') }}</div>
                  <ol>
                    <li>{{ t('ENSURE_SCHEMA_LABELS_EXIST') }}</li>
                    <li>{{ t('HARD_DELETE_ALL_SCHEMAS') }}</li>
                    <li>{{ t('UPDATE_SCHEMA_LABEL_KEYS') }}</li>
                  </ol>
                </v-col>
              </v-row>
            </v-col>
          </v-row>
        </v-card>
      </v-col>
    </v-row>

    <!-- License Knowledge Base -->
    <v-row>
      <v-col cols="12" xs="12">
        <v-card class="pa-3">
          <v-row>
            <v-col cols="12" xs="12">
              <h2 class="pb-3">License Knowledge Base</h2>
              <span class="caption">(Licenses, Classifications, PolicyRules, PolicyLabels)</span>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="3">
              <DCActionButton
                @click="exportLicenseKnowledgeBase"
                large
                class="pa-3"
                :hint="t('TT_ExportLicenseKnowledgeBase')"
                :text="t('BTN_ExportLicenseKnowledgeBase')" />
            </v-col>
            <v-col cols="5" v-if="!isProd">
              <DCActionButton
                @click="uploadImportLicenseKnowledgeBase"
                large
                class="pa-3"
                :hint="t('TT_ImportLicenseKnowledgeBase')"
                :text="t('BTN_ImportLicenseKnowledgeBase')" />
              <v-row>
                <v-col cols="12" class="caption">
                  <div class="pa-3">{{ t('IMPORT_WILL_DO_FOLLOWING') }}</div>
                  <ol>
                    <li>{{ t('ENSURE_POLICY_LABELS_EXIST') }}</li>
                    <li>{{ t('HARD_DELETE_LICENSES_CLASSIFICATIONS_POLICYRULES') }}</li>
                    <li>{{ t('IMPORT_ALL_LICENSES_CLASSIFICATIONS_POLICYRULES') }}</li>
                    <li>{{ t('UPDATE_POLICY_LABEL_KEYS') }}</li>
                  </ol>
                </v-col>
              </v-row>
            </v-col>
          </v-row>
        </v-card>
      </v-col>
    </v-row>

    <!-- Licenses Public Source -->
    <v-row>
      <v-col cols="12" xs="12">
        <v-card class="pa-3">
          <v-row>
            <v-col cols="12" xs="12">
              <h2 class="pb-3">Licenses Public Source</h2>
              <span>
                from
                <a
                  target="_blank"
                  href="https://raw.githubusercontent.com/spdx/license-list-data/master/json/licenses.json">
                  https://raw.githubusercontent.com/spdx/license-list-data/master/json/licenses.json
                </a>
              </span>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12" xs="12" md="3">
              <DCActionButton
                @click="refreshLicenses"
                large
                class="pa-3"
                :disabled="isReloadingLicenses"
                :hint="t('TT_refresh_licenses')"
                :text="t('BTN_RefreshLicenses')" />
              <v-row>
                <v-col cols="12" class="caption">
                  <div class="pa-3">{{ t('REFRESH_WILL_DO_FOLLOWING') }}</div>
                  <ol>
                    <li>{{ t('LOAD_ALL_LICENSES_PUBLIC_SERVER') }}</li>
                    <li>{{ t('UPDATE_EXISTING_LICENSES') }}</li>
                    <li>{{ t('DISPLAY_DIFFERENCE_EXISTING_LOADED_LICENSES') }}</li>
                    <li>{{ t('DISPLAY_ERRORS_LOADED_LICENSES') }}</li>
                    <li>{{ t('IMPORT_NEW_LICENSES') }}</li>
                  </ol>
                </v-col>
              </v-row>
            </v-col>
            <v-col cols="12" xs="12" md="9" v-if="refreshLicensesJob">
              <DCopyClipboardButton
                :tableButton="false"
                class="ml-1"
                :hint="t('TT_COPY_TO_CLIPBOARD')"
                :content="lrStatusText" />
              <pre class="small-text pt-2" v-if="refreshLicensesJob" v-text="lrStatusText" />
              <span v-if="spdxLicensesPresent">
                <DCActionButton
                  @click="openDiffDetails"
                  large
                  class="pa-3"
                  :disabled="isReloadingLicenses"
                  :hint="'Diff Details'"
                  :text="'Diff Details'" />
              </span>
            </v-col>
          </v-row>
        </v-card>
      </v-col>
    </v-row>

    <!-- Import Departments -->
    <v-row>
      <v-col cols="12" xs="12">
        <v-card class="pa-3">
          <v-row>
            <v-col cols="12" xs="12">
              <h2 class="pb-3">Departments</h2>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12" xs="12" md="3">
              <DCActionButton
                @click="refreshDepartments"
                large
                class="pa-3"
                :disabled="isRefreshingDepartments"
                :hint="t('TT_Refresh_Departments')"
                :text="t('BTN_Refresh_Departments')" />
              <v-row>
                <v-col cols="12" class="caption">
                  <div class="pa-3">{{ t('REFRESH_WILL_DO_FOLLOWING') }}</div>
                  <ol>
                    <li>{{ t('LOAD_ALL_DEPARTMENTS_MY_STRUCTURE_SERVER') }}</li>
                    <li>{{ t('ADD_NEW_DEPARTMENTS') }}</li>
                    <li>{{ t('UPDATE_EXISTING_DEPARTMENTS') }}</li>
                    <li>{{ t('MARK_REMOVED_DEPARTMENTS_SOFT_DELETED') }}</li>
                    <li>{{ t('DISPLAY_ERRORS') }}</li>
                  </ol>
                </v-col>
              </v-row>
            </v-col>
            <v-col cols="12" xs="12" md="9" v-if="refreshDepartmentsJob">
              <DCopyClipboardButton
                :tableButton="false"
                class="ml-1"
                :hint="t('TT_COPY_TO_CLIPBOARD')"
                :content="JSON.stringify(refreshDepartmentsJob.customRes, null, 2)" />
              <json-viewer
                :value="refreshDepartmentsJob.customRes ? refreshDepartmentsJob.customRes : {}"
                :expand-depth="2"
                aria-expanded="true"
                theme="jv-dark"
                sort />
            </v-col>
          </v-row>
        </v-card>
      </v-col>
    </v-row>

    <!-- File Upload Component -->
    <DiscoFileUpload
      ref="upload"
      :uploadTargetUrl="uploadURL"
      acceptTypes=".json"
      @reqFinished="fileUploaded"
      @reqProgress="uploadProgress" />

    <!-- Error Dialog -->
    <ErrorDialog ref="errorDialog" />

    <!-- Confirmation Dialog -->
    <ConfirmationDialog ref="dlgConfirmationDialog" v-model="showConfirmationDialog" @confirm="onConfirm" />

    <!-- Diff Dialog -->
    <DiffDialog ref="diffDlg" @refreshSpdxLicensesDiffPresent="refreshSpdxLicensesDiffPresent" />
  </v-card>
</template>
