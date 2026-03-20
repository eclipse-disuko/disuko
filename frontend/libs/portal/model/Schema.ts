// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export interface ISchema {
  _key: string;
  type: number;
  content: string;
  description: string;
  created: number;
  active: boolean;
  label: string;
}

class SchemaDTO implements ISchema {
  public _key: string;
  public type: number;
  public name: string;
  public description: string;
  public content: string;
  public version: string;
  public active: boolean;
  public created: number;
  public updated: number;
  public label: string;

  public constructor() {
    this._key = '';
    this.type = 1;
    this.description = '';
    this.content = '';
    this.name = '';
    this.version = '';
    this.active = false;
    this.created = new Date().getTime();
    this.updated = new Date().getTime();
    this.label = '';
  }
}

export default class Schema extends SchemaDTO {
  constructor(dto: SchemaDTO) {
    super();
    if (dto !== null) {
      Object.assign(this, dto);
    }
  }
}
