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
    currentVersion: {} as VersionSlim,
    channelSpdxs: [] as SpdxFile[],
    selectedSpdx: {} as SpdxFile,
    allSBOMSFlat: [] as VersionSbomsFlat[],
    allSBOMS: [] as VersionSboms[],
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
  const setCurrentVersion = (version: VersionSlim) => {
    state.currentVersion = version;
    clearSbomStats();
    clearGeneralStats();
  };

  const resetCurrentVersion = () => {
    if (!state.currentVersion?._key) return;
    const project = projectStore.currentProject;
    if (!project) return;
    state.currentVersion = project.versions[state.currentVersion._key];
  };

  const setSelectedSpdx = (spdx: SpdxFile) => {
    state.selectedSpdx = spdx;
    clearSbomStats();
  };

  const setChannelSpdxs = (spdxs: SpdxFile[]) => {
    state.channelSpdxs = spdxs;
  };

  const fetchAllSBOMsFlat = async () => {
    const projectKey = projectStore.currentProject?._key;
    if (!projectKey) return;
    const data = await ProjectService.getAllSbomsFlat(projectKey);
    state.allSBOMSFlat = data.items;
    state.allVersions = data.versions;
  };

  const fetchAllSBOMs = async () => {
    const projectKey = projectStore.currentProject?._key;
    if (!projectKey) return;
    state.allSBOMS = await ProjectService.getAllSboms(projectKey);
  };

  const fetchSBOMHistory = async (versionKey?: string) => {
    const projectKey = projectStore.currentProject?._key;
    const resolvedVersionKey = versionKey || state.currentVersion._key;
    if (!projectKey || !resolvedVersionKey) return;
    const spdxFileHistory = (await versionService.getSbomHistory(projectKey, resolvedVersionKey)).data;
    if (spdxFileHistory[0]) {
      spdxFileHistory[0].isRecent = true;
    }
    setChannelSpdxs(spdxFileHistory);
  };

  const fetchSBOMStats = async (spdxKey: string) => {
    const projectKey = projectStore.currentProject?._key;
    const versionKey = state.currentVersion._key;
    if (!projectKey || !versionKey || !spdxKey) return;
    if (Object.keys(state.sbomStats).length > 0 && state.selectedSpdx?._key === spdxKey) return;
    return versionService.getSBOMStats(projectKey, versionKey, spdxKey).then((data) => {
      if (state.currentVersion._key === versionKey && state.selectedSpdx?._key === spdxKey) {
        state.sbomStats = data.data;
      }
    });
  };

  const fetchGeneralVersionStats = async () => {
    const projectKey = projectStore.currentProject?._key;
    const versionKey = state.currentVersion._key;
    if (!projectKey || !versionKey) return;
    if (Object.keys(state.generalStats).length > 0) return;
    return versionService.getGeneralVersionStats(projectKey, versionKey).then((data) => {
      if (state.currentVersion._key === versionKey) {
        state.generalStats = data.data;
      }
    });
  };

  const reset = () => {
    state.currentVersion = {} as VersionSlim;
    state.channelSpdxs = [];
    state.selectedSpdx = {} as SpdxFile;
    state.allSBOMSFlat = [];
    state.allSBOMS = [];
    state.allVersions = [];
    clearSbomStats();
    clearGeneralStats();
  };

  // Getters
  const getCurrentVersion = computed(() => state.currentVersion);
  const getChannelSpdxs = computed(() => state.channelSpdxs);
  const getSelectedSpdx = computed(() => state.selectedSpdx);
  const getAllSBOMsFlat = computed(() => state.allSBOMSFlat);
  const getAllSBOMs = computed(() => state.allSBOMS);
  const getSbomStats = computed(() => state.sbomStats);
  const getGeneralStats = computed(() => state.generalStats);

  return {
    ...toRefs(state),

    // Actions
    setCurrentVersion,
    resetCurrentVersion,
    setSelectedSpdx,
    setChannelSpdxs,
    fetchAllSBOMsFlat,
    fetchAllSBOMs,
    fetchSBOMHistory,
    fetchSBOMStats,
    fetchGeneralVersionStats,
    reset,

    // Getters
    getCurrentVersion,
    getChannelSpdxs,
    getSelectedSpdx,
    getAllSBOMsFlat,
    getAllSBOMs,
    getSbomStats,
    getGeneralStats,
  };
});
