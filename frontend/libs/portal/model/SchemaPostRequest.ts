// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export default class SchemaPostRequest {
  public name: string;
  public version: string;
  public description: string;
  public label: string;

  public constructor() {
    this.name = '';
    this.version = '';
    this.description = '';
    this.label = '';
  }
}
