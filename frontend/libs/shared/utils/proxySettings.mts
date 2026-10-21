// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export const apiProxySettings = {
  development: {
    '/api': {
      target: 'https://localhost:3333',
      secure: false,
      changeOrigin: true,
    },
  },
};
