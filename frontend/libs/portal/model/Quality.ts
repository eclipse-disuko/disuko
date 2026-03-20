// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {PolicyRuleStatus} from '@disclosure-portal/model/VersionDetails';
import {ObligationDTO} from './IObligation';
import {UnmatchedLicense} from './Project';

export enum ScanRemarkLevel {
  NOT_SET = '',
  PROBLEM = 'PROBLEM',
  WARNING = 'WARNING',
  INFORMATION = 'INFORMATION',
}

export class ScanRemark {
  public key = '';
  public name = '';
  public version = '';
  public type = '';
  public status: ScanRemarkLevel = ScanRemarkLevel.NOT_SET;
  public remarkKey = '';
  public descriptionKey = '';
  public spdxId = '';
  public policyRuleStatus: PolicyRuleStatus[] = [];
  public unmatchedLicenses: UnmatchedLicense[] = [];
}

export class LicenseRemark {
  public key = '';
  public name = '';
  public version = '';
  public status = '';
  public license = '';
  public remark = '';
  public type = '';
  public description = '';
  public spdxId = '';
  public policyRuleStatus: PolicyRuleStatus[] = [];
}

export class AffectedComponent {
  public spdxid = '';
  public name = '';
  public version = '';
  public policyRuleStatus: PolicyRuleStatus[] = [];
}

export class LicenseRemarks {
  public license = '';
  public warnings = false;
  public alarms = false;
  public obligations: ObligationDTO[] = [];
  public affected: AffectedComponent[] = [];
}

export enum EventType {
  NOT_SET = '',
  CLOSED = 'CLOSED',
  CANCELLED = 'CANCELLED',
  IN_PROGRESS = 'IN_PROGRESS',
  REOPENED = 'REOPENED',
  COMMENT = 'COMMENT',
  CHANGED_LEVEL = 'CHANGED_LEVEL',
  CHANGED_TITLE = 'CHANGED_TITLE',
  CHANGED_DESCRIPTION = 'CHANGED_DESCRIPTION',
}
export interface LevelChange {
  before: ReviewRemarkLevel;
  after: ReviewRemarkLevel;
}

export interface TitleChange {
  before: string;
  after: string;
}

export interface DescriptionChange {
  before: string;
  after: string;
}

export type Comment = string;

export class Event {
  public key = '';
  public created = '';
  public updated = '';

  public type: EventType = EventType.NOT_SET;
  public author = '';
  public authorFullName = '';
  public content?: LevelChange | TitleChange | DescriptionChange | Comment;
}

export class ComponentMeta {
  public componentId = '';
  public componentName = '';
  public componentVersion = '';
}

export class LicenseMeta {
  public licenseId = '';
  public licenseName = '';
}

export enum ReviewRemarkLevel {
  NOT_SET = '',
  GREEN = 'GREEN',
  YELLOW = 'YELLOW',
  RED = 'RED',
}

export enum ReviewRemarkStatus {
  NOT_SET = '',
  OPEN = 'OPEN',
  IN_PROGRESS = 'IN_PROGRESS',
  CLOSED = 'CLOSED',
  CANCELLED = 'CANCELLED',
}

export class ReviewRemark {
  public key = '';
  public created = '';
  public updated = '';

  public author = '';
  public title = '';
  public closed = '';
  public level: ReviewRemarkLevel = ReviewRemarkLevel.NOT_SET;
  public description = '';
  public status: ReviewRemarkStatus = ReviewRemarkStatus.NOT_SET;
  public events: Event[] = [];
  public sbomId = '';
  public sbomName = '';
  public sbomUploaded: Date | null = null;

  public components: ComponentMeta[] = [];
  public licenses: LicenseMeta[] = [];
  public closing = '';
}

export class ReviewRemarkRequest {
  public title = '';
  public level: ReviewRemarkLevel = ReviewRemarkLevel.NOT_SET;
  public description = '';
  public sbomId = '';
  public components: string[] = [];
  public licenses: string[] = [];

  public static toRequest(item: ReviewRemark): ReviewRemarkRequest {
    const req = new ReviewRemarkRequest();
    req.description = item.description;
    req.level = item.level;
    req.title = item.title;
    req.sbomId = item.sbomId;
    return req;
  }
}

export class SetReviewRemarkStatusRequest {
  public status: ReviewRemarkStatus = ReviewRemarkStatus.NOT_SET;
}

export class CommentReviewRemarkRequest {
  public content = '';
}

export const levelWeight: Map<string, number> = new Map<string, number>([
  ['INFORMATION', 0],
  ['WARNING', 1],
  ['ALARM', 2],
]);

export function compareLevel(a: string, b: string): number {
  return (levelWeight.get(a) ?? 0) - (levelWeight.get(b) ?? 0);
}

export function compareRRLevel(a: ReviewRemarkLevel, b: ReviewRemarkLevel): number {
  const levelWeight: Map<ReviewRemarkLevel, number> = new Map<ReviewRemarkLevel, number>([
    [ReviewRemarkLevel.GREEN, 0],
    [ReviewRemarkLevel.YELLOW, 1],
    [ReviewRemarkLevel.RED, 2],
  ]);
  return levelWeight.get(a)! - levelWeight.get(b)!;
}
