// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import Icons from '@shared/constants/icons';
import type {IObligation} from '@shared/model/IObligation';
import {PolicyState} from '@shared/model/PolicyRule';
import {DATE_FORMAT, DATE_FORMAT_SHORT, DATETIME_FORMAT, DATETIME_FORMAT_FILE} from '@shared/utils/constant';
import dayjs from 'dayjs';
import utc from 'dayjs/plugin/utc';
import {useLanguageStore} from '@shared/stores/language.store';
import {storeToRefs} from 'pinia';

dayjs.extend(utc);

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
  const date = dayjs(dateStr.toString());
  return (asUTC ? date.utc() : date).format(DATE_FORMAT_SHORT);
}

export function formatDateTimeForFile(dateStr: string | Date, asUTC = false) {
  if (!dateStr) {
    return '';
  }
  const date = dayjs(dateStr.toString());
  return (asUTC ? date.utc() : date).format(DATETIME_FORMAT_FILE);
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
  const firstOrder = getOrderForPolicyState(policyState1);
  const secondOrder = getOrderForPolicyState(policyState2);
  if (firstOrder > secondOrder) {
    return 1;
  }
  if (firstOrder < secondOrder) {
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
  const languageStore = useLanguageStore();
  const {appLanguage} = storeToRefs(languageStore);
  const getNameForLanguage = (obligation: IObligation): string => {
    if (!obligation) {
      return '';
    }
    if (appLanguage.value === 'de') {
      return obligation.nameDe;
    }
    return obligation.name;
  };

  const getDescriptionForLanguage = (obligation: IObligation, short = false): string => {
    let description;
    if (appLanguage.value === 'de') {
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
    return appLanguage.value === 'en' ? 'padding-config-rules-header' : 'padding-config-rules-header-de';
  };

  const gridPolicyRulesAssignmentsRowClassByLanguage = (): string => {
    return appLanguage.value === 'en' ? 'padding-config-rules-row' : 'padding-config-rules-row-de';
  };

  return {
    gridPolicyRulesAssignmentsHeaderClassByLanguage,
    gridPolicyRulesAssignmentsRowClassByLanguage,
    getDescriptionForLanguage,
    getNameForLanguage,
  };
}
