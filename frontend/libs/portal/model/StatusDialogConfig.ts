// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export default class StatusDialogConfig {
  public title: string;
  public description: string;
  public errorDescription: string;
  public status: number;

  public constructor() {
    this.title = '';
    this.status = 0;
    this.description = '';
    this.errorDescription = '';
  }
}
