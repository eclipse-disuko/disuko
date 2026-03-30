// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import {PolicyLabels} from '@disclosure-portal/constants/policyLabels';
import {NoticeFileFormat, SbomStats} from '@disclosure-portal/model/VersionDetails';
import ProjectService from '@disclosure-portal/services/projects';
import VersionService from '@disclosure-portal/services/version';
import {useAppStore} from '@disclosure-portal/stores/app';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useSbomStore} from '@disclosure-portal/stores/sbom.store';
import {downloadFile, getIconColorScanRemarkLevel} from '@disclosure-portal/utils/View';
import {TOOLTIP_OPEN_DELAY_IN_MS} from '@shared/utils/constant';
import {computed, nextTick, onMounted, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import JsonViewer3 from 'vue-json-viewer';
import {VIcon} from 'vuetify/components';

const {t} = useI18n();
const appStore = useAppStore();
const projectStore = useProjectStore();
const sbomStore = useSbomStore();
const currentProject = computed(() => projectStore.currentProject!);
const currentVersionId = computed(() => sbomStore.getCurrentVersion._key);
const currentProjectId = computed(() => currentProject.value._key);
const spdx = computed(() => sbomStore.selectedSpdx);
const spdxFileHistory = computed(() => sbomStore.getChannelSpdxs);
const labelTools = computed(() => appStore.getLabelsTools);
const isVehicleProject = computed(() =>
  currentProject.value.policyLabels.some(
    (lbl) => labelTools.value.policyLabelsMap[lbl]?.name === PolicyLabels.VEHICLE_PLATFORM,
  ),
);
const isHTMLSelected = computed(() => selectedFormat.value === NoticeFileFormat.html);
const isTextSelected = computed(() => selectedFormat.value === NoticeFileFormat.plain);
const isJSONFormat = computed(() => selectedFormat.value === NoticeFileFormat.json);

const selectedFormat = ref<NoticeFileFormat | null>(null);
const jsonContent = ref('');
const htmlOrPlainContent = ref('');
const downloadContent = ref('');
const showPreview = ref(false);
const dataAreLoaded = ref(false);
const showProgressBar = ref(false);
const sbomStats = ref<SbomStats>({} as SbomStats);
const projectSettingsDialog = ref();

const loadContent = async (format: NoticeFileFormat) => {
  showPreview.value = false;
  selectedFormat.value = format;

  if (!currentProjectId.value || !currentVersionId.value || spdxFileHistory.value.length < 1 || !spdx.value) {
    return;
  }

  const response = await ProjectService.downloadNoticeFileInFormatForSbom(
    format,
    currentProjectId.value,
    currentVersionId.value,
    spdx.value._key,
  );

  updateContentsByType(response.data);
};

const updateContentsByType = (content: string) => {
  if (content === '') {
    return;
  }
  showPreview.value = true;

  switch (selectedFormat.value) {
    case NoticeFileFormat.plain:
      htmlOrPlainContent.value = content.replaceAll('\n', '<br>');
      downloadContent.value = content;
      break;
    case NoticeFileFormat.html:
      htmlOrPlainContent.value = content;
      downloadContent.value = content;
      break;
    case NoticeFileFormat.json:
      jsonContent.value = content;
      downloadContent.value = JSON.stringify(content, null, '  ');
      break;
  }
};

const reload = async (forceReload = false) => {
  dataAreLoaded.value = false;
  showProgressBar.value = false;

  if (!forceReload && dataAreLoaded.value) return;

  if (spdxFileHistory.value.length === 0) {
    showPreview.value = false;
    dataAreLoaded.value = true;
    return;
  }

  await loadContent(NoticeFileFormat.html);

  const delayTimeout = setTimeout(() => {
    showProgressBar.value = true;
  }, 600);

  sbomStats.value = (
    await VersionService.getSBOMStats(currentProjectId.value, currentVersionId.value, spdx.value._key, true)
  ).data;

  clearTimeout(delayTimeout);
  dataAreLoaded.value = true;
  showProgressBar.value = false;
};

const loadStyledHTML = () => {
  loadContent(NoticeFileFormat.html);
};

const loadPlainText = () => {
  loadContent(NoticeFileFormat.plain);
};

const loadJSON = () => {
  loadContent(NoticeFileFormat.json);
};

const getFileName = computed(() => {
  let ending = '';
  switch (selectedFormat.value) {
    case NoticeFileFormat.plain:
      ending = 'txt';
      break;
    case NoticeFileFormat.html:
      ending = 'html';
      break;
    case NoticeFileFormat.json:
      ending = 'json';
      break;
  }
  return `${currentProject.value.name}_${currentVersionId.value}_notice.${ending}`;
});

const getContentType = () => {
  switch (selectedFormat.value) {
    case NoticeFileFormat.plain:
      return 'text/plain';
    case NoticeFileFormat.html:
      return 'text/html';
    case NoticeFileFormat.json:
      return 'application/json';
  }
};

const downloadNoticeFile = () => {
  downloadFile(downloadContent.value, getFileName.value, getContentType()!);
};

const showProjectSettingsDialog = async () => {
  if (projectSettingsDialog.value) {
    projectSettingsDialog.value?.showDialog(currentProject.value);
    projectSettingsDialog.value.activeTab = 'owner';
    await nextTick();
    const element = document.querySelector('#thirdparty-address') as HTMLElement;
    if (element) {
      element.focus();
      element.scrollIntoView({behavior: 'smooth', block: 'center'});
    }
  }
};

watch(
  () => spdx.value,
  async () => {
    await reload(true);
  },
);

watch(
  () => currentProject.value,
  async () => {
    await reload(true);
  },
);

onMounted(async () => {
  await reload(true);
});
</script>

<template>
  <TableLayout has-title has-tab>
    <template #description>
      <template v-if="currentVersionId && spdxFileHistory.length > 0">
        <v-row
          class="shrink"
          style="height: 52px"
          v-if="showProgressBar && !dataAreLoaded"
          align="center"
          justify="center">
          <v-col cols="12" sm="8" md="6">
            <v-progress-linear indeterminate color="primary" height="3"></v-progress-linear>
          </v-col>
        </v-row>
        <v-row v-if="dataAreLoaded && sbomStats">
          <v-col class="d-flex ga-2">
            <v-tooltip
              v-if="sbomStats.PolicyState.NoAssertion > 0"
              :open-delay="TOOLTIP_OPEN_DELAY_IN_MS"
              location="bottom"
              content-class="dpTooltip">
              <template v-slot:activator="{props}">
                <v-btn class="text-none card-border font-weight-light" variant="text" size="small" v-bind="props">
                  <v-icon color="red" icon="mdi-lightning-bolt-circle" class="mr-2"></v-icon>
                  {{ `${sbomStats.PolicyState.NoAssertion} ${t('UNASSERTED')}` }}
                </v-btn>
              </template>
              <span>{{ t('TT_UNASSERTED') }}</span>
            </v-tooltip>
            <v-tooltip
              v-if="sbomStats.PolicyState.Denied > 0"
              :open-delay="TOOLTIP_OPEN_DELAY_IN_MS"
              location="bottom"
              content-class="dpTooltip">
              <template v-slot:activator="{props}">
                <v-btn class="text-none card-border font-weight-light" variant="text" size="small" v-bind="props">
                  <v-icon color="red" icon="mdi-minus-circle" class="mr-2"></v-icon>
                  {{ `${sbomStats.PolicyState.Denied} ${t('DENIED')}` }}
                </v-btn>
              </template>
              <span>{{ t('TT_DENIED') }}</span>
            </v-tooltip>
            <v-tooltip
              v-if="sbomStats.PolicyState.Warned > 0"
              :open-delay="TOOLTIP_OPEN_DELAY_IN_MS"
              location="bottom"
              content-class="dpTooltip">
              <template v-slot:activator="{props}">
                <v-btn class="text-none card-border font-weight-light" variant="text" size="small" v-bind="props">
                  <v-icon color="warning" icon="mdi-alert" class="mr-2"></v-icon>
                  {{ `${sbomStats.PolicyState.Warned} ${t('WARNED')}` }}
                </v-btn>
              </template>
              <span>{{ t('TT_WARNED') }}</span>
            </v-tooltip>
            <v-tooltip
              v-if="sbomStats.notChartFossLicense.total > 0"
              :open-delay="TOOLTIP_OPEN_DELAY_IN_MS"
              location="bottom"
              content-class="dpTooltip">
              <template v-slot:activator="{props}">
                <v-btn class="text-none card-border font-weight-light" variant="text" size="small" v-bind="props">
                  <v-icon color="red" icon="mdi-shield-off-outline" class="mr-2"></v-icon>
                  {{ `${sbomStats.notChartFossLicense.total} ${t('BTN_NON_CHART')}` }}
                </v-btn>
              </template>
              <span>{{ t('TT_NON_CHART') }}</span>
            </v-tooltip>
            <v-tooltip
              v-if="sbomStats.scanRemarkType.missingCopyrights > 0"
              :open-delay="TOOLTIP_OPEN_DELAY_IN_MS"
              location="bottom"
              content-class="dpTooltip">
              <template v-slot:activator="{props}">
                <v-btn class="text-none card-border font-weight-light" variant="text" size="small" v-bind="props">
                  <v-icon
                    :color="getIconColorScanRemarkLevel(sbomStats.scanRemarkType.missingCopyrightsLevel)"
                    icon="mdi-circle"
                    class="mr-2"></v-icon>
                  {{ `${sbomStats.scanRemarkType.missingCopyrights} ${t('BTN_MISSING_COPYRIGHT')}` }}
                </v-btn>
              </template>
              <span>{{ t('TT_MISSING_COPYRIGHT') }}</span>
            </v-tooltip>
            <v-tooltip
              v-if="sbomStats.scanRemarkType.malformedCopyrights > 0"
              :open-delay="TOOLTIP_OPEN_DELAY_IN_MS"
              location="bottom"
              content-class="dpTooltip">
              <template v-slot:activator="{props}">
                <v-btn class="text-none card-border font-weight-light" variant="text" size="small" v-bind="props">
                  <v-icon color="grey" icon="mdi-circle" class="mr-2"></v-icon>
                  {{ `${sbomStats.scanRemarkType.malformedCopyrights} ${t('BTN_MALFORMED_COPYRIGHTS')}` }}
                </v-btn>
              </template>
              <span>{{ t('TT_MALFORMED_COPYRIGHTS') }}</span>
            </v-tooltip>
          </v-col>
        </v-row>
        <v-row class="mt-0">
          <v-col>
            <v-btn
              size="small"
              :variant="isHTMLSelected ? 'tonal' : 'text'"
              :class="{active: isHTMLSelected}"
              class="ma-2 text-none card-border ml-0"
              @click="loadStyledHTML">
              <v-icon color="primary">mdi-code-brackets</v-icon>
              HTML
            </v-btn>
            <v-btn
              size="small"
              :variant="isTextSelected ? 'tonal' : 'text'"
              :class="{active: isTextSelected}"
              class="ma-2 text-none card-border"
              @click="loadPlainText">
              <v-icon color="primary">mdi-format-text</v-icon>
              Plain Text
            </v-btn>
            <v-btn
              size="small"
              :variant="isJSONFormat ? 'tonal' : 'text'"
              :class="{active: isJSONFormat}"
              class="ma-2 text-none card-border"
              @click="loadJSON">
              <v-icon color="primary">mdi-code-json</v-icon>
              JSON
            </v-btn>

            <DCActionButton
              v-if="!isVehicleProject && currentProject.isProjectOwner"
              icon="mdi-pencil"
              :hint="t('TT_BTN_EDIT_3RD_ADDRESS')"
              :text="t('BTN_EDIT_3RD_ADDRESS')"
              size="small"
              class="mr-2"
              @click.stop="showProjectSettingsDialog" />
          </v-col>
        </v-row>
      </template>
      <template v-else>
        <span class="d-subtitle-2">{{ t('NO_DATA_AVAILABLE') }}</span>
      </template>
    </template>
    <template #table>
      <v-card class="card-border fill-height pa-4 overflow-auto" v-if="showPreview">
        <v-row>
          <v-col md="auto">
            <h3 class="d-subtitle-2 d-secondary-text mb-2 mt-2">{{ t('PREVIEW') }}</h3>
          </v-col>
          <DCopyClipboardButton
            class="mt-3"
            :tableButton="true"
            :hint="t('TT_noticeCopyText')"
            :content="downloadContent" />
          <v-spacer />
          <v-col md="auto" class="d-flex justify-end">
            <DCActionButton
              large
              icon="mdi-download"
              :text="t('BTN_DOWNLOAD')"
              :hint="t('TT_download_notice')"
              class="mr-2"
              @click="downloadNoticeFile" />
          </v-col>
        </v-row>
        <v-row v-if="isJSONFormat">
          <v-col style="height: 400px; width: 400px">
            <JsonViewer3 :value="jsonContent" :expand-depth="1" aria-expanded="true" theme="jv-dark" sort />
          </v-col>
        </v-row>
        <v-row v-else>
          <v-col cols="12" xs="12" sm="12" md="10" lg="8">
            <div
              class="d-text pt-2 html-notice-file"
              v-html="htmlOrPlainContent"
              v-if="htmlOrPlainContent && htmlOrPlainContent.length > 0" />
          </v-col>
        </v-row>
      </v-card>
    </template>
  </TableLayout>
  <ProjectSettings ref="projectSettingsDialog"> </ProjectSettings>
</template>
