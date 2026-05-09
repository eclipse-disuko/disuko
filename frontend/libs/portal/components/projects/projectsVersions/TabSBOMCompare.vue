<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import ComponentCompareDialog from '@disclosure-portal/components/dialog/ComponentCompareDialog.vue';
import Icons from '@disclosure-portal/constants/icons';
import {PolicyState} from '@disclosure-portal/model/PolicyRule';
import {SpdxIdentifier} from '@disclosure-portal/model/Spdx';
import {
  ComponentDiffType,
  ComponentInfo,
  ComponentMultiDiff,
  PolicyRuleStatus,
} from '@disclosure-portal/model/VersionDetails';
import ProjectService from '@disclosure-portal/services/projects';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useSbomStore} from '@disclosure-portal/stores/sbom.store';
import {escapeHtml} from '@disclosure-portal/utils/Validation';
import {
  formatDateTimeShort,
  getIconColorForPolicyType,
  getIconForDiffType,
  getIconForPolicyType,
  getPrStatusSortIndex,
  sortPolicyStatesByOrder,
} from '@disclosure-portal/utils/View';
import {DataTableHeader, DataTableItem, SortItem} from '@shared/types/table';
import _ from 'lodash';
import {computed, onBeforeMount, ref} from 'vue';
import {useI18n} from 'vue-i18n';

class CellValueVersion {
  public version: string;
  public icon: string;

  constructor(version: string, icon: string) {
    this.version = version;
    this.icon = icon;
  }
}

class CellValue {
  public value: string;
  public version: string;
  public icon: string;

  constructor(value: string, version: string, icon: string) {
    this.value = value;
    this.version = version;
    this.icon = icon;
  }
}

class CellValueWithPrStatus {
  public value: string;
  public version: string;
  public icon: string;
  public unasserted: boolean;
  public questioned: boolean;
  public policyRuleStatus: PolicyRuleStatus[];

  constructor(
    value: string,
    version: string,
    icon: string,
    unasserted: boolean,
    questioned: boolean,
    policyRuleStatus: PolicyRuleStatus[],
  ) {
    this.value = value;
    this.version = version;
    this.icon = icon;
    this.unasserted = unasserted;
    this.questioned = questioned;
    this.policyRuleStatus = policyRuleStatus;
  }
}

class CellValueWithLicenseInformation {
  public value: string;
  public version: string;
  public icon: string;

  public licenseEffective: string;
  public licenseDeclared: string;
  public license: string;
  public usedPolicyRule: string;
  public licenseApplied: string;

  constructor(
    value: string,
    version: string,
    icon: string,
    licenseEffective: string,
    licenseDeclared: string,
    license: string,
    usedPolicyRule: string,
    licenseApplied: string,
  ) {
    this.value = value;
    this.version = version;
    this.icon = icon;
    this.licenseEffective = licenseEffective;
    this.licenseDeclared = licenseDeclared;
    this.license = license;
    this.usedPolicyRule = usedPolicyRule;
    this.licenseApplied = licenseApplied;
  }
}

const currentProject = computed(() => useProjectStore().currentProject!);
const sbomStore = useSbomStore();
const {t} = useI18n();

const version = computed(() => sbomStore.getCurrentVersion);
const allSboms = computed(() => sbomStore.getAllSBOMs);
const groupedSpdxs = computed(() => {
  const res: SpdxIdentifier[] = [];
  allSboms.value.forEach((vs) => {
    const newHeader = new SpdxIdentifier('header', '', '', '', vs.VersionName, '');
    res.push(newHeader);
    for (const spdx of vs.SpdxFileHistory) {
      const uploaded = formatDateTimeShort(spdx.Uploaded);
      const ident = new SpdxIdentifier(spdx._key, spdx.MetaInfo.Name, uploaded, vs.VersionKey, '', spdx.Tag);
      res.push(ident);
      ident.versionName = vs.VersionName;
    }
  });
  return res;
});
const componentHeaders = computed<DataTableHeader[]>(() => {
  return [
    {
      title: t('COL_COMPARE_STATUS'),
      align: 'start',
      width: 80,
      value: 'DiffType',
      sortable: true,
      sortRaw: compareDiffType,
    },
    {
      title: t('COL_SPDX_STATUS'),
      align: 'center',
      value: 'prStatus',
      width: '130',
      sortable: true,
      sortRaw: comparePrStatus,
    },
    {
      title: t('COL_NAME'),
      align: 'start',
      value: 'Name',
      sortable: true,
    },
    {
      title: t('COL_SPDX_VERSION'),
      align: 'start',
      value: 'Version',
      sortable: true,
      sortRaw: compareVersions,
    },
    {
      title: t('COL_SPDX_TYPE'),
      align: 'center',
      value: 'Type',
      width: '150',
      sortable: true,
      sortRaw: compareType,
    },
    {
      title: t('COL_SPDX_LICENSE_EFFECTIVE'),
      align: 'start',
      class: 'tableHeaderCell',
      value: 'LicenseEffective',
      width: '250',
      sortable: true,
      sortRaw: compareLicense,
    },
  ];
});
const filteredList = computed(() => {
  return componentMultiDiffList.value.filter((info: ComponentMultiDiff) => {
    return filterOnLicenseMulti(info) && filterOnPolicyTypeMulti(info);
  });
});

const allLicenses = ref<string[]>([]);
const componentMultiDiffList = ref<ComponentMultiDiff[]>([]);
const search = ref<string>('');
const allPolicyTypes = ref<string[]>([]);
const selectedFilterLicenses = ref<string[]>([]);
const selectedFilterPolicyTypes = ref<string[]>([]);
const selectedSpdxPrevious = ref<SpdxIdentifier | null>(null);
const selectedSpdxCurrent = ref<SpdxIdentifier | null>(null);
const dataAreLoaded = ref(false);
const licenseFilterOpened = ref(false);
const prStatusFilterOpened = ref(false);
const prevError = ref('');
const currentError = ref('');
const compDia = ref<InstanceType<typeof ComponentCompareDialog>>();

const myMenuProps = {
  closeOnClick: false,
  closeOnContentClick: false,
  disableKeys: true,
  openOnClick: false,
  maxHeight: 500,
  bottom: true,
  offsetY: true,
};

const getIconDiffType = (type: ComponentDiffType): string => {
  return getIconForDiffType(type);
};

const compareType = (a: ComponentMultiDiff, b: ComponentMultiDiff): number => {
  const aTypes = getTableViewDataForType(a);
  const bTypes = getTableViewDataForType(b);
  return (bTypes[bTypes.length - 1]?.value || '').localeCompare(aTypes[aTypes.length - 1]?.value || '');
};

const comparePrStatus = (a: ComponentMultiDiff, b: ComponentMultiDiff): number => {
  const aStatus = getTableViewDataForPrStatus(a);
  const bStatus = getTableViewDataForPrStatus(b);
  return (
    getPrStatusSortIndex(bStatus[bStatus.length - 1]?.value || '') -
    getPrStatusSortIndex(aStatus[aStatus.length - 1]?.value || '')
  );
};

const compareDiffType = (a: ComponentMultiDiff, b: ComponentMultiDiff): number => {
  const diffWeight: Map<ComponentDiffType, number> = new Map<ComponentDiffType, number>([
    [ComponentDiffType.UNCHANGED, 0],
    [ComponentDiffType.CHANGED, 1],
    [ComponentDiffType.REMOVED, 2],
    [ComponentDiffType.NEW, 3],
  ]);
  return diffWeight.get(b.DiffType)! - diffWeight.get(a.DiffType)!;
};

const compareLicense = (a: ComponentMultiDiff, b: ComponentMultiDiff): number => {
  const aLic = getTableViewDataForLicenseEffective(a);
  const bLic = getTableViewDataForLicenseEffective(b);
  return (bLic[bLic.length - 1]?.value || '').localeCompare(aLic[aLic.length - 1]?.value || '');
};

const compareVersions = (a: ComponentMultiDiff, b: ComponentMultiDiff): number => {
  return plainVersion(b).localeCompare(plainVersion(a));
};

const plainVersion = (diff: ComponentMultiDiff) => {
  if (!diff) {
    return '';
  }
  if (diff.DiffType === 'UNCHANGED' || diff.DiffType === 'REMOVED') {
    return diff.ComponentsOld[0].version;
  } else if (diff.DiffType === 'NEW') {
    return diff.ComponentsNew[0].version;
  } else if (diff.DiffType === 'CHANGED') {
    return diff.ComponentsNew[0].version;
  }
  return '';
};

const filterOnLicenseMulti = (info: ComponentMultiDiff): boolean => {
  if (selectedFilterLicenses.value.length > 0) {
    if (hasMultipleLicense(info)) {
      let found = false;
      separateMultipleLicenses(info).forEach((lic: string) => {
        if (!found && selectedFilterLicenses.value.includes(lic)) {
          found = true;
        }
      });
      return found;
    } else {
      let found = false;
      const singleLicenses = info.ComponentsOld.map((comp) => comp.licenseEffective);
      singleLicenses.push(...info.ComponentsNew.map((comp) => comp.licenseEffective));
      singleLicenses.forEach((license) => {
        if (!found && selectedFilterLicenses.value.includes(license)) {
          found = true;
        }
      });
      return found;
    }
  } else {
    return true;
  }
};

const checkMultipleLicenses = (comp: ComponentInfo): boolean =>
  comp.licenseEffective.indexOf(' AND ') > 0 || comp.license.indexOf(' OR ') > 0;

const hasMultipleLicense = (info: ComponentMultiDiff): boolean =>
  _.some(info.ComponentsNew, checkMultipleLicenses) || _.some(info.ComponentsOld, checkMultipleLicenses);

const multipleLicensesToArray = (comp: ComponentInfo): string[] =>
  comp.licenseEffective
    .replace('(', '')
    .replace(')', '')
    .split(/ OR | AND /);

const separateMultipleLicenses = (info: ComponentMultiDiff): string[] => {
  const licensesNew = info.ComponentsNew.flatMap(multipleLicensesToArray);
  const licensesOld = info.ComponentsOld.flatMap(multipleLicensesToArray);
  return _.union(licensesNew, licensesOld);
};

const computeMultiLicensesFilter = () => {
  const computedLicenses: string[] = [];
  componentMultiDiffList.value.forEach((info) => {
    if (hasMultipleLicense(info)) {
      separateMultipleLicenses(info).forEach((lic: string) => {
        if (!computedLicenses.includes(lic)) {
          computedLicenses.push(lic.trim());
        }
      });
    } else {
      const pushIfNotIncludes = (comp: ComponentInfo): void => {
        if (!computedLicenses.includes(comp.licenseEffective)) {
          computedLicenses.push(comp.licenseEffective);
        }
      };
      info.ComponentsOld.forEach(pushIfNotIncludes);
      info.ComponentsNew.forEach(pushIfNotIncludes);
    }
  });
  allLicenses.value = computedLicenses.sort();
};

const filterOnPolicyTypeMulti = (info: ComponentMultiDiff): boolean => {
  if (selectedFilterPolicyTypes.value.length > 0) {
    if (
      _.some(info.ComponentsNew, (comp) => !!comp.prStatus) ||
      _.some(info.ComponentsOld, (comp) => !!comp.prStatus)
    ) {
      const allPrStatus = info.ComponentsNew.filter((comp) => !!comp.prStatus).map((comp) => comp.prStatus);
      allPrStatus.push(...info.ComponentsOld.filter((comp) => !!comp.prStatus).map((comp) => comp.prStatus));
      for (const filterType of selectedFilterPolicyTypes.value) {
        const included = allPrStatus.includes(filterType);
        if (included) {
          return true;
        }
      }
    }
    return false;
  } else {
    return true;
  }
};

const computeMultiPolicyTypesFilter = () => {
  const computedPolicyTypes: string[] = [];
  const computePolicyTypesFilter = (info: ComponentInfo): void => {
    if (info.policyRuleStatus) {
      info.policyRuleStatus.forEach((policyRuleStatus: PolicyRuleStatus) => {
        if (!computedPolicyTypes.includes(policyRuleStatus.type) && policyRuleStatus.used) {
          computedPolicyTypes.push(policyRuleStatus.type);
        }
      });
    }
    if (info.questioned && !computedPolicyTypes.includes(PolicyState.QUESTIONED)) {
      computedPolicyTypes.push(PolicyState.QUESTIONED);
    }
    if (info.unasserted && !computedPolicyTypes.includes(PolicyState.NOASSERTION)) {
      computedPolicyTypes.push(PolicyState.NOASSERTION);
    }
  };
  componentMultiDiffList.value.forEach((info) => {
    info.ComponentsOld.forEach(computePolicyTypesFilter);
    info.ComponentsNew.forEach(computePolicyTypesFilter);
  });
  allPolicyTypes.value = computedPolicyTypes.sort(sortPolicyStatesByOrder);
};

const showDetails = (e: Event, row: DataTableItem<ComponentMultiDiff>) => {
  if (!selectedSpdxPrevious.value || !selectedSpdxCurrent.value) {
    return;
  }
  compDia.value?.open(row.item, row.item.Name, selectedSpdxPrevious.value, selectedSpdxCurrent.value);
};

const resetErrors = () => {
  prevError.value = '';
  currentError.value = '';
};

const reloadInternal = async (forceReload: boolean) => {
  if (!forceReload && dataAreLoaded.value) return;
  if (!currentProject.value._key) {
    return;
  }
  await sbomStore.fetchAllSBOMsFlat();
  selectedFilterLicenses.value = [];
  resetErrors();

  setSelection();
};

const setSelection = () => {
  const selectedSpdx = sbomStore.getSelectedSBOM;
  if (selectedSpdx === null) {
    selectedSpdxCurrent.value = groupedSpdxs.value.find((s) => s.versionKey === version.value._key) || null;
  } else {
    selectedSpdxCurrent.value = groupedSpdxs.value.find((s) => s.spdxFileId === selectedSpdx?._key) || null;
  }
  selectedSpdxPrevious.value = null;
};

const compare = async (): Promise<void> => {
  if (!selectedSpdxPrevious.value) {
    prevError.value = 'Field is required!';
    return;
  }
  if (!selectedSpdxCurrent.value) {
    currentError.value = 'Field is required!';
    return;
  }
  componentMultiDiffList.value = (
    await ProjectService.compareSpdxFiles(
      currentProject.value._key,
      selectedSpdxPrevious.value.versionKey,
      selectedSpdxPrevious.value.spdxFileId,
      selectedSpdxCurrent.value.versionKey,
      selectedSpdxCurrent.value.spdxFileId,
    )
  ).data;

  computeMultiLicensesFilter();
  computeMultiPolicyTypesFilter();
};

const getTableViewDataForVersion = (componentMultiDiff: ComponentMultiDiff): CellValueVersion[] => {
  let result: CellValueVersion[] = [];
  if (componentMultiDiff.DiffType === 'CHANGED') {
    const oldValues = componentMultiDiff.ComponentsOld.map((comp) => new CellValueVersion(comp.version, 'remove'));
    const newValues: CellValueVersion[] = [];
    for (const comp of componentMultiDiff.ComponentsNew) {
      if (_.some(oldValues, {version: comp.version})) {
        const componentChanges = componentMultiDiff.Changes[`${comp.version}_${comp.version}`];
        if (componentChanges) {
          if (
            componentChanges.Version ||
            componentChanges.Type ||
            componentChanges.LicenseEffective ||
            componentChanges.prStatus
          ) {
            // one of the in-table visible changes => both 'add'/'remove' needed as placeholder to show diff of a cell in two local rows
            newValues.push(new CellValueVersion(comp.version, Icons.ADD));
          } else {
            // changes are present but visible in Details Dialog only => remove already added one and notify per 'change' icon with single local row
            const findIndex = _.findIndex(oldValues, {version: comp.version});
            oldValues.splice(findIndex, 1);
            newValues.push(new CellValueVersion(comp.version, Icons.CHANGED));
          }
        } else {
          // two same versions which is the same locally without any diff
          const findIndex = _.findIndex(oldValues, {version: comp.version});
          oldValues.splice(findIndex, 1);
          newValues.push(new CellValueVersion(comp.version, ''));
        }
      } else {
        newValues.push(new CellValueVersion(comp.version, Icons.ADD));
      }
    }
    oldValues.push(...newValues);
    result = oldValues;
  } else if (componentMultiDiff.DiffType === 'UNCHANGED' || componentMultiDiff.DiffType === 'REMOVED') {
    result = componentMultiDiff.ComponentsOld.map((comp) => new CellValueVersion(comp.version, ''));
  } else if (componentMultiDiff.DiffType === 'NEW') {
    result = componentMultiDiff.ComponentsNew.map((comp) => new CellValueVersion(comp.version, ''));
  }

  if (result.length === 1) {
    return result.map((value) => new CellValueVersion(value.version, ''));
  }
  return _.orderBy(result, (i) => i.version, 'asc');
};

const getTableViewDataForType = (componentMultiDiff: ComponentMultiDiff): CellValue[] => {
  let result: CellValue[] = [];
  if (componentMultiDiff.DiffType === 'CHANGED') {
    const oldValues = componentMultiDiff.ComponentsOld.map((comp) => componentInfoToCellValueType(comp, ''));
    const newValues: CellValue[] = [];
    for (const comp of componentMultiDiff.ComponentsNew) {
      if (_.some(oldValues, {version: comp.version})) {
        const componentChanges = componentMultiDiff.Changes[`${comp.version}_${comp.version}`];
        if (componentChanges) {
          if (componentChanges.Type) {
            // one of the in-table visible changes exactly for the field => both '+'/'-' needed to show diff of a field in two local rows
            const findIndex = _.findIndex(oldValues, {version: comp.version});
            oldValues[findIndex].icon = Icons.REMOVED;
            newValues.push(componentInfoToCellValueType(comp, Icons.ADD));
          } else if (componentChanges.Version || componentChanges.LicenseEffective || componentChanges.prStatus) {
            // other in-table visible changes not for the field => both without icons needed to reserve two local rows to show the diff of other field in other column with '+'/'-'
            newValues.push(componentInfoToCellValueType(comp, ''));
          }
        }
      } else {
        newValues.push(componentInfoToCellValueType(comp, ''));
      }
    }
    oldValues.push(...newValues);
    result = oldValues;
  } else if (componentMultiDiff.DiffType === 'UNCHANGED' || componentMultiDiff.DiffType === 'REMOVED') {
    result = componentMultiDiff.ComponentsOld.map((comp) => componentInfoToCellValueType(comp, ''));
  } else if (componentMultiDiff.DiffType === 'NEW') {
    result = componentMultiDiff.ComponentsNew.map((comp) => componentInfoToCellValueType(comp, ''));
  }

  if (result.length === 1) {
    return result.map((value) => new CellValue(value.value, value.version, ''));
  }
  return _.orderBy(result, (i) => i.version, 'asc');
};

const componentInfoToCellValueType = (comp: ComponentInfo, icon: string): CellValue => {
  return new CellValue(comp.type, comp.version, icon);
};

const getTableViewDataForLicenseEffective = (
  componentMultiDiff: ComponentMultiDiff,
): CellValueWithLicenseInformation[] => {
  let result: CellValueWithLicenseInformation[] = [];
  if (componentMultiDiff.DiffType === 'CHANGED') {
    const oldValues = componentMultiDiff.ComponentsOld.map((comp) => componentInfoToCellValueLicense(comp, ''));
    const newValues: CellValueWithLicenseInformation[] = [];
    for (const comp of componentMultiDiff.ComponentsNew) {
      if (_.some(oldValues, {version: comp.version})) {
        const componentChanges = componentMultiDiff.Changes[`${comp.version}_${comp.version}`];
        if (componentChanges) {
          if (componentChanges.LicenseEffective) {
            // one of the in-table visible changes exactly for the field => both '+'/'-' needed to show diff of a field in two local rows
            const findIndex = _.findIndex(oldValues, {version: comp.version});
            oldValues[findIndex].icon = Icons.REMOVED;
            newValues.push(componentInfoToCellValueLicense(comp, Icons.ADD));
          } else if (componentChanges.Version || componentChanges.Type || componentChanges.prStatus) {
            // other in-table visible changes not for the field => both without icons needed to reserve two local rows to show the diff of other field in other column with '+'/'-'
            newValues.push(componentInfoToCellValueLicense(comp, ''));
          }
        }
      } else {
        newValues.push(componentInfoToCellValueLicense(comp, ''));
      }
    }
    oldValues.push(...newValues);
    result = oldValues;
  } else if (componentMultiDiff.DiffType === 'UNCHANGED' || componentMultiDiff.DiffType === 'REMOVED') {
    result = componentMultiDiff.ComponentsOld.map((comp) => componentInfoToCellValueLicense(comp, ''));
  } else if (componentMultiDiff.DiffType === 'NEW') {
    result = componentMultiDiff.ComponentsNew.map((comp) => componentInfoToCellValueLicense(comp, ''));
  }

  if (result.length === 1) {
    return result.map(
      (value) =>
        new CellValueWithLicenseInformation(
          value.value,
          value.version,
          '',
          value.licenseEffective,
          value.licenseDeclared,
          value.license,
          value.usedPolicyRule,
          value.licenseApplied,
        ),
    );
  }
  return _.orderBy(result, (i) => i.version, 'asc');
};

const componentInfoToCellValueLicense = (comp: ComponentInfo, icon: string): CellValueWithLicenseInformation => {
  const licenseEffective = comp.licenseEffective || 'NOASSERTION';
  const licenseDeclared = comp.licenseDeclared || 'NOASSERTION';
  const license = comp.license || 'NOASSERTION';
  return new CellValueWithLicenseInformation(
    licenseEffective,
    comp.version,
    icon,
    licenseEffective,
    licenseDeclared,
    license,
    comp.usedPolicyRule,
    comp.licenseApplied,
  );
};

const getTableViewDataForPrStatus = (componentMultiDiff: ComponentMultiDiff): CellValueWithPrStatus[] => {
  let result: CellValueWithPrStatus[] = [];
  if (componentMultiDiff.DiffType === 'CHANGED') {
    const oldValues = componentMultiDiff.ComponentsOld.map((comp) => componentInfoToCellValuePrStatus(comp, ''));
    const newValues: CellValueWithPrStatus[] = [];
    for (const comp of componentMultiDiff.ComponentsNew) {
      if (_.some(oldValues, {version: comp.version})) {
        const componentChanges = componentMultiDiff.Changes[`${comp.version}_${comp.version}`];
        if (componentChanges) {
          if (componentChanges.prStatus) {
            // one of the in-table visible changes exactly for the field => both '+'/'-' needed to show diff of a field in two local rows
            const findIndex = _.findIndex(oldValues, {version: comp.version});
            oldValues[findIndex].icon = Icons.REMOVED;
            newValues.push(componentInfoToCellValuePrStatus(comp, Icons.ADD));
          } else if (componentChanges.Version || componentChanges.Type || componentChanges.LicenseEffective) {
            // other in-table visible changes not for the field => both without icons needed to reserve two local rows to show the diff of other field in other column with '+'/'-'
            newValues.push(componentInfoToCellValuePrStatus(comp, ''));
          }
        }
      } else {
        newValues.push(componentInfoToCellValuePrStatus(comp, ''));
      }
    }
    oldValues.push(...newValues);
    result = oldValues;
  } else if (componentMultiDiff.DiffType === 'UNCHANGED' || componentMultiDiff.DiffType === 'REMOVED') {
    result = componentMultiDiff.ComponentsOld.map((comp) => componentInfoToCellValuePrStatus(comp, ''));
  } else if (componentMultiDiff.DiffType === 'NEW') {
    result = componentMultiDiff.ComponentsNew.map((comp) => componentInfoToCellValuePrStatus(comp, ''));
  }

  if (result.length === 1) {
    return result.map(
      (value) =>
        new CellValueWithPrStatus(
          value.value,
          value.version,
          '',
          value.unasserted,
          value.questioned,
          value.policyRuleStatus,
        ),
    );
  }
  return _.orderBy(result, (i) => i.version, 'asc');
};

const componentInfoToCellValuePrStatus = (comp: ComponentInfo, icon: string): CellValueWithPrStatus => {
  return new CellValueWithPrStatus(
    comp.prStatus,
    comp.version,
    icon,
    comp.unasserted,
    comp.questioned,
    comp.policyRuleStatus,
  );
};

const onSelectionChanged = () => {
  componentMultiDiffList.value = [];
};
const formatText = (text: string): string => {
  return escapeHtml(text)
    .replace(/ AND /gi, ' <strong class="db-highlight">AND</strong> ')
    .replace(/ OR /gi, ' <strong class="db-highlight">OR</strong> ')
    .replace(/ and /gi, ' <strong class="db-highlight">AND</strong> ')
    .replace(/ or /gi, ' <strong class="db-highlight">OR</strong> ');
};

const sortBy: SortItem[] = [{key: 'DiffType', order: 'asc'}];

onBeforeMount(async () => {
  await reloadInternal(true);
});
</script>

<template>
  <v-form ref="compareForm">
    <TableLayout has-title has-tab>
      <template #description>
        <div class="d-flex ga-4 align-center flex-row flex-wrap">
          <v-row>
            <v-col xs="12" md="6">
              <v-select
                v-model="selectedSpdxCurrent"
                variant="outlined"
                clearable
                density="compact"
                hide-details
                :items="groupedSpdxs"
                :label="t('SBOM_COMPARE_CURRENT')"
                item-text="label"
                item-value="spdxFileId"
                return-object
                :error-messages="currentError"
                @click:clear="resetErrors"
                v-on:change="resetErrors"
                @change="onSelectionChanged()">
                <template #item="{item, props}">
                  <v-list-item v-if="item.raw.spdxFileId === 'header'" title="">
                    <span class="d-subtitle-2">{{ item.raw.header }}</span>
                  </v-list-item>
                  <v-list-item v-bind="props" title="" v-else>
                    <v-icon
                      color="primary"
                      v-if="currentProject && currentProject.approvablespdx.spdxkey == item.raw.spdxFileId"
                      size="small"
                      class="pb-1"
                      >mdi-star
                    </v-icon>
                    <span class="d-subtitle-2 ml-2">{{ item.raw.uploaded }}</span>
                    <span class="d-text d-secondary-text">&nbsp;-&nbsp;{{ item.raw.label }}</span>
                    <span class="d-text d-secondary-text" v-if="item.raw.tag">&nbsp;({{ item.raw.tag }})</span>
                  </v-list-item>
                </template>
                <template #selection="{item}">
                  <div class="d-inline">
                    <v-icon
                      color="primary"
                      v-if="currentProject && currentProject.approvablespdx.spdxkey == item.raw.spdxFileId"
                      size="small"
                      class="pr-2 pb-1"
                      >mdi-star
                    </v-icon>
                    <span class="d-subtitle-2">{{ item.raw.uploaded }}</span>
                    <span class="d-text d-secondary-text">&nbsp;-&nbsp;{{ item.raw.label }}</span>
                    <span class="d-text d-secondary-text" v-if="item.raw.tag">&nbsp;({{ item.raw.tag }})</span>
                  </div>
                </template>
              </v-select>
            </v-col>
            <v-col xs="12" md="6">
              <v-select
                v-model="selectedSpdxPrevious"
                variant="outlined"
                hide-details="auto"
                clearable
                density="compact"
                :items="groupedSpdxs"
                :label="t('SBOM_COMPARE_PREVIOUS')"
                item-text="label"
                item-value="spdxFileId"
                return-object
                :menu-props="myMenuProps"
                :error-messages="prevError"
                @click:clear="resetErrors"
                v-on:change="resetErrors"
                @change="onSelectionChanged()">
                <template #item="{item, props}">
                  <v-list-item v-if="item.raw.spdxFileId === 'header'" title="">
                    <span class="d-subtitle-2">{{ item.raw.header }}</span>
                  </v-list-item>
                  <v-list-item v-bind="props" title="" v-else>
                    <v-icon
                      color="primary"
                      v-if="currentProject && currentProject.approvablespdx.spdxkey == item.raw.spdxFileId"
                      size="small"
                      class="pb-1"
                      >mdi-star
                    </v-icon>
                    <span class="d-subtitle-2 ml-2">{{ item.raw.uploaded }}</span>
                    <span class="d-text d-secondary-text">&nbsp;-&nbsp;{{ item.raw.label }}</span>
                    <span class="d-text d-secondary-text" v-if="item.raw.tag">&nbsp;({{ item.raw.tag }})</span>
                  </v-list-item>
                </template>
                <template #selection="{item}">
                  <div class="d-inline">
                    <v-icon
                      color="primary"
                      v-if="currentProject && currentProject.approvablespdx.spdxkey == item.raw.spdxFileId"
                      size="small"
                      class="pr-2 pb-1"
                      >mdi-star
                    </v-icon>
                    <span class="d-subtitle-2">{{ item.raw.uploaded }}</span>
                    <span class="d-text d-secondary-text">&nbsp;-&nbsp;{{ item.raw.label }}</span>
                    <span class="d-text d-secondary-text" v-if="item.raw.tag">&nbsp;({{ item.raw.tag }})</span>
                  </div>
                </template>
              </v-select>
            </v-col>
          </v-row>
          <DCActionButton
            large
            icon="mdi-compare-horizontal"
            :text="t('COMPARE')"
            :hint="t('TT_compare_sboms')"
            @click="compare" />
        </div>
      </template>
      <template #buttons>
        <v-spacer></v-spacer>
        <DSearchField v-model="search" class="mr-1" />
      </template>
      <template #table>
        <div ref="sgrid" class="fill-height">
          <v-data-table
            fixed-header
            :headers="componentHeaders"
            density="compact"
            class="striped-table fill-height"
            :items-per-page="100"
            :footer-props="{
              'items-per-page-options': [10, 50, 100, -1],
            }"
            :items="filteredList"
            :search="search"
            @click:row="showDetails"
            :sort-by="sortBy">
            <template #[`header.LicenseEffective`]="{column, getSortIcon, toggleSort}">
              <div class="v-data-table-header__content">
                <span>{{ column.title }}</span>
                <v-menu :close-on-content-click="false" v-model="licenseFilterOpened">
                  <template #activator="{props}">
                    <DIconButton
                      :parentProps="props"
                      icon="mdi-filter-variant"
                      :hint="t('TT_SHOW_FILTER')"
                      :color="selectedFilterLicenses.length > 0 ? 'primary' : 'default'"
                      location="top" />
                  </template>
                  <div class="bg-background" style="width: 280px">
                    <v-row class="d-flex ma-1 mr-2 justify-end">
                      <DIconButton icon="mdi-close" @clicked="licenseFilterOpened = false" color="default" />
                    </v-row>
                    <v-autocomplete
                      v-model="selectedFilterLicenses"
                      :items="allLicenses"
                      class="pa-2 mx-2 pb-4"
                      :label="t('Lbl_filter_license')"
                      clearable
                      multiple
                      variant="outlined"
                      density="compact"
                      menu
                      transition="scale-transition"
                      persistent-clear
                      :list-props="{class: 'striped-filter-dd py-0'}">
                      <template #item="{props, item}">
                        <v-list-item v-bind="props" class="px-2 py-0" title="">
                          <template #prepend="{isSelected}">
                            <v-checkbox hide-details :model-value="isSelected" />
                          </template>
                          <span class="pFilterEntry">
                            {{ item.title }}
                          </span>
                        </v-list-item>
                      </template>
                      <template #selection="{item, index}">
                        <div v-if="index === 0" class="pFilterEntry">
                          {{ item.raw }}
                        </div>
                        <span v-if="index === 1" class="pAdditionalFilter">
                          +{{ selectedFilterLicenses.length - 1 }} others
                        </span>
                      </template>
                    </v-autocomplete>
                  </div>
                </v-menu>
                <v-icon
                  class="v-data-table-header__sort-icon"
                  :icon="getSortIcon(column)"
                  @click="toggleSort(column)" />
              </div>
            </template>
            <template #[`header.prStatus`]="{column, getSortIcon, toggleSort}">
              <div class="v-data-table-header__content">
                <span>{{ column.title }}</span>
                <v-menu :close-on-content-click="false" v-model="prStatusFilterOpened">
                  <template #activator="{props}">
                    <DIconButton
                      :parentProps="props"
                      icon="mdi-filter-variant"
                      :hint="t('TT_SHOW_FILTER')"
                      :color="selectedFilterPolicyTypes.length > 0 ? 'primary' : 'default'" />
                  </template>
                  <div class="bg-background" style="width: 280px">
                    <v-row class="d-flex ma-1 mr-2 justify-end">
                      <DIconButton icon="mdi-close" @clicked="prStatusFilterOpened = false" color="default" />
                    </v-row>
                    <v-select
                      v-model="selectedFilterPolicyTypes"
                      :items="allPolicyTypes"
                      class="pa-2 mx-2"
                      :label="t('Lbl_filter_policyType')"
                      clearable
                      multiple
                      item-title="text"
                      item-value="value"
                      variant="outlined"
                      density="compact"
                      clearalbe
                      menu
                      transition="scale-transition"
                      persistent-clea
                      :list-props="{class: 'striped-filter-dd py-0'}">
                      <template #item="{item, props}">
                        <v-list-item v-bind="props" class="px-2 py-0">
                          <template #prepend="{isSelected}">
                            <v-checkbox hide-details :model-value="isSelected" />
                          </template>
                          <template #title>
                            <v-icon small :color="getIconColorForPolicyType(item.raw)"
                              >{{ getIconForPolicyType(item.raw) }}
                            </v-icon>
                            <span class="pFilterEntry pl-1">{{ item.raw }}</span>
                          </template>
                        </v-list-item>
                      </template>
                      <template #selection="{item, index}">
                        <div v-if="index === 0" class="d-flex align-center">
                          <v-icon small :color="getIconColorForPolicyType(item.raw)"
                            >{{ getIconForPolicyType(item.raw) }}
                          </v-icon>
                          <span class="pFilterEntry pl-1">{{ item.raw }}</span>
                        </div>
                        <span v-if="index === 1" class="pAdditionalFilter">
                          +{{ selectedFilterPolicyTypes.length - 1 }} others
                        </span>
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
            <template #[`item.DiffType`]="{item}">
              <span>
                <v-icon small>{{ getIconDiffType(item.DiffType) }}</v-icon>
                <Tooltip>
                  {{ t(item.DiffType) }}
                </Tooltip>
              </span>
            </template>
            <template #[`item.Version`]="{item}">
              <div class="flex flex-col" style="gap: 5px">
                <div v-for="(entry, i) in getTableViewDataForVersion(item)" :key="i">
                  <v-icon style="width: 16px; height: 16px" class="mr-2" small>{{ entry.icon }}</v-icon>
                  <span>{{ entry.version }}</span>
                </div>
              </div>
            </template>
            <template #[`item.Type`]="{item}">
              <div class="flex flex-col" style="gap: 5px">
                <div v-for="(entry, i) in getTableViewDataForType(item)" :key="i">
                  <v-icon style="width: 16px; height: 16px" class="mr-2" small>{{ entry.icon }}</v-icon>
                  <span>{{ entry.value }}</span>
                </div>
              </div>
            </template>
            <template #[`item.LicenseEffective`]="{item}">
              <div class="flex flex-col gap-2">
                <div v-for="(entry, i) in getTableViewDataForLicenseEffective(item)" :key="i">
                  <span>
                    <v-icon v-if="entry.icon" class="mr-2 size-4" small>{{ entry.icon }}</v-icon>
                    <span v-html="formatText(entry.value)"></span>

                    <Tooltip>
                      <div>
                        <div>{{ t('COL_SPDX_LICENSE_EFFECTIVE') }}:</div>
                        <div class="mr-2 mb-5" v-html="formatText(entry.licenseEffective)"></div>
                        <div>{{ t('COL_SPDX_LICENSE_DECLARED') }}:</div>
                        <div class="mr-2 mb-5" v-html="formatText(entry.licenseDeclared)"></div>
                        <div>{{ t('COL_SPDX_LICENSE_CONCLUDED') }}:</div>
                        <div class="mr-2 mb-5" v-html="formatText(entry.license)"></div>
                        <span>{{ entry.usedPolicyRule }} (based on {{ entry.licenseApplied }})</span>
                      </div>
                    </Tooltip>
                  </span>
                </div>
              </div>
            </template>
            <template #[`item.prStatus`]="{item}">
              <div class="flex flex-col gap-2">
                <div v-for="(entry, i) in getTableViewDataForPrStatus(item)" :key="i">
                  <span>
                    <v-icon small style="width: 16px; height: 16px" class="mr-2">{{ entry.icon }}</v-icon>
                    <v-icon small :color="getIconColorForPolicyType(entry.value)">
                      {{ getIconForPolicyType(entry.value) }}
                    </v-icon>

                    <Tooltip>
                      <span v-if="entry.unasserted">{{ t('UNASSERTED') }}</span>
                      <span v-else-if="entry.questioned">{{ t('TT_questioned') }}</span>
                      <span v-else-if="entry.policyRuleStatus">
                        <span v-for="pr in entry.policyRuleStatus" :key="pr.name">{{ pr.name }}</span>
                      </span>
                      <span v-else></span>
                    </Tooltip>
                  </span>
                </div>
              </div>
            </template>
          </v-data-table>
        </div>
      </template>
    </TableLayout>
  </v-form>
  <ComponentCompareDialog ref="compDia"></ComponentCompareDialog>
</template>
