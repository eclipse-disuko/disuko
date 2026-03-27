// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import Label from '@disclosure-portal/model/Label';
import AdminService from '@disclosure-portal/services/admin';
import {DiscoForm} from '@disclosure-portal/types/discobasics';
import useRules from '@disclosure-portal/utils/Rules';
import DCloseButton from '@shared/components/disco/DCloseButton.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import {computed, nextTick, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const props = defineProps<{
  initialData?: Label;
  mode: 'create' | 'edit';
  errorMessage?: string;
  isOpen: boolean;
  type?: 'SCHEMA' | 'POLICY' | 'PROJECT';
}>();

const emit = defineEmits(['update:isOpen', 'reload']);

const isDialogOpen = computed({
  get() {
    return props.isOpen;
  },
  set(value) {
    emit('update:isOpen', value);
  },
});

const defaultLabel = {
  name: '',
  description: '',
  type: props.type || 'SCHEMA',
} as Label;

const item = ref<Label>(defaultLabel);

const activeRules = ref({
  required: useRules().minMax('Name', 3, 80, false),
  description: useRules().longText('Description'),
});

const closeDialog = () => {
  isDialogOpen.value = false;
};

const labelNameLowercase = () => {
  item.value.name = item.value.name.toLowerCase();
};

const labelDialog = ref<DiscoForm | null>(null);
const doDialogAction = async () => {
  await nextTick(async () => {
    labelDialog.value?.validate().then(async (info) => {
      if (!info.valid) {
        return;
      }
      item.value.name = item.value.name.trim();
      if (props.mode === 'create') {
        await AdminService.createLabel(item.value);
      } else if (props.mode === 'edit') {
        await AdminService.editLabel(item.value);
      }
      closeAndReload();
    });
  });
};

const snack = useSnackbar();
const closeAndReload = () => {
  labelDialog.value?.reset();
  snack.info(t(`DIALOG_${props.type}_LABEL_${props.mode}_SUCCESS`));
  emit('reload');
  closeDialog();
};

watch(
  () => props.initialData,
  (newData) => {
    if (props.mode === 'edit' && newData) {
      item.value = {...newData};
    }
  },
  {immediate: true},
);

watch(
  () => props.type,
  (newType) => {
    item.value.type = newType || 'SCHEMA';
  },
  {immediate: true},
);
</script>

<template>
  <v-dialog v-model="isDialogOpen" max-width="600px" persistent>
    <v-form ref="labelDialog">
      <v-card class="pa-8">
        <v-card-title>
          <v-row>
            <v-col cols="10" class="d-flex align-center">
              <span class="text-h5">
                <template v-if="item.type === 'SCHEMA'">
                  {{ mode === 'edit' ? t('LM_DIALOG_TITLE_EDIT_SCHEMA') : t('LM_DIALOG_TITLE_SCHEMA') }}
                </template>
                <template v-else-if="item.type === 'PROJECT'">
                  {{ mode === 'edit' ? t('LM_DIALOG_TITLE_EDIT_PROJECT') : t('LM_DIALOG_TITLE_PROJECT') }}
                </template>
                <template v-else>
                  {{ mode === 'edit' ? t('LM_DIALOG_TITLE_EDIT_POLICY') : t('LM_DIALOG_TITLE_POLICY') }}
                </template>
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
                variant="outlined"
                v-model="item.name"
                class="required"
                :rules="activeRules.required"
                :label="t('AL_DIALOG_TF_NAME')"
                @keyup="labelNameLowercase"
                hide-details="auto"
                autofocus />
            </v-col>
          </v-row>
          <v-row>
            <v-col class="px-0">
              <v-textarea
                no-resize
                variant="outlined"
                counter="1000"
                :rules="activeRules.description"
                v-model="item.description"
                :label="t('AL_DIALOG_TF_DESCRIPTION')"
                hide-details="auto" />
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions class="justify-end">
          <v-btn isDialogButton size="small" variant="text" @click="closeDialog" color="primary" class="mr-5">
            {{ t('BTN_CANCEL') }}
          </v-btn>
          <v-btn isDialogButton size="small" variant="flat" @click="doDialogAction" color="primary" class="mr-1">
            <span v-if="props.mode === 'create'">{{ t('NP_DIALOG_BTN_CREATE') }}</span>
            <span v-else>{{ t('NP_DIALOG_BTN_SAVE') }}</span>
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-form>
  </v-dialog>
</template>
