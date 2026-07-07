// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export interface MailTemplate {
  _key: string;
  subject: string;
  message: string;
  bcc: string;
  cc: string;
  values: Record<string, string>;
}

export interface UpdateMailTemplate {
  subject: string;
  message: string;
  bcc: string;
  cc: string;
}
