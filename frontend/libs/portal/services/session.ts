// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {getApi} from '@disclosure-portal/api';
import {RefreshTokenRequestDto, RefreshTokenResponseDto} from '@disclosure-portal/model/Credentials';
import {AxiosResponse} from 'axios';

const {api} = getApi();

class SessionService {
  public async getRefreshAccessToken(): Promise<AxiosResponse<RefreshTokenResponseDto>> {
    const refreshTokenRequestDto = new RefreshTokenRequestDto();
    return await api.post('/api/v1/refreshToken', refreshTokenRequestDto);
  }
}

const sessionService = new SessionService();

export default sessionService;
