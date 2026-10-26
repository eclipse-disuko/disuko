// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import Icons from '@disclosure-portal/constants/icons';
import {ReviewRemarkLevel, ScanRemarkLevel} from '@disclosure-portal/model/Quality';
import {ComponentDiffType} from '@disclosure-portal/model/VersionDetails';
import {BlobPart} from '@disclosure-portal/types/discobasics';

export function getIconColorScanRemarkLevel(level: ScanRemarkLevel) {
  if (level === ScanRemarkLevel.PROBLEM) {
    return 'red';
  } else if (level === ScanRemarkLevel.WARNING) {
    return 'orange';
  } else if (level === ScanRemarkLevel.INFORMATION) {
    return 'gray';
  }
  return '';
}

export function getIconColorReviewRemarkLevel(level: ReviewRemarkLevel) {
  if (level.toLowerCase() === 'yellow') {
    return 'brand'; // more subtle color
  } else if (level === ReviewRemarkLevel.GRAY) {
    return 'grey';
  }
  return level.toLowerCase();
}

export function getIconReviewRemarkLevel(level: ReviewRemarkLevel) {
  switch (level) {
    case ReviewRemarkLevel.GREEN: {
      return 'mdi-comment-check';
    }
    case ReviewRemarkLevel.YELLOW: {
      return 'mdi-comment-alert';
    }
    case ReviewRemarkLevel.RED: {
      return 'mdi-comment-remove';
    }
    case ReviewRemarkLevel.GRAY: {
      return 'mdi-comment-remove';
    }
  }
}

export function getScanRemarkStatusSortIndex(prStatus: ScanRemarkLevel): number {
  switch (prStatus) {
    case ScanRemarkLevel.PROBLEM:
      return 0;
    case ScanRemarkLevel.WARNING:
      return 1;
    case ScanRemarkLevel.INFORMATION:
      return 2;
    default:
      return 10000;
  }
}

export function getWorstScanRemarkLevel(remarks?: {status: string}[] | null): ScanRemarkLevel {
  if (!remarks || remarks.length === 0) {
    return ScanRemarkLevel.NOT_SET;
  }

  let worst = ScanRemarkLevel.NOT_SET;

  for (const remark of remarks) {
    const status = remark.status as ScanRemarkLevel;
    const statusRank = getScanRemarkStatusSortIndex(status);
    const worstRank = getScanRemarkStatusSortIndex(worst);

    if (statusRank < worstRank) {
      worst = status;
    }
  }

  return worst;
}

export function getIconForDiffType(type: ComponentDiffType): string {
  switch (type) {
    case ComponentDiffType.NEW:
      return Icons.ADD;
    case ComponentDiffType.REMOVED:
      return Icons.REMOVED;
    case ComponentDiffType.CHANGED:
      return Icons.CHANGED;
    default:
      return '';
  }
}

export function downloadFile(content: string, filename: string, contentType: string) {
  const link = document.createElement('a');
  link.click();
  link.target = '_blank';
  link.rel = 'noopener noreferrer';
  link.download = filename;
  link.href = URL.createObjectURL(new Blob([content as BlobPart], {type: contentType}));
  link.click();
}
