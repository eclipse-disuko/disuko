<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {DocumentMeta} from '@disclosure-portal/model/ApprovalRequest';
import {useI18n} from 'vue-i18n';

const model = defineModel<DocumentMeta>({required: true});

const {t} = useI18n();
</script>

<template>
  <Stack class="gap-1">
    <Stack v-if="model.c6" direction="row" align="center">
      <v-icon size="small" color="warning">mdi-alert</v-icon>
      <span class="d-block">{{ t('NO_FOSS_WARNING') }}</span>
    </Stack>
    <Stack direction="row" align="center">
      <v-icon size="small" color="warning">mdi-alert</v-icon>
      <span class="d-block">{{ t('NO_FOSS_DISABLED_TOOLTIP') }}</span>
    </Stack>
    <v-switch
      class="ml-2"
      v-model="model.c6"
      color="primary"
      density="compact"
      :label="t('NO_FOSS_MARKER')"
      hide-details
      disabled></v-switch>
  </Stack>

  <ExpansionPanel :title="t('SBOM_APPROVAL_ATTRIBUTES')">
    <template #body>
      <v-checkbox v-model="model.c1" :readonly="model.c6" :label="t('SBOM_APPROVAL_CHECK1')" hide-details />
      <v-checkbox v-model="model.c2" :readonly="model.c6" :label="t('SBOM_APPROVAL_CHECK2')" hide-details />
      <v-checkbox v-model="model.c3" :readonly="model.c6" :label="t('SBOM_APPROVAL_CHECK3')" hide-details />
      <v-checkbox v-model="model.c4" :readonly="model.c6" :label="t('SBOM_APPROVAL_CHECK4')" hide-details />
      <v-checkbox v-model="model.c5" :readonly="model.c6" :label="t('SBOM_APPROVAL_CHECK5')" hide-details />
      <v-checkbox v-model="model.c6" :label="t('SBOM_APPROVAL_CHECK6')" disabled></v-checkbox>
    </template>
  </ExpansionPanel>
</template>
