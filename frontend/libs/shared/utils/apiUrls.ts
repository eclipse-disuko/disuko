// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

const DASHBOARD_PATH = '/dashboard';

function createProjectBaseURL(project: string): string {
  return `${DASHBOARD_PATH}/projects/${encodeURIComponent(project)}`;
}

export function createProjectURL(project: string): string {
  return `${createProjectBaseURL(project)}/overview`;
}

export function createGroupURL(group: string): string {
  return `${DASHBOARD_PATH}/groups/${encodeURIComponent(group)}/overview`;
}

export function createVersionURL(project: string, version: string): string {
  return `${createProjectBaseURL(project)}/versions/${encodeURIComponent(version)}`;
}

export function createSBOMURL(project: string, version: string, sbom: string): string {
  return `${createVersionURL(project, version)}/component/${encodeURIComponent(sbom)}`;
}
