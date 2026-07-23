<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {useApprovalFormBase} from '@disclosure-portal/composables/useApprovalFormBase';
import {DocumentMeta, InternalApprovalRequest} from '@disclosure-portal/model/ApprovalRequest';
import ErrorDialogConfig from '@shared/types/ErrorDialogConfig';
import {ApprovableSPDXDto} from '@disclosure-portal/model/Project';
import {UserDto} from '@shared/types/Users';
import projectService from '@disclosure-portal/services/projects';
import {useIdleStore} from '@shared/stores/idle.store';
import {useJobStore} from '@disclosure-portal/stores/jobs';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import eventBus from '@shared/utils/eventbus';
import useRules from '@disclosure-portal/utils/Rules';
import useSnackbar from '@shared/composables/useSnackbar';
import config from '@shared/utils/config';
import {computed, nextTick, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {VForm} from 'vuetify/components';
import {useAppStore} from '@disclosure-portal/stores/app';
import {useApprovableInfoStore} from '@disclosure-portal/stores/approvableInfo.store';

const projectStore = useProjectStore();
const appStore = useAppStore();
const {longText} = useRules();
const {t} = useI18n();
const snackbar = useSnackbar();
const idle = useIdleStore();
const approvableInfoStore = useApprovableInfoStore();

const approverTab = ref<'developer' | 'owner'>('developer');
const form = ref<VForm | null>(null);
const ownerApprover1 = ref('');
const ownerApprover2 = ref('');
const developerApprover1 = ref('');
const developerApprover2 = ref('');
const ownerApproverPre1 = ref<UserDto>();
const ownerApproverPre2 = ref<UserDto>();
const developerApproverIn1 = ref();
const developerApproverPre1 = ref<UserDto>();
const developerApproverIn2 = ref();
const developerApproverPre2 = ref<UserDto>();
const isVehicle = ref(false);

const {
  isVisible,
  selectedChannel,
  sboms,
  selectedSbom,
  approvableInfo,
  comment,
  withZip,
  noFOSS,
  fossVersion,
  mixedFOSS,
  c1,
  c2,
  c3,
  c4,
  c5,
  selectedProjects,
  tab,
  projectModel,
  channels,
  countApprovables,
  stats,
  selectedProjectsContainEmptySbom,
  updateSelectedProjects,
  checkFossMixedStatus,
  loadStats,
  loadSBOMHist,
  autoSelect,
} = useApprovalFormBase({
  setDefaultFlags: () => {
    c1.value = false;
    c2.value = noFOSS.value ? false : countApprovables.value > 0 || selectedSbom.value != null;
    c3.value = noFOSS.value ? false : !(countApprovables.value > 0);
    c4.value = !noFOSS.value;
    c5.value = false;
  },
  resetExtraState: () => {
    approverTab.value = 'developer';
    ownerApprover1.value = '';
    ownerApprover2.value = '';
    developerApprover1.value = '';
    developerApprover2.value = '';
    ownerApproverPre1.value = undefined;
    ownerApproverPre2.value = undefined;
    developerApproverPre1.value = undefined;
    developerApproverPre2.value = undefined;
  },
  fetchFlat: false,
});

const commentRule = longText(t('TAD_COMMENT'));

const open = async (isVehicleProject: boolean) => {
  idle.showIdle = true;
  await approvableInfoStore.fetchApprovableInfo();

  isVehicle.value = isVehicleProject;
  if (config.useFutureFoss) {
    fossVersion.value = 'default';
  } else {
    fossVersion.value = 'legacy';
  }
  if (projectModel.value.customerMeta.userFRI) {
    ownerApproverPre1.value = projectModel.value.customerMeta.userFRI;
  }
  if (projectModel.value.customerMeta.userSRI) {
    ownerApproverPre2.value = projectModel.value.customerMeta.userSRI;
  }
  if (projectModel.value.supplierExtraData.userFRI) {
    developerApproverPre1.value = projectModel.value.supplierExtraData.userFRI;
  }
  if (projectModel.value.supplierExtraData.userSRI) {
    developerApproverPre2.value = projectModel.value.supplierExtraData.userSRI;
  }
  noFOSS.value = projectModel.value.isNoFoss;
  updateSelectedProjects();
  checkFossMixedStatus();

  await autoSelect();
  developerApproverIn1.value?.resetForm();
  developerApproverIn2.value?.resetForm();
  idle.showIdle = false;
  isVisible.value = true;
};

const jobStore = useJobStore();
const doDialogAction = async () => {
  await nextTick(async () => {
    form.value?.validate().then(async (info) => {
      if (!info.valid) {
        return;
      }
      const isDev1Valid = await developerApproverIn1.value?.validateOnCreate();
      const isDev2Valid = await developerApproverIn2.value?.validateOnCreate();
      if (!isDev1Valid || !isDev2Valid) {
        approverTab.value = 'developer';
        return;
      }
      if (
        (ownerApprover1.value !== '' || ownerApprover2.value !== '') &&
        ownerApprover1.value === ownerApprover2.value
      ) {
        approverTab.value = 'owner';
        const d = new ErrorDialogConfig();
        d.title = '' + t('SBOM_REQUEST_INTERNAL_APPROVAL');
        d.description = '' + t('EQUAL_OWNER_APPROVERS_ERROR_MESSAGE');
        eventBus.emit('on-error', {error: d});
        return;
      }
      if (
        (ownerApprover1.value !== '' && ownerApprover2.value === '') ||
        (ownerApprover1.value === '' && ownerApprover2.value !== '')
      ) {
        approverTab.value = 'owner';
        const d = new ErrorDialogConfig();
        d.title = '' + t('SBOM_REQUEST_INTERNAL_APPROVAL');
        d.description = '' + t('BOTH_OR_NONE_OWNER_APPROVERS_ALLOWED_ERROR_MESSAGE');
        eventBus.emit('on-error', {error: d});
        return;
      }
      if (developerApprover1.value === developerApprover2.value) {
        approverTab.value = 'developer';
        const d = new ErrorDialogConfig();
        d.title = '' + t('SBOM_REQUEST_INTERNAL_APPROVAL');
        d.description = '' + t('EQUAL_DEVELOPER_APPROVERS_ERROR_MESSAGE');
        eventBus.emit('on-error', {error: d});
        return;
      }
      idle.showIdle = true;
      idle.idleMessage = t('SBOM_REQUEST_APPROVAL_PROGRESS');
      if (!projectModel.value.isGroup && selectedSbom.value) {
        const approvableSpdx = {
          spdxkey: '',
          versionkey: '',
        } as ApprovableSPDXDto;
        approvableSpdx.spdxkey = selectedSbom.value?._key ?? '';
        approvableSpdx.versionkey = selectedChannel.value?._key ?? '';
        await projectService.updateApprovableSpdx(approvableSpdx, projectModel.value._key);
      }

      const metaDoc = new DocumentMeta();
      metaDoc.c1 = c1.value;
      metaDoc.c2 = c2.value;
      metaDoc.c3 = c3.value;
      metaDoc.c4 = c4.value;
      metaDoc.c5 = c5.value;
      metaDoc.c6 = noFOSS.value;

      const req: InternalApprovalRequest = {
        withZip: withZip.value,
        comment: comment.value,
        guidProject: projectModel.value._key,
        metaDoc: metaDoc,
        customerApprover1: ownerApprover1.value,
        customerApprover2: ownerApprover2.value,
        supplierApprover1: developerApprover1.value,
        supplierApprover2: developerApprover2.value,
        fossVersion: 'default',
      };

      projectService.createInternalApproval(req, projectModel.value._key).then(async (response) => {
        if (response) {
          await jobStore.pollJobStatus(projectModel.value._key, response.jobKey);
          isVisible.value = false;
          snackbar.info(t('DIALOG_request_internal_approval_success'));
          appStore.setShouldReloadApprovals(true);
          if (!projectModel.value.isGroup) {
            await projectStore.fetchProjectByKey(projectModel.value._key);
          }
        } else {
          idle.showIdle = false;
          idle.idleMessage = '';
        }
      });
    });
  });
};

const isEnterpriseOrMobileOrOther = computed(() => {
  return !isVehicle.value && (stats.value.denied > 0 || stats.value.noAssertion > 0);
});

const close = () => {
  isVisible.value = false;
};
const dialogConfig = {
  title: t('SBOM_REQUEST_INTERNAL_APPROVAL'),
  secondaryButton: {text: t('BTN_CANCEL')},
  primaryButton: {text: t('BTN_REQUEST')},
};
defineExpose({open});
</script>

<template>
  <v-form ref="form">
    <v-dialog v-model="isVisible" content-class="large" scrollable width="850">
      <DialogLayout :config="dialogConfig" @close="close" @secondary-action="close" @primary-action="doDialogAction">
        <Stack class="gap-4">
          <v-tabs v-model="approverTab" slider-color="mbti" show-arrows bg-color="tabsHeader">
            <v-tab value="developer">{{ t('TAB_TITLE_DEVELOPER_APPROVER') }}</v-tab>
            <v-tab value="owner">{{ t('TAB_TITLE_OWNER_APPROVER') }}</v-tab>
          </v-tabs>
          <v-tabs-window v-model="approverTab" eager>
            <v-tabs-window-item value="developer" eager>
              <Stack class="gap-4">
                <Stack direction="row">
                  <v-icon size="small" color="warning">mdi-alert</v-icon>
                  <span class="text-body-2">{{ t('REPORTER_REMARK') }}</span>
                </Stack>
                <DAutocompleteUser
                  ref="developerApproverIn1"
                  v-model="developerApprover1"
                  :preselect="developerApproverPre1"
                  :project-key="projectModel._key"
                  :label="t('FIRST_REPORTER_LABEL')"
                  data-testid="developerApprover1"
                  only-internal-users
                  required />
                <DAutocompleteUser
                  ref="developerApproverIn2"
                  v-model="developerApprover2"
                  :preselect="developerApproverPre2"
                  :project-key="projectModel._key"
                  :label="t('SECOND_REPORTER_LABEL')"
                  data-testid="developerApprover2"
                  only-internal-users
                  required />
              </Stack>
            </v-tabs-window-item>
            <v-tabs-window-item value="owner">
              <Stack class="gap-4">
                <Stack direction="row">
                  <v-icon size="small" color="warning">mdi-alert</v-icon>
                  <span class="text-body-2">{{ t('REPORTER_REMARK') }}</span>
                </Stack>
                <DAutocompleteUser
                  ref="ownerApproverIn1"
                  v-model="ownerApprover1"
                  :preselect="ownerApproverPre1"
                  :project-key="projectModel._key"
                  :label="t('FIRST_REPORTER_LABEL')"
                  data-testid="ownerApprover1"
                  only-internal-users />
                <DAutocompleteUser
                  ref="ownerApproverIn2"
                  v-model="ownerApprover2"
                  :preselect="ownerApproverPre2"
                  :project-key="projectModel._key"
                  :label="t('SECOND_REPORTER_LABEL')"
                  data-testid="ownerApprover2"
                  only-internal-users />
              </Stack>
            </v-tabs-window-item>
          </v-tabs-window>

          <SbomChannelSelector
            v-if="!projectModel.isGroup"
            :channels="channels"
            :sboms="sboms"
            :selected-channel="selectedChannel"
            :selected-sbom="selectedSbom"
            :no-f-o-s-s="noFOSS"
            :is-vehicle="isVehicle"
            :approvable-spdx-key="projectModel.approvablespdx.spdxkey"
            @update:selected-channel="
              selectedChannel = $event;
              loadSBOMHist();
            "
            @update:selected-sbom="
              selectedSbom = $event;
              loadStats();
            " />

          <ApprovalWarnings
            :is-denied-or-unasserted="false"
            :is-enterprise-or-mobile-or-other="isEnterpriseOrMobileOrOther && !noFOSS"
            :mixed-f-o-s-s="mixedFOSS"
            :no-f-o-s-s="noFOSS"
            :foss-version="fossVersion"
            :selected-projects-contain-empty-sbom="selectedProjectsContainEmptySbom" />

          <FossVersionSelector v-if="config.useFutureFoss" v-model="fossVersion" :disabled="true" />

          <ApprovalContentTabs
            v-model:tab="tab"
            :stats="stats"
            :show-red-warn-denied-decisions-message="approvableInfo.hasDeniedDecisions"
            :projects="approvableInfo.projects"
            :channels="projectModel.versions"
            :is-group="projectModel.isGroup"
            :no-f-o-s-s="noFOSS"
            :foss-version="fossVersion"
            :selected-projects="selectedProjects"
            do-filter
            @update:selectedProjects="selectedProjects = $event" />

          <v-textarea
            v-model="comment"
            :rules="commentRule"
            :label="t('TAD_COMMENT')"
            variant="outlined"
            counter="1000"
            hide-details
            no-resize />

          <v-switch v-model="withZip" color="primary" :label="t('WITH_ZIP_MARKER')" hide-details></v-switch>

          <LegacyApprovalSection
            v-if="fossVersion === 'legacy'"
            :no-f-o-s-s="noFOSS"
            :c1="c1"
            :c2="c2"
            :c3="c3"
            :c4="c4"
            :c5="c5"
            @update:c1="c1 = $event"
            @update:c2="c2 = $event"
            @update:c3="c3 = $event"
            @update:c4="c4 = $event"
            @update:c5="c5 = $event" />
        </Stack>
      </DialogLayout>
    </v-dialog>
  </v-form>
</template>
