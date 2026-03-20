// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export interface IDefaultSelectItem {
  text: string;
  value: string;
}

export interface ISelectItemWithCount extends IDefaultSelectItem {
  count: number;
}
