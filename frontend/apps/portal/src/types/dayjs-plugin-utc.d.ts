// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {PluginFunc} from 'dayjs';

declare module 'dayjs' {
  interface Dayjs {
    utc(): Dayjs;
    local(): Dayjs;
    isUTC(): boolean;
    utcOffset(offset: number | string, keepLocalTime?: boolean): Dayjs;
  }
}
declare module 'dayjs-plugin-utc' {
  const plugin: PluginFunc;
  export = plugin;
}
