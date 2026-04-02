<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<template>
  <v-form ref="form">
    <v-card class="pa-8">
      <v-card-title>
        <v-row>
          <v-col cols="10" class="d-flex align-center">
            <span class="text-h5">
              {{ title }}
            </span>
          </v-col>
          <v-col cols="2" class="text-right px-0">
            <DCloseButton @click="$emit('close')" />
          </v-col>
        </v-row>
      </v-card-title>
      <v-card-text class="pt-2">
        <v-row>
          <v-col cols="12" xs="12">
            <v-switch
              v-model="localItem.spdx"
              hide-details
              color="primary"
              :label="t('COL_DELIVERY')"
              class="shrink mr-2 mt-0" />
            <div class="d-text ml-12">{{ t('SBOM_SUBSCRIPTION_TEXT') }}</div>
          </v-col>
        </v-row>
        <v-row>
          <v-col cols="12" xs="12">
            <v-switch
              v-model="localItem.overallReview"
              hide-details
              color="primary"
              :label="t('HEADLINE_OVERALL_REVIEW')"
              class="shrink mr-2 mt-0" />
            <div class="d-text ml-12">{{ t('OVERALL_REVIEW_SUBSCRIPTION_TEXT') }}</div>
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-actions class="justify-end">
        <DCActionButton
          isDialogButton
          size="small"
          variant="text"
          @click="emit('close')"
          class="mr-5"
          :text="t('BTN_CANCEL')" />

        <DCActionButton isDialogButton size="small" variant="flat" @click="confirm" :text="confirmText" />
      </v-card-actions>
    </v-card>
  </v-form>
</template>

<script lang="ts" setup>
import {ProjectSubscriptions} from '@disclosure-portal/model/Project';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import DCloseButton from '@shared/components/disco/DCloseButton.vue';
import _ from 'lodash';
import {ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';

const props = defineProps<{
  title: string;
  confirmText: string;
  item: ProjectSubscriptions | null;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'confirm', value: ProjectSubscriptions): void;
}>();

const {t} = useI18n();

const localItem = ref<ProjectSubscriptions>({
  spdx: false,
  overallReview: false,
});

watch(
  () => props.item,
  (newItem) => {
    if (newItem) {
      localItem.value = _.cloneDeep(newItem);
    }
  },
  {immediate: true},
);

function confirm() {
  emit('confirm', localItem.value);
}
</script>
