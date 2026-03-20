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
