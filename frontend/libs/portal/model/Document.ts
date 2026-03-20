// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export class ChainDetails {
  public updated = '';
  public hash = '';
}

export class DocumentDto {
  public _key = '';
  public created = '';
  public updated = '';
  public approvalId = '';
  public type = '';
  public lang = '';
  public description = '';
  public hash = '';
  public fileName = '';
  public version = '';
  public chain: ChainDetails[] = [];
}
