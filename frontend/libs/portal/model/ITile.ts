// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export default interface ITile {
  cnt: number;
  title: string;
  url: string;
  color: string;
  visible: boolean;
  icon: string;
  expand: boolean;
  expandGroup: boolean;
}
