// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import type {SortItem} from '@shared/types/table';
import {DATE_FORMAT, DATE_FORMAT_SHORT} from '@shared/utils/constant';
import dayjs from 'dayjs';

export function getCssClassForTableRow(item: {Created: string}): string {
  return isNew(item) ? 'animation-new-row' : '';
}

export function getCssClassForReadonlyRow(): string {
  return 'default';
}

export function isNew(item: {Created: string} | {Uploaded: string} | {created: string}) {
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
    if ((code > 64 && code < 91) || (code > 96 && code < 123)) {
      return 'Invalid date';
    }
    const date = dayjs(dateStr);
    if (date.isValid()) {
      return date.format(DATE_FORMAT_SHORT).replace(/([+-])(\d{2}):(\d{2})/, (_, sign, hours, minutes) => {
        const offsetHours = parseInt(hours, 10);
        const offsetMinutes = parseInt(minutes, 10);
        return `${sign}${offsetHours}${offsetMinutes ? ':' + minutes : ''}`;
      });
    } else {
      return 'Invalid date';
    }
  } catch (error) {
    return 'Invalid date ' + error;
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
