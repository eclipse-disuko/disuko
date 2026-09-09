// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export interface Version {
  name: string;
  description?: string;
  status?: string;
  lastSbomUploaded?: string;
}

export interface SBOM {
  name: string;
  updated: string;
  valid: boolean;
  id: string;
  version: string;
  details: {
    name: string;
    id: string;
    version: string;
    creators: string;
    created: string;
    uploaded: string;
    status: boolean;
    tag?: string;
    isLocked?: boolean;
    isRetain?: boolean;
  };
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
  sboms?: SBOM[];
}
