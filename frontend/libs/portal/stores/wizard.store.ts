// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {isValidByRules, removeStep, useNewWizard} from '@disclosure-portal/composables/useNewWizard';
import {
  architectures,
  createWizardProjectFromProject,
  CustomerMeta,
  developments,
  DisclosureDocumentMeta,
  StepId,
  stepIds,
  StepType,
  targetPlatforms,
  WizardProject,
} from '@disclosure-portal/model/NewWizard';
import type {Project} from '@disclosure-portal/model/Project';
import companyService from '@disclosure-portal/services/companies';
import projectService from '@disclosure-portal/services/projects';
import {useIdleStore} from '@disclosure-portal/stores/idle.store';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useUserStore} from '@disclosure-portal/stores/user';
import {deepCopy} from '@disclosure-portal/utils/Tools';
import useSnackbar from '@shared/composables/useSnackbar';
import {defineStore} from 'pinia';
import {computed, reactive, toRefs, watch} from 'vue';
import {useI18n} from 'vue-i18n';

type WizardState = 'create' | 'edit';

const createInitProjectWizard = (): WizardProject => ({
  targetPlatform: null,
  name: '',
  description: '',
  architecture: null,
  targetUsers: null,
  distributionTarget: null,
  development: null,
  owner: '',
  id: '',
  isDummy: false,
  isGroup: false,
  parentKey: '',
  projectSettings: {
    customerMeta: {},
    noticeContactMeta: {},
    documentMeta: {},
    supplierExtraData: {},
    noFossProject: false,
    customIds: [],
  },
  applicationMeta: null,
});

type CurrentStep = Pick<StepType, 'id' | 'index'>;

export const useWizardStore = defineStore('wizard', () => {
  const {initSteps, enterpriseOrMobileSteps, vehicleSteps, otherSteps, mergeSteps, validationRules, groupSteps} =
    useNewWizard();

  const state = reactive({
    isWizardOpen: false,
    project: createInitProjectWizard() as WizardProject,
    mode: 'create' as WizardState,
    currentStep: {id: stepIds.platform, index: 0} as CurrentStep,
    steps: initSteps as StepType[],
    previewLoading: false,
  });
  const {t} = useI18n();
  const idle = useIdleStore();
  const snack = useSnackbar();
  const projectStore = useProjectStore();
  const userStore = useUserStore();

  const loadWizardCards = async (projectKey: string) => {
    let cards = {};
    try {
      idle.show();
      cards = await projectService.getWizardByProjectKey(projectKey);
    } catch {
      snack.error(t('APPLICATION_MESSAGE_ERROR'));
      return;
    } finally {
      idle.hide();
    }
    return cards;
  };

  const prefillOwnerDepartment = async () => {
    const dept = await companyService.find(
      `${userStore.getProfile.metaData.department} ${userStore.getProfile.metaData.departmentDescription}`,
    );
    if (dept && dept.length > 0 && dept.find((dep) => dep.level === 0)) {
      state.project.projectSettings.customerMeta.dept = dept.find((dep) => dep.level === 0)!;
    }
  };

  const setDevelopmentFromParent = (parentProject: Project) => {
    if (parentProject.supplierExtraData.external) {
      state.project.development = developments.external;
      state.project.projectSettings.documentMeta.supplierName = parentProject?.documentMeta.supplierName;
      state.project.projectSettings.supplierExtraData.external = true;
    } else if (parentProject.documentMeta.supplierDept.companyName === parentProject.customerMeta.dept.companyName) {
      state.project.development = developments.inhouse;
      state.project.projectSettings.documentMeta = parentProject?.documentMeta as DisclosureDocumentMeta;
    } else {
      state.project.development = developments.internal;
      state.project.projectSettings.documentMeta = parentProject?.documentMeta as DisclosureDocumentMeta;
    }

    state.project.projectSettings.customerMeta = parentProject?.customerMeta as CustomerMeta;
  };

  const openWizard = async (params?: {
    parentProject?: Project;
    project?: Project;
    mode?: WizardState;
    isGroup?: boolean;
  }) => {
    // Group Project
    if (params?.isGroup) {
      state.project.isGroup = true;
      state.currentStep = {id: stepIds.details, index: 0};
      state.steps = deepCopy(groupSteps);
      setStepSeen(stepIds.details);
    }
    // Edit existing project
    if (params?.project && params?.mode === 'edit') {
      const wizardProject = createWizardProjectFromProject(params.project);
      let cards = await loadWizardCards(params.project._key);
      state.project = {...wizardProject, ...cards};
      setAvailableSteps();
    }

    // Create child project in a group
    const parentProject = params?.parentProject ?? null;
    if (parentProject) {
      state.project.parentKey = parentProject._key;
      setDevelopmentFromParent(parentProject);
    }

    // Create new project
    state.isWizardOpen = true;
    state.mode = params?.mode ?? 'create';
    if (state.mode === 'create') {
      await prefillOwnerDepartment();
    }
  };

  const close = () => {
    state.isWizardOpen = false;
    state.project = createInitProjectWizard();
    state.currentStep = {id: stepIds.platform, index: 0};
    state.steps = initSteps;
  };

  const preview = async () => {
    try {
      state.previewLoading = true;
      if (state.project.isGroup) {
        state.project = await projectService.previewGroupWizard(state.project);
      } else {
        state.project = await projectService.previewProjectWizard(state.project);
      }
    } catch (error) {
      console.error('Error previewing project:', error);
    } finally {
      state.previewLoading = false;
    }
  };

  const createProject = async () => {
    try {
      state.previewLoading = true;
      let project;
      if (state.project.isGroup) {
        project = await projectService.createGroupWizard(state.project);
      } else {
        project = await projectService.createProjectWizard(state.project);
      }
      snack.info(t('DIALOG_project_create_success'));
      return project;
    } catch (error) {
      console.error('Error creating project:', error);
      snack.info(t('WIZARD_error_create_project'));
    } finally {
      state.previewLoading = false;
    }
  };

  const updateProject = async () => {
    try {
      state.previewLoading = true;
      await projectService.updateProjectWizard(state.project, projectStore.currentProject!._key);
      await projectStore.fetchProjectByKey(projectStore.currentProject!._key);
      snack.info(t('DIALOG_project_edit_success'));
    } catch (error) {
      console.error('Error updating project:', error);
      snack.info(t('WIZARD_error_edit_project'));
    } finally {
      state.previewLoading = false;
    }
  };

  const setAvailableSteps = () => {
    let newSteps: StepType[];

    if (!state.project.isGroup) {
      if (isEnterpriseOrMobilePlatform.value) {
        newSteps = deepCopy(enterpriseOrMobileSteps);
      } else if (state.project.targetPlatform === targetPlatforms.vehicle) {
        newSteps = deepCopy(vehicleSteps);
      } else {
        newSteps = deepCopy(otherSteps);
      }

      state.steps = mergeSteps(newSteps, state.steps);
    } else {
      state.steps = mergeSteps(deepCopy(groupSteps), state.steps);
    }

    if (state.project.development === developments.inhouse) {
      state.steps = removeStep(state.steps, stepIds.developer);
    }
  };

  const getStepById = (id: StepId): StepType | undefined => {
    return state.steps.find((step) => step.id === id);
  };

  const setStepSeen = (stepId: StepId) => {
    const step = getStepById(stepId);
    if (!step || step.index === state.steps.length - 1) return;
    step.seen = true;
  };

  const selectStep = (stepIndex: number) => {
    const currentStep = getStepById(state.currentStep.id);
    if (currentStep && !currentStep?.isCompleted) {
      currentStep.errorText = t('VALIDATION_missing_selection');
    }
    if (stepIndex >= 0 && stepIndex < state.steps.length) {
      const stepId = state.steps[stepIndex].id;
      if (stepId) {
        state.currentStep = {
          id: stepId,
          index: stepIndex,
        };
      }

      setStepSeen(stepId);
    }
  };

  const nextStep = () => {
    selectStep(state.currentStep.index + 1);
  };

  const nextTwoSteps = () => {
    selectStep(state.currentStep.index + 2);
  };

  const previousStep = () => {
    selectStep(state.currentStep.index - 1);
  };

  const setStepCompleted = (id: StepId) => {
    setStepSeen(id);
    const step = getStepById(id);
    if (!step) return;
    step.isCompleted = true;
    step.errorText = '';
  };

  const setStepIncomplete = (id: StepId) => {
    const step = getStepById(id);
    if (!step) return;
    step.isCompleted = false;
  };

  const setStepCompletionStatus = (id: StepId, condition: boolean) => {
    if (condition) {
      setStepCompleted(id);
    } else {
      setStepIncomplete(id);
    }
  };

  const updateProjectSettingsBasedOnDevelopment = () => {
    if (!state.project.development) return;

    const settings = state.project.projectSettings;

    if (state.project.development === developments.inhouse) {
      settings.supplierExtraData.external = false;
      settings.documentMeta.supplierDept = settings.customerMeta.dept;
      settings.documentMeta.supplierAddress = settings.customerMeta.address;
      settings.documentMeta.supplierName = settings.customerMeta.dept?.companyName;
      settings.documentMeta.supplierNr = '';
    } else if (state.project.development === developments.external) {
      settings.supplierExtraData.external = true;
    } else if (state.project.development === developments.internal) {
      settings.supplierExtraData.external = false;
    }
  };

  const isCurrentStepCompleted = computed(() => {
    return getStepById(state.currentStep.id)?.isCompleted ?? false;
  });

  const isStepBeforeFinal = computed(() => {
    return state.currentStep.index === state.steps.length - 2;
  });

  const isFinalStep = computed(() => {
    return state.currentStep.index === state.steps.length - 1 && state.currentStep.index !== 0;
  });

  const allStepsCompleted = computed(() =>
    state.steps.filter((step) => step.id !== stepIds.summary).every((step) => step.isCompleted),
  );

  const isEnterpriseOrMobilePlatform = computed(
    () =>
      state.project.targetPlatform === targetPlatforms.enterprise ||
      state.project.targetPlatform === targetPlatforms.mobile,
  );
  const isVehiclePlatform = computed(() => state.project.targetPlatform === targetPlatforms.vehicle);
  const isOtherPlatform = computed(() => state.project.targetPlatform === targetPlatforms.other);

  const isVehicleArchitectures = computed(() => {
    return (
      state.project.architecture === architectures.vehicleOnboard ||
      state.project.architecture === architectures.vehicleOffboard
    );
  });
  const isVehicleOnboardArchitecture = computed(() => {
    return state.project.architecture === architectures.vehicleOnboard;
  });

  // Watcher to update step completion status based on project changes
  watch(
    () => state.project,
    () => {
      setStepCompletionStatus(stepIds.platform, !!state.project.targetPlatform);
      setStepCompletionStatus(stepIds.details, isValidByRules(state.project?.name, validationRules.name));
      setStepCompletionStatus(stepIds.architecture, !!state.project.architecture);
      setStepCompletionStatus(stepIds.targetUsers, !!state.project.targetUsers);
      setStepCompletionStatus(stepIds.distributionTarget, !!state.project.distributionTarget);
      setStepCompletionStatus(stepIds.development, !!state.project.development);
      setStepCompletionStatus(
        stepIds.owner,
        Object.keys(state.project.projectSettings?.customerMeta?.dept ?? {})?.length > 0 &&
          isValidByRules(state.project.projectSettings?.customerMeta?.address, validationRules.address) &&
          isValidByRules(state.project.projectSettings?.noticeContactMeta?.address, validationRules.address),
      );
      setStepCompletionStatus(
        stepIds.developer,
        Object.keys(state.project.projectSettings?.documentMeta?.supplierDept ?? {})?.length > 0 ||
          (isValidByRules(state.project.projectSettings?.documentMeta?.supplierName, validationRules.name) &&
            isValidByRules(state.project.projectSettings?.documentMeta?.supplierAddress, validationRules.address) &&
            isValidByRules(state.project.projectSettings?.documentMeta?.supplierNr, validationRules.supplierNr)),
      );
      setStepCompletionStatus(stepIds.summary, allStepsCompleted.value);
    },
    {deep: true, immediate: true},
  );

  return {
    ...toRefs(state),
    isFinalStep,
    isStepBeforeFinal,
    isCurrentStepCompleted,
    allStepsCompleted,
    isEnterpriseOrMobilePlatform,
    isVehiclePlatform,
    isOtherPlatform,
    isVehicleArchitectures,
    isVehicleOnboardArchitecture,
    //actions
    updateProject,
    createProject,
    preview,
    openWizard,
    close,
    setAvailableSteps,
    nextStep,
    previousStep,
    nextTwoSteps,
    selectStep,
    updateProjectSettingsBasedOnDevelopment,
  };
});
