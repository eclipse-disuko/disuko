// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {Rights} from '@shared/types/Rights';
import {UserDto} from '@shared/types/Users';

export default interface SimpleProfileData {
  rights: Rights;
  profile: UserDto;
  allowed: boolean;
}
