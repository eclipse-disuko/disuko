// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useLicense} from '@disclosure-portal/composables/useLicense';
import {IDefaultSelectItem, IObligation} from '@disclosure-portal/model/IObligation';
import {compareFamily} from '@disclosure-portal/model/License';
import PolicyRule from '@disclosure-portal/model/PolicyRule';
import {CalculatedRuleConfigType} from '@disclosure-portal/model/CalculatedPolicyRules';
import AdminService from '@disclosure-portal/services/admin';
import useViewTools from '@disclosure-portal/utils/View';
import {defineStore} from 'pinia';
import {computed, reactive, toRefs} from 'vue';
import {useI18n} from 'vue-i18n';

const toBoolArray = (values: Array<string | boolean>): boolean[] => values.map((v) => v === true || v === 'true');

type ScopeFilterName = 'isLicenseChart' | 'approvalState' | 'family' | 'licenseType' | 'source';
type CalculatedBucketName = 'deniedClassifications' | 'warnedClassifications' | 'allowedClassifications';

export const useCalculatedPolicyRuleStore = defineStore('calculatedPolicyRule', () => {
  const {t} = useI18n();
  const {getLicenseApprovalTypes, getLicenseFamily, getLicenseTypes} = useLicense();
  const {getNameForLanguage} = useViewTools();

  const state = reactive({
    rule: new PolicyRule(),
    classifications: [] as IObligation[],
    classificationsLoaded: false,
  });

  const classificationOptions = computed(() =>
    state.classifications.filter((c) => c?._key).map((c) => ({text: getNameForLanguage(c), value: c._key})),
  );

  const licenseChartOptions = computed<IDefaultSelectItem[]>(() => [
    {text: t('TABLE_LICENSE_CHART_STATUS_IS'), value: 'true'},
    {text: t('TABLE_LICENSE_CHART_STATUS_IS_NOT'), value: 'false'},
  ]);

  const approvalStateOptions = computed<IDefaultSelectItem[]>(() =>
    getLicenseApprovalTypes().map((approval) => ({text: approval.text, value: approval.value || 'not set'})),
  );

  const familyOptions = computed<IDefaultSelectItem[]>(() =>
    getLicenseFamily()
      .map((family) => ({text: family.text, value: family.value || 'not declared'}))
      .sort((a, b) => compareFamily(a.value, b.value)),
  );

  const licenseTypeOptions = computed<IDefaultSelectItem[]>(() =>
    getLicenseTypes().map((type) => ({text: type.text, value: type.value || 'not declared'})),
  );

  const sourceOptions = computed<IDefaultSelectItem[]>(() => [
    {text: 'spdx', value: 'spdx'},
    {text: 'custom', value: 'custom'},
  ]);

  const getScopeFilterValues = (filterName: ScopeFilterName): Array<string | boolean> => {
    const value = state.rule.CalculatedConfig.licenseScope[filterName];
    return filterName === 'isLicenseChart' ? value.map(String) : (value as Array<string | boolean>);
  };

  const calculatedRuleConfig = computed<CalculatedRuleConfigType>(() => ({
    calculated: state.rule.Calculated,
    buckets: state.rule.CalculatedConfig.bucketDefinition,
    classificationOptions: classificationOptions.value,
    scopeConfig: {
      isLicenseChart: {
        options: licenseChartOptions.value,
        values: getScopeFilterValues('isLicenseChart'),
      },
      approvalState: {
        options: approvalStateOptions.value,
        values: getScopeFilterValues('approvalState'),
      },
      family: {
        options: familyOptions.value,
        values: getScopeFilterValues('family'),
      },
      licenseType: {
        options: licenseTypeOptions.value,
        values: getScopeFilterValues('licenseType'),
      },
      source: {
        options: sourceOptions.value,
        values: getScopeFilterValues('source'),
      },
    },
  }));

  const setRule = (newRule: PolicyRule): void => {
    state.rule = newRule;
  };

  const setCalculated = (value: boolean): void => {
    state.rule.Calculated = value;
  };

  const setBucketClassifications = (bucketName: CalculatedBucketName, value: string[]): void => {
    state.rule.CalculatedConfig.bucketDefinition[bucketName] = value;
  };

  const setScopeFilterValues = (filterName: ScopeFilterName, values: Array<string | boolean>): void => {
    const scope = state.rule.CalculatedConfig.licenseScope;
    if (filterName === 'isLicenseChart') {
      scope.isLicenseChart = toBoolArray(values);
    } else {
      scope[filterName] = values as string[];
    }
  };

  const retrieveClassifications = async (): Promise<void> => {
    try {
      state.classificationsLoaded = false;
      const response = (await AdminService.getAllObligations()).data;
      state.classifications = response.items;
    } catch (error) {
      console.error('Failed to retrieve classifications:', error);
      state.classifications = [];
    } finally {
      state.classificationsLoaded = true;
    }
  };

  const reset = (): void => {
    state.rule = new PolicyRule();
    state.classifications = [];
    state.classificationsLoaded = false;
  };

  return {
    ...toRefs(state),
    classificationOptions,
    licenseChartOptions,
    approvalStateOptions,
    familyOptions,
    licenseTypeOptions,
    sourceOptions,
    calculatedRuleConfig,
    setRule,
    setCalculated,
    setBucketClassifications,
    setScopeFilterValues,
    getScopeFilterValues,
    retrieveClassifications,
    reset,
  };
});
