// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export default class ExternalSourcePostRequest {
  public URL: string;
  public Comment: string;

  public constructor() {
    this.URL = '';
    this.Comment = '';
  }
}
