// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {getColorRGB} from '@disclosure-portal/utils/Tools';
import {Chart as ChartJS} from 'chart.js';

export function applyChartDefaults() {
  ChartJS.defaults.font.family = "'Roboto', 'sans serif'";
  ChartJS.defaults.color = getColorRGB('--v-theme-textColor');
}
