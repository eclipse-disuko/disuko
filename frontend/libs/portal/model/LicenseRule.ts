// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export interface LicenseRuleRequest {
  sbomId: string;
  sbomName: string;
  sbomUploaded: string;
  componentSpdxId: string;
  componentName: string;
  componentVersion: string;
  licenseExpression: string;
  licenseDecisionId: string;
  licenseDecisionName: string;
  comment: string;
  creator: string;
  active?: boolean;
}

export interface LicenseRuleSlim {
  created: string;
  componentName: string;
  licenseExpression: string;
  licenseDecisionId: string;
  licenseDecisionName: string;
  creator: string;
  previewMode: boolean;
}
