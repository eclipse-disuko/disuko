// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import { useI18n } from 'vue-i18n';

/**
 * Composable for managing page titles throughout the application
 * Provides consistent title formatting: "Page Title - Section - Disclosure Portal"
 */
export function usePageTitle() {
  const { t } = useI18n();
  const baseTitle = import.meta.env.DEV ? 'Dev Disuko' : t('APP_NAME');

  /**
   * Sets a reactive title directly with a string
   * @param titleString - The title string to set
   */
  const useReactiveTitle = (titleString: string) => {
    document.title = titleString + ' | ' + baseTitle;
  };

  return {
    useReactiveTitle,
    baseTitle
  };
}
