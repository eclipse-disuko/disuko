// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import IObligation from '@disclosure-portal/model/IObligation';

interface ILicenseSlim {
  _key: string;
  licenseId: string;
  name: string;
  source: string;
  created: string;
  updated: string;
}

interface ILicense {
  _key: string;
  type: number;
  text: string;
  licenseId: string;
}

export interface ILicenseMetaSlim {
  family?: string;
  approvalState: string;
  licenseType: string;
  isLicenseChart: boolean;
  classifications: IObligation[];
  prevalentClassificationLevel: string;
}

export interface ILicenseMeta {
  approvalState: string;
  licenseType: string;
  reviewState: string;
  reviewDate: string;
  licenseUrl: string;
  sourceUrl: string;
  family?: string;
  obligationsList: IObligation[];
  changelog: string;
  evaluation: string;
  legalComments: string;
  isLicenseChart: boolean;
  classifications: IObligation[];
  prevalentClassificationLevel: string;
}

export interface ITextValue {
  value: string;
  text: string;
}

export class AliasDTO {
  public _key: string;
  public licenseId: string;
  public description: string;

  public constructor() {
    this._key = '';
    this.licenseId = '';
    this.description = '';
  }
}

class LicenseSlimDTO implements ILicenseSlim {
  public _key: string;
  public licenseId: string;
  public name: string;
  public source: string;
  public created: string;
  public updated: string;
  public meta: ILicenseMetaSlim;
  public aliases: AliasDTO[];

  public constructor() {
    this._key = '';
    this.created = '';
    this.updated = '';
    this.licenseId = '';
    this.name = '';
    this.source = 'custom';
    this.meta = {
      approvalState: '',
      licenseType: '',
      family: '',
      isLicenseChart: false,
      classifications: [],
      prevalentClassificationLevel: '',
    };
    this.aliases = [];
  }
}

class LicenseDTO implements ILicense {
  public _key: string;
  public isDeprecatedLicenseId: boolean;
  public type: number;
  public name: string;
  public text: string;
  public active: boolean;
  public licenseId: string;
  public source: string;
  public aliases: AliasDTO[];
  public meta: ILicenseMeta;
  public created: string;
  public updated: string;

  public constructor() {
    this._key = '';
    this.isDeprecatedLicenseId = false;
    this.type = 1;
    this.text = '';
    this.name = '';
    this.active = false;
    this.licenseId = '';
    this.source = 'custom';
    this.aliases = [];
    this.created = '';
    this.updated = '';
    this.meta = {
      evaluation: '',
      licenseUrl: '',
      sourceUrl: '',
      approvalState: '',
      reviewDate: '',
      reviewState: '',
      licenseType: '',
      family: '',
      changelog: '',
      obligationsList: [],
      legalComments: '',
      isLicenseChart: false,
      prevalentClassificationLevel: '',
      classifications: [],
    };
  }
}

export class LicenseSlim extends LicenseSlimDTO {
  constructor(dto: LicenseSlimDTO) {
    super();
    if (dto !== null) {
      Object.assign(this, dto);
    }
  }
}

export default class License extends LicenseDTO {
  constructor(dto: LicenseDTO) {
    super();
    if (dto !== null) {
      Object.assign(this, dto);
    }
  }
}

export class LicenseDiff {
  public licenseId = '';
  public oldLicense: License = {} as License;
  public newLicense: License = {} as License;
}

export class LicenseWithSimilarity {
  public license: License = {} as License;
  public similarity = 0.0;
}

export class ClassificationWithCount {
  public classification: IObligation = {} as IObligation;
  public count = 0;
}

export interface PossibleFilterValues {
  possibleCharts: Record<string, number>;
  possibleSources: Record<string, number>;
  possibleFamilies: Record<string, number>;
  possibleApproval: Record<string, number>;
  possibleType: Record<string, number>;
  possibleClassifications: ClassificationWithCount[];
}

interface ILicensesResponse {
  licenses: LicenseSlim[];
  count: number;
  meta: PossibleFilterValues;
}

export class LicensesResponseDTO implements ILicensesResponse {
  public licenses: LicenseSlim[];
  public count: number;
  public meta: PossibleFilterValues;

  public constructor() {
    this.licenses = [];
    this.count = 0;
    this.meta = {
      possibleApproval: {},
      possibleCharts: {},
      possibleClassifications: [],
      possibleFamilies: {},
      possibleSources: {},
      possibleType: {},
    };
  }
}

export class LicensesResponse extends LicensesResponseDTO {
  constructor(dto: LicensesResponseDTO) {
    super();
    Object.assign(this, dto);
  }
}

export function getLicenseApprovalTypeKeys(): string[] {
  return ['', 'pending', 'check', 'assigning', 'approved', 'forbidden', 'deprecated'] as string[];
}

export enum LicenseFamily {
  PERMISSIVE = 'permissive',
  WEAKCOPYLEFT = 'weak copyleft',
  STRONGCOPYLEFT = 'strong copyleft',
  NETWORKCOPYLEFT = 'network copyleft',
  NOTDECLARED = 'not declared',
}

export const familyWeight: Map<string, number> = new Map<string, number>([
  ['', -1],
  ['unknown', 0],
  [LicenseFamily.NOTDECLARED, 1],
  [LicenseFamily.NETWORKCOPYLEFT, 2],
  [LicenseFamily.STRONGCOPYLEFT, 3],
  [LicenseFamily.WEAKCOPYLEFT, 4],
  [LicenseFamily.PERMISSIVE, 5],
]);

export function compareFamily(aRaw: string, bRaw: string): number {
  const a = aRaw.toLowerCase().replace('_', ' ');
  const b = bRaw.toLowerCase().replace('_', ' ');

  const weightA = familyWeight.get(a);
  if (!weightA) {
    console.warn(`Unknown license family: ${a}`);
  }

  const weightB = familyWeight.get(b);
  if (!weightB) {
    console.warn(`Unknown license family: ${b}`);
  }

  return (weightA ?? 0) - (weightB ?? 0);
}

export interface lookupRequest {
  ids: string[];
}

export interface lookupResponse {
  items: LicenseSlimDTO[];
}
