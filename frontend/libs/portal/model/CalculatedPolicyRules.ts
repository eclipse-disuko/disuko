// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {IDefaultSelectItem} from '@disclosure-portal/model/IObligation';

export interface BucketDefinition {
  deniedClassifications: string[];
  warnedClassifications: string[];
  allowedClassifications: string[];
}

export interface CalculatedPolicyScope {
  isLicenseChart: boolean[];
  approvalState: string[];
  family: string[];
  licenseType: string[];
  source: string[];
}

export interface CalculatedPolicyConfig {
  bucketDefinition: BucketDefinition;
  licenseScope: CalculatedPolicyScope;
}

export interface ScopeFilterConfig {
  isLicenseChart: {options: IDefaultSelectItem[]; values: Array<string | boolean>};
  approvalState: {options: IDefaultSelectItem[]; values: Array<string | boolean>};
  family: {options: IDefaultSelectItem[]; values: Array<string | boolean>};
  licenseType: {options: IDefaultSelectItem[]; values: Array<string | boolean>};
  source: {options: IDefaultSelectItem[]; values: Array<string | boolean>};
}

export interface CalculatedRuleConfigType {
  calculated: boolean;
  buckets: BucketDefinition;
  classificationOptions: IDefaultSelectItem[];
  scopeConfig: ScopeFilterConfig;
}
