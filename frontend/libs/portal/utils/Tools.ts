// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export function deepCopy<Type>(source: Type): Type {
  return JSON.parse(JSON.stringify(source));
}

export function getColor(colorVariable: string) {
  return getComputedStyle(document.documentElement).getPropertyValue(colorVariable).trim();
}

export function getColorRGB(colorVariable: string) {
  return 'rgb(' + getComputedStyle(document.documentElement).getPropertyValue(colorVariable).trim() + ')';
}

export function capitalizeFirstLetter(text: string) {
  if (!text) return '';
  return text.charAt(0).toUpperCase() + text.slice(1);
}

export function truncateText(text: string, maxLength: number) {
  if (text.length > maxLength) {
    return text.substring(0, maxLength) + '...';
  }
  return text;
}
