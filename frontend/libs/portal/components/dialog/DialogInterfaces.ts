// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {Approval} from '@disclosure-portal/model/Approval';
import ErrorDialogConfig from '@shared/types/ErrorDialogConfig';
import {IObligation} from '@disclosure-portal/model/IObligation';
import Label from '@disclosure-portal/model/Label';
import PolicyRule from '@disclosure-portal/model/PolicyRule';
import {ComponentDetails, ProjectModel, UnmatchedLicense} from '@disclosure-portal/model/Project';
import {SpdxIdentifier} from '@disclosure-portal/model/Spdx';
import StatusDialogConfig from '@disclosure-portal/model/StatusDialogConfig';
import {TaskDto} from '@disclosure-portal/model/Users';
import {ComponentMultiDiff, ExternalSource, PolicyRuleStatus} from '@disclosure-portal/model/VersionDetails';

export interface ProviderPrivacyDialogInterface {
  open(): void;
}

export interface StatusDialogInterface {
  open(config: StatusDialogConfig): void;
}

export interface TermsOfUseDialogInterface {
  open(): void;
}

export interface DisabledUserDialogInterface {
  open(): void;
}

export interface TaskApprovalDialogInterface {
  open(item: TaskDto): void;
}

export interface RequestApprovalDialogInterface {
  open(projectId: string): void;
  close(): void;
}

export interface RequestExternalApprovalDialogInterface {
  open(projectId: string, vehiclePlatform: boolean): void;
  close(): void;
}

export interface AddChildrenProjectDialogInterface {
  open(projectModel: ProjectModel): void;
}

export interface SettingsDialogInterface {
  open(projectModel: ProjectModel, labels: Label[]): void;
}

export interface ErrorDialogInterface {
  open(config: ErrorDialogConfig): void;
  close(): void;
}

export interface NewSchemaDialogInterface {
  open(labels: Label[]): void;
  close(): void;
}

export interface ComponentCompareInterface {
  open(
    details: ComponentMultiDiff,
    name: string,
    spdxMetaPrevious: SpdxIdentifier,
    spdxMetaCurrent: SpdxIdentifier,
  ): void;
}

export interface ComponentDetailsInterface {
  open(
    project: ProjectModel,
    versionKey: string,
    sbomId: string,
    details: ComponentDetails,
    policyRuleStatus?: PolicyRuleStatus[],
    unmatched?: UnmatchedLicense[],
  ): void;
}

export interface NewPolicyRuleDialogInterface {
  open(labels: Label[]): void;
  edit(model: PolicyRule, labels: Label[]): void;
}

export interface LicenseDialogInterface {
  open(): void;
  close(): void;
}

export interface LicenseCompareDialogInterface {
  open(): void;
  search(licenseText: string): void;
  close(): void;
}

export interface NewObligationDialogInterface {
  open(): void;
  edit(item: IObligation): void;
  close(): void;
}

export interface TabInterface {
  open(): void;
}

export interface NewExternalSourceDialogInterface {
  open(projectKey: string, versionKey: string): void;
  edit(projectKey: string, versionKey: string, item: ExternalSource): void;
  close(): void;
}

export interface ConfigurePoliciesForLicenseDialogInterface {
  open(licenseId: string, licenseName: string): void;
  close(): void;
}

export interface DiffDialogInterface {
  open(): void;
}

export interface DialogInterface {
  open(): void;
  close(): void;
}

export interface EditApprovalReviewExternalDialogInterface {
  open(item: Approval): void;
  close(): void;
}

export interface AuditDialogInterface {
  open(id: string, name: string): void;
}

export interface SpdxTagDialogInterface {
  setOrUpdate(projectUuid: string, versionUuid: string, spdxUuid: string, spdxName: string, tag: string): void;
}
