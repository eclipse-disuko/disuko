// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {getApi} from '@disclosure-portal/api';
import {NewsboxItem, NewsboxItemCreateDto, default as Newsbox} from '@disclosure-portal/model/Newsbox';
import {UserLastSeenDto} from '@disclosure-portal/model/Users';
import {AxiosResponse} from 'axios';

const {api} = getApi();

const basePath = 'newsbox/items';

class NewsboxService {
  public async getNewsboxItems(): Promise<AxiosResponse<Newsbox>> {
    return api.get<Newsbox>(`/api/v1/newsbox/items`);
  }

  public async getAllNewsboxItems(): Promise<AxiosResponse<NewsboxItem[]>> {
    return api.get<NewsboxItem[]>(`/api/v1/admin/newsbox/items`);
  }

  public async updateLastSeen(userId: string, data: UserLastSeenDto) {
    userId = encodeURIComponent('' + userId).replace(/\./g, '%2E');
    return api.put(`/api/v1/newsbox/items/${userId}`, data);
  }

  public createNewsboxItem = (item: NewsboxItemCreateDto): Promise<AxiosResponse<string>> => {
    return api.post<string>(`/api/v1/admin/${basePath}`, item);
  };

  public updateNewsboxItem = (id: string, item: NewsboxItem): Promise<AxiosResponse<string>> => {
    const processedItem = {
      ...item,
      expiry: item.expiry && item.expiry !== '' ? item.expiry : null,
    };
    return api.put<string>(`/api/v1/admin/${basePath}/${encodeURIComponent(id)}`, processedItem);
  };

  public deleteItemsAdmin = (id: string): Promise<AxiosResponse<string>> => {
    return api.delete<string>(`/api/v1/admin/${basePath}/${encodeURIComponent(id)}`);
  };
}

const newsboxService = new NewsboxService();
export default newsboxService;
