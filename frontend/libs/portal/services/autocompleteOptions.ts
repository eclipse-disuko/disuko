// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import companyService from '@disclosure-portal/services/companies';
import projectService from '@disclosure-portal/services/projects';
import profileService from '@shared/user/services/profile.service';

export const fetchCompanyOptions = (query: string) => companyService.find(query);

export const fetchProfileUserOptions = async (query: string, active?: boolean) =>
  (await profileService.getUsersBySearchFragment(query, active)).data;

export const fetchProjectUserOptions = (projectKey: string) => async (query: string, active?: boolean) =>
  (await projectService.getUsersBySearchFragment(projectKey, query, active)).data;
