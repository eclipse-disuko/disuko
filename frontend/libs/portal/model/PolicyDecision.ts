// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export interface PolicyDecisionRequest {
  sbomId: string;
  sbomName: string;
  sbomUploaded: string;
  componentSpdxId: string;
  componentName: string;
  componentVersion: string;
  licenseExpression: string;
  licenseId: string;
  policyId: string;
  policyEvaluated: string;
  policyDecision: string;
  comment: string;
  creator: string;
}

export interface PolicyDecisionSlim {
  created: string;
  componentName: string;
  componentVersion: string;
  licenseExpression: string;
  licenseId: string;
  policyId: string;
  policyEvaluated: string;
  policyDecision: string;
  creator: string;
  previewMode: boolean;
}
