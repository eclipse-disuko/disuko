// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import axios from 'axios';
/**
 * Run-time configuration for the Disco Portal.
 */
interface RuntimeConfig {
  SERVER_URL: string;
  PUBLIC_API_ENDPOINT: string;
  PUBLIC_UI_URL?: string;
  OAUTH: {
    LOGIN: string;
    LOGOUT: string;
  };
  isProd: boolean;
  useFutureFoss?: boolean;
  /**
   * Future FOSS config for Vehicle projects
   */
  useFutureProduct?: boolean;
  /**
   * Future FOSS config for Enterprise, Mobile and Other projects
   */
  useFutureIt?: boolean;
  enforceFOSSOfficeConfirmation?: boolean;
  useWinterTheme?: boolean;
}
// This configuration is used for local development, usually to run directly in dev machine
let config = {
  SERVER_URL: import.meta.env.VITE_SERVER_URL,
  PUBLIC_API_ENDPOINT: import.meta.env.VITE_PUBLIC_API_ENDPOINT,
  PUBLIC_UI_URL: import.meta.env.VITE_PUBLIC_UI_URL,
  OAUTH: {
    LOGIN: import.meta.env.VITE_OAUTH_LOGIN,
    LOGOUT: import.meta.env.VITE_OAUTH_LOGOUT,
  },
  isProd: false,
  useFutureFoss: import.meta.env.VITE_USE_FUTURE_FOSS === 'true',
  useFutureProduct: import.meta.env.VITE_USE_FUTURE_PRODUCT === 'true',
  useFutureIt: import.meta.env.VITE_USE_FUTURE_IT === 'true',
  enforceFOSSOfficeConfirmation: false,
  useWinterTheme: false,
} satisfies RuntimeConfig;

/**
 * Determine the config.json path based on the deployment base path.
 * Uses Vite's BASE_URL which is configured per app in vite.config.
 */
function getConfigPath(): string {
  const base = import.meta.env.BASE_URL;
  const normalizedBase = base.replace(/\/$/, '');
  return `${normalizedBase}/config.json`;
}

export const setup = async () => {
  function getVersion() {
    const style = ['color: white', 'background: green', 'font-size:14px', 'padding: 2px'].join(';');
    console.info(`%c Build Date: ${import.meta.env.VITE_BUILD_DATE}`, style);
    console.info(`%c Version: ${import.meta.env.VITE_VERSION}`, style);
    console.info(`%c Commit: ${import.meta.env.VITE_COMMIT}`, style);
    console.info(`%c NODE_ENV=production: ${import.meta.env.PROD}`, style);
    console.info(`%c mode: ${import.meta.env.MODE}`, style);
  }
  // @ts-expect-error this adds getVersion() to the global window object so devs can see the version in the console
  window.getVersion = getVersion;

  // In deployed environments, the configuration is replaced by helm charts config map
  if (import.meta.env.PROD) {
    const configPath = getConfigPath();
    const runtimeConfig = await axios.get(configPath);
    config = {
      SERVER_URL: runtimeConfig.data.VITE_SERVER_URL,
      PUBLIC_API_ENDPOINT: runtimeConfig.data.VITE_PUBLIC_API_ENDPOINT,
      PUBLIC_UI_URL: runtimeConfig.data.VITE_PUBLIC_UI_URL,
      OAUTH: {
        LOGIN: runtimeConfig.data.VITE_OAUTH_LOGIN,
        LOGOUT: runtimeConfig.data.VITE_OAUTH_LOGOUT,
      },
      isProd: runtimeConfig.data.IS_PROD,
      useFutureFoss: runtimeConfig.data.VITE_USE_FUTURE_FOSS,
      useFutureProduct: runtimeConfig.data.VITE_USE_FUTURE_PRODUCT,
      useFutureIt: runtimeConfig.data.VITE_USE_FUTURE_IT,
      enforceFOSSOfficeConfirmation: runtimeConfig.data.ENFORCE_CONFIRMATION,
      useWinterTheme: false,
    };
  }
};

// TODO: call setup after app init if services are refactored
await setup();
export default config;
