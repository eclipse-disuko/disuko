// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useApi} from '@shared/api/useApi';
import {Department} from '@disclosure-portal/model/Department';

const {api} = useApi();

class CompanyService {
  public async find(searchStr: string): Promise<Department[]> {
    if (searchStr.length < 3) {
      return [];
    }
    searchStr = encodeURIComponent('' + searchStr);
    return (await api.get<Department[]>(`/api/v1/departments/find/${searchStr}`)).data;
  }
}

const companyService = new CompanyService();
export default companyService;
