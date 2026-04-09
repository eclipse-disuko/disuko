<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {useLicense} from '@disclosure-portal/composables/useLicense';
import LicenseModel from '@disclosure-portal/model/License';
import {formatDate, formatDateTime, getStrWithMaxLength} from '@disclosure-portal/utils/View';
import {TOOLTIP_OPEN_DELAY_IN_MS} from '@shared/utils/constant';
import {useI18n} from 'vue-i18n';

defineProps<{
  license: LicenseModel;
}>();

const {t} = useI18n();
const {getI18NTextOfPrefixKey} = useLicense();
</script>

<template>
  <TableLayout has-title has-tab>
    <template #table>
      <v-row class="pa-4" v-if="license && license.meta">
        <v-col cols="12" xs="6" sm="6" md="3" lg="2">
          <span class="text-caption text-grey-darken-1">{{ t('CD_LICENSE_ID') }}</span
          ><br />
          <span class="text-body-2">{{ license.licenseId }}</span>
        </v-col>
        <v-col cols="12" xs="6" sm="6" md="3" lg="2">
          <span class="text-caption text-grey-darken-1">{{ t('COL_LICENSE_FAMILY') }}</span
          ><br />
          <span class="text-body-2">{{ getI18NTextOfPrefixKey('LIC_FAMILY_', license.meta.family as string) }}</span>
        </v-col>
        <v-col cols="12" xs="6" sm="6" md="3" lg="1">
          <span class="text-caption text-grey-darken-1">{{ t('COL_LICENSE_SOURCE') }}</span
          ><br />
          <span class="text-body-2">{{ license.source }}</span>
        </v-col>
        <v-col v-if="license.isDeprecatedLicenseId" cols="12" xs="6" sm="6" md="3" lg="1">
          <span class="text-caption text-grey-darken-1">{{ t('COL_LICENSE_SPDX_STATUS') }}</span
          ><br />
          <span class="text-body-2">{{ t('SPDX_STATUS_VALUE') }}</span>
        </v-col>
        <v-col cols="12" xs="6" sm="6" md="3" lg="1">
          <span class="text-caption text-grey-darken-1">{{ t('LICENSE_CHART_TITLE') }}</span
          ><br />
          <span class="text-body-2">{{ t(`LICENSE_CHART_STATUS_${license.meta.isLicenseChart}`) }}</span>
        </v-col>
        <v-col cols="12" xs="6" sm="6" md="3" lg="1">
          <span class="text-caption text-grey-darken-1">{{ t('Lbl_created') }}</span
          ><br />
          <v-tooltip :open-delay="TOOLTIP_OPEN_DELAY_IN_MS" location="bottom" content-class="dpTooltip">
            <template v-slot:activator="{props}">
              <span class="text-body-2" v-bind="props">{{ formatDate(license.created) }}</span>
            </template>
            <span class="text-body-2">{{ formatDateTime(license.created) }}</span
            ><br />
          </v-tooltip>
        </v-col>
        <v-col cols="12" xs="6" sm="6" md="3" lg="1">
          <span class="text-caption text-grey-darken-1">{{ t('Lbl_updated') }}</span
          ><br />
          <v-tooltip :open-delay="TOOLTIP_OPEN_DELAY_IN_MS" location="bottom" content-class="dpTooltip">
            <template v-slot:activator="{props}">
              <span class="text-body-2" v-bind="props">{{ formatDate(license.updated) }}</span>
            </template>
            <span class="text-body-2">{{ formatDateTime(license.updated) }}</span
            ><br />
          </v-tooltip>
        </v-col>
        <v-col cols="12" xs="6" sm="6" md="3" lg="3">
          <span class="text-caption text-grey-darken-1">{{ t('COL_LICENSE_TYPE') }}</span
          ><br />
          <span class="text-body-2">{{ getI18NTextOfPrefixKey('LT_', license.meta.licenseType) }}</span>
        </v-col>
        <v-col cols="12" xs="6" sm="6" md="3" lg="2">
          <span class="text-caption text-grey-darken-1">{{ t('COL_APPROVAL_STATUS') }}</span
          ><br />
          <span
            class="text-body-2"
            :style="license.meta.approvalState === 'forbidden' ? 'color: var(--v-licenceForbidden-base)' : ''">
            {{ getI18NTextOfPrefixKey('LT_APP_', license.meta.approvalState) }}</span
          >
        </v-col>
        <v-col cols="12" xs="6" sm="6" md="3" lg="2">
          <span class="text-caption text-grey-darken-1">{{ t('COL_REVIEW_STATUS') }}</span
          ><br />
          <span class="text-body-2">{{ getI18NTextOfPrefixKey('LT_RS_', license.meta.reviewState) }}</span>
        </v-col>
        <v-col cols="12" v-if="license.meta.reviewState === 'reviewed'" xs="6" sm="6" md="3" lg="1">
          <span class="text-caption text-grey-darken-1">{{ t('COL_REVIEW_DATE') }}</span
          ><br />
          <span class="text-body-2">{{ formatDate(license.meta.reviewDate) }}</span
          ><br />
        </v-col>
        <v-col cols="12">
          <span class="text-caption text-grey-darken-1">{{ t('COL_LICENSE_URL') }}</span
          ><br />
          <DExternalLink
            :text="getStrWithMaxLength(45, license.meta.licenseUrl)"
            :url="license.meta.licenseUrl"
            tooltip></DExternalLink>
          <br />
        </v-col>
        <v-col cols="12">
          <span class="text-caption text-grey-darken-1">{{ t('COL_SOURCE_URL') }}</span
          ><br />
          <DExternalLink
            :text="getStrWithMaxLength(45, license.meta.sourceUrl)"
            :url="license.meta.sourceUrl"
            tooltip></DExternalLink>
          <br />
        </v-col>
      </v-row>
    </template>
  </TableLayout>
</template>
