import type {ExternalSourceCodeReference} from '@cli/models/ExternalSourceCode';
import type {PolicyRule} from '@cli/models/PolicyRule';
import type {Project} from '@cli/models/Project';
import type {Remark, RRCommentExternDTO} from '@cli/models/ReviewRemark';
import type {
  SBOM,
  SbomDetails,
  SbomStatusPublicResponseDto,
  SpdxStatusInformation,
  SpdxUploadResponse,
} from '@cli/models/Sbom';
import type {VersionDetails, VersionRequest, VersionStatusPublicResponse} from '@cli/models/Version';
import {NoticeFileFormat} from '@cli/models/Version';
import {getApi} from '@cli/api';
import {isAxiosError} from 'axios';
import {authService} from '@cli/services/authService';

const {api: axios} = getApi();

class ProjectService {
  private static instance: ProjectService;

  private constructor() {}

  public static getInstance(): ProjectService {
    if (!ProjectService.instance) {
      ProjectService.instance = new ProjectService();
    }
    return ProjectService.instance;
  }

  private async fetchAuthInfo(): Promise<{projectUuid?: string; isGroup?: boolean} | null> {
    // Prefer cached auth info from AuthService (avoids duplicate /auth/info calls during login)
    const cachedUuid = authService.getAuth();
    const cachedIsGroup = authService.getIsGroup();
    if (cachedUuid && cachedIsGroup !== null) {
      return {projectUuid: cachedUuid, isGroup: cachedIsGroup};
    }

    // Fall back to backend once (AuthService will cache the result)
    const fetchedUuid = await authService.fetchAuthInfo();
    if (!fetchedUuid) {
      return null;
    }
    const fetchedIsGroup = authService.getIsGroup();
    return {
      projectUuid: fetchedUuid,
      isGroup: fetchedIsGroup === null ? undefined : fetchedIsGroup,
    };
  }

  /**
   * Fetches a project by its UUID
   * @param projectUuid The project UUID
   * @returns Project or null if there's an error
   */
  /**
   * Fetch basic project information
   */
  private async fetchProjectInfo(projectUuid: string): Promise<Project> {
    const response = await axios.get<Project>(`/v1/projects/${projectUuid}`);

    // Check if response exists and has data
    if (!response || !response.data) {
      throw new Error('No project data received');
    }

    // Validate essential project fields
    const {name, uuid} = response.data;
    if (!name || !uuid) {
      throw new Error('Invalid project data');
    }

    return response.data;
  }

  /**
   * Fetch project status and version statuses
   */
  private async fetchProjectStatus(projectUuid: string): Promise<{
    status: string;
    versionStatus: VersionStatusPublicResponse[];
  }> {
    const response = await axios.get<{status: string; versionStatus: VersionStatusPublicResponse[]}>(
      `/v1/projects/${projectUuid}/status`,
    );
    return response.data || {status: 'inactive', versionStatus: []};
  }

  private async fetchProjectVersionNames(projectUuid: string): Promise<string[]> {
    const response = await axios.get<string[]>(`/v1/projects/${projectUuid}/versions`);
    return response.data || [];
  }

  /**
   * Fetch and prepare version details for a project
   */
  private async fetchAndPrepareVersions(
    projectUuid: string,
    versionNames: string[],
    versionStatus: VersionStatusPublicResponse[],
  ): Promise<VersionDetails[]> {
    const versionDetailsPromises = versionNames.map(async (versionName) => {
      const detailsResponse = await axios.get<VersionDetails>(
        `/v1/projects/${projectUuid}/versions/${encodeURIComponent(versionName)}`,
      );
      // Optionally, merge status info from versionStatus
      const statusInfo = versionStatus?.find((vs) => vs.name === versionName);
      return {
        ...detailsResponse.data,
        status: statusInfo?.status || detailsResponse.data.status || 'unreviewed',
        lastSbomUploaded: statusInfo?.lastSbomUploaded || detailsResponse.data.lastSbomUploaded,
      };
    });
    return await Promise.all(versionDetailsPromises);
  }

  /**
   * Helper for consistent error logging and processing
   */
  private logError(context: string, error: unknown): string {
    console.error(`Error in ${context}:`, error);
    if (isAxiosError(error)) {
      if (error.response?.data) {
        const apiError = error.response.data;
        if (error.response.status === 404 && context.includes('project')) {
          return 'Project uuid wrong';
        }
        return apiError.message || apiError.raw || `Error in ${context}`;
      }
      // Handle network errors
      return error.message || `Network error in ${context}`;
    }
    return error instanceof Error ? error.message : `Error in ${context}`;
  }

  public async getProject(
    projectUuid: string,
    options?: {parentProjectUuid?: string | null},
  ): Promise<(Project & {status: string; versions: VersionDetails[]; children?: Project[]}) | null> {
    try {
      const authInfo = await this.fetchAuthInfo();

      let project: Project;
      try {
        project = await this.fetchProjectInfo(projectUuid);
      } catch (error) {
        this.logError('getProject - fetchProjectInfo', error);
        throw error;
      }

      const projectStatus = await this.fetchProjectStatus(projectUuid);
      const isDeprecatedProject = (projectStatus.status || '').toLowerCase() === 'deprecated';
      const isChildProject = Boolean(options?.parentProjectUuid && options.parentProjectUuid !== projectUuid);
      const authProjectMismatch = Boolean(authInfo?.projectUuid && authInfo.projectUuid !== projectUuid);
      const isGroup = isChildProject || authProjectMismatch ? false : (authInfo?.isGroup ?? project.isGroup);

      if (!project || !projectStatus) {
        this.logError('getProject', 'Missing project or status data');
        return null;
      }

      // If no versions, return project with empty versions
      let versions: VersionDetails[] = [];
      let allProjectSboms: SBOM[] = [];

      if (!isGroup && !isDeprecatedProject) {
        const versionNames = await this.fetchProjectVersionNames(projectUuid);
        if (versionNames.length === 0) {
          project.sboms = allProjectSboms;
        } else {
          versions = await this.fetchAndPrepareVersions(projectUuid, versionNames, projectStatus.versionStatus);

          // Fetch SBOM list for each version
          const versionsWithSboms = await Promise.all(
            versions.map(async (version) => {
              const sbomList = await this.getProjectVersionSBOMs(projectUuid, version.name);
              // Add version information and fetch details for each SBOM
              if (sbomList && sbomList.length > 0) {
                // Fetch details for all SBOMs in parallel
                type RawSbom = {
                  id?: string;
                  name?: string;
                  updated?: string;
                  valid?: boolean;
                  spdxVersion?: string;
                };
                const sbomsWithDetails: SBOM[] = await Promise.all(
                  (sbomList as RawSbom[]).map(async (sbom) => {
                    const sbomId = sbom.id;
                    let details: SbomDetails | null = null;
                    if (sbomId) {
                      details = await this.getSBOMDetails(projectUuid, version.name, sbomId);
                    }
                    const mapped = {
                      name: sbom.name || details?.name || sbomId || 'unknown',
                      updated: sbom.updated || details?.uploaded || new Date().toISOString(),
                      valid: sbom.valid ?? true,
                      id: sbomId || details?.id || 'unknown',
                      version: version.name,
                      details: {
                        name: details?.name || sbom.name || 'unknown',
                        id: details?.id || sbomId || 'unknown',
                        version: details?.version || sbom.spdxVersion || 'SPDX',
                        creators: details?.creators || 'unknown',
                        created: details?.created || new Date().toISOString(),
                        uploaded: details?.uploaded || new Date().toISOString(),
                        status: details?.status ?? true,
                        tag: details?.tag,
                        isLocked: details?.isLocked,
                        isRetain: (details as any)?.isRetain ?? (details as any)?.isToRetain,
                      },
                    };
                    return mapped;
                  }),
                );
                allProjectSboms.push(...sbomsWithDetails);
              }
              return {
                ...version,
              };
            }),
          );
          versions = versionsWithSboms;
          // Add consolidated SBOM list to project
          project.sboms = allProjectSboms; // aggregated SBOM list across versions
        }
      }

      let children: Project[] | undefined = undefined;
      if (isGroup) {
        try {
          children = (await this.getChildrenProjects(projectUuid)) || undefined;
        } catch (err) {
          this.logError('getChildrenProjects (called from getProject)', {
            error: err,
            uuid: projectUuid,
            isGroup,
          });
        }
      }

      const projectResult = {
        ...project,
        status: projectStatus.status,
        versions,
        ...(isGroup ? {isGroup: true} : {}),
        ...(children ? {children: children.map((child) => ({...child}))} : {}),
      } as Project & {status: string; versions: VersionDetails[]; children?: Project[]};
      return projectResult;
    } catch (error) {
      this.logError('getProject', error);
      // Re-throw the original error to preserve error structure
      throw error;
    }
  }

  /**
   * Get project versions (channels)
   * @param projectUuid The project UUID
   * @returns Array of VersionDetails or null if there's an error
   */
  public async getProjectVersions(projectUuid: string): Promise<VersionDetails[] | null> {
    try {
      const response = await axios.get<VersionDetails[]>(`/v1/projects/${projectUuid}/versions`);
      return response.data;
    } catch (error) {
      console.error('Error fetching project versions:', error);
      return null;
    }
  }

  /**
   * Create a new project version (channel)
   * @param projectUuid The project UUID
   * @param versionRequest The version details to create
   * @returns VersionDetails or null if there's an error
   */
  public async createProjectVersion(
    projectUuid: string,
    versionRequest: VersionRequest,
  ): Promise<VersionDetails | null> {
    try {
      const response = await axios.post<VersionDetails>(
        `/v1/projects/${projectUuid}/versions`,
        versionRequest,
      );
      return response.data;
    } catch (error) {
      console.error('Error creating project version:', error);
      return null;
    }
  }

  /**
   * Upload SBOM as SPDX
   * @param projectUuid The project UUID
   * @param version The version identifier
   * @param file The SPDX file to upload
   * @returns SpdxUploadResponse or null if there's an error
   */
  public async uploadSBOM(projectUuid: string, version: string, file: File): Promise<SpdxUploadResponse | null> {
    try {
      const formData = new FormData();
      formData.append('file', file);

      const response = await axios.post<SpdxUploadResponse>(
        `/v1/projects/${projectUuid}/versions/${version}/sboms`,
        formData,
        {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
        },
      );
      return response.data;
    } catch (error) {
      console.error('Error uploading SBOM:', error);
      if (isAxiosError(error) && error.response?.data) {
        const data = error.response.data as Partial<SpdxUploadResponse> & {
          message?: string;
          raw?: string;
          code?: string;
          reqID?: string;
        };
        // Ensure we return an object with docIsValid false and captured validationFailedMessage if provided
        return {
          docIsValid: false,
          validationFailedMessage: data.validationFailedMessage, // only treat as validation details if explicitly provided
          hash: data.hash,
          fileUploaded: data.fileUploaded,
          id: data.id,
          sbomguid: data.sbomguid,
          code: data.code,
          message: data.message || data.raw || 'Upload failed',
          reqID: data.reqID,
          raw: data.raw,
        } as SpdxUploadResponse;
      }
      return null;
    }
  }

  /**
   * Get SBOM status information
   * @param projectUuid The project UUID
   * @param version The version identifier
   * @param sbomUuid The SBOM UUID
   * @returns SpdxStatusInformation or null if there's an error
   */
  public async getSBOMStatus(
    projectUuid: string,
    version: string,
    sbomUuid: string,
  ): Promise<SbomStatusPublicResponseDto | null> {
    try {
      const response = await axios.get<SbomStatusPublicResponseDto>(
        `/v1/projects/${projectUuid}/versions/${version}/sboms/${sbomUuid}/status`,
      );
      return response.data;
    } catch (error) {
      console.error('Error fetching SBOM status:', error);
      return null;
    }
  }

  /**
   * Get list of all SBOM deliveries for a project version
   * @param projectUuid The project UUID
   * @param version The version identifier
   * @returns Array of SBOM deliveries or null if there's an error
   */
  public async getProjectVersionSBOMs(
    projectUuid: string,
    version: string,
  ): Promise<SbomStatusPublicResponseDto[] | null> {
    try {
      const response = await axios.get<SbomStatusPublicResponseDto[]>(
        `/v1/projects/${projectUuid}/versions/${version}/sboms`,
      );

      const sbomsWithVersion = response.data.map((sbom) => ({
        ...sbom,
        version,
      }));

      return sbomsWithVersion;
    } catch (error) {
      console.error('Error fetching SBOM deliveries:', error);
      return null;
    }
  }

  /**
   * Get detailed information for a specific SBOM delivery
   */
  private async getSBOMDetails(projectUuid: string, version: string, sbomUuid: string): Promise<SbomDetails | null> {
    try {
      const response = await axios.get<SbomDetails>(
        `/v1/projects/${projectUuid}/versions/${version}/sboms/${sbomUuid}`,
      );
      return response.data;
    } catch (error) {
      console.error('Error fetching SBOM details:', error);
      return null;
    }
  }

  /**
   * Get policy rules of project
   */
  public async getProjectPolicyRules(projectUuid: string): Promise<PolicyRule[] | null> {
    try {
      const response = await axios.get<PolicyRule[]>(`/v1/projects/${projectUuid}/policyrules`);
      return response.data;
    } catch (error) {
      console.error('Error fetching policy rules:', error);
      return null;
    }
  }

  /**
   * Check SBOM status information before upload
   * @param projectUuid The project UUID
   * @param file The SPDX file to check
   * @returns SpdxStatusInformation or null if there's an error
   */
  public async checkSBOM(projectUuid: string, file: File): Promise<SpdxStatusInformation | null> {
    try {
      const formData = new FormData();
      formData.append('file', file);

      const response = await axios.post<SpdxStatusInformation>(
        `/v1/projects/${projectUuid}/sbomcheck`,
        formData,
        {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
        },
      );
      return response.data;
    } catch (error) {
      console.error('Error checking SBOM:', error);
      return null;
    }
  }

  /**
   * Get review remarks for version
   */
  public async getVersionReviewRemarks(projectUuid: string, version: string): Promise<Remark[] | null> {
    try {
      const response = await axios.get<Remark[]>(
        `/v2/projects/${projectUuid}/versions/${version}/reviewremarks`,
      );
      return response.data;
    } catch (error) {
      console.error('Error fetching review remarks:', error);
      return null;
    }
  }

  /**
   * Comment on a review remark
   */
  public async commentOnReviewRemark(
    projectUuid: string,
    version: string,
    reviewRemarkUuid: string,
    comment: RRCommentExternDTO,
  ): Promise<void> {
    try {
      await axios.post(
        `/v1/projects/${projectUuid}/versions/${version}/reviewremarks/${reviewRemarkUuid}`,
        comment,
      );
    } catch (error) {
      console.error('Error commenting on review remark:', error);
      throw error;
    }
  }

  /**
   * Delete version of project
   */
  public async deleteProjectVersion(projectUuid: string, version: string): Promise<void> {
    try {
      await axios.delete(`/v1/projects/${projectUuid}/versions/${version}`);
    } catch (error) {
      console.error('Error deleting project version:', error);
      throw error;
    }
  }

  /**
   * Lock an SBOM
   * @param projectUuid The project UUID
   * @param version The version identifier
   * @param sbomId The SBOM identifier
   * @returns Promise<void>
   */
  public async lockSbom(projectUuid: string, version: string, sbomId: string): Promise<void> {
    try {
      await axios.put(`/v1/projects/${projectUuid}/versions/${version}/sboms/${sbomId}/lock`, {});
    } catch (error) {
      console.error('Error locking SBOM:', error);
      throw error;
    }
  }

  /**
   * Unlock an SBOM
   * @param projectUuid The project UUID
   * @param version The version identifier
   * @param sbomId The SBOM identifier
   * @returns Promise<void>
   */
  public async unlockSbom(projectUuid: string, version: string, sbomId: string): Promise<void> {
    try {
      await axios.put(`/v1/projects/${projectUuid}/versions/${version}/sboms/${sbomId}/unlock`, {});
    } catch (error) {
      console.error('Error unlocking SBOM:', error);
      throw error;
    }
  }

  /**
   * Update the tag of an SBOM
   * @param projectUuid The project UUID
   * @param version The version identifier
   * @param sbomId The SBOM identifier
   * @param tag The new tag value
   * @returns Promise<void>
   */
  public async updateSbomTag(projectUuid: string, version: string, sbomId: string, tag: string): Promise<void> {
    try {
      await axios.put(`/v1/projects/${projectUuid}/versions/${version}/sboms/${sbomId}/tag`, {tag});
    } catch (error) {
      console.error('Error updating SBOM tag:', error);
      throw error;
    }
  }

  /**
   * Get SBOM status information for a specific version/channel
   * @param projectUuid The project UUID
   * @param version The version/channel identifier
   * @param sbomId The SBOM identifier
   * @returns Promise<SpdxStatusInformation> SBOM status information
   */
  public async getSbomStatus(projectUuid: string, version: string, sbomId: string): Promise<SpdxStatusInformation> {
    try {
      const response = await axios.get<SpdxStatusInformation>(
        `/v1/projects/${projectUuid}/versions/${version}/sboms/${sbomId}/check`,
      );

      if (!response || !response.data) {
        throw new Error('No SBOM status data received');
      }

      return response.data;
    } catch (error) {
      console.error('Error getting SBOM status:', error);
      throw error;
    }
  }

  /**
   * Download notice file for a specific SBOM in the requested format.
   * Backend endpoints:
   *  /notice/html  -> HTML formatted (text/html)
   *  /notice/json  -> JSON structure (application/json)
   *  /notice/text  -> Plain text (text/plain)
   * @param projectUuid Project UUID
   * @param version Version/channel identifier
   * @param sbomId SBOM identifier
   * @param format NoticeFileFormat (html | json | plain)
   * @returns Raw string content of the notice file or null if failed
   */
  public async downloadNoticeFile(
    projectUuid: string,
    version: string,
    sbomId: string,
    format: NoticeFileFormat,
  ): Promise<string | null> {
    try {
      // Map 'plain' format to backend 'text' endpoint
      const endpointFormat = format === NoticeFileFormat.plain ? 'text' : format;
      const response = await axios.get<string>(
        `/v1/projects/${projectUuid}/versions/${encodeURIComponent(
          version,
        )}/sboms/${sbomId}/notice/${endpointFormat}`,
      );
      // For json we ensure we return the raw string (backend may already send parsed object if axios interprets)
      if (format === NoticeFileFormat.json && typeof response.data !== 'string') {
        return JSON.stringify(response.data);
      }
      return response.data;
    } catch (error) {
      console.error('Error downloading notice file:', error);
      return null;
    }
  }

  /**
   * Get children projects of a group
   */
  public async getChildrenProjects(projectUuid: string): Promise<(Project & {versions?: VersionDetails[]})[] | null> {
    try {
      const response = await axios.get<Project[]>(`/v1/groups/${projectUuid}/children`);

      // For each child project, fetch its versions and status
      const childrenWithVersions = await Promise.all(
        response.data.map(async (child) => {
          try {
            // Get version status for this child
            const childStatus = await this.fetchProjectStatus(child.uuid);
            const isDeprecatedChild = (childStatus.status || '').toLowerCase() === 'deprecated';
            if (isDeprecatedChild) {
              return null;
            }

            // Get version names for this child
            const versionNames = await this.fetchProjectVersionNames(child.uuid);

            // If there are versions, fetch their details
            let versions: VersionDetails[] = [];
            if (versionNames.length > 0) {
              versions = await this.fetchAndPrepareVersions(child.uuid, versionNames, childStatus.versionStatus);
            }

            return {
              ...child,
              versions,
            } as Project & {versions?: VersionDetails[]};
          } catch (error) {
            console.error(`Error fetching versions for child project ${child.name}:`, error);
            return {...child} as Project & {versions?: VersionDetails[]};
          }
        }),
      );

      return childrenWithVersions.filter((child): child is Project & {versions?: VersionDetails[]} => child !== null);
    } catch (error) {
      this.logError('getChildrenProjects', error);
      return null;
    }
  }

  /**
   * Get external source code references for a project version
   * @param projectUuid The project UUID
   * @param version The version identifier
   * @returns Array of external source code references or null if there's an error
   */
  public async getExternalSourceCodeReferences(
    projectUuid: string,
    version: string,
  ): Promise<ExternalSourceCodeReference[] | null> {
    try {
      const response = await axios.get<ExternalSourceCodeReference[]>(
        `/v1/projects/${projectUuid}/versions/${version}/ccs`,
      );
      return response.data;
    } catch (error) {
      console.error('Error fetching external source code references:', error);
      return null;
    }
  }

  /**
   * Create external source code reference for a project version
   * @param projectUuid The project UUID
   * @param version The version identifier
   * @param sourceData The external source code reference data
   * @returns Success response or null if there's an error
   */
  public async createExternalSourceCodeReference(
    projectUuid: string,
    version: string,
    sourceData: {url: string; comment?: string},
  ): Promise<ExternalSourceCodeReference | null> {
    try {
      const response = await axios.post<ExternalSourceCodeReference>(
        `/v1/projects/${projectUuid}/versions/${version}/ccs`,
        sourceData,
      );
      return response.data;
    } catch (error) {
      console.error('Error creating external source code reference:', error);
      if (isAxiosError(error) && error.response?.data) {
        const data = error.response.data as Record<string, unknown>;
        throw new Error((data.message as string) || (data.raw as string) || 'Failed to create reference');
      }
      throw error;
    }
  }
}

// Export singleton instance
export const projectService = ProjectService.getInstance();
