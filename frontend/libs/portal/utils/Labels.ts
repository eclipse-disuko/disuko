// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import Label from '@disclosure-portal/model/Label';
import AdminService from '@disclosure-portal/services/admin';
import {IMap} from '@disclosure-portal/utils/View';

export class LabelsTools {
  public schemaLabelItems: Label[] = [] as Label[];
  public schemaLabelsMapByName: IMap<Label> = {} as IMap<Label>;
  public schemaLabelsMap: IMap<Label> = {} as IMap<Label>;

  public policyLabelsMap: IMap<Label> = {} as IMap<Label>;
  public policyLabelItems: Label[] = [] as Label[];
  public policyLabelsMapByName: IMap<Label> = {} as IMap<Label>;

  public projectLabelsMap: IMap<Label> = {} as IMap<Label>;
  public projectLabelItems: Label[] = [] as Label[];
  public projectLabelsMapByName: IMap<Label> = {} as IMap<Label>;

  public async loadLabels() {
    this.schemaLabelItems = (await AdminService.getSchemaLabels()).data;
    this.createSchemaLabelsMap();
    this.policyLabelItems = (await AdminService.getPolicyLabels()).data;
    this.createPolicyLabelsMap();
    this.projectLabelItems = (await AdminService.getProjectLabels()).data;
    this.createProjectLabelsMap();
  }

  public createSchemaLabelsMap() {
    if (!(this.schemaLabelItems && this.schemaLabelsMapByName && this.schemaLabelsMap)) return;
    for (const lbl of this.schemaLabelItems) {
      this.schemaLabelsMap[lbl._key] = lbl;
      this.schemaLabelsMapByName[lbl.name] = lbl;
    }
  }

  public createPolicyLabelsMap() {
    if (!(this.policyLabelItems && this.policyLabelsMapByName && this.policyLabelsMap)) return;
    for (const lbl of this.policyLabelItems) {
      this.policyLabelsMap[lbl._key] = lbl;
      this.policyLabelsMapByName[lbl.name] = lbl;
    }
  }

  public createProjectLabelsMap() {
    if (!(this.projectLabelItems && this.projectLabelsMapByName && this.projectLabelsMap)) return;
    for (const lbl of this.projectLabelItems) {
      this.projectLabelsMap[lbl._key] = lbl;
      this.projectLabelsMapByName[lbl.name] = lbl;
    }
  }

  public camelCaseToLabel(str: string): string {
    return str
      .replace(/([a-z0-9])([A-Z])/g, '$1 $2')
      .replace(/([A-Z])([A-Z][a-z])/g, '$1 $2')
      .replace(/^./, function (str) {
        return str.toUpperCase();
      });
  }
}
