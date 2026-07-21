<script setup lang="ts">
import {NoticeFileFormat} from '@cli/models/Version';
import {projectService} from '@cli/services/projectService';
import {downloadFile} from '@disclosure-portal/utils/View';
import {computed, reactive, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import JsonViewer3 from 'vue-json-viewer';
import {useRoute} from 'vue-router';

const {t} = useI18n();

const selectedFormat = ref<NoticeFileFormat | null>(null);
const isLoading = ref<boolean>(false);
const noticeState = reactive({
  rawContent: '',
  htmlOrPlainContent: '',
  parsedJson: null as object | unknown[] | null,
  rawJson: '',
  downloadContent: '',
  hasError: null as string | null,
  showPreview: false,
});

const route = useRoute();
const routeProjectId = computed(() => (typeof route.params.id === 'string' ? route.params.id : null));
const routeVersion = computed(() => (typeof route.params.version === 'string' ? route.params.version : null));
const routeSbom = computed(() => (typeof route.params.spdx === 'string' ? route.params.spdx : null));
const hasNoticeFileContext = computed(() => !!routeProjectId.value && !!routeVersion.value && !!routeSbom.value);

const isHTMLSelected = computed(() => selectedFormat.value === NoticeFileFormat.html);
const isTextSelected = computed(() => selectedFormat.value === NoticeFileFormat.plain);
const isJSONFormat = computed(() => selectedFormat.value === NoticeFileFormat.json);

function resetContent() {
  noticeState.rawContent = '';
  noticeState.htmlOrPlainContent = '';
  noticeState.parsedJson = null;
  noticeState.rawJson = '';
  noticeState.downloadContent = '';
  noticeState.hasError = null;
  noticeState.showPreview = false;
}

async function fetchNotice(format: NoticeFileFormat) {
  if (!hasNoticeFileContext.value) return;
  resetContent();
  selectedFormat.value = format;
  isLoading.value = true;
  try {
    const content = await projectService.downloadNoticeFile(
      routeProjectId.value!,
      routeVersion.value!,
      routeSbom.value!,
      format,
    );
    if (content == null) {
      noticeState.hasError = t('NO_DATA_AVAILABLE');
      return;
    }
    noticeState.rawContent = content;
    noticeState.showPreview = true;
    if (isTextSelected.value) {
      handlePlainFormat(content);
    } else if (isHTMLSelected.value) {
      handleHtmlFormat(content);
    } else if (isJSONFormat.value) {
      handleJsonFormat(content);
    }
  } catch (e) {
    console.error('Error loading notice file', e);
    noticeState.hasError = t('ERROR_LOADING_PROJECT');
  } finally {
    isLoading.value = false;
  }
}

function handlePlainFormat(content: string) {
  noticeState.htmlOrPlainContent = content.replace(/\n/g, '<br/>');
  noticeState.downloadContent = content;
}

function handleHtmlFormat(content: string) {
  noticeState.htmlOrPlainContent = content;
  noticeState.downloadContent = content;
}

function handleJsonFormat(content: string) {
  noticeState.rawJson = content;
  try {
    const parsed = JSON.parse(content);
    if (typeof parsed === 'object' && parsed !== null) {
      noticeState.parsedJson = parsed as object | unknown[];
      noticeState.downloadContent = JSON.stringify(parsed, null, '  ');
    } else {
      noticeState.parsedJson = null;
      noticeState.downloadContent = content;
    }
  } catch {
    noticeState.parsedJson = null;
    noticeState.downloadContent = content;
  }
}

function downloadCurrent() {
  if (!noticeState.downloadContent || !selectedFormat.value) return;
  let contentType = 'text/plain';
  let fileEnding = 'txt';
  if (isHTMLSelected.value) {
    contentType = 'text/html';
    fileEnding = 'html';
  } else if (isJSONFormat.value) {
    contentType = 'application/json';
    fileEnding = 'json';
  } else if (isTextSelected.value) {
    contentType = 'text/plain';
    fileEnding = 'txt';
  }
  const versionName = routeVersion.value || 'version';
  const projectName = routeProjectId.value || 'project';
  const filename = `${projectName}_${versionName}_notice.${fileEnding}`;
  downloadFile(noticeState.downloadContent, filename, contentType);
}

watch(
  () => hasNoticeFileContext.value,
  (val) => {
    if (val && !selectedFormat.value) {
      fetchNotice(NoticeFileFormat.html);
    }
  },
  {immediate: true},
);
</script>

<template>
  <TableLayout has-title has-tab>
    <template #description>
      <template v-if="hasNoticeFileContext">
        <Stack direction="row" class="mt-0 gap-0">
          <div>
            <v-btn
              size="small"
              :variant="isHTMLSelected ? 'tonal' : 'text'"
              :class="{active: isHTMLSelected}"
              class="ma-2 text-none card-border ml-0"
              @click="fetchNotice(NoticeFileFormat.html)">
              <v-icon color="primary">mdi-code-brackets</v-icon>
              HTML
            </v-btn>
            <v-btn
              size="small"
              :variant="isTextSelected ? 'tonal' : 'text'"
              :class="{active: isTextSelected}"
              class="ma-2 text-none card-border"
              @click="fetchNotice(NoticeFileFormat.plain)">
              <v-icon color="primary">mdi-format-text</v-icon>
              Plain Text
            </v-btn>
            <v-btn
              size="small"
              :variant="isJSONFormat ? 'tonal' : 'text'"
              :class="{active: isJSONFormat}"
              class="ma-2 text-none card-border"
              @click="fetchNotice(NoticeFileFormat.json)">
              <v-icon color="primary">mdi-code-json</v-icon>
              JSON
            </v-btn>
          </div>
        </Stack>
      </template>
      <template v-else>
        <v-alert type="info" variant="tonal" density="compact">{{ t('SELECT_PROJECT_TOKEN') }}</v-alert>
      </template>
    </template>
    <template #table>
      <v-card class="card-border fill-height pa-4 overflow-auto" v-if="noticeState.showPreview && !isLoading">
        <Stack direction="row" justify="between" align="center" class="mb-4">
          <h3 class="d-subtitle-2 d-secondary-text mt-2 mb-2">{{ t('PREVIEW') }}</h3>
          <div class="d-flex gap-2">
            <DCopyClipboardButton
              :tableButton="true"
              :hint="t('TT_noticeCopyText')"
              :content="noticeState.downloadContent" />
            <DCActionButton
              large
              icon="mdi-download"
              :text="t('BTN_DOWNLOAD')"
              :hint="t('TT_download_notice')"
              class="mr-2"
              @click="downloadCurrent" />
          </div>
        </Stack>
        <Stack v-if="isJSONFormat" direction="row" class="gap-0">
          <div style="min-height: 400px; max-height: 400px; overflow: auto; width: 100%">
            <template v-if="noticeState.parsedJson">
              <JsonViewer3
                :value="noticeState.parsedJson"
                :expand-depth="1"
                aria-expanded="true"
                theme="jv-dark"
                sort />
            </template>
            <template v-else>
              <pre class="json-fallback" style="white-space: pre-wrap; font-family: monospace; font-size: 13px">{{
                noticeState.rawJson
              }}</pre>
            </template>
          </div>
        </Stack>
        <Stack v-else direction="row" class="gap-0">
          <div class="w-1/2 overflow-auto">
            <div
              class="d-text html-notice-file pt-2"
              v-html="noticeState.htmlOrPlainContent"
              v-if="noticeState.htmlOrPlainContent && noticeState.htmlOrPlainContent.length > 0" />
          </div>
        </Stack>
      </v-card>
      <div v-else-if="isLoading" class="d-flex my-6 justify-center">
        <v-progress-circular indeterminate color="primary" />
      </div>
      <v-alert v-else-if="noticeState.hasError" type="error" variant="tonal" density="compact" class="mb-4">{{
        noticeState.hasError
      }}</v-alert>
      <v-alert
        v-else-if="hasNoticeFileContext && !noticeState.showPreview && !isLoading"
        type="info"
        variant="tonal"
        density="compact">
        {{ t('NO_DATA_AVAILABLE') }}
      </v-alert>
    </template>
  </TableLayout>
</template>
