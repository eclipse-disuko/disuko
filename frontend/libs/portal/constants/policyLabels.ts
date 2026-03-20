// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export enum PolicyLabels {
  ENTERPRISE_PLATFORM = 'enterprise platform',
  MOBILE_PLATFORM = 'mobile platform',
  OTHER_PLATFORM = 'other platform',
  VEHICLE_PLATFORM = 'vehicle platform',
  // ---
  FRONTEND_LAYER = 'frontend layer',
  BACKEND_LAYER = 'backend layer',
  COMBINED_LAYER = 'combined layer', // for old projects
  // ---
  ONBOARD = 'onboard',
  OFFBOARD = 'offboard',
  // ---
  ENTITY_USERS = 'entity users',
  GROUP_USERS = 'group users',
  EXTERNAL_USERS = 'external users',
  // ---
  ENTITY_TARGET = 'entity target',
  EXTERNAL_TARGET = 'external target',
}

export enum ProjectLabels {
  DEVELOPMENT_INHOUSE = 'develop inhouse',
  DEVELOPMENT_INTERNAL = 'develop internal',
  DEVELOPMENT_EXTERNAL = 'develop external',
  DUMMY = 'dummy',
}
