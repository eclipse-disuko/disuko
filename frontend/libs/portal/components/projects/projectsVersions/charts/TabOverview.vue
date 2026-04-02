<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {LicenseFamily} from '@disclosure-portal/model/License';
import {PolicyState} from '@disclosure-portal/model/PolicyRule';
import {ReviewRemarkLevel, ScanRemarkLevel} from '@disclosure-portal/model/Quality';
import {
  ComponentStats,
  GeneralStats,
  LicenseFamilyStats,
  LicenseRemarkStats,
  ReviewRemarkStats,
  SbomStats,
  ScanRemarkStats,
} from '@disclosure-portal/model/VersionDetails';
import versionService from '@disclosure-portal/services/version';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useSbomStore} from '@disclosure-portal/stores/sbom.store';
import {getColor, getColorRGB} from '@disclosure-portal/utils/Tools';
import {
  ArcElement,
  BarElement,
  CategoryScale,
  ChartEvent,
  Chart as ChartJS,
  ChartOptions,
  ChartType,
  CoreChartOptions,
  Legend,
  LinearScale,
  Title,
  Tooltip,
} from 'chart.js';
import {computed, onMounted, ref, watch} from 'vue';
import {Bar, Doughnut} from 'vue-chartjs';
import {useI18n} from 'vue-i18n';
import {useRoute, useRouter} from 'vue-router';
import {applyChartDefaults} from './ChartSettings';
import {useAppStore} from '@disclosure-portal/stores/app';

ChartJS.register(Title, Tooltip, Legend, BarElement, ArcElement, CategoryScale, LinearScale);

type ChartData = {
  labels: string[];
  datasets: {
    data: (number | null)[];
    backgroundColor: string[];
    borderColor: string[];
    borderWidth: number;
    minBarLength: number;
    barPercentage: number;
  }[];
};

const {t} = useI18n();
const router = useRouter();
const route = useRoute();
const sbomStore = useSbomStore();
applyChartDefaults();

const projectModel = computed(() => useProjectStore().currentProject);
const versionDetails = computed(() => sbomStore.getCurrentVersion);
const spdxFileHistory = computed(() => sbomStore.getChannelSpdxs);
const currentSpdx = computed(() => sbomStore.getSelectedSpdx);
const sbomStats = ref<SbomStats>({} as SbomStats);
const generalStats = ref<GeneralStats>({} as GeneralStats);
const policyStateStats = ref<ComponentStats>({} as ComponentStats);
const licenseFamilyStats = ref<LicenseFamilyStats>({} as LicenseFamilyStats);
const reviewRemarksStats = ref<ReviewRemarkStats>({} as ReviewRemarkStats);
const scanRemarkStats = ref<ScanRemarkStats>({} as ScanRemarkStats);
const licenseRemarkStats = ref<LicenseRemarkStats>({} as LicenseRemarkStats);
const isLoaded = ref(false);
const legendItemsReviewRemarks = ref<any[]>([]);
const legendItemsScanRemarks = ref<any[]>([]);
const legendItemsLicenseRemarks = ref<any[]>([]);
const testColor = ref('#ffffff');
const legendItemsPolicyStates = ref<any[]>([]);
const legendItemsLicenseFamily = ref<any[]>([]);
const alternateRender = ref(useAppStore().alternateRender);

onMounted(async () => {
  testColor.value = getColor('--v-theme-licenceChartIcon');
  if (spdxFileHistory.value.length === 0) {
    return;
  }
  await loadChartData();
  isLoaded.value = true;
  updateLegends();
});

watch(
  () => currentSpdx.value,
  async () => {
    if (!currentSpdx.value) {
      return;
    }
    isLoaded.value = false;
    await loadChartData();
    isLoaded.value = true;
    updateLegends();
  },
);

const chartDataLicenseFamily = computed(() => {
  return createChartData(
    ['Other', 'NetworkCopyLeft', 'StrongCopyLeft', 'WeakCopyLeft', 'Permissive'] as (keyof LicenseFamilyStats)[],
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

function generateLegend(data: ChartData) {
  return data.labels.map((label: string, index: number) => {
    return {
      text: label,
      color: data.datasets[0].backgroundColor[index],
      borderColor: data.datasets[0].borderColor[index],
      value: data.datasets[0],
    };
  });
}

const chartOptionsLicenseFamily = createChartOptions('bar', openFamilyFilteredComponents);

const chartDataPolicyStates = computed(() => {
  return createChartData(
    ['Denied', 'NoAssertion', 'Warned', 'Questioned', 'Allowed'] as (keyof ComponentStats)[],
    'PS_CHART',
    policyStateStats.value,
    (label) => getColorForLabel('policyState', label),
    (label) => getColorForLabel('policyState', label, false),
  );
});

const chartOptionsPolicyStates = createChartOptions('bar', openFilteredComponents);

const chartDataReviewRemarks = computed(() => {
  return createChartData(
    ['NotAcceptable', 'AcceptableAfterChanges', 'Acceptable'] as (keyof ReviewRemarkStats)[],
    'RR_CHART',
    reviewRemarksStats.value,
    (label) => getColorForLabel('reviewRemark', label),
    (label) => getColorForLabel('reviewRemark', label, false),
  );
});

const chartOptionsReviewRemark = createChartOptions('doughnut', openFilteredReviewRemarks);

const chartDataScanRemarks = computed(() => {
  return createChartData(
    ['Problem', 'Warning', 'Information'] as (keyof ScanRemarkStats)[],
    'SR_CHART',
    scanRemarkStats.value,
    (label) => getColorForLabel('scanRemark', label),
    (label) => getColorForLabel('scanRemark', label, false),
  );
});

const chartOptionsScanRemark = createChartOptions('doughnut', openFilteredScanRemarks);

const chartDataLicenseRemarks = computed(() => {
  return createChartData(
    ['Alarm', 'Warning', 'Information'] as (keyof LicenseRemarkStats)[],
    'LR_CHART',
    licenseRemarkStats.value,
    (label) => getColorForLabel('licenseRemark', label),
    (label) => getColorForLabel('licenseRemark', label, false),
  );
});

const chartOptionsLicenseRemark = createChartOptions('doughnut', openLicenseRemarks);

function createPattern(color: string, patternType: string): CanvasPattern | null {
  const canvas = document.createElement('canvas');
  canvas.width = 20;
  canvas.height = 20;
  const ctx = canvas.getContext('2d');
  if (!ctx) return null;

  if (patternType === 'diagonal-stripes') {
    ctx.strokeStyle = color;
    ctx.lineWidth = 2;
    ctx.beginPath();
    ctx.moveTo(0, 0);
    ctx.lineTo(20, 20);
    ctx.stroke();
  } else if (patternType === 'dots') {
    ctx.fillStyle = color;
    ctx.beginPath();
    ctx.arc(10, 10, 4, 0, Math.PI * 2);
    ctx.fill();
  } else if (patternType === 'crosshatch') {
    ctx.strokeStyle = color;
    ctx.lineWidth = 2;
    ctx.beginPath();
    ctx.moveTo(0, 0);
    ctx.lineTo(20, 20);
    ctx.stroke();
    ctx.beginPath();
    ctx.moveTo(20, 0);
    ctx.lineTo(0, 20);
    ctx.stroke();
  }

  return ctx.createPattern(canvas, 'repeat');
}

function createChartData<T>(
  keys: (keyof T)[],
  chartLabelPrefix: string,
  stats: T,
  colorFunction: (key: string) => string,
  borderColorFunction: (key: string) => string,
) {
  const labels = keys.map((key) => t(`${chartLabelPrefix}_${key.toString().toUpperCase()}`));
  let backgroundColors = keys.map((key) => colorFunction(key.toString()));
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
      openLicenseRemarks();
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
  chartType: ChartType = 'bar',
  onClick?: (event: ChartEvent, elements: any[]) => void,
  addHoverForChartSegments = true,
) {
  let options: ChartOptions = {
    maintainAspectRatio: false,
    plugins: {
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
    };
  }

  if (chartType === 'doughnut') {
    const doptions = options as ChartOptions<'doughnut'>;
    doptions.cutout = '70%';
    return doptions;
  }
  if (chartType === 'bar') {
    const bOptions = options as ChartOptions<'bar'>;
    bOptions.indexAxis = 'y';
    bOptions.scales = {
      x: {
        grid: {
          color: getColorRGB('--v-theme-chartGrey'),
          lineWidth: 0.4,
        },
        beginAtZero: true,
        ticks: {
          font: {size: 12},
          stepSize: 40,
          color: getColorRGB('--v-theme-chartLabelColor'),
        },
      },
      y: {
        display: false,
        grid: {
          color: getColorRGB('--v-theme-chartGrey'),
          lineWidth: 0.4,
        },
        ticks: {
          font: {size: 12},
          autoSkip: false,
          color: getColorRGB('--v-theme-chartLabelColor'),
        },
      },
    };
    return bOptions;
  }
  return options;
}

function getColorForLabel(type: string, label: string, transparent = true) {
  const colorMap: {[key: string]: {[key: string]: string}} = {
    licenseFamily: {
      Other: '--v-theme-chartFLRed',
      NetworkCopyLeft: '--v-theme-chartFLRed',
      StrongCopyLeft: '--v-theme-chartFLRed',
      WeakCopyLeft: '--v-theme-chartFLYellow',
      Permissive: '--v-theme-chartFLGreen',
    },
    policyState: {
      Allowed: '--v-theme-chartGreen',
      Warned: '--v-theme-chartYellow',
      NoAssertion: '--v-theme-chartRed',
      Denied: '--v-theme-chartRed',
      Questioned: '--v-theme-chartGreen',
    },
    reviewRemark: {
      Acceptable: '--v-theme-chartGreen',
      AcceptableAfterChanges: '--v-theme-chartYellow',
      NotAcceptable: '--v-theme-chartRed',
    },
    scanRemark: {
      Information: '--v-theme-chartGrey',
      Warning: '--v-theme-chartYellow',
      Problem: '--v-theme-chartRed',
    },
    licenseRemark: {
      Information: '--v-theme-chartGrey',
      Warning: '--v-theme-chartYellow',
      Alarm: '--v-theme-chartRed',
    },
  };
  const colorVariable = (colorMap[type] && colorMap[type][label]) || '--v-theme-chartGrey';
  const color = getColor(colorVariable);
  if (alternateRender.value) {
    return `rgb(${color})`;
  } else {
    const hexColor = rgbToHex(color);
    return transparent ? addTransparencyToHex(hexColor) : `rgb(${color})`;
  }
}

function rgbToHex(rgb: string): string {
  const [r, g, b] = rgb.split(',').map((v) => parseInt(v.trim(), 10));
  return (
    '#' +
    [r, g, b]
      .map((x) => x.toString(16).padStart(2, '0'))
      .join('')
      .toUpperCase()
  );
}
function addTransparencyToHex(hex: string, alpha = 0.3) {
  const opacity = Math.round(alpha * 255)
    .toString(16)
    .padStart(2, '0');
  return `${hex}${opacity}`;
}

function openLicenseRemarks(): void {
  router.push({
    path: `/dashboard/projects/${encodeURIComponent(route.params.uuid as string)}/versions/${encodeURIComponent(
      versionDetails.value._key,
    )}/sbomQuality/${currentSpdx.value._key}/licenseRemarks`,
  });
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
  router.push({
    path: `/dashboard/projects/${encodeURIComponent(route.params.uuid as string)}/versions/${encodeURIComponent(
      versionDetails.value._key,
    )}/component/${encodeURIComponent(currentSpdx.value._key)}`,
    query: {policyFilter: mapped.toLowerCase()},
  });
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
  router.push({
    path: `/dashboard/projects/${encodeURIComponent(route.params.uuid as string)}/versions/${encodeURIComponent(
      versionDetails.value._key,
    )}/component/${encodeURIComponent(currentSpdx.value._key)}`,
    query: {family: mapped},
  });
}

function openFilteredReviewRemarks(event: ChartEvent, elements: ArcElement[]) {
  const labelMapping = {
    [t('RR_CHART_ACCEPTABLE')]: ReviewRemarkLevel.GREEN,
    [t('RR_CHART_ACCEPTABLEAFTERCHANGES')]: ReviewRemarkLevel.YELLOW,
    [t('RR_CHART_NOTACCEPTABLE')]: ReviewRemarkLevel.RED,
  };

  const mapped = resolveBarLabel<ReviewRemarkLevel>(elements, chartDataReviewRemarks.value, labelMapping);
  if (!mapped) {
    return;
  }
  router.push({
    path: `/dashboard/projects/${encodeURIComponent(route.params.uuid as string)}/versions/${encodeURIComponent(
      versionDetails.value._key,
    )}/sbomQuality/${encodeURIComponent(currentSpdx.value._key)}/reviewRemarks`,
    query: {reviewRemarkLevel: mapped},
  });
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
  router.push({
    path: `/dashboard/projects/${encodeURIComponent(route.params.uuid as string)}/versions/${encodeURIComponent(
      versionDetails.value._key,
    )}/sbomQuality/${encodeURIComponent(currentSpdx.value._key)}/scanRemarks`,
    query: {scanRemarkLevel: mapped},
  });
}

async function loadChartData() {
  const [sbomStatsData, generalStatsData] = await Promise.all([
    versionService.getSBOMStats(projectModel.value._key, versionDetails.value._key, currentSpdx.value._key),
    versionService.getGeneralVersionStats(projectModel.value._key, versionDetails.value._key),
  ]);
  sbomStats.value = sbomStatsData.data;
  generalStats.value = generalStatsData.data;

  policyStateStats.value = sbomStats.value.PolicyState;
  licenseFamilyStats.value = sbomStats.value.LicenseFamily;
  reviewRemarksStats.value = generalStats.value.ReviewRemark;
  scanRemarkStats.value = sbomStats.value.ScanRemark;
  licenseRemarkStats.value = sbomStats.value.LicenseRemark;
  legendItemsReviewRemarks.value = generateLegend(chartDataReviewRemarks.value);
  legendItemsLicenseRemarks.value = generateLegend(chartDataLicenseRemarks.value);
  legendItemsScanRemarks.value = generateLegend(chartDataScanRemarks.value);
}

function unfilteredComponentsPath() {
  return `/dashboard/projects/${encodeURIComponent(route.params.uuid as string)}/versions/${encodeURIComponent(
    versionDetails.value._key,
  )}/component/${currentSpdx.value._key}`;
}

function unfilteredScanRemarksPath() {
  return `/dashboard/projects/${encodeURIComponent(route.params.uuid as string)}/versions/${encodeURIComponent(
    versionDetails.value._key,
  )}/sbomQuality/${currentSpdx.value._key}/scanRemarks`;
}

function unfilteredLicenseRemarksPath() {
  return `/dashboard/projects/${encodeURIComponent(route.params.uuid as string)}/versions/${encodeURIComponent(
    versionDetails.value._key,
  )}/sbomQuality/${currentSpdx.value._key}/licenseRemarks`;
}

function unfilteredReviewRemarksPath() {
  return `/dashboard/projects/${encodeURIComponent(route.params.uuid as string)}/versions/${encodeURIComponent(
    versionDetails.value._key,
  )}/sbomQuality/${currentSpdx.value._key}/reviewRemarks`;
}
</script>
<template>
  <TableLayout has-title has-tab>
    <template #table>
      <v-row v-if="spdxFileHistory.length <= 0" class="m-0">
        <v-col cols="12" xs="12" align="center">
          <span>{{ t('MISSING_SBOM_UPLOAD_INFO') }}</span>
        </v-col>
      </v-row>
      <v-row v-if="spdxFileHistory.length > 0 && isLoaded" class="d-flex justify-start discoCol m-0">
        <v-col xs="12" md="6" lg="4">
          <v-card class="bar-container pa-4 card-border">
            <v-card-text>
              <ChartHeader
                :header-text="'POLICY_STATE'"
                :help-text="'PS_CHART_HELP'"
                :navigation-path="unfilteredComponentsPath()"></ChartHeader>
            </v-card-text>
            <v-card v-if="policyStateStats.Total === 0" class="empty-container text-center justify-center">
              <span>{{ t('NO_DATA') }}</span>
            </v-card>
            <v-sheet v-else color="transparent">
              <v-row>
                <v-col cols="4" class="pr-0">
                  <ul class="chart-legend">
                    <li
                      v-for="(item, index) in legendItemsPolicyStates"
                      v-bind:key="index"
                      class="legend-item"
                      @click="handleLegendClick(index, 'policyStates')">
                      <span class="text-caption">{{ item.text }}</span>
                    </li>
                  </ul>
                </v-col>
                <v-col cols="8" class="pl-0">
                  <Bar
                    :data="chartDataPolicyStates"
                    :options="chartOptionsPolicyStates as _DeepPartialObject<CoreChartOptions<'bar'>>"></Bar>
                </v-col>
              </v-row>
            </v-sheet>
          </v-card>
        </v-col>
        <v-col xs="12" sm="8" md="4" lg="3" xl="2" class="d-none d-md-block">
          <v-card class="doughnut-container pa-4 card-border">
            <v-card-text class="mb-4">
              <ChartHeader
                :header-text="'TAB_SCAN_REMARKS'"
                :help-text="'SR_CHART_HELP'"
                :navigation-path="unfilteredScanRemarksPath()"></ChartHeader>
            </v-card-text>
            <v-card-text v-if="scanRemarkStats.Total === 0" class="text-center justify-center">
              <span>{{ t('NO_DATA') }}</span>
            </v-card-text>
            <v-row v-else>
              <v-sheet class="d-flex d-inline-flex align-center pa-2" color="transparent">
                <div style="width: 90px; height: 90px">
                  <Doughnut :data="chartDataScanRemarks" :options="chartOptionsScanRemark"></Doughnut>
                </div>
                <div>
                  <ul class="chart-legend">
                    <li
                      v-for="(item, index) in legendItemsScanRemarks"
                      v-bind:key="index"
                      class="legend-item"
                      @click="handleLegendClick(index, 'scanRemarks')">
                      <span
                        class="legend-indicator"
                        :style="{backgroundColor: item.color, borderColor: item.borderColor}"></span>
                      <span class="text-caption"
                        >{{ item.text }}: {{ item.value.data[index] ? item.value.data[index] : '0' }}</span
                      >
                    </li>
                  </ul>
                </div>
              </v-sheet>
            </v-row>
          </v-card>
        </v-col>
        <v-col xs="12" sm="6" md="4" lg="3" xl="2" class="hidden-md-and-down">
          <v-card class="doughnut-container pa-4 card-border">
            <v-card-text class="mb-4">
              <ChartHeader
                :header-text="'TAB_REVIEW_REMARKS'"
                :help-text="'RR_CHART_HELP'"
                :navigation-path="unfilteredReviewRemarksPath()"></ChartHeader>
            </v-card-text>
            <v-card-text v-if="reviewRemarksStats.Total === 0" class="text-center justify-center">
              <span>{{ t('NO_DATA') }}</span>
            </v-card-text>
            <v-row v-else>
              <v-sheet class="d-flex d-inline-flex align-center pa-2" color="transparent">
                <div style="width: 85px; height: 85px">
                  <Doughnut :data="chartDataReviewRemarks" :options="chartOptionsReviewRemark"></Doughnut>
                </div>
                <div>
                  <ul class="chart-legend">
                    <li
                      v-for="(item, index) in legendItemsReviewRemarks"
                      v-bind:key="index"
                      class="legend-item"
                      @click="handleLegendClick(index, 'reviewRemarks')">
                      <span
                        class="legend-indicator"
                        :style="{backgroundColor: item.color, borderColor: item.borderColor}"></span>
                      <span class="text-caption"
                        >{{ item.text }}: {{ item.value.data[index] ? item.value.data[index] : '0' }}</span
                      >
                    </li>
                  </ul>
                </div>
              </v-sheet>
            </v-row>
          </v-card>
        </v-col>
      </v-row>
      <v-row v-if="spdxFileHistory.length > 0 && isLoaded" class="d-flex justify-start discoCol m-0">
        <v-col xs="12" md="6" lg="4">
          <v-card class="bar-container pa-4 card-border">
            <v-card-text>
              <ChartHeader
                :navigation-path="unfilteredComponentsPath()"
                :header-text="'LICENSE_FAMILY'"
                :help-text="'LF_CHART_HELP'"></ChartHeader>
            </v-card-text>
            <v-card v-if="licenseFamilyStats.Total === 0" class="empty-container text-center justify-center">
              <span>{{ t('NO_DATA') }}</span>
            </v-card>
            <v-sheet v-else color="transparent">
              <v-row>
                <v-col cols="4" class="pr-0">
                  <ul class="chart-legend">
                    <li
                      v-for="(item, index) in legendItemsLicenseFamily"
                      v-bind:key="index"
                      class="legend-item"
                      @click="handleLegendClick(index, 'licenseFamily')">
                      <span class="text-caption">{{ item.text }}</span>
                    </li>
                  </ul>
                </v-col>
                <v-col cols="8" class="pl-0">
                  <Bar :data="chartDataLicenseFamily" :options="chartOptionsLicenseFamily"></Bar>
                </v-col>
              </v-row>
            </v-sheet>
          </v-card>
        </v-col>
        <v-col xs="12" sm="8" md="4" lg="3" xl="2" class="d-none d-md-block">
          <v-card class="doughnut-container pa-4 card-border">
            <v-card-text class="mb-4">
              <ChartHeader
                :header-text="'TAB_LICENSE_REMARKS'"
                :help-text="'LR_CHART_HELP'"
                :navigation-path="unfilteredLicenseRemarksPath()"></ChartHeader>
            </v-card-text>
            <v-card-text v-if="licenseRemarkStats.Total === 0" class="text-center justify-center">
              <span>{{ t('NO_DATA') }}</span>
            </v-card-text>
            <v-row v-else>
              <v-sheet class="d-flex d-inline-flex align-center pa-2" color="transparent">
                <div style="width: 90px; height: 90px">
                  <Doughnut :data="chartDataLicenseRemarks" :options="chartOptionsLicenseRemark"></Doughnut>
                </div>
                <div>
                  <ul class="chart-legend">
                    <li
                      v-for="(item, index) in legendItemsLicenseRemarks"
                      v-bind:key="index"
                      class="legend-item"
                      @click="handleLegendClick(index, 'licenseRemarks')">
                      <span
                        class="legend-indicator"
                        :style="{backgroundColor: item.color, borderColor: item.borderColor}"></span>
                      <span class="text-caption"
                        >{{ item.text }}: {{ item.value.data[index] ? item.value.data[index] : '0' }}</span
                      >
                    </li>
                  </ul>
                </div>
              </v-sheet>
            </v-row>
          </v-card>
        </v-col>
      </v-row>
      <v-row v-if="spdxFileHistory.length > 0 && isLoaded" class="d-md-none m-0">
        <v-col xs="12" sm="6">
          <v-card class="doughnut-container pa-4 card-border">
            <v-card-text class="mb-4">
              <ChartHeader
                :header-text="'TAB_SCAN_REMARKS'"
                :help-text="'SR_CHART_HELP'"
                :navigation-path="unfilteredScanRemarksPath()"></ChartHeader>
            </v-card-text>
            <v-card-text v-if="scanRemarkStats.Total === 0" class="text-center justify-center">
              <span>{{ t('NO_DATA') }}</span>
            </v-card-text>
            <v-row v-else>
              <v-sheet class="d-flex d-inline-flex align-center pa-2" color="transparent">
                <div style="width: 90px; height: 90px">
                  <Doughnut :data="chartDataScanRemarks" :options="chartOptionsScanRemark"></Doughnut>
                </div>
                <div>
                  <ul class="chart-legend">
                    <li
                      v-for="(item, index) in legendItemsScanRemarks"
                      v-bind:key="index"
                      class="legend-item"
                      @click="handleLegendClick(index, 'scanRemarks')">
                      <span
                        class="legend-indicator"
                        :style="{backgroundColor: item.color, borderColor: item.borderColor}"></span>
                      <span class="text-caption"
                        >{{ item.text }}: {{ item.value.data[index] ? item.value.data[index] : '0' }}</span
                      >
                    </li>
                  </ul>
                </div>
              </v-sheet>
            </v-row>
          </v-card>
        </v-col>
        <v-col xs="12" sm="6">
          <v-card class="doughnut-container pa-4 card-border">
            <v-card-text class="mb-4">
              <ChartHeader
                :header-text="'TAB_LICENSE_REMARKS'"
                :help-text="'LR_CHART_HELP'"
                :navigation-path="unfilteredLicenseRemarksPath()"></ChartHeader>
            </v-card-text>
            <v-card-text v-if="licenseRemarkStats.Total === 0" class="text-center justify-center">
              <span>{{ t('NO_DATA') }}</span>
            </v-card-text>
            <v-row v-else>
              <v-sheet class="d-flex d-inline-flex align-center pa-2" color="transparent">
                <div style="width: 90px; height: 90px">
                  <Doughnut :data="chartDataLicenseRemarks" :options="chartOptionsLicenseRemark"></Doughnut>
                </div>
                <div>
                  <ul class="chart-legend">
                    <li
                      v-for="(item, index) in legendItemsLicenseRemarks"
                      v-bind:key="index"
                      class="legend-item"
                      @click="handleLegendClick(index, 'licenseRemarks')">
                      <span
                        class="legend-indicator"
                        :style="{backgroundColor: item.color, borderColor: item.borderColor}"></span>
                      <span class="text-caption"
                        >{{ item.text }}: {{ item.value.data[index] ? item.value.data[index] : '0' }}</span
                      >
                    </li>
                  </ul>
                </div>
              </v-sheet>
            </v-row>
          </v-card>
        </v-col>
      </v-row>
      <v-row v-if="spdxFileHistory.length > 0 && isLoaded" class="d-flex justify-start discoCol d-lg-none m-0">
        <v-col xs="12" md="6" lg="3" class="d-none d-md-block">
          <div style="min-height: 230px"></div>
        </v-col>
        <v-col xs="12" sm="6" md="4" lg="3" class="d-lg-none">
          <v-card class="doughnut-container pa-4 card-border">
            <v-card-text class="mb-4">
              <ChartHeader
                :header-text="'TAB_REVIEW_REMARKS'"
                :help-text="'RR_CHART_HELP'"
                :navigation-path="unfilteredReviewRemarksPath()"></ChartHeader>
            </v-card-text>
            <v-card-text v-if="reviewRemarksStats.Total === 0" class="text-center justify-center">
              <span>{{ t('NO_DATA') }}</span>
            </v-card-text>
            <v-row v-else>
              <v-sheet class="d-flex d-inline-flex align-center pa-2" color="transparent">
                <div class="w-[85px] h-[85px]">
                  <Doughnut :data="chartDataReviewRemarks" :options="chartOptionsReviewRemark"></Doughnut>
                </div>
                <div>
                  <ul class="chart-legend">
                    <li
                      v-for="(item, index) in legendItemsReviewRemarks"
                      v-bind:key="index"
                      class="legend-item"
                      @click="handleLegendClick(index, 'reviewRemarks')">
                      <span
                        class="legend-indicator"
                        :style="{backgroundColor: item.color, borderColor: item.borderColor}"></span>
                      <span class="text-caption"
                        >{{ item.text }}: {{ item.value.data[index] ? item.value.data[index] : '0' }}</span
                      >
                    </li>
                  </ul>
                </div>
              </v-sheet>
            </v-row>
          </v-card>
        </v-col>
      </v-row>
      <v-row class="justify-center m-0" v-if="spdxFileHistory.length > 0">
        <v-col v-if="!isLoaded" cols="12" lg="6" md="6" sm="12" class="text-center pt-16 px-16">
          <v-progress-circular indeterminate color="primary" size="64"></v-progress-circular>
        </v-col>
      </v-row>
    </template>
  </TableLayout>
</template>
<style scoped>
.doughnut-container {
  min-height: 236px;
}
.chart-legend {
  list-style-type: none;
  margin: 0;
  padding: 0;
}

.legend-item {
  padding-left: 14px;
  display: flex;
  align-items: center;
  margin-bottom: 4px;
  cursor: pointer;
}

.legend-indicator {
  display: inline-block;
  width: 8px;
  height: 8px;
  margin-right: 6px;
  border: 2px solid;
  border-radius: 50%;
}
</style>
