// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export class Filter {
  public name = '';
  public values: string[] = [];

  public constructor(name: string, values: string[]) {
    this.name = name;
    this.values = values;
  }
}

export class FilterSetDto {
  public _key = '';
  public name = '';
  public includedFilters: Filter[] = [];
  public excludedFilters: Filter[] = [];
  public tableName = ''; // do I need this here?
}

export class FilterSetRequestDto {
  public name = '';
  public includedFilters: Filter[] = [];
  public excludedFilters: Filter[] = [];
  public tableName = '';
}
