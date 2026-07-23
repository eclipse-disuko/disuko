// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {ComponentStats, SpdxFile, VersionSlim} from '@disclosure-portal/model/VersionDetails';
import versionService from '@disclosure-portal/services/version';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useSbomStore} from '@disclosure-portal/stores/sbom.store';
import dayjs from 'dayjs';
import {computed, ref, watch} from 'vue';
import {useApprovableInfoStore} from '@disclosure-portal/stores/approvableInfo.store';

export interface UseApprovalFormBaseOptions {
  setDefaultFlags: () => void;
  resetExtraState?: () => void;
  fetchFlat?: boolean;
}

export function useApprovalFormBase(options: UseApprovalFormBaseOptions) {
  const projectStore = useProjectStore();
  const sbomStore = useSbomStore();
  const approvableInfoStore = useApprovableInfoStore();

  const isVisible = ref(false);
  const selectedChannel = ref<VersionSlim | null>(null);
  const sboms = ref<SpdxFile[]>([]);
  const selectedSbom = ref<SpdxFile | null>(null);
  const sbomStats = ref<ComponentStats>(new ComponentStats());
  const approvableInfo = computed(() => approvableInfoStore.approvableInfo);
  const comment = ref('');
  const withZip = ref(false);
  const noFOSS = ref(false);
  const fossVersion = ref<'default' | 'legacy'>('legacy');
  const mixedFOSS = ref(false);
  const c1 = ref(false);
  const c2 = ref(false);
  const c3 = ref(false);
  const c4 = ref(false);
  const c5 = ref(false);
  const selectedProjects = ref<string[]>([]);
  const tab = ref<'general' | 'approvable'>('general');

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
      if (!approvableInfo.value?.projects || selectedProjects.value.length === 0) {
        return new ComponentStats();
      }

      const selectedProjectsSet = new Set(selectedProjects.value);
      const aggregatedStats = new ComponentStats();

      approvableInfo.value.projects
        .filter((p) => selectedProjectsSet.has(p.projectKey))
        .forEach((project) => {
          const s = project.stats;
          if (s) {
            aggregatedStats.total += s.total || 0;
            aggregatedStats.allowed += s.allowed || 0;
            aggregatedStats.warned += s.warned || 0;
            aggregatedStats.denied += s.denied || 0;
            aggregatedStats.questioned += s.questioned || 0;
            aggregatedStats.noAssertion += s.noAssertion || 0;
          }
        });

      return aggregatedStats;
    }
    return sbomStats.value;
  });

  const selectedProjectsContainEmptySbom = computed(() => {
    if (fossVersion.value === 'legacy') {
      return false;
    }

    if (selectedChannel.value && selectedSbom.value) {
      return false;
    }

    return (
      approvableInfo.value.projects?.some(
        (project) =>
          selectedProjects.value.includes(project.projectKey) &&
          (!project.approvablespdx.spdxkey || !project.approvablespdx.versionkey),
      ) ?? false
    );
  });

  const updateSelectedProjects = () => {
    if (!approvableInfo.value?.projects) {
      return;
    }

    if (fossVersion.value === 'legacy') {
      selectedProjects.value = approvableInfo.value.projects
        .filter((p) => p.isNonFoss === noFOSS.value)
        .map((p) => p.projectKey);
    } else {
      selectedProjects.value = approvableInfo.value.projects.map((p) => p.projectKey);
    }
  };

  const checkFossMixedStatus = () => {
    if (projectModel.value.isGroup && fossVersion.value === 'legacy') {
      for (const project of approvableInfo.value.projects) {
        if (project.isNonFoss !== projectModel.value.isNoFoss) {
          mixedFOSS.value = true;
          return;
        }
      }
    }
    mixedFOSS.value = false;
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

  const loadSBOMHist = async () => {
    selectedSbom.value = null;
    if (!selectedChannel.value?._key) return;
    if (options.fetchFlat) {
      await sbomStore.fetchAllSBOMsFlat();
    }
    const versionEntry = sbomStore.getAllSBOMs.find((v) => v.versionKey === selectedChannel.value!._key);
    const spdxFileHistory = (versionEntry?.spdxFileHistory ?? []).slice(0, 5);
    if (spdxFileHistory[0]) {
      spdxFileHistory[0].isRecent = true;
    }
    sboms.value = spdxFileHistory;
  };

  const autoSelect = async () => {
    if (channels.value.length === 0) {
      return;
    }

    if (approvableInfo.value.projects.length === 0) {
      return;
    }

    if (!!sbomStore.selectedSBOMKey && !projectModel.value.isGroup) {
      selectedChannel.value = sbomStore.currentVersion;
    } else if (!noFOSS.value) {
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
      if (!!sbomStore.selectedSBOMKey) {
        selectedSbom.value = sbomStore.getSelectedSBOM ?? null;
      }
      await loadStats();
    }
  };

  const resetBase = () => {
    selectedChannel.value = null;
    sboms.value = [];
    selectedSbom.value = null;
    sbomStats.value = new ComponentStats();
    comment.value = '';
    withZip.value = false;
    c1.value = false;
    c2.value = false;
    c3.value = false;
    c4.value = false;
    c5.value = false;
    noFOSS.value = false;
    tab.value = 'general';
    selectedProjects.value = [];
    mixedFOSS.value = false;
  };

  watch(isVisible, (v) => {
    if (!v) {
      resetBase();
      options.resetExtraState?.();
    }
  });

  watch(noFOSS, () => {
    options.setDefaultFlags();
    selectedChannel.value = null;
    selectedSbom.value = null;
    sbomStats.value = new ComponentStats();
    updateSelectedProjects();
  });

  watch(selectedSbom, () => {
    options.setDefaultFlags();
  });

  watch(
    fossVersion,
    () => {
      updateSelectedProjects();
      checkFossMixedStatus();
    },
    {immediate: true},
  );

  return {
    isVisible,
    selectedChannel,
    sboms,
    selectedSbom,
    sbomStats,
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
    resetBase,
  };
}
