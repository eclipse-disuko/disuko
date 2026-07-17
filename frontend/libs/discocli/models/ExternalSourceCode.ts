export interface ExternalSourceCodeReference {
  _key?: string;
  url: string;
  comment?: string;
  created?: string;
  origin?: string;
  uploader?: string;
  fileSize?: number;
  hash?: string;
  sourceType?: string;
}

export interface ExtendedExternalSource extends ExternalSourceCodeReference {
  _displayUrl: string;
  _isFileUrl: boolean;
  _displayComment: string;
  _displayOrigin: string;
  _displayUploader: string;
  _referenceInfo: string;
  _createdString: string;
}
