// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {IDefaultSelectItem} from '@disclosure-portal/model/ISelectItem';
import {createProjectModel, type Project, ProjectSubscriptions} from '@disclosure-portal/model/Project';
import ProjectPostRequest from '@disclosure-portal/model/ProjectPostRequest';
import {ProjectSlim} from '@disclosure-portal/model/ProjectsResponse';
import {WizardProjectPostRequest} from '@disclosure-portal/model/Wizard';
import {default as projectService} from '@disclosure-portal/services/projects';
import {useAppStore} from '@disclosure-portal/stores/app';
import {SearchOptions} from '@disclosure-portal/utils/Table';
import useSnackbar from '@shared/composables/useSnackbar';
import {defineStore} from 'pinia';
import {computed, reactive, toRefs, watch} from 'vue';
import {useI18n} from 'vue-i18n';

export enum ProjectStatusType {
  statusNew = 'ready',
  active = 'active',
  deprecated = 'deprecated',
}

export const useProjectStore = defineStore('project', () => {
  const {info} = useSnackbar();
  const {t} = useI18n();
  // TODO backwards-compatibility remove when appStore.currentProject is removed
  const appStore = useAppStore();

  const state = reactive({
    projects: [] as ProjectSlim[],
    projectsCount: 0,
    loading: false,
    currentProject: null as Project | null,
    hasVehiclePlatformChildren: false,
    hasOnlyVehiclePlatformChildren: false,
  });

  const projectPossibleStatuses = computed((): IDefaultSelectItem[] => {
    return Object.keys(ProjectStatusType).map((key) => ({
      text: key,
      value: ProjectStatusType[key as keyof typeof ProjectStatusType],
    }));
  });

  const resetCurrentProject = () => {
    state.currentProject = null;
    appStore.selectedSpdx = {};
    appStore.currentVersion = {};
    appStore.resetCurrentProject();
  };

  watch(
    () => state.currentProject,
    async () => {
      if (!state.currentProject) return;
      if (!state.currentProject.isGroup) return;
      state.hasVehiclePlatformChildren = (
        await projectService.getVehiclePlatform(state.currentProject._key)
      ).data.found;

      state.hasOnlyVehiclePlatformChildren = (
        await projectService.getVehiclePlatformOnly(state.currentProject._key)
      ).data.found;
    },
  );

  // API Calls
  const fetchProjects = async (options?: SearchOptions) => {
    if (state.projects.length === 0 || options) {
      state.loading = true;
    }
    try {
      const projectsResponse = options
        ? await projectService.getAllWithOptions(options)
        : await projectService.getAll();
      state.projects = projectsResponse.data.projects;
      state.projectsCount = projectsResponse.data.count;
    } catch (error) {
      console.error('Error fetching projects:', error);
    } finally {
      state.loading = false;
    }
  };

  const fetchProjectByKey = async (projectKey: string) => {
    try {
      state.loading = true;

      if (projectKey !== state.currentProject?._key) {
        resetCurrentProject();
      }

      const project = createProjectModel(await projectService.get(projectKey));

      appStore.currentProject = project;
      state.currentProject = project;
      if (project.isGroup) {
        state.currentProject.projectChildren = await projectService.getChildren(projectKey);
      }
    } catch (error) {
      console.error('Error fetching project:', error);
    } finally {
      state.loading = false;
    }
  };

  const createProject = async (project: WizardProjectPostRequest) => {
    try {
      state.loading = true;

      const response = (await projectService.create(project)).data;
      await fetchProjectByKey(response.id);

      info(t('DIALOG_project_create_success'));
      return response;
    } catch (error) {
      console.error('Error creating project:', error);
    } finally {
      state.loading = false;
    }
  };

  const updateProject = async (project: ProjectPostRequest) => {
    try {
      state.loading = true;

      await projectService.updateProject(project, project.id);
      await fetchProjectByKey(project.id);

      info(t('DIALOG_project_edit_success'));
    } catch (error) {
      console.error('Error updating project:', error);
      info(t('ERROR_SAVING_SETTINGS'));
    } finally {
      state.loading = false;
    }
  };

  const deleteProject = async (projectKey: string) => {
    try {
      await projectService.delete(projectKey);

      if (projectKey === state.currentProject?._key) {
        resetCurrentProject();
      }

      info(t('DIALOG_project_delete_success'));
    } catch (error) {
      console.error('Error deleting project:', error);
    }
    await fetchProjects();
  };

  const deprecateProject = async (projectKey: string) => {
    try {
      await projectService.deprecate(projectKey);

      if (projectKey === state.currentProject?._key) {
        await fetchProjectByKey(projectKey);
      }

      info(t('DIALOG_project_deprecate_success'));
    } catch (error) {
      console.error('Error deprecating project:', error);
    }
  };

  const updateProjectSubscriptions = async (projectKey: string, subscription: ProjectSubscriptions) => {
    if (!state.currentProject) return;
    try {
      state.currentProject.subscriptions = await projectService.saveProjectSubscriptions(projectKey, subscription);
      info(t('SUBSCRIPTIONS_UPDATED'));
    } catch (error) {
      console.error('Error updating project subscriptions:', error);
    }
  };

  const fetchProjectPossibleChildren = async (projectKey: string) => {
    try {
      state.projects = (await projectService.getPossibleChildren(projectKey)).projects;
    } catch (error) {
      console.error('Error fetching possible children for project:', error);
    }
  };

  const areMandatoryProjectSettingsSet = computed(() => {
    return !!(
      (state.currentProject?.customerMeta?.dept && state.currentProject?.documentMeta?.supplierDept) ||
      (state.currentProject?.customerMeta?.dept && state.currentProject?.supplierExtraData?.external)
    );
  });

  return {
    ...toRefs(state),
    projectPossibleStatuses,

    //Actions
    fetchProjects,
    fetchProjectByKey,
    createProject,
    updateProject,
    deleteProject,
    deprecateProject,
    resetCurrentProject,
    updateProjectSubscriptions,
    fetchProjectPossibleChildren,
    // Getters
    areMandatoryProjectSettingsSet,
  };
});
