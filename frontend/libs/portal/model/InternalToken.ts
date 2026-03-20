// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export const capabilityMap = {
  StatisticsCSV: 0,
  CustomLicenses: 1,
};
export type Capability = keyof typeof capabilityMap;

export interface InternalToken {
  _key?: string;
  updated?: string;
  created?: string;
  name: string;
  revoked: boolean;
  description: string;
  expiry: string;
  token: string;
  capabilities: Capability[];
}

// Shared interface for create/update operations
export interface InternalTokenRequest {
  _key?: string;
  name: string;
  revoked?: boolean;
  token?: string;
  description: string;
  expiry: string;
  capabilities: Capability[];
}
