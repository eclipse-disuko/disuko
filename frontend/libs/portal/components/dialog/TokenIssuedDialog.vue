<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {Token} from '@disclosure-portal/model/Project';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {formatDate} from '@disclosure-portal/utils/View';
import useSnackbar from '@shared/composables/useSnackbar';
import type {DialogLayoutConfig} from '@shared/layouts/DialogLayout.vue';
import config from '@shared/utils/config';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const {info} = useSnackbar();

const showDialog = defineModel<boolean>('showDialog');
const props = defineProps<{
  token: Token;
  renewed: boolean;
}>();

const endpoint = ref(config.PUBLIC_API_ENDPOINT);
const publicUrl = ref(config.PUBLIC_UI_URL);
const snackBarText = ref<string | undefined>('');
const isProd = computed(() => config.isProd);

const projectName = computed(() => useProjectStore().currentProject?.name);
const tokenExpiry = computed(() => formatDate(props.token!.expiry));
const warningHintText = computed(() => t('TOKEN_INFO_WARN_HINT'));

const close = () => {
  showDialog.value = false;
};

function getInfoToClipboard() {
  const publicUrlLine = isProd.value ? '' : `${t('TOKEN_COPY_TEXT_PUBLIC_URL')}: ${publicUrl.value}\n`;
  return `${t('TOKEN_COPY_TEXT_PROJECT_NAME')}: ${projectName.value}
${t('TOKEN_COPY_TEXT_PROJECT_TOKEN')}: ${props.token.tokenSecret}
${t('TOKEN_COPY_TEXT_TOKEN_EXPIRY')}: ${props.token.expiry}
${t('TOKEN_COPY_TEXT_TOKEN_NAME')}: ${props.token.company}
${t('TOKEN_COPY_TEXT_TOKEN_DESCRIPTION')}: ${props.token.description}
${t('TOKEN_COPY_TEXT_ENDPOINT')}: ${endpoint.value}
${publicUrlLine}

${warningHintText.value}`;
}

function copyTokenToClipboard() {
  navigator.clipboard
    .writeText(props.token.tokenSecret)
    .then(() => {
      snackBarText.value = t('SNACK_copied_to_clipboard');
      info(snackBarText.value);
    })
    .catch(() => {
      snackBarText.value = 'Something went wrong';
    });
}

function copyInfoToClipboard() {
  navigator.clipboard
    .writeText(getInfoToClipboard())
    .then(() => {
      snackBarText.value = t('SNACK_copied_to_clipboard');
      info(snackBarText.value);
    })
    .catch(() => {
      snackBarText.value = 'Something went wrong';
    });
}
const dialogConfig: DialogLayoutConfig = {
  title: t(props.renewed ? 'DLG_RENEWAL_TITLE' : 'UM_DIALOG_TITLE_TOKEN_ISSUED'),
  secondaryButton: {text: t('BTN_COPY_TOKEN_CONTENT')},
  primaryButton: {text: t('BTN_CLOSE')},
};
</script>

<template>
  <v-dialog v-model="showDialog" content-class="small" persistent width="900">
    <DialogLayout :config="dialogConfig" @close="close" @secondary-action="copyInfoToClipboard" @primary-action="close">
      <Stack>
        <v-text-field
          autocomplete="off"
          readonly
          variant="outlined"
          :label="t('PROJECT')"
          :model-value="projectName"
          hide-details="auto">
        </v-text-field>

        <v-text-field
          autocomplete="off"
          readonly
          :label="t('NP_DIALOG_HINT_TOKEN')"
          :model-value="token.tokenSecret"
          hide-details="auto"
          variant="outlined">
          <template v-slot:append-inner>
            <v-icon v-bind="props" color="primary" @click="copyTokenToClipboard">mdi-content-copy</v-icon>
            <Tooltip>{{ t('TT_COPY_TO_CLIPBOARD') }}</Tooltip>
          </template>
        </v-text-field>

        <v-text-field
          autocomplete="off"
          readonly
          variant="outlined"
          :label="t('COL_EXPIRY')"
          :model-value="tokenExpiry"
          hide-details="auto">
        </v-text-field>

        <v-text-field
          autocomplete="off"
          readonly
          variant="outlined"
          :label="t('TOKEN_NAME')"
          :model-value="token.company"
          hide-details="auto">
        </v-text-field>

        <v-text-field
          autocomplete="off"
          readonly
          variant="outlined"
          :label="t('DESCRIPTION')"
          :model-value="token.description"
          hide-details="auto">
        </v-text-field>

        <v-text-field
          autocomplete="off"
          readonly
          variant="outlined"
          :label="t('ENDPOINT')"
          :model-value="endpoint"
          hide-details="auto">
        </v-text-field>

        <v-text-field
          v-if="!isProd"
          autocomplete="off"
          readonly
          variant="outlined"
          :label="t('TOKEN_COPY_TEXT_PUBLIC_URL')"
          :model-value="publicUrl"
          hide-details="auto">
        </v-text-field>

        <v-alert color="warning" type="warning">
          <span>
            {{ warningHintText }}
          </span>
        </v-alert>
      </Stack>
    </DialogLayout>
  </v-dialog>
</template>
