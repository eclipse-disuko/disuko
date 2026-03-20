// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {getApi} from '@disclosure-portal/api';
import {ChangeLogResponse} from '@disclosure-portal/model/ChangeLog';
import {Checklist, ChecklistItem} from '@disclosure-portal/model/Checklist';
import {ICountResponse, ISuccessRsponse} from '@disclosure-portal/model/Common';
import {CustomId, CustomIdUsage} from '@disclosure-portal/model/CustomId';
import {InternalToken} from '@disclosure-portal/model/InternalToken';
import IObligation from '@disclosure-portal/model/IObligation';
import {JobDto, SetConfigDto} from '@disclosure-portal/model/Job';
import Label from '@disclosure-portal/model/Label';
import License, {LicenseDiff} from '@disclosure-portal/model/License';
import {MailData} from '@disclosure-portal/model/MailData';
import {MemStats} from '@disclosure-portal/model/Memstats';
import {Notification, NotificationDto} from '@disclosure-portal/model/Notification';
import PolicyRule from '@disclosure-portal/model/PolicyRule';
import SimpleProfileData from '@disclosure-portal/model/ProfileData';
import {ReviewTemplate} from '@disclosure-portal/model/ReviewTemplate';
import SampleDataCreationState from '@disclosure-portal/model/SampleData';
import SchemaModel from '@disclosure-portal/model/Schema';
import SystemStatsResponse from '@disclosure-portal/model/Statistic';
import {DeletePersonalDataResponse} from '@disclosure-portal/model/UserDeletion';
import {
  ProjectRoleDto,
  TaskDto,
  TermsOfUseVersionResponse,
  UserDto,
  UserList,
  UserMailDto,
  UserRequestDto,
  UserRolesRequestDto,
} from '@disclosure-portal/model/Users';
import {AuditLog} from '@disclosure-portal/model/VersionDetails';
import {SearchOptions} from '@disclosure-portal/utils/Table';
import {AxiosResponse} from 'axios';

const modelName = 'admin';

export interface GetDashboardCounts {
  projectCount: number;
  licenseCount: number;
  policyRuleCount: number;
  labelCount: number;
  schemaCount: number;
  obligationCount: number;
  userCount: number;
  disclosureCount: number;
  reviewTemplateCount: number;
  activeJobCount: number;
  hasNewNewsboxItem: boolean;
}

const {api} = getApi();

class AdminService {
  public async getDashboardCounts(): Promise<GetDashboardCounts> {
    const result = await api.get('/api/v1/admin/counts/dashboard');
    return result.data;
  }

  public getAllSchemas = () => api.get<SchemaModel[]>(`/api/v1/${modelName}/schemas`);

  public getSchema = (key: string) => api.get<SchemaModel>(`/api/v1/${modelName}/schemas/${key}`);

  public async createSchema(data: FormData) {
    return api.post(`/api/v1/${modelName}/schemas`, data);
  }

  public async sendEmail(data: MailData) {
    return api.post(`/api/v1/${modelName}/mail/send`, data);
  }

  public activate(id: string) {
    return api.post(`/api/v1/${modelName}/schemas/${id}/activation`);
  }

  public downloadSchema(itemId: string) {
    return api.get(`/api/v1/${modelName}/schemas/${itemId}/download`, {
      responseType: 'blob',
    });
  }

  public async getAuditTrail(id: string): Promise<AuditLog[]> {
    return (await api.get<AuditLog[]>(`/api/v1/${modelName}/policyrules/${id}/audit`)).data;
  }

  public async getChangeLog(id: string): Promise<ChangeLogResponse[]> {
    return (await api.get<ChangeLogResponse[]>(`/api/v1/${modelName}/policyrules/${id}/changelog`)).data;
  }

  public deletePolicyRule(id: string) {
    return api.delete(`/api/v1/${modelName}/policyrules/${id}`);
  }

  public deprecatePolicyRule(id: string) {
    return api.put(`/api/v1/${modelName}/policyrules/${id}/deprecate`);
  }

  public copyPolicyRule(id: string) {
    return api.put(`/api/v1/${modelName}/policyrules/${id}/copy`);
  }

  public getAllObligations() {
    return api.get(`/api/v1/${modelName}/obligations`);
  }

  public postObligation(item: IObligation) {
    return api.post(`/api/v1/${modelName}/obligations`, item);
  }

  public putObligation(item: IObligation) {
    return api.put(`/api/v1/${modelName}/obligations/${item._key}`, item);
  }

  public deleteObligation(id: string) {
    return api.delete(`/api/v1/${modelName}/obligations/${id}`);
  }

  public postPolicyRule(rule: PolicyRule): Promise<AxiosResponse<PolicyRule>> {
    return api.post(`/api/v1/${modelName}/policyrules`, rule);
  }

  public putPolicyRule(rule: PolicyRule): Promise<AxiosResponse<PolicyRule>> {
    return api.put(`/api/v1/${modelName}/policyrules/${rule._key}`, rule);
  }

  /**
   * @deprecated referenced in probably unused view.
   */
  public getPolicyRuleCount(): Promise<AxiosResponse<ICountResponse>> {
    return api.get(`/api/v1/${modelName}/policyrules/count`);
  }

  public deleteLabel(key: string) {
    return api.delete(`/api/v1/${modelName}/labels/${key}`);
  }

  public editLabel(lbl: Label) {
    return api.put(`/api/v1/${modelName}/labels/${lbl._key}`, lbl);
  }

  public createLabel(lbl: Label) {
    return api.post(`/api/v1/${modelName}/labels`, lbl);
  }

  public getLabels = () => api.get<Label[]>(`/api/v1/${modelName}/labels`);

  public getPolicyLabels = () => api.get<Label[]>(`/api/v1/${modelName}/labels?type=POLICY`);

  public getProjectLabels = () => api.get<Label[]>(`/api/v1/${modelName}/labels?type=PROJECT`);

  public async getSchemaLabels() {
    return api.get<Label[]>(`/api/v1/${modelName}/labels?type=SCHEMA`);
  }

  public async triggerReloadLicensesJob() {
    return api.post(`/api/v1/${modelName}/licenses/import`);
  }

  public async exportLicenseKnowledgeBase(): Promise<AxiosResponse<string>> {
    return api.get(`/api/v1/${modelName}/licenses/knowledgebase/export`);
  }

  public async getSpdxLicensesCount(): Promise<ICountResponse> {
    return (await api.get<ICountResponse>(`/api/v1/${modelName}/licenses/spdx/count`)).data;
  }

  public async getLicensesDiffs(): Promise<LicenseDiff[]> {
    return (await api.get<LicenseDiff[]>(`/api/v1/${modelName}/licenses/spdx/diffs`)).data;
  }

  public async updateLicense(data: License, uid: string) {
    return api.put<string>(`api/v1/${modelName}/licenses/${encodeURIComponent(uid)}`, data);
  }

  public async deleteSpdxLicense(uid: string) {
    return api.delete(`/api/v1/${modelName}/licenses/spdx/${encodeURIComponent(uid)}`);
  }

  public async exportSchemaKnowledgeBase(): Promise<AxiosResponse<string>> {
    return api.get(`/api/v1/${modelName}/schemas/knowledgebase/export`);
  }

  public async triggerCreateSampleData(cnt: number, withFileUpload: boolean): Promise<AxiosResponse<string>> {
    return api.post(`/api/v1/${modelName}/utils/sampledata?cnt=` + cnt + '&fileUpload=' + withFileUpload);
  }

  public async stopCreateSampleData(): Promise<AxiosResponse<string>> {
    return api.delete(`/api/v1/${modelName}/utils/sampledata`);
  }

  public async getCreateSampleDataState(): Promise<AxiosResponse<SampleDataCreationState>> {
    return api.get(`/api/v1/${modelName}/utils/sampledata`);
  }

  public async getStats(): Promise<AxiosResponse<SystemStatsResponse>> {
    return api.get(`/api/v1/${modelName}/utils/stats`);
  }

  public async updateStats(): Promise<AxiosResponse<SystemStatsResponse>> {
    return api.put(`/api/v1/${modelName}/utils/stats`);
  }

  public async getJobsAll() {
    return (await api.get<JobDto[]>('/api/v1/admin/jobs')).data;
  }

  public async startJob(jobId: number) {
    return api.put(`/api/v1/admin/jobs/${jobId}`);
  }

  public async rerunOnetimeJob(key: string) {
    return api.put(`/api/v1/admin/jobs/onetime/${key}`);
  }

  public async setJobConfig(jobId: number, data: SetConfigDto): Promise<AxiosResponse<ISuccessRsponse>> {
    return api.put(`/api/v1/admin/jobs/${jobId}/config`, data);
  }

  public async getJobLatest(jobType: number) {
    return api.get(`/api/v1/admin/jobs/latest/${jobType}`);
  }

  // DB / S3 check
  public async getDbS3CheckStart() {
    return api.get('/api/v1/analyse/files/start');
  }

  public async getDbS3CheckStop() {
    return api.get('/api/v1/analyse/files/stop');
  }

  public async getDbS3CheckGetResult(): Promise<AxiosResponse<string>> {
    return api.get('/api/v1/analyse/files/status');
  }

  public async getAllProjectAccessRights() {
    return api.get(`/api/v1/${modelName}/utils/rightsProject`);
  }

  public async getAllAccessRights() {
    return api.get(`/api/v1/${modelName}/utils/rights`);
  }

  public getUsers = () => api.get<UserList>(`/api/v1/${modelName}/users/`);

  public async getUsersWithOptions(options: SearchOptions, signal?: AbortSignal) {
    return api.post<UserList>(`/api/v1/${modelName}/users/search`, options, signal ? {signal} : undefined);
  }

  public getUser = (uuid: string) => api.get<UserDto>(`/api/v1/${modelName}/users/${uuid}`);

  public async updateUserRoles(userId: string, data: UserRolesRequestDto): Promise<AxiosResponse<UserDto>> {
    userId = encodeURIComponent('' + userId).replace(/\./g, '%2E');
    return api.put<UserDto>(`/api/v1/${modelName}/users/${userId}/roles`, data);
  }

  public async getUserTokens(uuid: string): Promise<AxiosResponse<SimpleProfileData>> {
    return await api.get(`/api/v1/${modelName}/users/${uuid}/tokens`);
  }

  public async getUserTokensForNonInternal(uuid: string): Promise<AxiosResponse<SimpleProfileData>> {
    return await api.get(`/api/v1/${modelName}/users/${uuid}/tokensNonInternal`);
  }

  public async enableDisableUser(userId: string, data: UserRequestDto): Promise<AxiosResponse<UserDto>> {
    userId = encodeURIComponent('' + userId).replace(/\./g, '%2E');
    return api.put<UserDto>(`/api/v1/${modelName}/users/${userId}/active`, data);
  }

  public async getTermsOfUseCurrentVersion(): Promise<TermsOfUseVersionResponse> {
    return (await api.get<TermsOfUseVersionResponse>(`/api/v1/${modelName}/users/termsOfUseCurrentVersion`)).data;
  }

  public async getUserAuditTrail(uuid: string): Promise<AuditLog[]> {
    return (await api.get<AuditLog[]>(`/api/v1/${modelName}/users/${uuid}/audit`)).data;
  }

  public async getSystemProfile(): Promise<MemStats> {
    return (await api.get<MemStats>(`/api/v1/${modelName}/system/profile`)).data;
  }

  public async getClassificationAuditTrail(uuid: string): Promise<AuditLog[]> {
    return (await api.get<AuditLog[]>(`/api/v1/${modelName}/obligations/${uuid}/audit`)).data;
  }

  public async downloadLCcsv() {
    return api.get(`/api/v1/${modelName}/obligations/csv`);
  }

  public async downloadLPcsv() {
    return api.get(`/api/v1/${modelName}/policyrules/csv`);
  }

  public async downloadPLcsv() {
    return api.get(`/api/v1/${modelName}/labels/csv`);
  }

  public async downloadReviewTemplateCSV() {
    return api.get(`/api/v1/${modelName}/templates/review/csv`);
  }

  public getNotification = (): Promise<AxiosResponse<Notification>> =>
    api.get<Notification>(`/api/v1/${modelName}/notification`);

  public async setNotification(notification: NotificationDto) {
    return api.post(`api/v1/${modelName}/notification`, notification);
  }

  public async getUserProjectRoles(uuid: string): Promise<ProjectRoleDto[]> {
    const result = await api.get(`/api/v1/${modelName}/users/${uuid}/projectroles`);
    return result.data;
  }

  public async getUserTasks(uuid: string): Promise<TaskDto[]> {
    const result = await api.get(`/api/v1/${modelName}/users/${uuid}/tasks`);
    return result.data;
  }

  public async getUserMailById(userId: string): Promise<UserMailDto> {
    const result = await api.get(`/api/v1/${modelName}/users/mails/${userId}`);
    return result.data;
  }

  public createReviewTemplate(template: ReviewTemplate): Promise<AxiosResponse<ReviewTemplate>> {
    return api.post<ReviewTemplate>(`/api/v1/${modelName}/templates/review`, template);
  }

  public editReviewTemplate(template: ReviewTemplate): Promise<AxiosResponse<ReviewTemplate>> {
    return api.put<ReviewTemplate>(
      `/api/v1/${modelName}/templates/review/${encodeURIComponent(template._key)}`,
      template,
    );
  }

  public deleteReviewTemplate(id: string) {
    return api.delete<ISuccessRsponse>(`/api/v1/${modelName}/templates/review/${encodeURIComponent(id)}`);
  }

  public getReviewTemplates(): Promise<AxiosResponse<ReviewTemplate[]>> {
    return api.get<ReviewTemplate[]>(`/api/v1/${modelName}/templates/review`);
  }

  public getReviewTemplate(id: string): Promise<AxiosResponse<ReviewTemplate>> {
    return api.get<ReviewTemplate>(`/api/v1/${modelName}/templates/review/${encodeURIComponent(id)}`);
  }

  public getCustomIds(): Promise<AxiosResponse<CustomId[]>> {
    return api.get<CustomId[]>(`/api/v1/${modelName}/customid/`);
  }

  public getChecklist(): Promise<AxiosResponse<Checklist[]>> {
    return api.get<Checklist[]>(`/api/v1/${modelName}/checklist`);
  }

  public deleteChecklistById(id: string) {
    return api.delete<ISuccessRsponse>(`/api/v1/${modelName}/checklist/${encodeURIComponent(id)}`);
  }

  public createChecklist(item: Checklist) {
    return api.post<ISuccessRsponse>(`/api/v1/${modelName}/checklist`, item);
  }

  public editChecklist(item: Checklist) {
    return api.put<Checklist>(`/api/v1/${modelName}/checklist/${item._key}`, item);
  }

  public createChecklistItem(id: string, item: ChecklistItem) {
    return api.post<Checklist>(`/api/v1/${modelName}/checklist/${id}/items/`, item);
  }

  public editChecklistItem(id: string, item: ChecklistItem) {
    return api.put<Checklist>(`/api/v1/${modelName}/checklist/${id}/items/${item._key}`, item);
  }

  public deleteChecklistItem(id: string, itemId: string) {
    return api.delete<Checklist>(`/api/v1/${modelName}/checklist/${id}/items/${itemId}`);
  }

  public createCustomId(id: CustomId) {
    return api.post<ISuccessRsponse>(`/api/v1/${modelName}/customid/`, id);
  }

  public editCustomId(id: CustomId) {
    return api.put<ISuccessRsponse>(`/api/v1/${modelName}/customid/${id._key}`, id);
  }

  public deleteCustomId(id: string) {
    return api.delete<ISuccessRsponse>(`/api/v1/${modelName}/customid/${id}`);
  }

  public customIdUsage(id: string) {
    return api.get<CustomIdUsage>(`/api/v1/${modelName}/customid/${id}/usage`);
  }

  public getInternalTokens(): Promise<AxiosResponse<InternalToken[]>> {
    return api.get<InternalToken[]>(`/api/v1/${modelName}/internaltoken`);
  }

  public createInternalToken(user: InternalToken) {
    return api.post<ISuccessRsponse>(`/api/v1/${modelName}/internaltoken/`, user);
  }

  public renewInternalToken(tokenId: string) {
    return api.put<InternalToken>(`/api/v1/${modelName}/internaltoken/${tokenId}`);
  }

  public revokeInternalToken(user: string) {
    return api.delete<ISuccessRsponse>(`/api/v1/${modelName}/internaltoken/${user}`);
  }

  public executeDryRun(username: string) {
    return api.get<DeletePersonalDataResponse>(`/api/v1/${modelName}/users/delete-personal-data?username=${username}`);
  }

  public getPersonalDetails(username: string, entity: string) {
    return api.get<{
      success: boolean;
      message: string;
      data: Array<{
        entityID: string;
        entityType: string;
        entityStatus?: string;
        entityName: string;
        disableDeleteReason?: string;
      }>;
    }>(`/api/v1/${modelName}/users/get-personal-details/${username}?entity=${entity}`);
  }

  public deletePersonalDataByEntity(username: string, entity: string) {
    return api.delete<ISuccessRsponse>(
      `/api/v1/${modelName}/users/delete-personal-data/${entity}?username=${username}`,
    );
  }

  public deletePersonalDataByEntityId(entity: string, id: string) {
    return api.delete<ISuccessRsponse>(`/api/v1/${modelName}/users/delete-personal-data/${entity}/${id}`);
  }
}

const adminService = new AdminService();
export default adminService;
