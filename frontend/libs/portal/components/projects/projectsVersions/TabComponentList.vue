<script setup lang="ts">
import {DecisionType, DialogBulkPolicyDecisionEntry} from '@disclosure-portal/components/dialog/DialogConfigs';
import useDimensions from '@disclosure-portal/composables/useDimensions';
import {useLicense} from '@disclosure-portal/composables/useLicense';
import {compareFamily} from '@disclosure-portal/model/License';
import {PolicyState, PolicyStates} from '@disclosure-portal/model/PolicyRule';
import {
  ComponentInfo,
  ComponentInfoSlim,
  ComponentStats,
  PolicyRuleStatus,
} from '@disclosure-portal/model/VersionDetails';
import ProjectService from '@disclosure-portal/services/projects';
import VersionService from '@disclosure-portal/services/version';
import {useAppStore} from '@disclosure-portal/stores/app';
import {useIdleStore} from '@disclosure-portal/stores/idle.store';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import eventBus from '@disclosure-portal/utils/eventbus';
import {formatDateAndTime} from '@disclosure-portal/utils/Table';
import {escapeHtml} from '@disclosure-portal/utils/Validation';
import {
  getIconColorForPolicyType,
  getIconForPolicyType,
  getPrStatusSortIndex,
  policyStateToTranslationKey,
  sortPolicyStatesByOrder,
} from '@disclosure-portal/utils/View';
import {IRuleBtnCallbacks} from '@shared/components/disco/interfaces';
import {useHeaderSettingsStore} from '@shared/stores/headerSettings.store';
import {DataTableHeader, DataTableHeaderFilterItems, DataTableItem, SortItem} from '@shared/types/table';
import {storeToRefs} from 'pinia';
import {computed, nextTick, onMounted, onUnmounted, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRoute} from 'vue-router';

type TabelItem = ComponentInfo & {
  showPolicyDecision: boolean;
  showLicenseDecision: boolean;
};

const route = useRoute();
const {t} = useI18n();
const {getI18NTextOfPrefixKey} = useLicense();
const appStore = useAppStore();
const projectStore = useProjectStore();
const idle = useIdleStore();

const gridName = 'ComponentList';
const headerSettingsStore = useHeaderSettingsStore();
const {filteredHeaders} = storeToRefs(headerSettingsStore);

const projectModel = computed(() => projectStore.currentProject!);
const versionDetails = computed(() => appStore.getCurrentVersion);
const spdxFileHistory = computed(() => appStore.getChannelSpdxs);
const currentSpdx = computed(() => appStore.getSelectedSpdx);

const search = ref('');
const sortBy = ref<SortItem[]>([{key: 'prStatus', order: 'desc'}]);
const policies = ref(PolicyStates);
const selectedFilterPolicyTypes = ref<PolicyState[]>([]);
const selectedFilterLicenses = ref<string[]>([]);
const selectedFilterTypes = ref<string[]>([]);
const selectedFilterFamily = ref<string[]>([]);
const dataIsLoaded = ref(false);
const componentList = ref<TabelItem[]>([]);
const stats = ref<ComponentStats>(new ComponentStats());
const allLicenses = ref<DataTableHeaderFilterItems[]>([]);
const allPolicyTypes = ref<DataTableHeaderFilterItems[]>([]);
const allTypes = ref<DataTableHeaderFilterItems[]>([]);
const allFamilies = ref<DataTableHeaderFilterItems[]>([]);
const forceReload = ref(true);
const tableHeight = ref(0);
const {calculateHeight} = useDimensions();
const tableComponents = ref<HTMLElement | null>(null);
const newComponentDetailsDlg = ref();
const licenseRuleDialog = ref();
const policyDecisionDialog = ref();
const bulkPolicyDecisionsDialog = ref();
const bulkPolicyDecisionDeniedReason = ref('');

const componentId = computed(() =>
  Array.isArray(route.params?.componentId) ? route.params.componentId[0] : route.params?.componentId || '',
);

const updateTableHeight = async () => {
  await nextTick();
  if (tableComponents.value) {
    tableHeight.value = calculateHeight(tableComponents.value, true, true);
  }
};

const headers: DataTableHeader[] = [
  {
    title: 'COL_SPDX_STATUS',
    sortable: true,
    align: 'center',
    value: 'prStatus',
    width: 150,
    selectable: true,
  },
  {
    title: 'COL_POLICY_DECISION',
    sortable: true,
    align: 'center',
    value: 'showPolicyDecision',
    width: 100,
    selectable: true,
  },
  {
    title: 'COL_LICENSE_DECISION',
    sortable: true,
    align: 'center',
    value: 'showLicenseDecision',
    width: 100,
    selectable: true,
  },
  {
    title: 'COL_COMPONENT_NAME',
    align: 'start',
    value: 'name',
    selectable: true,
    sortable: true,
    width: 250,
  },
  {
    title: 'COL_COMPONENT_VERSION',
    align: 'start',
    value: 'version',
    width: 160,
    selectable: true,
    sortable: true,
  },
  {
    title: 'COL_SPDX_TYPE',
    align: 'start',
    sortable: true,
    value: 'type',
    width: 150,
    selectable: true,
  },
  {
    title: 'COL_PURL',
    align: 'start',
    width: 300,
    sortable: false,
    value: 'purl',
    selectable: true,
  },
  {
    title: 'COL_LICENSE_FAMILY',
    align: 'start',
    sortable: true,
    width: 180,
    value: 'worstFamily',
    selectable: true,
  },
  {
    title: 'COL_SPDX_LICENSE_EFFECTIVE',
    align: 'start',
    width: 200,
    value: 'licenseEffective',
    selectable: true,
    sortable: true,
  },
];

headerSettingsStore.setupStore(gridName, headers);

const filteredList = computed(() => {
  return componentList.value.filter((info: TabelItem) => {
    return filterOnLicense(info) && filterOnPolicyType(info) && filterOnType(info) && filterOnFamily(info);
  });
});

const filterOnLicense = (info: TabelItem): boolean => {
  if (selectedFilterLicenses.value.length <= 0) {
    return true;
  } else {
    return selectedFilterLicenses.value.some(
      (selectedFilterLicense) => info.licenseEffective.toUpperCase().indexOf(selectedFilterLicense) > -1,
    );
  }
};

const filterOnType = (info: TabelItem): boolean => {
  if (selectedFilterTypes.value.length <= 0) {
    return true;
  } else {
    return selectedFilterTypes.value.some((filterType) => info.type === filterType);
  }
};

const filterOnFamily = (info: TabelItem): boolean => {
  if (selectedFilterFamily.value.length <= 0) {
    return true;
  } else {
    return selectedFilterFamily.value.some(
      (filterFamily) => getI18NTextOfPrefixKey('LIC_FAMILY_', info.worstFamily) === filterFamily,
    );
  }
};

const filterOnPolicyType = (info: TabelItem): boolean => {
  if (selectedFilterPolicyTypes.value.length > 0 && !selectedFilterPolicyTypes.value.includes(PolicyState.NOT_SET)) {
    if (info.prStatus) {
      for (const filterType of selectedFilterPolicyTypes.value) {
        if (info.prStatus.includes(filterType)) {
          return true;
        }
      }
    }
    return false;
  } else {
    return true;
  }
};

const customKeySort = {
  prStatus: (a: string, b: string) => {
    const value1Str = a ?? '';
    const value2Str = b ?? '';
    const prStatus1Index = getPrStatusSortIndex(value1Str);
    const prStatus2Index = getPrStatusSortIndex(value2Str);
    return prStatus2Index - prStatus1Index;
  },
  worstFamily: (a: string, b: string) => {
    const value1Str = a ?? '';
    const value2Str = b ?? '';
    return compareFamily(value2Str, value1Str);
  },
  type: (a: string, b: string) => {
    if (a === 'Root' && b === 'Root') return 0;
    if (a === 'Root') return -1; // a < b
    if (b === 'Root') return 1; // a > b

    const value1Str = a ?? '';
    const value2Str = b ?? '';
    return value1Str.localeCompare(value2Str);
  },
};

const showDetails = async (item: TabelItem) => {
  idle.show();

  const response = await ProjectService.getComponentDetailsForSbom(
    projectModel.value._key,
    versionDetails.value._key,
    currentSpdx.value._key,
    item.spdxId,
  );

  if (newComponentDetailsDlg.value) {
    newComponentDetailsDlg.value?.open(
      projectModel.value,
      versionDetails.value._key,
      currentSpdx.value._key,
      response.data,
      item.policyRuleStatus,
      item.unmatchedLicenses,
      item.policyDecisionsApplied,
      item.policyDecisionDeniedReason,
    );
  }

  idle.hide();
};

const openLicenseRuleDialog = (item: TabelItem) => {
  const component = new ComponentInfoSlim();
  component.spdxId = item.spdxId;
  component.name = item.name;
  component.version = item.version;
  component.licenseExpression = item.licenseEffective;

  licenseRuleDialog.value?.open({
    licenseId: '',
    component: component,
    policyStatus: item.policyRuleStatus,
  });
};

const openPolicyDecisionDialog = (item: TabelItem, type: DecisionType): void => {
  const component = new ComponentInfoSlim();
  component.spdxId = item.spdxId;
  component.name = item.name;
  component.version = item.version;
  component.licenseExpression = item.licenseEffective;

  let policies: PolicyRuleStatus[];
  switch (type) {
    case 'warn':
      policies = item.policyRuleStatus.filter((pr) => pr.canMakeWarnedDecision);
      break;
    case 'deny':
      policies = item.policyRuleStatus.filter(
        (pr) => pr.canMakeDeniedDecision && !pr.deniedDecisionDeniedReason?.trim(),
      );
  }

  policyDecisionDialog.value?.open({
    component,
    policies,
    type,
  });
};

const canMakeWarnedDecisionComponents = computed(() =>
  componentList.value.filter(
    (item: TabelItem) => item.version.trim() !== '' && item.policyRuleStatus.some((pr) => pr.canMakeWarnedDecision),
  ),
);

const bulkPolicyDecisionDisabled = computed(
  () => canMakeWarnedDecisionComponents.value.length === 0 || bulkPolicyDecisionDeniedReason.value.trim().length > 0,
);

const bulkPolicyDecisionTooltip = computed(() =>
  bulkPolicyDecisionDeniedReason.value.trim().length > 0
    ? t('TT_' + bulkPolicyDecisionDeniedReason.value)
    : t('TT_BULK_POLICY_DECISION'),
);

const openBulkPolicyDecisionsDialog = (): void => {
  if (canMakeWarnedDecisionComponents.value.length === 0) return;

  const items: DialogBulkPolicyDecisionEntry[] = [];
  for (const cmp of canMakeWarnedDecisionComponents.value) {
    const component = new ComponentInfoSlim();
    component.spdxId = cmp.spdxId;
    component.name = cmp.name;
    component.version = cmp.version;
    component.licenseExpression = cmp.licenseEffective;

    for (const policy of cmp.policyRuleStatus.filter((pr) => pr.canMakeWarnedDecision)) {
      const item: DialogBulkPolicyDecisionEntry = {
        component,
        policy,
      };
      items.push(item);
    }
  }

  bulkPolicyDecisionsDialog.value?.open({items});
};

const reload = async () => {
  await load();
  computeLicensesFilter();
  computePolicyTypesFilter();
  computeTypeFilter();
  computeFamilyFilter();
};

const computeLicensesFilter = () => {
  allLicenses.value = [
    ...new Set(
      componentList.value.flatMap((info: TabelItem) => {
        const license = info.licenseEffective;
        const licenseCleaned = license
          .toUpperCase()
          .replaceAll('(', ' ')
          .replaceAll(')', ' ')
          .replaceAll(' OR ', ' ')
          .replaceAll(' AND ', ' ');
        return licenseCleaned
          .split(' ')
          .map((licenseSingle) => licenseSingle.trim())
          .filter((licenseSingle) => licenseSingle.length > 0);
      }),
    ),
  ]
    .sort()
    .map((license) => ({value: license}));
};

const computeTypeFilter = () => {
  const computedTypes = [...new Set(componentList.value.map((info) => info.type))].sort();

  allTypes.value = computedTypes.map((type) => ({
    value: type,
  }));
};

const computeFamilyFilter = () => {
  const computedFamilies = [...new Set(componentList.value.map((info: TabelItem) => info.worstFamily))].sort(
    compareFamily,
  );

  allFamilies.value = computedFamilies.map((family) => ({
    value: getI18NTextOfPrefixKey('LIC_FAMILY_', family),
  }));
};

const computePolicyTypesFilter = () => {
  const computedPolicyTypes: string[] = [];
  componentList.value.forEach((info: TabelItem) => {
    if (info.policyRuleStatus) {
      info.policyRuleStatus.forEach((policyRuleStatus: PolicyRuleStatus) => {
        if (!computedPolicyTypes.includes(policyRuleStatus.type) && policyRuleStatus.used) {
          computedPolicyTypes.push(policyRuleStatus.type);
        }
      });
    }
    if (info.questioned && !computedPolicyTypes.includes('questioned')) {
      computedPolicyTypes.push('questioned');
    }
    if (info.unasserted && !computedPolicyTypes.includes('NOASSERTION')) {
      computedPolicyTypes.push('noassertion');
    }
  });

  allPolicyTypes.value = [...new Set(computedPolicyTypes)].sort(sortPolicyStatesByOrder).map((policyType) => ({
    value: policyType,
    icon: getIconForPolicyType(policyType),
    iconColor: getIconColorForPolicyType(policyType),
  }));
};

const getTableItems = (componentInfo: ComponentInfo[]): TabelItem[] =>
  componentInfo.map((info) => {
    const canMakeDecision = info.policyRuleStatus.some((pr) => pr.canMakeWarnedDecision || pr.canMakeDeniedDecision);

    return {
      ...info,
      showPolicyDecision: canMakeDecision || info.policyDecisionsApplied.length > 0,
      showLicenseDecision: Boolean(info.licenseRuleApplied) || info.canChooseLicense,
    };
  });

const load = async () => {
  dataIsLoaded.value = false;

  const {
    componentInfo,
    componentStats,
    bulkPolicyDecisionDeniedReason: bulkDeniedReason,
  } = await VersionService.getVersionComponentsForSbom(
    projectModel.value._key,
    versionDetails.value._key,
    currentSpdx.value._key,
  );

  componentList.value = componentInfo ? getTableItems(componentInfo) : [];
  stats.value = componentStats ?? new ComponentStats();
  bulkPolicyDecisionDeniedReason.value = bulkDeniedReason;

  dataIsLoaded.value = true;
};

const onRouteQueryChange = () => {
  filterOnPolicyState();
  filterOnFamilyQuery();
};

const filterOnPolicyState = () => {
  const policyState = route.query.policyFilter as PolicyState;
  if (policyState) {
    selectedFilterPolicyTypes.value = [policyState];
  }
};

const filterOnFamilyQuery = () => {
  const fam = route.query.family as string;
  if (fam === 'not declared') {
    selectedFilterFamily.value = [''];
  } else if (fam) {
    selectedFilterFamily.value = [fam];
  }
};

const formatText = (text: string): string => {
  text = escapeHtml(text);
  if (text.includes(' AND ') || text.includes(' OR ')) {
    return text
      .replace(/ AND /g, ' <strong class="db-highlight">AND</strong> ')
      .replace(/ OR /g, ' <strong class="db-highlight">OR</strong> ');
  }
  return text;
};

const ruleCallback: IRuleBtnCallbacks = {
  getUrlToComponents: (policyState: PolicyState) => {
    const spdxFile = currentSpdx.value; // Verwende den Store direkt
    let spdxAndFilter = '';
    if (spdxFile) {
      spdxAndFilter = '/' + spdxFile._key;
      if (policyState !== PolicyState.NOT_SET) {
        spdxAndFilter += `?policyFilter=${policyState}`;
      }
    }
    return `/dashboard/projects/${encodeURIComponent(projectModel.value._key)}/versions/${encodeURIComponent(
      versionDetails.value._key,
    )}/component${spdxAndFilter}`;
  },
  handlePolicySelect: (filter: PolicyState) => {
    forceReload.value = false;
    if (filter.length < 1) {
      selectedFilterPolicyTypes.value = [];
      return;
    }
    if (selectedFilterPolicyTypes.value.length === 1 && selectedFilterPolicyTypes.value[0] === filter) {
      return;
    }
    selectedFilterPolicyTypes.value = [filter];
  },
  getCountForPolicyFilterBtn: (policy: PolicyState) => {
    if (!stats.value) return 0;
    switch (policy) {
      case PolicyState.NOT_SET:
        return stats.value.Total;
      case PolicyState.DENY:
        return stats.value.Denied;
      case PolicyState.NOASSERTION:
        return stats.value.NoAssertion;
      case PolicyState.QUESTIONED:
        return stats.value.Questioned;
      case PolicyState.WARN:
        return stats.value.Warned;
      case PolicyState.ALLOW:
        return stats.value.Allowed;
      default:
        return 0;
    }
  },
  getInitSelectedPolicy: () => {
    return PolicyState.NOT_SET;
  },
  getToolTipKeyForPolicyFilterBtn: (policy: PolicyState) => {
    switch (policy) {
      case PolicyState.NOT_SET:
        return 'TT_COMPONENTS_TOTAL';
      case PolicyState.DENY:
        return 'TT_COMPONENTS_DENIED';
      case PolicyState.NOASSERTION:
        return 'TT_COMPONENTS_NOASSERTION';
      case PolicyState.QUESTIONED:
        return 'TT_COMPONENTS_QUESTIONED';
      case PolicyState.WARN:
        return 'TT_COMPONENTS_WARNED';
      case PolicyState.ALLOW:
        return 'TT_COMPONENTS_ALLOWED';
      default:
        return 'unknown_policy';
    }
  },
  getActiveClassForPolicyFilterBtn: () => {
    return '';
  },
  setRuleButtons: () => {},
};

watch(() => currentSpdx.value, reload);

watch(filteredList, () => {
  const component = filteredList.value.find((item) => item.spdxId === componentId.value);
  if (component) {
    showDetails(component);
  }
});

watch(() => route.path, onRouteQueryChange);

onMounted(async () => {
  eventBus.on('window-resize', updateTableHeight);
  filterOnPolicyState();
  filterOnFamilyQuery();
  await updateTableHeight();
  await reload();
  if (spdxFileHistory.value.length === 0) {
    return;
  }
});

onUnmounted(async () => {
  eventBus.off('window-resize', updateTableHeight);
});
</script>

<template>
  <TableLayout has-tab has-title>
    <template #buttons>
      <DRuleButtons :policies="policies" :callbacks="ruleCallback" />
      <v-spacer></v-spacer>
      <DCActionButton
        :text="t('BTN_BULK_POLICY_DECISION')"
        large
        icon="mdi-checkbox-marked-circle-plus-outline"
        :hint="bulkPolicyDecisionTooltip"
        @click="openBulkPolicyDecisionsDialog"
        :disabled="bulkPolicyDecisionDisabled"
        class="pr-4">
      </DCActionButton>
      <v-text-field
        autocomplete="off"
        variant="outlined"
        v-model="search"
        append-inner-icon="mdi-magnify"
        :label="t('labelSearch')"
        density="compact"
        clearable
        single-line
        hide-details="auto" />
    </template>
    <template #table>
      <div ref="tableComponents" class="fill-height">
        <v-data-table
          :loading="!dataIsLoaded"
          fixed-header
          :headers="filteredHeaders"
          density="compact"
          class="striped-table custom-data-table fill-height"
          :items-per-page="100"
          :footer-props="{
            'items-per-page-options': [10, 50, 100, -1],
          }"
          :height="tableHeight"
          :items="filteredList"
          :search="search"
          :custom-key-sort="customKeySort"
          v-model:sort-by="sortBy"
          @click:row="(_: Event, dataItem: DataTableItem<TabelItem>) => showDetails(dataItem.item)">
          <template v-slot:[`header.type`]="{column, getSortIcon, toggleSort}">
            <span class="mr-1">{{ column.title }}</span>
            <GridHeaderFilterIcon
              v-model="selectedFilterTypes"
              :column="column"
              :label="t('TYPE')"
              :allItems="allTypes">
            </GridHeaderFilterIcon>
            <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
          </template>
          <template v-slot:[`header.licenseEffective`]="{column, getSortIcon, toggleSort}">
            <span class="mr-1">{{ column.title }}</span>
            <GridHeaderFilterIcon
              v-model="selectedFilterLicenses"
              :column="column"
              :label="t('LICENSE')"
              :allItems="allLicenses">
            </GridHeaderFilterIcon>
            <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
          </template>
          <template v-slot:[`header.prStatus`]="{column, getSortIcon, toggleSort}">
            <HeaderSettings :column="column" :grid-name="gridName" />
            <span class="ml-6 mr-1">{{ column.title }}</span>
            <GridHeaderFilterIcon
              v-model="selectedFilterPolicyTypes"
              :column="column"
              :label="t('POLICY_STATE')"
              :allItems="allPolicyTypes">
            </GridHeaderFilterIcon>
            <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
          </template>
          <template v-slot:[`header.worstFamily`]="{column, getSortIcon, toggleSort}">
            <span class="mr-1">{{ column.title }}</span>
            <GridHeaderFilterIcon
              v-model="selectedFilterFamily"
              :column="column"
              :label="t('LICENSE_FAMILY')"
              :allItems="allFamilies">
            </GridHeaderFilterIcon>
            <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
          </template>
          <template v-slot:[`item.prStatus`]="{item}">
            <span v-if="item.unasserted">
              <Tooltip>
                {{ t(policyStateToTranslationKey('noassertion')) + ` (${item.licenseApplied})` }}
              </Tooltip>
              <v-icon small :color="getIconColorForPolicyType('noassertion')">
                {{ getIconForPolicyType('noassertion') }}
              </v-icon>
            </span>
            <span v-else-if="item.questioned">
              <Tooltip>
                {{ t('TT_questioned') + ` (${item.licenseApplied})` }}
              </Tooltip>
              <v-icon small :color="getIconColorForPolicyType('questioned')">
                {{ getIconForPolicyType('questioned') }}
              </v-icon>
            </span>
            <span v-else-if="item.prStatus == ''">
              <Tooltip>{{ t('TT_rootpackage') }}</Tooltip>
              <v-icon small :color="getIconColorForPolicyType(item.prStatus)">
                {{ getIconForPolicyType(item.prStatus) }}
              </v-icon>
            </span>
            <span v-else-if="item.policyRuleStatus">
              <Tooltip>
                <span v-for="(prStatus, index) in item.policyRuleStatus" :key="index">
                  {{ prStatus.name + ' (' + item.licenseApplied + ')' }}
                </span>
              </Tooltip>
              <v-icon small :color="getIconColorForPolicyType(item.prStatus)">
                {{ getIconForPolicyType(item.prStatus) }}
              </v-icon>
            </span>
          </template>
          <template v-slot:[`item.showPolicyDecision`]="{item}">
            <PolicyDecisionCell :item="item" @open-policy-decision="(type) => openPolicyDecisionDialog(item, type)" />
          </template>
          <template v-slot:[`item.showLicenseDecision`]="{item}">
            <span v-if="item.licenseRuleApplied">
              <v-icon size="small" :color="item.licenseRuleApplied.previewMode ? 'grey' : ''">
                {{ item.licenseRuleApplied.previewMode ? 'mdi-progress-alert' : 'mdi-information-outline' }}
              </v-icon>
              <Tooltip>
                <span class="text-subtitle-1">{{
                  item.licenseRuleApplied.previewMode
                    ? t('TT_LICENSE_RULE_APPLIED_PREVIEW')
                    : t('TT_LICENSE_RULE_APPLIED')
                }}</span>
                <br />
                <span class="d-text d-secondary-text">{{ t('TT_LICENSE_RULE_EXPRESSION') }}</span>
                <br />
                <span
                  class="d-text d-secondary-text"
                  v-html="formatText(item.licenseRuleApplied.licenseExpression)"></span>
                <br />
                <span class="d-text d-secondary-text">{{
                  t('TT_LICENSE_RULE_DECISION', {
                    decision: item.licenseRuleApplied.licenseDecisionName,
                    decisionId: item.licenseRuleApplied.licenseDecisionId,
                  })
                }}</span>
                <br />
                <span class="d-text d-secondary-text">{{
                  t('TT_LICENSE_RULE_BY_AT', {
                    creator: item.licenseRuleApplied.creator,
                    created: formatDateAndTime(item.licenseRuleApplied.created),
                  })
                }}</span>
              </Tooltip>
            </span>
            <span
              v-else-if="item.canChooseLicense && !item.choiceDeniedReason"
              @click.stop="openLicenseRuleDialog(item)">
              <v-icon size="small" color="primary">mdi-text-box-edit-outline</v-icon>
              <Tooltip>
                <span>{{ t('TT_license_rule') }}</span>
              </Tooltip>
            </span>
            <span v-else-if="item.canChooseLicense && item.choiceDeniedReason">
              <v-icon size="small" color="primary" :style="'opacity: 0.38;'">mdi-text-box-edit-outline</v-icon>
              <Tooltip>
                <span>{{ t('TT_' + item.choiceDeniedReason) }}</span>
              </Tooltip>
            </span>
          </template>
          <template v-slot:[`item.licenseEffective`]="{item}">
            <span>
              <span v-html="formatText(item.licenseEffective)"></span>
              <Tooltip>
                <div class="max-w-[500px]">
                  <div>{{ t('COL_SPDX_LICENSE_EFFECTIVE') }}:</div>
                  <div class="mr-2" v-html="formatText(item.licenseEffective)"></div>
                  <br />
                  <div>{{ t('COL_SPDX_LICENSE_DECLARED') }}:</div>
                  <div class="mr-2">{{ item.licenseDeclared }}</div>
                  <br />
                  <div>{{ t('COL_SPDX_LICENSE_CONCLUDED') }}:</div>
                  <div class="mr-2">{{ item.license }}</div>
                  <br />
                  <span>{{ item.usedPolicyRule }} (based on {{ item.licenseApplied }})</span>
                </div>
              </Tooltip>
            </span>
          </template>
          <template v-slot:[`item.worstFamily`]="{item}">
            <span>
              {{ getI18NTextOfPrefixKey('LIC_FAMILY_', item.worstFamily) }}
            </span>
          </template>
        </v-data-table>
      </div>
    </template>
  </TableLayout>

  <ComponentDetailsDialog
    ref="newComponentDetailsDlg"
    @reloadAfterCreation="reload"
    @triggerBulk="openBulkPolicyDecisionsDialog" />
  <LicenseRuleDialog ref="licenseRuleDialog" @reload="reload" />
  <PolicyDecisionDialog ref="policyDecisionDialog" @reload="reload" @triggerBulk="openBulkPolicyDecisionsDialog" />
  <BulkPolicyDecisionsDialog ref="bulkPolicyDecisionsDialog" @reload="reload" />
</template>
