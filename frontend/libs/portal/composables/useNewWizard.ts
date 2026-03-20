// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {StepId, StepType, stepIds} from '@disclosure-portal/model/NewWizard';
import useRules from '@disclosure-portal/utils/Rules';
import {defineAsyncComponent} from 'vue';
import {useI18n} from 'vue-i18n';

/**
 * Validates a value against an array of Vuetify validation rules.
 * This function takes the validation rules (like those from useRules) and checks if the value passes all of them.
 * @param value - The value to validate
 * @param rules - Array of validation rule functions from Vuetify/useRules
 * @returns true if all validation rules pass, false otherwise
 */
export const isValidByRules = (
  value: string | undefined | null,
  rules: Array<(v: string) => boolean | string>,
): boolean => {
  if (!rules || rules.length === 0) return true;

  for (const rule of rules) {
    const result = rule(value || '');
    if (result !== true) {
      return false;
    }
  }
  return true;
};

/**
 * Converts an array of Pick<StepType, 'id' | 'i18nKey'> to StepType[] with sequential indices
 * @param partialSteps - Array of partial step definitions
 * @returns Array of complete StepType objects with sequential indices
 */
const buildSteps = (partialSteps: Pick<StepType, 'id' | 'i18nKey'>[]): StepType[] => {
  return partialSteps.map((step: Pick<StepType, 'id' | 'i18nKey'>, index: number) => ({
    id: step.id,
    i18nKey: step.i18nKey,
    index,
    isCompleted: false,
    errorText: '',
    seen: false,
  }));
};

/**
 * Removes a step from an existing steps array and updates indices of remaining steps
 * @param steps - Array of existing steps
 * @param stepId - ID of the step to remove
 * @returns New array of steps with updated indices
 */
export const removeStep = (steps: StepType[], stepId: StepId): StepType[] => {
  return steps
    .filter((step) => step.id !== stepId)
    .map((step, index) => ({
      ...step,
      index,
    }));
};

/**
 * Merges new steps with existing steps, preserving isCompleted and errorText from existing steps
 * @param newSteps - Array of new steps (from buildSteps)
 * @param existingSteps - Array of current steps with user progress
 * @returns Array of steps with updated indices but preserved completion status
 */
export const mergeSteps = (newSteps: StepType[], existingSteps: StepType[]): StepType[] => {
  const existingStepsMap = new Map(existingSteps.map((step) => [step.id, step]));
  return newSteps.map((newStep) => {
    const existingStep = existingStepsMap.get(newStep.id);

    if (existingStep) {
      return {
        ...newStep,
        isCompleted: existingStep.isCompleted,
        errorText: existingStep.errorText,
        seen: existingStep.seen,
      };
    }
    return newStep;
  });
};

export const useNewWizard = () => {
  const {t} = useI18n();
  const rules = useRules();

  // Define validation rules for wizard steps
  const validationRules = {
    name: rules.minMax(t('LBL_NAME'), 3, 80, false),
    description: rules.longText(t('NP_DIALOG_TF_DESCRIPTION')),
    address: rules.minMax(t('NP_DIALOG_TF_ADDRESS'), 3, 300, true),
    supplierNr: rules.minMax(t('NP_DIALOG_TF_SUPPLIER_NR'), 1, 25, true),
  };

  const stepPlatform: Pick<StepType, 'id' | 'i18nKey'> = {
    id: stepIds.platform,
    i18nKey: 'WIZARD_page_target_platform',
  };
  const stepDetails: Pick<StepType, 'id' | 'i18nKey'> = {
    id: stepIds.details,
    i18nKey: 'WIZARD_page_details',
  };
  const stepArchitecture: Pick<StepType, 'id' | 'i18nKey'> = {
    id: stepIds.architecture,
    i18nKey: 'WIZARD_page_architecture',
  };
  const stepTargetUsers: Pick<StepType, 'id' | 'i18nKey'> = {
    id: stepIds.targetUsers,
    i18nKey: 'WIZARD_page_target_users',
  };
  const stepDistributionTarget: Pick<StepType, 'id' | 'i18nKey'> = {
    id: stepIds.distributionTarget,
    i18nKey: 'WIZARD_page_distribution_target',
  };
  const stepDevelopment: Pick<StepType, 'id' | 'i18nKey'> = {
    id: stepIds.development,
    i18nKey: 'WIZARD_page_developer_target',
  };
  const stepOwner: Pick<StepType, 'id' | 'i18nKey'> = {
    id: stepIds.owner,
    i18nKey: 'WIZARD_page_development_owner',
  };
  const stepDeveloper: Pick<StepType, 'id' | 'i18nKey'> = {
    id: stepIds.developer,
    i18nKey: 'WIZARD_page_development_developer',
  };
  const stepSummary: Pick<StepType, 'id' | 'i18nKey'> = {
    id: stepIds.summary,
    i18nKey: 'WIZARD_page_summary',
  };
  const initSteps: StepType[] = buildSteps([stepPlatform]);
  const enterpriseOrMobileSteps = buildSteps([
    stepPlatform,
    stepDetails,
    stepArchitecture,
    stepTargetUsers,
    stepDistributionTarget,
    stepDevelopment,
    stepOwner,
    stepDeveloper,
    stepSummary,
  ]);

  const vehicleSteps = buildSteps([
    stepPlatform,
    stepDetails,
    stepArchitecture,
    stepDevelopment,
    stepOwner,
    stepDeveloper,
    stepSummary,
  ]);

  const otherSteps = buildSteps([stepPlatform, stepDetails, stepDevelopment, stepOwner, stepDeveloper, stepSummary]);

  const groupSteps = buildSteps([stepDetails, stepDevelopment, stepOwner, stepDeveloper, stepSummary]);

  const StepPlatform = defineAsyncComponent(
    () => import('@disclosure-portal/components/new-wizard/WizardStepPlatform.vue'),
  );
  const StepDetails = defineAsyncComponent(
    () => import('@disclosure-portal/components/new-wizard/WizardStepDetails.vue'),
  );
  const StepArchitecture = defineAsyncComponent(
    () => import('@disclosure-portal/components/new-wizard/WizardStepArchitecture.vue'),
  );
  const StepTargetUsers = defineAsyncComponent(
    () => import('@disclosure-portal/components/new-wizard/WizardStepTargetUsers.vue'),
  );
  const StepDistributionTarget = defineAsyncComponent(
    () => import('@disclosure-portal/components/new-wizard/WizardStepDistributionTarget.vue'),
  );
  const StepDevelopment = defineAsyncComponent(
    () => import('@disclosure-portal/components/new-wizard/WizardStepDevelopment.vue'),
  );
  const StepOwner = defineAsyncComponent(() => import('@disclosure-portal/components/new-wizard/WizardStepOwner.vue'));
  const StepDeveloper = defineAsyncComponent(
    () => import('@disclosure-portal/components/new-wizard/WizardStepDeveloper.vue'),
  );
  const StepSummary = defineAsyncComponent(
    () => import('@disclosure-portal/components/new-wizard/WizardStepSummary.vue'),
  );

  const componentMap = {
    [stepIds.platform]: StepPlatform,
    [stepIds.details]: StepDetails,
    [stepIds.architecture]: StepArchitecture,
    [stepIds.targetUsers]: StepTargetUsers,
    [stepIds.distributionTarget]: StepDistributionTarget,
    [stepIds.development]: StepDevelopment,
    [stepIds.owner]: StepOwner,
    [stepIds.developer]: StepDeveloper,
    [stepIds.summary]: StepSummary,
  };

  return {
    stepPlatform,
    stepDetails,
    stepArchitecture,
    stepTargetUsers,
    stepDistributionTarget,
    stepDevelopment,
    stepOwner,
    stepDeveloper,
    stepSummary,
    initSteps,
    enterpriseOrMobileSteps,
    vehicleSteps,
    otherSteps,
    buildSteps,
    removeStep,
    mergeSteps,
    componentMap,
    validationRules,
    groupSteps,
  };
};
