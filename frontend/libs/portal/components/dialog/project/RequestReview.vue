<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {DocumentMeta, PlausibilityCheckRequest} from '@disclosure-portal/model/ApprovalRequest';
import {ApprovableSPDXDto} from '@disclosure-portal/model/Project';
import {UserDto} from '@disclosure-portal/model/Users';
import {ComponentStats, SpdxFile, VersionSlim} from '@disclosure-portal/model/VersionDetails';
import profileService from '@disclosure-portal/services/profile';
import projectService from '@disclosure-portal/services/projects';
import versionService from '@disclosure-portal/services/version';
import {useIdleStore} from '@disclosure-portal/stores/idle.store';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useSbomStore} from '@disclosure-portal/stores/sbom.store';
import useRules from '@disclosure-portal/utils/Rules';
import {formatDateAndTime} from '@disclosure-portal/utils/Table';
import DAutocompleteUser from '@shared/components/disco/DAutocompleteUser.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import dayjs from 'dayjs';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {VForm} from 'vuetify/components';
import {ApprovableInfo} from '@disclosure-portal/model/Approval';

const projectStore = useProjectStore();
const sbomStore = useSbomStore();
const {longText} = useRules();
const {t} = useI18n();
const snackbar = useSnackbar();
const idle = useIdleStore();

const isVisible = ref(false);
const selectedChannel = ref<VersionSlim | null>(null);
const sboms = ref<SpdxFile[]>([]);
const selectedSbom = ref<SpdxFile | null>(null);
const sbomStats = ref<ComponentStats>({} as ComponentStats);
const tab = ref('');
const approvableInfo = ref<ApprovableInfo>({} as ApprovableInfo);
const comment = ref('');
const approver = ref('');
const approverPreselect = ref<UserDto | undefined>(undefined);
const selectUserField = ref<InstanceType<typeof DAutocompleteUser>>();
const form = ref<VForm | null>(null);

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
const stats = computed(() => {
  if (projectModel.value.isGroup) {
    return approvableInfo.value.stats;
  }
  return sbomStats.value;
});

const isLoading = ref(false);

const dialogConfig = computed(() => ({
  title: t('SBOM_REQUEST_PLAUSIBILITY_CHECK'),
  secondaryButton: {text: t('BTN_CANCEL'), disabled: isLoading.value},
  primaryButton: {text: t('BTN_REQUEST'), disabled: isLoading.value, loading: isLoading.value},
}));

const activeRules = ref({
  comment: longText(t('TAD_COMMENT')),
  channel: [(v: VersionSlim | null) => v !== null || t('VERSION_REQUIRED')],
  sbom: [(v: SpdxFile | null) => v !== null || t('SBOM_REQUIRED')],
});

const open = async (payload?: {comment: string; reviewer: string}) => {
  if (payload) {
    comment.value = t(payload.comment);
    const users = (await profileService.getUsersBySearchFragment(payload.reviewer, true)).data;
    if (users.length > 0) {
      approverPreselect.value = users[0];
    }
  }
  approvableInfo.value = await projectService.getApprovableInfo(projectModel.value._key);

  await autoSelect();
  isVisible.value = true;
};

const loadSBOMHist = async () => {
  selectedSbom.value = null;
  if (!selectedChannel.value?._key) {
    return;
  }
  const versionEntry = sbomStore.getAllSBOMs.find((v) => v.VersionKey === selectedChannel.value!._key);
  const spdxFileHistory = (versionEntry?.SpdxFileHistory ?? []).slice(0, 5);
  if (spdxFileHistory[0]) {
    spdxFileHistory[0].isRecent = true;
  }
  sboms.value = spdxFileHistory;
};
const loadStats = async () => {
  if (!selectedChannel.value?._key || !selectedSbom.value?._key) {
    return;
  }
  sbomStats.value = (
    await versionService.getVersionComponentsForSbom(
      projectModel.value._key,
      selectedChannel.value._key,
      selectedSbom.value._key,
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

  if (Object.keys(sbomStore.selectedSpdx).length > 0 && !projectModel.value.isGroup) {
    selectedChannel.value = sbomStore.currentVersion;
  } else {
    selectedChannel.value =
      channels.value.find((a) => a._key === approvableInfo.value.projects[0].approvablespdx.versionkey) ?? null;
  }
  if (selectedChannel.value) {
    await loadSBOMHist();
    if (sboms.value.length === 0) {
      return;
    }
    selectedSbom.value =
      sboms.value.find((a) => a._key === approvableInfo.value.projects[0].approvablespdx.spdxkey) ?? null;
    if (Object.keys(sbomStore.selectedSpdx).length > 0) {
      selectedSbom.value = sbomStore.selectedSpdx ?? null;
    }
    await loadStats();
  }
};

const doDialogAction = async () => {
  const info = await form.value?.validate();
  const userFieldValid = await selectUserField.value?.validateOnCreate();
  if (!info?.valid || !userFieldValid) {
    return;
  }

  isLoading.value = true;
  idle.show();

  if (!projectModel.value.isGroup) {
    const approvableSpdx = {
      spdxkey: '',
      versionkey: '',
    } as ApprovableSPDXDto;
    approvableSpdx.spdxkey = selectedSbom.value?._key ?? '';
    approvableSpdx.versionkey = selectedChannel.value?._key ?? '';
    await projectService
      .updateApprovableSpdx(approvableSpdx, projectModel.value._key)
      .then(async () => await projectStore.fetchProjectByKey(projectModel.value._key));
  }

  const req: PlausibilityCheckRequest = new PlausibilityCheckRequest();
  req.approver = approver.value;
  req.comment = comment.value;
  req.guidProject = projectModel.value._key;

  req.metaDoc.c2 = countApprovables.value > 0;
  req.metaDoc.c3 = !(countApprovables.value > 0);
  req.metaDoc.c4 = true;
  if (projectModel.value.isNoFoss) {
    req.metaDoc = new DocumentMeta();
    req.metaDoc.c6 = true;
  }

  const response = await projectService.createPlausibilityCheck(req, projectModel.value._key);

  isLoading.value = false;
  idle.hide();

  if (response) {
    isVisible.value = false;
    snackbar.info(t('DIALOG_request_review_success'));
  }
};

const close = () => {
  isVisible.value = false;
  approverPreselect.value = undefined;
  comment.value = '';
  approver.value = '';
};

defineExpose({open});
</script>

<template>
  <v-form ref="form">
    <v-dialog v-model="isVisible" content-class="large" scrollable width="800">
      <DialogLayout :config="dialogConfig" @close="close" @secondary-action="close" @primary-action="doDialogAction">
        <Stack>
          <Stack :direction="projectModel.isGroup ? undefined : 'row'">
            <DAutocompleteUser
              class="w-1/2"
              ref="selectUserField"
              v-model="approver"
              :preselect="approverPreselect"
              :readonly="!!approverPreselect"
              :project-key="projectModel._key"
              :label="t('APPROVER_LABEL')"
              only-internal-users
              required />
            <v-select
              class="mb-auto"
              v-if="!projectModel.isGroup"
              v-model="selectedChannel"
              variant="outlined"
              item-title="name"
              return-object
              :label="t('SELECT_VERSION')"
              :items="channels"
              :rules="activeRules.channel"
              @update:modelValue="loadSBOMHist"
              hide-details
              required />
          </Stack>

          <v-autocomplete
            v-if="!projectModel.isGroup"
            v-model="selectedSbom"
            @update:modelValue="loadStats"
            variant="outlined"
            item-title="name"
            :label="t('SELECT_SBOM_DELIVERY')"
            :rules="activeRules.sbom"
            :items="sboms"
            hide-details
            required>
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
              <div style="min-width: 13px" class="d-flex">
                <v-icon
                  color="primary"
                  v-if="projectModel.approvablespdx.spdxkey == item.raw._key"
                  size="small"
                  class="pb-1"
                  >mdi-star</v-icon
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

          <v-tabs v-model="tab" slider-color="mbti" show-arrows bg-color="tabsHeader">
            <v-tab value="general">{{ t('TAB_TITLE_GENERAL') }}</v-tab>
            <v-tab value="approvable" v-if="projectModel.isGroup">{{ t('TAB_TITLE_DETAILS') }}</v-tab>
          </v-tabs>
          <v-tabs-window v-model="tab">
            <v-tabs-window-item value="general">
              <DApprovalComponents :stats="stats!" />
            </v-tabs-window-item>
            <v-tabs-window-item value="approvable" v-if="projectModel.isGroup">
              <GridSPDXList
                v-if="approvableInfo.projects && approvableInfo.projects.length > 0"
                :projects="approvableInfo.projects as any"
                :channels="projectModel.versions"
                showSbomExtras />
            </v-tabs-window-item>
          </v-tabs-window>

          <v-textarea
            v-model="comment"
            :rules="activeRules.comment"
            :label="t('TAD_COMMENT')"
            variant="outlined"
            counter="1000"
            rows="3"
            no-resize />
        </Stack>
      </DialogLayout>
    </v-dialog>
  </v-form>
</template>
