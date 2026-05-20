// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {CustomId} from '@disclosure-portal/model/CustomId';
import {Department} from '@disclosure-portal/model/Department';
import type Label from '@disclosure-portal/model/Label';
import {Project} from '@disclosure-portal/model/Project';
import {UserDto} from '@shared/types/Users';
import {SupplierExtraData} from '@disclosure-portal/model/Wizard';

export const stepIds = {
  platform: 'platform',
  details: 'details',
  architecture: 'architecture',
  targetUsers: 'targetUsers',
  distributionTarget: 'distributionTarget',
  development: 'development',
  owner: 'owner',
  developer: 'developer',
  summary: 'summary',
} as const;

export type StepId = (typeof stepIds)[keyof typeof stepIds];

export interface StepType {
  id: StepId;
  i18nKey: string; // i18nKey as string
  index: number;
  isCompleted: boolean;
  errorText?: string;
  seen: boolean;
}

export const targetPlatforms = {
  enterprise: 'Enterprise IT',
  mobile: 'Mobile',
  vehicle: 'Product',
  other: 'Other',
} as const;
export type TargetPlatform = (typeof targetPlatforms)[keyof typeof targetPlatforms] | null;

export const architectures = {
  frontend: 'Frontend or Client',
  backend: 'Backend',
  vehicleOnboard: 'Product onboard',
  vehicleOffboard: 'Product offboard',
  none: 'None',
} as const;
export type Architecture = (typeof architectures)[keyof typeof architectures] | null;

export const targetUsers = {
  company: 'Company',
  businessPartner: 'Business Partner',
  customer: 'End Customer',
} as const;
export type TargetUsers = (typeof targetUsers)[keyof typeof targetUsers] | null;

export const distributionTargets = {
  company: 'Company',
  businessPartner: 'Business Partner',
} as const;
export type DistributionTarget = (typeof distributionTargets)[keyof typeof distributionTargets] | null;

export const developments = {
  inhouse: 'In-house',
  internal: 'Internal Developer',
  external: 'External Developer',
};
export type Development = (typeof developments)[keyof typeof developments] | null;

export interface DisclosureDocumentMeta {
  supplierName: string;
  supplierNr: string;
  supplierAddress: string;
  supplierDept: Department;
  deptMissing: boolean;
}

export interface CustomerMeta {
  dept: Department;
  deptMissing: boolean;
  address: string;
  fRI: string;
  sRI: string;
  acc: string;
  userFRI: UserDto;
  userSRI: UserDto;
}

export interface NoticeContactMeta {
  address: string;
  email: string;
}

export interface ApplicationMeta {
  id: string;
  secondaryId: string;
  name: string;
  externalLink: string;
}

export interface ProjectSettings {
  /**
   * Developer
   */
  documentMeta: DisclosureDocumentMeta;
  /**
   * Owner
   */
  customerMeta: CustomerMeta;
  /**
   * Third party notice address
   */
  noticeContactMeta: NoticeContactMeta;
  supplierExtraData: SupplierExtraData;
  projectKey?: string;
  noFossProject?: boolean;
  customIds?: CustomId[];
}

export interface WizardProject {
  name: string;
  description: string;
  targetPlatform: TargetPlatform;
  architecture: Architecture;
  distributionTarget: DistributionTarget;
  targetUsers: TargetUsers;
  development: Development;
  owner: string;
  id: string;
  isDummy?: boolean;
  isGroup: boolean;
  parentKey: string;
  /**
   * Owner and Developer Settings
   */
  projectSettings: ProjectSettings;
  /**
   * Application Reference (Details Step)
   */
  applicationMeta: ApplicationMeta | null;
  labels?: Label['_key'][];
  hasDeniedDecisions: boolean;
}

export function createWizardProjectFromProject(project: Project): WizardProject {
  return {
    id: project._key,
    name: project.name,
    description: project.description,
    isDummy: project.isDummy,
    isGroup: project.isGroup,
    parentKey: project.parent,
    targetPlatform: null,
    architecture: null,
    distributionTarget: null,
    targetUsers: null,
    development: null,
    owner: '',
    // Preserve existing project settings if available
    projectSettings: {
      documentMeta: project.documentMeta,
      customerMeta: project.customerMeta,
      noticeContactMeta: project.noticeContactMeta,
      supplierExtraData: project.supplierExtraData,
    },
    applicationMeta: project.applicationMeta,
  };
}
