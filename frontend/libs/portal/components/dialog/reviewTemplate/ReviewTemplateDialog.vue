<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {ReviewRemarkLevel} from '@disclosure-portal/model/Quality';
import {ReviewTemplate} from '@disclosure-portal/model/ReviewTemplate';
import {DiscoForm} from '@disclosure-portal/types/discobasics';
import useRules from '@disclosure-portal/utils/Rules';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import DCloseButton from '@shared/components/disco/DCloseButton.vue';
import {computed, nextTick, Ref, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';

const props = defineProps<{
  initialData: ReviewTemplate | null;
  mode: 'create' | 'edit';
  errorMessage?: string;
  levels: ReviewRemarkLevel[];
}>();

const dialogModel = defineModel('dialog', {type: Boolean});
const emit = defineEmits<{(event: 'save', formData: ReviewTemplate): void}>();

const {t} = useI18n();
const {minMax} = useRules();

const defaultFormData = {
  title: '',
  description: '',
  level: ReviewRemarkLevel.NOT_SET,
  source: '',
} as ReviewTemplate;

const formData = ref<ReviewTemplate>(defaultFormData);

watch(
  () => dialogModel.value,
  (newVal) => {
    if (newVal && props.mode === 'create') {
      formData.value = {...defaultFormData};
    }
  },
);

watch(
  () => props.initialData,
  (newData) => {
    if (props.mode === 'edit' && newData) {
      formData.value = {...newData};
    }
  },
  {immediate: true},
);

const activeRules = ref({
  title: minMax(t('NPV_DIALOG_TF_TITLE'), 5, 80, false),
  description: minMax(t('NP_DIALOG_TF_DESCRIPTION'), 10, 700, false),
  level: [(v: string) => v !== ReviewRemarkLevel.NOT_SET || t('DLG_LEVEL_REQUIRED')],
  source: minMax(t('NPV_DIALOG_TF_SOURCE'), 3, 100, true),
});

const closeDialog = () => {
  dialogModel.value = false;
};

const requiredDescriptionLabel = computed(() => {
  return `${t('NP_DIALOG_TF_DESCRIPTION')}*`;
});

const reviewTemplateDialog: Ref<DiscoForm | null> = ref(null);
const saveForm = async () => {
  await nextTick(async () => {
    reviewTemplateDialog.value?.validate().then(async (info) => {
      if (!info.valid) {
        return;
      }
      emit('save', formData.value);
      reviewTemplateDialog.value?.reset();
    });
  });
};
</script>

<template>
  <v-dialog v-model="dialogModel" width="600" persistent>
    <v-form ref="reviewTemplateDialog">
      <v-card class="pa-8">
        <v-card-title>
          <v-row>
            <v-col cols="10" class="d-flex align-center">
              <span class="text-h5">
                {{
                  mode === 'edit' ? t('DIALOG_TITLE_EDIT_REVIEW_TEMPLATE') : t('DIALOG_TITLE_CREATE_REVIEW_TEMPLATE')
                }}
              </span>
            </v-col>
            <v-col cols="2" class="text-right px-0">
              <DCloseButton @click="closeDialog" />
            </v-col>
          </v-row>
        </v-card-title>
        <v-card-text>
          <v-row>
            <v-col cols="12" class="px-0">
              <v-text-field
                autocomplete="off"
                required
                variant="outlined"
                hide-details="auto"
                v-model="formData.title"
                :rules="activeRules.title"
                :label="t('NPV_DIALOG_TF_TITLE')"></v-text-field>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12" class="px-0 pb-0">
              <v-textarea
                variant="outlined"
                hide-details="auto"
                no-resize
                v-model="formData.description"
                :counter="700"
                :label="requiredDescriptionLabel"
                :rules="activeRules.description"></v-textarea>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12" class="px-0 pt-0">
              <v-select
                required
                variant="outlined"
                hide-details="auto"
                v-model="formData.level"
                :item-title="(item) => (item ? t('REMARK_LEVEL_' + item) : '')"
                :item-value="(item) => item"
                :items="levels"
                :rules="activeRules.level"
                :label="t('COL_LEVEL')"
                v-bind:menu-props="{location: 'bottom'}">
              </v-select>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12" class="px-0">
              <v-text-field
                autocomplete="off"
                variant="outlined"
                hide-details="auto"
                v-model="formData.source"
                :rules="activeRules.source"
                :label="t('NPV_DIALOG_TF_SOURCE')"></v-text-field>
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions class="justify-end">
          <DCActionButton
            isDialogButton
            size="small"
            variant="text"
            @click="dialogModel = false"
            class="mr-5"
            :text="t('BTN_CANCEL')" />

          <DCActionButton
            isDialogButton
            size="small"
            variant="flat"
            @click="saveForm"
            :text="mode === 'create' ? t('NP_DIALOG_BTN_CREATE') : t('NP_DIALOG_BTN_SAVE')" />
        </v-card-actions>
      </v-card>
    </v-form>
  </v-dialog>
</template>
