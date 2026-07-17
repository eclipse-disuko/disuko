import type {SbomStatusPublicResponseDto} from './Sbom';

export interface VersionRequest {
  name: string;
  description?: string;
}

export interface VersionDetails {
  name: string;
  description: string;
  status: 'new' | 'unreviewed';
  lastSbomUploaded: string;
  sbomList?: SbomStatusPublicResponseDto[];
}

export interface VersionStatusPublicResponse {
  name: string;
  status: 'new' | 'unreviewed';
  lastSbomUploaded: string;
}

export enum NoticeFileFormat {
  html = 'html',
  json = 'json',
  plain = 'plain',
}
