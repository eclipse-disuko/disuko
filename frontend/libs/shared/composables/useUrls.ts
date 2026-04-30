// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {Router} from 'vue-router';
import {ICallback} from '@shared/types/Callback';
import {useEventKeysStore} from '@shared/stores/eventKeys.store';
import {storeToRefs} from 'pinia';
import {createProjectURL} from '@disclosure-portal/utils/url';

const openUrlInNewTab = (url: string, features?: string) => {
  window.open('#' + url, '_blank', features);
};

export const useUrls = () => {
  const eventKeysStore = useEventKeysStore();
  const {controlIsPressed, shiftIsPressed} = storeToRefs(eventKeysStore);

  const openUrl = (url: string, router: Router, callbackOnSameSite: ICallback | null = null) => {
    if (controlIsPressed.value) {
      openUrlInNewTab(url);
      return;
    }

    if (shiftIsPressed.value) {
      openUrlInNewTab(url, 'height=500,width=1024');
      return;
    }

    router.push(url);

    if (callbackOnSameSite) {
      callbackOnSameSite();
    }
  };

  const openProjectUrlByKey = (_key: string, router: Router) => {
    if (controlIsPressed.value) {
      openUrlInNewTab(createProjectURL(_key));
      return;
    }

    return router.push({path: createProjectURL(_key)});
  };

  return {
    openUrl,
    openUrlInNewTab,
    openProjectUrlByKey,
  };
};
