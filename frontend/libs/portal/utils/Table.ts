// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {ICreated, ICreatedSmall, IUploaded} from '@disclosure-portal/model/Project';
import {OverallReviewState} from '@disclosure-portal/model/VersionDetails';
import {SortItem} from '@shared/types/table';
import {DATE_FORMAT, DATE_FORMAT_SHORT} from '@shared/utils/constant';
import dayjs from 'dayjs';

export function getCssClassForTableRow(item: ICreated): string {
  return isNew(item) ? 'animation-new-row' : '';
}

export function getCssClassForReadonlyRow(): string {
  return 'default';
}

export function isNew(item: ICreated | IUploaded | ICreatedSmall) {
  if (!item) {
    return false;
  }

  let date;
  if ('Created' in item) {
    date = item.Created;
  } else if ('Uploaded' in item) {
    date = item.Uploaded;
  } else if ('created' in item) {
    date = item.created;
  } else {
    return false;
  }

  return dayjs(date).unix() + 1 > dayjs().unix();
}

export function formatDate(dateStr: string) {
  if (!dateStr || dateStr.length === 0) {
    return 'Invalid date';
  }
  return dayjs(dateStr).format(DATE_FORMAT);
}

export function formatDateAndTime(dateStr: string) {
  try {
    if (!dateStr || dateStr.length === 0) {
      return 'Invalid date';
    }
    const code = dateStr.charCodeAt(0);
    if (
      (code > 64 && code < 91) || // upper alpha (A-Z)
      (code > 96 && code < 123)
    ) {
      // lower alpha (a-z)
      return 'Invalid date';
    }
    const d = dayjs(dateStr);
    if (d.isValid()) {
      return d.format(DATE_FORMAT_SHORT);
    } else {
      return 'Invalid date';
    }
  } catch (e) {
    return 'Invalid date ' + e;
  }
}

export interface SearchOptions {
  page: number;
  itemsPerPage: number;
  search: string;
  sortBy: SortItem[];
  sortDesc?: boolean[];
  groupBy: string[];
  groupDesc?: boolean[];
  multiSort?: boolean;
  mustSort?: boolean;
  filterString: string;
  filterBy: Record<string, string[]>;
}

export function getOverallReviewIcon(state: OverallReviewState) {
  switch (state) {
    case OverallReviewState.UNREVIEWED:
      return 'mdi-comment';
    case OverallReviewState.ACCEPTABLE:
      return 'mdi-comment-check';
    case OverallReviewState.ACCEPTABLE_AFTER_CHANGES:
      return 'mdi-comment-alert';
    case OverallReviewState.NOT_ACCEPTABLE:
      return 'mdi-comment-remove';
    case OverallReviewState.AUDITED:
      return 'mdi-clipboard-check-outline';
  }
}

export function getOverallReviewColor(state: OverallReviewState) {
  switch (state) {
    case OverallReviewState.UNREVIEWED:
      return 'versionUnreviewed';
    case OverallReviewState.ACCEPTABLE:
      return 'versionApproved';
    case OverallReviewState.ACCEPTABLE_AFTER_CHANGES:
      return 'versionAcceptableAfterChanges';
    case OverallReviewState.NOT_ACCEPTABLE:
      return 'red';
    case OverallReviewState.AUDITED:
      return 'green';
  }
}

export function getVersionStateIcon(state: string) {
  switch (state) {
    case 'new':
      return 'mdi-comment';
    case 'unreviewed':
      return 'mdi-comment';
    case 'acceptable':
      return 'mdi-comment-check';
    case 'acceptable_after_changes':
      return 'mdi-comment-alert';
    case 'not_acceptable':
      return 'mdi-comment-remove';
    case 'audited':
      return 'mdi-clipboard-check-outline';
  }
}

export function getIconColor(state: string) {
  switch (state) {
    case 'new':
      return 'versionNew';
    case 'unreviewed':
      return 'versionUnreviewed';
    case 'acceptable':
      return 'versionApproved';
    case 'acceptable_after_changes':
      return 'versionAcceptableAfterChanges';
    case 'not_acceptable':
      return 'red';
    case 'audited':
      return 'green';
  }
}

const outdatedSbomDays = 30;

export function sbomOutdated(sbomUploaded: string) {
  if (sbomUploaded) {
    return dayjs().diff(dayjs(sbomUploaded), 'days') > outdatedSbomDays;
  }
  return false;
}

export function getOverallReviewTranslationKey(status: string): string {
  if (!status) return 'OVERALL_REVIEW_UNREVIEWED';
  const upperStatus = status.toUpperCase();
  return 'OVERALL_REVIEW_' + upperStatus;
}
