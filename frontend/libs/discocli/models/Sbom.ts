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

export interface SpdxStatusInformation {
  components: SpdxStatusComponent[];
  scanRemarks?: SpdxStatusScanRemarks[];
  licenseRemarks?: SpdxStatusLicenceRemarks[];
  policies?: SpdxStatusPolicy[];
}

export interface SpdxStatusComponent {
  spdxId: string;
  license: string;
  name: string;
  version: string;
  prStatus: string;
  scanRemarks: null | Array<{
    status: string;
    remark: string;
    type: string;
    description: string;
  }>;
  licenseRemarks: Array<{
    status: string;
    remark: string;
    type: string;
    licenseMatched: string;
    description: string;
  }>;
  policyRuleStatus: Array<{
    name: string;
    licenseMatched: string;
    type: string;
    used: boolean;
    description: string;
  }>;
  usedAliases: Array<{
    name: string;
    referencedName: string;
  }>;
  usedDecision?: {
    expression: string;
    licenseID: string;
    name: string;
  };
  packageUrl?: string;
  hash?: string;
  issues?: string[];
}

export interface SpdxStatusScanRemarks {
  description: string;
  type: string;
  severity: string;
  state?: string;
  component?: string;
  created?: string;
}

export interface SpdxStatusLicenceRemarks {
  description: string;
  type: string;
  component?: string;
  created?: string;
}

export interface SpdxStatusPolicy {
  description: string;
  name: string;
  type: string;
  key: string;
  created?: string;
}

export interface SpdxUploadResponse {
  docIsValid: boolean;
  validationFailedMessage?: string; // multi-line details from backend
  hash?: string;
  fileUploaded?: boolean;
  id?: string;
  sbomguid?: string;
  code?: string; // backend error code
  message?: string; // backend message
  reqID?: string; // backend request id
  raw?: string; // raw additional details
}

// Backend lightweight status response for SBOM listings. We extend it with optional
// fields so it can be safely widened into the richer SBOM model where data is available.
export interface SbomStatusPublicResponseDto {
  status: string;
  id?: string;
  name?: string;
  version?: string; // channel/version identifier
  updated?: string;
  valid?: boolean;
  details?: {
    name?: string;
    id?: string;
    version?: string;
    creators?: string;
    created?: string;
    uploaded?: string;
    status?: boolean;
    tag?: string;
    isLocked?: boolean;
    isRetain?: boolean;
  };
}

// Detailed SBOM information fetched per SBOM (subset aligns with SBOM.details)
export interface SbomDetails {
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
}

export interface PolicyRule {
  name: string;
  licenseMatched: string;
  type: string;
  used?: boolean;
  description?: string;
}

export interface LicenseRemark {
  status: string;
  licenseMatched?: string;
  type: string;
  description?: string;
}

export interface ScanRemark {
  status: string;
  type: string;
  description: string;
  severity?: string;
}

export interface AugmentedLicenseRemark extends LicenseRemark {
  _statusIcon: string;
  _statusColor: string;
  _statusText: string;
  _description: string;
}

export interface AugmentedScanRemark extends ScanRemark {
  _statusIcon: string;
  _statusColor: string;
  _statusText: string;
  _description: string;
}

export interface AugmentedPolicyRule extends PolicyRule {
  _iconColor: string;
  _icon: string;
  _description: string;
}

/**
 * Compare function for license status to enable proper sorting
 * Order: ALARM (0) < WARNING (1) < INFORMATION (2) < UNKNOWN (3)
 */
export function compareLicenseStatus(a: string | undefined, b: string | undefined): number {
  const statusWeight: Record<string, number> = {
    alarm: 0,
    warning: 1,
    information: 2,
    unknown: 3,
  };
  const aWeight = statusWeight[a?.toLowerCase() || ''] ?? 999;
  const bWeight = statusWeight[b?.toLowerCase() || ''] ?? 999;
  return aWeight - bWeight;
}

/**
 * Compare function for scan remark status to enable proper sorting
 * Order: PROBLEM (0) < WARNING (1) < INFORMATION (2)
 */
export function compareScanRemarkStatus(a: string | undefined, b: string | undefined): number {
  const statusWeight: Record<string, number> = {
    problem: 0,
    warning: 1,
    information: 2,
  };
  const aWeight = statusWeight[a?.toLowerCase() || ''] ?? 999;
  const bWeight = statusWeight[b?.toLowerCase() || ''] ?? 999;
  return aWeight - bWeight;
}
