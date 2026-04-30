<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import ClassificationsPerLicenseDialog from '@disclosure-portal/components/dialog/ClassificationsPerLicenseDialog.vue';
import CalculatedRuleConfig from '@disclosure-portal/components/policy-rules/CalculatedRuleConfig.vue';
import {useLicense} from '@disclosure-portal/composables/useLicense';
import {IDefaultSelectItem, IObligation} from '@disclosure-portal/model/IObligation';
import Label from '@disclosure-portal/model/Label';
import {compareFamily, LicenseSlim} from '@disclosure-portal/model/License';
import PolicyRule, {PolicyRules, PolicyState} from '@disclosure-portal/model/PolicyRule';
import {compareLevel, levelWeight} from '@disclosure-portal/model/Quality';
import {Rights} from '@disclosure-portal/model/Rights';
import AdminService from '@disclosure-portal/services/admin';
import licenseService from '@disclosure-portal/services/license';
import policyRuleService from '@disclosure-portal/services/policyrules';
import ProjectService from '@disclosure-portal/services/projects';
import {useUserStore} from '@disclosure-portal/stores/user';
import {removeFromList} from '@disclosure-portal/utils/List';
import {RightsUtils} from '@disclosure-portal/utils/Rights';
import useViewTools, {getIconColorOfLevel, getIconOfLevel, IMap, openUrlInNewTab} from '@disclosure-portal/utils/View';
import {useCalculatedPolicyRuleStore} from '@disclosure-portal/stores/calculatedPolicyRule.store';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import DSearchField from '@shared/components/disco/DSearchField.vue';
import DCloseButton from '@shared/components/disco/DCloseButton.vue';
import DIconButton from '@shared/components/disco/DIconButton.vue';
import DLicenseChartIcon from '@shared/components/disco/DLicenseChartIcon.vue';
import DRuleButtons from '@shared/components/disco/DRuleButtons.vue';
import {IRuleBtnCallbacks} from '@shared/components/disco/interfaces';
import useSnackbar from '@shared/composables/useSnackbar';
import TableLayout from '@shared/layouts/TableLayout.vue';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {DataTableHeader, DataTableItem} from '@shared/types/table';
import {TOOLTIP_OPEN_DELAY_IN_MS} from '@shared/utils/constant';
import _, {indexOf} from 'lodash';
import {computed, onMounted, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRouter} from 'vue-router';
import {storeToRefs} from 'pinia';

const {t} = useI18n();
const {getI18NTextOfPrefixKey} = useLicense();
const {getNameForLanguage} = useViewTools();
const router = useRouter();
const calculatedPolicyRuleStore = useCalculatedPolicyRuleStore();
const {rule, classifications, classificationsLoaded} = storeToRefs(calculatedPolicyRuleStore);
const isPolicyManager = ref(false);
const rights = ref(new Rights());
const userStore = useUserStore();
const ruleId = ref('');
const breadcrumbs = useBreadcrumbsStore();
const hasChanges = ref(false);
const filterUnSelected = ref('');
const filterSelected = ref('');
const licensesLoading = ref(true);
const ruleLoaded = ref(false);
const mode = ref(PolicyState.ALLOW);
const menuIsLicenseChartSelected = ref(false);
const notSelectedLicenses = ref(<LicenseSlim[]>[]);
const selectedFilterIsLicenseChartNotSelected = ref<string[]>([]);
const selectedFilterClassificationsNotSelected = ref<string[]>([]);
const selectedFilterFamilyNotSelected = ref<string[]>([]);
const selectedFilterApprovalNotSelected = ref<string[]>([]);
const selectedFilterTypeNotSelected = ref<string[]>([]);
const selectedFilterIsLicenseChartSelected = ref<string[]>([]);
const possibleIsLicenseChartSelected = ref<IDefaultSelectItem[]>([]);
const allLicenses = ref<LicenseSlim[]>([]);
const selectedLicenses = ref<LicenseSlim[]>([]);
const possibleFamilySelected = ref<IDefaultSelectItem[]>([]);
const selectedFilterFamilySelected = ref<string[]>([]);
const possibleApprovalSelected = ref<IDefaultSelectItem[]>([]);
const selectedFilterApprovalSelected = ref<string[]>([]);
const possibleTypeSelected = ref<IDefaultSelectItem[]>([]);
const selectedFilterTypeSelected = ref<string[]>([]);
const possibleClassificationsSelected = ref<IDefaultSelectItem[]>([]);
const selectedFilterClassificationsSelected = ref<string[]>([]);
const possibleIsLicenseChartNotSelected = ref<IDefaultSelectItem[]>([]);
const possibleFamilyNotSelected = ref<IDefaultSelectItem[]>([]);
const possibleApprovalNotSelected = ref<IDefaultSelectItem[]>([]);
const possibleTypeNotSelected = ref<IDefaultSelectItem[]>([]);
const possibleClassificationsNotSelected = ref<IDefaultSelectItem[]>([]);
const labelsMap = ref<IMap<Label>>({});
const policyLabels = ref<Label[]>([]);
const menu6 = ref(false);
const menu5 = ref(false);
const menu4 = ref(false);
const menu3 = ref(false);
const menu2 = ref(false);
const menu = ref(false);
const menuClassification = ref(false);
const menuClassificationNot = ref(false);
const menuIsLicenseChartNotSelected = ref(false);
const {info} = useSnackbar();
const classificationsDialogRef = ref<InstanceType<typeof ClassificationsPerLicenseDialog> | null>(null);

const canEditManual = computed(() => isPolicyManager.value && !rule.value.Deprecated && !rule.value.Calculated);
const canEditCalculated = computed(() => isPolicyManager.value && rule.value.Calculated);

const retrieveRule = async (policyRuleId: string) => {
  if (router.currentRoute.value.params.uuid) {
    rule.value = new PolicyRule(
      (await ProjectService.getProjectPolicyRule(<string>router.currentRoute.value.params.uuid, policyRuleId)).data,
    );
  } else {
    rule.value = new PolicyRule((await policyRuleService.getPolicyRule(policyRuleId)).data);
  }
  calculatedPolicyRuleStore.setRule(rule.value);
  ruleLoaded.value = true;
};

const initBreadcrumbs = () => {
  breadcrumbs.setCurrentBreadcrumbs([
    {title: t('BC_Dashboard'), disabled: false, href: '/dashboard/home'},
    {title: t('POLICY_RULES'), disabled: false, href: '/dashboard/policyrules'},
    {
      title: '' + rule.value.Name,
      disabled: false,
      href: '/dashboard/policyrules/' + encodeURIComponent(ruleId.value),
    },
  ]);
};

const reload = async () => {
  await retrieveRule(rule.value._key);
  initBreadcrumbs();
  hasChanges.value = false;
};
watch(selectedFilterClassificationsSelected, reload);

// Track changes to calculated policy rule state for UI dirty flag
watch(
  () => calculatedPolicyRuleStore.rule,
  () => {
    hasChanges.value = true;
  },
  {deep: true},
);

const getActiveClassForPolicyFilterBtn = (policy: PolicyState): string => {
  switch (policy) {
    case PolicyState.DENY:
      return 'deny-border';
    case PolicyState.WARN:
      return 'warning-border';
    case PolicyState.ALLOW:
      return 'allow-border';
    default:
      return '';
  }
};

const getCssClass = () => (canEditManual.value ? 'force-border' : 'force-border tableNoHandCursor');

const onACFocus = (a: FocusEvent) => {
  const clickEvent = new Event('click', {bubbles: true});
  if (a.target) {
    a.target.dispatchEvent(clickEvent);
  }
};

const filteredListNotSelected = computed(() =>
  _.chain(notSelectedLicenses.value)
    .filter(filterOnIsLicenseChartNotSelected)
    .filter(filterOnClassificationNotSelected)
    .filter(filterOnFamilyNotSelected)
    .filter(filterOnApprovalNotSelected)
    .filter(filterOnTypeNotSelected)
    .value(),
);

const filterOnIsLicenseChartNotSelected = (item: LicenseSlim): boolean =>
  selectedFilterIsLicenseChartNotSelected.value.length === 0 ||
  indexOf(selectedFilterIsLicenseChartNotSelected.value, item.meta.isLicenseChart + '') !== -1;

const filterOnClassificationNotSelected = (item: LicenseSlim): boolean => {
  if (selectedFilterClassificationsNotSelected.value.length === 0) {
    return true;
  }
  if (!item.meta.classifications || item.meta.classifications.length === 0) {
    return selectedFilterClassificationsNotSelected.value.includes('');
  }
  return item.meta.classifications.some((classification) => {
    const classificationName = classification ? useViewTools().getNameForLanguage(classification) : '';
    return classificationName && selectedFilterClassificationsNotSelected.value.includes(classificationName);
  });
};

const filterOnFamilyNotSelected = (item: LicenseSlim): boolean =>
  selectedFilterFamilyNotSelected.value.length === 0 ||
  indexOf(selectedFilterFamilyNotSelected.value, item.meta.family) !== -1;

const filterOnApprovalNotSelected = (item: LicenseSlim): boolean =>
  selectedFilterApprovalNotSelected.value.length === 0 ||
  indexOf(selectedFilterApprovalNotSelected.value, item.meta.approvalState) !== -1;

const filterOnTypeNotSelected = (item: LicenseSlim): boolean =>
  selectedFilterTypeNotSelected.value.length === 0 ||
  indexOf(selectedFilterTypeNotSelected.value, item.meta.licenseType) !== -1;

const retrieveSpdxLicenses = async () => {
  licensesLoading.value = true;
  allLicenses.value = (await licenseService.getAll()).data.map((license) => new LicenseSlim(license));
  fillTables(true);
  licensesLoading.value = false;
};

const fillTables = (updateNotSelectedLicenses: boolean) => {
  if (updateNotSelectedLicenses) {
    notSelectedLicenses.value = [];
  }
  selectedLicenses.value = [];

  allLicenses.value.forEach((license) => {
    let selected = false;
    for (const policyState of PolicyRules) {
      selected = fillLicenseTbl(policyState, license) || selected;
    }
    if (updateNotSelectedLicenses && !selected) {
      notSelectedLicenses.value.push(license);
    }
  });
  refreshPossible(selectedLicenses.value);
  refreshPossible(notSelectedLicenses.value, false);
};

const refreshPossible = (items: LicenseSlim[], selected = true) => {
  if (selected) {
    possibleIsLicenseChartSelected.value = getPossibleIsLicenseChart(items, true);
    possibleFamilySelected.value = getPossibleFamily(items, true);
    possibleApprovalSelected.value = getPossibleApproval(items, true);
    possibleTypeSelected.value = getPossibleType(items, true);
    possibleClassificationsSelected.value = getPossibleClassifications(items, true);
  } else {
    possibleIsLicenseChartNotSelected.value = getPossibleIsLicenseChart(items, false);
    possibleFamilyNotSelected.value = getPossibleFamily(items, false);
    possibleApprovalNotSelected.value = getPossibleApproval(items, false);
    possibleTypeNotSelected.value = getPossibleType(items, false);
    possibleClassificationsNotSelected.value = getPossibleClassifications(items, false);
  }
};

const getPossibleClassifications = (items: LicenseSlim[], selected = true): IDefaultSelectItem[] => {
  let alreadySelected = selectedFilterClassificationsSelected.value;
  if (!selected) {
    alreadySelected = selectedFilterClassificationsNotSelected.value;
  }
  return _.chain(items)
    .map((item: LicenseSlim) => {
      if (!item.meta.classifications || item.meta.classifications.length === 0) {
        return [
          {
            text: t('NO_CLASSIFICATIONS'),
            value: '',
          },
        ] as IDefaultSelectItem[];
      }
      return _.map(
        item.meta.classifications,
        (c: IObligation) =>
          ({
            text: getNameForLanguage(c),
            value: getNameForLanguage(c),
          }) as IDefaultSelectItem,
      );
    })
    .flatten()
    .union(
      _.map(
        alreadySelected,
        (selectedValue: string) =>
          ({
            text: selectedValue,
            value: selectedValue,
          }) as IDefaultSelectItem,
      ),
    )
    .uniqBy('value')
    .value();
};

const getPossibleType = (items: LicenseSlim[], selected = true): IDefaultSelectItem[] => {
  let alreadySelected = selectedFilterTypeSelected.value;
  if (!selected) {
    alreadySelected = selectedFilterTypeNotSelected.value;
  }

  return _.chain(items)
    .uniqBy('meta.licenseType')
    .map(
      (item: LicenseSlim) =>
        ({
          text: getI18NTextOfPrefixKey('LT_', item.meta.licenseType),
          value: item.meta.licenseType,
        }) as IDefaultSelectItem,
    )
    .union(
      _.map(
        alreadySelected,
        (selectedValue: string) =>
          ({
            text: getI18NTextOfPrefixKey('LT_', selectedValue),
            value: selectedValue,
          }) as IDefaultSelectItem,
      ),
    )
    .value();
};

const getPossibleApproval = (items: LicenseSlim[], selected = true): IDefaultSelectItem[] => {
  let alreadySelected = selectedFilterApprovalSelected.value;
  if (!selected) {
    alreadySelected = selectedFilterApprovalNotSelected.value;
  }

  return _.chain(items)
    .uniqBy('meta.approvalState')
    .map(
      (item: LicenseSlim) =>
        ({
          text: getI18NTextOfPrefixKey('LT_APP_', item.meta.approvalState),
          value: item.meta.approvalState,
        }) as IDefaultSelectItem,
    )
    .union(
      _.map(
        alreadySelected,
        (selectedValue: string) =>
          ({
            text: getI18NTextOfPrefixKey('LT_APP_', selectedValue),
            value: selectedValue,
          }) as IDefaultSelectItem,
      ),
    )
    .value();
};

const getPossibleFamily = (items: LicenseSlim[], selected = true): IDefaultSelectItem[] => {
  let alreadySelected = selectedFilterFamilySelected.value;
  if (!selected) {
    alreadySelected = selectedFilterFamilyNotSelected.value;
  }
  return _.chain(items)
    .uniqBy('meta.family')
    .sort((a, b) => compareFamily(a.meta.family!, b.meta.family!))
    .map(
      (item: LicenseSlim) =>
        ({
          text: getI18NTextOfPrefixKey('LIC_FAMILY_', item.meta.family!),
          value: item.meta.family,
        }) as IDefaultSelectItem,
    )
    .union(
      _.map(
        alreadySelected,
        (selectedValue: string) =>
          ({
            text: getI18NTextOfPrefixKey('LIC_FAMILY_', selectedValue),
            value: selectedValue,
          }) as IDefaultSelectItem,
      ),
    )
    .value();
};

const getPossibleIsLicenseChart = (items: LicenseSlim[], selected = true): IDefaultSelectItem[] => {
  let alreadySelected = selectedFilterIsLicenseChartSelected.value;
  if (!selected) {
    alreadySelected = selectedFilterIsLicenseChartNotSelected.value;
  }

  return _.chain(items)
    .uniqBy('meta.isLicenseChart')
    .map(
      (item: LicenseSlim) =>
        ({
          text: t(item.meta.isLicenseChart ? 'TABLE_LICENSE_CHART_STATUS_IS' : 'TABLE_LICENSE_CHART_STATUS_IS_NOT'),
          value: item.meta.isLicenseChart + '',
        }) as IDefaultSelectItem,
    )
    .union(
      _.map(
        alreadySelected,
        (selectedValue: string) =>
          ({
            text: t(selectedValue === 'true' ? 'TABLE_LICENSE_CHART_STATUS_IS' : 'TABLE_LICENSE_CHART_STATUS_IS_NOT'),
            value: selectedValue,
          }) as IDefaultSelectItem,
      ),
    )
    .value();
};

const fillLicenseTbl = (policyState: PolicyState, license: LicenseSlim): boolean => {
  for (const licenseId of getComponents(policyState)) {
    if (licenseId === license.licenseId) {
      if (mode.value === policyState) {
        selectedLicenses.value.push(license);
      }
      return true;
    }
  }
  return false;
};

const reloadLabels = async () => {
  policyLabels.value = (await AdminService.getPolicyLabels()).data;
  createLabelsMap();
};

const createLabelsMap = () => {
  labelsMap.value = {};
  for (const lbl of policyLabels.value) {
    labelsMap.value[lbl._key] = lbl;
  }
};

const getComponents = (policyState: PolicyState) => {
  switch (policyState) {
    case PolicyState.ALLOW:
      return rule.value.ComponentsAllow;
    case PolicyState.DENY:
      return rule.value.ComponentsDeny;
    case PolicyState.WARN:
      return rule.value.ComponentsWarn;
    default:
      throw new Error('Unknown rule state: ' + mode.value);
  }
};

const getWarnLevel = (name: string) => {
  const classification = classifications.value.find((c) => c.name === name || c.nameDe === name);
  return classification ? classification.warnLevel : 'INFORMATION';
};

const moveAllFilteredToSelectedList = () => {
  filteredListNotSelected.value.forEach((license: LicenseSlim) => {
    const index = notSelectedLicenses.value.indexOf(license);
    if (index !== -1) {
      notSelectedLicenses.value.splice(index, 1);
      selectedLicenses.value.push(license);
      getComponents(mode.value).push(license.licenseId);
    }
  });
  refreshPossible(selectedLicenses.value);
  refreshPossible(notSelectedLicenses.value, false);
  hasChanges.value = true;
};

const unselectLicense = (license: LicenseSlim) => {
  if (!canEditManual.value) {
    return;
  }
  removeFromList(selectedLicenses.value, license);
  removeFromList(getComponents(mode.value), license.licenseId);
  notSelectedLicenses.value = [license].concat(notSelectedLicenses.value);
  hasChanges.value = true;

  refreshPossible(selectedLicenses.value);
  refreshPossible(notSelectedLicenses.value, false);
};

const filteredListSelected = computed(() =>
  _.chain(selectedLicenses.value)
    .filter(filterOnIsLicenseChartSelected)
    .filter(filterOnClassificationsSelected)
    .filter(filterOnFamilySelected)
    .filter(filterOnApprovalSelected)
    .filter(filterOnTypeSelected)
    .value(),
);

const filterOnIsLicenseChartSelected = (item: LicenseSlim): boolean =>
  selectedFilterIsLicenseChartSelected.value.length === 0 ||
  indexOf(selectedFilterIsLicenseChartSelected.value, item.meta.isLicenseChart + '') !== -1;

const filterOnClassificationsSelected = (item: LicenseSlim): boolean => {
  if (selectedFilterClassificationsSelected.value.length === 0) {
    return true;
  }
  if (!item.meta.classifications || item.meta.classifications.length === 0) {
    return selectedFilterClassificationsSelected.value.includes('');
  }
  return item.meta.classifications.some((classification) => {
    const classificationName = classification ? useViewTools().getNameForLanguage(classification) : '';
    return classificationName && selectedFilterClassificationsSelected.value.includes(classificationName);
  });
};

const filterOnFamilySelected = (item: LicenseSlim): boolean =>
  selectedFilterFamilySelected.value.length === 0 ||
  indexOf(selectedFilterFamilySelected.value, item.meta.family) !== -1;

const filterOnApprovalSelected = (item: LicenseSlim): boolean =>
  selectedFilterApprovalSelected.value.length === 0 ||
  indexOf(selectedFilterApprovalSelected.value, item.meta.approvalState) !== -1;

const filterOnTypeSelected = (item: LicenseSlim): boolean =>
  selectedFilterTypeSelected.value.length === 0 ||
  indexOf(selectedFilterTypeSelected.value, item.meta.licenseType) !== -1;

const openLicense = (item: LicenseSlim) => {
  openUrlInNewTab(`/dashboard/licenses/${encodeURIComponent(item.licenseId)}`);
};

const selectLicense = (license: LicenseSlim) => {
  if (!canEditManual.value) {
    return;
  }
  const index = notSelectedLicenses.value.indexOf(license, 0);
  if (index > -1) {
    notSelectedLicenses.value.splice(index, 1);
  }
  selectedLicenses.value.push(license);
  getComponents(mode.value).push(license.licenseId);
  hasChanges.value = true;

  refreshPossible(selectedLicenses.value);
  refreshPossible(notSelectedLicenses.value, false);
};

const openClassifications = (licenseClassifications: IObligation[], licenseName: string, licenseId: string) => {
  if (classificationsDialogRef.value) {
    classificationsDialogRef.value?.open(licenseClassifications, licenseName, licenseId);
  }
};

const saveChanges = async () => {
  rule.value = (await AdminService.putPolicyRule(rule.value)).data;
  await retrieveSpdxLicenses();
  info(t('DESCRIPTION_POLICY_RULE_SAVED'));
  hasChanges.value = false;
};

const policies = ref(PolicyRules);

const ruleCallback: IRuleBtnCallbacks = {
  getUrlToComponents: (policy) => {
    return '';
  },
  handlePolicySelect: (policy, selectedFilterPolicyTypes) => {
    // change mode
    mode.value = policy;
    // update tables
    fillTables(false);
  },
  getCountForPolicyFilterBtn: (policy) => {
    switch (policy) {
      case PolicyState.ALLOW:
        return rule.value.ComponentsAllow.length;
      case PolicyState.DENY:
        return rule.value.ComponentsDeny.length;
      case PolicyState.WARN:
        return rule.value.ComponentsWarn.length;
      default:
        throw new Error('Method not implemented.');
    }
  },
  getInitSelectedPolicy: () => {
    return PolicyState.ALLOW;
  },
  getToolTipKeyForPolicyFilterBtn: (policy) => {
    switch (policy) {
      case PolicyState.DENY:
        return 'PR_COMPONENTS_DENIED';
      case PolicyState.WARN:
        return 'PR_COMPONENTS_WARNED';
      case PolicyState.ALLOW:
        return 'PR_COMPONENTS_ALLOWED';
      default:
        return 'unknown_policy';
    }
  },
  getActiveClassForPolicyFilterBtn: (policy) => {
    switch (policy) {
      case PolicyState.DENY:
        return 'deny-border';
      case PolicyState.WARN:
        return 'warning-border';
      case PolicyState.ALLOW:
        return 'allow-border';
      default:
        return '';
    }
  },
  setRuleButtons: () => {},
};

const componentHeaders = computed<DataTableHeader[]>(() => {
  const headers: DataTableHeader[] = [
    {
      title: t('COL_LICENSE_CHART_STATUS'),
      tooltipText: t('TABLE_LICENSE_CHART_STATUS_TOOLTIP'),
      align: 'center',
      value: 'meta.isLicenseChart',
      width: 100,
      class: 'licenseChartHeader tableHeaderCell',
      filterable: true,
      sortable: false,
    },
    {
      title: t('CLASSIFICATION'),
      tooltipText: t('LC_CLASSIFICATION_TT'),
      class: 'tableHeaderCell',
      filterable: true,
      width: 150,
      value: 'meta.classifications',
      sort: (a: IObligation[], b: IObligation[]) => {
        const getHighestWarnLevel = (obligations: IObligation[] = []) => {
          if (!obligations || obligations.length === 0) {
            return 'INFORMATION';
          }
          return obligations.reduce((highestLevel, obligation) => {
            const currentLevel = obligation.warnLevel?.toUpperCase() || 'INFORMATION';
            const currentWeight = levelWeight.get(currentLevel) ?? 0;
            const highestWeight = levelWeight.get(highestLevel) ?? 0;

            return currentWeight > highestWeight ? currentLevel : highestLevel;
          }, 'INFORMATION');
        };

        const levelA = getHighestWarnLevel(a);
        const levelB = getHighestWarnLevel(b);

        return compareLevel(levelA, levelB);
      },
    },
    {
      title: t('COL_NAME'),
      tooltipText: t('COL_LICENSE_NAME_TOOLTIP'),
      align: 'start',
      filterable: true,
      width: 400,
      class: 'tableHeaderCell',
      value: 'name',
      sortable: true,
    },
    {
      title: t('COL_LICENSE_ID'),
      tooltipText: t('COL_LICENSE_ID_TOOLTIP'),
      align: 'start',
      filterable: true,
      width: 400,
      class: 'tableHeaderCell',
      value: 'licenseId',
      sortable: true,
    },
    {
      title: t('COL_LICENSE_FAMILY'),
      tooltipText: t('COL_LICENSE_FAMILY_TOOLTIP'),
      align: 'start',
      filterable: true,
      width: 200,
      class: 'tableHeaderCell',
      value: 'meta.family',
      sortable: true,
    },
    {
      title: t('COL_TYPE'),
      tooltipText: t('COL_LICENSE_TYPE_TOOLTIP'),
      class: 'tableHeaderCell',
      filterable: true,
      width: 150,
      value: 'meta.licenseType',
      sortable: true,
    },
    {
      title: t('COL_APPROVAL_STATUS'),
      tooltipText: t('COL_APPROVAL_STATUS_TOOLTIP'),
      align: 'start',
      filterable: true,
      width: 170,
      class: 'tableHeaderCell',
      value: 'meta.approvalState',
      sortable: true,
    },
  ];

  if (canEditManual.value) {
    headers.push({
      title: '',
      align: 'start',
      filterable: false,
      sortable: false,
      width: 50,
      class: 'tableHeaderCell',
      value: 'remove',
    });
  }

  if (!canEditManual.value) {
    headers.push({
      title: t('COL_ACTIONS'),
      sortable: false,
      align: 'center',
      width: 100,
      class: 'tableHeaderCell',
      value: 'actions',
    });
  }

  return headers as DataTableHeader[];
});

const componentHeadersSelected = computed<DataTableHeader[]>(() => componentHeaders.value);
const componentHeadersUnSelected = computed<DataTableHeader[]>(() => [
  {
    title: '',
    align: 'start',
    filterable: false,
    width: 50,
    class: 'tableHeaderCell',
    value: 'add',
  },
  ...(componentHeaders.value as DataTableHeader[]),
]);

onMounted(async () => {
  policies.value = PolicyRules;

  rights.value = userStore.getRights;
  isPolicyManager.value = RightsUtils.isPolicyManager();
  ruleId.value = <string>router.currentRoute.value.params.id;
  await retrieveRule(ruleId.value);
  initBreadcrumbs();

  await retrieveSpdxLicenses();
  await calculatedPolicyRuleStore.retrieveClassifications();
  await reloadLabels();
});

const handleSetCalculatedEnabled = (value: boolean) => {
  calculatedPolicyRuleStore.setCalculated(value);
  hasChanges.value = true;
};
</script>

<template>
  <TableLayout has-tab has-title>
    <template #buttons>
      <div class="flex w-full flex-col gap-4">
        <div class="grid w-full basis-full gap-6" :class="{'grid-cols-2': canEditManual || canEditCalculated}">
          <div v-if="isPolicyManager" class="d-flex ga-2 align-center mt-2 h-9 flex-row">
            <h3 class="d-subtitle-2">{{ t('TABLE_HEADER_LICENSES') }}</h3>
            <DCActionButton
              :text="t('BTN_SAVE')"
              icon="mdi-content-save"
              :hint="t('BTN_SAVE')"
              @click="saveChanges"
              v-if="hasChanges && rule.Deprecated === false && !rule.Calculated" />
          </div>
          <div v-if="isPolicyManager" class="d-flex align-center justify-space-between mt-2 h-9 flex-row">
            <h3 v-if="!rule.Calculated && ruleLoaded" class="d-subtitle-2">
              {{ t('TABLE_HEADER_AVAILABLE_LICENSES') }}
            </h3>
            <div v-else></div>
            <template v-if="ruleLoaded">
              <DCActionButton
                v-if="!rule.Calculated"
                variant="outlined"
                :text="t('CALCULATED_POLICY_RULE_ENABLED')"
                icon="mdi-calculator-variant"
                :hint="t('CALCULATED_POLICY_RULE_ENABLED')"
                @click="handleSetCalculatedEnabled(true)" />
              <DCActionButton
                v-else
                variant="outlined"
                :text="t('MANUAL_RULES')"
                icon="mdi-cog-outline"
                :hint="t('MANUAL_RULES')"
                @click="handleSetCalculatedEnabled(false)" />
            </template>
          </div>

          <div :class="{'col-span-2': !canEditManual && !rule.Calculated, 'col-start-1': rule.Calculated}">
            <div class="d-flex ga-1 label-filter flex-row">
              <div class="overflow-auto">
                <DRuleButtons :policies="policies" :callbacks="ruleCallback" min-width="128px" :forceClickable="true" />
              </div>
              <v-spacer />
              <DSearchField v-model="filterSelected" />
            </div>
          </div>
          <div v-if="canEditManual && ruleLoaded">
            <div class="d-flex ga-1 label-filter flex-row">
              <DCActionButton
                large
                variant="outlined"
                :text="`${t('MOVE_TO_SELECTED')} (${filteredListNotSelected.length})`"
                icon="mdi-chevron-left"
                :hint="t('TT_MOVE_TO_SELECTED')"
                @click="moveAllFilteredToSelectedList"
                v-if="filteredListNotSelected.length > 0" />
              <v-spacer />
              <DSearchField v-model="filterUnSelected" />
            </div>
          </div>
        </div>
      </div>
    </template>
    <template #table>
      <v-row class="fill-height">
        <v-col :cols="canEditManual || canEditCalculated ? 6 : 12" class="fill-height">
          <div class="fill-height" :class="getActiveClassForPolicyFilterBtn(mode)">
            <v-data-table
              :loading="licensesLoading"
              fixed-header
              :headers="componentHeadersSelected"
              :class="getCssClass() + ' striped-table fill-height'"
              :search="filterSelected"
              :items-per-page="25"
              :items="filteredListSelected"
              @[canEditManual&&`click:row`]="
                (event: Event, dataItem: DataTableItem<LicenseSlim>) => unselectLicense(dataItem.item)
              "
              density="compact">
              <template v-slot:header.meta.isLicenseChart="{column}">
                <div class="v-data-table-header__content">
                  <span>{{ column.title }}</span>
                  <v-menu :close-on-content-click="false" v-model="menuIsLicenseChartSelected">
                    <template v-slot:activator="{props}">
                      <DIconButton
                        :parentProps="props"
                        icon="mdi-filter-variant"
                        :hint="t('TT_SHOW_FILTER')"
                        :color="selectedFilterIsLicenseChartSelected.length > 0 ? 'primary' : 'default'" />
                    </template>
                    <div style="width: 280px" class="bg-background">
                      <v-row class="d-flex ma-1 mr-2 justify-end">
                        <DCloseButton @click="menuIsLicenseChartSelected = false" />
                      </v-row>
                      <v-select
                        v-model="selectedFilterIsLicenseChartSelected"
                        class="pa-2 mx-2"
                        density="compact"
                        @focus="onACFocus"
                        variant="outlined"
                        autofocus
                        :items="possibleIsLicenseChartSelected"
                        :label="t('Lbl_filter_License_Chart_Status')"
                        hide-details
                        color="inputActiveBorderColor"
                        multiple
                        location="bottom"
                        item-title="text"
                        item-value="value"
                        menu
                        clearable
                        transition="scale-transition"
                        persistent-clear
                        :list-props="{class: 'striped-filter-dd py-0'}">
                        <template v-slot:selection="{item, index}">
                          <span v-if="index === 0" class="pFilterEntry">{{ item.title }}</span>
                          <span v-if="index === 1" class="pAdditionalFilter">
                            +{{ selectedFilterIsLicenseChartSelected.length - 1 }} others
                          </span>
                        </template>
                        <template v-slot:item="{item, props}">
                          <v-list-item v-bind="props" class="px-2 py-0">
                            <template v-slot:prepend="{isSelected}">
                              <v-checkbox :model-value="isSelected" hide-details></v-checkbox>
                            </template>
                            <template v-slot:title="{}">
                              <span class="pFilterEntry">{{ item.props.title }}</span>
                            </template>
                          </v-list-item>
                        </template>
                      </v-select>
                    </div>
                  </v-menu>
                </div>
              </template>
              <template v-slot:header.meta.licenseType="{column, getSortIcon, toggleSort}">
                <div class="v-data-table-header__content">
                  <span>{{ column.title }}</span>
                  <v-menu :close-on-content-click="false" v-model="menu5">
                    <template v-slot:activator="{props}">
                      <DIconButton
                        :parentProps="props"
                        icon="mdi-filter-variant"
                        :hint="t('TT_SHOW_FILTER')"
                        :color="selectedFilterTypeSelected.length > 0 ? 'primary' : 'default'" />
                    </template>
                    <div style="width: 280px" class="bg-background">
                      <v-row class="d-flex ma-1 mr-2 justify-end">
                        <DCloseButton @click="menu5 = false" />
                      </v-row>
                      <v-select
                        v-model="selectedFilterTypeSelected"
                        class="pa-2 mx-2"
                        density="compact"
                        clearable
                        @focus="onACFocus"
                        variant="outlined"
                        autofocus
                        :items="possibleTypeSelected"
                        :label="t('Lbl_filter_type')"
                        hide-details
                        color="inputActiveBorderColor"
                        multiple
                        v-bind:menu-props="{location: 'bottom'}"
                        item-title="text"
                        item-value="value"
                        menu
                        transition="scale-transition"
                        persistent-clear
                        :list-props="{class: 'striped-filter-dd py-0'}">
                        <template v-slot:selection="{item, index}">
                          <span v-if="index === 0" class="pFilterEntry">{{ item.title }}</span>
                          <span v-if="index === 1" class="pAdditionalFilter">
                            +{{ selectedFilterTypeSelected.length - 1 }} others
                          </span>
                        </template>
                        <template v-slot:item="{item, props}">
                          <v-list-item v-bind="props" class="px-2 py-0">
                            <template v-slot:prepend="{isSelected}">
                              <v-checkbox :model-value="isSelected" hide-details></v-checkbox>
                            </template>
                            <template v-slot:title="{}">
                              <span class="pFilterEntry">{{ item.props.title }}</span>
                            </template>
                          </v-list-item>
                        </template>
                      </v-select>
                    </div>
                  </v-menu>
                  <v-icon
                    class="v-data-table-header__sort-icon"
                    :icon="getSortIcon(column)"
                    @click="toggleSort(column)" />
                </div>
              </template>
              <template v-slot:header.meta.approvalState="{column, getSortIcon, toggleSort}">
                <div class="v-data-table-header__content">
                  <span>{{ column.title }}</span>
                  <v-menu :close-on-content-click="false" v-model="menu4">
                    <template v-slot:activator="{props}">
                      <DIconButton
                        :parentProps="props"
                        icon="mdi-filter-variant"
                        :hint="t('TT_SHOW_FILTER')"
                        :color="selectedFilterApprovalSelected.length > 0 ? 'primary' : 'default'" />
                    </template>
                    <div style="width: 280px" class="bg-background">
                      <v-row class="d-flex ma-1 mr-2 justify-end">
                        <DCloseButton @click="menu4 = false" />
                      </v-row>
                      <v-select
                        v-model="selectedFilterApprovalSelected"
                        class="pa-2 mx-2"
                        density="compact"
                        clearable
                        @focus="onACFocus"
                        variant="outlined"
                        autofocus
                        :items="possibleApprovalSelected"
                        :label="t('Lbl_filter_approval')"
                        hide-details
                        color="inputActiveBorderColor"
                        multiple
                        v-bind:menu-props="{location: 'bottom'}"
                        item-title="text"
                        item-value="value"
                        menu
                        transition="scale-transition"
                        persistent-clear
                        :list-props="{class: 'striped-filter-dd py-0'}">
                        <template v-slot:selection="{item, index}">
                          <span v-if="index === 0" class="pFilterEntry">{{ item.title }}</span>
                          <span v-if="index === 1" class="pAdditionalFilter">
                            +{{ selectedFilterApprovalSelected.length - 1 }} others
                          </span>
                        </template>
                        <template v-slot:item="{item, props}">
                          <v-list-item v-bind="props" class="px-2 py-0">
                            <template v-slot:prepend="{isSelected}">
                              <v-checkbox :model-value="isSelected" hide-details></v-checkbox>
                            </template>
                            <template v-slot:title="{}">
                              <span class="pFilterEntry">{{ item.props.title }}</span>
                            </template>
                          </v-list-item>
                        </template>
                      </v-select>
                    </div>
                  </v-menu>
                  <v-icon
                    class="v-data-table-header__sort-icon"
                    :icon="getSortIcon(column)"
                    @click="toggleSort(column)" />
                </div>
              </template>
              <template v-slot:header.meta.family="{column, getSortIcon, toggleSort}">
                <div class="v-data-table-header__content">
                  <span>{{ column.title }}</span>
                  <v-menu :close-on-content-click="false" v-model="menu">
                    <template v-slot:activator="{props}">
                      <DIconButton
                        :parentProps="props"
                        icon="mdi-filter-variant"
                        :hint="t('TT_SHOW_FILTER')"
                        :color="selectedFilterFamilySelected.length > 0 ? 'primary' : 'default'" />
                    </template>
                    <div style="width: 280px" class="bg-background">
                      <v-row class="d-flex ma-1 mr-2 justify-end">
                        <DCloseButton @click="menu = false" />
                      </v-row>
                      <v-select
                        v-model="selectedFilterFamilySelected"
                        class="pa-2 mx-2"
                        density="compact"
                        clearable
                        @focus="onACFocus"
                        variant="outlined"
                        autofocus
                        :items="possibleFamilySelected"
                        :label="t('Lbl_filter_family')"
                        hide-details
                        color="inputActiveBorderColor"
                        multiple
                        v-bind:menu-props="{location: 'bottom'}"
                        item-title="text"
                        item-value="value"
                        menu
                        transition="scale-transition"
                        persistent-clear
                        :list-props="{class: 'striped-filter-dd py-0'}">
                        <template v-slot:selection="{item, index}">
                          <span v-if="index === 0" class="pFilterEntry">{{ item.title }}</span>
                          <span v-if="index === 1" class="pAdditionalFilter">
                            +{{ selectedFilterFamilySelected.length - 1 }} others
                          </span>
                        </template>
                        <template v-slot:item="{item, props}">
                          <v-list-item v-bind="props" class="px-2 py-0">
                            <template v-slot:prepend="{isSelected}">
                              <v-checkbox :model-value="isSelected" hide-details></v-checkbox>
                            </template>
                            <template v-slot:title="{}">
                              <span class="pFilterEntry">{{ item.props.title }}</span>
                            </template>
                          </v-list-item>
                        </template>
                      </v-select>
                    </div>
                  </v-menu>
                  <v-icon
                    class="v-data-table-header__sort-icon"
                    :icon="getSortIcon(column)"
                    @click="toggleSort(column)" />
                </div>
              </template>
              <template v-slot:header.meta.classifications="{column, getSortIcon, toggleSort}">
                <div class="v-data-table-header__content">
                  <span>{{ column.title }}</span>
                  <v-menu :close-on-content-click="false" v-model="menuClassification">
                    <template v-slot:activator="{props}">
                      <DIconButton
                        :parentProps="props"
                        icon="mdi-filter-variant"
                        :hint="t('TT_SHOW_FILTER')"
                        :color="selectedFilterFamilySelected.length > 0 ? 'primary' : 'default'" />
                    </template>
                    <div style="width: 520px" class="bg-background">
                      <v-row class="d-flex ma-1 mr-2 justify-end">
                        <DCloseButton @click="menuClassification = false" />
                      </v-row>
                      <v-select
                        v-model="selectedFilterClassificationsSelected"
                        class="pa-2 mx-2"
                        density="compact"
                        clearable
                        @focus="onACFocus"
                        variant="outlined"
                        autofocus
                        :items="possibleClassificationsSelected"
                        :label="t('CLASSIFICATION')"
                        hide-details
                        color="inputActiveBorderColor"
                        multiple
                        v-bind:menu-props="{location: 'bottom'}"
                        item-title="text"
                        item-value="value"
                        menu
                        transition="scale-transition"
                        persistent-clear
                        :list-props="{class: 'striped-filter-dd py-0'}">
                        <template v-slot:selection="{item, index}">
                          <span v-if="index === 0" class="pFilterEntry">{{ item.title }}</span>
                          <span v-if="index === 1" class="pAddtionalFilter">
                            +{{ selectedFilterClassificationsSelected.length - 1 }} others
                          </span>
                        </template>
                        <template v-slot:item="{item, props}">
                          <v-list-item v-bind="props">
                            <template v-slot:prepend="{isSelected}">
                              <v-checkbox :model-value="isSelected" hide-details></v-checkbox>
                            </template>
                            <template v-slot:title="{}">
                              <v-icon
                                size="small"
                                :color="getIconColorOfLevel(getWarnLevel(item.value))"
                                class="mr-2"
                                :icon="getIconOfLevel(getWarnLevel(item.value).toUpperCase())">
                              </v-icon>
                              <span class="pFilterEntry">{{ item.props.title }}</span>
                            </template>
                          </v-list-item>
                        </template>
                      </v-select>
                    </div>
                  </v-menu>
                  <v-icon
                    class="v-data-table-header__sort-icon"
                    :icon="getSortIcon(column)"
                    @click="toggleSort(column)" />
                </div>
              </template>
              <template v-slot:item.remove>
                <v-icon color="primary" icon="mdi-close"></v-icon>
              </template>
              <template v-slot:item.meta.isLicenseChart="{item}">
                <DLicenseChartIcon :meta="item.meta" />
              </template>
              <template v-slot:item.meta.approvalState="{item}">
                {{ getI18NTextOfPrefixKey('LT_APP_', item.meta.approvalState) }}
              </template>
              <template v-slot:item.meta.licenseType="{item}">
                {{ getI18NTextOfPrefixKey('LT_', item.meta.licenseType) }}
              </template>
              <template v-slot:item.meta.family="{item}">
                {{ getI18NTextOfPrefixKey('LIC_FAMILY_', item.meta.family!) }}
              </template>

              <template v-slot:item.meta.classifications="{item}">
                <span @click.stop="openClassifications(item.meta.classifications, item.name, item.licenseId)">
                  <v-tooltip :open-delay="TOOLTIP_OPEN_DELAY_IN_MS" location="bottom">
                    <template v-slot:activator="{props}">
                      <v-icon
                        v-bind="props"
                        color="primary"
                        size="small"
                        icon="mdi-chevron-right"
                        :class="
                          item.meta.prevalentClassificationLevel.toUpperCase() === 'WARNING' ? 'mr-1' : 'mr-2'
                        "></v-icon>
                      <v-icon
                        v-bind="props"
                        density="compact"
                        style="font-size: 20px"
                        :icon="getIconOfLevel(item.meta.prevalentClassificationLevel)"
                        :color="getIconColorOfLevel(item.meta.prevalentClassificationLevel)"></v-icon>
                    </template>
                    <span>{{ t('TT_OPEN_CLASSIFICATIONS', {license: item.name}) }}</span>
                  </v-tooltip>
                </span>
              </template>
              <template v-slot:item.actions="{item}">
                <DIconButton
                  v-if="rights.isInternal"
                  icon="mdi-chevron-right"
                  :hint="t('TT_open_license')"
                  @clicked="openLicense(item)" />
              </template>
            </v-data-table>
          </div>
        </v-col>
        <v-col cols="6" v-if="ruleLoaded && canEditManual" class="fill-height">
          <v-data-table
            :loading="licensesLoading"
            fixed-header
            :headers="componentHeadersUnSelected"
            :items-per-page="25"
            :items="filteredListNotSelected"
            :search="filterUnSelected"
            @click:row="(event: Event, dataItem: DataTableItem<LicenseSlim>) => selectLicense(dataItem.item)"
            density="compact"
            class="striped-table fill-height">
            <template v-slot:header.meta.isLicenseChart="{column}">
              <div class="v-data-table-header__content">
                <span>{{ column.title }}</span>
                <v-menu :close-on-content-click="false" v-model="menuIsLicenseChartNotSelected">
                  <template v-slot:activator="{props}">
                    <DIconButton
                      :parentProps="props"
                      icon="mdi-filter-variant"
                      :hint="t('TT_SHOW_FILTER')"
                      :color="selectedFilterIsLicenseChartNotSelected.length > 0 ? 'primary' : 'default'" />
                  </template>
                  <div style="width: 280px" class="bg-background">
                    <v-row class="d-flex ma-1 mr-2 justify-end">
                      <DCloseButton @click="menuIsLicenseChartNotSelected = false" />
                    </v-row>
                    <v-select
                      v-model="selectedFilterIsLicenseChartNotSelected"
                      class="pa-2 mx-2"
                      density="compact"
                      @focus="onACFocus"
                      variant="outlined"
                      autofocus
                      :items="possibleIsLicenseChartNotSelected"
                      :label="t('Lbl_filter_License_Chart_Status')"
                      hide-details
                      color="inputActiveBorderColor"
                      multiple
                      location="bottom"
                      item-title="text"
                      item-value="value"
                      menu
                      clearable
                      transition="scale-transition"
                      persistent-clear
                      :list-props="{class: 'striped-filter-dd py-0'}">
                      <template v-slot:selection="{item, index}">
                        <span v-if="index === 0" class="pFilterEntry">{{ item.title }}</span>
                        <span v-if="index === 1" class="pAdditionalFilter">
                          +{{ selectedFilterIsLicenseChartNotSelected.length - 1 }} others
                        </span>
                      </template>
                      <template v-slot:item="{item, props}">
                        <v-list-item v-bind="props" class="px-2 py-0">
                          <template v-slot:prepend="{isSelected}">
                            <v-checkbox :model-value="isSelected" hide-details></v-checkbox>
                          </template>
                          <template v-slot:title="{}">
                            <span class="pFilterEntry">{{ item.props.title }}</span>
                          </template>
                        </v-list-item>
                      </template>
                    </v-select>
                  </div>
                </v-menu>
              </div>
            </template>
            <template v-slot:header.meta.licenseType="{column, getSortIcon, toggleSort}">
              <div class="v-data-table-header__content">
                <span>{{ column.title }}</span>
                <v-menu :close-on-content-click="false" v-model="menu6">
                  <template v-slot:activator="{props}">
                    <DIconButton
                      :parentProps="props"
                      icon="mdi-filter-variant"
                      :hint="t('TT_SHOW_FILTER')"
                      :color="selectedFilterTypeNotSelected.length > 0 ? 'primary' : 'default'" />
                  </template>
                  <div style="width: 280px" class="bg-background">
                    <v-row class="d-flex ma-1 mr-2 justify-end">
                      <DCloseButton @click="menu6 = false" />
                    </v-row>
                    <v-select
                      v-model="selectedFilterTypeNotSelected"
                      class="pa-2 mx-2"
                      density="compact"
                      clearable
                      @focus="onACFocus"
                      variant="outlined"
                      autofocus
                      :items="possibleTypeNotSelected"
                      :label="t('Lbl_filter_type')"
                      hide-details
                      color="inputActiveBorderColor"
                      multiple
                      v-bind:menu-props="{location: 'bottom'}"
                      item-title="text"
                      item-value="value"
                      menu
                      transition="scale-transition"
                      persistent-clear
                      :list-props="{class: 'striped-filter-dd py-0'}">
                      <template v-slot:selection="{item, index}">
                        <span v-if="index === 0" class="pFilterEntry">{{ item.title }}</span>
                        <span v-if="index === 1" class="pAdditionalFilter">
                          +{{ selectedFilterTypeNotSelected.length - 1 }} others
                        </span>
                      </template>
                      <template v-slot:item="{item, props}">
                        <v-list-item v-bind="props" class="px-2 py-0">
                          <template v-slot:prepend="{isSelected}">
                            <v-checkbox :model-value="isSelected" hide-details></v-checkbox>
                          </template>
                          <template v-slot:title="{}">
                            <span class="pFilterEntry">{{ item.props.title }}</span>
                          </template>
                        </v-list-item>
                      </template>
                    </v-select>
                  </div>
                </v-menu>
                <v-icon
                  class="v-data-table-header__sort-icon"
                  :icon="getSortIcon(column)"
                  @click="toggleSort(column)" />
              </div>
            </template>
            <template v-slot:header.meta.approvalState="{column, getSortIcon, toggleSort}">
              <div class="v-data-table-header__content">
                <span>{{ column.title }}</span>
                <v-menu :close-on-content-click="false" v-model="menu3">
                  <template v-slot:activator="{props}">
                    <DIconButton
                      :parentProps="props"
                      icon="mdi-filter-variant"
                      :hint="t('TT_SHOW_FILTER')"
                      :color="selectedFilterApprovalNotSelected.length > 0 ? 'primary' : 'default'" />
                  </template>
                  <div style="width: 280px" class="bg-background">
                    <v-row class="d-flex ma-1 mr-2 justify-end">
                      <DCloseButton @click="menu3 = false" />
                    </v-row>
                    <v-select
                      v-model="selectedFilterApprovalNotSelected"
                      class="pa-2 mx-2"
                      density="compact"
                      clearable
                      @focus="onACFocus"
                      variant="outlined"
                      autofocus
                      :items="possibleApprovalNotSelected"
                      :label="t('Lbl_filter_approval')"
                      hide-details
                      color="inputActiveBorderColor"
                      multiple
                      v-bind:menu-props="{location: 'bottom'}"
                      item-title="text"
                      item-value="value"
                      menu
                      transition="scale-transition"
                      persistent-clear
                      :list-props="{class: 'striped-filter-dd py-0'}">
                      <template v-slot:selection="{item, index}">
                        <span v-if="index === 0" class="pFilterEntry">{{ item.title }}</span>
                        <span v-if="index === 1" class="pAdditionalFilter">
                          +{{ selectedFilterApprovalNotSelected.length - 1 }} others
                        </span>
                      </template>
                      <template v-slot:item="{item, props}">
                        <v-list-item v-bind="props" class="px-2 py-0">
                          <template v-slot:prepend="{isSelected}">
                            <v-checkbox :model-value="isSelected" hide-details></v-checkbox>
                          </template>
                          <template v-slot:title="{}">
                            <span class="pFilterEntry">{{ item.props.title }}</span>
                          </template>
                        </v-list-item>
                      </template>
                    </v-select>
                  </div>
                </v-menu>
                <v-icon
                  class="v-data-table-header__sort-icon"
                  :icon="getSortIcon(column)"
                  @click="toggleSort(column)" />
              </div>
            </template>
            <template v-slot:header.meta.family="{column, getSortIcon, toggleSort}">
              <div class="v-data-table-header__content">
                <span>{{ column.title }}</span>
                <v-menu :close-on-content-click="false" v-model="menu2">
                  <template v-slot:activator="{props}">
                    <DIconButton
                      :parentProps="props"
                      icon="mdi-filter-variant"
                      :hint="t('TT_SHOW_FILTER')"
                      :color="selectedFilterFamilyNotSelected.length > 0 ? 'primary' : 'default'" />
                  </template>
                  <div style="width: 280px" class="bg-background">
                    <v-row class="d-flex ma-1 mr-2 justify-end">
                      <DCloseButton @click="menu2 = false" />
                    </v-row>
                    <v-select
                      v-model="selectedFilterFamilyNotSelected"
                      class="pa-2 mx-2"
                      density="compact"
                      clearable
                      @focus="onACFocus"
                      variant="outlined"
                      autofocus
                      :items="possibleFamilyNotSelected"
                      :label="t('Lbl_filter_family')"
                      hide-details
                      color="inputActiveBorderColor"
                      multiple
                      v-bind:menu-props="{location: 'bottom'}"
                      item-title="text"
                      item-value="value"
                      menu
                      transition="scale-transition"
                      persistent-clear
                      :list-props="{class: 'striped-filter-dd py-0'}">
                      <template v-slot:selection="{item, index}">
                        <span v-if="index === 0" class="pFilterEntry">{{ item.title }}</span>
                        <span v-if="index === 1" class="pAdditionalFilter">
                          +{{ selectedFilterFamilyNotSelected.length - 1 }} others
                        </span>
                      </template>
                      <template v-slot:item="{item, props}">
                        <v-list-item v-bind="props" class="px-2 py-0">
                          <template v-slot:prepend="{isSelected}">
                            <v-checkbox :model-value="isSelected" hide-details></v-checkbox>
                          </template>
                          <template v-slot:title="{}">
                            <span class="pFilterEntry">{{ item.props.title }}</span>
                          </template>
                        </v-list-item>
                      </template>
                    </v-select>
                  </div>
                </v-menu>
                <v-icon
                  class="v-data-table-header__sort-icon"
                  :icon="getSortIcon(column)"
                  @click="toggleSort(column)" />
              </div>
            </template>
            <template v-slot:header.meta.classifications="{column, getSortIcon, toggleSort}">
              <div class="v-data-table-header__content">
                <span>{{ column.title }}</span>
                <v-menu :close-on-content-click="false" v-model="menuClassificationNot">
                  <template v-slot:activator="{props}">
                    <DIconButton
                      :parentProps="props"
                      icon="mdi-filter-variant"
                      :hint="t('TT_SHOW_FILTER')"
                      :color="selectedFilterClassificationsNotSelected.length > 0 ? 'primary' : 'default'" />
                  </template>
                  <div style="width: 520px" class="bg-background">
                    <v-row class="d-flex ma-1 mr-2 justify-end">
                      <DCloseButton @click="menuClassificationNot = false" />
                    </v-row>
                    <v-select
                      v-model="selectedFilterClassificationsNotSelected"
                      class="pa-2 mx-2"
                      density="compact"
                      clearable
                      @focus="onACFocus"
                      variant="outlined"
                      autofocus
                      :items="possibleClassificationsNotSelected"
                      :label="t('CLASSIFICATION')"
                      hide-details
                      color="inputActiveBorderColor"
                      multiple
                      v-bind:menu-props="{location: 'bottom'}"
                      item-title="text"
                      item-value="value"
                      menu
                      transition="scale-transition"
                      persistent-clear
                      :list-props="{class: 'striped-filter-dd py-0'}">
                      <template v-slot:selection="{item, index}">
                        <span v-if="index === 0" class="pFilterEntry">{{ item.title }}</span>
                        <span v-if="index === 1" class="pAddtionalFilter">
                          +{{ selectedFilterClassificationsNotSelected.length - 1 }} others
                        </span>
                      </template>
                      <template v-slot:item="{item, props}">
                        <v-list-item v-bind="props">
                          <template v-slot:prepend="{isSelected}">
                            <v-checkbox :model-value="isSelected" hide-details></v-checkbox>
                          </template>
                          <template v-slot:title="{}">
                            <v-icon
                              size="small"
                              :color="getIconColorOfLevel(getWarnLevel(item.value))"
                              class="mr-2"
                              :icon="getIconOfLevel(getWarnLevel(item.value).toUpperCase())">
                            </v-icon>
                            <span class="pFilterEntry">{{ item.props.title }}</span>
                          </template>
                        </v-list-item>
                      </template>
                    </v-select>
                  </div>
                </v-menu>
                <v-icon
                  class="v-data-table-header__sort-icon"
                  :icon="getSortIcon(column)"
                  @click="toggleSort(column)" />
              </div>
            </template>
            <template v-slot:item.add>
              <span style="float: left" v-if="canEditManual">
                <v-icon color="primary" icon="mdi-chevron-left"></v-icon>
              </span>
            </template>
            <template v-slot:item.meta.isLicenseChart="{item}">
              <DLicenseChartIcon :meta="item.meta" />
            </template>
            <template v-slot:item.meta.approvalState="{item}">
              {{ getI18NTextOfPrefixKey('LT_APP_', item.meta.approvalState) }}
            </template>
            <template v-slot:item.meta.licenseType="{item}">
              {{ getI18NTextOfPrefixKey('LT_', item.meta.licenseType) }}
            </template>
            <template v-slot:item.meta.family="{item}">
              {{ getI18NTextOfPrefixKey('LIC_FAMILY_', item.meta.family!) }}
            </template>
            <template v-slot:item.meta.classifications="{item}">
              <span @click.stop="openClassifications(item.meta.classifications, item.name, item.licenseId)">
                <v-tooltip :open-delay="TOOLTIP_OPEN_DELAY_IN_MS" location="bottom">
                  <template v-slot:activator="{props}">
                    <v-icon
                      v-bind="props"
                      color="primary"
                      small
                      :class="item.meta.prevalentClassificationLevel.toUpperCase() === 'WARNING' ? 'mr-1' : 'mr-2'"
                      >mdi-chevron-right</v-icon
                    >
                    <v-icon
                      v-bind="props"
                      style="font-size: 20px"
                      :color="getIconColorOfLevel(item.meta.prevalentClassificationLevel)"
                      >{{ getIconOfLevel(item.meta.prevalentClassificationLevel) }}</v-icon
                    >
                  </template>
                  <span>{{ t('TT_OPEN_CLASSIFICATIONS', {license: item.name}) }}</span>
                </v-tooltip>
              </span>
            </template>
          </v-data-table>
        </v-col>
        <v-col cols="6" v-if="canEditCalculated" class="fill-height">
          <div class="flex h-full flex-col">
            <div v-show="classificationsLoaded" class="flex-1 overflow-auto">
              <CalculatedRuleConfig />
              <div class="d-flex mt-8 justify-end">
                <DCActionButton
                  :text="t('BTN_SAVE')"
                  icon="mdi-content-save"
                  :hint="t('BTN_SAVE')"
                  @click="saveChanges"
                  v-if="hasChanges && rule.Deprecated === false" />
              </div>
            </div>
          </div>
        </v-col>
      </v-row>
    </template>
  </TableLayout>
  <ClassificationsPerLicenseDialog ref="classificationsDialogRef"></ClassificationsPerLicenseDialog>
</template>
<style scoped>
.label-filter {
  @media (width < 1450px) {
    flex-direction: column-reverse !important;
  }
}
</style>
