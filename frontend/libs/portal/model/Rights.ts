// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {ActionRights} from '@disclosure-portal/model/Credentials';

export interface IRights {
  read: boolean;
  update: boolean;
  create: boolean;
  delete: boolean;
}

export class CRUDRights implements IRights {
  public read = true;
  public update = false;
  public create = false;
  public delete = false;
}

export enum Group {
  UserInternal = 'Internal',
  UserNonInternal = 'NonInternal',
  UserLicenseManager = 'LicenseManager',
  UserPolicyManager = 'PolicyManager',
  UserProjectAnalyst = 'ProjectAnalyst',
  UserDomainAdmin = 'DomainAdmin',
  UserApplicationAdmin = 'ApplicationAdmin',
  UserFOSSOffice = 'FOSSOffice',
  ProjectOwner = 'Owner',
  ProjectViewer = 'Viewer',
  ProjectSupplier = 'Supplier',
}

export class Rights {
  public allowTools!: CRUDRights;
  public allowStyleguide!: CRUDRights;
  public allowS3Tests!: CRUDRights;
  public allowSampleData!: CRUDRights;
  public allowUsers!: CRUDRights;
  public allowSchema!: CRUDRights;
  public allowLabel!: CRUDRights;
  public allowPolicy!: CRUDRights;
  public allowObligation!: CRUDRights;
  public allowLicense!: CRUDRights;
  public allowProject!: CRUDRights;
  public allowProjectGroup!: CRUDRights;
  public allowRequestApproval!: CRUDRights;
  public allowRequestPlausi!: CRUDRights;
  public allowProjectVersion!: CRUDRights;
  public allowDisclosureDocument!: ActionRights;
  public allowAllProjectUserManagement!: CRUDRights;
  public allowAllProjectTokenManagement!: CRUDRights;
  public allowProjectAudit!: CRUDRights;
  public allowProjectUserManagement!: CRUDRights;
  public allowProjectTokenManagement!: CRUDRights;
  public allowSBOMAction!: ActionRights;
  public allowCCSAction!: ActionRights;
  public allowReviewTemplates!: CRUDRights;
  public allowExecuteChecklist!: boolean;
  public isInternal = false;
  public groups!: Group[];

  allgroups() {
    return `${this.groups.join(', ')}`;
  }

  isDomainAdmin() {
    return this.groups.includes(Group.UserDomainAdmin);
  }

  isApplicationAdmin() {
    return this.groups.includes(Group.UserApplicationAdmin);
  }

  isProjectAnalyst() {
    return this.groups.includes(Group.UserProjectAnalyst);
  }

  isFOSSOffice() {
    return this.groups.includes(Group.UserFOSSOffice);
  }

  isAnyOfAdmin() {
    return (
      this.hasClassificationsAccess() ||
      this.hasPolicyAccess() ||
      this.hasAllProjectsReadonly() ||
      this.hasLabelAccess() ||
      this.hasSchemaAccess() ||
      this.hasToolsAccess() ||
      this.hasSampleDataAccess() ||
      this.hasStyleguideAccess() ||
      this.hasUsersAccess()
    );
  }

  hasClassificationsAccess() {
    return (
      this.allowObligation &&
      (this.allowObligation.create || this.allowObligation.update || this.allowObligation.delete)
    );
  }

  hasReviewTemplatesAcces() {
    return (
      this.allowReviewTemplates &&
      (this.allowReviewTemplates.create || this.allowReviewTemplates.update || this.allowReviewTemplates.delete)
    );
  }

  hasPolicyAccess() {
    return this.allowPolicy && (this.allowPolicy.create || this.allowPolicy.update || this.allowPolicy.delete);
  }

  hasAllProjectsReadonly() {
    return this.allowProject && this.allowProject.read;
  }

  hasLabelAccess() {
    return this.allowLabel && (this.allowLabel.create || this.allowLabel.update || this.allowLabel.delete);
  }

  hasSchemaAccess() {
    return this.allowSchema && (this.allowSchema.create || this.allowSchema.update || this.allowSchema.delete);
  }

  hasToolsAccess() {
    return this.allowTools && (this.allowTools.create || this.allowTools.update || this.allowTools.delete);
  }

  hasSampleDataAccess() {
    return (
      this.allowSampleData &&
      (this.allowSampleData.create || this.allowSampleData.update || this.allowSampleData.delete)
    );
  }

  hasStyleguideAccess() {
    return (
      this.allowStyleguide &&
      (this.allowStyleguide.create || this.allowStyleguide.update || this.allowStyleguide.delete)
    );
  }

  hasUsersAccess() {
    return (
      this.allowUsers &&
      this.allowUsers.create &&
      this.allowUsers.read &&
      this.allowUsers.update &&
      this.allowUsers.delete
    );
  }

  hasLicenseAccess() {
    return this.allowLicense && (this.allowLicense.create || this.allowLicense.update || this.allowLicense.delete);
  }
}
