// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {getLicenseApprovalTypeKeys, ITextValue} from '@disclosure-portal/model/License';
import {useI18n} from 'vue-i18n';

/**
 * Composable hook for license-related utilities
 * @returns Object containing all license utility functions
 */
export function useLicense() {
  const {t} = useI18n();

  /**
   * Generates an internationalized text for a given prefix and key
   * @param prefix - The i18n prefix to use
   * @param key - The key to translate
   * @returns The translated text
   */
  function getI18NTextOfPrefixKey(prefix: string, key: string): string {
    if (!key || key.length === 0) {
      key = 'unknown';
    }
    key = key.replaceAll(' ', '_').toUpperCase();
    return t(prefix + key);
  }

  /**
   * Creates an array of ITextValue objects from an array of keys with i18n translations
   * @param prefix - The i18n prefix to use for translations
   * @param keys - Array of keys to convert
   * @returns Array of ITextValue objects
   */
  function createITextValueArrayOfKeys(prefix: string, keys: string[]): ITextValue[] {
    const result = [] as ITextValue[];
    keys.forEach((key) => {
      result.push({value: key, text: getI18NTextOfPrefixKey(prefix, key)});
    });
    return result;
  }

  /**
   * Get license types with their i18n translations
   * @returns Array of license types as ITextValue
   */
  function getLicenseTypes(): ITextValue[] {
    const keys = ['', 'open source', 'non foss', 'freeware', 'proprietary', 'public domain'] as string[];
    return createITextValueArrayOfKeys('LT_', keys);
  }

  /**
   * Get license review states with their i18n translations
   * @returns Array of review states as ITextValue
   */
  function getLicenseReviewStates(): ITextValue[] {
    const keys = ['', 'requested', 'in review', 'reviewed'] as string[];
    return createITextValueArrayOfKeys('LT_RS_', keys);
  }

  /**
   * Get license approval types with their i18n translations
   * @returns Array of approval types as ITextValue
   */
  function getLicenseApprovalTypes(): ITextValue[] {
    return createITextValueArrayOfKeys('LT_APP_', getLicenseApprovalTypeKeys());
  }

  /**
   * Get license families with their i18n translations
   * @returns Array of license families as ITextValue
   */
  function getLicenseFamily(): ITextValue[] {
    const keys = ['', 'permissive', 'weak copyleft', 'strong copyleft', 'network copyleft'] as string[];
    return createITextValueArrayOfKeys('LIC_FAMILY_', keys);
  }

  return {
    getI18NTextOfPrefixKey,
    createITextValueArrayOfKeys,
    getLicenseTypes,
    getLicenseReviewStates,
    getLicenseApprovalTypes,
    getLicenseFamily,
  };
}
