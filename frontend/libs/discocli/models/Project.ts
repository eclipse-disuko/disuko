import type {SBOM} from './Sbom';
export interface Version {
  name: string;
  description?: string;
  status?: string;
  lastSbomUploaded?: string;
}

export interface RemarkComment {
  author: string;
  content: string;
  created: string;
}

export interface Remark {
  author: string;
  content: string;
  severity: string;
  created: string;
  uuid?: string;
  comments?: RemarkComment[];
}

export interface RemarkWithVersion extends Remark {
  version: string;
}

export interface Project {
  uuid: string;
  name: string;
  description: string;
  schema?: string;
  created?: string;
  updated?: string;
  isGroup?: boolean;
  versions?: Version[];
  children?: string[];
  status?: string;
  // Consolidated SBOM deliveries fetched across versions (added by service)
  sboms?: SBOM[];
}
