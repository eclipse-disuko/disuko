<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import ApprovalContentTabs from '@disclosure-portal/components/approval/ApprovalContentTabs.vue';
import FossVersionSelector from '@disclosure-portal/components/approval/FossVersionSelector.vue';
import SbomChannelSelector from '@disclosure-portal/components/approval/SbomChannelSelector.vue';
import {useInternalApprovalForm} from '@disclosure-portal/composables/approval/useInternalApprovalForm';
import config from '@shared/utils/config';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();

const {
  isVisible,
  projectModel,
  channels,
  sboms,
  stats,
  selectedChannel,
  selectedSbom,
  noFOSS,
  withZip,
  tab,
  approverTab,
  comment,
  commentRule,
  fossVersion,
  selectedProjects,
  form,
  isVehicle,
  approvableInfo,
  ownerApprover1,
  ownerApprover2,
  developerApprover1,
  developerApprover2,
  ownerApproverIn1,
  ownerApproverPre1,
  ownerApproverIn2,
  ownerApproverPre2,
  developerApproverIn1,
  developerApproverPre1,
  developerApproverIn2,
  developerApproverPre2,
  loadSBOMHist,
  loadStats,
  doDialogAction,
  close,
  open,
  dialogConfig,
} = useInternalApprovalForm();

defineExpose({open});
</script>

<template>
  <v-form ref="form">
    <v-dialog v-model="isVisible" content-class="large" scrollable width="850">
      <DialogLayout :config="dialogConfig" @close="close" @secondary-action="close" @primary-action="doDialogAction">
        <Stack class="gap-4">
          <v-tabs v-model="approverTab" slider-color="brand" show-arrows bg-color="tabsHeader">
            <v-tab value="owner">{{ t('TAB_TITLE_OWNER_APPROVER') }}</v-tab>
            <v-tab value="developer">{{ t('TAB_TITLE_DEVELOPER_APPROVER') }}</v-tab>
          </v-tabs>
          <v-tabs-window v-model="approverTab" eager>
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
          </v-tabs-window>

          <SbomChannelSelector
            v-if="!projectModel.isGroup"
            :channels="channels"
            :sboms="sboms"
            v-model:selected-channel="selectedChannel"
            v-model:selected-sbom="selectedSbom"
            :no-f-o-s-s="noFOSS"
            :is-vehicle="isVehicle"
            :approvable-spdx-key="projectModel.approvablespdx.spdxkey"
            @update:selected-channel="loadSBOMHist"
            @update:selected-sbom="loadStats" />

          <FossVersionSelector v-if="config.useFutureFoss" v-model="fossVersion" />

          <ApprovalContentTabs
            v-model:tab="tab"
            :stats="stats!"
            :show-red-warn-denied-decisions-message="approvableInfo.hasDeniedDecisions"
            :projects="approvableInfo.projects"
            :channels="projectModel.versions"
            :is-group="projectModel.isGroup"
            :no-f-o-s-s="noFOSS"
            :foss-version="fossVersion"
            v-model:selected-projects="selectedProjects" />

          <v-textarea
            v-model="comment"
            :rules="commentRule"
            :label="t('TAD_COMMENT')"
            variant="outlined"
            counter="1000"
            hide-details
            no-resize />

          <v-switch v-model="withZip" color="primary" :label="t('WITH_ZIP_MARKER')" hide-details></v-switch>
          <div>
            <Stack direction="row" align="center">
              <v-icon v-if="noFOSS" size="small">mdi-alert</v-icon>
              <span class="d-block" v-if="noFOSS">{{ t('NO_FOSS_WARNING') }}</span>
            </Stack>
            <v-switch v-model="noFOSS" color="primary" :label="t('NO_FOSS_MARKER')" hide-details></v-switch>
          </div>
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
