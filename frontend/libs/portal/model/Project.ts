// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {CustomId} from '@disclosure-portal/model/CustomId';
import {Department} from '@disclosure-portal/model/Department';
import {LicenseRuleSlim} from '@disclosure-portal/model/LicenseRule';
import {PolicyDecisionSlim} from '@disclosure-portal/model/PolicyDecision';
import ProjectPostRequest from '@disclosure-portal/model/ProjectPostRequest';
import {ProjectChildren} from '@disclosure-portal/model/ProjectsResponse';
import {Group, Rights} from '@disclosure-portal/model/Rights';
import {UserDto} from '@disclosure-portal/model/Users';
import {SupplierExtraData} from '@disclosure-portal/model/Wizard';
import {IMap} from '@disclosure-portal/utils/View';
import {Application} from './Application';
import License from './License';
import Schema from './Schema';
import {ComponentStats, PolicyRuleStatus, VersionSlim} from './VersionDetails';

export interface ICreated {
  Created: string;
}

export interface ICreatedSmall {
  created: string;
}

export interface IUploaded {
  Uploaded: string;
}

export class CreateProjectsResponse {
  id = '';
  groupId = '';
  name = '';
  taskGuid = '';
}

export class Token {
  public _key = '';
  public company = '';
  public description = '';
  public status = '';
  public created = '';
  public expiry = '';
  public tokenSecret = '';
}

export class UnmatchedLicense {
  public orig = '';
  public referenced = '';
  public known = false;
}

export enum UserType {
  OWNER = 'Owner',
  SUPPLIER = 'Supplier',
  VIEWER = 'Viewer',
}

export class ProjectUser {
  public userId = '';
  public userType = '';
  public created = '';
  public userProfile = {} as UserDto;
  public responsible = false;
  public comment = '';
}

export class UserManagement {
  public users: ProjectUser[] = [];
}

export interface IHash {
  [key: string]: VersionSlim;
}

export class DisclosureDocumentMetaDTO {
  public supplierName = '';
  public supplierNr = '';
  public supplierAddress = '';
  public supplierDept = {} as Department;
  public deptMissing = false;
}

export class DisclosureDocumentMeta extends DisclosureDocumentMetaDTO {
  constructor(dto: DisclosureDocumentMetaDTO | null) {
    super();
    if (dto !== null) {
      Object.assign(this, dto);
    }
  }

  public fill(dto: DisclosureDocumentMetaDTO | null) {
    if (dto !== null) {
      Object.assign(this, dto);
    }
  }
}

export class Company {
  public name = '';
  public code = '';
}

export class CustomerMetaDTO {
  public dept = {} as Department;
  public deptMissing = false;
  public address = '';
  public fRI = '';
  public sRI = '';
  public acc = '';
  public userFRI!: UserDto;
  public userSRI!: UserDto;

  public fill(dto: CustomerMetaDTO | null) {
    if (dto !== null) {
      Object.assign(this, dto);
    }
  }
}

export class NoticeContactMetaDTO {
  public address = '';
  public email = '';

  public fill(dto: NoticeContactMetaDTO | null) {
    if (dto !== null) {
      Object.assign(this, dto);
    }
  }
}

export class ProjectSettingsModel {
  public documentMeta = new DisclosureDocumentMeta(null);
  public customerMeta = new CustomerMetaDTO();
  public noticeContactMeta = new NoticeContactMetaDTO();
  public supplierExtraData = new SupplierExtraData();
  public projectKey = '';
  public noFossProject = false;
  public customIds = [] as CustomId[];

  public fill(projectModel: ProjectModel) {
    this.projectKey = projectModel._key;
    this.documentMeta.fill(projectModel.documentMeta);
    this.customerMeta.fill(projectModel.customerMeta);
    this.noFossProject = projectModel.isNoFoss;

    if (projectModel.documentMeta.supplierDept && projectModel.documentMeta.supplierDept.deptId) {
      this.documentMeta.supplierDept = new Department();
      this.documentMeta.supplierDept.fill(projectModel.documentMeta.supplierDept);
    }
    if (projectModel.customerMeta.dept && projectModel.customerMeta.dept.deptId) {
      this.customerMeta.dept = new Department();
      this.customerMeta.dept.fill(projectModel.customerMeta.dept);
    }

    this.noticeContactMeta.fill(projectModel.noticeContactMeta);
    this.supplierExtraData = projectModel.supplierExtraData;
    this.customIds = [...(projectModel.customIds || [])];
    // this.applicationMeta.fill(projectModel.applicationMeta);
  }

  fill2(projectModel: ProjectPostRequest) {
    if (projectModel.projectSettings == null) {
      return;
    }
    this.documentMeta.fill(projectModel.projectSettings.documentMeta);
    this.customerMeta.fill(projectModel.projectSettings.customerMeta);
    this.documentMeta.supplierDept = new Department();
    this.documentMeta.supplierDept.fill(projectModel.projectSettings.documentMeta.supplierDept);
    this.customerMeta.dept = new Department();
    this.customerMeta.dept.fill(projectModel.projectSettings.customerMeta.dept);
    this.noticeContactMeta.fill(projectModel.projectSettings.noticeContactMeta);
    this.supplierExtraData = projectModel.projectSettings.supplierExtraData;
    this.noFossProject = projectModel.projectSettings.noFossProject;
  }

  fill3(projectSettings: ProjectSettingsModel) {
    this.projectKey = projectSettings.projectKey;
    this.documentMeta.fill(projectSettings.documentMeta);
    this.customerMeta.fill(projectSettings.customerMeta);
    this.noticeContactMeta.fill(projectSettings.noticeContactMeta);
    this.supplierExtraData = projectSettings.supplierExtraData;
  }

  handleExternalSupplierData() {
    if (this.supplierExtraData?.external) {
      this.documentMeta.supplierDept = {} as Department;
      this.documentMeta.supplierNr = '';
      this.supplierExtraData.fRI = '';
      this.supplierExtraData.sRI = '';
      this.supplierExtraData.userFRI = null;
      this.supplierExtraData.userSRI = null;
    }
  }
}

export class ParentProjectSettingsModel {
  public projectKey = '';
  public documentMeta = new DisclosureDocumentMeta(null);
  public customerMeta = new CustomerMetaDTO();
  public supplierExtraData = new SupplierExtraData();
  public noticeContactMeta = new NoticeContactMetaDTO();

  public fill(projectKey: string, parentProjectSettingsModel: ParentProjectSettingsModel) {
    this.projectKey = projectKey;
    this.documentMeta.fill(parentProjectSettingsModel.documentMeta);
    this.customerMeta.fill(parentProjectSettingsModel.customerMeta);

    if (
      parentProjectSettingsModel.documentMeta.supplierDept &&
      parentProjectSettingsModel.documentMeta.supplierDept.deptId
    ) {
      this.documentMeta.supplierDept = new Department();
      this.documentMeta.supplierDept.fill(parentProjectSettingsModel.documentMeta.supplierDept);
    }
    if (parentProjectSettingsModel.customerMeta.dept && parentProjectSettingsModel.customerMeta.dept.deptId) {
      this.customerMeta.dept = new Department();
      this.customerMeta.dept.fill(parentProjectSettingsModel.customerMeta.dept);
    }

    this.supplierExtraData = parentProjectSettingsModel.supplierExtraData;
    this.noticeContactMeta.fill(parentProjectSettingsModel.noticeContactMeta);
  }
}

export class ApprovableSPDXDto {
  public versionName = '';
  public versionkey = '';
  public spdxkey = '';
}

export class ProjectSubscriptions {
  public spdx = false;
  public overallReview = false;
}

export class ProjectDTO {
  public _key: string;
  public name: string;
  public description: string;
  public schemaLabel: string;
  public policyLabels: string[];
  public projectLabels: string[];
  public children: string[];
  public freeLabels: string[];
  public applicationId: string;
  public versions: IHash = {};
  public created: string;
  public updated: string;
  public correspondingSchema: Schema = {} as Schema;
  public token: Token[];
  public documentMeta: DisclosureDocumentMeta;
  public customerMeta: CustomerMetaDTO;
  public noticeContactMeta: NoticeContactMetaDTO;
  public accessRights: Rights;
  public parentName: string;
  public parent: string;
  public status: string;
  public isGroup = false;
  public approvablespdx = {} as ApprovableSPDXDto;
  public parentProjectSettings = {} as ParentProjectSettingsModel;
  public isNoFoss = false;
  public supplierMissing = false;
  public isMissing = false;
  public applicationMeta = {} as Application;
  public subscriptions = {} as ProjectSubscriptions;
  public responsible: string;
  public customIds = [] as CustomId[];
  public isDummy = false;
  public dummyDeletionDate: string;
  public hasChildren = false;
  public hasApproval = false;
  public hasSBOMToRetain = false;

  public constructor() {
    this._key = '';
    this.name = '';
    this.applicationId = '';
    this.parent = '';
    this.parentName = '';
    this.description = '';
    this.schemaLabel = '';
    this.status = '';
    this.policyLabels = [];
    this.projectLabels = [];
    this.freeLabels = [];
    this.children = [];
    this.updated = '';
    this.created = '';
    this.correspondingSchema = {} as Schema;
    this.documentMeta = {} as DisclosureDocumentMeta;
    this.customerMeta = {} as CustomerMetaDTO;
    this.noticeContactMeta = {} as NoticeContactMetaDTO;
    this.token = [];
    this.accessRights = {} as Rights;
    this.isNoFoss = false;
    this.supplierMissing = false;
    this.isMissing = false;
    this.responsible = '';
    this.customIds = [];
    this.isDummy = false;
    this.dummyDeletionDate = '';
    this.hasChildren = false;
    this.hasApproval = false;
    this.hasSBOMToRetain = false;
  }
}

export class ProjectMemberRequest {
  public targetUser: string;
  public userType: string;
  public comment: string;
  public responsible: boolean;

  constructor(targetUser: string, userType: string, comment: string, responsible: boolean) {
    this.targetUser = targetUser;
    this.userType = userType;
    this.comment = comment;
    this.responsible = responsible;
  }
}

export class MultiProjectMemberRequest extends ProjectMemberRequest {
  public targetProjects: string[];

  constructor(targetUser: string, userType: string, comment: string, responsible: boolean, targetProjects: string[]) {
    super(targetUser, userType, comment, responsible);
    this.targetProjects = targetProjects;
  }
}

export class FillCustomerReq {
  public customerApprover1 = '';
  public customerApprover2 = '';
  constructor(cust1: string, cust2: string) {
    this.customerApprover1 = cust1;
    this.customerApprover2 = cust2;
  }
}

export class ApprovableDto {
  public projectName = '';
  public projectKey = '';
  public approvablespdx: ApprovableSPDXDto = {} as ApprovableSPDXDto;
  public spdxname = '';
  public spdxtag = '';
  public customerdiff = false;
  public supplierdiff = false;
  public stats: ComponentStats = {} as ComponentStats;
  public spdxUploaded = '';
  public IsSpdxRecent = false;
}

export class TokenRequest {
  public Company: string;
  public Description: string;
  public Expiry: string;
  public Status: string;

  constructor(company: string, description: string, expiry: string, status: string) {
    this.Company = company;
    this.Description = description;
    this.Expiry = expiry;
    this.Status = status;
  }
}

export class CommonLicenses {
  public LicenseId = '';
  public Name = '';
  public SpdxDocument = '';
  public ExternalDocumentId = '';
  public SeeAlsos: string[] = [];
  public Comment = '';
  public ExtractedText = '';
}

export class IdentifiedLicense {
  public License = {} as CommonLicenses;
  public AliasTargetId = '';
}

export class DetailedLicense {
  public License: License = {} as License;
  public OrigName = '';
  public ReferencedName = '';
}

export class LicenseNameId {
  public name = '';
  public licenseId = '';
}

export class DetailedLicenseSlim {
  public License: LicenseNameId = {} as LicenseNameId;
  public OrigName = '';
  public ReferencedName = '';
}

export class Details {
  public Key = '';
  public Value = '';

  constructor(key: string, value: string) {
    this.Key = key;
    this.Value = value;
  }
}

export class ComponentDetails {
  public UnassertedLicenseText = false;
  public PolicyStatus = [] as PolicyRuleStatus[];
  public UnmatchedLicenses = [] as UnmatchedLicense[];
  public PolicyDecisionsApplied: PolicyDecisionSlim[] = [];
  public PolicyDecisionDeniedReason: string = '';

  public RawInfo = {} as IMap<string>;
  public Attributes = [] as Details[];

  public UnknownLicenses = [] as string[];
  public ExtractedLicenses = [] as CommonLicenses[];
  public IdentifiedViaAlias = [] as IdentifiedLicense[];
  public KnownLicenses = [] as DetailedLicense[];

  public Problems = [] as string[];

  public CanChooseLicense = false;
  public ChoiceDeniedReason = '';
  public LicenseRuleApplied?: LicenseRuleSlim;
  public ContainsOr = false;
}

export class ComponentLicenses {
  public UnknownLicenses = [] as string[];
  public KnownLicenses = [] as DetailedLicenseSlim[];
}

export class SbomLicense {
  public id = '';
  public origId = '';
  public name = '';
}

export class SbomLicenses {
  public unknown = [] as string[];
  public known = [] as SbomLicense[];
}

export class ProjectModel extends ProjectDTO {
  public supplierExtraData: SupplierExtraData = new SupplierExtraData();
}

export function getProjectUserTypes(): string[] {
  return Object.values(UserType);
}

export interface ProjectChildMemberCombi {
  projectKey: string;
  projectName: string;
  projectMember: ProjectUser;
}

export interface IProjectChildrenMembers {
  list: ProjectChildMemberCombi[];
}

export interface ProjectChildrenMemberSuccessResponse {
  success: boolean;
  userId: string;
  projectKey: string;
  projectName: string;
  message: string;
}

export class ProjectKeyName {
  key: string;
  name: string;

  constructor(key: string, name: string) {
    this.key = key;
    this.name = name;
  }
}

export function createProjectModel(project: ProjectModel) {
  const isProjectMember =
    project.accessRights &&
    (project.accessRights.groups.includes(Group.ProjectOwner) ||
      project.accessRights.groups.includes(Group.ProjectSupplier) ||
      project.accessRights.groups.includes(Group.ProjectViewer));
  const deptMissing =
    project.accessRights.allowProject.update && (project.documentMeta.deptMissing || project.customerMeta.deptMissing);

  const allowProjectCreate = project.accessRights && project.accessRights.allowProject.create;
  const allowProjectRead = project.accessRights && project.accessRights.allowProject.read;
  const allowProjectEdit = project.accessRights && project.accessRights.allowProject.update;
  const allowProjectDelete = project.accessRights && project.accessRights.allowProject.delete;
  const allowUserManagementCreate =
    project.accessRights.allowProjectUserManagement.create || project.accessRights.allowAllProjectUserManagement.create;
  const allowUserManagementUpdate =
    project.accessRights.allowProjectUserManagement.update || project.accessRights.allowAllProjectUserManagement.update;

  const allowUserManagementDelete =
    project.accessRights.allowProjectUserManagement.delete || project.accessRights.allowAllProjectUserManagement.delete;

  return {
    ...project,
    //Group
    allowGroupEdit: allowProjectEdit && project.isGroup,
    allowGroupCreate: allowProjectCreate && project.isGroup,
    //Project
    isProjectOwner: project.accessRights.groups.includes(Group.ProjectOwner),
    isProjectMember,
    isUserManagementAllowed:
      project.accessRights &&
      (project.accessRights.allowProjectUserManagement.read || project.accessRights.allowAllProjectUserManagement.read),
    isTokenManagementAllowed:
      project.accessRights &&
      (project.accessRights.allowProjectTokenManagement.read ||
        project.accessRights.allowAllProjectTokenManagement.read),
    allowProjectRead,
    allowProjectCreate,
    allowProjectEdit,
    allowProjectDelete,
    allowUserManagementCreate,
    allowUserManagementUpdate,
    allowUserManagementDelete,
    hasSubscriptions: (project.subscriptions?.spdx ?? false) || (project.subscriptions?.overallReview ?? false),
    showSubscriptionButton: !project.isGroup && isProjectMember && !(project.status === 'deprecated'),
    appLinkText: !project.applicationMeta.id
      ? ''
      : project.applicationMeta.secondaryId
        ? `${project.applicationMeta.name} (${project.applicationMeta.secondaryId})`
        : project.applicationMeta.name,
    deptMissing,
    isApprovalDisabled: Boolean(project.parent) || deptMissing || project.isDummy,
    isApprovalAllowed: project.accessRights.allowProject.update && project.accessRights.allowRequestApproval.create,
    isDeprecated: project.status === 'deprecated',
    projectChildren: {} as ProjectChildren,
    hasParent: Boolean(project.parent),
  };
}

export type Project = ReturnType<typeof createProjectModel>;
