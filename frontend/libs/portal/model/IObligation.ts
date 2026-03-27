// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export interface IObligation {
  _key: string;
  name: string;
  nameDe: string;
  type: string;
  warnLevel: string;
  description: string;
  descriptionDe: string;
  autoApproved: boolean;
  created: string;
  updated: string;
  spdxid: string;
  remark: string;
}

export interface IObligationResponseAll {
  count: number;
  items: IObligation[];
}

export interface IDefaultSelectItem {
  text: string;
  value: string;
}

export interface ISelectItemWithCount extends IDefaultSelectItem {
  count: number;
}

export class ObligationDTO implements IObligation {
  _key: string;
  autoApproved: boolean;
  created: string;
  description: string;
  descriptionDe: string;
  name: string;
  nameDe: string;
  type: string;
  updated: string;
  warnLevel: string;
  spdxid: string;
  remark: string;
  public constructor() {
    this._key = '';
    this.autoApproved = false;
    this.created = '';
    this.description = '';
    this.descriptionDe = '';
    this.name = '';
    this.nameDe = '';
    this.type = '';
    this.updated = '';
    this.warnLevel = '';
    this.spdxid = '';
    this.remark = '';
  }
}

export default class Obligation extends ObligationDTO {
  constructor(dto: ObligationDTO) {
    super();
    Object.assign(this, dto);
  }
}
