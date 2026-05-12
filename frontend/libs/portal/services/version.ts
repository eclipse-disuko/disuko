// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useApi} from '@shared/api/useApi';
import {ISuccessRsponse} from '@disclosure-portal/model/Common';
import ExternalSourcePostRequest from '@disclosure-portal/model/ExternalSourcePostRequest';
import {ComponentLicenses, ProjectModel, SbomLicenses} from '@disclosure-portal/model/Project';
import ProjectVersionPostRequest from '@disclosure-portal/model/ProjectVersionPostRequest';
import {
  CommentReviewRemarkRequest,
  LicenseRemarks,
  ReviewRemark,
  ReviewRemarkRequest,
  ScanRemark,
  SetReviewRemarkStatusRequest,
} from '@disclosure-portal/model/Quality';
import {
  AuditLog,
  ComponentInfoSlim,
  ComponentsInfoResponse,
  ExternalSource,
  GeneralStats,
  OverallReviewRequest,
  SbomStats,
  SpdxFile,
  VersionSlimDto,
} from '@disclosure-portal/model/VersionDetails';
import {AxiosResponse} from 'axios';
import {BulkSetReviewRemarkStatusRequest} from '../model/ReviewRemarkBulkOperations';

const modelName = 'projects';
const {api} = useApi();

class VersionService {
  public async getVersion(projectUid: string, versionKey: string | null) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    versionKey = encodeURIComponent('' + versionKey);
    return api.get<VersionSlimDto>(`/api/v1/${modelName}/${projectUid}/versions/${versionKey}`);
  }

  public async getVersionWithProject(projectUid: string, versionKey: string | null) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    versionKey = encodeURIComponent('' + versionKey);
    return api.get<ProjectModel>(`/api/v1/${modelName}/${projectUid}/versions/${versionKey}?withProject=true`);
  }

  public async getSbomHistory(projectUid: string, versionKey: string | null, limit = -1) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    versionKey = encodeURIComponent('' + versionKey);
    return api.get<SpdxFile[]>(`/api/v1/${modelName}/${projectUid}/versions/${versionKey}/sbomhistory?limit=${limit}`);
  }

  public async getVersionComponentsForSbom(projectUid: string, versionKey: string | null, sbomUuid: string | null) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    versionKey = encodeURIComponent('' + versionKey);
    sbomUuid = encodeURIComponent('' + sbomUuid);
    const response = await api.get<ComponentsInfoResponse>(
      `/api/v1/${modelName}/${projectUid}/versions/${versionKey}/components/${sbomUuid}`,
    );
    return response.data;
  }

  public async getVersionComponentsBySearch(
    projectUid: string,
    versionKey: string | null,
    sbomUuid: string | null,
    searchFragment: string | null,
  ) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    versionKey = encodeURIComponent('' + versionKey);
    sbomUuid = encodeURIComponent('' + sbomUuid);
    if (!searchFragment) {
      const response = await api.get<ComponentInfoSlim[]>(
        `/api/v1/${modelName}/${projectUid}/versions/${versionKey}/components/${sbomUuid}`,
      );
      return response.data;
    } else {
      searchFragment = encodeURIComponent('' + searchFragment);
      const response = await api.get<ComponentInfoSlim[]>(
        `/api/v1/${modelName}/${projectUid}/versions/${versionKey}/components/${sbomUuid}/${searchFragment}`,
      );
      return response.data;
    }
  }

  public async getVersionComponentsLicenses(
    projectUid: string,
    versionKey: string | null,
    sbomUuid: string | null,
    spdxId: string | null,
  ) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    versionKey = encodeURIComponent('' + versionKey);
    sbomUuid = encodeURIComponent('' + sbomUuid);
    spdxId = encodeURIComponent('' + spdxId);
    const response = await api.get<ComponentLicenses>(
      `/api/v1/${modelName}/${projectUid}/versions/${versionKey}/component/${sbomUuid}/${spdxId}/licenses`,
    );
    return response.data;
  }

  public async getVersionSbomAllLicenses(projectUid: string, versionKey: string | null, sbomUuid: string | null) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    versionKey = encodeURIComponent('' + versionKey);
    sbomUuid = encodeURIComponent('' + sbomUuid);
    const response = await api.get<SbomLicenses>(
      `/api/v1/${modelName}/${projectUid}/versions/${versionKey}/components/${sbomUuid}/licenses`,
    );
    return response.data;
  }

  public async getScanRemarksForSbom(projectUid: string, versionKey: string, sbomUuid: string) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    versionKey = encodeURIComponent('' + versionKey);
    sbomUuid = encodeURIComponent('' + sbomUuid);
    const response = await api.get<ScanRemark[]>(
      `/api/v1/${modelName}/${projectUid}/versions/${versionKey}/quality/scanremarks/${sbomUuid}`,
    );
    return response.data;
  }

  public async getLicenseRemarksForSbom(projectUid: string, versionKey: string, sbomUuid: string) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    versionKey = encodeURIComponent('' + versionKey);
    sbomUuid = encodeURIComponent('' + sbomUuid);
    const response = await api.get<LicenseRemarks[]>(
      `/api/v1/${modelName}/${projectUid}/versions/${versionKey}/quality/licenseremarks/${sbomUuid}`,
    );
    return response.data;
  }

  public async getReviewRemarks(projectUid: string, versionKey: string) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    versionKey = encodeURIComponent('' + versionKey);
    return api.get<ReviewRemark[]>(`/api/v1/${modelName}/${projectUid}/versions/${versionKey}/quality/reviewremarks`);
  }

  public async getReviewRemarksForComponent(projectUid: string, versionKey: string, sbomUuid: string, spdxId: string) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    versionKey = encodeURIComponent('' + versionKey);
    sbomUuid = encodeURIComponent('' + sbomUuid);
    spdxId = encodeURIComponent('' + spdxId);
    return api.get<ReviewRemark[]>(
      `/api/v1/${modelName}/${projectUid}/versions/${versionKey}/component/${sbomUuid}/${spdxId}/reviewremarks`,
    );
  }

  public async createReviewRemark(projectUid: string, versionKey: string, data: ReviewRemarkRequest) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    versionKey = encodeURIComponent('' + versionKey);
    return api.post<ISuccessRsponse>(
      `/api/v1/${modelName}/${projectUid}/versions/${versionKey}/quality/reviewremarks`,
      data,
    );
  }

  public async editReviewRemark(projectUid: string, versionKey: string, remarkKey: string, data: ReviewRemarkRequest) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    versionKey = encodeURIComponent('' + versionKey);
    remarkKey = encodeURIComponent('' + remarkKey);
    return api.put<ISuccessRsponse>(
      `/api/v1/${modelName}/${projectUid}/versions/${versionKey}/quality/reviewremarks/${remarkKey}`,
      data,
    );
  }

  public async setReviewRemarkStatus(
    projectUid: string,
    versionKey: string,
    remarkKey: string,
    data: SetReviewRemarkStatusRequest,
  ) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    versionKey = encodeURIComponent('' + versionKey);
    remarkKey = encodeURIComponent('' + remarkKey);
    return api.put<ISuccessRsponse>(
      `/api/v1/${modelName}/${projectUid}/versions/${versionKey}/quality/reviewremarks/${remarkKey}/status`,
      data,
    );
  }

  public async bulkSetReviewRemarkStatus(
    projectUid: string,
    versionKey: string,
    data: BulkSetReviewRemarkStatusRequest,
  ) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    versionKey = encodeURIComponent('' + versionKey);
    return api.post<ISuccessRsponse>(
      `/api/v1/${modelName}/${projectUid}/versions/${versionKey}/quality/reviewremarks/bulk-status`,
      data,
    );
  }

  public async commentReviewRemark(
    projectUid: string,
    versionKey: string,
    remarkKey: string,
    data: CommentReviewRemarkRequest,
  ) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    versionKey = encodeURIComponent('' + versionKey);
    remarkKey = encodeURIComponent('' + remarkKey);
    return api.post<ISuccessRsponse>(
      `/api/v1/${modelName}/${projectUid}/versions/${versionKey}/quality/reviewremarks/${remarkKey}/comments`,
      data,
    );
  }

  public async createOverallReview(projectUid: string, versionKey: string, data: OverallReviewRequest) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    versionKey = encodeURIComponent('' + versionKey);
    return api.post<ISuccessRsponse>(`/api/v1/${modelName}/${projectUid}/versions/${versionKey}/overallreview`, data);
  }

  public async createVersion(projectUid: string, data: ProjectVersionPostRequest) {
    return api.post<ProjectModel>(`/api/v1/${modelName}/${encodeURIComponent(projectUid)}/versions`, data);
  }

  public async updateProjectVersion(projectUid: string, versionKey: string, data: ProjectVersionPostRequest) {
    return api.put<ProjectModel>(
      `/api/v1/${modelName}/${encodeURIComponent(projectUid)}/versions/${encodeURIComponent(versionKey)}`,
      data,
    );
  }

  public getApprovalOrReviewUsage(projectUid: string, versionKey: string): Promise<AxiosResponse<ISuccessRsponse>> {
    return api.get<ISuccessRsponse>(
      `/api/v1/${modelName}/${encodeURIComponent(projectUid)}/versions/${encodeURIComponent(versionKey)}/approvalOrReviewUsage`,
    );
  }

  public async deleteVersion(projectUid: string, versionKey: string) {
    projectUid = encodeURIComponent('' + projectUid);
    return api.delete(`/api/v1/${modelName}/${projectUid}/versions/${encodeURIComponent(versionKey)}`);
  }

  public async getAuditTrail(projectUuid: string, versionUuid: string): Promise<AuditLog[]> {
    projectUuid = encodeURIComponent('' + projectUuid).replace(/\./g, '%2E');
    versionUuid = encodeURIComponent('' + versionUuid);
    return (await api.get<AuditLog[]>(`/api/v1/${modelName}/${projectUuid}/versions/${versionUuid}/audit`)).data;
  }

  public async getExternalSources(projectUid: string, versionKey: string | null) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    versionKey = encodeURIComponent('' + versionKey);
    return api.get<ExternalSource[]>(`/api/v1/${modelName}/${projectUid}/versions/${versionKey}/externalsources`);
  }

  public async deleteExternalSource(externalSource: string, projectUid: string, versionKey: string) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    versionKey = encodeURIComponent('' + versionKey);
    return api.delete(`/api/v1/${modelName}/${projectUid}/versions/${versionKey}/externalsources/${externalSource}`);
  }

  public async createExternalSource(data: ExternalSourcePostRequest, projectUid: string, versionKey: string) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    versionKey = encodeURIComponent('' + versionKey);
    return api.post<ProjectModel>(`/api/v1/${modelName}/${projectUid}/versions/${versionKey}/externalsources`, data);
  }

  public async updateExternalSource(
    externalSource: string,
    data: ExternalSourcePostRequest,
    projectUid: string,
    versionKey: string,
  ) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    versionKey = encodeURIComponent('' + versionKey);
    return api.put<ProjectModel>(
      `/api/v1/${modelName}/${projectUid}/versions/${versionKey}/externalsources/${externalSource}`,
      data,
    );
  }

  public async getSBOMStats(projectUid: string, versionKey: string, sbomUuid: string, fossOnly = false) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    versionKey = encodeURIComponent('' + versionKey);
    sbomUuid = encodeURIComponent('' + sbomUuid);
    return api.get<SbomStats>(
      `/api/v1/${modelName}/${projectUid}/versions/${versionKey}/components/${sbomUuid}/stats?${fossOnly ? `&fossOnly=${fossOnly}` : ''}`,
    );
  }

  public async getGeneralVersionStats(projectUid: string, versionKey: string) {
    projectUid = encodeURIComponent('' + projectUid).replace(/\./g, '%2E');
    versionKey = encodeURIComponent('' + versionKey);
    return api.get<GeneralStats>(`/api/v1/${modelName}/${projectUid}/versions/${versionKey}/stats`);
  }
}

const versionService = new VersionService();
export default versionService;
