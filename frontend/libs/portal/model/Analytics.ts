// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {SearchOptions} from '@disclosure-portal/utils/Table';
import {LicenseSlim, PossibleFilterValues} from './License';

export interface IAnalyticsSearchRequest {
  component: string;
  license: string;
  exactComponent: boolean;
  exactLicense: boolean;
}

export interface SearchResponseItem {
  componentName: string;
  componentVersion: string;
  entryLicense: string;
  key: string;
  lastUpdate: string;
  licenseConcluded: string;
  licenseDeclared: string;
  name: string;
  ownerCompany: string;
  ownerDep: string;
  ownerDeptMissing: boolean;
  projectVersionKey: string;
  projectVersionName: string;
  responsible: string;
  sbomName: string;
  sbomStatus: string;
  type: string;
  ownerCompanyMissing: boolean;
  ownerCompanyName: string;
}

export interface LicenseItemMeta {
  family: string;
  approvalState: string;
  licenseType: string;
  isLicenseChart: boolean;
}
export interface SearchOccurrenciesItem {
  count: number;
  license: LicenseSlim | undefined;
  origName: string;
}

export interface OccurencesRes {
  list: SearchOccurrenciesItem[];
  possibleValues: PossibleFilterValues;
}

export interface IAnalyticsComponentsSearchResponse {
  result: string[];
}

export interface IAnalyticsLicenseSearchResponse {
  result: string[];
}

export interface IAnalyticsSearchResponse {
  result: SearchResponseItem[];
  success: boolean;
  count: number;
}

export interface IComponentsSearchResponse {
  result: string[];
}

export interface ILicensesSearchResponse {
  result: string[];
}

export interface CombinedSearchOptions {
  analyticsRequestSearchOptions: IAnalyticsSearchRequest;
  requestSearchOptions: SearchOptions;
}

export class Stats {
  public projectCount = 0;
  public projectActiveCount = 0;
  public projectDeletedCount = 0;
  public licenseActiveCount = 0;
  public licenseCount = 0;
  public licenseChartCount = 0;
  public licenseDeletedCount = 0;
  public licenseForbiddenCount = 0;
  public licenseUnknownCount = 0;
  public uploadFileCntSBOM = 0;

  public userCount = 0;
  public userActiveCount = 0;
  public userDeactivateCount = 0;
  public userTermsNotAcceptedCount = 0;

  public completedTrainings = false;
}
