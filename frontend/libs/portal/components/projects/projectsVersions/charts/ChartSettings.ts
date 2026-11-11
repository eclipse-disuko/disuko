// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {getColorRGB} from '@disclosure-portal/utils/Tools';
import {Chart as ChartJS} from 'chart.js';

export const chartColorMap: {[key: string]: {[key: string]: string}} = {
  licenseFamily: {
    other: '--v-theme-chartFLRed',
    networkCopyLeft: '--v-theme-chartFLRed',
    strongCopyLeft: '--v-theme-chartFLRed',
    weakCopyLeft: '--v-theme-chartFLYellow',
    permissive: '--v-theme-chartFLGreen',
  },
  policyState: {
    allowed: '--v-theme-chartGreen',
    warned: '--v-theme-chartYellow',
    noAssertion: '--v-theme-chartRed',
    denied: '--v-theme-chartRed',
    questioned: '--v-theme-chartGreen',
  },
  reviewRemark: {
    information: '--v-theme-info',
    acceptable: '--v-theme-chartGreen',
    acceptableAfterChanges: '--v-theme-chartYellow',
    notAcceptable: '--v-theme-chartRed',
    cancelled: '--v-theme-chartGrey',
  },
  scanRemark: {
    information: '--v-theme-chartGrey',
    warning: '--v-theme-chartYellow',
    problem: '--v-theme-chartRed',
  },
  licenseRemark: {
    information: '--v-theme-chartGrey',
    warning: '--v-theme-chartYellow',
    alarm: '--v-theme-chartRed',
  },
};

export function applyChartDefaults() {
  ChartJS.defaults.font.family = "'Roboto', 'sans serif'";
  ChartJS.defaults.color = getColorRGB('--v-theme-textColor');
}

export function createPattern(color: string, patternType: string): CanvasPattern | null {
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

export function rgbToHex(rgb: string): string {
  const [r, g, b] = rgb.split(',').map((v) => parseInt(v.trim(), 10));
  return (
    '#' +
    [r, g, b]
      .map((x) => x.toString(16).padStart(2, '0'))
      .join('')
      .toUpperCase()
  );
}

export function addTransparencyToHex(hex: string, alpha = 0.3) {
  const opacity = Math.round(alpha * 255)
    .toString(16)
    .padStart(2, '0');
  return `${hex}${opacity}`;
}
