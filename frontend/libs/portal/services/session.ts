// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useApi} from '@shared/api/useApi';
import {RefreshTokenRequestDto, RefreshTokenResponseDto} from '@shared/types/Credentials';
import {AxiosResponse} from 'axios';

const {api} = useApi();

class SessionService {
  public async getRefreshAccessToken(): Promise<AxiosResponse<RefreshTokenResponseDto>> {
    const refreshTokenRequestDto = new RefreshTokenRequestDto();
    return await api.post('/api/v1/refreshToken', refreshTokenRequestDto);
  }
}

const sessionService = new SessionService();

export default sessionService;
