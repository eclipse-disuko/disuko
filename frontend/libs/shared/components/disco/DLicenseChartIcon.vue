<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script lang="ts">
import {ILicenseMetaSlim} from '@disclosure-portal/model/License';
import {TOOLTIP_OPEN_DELAY_IN_MS} from '@shared/utils/constant';
import {defineComponent, toRefs} from 'vue';
import {useI18n} from 'vue-i18n';

export default defineComponent({
  name: 'DLicenseChartIcon',
  props: {
    meta: {
      type: Object as () => ILicenseMetaSlim,
      required: false,
      default: null,
    },
  },
  setup(props) {
    const {meta} = toRefs(props);
    const {t} = useI18n();

    function getLicenseSlimIconText(licMeta: ILicenseMetaSlim): string {
      let result = '';
      if (!licMeta) return '';
      if (licMeta.approvalState === 'forbidden') {
        result = '' + t('ICON_LICENSE_FORBIDDEN_TOOLTIP');
      } else if (licMeta.isLicenseChart) {
        result = '' + t('ICON_LICENSE_CHART_STATUS_TOOLTIP');
      } else if (!licMeta.isLicenseChart) {
        result = '' + t('TABLE_LICENSE_CHART_STATUS_IS_NOT');
      } else if (licMeta.approvalState === 'pending') {
        result = '' + t('ICON_LICENSE_PENDING_STATUS_TOOLTIP');
      } else if (!licMeta.approvalState) {
        result = '' + t('ICON_LICENSE_NOT_SET_STATUS_TOOLTIP');
      }
      return result;
    }
    const getText = (licenseMeta: ILicenseMetaSlim) => {
      return getLicenseSlimIconText(licenseMeta);
    };
    return {
      meta,
      TOOLTIP_OPEN_DELAY_IN_MS,
      getText,
    };
  },
});
</script>

<template>
  <v-tooltip :open-delay="TOOLTIP_OPEN_DELAY_IN_MS" location="bottom" content-class="dpTooltip">
    <template v-slot:activator="{props}">
      <span>
        <span
          v-if="meta.approvalState === 'forbidden'"
          v-bind="props"
          class="dashboard-card-icon material-symbols-outlined"
          style="font-size: 20px; height: 25px; color: rgb(var(--v-theme-licenceForbidden)) !important">
          <v-icon size="x-small" icon="mdi-alert-decagram"></v-icon>
        </span>
        <span
          v-else-if="meta.isLicenseChart"
          v-bind="props"
          class="dashboard-card-icon material-symbols-outlined"
          style="font-size: 20px; height: 25px; color: rgb(var(--v-theme-licenceChartIcon)) !important">
          <v-icon size="x-small" icon="mdi-shield-check-outline"></v-icon>
        </span>
        <v-icon
          v-else-if="!meta.approvalState || meta.approvalState === 'pending' || !meta.isLicenseChart"
          v-bind="props"
          style="height: 25px"
          color="licenceNotApproved"
          size="small"
          tag="span"
          icon="mdi-shield-off-outline">
        </v-icon>
      </span>
    </template>
    <span>{{ getText(meta) }}</span>
  </v-tooltip>
</template>
