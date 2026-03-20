import {Department} from '@disclosure-portal/model/Department';
import {
  CustomerMetaDTO,
  DisclosureDocumentMeta,
  NoticeContactMetaDTO,
  ProjectModel,
} from '@disclosure-portal/model/Project';
import {ProjectSlim} from '@disclosure-portal/model/ProjectsResponse';
import {SupplierExtraData, WizardProjectPostRequest} from '@disclosure-portal/model/Wizard';
import {Application} from './Application';

export class ProjectPostResponse {
  public name = '';
  public id = '';
  public taskGuid = '';
}

export class ProjectSettingsRequest {
  public documentMeta = new DisclosureDocumentMeta(null);
  public customerMeta = new CustomerMetaDTO();
  public noticeContactMeta = new NoticeContactMetaDTO();
  public supplierExtraData = new SupplierExtraData();
  public noFossProject = false;
}

export default class ProjectPostRequest {
  public name: string;
  public freeLabels: string[];
  public schemaLabel: string;
  public policyLabels: string[];
  public projectLabels: string[];
  public children: string[];
  public parent: string;
  public description: string;
  public Updated: string;
  public owner: string;
  public id: string;
  public isGroup = false;
  public isNoFoss = false;
  public projectSettings: ProjectSettingsRequest | null = null;
  public applicationMeta: Application | undefined = undefined;

  public constructor() {
    this.name = '';
    this.children = [];
    this.applicationMeta = new Application();
    this.parent = '';
    this.schemaLabel = '';
    this.policyLabels = [];
    this.projectLabels = [];
    this.freeLabels = [];
    this.description = '';
    this.Updated = '';
    this.owner = '';
    this.id = '';
    this.isNoFoss = false;
  }

  public fillWithProjectSlim(model: ProjectSlim) {
    if (!model) return;
    this.name = model.name;
    this.description = model.description;
    this.schemaLabel = model.schemaLabel;
    this.policyLabels = model.policyLabels;
    this.freeLabels = model.freeLabels;
    this.projectLabels = model.projectLabels;
    this.children = model.children;
    this.isGroup = model.isGroup;
    this.id = model._key;
    this.isNoFoss = model.isNoFoss;
    if (model.applicationMeta) {
      this.applicationMeta = {
        id: model.applicationMeta.id,
        name: model.applicationMeta.name,
        secondaryId: model.applicationMeta.secondaryId,
        externalLink: model.applicationMeta.externalLink,
      };
    }
  }

  public fillWithProjectModel(model: ProjectModel) {
    this.name = model.name;
    this.description = model.description;
    this.schemaLabel = model.schemaLabel;
    this.policyLabels = model.policyLabels;
    this.projectLabels = model.projectLabels;
    this.freeLabels = model.freeLabels;
    this.children = model.children;
    this.id = model._key;
    this.isGroup = model.isGroup;
    if (!this.children) {
      this.children = [];
    }
    this.projectSettings = new ProjectSettingsRequest();
    this.projectSettings.documentMeta.fill(model.documentMeta);
    this.projectSettings.customerMeta.fill(model.customerMeta);
    this.projectSettings.noticeContactMeta.fill(model.noticeContactMeta);
    if (model.supplierExtraData) {
      this.projectSettings.supplierExtraData = model.supplierExtraData;
    }
    this.isNoFoss = model.isNoFoss;
    this.projectSettings.noFossProject = model.isNoFoss;
    if (model.applicationMeta) {
      this.applicationMeta = {
        id: model.applicationMeta.id,
        name: model.applicationMeta.name,
        secondaryId: model.applicationMeta.secondaryId,
        externalLink: model.applicationMeta.externalLink,
      };
    }
  }

  public fillFromWizard(wizard: WizardProjectPostRequest) {
    this.name = wizard.name;
    this.freeLabels = wizard.freeLabels;
    this.schemaLabel = wizard.schemaLabel;
    this.policyLabels = wizard.policyLabels;
    this.projectLabels = wizard.projectLabels;
    this.children = wizard.children;
    this.description = wizard.description;
    this.owner = wizard.owner;
    this.id = wizard.id;
    this.isGroup = wizard.isGroup;
    this.projectSettings = new ProjectSettingsRequest();
    this.projectSettings.documentMeta.fill(wizard.projectSettings.documentMeta);

    const supplierDept = this.projectSettings.documentMeta.supplierDept;
    if (!supplierDept || !supplierDept.deptId || supplierDept.deptId.trim() === '') {
      this.projectSettings.documentMeta.supplierDept = {} as Department;
    }

    this.projectSettings.customerMeta.fill(wizard.projectSettings.customerMeta);
    this.projectSettings.noticeContactMeta.fill(wizard.projectSettings.noticeContactMeta);
    if (wizard.projectSettings.supplierExtraData) {
      this.projectSettings.supplierExtraData = wizard.projectSettings.supplierExtraData;
    }
    this.isNoFoss = wizard.projectSettings.noFossProject;
    this.projectSettings.noFossProject = wizard.projectSettings.noFossProject;
    this.applicationMeta = {
      id: wizard.applicationMeta.id,
      name: wizard.applicationMeta.name,
      secondaryId: wizard.applicationMeta.secondaryId,
      externalLink: wizard.applicationMeta.externalLink,
    };
  }
}
