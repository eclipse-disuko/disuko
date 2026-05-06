<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {useApprovalCheck} from '@disclosure-portal/composables/useApprovalCheck';
import {DocumentMeta, ExternalApprovalRequest} from '@disclosure-portal/model/ApprovalRequest';
import {ApprovableSPDXDto} from '@disclosure-portal/model/Project';
import {ComponentStats, OverallReviewState, SpdxFile, VersionSlim} from '@disclosure-portal/model/VersionDetails';
import projectService from '@disclosure-portal/services/projects';
import versionService from '@disclosure-portal/services/version';
import {useIdleStore} from '@shared/stores/idle.store';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useJobStore} from '@disclosure-portal/stores/jobs';
import {useSbomStore} from '@disclosure-portal/stores/sbom.store';
import useRules from '@disclosure-portal/utils/Rules';
import {formatDateAndTime} from '@disclosure-portal/utils/Table';
import config from '@shared/utils/config';
import dayjs from 'dayjs';
import {computed, nextTick, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {VForm} from 'vuetify/components';
import {ApprovableInfo} from '@disclosure-portal/model/Approval';

const projectStore = useProjectStore();
const sbomStore = useSbomStore();
const {longText} = useRules();
const {t} = useI18n();
const idle = useIdleStore();
const {isAudited} = useApprovalCheck();

const isVisible = ref(false);
const selectedChannel = ref<VersionSlim | null>(null);
const sboms = ref<SpdxFile[]>([]);
const selectedSbom = ref<SpdxFile | null>(null);
const sbomStats = ref<ComponentStats>({} as ComponentStats);
const tab = ref('');
const approvableInfo = ref<ApprovableInfo>({} as ApprovableInfo);
const childProjectChannels = ref<Map<string, VersionSlim>>(new Map());
const comment = ref('');
const radioGroup = ref(0);
const c1 = ref(false);
const c2 = ref(false);
const c3 = ref(false);
const c4 = ref(false);
const c5 = ref(false);
const noFOSS = ref(false);
const withZip = ref(false);
const form = ref<VForm | null>(null);
const dd = ref();
const vehicle = ref(false);
const activePanel = ref<number | null>(null);
const fossVersion = ref<'default' | 'legacy'>('legacy');
const selectedProjects = ref<string[]>([]);
const allChannelSboms = ref<Map<string, SpdxFile[]>>(new Map());

const projectModel = computed(() => projectStore.currentProject!);
const channels = computed(() => {
  const res = Object.values(projectModel.value.versions);
  res.sort((a, b) => (dayjs(a.updated).isBefore(b.updated) ? 1 : -1));
  return res;
});
const countApprovables = computed(() => {
  if (!approvableInfo.value.projects) {
    return 0;
  }
  return approvableInfo.value.projects.filter(
    (p) => p.approvablespdx.spdxkey !== '' && p.approvablespdx.versionkey !== '',
  ).length;
});

const isRdConfirmationMissing = computed(() => {
  if (!vehicle.value) {
    return false;
  }

  if (!projectModel.value.isGroup) {
    const approvableSpdx = approvableInfo.value.projects?.[0]?.approvablespdx;

    if (!approvableSpdx?.spdxkey || !approvableSpdx?.versionkey) {
      return false;
    }

    const channel = channels.value.find((c) => c._key === approvableSpdx.versionkey);

    if (!channel) {
      return false;
    }

    const hasAuditedReview = channel.overallReviews?.some(
      (review) => review.sbomId === approvableSpdx.spdxkey && review.state === OverallReviewState.AUDITED,
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

const defaultC1 = () => {
  if (noFOSS.value) {
    return false;
  } else {
    return vehicle.value;
  }
};
const defaultC2 = () => {
  if (noFOSS.value) {
    return false;
  } else {
    if (vehicle.value) {
      return false;
    } else {
      return countApprovables.value > 0 || selectedSbom.value != null;
    }
  }
};
const defaultC3 = () => {
  if (noFOSS.value) {
    return false;
  } else {
    if (vehicle.value) {
      return false;
    } else {
      return !(countApprovables.value > 0);
    }
  }
};
const defaultC4 = () => {
  if (noFOSS.value) {
    return false;
  } else {
    return !vehicle.value;
  }
};
const defaultRadioGroup = () => {
  if (noFOSS.value) {
    return 3;
  } else {
    return 1;
  }
};

const setDefaultFlags = () => {
  radioGroup.value = defaultRadioGroup();
  c1.value = defaultC1();
  c2.value = defaultC2();
  c3.value = defaultC3();
  c4.value = defaultC4();
  c5.value = false;
};

const resetFormState = () => {
  selectedChannel.value = null;
  selectedSbom.value = null;
  sboms.value = [];
  sbomStats.value = new ComponentStats();
  comment.value = '';
  activePanel.value = null;
  radioGroup.value = 0;
  c1.value = false;
  c2.value = false;
  c3.value = false;
  c4.value = false;
  c5.value = false;
  noFOSS.value = false;
  withZip.value = false;
  childProjectChannels.value.clear();
  allChannelSboms.value.clear();
};

watch(isVisible, (newValue) => {
  if (!newValue) {
    resetFormState();
  }
});

watch(noFOSS, () => {
  setDefaultFlags();
  selectedChannel.value = null;
  selectedSbom.value = null;
  sbomStats.value = new ComponentStats();
});
watch(selectedSbom, () => {
  setDefaultFlags();
});
watch(radioGroup, () => {
  if (radioGroup.value == 3) {
    noFOSS.value = true;
  }
});

const stats = computed(() => {
  if (projectModel.value.isGroup) {
    return approvableInfo.value.stats;
  }
  return sbomStats.value;
});

const commentRule = longText(t('TAD_COMMENT'));

const open = async (isVehicle: boolean) => {
  idle.showIdle = true;
  approvableInfo.value = await projectService.getApprovableInfo(projectModel.value._key);

  vehicle.value = isVehicle;
  if (vehicle.value) {
    withZip.value = true;
  }
  noFOSS.value = projectModel.value.isNoFoss;
  setDefaultFlags();
  await autoSelect();

  if (!projectModel.value.isGroup) {
    allChannelSboms.value.clear();
    for (const channel of channels.value) {
      const versionEntry = sbomStore.getAllSBOMs.find((v) => v.versionKey === channel._key);
      allChannelSboms.value.set(channel._key, versionEntry?.spdxFileHistory ?? []);
    }
  }

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

const loadSBOMHist = async () => {
  selectedSbom.value = null;
  if (!selectedChannel.value?._key) return;
  await sbomStore.fetchAllSBOMsFlat();
  const versionEntry = sbomStore.getAllSBOMs.find((v) => v.versionKey === selectedChannel.value!._key);
  const spdxFileHistory = (versionEntry?.spdxFileHistory ?? []).slice(0, 5);
  if (spdxFileHistory[0]) {
    spdxFileHistory[0].isRecent = true;
  }
  sboms.value = spdxFileHistory;
};
const loadStats = async () => {
  if (!selectedChannel.value || !selectedSbom.value) {
    sbomStats.value = new ComponentStats();
    return;
  }
  sbomStats.value = (
    await versionService.getVersionComponentsForSbom(
      projectModel.value._key,
      selectedChannel.value?._key ?? '',
      selectedSbom.value?._key ?? '',
    )
  ).componentStats;
};
const autoSelect = async () => {
  if (channels.value.length === 0) {
    return;
  }

  if (approvableInfo.value.projects.length === 0) {
    return;
  }
  if (!noFOSS.value) {
    selectedChannel.value =
      channels.value.find((a) => a._key === approvableInfo.value.projects[0].approvablespdx.versionkey) ?? null;
  }
  if (!!sbomStore.selectedSBOMKey && !projectModel.value.isGroup) {
    selectedChannel.value = sbomStore.currentVersion;
  }
  if (selectedChannel.value) {
    await loadSBOMHist();
    if (sboms.value.length === 0) {
      return;
    }
    selectedSbom.value =
      sboms.value.find((a) => a._key === approvableInfo.value.projects[0].approvablespdx.spdxkey) ?? null;
    if (selectedSbom.value === null) {
      selectedSbom.value = sbomStore.getSelectedSBOM ?? null;
    }
    await loadStats();
  }
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
  metaDoc.c6 = noFOSS.value;

  const req: ExternalApprovalRequest = {
    comment: comment.value,
    guidProject: projectModel.value._key,
    metaDoc: metaDoc,
    withZip: withZip.value,
    fossVersion: 'vanilla',
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
  return vehicle.value && (stats.value.Denied > 0 || stats.value.NoAssertion > 0);
});

/**
 * Checks if the project is either an enterprise, mobile, or other platform project
 * aka **not** a vehicle project
 * and also checks if project is a group without vehicle
 * and contains denied or unasserted components.
 */
const isEnterpriseOrMobileOrOther = computed(() => {
  return !vehicle.value && (stats.value.Denied > 0 || stats.value.NoAssertion > 0);
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
            <Stack v-if="!projectModel.isGroup">
              <v-select
                v-model="selectedChannel"
                variant="outlined"
                item-title="name"
                return-object
                :label="t('SELECT_VERSION')"
                :items="channels"
                :disabled="noFOSS"
                hide-details
                @update:modelValue="loadSBOMHist" />
              <v-autocomplete
                v-model="selectedSbom"
                @update:modelValue="loadStats"
                :disabled="noFOSS"
                variant="outlined"
                item-title="name"
                :label="t('SELECT_SBOM_DELIVERY')"
                hide-details
                :items="sboms">
                <template v-slot:item="{item, props}">
                  <v-list-item v-bind="props" title="">
                    <div class="d-flex">
                      <div>
                        <v-icon
                          color="primary"
                          v-if="projectModel.approvablespdx.spdxkey == item.raw._key"
                          size="small"
                          class="pb-1"
                          >mdi-star</v-icon
                        >
                      </div>
                      <div>
                        <v-icon
                          color="green"
                          v-if="vehicle && isAudited(selectedChannel, item?.raw?._key)"
                          size="small"
                          class="ml-1 pb-1"
                          >mdi-clipboard-check-outline</v-icon
                        >
                      </div>
                      <span class="d-subtitle-2 ml-5">{{ formatDateAndTime(item.raw.uploaded) }}&nbsp;</span>
                      <span class="d-text d-secondary-text">&nbsp;-&nbsp;{{ item.raw.metaInfo.name }}</span>
                      <span class="d-text d-secondary-text" v-if="item.raw.tag">&nbsp;({{ item.raw.tag }})</span>
                      <span class="d-text d-secondary-text" v-if="item.raw.isRecent"
                        >&nbsp;{{ '[' + t('SBOM_LATEST') + ']' }}</span
                      >
                      <span class="d-text d-secondary-text" v-else>&nbsp;{{ '[' + t('SBOM_FORMER') + ']' }}</span>
                    </div>
                  </v-list-item>
                </template>
                <template v-slot:selection="{item}">
                  <div style="min-width: 13px">
                    <v-icon
                      color="primary"
                      v-if="projectModel.approvablespdx.spdxkey == item.raw._key"
                      size="small"
                      class="pb-1"
                      >mdi-star</v-icon
                    >
                  </div>
                  <div>
                    <v-icon
                      color="green"
                      v-if="vehicle && isAudited(selectedChannel, item?.raw?._key)"
                      size="small"
                      class="ml-1 pb-1"
                      >mdi-clipboard-check-outline</v-icon
                    >
                  </div>
                  <span class="d-subtitle-2 ml-5">{{ formatDateAndTime(item.raw.uploaded) }}&nbsp;</span>
                  <span class="d-text d-secondary-text">&nbsp;-&nbsp;{{ item.raw.metaInfo.name }}</span>
                  <span class="d-text d-secondary-text" v-if="item.raw.tag">&nbsp;({{ item.raw.tag }})</span>
                  <span class="d-text d-secondary-text" v-if="item.raw.isRecent"
                    >&nbsp;{{ '[' + t('SBOM_LATEST') + ']' }}</span
                  >
                  <span class="d-text d-secondary-text" v-else>&nbsp;{{ '[' + t('SBOM_FORMER') + ']' }}</span>
                </template>
              </v-autocomplete>
            </Stack>

            <section
              id="warning"
              v-if="isDeniedOrUnasserted || isEnterpriseOrMobileOrOther || noFOSS || isRdConfirmationMissing">
              <v-alert color="warning" type="warning">
                <span v-if="isDeniedOrUnasserted">
                  {{ t('DENIED_OR_UNASSARETED_MESSAGE') }}
                </span>
                <span v-else-if="isRdConfirmationMissing">
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

            <Stack v-if="config.useFutureFoss" direction="row" align="center" class="rounded bg-gray-500/20 py-1">
              <v-radio-group inline hide-details v-model="fossVersion">
                <v-radio :label="t('FOSSDD_STANDARD')" value="default"></v-radio>
                <v-radio :label="t('FOSSDD_LEGACY')" value="legacy"></v-radio>
              </v-radio-group>
              <v-spacer></v-spacer>
              <DIconButton icon="mdi-information-outline" :hint="t('FOSSDD_VERSION_TOOLTIP')" />
            </Stack>

            <v-tabs v-model="tab" slider-color="brand" show-arrows bg-color="tabsHeader">
              <v-tab value="general">{{ t('TAB_TITLE_GENERAL') }}</v-tab>
              <v-tab value="approvable" v-if="projectModel.isGroup">{{ t('TAB_TITLE_DETAILS') }}</v-tab>
            </v-tabs>
            <v-tabs-window v-model="tab">
              <v-tabs-window-item value="general">
                <DApprovalComponents
                  :stats="stats!"
                  :showRedWarnDeniedDecisionsMessage="showRedWarnDeniedDecisionsMessage" />
              </v-tabs-window-item>
              <v-tabs-window-item value="approvable" v-if="projectModel.isGroup">
                <GridSPDXList
                  :projects="approvableInfo.projects"
                  :channels="childProjectChannels"
                  showSbomExtras
                  selectable
                  @update:selectedProjects="selectedProjects = $event" />
              </v-tabs-window-item>
            </v-tabs-window>

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
