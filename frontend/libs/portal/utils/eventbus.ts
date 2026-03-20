// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import DHTTPError from '@disclosure-portal/model/DHTTPError';
import ErrorDialogConfig from '@disclosure-portal/model/ErrorDialogConfig';
import IdleInfo, {INotificationMeta} from '@disclosure-portal/model/IdleInfo';
import mitt from 'mitt';

type Events = {
  'show-snackbar': {message: string; timeout?: number; level: string};
  'on-api-error': DHTTPError;
  'on-error': {error: ErrorDialogConfig};
  'on-idle': {idle: IdleInfo};
  'window-resize': {};
  'set-notification': {config: INotificationMeta};
  'tab-change': {tabIndex: number};
};

const eventBus = mitt<Events>();

export default eventBus;
