<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {useView} from '@disclosure-portal/composables/useView';
import {IDefaultSelectItem, IObligation} from '@disclosure-portal/model/IObligation';
import AdminService from '@disclosure-portal/services/admin';
import useRules from '@disclosure-portal/utils/Rules';
import {getClassificationLevels, getClassificationTypes} from '@disclosure-portal/utils/View';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import DCloseButton from '@shared/components/disco/DCloseButton.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import {nextTick, reactive, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const {info: infoSnackbar} = useSnackbar();
const emit = defineEmits(['reload']);
const {getTextOfLevel, getTextOfType} = useView();

const show = ref(false);
const mode = ref<'create' | 'edit'>('create');
const title = ref('');
const obligationTypes = ref<IDefaultSelectItem[]>([]);
const obligationWarnLevels = ref<IDefaultSelectItem[]>([]);

const item = reactive<IObligation>({
  _key: '',
  autoApproved: false,
  created: new Date().toISOString(),
  updated: new Date().toISOString(),
  name: '',
  nameDe: '',
  type: '',
  warnLevel: '',
  description: '',
  descriptionDe: '',
});

const {minMax, longText} = useRules();

const activeRules = ref({
  name: minMax(t('NPV_DIALOG_TF_NAME'), 5, 80, false),
  required: [(value: string) => !!value || t('VALIDATION_required')],
  description: longText(t('AL_DIALOG_TF_DESCRIPTION_EN')),
  descriptionDe: longText(t('AL_DIALOG_TF_DESCRIPTION_DE')),
});

const classificationDialog = ref();

const initializeDropdowns = () => {
  obligationTypes.value = getClassificationTypes().map((type) => ({
    value: type,
    text: getTextOfType(type),
  }));

  obligationWarnLevels.value = getClassificationLevels().map((level) => ({
    value: level,
    text: getTextOfLevel(level),
  }));
};

const open = (model?: IObligation) => {
  if (model) {
    mode.value = 'edit';
    title.value = t('DIALOG_TITLE_EDIT_CLASSIFICATION');
    Object.assign(item, model);
  } else {
    mode.value = 'create';
    title.value = t('DIALOG_TITLE_CREATE_CLASSIFICATION');
    resetForm();
  }
  show.value = true;
};

const resetForm = () => {
  Object.assign(item, {
    name: '',
    nameDe: '',
    type: '',
    warnLevel: '',
    description: '',
    descriptionDe: '',
  });
  if (classificationDialog.value) classificationDialog.value.reset();
};

const close = () => {
  show.value = false;
  resetForm();
};

const doDialogAction = async () => {
  await nextTick(async () => {
    const validationResult = await classificationDialog.value?.validate();
    if (!validationResult?.valid) {
      return;
    }

    if (mode.value === 'create') {
      await AdminService.postObligation(item);
      infoSnackbar(t('DIALOG_classification_create_success'));
    } else if (mode.value === 'edit') {
      await AdminService.putObligation(item);
      infoSnackbar(t('DIALOG_classification_edit_success'));
    }
    close();
    emit('reload');
  });
};

watch(show, (value) => {
  if (!value) {
    resetForm();
  }
});

initializeDropdowns();

defineExpose({
  open,
});
</script>

<template>
  <v-dialog v-model="show" width="600" persistent scrollable>
    <v-form ref="classificationDialog">
      <v-card class="pa-8 dDialog">
        <v-card-title>
          <v-row>
            <v-col cols="10" class="d-flex align-center">
              <span class="text-h5">{{ title }}</span>
            </v-col>
            <v-col cols="2" class="text-right px-0">
              <DCloseButton @click="close" />
            </v-col>
          </v-row>
        </v-card-title>
        <v-card-text>
          <v-row>
            <v-col cols="12" xs="12" class="errorBorder">
              <v-text-field
                autocomplete="off"
                variant="outlined"
                class="required"
                hide-details="auto"
                v-model="item.name"
                :rules="activeRules.name"
                :label="t('AL_DIALOG_TF_NAME_EN')"
                clearable />
            </v-col>
            <v-col cols="12" xs="12" class="errorBorder">
              <v-text-field
                autocomplete="off"
                variant="outlined"
                clearable
                class="required"
                v-model="item.nameDe"
                :rules="activeRules.name"
                :label="t('AL_DIALOG_TF_NAME_DE')"
                hide-details="auto" />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12" xs="12" class="errorBorder">
              <v-select
                v-model="item.type"
                variant="outlined"
                class="required"
                clearable
                item-title="text"
                :rules="activeRules.required"
                :items="obligationTypes"
                :label="t('LABEL_OBLIGATION_TYPE')"
                hide-details="auto"></v-select>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12" xs="12" class="errorBorder">
              <v-select
                v-model="item.warnLevel"
                variant="outlined"
                class="required"
                clearable
                item-title="text"
                :rules="activeRules.required"
                :items="obligationWarnLevels"
                :label="t('LABEL_OBLIGATION_WARN_LEVEL')"
                hide-details="auto"></v-select>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12" xs="12">
              <v-textarea
                variant="outlined"
                clearable
                v-model="item.description"
                :rules="activeRules.description"
                :label="t('AL_DIALOG_TF_DESCRIPTION_EN')"
                counter="1000"
                no-resize
                hide-details="auto" />
            </v-col>
            <v-col cols="12" xs="12">
              <v-textarea
                variant="outlined"
                clearable
                v-model="item.descriptionDe"
                :rules="activeRules.descriptionDe"
                :label="t('AL_DIALOG_TF_DESCRIPTION_DE')"
                counter="1000"
                no-resize
                hide-details="auto" />
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <DCActionButton
            isDialogButton
            size="small"
            variant="text"
            @click="close"
            class="mr-5"
            :text="t('BTN_CANCEL')" />

          <DCActionButton
            isDialogButton
            size="small"
            variant="flat"
            @click="doDialogAction"
            :text="mode === 'create' ? t('NP_DIALOG_BTN_CREATE') : t('NP_DIALOG_BTN_SAVE')" />
        </v-card-actions>
      </v-card>
    </v-form>
  </v-dialog>
</template>
