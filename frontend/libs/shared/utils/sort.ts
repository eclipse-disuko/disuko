// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export const sortByText = <T extends {text: string}>(items: T[]): T[] =>
  items.toSorted((left, right) => left.text.localeCompare(right.text));
