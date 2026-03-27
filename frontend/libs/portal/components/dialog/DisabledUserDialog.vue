// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<template>
  <v-dialog v-model="show" content-class="small" persistent>
    <v-card class="pa-8 dDialog">
      <v-card-title>
        <v-row class="align-center">
          <v-col>
            <span class="text-h5">{{ t('DLG_DISABLED_USER_TITLE') }}</span>
          </v-col>
        </v-row>
      </v-card-title>
      <v-card-text>
        <v-row>
          <v-col>
            <span>{{ t('DLG_DISABLED_USER_TEXT') }}</span>
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <DCActionButton
          isDialogButton
          :text="t('BTN_LOGOUT')"
          icon="logout"
          :hint="t('TT_DISABLED_USER_LOGOUT')"
          @click="logoutUser"
        ></DCActionButton>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script>
import {useUserStore} from '@disclosure-portal/stores/user';
import {logout} from '@disclosure-portal/utils/logout';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import {defineComponent, ref} from 'vue';
import {useI18n} from 'vue-i18n';

export default defineComponent({
  name: 'DisabledUserDialog',
  components: {
    DCActionButton,
  },
  setup() {
    const {t} = useI18n();
    const userStore = useUserStore();
    const show = ref(false);

    const open = () => {
      show.value = true;
    };

    const logoutUser = () => {
      userStore.clear();
      logout();
    };

    return {
      open,
      show,
      t,
      logoutUser,
    };
  },
});
</script>
