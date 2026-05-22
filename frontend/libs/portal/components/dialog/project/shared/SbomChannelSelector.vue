<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {SpdxFile, VersionSlim} from '@disclosure-portal/model/VersionDetails';
import {formatDateAndTime} from '@disclosure-portal/utils/Table';
import {useI18n} from 'vue-i18n';

interface Props {
  channels: VersionSlim[];
  sboms: SpdxFile[];
  noFOSS: boolean;
  approvableSpdxKey: string;
}

const props = defineProps<Props>();
const selectedChannel = defineModel<VersionSlim | null>('selectedChannel', {required: true});
const selectedSbom = defineModel<SpdxFile | null>('selectedSbom', {required: true});

const {t} = useI18n();
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
      autocomplete="off" />
    <v-autocomplete
      v-model="selectedSbom"
      :disabled="props.noFOSS"
      variant="outlined"
      item-title="name"
      :label="t('SELECT_SBOM_DELIVERY')"
      hide-details
      autocomplete="off"
      :items="props.sboms">
      <template v-slot:item="{item, props: itemProps}">
        <v-list-item v-bind="itemProps" title="">
          <div class="d-flex">
            <div>
              <v-icon color="primary" v-if="approvableSpdxKey == item.raw._key" size="small" class="pb-1"
                >mdi-star</v-icon
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
        <span class="d-subtitle-2 ml-5">{{ formatDateAndTime(item.raw.uploaded) }}&nbsp;</span>
        <span class="d-text d-secondary-text">&nbsp;-&nbsp;{{ item.raw.metaInfo.name }}</span>
        <span class="d-text d-secondary-text" v-if="item.raw.tag">&nbsp;({{ item.raw.tag }})</span>
        <span class="d-text d-secondary-text" v-if="item.raw.isRecent">&nbsp;{{ '[' + t('SBOM_LATEST') + ']' }}</span>
        <span class="d-text d-secondary-text" v-else>&nbsp;{{ '[' + t('SBOM_FORMER') + ']' }}</span>
      </template>
    </v-autocomplete>
  </Stack>
</template>
