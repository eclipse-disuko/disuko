// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export class UserRequestDto {
  public _key = '';
  public user = '';
  public forename = '';
  public lastname = '';
  public email = '';
  public termsOfUse = false;
  public termsOfUseDate = '';
  public active = false;
}

export class UserMetaData {
  public companyIdentifier = '';
  public department = '';
  public departmentDescription = '';
}

export class UserDto {
  public user = '';
  public forename = '';
  public lastname = '';
  public email = '';
  public termsOfUse = false;
  public termsOfUseDate = '';
  public termsOfUseVersion = '';
  public created = '';
  public updated = '';
  public isSelectable = false;
  public _key = '';
  public roles: string[] = [];
  public metaData = new UserMetaData();
  public active = false;
  public isInternal = false;
  public deprovisioned = '';
}

export class UserLastSeenDto {
  public newsboxLastSeenId = '';
}

export class UserRolesRequestDto {
  public roles: string[] = [];

  constructor(roles: string[]) {
    this.roles = roles;
  }
}

export class UserList {
  public items: UserDto[] = [];
  public count = 0;
}

export class TaskDto {
  public id = '';
  public created = '';
  public updated = '';
  public approvalGuid = '';
  public approvalType = '';
  public creator = '';
  public creatorDepartment = '';
  public creatorDepartmentDescription = '';
  public creatorFullName = '';
  public delegatedTo = '';
  public delegatedToFullName = '';
  public projectGuid = '';
  public status = '';
  public resultStatus = '';
  public type = '';
  public projectName = '';
  public isProjectGroup = false;
  public projectType = '';
}
export class TermsOfUseVersionResponse {
  public termsOfUseCurrentVersion = '';
}

export class ProjectRoleDto {
  public projectName = '';
  public projectKey = '';
  public userId = '';
  public userType = '';
  public responsible = false;
}

export class UserMailDto {
  public user = '';
  public forename = '';
  public lastname = '';
  public email = '';
}
