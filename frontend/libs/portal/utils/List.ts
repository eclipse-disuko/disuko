// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export function removeFromList<T>(list: T[], value: T) {
  const index = list.indexOf(value, 0);
  if (index > -1) {
    list.splice(index, 1);
  }
}
