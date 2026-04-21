// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {DocumentDto} from '@disclosure-portal/model/Document';
import {ApprovableSPDXDto} from '@disclosure-portal/model/Project';
import {ComponentStats} from './VersionDetails';

export class ProjectApprovable {
  public projectKey = '';
  public projectName = '';
  public spdxname = '';
  public spdxtag = '';
  public spdxUploaded = '';
  public isSpdxRecent = '';
  public customerdiff = false;
  public supplierdiff = false;
  public stats: ComponentStats = {} as ComponentStats;
  public approvablespdx: ApprovableSPDXDto = {} as ApprovableSPDXDto;
}

export class ApprovableInfo {
  public stats: ComponentStats = {} as ComponentStats;
  public projects: ProjectApprovable[] = [];
  public hasDeniedDecisions = false;
}

export enum ApprovalStates {
  Pending = 'PENDING',
  Declined = 'DECLINED',
  Approved = 'APPROVED',
  SupplierApproved = 'SUPPLIER_APPROVED',
  CustomerApproved = 'CUSTOMER_APPROVED',
  Aborted = 'ABORTED',
  Generating = 'GENERATING',
  GenerationFailed = 'GENERATION_FAILED',
}

export enum ApproverRoles {
  Supplier1 = 0,
  Supplier2 = 1,
  Customer1 = 2,
  Customer2 = 3,
}

export class ApproveState {
  public updated = '';
  public state: ApprovalStates = ApprovalStates.Pending;
}

export class InternalApproval {
  public states: ApproveState[] = [];
  public approver: string[] = [];
  public approverFullName: string[] = [];
  public comments: string[] = [];
  public docVersion = 0;
  public aborted = false;
}

export class ExternalApproval {
  public vehicle = false;
  public state: ApprovalStates = ApprovalStates.Pending;
  public comment = '';
}

export class Plausibility {
  public state: ApproveState = {state: ApprovalStates.Pending, updated: ''};
  public comment = '';
  public approver = '';
  public approverFullName = '';
}

export enum ApprovalType {
  Internal = 'INTERNAL',
  Plausibility = 'PLAUSIBILITY',
  External = 'EXTERNAL',
}

export class DocumentFlags {
  public c1 = false;
  public c2 = false;
  public c3 = false;
  public c4 = false;
  public c5 = false;
  public c6 = false;
}

export class Approval {
  public key = '';
  public created = '';
  public updated = '';

  public projectKey = '';
  public creator = '';
  public creatorFullName = '';
  public comment = '';

  public info: ApprovableInfo = {} as ApprovableInfo;
  public documents: DocumentDto[] = [];
  public flags: DocumentFlags = {} as DocumentFlags;

  public type: ApprovalType = ApprovalType.Internal;

  public internal: InternalApproval = {} as InternalApproval;
  public plausibility: Plausibility = {} as Plausibility;
  public external: ExternalApproval = {} as ExternalApproval;

  public status: ApprovalStates = ApprovalStates.Pending;
}

export enum PowerOfAttorneyType {
  iV = 'iV',
  ppa = 'ppa',
  other = 'other',
}

export class ApprovalUpdate {
  public state = ApprovalStates.Pending;
  public comment = '';

  public powerOfAttorney!: PowerOfAttorneyType;
}

export class CreateApprovalResponse {
  public approvalGuid = '';
  public success = false;
  public jobKey = '';
}
