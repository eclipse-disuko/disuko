// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export default class IdleInfo {
  public message: string;
  public progressUnit: string;
  public show: boolean;
  public progress: number;
  public constructor(showIdle: boolean) {
    this.message = '';
    this.progress = -1;
    this.show = showIdle;
    this.progressUnit = '%';
  }
}

export interface INotificationMeta {
  enabled: boolean;
  text: string;
}
