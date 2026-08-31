// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export default interface INavItem {
  title: string;
  path: string;
  iconName: string;
  condition: boolean;
  active: boolean;
  tooltip?: string;
  subItems: INavItem[];
  externalPath?: string;
  // eslint-disable-next-line semi
}

export interface INavItemGroup {
  items: INavItem[];
  adminItem: INavItem;
}
