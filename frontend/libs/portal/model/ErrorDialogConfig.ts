// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export default class ErrorDialogConfig {
  public title: string;
  public titleKeyOrCode: string;
  public description: string;
  public errorCode: string;
  public reqId: string;
  public copyDesc: boolean;
  public stackTrace: any;

  public constructor() {
    this.title = '';
    this.titleKeyOrCode = '';
    this.description = '';
    this.errorCode = '';
    this.stackTrace = '';
    this.reqId = '';
    this.copyDesc = false;
  }
}
