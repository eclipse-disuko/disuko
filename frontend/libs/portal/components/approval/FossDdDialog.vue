<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import ApprovalContentTabs from '@disclosure-portal/components/approval/ApprovalContentTabs.vue';
import FossVersionSelector from '@disclosure-portal/components/approval/FossVersionSelector.vue';
import SbomChannelSelector from '@disclosure-portal/components/approval/SbomChannelSelector.vue';
import DocumentDownloadDialog from '@disclosure-portal/components/dialog/DocumentDownloadDialog.vue';
import {useFossDdForm} from '@disclosure-portal/composables/approval/useFossDdForm';
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
  comment,
  commentRule,
  fossVersion,
  selectedProjects,
  form,
  vehicle,
  approvableInfo,
  childProjectChannels,
  dd,
  isDeniedOrUnasserted,
  isEnterpriseOrMobileOrOther,
  isFossOfficeConfirmationMissing,
  showRedWarnDeniedDecisionsMessage,
  loadSBOMHist,
  loadStats,
  doDialogAction,
  open,
} = useFossDdForm();

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
              :channels="channels"
              :sboms="sboms"
              v-model:selected-channel="selectedChannel"
              v-model:selected-sbom="selectedSbom"
              :no-f-o-s-s="noFOSS"
              :is-vehicle="vehicle"
              :approvable-spdx-key="projectModel.approvablespdx.spdxkey"
              @update:selected-channel="loadSBOMHist"
              @update:selected-sbom="loadStats" />

            <section
              id="warning"
              v-if="isDeniedOrUnasserted || isEnterpriseOrMobileOrOther || noFOSS || isFossOfficeConfirmationMissing">
              <v-alert color="warning" type="warning">
                <span v-if="isDeniedOrUnasserted">
                  {{ t('DENIED_OR_UNASSARETED_MESSAGE') }}
                </span>
                <span v-else-if="isFossOfficeConfirmationMissing">
                  {{ t('CONFIRMATION_MISSING') }}
                </span>
                <span v-else-if="isEnterpriseOrMobileOrOther">
                  {{ t('ENTERPRISE_MOBILE_OTHER_MESSAGE') }}
                  <a :href="t('ENTERPRISE_MOBILE_OTHER_MESSAGE_CTA')" target="_blank">
                    <v-icon>mdi mdi-chevron-right</v-icon>
                    <span>{{ t('LINK_CLICK_HERE') }} </span>
                  </a>
                </span>
                <span v-else-if="noFOSS">
                  {{ t('NO_FOSS_MESSAGE') }}
                </span>
              </v-alert>
            </section>

            <FossVersionSelector v-if="config.useFutureFoss" v-model="fossVersion" />

            <ApprovalContentTabs
              v-model:tab="tab"
              :stats="stats!"
              :show-red-warn-denied-decisions-message="showRedWarnDeniedDecisionsMessage"
              :projects="approvableInfo.projects"
              :channels="childProjectChannels"
              :is-group="projectModel.isGroup"
              :no-f-o-s-s="noFOSS"
              :foss-version="fossVersion"
              v-model:selected-projects="selectedProjects"
              selectable />

            <v-textarea
              v-model="comment"
              :rules="commentRule"
              :label="t('TAD_COMMENT')"
              variant="outlined"
              counter="1000"
              hide-details
              no-resize />

            <v-switch
              v-model="withZip"
              color="primary"
              :readonly="vehicle"
              :label="t('WITH_ZIP_MARKER')"
              hide-details></v-switch>
            <div>
              <Stack direction="row" align="center">
                <v-icon v-if="noFOSS" size="small">mdi-alert</v-icon>
                <span class="d-block" v-if="noFOSS">{{ t('NO_FOSS_WARNING') }}</span>
              </Stack>
              <v-switch v-model="noFOSS" color="primary" :label="t('NO_FOSS_MARKER')" hide-details></v-switch>
            </div>
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
            v-if="!isDeniedOrUnasserted"
            size="small"
            variant="flat"
            @click="doDialogAction"
            :text="t('BTN_GENERATE_FOSS_DD')" />

          <DCActionButton
            isDialogButton
            v-else
            size="small"
            variant="flat"
            color="primary"
            @click="isVisible = false"
            :text="t('BTN_CLOSE')" />
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
