// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import {UserDto} from '@disclosure-portal/model/Users';
import AdminService from '@disclosure-portal/services/admin';
import {useUserStore} from '@disclosure-portal/stores/user';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {computed, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRoute} from 'vue-router';

const {t} = useI18n();
const route = useRoute();

const userProfile = ref<UserDto>({} as UserDto);
const rolesSwitchable = ref(false);
const hasUserAccess = ref(false);
const {dashboardCrumbs, ...breadcrumbs} = useBreadcrumbsStore();
const userStore = useUserStore();
const user = computed(() => userStore.getProfile.user);

const initBreadcrumbs = () => {
  breadcrumbs.setCurrentBreadcrumbs([
    ...dashboardCrumbs,
    {
      title: t('TITLE_USERS'),
      disabled: false,
      href: '/dashboard/admin/users',
    },
    {
      title: userProfile.value.user || '',
      disabled: false,
      href: `/dashboard/admin/users/${encodeURIComponent(userProfile.value._key)}`,
    },
  ]);
};

const reloadUserProfile = async (uuid: string) => {
  const response = await AdminService.getUser(uuid);
  userProfile.value = response.data;
  rolesSwitchable.value = userProfile.value.user === user.value;
  hasUserAccess.value = true;
};

onMounted(async () => {
  await reloadUserProfile(route.params.uuid as string);
  initBreadcrumbs();
});
</script>

<template>
  <UserMain
    :has-users-access="hasUserAccess"
    :roles-switchable="rolesSwitchable"
    :user-profile="userProfile"
    @reloadUserProfile="reloadUserProfile" />
</template>
