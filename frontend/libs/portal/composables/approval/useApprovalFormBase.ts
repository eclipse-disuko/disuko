// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {ApprovableInfo} from '@disclosure-portal/model/Approval';
import {ComponentStats, SpdxFile, VersionSlim} from '@disclosure-portal/model/VersionDetails';
import versionService from '@disclosure-portal/services/version';
import {useIdleStore} from '@shared/stores/idle.store';
import {useJobStore} from '@disclosure-portal/stores/jobs';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useSbomStore} from '@disclosure-portal/stores/sbom.store';
import useRules from '@disclosure-portal/utils/Rules';
import dayjs from 'dayjs';
import {computed, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {VForm} from 'vuetify/components';

export interface ApprovalFormHooks {
  afterSetDefaultFlags?: () => void;
  afterResetFormState?: () => void;
  afterNoFOSSChange?: () => void;
  afterSelectedSbomChange?: () => void;
  afterAutoSelect?: () => void;
  computeGroupStats?: (approvableInfo: ApprovableInfo, selectedProjects: string[]) => ComponentStats;
  initFossVersion?: (isVehicle: boolean) => 'default' | 'legacy';
  defaultC1?: () => boolean;
  defaultC2?: (countApprovables: number, selectedSbom: SpdxFile | null) => boolean;
  defaultC3?: (countApprovables: number) => boolean;
  defaultC4?: () => boolean;
  defaultRadioGroup?: () => number;
  loadSBOMHistPreHook?: () => Promise<void>;
}

export function useApprovalFormBase(hooks: ApprovalFormHooks = {}) {
  const projectStore = useProjectStore();
  const sbomStore = useSbomStore();
  const idle = useIdleStore();
  const jobStore = useJobStore();
  const {longText} = useRules();
  const {t} = useI18n();

  const isVisible = ref(false);
  const selectedChannel = ref<VersionSlim | null>(null);
  const sboms = ref<SpdxFile[]>([]);
  const selectedSbom = ref<SpdxFile | null>(null);
  const sbomStats = ref<ComponentStats>({} as ComponentStats);
  const tab = ref<number | string>(0);
  const approvableInfo = ref<ApprovableInfo>({} as ApprovableInfo);
  const comment = ref('');
  const c1 = ref(false);
  const c2 = ref(false);
  const c3 = ref(false);
  const c4 = ref(false);
  const c5 = ref(false);
  const noFOSS = ref(false);
  const withZip = ref(false);
  const form = ref<VForm | null>(null);
  const fossVersion = ref<'default' | 'legacy'>('legacy');
  const selectedProjects = ref<string[]>([]);

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

  const stats = computed(() => {
    if (projectModel.value.isGroup) {
      if (hooks.computeGroupStats) {
        return hooks.computeGroupStats(approvableInfo.value, selectedProjects.value);
      }
      return approvableInfo.value.stats;
    }
    return sbomStats.value;
  });

  const commentRule = longText(t('TAD_COMMENT'));

  const setDefaultFlags = () => {
    c1.value = hooks.defaultC1 ? hooks.defaultC1() : false;
    c2.value = hooks.defaultC2
      ? hooks.defaultC2(countApprovables.value, selectedSbom.value)
      : noFOSS.value
        ? false
        : countApprovables.value > 0 || selectedSbom.value != null;
    c3.value = hooks.defaultC3
      ? hooks.defaultC3(countApprovables.value)
      : noFOSS.value
        ? false
        : !(countApprovables.value > 0);
    c4.value = hooks.defaultC4 ? hooks.defaultC4() : !noFOSS.value;
    c5.value = false;
    hooks.afterSetDefaultFlags?.();
  };

  const resetFormState = () => {
    selectedChannel.value = null;
    selectedSbom.value = null;
    sboms.value = [];
    sbomStats.value = new ComponentStats();
    comment.value = '';
    tab.value = 0;
    c1.value = false;
    c2.value = false;
    c3.value = false;
    c4.value = false;
    c5.value = false;
    noFOSS.value = false;
    withZip.value = false;
    hooks.afterResetFormState?.();
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
    hooks.afterNoFOSSChange?.();
  });

  watch(selectedSbom, () => {
    setDefaultFlags();
    hooks.afterSelectedSbomChange?.();
  });

  const loadSBOMHist = async () => {
    selectedSbom.value = null;
    if (!selectedChannel.value?._key) return;
    if (hooks.loadSBOMHistPreHook) {
      await hooks.loadSBOMHistPreHook();
    }
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
    hooks.afterAutoSelect?.();
  };

  return {
    isVisible,
    selectedChannel,
    sboms,
    selectedSbom,
    sbomStats,
    tab,
    approvableInfo,
    comment,
    c1,
    c2,
    c3,
    c4,
    c5,
    noFOSS,
    withZip,
    form,
    fossVersion,
    selectedProjects,
    projectModel,
    channels,
    countApprovables,
    stats,
    commentRule,
    setDefaultFlags,
    resetFormState,
    loadSBOMHist,
    loadStats,
    autoSelect,
    projectStore,
    sbomStore,
    idle,
    jobStore,
    t,
  };
}
