// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import IdleInfo from '@shared/types/IdleInfo';
import {BlobPart} from '@disclosure-portal/types/discobasics';
import {AxiosResponse} from 'axios';
import eventBus from '@shared/utils/eventbus';

export function downloadFile(fileName: string, downloadPromise: Promise<AxiosResponse>, withIdle: boolean = false) {
  const link = document.createElement('a');
  link.click();
  link.target = '_blank';
  link.rel = 'noopener noreferrer';
  if (withIdle) {
    eventBus.emit('on-idle', {idle: new IdleInfo(true)});
  }
  downloadPromise
    .then((response) => {
      if (!response) {
        return;
      }
      link.download = fileName;

      let content;
      if (
        typeof response.data === 'object' &&
        response.data !== null &&
        response.data.type !== 'application/octet-stream'
      ) {
        content = JSON.stringify(response.data, null, 4);
      } else {
        content = response.data;
      }
      if (withIdle) {
        eventBus.emit('on-idle', {idle: new IdleInfo(false)});
      }

      link.href = URL.createObjectURL(new Blob([content as unknown as BlobPart]));
      link.click();
    })
    .catch((e) => {
      console.error('cannot download file: ' + fileName, e);
    });
}
