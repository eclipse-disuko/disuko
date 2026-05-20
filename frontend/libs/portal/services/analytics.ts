// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useApi} from '@shared/api/useApi';
import {
  CombinedSearchOptions,
  IAnalyticsComponentsSearchResponse,
  IAnalyticsLicenseSearchResponse,
  IAnalyticsSearchRequest,
  IAnalyticsSearchResponse,
  OccurencesRes,
} from '@disclosure-portal/model/Analytics';

const {api} = useApi();

const modelName = 'analytics';

class AnalyticsService {
  public async searchAnalytics(data: CombinedSearchOptions, all: boolean) {
    return api.post<IAnalyticsSearchResponse>(`/api/v1/${modelName}/search?all=${all}`, data);
  }

  public async searchOccurrencies() {
    return api.get<OccurencesRes>(`/api/v1/${modelName}/occurrences`);
  }

  public async searchComponents(data: IAnalyticsSearchRequest) {
    return api.post<IAnalyticsComponentsSearchResponse>(`/api/v1/${modelName}/components/search`, data);
  }

  public downloadReport() {
    return api.get(`/api/v1/${modelName}/report`, {
      responseType: 'blob',
    });
  }

  public async searchLicenses(data: IAnalyticsSearchRequest) {
    return api.post<IAnalyticsLicenseSearchResponse>(`/api/v1/${modelName}/licenses/search`, data);
  }

  public async getStats() {
    return api.get(`/api/v1/${modelName}/stats`);
  }

  public async export() {
    return api.post(`/api/v1/${modelName}/export`);
  }
}

const analyticsService = new AnalyticsService();
export default analyticsService;
