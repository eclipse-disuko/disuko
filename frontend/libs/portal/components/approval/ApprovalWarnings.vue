<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import config from '@shared/utils/config';
import {useI18n} from 'vue-i18n';

interface Props {
  isDeniedOrUnasserted?: boolean;
  isEitherFutureFoss?: boolean;
  isRdConfirmationMissing?: boolean;
  isWarned?: boolean;
  isEnterpriseOrMobileOrOther?: boolean;
  noFOSS?: boolean;
  mixedFOSS?: boolean;
  fossVersion?: 'default' | 'legacy';
  selectedProjectsContainEmptySbom?: boolean;
  isNoSbomNoFossWarning?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  isDeniedOrUnasserted: false,
  isEitherFutureFoss: false,
  isRdConfirmationMissing: false,
  isWarned: false,
  isEnterpriseOrMobileOrOther: false,
  noFOSS: false,
  mixedFOSS: false,
  fossVersion: 'legacy',
  selectedProjectsContainEmptySbom: false,
  isNoSbomNoFossWarning: false,
});

const {t} = useI18n();
</script>

<template>
  <section
    v-if="
      props.isDeniedOrUnasserted ||
      (props.isEitherFutureFoss && props.isRdConfirmationMissing) ||
      (props.isEitherFutureFoss && props.isWarned) ||
      props.isEnterpriseOrMobileOrOther
    ">
    <v-alert color="warning" type="warning">
      <span v-if="props.isDeniedOrUnasserted">
        {{ t('DENIED_OR_UNASSARETED_MESSAGE') }}
      </span>
      <span v-else-if="props.isEitherFutureFoss && props.isRdConfirmationMissing">
        {{ t('CONFIRMATION_MISSING') }}
      </span>
      <span v-else-if="props.isEitherFutureFoss && props.isWarned">
        {{ t('WARNING_MESSAGE') }}
      </span>
      <span v-else-if="props.isEnterpriseOrMobileOrOther">
        {{ t('ENTERPRISE_MOBILE_OTHER_MESSAGE') }}
        <a :href="t('ENTERPRISE_MOBILE_OTHER_MESSAGE_CTA')" target="_blank">
          <v-icon>mdi mdi-chevron-right</v-icon>
          <span>{{ t('LINK_CLICK_HERE') }} </span>
        </a>
      </span>
    </v-alert>
  </section>

  <section v-if="props.mixedFOSS">
    <v-alert color="warning" type="warning">
      {{ t('MIXED_FOSS_MESSAGE') }}
    </v-alert>
  </section>

  <section v-if="props.noFOSS && props.fossVersion === 'legacy'">
    <v-alert color="warning" type="warning">
      {{ t('NO_FOSS_MESSAGE') }}
    </v-alert>
  </section>

  <section v-if="(config.useFutureProduct || config.useFutureIt) && props.selectedProjectsContainEmptySbom">
    <v-alert color="warning" type="warning">
      {{ t('NO_PROJECT_NO_FOSS') }}
    </v-alert>
  </section>

  <section v-if="props.isNoSbomNoFossWarning && props.fossVersion === 'default'">
    <v-alert color="warning" type="warning">
      {{ t('NO_PROJECT_NO_FOSS') }}
    </v-alert>
  </section>

  <section v-if="props.isNoSbomNoFossWarning && props.fossVersion === 'legacy'">
    <v-alert color="warning" type="warning">
      {{ t('MISSING_NO_FOSS_FLAG_AND_SBOM') }}
    </v-alert>
  </section>
</template>

<style scoped lang="scss">
a {
  color: var(--text-color);
  display: block;

  &:hover {
    text-decoration: underline;
  }
}
</style>
