// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import ConfirmationDialog from '@disclosure-portal/components/dialog/ConfirmationDialog.vue';
import {JOB_STATUS_SUCCESS, JOB_TYPE_TERMS_OF_USE, JobDto, SetConfigDto} from '@disclosure-portal/model/Job';
import AdminService from '@disclosure-portal/services/admin';
import {useUserStore} from '@disclosure-portal/stores/user';
import {formatDateAndTime} from '@disclosure-portal/utils/Table';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import {computed, nextTick, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const userStore = useUserStore();
const termsOfUseVersion = ref('');
const termsOfUseCurrentVersion = ref('');
const job = ref<JobDto>({} as JobDto);
const isInProgress = ref(false);
const dlgConfirmationDialog = ref<InstanceType<typeof ConfirmationDialog> | null>(null);
const getJobStatus = async (): Promise<void> => {
  const refreshStatus = setInterval(async () => {
    const serverResponse = await AdminService.getJobLatest(JOB_TYPE_TERMS_OF_USE);
    job.value = serverResponse.data;
    if (job.value.status === JOB_STATUS_SUCCESS) {
      clearInterval(refreshStatus);
      isInProgress.value = false;
    }
  }, 1000);
};

const init = async () => {
  await getJobStatus();
  const response = await AdminService.getTermsOfUseCurrentVersion();
  termsOfUseCurrentVersion.value = response.termsOfUseCurrentVersion;
};

const triggerResetTermsAcceptanceJob = () => {
  dlgConfirmationDialog.value?.openWithoutDetails('DLG_CONFIRMATION_DESCRIPTION_RESET_TERMS_ACCEPTANCE', 'BTN_RESET');
};

const doTriggerResetTermsAcceptanceJob = async () => {
  isInProgress.value = true;
  await AdminService.setJobConfig(JOB_TYPE_TERMS_OF_USE, new SetConfigDto(termsOfUseVersion.value)).then(async () => {
    await AdminService.startJob(JOB_TYPE_TERMS_OF_USE);
    await getJobStatus();
  });
};

onMounted(() => {
  nextTick(init);
});

const userHasPermission = computed(() => {
  return userStore.getRights.isApplicationAdmin();
});
</script>

<template>
  <v-card class="pa-4">
    <v-col cols="12" xs="12">
      <v-row class="align-start mb-3">
        <v-col cols="12" xs="6" md="4">
          <v-text-field
            autocomplete="off"
            disabled
            variant="outlined"
            density="compact"
            :label="t('LBL_TERMS_VERSION_CURRENT')"
            v-model="termsOfUseCurrentVersion"></v-text-field>
        </v-col>
        <v-col cols="12" xs="6" md="4">
          <v-text-field
            autocomplete="off"
            variant="outlined"
            density="compact"
            :label="t('LBL_TERMS_VERSION_NEXT')"
            v-model="termsOfUseVersion"
            class="required"></v-text-field>
        </v-col>
        <v-col cols="12" class="justify-center mr-2" xs="6" md="4">
          <DCActionButton
            :disabled="!termsOfUseVersion || isInProgress || !userHasPermission"
            @click="triggerResetTermsAcceptanceJob"
            :hint="t('BTN_RESET_TERMS_ACCEPTANCE')"
            :text="t('BTN_RESET_TERMS_ACCEPTANCE')"
            large></DCActionButton>
        </v-col>
      </v-row>
      <v-row class="d-flex align-start mb-3">
        <v-col cols="12" xs="12" md="9" v-if="job">
          <span class="d-subtitle-2">Last ran: </span>
          <span class="d-subtitle-2">{{ formatDateAndTime(job.updated) }}</span>
          <pre class="d-subtitle-2 pt-2" v-if="job" v-text="job.customRes" />
        </v-col>
      </v-row>
    </v-col>
    <ConfirmationDialog ref="dlgConfirmationDialog" @confirm="doTriggerResetTermsAcceptanceJob"></ConfirmationDialog>
  </v-card>
</template>
