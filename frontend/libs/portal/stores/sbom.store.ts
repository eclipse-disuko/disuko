// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {NameKeyIdentifier, VersionSboms, VersionSbomsFlat} from '@disclosure-portal/model/ProjectsResponse';
import {SpdxFile, VersionSlim} from '@disclosure-portal/model/VersionDetails';
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
  });

  // Actions
  const setCurrentVersion = (version: VersionSlim) => {
    state.currentVersion = version;
  };

  const resetCurrentVersion = () => {
    if (!state.currentVersion?._key) return;
    const project = projectStore.currentProject;
    if (!project) return;
    state.currentVersion = project.versions[state.currentVersion._key];
  };

  const setSelectedSpdx = (spdx: SpdxFile) => {
    state.selectedSpdx = spdx;
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

  const fetchSBOMHistory = async () => {
    const projectKey = projectStore.currentProject?._key;
    if (!projectKey || !state.currentVersion._key) return;
    const spdxFileHistory = (await versionService.getSbomHistory(projectKey, state.currentVersion._key)).data;
    if (spdxFileHistory[0]) {
      spdxFileHistory[0].isRecent = true;
    }
    setChannelSpdxs(spdxFileHistory);
  };

  const reset = () => {
    state.currentVersion = {} as VersionSlim;
    state.channelSpdxs = [];
    state.selectedSpdx = {} as SpdxFile;
    state.allSBOMSFlat = [];
    state.allSBOMS = [];
    state.allVersions = [];
  };

  // Getters
  const getCurrentVersion = computed(() => state.currentVersion);
  const getChannelSpdxs = computed(() => state.channelSpdxs);
  const getSelectedSpdx = computed(() => state.selectedSpdx);
  const getAllSBOMsFlat = computed(() => state.allSBOMSFlat);
  const getAllSBOMs = computed(() => state.allSBOMS);

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
    reset,

    // Getters
    getCurrentVersion,
    getChannelSpdxs,
    getSelectedSpdx,
    getAllSBOMsFlat,
    getAllSBOMs,
  };
});
