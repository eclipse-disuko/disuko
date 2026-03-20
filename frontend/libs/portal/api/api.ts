// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useAxios} from '@shared/api/useAxios';
import {AxiosResponse} from 'axios';
import {initInterceptors} from './interceptors';

export const getApi = () => {
  const {instance, NO_IDLE_PARAM} = useAxios();

  initInterceptors(instance, NO_IDLE_PARAM);

  const getData = async <T>(promise: Promise<AxiosResponse<T>>): Promise<T | null> => {
    const response: AxiosResponse<T> = await promise;
    return response?.data ?? null;
  };

  return {
    NO_IDLE_PARAM,
    getData,
    api: instance,
  };
};
