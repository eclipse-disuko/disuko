// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useApi} from '@shared/api/useApi';
import {I18nImportResponse, I18nLocaleListItem, I18nLocaleResponse} from '@disclosure-portal/model/I18n';
import {AxiosResponse} from 'axios';

const {api} = useApi();

class I18nService {
  public getLocales(): Promise<AxiosResponse<I18nLocaleListItem[]>> {
    return api.get<I18nLocaleListItem[]>('/api/v1/i18n');
  }

  public getLocale(code: string): Promise<AxiosResponse<I18nLocaleResponse>> {
    return api.get<I18nLocaleResponse>(`/api/v1/i18n/${encodeURIComponent(code)}`);
  }

  public upsertTranslation(locale: string, key: string, value: string): Promise<AxiosResponse> {
    return api.put(`/api/v1/i18n/${encodeURIComponent(locale)}/${encodeURIComponent(key)}`, {
      value,
      description: '',
    });
  }

  public deleteTranslation(locale: string, key: string): Promise<AxiosResponse> {
    return api.delete(`/api/v1/i18n/${encodeURIComponent(locale)}/${encodeURIComponent(key)}`);
  }

  public exportLocale(locale: string): Promise<AxiosResponse<Blob>> {
    return api.get<Blob>(`/api/v1/i18n/export/${encodeURIComponent(locale)}`, {responseType: 'blob'});
  }

  public importLocale(locale: string, formData: FormData): Promise<AxiosResponse<I18nImportResponse>> {
    return api.post<I18nImportResponse>(`/api/v1/i18n/${encodeURIComponent(locale)}/import`, formData, {
      headers: {'Content-Type': 'multipart/form-data'},
    });
  }
}

const i18nService = new I18nService();
export default i18nService;
