<script lang="ts" setup>
import {getApi} from '@disclosure-portal/api';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import DCloseButton from '@shared/components/disco/DCloseButton.vue';
import DIconButton from '@shared/components/disco/DIconButton.vue';
import TableLayout from '@shared/layouts/TableLayout.vue';
import {DiscoForm} from '@disclosure-portal/types/discobasics';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {DataTableHeader, SortItem} from '@shared/types/table';
import {computed, onMounted, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRoute, useRouter} from 'vue-router';

interface LocaleEntry {
  key: string;
  translation: string;
}

interface I18nLocaleResponse {
  localeCode: string;
  displayName: string;
  nativeName: string;
  isDefault: boolean;
  scope: string;
  entryCount: number;
  entries: Record<string, string>;
  fallbackUsed: boolean;
}

interface I18nImportIssue {
  fileName: string;
  key?: string;
  code: string;
  message: string;
}

interface I18nImportResponse {
  success: boolean;
  validationPassed: boolean;
  locale: string;
  filesProcessed: number;
  totalKeysParsed: number;
  appended: number;
  updated: number;
  unchanged: number;
  errors?: I18nImportIssue[];
}

interface I18nLocaleListItem {
  localeCode: string;
}

const {t} = useI18n();
const route = useRoute();
const router = useRouter();
const {dashboardCrumbs, ...breadcrumbs} = useBreadcrumbsStore();
const {api} = getApi();
const knownLocalesForGlobalDelete = ['en', 'de'];

const search = ref('');
const sortItems = ref<SortItem[]>([{key: 'key', order: 'asc'}]);
const entries = ref<LocaleEntry[]>([]);
const editRowKey = ref<string | null>(null);
const draftTranslation = ref('');
const showAddEntryDialog = ref(false);
const newEntryKey = ref('');
const newEntryTranslation = ref('');
const addEntryToAllLocales = ref(false);
const addEntryError = ref('');
const addEntryFormRef = ref<DiscoForm | null>(null);

const addEntryRules = {
  key: [(v: string) => !!v?.trim() || t('VALIDATION_required')],
  translation: [(v: string) => !!v?.trim() || t('VALIDATION_required')],
};
const showDeleteDialog = ref(false);
const deleteEntryKey = ref<string | null>(null);
const deleteGlobally = ref(false);
const importInputRef = ref<HTMLInputElement | null>(null);
const isImporting = ref(false);
const isExporting = ref(false);
const importToAllLocales = ref(false);
const showImportDialog = ref(false);
const showImportResultDialog = ref(false);
const importResult = ref<I18nImportResponse | null>(null);
const isLoading = ref(false);
const actionError = ref('');
const localeMeta = ref<I18nLocaleResponse | null>(null);

const localeCode = computed(() => String(route.params.localeCode || '').toLowerCase());

const selectedLocaleLabel = computed(() => {
  if (localeMeta.value?.nativeName) {
    return localeMeta.value.nativeName;
  }
  const translated = t(`LANG_${localeCode.value}`);
  if (translated !== `LANG_${localeCode.value}`) {
    return translated;
  }
  return localeCode.value.toUpperCase();
});

const pageTitle = computed(() => {
  return `${t('ADMIN_I18N_PAGE_TITLE')} - ${selectedLocaleLabel.value}`;
});

const pageDescription = computed(() => t('ADMIN_I18N_DETAIL_PAGE_DESCRIPTION'));
const importResultHasErrors = computed(() => (importResult.value?.errors?.length || 0) > 0);

const normalizeValue = (value: unknown): string => {
  if (value === null || value === undefined) {
    return '';
  }
  if (typeof value === 'string') {
    return value;
  }
  if (typeof value === 'number' || typeof value === 'boolean') {
    return String(value);
  }
  try {
    return JSON.stringify(value);
  } catch {
    return String(value);
  }
};

const getRowTranslation = (item: {raw?: LocaleEntry} | LocaleEntry): string => {
  const row = getRow(item);
  return row.translation;
};

const getRow = (item: {raw?: LocaleEntry} | LocaleEntry): LocaleEntry => {
  return 'raw' in item && item.raw ? item.raw : item as LocaleEntry;
};

const isEditingRow = (item: {raw?: LocaleEntry} | LocaleEntry): boolean => {
  return getRow(item).key === editRowKey.value;
};

const isMultilineValue = (value: string): boolean => value.includes('\n');

const extractApiErrorMessage = (error: any): string => {
  const message = error?.response?.data?.message;
  if (typeof message === 'string' && message.trim()) {
    return message;
  }
  return t('ERROR_500_TITLE');
};

const syncEntries = (data: Record<string, string> = {}) => {
  entries.value = Object.entries(data).map(([key, value]) => ({
    key,
    translation: normalizeValue(value),
  }));
};

const fetchLocale = async () => {
  if (!localeCode.value) {
    entries.value = [];
    return;
  }

  isLoading.value = true;
  actionError.value = '';
  try {
    const response = await api.get<I18nLocaleResponse>(`/api/v1/i18n/${encodeURIComponent(localeCode.value)}`);
    localeMeta.value = response.data;
    syncEntries(response.data?.entries || {});
  } catch {
    actionError.value = t('ERROR_500_TITLE');
    entries.value = [];
  } finally {
    isLoading.value = false;
  }
};

const upsertTranslation = async (key: string, value: string): Promise<boolean> => {
  actionError.value = '';
  try {
    await api.put(`/api/v1/i18n/${encodeURIComponent(localeCode.value)}/${encodeURIComponent(key)}`, {
      value,
      description: '',
    });
    return true;
  } catch {
    actionError.value = t('ERROR_500_TITLE');
    return false;
  }
};

const exportAsJson = async () => {
  isExporting.value = true;
  actionError.value = '';
  try {
		const response = await api.get(`/api/v1/i18n/export/${encodeURIComponent(localeCode.value)}`, {
      responseType: 'blob',
    });

    const fileName = `locale.${localeCode.value}.json`;
    const url = URL.createObjectURL(response.data);
    const anchor = document.createElement('a');
    anchor.href = url;
    anchor.download = fileName;
    document.body.appendChild(anchor);
    anchor.click();
    document.body.removeChild(anchor);
    URL.revokeObjectURL(url);
  } catch {
    actionError.value = t('ERROR_500_TITLE');
  } finally {
    isExporting.value = false;
  }
};

const openImportPicker = () => {
  showImportDialog.value = true;
};

const selectImportFiles = () => {
  importInputRef.value?.click();
};

const closeImportDialog = () => {
  showImportDialog.value = false;
  importToAllLocales.value = false;
};

const resetImportResultDialog = () => {
  showImportResultDialog.value = false;
  importResult.value = null;
};

const onImportFilesSelected = async (event: Event) => {
  const input = event.target as HTMLInputElement;
  const selectedFiles = input.files;
  if (!selectedFiles || selectedFiles.length === 0) {
    return;
  }

  isImporting.value = true;
  actionError.value = '';
  importResult.value = null;

  try {
    // Determine which locales to import to
    let targetLocales = [localeCode.value];
    if (importToAllLocales.value) {
      const localeResponse = await api.get<I18nLocaleListItem[]>('/api/v1/i18n');
      targetLocales = Array.from(new Set(localeResponse.data.map((item) => item.localeCode)));
    }

    // Import files to each target locale
    const importPromises = targetLocales.map(async (targetLocale) => {
      const formData = new FormData();
      Array.from(selectedFiles).forEach((file) => {
        formData.append('files', file);
      });

      return api.post<I18nImportResponse>(
        `/api/v1/i18n/${encodeURIComponent(targetLocale)}/import`,
        formData,
        {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
        },
      );
    });

    const responses = await Promise.all(importPromises);
    const currentLocaleResponse = responses[0];

    importResult.value = currentLocaleResponse.data;
    showImportResultDialog.value = true;
    showImportDialog.value = false;
    await fetchLocale();
  } catch (error: any) {
    const responseData = error?.response?.data as I18nImportResponse | undefined;
    if (responseData) {
      importResult.value = responseData;
      showImportResultDialog.value = true;
    } else {
      actionError.value = extractApiErrorMessage(error);
    }
  } finally {
    if (input) {
      input.value = '';
    }
    importToAllLocales.value = false;
    isImporting.value = false;
  }
};

const startEdit = (item: LocaleEntry) => {
  editRowKey.value = item.key;
  draftTranslation.value = item.translation;
};

const cancelEdit = () => {
  editRowKey.value = null;
  draftTranslation.value = '';
};

const resetAddEntryDialog = () => {
  showAddEntryDialog.value = false;
  newEntryKey.value = '';
  newEntryTranslation.value = '';
  addEntryToAllLocales.value = false;
  addEntryError.value = '';
  addEntryFormRef.value?.reset();
};

const addEntry = async () => {
  const key = newEntryKey.value.trim();

  if (!key) {
    addEntryError.value = t('VALIDATION_required');
    return;
  }

  if (!newEntryTranslation.value.trim()) {
    addEntryError.value = t('VALIDATION_required');
    return;
  }

  if (entries.value.some((entry) => entry.key === key)) {
    addEntryError.value = t('ERROR_LABEL_USED');
    return;
  }

  actionError.value = '';
  addEntryError.value = '';
  try {
    let targetLocales = [localeCode.value];

    if (addEntryToAllLocales.value) {
      const localeResponse = await api.get<I18nLocaleListItem[]>('/api/v1/i18n');
      targetLocales = Array.from(new Set((localeResponse.data || [])
        .map((item) => String(item.localeCode || '').toLowerCase())
        .filter((code) => !!code)));
      if (!targetLocales.includes(localeCode.value)) {
        targetLocales.push(localeCode.value);
      }
    }

    await Promise.all(
      targetLocales.map((targetLocale) =>
        api.put(`/api/v1/i18n/${encodeURIComponent(targetLocale)}/${encodeURIComponent(key)}`, {
          value: newEntryTranslation.value,
          description: '',
        })),
    );
  } catch {
    addEntryError.value = t('ERROR_500_TITLE');
    return;
  }

  entries.value = [
    {
      key,
      translation: newEntryTranslation.value,
    },
    ...entries.value,
  ];

  resetAddEntryDialog();
};

const saveEdit = async () => {
  if (!editRowKey.value) {
    return;
  }

  const key = editRowKey.value;
  const persisted = await upsertTranslation(key, draftTranslation.value);
  if (!persisted) {
    return;
  }

  entries.value = entries.value.map((item) => {
    if (item.key !== key) {
      return item;
    }

    return {
      ...item,
      translation: draftTranslation.value,
    };
  });

  cancelEdit();
};

const onDeleteConfirm = async () => {
  if (!deleteEntryKey.value) {
    return;
  }

  const keyToDelete = deleteEntryKey.value;

  actionError.value = '';
  const targetLocales = deleteGlobally.value
    ? Array.from(new Set([localeCode.value, ...knownLocalesForGlobalDelete]))
    : [localeCode.value];

  const deleteResults = await Promise.allSettled(
    targetLocales.map((targetLocale) =>
      api.delete(`/api/v1/i18n/${encodeURIComponent(targetLocale)}/${encodeURIComponent(keyToDelete)}`),
    ),
  );

  const currentLocaleIndex = targetLocales.findIndex((code) => code === localeCode.value);
  const currentLocaleResult = deleteResults[currentLocaleIndex];
  if (currentLocaleResult?.status === 'rejected') {
    actionError.value = extractApiErrorMessage(currentLocaleResult.reason);
    return;
  }

  const failedSecondaryLocales = deleteResults
    .map((result, index) => ({result, locale: targetLocales[index]}))
    .filter((item) => item.locale !== localeCode.value && item.result.status === 'rejected')
    .map((item) => item.locale);

  if (failedSecondaryLocales.length > 0) {
    actionError.value = `Deleted in ${localeCode.value.toUpperCase()}, but failed in: ${failedSecondaryLocales.join(', ').toUpperCase()}`;
  }

  entries.value = entries.value.filter((entry) => entry.key !== keyToDelete);

  if (editRowKey.value === keyToDelete) {
    cancelEdit();
  }

  showDeleteDialog.value = false;
  deleteEntryKey.value = null;
  deleteGlobally.value = false;
};

const resetDeleteDialog = () => {
  showDeleteDialog.value = false;
  deleteEntryKey.value = null;
  deleteGlobally.value = false;
};

const openDeleteDialog = (item: LocaleEntry) => {
  deleteEntryKey.value = item.key;
  deleteGlobally.value = false;
  showDeleteDialog.value = true;
};

const headers = computed((): DataTableHeader[] => [
  {
    title: t('KEY'),
    align: 'start',
    value: 'key',
    sortable: true,
    width: 360,
    minWidth: 280,
  },
  {
    title: t('VALUE'),
    align: 'start',
    value: 'translation',
    sortable: true,
    minWidth: 500,
  },
  {
    title: t('COL_ACTIONS').trim(),
    align: 'center',
    value: 'actions',
    sortable: false,
    width: 120,
  },
]);

const initBreadcrumbs = () => {
  breadcrumbs.setCurrentBreadcrumbs([
    ...dashboardCrumbs,
    {
      title: t('ADMIN_I18N_PAGE_TITLE'),
      disabled: false,
      href: '/dashboard/admin/i18n',
    },
    {
      title: selectedLocaleLabel.value,
    },
  ]);
};

onMounted(async () => {
  if (!localeCode.value) {
    router.replace({name: 'I18nAdmin'});
    return;
  }
  await fetchLocale();
  initBreadcrumbs();
});

watch(
  () => route.params.localeCode,
  async () => {
    if (!localeCode.value) {
      router.replace({name: 'I18nAdmin'});
      return;
    }
    cancelEdit();
    resetAddEntryDialog();
    await fetchLocale();
    initBreadcrumbs();
  },
);
</script>

<template>
  <TableLayout>
    <template #description>
      <h1 class="text-h5">{{ pageTitle }}</h1>
      <p class="text-body-2 text-medium-emphasis mt-1">{{ pageDescription }}</p>
    </template>
    <template #buttons>
      <span class="text-h6">{{ t('ADMIN_I18N_DETAIL_SECTION_TITLE') }}</span>
      <DCActionButton
        large
        class="mx-2"
        icon="mdi-plus"
        :hint="`${t('BTN_ADD')} ${t('KEY')}`"
        :text="t('BTN_ADD')"
        :disabled="isLoading || isImporting"
        @click="showAddEntryDialog = true" />
      <DCActionButton
        large
        class="mx-2"
        icon="mdi-file-export-outline"
        :text="t('BTN_EXPORT_JSON')"
        :disabled="isLoading || isExporting || isImporting"
        @click="exportAsJson" />
      <DCActionButton
        large
        class="mx-2"
        icon="mdi-file-upload-outline"
        :text="t('BTN_UPLOAD_JSON')"
        :disabled="isLoading || isExporting || isImporting"
        @click="openImportPicker" />
      <v-spacer></v-spacer>
      <DSearchField v-model="search" />
      <input
        ref="importInputRef"
        type="file"
        accept="application/json,.json"
        multiple
        class="d-none"
        @change="onImportFilesSelected" />
    </template>
    <template #table>
      <div class="fill-height">
        <v-data-table
          density="compact"
          class="striped-table fill-height"
          fixed-header
          :headers="headers"
          :items="entries"
          item-value="key"
          :search="search"
          :sort-by="sortItems"
          :loading="isLoading"
          :items-per-page="25">
          <template #[`item.translation`]="{item}">
            <div v-if="isEditingRow(item)">
              <v-textarea
                v-if="isMultilineValue(draftTranslation)"
                v-model="draftTranslation"
                auto-grow
                rows="3"
                density="compact"
                variant="outlined"
                hide-details />
              <v-text-field
                v-else
                v-model="draftTranslation"
                density="compact"
                variant="outlined"
                hide-details />
            </div>
            <span v-else class="text-body-2" style="white-space: pre-wrap">{{ getRowTranslation(item) }}</span>
          </template>
          <template #[`item.actions`]="{item}">
            <div class="d-flex justify-center ga-1">
              <template v-if="isEditingRow(item)">
                <DIconButton icon="mdi-check" color="success" @clicked="saveEdit" />
                <DIconButton icon="mdi-close" color="secondary" @clicked="cancelEdit" />
              </template>
              <template v-else>
                <DIconButton icon="mdi-pencil" :hint="t('BTN_EDIT')" @clicked="startEdit(getRow(item))" />
                <DIconButton
                  icon="mdi-delete"
                  color="error"
                  :hint="t('BTN_DELETE')"
                  @clicked="openDeleteDialog(getRow(item))" />
              </template>
            </div>
          </template>
        </v-data-table>
        <p v-if="actionError" class="text-error text-caption mt-2 mb-0">{{ actionError }}</p>
      </div>
    </template>
  </TableLayout>

  <v-dialog v-model="showAddEntryDialog" max-width="600px" persistent>
    <v-form ref="addEntryFormRef">
      <v-card class="pa-8">
        <v-card-title>
          <v-row>
            <v-col cols="10" class="d-flex align-center">
              <span class="text-h5">{{ `${t('BTN_ADD')} ${t('KEY')}` }}</span>
            </v-col>
            <v-col cols="2" class="px-0 text-right">
              <DCloseButton @click="resetAddEntryDialog" />
            </v-col>
          </v-row>
        </v-card-title>
        <v-card-text>
          <v-row>
            <v-col cols="12" class="px-0">
              <v-text-field
                v-model="newEntryKey"
                variant="outlined"
                class="required"
                :rules="addEntryRules.key"
                :label="t('KEY')"
                placeholder="e.g. ADMIN_I18N_PAGE_TITLE"
                hide-details="auto"
                autofocus />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12" class="px-0">
              <v-textarea
                v-model="newEntryTranslation"
                auto-grow
                rows="3"
                variant="outlined"
                class="required"
                :rules="addEntryRules.translation"
                :label="t('VALUE')"
                placeholder="e.g. Internationalization"
                hide-details="auto" />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12" class="px-0">
              <v-checkbox
                v-model="addEntryToAllLocales"
                density="compact"
                :label="t('ADMIN_I18N_ADD_TO_ALL_LOCALES')"
                hide-details
                class="mt-2" />
              <p v-if="addEntryError" class="text-error text-caption mt-3 mb-0">{{ addEntryError }}</p>
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions class="justify-end">
          <v-btn size="small" variant="text" color="primary" class="mr-5" @click="resetAddEntryDialog">
            {{ t('BTN_CANCEL') }}
          </v-btn>
          <v-btn size="small" variant="flat" color="primary" class="mr-1" @click="addEntry">
            {{ t('BTN_ADD') }}
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-form>
  </v-dialog>

  <v-dialog v-model="showDeleteDialog" max-width="500px">
    <v-card class="pa-8">
      <v-card-title>
        <v-row>
          <v-col cols="10" class="d-flex align-center">
            <span class="text-h5">{{ t('DLG_CONFIRMATION_TITLE') }}</span>
          </v-col>
          <v-col cols="2" class="px-0 text-right">
            <DCloseButton @click="resetDeleteDialog" />
          </v-col>
        </v-row>
      </v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="12" class="px-0">
            {{ t('DLG_CONFIRMATION_DESCRIPTION') }}<strong>{{ deleteEntryKey }}</strong>?
          </v-col>
        </v-row>
        <v-row>
          <v-col cols="12" class="px-0">
            <v-checkbox
              v-model="deleteGlobally"
              density="compact"
              :label="t('ADMIN_I18N_DELETE_FROM_ALL_LOCALES')"
              hide-details
              class="mt-2" />
          </v-col>
        </v-row>
        <v-row>
          <div class="f-modal-alert">
            <div class="f-modal-icon f-modal-warning scaleWarning">
              <span class="f-modal-body pulseWarningIns"></span>
              <span class="f-modal-dot pulseWarningIns"></span>
            </div>
          </div>
        </v-row>
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn size="small" variant="text" color="primary" class="mr-5" @click="resetDeleteDialog">
          {{ t('BTN_CANCEL') }}
        </v-btn>
        <v-btn size="small" variant="flat" color="primary" class="mr-1" @click="onDeleteConfirm">
          {{ t('BTN_DELETE') }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog v-model="showImportDialog" max-width="500px" persistent>
    <v-card class="pa-8">
      <v-card-title>
        <v-row>
          <v-col cols="10" class="d-flex align-center">
            <span class="text-h5">{{ t('BTN_UPLOAD_JSON') }}</span>
          </v-col>
          <v-col cols="2" class="px-0 text-right">
            <DCloseButton @click="closeImportDialog" />
          </v-col>
        </v-row>
      </v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="12" class="px-0">
            <v-checkbox
              v-model="importToAllLocales"
              density="compact"
              :label="t('ADMIN_I18N_IMPORT_TO_ALL_LOCALES')"
              hide-details
              class="mt-2" />
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn size="small" variant="text" color="primary" class="mr-5" @click="closeImportDialog">
          {{ t('BTN_CANCEL') }}
        </v-btn>
        <v-btn size="small" variant="flat" color="primary" class="mr-1" :disabled="isImporting" @click="selectImportFiles">
          {{ t('BTN_UPLOAD_JSON') }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog v-model="showImportResultDialog" max-width="800px">
    <v-card class="pa-8">
      <v-card-title>
        <v-row>
          <v-col cols="10" class="d-flex align-center">
            <span class="text-h5">{{ t('ADMIN_I18N_IMPORT_RESULT_TITLE') }}</span>
          </v-col>
          <v-col cols="2" class="px-0 text-right">
            <DCloseButton @click="resetImportResultDialog" />
          </v-col>
        </v-row>
      </v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="12" class="px-0">
            <p class="mb-2"><strong>{{ t('ADMIN_I18N_IMPORT_RESULT_LOCALE') }}:</strong> {{ importResult?.locale }}</p>
            <p class="mb-2"><strong>{{ t('ADMIN_I18N_IMPORT_RESULT_FILES') }}:</strong> {{ importResult?.filesProcessed ?? 0 }}</p>
            <p class="mb-2"><strong>{{ t('ADMIN_I18N_IMPORT_RESULT_KEYS') }}:</strong> {{ importResult?.totalKeysParsed ?? 0 }}</p>
            <p class="mb-2"><strong>{{ t('ADMIN_I18N_IMPORT_RESULT_APPENDED') }}:</strong> {{ importResult?.appended ?? 0 }}</p>
            <p class="mb-2"><strong>{{ t('ADMIN_I18N_IMPORT_RESULT_UPDATED') }}:</strong> {{ importResult?.updated ?? 0 }}</p>
            <p class="mb-2"><strong>{{ t('ADMIN_I18N_IMPORT_RESULT_UNCHANGED') }}:</strong> {{ importResult?.unchanged ?? 0 }}</p>

            <v-alert
              v-if="importResultHasErrors"
              type="error"
              variant="tonal"
              class="mt-4"
              :title="t('ADMIN_I18N_IMPORT_RESULT_ERROR_TITLE')">
              <ul class="mb-0 pl-4">
                <li v-for="(issue, index) in importResult?.errors || []" :key="`${issue.fileName}-${issue.key || ''}-${index}`">
                  {{ issue.fileName }}
                  <span v-if="issue.key"> ({{ issue.key }})</span>
                  : {{ issue.message }}
                </li>
              </ul>
            </v-alert>
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-actions class="justify-end">
        <v-btn size="small" variant="flat" color="primary" class="mr-1" @click="resetImportResultDialog">
          {{ t('BTN_CLOSE') }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
