import type {AxiosInstance} from 'axios';
import {AxiosError, AxiosRequestConfig, AxiosResponse} from 'axios';
import {authService} from './authService';

let isRefreshing = false;
let refreshPromise: Promise<void> | null = null;
const installedApis = new WeakSet<AxiosInstance>();

async function refreshToken(api: AxiosInstance): Promise<void> {
  if (isRefreshing && refreshPromise) {
    return refreshPromise;
  }
  isRefreshing = true;
  refreshPromise = api
    .get('/auth/refresh', {withCredentials: true})
    .then(() => {
      isRefreshing = false;
      refreshPromise = null;
    })
    .catch((err) => {
      isRefreshing = false;
      refreshPromise = null;
      void authService.logout();
      throw err;
    });
  return refreshPromise;
}

export function setupAuthInterceptors(api: AxiosInstance): void {
  if (installedApis.has(api)) {
    return;
  }
  installedApis.add(api);

  api.defaults.withCredentials = true;

  api.interceptors.response.use(
    (response: AxiosResponse) => response,
    async (error: AxiosError) => {
      if (!error.config) {
        return Promise.reject(error);
      }

      const status = error.response?.status ?? 0;
      const originalRequest = error.config as AxiosRequestConfig & {_retry?: boolean};
      const url = (originalRequest.url || '').toString();

      const isRefreshCall = url.includes('/auth/refresh') || url.includes('/auth/login');
      const shouldRefresh = status === 401 && !originalRequest._retry && !isRefreshCall;

      if (shouldRefresh) {
        originalRequest._retry = true;
        try {
          await refreshToken(api);
          return api(originalRequest);
        } catch (refreshErr) {
          return Promise.reject(refreshErr);
        }
      }

      return Promise.reject(error);
    },
  );
}

export default setupAuthInterceptors;
