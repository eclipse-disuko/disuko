// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import axios from 'axios';
import myConfig from '../utils/config';

export const useAxios = () => {
  const instance = axios.create({
    withCredentials: true,
    baseURL: myConfig.SERVER_URL,
  });

  const NO_IDLE_PARAM = 'noIdle';

  return {
    NO_IDLE_PARAM,
    instance,
  };
};
