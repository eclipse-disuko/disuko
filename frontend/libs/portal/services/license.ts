// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useApi} from '@shared/api/useApi';
import {ICountResponse, IFoundResponse} from '@disclosure-portal/model/Common';
import License, {
  LicenseSlim,
  LicensesResponse,
  LicenseWithSimilarity,
  lookupRequest,
  lookupResponse,
} from '@disclosure-portal/model/License';
import {PolicyRulesForLicenseDto} from '@disclosure-portal/model/PolicyRule';
import {AuditLog} from '@disclosure-portal/model/VersionDetails';
import {SearchOptions} from '@disclosure-portal/utils/Table';
import {AxiosResponse} from 'axios';

const {api} = useApi();

const modelName = 'licenses';

class LicenseService {
  public async create(data: License) {
    return api.post(`/api/v1/${modelName}`, data);
  }

  public get = (id: string) => api.get<License>(`/api/v1/${modelName}/${encodeURIComponent(id)}`);

  public async getAuditTrail(id: string): Promise<AuditLog[]> {
    return (await api.get<AuditLog[]>(`/api/v1/${modelName}/${encodeURIComponent(id)}/audit`)).data;
  }

  public async head(id: string): Promise<AxiosResponse<IFoundResponse>> {
    return api.get(`/api/v1/${modelName}/exists/${encodeURIComponent(id)}`);
  }

  public async delete(id: string) {
    return api.delete(`/api/v1/${modelName}/${encodeURIComponent(id)}`);
  }

  public async getAll() {
    return api.get<LicenseSlim[]>(`/api/v1/${modelName}`);
  }

  public async search(options: SearchOptions, signal?: AbortSignal) {
    return api.post<LicensesResponse>(`/api/v1/${modelName}/search`, options, signal ? {signal} : undefined);
  }

  public async headAlias(alias: string): Promise<AxiosResponse<IFoundResponse>> {
    return api.get(`/api/v1/${modelName}/aliases/${encodeURIComponent(alias)}`);
  }

  public async getList(idList: string) {
    return api.get<License[]>(`/api/v1/${modelName}/list/${encodeURIComponent(idList)}`);
  }

  public getCountOfLicencesUsingThisObligation(obligationKey: string): Promise<AxiosResponse<ICountResponse>> {
    return api.get<ICountResponse>(`/api/v1/${modelName}/obligation/${obligationKey}/usagecount`);
  }

  public getCountOfPolicyRuleUsingThisLicence(licenceKey: string): Promise<AxiosResponse<ICountResponse>> {
    return api.get<ICountResponse>(`/api/v1/${modelName}/policyrules/${licenceKey}/usagecount`);
  }

  public async getPolicyRuleAssignmentsForThisLicence(licenceKey: string): Promise<PolicyRulesForLicenseDto> {
    return (await api.get<PolicyRulesForLicenseDto>(`/api/v1/${modelName}/policyrules/${licenceKey}`)).data;
  }

  public async updatePolicyRulesAssignmentsForLicense(
    policiesForLicense: PolicyRulesForLicenseDto,
  ): Promise<AxiosResponse<string>> {
    return api.put<string>(`/api/v1/${modelName}/policyrules/${policiesForLicense.id}`, policiesForLicense);
  }

  public async update(data: License, uid: string) {
    return api.put<string>(`api/v1/${modelName}/${encodeURIComponent(uid)}`, data);
  }

  public async searchForSimilarLicenseText(licenseText: string) {
    return api.post<LicenseWithSimilarity[]>(`/api/v1/${modelName}/text/compare`, licenseText);
  }

  public async deleteAlias(id: string, alias: string) {
    return api.delete(`/api/v1/${modelName}/${id}/aliases/${alias}`);
  }

  public async getAllWithOptions(options: SearchOptions) {
    return api.post<LicensesResponse>(`/api/v1/${modelName}/search`, options);
  }

  public async lookup(ids: string[]) {
    const req: lookupRequest = {ids: ids};
    return api.post<lookupResponse>(`/api/v1/${modelName}/lookup`, req);
  }
}

const licenseService = new LicenseService();

export default licenseService;
