// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {getApi} from '@disclosure-portal/api';
import {CustomId} from '@disclosure-portal/model/CustomId';

export const customId = () => {
  const {api, getData} = getApi();

  const getCustomIds = () => getData<CustomId[]>(api.get(`/api/v1/customids`));

  return {getCustomIds};
};
