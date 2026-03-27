// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import GridAuditLog from '@disclosure-portal/components/grids/GridAuditLog.vue';
import adminService from '@disclosure-portal/services/admin';
import {ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const show = ref(false);
const id = ref('');
const title = ref('');

const open = (newId: string, name: string) => {
  id.value = newId;
  title.value = `${t('TITLE_AUDIT_LOG')}: ${name}`;
  show.value = true;
};

const close = () => {
  show.value = false;
};
</script>

<template>
  <slot name="default" :open="open">
    <v-btn text="Replace me" size="small" color="primary" @click.stop="open"></v-btn>
  </slot>
  <v-dialog v-model="show" content-class="large" scrollable>
    <v-card class="pa-8 dDialog" flat>
      <v-card-title>
        <v-row>
          <v-col cols="10">
            <span class="text-h5">{{ title }}</span>
          </v-col>
          <v-col cols="2" align="right">
            <v-btn icon @click="close">
              <v-icon>mdi-close</v-icon>
            </v-btn>
          </v-col>
        </v-row>
      </v-card-title>
      <v-card-text>
        <GridAuditLog :fetch-method="() => adminService.getClassificationAuditTrail(id)" :key="id" />
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn @click="close" depressed color="primary">
          {{ t('BTN_CLOSE') }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
