// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useApi} from '@shared/api/useApi';
import {FilterSetDto, FilterSetRequestDto} from '@disclosure-portal/model/FilterSet';

const {api} = useApi();

class FilterSetService {
  public async getFilterSets(tablename: string): Promise<FilterSetDto[]> {
    const result = await api.get(`api/v1/filtersets/tables/${tablename}`);
    return result.data;
  }

  public async getFilterSet(key: string): Promise<FilterSetDto> {
    key = encodeURIComponent('' + key).replace(/\./g, '%2E');
    const result = await api.get(`api/v1/filtersets/${key}`);
    return result.data;
  }

  public async create(data: FilterSetRequestDto) {
    return api.post<FilterSetDto>(`api/v1/filtersets`, data);
  }

  public async update(data: FilterSetRequestDto, key: string) {
    return api.put<FilterSetRequestDto>(`api/v1/filtersets/${key}`, data);
  }

  public async delete(key: string) {
    key = encodeURIComponent('' + key).replace(/\./g, '%2E');
    return api.delete(`api/v1/filtersets/${key}`);
  }
}
const filterSetService = new FilterSetService();
export default filterSetService;
