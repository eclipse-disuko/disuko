// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useApi} from '@shared/api/useApi';
import SimpleProfileData from '@shared/types/ProfileData';
import {ProjectRoleDto, TaskDto, UserDto, UserRequestDto} from '@shared/types/Users';
import {AxiosResponse} from 'axios';
import {DashboardCounts} from '@shared/types/DashboardCounts';

const {api} = useApi();

const modelName = 'profile';

class ProfileService {
  public async getDashboardCounts(): Promise<DashboardCounts> {
    const result = await api.get('/api/v1/counts/dashboard');
    return result.data;
  }

  public async getProfileData(): Promise<SimpleProfileData> {
    const result: AxiosResponse<SimpleProfileData> = await api.get(`/api/v1/${modelName}`, {withCredentials: true});
    return result.data;
  }

  public async update(userId: string, data: UserRequestDto) {
    userId = encodeURIComponent('' + userId).replace(/\./g, '%2E');
    return api.put<UserRequestDto>(`/api/v1/${modelName}/${userId}`, data);
  }

  public async getTasks(): Promise<TaskDto[]> {
    const result = await api.get(`/api/v1/${modelName}/tasks`);
    return result.data;
  }

  public async getTask(taskId: string): Promise<TaskDto> {
    const result = await api.get(`/api/v1/${modelName}/tasks/${taskId}`);
    return result.data;
  }

  public async getProjectRoles(): Promise<ProjectRoleDto[]> {
    const result = await api.get(`/api/v1/${modelName}/projectroles`);
    return result.data;
  }
  public getUsersBySearchFragment(searchFragment: string, active?: boolean) {
    searchFragment = encodeURIComponent('' + searchFragment).replace(/\./g, '%2E');
    const params = active !== undefined ? `?active=${active}` : '';
    return api.get<UserDto[]>(`/api/v1/${modelName}/search/${searchFragment}${params}`);
  }

  public downloadTasksCsv() {
    return api.get(`/api/v1/${modelName}/tasks/csv`);
  }

  public async delegateTask(taskId: string, delegateUserId: string): Promise<AxiosResponse<TaskDto>> {
    return api.put(`/api/v1/${modelName}/tasks/${taskId}/delegate`, {delegateUserId});
  }
}

const profileService = new ProfileService();

export default profileService;
