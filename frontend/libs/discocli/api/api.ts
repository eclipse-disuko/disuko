import config from '@shared/utils/config';
import axios, {AxiosResponse} from 'axios';
import {setupAuthInterceptors} from '@cli/services/authInterceptor';
import {initInterceptors} from './interceptors';

type ApiBundle = {
  NO_IDLE_PARAM: string;
  getData: <T>(promise: Promise<AxiosResponse<T>>) => Promise<T | null>;
  api: ReturnType<typeof axios.create>;
};

let cachedApi: ApiBundle | null = null;

export const getApi = (): ApiBundle => {
  if (cachedApi) {
    return cachedApi;
  }

  const instance = axios.create({
    withCredentials: true,
    baseURL: config.PUBLIC_API_ENDPOINT,
  });
  const NO_IDLE_PARAM = 'noIdle';

  initInterceptors(instance);
  setupAuthInterceptors(instance);

  const getData = async <T>(promise: Promise<AxiosResponse<T>>): Promise<T | null> => {
    const response: AxiosResponse<T> = await promise;
    return response?.data ?? null;
  };

  cachedApi = {
    NO_IDLE_PARAM,
    getData,
    api: instance,
  };

  return cachedApi;
};
