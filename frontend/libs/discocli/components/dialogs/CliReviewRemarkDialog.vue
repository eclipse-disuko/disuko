<template>
  <v-dialog v-model="show" content-class="large" scrollable width="1200" @after-leave="resetDialog">
    <v-card class="pa-8 dDialog" flat>
      <v-card-title>
        <div style="display: flex; align-items: center; width: 100%">
          <span class="text-h5">{{ t('TITLE_REVIEW_REMARK') }}</span>
          <v-spacer></v-spacer>
          <DCloseButton @click="close" />
        </div>
      </v-card-title>
      <v-card-text v-if="item">
        <v-tabs v-model="currentTab" slider-color="mbti" active-class="active" show-arrows bg-color="tabsHeader">
          <v-tab value="details">{{ t('TAB_TITLE_DETAILS') }}</v-tab>
          <v-tab value="tabComment">{{ t('COL_COMMENTS') }}</v-tab>
        </v-tabs>
        <v-tabs-window v-model="currentTab" class="h-[600px] overflow-auto">
          <v-tabs-window-item value="details">
            <div class="mt-4 grid grid-cols-2 gap-4">
              <div class="text-base leading-relaxed">
                <p class="mb-2 font-semibold">{{ t('COL_DESCRIPTION') }}:</p>
                <div class="text-base leading-relaxed">
                  <template v-for="(part, index) in descriptionParts" :key="index">
                    <DExternalLink v-if="part.isUrl" :url="part.text" :text="part.text" />
                    <span v-else-if="part.text">{{ part.text }}</span>
                  </template>
                </div>
              </div>
              <div class="text-base leading-relaxed">
                <p class="mb-2 font-semibold">{{ t('COL_REFERENCES') }}:</p>

                <!-- SBOM Reference -->
                <div v-if="item.sbomId !== ''" class="mb-3">
                  <span class="text-sm font-medium">{{ t('UM_DIALOG_REVIEW_REMARK_SBOM') }}:</span>
                  <span class="ml-1">{{ item.sbomName }}</span>
                  <span class="ml-1">(<DDateCellWithTooltip :value="item.sbomUploaded || ''" />)</span>
                </div>

                <!-- Components -->
                <div v-if="componentsDisplay.length > 0" class="mb-3">
                  <p class="mb-2 text-sm font-medium">{{ t('COMPONENTS') }}:</p>
                  <div class="flex flex-wrap gap-1">
                    <v-chip
                      v-for="component in componentsDisplay"
                      :key="component.key"
                      size="small"
                      variant="outlined"
                      color="secondary">
                      <span class="font-semibold">{{ component.name }}</span>
                      <span v-if="component.version" class="ml-1 text-gray-500">({{ component.version }})</span>
                    </v-chip>
                  </div>
                </div>

                <!-- Licenses -->
                <div v-if="licensesDisplay.length > 0">
                  <p class="mb-2 text-sm font-medium">{{ t('LICENSES') }}:</p>
                  <div class="flex flex-wrap gap-1">
                    <v-chip
                      v-for="license in licensesDisplay"
                      :key="license.key"
                      size="small"
                      color="secondary"
                      variant="outlined">
                      {{ license.label }}
                    </v-chip>
                  </div>
                </div>
              </div>
            </div>
          </v-tabs-window-item>
          <v-tabs-window-item value="tabComment">
            <GridReviewComments
              :readonly="!isOpen(item)"
              :reviewId="item.key || ''"
              :events="item.events || []"
              @comment="handleComment" />
          </v-tabs-window-item>
        </v-tabs-window>
      </v-card-text>
      <v-card-text class="text-body-3 text-medium-emphasis pt-0">
        {{ t('LINKS_NETWORK_DISCLAIMER') }}
        <br />
        {{ t('RR_WARN_TEXT') }}
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <DCActionButton size="small" variant="flat" @click="close" :text="t('BTN_CLOSE')" />
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script setup lang="ts">
import type {RemarkWithVersion, RRCommentExternDTO} from '@cli/models/ReviewRemark';
import {projectService} from '@cli/services/projectService';
import DExternalLink from '@shared/components/disco/DExternalLink.vue';
import {computed, defineExpose, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRoute} from 'vue-router';

const {t} = useI18n();
const route = useRoute();
const show = ref(false);
const currentTab = ref('details');
const item = ref<RemarkWithVersion | null>(null);
const submitting = ref(false);
const emit = defineEmits(['refresh']);

const componentsDisplay = computed(() => {
  if (!item.value?.components || !Array.isArray(item.value.components)) {
    return [];
  }

  return item.value.components.map((c: any, index: number) => ({
    key: `${c.componentName}-${c.componentVersion}-${index}`,
    id: c.componentId,
    name: c.componentName,
    version: c.componentVersion,
  }));
});

const licensesDisplay = computed(() => {
  if (!item.value?.licenses || !Array.isArray(item.value.licenses)) {
    return [];
  }

  return item.value.licenses.map((l: any, index: number) => ({
    key: `${l.licenseId}-${index}`,
    label: l.licenseName === '' ? `${l.licenseId} (${t('TT_REVIEW_REMARK_DIALOG_LICENSE_UNKNOWN')})` : l.licenseName,
  }));
});

const descriptionParts = computed(() => {
  const description = item.value?.description || '';
  const urlRegex = /(https?:\/\/[^\s]+)/g;
  const parts = description.split(urlRegex);
  return parts.map((part) => ({
    text: part,
    isUrl: /^https?:\/\//.test(part),
  }));
});

function isOpen(item: any) {
  return item && ['OPEN', 'IN_PROGRESS'].includes(item.status);
}
function open(remark?: RemarkWithVersion) {
  item.value = remark || null;
  show.value = true;
  currentTab.value = 'details';
}
function close() {
  show.value = false;
}
function resetDialog() {
  item.value = null;
  currentTab.value = 'details';
}
function handleComment(content: string) {
  submitComment(content);
}
async function submitComment(content?: string) {
  if (!content?.trim() || !item.value || submitting.value) return;

  const commentText = content.trim();

  const projectUuid = typeof route.params.id === 'string' ? route.params.id : null;
  const version = item.value.version;
  const reviewRemarkUuid = item.value.key;

  if (!projectUuid || !version || !reviewRemarkUuid) {
    console.error('Missing required data for comment submission');
    return;
  }

  submitting.value = true;
  try {
    const comment: RRCommentExternDTO = {content: commentText};
    await projectService.commentOnReviewRemark(projectUuid, version, reviewRemarkUuid, comment);

    const remarks = await projectService.getVersionReviewRemarks(projectUuid, version);
    const updatedRemark = remarks?.find((remark) => remark.key === reviewRemarkUuid);

    if (updatedRemark) {
      item.value.events = updatedRemark.events;
    }
    emit('refresh');
  } catch (e) {
    console.error('Error submitting comment:', e);
  } finally {
    submitting.value = false;
  }
}
defineExpose({open});
</script>
