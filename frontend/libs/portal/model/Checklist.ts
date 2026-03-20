// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export class Checklist {
  _key = '';
  updated?: string;
  created?: string;
  name = '';
  nameDE = '';
  description = '';
  descriptionDE = '';
  policyLabels: string[] = [];
  items: ChecklistItem[] = [];
  active: boolean = false;
}

export class ChecklistItem {
  _key = '';
  updated?: string;
  created?: string;
  name = '';
  triggerType = '';
  policyStatus: string[] = [];
  policyLabels: string[] = [];
  classifications: string[] = [];
  scanRemarks = '';
  licenseIds: string[] = [];
  componentNames: string[] = [];
  targetTemplateName = '';
  targetTemplateKey = '';
}

export enum TriggerTypes {
  DEFAULT = 'DEFAULT',
  CLASS_AND = 'CLASS_AND',
  CLASS_OR = 'CLASS_OR',
  POLICY_STATUS = 'POLICY_STATUS',
  SCAN_REMARK = 'SCAN_REMARK',
  LICENSE = 'LICENSE',
  COMPONENT_NAME = 'COMPONENT_NAME',
}

export enum PolicyStatusTypes {
  ALLOWED = 'ALLOWED',
  DENIED = 'DENIED',
  UNASSERTED = 'UNASSERTED',
  WARNED = 'WARNED',
  QUESTIONED = 'QUESTIONED',
}

export class ExecuteRequest {
  ids: string[] = [];
}
