// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<template>
  <TableLayout has-tab has-title>
    <template #buttons>
      <DCActionButton
        :text="t('BTN_DOWNLOAD')"
        large
        v-if="RightsUtils.rights().isProjectAnalyst()"
        icon="mdi-download"
        :hint="t('TT_download_label_report')"
        @click="downloadReport" />
    </template>
    <template #table>
      <v-table fixed-header density="compact" class="striped-table fill-height">
        <thead>
          <tr>
            <th class="text-left">
              {{ t('COL_NAME') }}
            </th>
            <th class="text-center">
              {{ t('COL_COUNT') }}
            </th>
          </tr>
        </thead>
        <tbody v-if="stats">
          <tr>
            <td>{{ t('COL_KEY_PROJECTS') }}</td>
            <td class="text-center">
              {{ stats.projectCount }}
            </td>
          </tr>
          <tr>
            <td>{{ t('COL_KEY_PROJECTS_ACTIVE') }}</td>
            <td class="text-center">
              {{ stats.projectActiveCount }}
            </td>
          </tr>
          <tr>
            <td>{{ t('COL_KEY_PROJECTS_DELETED') }}</td>
            <td class="text-center">
              {{ stats.projectDeletedCount }}
            </td>
          </tr>
          <tr>
            <td>{{ t('COL_KEY_LICENSE_COUNT') }}</td>
            <td class="text-center">
              {{ stats.licenseCount }}
            </td>
          </tr>
          <tr>
            <td>{{ t('COL_KEY_LICENSE_CHART_COUNT') }}</td>
            <td class="text-center">
              {{ stats.licenseChartCount }}
            </td>
          </tr>
          <tr>
            <td>{{ t('COL_KEY_LICENSE_ACTIVE_COUNT') }}</td>
            <td class="text-center">
              {{ stats.licenseActiveCount }}
            </td>
          </tr>
          <tr>
            <td>{{ t('COL_KEY_LICENSE_DELETED_COUNT') }}</td>
            <td class="text-center">
              {{ stats.licenseDeletedCount }}
            </td>
          </tr>
          <tr>
            <td>{{ t('COL_KEY_LICENSE_FORBIDDEN_COUNT') }}</td>
            <td class="text-center">
              {{ stats.licenseForbiddenCount }}
            </td>
          </tr>
          <tr>
            <td>{{ t('COL_KEY_LICENSE_UNKNOWN_COUNT') }}</td>
            <td class="text-center">
              {{ stats.licenseUnknownCount }}
            </td>
          </tr>
          <tr>
            <td>{{ t('COL_KEY_SBOM_UPLOAD_COUNT') }}</td>
            <td class="text-center">
              {{ stats.uploadFileCntSBOM }}
            </td>
          </tr>
          <tr>
            <td>{{ t('COL_KEY_USERS_COUNT') }}</td>
            <td class="text-center">
              {{ stats.userCount }}
            </td>
          </tr>
          <tr>
            <td>{{ t('COL_KEY_USERS_ACTIVE_COUNT') }}</td>
            <td class="text-center">
              {{ stats.userActiveCount }}
            </td>
          </tr>
          <tr>
            <td>{{ t('COL_KEY_USERS_DEACTIVE_COUNT') }}</td>
            <td class="text-center">
              {{ stats.userDeactivateCount }}
            </td>
          </tr>
          <tr>
            <td>{{ t('COL_KEY_USERS_TOSACCEETED_COUNT') }}</td>
            <td class="text-center">
              {{ stats.userTermsNotAcceptedCount }}
            </td>
          </tr>
          <tr>
            <td>{{ t('COL_KEY_TRAININGS_COMPLETED_COUNT') }}</td>
            <td class="text-center">
              {{ stats.completedTrainings }}
            </td>
          </tr>
        </tbody>
      </v-table>
    </template>
  </TableLayout>
</template>

<script lang="ts" setup>
import {Stats} from '@disclosure-portal/model/Analytics';
import AnalyticsService from '@disclosure-portal/services/analytics';
import {downloadFile} from '@disclosure-portal/utils/download';
import {RightsUtils} from '@disclosure-portal/utils/Rights';
import dayjs from 'dayjs';
import {onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const stats = ref<Stats | null>(null);

const downloadReport = async () => {
  const filename = `report_${dayjs().format('YYYY-MM-DD_hh_mm_ss')}.csv`;
  downloadFile(filename, AnalyticsService.downloadReport(), true);
};

onMounted(async () => {
  stats.value = (await AnalyticsService.getStats()).data as Stats;
});
</script>
