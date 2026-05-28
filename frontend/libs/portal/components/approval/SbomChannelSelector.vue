<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {useApprovalCheck} from '@disclosure-portal/composables/useApprovalCheck';
import {SpdxFile, VersionSlim} from '@disclosure-portal/model/VersionDetails';
import {formatDateAndTime} from '@disclosure-portal/utils/Table';
import {useI18n} from 'vue-i18n';

defineProps<{
  channels: VersionSlim[];
  sboms: SpdxFile[];
  noFOSS: boolean;
  isVehicle?: boolean;
  approvableSpdxKey: string;
}>();

const selectedChannel = defineModel<VersionSlim | null | undefined>('selectedChannel');
const selectedSbom = defineModel<SpdxFile | null>('selectedSbom');

const {t} = useI18n();
const {isAudited} = useApprovalCheck();
</script>

<template>
  <Stack>
    <v-select
      v-model="selectedChannel"
      variant="outlined"
      item-title="name"
      return-object
      autocomplete="off"
      :label="t('SELECT_VERSION')"
      :items="channels"
      :disabled="noFOSS"
      hide-details />
    <v-autocomplete
      v-model="selectedSbom"
      variant="outlined"
      item-title="name"
      autocomplete="off"
      :label="t('SELECT_SBOM_DELIVERY')"
      :items="sboms"
      :disabled="noFOSS"
      hide-details>
      <template v-slot:item="{item, props}">
        <v-list-item v-bind="props" title="">
          <div class="d-flex">
            <v-icon color="primary" v-if="approvableSpdxKey == item.raw._key" size="small" class="pb-1">
              mdi-star
            </v-icon>
            <div>
              <v-icon
                color="green"
                v-if="isVehicle && isAudited(selectedChannel ?? null, item?.raw?._key)"
                size="small"
                class="ml-1 pb-1"
                >mdi-clipboard-check-outline</v-icon
              >
            </div>
            <span class="d-subtitle-2 ml-5">{{ formatDateAndTime(item.raw.uploaded) }}&nbsp;</span>
            <span class="d-text d-secondary-text">&nbsp;-&nbsp;{{ item.raw.metaInfo.name }}</span>
            <span class="d-text d-secondary-text" v-if="item.raw.tag">&nbsp;({{ item.raw.tag }})</span>
            <span class="d-text d-secondary-text" v-if="item.raw.isRecent"
              >&nbsp;{{ '[' + t('SBOM_LATEST') + ']' }}</span
            >
            <span class="d-text d-secondary-text" v-else>&nbsp;{{ '[' + t('SBOM_FORMER') + ']' }}</span>
          </div>
        </v-list-item>
      </template>
      <template v-slot:selection="{item}">
        <div style="min-width: 13px">
          <v-icon color="primary" v-if="approvableSpdxKey == item.raw._key" size="small" class="pb-1">mdi-star</v-icon>
        </div>
        <div>
          <v-icon
            color="green"
            v-if="isVehicle && isAudited(selectedChannel ?? null, item?.raw?._key)"
            size="small"
            class="ml-1 pb-1"
            >mdi-clipboard-check-outline</v-icon
          >
        </div>
        <span class="d-subtitle-2 ml-5">{{ formatDateAndTime(item.raw.uploaded) }}&nbsp;</span>
        <span class="d-text d-secondary-text">&nbsp;-&nbsp;{{ item.raw.metaInfo.name }}</span>
        <span class="d-text d-secondary-text" v-if="item.raw.tag">&nbsp;({{ item.raw.tag }})</span>
        <span class="d-text d-secondary-text" v-if="item.raw.isRecent">&nbsp;{{ '[' + t('SBOM_LATEST') + ']' }}</span>
        <span class="d-text d-secondary-text" v-else>&nbsp;{{ '[' + t('SBOM_FORMER') + ']' }}</span>
      </template>
    </v-autocomplete>
  </Stack>
</template>
