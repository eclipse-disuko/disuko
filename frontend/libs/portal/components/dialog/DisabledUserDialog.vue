<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {useUserStore} from '@shared/user/stores/user.store';
import {logout} from '@disclosure-portal/utils/logout';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const userStore = useUserStore();
const show = ref(false);

const dialogConfig = computed(() => ({
  title: t('DLG_DISABLED_USER_TITLE'),
  primaryButton: {
    text: t('BTN_LOGOUT'),
    hint: t('TT_DISABLED_USER_LOGOUT'),
    icon: 'mdi mdi-logout',
    iconColor: 'default',
  },
  showCloseButton: false,
}));

const open = () => {
  show.value = true;
};

const logoutUser = () => {
  userStore.clear();
  logout();
};

defineExpose({open});
</script>
<template>
  <v-dialog v-model="show" width="600" persistent>
    <DialogLayout :config="dialogConfig" @primary-action="logoutUser">
      <Stack>
        <span>{{ t('DLG_DISABLED_USER_TEXT') }}</span>
      </Stack>
    </DialogLayout>
  </v-dialog>
</template>
