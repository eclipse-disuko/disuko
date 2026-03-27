// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script lang="ts">
import {OverallReviewRequest, OverallReviewState, SpdxFile} from '@disclosure-portal/model/VersionDetails';
import versionService from '@disclosure-portal/services/version';
import {useAppStore} from '@disclosure-portal/stores/app';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import useRules from '@disclosure-portal/utils/Rules';
import {formatDateAndTime, getOverallReviewTranslationKey} from '@disclosure-portal/utils/Table';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import DCloseButton from '@shared/components/disco/DCloseButton.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import {defineComponent, nextTick, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {VForm} from 'vuetify/components';

export default defineComponent({
  name: 'OverallReviewDialog',
  components: {
    DCActionButton,
    DCloseButton,
  },
  setup(_, {emit}) {
    const {minMax} = useRules();
    const {t} = useI18n();
    const {info: snack} = useSnackbar();
    const projectStore = useProjectStore();
    const appStore = useAppStore();
    const isVisible = ref(false);
    const selectedState = ref<OverallReviewState>(OverallReviewState.UNREVIEWED);
    const selectedSBOM = ref<SpdxFile | null>(null);
    const sbomHistory = ref<SpdxFile[]>([]);
    const comment = ref('');
    const approvableSbomId = ref('');
    const possibleStates = [
      OverallReviewState.UNREVIEWED,
      OverallReviewState.ACCEPTABLE,
      OverallReviewState.ACCEPTABLE_AFTER_CHANGES,
      OverallReviewState.NOT_ACCEPTABLE,
    ];
    const prId = ref('');
    const vId = ref('');
    const form = ref<VForm | null>(null);

    const rules = {
      comment: minMax(t('ATTR_COMMENT'), 0, 500, false),
    };

    const open = (projectId: string, versionKey: string, spdxFileHistory: SpdxFile[], selectedSpdxKey: SpdxFile, approvableSpdxKey: string) => {
      prId.value = projectId;
      vId.value = versionKey;
      sbomHistory.value = spdxFileHistory;
      approvableSbomId.value = approvableSpdxKey;
      selectedSBOM.value = selectedSpdxKey;
      isVisible.value = true;
    };

    const close = () => {
      isVisible.value = false;
    };

    const save = async () => {
      await nextTick();
      form.value?.validate().then(async (info) => {
        if (!info.valid) {
          return;
        }
        const req = {
          state: OverallReviewState.UNREVIEWED,
          comment: '',
          sbomId: '',
          sbomName: '',
          sbomUploaded: '',
        } as OverallReviewRequest;

        req.state = selectedState.value;
        if (selectedSBOM.value?._key) {
          req.sbomId = selectedSBOM.value._key; // Use the key of the currently selected SBOM
        } else {
          req.sbomId = selectedSBOM.value as unknown as string; // Handle null or fallback assignment
        }
        selectedSBOM.value = sbomHistory.value.find((sbom) => sbom._key === req.sbomId) || null;

        if (selectedSBOM.value) {
          req.sbomId = selectedSBOM.value._key;
          req.comment = comment.value;
          req.sbomName = selectedSBOM.value.MetaInfo?.Name || '';
          req.sbomUploaded = selectedSBOM.value.Uploaded;
        }
        await versionService.createOverallReview(prId.value, vId.value, req); // Replace with actual IDs
        await projectStore.fetchProjectByKey(projectStore.currentProject!._key);
        await appStore.fetchSBOMHistory();
        appStore.resetCurrentVersion();
        emit('reload');
        close();
        snack(t('DIALOG_overallreview_create_success'));
      });
    };

    return {
      t,
      isVisible,
      selectedState,
      selectedSBOM,
      sbomHistory,
      approvableSbomId,
      possibleStates,
      comment,
      rules,
      open,
      close,
      save,
      formatDateAndTime,
      form,
      getOverallReviewTranslationKey,
    };
  },
});
</script>

<template>
  <v-dialog v-model="isVisible" width="500">
    <v-form ref="form">
      <v-card class="pa-8 dDialog" flat>
        <v-card-title>
          <v-row>
            <v-col cols="10">
              <span class="text-h5">{{ t('HEADLINE_OVERALL_REVIEW') }}</span>
            </v-col>
            <v-col cols="2" align="right">
              <DCloseButton @click="close" />
            </v-col>
          </v-row>
        </v-card-title>

        <v-card-text class="pt-2">
          <v-row dense>
            <v-col cols="12">
              <v-select
                variant="outlined"
                density="compact"
                :items="possibleStates"
                v-model="selectedState"
                :label="t('SELECT_OVERALL_REVIEW_STATE')">
                <template v-slot:item="{item, props}">
                  <v-list-item v-bind="props" title="">
                    <span class="d-subtitle-2 ml-2">{{ t(getOverallReviewTranslationKey(item.raw)) }}</span>
                  </v-list-item>
                </template>
                <template v-slot:selection="{item}">
                  <span class="d-subtitle-2 ml-2">{{ t(getOverallReviewTranslationKey(item.raw)) }}</span>
                </template>
              </v-select>
            </v-col>
          </v-row>

          <v-row>
            <v-col cols="12">
              <v-textarea
                auto-grow
                variant="outlined"
                density="compact"
                :label="t('OVERALL_REVIEW_COMMENT')"
                v-model="comment"
                :counter="500"
                :rules="rules.comment"></v-textarea>
            </v-col>
          </v-row>

          <v-row>
            <v-col cols="12">
              <v-select
                variant="outlined"
                density="compact"
                :items="sbomHistory"
                v-model="selectedSBOM"
                item-text="_key"
                item-value="_key"
                :label="t('SBOM_DELIVERIES')">
                <template v-slot:item="{item, props}">
                  <v-list-item v-bind="props" title="">
                    <v-icon color="primary" v-if="approvableSbomId === item.raw._key" size="small" class="pr-2"
                      >mdi-star</v-icon
                    >
                    <span class="d-subtitle-2">{{ formatDateAndTime(item.raw.Uploaded) }}</span>
                    <span class="d-text d-secondary-text"> - {{ item.raw.MetaInfo.Name }}</span>
                    <span class="d-text d-secondary-text ml-1" v-if="item.raw.Tag">({{ item.raw.Tag }})</span>
                    <span class="d-text d-secondary-text" v-if="item.raw.isRecent"> [{{ t('SBOM_LATEST') }}] </span>
                    <span class="d-text d-secondary-text" v-else> [{{ t('SBOM_FORMER') }}] </span>
                  </v-list-item>
                </template>
                <template v-slot:selection="{item}">
                  <div class="d-inline">
                    <v-icon color="primary" v-if="approvableSbomId === item.raw._key" size="small" class="pr-2"
                      >mdi-star</v-icon
                    >
                    <span class="d-subtitle-2">{{ formatDateAndTime(item.raw.Uploaded) }}</span>
                    <span class="d-text d-secondary-text"> - {{ item.raw.MetaInfo.Name }}</span>
                    <span class="d-text d-secondary-text ml-1" v-if="item.raw.Tag">({{ item.raw.Tag }})</span>
                    <span class="d-text d-secondary-text" v-if="item.raw.isRecent"> [{{ t('SBOM_LATEST') }}] </span>
                    <span class="d-text d-secondary-text" v-else> [{{ t('SBOM_FORMER') }}] </span>
                  </div>
                </template>
              </v-select>
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
          <DCActionButton isDialogButton size="small" variant="flat" @click="save" :text="t('Btn_submit')" />
        </v-card-actions>
      </v-card>
    </v-form>
  </v-dialog>
</template>
