// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export default class DHTTPError {
  public code: string;
  public title: string;
  public message: string;
  public raw: string;
  public reqId: string;

  public constructor() {
    this.code = '';
    this.title = '';
    this.message = '';
    this.raw = '';
    this.reqId = '';
  }
}
