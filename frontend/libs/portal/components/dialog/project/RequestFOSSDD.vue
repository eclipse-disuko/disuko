<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {useApprovalFormBase} from '@disclosure-portal/composables/useApprovalFormBase';
import {DocumentMeta, ExternalApprovalRequest} from '@disclosure-portal/model/ApprovalRequest';
import {ApprovableSPDXDto} from '@disclosure-portal/model/Project';
import ErrorDialogConfig from '@shared/types/ErrorDialogConfig';
import projectService from '@disclosure-portal/services/projects';
import {useIdleStore} from '@shared/stores/idle.store';
import {useJobStore} from '@disclosure-portal/stores/jobs';
import eventBus from '@shared/utils/eventbus';
import useRules from '@disclosure-portal/utils/Rules';
import config from '@shared/utils/config';
import {computed, nextTick, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {VForm} from 'vuetify/components';
import SbomChannelSelector from './shared/SbomChannelSelector.vue';
import FossVersionSelector from './shared/FossVersionSelector.vue';
import ApprovalContentTabs from './shared/ApprovalContentTabs.vue';
import LegacyApprovalSection from './shared/LegacyApprovalSection.vue';

const {longText} = useRules();
const {t} = useI18n();
const idle = useIdleStore();

const form = ref<VForm | null>(null);
const dd = ref();

const {
  isVisible,
  selectedChannel,
  sboms,
  selectedSbom,
  approvableInfo,
  documentFlags,
  comment,
  withZip,
  noFOSS,
  fossVersion,
  mixedFOSS,
  selectedProjects,
  tab,
  projectModel,
  channels,
  countApprovables,
  stats,
  selectedProjectsContainEmptySbom,
  updateSelectedProjects,
  checkFossMixedStatus,
  loadSBOMHist,
  loadStats,
  autoSelect,
  fetchApprovableInfo,
} = useApprovalFormBase({
  setDefaultFlags: () => {
    documentFlags.value.c1 = false;
    documentFlags.value.c2 = noFOSS.value ? false : countApprovables.value > 0 || selectedSbom.value != null;
    documentFlags.value.c3 = noFOSS.value ? false : !(countApprovables.value > 0);
    documentFlags.value.c4 = !noFOSS.value;
    documentFlags.value.c5 = false;
  },
});

const showRedWarnDeniedDecisionsMessage = computed(() => approvableInfo.value.hasDeniedDecisions);

const commentRule = longText(t('TAD_COMMENT'));

const open = async () => {
  idle.showIdle = true;
  await fetchApprovableInfo();

  checkFossMixedStatus();
  fossVersion.value = config.useFutureFoss ? 'default' : 'legacy';
  noFOSS.value = projectModel.value.isNoFoss;
  updateSelectedProjects();
  await autoSelect();

  idle.showIdle = false;
  isVisible.value = true;
};

const jobStore = useJobStore();
const doDialogAction = async () => {
  await nextTick();
  const info = await form.value?.validate();
  if (!info?.valid) {
    return;
  }

  if (!projectModel.value.isGroup && selectedChannel.value !== null && selectedSbom.value === null) {
    const d = new ErrorDialogConfig();
    d.title = '' + t('TITLE_GENERATE_FOSS_DD');
    d.description = '' + t('BOTH_OR_NONE_CHANNEL_AND_SBOM_ALLOWED_ERROR_MESSAGE');
    eventBus.emit('on-error', {error: d});
    return;
  }

  const metaDoc = Object.assign(new DocumentMeta(), documentFlags.value, {c6: noFOSS.value});

  const req: ExternalApprovalRequest = {
    comment: comment.value,
    guidProject: projectModel.value._key,
    metaDoc: metaDoc,
    withZip: withZip.value,
    fossVersion: fossVersion.value,
    selectedProjects: selectedProjects.value,
  };

  idle.showIdle = true;

  if (!projectModel.value.isGroup) {
    const approvableSpdx = {
      spdxkey: '',
      versionkey: '',
    } as ApprovableSPDXDto;
    approvableSpdx.spdxkey = selectedSbom.value?._key ?? '';
    approvableSpdx.versionkey = selectedChannel.value?._key ?? '';
    await projectService.updateApprovableSpdx(approvableSpdx, projectModel.value._key);
  }

  const response = await projectService.createExternalApproval(req, projectModel.value._key);

  if (response) {
    await jobStore.pollJobStatus(projectModel.value._key, response.jobKey);
    isVisible.value = false;
    dd.value?.open(response.approvalGuid);
  } else {
    idle.showIdle = false;
  }
};

defineExpose({open});
</script>

<template>
  <v-form ref="form">
    <v-dialog v-model="isVisible" content-class="large" scrollable width="850">
      <v-card class="pa-8">
        <v-card-title>
          <Stack direction="row" align="center">
            <span class="text-h5">
              {{ t('TITLE_GENERATE_FOSS_DD') }}
            </span>
            <span class="flex-grow"></span>
            <span>
              <DCloseButton @click="isVisible = false" />
            </span>
          </Stack>
        </v-card-title>

        <v-card-text>
          <Stack class="gap-4">
            <SbomChannelSelector
              v-if="!projectModel.isGroup"
              v-model:selected-channel="selectedChannel"
              v-model:selected-sbom="selectedSbom"
              :channels="channels"
              :sboms="sboms"
              :no-f-o-s-s="noFOSS"
              :approvable-spdx-key="projectModel.approvablespdx.spdxkey"
              @update:selected-channel="loadSBOMHist()"
              @update:selected-sbom="loadStats()" />

            <section v-if="mixedFOSS">
              <v-alert color="warning" type="warning">
                <span>{{ t('MIXED_FOSS_MESSAGE') }}</span>
              </v-alert>
            </section>

            <section v-if="noFOSS && fossVersion === 'legacy'">
              <v-alert color="warning" type="warning">
                <span>{{ t('NO_FOSS_MESSAGE') }}</span>
              </v-alert>
            </section>

            <section v-if="selectedProjectsContainEmptySbom">
              <v-alert color="warning" type="warning">
                <span>{{ t('NO_PROJECT_NO_FOSS') }}</span>
              </v-alert>
            </section>

            <FossVersionSelector v-if="config.useFutureFoss" v-model="fossVersion" />

            <ApprovalContentTabs
              v-model:tab="tab"
              v-model:selected-projects="selectedProjects"
              :stats="stats"
              :show-red-warn-denied-decisions-message="showRedWarnDeniedDecisionsMessage"
              :projects="approvableInfo.projects ?? []"
              :channels="projectModel.versions"
              :is-group="projectModel.isGroup"
              :no-f-o-s-s="noFOSS"
              :foss-version="fossVersion" />

            <v-textarea
              v-model="comment"
              :rules="commentRule"
              :label="t('TAD_COMMENT')"
              variant="outlined"
              counter="1000"
              hide-details
              no-resize />

            <v-switch v-model="withZip" color="primary" :label="t('WITH_ZIP_MARKER')" hide-details></v-switch>

            <LegacyApprovalSection v-if="fossVersion === 'legacy'" v-model="documentFlags" />
          </Stack>
        </v-card-text>

        <v-card-actions>
          <v-spacer></v-spacer>
          <DCActionButton
            isDialogButton
            size="small"
            variant="text"
            @click="isVisible = false"
            class="mr-4"
            :text="t('BTN_CANCEL')" />

          <DCActionButton
            isDialogButton
            size="small"
            variant="flat"
            @click="doDialogAction"
            :text="t('BTN_GENERATE_FOSS_DD')" />
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-form>
  <DocumentDownloadDialog ref="dd" />
</template>
<style scoped lang="scss">
a {
  color: var(--text-color);
  display: block;
  &:hover {
    text-decoration: underline;
  }
}
</style>
