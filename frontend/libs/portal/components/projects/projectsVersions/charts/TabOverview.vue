<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {LicenseFamily} from '@disclosure-portal/model/License';
import {PolicyState} from '@disclosure-portal/model/PolicyRule';
import {ReviewRemarkLevel, ScanRemarkLevel} from '@disclosure-portal/model/Quality';
import {
  ComponentStats,
  LicenseFamilyStats,
  LicenseRemarkStats,
  ReviewRemarkStats,
  ScanRemarkStats,
} from '@disclosure-portal/model/VersionDetails';
import {useSbomStore} from '@disclosure-portal/stores/sbom.store';
import {getColor} from '@disclosure-portal/utils/Tools';
import useSnackbar from '@shared/composables/useSnackbar';
import Stack from '@shared/layouts/Stack.vue';
import {ArcElement, ChartEvent, Chart as ChartJS, ChartOptions, ChartType, Legend, Title, Tooltip} from 'chart.js';
import {computed, onMounted, ref, watch} from 'vue';
import {Doughnut} from 'vue-chartjs';
import {useI18n} from 'vue-i18n';
import {useRoute, useRouter} from 'vue-router';
import {addTransparencyToHex, applyChartDefaults, chartColorMap, createPattern, rgbToHex} from './ChartSettings';
import {useAppStore} from '@disclosure-portal/stores/app';

ChartJS.register(Title, Tooltip, Legend, ArcElement);

type ChartData = {
  labels: string[];
  datasets: {
    data: (number | null)[];
    backgroundColor: (string | CanvasPattern | null)[];
    borderColor: string[];
    borderWidth: number;
    minBarLength: number;
    barPercentage: number;
  }[];
};

type LegendItem = {
  text: string;
  color: string | CanvasPattern | null;
  borderColor: string;
  count: number;
};

type OverviewCard = {
  id: 'policyStates' | 'licenseFamily' | 'scanRemarks' | 'reviewRemarks' | 'licenseRemarks';
  type: 'bar' | 'doughnut';
  header: string;
  help: string;
  path: string;
  total: number | undefined;
  legend: LegendItem[];
  max?: number;
  data?: any;
  options?: any;
};

const {t} = useI18n();
const router = useRouter();
const route = useRoute();
const sbomStore = useSbomStore();
const snackbar = useSnackbar();
applyChartDefaults();

const versionDetails = computed(() => sbomStore.getCurrentVersion);
const spdxFileHistory = computed(() => sbomStore.getChannelSpdxs);
const currentSpdx = computed(() => sbomStore.getSelectedSBOM);
const policyStateStats = ref<ComponentStats>({} as ComponentStats);
const licenseFamilyStats = ref<LicenseFamilyStats>({} as LicenseFamilyStats);
const reviewRemarksStats = ref<ReviewRemarkStats>({} as ReviewRemarkStats);
const scanRemarkStats = ref<ScanRemarkStats>({} as ScanRemarkStats);
const licenseRemarkStats = ref<LicenseRemarkStats>({} as LicenseRemarkStats);
const isLoaded = ref(false);
const hasLoadError = ref(false);
const legendItemsReviewRemarks = ref<LegendItem[]>([]);
const legendItemsScanRemarks = ref<LegendItem[]>([]);
const legendItemsLicenseRemarks = ref<LegendItem[]>([]);
const legendItemsPolicyStates = ref<LegendItem[]>([]);
const legendItemsLicenseFamily = ref<LegendItem[]>([]);
const alternateRender = computed(() => useAppStore().alternateRender);

onMounted(async () => {
  if (spdxFileHistory.value.length === 0) {
    return;
  }
  await loadChartData();
});

watch(
  () => currentSpdx.value,
  async () => {
    if (!currentSpdx.value) {
      return;
    }
    await loadChartData();
  },
);

const chartDataLicenseFamily = computed(() => {
  return createChartData(
    ['other', 'networkCopyLeft', 'strongCopyLeft', 'weakCopyLeft', 'permissive'] as (keyof LicenseFamilyStats)[],
    'LF_CHART',
    licenseFamilyStats.value,
    (label) => getColorForLabel('licenseFamily', label),
    (label) => getColorForLabel('licenseFamily', label, false),
  );
});

function updateLegends() {
  legendItemsPolicyStates.value = generateLegend(chartDataPolicyStates.value);
  legendItemsLicenseFamily.value = generateLegend(chartDataLicenseFamily.value);
  legendItemsReviewRemarks.value = generateLegend(chartDataReviewRemarks.value);
  legendItemsScanRemarks.value = generateLegend(chartDataScanRemarks.value);
  legendItemsLicenseRemarks.value = generateLegend(chartDataLicenseRemarks.value);
}

function generateLegend(data: ChartData): LegendItem[] {
  return data.labels.map((label: string, index: number) => ({
    text: label,
    color: data.datasets[0].backgroundColor[index],
    borderColor: data.datasets[0].borderColor[index],
    count: data.datasets[0].data[index] ?? 0,
  }));
}

function maxCount(items: LegendItem[]) {
  return Math.max(0, ...items.map((i) => i.count));
}

function barStyle(item: LegendItem, max = 0) {
  const width = max > 0 && item.count > 0 ? Math.max((item.count / max) * 100, 2) : 0;
  const background = alternateRender.value
    ? `repeating-linear-gradient(45deg, ${item.borderColor} 0 2px, transparent 2px 6px)`
    : item.borderColor;
  return {width: `${width}%`, background, borderColor: item.borderColor};
}

const gridClass = 'grid w-full max-w-[80rem] grid-cols-1 gap-4 p-3 sm:p-5 md:grid-cols-2 lg:grid-cols-4';
const spanClass = (type: OverviewCard['type']) => (type === 'bar' ? 'md:col-span-2' : 'col-span-1');
const cardClass =
  'overview-card @container relative flex min-w-0 min-h-60 flex-col gap-5 overflow-hidden rounded-xl px-5 pt-4 pb-5';
const rowClass =
  'grid w-full cursor-pointer items-center rounded border-0 px-2 text-left text-inherit transition-colors duration-200 hover:bg-[rgba(var(--v-theme-on-surface),0.06)] focus-visible:outline-2 focus-visible:outline-offset-1 focus-visible:outline-[rgb(var(--v-theme-primary))] active:translate-y-px motion-reduce:transform-none motion-reduce:transition-none';

const hovered = ref<{cardId: OverviewCard['id']; index: number} | null>(null);

function setHovered(cardId: OverviewCard['id'], index: number | null) {
  hovered.value = index === null ? null : {cardId, index};
}

function hoveredItem(card: OverviewCard): LegendItem | null {
  return hovered.value?.cardId === card.id ? (card.legend[hovered.value.index] ?? null) : null;
}

function isHovered(card: OverviewCard, index: number) {
  return hovered.value?.cardId === card.id && hovered.value.index === index;
}

function countClass(count: number) {
  return [
    'text-right text-sm tabular-nums',
    count ? 'font-semibold' : 'font-normal text-[color:rgba(var(--v-theme-on-surface),0.5)]',
  ];
}

function labelClass(count: number) {
  return [
    'text-sm leading-snug break-words',
    count ? 'text-[color:rgba(var(--v-theme-on-surface),0.87)]' : 'text-[color:rgba(var(--v-theme-on-surface),0.5)]',
  ];
}

// Legends are ordered worst-first, so the first non-zero entry is the card's most severe state.
function accentColor(card: OverviewCard) {
  return card.legend.find((i) => i.count > 0)?.borderColor ?? 'transparent';
}

function dotStyle(item: LegendItem) {
  return {
    borderColor: item.borderColor,
    backgroundColor: typeof item.color === 'string' ? item.color : 'transparent',
  };
}

const cards = computed<OverviewCard[]>(() => [
  {
    id: 'policyStates',
    type: 'bar',
    header: 'POLICY_STATE',
    help: 'PS_CHART_HELP',
    path: unfilteredComponentsPath(),
    total: policyStateStats.value.total,
    legend: legendItemsPolicyStates.value,
    max: maxCount(legendItemsPolicyStates.value),
  },
  {
    id: 'scanRemarks',
    type: 'doughnut',
    header: 'TAB_SCAN_REMARKS',
    help: 'SR_CHART_HELP',
    path: unfilteredScanRemarksPath(),
    total: scanRemarkStats.value.total,
    legend: legendItemsScanRemarks.value,
    data: chartDataScanRemarks.value,
    options: chartOptionsScanRemark,
  },
  {
    id: 'reviewRemarks',
    type: 'doughnut',
    header: 'TAB_REVIEW_REMARKS',
    help: 'RR_CHART_HELP',
    path: unfilteredReviewRemarksPath(),
    total: reviewRemarksStats.value.total,
    legend: legendItemsReviewRemarks.value,
    data: chartDataReviewRemarks.value,
    options: chartOptionsReviewRemark,
  },
  {
    id: 'licenseFamily',
    type: 'bar',
    header: 'LICENSE_FAMILY',
    help: 'LF_CHART_HELP',
    path: unfilteredComponentsPath(),
    total: licenseFamilyStats.value.total,
    legend: legendItemsLicenseFamily.value,
    max: maxCount(legendItemsLicenseFamily.value),
  },
  {
    id: 'licenseRemarks',
    type: 'doughnut',
    header: 'TAB_LICENSE_REMARKS',
    help: 'LR_CHART_HELP',
    path: unfilteredLicenseRemarksPath(),
    total: licenseRemarkStats.value.total,
    legend: legendItemsLicenseRemarks.value,
    data: chartDataLicenseRemarks.value,
    options: chartOptionsLicenseRemark,
  },
]);

const chartDataPolicyStates = computed(() => {
  return createChartData(
    ['denied', 'noAssertion', 'warned', 'questioned', 'allowed'] as (keyof ComponentStats)[],
    'PS_CHART',
    policyStateStats.value,
    (label) => getColorForLabel('policyState', label),
    (label) => getColorForLabel('policyState', label, false),
  );
});

const chartDataReviewRemarks = computed(() => {
  return createChartData(
    [
      'notAcceptable',
      'acceptableAfterChanges',
      'information',
      'acceptable',
      'cancelled',
    ] as (keyof ReviewRemarkStats)[],
    'RR_CHART',
    reviewRemarksStats.value,
    (label) => getColorForLabel('reviewRemark', label),
    (label) => getColorForLabel('reviewRemark', label, false),
  );
});

const chartOptionsReviewRemark = createChartOptions('doughnut', openFilteredReviewRemarks, true, 'reviewRemarks');

const chartDataScanRemarks = computed(() => {
  return createChartData(
    ['problem', 'warning', 'information'] as (keyof ScanRemarkStats)[],
    'SR_CHART',
    scanRemarkStats.value,
    (label) => getColorForLabel('scanRemark', label),
    (label) => getColorForLabel('scanRemark', label, false),
  );
});

const chartOptionsScanRemark = createChartOptions('doughnut', openFilteredScanRemarks, true, 'scanRemarks');

const chartDataLicenseRemarks = computed(() => {
  return createChartData(
    ['alarm', 'warning', 'information'] as (keyof LicenseRemarkStats)[],
    'LR_CHART',
    licenseRemarkStats.value,
    (label) => getColorForLabel('licenseRemark', label),
    (label) => getColorForLabel('licenseRemark', label, false),
  );
});

const chartOptionsLicenseRemark = createChartOptions('doughnut', openLicenseRemarks, true, 'licenseRemarks');

function createChartData<T>(
  keys: (keyof T)[],
  chartLabelPrefix: string,
  stats: T,
  colorFunction: (key: string) => string,
  borderColorFunction: (key: string) => string,
) {
  const labels = keys.map((key) => t(`${chartLabelPrefix}_${key.toString().toUpperCase()}`));
  let backgroundColors: (string | CanvasPattern | null)[] = keys.map((key) => colorFunction(key.toString()));
  if (alternateRender.value) {
    backgroundColors = keys.map((key, index) => {
      const color = colorFunction(key.toString());
      const patternType = index % 3 === 0 ? 'diagonal-stripes' : index % 3 === 1 ? 'dots' : 'crosshatch';
      return createPattern(color, patternType);
    });
  }

  const borderColors = keys.map((key) => borderColorFunction(key.toString()));
  const data = keys.map((key) => {
    const value = (stats as any)[key] as unknown as number;
    return value === 0 ? null : value;
  });

  return {
    labels,
    datasets: [
      {
        data,
        backgroundColor: backgroundColors,
        borderColor: borderColors,
        borderWidth: 2,
        minBarLength: 10,
        barPercentage: 0.6,
      },
    ],
  };
}

function handleLegendClick(index: number, chartType: string) {
  switch (chartType) {
    case 'scanRemarks':
      openFilteredScanRemarks({} as ChartEvent, [{index} as unknown as ArcElement]);
      break;
    case 'reviewRemarks':
      openFilteredReviewRemarks({} as ChartEvent, [{index} as unknown as ArcElement]);
      break;
    case 'licenseRemarks':
      openLicenseRemarks({} as ChartEvent, [{index} as unknown as ArcElement]);
      break;
    case 'policyStates':
      openFilteredComponents({} as ChartEvent, [{index}]);
      break;
    case 'licenseFamily':
      openFamilyFilteredComponents({} as ChartEvent, [{index}]);
      break;
  }
}

function createChartOptions(
  chartType: ChartType,
  onClick?: (event: ChartEvent, elements: any[]) => void,
  addHoverForChartSegments = true,
  cardId?: OverviewCard['id'],
) {
  let options: ChartOptions = {
    maintainAspectRatio: false,
    plugins: {
      // Canvas is too small for the built-in tooltip; the hovered value is shown in the donut center instead.
      tooltip: {enabled: false},
      legend: {
        display: false,
        onHover: (event) => {
          if (event.native && event.native.target) {
            (event.native.target as HTMLElement).style.cursor = 'pointer';
          }
        },
      },
    },
    onClick,
  };

  if (addHoverForChartSegments) {
    options.onHover = (event, chartElement) => {
      if (event.native && event.native.target) {
        (event.native.target as HTMLElement).style.cursor = chartElement[0] ? 'pointer' : 'default';
      }
      if (cardId) {
        setHovered(cardId, chartElement[0]?.index ?? null);
      }
    };
  }

  if (chartType === 'doughnut') {
    const doptions = options as ChartOptions<'doughnut'>;
    doptions.cutout = '74%';
    return doptions;
  }
  return options;
}

function getColorForLabel(type: string, label: string, transparent = true) {
  const colorVariable = (chartColorMap[type] && chartColorMap[type][label]) || '--v-theme-chartGrey';
  const color = getColor(colorVariable);
  if (alternateRender.value) {
    return `rgb(${color})`;
  } else {
    const hexColor = rgbToHex(color);
    return transparent ? addTransparencyToHex(hexColor) : `rgb(${color})`;
  }
}

function openLicenseRemarks(event: ChartEvent, elements: ArcElement[]): void {
  const labelMapping = {
    [t('LR_CHART_ALARM')]: 'ALARM',
    [t('LR_CHART_WARNING')]: 'WARNING',
    [t('LR_CHART_INFORMATION')]: 'INFORMATION',
  };

  const mapped = resolveBarLabel<string>(elements, chartDataLicenseRemarks.value, labelMapping);
  navigateToPath(unfilteredLicenseRemarksPath(), mapped ? {licenseRemarkLevel: mapped} : undefined);
}

function resolveBarLabel<T extends string | (string | null)[] | null | undefined>(
  elements: any[],
  chartData: any,
  labelMapping: {[key: string]: T},
): T {
  if (elements.length === 0) {
    return null as T;
  }
  const chartElement = elements[0];
  const index = chartElement.index;
  const label = chartData.labels[index];

  return labelMapping[label];
}

function openFilteredComponents(_: ChartEvent, elements: {index: number}[]) {
  const labelMapping = {
    [t('PS_CHART_ALLOWED')]: PolicyState.ALLOW,
    [t('PS_CHART_WARNED')]: PolicyState.WARN,
    [t('PS_CHART_QUESTIONED')]: PolicyState.QUESTIONED,
    [t('PS_CHART_DENIED')]: PolicyState.DENY,
    [t('PS_CHART_NOASSERTION')]: PolicyState.NOASSERTION,
  };

  const mapped = resolveBarLabel<PolicyState>(elements, chartDataPolicyStates.value, labelMapping);
  if (!mapped) {
    return;
  }
  navigateToPath(unfilteredComponentsPath(), {policyFilter: mapped.toLowerCase()});
}

function openFamilyFilteredComponents(_: ChartEvent, elements: {index: number}[]) {
  const labelMapping = {
    [t('LF_CHART_PERMISSIVE')]: LicenseFamily.PERMISSIVE,
    [t('LF_CHART_NETWORKCOPYLEFT')]: LicenseFamily.NETWORKCOPYLEFT,
    [t('LF_CHART_STRONGCOPYLEFT')]: LicenseFamily.STRONGCOPYLEFT,
    [t('LF_CHART_WEAKCOPYLEFT')]: LicenseFamily.WEAKCOPYLEFT,
    [t('LF_CHART_OTHER')]: LicenseFamily.NOTDECLARED,
  };

  const mapped = resolveBarLabel<LicenseFamily>(elements, chartDataLicenseFamily.value, labelMapping);
  if (!mapped) {
    return;
  }
  navigateToPath(unfilteredComponentsPath(), {family: mapped});
}

function openFilteredReviewRemarks(event: ChartEvent, elements: ArcElement[]) {
  const labelMapping = {
    [t('RR_CHART_INFORMATION')]: ReviewRemarkLevel.BLUE,
    [t('RR_CHART_ACCEPTABLE')]: ReviewRemarkLevel.GREEN,
    [t('RR_CHART_ACCEPTABLEAFTERCHANGES')]: ReviewRemarkLevel.YELLOW,
    [t('RR_CHART_NOTACCEPTABLE')]: ReviewRemarkLevel.RED,
    [t('RR_CHART_CANCELLED')]: ReviewRemarkLevel.GRAY,
  };

  const mapped = resolveBarLabel<ReviewRemarkLevel>(elements, chartDataReviewRemarks.value, labelMapping);
  if (!mapped) {
    return;
  }
  navigateToPath(unfilteredReviewRemarksPath(), {reviewRemarkLevel: mapped});
}

function openFilteredScanRemarks(event: ChartEvent, elements: ArcElement[]) {
  const labelMapping = {
    [t('SR_CHART_INFORMATION')]: ScanRemarkLevel.INFORMATION,
    [t('SR_CHART_WARNING')]: ScanRemarkLevel.WARNING,
    [t('SR_CHART_PROBLEM')]: ScanRemarkLevel.PROBLEM,
  };

  const mapped = resolveBarLabel<ScanRemarkLevel>(elements, chartDataScanRemarks.value, labelMapping);
  if (!mapped) {
    return;
  }
  navigateToPath(unfilteredScanRemarksPath(), {scanRemarkLevel: mapped});
}

async function loadChartData() {
  const spdxKey = currentSpdx.value?._key;
  if (!spdxKey) {
    return;
  }
  const isStale = () => currentSpdx.value?._key !== spdxKey;
  isLoaded.value = false;
  hasLoadError.value = false;
  try {
    await Promise.all([sbomStore.fetchSBOMStats(spdxKey), sbomStore.fetchGeneralVersionStats()]);
    if (isStale()) {
      return;
    }

    const sbomStatsData = sbomStore.getSbomStats;
    const generalStatsData = sbomStore.getGeneralStats;

    policyStateStats.value = sbomStatsData.policyState;
    licenseFamilyStats.value = sbomStatsData.licenseFamily;
    reviewRemarksStats.value = generalStatsData.reviewRemark;
    scanRemarkStats.value = sbomStatsData.scanRemark;
    licenseRemarkStats.value = sbomStatsData.licenseRemark;
    updateLegends();
  } catch {
    if (isStale()) {
      return;
    }
    hasLoadError.value = true;
    snackbar.error(t('OVERVIEW_LOAD_ERROR'));
  } finally {
    if (!isStale()) {
      isLoaded.value = true;
    }
  }
}

function navigateToPath(path: string, query?: Record<string, string>) {
  if (!path) {
    return;
  }
  router.push({path, query});
}

function overviewPath(tab: 'component' | 'scanRemarks' | 'licenseRemarks' | 'reviewRemarks') {
  const spdx = currentSpdx.value;
  if (!spdx) {
    return '';
  }
  const suffix =
    tab === 'component'
      ? `component/${encodeURIComponent(spdx._key)}`
      : `sbomQuality/${encodeURIComponent(spdx._key)}/${tab}`;
  return `/dashboard/projects/${encodeURIComponent(route.params.uuid as string)}/versions/${encodeURIComponent(
    versionDetails.value._key,
  )}/${suffix}`;
}

function unfilteredComponentsPath() {
  return overviewPath('component');
}

function unfilteredScanRemarksPath() {
  return overviewPath('scanRemarks');
}

function unfilteredLicenseRemarksPath() {
  return overviewPath('licenseRemarks');
}

function unfilteredReviewRemarksPath() {
  return overviewPath('reviewRemarks');
}
</script>
<template>
  <TableLayout has-title has-tab>
    <template #table>
      <Stack
        v-if="spdxFileHistory.length <= 0 || !currentSpdx"
        align="center"
        justify="center"
        class="gap-2 px-4 py-16 text-center text-[13px] text-[color:rgba(var(--v-theme-on-surface),0.5)]">
        <v-icon icon="mdi-file-upload-outline" size="32" />
        <span>{{ t('MISSING_SBOM_UPLOAD_INFO') }}</span>
      </Stack>

      <Stack v-else-if="hasLoadError" role="alert" align="center" class="gap-4 px-4 py-16">
        <span>{{ t('OVERVIEW_LOAD_ERROR') }}</span>
        <v-btn prepend-icon="mdi-reload" variant="text" @click="loadChartData">{{ t('BTN_RELOAD') }}</v-btn>
      </Stack>

      <section v-else-if="!isLoaded" :class="gridClass" aria-busy="true">
        <div
          v-for="n in 5"
          :key="n"
          :class="[
            cardClass,
            spanClass(n === 1 || n === 4 ? 'bar' : 'doughnut'),
            'animate-pulse motion-reduce:animate-none',
          ]">
          <span class="block h-4 w-2/5 rounded bg-[rgba(var(--v-theme-on-surface),0.08)]"></span>
          <div v-if="n === 1 || n === 4" class="flex flex-col gap-4 pt-1">
            <div v-for="r in 5" :key="r" class="flex items-center gap-3">
              <span class="block h-3 w-28 rounded bg-[rgba(var(--v-theme-on-surface),0.08)]"></span>
              <span class="block h-2 flex-1 rounded-sm bg-[rgba(var(--v-theme-on-surface),0.05)]"></span>
            </div>
          </div>
          <div v-else class="flex flex-1 items-center gap-4">
            <span
              class="block size-28 flex-none rounded-full border-[0.9rem] border-solid border-[rgba(var(--v-theme-on-surface),0.08)]"></span>
            <div class="flex flex-1 flex-col gap-3">
              <span v-for="r in 3" :key="r" class="block h-3 rounded bg-[rgba(var(--v-theme-on-surface),0.08)]"></span>
            </div>
          </div>
        </div>
      </section>

      <section v-else :class="gridClass">
        <article
          v-for="(card, cardIndex) in cards"
          :key="card.id"
          :class="[
            cardClass,
            'overview-card--interactive starting:translate-y-2 starting:opacity-0',
            spanClass(card.type),
          ]"
          :style="{'--enter-delay': `${cardIndex * 50}ms`}">
          <span
            class="pointer-events-none absolute inset-x-0 top-0 h-0.5"
            :style="{background: `linear-gradient(90deg, transparent, ${accentColor(card)} 30%, transparent)`}"
            aria-hidden="true"></span>
          <ChartHeader
            :header-text="card.header"
            :help-text="card.help"
            :navigation-path="card.path"
            :total="card.type === 'bar' ? card.total : undefined" />

          <Stack
            v-if="card.total === 0"
            align="center"
            justify="center"
            class="min-h-36 flex-1 gap-2 text-sm text-[color:rgba(var(--v-theme-on-surface),0.65)]">
            <template v-if="card.type === 'doughnut'">
              <v-icon icon="mdi-check-circle-outline" size="32" color="chartGreen" />
              <span>{{ t('OVERVIEW_NO_REMARKS') }}</span>
            </template>
            <template v-else>
              <v-icon icon="mdi-chart-bar" size="32" />
              <span>{{ t('NO_DATA') }}</span>
            </template>
          </Stack>

          <ul v-else-if="card.type === 'bar'" class="m-0 flex list-none flex-col gap-0.5 p-0">
            <li v-for="(item, index) in card.legend" :key="item.text">
              <button
                type="button"
                :class="[rowClass, 'grid-cols-[minmax(0,11rem)_minmax(0,1fr)_2.5rem] gap-3 py-2']"
                :aria-label="`${item.text}: ${item.count}`"
                @click="handleLegendClick(index, card.id)">
                <span :class="labelClass(item.count)">{{ item.text }}</span>
                <span class="relative h-2 overflow-hidden rounded-sm bg-[rgba(var(--v-theme-on-surface),0.06)]">
                  <span
                    v-if="item.count > 0"
                    class="absolute inset-y-0 left-0 origin-left rounded-sm border border-solid opacity-70 transition-transform duration-500 ease-[cubic-bezier(0.2,0.8,0.2,1)] motion-reduce:transition-none starting:scale-x-0"
                    :style="barStyle(item, card.max)"></span>
                </span>
                <span :class="countClass(item.count)">{{ item.count }}</span>
              </button>
            </li>
          </ul>

          <Stack v-else direction="row" align="center" class="flex-1 gap-3 py-1 @[18rem]:gap-4">
            <div class="relative size-24 flex-none @[18rem]:size-28">
              <Doughnut :data="card.data" :options="card.options" />
              <div
                class="pointer-events-none absolute inset-0 flex flex-col items-center justify-center"
                aria-live="polite">
                <span
                  class="text-2xl leading-none font-semibold tracking-normal tabular-nums @[18rem]:text-[1.75rem]"
                  >{{ hoveredItem(card)?.count ?? card.total }}</span
                >
                <span
                  class="mt-1 max-w-16 text-center text-[11px] leading-tight tracking-normal break-words text-[color:rgba(var(--v-theme-on-surface),0.7)]"
                  >{{ hoveredItem(card)?.text ?? t('TOTAL') }}</span
                >
              </div>
            </div>
            <ul class="m-0 min-w-0 flex-1 list-none p-0">
              <li v-for="(item, index) in card.legend" :key="item.text">
                <button
                  type="button"
                  :class="[
                    rowClass,
                    'grid-cols-[8px_minmax(0,1fr)_auto] gap-2.5 py-2',
                    {'bg-[rgba(var(--v-theme-on-surface),0.08)]': isHovered(card, index)},
                  ]"
                  :aria-label="`${item.text}: ${item.count}`"
                  @mouseenter="setHovered(card.id, index)"
                  @mouseleave="setHovered(card.id, null)"
                  @focus="setHovered(card.id, index)"
                  @blur="setHovered(card.id, null)"
                  @click="handleLegendClick(index, card.id)">
                  <span
                    :class="['size-2 rounded-full border-2 border-solid', {'opacity-40': !item.count}]"
                    :style="dotStyle(item)"></span>
                  <span :class="labelClass(item.count)">{{ item.text }}</span>
                  <span :class="countClass(item.count)">{{ item.count }}</span>
                </button>
              </li>
            </ul>
          </Stack>
        </article>
      </section>
    </template>
  </TableLayout>
</template>

<style scoped>
/* Gradient hairline border: transparent border reveals the border-box gradient layer. */
.overview-card {
  border: 1px solid transparent;
  background:
    linear-gradient(180deg, rgba(var(--v-theme-on-surface), 0.03), transparent 45%) padding-box,
    linear-gradient(rgb(var(--v-theme-surface)), rgb(var(--v-theme-surface))) padding-box,
    linear-gradient(180deg, rgba(var(--v-theme-on-surface), 0.18), rgba(var(--v-theme-on-surface), 0.05)) border-box;
  box-shadow:
    inset 0 1px 0 rgba(var(--v-theme-on-surface), 0.04),
    0 1px 2px rgba(0, 0, 0, 0.18);
}

.overview-card--interactive {
  transition:
    opacity 300ms ease-out var(--enter-delay, 0ms),
    translate 300ms ease-out var(--enter-delay, 0ms),
    box-shadow 250ms ease-out;
}

.overview-card--interactive:hover,
.overview-card--interactive:focus-within {
  box-shadow:
    inset 0 1px 0 rgba(var(--v-theme-on-surface), 0.06),
    0 0 0 1px rgba(var(--v-theme-on-surface), 0.08),
    0 12px 28px -16px rgba(0, 0, 0, 0.6);
}

@media (prefers-reduced-motion: reduce) {
  .overview-card--interactive {
    transition: none;
  }
}
</style>
