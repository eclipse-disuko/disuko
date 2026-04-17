<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<template>
  <v-dialog v-model="show" content-class="large" scrollable width="1200" @after-leave="onAfterLeave">
    <DialogLayout
      :config="{
        title: item?.title,
        secondaryButton: {text: t('BTN_CLOSE')},
      }"
      @secondary-action="close"
      @close="close">
      <template #title-right>
        <div v-if="isOpen(item) || isInProgress(item)" class="ml-4 flex gap-2">
          <DCActionButton
            v-if="isOpen(item)"
            is-dialog-button
            size="small"
            variant="outlined"
            @click="markInProgress"
            :text="t('RR_MARK_IN_PROGRESS')" />
          <DCActionButton
            is-dialog-button
            size="small"
            variant="outlined"
            @click="closeRemark"
            :text="t('TT_BULK_CLOSE_REMARK')" />
        </div>
      </template>
      <v-tabs v-model="currentTab" slider-color="mbti" active-class="active" show-arrows bg-color="tabsHeader">
        <v-tab value="details">{{ t('TAB_TITLE_DETAILS') }}</v-tab>
        <v-tab value="tabComment">{{ t('COL_COMMENTS') }}</v-tab>
      </v-tabs>
      <v-tabs-window v-model="currentTab" class="h-100 overflow-auto">
        <v-tabs-window-item value="details">
          <div class="mt-4 grid grid-cols-2 gap-4">
            <div class="d-text label-text-field-color">
              <p class="mb-2 font-semibold">{{ t('COL_DESCRIPTION') }}:</p>
              <div>
                <template v-for="(part, index) in descriptionParts" :key="index">
                  <DExternalLink v-if="part.isUrl" :url="part.text" :text="part.text" />
                  <span v-else>{{ part.text }}</span>
                </template>
              </div>
            </div>
            <div class="d-text label-text-field-color">
              <p class="mb-2 font-semibold">{{ t('COL_REFERENCES') }}:</p>

              <!-- SBOM Reference -->
              <div v-if="item?.sbomId !== ''" class="mb-3">
                <span class="text-sm font-medium">{{ t('UM_DIALOG_REVIEW_REMARK_SBOM') }}:</span>
                <span class="ml-1">{{ item?.sbomName }}</span>
                <span class="ml-1">(<DDateCellWithTooltip :value="item?.sbomUploaded || ''" />)</span>
              </div>

              <!-- Components -->
              <div v-if="componentsDisplay.length > 0" class="mb-3">
                <p class="mb-2 text-sm font-medium">{{ t('COMPONENTS') }}:</p>
                <div class="flex flex-wrap gap-1">
                  <router-link
                    v-for="component in componentsDisplay"
                    :to="component.to"
                    target="_blank"
                    :key="component.key">
                    <v-chip size="small" variant="outlined" color="secondary">
                      <span class="font-semibold">{{ component.name }}</span>
                      <span class="ml-1 text-gray-500">@{{ component.version }}</span>
                    </v-chip>
                  </router-link>
                </div>
              </div>

              <!-- Licenses -->
              <div v-if="licensesDisplay.length > 0">
                <p class="mb-2 text-sm font-medium">{{ t('LICENSES') }}:</p>
                <div class="flex flex-wrap gap-1">
                  <template v-if="user.getRights.allowLicense.read">
                    <router-link v-for="license in licensesDisplay" :to="license.to" :key="license.key" target="_blank">
                      <v-chip size="small" color="secondary" variant="outlined">
                        {{ license.label }}
                      </v-chip>
                    </router-link>
                  </template>
                  <template v-else>
                    <v-chip
                      v-for="license in licensesDisplay"
                      :key="license.key"
                      size="small"
                      color="secondary"
                      variant="outlined">
                      {{ license.label }}
                    </v-chip>
                  </template>
                </div>
              </div>
            </div>
          </div>
        </v-tabs-window-item>
        <v-tabs-window-item value="tabComment">
          <GridReviewComments
            v-if="item"
            :readonly="!(isOpen(item) || isInProgress(item))"
            :reviewId="item.key"
            :events="item.events"
            :disableClick="isSubmittingComment"
            @comment="comment" />
        </v-tabs-window-item>
      </v-tabs-window>
      <template #left>
        <div class="text-sm italic">
          <p>{{ t('RR_WARN_TEXT') }}</p>
        </div>
      </template>
    </DialogLayout>
  </v-dialog>
</template>

<script lang="ts" setup>
import {
  CommentReviewRemarkRequest,
  ReviewRemark,
  ReviewRemarkStatus,
  SetReviewRemarkStatusRequest,
} from '@disclosure-portal/model/Quality';
import versionService from '@disclosure-portal/services/version';
import {useUserStore} from '@disclosure-portal/stores/user';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRoute} from 'vue-router';
import useSnackbar from '@shared/composables/useSnackbar';

const props = defineProps<{
  projectUuid: string;
  versionUuid: string;
}>();
const emit = defineEmits(['reload', 'closeRemark']);

const {t} = useI18n();
const route = useRoute();
const user = useUserStore();

const show = ref(false);
const currentTab = ref(0);
const item = ref<ReviewRemark>();
const isSubmittingComment = ref(false);
const {info: snack} = useSnackbar();

const currentSbom = computed(() =>
  Array.isArray(route.params.currentSbom) ? route.params.currentSbom[0] : route.params.currentSbom,
);

const componentsDisplay = computed(() => {
  if (!item.value?.components) {
    return [];
  }

  return item.value.components.map((c, index) => ({
    key: `${c.componentName}-${c.componentVersion}-${index}`,
    id: c.componentId,
    to: `/dashboard/projects/${props.projectUuid}/versions/${props.versionUuid}/component/${currentSbom.value}/${c.componentId}`,
    name: c.componentName,
    version: c.componentVersion,
  }));
});

const licensesDisplay = computed(() => {
  if (!item.value?.licenses) {
    return [];
  }

  return item.value.licenses.map((l, index) => ({
    key: `${l.licenseId}-${index}`,
    to: `/dashboard/licenses/${l.licenseId}/details`,
    label: l.licenseName === '' ? `${l.licenseId} (${t('TT_REVIEW_REMARK_DIALOG_LICENSE_UNKNOWN')})` : l.licenseName,
  }));
});

const descriptionParts = computed(() => {
  const description = item.value?.description || '';
  const urlRegex = /(https?:\/\/[^\s]+)/g;
  return description.split(urlRegex).map((part) => ({
    text: part,
    isUrl: urlRegex.test(part),
  }));
});

const open = (reviewRemark: ReviewRemark): void => {
  show.value = true;
  item.value = reviewRemark;
};

const close = (): void => {
  show.value = false;
};
const onAfterLeave = (): void => {
  currentTab.value = 0;
};

const isOpen = (item: ReviewRemark): boolean => {
  return item.status === ReviewRemarkStatus.OPEN;
};

const isInProgress = (item: ReviewRemark): boolean => {
  return item.status === ReviewRemarkStatus.IN_PROGRESS;
};

const comment = async (content: string) => {
  if (isSubmittingComment.value) {
    return;
  }

  isSubmittingComment.value = true;
  try {
    const req: CommentReviewRemarkRequest = {
      content: content,
    };
    await versionService.commentReviewRemark(props.projectUuid, props.versionUuid, item.value!.key, req);

    // Fetch updated review remarks
    const response = await versionService.getReviewRemarks(props.projectUuid, props.versionUuid);
    const updatedRemark = response.data.find((remark) => remark.key === item.value!.key);
    if (updatedRemark) {
      item.value = updatedRemark;
    }

    emit('reload');
  } finally {
    isSubmittingComment.value = false;
  }
};

const markInProgress = async () => {
  if (!item.value) return;

  const req: SetReviewRemarkStatusRequest = {
    status: ReviewRemarkStatus.IN_PROGRESS,
  };

  try {
    await versionService.setReviewRemarkStatus(props.projectUuid, props.versionUuid, item.value.key, req);

    // Fetch updated review remarks
    const response = await versionService.getReviewRemarks(props.projectUuid, props.versionUuid);
    const updatedRemark = response.data.find((remark) => remark.key === item.value!.key);
    if (updatedRemark) {
      item.value = updatedRemark;
    }
    snack(t('DIALOG_remark_in_progress'));
    emit('reload');
  } catch (error) {
    console.error('Failed to mark remark as in progress:', error);
  }
};

const closeRemark = () => {
  if (item.value) {
    emit('closeRemark', item.value);
  }
  close();
};

defineExpose({
  open,
});
</script>
