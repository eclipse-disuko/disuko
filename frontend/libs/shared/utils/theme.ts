// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export function getDefaultTheme(): string {
  const theme = localStorage.getItem('disco-theme');
  if (theme === 'dark' || theme === 'light') {
    return theme;
  }
  return 'dark'; // Default theme
}
