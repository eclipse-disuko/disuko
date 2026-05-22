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
import {useAppStore} from '@disclosure-portal/stores/app';
import eventBus from '@shared/utils/eventbus';
import useRules from '@disclosure-portal/utils/Rules';
import useSnackbar from '@shared/composables/useSnackbar';
import config from '@shared/utils/config';
import {nextTick, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {VForm} from 'vuetify/components';
import SbomChannelSelector from './shared/SbomChannelSelector.vue';
import FossVersionSelector from './shared/FossVersionSelector.vue';
import ApprovalContentTabs from './shared/ApprovalContentTabs.vue';
import LegacyApprovalSection from './shared/LegacyApprovalSection.vue';

const projectStore = useProjectStore();
const appStore = useAppStore();
const {longText} = useRules();
const {t} = useI18n();
const snackbar = useSnackbar();
const idle = useIdleStore();

const approverTab = ref<'developer' | 'owner'>('developer');
const form = ref<VForm | null>(null);
const ownerApprover1 = ref('');
const ownerApprover2 = ref('');
const developerApprover1 = ref('');
const developerApprover2 = ref('');
const ownerApproverIn1 = ref();
const ownerApproverPre1 = ref<UserDto>();
const ownerApproverIn2 = ref();
const ownerApproverPre2 = ref<UserDto>();
const developerApproverIn1 = ref();
const developerApproverPre1 = ref<UserDto>();
const developerApproverIn2 = ref();
const developerApproverPre2 = ref<UserDto>();

const resetExtraState = () => {
  approverTab.value = 'developer';
  ownerApprover1.value = '';
  ownerApprover2.value = '';
  developerApprover1.value = '';
  developerApprover2.value = '';
  ownerApproverPre1.value = undefined;
  ownerApproverPre2.value = undefined;
  developerApproverPre1.value = undefined;
  developerApproverPre2.value = undefined;
};

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
  resetExtraState,
});

const commentRule = longText(t('TAD_COMMENT'));

const open = async () => {
  idle.showIdle = true;
  await fetchApprovableInfo();

  fossVersion.value = config.useFutureFoss ? 'default' : 'legacy';
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
      const isOwner1Valid = await ownerApproverIn1.value?.validateOnCreate();
      const isOwner2Valid = await ownerApproverIn2.value?.validateOnCreate();
      if (!isDev1Valid || !isDev2Valid) {
        approverTab.value = 'developer';
        return;
      }
      if (!isOwner1Valid || !isOwner2Valid) {
        approverTab.value = 'owner';
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

      const metaDoc = Object.assign(new DocumentMeta(), documentFlags.value, {c6: noFOSS.value});

      const req: InternalApprovalRequest = {
        withZip: withZip.value,
        comment: comment.value,
        guidProject: projectModel.value._key,
        metaDoc: metaDoc,
        customerApprover1: ownerApprover1.value,
        customerApprover2: ownerApprover2.value,
        supplierApprover1: developerApprover1.value,
        supplierApprover2: developerApprover2.value,
        fossVersion: fossVersion.value,
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
          <v-tabs v-model="approverTab" slider-color="brand" show-arrows bg-color="tabsHeader">
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
            :show-red-warn-denied-decisions-message="approvableInfo.hasDeniedDecisions"
            :projects="approvableInfo.projects ?? []"
            :channels="projectModel.versions"
            :is-group="projectModel.isGroup"
            :no-f-o-s-s="noFOSS"
            :foss-version="fossVersion"
            :do-filter="true" />

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
      </DialogLayout>
    </v-dialog>
  </v-form>
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
