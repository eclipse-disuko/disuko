// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG // // SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import {useApprovalCheck} from '@disclosure-portal/composables/useApprovalCheck';
import {DocumentMeta, InternalApprovalRequest} from '@disclosure-portal/model/ApprovalRequest';
import ErrorDialogConfig from '@disclosure-portal/model/ErrorDialogConfig';
import {ApprovableInfoDto, ApprovableSPDXDto} from '@disclosure-portal/model/Project';
import {UserDto} from '@disclosure-portal/model/Users';
import {ComponentStats, SpdxFile, VersionSlim} from '@disclosure-portal/model/VersionDetails';
import projectService from '@disclosure-portal/services/projects';
import versionService from '@disclosure-portal/services/version';
import {useIdleStore} from '@disclosure-portal/stores/idle.store';
import {useJobStore} from '@disclosure-portal/stores/jobs';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useSbomStore} from '@disclosure-portal/stores/sbom.store';
import eventBus from '@disclosure-portal/utils/eventbus';
import useRules from '@disclosure-portal/utils/Rules';
import {formatDateAndTime} from '@disclosure-portal/utils/Table';
import useSnackbar from '@shared/composables/useSnackbar';
import config from '@shared/utils/config';
import dayjs from 'dayjs';
import {computed, nextTick, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {VForm} from 'vuetify/components';
import {useAppStore} from '@disclosure-portal/stores/app';

const projectStore = useProjectStore();
const appStore = useAppStore();
const sbomStore = useSbomStore();
const {longText} = useRules();
const {t} = useI18n();
const snackbar = useSnackbar();
const idle = useIdleStore();
const {isAudited} = useApprovalCheck();

const isVisible = ref(false);
const selectedChannel = ref<VersionSlim | null>(null);
const sboms = ref<SpdxFile[]>([]);
const selectedSbom = ref<SpdxFile | null>(null);
const sbomStats = ref<ComponentStats>({} as ComponentStats);
const tab = ref(0);
const approverTab = ref(0);
const approvableInfo = ref<ApprovableInfoDto>({} as ApprovableInfoDto);
const comment = ref('');
const c1 = ref(false);
const c2 = ref(false);
const c3 = ref(false);
const c4 = ref(false);
const c5 = ref(false);
const noFOSS = ref(false);
const withZip = ref(false);
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
const isVehicle = ref(false);
const fossVersion = ref<'default' | 'legacy'>('legacy');

const projectModel = computed(() => projectStore.currentProject!);
const channels = computed(() => {
  const res = Object.values(projectModel.value.versions);
  res.sort((a, b) => (dayjs(a.updated).isBefore(b.updated) ? 1 : -1));
  return res;
});
const countApprovables = computed(() => {
  if (!Array.isArray(approvableInfo.value.projects)) {
    return 0;
  }
  return approvableInfo.value.projects.filter((p) => {
    if (!p.approvablespdx) {
      return false;
    }
    const hasSpdxKey = p.approvablespdx.spdxkey !== '';
    const hasVersionKey = p.approvablespdx.versionkey !== '';
    return hasSpdxKey && hasVersionKey;
  }).length;
});
const defaultC2 = () => {
  if (noFOSS.value) {
    return false;
  } else {
    return countApprovables.value > 0 || selectedSbom.value != null;
  }
};
const defaultC3 = () => {
  if (noFOSS.value) {
    return false;
  } else {
    return !(countApprovables.value > 0);
  }
};
const defaultC4 = () => !noFOSS.value;

const setDefaultFlags = () => {
  c1.value = false;
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
  tab.value = 0;
  approverTab.value = 0;
  c1.value = false;
  c2.value = false;
  c3.value = false;
  c4.value = false;
  c5.value = false;
  noFOSS.value = false;
  withZip.value = false;
  ownerApprover1.value = '';
  ownerApprover2.value = '';
  developerApprover1.value = '';
  developerApprover2.value = '';
  ownerApproverPre1.value = undefined;
  ownerApproverPre2.value = undefined;
  developerApproverPre1.value = undefined;
  developerApproverPre2.value = undefined;
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

const stats = computed(() => {
  if (projectModel.value.isGroup) {
    return approvableInfo.value.stats;
  }
  return sbomStats.value;
});

const commentRule = longText(t('TAD_COMMENT'));

const open = async (isVehicleProject: boolean) => {
  idle.showIdle = true;
  isVehicle.value = isVehicleProject;
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
  approvableInfo.value = await projectService.getApprovableInfo(projectModel.value._key);

  await autoSelect();
  setDefaultFlags();
  developerApproverIn1.value?.resetForm();
  developerApproverIn2.value?.resetForm();
  idle.showIdle = false;
  isVisible.value = true;
};

const loadSBOMHist = async () => {
  selectedSbom.value = null;
  const spdxFileHistory = (await versionService.getSbomHistory(projectModel.value._key, selectedChannel.value!._key, 5))
    .data;
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
  if (Object.keys(sbomStore.selectedSpdx).length > 0 && !projectModel.value.isGroup) {
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
      selectedSbom.value = sbomStore.selectedSpdx ?? null;
    }
    await loadStats();
  }
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
        approverTab.value = 1;
        return;
      }
      if (!isOwner1Valid || !isOwner2Valid) {
        approverTab.value = 0;
        return;
      }
      if (
        (ownerApprover1.value !== '' || ownerApprover2.value !== '') &&
        ownerApprover1.value === ownerApprover2.value
      ) {
        approverTab.value = 0;
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
        approverTab.value = 0;
        const d = new ErrorDialogConfig();
        d.title = '' + t('SBOM_REQUEST_INTERNAL_APPROVAL');
        d.description = '' + t('BOTH_OR_NONE_OWNER_APPROVERS_ALLOWED_ERROR_MESSAGE');
        eventBus.emit('on-error', {error: d});
        return;
      }
      if (developerApprover1.value === developerApprover2.value) {
        approverTab.value = 1;
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
        fossVersion: 'vanilla',
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
          <v-tabs v-model="approverTab" slider-color="mbti" show-arrows bg-color="tabsHeader">
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
              variant="outlined"
              item-title="name"
              :label="t('SELECT_SBOM_DELIVERY')"
              :items="sboms"
              :disabled="noFOSS"
              hide-details>
              <template v-slot:item="{item, props}">
                <v-list-item v-bind="props" title="">
                  <div class="d-flex">
                    <v-icon
                      color="primary"
                      v-if="projectModel.approvablespdx.spdxkey == item.raw._key"
                      size="small"
                      class="pb-1">
                      mdi-star
                    </v-icon>
                    <div>
                      <v-icon
                        color="green"
                        v-if="isVehicle && isAudited(selectedChannel, item?.raw?._key)"
                        size="small"
                        class="pb-1 ml-1"
                        >mdi-clipboard-check-outline</v-icon
                      >
                    </div>
                    <span class="d-subtitle-2 ml-5">{{ formatDateAndTime(item.raw.Uploaded) }}&nbsp;</span>
                    <span class="d-text d-secondary-text">&nbsp;-&nbsp;{{ item.raw.MetaInfo.Name }}</span>
                    <span class="d-text d-secondary-text" v-if="item.raw.Tag">&nbsp;({{ item.raw.Tag }})</span>
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
                    v-if="isVehicle && isAudited(selectedChannel, item?.raw?._key)"
                    size="small"
                    class="pb-1 ml-1"
                    >mdi-clipboard-check-outline</v-icon
                  >
                </div>
                <span class="d-subtitle-2 ml-5">{{ formatDateAndTime(item.raw.Uploaded) }}&nbsp;</span>
                <span class="d-text d-secondary-text">&nbsp;-&nbsp;{{ item.raw.MetaInfo.Name }}</span>
                <span class="d-text d-secondary-text" v-if="item.raw.Tag">&nbsp;({{ item.raw.Tag }})</span>
                <span class="d-text d-secondary-text" v-if="item.raw.isRecent"
                  >&nbsp;{{ '[' + t('SBOM_LATEST') + ']' }}</span
                >
                <span class="d-text d-secondary-text" v-else>&nbsp;{{ '[' + t('SBOM_FORMER') + ']' }}</span>
              </template>
            </v-autocomplete>
          </Stack>

          <Stack v-if="config.useFutureFoss" direction="row" align="center" class="bg-gray-500/20 rounded py-1">
            <v-radio-group inline hide-details v-model="fossVersion">
              <v-radio :label="t('FOSSDD_STANDARD')" value="default"></v-radio>
              <v-radio :label="t('FOSSDD_LEGACY')" value="legacy"></v-radio>
            </v-radio-group>
            <v-spacer></v-spacer>
            <DIconButton icon="mdi-information-outline" :hint="t('FOSSDD_VERSION_TOOLTIP')" />
          </Stack>

          <v-tabs v-model="tab" slider-color="mbti" show-arrows bg-color="tabsHeader">
            <v-tab value="general">{{ t('TAB_TITLE_GENERAL') }}</v-tab>
            <v-tab value="approvable" v-if="projectModel.isGroup">{{ t('TAB_TITLE_DETAILS') }}</v-tab>
          </v-tabs>
          <v-tabs-window v-model="tab">
            <v-tabs-window-item value="general">
              <DApprovalComponents
                :stats="stats!"
                :showRedWarnDeniedDecisionsMessage="approvableInfo.hasDeniedDecisions" />
            </v-tabs-window-item>
            <v-tabs-window-item eager value="approvable" v-if="projectModel.isGroup">
              <GridSPDXList :projects="approvableInfo.projects" :channels="projectModel.versions" showSbomExtras />
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
