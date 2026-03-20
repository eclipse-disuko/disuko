// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

// Deletion action details
export interface DeletionAction {
  action_type: string;
  entity_id: string;
  entity_type: string;
  status: string;
  reason: string;
  project_id?: string;
  project_name?: string;
  task_id?: string;
  task_type?: string;
  priority?: string;
  role_name?: string;
  trace_type?: string;
}


export interface DeletionPlan {
  username: string;
  total_actions: number;
  task_actions: DeletionAction[];
  role_actions: DeletionAction[];
  trace_actions: DeletionAction[];
  profile_actions: DeletionAction[];
}

export interface DeletePersonalDataResponse {
  success: boolean;
  message: string;
  entities_effected: {
    user_tasks_count: number;
    user_roles_count: number;
    data_traces_count: number;
  }
}
