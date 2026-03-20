// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import config from '@shared/utils/config';

export function logout() {
  const url = config.SERVER_URL + config.OAUTH.LOGOUT;
  window.location.replace(url);
}
