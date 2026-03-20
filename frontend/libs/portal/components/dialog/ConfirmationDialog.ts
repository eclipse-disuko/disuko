export enum ConfirmationType {
  NOT_SET = '',
  REVOKE = 'revoke',
  RENEW = 'RENEW',
  DELETE = 'DELETE',
  CONFIRM = 'CONFIRM',
  DEPRECATE = 'DEPRECATE',
}

export interface IConfirmationDialogConfig {
  title?: string;
  type: ConfirmationType;
  contextKey?: string;
  key: string;
  name: string;
  description: string;
  extendedDetails?: string;
  okButton: string;
  okButtonIsDisabled?: boolean;
  emphasiseText?: string;
  emphasiseConfirmationText?: string;
}
