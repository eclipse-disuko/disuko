import {getColorRGB} from '@disclosure-portal/utils/Tools';
import {Chart as ChartJS} from 'chart.js';

export function applyChartDefaults() {
  ChartJS.defaults.font.family = "'Roboto', 'sans serif'";
  ChartJS.defaults.color = getColorRGB('--v-theme-textColor');
}
