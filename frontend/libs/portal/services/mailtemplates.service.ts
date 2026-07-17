// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {MailTemplate, UpdateMailTemplate} from '@disclosure-portal/model/MailTemplate';
import {useApi} from '@shared/api/useApi';
import {AxiosResponse} from 'axios';

const {api} = useApi();

class MailTemplatesService {
  public getAll(): Promise<AxiosResponse<MailTemplate[]>> {
    return api.get<MailTemplate[]>('/api/v1/admin/mailtemplates/');
  }

  public getById(id: string): Promise<AxiosResponse<MailTemplate>> {
    return api.get<MailTemplate>(`/api/v1/admin/mailtemplates/${encodeURIComponent(id)}`);
  }

  public update(id: string, data: UpdateMailTemplate): Promise<AxiosResponse<MailTemplate>> {
    return api.put<MailTemplate>(`/api/v1/admin/mailtemplates/${encodeURIComponent(id)}`, data);
  }

  public test(id: string, message: string): Promise<AxiosResponse<void>> {
    return api.post<void>(`/api/v1/admin/mailtemplates/${encodeURIComponent(id)}/test`, {message});
  }
}

const mailTemplatesService = new MailTemplatesService();
export default mailTemplatesService;
