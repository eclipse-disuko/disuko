// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useApi} from '@shared/api/useApi';
import {Application} from '@disclosure-portal/model/Application';

const {api} = useApi();

class ApplicationService {
  public async searchApplicationByQuery(query: string) {
    const response = await api.get<Application[]>(`/api/v1/api/applications/search?query=${query}`);
    return response.data;
  }
}

const applicationService = new ApplicationService();

export default applicationService;
