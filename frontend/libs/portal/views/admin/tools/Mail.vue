<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {MailData} from '@disclosure-portal/model/MailData';
import SampleDataCreationState from '@disclosure-portal/model/SampleData';
import AdminService from '@disclosure-portal/services/admin';
import {ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();

const mailData = ref<MailData>({} as MailData);
const serverResponse = ref<string>('');
const serverResponseDetails = ref<SampleDataCreationState | null>(null);

const sendMail = async () => {
  serverResponseDetails.value = (await AdminService.sendEmail(mailData.value)).data;
};
</script>

<template>
  <v-card class="pa-4 mb-3" style="height: auto; border: none !important">
    <v-row class="mb-3">
      <v-col cols="12">
        <h1 class="text-h5 pr-2">{{ t('TITLE_MAIL') }}</h1>
      </v-col>
    </v-row>
    <v-col cols="12" xs="12">
      <v-row class="d-flex align-start mb-3">
        <v-col cols="12" xs="6" md="4">
          <v-text-field
            autocomplete="off"
            variant="outlined"
            density="compact"
            :label="t('MAIL_TYPE')"
            v-model="mailData.type"></v-text-field>
        </v-col>
        <v-col cols="12" xs="6" md="4">
          <DCActionButton large @click="sendMail" :text="t('BTN_SEND_MAIL')"></DCActionButton>
        </v-col>
        <v-col cols="2">
          <div class="d-text" v-if="serverResponse" v-html="serverResponse" />
        </v-col>
      </v-row>
    </v-col>
  </v-card>
</template>
