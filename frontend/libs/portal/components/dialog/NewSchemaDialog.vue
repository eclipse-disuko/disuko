// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import ErrorDialog from '@disclosure-portal/components/dialog/ErrorDialog.vue';
import Label from '@disclosure-portal/model/Label';
import SchemaPostRequest from '@disclosure-portal/model/SchemaPostRequest';
import {DiscoForm} from '@disclosure-portal/types/discobasics';
import useRules from '@disclosure-portal/utils/Rules';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import DCloseButton from '@shared/components/disco/DCloseButton.vue';
import DLabel from '@shared/components/disco/DLabel.vue';
import DiscoFileUpload from '@shared/components/widgets/DiscoFileUpload.vue';
import config from '@shared/utils/config';
import {computed, nextTick, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const props = defineProps<{
  isOpen: boolean;
  labels: Label[];
}>();

const emit = defineEmits(['update:isOpen', 'onSave']);

const isDialogOpen = computed({
  get() {
    return props.isOpen;
  },
  set(value) {
    emit('update:isOpen', value);
  },
});

const item = ref<SchemaPostRequest>(new SchemaPostRequest());
const files = ref<File[]>([]);
const uploadURL = ref(config.SERVER_URL + '/api/v1/admin/schemas');
const sbomFileError = ref('');
const predefinedLabels = ref<Label[]>(props.labels);
const schemaFormValid = ref<boolean>(false);
const upload = ref<InstanceType<typeof DiscoFileUpload> | null>(null);

const activeRules = ref({
  required: useRules().minMax('Name', 3, 80, false),
  description: useRules().longText('Description'),
});

const schemaForm = ref<DiscoForm | null>(null);

const closeDialog = () => {
  isDialogOpen.value = false;
};

const doDialogAction = async () => {
  await nextTick(async () => {
    schemaForm.value?.validate().then(async (info) => {
      if (!info.valid) {
        return;
      }
      const formData = new FormData();
      formData.append('schema', JSON.stringify(item.value));

      for (const file of files.value) {
        if (file instanceof Blob && file.name) {
          formData.append('file', file, file.name);
        }
      }

      emit('onSave', formData);
      closeDialog();
    });
  });
};

watch(
  () => props.labels,
  (newLabels) => {
    predefinedLabels.value = newLabels;
  },
  {immediate: true},
);
</script>
<template>
  <v-dialog v-model="isDialogOpen" max-width="600px" persistent>
    <v-form ref="schemaForm" v-model="schemaFormValid">
      <v-card flat class="pa-8 dDialog">
        <v-card-title>
          <v-row>
            <v-col cols="10">
              <span class="text-h5">
                {{ t('SCHEMA_DIALOG_TITLE') }}
              </span>
            </v-col>
            <v-col cols="2" class="text-right px-0">
              <DCloseButton @click="closeDialog" />
            </v-col>
          </v-row>
        </v-card-title>
        <v-card-text>
          <v-row align="center" class="justify-center">
            <v-col cols="12" xs="12" class="errorBorder">
              <v-text-field
                autocomplete="off"
                required
                hide-details="auto"
                variant="outlined"
                density="compact"
                v-model="item.name"
                :rules="activeRules.required"
                :label="t('NP_DIALOG_TF_SCHEMA_NAME')"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12" xs="12" class="errorBorder">
              <v-text-field
                autocomplete="off"
                required
                :rules="activeRules.required"
                variant="outlined"
                density="compact"
                v-model="item.version"
                :label="t('NP_DIALOG_TF_VERSION')"
                hide-details="auto"
              />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12" xs="12">
              <v-select
                required
                variant="outlined"
                density="compact"
                :rules="activeRules.required"
                :items="predefinedLabels"
                v-model="item.label"
                :label="t('NP_DIALOG_SB_SCHEMALABELS')"
                v-bind:menu-props="{location: 'bottom'}"
                hide-details="auto"
                item-title="name"
                item-value="_key"
                clearable
              >
                <template v-slot:selection="{item}">
                  <DLabel v-if="item.raw.name" :labelName="item.raw.name" closable />
                </template>
                <template v-slot:label="{label}">
                  {{ label }}
                </template>
              </v-select>
            </v-col>
          </v-row>
          <v-row>
            <v-col>
              <div class="d-flex align-center border-md border-dashed border-opacity-25 px-3 py-3">
                <DiscoFileUpload
                  ref="upload"
                  :uploadTargetUrl="uploadURL"
                  acceptTypes=".json"
                  :directUpload="false"
                  @filesChanged="(f) => (files = f)"
                />
                <template v-if="files.length > 0">
                  <v-row>
                    <v-col class="d-flex align-center">
                      <h4>{{ t('selectedFile') }}{{ files.length > 0 ? files[0].name : '' }}</h4>
                      <v-icon class="pl-2" @click="() => upload!.clearFiles()">cancel</v-icon>
                    </v-col>
                  </v-row>
                </template>
                <template v-else>
                  <v-row>
                    <v-col class="d-flex align-center">
                      <span class="pr-1 opacity-50">{{ t('SCHEMA_DIALOG_DRAGDROP') }}</span>
                      <span>
                        <v-icon
                          variant="text"
                          size="small"
                          plain
                          color="primary"
                          icon="mdi-file-document-outline"
                        ></v-icon>
                      </span>
                      <span class="pl-1 pr-2 opacity-50">or</span>
                      <DCActionButton
                        icon="mdi-upload"
                        :text="t('BTN_UPLOAD')"
                        :hint="t('Btn_uploadSchema')"
                        @click="() => upload!.uploadClick()"
                        :disabled="!schemaFormValid"
                      />
                    </v-col>
                  </v-row>
                  <v-row v-if="sbomFileError">
                    <v-col>
                      <span class="error-text">{{ t(sbomFileError) }}</span>
                    </v-col>
                  </v-row>
                </template>
              </div>
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12" xs="12">
              <v-textarea
                required
                counter
                variant="outlined"
                density="compact"
                v-model="item.description"
                :label="t('NP_DIALOG_TF_DESCRIPTION')"
                no-resize
                hide-details="auto"
                :rules="activeRules.description"
              />
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions class="justify-end">
          <DCActionButton isDialogButton size="small" variant="text" @click="closeDialog" class="mr-5" :text="t('BTN_CANCEL')" />
          <DCActionButton isDialogButton size="small" variant="flat" @click="doDialogAction" :text="t('NP_DIALOG_BTN_CREATE')" />
        </v-card-actions>
      </v-card>
    </v-form>
    <ErrorDialog ref="errorDialog"></ErrorDialog>
  </v-dialog>
</template>

<style scoped>
.error-text {
  color: rgb(var(--v-theme-error));
}
</style>
