<script lang="ts">
import ExternalSourcePostRequest from '@disclosure-portal/model/ExternalSourcePostRequest';
import VersionService from '@disclosure-portal/services/version';
import {isURL} from '@disclosure-portal/utils/Validation';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import DCloseButton from '@shared/components/disco/DCloseButton.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import {defineComponent, nextTick, onMounted, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';

export default defineComponent({
  components: {
    DCActionButton,
    DCloseButton,
  },
  props: {
    show: {
      type: Boolean,
      default: false,
    },
  },
  setup(props, {emit}) {
    const {t} = useI18n();
    const snackbar = useSnackbar();

    // Reaktive Variablen
    const show = ref(false);
    const item = ref(new ExternalSourcePostRequest());
    const externalSourceForm = ref<HTMLFormElement | null>(null);
    const mode = ref('');
    const projectKey = ref('');
    const versionKey = ref('');
    const itemUid = ref('');
    const formError = ref('');
    const title = ref('');

    const rules = {
      required: (value: string) => !!value || t('VALIDATION_required'),
      validURL: (value: string) => (!!value && isURL(value)) || t('VALIDATION_url'),
      maxLength: (value: string) =>
        (!!value && value.length <= 2000) ||
        `${t('VALIDATION_max_length')} [${value != null ? value.length.toString() : '0'}/2000]`,
    };

    // Methoden
    const open = (projectKeyValue: string, versionKeyValue: string) => {
      mode.value = 'create';
      projectKey.value = projectKeyValue;
      versionKey.value = versionKeyValue;
      title.value = t('NP_DIALOG_CCS_TITLE');
      show.value = true;
      item.value = new ExternalSourcePostRequest();
      externalSourceForm.value?.reset();
    };

    const edit = (projectKeyValue: string, versionKeyValue: string, externalSource: any) => {
      mode.value = 'edit';
      projectKey.value = projectKeyValue;
      versionKey.value = versionKeyValue;
      title.value = t('NP_DIALOG_CCS_TITLE_EDIT');
      item.value.Comment = externalSource.comment;
      item.value.URL = externalSource.url;
      itemUid.value = externalSource._key;
      show.value = true;
    };

    const close = () => {
      show.value = false;
    };

    const convertFileUrl = () => {
      let url = item.value.URL;
      if (url.startsWith('\\\\')) {
        url = url.replace('\\\\', '/').replaceAll('\\', '/');
        url = `file://${url}`;
      } else if (url.startsWith('file://')) {
        url = url.replaceAll('\\', '/');
      }
      item.value.URL = url;
    };

    const doDialogAction = async () => {
      await nextTick(async () => {
        const validForm = (await validate())?.valid;
        if (validForm) {
          convertFileUrl();
          if (mode.value === 'create') {
            await VersionService.createExternalSource(item.value, projectKey.value, versionKey.value);
          } else if (mode.value === 'edit') {
            await VersionService.updateExternalSource(itemUid.value, item.value, projectKey.value, versionKey.value);
          }
          closeAndReload();
        }
      });
    };

    const closeAndReload = () => {
      snackbar.info(t(`DIALOG_source_code_${mode.value}_success`));
      close();
      emit('reload');
    };

    const validate = () => {
      return externalSourceForm.value?.validate() || false;
    };

    // Watchers
    watch(
      () => props.show,
      (newVal) => {
        if (!newVal) externalSourceForm.value?.reset();
      },
    );

    onMounted(() => {
      // externalSourceForm.value = document.querySelector('#externalSourceForm') as HTMLFormElement;
    });

    return {
      t,
      show,
      item,
      externalSourceForm,
      mode,
      projectKey,
      versionKey,
      itemUid,
      formError,
      title,
      rules,
      open,
      edit,
      close,
      doDialogAction,
      closeAndReload,
      validate,
    };
  },
});
</script>

<template>
  <v-form ref="externalSourceForm" id="externalSourceForm">
    <v-dialog v-model="show" scrollable width="700">
      <v-card class="pa-8 dDialog">
        <v-card-title>
          <v-row>
            <v-col cols="10">
              <span class="text-h5">{{ title }}</span>
            </v-col>
            <v-col cols="2" align="right">
              <DCloseButton @click="close" />
            </v-col>
          </v-row>
        </v-card-title>
        <v-card-text class="pt-2">
          <v-row dense>
            <v-col cols="12" xs="12" class="errorBorder">
              <v-text-field
                autocomplete="off"
                v-model="item.URL"
                variant="outlined"
                density="compact"
                class="required"
                :label="t('NES_DIALOG_URL')"
                :rules="[rules.required, rules.validURL, rules.maxLength]"
                hide-details="auto" />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12" xs="12" class="errorBorder">
              <v-textarea
                v-model="item.Comment"
                variant="outlined"
                density="compact"
                class="required"
                :label="t('NES_DIALOG_COMMENT')"
                :rules="[rules.required, rules.maxLength]"
                hide-details="auto" />
            </v-col>
          </v-row>
          <v-row v-if="formError">
            <v-col cols="12">
              <span class="text-[rgb(var(--v-theme-error))]">{{ formError }}</span>
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions class="justify-end">
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
            :text="t('NP_DIALOG_BTN_EDIT')" />
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-form>
</template>
