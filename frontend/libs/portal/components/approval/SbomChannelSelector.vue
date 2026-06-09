<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {useApprovalCheck} from '@disclosure-portal/composables/useApprovalCheck';
import {SpdxFile, VersionSlim} from '@disclosure-portal/model/VersionDetails';
import {formatDateAndTime} from '@disclosure-portal/utils/Table';
import config from '@shared/utils/config';
import {useI18n} from 'vue-i18n';

interface Props {
  channels: VersionSlim[];
  sboms: SpdxFile[];
  noFOSS: boolean;
  isVehicle: boolean;
  approvableSpdxKey: string;
}

const props = defineProps<Props>();
const selectedChannel = defineModel<VersionSlim | null>('selectedChannel', {default: null});
const selectedSbom = defineModel<SpdxFile | null>('selectedSbom', {default: null});

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
      :label="t('SELECT_VERSION')"
      :items="props.channels"
      :disabled="props.noFOSS"
      hide-details
      clearable
      autocomplete="off" />
    <v-autocomplete
      v-model="selectedSbom"
      :disabled="props.noFOSS"
      variant="outlined"
      item-title="name"
      :label="t('SELECT_SBOM_DELIVERY')"
      hide-details
      clearable
      autocomplete="off"
      :items="props.sboms">
      <template v-slot:item="{item, props: itemProps}">
        <v-list-item v-bind="itemProps" title="">
          <div class="d-flex">
            <div>
              <v-icon color="primary" v-if="approvableSpdxKey === item.raw._key" size="small" class="pb-1"
                >mdi-star</v-icon
              >
            </div>
            <div>
              <v-icon
                color="green"
                v-if="config.useFutureProduct && isVehicle && isAudited(selectedChannel, item?.raw?._key)"
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
          <v-icon color="primary" v-if="approvableSpdxKey === item.raw._key" size="small" class="pb-1">mdi-star</v-icon>
        </div>
        <div>
          <v-icon
            color="green"
            v-if="config.useFutureProduct && isVehicle && isAudited(selectedChannel, item?.raw?._key)"
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
