// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export const sortByAttribute = <T, K extends keyof T>(items: T[], key: K): T[] =>
  items.toSorted((left, right) => String(left[key]).localeCompare(String(right[key])));