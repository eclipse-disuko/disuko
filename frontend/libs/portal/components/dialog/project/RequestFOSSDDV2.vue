<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {useApprovalFormBase} from '@disclosure-portal/composables/useApprovalFormBase';
import {DocumentMeta, ExternalApprovalRequest} from '@disclosure-portal/model/ApprovalRequest';
import {ApprovableSPDXDto} from '@disclosure-portal/model/Project';
import {OverallReviewState, SpdxFile, VersionSlim} from '@disclosure-portal/model/VersionDetails';
import projectService from '@disclosure-portal/services/projects';
import versionService from '@disclosure-portal/services/version';
import {useIdleStore} from '@shared/stores/idle.store';
import {useSbomStore} from '@disclosure-portal/stores/sbom.store';
import {useJobStore} from '@disclosure-portal/stores/jobs';
import useRules from '@disclosure-portal/utils/Rules';
import config from '@shared/utils/config';
import dayjs from 'dayjs';
import {computed, nextTick, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {VForm} from 'vuetify/components';

const sbomStore = useSbomStore();
const {longText} = useRules();
const {t} = useI18n();
const idle = useIdleStore();

const form = ref<VForm | null>(null);
const dd = ref();
const vehicle = ref(false);
const radioGroup = ref(0);
const childProjectChannels = ref<Map<string, VersionSlim>>(new Map());

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
    radioGroup.value = noFOSS.value ? 3 : 2;
    c1.value = noFOSS.value ? false : vehicle.value;
    c2.value = noFOSS.value ? false : vehicle.value ? false : countApprovables.value > 0 || selectedSbom.value != null;
    c3.value = noFOSS.value ? false : vehicle.value ? false : !(countApprovables.value > 0);
    c4.value = noFOSS.value ? false : !vehicle.value;
    c5.value = false;
  },
  resetExtraState: () => {
    radioGroup.value = 0;
    childProjectChannels.value.clear();
  },
  fetchFlat: true,
});

const markedApprovableSpdx = computed(() => projectModel.value.approvablespdx);

const isNoSbomNoFossWarning = computed(() => !projectModel.value.isGroup && !selectedSbom.value && !noFOSS.value);

const getChannelSboms = (versionKey: string): SpdxFile[] => {
  const versionEntry = sbomStore.getAllSBOMs.find((entry) => entry.versionKey === versionKey);
  return (versionEntry?.spdxFileHistory ?? []).map((sbom, index) => ({...sbom, isRecent: index === 0}));
};

const getSelectableSbom = (versionKey: string, sbomKey: string) => {
  const visibleSbom = sboms.value.find((sbom) => sbom._key === sbomKey);
  if (visibleSbom) {
    return visibleSbom;
  }

  const channelSboms = getChannelSboms(versionKey);
  const exactSbom = channelSboms.find((sbom) => sbom._key === sbomKey) ?? null;
  if (exactSbom) {
    sboms.value = channelSboms;
  }

  return exactSbom;
};

const selectChannelAndSbom = async (versionKey: string, sbomKey: string) => {
  selectedChannel.value = channels.value.find((channel) => channel._key === versionKey) ?? null;
  if (!selectedChannel.value) {
    return false;
  }

  await loadSBOMHist();
  selectedSbom.value = getSelectableSbom(versionKey, sbomKey);
  if (!selectedSbom.value) {
    return false;
  }

  await loadStats();
  return true;
};

const isRdConfirmationMissing = computed(() => {
  if (!vehicle.value || noFOSS.value) {
    return false;
  }

  if (!projectModel.value.isGroup && !selectedSbom.value) {
    return false;
  }

  if (!projectModel.value.isGroup) {
    const approvableSpdx = markedApprovableSpdx.value;
    const sbomKey = selectedSbom.value?._key || approvableSpdx?.spdxkey;
    const channelKey = selectedChannel.value?._key || approvableSpdx?.versionkey;

    if (!sbomKey || !channelKey) {
      return false;
    }

    const channel = channels.value.find((c) => c._key === channelKey);
    if (!channel) {
      return false;
    }

    const hasAuditedReview = channel.overallReviews?.some(
      (review) => review.sbomId === sbomKey && review.state === OverallReviewState.AUDITED,
    );
    return !hasAuditedReview;
  }

  if (projectModel.value.isGroup && approvableInfo.value.projects) {
    const selectedProjectsSet = new Set(selectedProjects.value);

    for (const project of approvableInfo.value.projects) {
      if (selectedProjectsSet.size > 0 && !selectedProjectsSet.has(project.projectKey)) {
        continue;
      }

      if (!project.approvablespdx?.spdxkey || !project.approvablespdx?.versionkey) {
        continue;
      }

      const channel = childProjectChannels.value.get(project.approvablespdx.versionkey);

      if (!channel) {
        return true;
      }

      const hasAuditedReview = channel.overallReviews?.some(
        (review) => review.sbomId === project.approvablespdx.spdxkey && review.state === OverallReviewState.AUDITED,
      );

      if (!hasAuditedReview) {
        return true;
      }
    }
  }

  return false;
});

watch(radioGroup, () => {
  if (radioGroup.value == 3) {
    noFOSS.value = true;
  }
});

const commentRule = longText(t('TAD_COMMENT'));

const smartAutoSelect = async () => {
  if (sbomStore.selectedSBOMKey) {
    await autoSelect();
    if (selectedChannel.value && selectedSbom.value) {
      selectedSbom.value = getSelectableSbom(selectedChannel.value._key, selectedSbom.value._key) ?? selectedSbom.value;
    }
    return;
  }

  if (noFOSS.value) return;

  await sbomStore.fetchAllSBOMsFlat();

  const markedVersionKey = markedApprovableSpdx.value.versionkey;
  const markedSbomKey = markedApprovableSpdx.value.spdxkey;
  if (markedVersionKey && markedSbomKey && (await selectChannelAndSbom(markedVersionKey, markedSbomKey))) {
    return;
  }

  if (vehicle.value) {
    const candidates: {channel: VersionSlim; sbomId: string; reviewUpdated: string}[] = [];
    for (const channel of channels.value) {
      for (const review of channel.overallReviews ?? []) {
        if (review.state === OverallReviewState.AUDITED) {
          candidates.push({channel, sbomId: review.sbomId, reviewUpdated: review.updated});
        }
      }
    }
    candidates.sort((a, b) => (dayjs(a.reviewUpdated).isBefore(b.reviewUpdated) ? 1 : -1));
    if (candidates[0]) {
      await selectChannelAndSbom(candidates[0].channel._key, candidates[0].sbomId);
    }
    return;
  }

  const sortedFlat = [...sbomStore.getAllSBOMsFlat].sort((a, b) => (dayjs(a.updated).isBefore(b.updated) ? 1 : -1));
  const best = sortedFlat[0];
  if (best) {
    await selectChannelAndSbom(best.versionKey, best._key);
  }
};

const open = async (isVehicle: boolean) => {
  idle.showIdle = true;
  approvableInfo.value = await projectService.getApprovableInfo(projectModel.value._key);

  checkFossMixedStatus();
  vehicle.value = isVehicle;
  if (vehicle.value) {
    withZip.value = true;
  }
  if (config.useFutureIt && !isVehicle) {
    fossVersion.value = 'default';
  } else if (config.useFutureProduct && isVehicle) {
    fossVersion.value = 'default';
  } else {
    fossVersion.value = 'legacy';
  }
  noFOSS.value = projectModel.value.isNoFoss;

  updateSelectedProjects();

  await smartAutoSelect();

  if (projectModel.value.isGroup && approvableInfo.value.projects) {
    childProjectChannels.value.clear();

    const versionFetchPromises = approvableInfo.value.projects
      .filter((p) => p.approvablespdx.versionkey)
      .map(async (project) => {
        try {
          const versionDetails = await versionService.getVersion(project.projectKey, project.approvablespdx.versionkey);
          childProjectChannels.value.set(project.approvablespdx.versionkey, versionDetails.data);
        } catch (error) {
          console.error(`Failed to fetch version details for project ${project.projectKey}:`, error);
        }
      });

    await Promise.all(versionFetchPromises);
  }

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

  if (isRdConfirmationMissing.value && config.enforceFOSSOfficeConfirmation) {
    return;
  }

  const metaDoc: DocumentMeta = new DocumentMeta();
  if (vehicle.value) {
    metaDoc.c1 = radioGroup.value == 1;
    metaDoc.c2 = radioGroup.value == 2;
    metaDoc.c3 = false;
    metaDoc.c4 = false;
    metaDoc.c5 = false;
  } else {
    metaDoc.c1 = c1.value;
    metaDoc.c2 = c2.value;
    metaDoc.c3 = c3.value;
    metaDoc.c4 = c4.value;
    metaDoc.c5 = c5.value;
  }
  metaDoc.c6 = noFOSS.value || !selectedSbom.value;

  let determinedFossVersion: 'default' | 'legacy' | 'vehicle-legacy';

  if (config.useFutureIt && !vehicle.value) {
    determinedFossVersion = fossVersion.value === 'default' ? 'default' : 'legacy';
  } else if (config.useFutureProduct && vehicle.value) {
    determinedFossVersion = fossVersion.value === 'default' ? 'default' : 'vehicle-legacy';
  } else {
    determinedFossVersion = vehicle.value ? 'vehicle-legacy' : 'legacy';
  }

  const req: ExternalApprovalRequest = {
    comment: comment.value,
    guidProject: projectModel.value._key,
    metaDoc: metaDoc,
    withZip: withZip.value,
    fossVersion: determinedFossVersion,
    selectedProjects: selectedProjects.value,
  };

  idle.showIdle = true;

  if (!projectModel.value.isGroup && selectedSbom.value) {
    const approvableSpdx = {
      spdxkey: '',
      versionkey: '',
    } as ApprovableSPDXDto;
    approvableSpdx.spdxkey = selectedSbom.value._key;
    approvableSpdx.versionkey = selectedChannel.value?._key ?? '';
    await projectService.updateApprovableSpdx(approvableSpdx, projectModel.value._key);
  }

  const response = await (vehicle.value
    ? projectService.createVehicleApproval(req, projectModel.value._key)
    : projectService.createExternalApproval(req, projectModel.value._key));

  if (response) {
    await jobStore.pollJobStatus(projectModel.value._key, response.jobKey);
    isVisible.value = false;
    dd.value?.open(response.approvalGuid);
  } else {
    idle.showIdle = false;
  }
};

const isDeniedOrUnasserted = computed(() => {
  return vehicle.value && (stats.value.denied > 0 || stats.value.noAssertion > 0);
});

const isWarned = computed(() => {
  return vehicle.value && stats.value.warned > 0;
});

const isEitherFutureFoss = computed(() => {
  return (config.useFutureIt && !vehicle.value) || (config.useFutureProduct && vehicle.value);
});

const canGenerateFoss = computed(() => {
  const rdConfirmationCondition = vehicle.value ? !isRdConfirmationMissing.value : true;
  const futureFossCondition = isEitherFutureFoss.value && fossVersion.value === 'default' ? !isWarned.value : true;
  const noSbomLegacyCondition = !(isNoSbomNoFossWarning.value && fossVersion.value === 'legacy');
  return (
    !isDeniedOrUnasserted.value &&
    rdConfirmationCondition &&
    futureFossCondition &&
    noSbomLegacyCondition &&
    selectedProjects.value?.length > 0
  );
});

const isEnterpriseOrMobileOrOther = computed(() => {
  return !vehicle.value && (stats.value.denied > 0 || stats.value.noAssertion > 0);
});

const showRedWarnDeniedDecisionsMessage = computed(
  () => !isDeniedOrUnasserted.value && approvableInfo.value.hasDeniedDecisions,
);

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
              :selected-channel="selectedChannel"
              :selected-sbom="selectedSbom"
              :no-f-o-s-s="noFOSS"
              :is-vehicle="vehicle"
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
              :is-denied-or-unasserted="isDeniedOrUnasserted"
              :is-either-future-foss="isEitherFutureFoss"
              :is-rd-confirmation-missing="isRdConfirmationMissing"
              :is-warned="isWarned"
              :is-enterprise-or-mobile-or-other="isEnterpriseOrMobileOrOther"
              :mixed-f-o-s-s="mixedFOSS"
              :no-f-o-s-s="noFOSS"
              :foss-version="fossVersion"
              :selected-projects-contain-empty-sbom="selectedProjectsContainEmptySbom"
              :is-no-sbom-no-foss-warning="isNoSbomNoFossWarning" />

            <FossVersionSelector v-model="fossVersion" :disabled="!(config.useFutureProduct && vehicle)" />

            <ApprovalContentTabs
              v-model:tab="tab"
              :stats="stats"
              :show-red-warn-denied-decisions-message="showRedWarnDeniedDecisionsMessage"
              :projects="approvableInfo.projects"
              :channels="childProjectChannels"
              :is-group="projectModel.isGroup"
              :no-f-o-s-s="noFOSS"
              :foss-version="fossVersion"
              :selected-projects="selectedProjects"
              @update:selectedProjects="selectedProjects = $event" />

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

            <LegacyApprovalSection
              v-if="fossVersion === 'legacy'"
              :no-f-o-s-s="noFOSS"
              :is-vehicle="vehicle"
              v-model:c1="c1"
              v-model:c2="c2"
              v-model:c3="c3"
              v-model:c4="c4"
              v-model:c5="c5"
              v-model:radio-group="radioGroup"
              @update:noFOSS="noFOSS = $event" />
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
            v-if="canGenerateFoss"
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
