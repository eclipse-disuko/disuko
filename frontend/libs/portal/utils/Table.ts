// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {OverallReviewState} from '@disclosure-portal/model/VersionDetails';

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
