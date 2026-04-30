// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export function createVersionURL(project: string, version: string) {
  return `/dashboard/projects/${encodeURIComponent(project)}/versions/${encodeURIComponent(version)}`;
}
export function createProjectURL(project: string) {
  return `/dashboard/projects/${encodeURIComponent(project)}/overview`;
}
export function createSBOMURL(project: string, version: string, sbom: string) {
  return `/dashboard/projects/${encodeURIComponent(project)}/versions/${encodeURIComponent(version)}/component/${encodeURIComponent(sbom)}`;
}
