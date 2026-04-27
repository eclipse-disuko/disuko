// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import Icons from '@disclosure-portal/constants/icons';
import {CONTROL_IS_PRESSED, releaseKeys, SHIFT_IS_PRESSED} from '@disclosure-portal/keyState';
import {IObligation} from '@disclosure-portal/model/IObligation';
import {PolicyState} from '@disclosure-portal/model/PolicyRule';
import {ReviewRemarkLevel, ScanRemarkLevel} from '@disclosure-portal/model/Quality';
import {ComponentDiffType} from '@disclosure-portal/model/VersionDetails';
import {useAppStore} from '@disclosure-portal/stores/app';
import {BlobPart} from '@disclosure-portal/types/discobasics';
import {DATE_FORMAT, DATE_FORMAT_SHORT, DATETIME_FORMAT, DATETIME_FORMAT_FILE} from '@shared/utils/constant';
import dayjs from 'dayjs';
import {Router} from 'vue-router';

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
  }
}

export function getIconForPolicyType(prStatus: string) {
  if (prStatus === 'allow') {
    return Icons.ALLOW;
  } else if (prStatus === 'warn') {
    return Icons.WARNING;
  } else if (prStatus === 'deny') {
    return Icons.DENY;
  } else if (prStatus === 'questioned') {
    return Icons.QUESTIONED;
  } else if (prStatus === 'noassertion') {
    return Icons.NO_ASSERTION;
  }
  return Icons.COMPONENTS;
}

export function getIconForChange(change: 'Added' | 'Removed') {
  if (change === 'Added') {
    return Icons.ADD;
  } else if (change === 'Removed') {
    return Icons.REMOVED;
  }
}

export function getIconColorForPolicyType(prStatus: string) {
  if (prStatus === 'allow') {
    return 'green';
  } else if (prStatus === 'warn') {
    return 'policyStatusWarnedColor';
  } else if (prStatus === 'deny') {
    return 'policyStatusDeniedColor';
  } else if (prStatus === 'questioned') {
    return 'green';
  } else if (prStatus === 'noassertion') {
    return 'policyStatusUnassertedColor';
  }
  return 'var(--v-textColor-base)';
}

export function getIconColorForPolicyTypeHighlighted(prStatus: string) {
  if (prStatus === 'allow') {
    return 'green';
  } else if (prStatus === 'warn') {
    return 'warning';
  } else if (prStatus === 'deny') {
    return 'policyStatusDeniedColor';
  } else if (prStatus === 'questioned') {
    return 'var(--v-textColor-base)';
  } else if (prStatus === 'noassertion') {
    return 'policyStatusUnassertedColor';
  }
  return 'var(--v-textColor-base)';
}

export function getPrStatusSortIndex(prStatus: string): number {
  switch (prStatus) {
    case 'deny':
      return 0;
    case 'noassertion':
      return 1;
    case 'warn':
      return 2;
    case 'questioned':
      return 3;
    case '':
      return 4;
    case 'allow':
      return 5;
    default:
      return 10000;
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

export type ICallback = () => void;

export function openUrl(url: string, router: Router, callbackOnSameSite: ICallback | null = null) {
  if (CONTROL_IS_PRESSED) {
    window.open('#' + url, '_blank');
    return;
  }
  if (SHIFT_IS_PRESSED) {
    releaseKeys();
    window.open('#' + url, '_blank', 'height=500,width=1024');
    return;
  }
  router.push(url);
  if (callbackOnSameSite) {
    callbackOnSameSite();
  }
}

export function openUrlInNewTab(url: string) {
  window.open('#' + url, '_blank');
}

export function formatDate(dateStr: string | Date) {
  if (!dateStr) {
    return '';
  }
  return dayjs(dateStr.toString()).format(DATE_FORMAT);
}

export function formatDateTime(dateStr: string | Date) {
  if (!dateStr) {
    return '';
  }
  return dayjs(dateStr.toString()).format(DATETIME_FORMAT);
}

export function formatDateTimeShort(dateStr: string | Date, asUTC = false) {
  if (!dateStr) {
    return '';
  }
  const d = dayjs(dateStr.toString());
  return (asUTC ? d.utc() : d).format(DATE_FORMAT_SHORT);
}

export function formatDateTimeForFile(dateStr: string | Date, asUTC = false) {
  if (!dateStr) {
    return '';
  }
  const d = dayjs(dateStr.toString());
  return (asUTC ? d.utc() : d).format(DATETIME_FORMAT_FILE);
}

export interface IMap<VALUE> {
  [key: string]: VALUE;
}

export function getStrWithMaxLength(length: number, str: string): string {
  if (!str) {
    return str;
  }
  if (str.length < length) {
    return str;
  }
  return str.substring(0, length - 3) + '...';
}

export function policyStateToTranslationKey(policyState: string): string {
  switch (policyState) {
    case PolicyState.ALLOW:
    case 'allowed':
      return 'ALLOWED';
    case PolicyState.WARN:
    case 'warned':
      return 'WARNED';
    case PolicyState.QUESTIONED:
      return 'QUESTIONED';
    case PolicyState.NOASSERTION:
      return 'UNASSERTED';
    case PolicyState.DENY:
    case 'denied':
      return 'DENIED';
    default:
      return 'NOT_SET';
  }
}

export function getOrderForPolicyState(policyState: string): number {
  switch (policyState) {
    case PolicyState.DENY:
      return 0;
    case PolicyState.NOASSERTION:
      return 1;
    case PolicyState.QUESTIONED:
      return 2;
    case PolicyState.WARN:
      return 3;
    case PolicyState.ALLOW:
      return 4;
    default:
      return -1;
  }
}

export function sortPolicyStatesByOrder(policyState1: string, policyState2: string): number {
  const o1 = getOrderForPolicyState(policyState1);
  const o2 = getOrderForPolicyState(policyState2);
  if (o1 > o2) {
    return 1;
  }
  if (o1 < o2) {
    return -1;
  }
  return 0;
}

export function getIconColorOfLevel(level: string) {
  level = level.toUpperCase();
  if (level === 'WARNING') {
    return 'warning';
  } else if (level === 'INFORMATION') {
    return 'textColor';
  } else if (level === 'ALARM') {
    return 'red';
  }

  return '';
}

export function getIconOfLevel(level: string) {
  level = level.toUpperCase();
  if (level === 'WARNING') {
    return Icons.OBLIGATION_WARNING;
  } else if (level === 'INFORMATION') {
    return Icons.OBLIGATION_INFORMATION;
  } else if (level === 'ALARM') {
    return Icons.OBLIGATION_ALARM;
  }
  return '';
}

export function getClassificationLevels(): string[] {
  return ['warning'.toUpperCase(), 'information'.toUpperCase(), 'alarm'.toUpperCase()];
}

export function getClassificationTypes(): string[] {
  return ['obligation', 'right', 'exception', 'prohibition', 'limitation', 'liability', 'other'];
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

export function originShort(origin: string) {
  if (origin) {
    return origin.split(' ')[0];
  }
}

export function originTooltip(origin: string) {
  if (origin) {
    const originParts = origin.split(' ');
    if (originParts.length > 1) {
      originParts.splice(0, 1);
      return originParts.join(' ');
    }
  }
}

export default function useViewTools() {
  const appStore = useAppStore();
  const getNameForLanguage = (obligation: IObligation): string => {
    if (!obligation) {
      return '';
    }
    if (appStore.getAppLanguage === 'de') {
      return obligation.nameDe;
    }
    return obligation.name;
  };

  const getDescriptionForLanguage = (obligation: IObligation, short = false): string => {
    let description;
    if (appStore.getAppLanguage === 'de') {
      description = obligation.descriptionDe;
    } else {
      description = obligation.description;
    }
    if (short) {
      description = description.substring(0, 120);
    }
    return description;
  };

  const gridPolicyRulesAssignmentsHeaderClassByLanguage = (): string => {
    return appStore.getAppLanguage === 'en' ? 'padding-config-rules-header' : 'padding-config-rules-header-de';
  };

  const gridPolicyRulesAssignmentsRowClassByLanguage = (): string => {
    return appStore.getAppLanguage === 'en' ? 'padding-config-rules-row' : 'padding-config-rules-row-de';
  };

  return {
    gridPolicyRulesAssignmentsHeaderClassByLanguage,
    gridPolicyRulesAssignmentsRowClassByLanguage,
    getDescriptionForLanguage,
    getNameForLanguage,
  };
}
