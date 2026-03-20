// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {CONTROL_IS_PRESSED, releaseKeys, SHIFT_IS_PRESSED} from '@disclosure-portal/keyState';
import {ICallback} from '@disclosure-portal/model/Callback';
import {Router} from 'vue-router';

// TODO: Remove overhanding the router, use a composable instead
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

export function createVersionURL(project: string, version: string) {
  return `/dashboard/projects/${encodeURIComponent(project)}/versions/${encodeURIComponent(version)}`;
}
export function createProjectURL(project: string) {
  return `/dashboard/projects/${encodeURIComponent(project)}/overview`;
}
export function createSBOMURL(project: string, version: string, sbom: string) {
  return `/dashboard/projects/${encodeURIComponent(project)}/versions/${encodeURIComponent(version)}/component/${encodeURIComponent(sbom)}`;
}

export function openProjectUrlByKey(_key: string, router: Router) {
  if (CONTROL_IS_PRESSED) {
    openUrlInNewTab(createProjectURL(_key));
    return;
  }

  return router.push({path: createProjectURL(_key)});
}

export function openUrlInNewTab(url: string) {
  window.open('#' + url, '_blank');
}
