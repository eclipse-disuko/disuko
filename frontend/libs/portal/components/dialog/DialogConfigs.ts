// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {ProjectModel} from '@disclosure-portal/model/Project';
import {ProjectSlim} from '@disclosure-portal/model/ProjectsResponse';
import {LicenseMeta, ReviewRemark} from '@disclosure-portal/model/Quality';
import {ComponentInfoSlim, PolicyRuleStatus, VersionSlim} from '@disclosure-portal/model/VersionDetails';

export interface DialogVersionFormConfig {
  projectID: string;
  version?: VersionSlim;
}
export interface DialogReviewRemarkConfig {
  presetItem?: ReviewRemark;

  versionID?: string;
  spdxID?: string;

  components?: ComponentInfoSlim[];
  licenses?: LicenseMeta[];
}

export interface DialogEditProjectConfig {
  projectSlim?: ProjectSlim;
  project?: ProjectModel;
}

export interface DialogLicenseRuleConfig {
  licenseId: string;
  component: ComponentInfoSlim;
  policyStatus?: PolicyRuleStatus[];
}

export type DecisionType = 'warn' | 'deny';

export interface DialogPolicyDecisionConfig {
  component: ComponentInfoSlim;
  policies: PolicyRuleStatus[];
  type: DecisionType;
}

export interface DialogBulkPolicyDecisionEntry {
  component: ComponentInfoSlim;
  policy: PolicyRuleStatus;
}

export interface DialogBulkPolicyDecisionsConfig {
  items: DialogBulkPolicyDecisionEntry[];
}
