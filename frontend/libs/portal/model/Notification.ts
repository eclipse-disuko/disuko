// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export interface Notification {
  enabled: boolean;
  text: string;
  rev: string;
}

export interface NotificationDto {
  enabled: boolean;
  text: string;
}
