// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import DHTTPError from '@disclosure-portal/model/DHTTPError';
import {INotificationMeta} from '@disclosure-portal/model/IdleInfo';
import eventBus from '@disclosure-portal/utils/eventbus';
import {AxiosError, AxiosInstance, AxiosResponse, InternalAxiosRequestConfig} from 'axios';

export const initInterceptors = (axiosInstance: AxiosInstance) => {
  axiosInstance.interceptors.request.use(
    (config: InternalAxiosRequestConfig) => {
      config.headers['X-Client-Version'] = '2.0';
      return config;
    },
    (error) => {
      return Promise.reject(error);
    },
  );

  // Interceptor: Response-Handling
  axiosInstance.interceptors.response.use(
    (response: AxiosResponse) => {
      readNotificationFromHeader(response);
      return response;
    },
    async (axiosError: AxiosError) => {
      if (
        axiosError &&
        (axiosError.code === 'ERR_NETWORK' ||
          axiosError.code === 'ERR_CONNECTION_REFUSED' ||
          axiosError.code === 'ERR_SSL_PROTOCOL_ERROR')
      ) {
        eventBus.emit('on-api-error', createErrorNetwork());
        return Promise.reject(axiosError);
      }
      const axiosErrorResponse = axiosError.response;
      if (axiosErrorResponse) {
        const status = axiosErrorResponse.status;
        let httpError: DHTTPError;

        if (status === 503) {
          httpError = createError503();
        } else if (status === 401) {
          httpError = createError401();
        } else {
          if (axiosErrorResponse.data) {
            if (axiosErrorResponse.data instanceof Blob) {
              const rawError = await axiosErrorResponse.data.text();
              httpError = await extractErrorDtoAndCreateError(rawError, axiosErrorResponse);
            } else {
              httpError = await extractErrorDtoAndCreateError(axiosErrorResponse.data, axiosErrorResponse);
            }
          } else {
            httpError = createAndLogErrorWithNoDetails(axiosErrorResponse);
          }
        }
        eventBus.emit('on-api-error', httpError);
      }
      return Promise.reject(axiosError);
    },
  );

  const StatusPreconditionRequired = 428;

  const readNotificationFromHeader = (response: AxiosResponse) => {
    if (response.headers['x-notification']) {
      try {
        const notification = JSON.parse(response.headers['x-notification']);
        if (notification) {
          const config = {
            enabled: notification.enabled,
            text: notification.text,
          } as INotificationMeta;
          eventBus.emit('set-notification', {config: config});
        }
      } catch (e) {
        console.error(e);
      }
    }
  };

  const createError = (statusCode: number, title: string, raw: string, message: string, reqId: string): DHTTPError => {
    const error = new DHTTPError();
    if (statusCode !== StatusPreconditionRequired) {
      error.reqId = reqId;
    }
    error.title = title;
    error.raw = raw;
    error.code = '' + statusCode;
    error.message = message;
    return error;
  };

  const createErrorNetwork = (): DHTTPError => createError(503, 'ERROR_NETWORK', '', 'ERROR_NETWORK_MESSAGE', '');

  const createError503 = (): DHTTPError =>
    createError(503, 'ERROR_503', 'ERROR_TITLE_MAINTENANCE', 'ERROR_DESCRIPTION_MAINTENANCE', '');

  const createError401 = (): DHTTPError =>
    createError(401, 'ERROR_401', 'MISSING_AUTHENTICATION', 'MISSING_AUTHENTICATION', '');

  const extractErrorDtoAndCreateError = async (
    rawError: any,
    axiosErrorResponse: AxiosResponse,
  ): Promise<DHTTPError> => {
    const errorDto = extractErrorDto(rawError);
    if (errorDto) {
      return createErrorDetailsWithNoRaw(errorDto, axiosErrorResponse);
    }
    return createAndLogErrorWithNoDetails(axiosErrorResponse);
  };

  const extractErrorDto = (rawError: any) => {
    if (rawError) {
      if (rawError.message) {
        return rawError;
      } else {
        return tryParseRawError(rawError);
      }
    }
    return null;
  };

  const tryParseRawError = (rawError: any) => {
    try {
      const start = rawError.indexOf('{"');
      const end = rawError.indexOf('}');
      const cleanJson = rawError.substring(start, end + 1);
      return JSON.parse(cleanJson);
    } catch {
      return null;
    }
  };

  const createErrorDetailsWithNoRaw = (errorDto: any, axiosErrorResponse: AxiosResponse): DHTTPError => {
    const code = errorDto.code ? errorDto.code : axiosErrorResponse.statusText;
    const message = errorDto.message
      ? errorDto.message
      : axiosErrorResponse.statusText + ' (' + axiosErrorResponse.status + ')';
    return createError(axiosErrorResponse.status, code, '', message, errorDto.reqID);
  };

  const createAndLogErrorWithNoDetails = (axiosErrorResponse: AxiosResponse): DHTTPError =>
    createErrorWithNoDetails(axiosErrorResponse);

  const createErrorWithNoDetails = (axiosErrorResponse: AxiosResponse): DHTTPError =>
    createError(
      axiosErrorResponse.status,
      axiosErrorResponse.statusText,
      axiosErrorResponse.statusText + ' (' + axiosErrorResponse.status + ')',
      'SEE_CONSOLE_LOGS',
      '',
    );
};
