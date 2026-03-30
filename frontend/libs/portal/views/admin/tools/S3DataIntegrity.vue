// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<template>
  <v-card class="pa-4 mb-3" style="height: auto; border: none !important">
    <v-col cols="12" xs="12">
      <v-row class="d-flex align-start mb-3">
        <v-col cols="2">
          <DCActionButton
            large
            @click="getDbS3CheckGetResult"
            :hint="t('DB_S3_CHECK_GET')"
            :text="t('DB_S3_CHECK_GET')"></DCActionButton>
        </v-col>
        <v-col cols="2" class="mx-2">
          <DCActionButton
            large
            @click="getDbS3CheckStart"
            :hint="t('DB_S3_CHECK_START')"
            :text="t('DB_S3_CHECK_START')"
            class="mr-4"></DCActionButton>
        </v-col>
        <v-col cols="2" v-if="started" class="mx-2">
          <DCActionButton
            large
            @click="getDbS3CheckStop"
            :hint="t('DB_S3_CHECK_STOP')"
            :text="t('DB_S3_CHECK_STOP')"></DCActionButton>
        </v-col>
        <v-col cols="2" v-if="dbS3CheckResult" class="mx-2">
          <DCActionButton
            large
            icon="mdi-download"
            :text="t('BTN_DOWNLOAD')"
            :hint="t('TT_download_notice')"
            class="mr-2"
            @click="downloadDbS3CheckGetResult" />
        </v-col>
        <v-spacer />
      </v-row>
      <v-row v-if="dbS3CheckResult">
        <v-col>
          <json-viewer :value="dbS3CheckResult" :expand-depth="2" aria-expanded="true" theme="jv-dark" sort />
        </v-col>
      </v-row>
    </v-col>
  </v-card>
</template>

<script setup lang="ts">
import AdminService from '@disclosure-portal/services/admin';
import {downloadFile} from '@disclosure-portal/utils/View';
import dayjs from 'dayjs';
import {onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import JsonViewer from 'vue-json-viewer';

const {t} = useI18n();
const dbS3CheckResult = ref<string | null>(null);
const started = ref(false);

const open = async () => {
  await getDbS3CheckGetResult();
};

const getDbS3CheckStart = async () => {
  started.value = true;
  await AdminService.getDbS3CheckStart();
  await getDbS3CheckGetResult();
};

const getDbS3CheckStop = async () => {
  await AdminService.getDbS3CheckStop();
  await getDbS3CheckGetResult();
};

const getDbS3CheckGetResult = async () => {
  dbS3CheckResult.value = null;
  dbS3CheckResult.value = (await AdminService.getDbS3CheckGetResult()).data;
};

const downloadDbS3CheckGetResult = async () => {
  if (!dbS3CheckResult.value) {
    await getDbS3CheckGetResult();
  }
  const filename = 'DbS3Check_' + dayjs(new Date()).format('YYYY-MM-DD_hh_mm_ss') + '.json';
  downloadFile(JSON.stringify(dbS3CheckResult.value, null, 2).replaceAll(',', ',\r'), filename, 'application/json');
};

onMounted(() => {
  open();
  started.value = false;
});
</script>

<style scoped></style>
