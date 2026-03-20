import {ProjectSettingsModel} from '@disclosure-portal/model/Project';
import ProjectPostRequest from '@disclosure-portal/model/ProjectPostRequest';
import {UserDto} from '@disclosure-portal/model/Users';
import {Application} from './Application';

export enum ProjectCreationMode {
  default = 'default',
}

export class SupplierExtraData {
  public fRI = '';
  public sRI = '';
  public userFRI: UserDto | null = null;
  public userSRI: UserDto | null = null;
  public external = false;
}

export interface IProjectCreatedResponse {
  groupId: string;
  id: string;
  name: string;
  taskGuid: string;
}

export class WizardProjectPostRequest {
  public name: string;
  public description: string;
  public targetPlatform: string;
  public architecture: string;
  public distributionTarget: string;
  public targetUsers: string;
  public owner: string;
  public id: string;
  public freeLabels: string[];
  public schemaLabel: string;
  public policyLabels: string[];
  public projectLabels: string[];
  public children!: string[];
  public isGroup = false;
  public creationMode = ProjectCreationMode.default;
  public projectSettings = new ProjectSettingsModel();
  public applicationMeta = new Application();
  public parent: string;
  public parentName: string;
  public constructor() {
    this.owner = '';
    this.id = '';
    this.schemaLabel = '';
    this.policyLabels = [];
    this.projectLabels = [];
    this.freeLabels = [];
    this.name = '';
    this.description = '';
    this.targetPlatform = '';
    this.architecture = '';
    this.distributionTarget = '';
    this.targetUsers = '';
    this.creationMode = ProjectCreationMode.default;
    this.parent = '';
    this.parentName = '';
  }

  public fillWithData(project: ProjectPostRequest) {
    this.owner = project.owner;
    this.id = project.id;
    this.name = project.name;
    this.description = project.description;
    this.freeLabels = project.freeLabels;
    this.policyLabels = project.policyLabels;
    this.projectLabels = project.projectLabels;
    this.projectSettings.fill2(project);
  }
}

export interface WizardCard {
  key: string;
  image: string;
  class: string;
  title: string;
  subtitle: string;
  helptext: string;
  isFlipped: boolean;
  isActive?: boolean;
}
