// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export const JOB_TYPE_LICENSE = 0;
export const JOB_TYPE_TERMS_OF_USE = 1;
export const JOB_TYPE_DEPARTMENT = 2;
export const JOB_TYPE_DEPARTMENT_LOAD_DB = 3;
export const JOB_TYPE_LICENSE_ANNOUNCEMENT_GENERATION = 4;
export const JOB_TYPE_ANALYZE = 5;
export const JOB_TYPE_NOTIFICATION = 6;
export const JOB_TYPE_DEPROVISIONING = 7;
export const JOB_TYPE_POLICY_RULE_CHANGE_LOGS = 8;
export const JOB_TYPE_FOSS_DD = 10;

export const JOB_STATUS_IDLE = 0;
export const JOB_STATUS_IN_PROGRESS = 1;
export const JOB_STATUS_SUCCESS = 2;
export const JOB_STATUS_FAILURE = 3;

export const JOB_EXECUTION_MANUAL = 0;
export const JOB_EXECUTION_PERIODIC = 1;
export const JOB_EXECUTION_ONETIME = 2;

export class JobLogEntryDto {
  created: string;
  msg: string;
  level: string;
  instance: string;

  constructor() {
    this.created = '';
    this.msg = '';
    this.level = '';
    this.instance = '';
  }
}

export class RefreshRes {
  added: number;
  changed: number;
  differences: number;
  errors: number;
  handled: number;
  total: number;
  unchanged: number;
  addedLics: string[];
  updatedLics: string[];
  diffLics: string[];
  errorLics: string[];

  constructor() {
    this.added = 0;
    this.changed = 0;
    this.differences = 0;
    this.errors = 0;
    this.handled = 0;
    this.total = 0;
    this.unchanged = 0;
    this.addedLics = [];
    this.updatedLics = [];
    this.diffLics = [];
    this.errorLics = [];
  }
}
export class JobDto {
  _key: string;
  updated: string;
  created: string;
  name: string;
  jobType: number;
  execution: number;
  status: number;
  config: string;
  log: JobLogEntryDto[];
  customRes: unknown;
  nextScheduledExecution?: string;

  constructor() {
    this._key = '';
    this.updated = '';
    this.created = '';
    this.name = '';
    this.jobType = -1;
    this.execution = -1;
    this.status = -1;
    this.config = '';
    this.log = [];
  }
}

export class SetConfigDto {
  config: string;

  constructor(config: string) {
    this.config = config;
  }
}

export function jobStatusToString(jobStatus: number): string {
  switch (jobStatus) {
    case JOB_STATUS_IDLE:
      return 'JOB_STATUS_IDLE';
    case JOB_STATUS_IN_PROGRESS:
      return 'JOB_STATUS_ON_PROGRESS';
    case JOB_STATUS_SUCCESS:
      return 'JOB_STATUS_SUCCESS';
    case JOB_STATUS_FAILURE:
      return 'JOB_STATUS_FAILURE';
    default:
      return 'JOB_STATUS_UNKNOWN';
  }
}

export function jobTypeToString(jobType: number): string {
  switch (jobType) {
    case JOB_TYPE_LICENSE:
      return 'JOB_TYPE_LICENSE';
    case JOB_TYPE_TERMS_OF_USE:
      return 'JOB_TYPE_TERMS_OF_USE';
    case JOB_TYPE_DEPARTMENT:
      return 'JOB_TYPE_DEPARTMENT';
    case JOB_TYPE_DEPARTMENT_LOAD_DB:
      return 'JOB_TYPE_DEPARTMENT';
    case JOB_TYPE_LICENSE_ANNOUNCEMENT_GENERATION:
      return 'JOB_TYPE_LICENSE_ANNOUNCEMENT_GENERATION';
    case JOB_TYPE_ANALYZE:
      return 'JOB_TYPE_ANALYZE';
    case JOB_TYPE_NOTIFICATION:
      return 'JOB_TYPE_NOTIFICATION';
    case JOB_TYPE_DEPROVISIONING:
      return 'JOB_TYPE_DEPROVISIONING';
    case JOB_TYPE_POLICY_RULE_CHANGE_LOGS:
      return 'JOB_TYPE_POLICY_RULE_CHANGE_LOGS';
    case JOB_TYPE_FOSS_DD:
      return 'JOB_TYPE_FOSS_DD';
    default:
      return 'JOB_TYPE_UNKNOWN';
  }
}

export function jobExecutionToString(jobExecution: number): string {
  switch (jobExecution) {
    case JOB_EXECUTION_MANUAL:
      return 'JOB_EXECUTION_MANUAL';
    case JOB_EXECUTION_PERIODIC:
      return 'JOB_EXECUTION_PERIODIC';
    case JOB_EXECUTION_ONETIME:
      return 'JOB_EXECUTION_ONETIME';
    default:
      return 'JOB_EXECUTION_UNKNOWN';
  }
}
