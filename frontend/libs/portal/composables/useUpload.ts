// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useApi} from '@shared/api/useApi';
import SchemaPostRequest from '@disclosure-portal/model/SchemaPostRequest';
import {AxiosProgressEvent} from 'axios';

export const useUpload = () => {
  const {api} = useApi();

  const uploadFormDataFile = async ({
    uploadUrl,
    file,
    onUploadProgress,
    schema,
  }: {
    uploadUrl: string;
    file: File;
    onUploadProgress: (progressEvent: AxiosProgressEvent) => void;
    schema?: SchemaPostRequest;
  }) => {
    const formData = new FormData();
    formData.append('file', file as Blob);

    if (schema) {
      formData.append('schema', JSON.stringify(schema));
    }

    return api.post(uploadUrl, formData, {
      withCredentials: true,
      headers: {'Content-Type': 'multipart/form-data'},
      onUploadProgress,
    });
  };

  return {
    uploadFormDataFile,
  };
};
