// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export function isURL(str: string) {
  if (str.startsWith('\\\\')) {
    str = 'file://' + str;
  }
  let url;
  try {
    url = new URL(str);
  } catch {
    return false;
  }
  return url.protocol === 'file:' || url.protocol === 'http:' || url.protocol === 'https:' || url.protocol === 'ssh:';
}

export function isURLOrEmpty(str: string) {
  let url;
  if (!str || str.length === 0) {
    return true;
  }
  try {
    url = new URL(str);
  } catch {
    return false;
  }
  return url.protocol === 'http:' || url.protocol === 'https:' || url.protocol === 'ssh:';
}

export function isSpdxIdentifier(str: string, trim = true): boolean {
  if (str) {
    if (trim) {
      str = str.trim();
    }
    if (str.length > 80 || str.length < 2) {
      return false;
    }
    if (str.length > 0) {
      const regex = /^[a-zA-Z0-9\-._+]*$/;
      return regex.test(str);
    }
  }
  return false;
}

export function isSpdxAliasIdentifier(str: string, trim = true): boolean {
  if (str) {
    if (trim) {
      str = str.trim();
    }
    if (str.length > 384 || str.length < 2) {
      return false;
    }
    if (str.length > 0) {
      const regex = /^[a-zA-Z0-9\-._+ ]*$/;
      return regex.test(str);
    }
  }
  return false;
}

export function escapeHtml(unsafe: string) {
  return unsafe
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
    .replace(/'/g, '&#039;');
}
