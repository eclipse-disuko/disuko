// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import eventBus from '@disclosure-portal/utils/eventbus';

const timeout = 3000;

interface SnackEvent {
  message: string;
  timeout?: number;
  level: 'info' | 'error';
}

export default function useSnackbar() {
  const info = (message: string) => {
    eventBus.emit('show-snackbar', {
      message: message,
      timeout: timeout,
      level: 'info',
    } as SnackEvent);
  };

  const error = (message: string) => {
    eventBus.emit('show-snackbar', {
      message: message,
      timeout: timeout,
      level: 'error',
    } as SnackEvent);
  };

  return {
    info,
    error,
  };
}
