// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import SimpleProfileData from '@shared/user/models/ProfileData';
import {Rights} from '@shared/user/models/Rights';
import {UserDto} from '@shared/types/Users';
import {defineStore} from 'pinia';

export const useUserStore = defineStore('user', {
  state: () => ({
    simpleProfileData: {
      rights: {} as Rights,
      profile: {} as UserDto,
      allowed: true,
    } as SimpleProfileData,
  }),
  actions: {
    setSimpleProfileData(simpleProfileData: SimpleProfileData) {
      Object.assign(this.simpleProfileData, simpleProfileData);
      this.simpleProfileData.rights = new Rights();
      Object.assign(this.simpleProfileData.rights, simpleProfileData.rights);
      this.simpleProfileData.allowed = true;
    },
    clear() {
      this.simpleProfileData.allowed = false;
    },
    updateTermsOfUse(termsOfUse: boolean, termsOfUseDate: string) {
      this.simpleProfileData.profile.termsOfUse = termsOfUse;
      this.simpleProfileData.profile.termsOfUseDate = termsOfUseDate;
    },
  },
  getters: {
    getRights(): Rights {
      return this.simpleProfileData.rights;
    },
    getProfile(): UserDto {
      return this.simpleProfileData.profile;
    },
  },
});
