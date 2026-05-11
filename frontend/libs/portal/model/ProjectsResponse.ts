// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {CustomId} from '@disclosure-portal/model/CustomId';
import {ProjectModel} from '@disclosure-portal/model/Project';
import {Rights} from '@disclosure-portal/model/Rights';
import {SpdxFile, VersionSlimDto} from '@disclosure-portal/model/VersionDetails';
import {Application} from './Application';

export class ProjectSlimDto {
  public _key: string;
  public name: string;
  public description: string;
  public applicationId: string;
  public schemaLabel: string;
  public policyLabels: string[];
  public projectLabels: string[];
  public children: string[];
  public freeLabels: string[];
  public parent: string;
  public created?: string;
  public updated: string;
  public accessRights: Rights;
  public status: string;
  public isGroup = false;
  public isDeleted = false;
  public supplier: string;
  public supplierMissing = false;
  public company: string;
  public department: string;
  public missing = false;
  public isNoFoss = false;
  public applicationMeta = new Application();
  public isInGroupApproval = false;
  public responsible: string;
  public customIds: CustomId[];
  public isDummy: boolean;
  public dummyDeletionDate: string;
  public deleteDisabledReason?: string;
  public hasChildren: boolean;
  public hasApproval: boolean;
  public hasSBOMToRetain: boolean;

  constructor() {
    this._key = '';
    this.name = '';
    this.description = '';
    this.applicationId = '';
    this.schemaLabel = '';
    this.policyLabels = [];
    this.projectLabels = [];
    this.parent = '';
    this.children = [];
    this.freeLabels = [];
    this.updated = '';
    this.accessRights = {} as Rights;
    this.status = '';
    this.isGroup = false;
    this.supplier = '';
    this.company = '';
    this.department = '';
    this.applicationMeta = new Application();
    this.responsible = '';
    this.customIds = [];
    this.isDummy = false;
    this.dummyDeletionDate = '';
    this.hasChildren = false;
    this.hasApproval = false;
    this.hasSBOMToRetain = false;
  }
}

export interface ProjectChildrenCombiDto {
  projectKey: string;
  project: ProjectSlimDto;
  version: VersionSlimDto;
  hasProjectReadAccess: boolean;
}

export class ProjectSlim extends ProjectSlimDto {
  public newAdded = false;

  constructor(dto: ProjectSlimDto) {
    super();
    Object.assign(this, dto);
  }

  fromProjectModel(model: ProjectModel) {
    this._key = model._key;
    this.name = model.name;
    this.description = model.description;
    this.schemaLabel = model.schemaLabel;
    this.policyLabels = model.policyLabels;
    this.projectLabels = model.projectLabels;
    this.freeLabels = model.freeLabels;
    this.children = model.children;
    this.isGroup = model.isGroup;
    this.isNoFoss = model.isNoFoss;
    this.isDummy = model.isDummy;
    this.dummyDeletionDate = model.dummyDeletionDate;
    this.hasApproval = model.hasApproval;
    this.hasChildren = model.hasChildren;
    this.hasSBOMToRetain = model.hasSBOMToRetain;
    if (model.applicationMeta) {
      this.applicationMeta = {
        id: model.applicationMeta.id,
        name: model.applicationMeta.name,
        secondaryId: model.applicationMeta.secondaryId,
        externalLink: model.applicationMeta.externalLink,
      };
    }
  }
}

interface IProjectsResponse {
  projects: ProjectSlim[];
  count: number;
}

export class ProjectsResponseDTO implements IProjectsResponse {
  public projects: ProjectSlim[];
  public count: number;

  public constructor() {
    this.projects = [];
    this.count = 0;
  }
}

export class ProjectsResponse extends ProjectsResponseDTO {
  constructor(dto: ProjectsResponseDTO) {
    super();
    Object.assign(this, dto);
  }
}

export interface ProjectChildren {
  list: ProjectChildrenCombiDto[];
  projects: N[];
}

export interface VersionSboms {
  versionName: string;
  versionKey: string;
  spdxFileHistory: SpdxFile[];
}

export class VersionSbomsFlat extends SpdxFile {
  public versionName = '';
  public versionKey = '';
}

export class NameKeyIdentifier {
  public name = '';
  public key = '';
}

export class ProjectSbomsFlat {
  public items = [] as VersionSbomsFlat[];
  public versions = [] as NameKeyIdentifier[];
}
