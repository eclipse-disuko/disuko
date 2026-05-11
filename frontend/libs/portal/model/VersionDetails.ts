// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {LicenseRuleSlim} from '@disclosure-portal/model/LicenseRule';
import {IUploaded, UnmatchedLicense} from '@disclosure-portal/model/Project';
import {ScanRemarkLevel} from '@disclosure-portal/model/Quality';
import {PolicyDecisionSlim} from './PolicyDecision';

export class MetaInfo {
  public name = '';
  public spdxId = '';
  public spdxVersion = '';
  // public DataLicense = '';
  // public NameWithUploaded = '';
  public comment = '';
  // public Creators: string[] = [];
  // public CreationData = '';
}

export class SourceCode {
  public _key = '';
  public Hash = '';
  public FileName = '';
  public FileSize = 0;
  public Created: Date = new Date();
  public Updated: Date = new Date();
}

export class AuditLog {
  public _key = '';
  public message = '';
  public meta = '';
  public created: Date = new Date();
  public user = '';
}

export class ExternalSource {
  public _key = '';
  public sourceType = '';
  public url = '';
  public comment = '';
  public hash = '';
  public fileSize = 0;
  public created: Date = new Date();
  public origin = '';
  public uploader = '';
}

export class PolicyRuleStatus {
  public key = '';
  public name = '';
  public type = '';
  public licenseMatched = '';
  public used = true;
  public description = '';
  public isDecisionMade = false;
  public canMakeWarnedDecision = false;
  public canMakeDeniedDecision = false;
  public deniedDecisionDeniedReason = '';
  public licenseRecommendationWeight: number | null = null;
}

export interface IComponentInfo {
  spdxId: string;
  name: string;
  version: string;
  licenseEffective: string;
  license: string;
  licenseDeclared: string;
  licenseComments: string;
  copyrightText: string;
  description: string;
  downloadLocation: string;
  type: string;
  modified: boolean;
  questioned: boolean;
  unasserted: boolean;
  policyRuleStatus: PolicyRuleStatus[];
  licenseApplied: string;
  prStatus: string;
  usedPolicyRule: string;
  purl: string;
}

export class ComponentStats {
  public Total = 0;
  public Allowed = 0;
  public Warned = 0;
  public Denied = 0;
  public Questioned = 0;
  public NoAssertion = 0;
}

export class LicenseFamilyStats {
  public Total = 0;
  public Permissive = 0;
  public WeakCopyLeft = 0;
  public StrongCopyLeft = 0;
  public NetworkCopyLeft = 0;
  public Other = 0;
}
export class ReviewRemarkStats {
  public Total = 0;
  public Acceptable = 0;
  public AcceptableAfterChanges = 0;
  public NotAcceptable = 0;
}

export class ScanRemarkStats {
  public Total = 0;
  public Information = 0;
  public Warning = 0;
  public Problem = 0;
}

export class ScanRemarkTypeStats {
  public total = 0;
  public missingCopyrights = 0;
  public missingCopyrightsLevel = ScanRemarkLevel.NOT_SET;
  public malformedCopyrights = 0;
}

export class NotChartFossLicenseStats {
  public total = 0;
}

export class LicenseRemarkStats {
  public Total = 0;
  public Information = 0;
  public Warning = 0;
  public Alarm = 0;
}

export class InApproval {
  public IsInApproval = false;
  public ApprovalGuid = '';
  public Status = '';
}

export class GeneralStats {
  public SBOMDelivered = false;
  public SourceUploaded = false;
  public ReviewRemark = new ReviewRemarkStats();
}

export class SbomStats {
  public PolicyState = new ComponentStats();
  public LicenseFamily = new LicenseFamilyStats();
  public ScanRemark = new ScanRemarkStats();
  public LicenseRemark = new LicenseRemarkStats();
  public ApprovalInfo = new InApproval();
  public scanRemarkType = new ScanRemarkTypeStats();
  public notChartFossLicense = new NotChartFossLicenseStats();
}

export class ComponentInfo implements IComponentInfo {
  public spdxId = '';
  public name = '';
  public version = '';
  public licenseEffective = '';
  public license = '';
  public licenseDeclared = '';
  public licenseComments = '';
  public worstFamily = '';
  public copyrightText = '';
  public description = '';
  public downloadLocation = '';
  public type = '';
  public modified = false;
  public questioned = false;
  public unasserted = false;
  public policyRuleStatus: PolicyRuleStatus[] = [];
  public unmatchedLicenses: UnmatchedLicense[] = [];
  public prStatus = '';
  public usedPolicyRule = '';
  public licenseApplied = '';
  public purl = '';
  public canChooseLicense = false;
  public choiceDeniedReason = '';
  public licenseRuleApplied?: LicenseRuleSlim;
  public policyDecisionsApplied: PolicyDecisionSlim[] = [];
  public policyDecisionDeniedReason = '';
  public licenseRecommended = '';
  public licenseRecommendedMsg = '';
}

export class ComponentInfoSlim {
  public spdxId = '';
  public name = '';
  public version = '';
  public licenseExpression = '';
  public componentInfo: ComponentInfo[] = [];
}

export class ApprovalInfo {
  public isInApproval = false;
  public comment = '';
  // public ApprovalGuid = '';
  public status = '';
}

export type Nullable<T> = T | null;

export class SpdxFile {
  public _key = '';
  public created = '';
  public hash = '';
  // public Content = '';
  // public ContentValid = false;
  // public SchemaValid = false;
  // public Type = 0;
  // public ValidationErrors = '';
  // public SchemaId = '';
  // public SchemaName = '';
  // public Stats: ComponentStats = new ComponentStats();
  // public ComponentInfo: ComponentInfo[] = [];
  public metaInfo: MetaInfo = new MetaInfo();
  public uploaded = '';
  public updated: Date = new Date();
  public origin = '';
  public uploader = '';
  public tag = '';
  public approvalInfo = new ApprovalInfo();
  public isInUse = false;
  public isLocked = false;
  public isToDelete = false;
  public isToRetain = false;
  public isRecent = false;
  public overallReview?: OverallReview;
}

export class SpdxFileSlim {
  public _key = '';
  public projectVersionId = '';
  public uploaded: Date = new Date();
  public updated = '';
  public name = '';
}

export enum OverallReviewState {
  UNREVIEWED = 'UNREVIEWED',
  ACCEPTABLE = 'ACCEPTABLE',
  ACCEPTABLE_AFTER_CHANGES = 'ACCEPTABLE_AFTER_CHANGES',
  AUDITED = 'AUDITED',
  NOT_ACCEPTABLE = 'NOT_ACCEPTABLE',
}

export class OverallReview {
  public created = '';
  public updated = '';
  public state = OverallReviewState.UNREVIEWED;
  public comment = '';
  public sbomId = '';
  public sbomName = '';
  public sbomUploaded = '';
  public creator = '';
  public creatorFullName = '';
}

export class VersionSlimDto {
  public _key = '';
  public parentKey = '';
  public name = '1.0';
  public description = '';
  public created = '';
  public updated = '';
  public status = '';
  public currentSpdxFile: SpdxFileSlim = new SpdxFileSlim();
  public spdxFileHistory: SpdxFileSlim[] = [];
  public isDeleted = false;
  public overallReviews: OverallReview[] = [];
}

export class VersionSlim extends VersionSlimDto {
  constructor(dto: VersionSlimDto) {
    super();
    Object.assign(this, dto);
  }
}

export class ComponentsInfoResponse {
  componentInfo: ComponentInfo[] = [];
  componentStats: ComponentStats = new ComponentStats();
  bulkPolicyDecisionDeniedReason = '';
}

export enum ComponentDiffType {
  UNCHANGED = 'UNCHANGED',
  NEW = 'NEW',
  REMOVED = 'REMOVED',
  CHANGED = 'CHANGED',
}

export class ComponentDiff {
  public SpdxId = false;
  public Name = false;
  public Version = false;
  public LicenseComments = false;
  public LicenseDeclared = false;
  public License = false;
  public LicenseEffective = false;
  public CopyrightText = false;
  public Description = false;
  public DownloadLocation = false;
  public prStatus = false;
  public Type = false;
  public Modified = false;
  public Questioned = false;
  public Unasserted = false;
  public PURL = false;
  public DiffType = ComponentDiffType.UNCHANGED;
  public ComponentOld = new ComponentInfo();
  public ComponentNew = new ComponentInfo();
}

export class ComponentDiffWrapper implements IComponentInfo {
  public diff: ComponentDiff;

  constructor(diff: ComponentDiff) {
    this.diff = diff;
  }

  get name(): string {
    return this.getComponentForView().name;
  }

  get spdxId(): string {
    return this.getComponentForView().spdxId;
  }

  get type(): string {
    return this.getComponentForView().type;
  }

  get version(): string {
    return this.getComponentForView().version;
  }

  get licenseEffective(): string {
    return this.getComponentForView().licenseEffective;
  }

  get license(): string {
    return this.getComponentForView().license;
  }

  get licenseDeclared(): string {
    return this.getComponentForView().licenseDeclared;
  }

  get licenseComments(): string {
    return this.getComponentForView().licenseComments;
  }

  get copyrightText(): string {
    return this.getComponentForView().copyrightText;
  }

  get description(): string {
    return this.getComponentForView().description;
  }

  get downloadLocation(): string {
    return this.getComponentForView().downloadLocation;
  }

  get policyRuleStatus(): PolicyRuleStatus[] {
    return this.getComponentForView().policyRuleStatus;
  }

  get prStatus(): string {
    return this.getComponentForView().prStatus;
  }

  get usedPolicyRule(): string {
    return this.getComponentForView().usedPolicyRule;
  }

  get modified(): boolean {
    return this.getComponentForView().modified;
  }

  get questioned(): boolean {
    return this.getComponentForView().questioned;
  }

  get unasserted(): boolean {
    return this.getComponentForView().unasserted;
  }

  get licenseApplied(): string {
    return this.getComponentForView().licenseApplied;
  }

  get purl(): string {
    return this.getComponentForView().purl;
  }

  public getComponentForView(): ComponentInfo {
    if (this.diff.DiffType === ComponentDiffType.REMOVED) {
      return this.diff.ComponentOld;
    }
    return this.diff.ComponentNew;
  }

  public getOtherComponent(): ComponentInfo {
    if (this.diff.DiffType === ComponentDiffType.REMOVED) {
      return new ComponentInfo();
    }
    if (this.diff.DiffType === ComponentDiffType.NEW) {
      return new ComponentInfo();
    }
    return this.diff.ComponentOld;
  }
}

interface IVersionDetailsDTO {
  Name: string;
  Description: string;
}

export class VersionDetailsDTO implements IVersionDetailsDTO {
  public Name: string;
  public Description: string;
  public Labels: string[];
  public Created: string;
  public Updated: string;
  public CurrentSpdxFile: SpdxFile;

  public constructor() {
    this.Name = 'Neues Project';
    this.Description = '';
    this.Labels = [];
    this.Updated = '';
    this.Created = '';
    this.CurrentSpdxFile = {} as SpdxFile;
  }
}

export class VersionDetails extends VersionDetailsDTO {
  constructor(dto: VersionDetailsDTO) {
    super();
    Object.assign(this, dto);
  }
}

export enum NoticeFileFormat {
  plain = 'text',
  html = 'html',
  json = 'json',
}

export class ComponentChanges {
  public SpdxId = false;
  public Name = false;
  public Version = false;
  public LicenseComments = false;
  public LicenseDeclared = false;
  public License = false;
  public LicenseEffective = false;
  public CopyrightText = false;
  public Description = false;
  public DownloadLocation = false;
  public prStatus = false;
  public Type = false;
  public Modified = false;
  public Questioned = false;
  public Unasserted = false;
  public PURL = false;
}

export class ComponentMultiDiff {
  public DiffType = ComponentDiffType.UNCHANGED;
  public Name = '';
  public Changes: Record<string, ComponentChanges> = {};
  public ComponentsOld = [] as ComponentInfo[];
  public ComponentsNew = [] as ComponentInfo[];
}

export class OverallReviewRequest {
  public state = OverallReviewState.UNREVIEWED;
  public comment = '';
  public sbomId = '';
  public sbomName = '';
  public sbomUploaded = '';
}
