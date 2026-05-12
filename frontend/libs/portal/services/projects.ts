// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useApi} from '@shared/api/useApi';
import {
  Approval,
  ApprovableInfo,
  ApprovalUpdate,
  ApproverRoles,
  CreateApprovalResponse,
} from '@disclosure-portal/model/Approval';
import {
  ExternalApprovalRequest,
  InternalApprovalRequest,
  PlausibilityCheckRequest,
} from '@disclosure-portal/model/ApprovalRequest';
import {Checklist, ExecuteRequest} from '@disclosure-portal/model/Checklist';
import {IFoundResponse, ISuccessRsponse} from '@disclosure-portal/model/Common';
import {Decision} from '@disclosure-portal/model/Decision';
import {DocumentDto} from '@disclosure-portal/model/Document';
import {LicenseRuleRequest} from '@disclosure-portal/model/LicenseRule';
import type {WizardProject} from '@disclosure-portal/model/NewWizard';
import {PolicyDecisionRequest} from '@disclosure-portal/model/PolicyDecision';
import PolicyRule, {PolicyRuleDto} from '@disclosure-portal/model/PolicyRule';
import {
  ApprovableSPDXDto,
  ComponentDetails,
  FillCustomerReq,
  IProjectChildrenMembers,
  MultiProjectMemberRequest,
  ProjectChildrenMemberSuccessResponse,
  ProjectMemberRequest,
  ProjectModel,
  ProjectSettingsModel,
  ProjectSubscriptions,
  ProjectUser,
  Token,
  TokenRequest,
  UserManagement,
} from '@disclosure-portal/model/Project';
import ProjectPostRequest from '@disclosure-portal/model/ProjectPostRequest';
import {
  ProjectChildren,
  ProjectSbomsFlat,
  ProjectsResponse,
  VersionSboms,
} from '@disclosure-portal/model/ProjectsResponse';
import {ReviewTemplate} from '@disclosure-portal/model/ReviewTemplate';
import {UserDto} from '@shared/types/Users';
import {AuditLog, ComponentMultiDiff, NoticeFileFormat} from '@disclosure-portal/model/VersionDetails';
import {IProjectCreatedResponse, WizardProjectPostRequest} from '@disclosure-portal/model/Wizard';
import {SearchOptions} from '@disclosure-portal/utils/Table';
import {AxiosResponse} from 'axios';

const {api, NO_IDLE_PARAM} = useApi();

const modelName = 'projects';

export enum RemarkTypes {
  license = 'license',
  review = 'review',
  scan = 'scan',
}

export enum DocumentDownloadVersion {
  Supplier1 = '0',
  Supplier2 = '1',
  Customer1 = '2',
  Customer2 = '3',
  None = '4',
}

class ProjectService {
  /*
  private validHeadStates(status: number) {
    return status === 200 || status === 404;
  }

   */

  public getAll = (signal?: AbortSignal) =>
    api.get<ProjectsResponse>(`/api/v1/${modelName}`, signal ? {signal} : undefined);

  public getDisclosures = () => api.get<ProjectsResponse>(`/api/v1/disclosures`);

  public getRecent = () => api.get<ProjectsResponse>(`/api/v1/${modelName}/recent`);

  public async getPossibleChildren(projectUid: string) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    return (await api.get<ProjectsResponse>(`/api/v1/${modelName}/${projectUid}/possibleChildren`)).data;
  }

  public async getAllWithOptions(options: SearchOptions, signal?: AbortSignal) {
    return api.post<ProjectsResponse>(`/api/v1/${modelName}/search`, options, signal ? {signal} : undefined);
  }

  public async get(projectUid: string) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    return (await api.get<ProjectModel>(`/api/v1/${modelName}/${projectUid}`)).data;
  }

  public async getChildren(projectUid: string) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    return (await api.get<ProjectChildren>(`/api/v1/${modelName}/${projectUid}/children`)).data;
  }

  public async getApprovableInfo(projectUid: string) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    return (await api.get<ApprovableInfo>(`/api/v1/${modelName}/${projectUid}/approvableinfo`)).data;
  }

  public async getAllApprovals(projectUid: string) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    return (await api.get<Approval[]>(`/api/v1/${modelName}/${projectUid}/approval/list`)).data;
  }

  public async getAllSboms(projectUid: string) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    return (await api.get<VersionSboms[]>(`/api/v1/${modelName}/${projectUid}/allSbomLists`)).data;
  }

  public async getAllSbomsFlat(projectUid: string) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    return (await api.get<ProjectSbomsFlat>(`/api/v1/${modelName}/${projectUid}/allSBOM`)).data;
  }

  public async compareSpdxFiles(
    projectUid: string,
    versionKeyOld: string,
    spdxIdOld: string,
    versionKeyNew: string,
    spdxIdNew: string,
  ) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    return api.get<ComponentMultiDiff[]>(
      `/api/v1/${modelName}/${projectUid}/spdxCompare/${versionKeyOld}/${spdxIdOld}/${versionKeyNew}/${spdxIdNew}`,
    );
  }

  public async create(data: WizardProjectPostRequest): Promise<AxiosResponse<IProjectCreatedResponse>> {
    return api.post<IProjectCreatedResponse>(`/api/v1/${modelName}`, data);
  }

  public async createGroup(data: ProjectPostRequest) {
    return api.post<ProjectPostRequest>(`/api/v1/${modelName}`, data);
  }

  public async createInternalApproval(
    data: InternalApprovalRequest,
    projectUid: string,
  ): Promise<CreateApprovalResponse> {
    return (
      await api.post<CreateApprovalResponse>(
        `api/v1/${modelName}/${encodeURIComponent(projectUid)}/approval/create/internal`,
        data,
      )
    ).data;
  }

  public async createExternalApproval(
    data: ExternalApprovalRequest,
    projectUid: string,
  ): Promise<CreateApprovalResponse> {
    return (
      await api.post<CreateApprovalResponse>(
        `api/v1/${modelName}/${encodeURIComponent(projectUid)}/approval/create/external`,
        data,
      )
    ).data;
  }

  public async createVehicleApproval(
    data: ExternalApprovalRequest,
    projectUid: string,
  ): Promise<CreateApprovalResponse> {
    return (
      await api.post<CreateApprovalResponse>(
        `api/v1/${modelName}/${encodeURIComponent(projectUid)}/approval/create/vehicle`,
        data,
      )
    ).data;
  }

  public async createPlausibilityCheck(
    data: PlausibilityCheckRequest,
    projectUid: string,
  ): Promise<CreateApprovalResponse> {
    return (
      await api.post<CreateApprovalResponse>(
        `api/v1/${modelName}/${encodeURIComponent(projectUid)}/approval/create/plausi`,
        data,
      )
    ).data;
  }

  public async updateApproval(data: ApprovalUpdate, projectUid: string, appId: string) {
    return api.put<ProjectPostRequest>(
      `api/v1/${modelName}/${encodeURIComponent(projectUid)}/approval/${encodeURIComponent(appId)}`,
      data,
    );
  }

  public async getApproval(approvalId: string, projectUid: string) {
    return api.get<Approval>(
      `api/v1/${modelName}/${encodeURIComponent(projectUid)}/approval/${encodeURIComponent(approvalId)}`,
    );
  }

  public async getVehiclePlatform(projectUid: string): Promise<AxiosResponse<IFoundResponse>> {
    return api.get(`api/v1/${modelName}/${encodeURIComponent(projectUid)}/approval/vehiclechildren`);
  }

  public async getVehiclePlatformOnly(projectUid: string): Promise<AxiosResponse<IFoundResponse>> {
    return api.get(`api/v1/${modelName}/${encodeURIComponent(projectUid)}/approval/vehiclechildrenonly`);
  }

  public async getApproverUser(
    approvalId: string,
    projectUid: string,
    approver: ApproverRoles,
  ): Promise<AxiosResponse<UserDto>> {
    return api.get<UserDto>(
      `api/v1/${modelName}/${encodeURIComponent(projectUid)}/approval/${encodeURIComponent(approvalId)}/approver/${encodeURIComponent(approver)}`,
    );
  }

  public async updateProject(data: ProjectPostRequest, projectUid: string) {
    return (await api.put<ProjectModel>(`api/v1/${modelName}/${encodeURIComponent(projectUid)}`, data)).data;
  }

  public getApprovalOrReviewUsage(projectUid: string): Promise<AxiosResponse<ISuccessRsponse>> {
    return api.get<ISuccessRsponse>(`/api/v1/${modelName}/${encodeURIComponent(projectUid)}/approvalOrReviewUsage`);
  }

  public getPendingApprovalOrReviewUsage(projectUid: string, userId: string): Promise<AxiosResponse<ISuccessRsponse>> {
    return api.get<ISuccessRsponse>(
      `/api/v1/${modelName}/${encodeURIComponent(projectUid)}/users/${userId}/pendingApprovalOrReviewUsage`,
    );
  }

  public async delete(projectUid: string) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    return api.delete(`/api/v1/${modelName}/${projectUid}`);
  }

  public async deprecate(projectUid: string) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    return api.put(`/api/v1/${modelName}/${projectUid}/deprecate`);
  }

  public downloadNoticeFileInFormatForSbom(
    format: NoticeFileFormat,
    projectUid: string,
    versionKey: string,
    sbomUuid: string,
  ) {
    return api.get(
      `/api/v1/${modelName}/${encodeURIComponent(projectUid)}/versions/${encodeURIComponent(versionKey)}/notice/${encodeURIComponent(sbomUuid)}/` +
        format,
    );
  }

  public downloadScanOrLicenseRemarksForSbomCsv(
    projectUid: string,
    versionKey: string,
    type: RemarkTypes,
    sbomUuid: string,
  ): Promise<AxiosResponse> {
    return api.get(
      `/api/v1/${modelName}/${encodeURIComponent(projectUid)}/versions/${encodeURIComponent(versionKey)}/quality/${type}remarks/${encodeURIComponent(sbomUuid)}/download`,
      {
        responseType: 'text',
      },
    );
  }

  public downloadReviewRemarksForSbomCsv(projectUid: string, versionKey: string): Promise<AxiosResponse> {
    return api.get(
      `/api/v1/${modelName}/${encodeURIComponent(projectUid)}/versions/${encodeURIComponent(versionKey)}/quality/reviewremarks/download`,
      {
        responseType: 'text',
      },
    );
  }

  public downloadSpdxHistoryFile(projectUid: string, versionKey: string, fileKey: string) {
    return api.get(
      `/api/v1/${modelName}/${encodeURIComponent(projectUid)}/versions/${encodeURIComponent(versionKey)}/spdx/${fileKey}`,
      {
        responseType: 'blob',
      },
    );
  }

  public updateSpdxTag(projectUid: string, versionKey: string, fileKey: string, tag: string) {
    return api.put<ISuccessRsponse>(
      `/api/v1/${modelName}/${encodeURIComponent(projectUid)}/versions/${encodeURIComponent(versionKey)}/spdx/${fileKey}/tag`,
      {tag: tag},
    );
  }

  public toggleSpdxLock(projectUid: string, versionKey: string, fileKey: string) {
    return api.put<ISuccessRsponse>(
      `/api/v1/${modelName}/${encodeURIComponent(projectUid)}/versions/${encodeURIComponent(versionKey)}/spdx/${fileKey}/toggleLock`,
    );
  }

  public deleteSpdx(projectUid: string, versionKey: string, fileKey: string) {
    return api.delete<ISuccessRsponse>(
      `/api/v1/${modelName}/${encodeURIComponent(projectUid)}/versions/${encodeURIComponent(versionKey)}/spdx/${fileKey}`,
    );
  }

  public downloadDocumentByTask(projectUid: string, taskKey: string, type: string, lang: string, docVersion: string) {
    return api.get(
      `/api/v1/${modelName}/${encodeURIComponent(projectUid)}/documents/downloadTask/${encodeURIComponent(taskKey)}/${encodeURIComponent(type)}/${encodeURIComponent(lang)}/${docVersion}`,
      {
        responseType: 'blob',
      },
    );
  }

  public async getPolicyRules(projectUid: string) {
    return api.get<PolicyRuleDto[]>(`/api/v1/${modelName}/${encodeURIComponent(projectUid)}/policyrules`);
  }

  public async getProjectPolicyRule(projectUid: string, ruleId: string) {
    return api.get<PolicyRule>(`/api/v1/${modelName}/${encodeURIComponent(projectUid)}/policyrules/${ruleId}`);
  }

  public async getComponentDetailsForSbom(projectUid: string, versionKey: string, sbomUuid: string, spdxId: string) {
    return api.get<ComponentDetails>(
      `/api/v1/${modelName}/${encodeURIComponent(projectUid)}/versions/${encodeURIComponent(versionKey)}/component/${encodeURIComponent(sbomUuid)}/${encodeURIComponent(spdxId)}`,
    );
  }

  public async getUserManagement(projectUuid: string): Promise<UserManagement> {
    return (await api.get<UserManagement>(`/api/v1/${modelName}/${encodeURIComponent(projectUuid)}/users`)).data;
  }

  public async getChildrenMembers(projectUid: string): Promise<IProjectChildrenMembers> {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    return (await api.get<IProjectChildrenMembers>(`/api/v1/${modelName}/${projectUid}/children/users`)).data;
  }

  public async addProjectMember(projectUid: string, user: ProjectUser, comment: string, responsible: boolean) {
    return api.post(
      `/api/v1/${modelName}/${encodeURIComponent(projectUid)}/users`,
      new ProjectMemberRequest(user.userId, user.userType, comment, responsible),
    );
  }

  public async addProjectChildrenMember(
    projectUid: string,
    user: ProjectUser,
    comment: string,
    responsible: boolean,
    targetProjectKeys: string[],
  ) {
    return api.post<ProjectChildrenMemberSuccessResponse[]>(
      `/api/v1/${modelName}/${encodeURIComponent(projectUid)}/children/users`,
      new MultiProjectMemberRequest(user.userId, user.userType, comment, responsible, targetProjectKeys),
    );
  }

  public async editProjectMember(
    projectUid: string,
    user: ProjectUser,
    oldId: string,
    comment: string,
    responsible: boolean,
  ) {
    return api.put(
      `/api/v1/${modelName}/${encodeURIComponent(projectUid)}/users/${oldId}`,
      new ProjectMemberRequest(user.userId, user.userType, comment, responsible),
    );
  }

  public async deleteProjectMember(projectUid: string, userId: string) {
    return api.delete(`/api/v1/${modelName}/${encodeURIComponent(projectUid)}/users/${userId}`);
  }

  public getUsersBySearchFragment(projectUid: string, searchFragment: string, active?: boolean) {
    searchFragment = encodeURIComponent('' + searchFragment).replace(/\./g, '%2E');
    const params = active !== undefined ? `?active=${active}` : '';
    return api.get<UserDto[]>(
      `/api/v1/${modelName}/${encodeURIComponent(projectUid)}/users/profile/search/${searchFragment}${params}`,
    );
  }

  public async getTokens(projectUuid: string): Promise<Token[]> {
    return (await api.get<Token[]>(`/api/v1/${modelName}/${encodeURIComponent(projectUuid)}/tokens`)).data;
  }

  public async getAuditTrail(projectUuid: string): Promise<AuditLog[]> {
    return (await api.get<AuditLog[]>(`/api/v1/${modelName}/${encodeURIComponent(projectUuid)}/audit`)).data;
  }

  public async getDocuments(projectUuid: string, noIdleAnimation?: boolean): Promise<DocumentDto[]> {
    const noIdle = noIdleAnimation ? `?${NO_IDLE_PARAM}=true` : '';
    return (await api.get<DocumentDto[]>(`/api/v1/${modelName}/${encodeURIComponent(projectUuid)}/documents${noIdle}`))
      .data;
  }

  public async addProjectToken(projectUid: string, token: Token) {
    return api.post<Token>(
      `/api/v1/${modelName}/${encodeURIComponent(projectUid)}/tokens`,
      new TokenRequest(token.company, token.description, token.expiry, ''),
    );
  }

  public async revokeProjectToken(projectUid: string, token: string) {
    return api.delete(`/api/v1/${modelName}/${encodeURIComponent(projectUid)}/tokens/${token}`);
  }

  public async renewProjectToken(projectUid: string, token: string) {
    return api.put<Token>(`/api/v1/${modelName}/${encodeURIComponent(projectUid)}/tokens/${token}`);
  }

  public async updateProjectDocumentMeta(settings: ProjectSettingsModel, projectUid: string) {
    return api.put(`api/v1/${modelName}/${encodeURIComponent(projectUid)}/settings`, settings);
  }

  public async updateApprovableSpdx(approvableSpdx: ApprovableSPDXDto, projectUid: string) {
    return api.put(`api/v1/${modelName}/${encodeURIComponent(projectUid)}/approvableSPDX`, approvableSpdx);
  }

  public async fillCustomer(data: FillCustomerReq, projectUid: string, appId: string) {
    return api.post(
      `api/v1/${modelName}/${encodeURIComponent(projectUid)}/approval/${encodeURIComponent(appId)}/fillCustomer`,
      data,
    );
  }

  public async saveProjectSubscriptions(projectUid: string, data: ProjectSubscriptions): Promise<ProjectSubscriptions> {
    return (
      await api.put<ProjectSubscriptions>(`api/v1/${modelName}/${encodeURIComponent(projectUid)}/subscriptions`, data)
    ).data;
  }

  public async getReviewTemplates(projectUid: string): Promise<AxiosResponse<ReviewTemplate[]>> {
    return await api.get<ReviewTemplate[]>(`/api/v1/${modelName}/${encodeURIComponent(projectUid)}/templates/review`);
  }

  public async getReviewTemplate(projectUid: string, id: string): Promise<AxiosResponse<ReviewTemplate>> {
    return await api.get<ReviewTemplate>(
      `/api/v1/${modelName}/${encodeURIComponent(projectUid)}/templates/review/${encodeURIComponent(id)}`,
    );
  }

  async getDecisions(projectUid: string) {
    return api.get<Decision[]>(`/api/v1/${modelName}/${encodeURIComponent(projectUid)}/decisions`);
  }

  async createLicenseRule(
    projectUid: string,
    versionKey: string,
    data: LicenseRuleRequest,
  ): Promise<AxiosResponse<ISuccessRsponse>> {
    return api.post<ISuccessRsponse>(
      `/api/v1/${modelName}/${encodeURIComponent(projectUid)}/versions/${encodeURIComponent(versionKey)}/createLicenseRule`,
      data,
    );
  }

  async cancelLicenseRule(projectUid: string, licenseRuleId: string) {
    return api.put<ISuccessRsponse>(
      `/api/v1/${modelName}/${encodeURIComponent(projectUid)}/licenserules/${encodeURIComponent(licenseRuleId)}/cancel`,
    );
  }

  async getJobByKey(projectUid: string, jobId: string) {
    return api.get<{status: number}>(
      `/api/v1/${modelName}/${encodeURIComponent(projectUid)}/jobs/${encodeURIComponent(jobId)}`,
    );
  }

  public async cloneProject(projectUid: string, count: number = 1): Promise<IProjectCreatedResponse> {
    return (
      await api.post<IProjectCreatedResponse>(`/api/v1/${modelName}/${encodeURIComponent(projectUid)}/clone`, {count})
    ).data;
  }

  async executeChecklistChecks(
    projectUid: string,
    versionKey: string,
    sbomUuid: string,
    req: ExecuteRequest,
  ): Promise<AxiosResponse<ISuccessRsponse>> {
    return api.post<ISuccessRsponse>(
      `/api/v1/${modelName}/${encodeURIComponent(projectUid)}/versions/${encodeURIComponent(versionKey)}/checklists/${encodeURIComponent(sbomUuid)}`,
      req,
    );
  }

  async getApplicableChecklists(projectUid: string) {
    return api.get<Checklist[]>(`/api/v1/${modelName}/${encodeURIComponent(projectUid)}/checklists`);
  }

  public async previewProjectWizard(data: WizardProject) {
    return (await api.post<WizardProject>(`/api/v1/${modelName}/wizard/preview`, data)).data;
  }

  public async createProjectWizard(data: WizardProject) {
    return (await api.post<ProjectModel>(`/api/v1/${modelName}/wizard`, data)).data;
  }

  public async updateProjectWizard(data: WizardProject, projectUid: string) {
    return api.put<ProjectModel>(`/api/v1/${modelName}/wizard/${encodeURIComponent(projectUid)}`, data);
  }

  public async getWizardByProjectKey(projectUid: string) {
    return (await api.get<Partial<WizardProject>>(`/api/v1/${modelName}/wizard/${encodeURIComponent(projectUid)}`))
      .data;
  }

  public async previewGroupWizard(data: WizardProject) {
    return (await api.post<WizardProject>(`/api/v1/${modelName}/group-wizard/preview`, data)).data;
  }

  public async createGroupWizard(data: WizardProject) {
    return (await api.post<ProjectModel>(`/api/v1/${modelName}/group-wizard`, data)).data;
  }

  async createPolicyDecision(
    projectUid: string,
    versionKey: string,
    data: PolicyDecisionRequest,
  ): Promise<AxiosResponse<ISuccessRsponse>> {
    return api.post<ISuccessRsponse>(
      `/api/v1/${modelName}/${encodeURIComponent(projectUid)}/versions/${encodeURIComponent(versionKey)}/createPolicyDecision`,
      data,
    );
  }

  async createBulkPolicyDecision(
    projectUid: string,
    versionKey: string,
    data: PolicyDecisionRequest[],
  ): Promise<AxiosResponse<ISuccessRsponse>> {
    return api.post<ISuccessRsponse>(
      `/api/v1/${modelName}/${encodeURIComponent(projectUid)}/versions/${encodeURIComponent(versionKey)}/createBulkPolicyDecision`,
      data,
    );
  }

  async cancelPolicyDecision(projectUid: string, policyDecisionId: string) {
    return api.put<ISuccessRsponse>(
      `/api/v1/${modelName}/${encodeURIComponent(projectUid)}/policyDecisions/${encodeURIComponent(policyDecisionId)}/cancel`,
    );
  }
}
const projectService = new ProjectService();
export default projectService;
