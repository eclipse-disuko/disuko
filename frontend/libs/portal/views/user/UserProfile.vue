<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script lang="ts">
import UserMain from '@disclosure-portal/components/user/UserMain.vue';
import {UserDto} from '@shared/types/Users';
import {RightsUtils} from '@disclosure-portal/utils/Rights';
import {defineComponent, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';

import {useUserStore} from '@disclosure-portal/stores/user';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';

export default defineComponent({
  name: 'UserProfile',
  components: {
    UserMain,
  },

  setup() {
    const {t} = useI18n();
    const userStore = useUserStore();
    const {dashboardCrumbs, ...breadcrumbs} = useBreadcrumbsStore();

    const userProfile = ref<UserDto>({} as UserDto);
    const rolesSwitchable = ref(false);
    const reloadUserProfile = () => {
      userProfile.value = userStore.getProfile;
      rolesSwitchable.value = RightsUtils.hasUsersAccess();
    };

    const initBreadcrumbs = () => {
      breadcrumbs.setCurrentBreadcrumbs([
        ...dashboardCrumbs,
        {
          title: t('BTN_profile'),
          disabled: false,
          href: '/dashboard/user',
        },
      ]);
    };

    onMounted(() => {
      reloadUserProfile();
      initBreadcrumbs();
    });

    return {
      userProfile,
      rolesSwitchable,
      reloadUserProfile,
    };
  },
});
</script>

<template>
  <UserMain
    :has-users-access="false"
    :roles-switchable="rolesSwitchable"
    :user-profile="userProfile"
    @reloadUserProfile="reloadUserProfile" />
</template>
