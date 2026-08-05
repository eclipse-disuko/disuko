// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

/**
 * Sorts items by a given attribute using localeCompare.
 * @param items - Array to sort
 * @param key - Attribute key to sort by
 * @param transform - Optional function to transform the value before comparison (e.g., i18n translation)
 */
export const sortByAttribute = <T, K extends keyof T>(items: T[], key: K, transform?: (value: T[K]) => string): T[] =>
  items.toSorted((left, right) => {
    const leftVal = transform ? transform(left[key]) : String(left[key]);
    const rightVal = transform ? transform(right[key]) : String(right[key]);
    return leftVal.localeCompare(rightVal);
  });
