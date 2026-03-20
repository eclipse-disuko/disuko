// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import useSnackbar from '@shared/composables/useSnackbar';
import {useI18n} from 'vue-i18n';

/**
 * Composable for clipboard operations with automatic snackbar notifications
 * @returns Object with copyToClipboard function
 */
export function useClipboard() {
  const {t} = useI18n();
  const {info} = useSnackbar();

  /**
   * Copy text to clipboard with automatic snackbar notifications
   * @param content The text content to copy to clipboard
   * @returns Promise that resolves when copy is successful or rejects on error
   */
  const copyToClipboard = async (content: string): Promise<void> => {
    try {
      await navigator.clipboard.writeText(content);
      info(t('SNACK_copied_to_clipboard'));
    } catch (error) {
      info(t('SNACK_WENT_WRONG'));
      throw error;
    }
  };

  return {
    copyToClipboard,
  };
}
