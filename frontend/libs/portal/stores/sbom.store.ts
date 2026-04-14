// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {NameKeyIdentifier, VersionSboms, VersionSbomsFlat} from '@disclosure-portal/model/ProjectsResponse';
import {GeneralStats, SbomStats, SpdxFile, VersionSlim} from '@disclosure-portal/model/VersionDetails';
import ProjectService from '@disclosure-portal/services/projects';
import versionService from '@disclosure-portal/services/version';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {defineStore} from 'pinia';
import {computed, reactive, toRefs} from 'vue';

export const useSbomStore = defineStore('sbom', () => {
  const projectStore = useProjectStore();

  const state = reactive({
    currentVersionKey: '' as string,
    selectedSBOMKey: '' as string,
    allSBOMSFlat: [] as VersionSbomsFlat[],
    allVersions: [] as NameKeyIdentifier[],
    sbomStats: {} as SbomStats,
    generalStats: {} as GeneralStats,
  });

  const clearSbomStats = () => {
    state.sbomStats = {} as SbomStats;
  };

  const clearGeneralStats = () => {
    state.generalStats = {} as GeneralStats;
  };

  // Actions
  const setCurrentVersion = (versionKey: string) => {
    state.currentVersionKey = versionKey;
    clearSbomStats();
    clearGeneralStats();
  };

  const setSelectedSBOMKey = (key: string) => {
    state.selectedSBOMKey = key;
    clearSbomStats();
  };

  const fetchAllSBOMsFlat = async (force?: boolean) => {
    const projectKey = projectStore.currentProject?._key;
    if (!projectKey) return;
    if (state.allSBOMSFlat.length > 0 && !force) return;
    const data = await ProjectService.getAllSbomsFlat(projectKey);
    state.allSBOMSFlat = data.items;
    state.allVersions = data.versions;
  };

  const fetchSBOMStats = async (spdxKey: string) => {
    const projectKey = projectStore.currentProject?._key;
    const versionKey = state.currentVersionKey;
    if (!projectKey || !versionKey || !spdxKey) return;
    if (Object.keys(state.sbomStats).length > 0 && state.selectedSBOMKey === spdxKey) return;
    return versionService.getSBOMStats(projectKey, versionKey, spdxKey).then((data) => {
      if (state.currentVersionKey === versionKey && state.selectedSBOMKey === spdxKey) {
        state.sbomStats = data.data;
      }
    });
  };

  const fetchGeneralVersionStats = async () => {
    const projectKey = projectStore.currentProject?._key;
    const versionKey = state.currentVersionKey;
    if (!projectKey || !versionKey) return;
    if (Object.keys(state.generalStats).length > 0) return;
    return versionService.getGeneralVersionStats(projectKey, versionKey).then((data) => {
      if (state.currentVersionKey === versionKey) {
        state.generalStats = data.data;
      }
    });
  };

  const reset = () => {
    state.currentVersionKey = '';
    state.selectedSBOMKey = '';
    state.allSBOMSFlat = [];
    state.allVersions = [];
    clearSbomStats();
    clearGeneralStats();
  };

  // Getters
  const currentVersion = computed((): VersionSlim =>
    projectStore.currentProject?.versions[state.currentVersionKey] ?? ({} as VersionSlim),
  );
  const getCurrentVersion = computed(() => currentVersion.value);
  const channelSpdxs = computed((): SpdxFile[] =>
    state.allSBOMSFlat
      .filter((item) => item.versionKey === state.currentVersionKey)
      .map((item, index) => ({...item, isRecent: index === 0})),
  );
  const getChannelSpdxs = computed(() => channelSpdxs.value);
  const getSelectedSBOM = computed(() => state.allSBOMSFlat.find((item) => item._key === state.selectedSBOMKey));
  const getAllSBOMsFlat = computed(() => state.allSBOMSFlat);
  const getAllSBOMs = computed((): VersionSboms[] => {
    const map = new Map<string, VersionSboms>();
    for (const item of state.allSBOMSFlat) {
      if (!map.has(item.versionKey)) {
        const entry = new VersionSboms();
        entry.VersionKey = item.versionKey;
        entry.VersionName = item.versionName;
        map.set(item.versionKey, entry);
      }
      map.get(item.versionKey)!.SpdxFileHistory.push(item);
    }
    return [...map.values()];
  });
  const getSbomStats = computed(() => state.sbomStats);
  const getGeneralStats = computed(() => state.generalStats);

  return {
    ...toRefs(state),

    // Actions
    setCurrentVersion,
    setSelectedSBOMKey,
    fetchAllSBOMsFlat,
    fetchSBOMStats,
    fetchGeneralVersionStats,
    reset,

    // Getters
    currentVersion,
    getCurrentVersion,
    channelSpdxs,
    getChannelSpdxs,
    getSelectedSBOM,
    getAllSBOMsFlat,
    getAllSBOMs,
    getSbomStats,
    getGeneralStats,
  };
});
