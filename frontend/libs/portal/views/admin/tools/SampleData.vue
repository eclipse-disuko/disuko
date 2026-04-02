<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import SampleDataCreationState from '@disclosure-portal/model/SampleData';
import AdminService from '@disclosure-portal/services/admin';
import {formatDateAndTime} from '@disclosure-portal/utils/Table';
import {onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const serverResponse = ref('');
const started = ref(false);
const serverResponseDetails = ref<SampleDataCreationState | null>(null);
const cntSampleData = ref(1);

const startSampleData = async (withFileUpload: boolean) => {
  started.value = true;
  serverResponse.value = '';
  serverResponseDetails.value = null;
  serverResponse.value = (await AdminService.triggerCreateSampleData(cntSampleData.value, withFileUpload)).data;
  await getSampleDataState();
};

const stopSampleData = async () => {
  serverResponse.value = '';
  serverResponseDetails.value = null;
  serverResponse.value = (await AdminService.stopCreateSampleData()).data;
  started.value = false;
};

const getSampleDataState = async () => {
  serverResponse.value = '';
  serverResponseDetails.value = (await AdminService.getCreateSampleDataState()).data;
};

onMounted(() => {
  started.value = false;
});
</script>

<template>
  <v-card class="pa-4 mb-3" style="height: auto; border: none !important">
    <v-col cols="12" xs="12">
      <v-row class="d-flex align-start mb-3">
        <v-col cols="2" class="d-flex">
          <v-text-field
            autocomplete="off"
            variant="outlined"
            density="compact"
            label="Count projects"
            v-model="cntSampleData"></v-text-field>
        </v-col>
        <v-col cols="1" class="d-flex justify-center mr-2">
          <DCActionButton
            @click="startSampleData(false)"
            :hint="t('TT_CREATE_SAMPLE_DATA')"
            :text="t('TT_CREATE_SAMPLE_DATA')"
            large></DCActionButton>
        </v-col>
        <v-col cols="2" class="d-flex justify-center">
          <DCActionButton
            @click="startSampleData(true)"
            :hint="t('TT_CREATE_SAMPLE_DATA_WITH_FILE')"
            :text="t('TT_CREATE_SAMPLE_DATA_WITH_FILE')"
            large></DCActionButton>
        </v-col>
        <v-col cols="2" class="d-flex justify-center">
          <DCActionButton
            large
            @click="getSampleDataState"
            :hint="t('TT_GET_SAMPLE_DATA_STATE')"
            :text="t('TT_GET_SAMPLE_DATA_STATE')"></DCActionButton>
        </v-col>
        <v-col cols="1" class="d-flex justify-center">
          <DCActionButton
            large
            @click="stopSampleData"
            :hint="t('TT_CLEAR_SAMPLE_DATA_STATE')"
            :text="t('TT_CLEAR_SAMPLE_DATA_STATE')"></DCActionButton>
        </v-col>
        <v-col cols="2">
          <div class="d-text" v-if="serverResponse" v-html="serverResponse" />
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="8">
          <v-table v-if="serverResponseDetails" fixed-header density="compact" class="borderTable">
            <template v-slot:default>
              <thead height="50">
                <tr>
                  <th class="text-left">
                    {{ t('COL_NAME') }}
                  </th>
                  <th class="text-left">
                    {{ t('COL_VALUES') }}
                  </th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td>Is Running</td>
                  <td>{{ serverResponseDetails.isRunning }}</td>
                </tr>
                <tr>
                  <td>Start Time</td>
                  <td>{{ formatDateAndTime(serverResponseDetails.startTime) }}</td>
                </tr>
                <tr>
                  <td>End Time</td>
                  <td>{{ formatDateAndTime(serverResponseDetails.endTime) }}</td>
                </tr>
                <tr>
                  <td>Has Errors</td>
                  <td>{{ serverResponseDetails.hasErrors }}</td>
                </tr>
                <tr>
                  <td>Last Error</td>
                  <td>{{ serverResponseDetails.lastError }}</td>
                </tr>
                <tr>
                  <td>With File Upload</td>
                  <td>{{ serverResponseDetails.withFileUpload }}</td>
                </tr>
                <tr>
                  <td>Created Count</td>
                  <td>{{ serverResponseDetails.createdCount }}</td>
                </tr>
                <tr>
                  <td>Target Count</td>
                  <td>{{ serverResponseDetails.targetCount }}</td>
                </tr>
                <tr>
                  <td>RequestId</td>
                  <td>{{ serverResponseDetails.reqID }}</td>
                </tr>
              </tbody>
            </template>
          </v-table>
        </v-col>
      </v-row>
    </v-col>
  </v-card>
</template>

<style scoped></style>
