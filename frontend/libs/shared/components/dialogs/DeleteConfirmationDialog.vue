// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script lang="ts" setup>
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import DCloseButton from '@shared/components/disco/DCloseButton.vue';
import {ref} from 'vue';
import {useI18n} from 'vue-i18n';

defineProps<{
  title: string;
  message?: string;
  buttonText?: string;
}>();
const emit = defineEmits(['confirmed']);

const {t} = useI18n();
const isDialogVisible = ref(false);

const showDialog = () => (isDialogVisible.value = true);

const cancelDialog = () => (isDialogVisible.value = false);
const confirmDialog = () => {
  emit('confirmed');
  isDialogVisible.value = false;
};
</script>
<template>
  <slot name="default" :showDialog="showDialog">
    <v-btn text="Replace me" size="small" color="primary" @click.stop="showDialog"></v-btn>
  </slot>
  <v-dialog v-model="isDialogVisible" content-class="small" width="800">
    <v-card class="pa-8 dDialog" flat>
      <v-card-title>
        <v-row>
          <v-col cols="10" class="d-flex align-center">
            <span class="text-h5">{{ t('DLG_CONFIRMATION_TITLE') }}</span>
          </v-col>
          <v-col cols="2" align="right">
            <DCloseButton @click="cancelDialog" />
          </v-col>
        </v-row>
      </v-card-title>
      <v-card-text>
        <v-row>
          <v-col cols="12">
            {{ t('DLG_CONFIRMATION_DESCRIPTION') }}
            <span v-if="title">{{ title }}</span>
          </v-col>
          <v-col v-if="message" cols="12">
            <div class="text-body-1 text-medium-emphasis mt-2">{{ message }}</div>
          </v-col>
        </v-row>
        <v-row>
          <div class="f-modal-alert">
            <div class="f-modal-icon f-modal-warning scaleWarning">
              <span class="f-modal-body pulseWarningIns"></span>
              <span class="f-modal-dot pulseWarningIns"></span>
            </div>
          </div>
        </v-row>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <DCActionButton isDialogButton :text="t('BTN_CANCEL')" class="mr-8" variant="text" @click="cancelDialog" />
        <DCActionButton
          isDialogButton
          :text="buttonText ? buttonText : t('BTN_DELETE')"
          variant="elevated"
          @click="confirmDialog" />
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
