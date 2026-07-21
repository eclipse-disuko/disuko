import {SBOM} from '@cli/models/Sbom';
import {defineStore} from 'pinia';
import {computed, reactive, toRefs} from 'vue';
import type {Router} from 'vue-router';
import type {Project} from '../models/Project';
import {authService} from '../services/authService';
import {projectService} from '../services/projectService';

export const useAppStore = defineStore('app', () => {
  const initialAuthProjectUuid = authService.getAuth();

  const state = reactive({
    currentProject: null as Project | null,
    authProjectUUID: initialAuthProjectUuid ?? null,
    parentProject: initialAuthProjectUuid
      ? ({uuid: initialAuthProjectUuid, name: 'Parent'} as {uuid: string; name: string})
      : null,
  });

  // Getters
  const isAuthenticated = computed(() => Boolean(state.authProjectUUID));

  const setCurrentProject = (project: Project | null) => {
    state.currentProject = project;
  };

  const setParentProject = (parent: {uuid: string; name: string} | null) => {
    state.parentProject = parent;
  };

  const setAuth = (projectUuid: string) => {
    state.authProjectUUID = projectUuid;
    state.parentProject = {uuid: projectUuid, name: 'Parent'};
  };

  const setAuthProjectUuid = (projectUuid: string | null) => {
    state.authProjectUUID = projectUuid;
    state.parentProject = projectUuid ? {uuid: projectUuid, name: 'Parent'} : null;
  };

  const clearAuth = async () => {
    state.authProjectUUID = null;
    state.parentProject = null;
    await authService.logout();
  };

  const navigateToParent = async (router: Router) => {
    if (!state.parentProject) return;
    const parentUuid = state.parentProject.uuid;
    setParentProject(null);
    await router.push(`/projects/${parentUuid}/overview`);
  };

  const getLatestSbomByVersion = (project: Project | null, version: string): SBOM | null => {
    if (!project) return null;
    const versionSboms = (project.sboms || []).filter((sbom: SBOM) => {
      const sbomVersion = sbom.version || sbom.details?.version;
      return sbomVersion === version;
    });
    if (versionSboms.length === 0) {
      return null;
    }

    let latest: SBOM | null = null;
    for (const sbom of versionSboms) {
      if (!latest) {
        latest = sbom;
        continue;
      }
      const latestDate = new Date(latest.details?.uploaded || latest.updated || 0).getTime();
      const currentDate = new Date(sbom.details?.uploaded || sbom.updated || 0).getTime();
      if (currentDate > latestDate) {
        latest = sbom;
      }
    }
    return latest;
  };

  const navigateToSbom = async (project: Project, version: string, router: Router): Promise<boolean> => {
    if (!project?.uuid || !version) {
      return false;
    }

    const shouldFetch =
      !state.currentProject || state.currentProject.uuid !== project.uuid || !state.currentProject.sboms;
    if (shouldFetch) {
      await fetchCurrentProject(project.uuid);
    }

    const projectWithSboms = state.currentProject ?? project;
    const latestSbom = getLatestSbomByVersion(projectWithSboms, version);
    const sbomId = latestSbom?.id || latestSbom?.details?.id;
    const spdxParam = sbomId ? `/${encodeURIComponent(sbomId)}` : '';
    const path = `/projects/${encodeURIComponent(project.uuid)}/versions/${encodeURIComponent(version)}/components${spdxParam}`;
    await router.push(path);
    return true;
  };

  const getCurrentProject = (): Project | null => {
    return state.currentProject;
  };

  const fetchCurrentProject = async (projectUuid: string) => {
    const authProjectUuid = state.authProjectUUID ?? authService.getAuth();
    const authIsGroup = authService.getIsGroup();

    if (authIsGroup === true && authProjectUuid) {
      setParentProject({uuid: authProjectUuid, name: 'Parent'});
    } else if (authIsGroup === false) {
      setParentProject(null);
    }
    state.currentProject = await projectService.getProject(projectUuid);
  };

  const refetchCurrentProject = async (projectUuid: string) => {
    await fetchCurrentProject(projectUuid);
  };

  return {
    // State
    ...toRefs(state),

    // Actions
    setCurrentProject,
    setParentProject,
    navigateToParent,
    getCurrentProject,
    fetchCurrentProject,
    refetchCurrentProject,
    navigateToSbom,
    setAuth,
    setAuthProjectUuid,
    clearAuth,
    isAuthenticated,
  };
});
