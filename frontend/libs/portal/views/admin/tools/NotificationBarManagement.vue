<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import adminService from '@disclosure-portal/services/admin';
import {useAppStore} from '@disclosure-portal/stores/app';
import {computed, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const appStore = useAppStore();

const loading = ref(false);
const enabled = ref(false);
const notificationText = ref('');

const notificationButtonLabel = computed(() => {
  return enabled.value ? 'NOTIFICATION_DISABLE_LABEL' : 'NOTIFICATION_ENABLE_LABEL';
});

const open = async () => {
  const response = await adminService.getNotification();
  enabled.value = response.data.enabled;
  notificationText.value = response.data.text;
};

const previewNotification = () => {
  appStore.notificationMessage = notificationText.value;
  appStore.notificationClosed = false;
};
const postNotification = async () => {
  loading.value = true;
  try {
    enabled.value = !enabled.value;
    await adminService.setNotification({enabled: enabled.value, text: notificationText.value});
    appStore.notificationMessage = notificationText.value;
    appStore.notificationClosed = !enabled.value;
  } finally {
    loading.value = false;
  }
};

onMounted(async () => {
  await open();
});
</script>

<template>
  <v-card class="pa-4 mb-3" style="height: auto; border: none !important">
    <v-col cols="12" xs="12">
      <v-row class="d-flex align-start mb-3">
        <v-col cols="12">
          <v-textarea
            variant="outlined"
            density="compact"
            :label="t('NOTIFICATION_TEXT_LABEL')"
            auto-grow
            v-model="notificationText"
            :disabled="enabled" />
        </v-col>
      </v-row>
      <div class="d-flex align-center" style="gap: 24px">
        <DCActionButton :disabled="enabled" @click="previewNotification" hint="Preview" text="Preview" large />
        <DCActionButton
          @click="postNotification"
          :hint="t(notificationButtonLabel)"
          :text="t(notificationButtonLabel)"
          large
          :loading="loading" />
        <span>{{ t('NOTIFICATION_STATUS') }}: {{ loading ? 'loading' : enabled }}</span>
      </div>
      <div class="mt-4">
        <span>Note: It can take up to 30 seconds for the changes to take effect.</span>
      </div>
    </v-col>
  </v-card>
</template>
