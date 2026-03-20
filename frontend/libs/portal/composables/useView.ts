// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useI18n} from 'vue-i18n';

export const useView = () => {
  const {t} = useI18n();

  function getTextOfLevel(level: string) {
    level = level.toUpperCase();
    return t('CLASSIFICATION_LEVEL_' + level);
  }
  function getTextOfType(type: string): string {
    return t('CLASSIFICATION_TYPE_' + type.toUpperCase());
  }

  return {
    getTextOfLevel,
    getTextOfType,
  };
};
