// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export class ApprovalStats {
  public total = 0;
  public allowed = 0;
  public warned = 0;
  public denied = 0;
  public questioned = 0;
  public noAssertion = 0;
}

export class DocumentMeta {
  public c1 = false;
  public c2 = false;
  public c3 = false;
  public c4 = false;
  public c5 = false;
  public c6 = false;
}

export interface InternalApprovalRequest {
  customerApprover1: string;
  customerApprover2: string;
  supplierApprover1: string;
  supplierApprover2: string;
  comment: string;
  guidProject: string;
  metaDoc: DocumentMeta;
  withZip: boolean;
  fossVersion: 'default' | 'legacy' | 'vehicle-legacy' | 'vanilla';
}

export interface ExternalApprovalRequest {
  comment: string;
  guidProject: string;
  metaDoc: DocumentMeta;
  withZip: boolean;
  fossVersion: 'default' | 'legacy' | 'vehicle-legacy' | 'vanilla';
  selectedProjects: string[];
}

export class PlausibilityCheckRequest {
  public comment = '';
  public guidProject = '';
  public metaDoc = new DocumentMeta();
  public approver = '';
}

export class ApprovalUpdate {
  public customerApprover1 = '';
  public customerApprover2 = '';
  public supplierApprover1 = '';
  public supplierApprover2 = '';
  public guidProject = '';
  public guidApproval = '';
  public requestUser = '';
  public guidVersion = '';
  public guidSBOM = '';
  public accepted = false;
  public comment = '';
}

export interface ApprovalResponse {
  success: boolean;
}
