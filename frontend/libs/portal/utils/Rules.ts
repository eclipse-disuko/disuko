// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useI18n} from 'vue-i18n';

export default function useRules() {
  const {t} = useI18n();

  const required = (fieldName: string): Array<(v: string) => boolean | string> => {
    return [(v: string) => !!v || t('IS_REQUIRED', {fieldName: fieldName})];
  };

  const requiredOrEmpty = (fieldName: string) => {
    return [(v: string) => !!v || v === '' || t('VALIDATION_required', {fieldName: fieldName})];
  };

  const minMaxArray = (fieldName: string, min: number, max: number) => {
    return [
      (v: string[]) => {
        if (!v) return true;
        const invalidElements: string[] = [];
        v.forEach((value) => {
          if (value.length < min || value.length > max) {
            invalidElements.push(value);
          }
        });
        if (invalidElements.length > 0) {
          return t('RULE_MIN_MAX_ARRAY', {fieldName, fieldValue: invalidElements.join(', '), min, max});
        }
        return true;
      },
    ];
  };

  const min = (fieldName: string, min: number, allowEmpty = true) => {
    if (allowEmpty) {
      return [(v: string) => !v || (v && v.length >= min) || t('RULE_MIN_OR_EMPTY', {fieldName, min})];
    } else {
      return [
        (v: string) => !!v || t('IS_REQUIRED', {fieldName}),
        (v: string) => (v && v.length >= min) || t('RULE_MIN', {fieldName, min}),
      ];
    }
  };

  const minMax = (fieldName: string, min: number, max: number, allowEmpty = true) => {
    if (allowEmpty) {
      return [
        (v: string) =>
          !v || (v && v.length >= min && v.length <= max) || t('RULE_MIN_MAX_OR_EMPTY', {fieldName, min, max}),
      ];
    } else {
      return [
        (v: string) => !!v || t('IS_REQUIRED', {fieldName}),
        (v: string) => (v && v.length >= min && v.length <= max) || t('RULE_MIN_MAX', {fieldName, min, max}),
      ];
    }
  };

  const equalsLength = (fieldName: string, equalsTo: number) => {
    return [
      (v: string) => !!v || t('IS_REQUIRED', {fieldName}),
      (v: string) => (v && v.length === equalsTo) || t('RULE_EQUALS', {fieldName, equal: equalsTo}),
    ];
  };

  const longText = (fieldName: string, max = 1000) => {
    return [(v: string) => !v || (v && v.length <= max) || t('RULE_LONGTEXT', {fieldName, max})];
  };

  return {
    required,
    requiredOrEmpty,
    minMaxArray,
    minMax,
    min,
    equalsLength,
    longText,
  };
}
