// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {VersionSlim} from '@disclosure-portal/model/VersionDetails';

export default class ProjectVersionPostRequest {
  public static toProjectVersionPostRequest(version: VersionSlim): ProjectVersionPostRequest {
    const pvpr = new ProjectVersionPostRequest();
    pvpr.name = version.name;
    pvpr.description = version.description;
    return pvpr;
  }

  public name: string;
  public description: string;

  public constructor() {
    this.name = '';
    this.description = '';
  }
}
