// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

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
