// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script lang="ts">
import {Tags} from '@disclosure-portal/constants/ruleValidations';
import {VersionSlim} from '@disclosure-portal/model/VersionDetails';
import projectService from '@disclosure-portal/services/projects';
import versionService from '@disclosure-portal/services/version';
import {useAppStore} from '@disclosure-portal/stores/app';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import useRules from '@disclosure-portal/utils/Rules';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import DCloseButton from '@shared/components/disco/DCloseButton.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import {computed, defineComponent, nextTick, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {VForm} from 'vuetify/components';

export default defineComponent({
  name: 'DSpdxTagDialog',
  components: {DCActionButton, DCloseButton},
  props: {
    presetTag: {
      type: String,
      required: false,
    },
    versionID: {
      type: String,
      required: true,
    },
    spdxID: {
      type: String,
      required: true,
    },
    spdxName: {
      type: String,
      required: true,
    },
    channelView: {
      type: Boolean,
      required: false,
      default: false,
    },
  },
  setup(props) {
    const {t} = useI18n();
    const isVisible = ref(false);
    const tag = ref('');
    const appStore = useAppStore();
    const projectStore = useProjectStore();
    const dialog = ref<VForm | null>(null);
    const {info: snack} = useSnackbar();
    const showDialog = () => {
      if (props.presetTag) {
        tag.value = props.presetTag;
      }
      isVisible.value = true;
    };
    const reset = () => {
      if (props.presetTag) {
        tag.value = props.presetTag;
      } else {
        dialog.value?.reset();
      }
    };
    const close = () => {
      dialog.value?.reset();
      isVisible.value = false;
    };
    const doDialogAction = async () => {
      await nextTick(async () => {
        dialog.value?.validate().then(async (info) => {
          if (!info.valid) {
            return;
          }
          await projectService.updateSpdxTag(projectModel.value._key, props.versionID, props.spdxID, tag.value);
          snack(t('DIALOG_SPDX_TAG_UPDATE_SUCCESS'));
          if (props.channelView) {
            const spdxFileHistory = (
              await versionService.getSbomHistory(projectModel.value._key, versionDetails.value._key)
            ).data;
            if (spdxFileHistory[0]) {
              spdxFileHistory[0].isRecent = true;
            }
            appStore.setChannelSpdxs(spdxFileHistory);
            await appStore.fetchAllSBOMs();
          } else {
            await appStore.fetchAllSBOMsFlat();
          }
          dialog.value?.reset();
          isVisible.value = false;
        });
      });
    };
    const projectModel = computed(() => projectStore.currentProject!);
    const versionDetails = computed((): VersionSlim => appStore.getCurrentVersion);

    const activeRules = ref({
      tag: useRules().minMax(t('COL_SBOM_TAG'), Tags.TAG_MIN_LENGTH, Tags.TAG_MAX_LENGTH, false),
    });

    return {
      isVisible,
      showDialog,
      reset,
      doDialogAction,
      dialog,
      close,
      tag,
      activeRules,
      t,
    };
  },
});
</script>

<template>
  <slot :showDialog="showDialog">
    <v-btn text="Replace me" size="small" color="primary" @click.stop="showDialog"></v-btn>
  </slot>
  <v-dialog v-model="isVisible" content-class="msmall" scrollable width="500">
    <v-form ref="dialog" @submit.prevent="doDialogAction">
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
        <v-card-text class="pt-2">
          <v-row dense>
            <v-col cols="12" xs="12" class="errorBorder">
              <v-text-field
                autocomplete="off"
                variant="outlined"
                :rules="activeRules.tag"
                v-model="tag"
                :label="t('COL_SBOM_TAG')"
                autofocus />
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <DCActionButton size="small" variant="text" @click="close" class="mr-5" :text="t('BTN_CLOSE')" />
          <DCActionButton size="small" variant="flat" @click="doDialogAction" :text="t('NP_DIALOG_BTN_EDIT')" />
        </v-card-actions>
      </v-card>
    </v-form>
  </v-dialog>
</template>
