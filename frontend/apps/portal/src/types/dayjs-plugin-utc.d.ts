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
