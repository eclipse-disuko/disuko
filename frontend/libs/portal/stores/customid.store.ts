// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {CustomIds} from '@disclosure-portal/model/CustomId';
import {customId} from '@disclosure-portal/services/customid.service';
import {defineStore} from 'pinia';
import {ref} from 'vue';

export const useCustomIdStore = defineStore('customId', () => {
  const {getCustomIds} = customId();

  const customIds = ref<CustomIds>({} as CustomIds);

  const updateCustomIds = async () => {
    try {
      const customIdsRaw = await getCustomIds();
      if (customIdsRaw) {
        customIds.value = new CustomIds(customIdsRaw);
      }
    } catch (error) {
      console.error(error);
    }
  };

  return {
    customIds,
    updateCustomIds,
  };
});
