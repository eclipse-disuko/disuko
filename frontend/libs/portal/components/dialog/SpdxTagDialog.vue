<script lang="ts">
import {Tags} from '@disclosure-portal/constants/ruleValidations';
import projectService from '@disclosure-portal/services/projects';
import DCloseButton from '@shared/components/disco/DCloseButton.vue';
import {computed, defineComponent, ref} from 'vue';
import {useI18n} from 'vue-i18n';

export default defineComponent({
  name: 'SpdxTagDialog',
  components: {DCloseButton},
  props: {
    reload: {
      type: Function,
      required: true,
    },
  },
  setup(_, {emit}) {
    const {t} = useI18n();
    const show = ref(false);
    const projectUuid = ref('');
    const versionUuid = ref('');
    const spdxUuid = ref('');
    const spdxName = ref('');
    const tag = ref('');

    const rules = {
      tag: [
        (v: string) => !!v || `${Tags.TAG_MIN_LENGTH}-${Tags.TAG_MAX_LENGTH} characters required.`,
        (v: string) => v.length >= Tags.TAG_MIN_LENGTH || `Minimum ${Tags.TAG_MIN_LENGTH} characters required.`,
        (v: string) => v.length <= Tags.TAG_MAX_LENGTH || `Maximum ${Tags.TAG_MAX_LENGTH} characters allowed.`,
      ],
    };

    const isTagValid = computed(
      () => tag.value && tag.value.length >= Tags.TAG_MIN_LENGTH && tag.value.length <= Tags.TAG_MAX_LENGTH,
    );

    const setOrUpdate = (pUuid: string, vUuid: string, sUuid: string, sName: string, sTag: string) => {
      projectUuid.value = pUuid;
      versionUuid.value = vUuid;
      spdxUuid.value = sUuid;
      spdxName.value = sName;
      tag.value = sTag;
      show.value = true;
    };

    const close = () => {
      show.value = false;
    };

    const doDialogAction = async () => {
      if (isTagValid.value) {
        await projectService.updateSpdxTag(projectUuid.value, versionUuid.value, spdxUuid.value, tag.value);
        closeAndReload();
      }
    };

    const closeAndReload = () => {
      emit('reload');
      show.value = false;
    };

    return {
      t,
      show,
      projectUuid,
      versionUuid,
      spdxUuid,
      spdxName,
      tag,
      rules,
      isTagValid,
      setOrUpdate,
      close,
      doDialogAction,
    };
  },
});
</script>
<template>
  <v-dialog v-model="show" content-class="msmall" scrollable>
    <v-card class="pa-8 dDialog" flat>
      <v-card-title>
        <v-row>
          <v-col cols="10">
            <span class="text-h5">{{ t('SPDX_TAG_TITLE') + spdxName }}</span>
          </v-col>
          <v-col cols="2" align="right">
            <DCloseButton @click="close" />
          </v-col>
        </v-row>
      </v-card-title>
      <v-card-text>
        <v-row dense>
          <v-col cols="12" xs="12" class="errorBorder">
            <v-text-field
              autocomplete="off"
              :rules="rules.tag"
              variant="outlined"
              v-model="tag"
              :label="t('COL_SBOM_TAG')"
              autofocus
            />
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn @click="close" plain class="secondary mr-8">
          {{ t('BTN_CLOSE') }}
        </v-btn>
        <v-btn @click="doDialogAction" :disabled="!isTagValid" color="primary">
          <span>{{ t('NP_DIALOG_BTN_EDIT') }}</span>
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
