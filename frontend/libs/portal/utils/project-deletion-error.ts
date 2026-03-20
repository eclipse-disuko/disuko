// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import { ProjectSlim } from '@disclosure-portal/model/ProjectsResponse';
import { Project } from '@disclosure-portal/model/Project';



export function canDeleteProject(project: ProjectSlim | Project): boolean {
  return !project.hasChildren && !project.hasApproval && !project.hasSBOMToRetain && project.accessRights.allowProject.delete
}

export function getDeleteTooltip(project: ProjectSlim | Project): string {

  if (!project.accessRights.allowProject.delete) {
    return "TT_not_project_user_right";
  }
  if (project.hasChildren) {
    return "TT_not_delete_project_has_children";
  }
  if (project.hasApproval) {
    return "TT_not_delete_project_has_Approval";
  }
  if (project.hasSBOMToRetain) {
    return "TT_not_delete_project_has_SBOM_retained";
  }
  return "TT_delete_project";
}
