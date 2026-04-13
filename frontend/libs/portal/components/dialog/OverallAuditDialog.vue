<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {OverallReviewRequest, OverallReviewState, SpdxFile} from '@disclosure-portal/model/VersionDetails';
import versionService from '@disclosure-portal/services/version';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useSbomStore} from '@disclosure-portal/stores/sbom.store';
import useRules from '@disclosure-portal/utils/Rules';
import {formatDateAndTime, getOverallReviewTranslationKey} from '@disclosure-portal/utils/Table';
import useSnackbar from '@shared/composables/useSnackbar';
import {computed, nextTick, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {VForm} from 'vuetify/components';

const emit = defineEmits(['reload']);

const {minMax} = useRules();
const {t} = useI18n();
const {info: snack} = useSnackbar();
const projectStore = useProjectStore();
const sbomStore = useSbomStore();

const isVisible = ref(false);
const selectedState = ref<OverallReviewState>(OverallReviewState.AUDITED);
const selectedSBOM = ref<SpdxFile | null>(null);
const comment = ref('');
const possibleStates = [OverallReviewState.AUDITED];
const form = ref<VForm | null>(null);

const currentProject = computed(() => projectStore.currentProject!);
const channelSpdxs = computed(() => sbomStore.channelSpdxs);
const selectedSpdx = computed(() => sbomStore.selectedSpdx);

const rules = {
  comment: minMax(t('ATTR_COMMENT'), 0, 500, false),
};

const open = () => {
  selectedSBOM.value = selectedSpdx.value;
  comment.value = t('AUDIT_ATTR_COMMENT');
  isVisible.value = true;
};

const close = () => {
  isVisible.value = false;
};

const save = async () => {
  await nextTick();
  form.value?.validate().then(async (info) => {
    if (!info.valid) return;

    const req = {
      state: OverallReviewState.UNREVIEWED,
      comment: '',
      sbomId: '',
      sbomName: '',
      sbomUploaded: '',
    } as OverallReviewRequest;

    req.state = selectedState.value;
    if (selectedSBOM.value?._key) {
      req.sbomId = selectedSBOM.value._key;
    } else {
      req.sbomId = selectedSBOM.value as unknown as string;
    }
    selectedSBOM.value = channelSpdxs.value.find((sbom) => sbom._key === req.sbomId) || null;

    if (selectedSBOM.value) {
      req.sbomId = selectedSBOM.value._key;
      req.comment = comment.value;
      req.sbomName = selectedSBOM.value.MetaInfo?.Name || '';
      req.sbomUploaded = selectedSBOM.value.Uploaded;
    }
    await versionService.createOverallReview(currentProject.value._key, sbomStore.currentVersion._key, req);
    await projectStore.fetchProjectByKey(currentProject.value._key);
    emit('reload');
    close();
    snack(t('DIALOG_overallreview_create_success'));
  });
};

defineExpose({open});
</script>

<template>
  <v-dialog v-model="isVisible" width="500">
    <DialogLayout
      :config="{
        title: t('HEADLINE_OVERALL_AUDIT'),
        secondaryButton: {text: t('BTN_CANCEL')},
        primaryButton: {text: t('Btn_submit')},
      }"
      @close="close"
      @secondary-action="close"
      @primary-action="save">
      <v-form ref="form">
        <Stack>
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
          <v-textarea
            auto-grow
            variant="outlined"
            density="compact"
            :label="t('OVERALL_REVIEW_COMMENT')"
            v-model="comment"
            :counter="500"
            :rules="rules.comment" />
          <v-select
            variant="outlined"
            density="compact"
            :items="channelSpdxs"
            v-model="selectedSBOM"
            item-text="_key"
            item-value="_key"
            :label="t('SBOM_DELIVERIES')">
            <template v-slot:item="{item, props}">
              <v-list-item v-bind="props" title="">
                <v-icon
                  color="primary"
                  v-if="currentProject.approvablespdx.spdxkey === item.raw._key"
                  size="small"
                  class="pr-2"
                  icon="mdi-star"></v-icon>
                <span class="d-subtitle-2">{{ formatDateAndTime(item.raw.Uploaded) }}</span>
                <span class="d-text d-secondary-text"> - {{ item.raw.MetaInfo.Name }}</span>
                <span class="d-text d-secondary-text ml-1" v-if="item.raw.Tag">({{ item.raw.Tag }})</span>
                <span class="d-text d-secondary-text" v-if="item.raw.isRecent"> [{{ t('SBOM_LATEST') }}] </span>
                <span class="d-text d-secondary-text" v-else> [{{ t('SBOM_FORMER') }}] </span>
              </v-list-item>
            </template>
            <template v-slot:selection="{item}">
              <div class="d-inline">
                <v-icon
                  color="primary"
                  v-if="currentProject.approvablespdx.spdxkey === item.raw._key"
                  size="small"
                  class="pr-2"
                  icon="mdi-star"></v-icon>
                <span class="d-subtitle-2">{{ formatDateAndTime(item.raw.Uploaded) }}</span>
                <span class="d-text d-secondary-text"> - {{ item.raw.MetaInfo.Name }}</span>
                <span class="d-text d-secondary-text ml-1" v-if="item.raw.Tag">({{ item.raw.Tag }})</span>
                <span class="d-text d-secondary-text" v-if="item.raw.isRecent"> [{{ t('SBOM_LATEST') }}] </span>
                <span class="d-text d-secondary-text" v-else> [{{ t('SBOM_FORMER') }}] </span>
              </div>
            </template>
          </v-select>
        </Stack>
      </v-form>
    </DialogLayout>
  </v-dialog>
</template>
