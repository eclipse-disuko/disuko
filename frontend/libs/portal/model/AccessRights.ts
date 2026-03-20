// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {ActionRights} from '@disclosure-portal/model/Credentials';
import {CRUDRights} from '@disclosure-portal/model/Rights';

export class AccessRights {
  public allowSchema!: CRUDRights;
  public allowLabel!: CRUDRights;
  public allowPolicy!: CRUDRights;
  public allowObligation!: CRUDRights;
  public allowLicense!: CRUDRights;
  public allowTools!: CRUDRights;
  public allowStyleguide!: CRUDRights;
  public allowS3Tests!: CRUDRights;
  public allowSampleData!: CRUDRights;
  public allowUsers!: CRUDRights;
  public allowAllProjectUserManagement!: CRUDRights;
  public allowAllProjectTokenManagement!: CRUDRights;
  public allowProjectAudit!: CRUDRights;
  public allowProject!: CRUDRights;
  public AllowProjectGroup!: CRUDRights;
  public allowRequestApproval!: CRUDRights;
  public allowTask!: CRUDRights;
  public allowAnnouncement!: CRUDRights;
  public allowRequestPlausi!: CRUDRights;
}

export class ProjectAccessRights {
  public allowProject!: CRUDRights;
  public AllowProjectGroup!: CRUDRights;
  public allowProjectVersion!: CRUDRights;
  public allowProjectPolicy!: CRUDRights;
  public allowProjectUserManagement!: CRUDRights;
  public allowProjectTokenManagement!: CRUDRights;
  public allowSBOMAction!: ActionRights;
  public allowCCSAction!: ActionRights;
  public allowDisclosureDocument!: ActionRights;
}
