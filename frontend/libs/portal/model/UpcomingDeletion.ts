// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export interface BlockingProject {
  key: string;
  name: string;
  projectLabels: string[];
  policyLabels: string[];
  freeLabels: string[];
}

export interface UpcomingDeletion {
  user: string;
  forename: string;
  lastname: string;
  department: string;
  departmentDescription: string;
  deprovisioned: string;
  deletionDate: string;
  overdue: boolean;
  blockingProjects: BlockingProject[];
}
