// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export interface DiscoForm {
  validate(): Promise<{
    valid: boolean;
    errors: {id: string | number; errorMessages: string[]}[];
  }>;
  reset(): void;

  resetValidation(): void;
}

export interface UIElementDimension {
  clientHeight: string;
  clientWidth: string;
}

export type BufferSource = ArrayBufferView | ArrayBuffer;
export type BlobPart = BufferSource | Blob | string;
