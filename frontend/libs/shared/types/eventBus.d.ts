// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import DHTTPError from '@shared/types/DHTTPError';
import ErrorDialogConfig from '@shared/types/ErrorDialogConfig';
import IdleInfo, {INotificationMeta} from '@shared/types/IdleInfo';

export type Events = {
  'show-snackbar': {message: string; timeout?: number; level: string};
  'on-api-error': DHTTPError;
  'on-error': {error: ErrorDialogConfig};
  'on-idle': {idle: IdleInfo};
  'window-resize': unknown;
  'set-notification': {config: INotificationMeta};
  'tab-change': {tabIndex: number};
};
