<script setup lang="ts">
import {formatDateAndTime} from '@disclosure-portal/utils/Table';
import {computed} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRouter} from 'vue-router';

const {t} = useI18n();
const router = useRouter();

const props = defineProps<{
  sboms: any[];
  spdx?: string;
  projectId?: string;
  version?: string;
}>();

// Get the current SBOM object based on spdx ID
const currentSBOM = computed(() => {
  if (!props.spdx || !props.sboms || props.sboms.length === 0) {
    return null;
  }
  return props.sboms.find((sbom) => sbom.id === props.spdx);
});

// Get all SBOMs for the current version
const versionSBOMs = computed(() => {
  if (!props.version || !props.sboms || props.sboms.length === 0) {
    return [];
  }
  return props.sboms.filter((sbom) => sbom.version === props.version);
});

// Handle SBOM selection
const handleSBOMSelect = (sbom: any) => {
  if (sbom && props.projectId) {
    router.push(`/projects/${props.projectId}/versions/${props.version}/components/${sbom.id}`);
  }
};
</script>

<template>
  <div v-if="versionSBOMs.length" class="d-flex align-center pa-3">
    <v-select
      :model-value="spdx"
      :items="versionSBOMs"
      item-title="id"
      item-value="id"
      return-object
      variant="outlined"
      clearable
      density="compact"
      hide-details
      :label="`${t('COL_CHANNEL')} '${version}' > ${t('TAB_PROJECT_SBOM_LIST')}`"
      @update:model-value="handleSBOMSelect"
      style="min-width: 400px">
      <template v-slot:item="{internalItem, props: itemProps}">
        <v-list-item v-bind="itemProps" title="">
          <span class="d-subtitle-2" v-if="internalItem.raw.details?.uploaded">
            {{ formatDateAndTime(internalItem.raw.details.uploaded) }}
          </span>
          <span class="d-text d-secondary-text">&nbsp;-&nbsp;{{ internalItem.raw.details?.name }}</span>
          <span class="d-text d-secondary-text">/{{ version }}</span>
          <span class="d-text d-secondary-text" v-if="internalItem.raw.tag">&nbsp;[{{ internalItem.raw.tag }}]</span>
        </v-list-item>
      </template>
      <template v-slot:selection>
        <div v-if="currentSBOM" class="d-inline">
          <span class="d-subtitle-2" v-if="currentSBOM.details?.uploaded">
            {{ formatDateAndTime(currentSBOM.details.uploaded) }}
          </span>
          <span class="d-text d-secondary-text">&nbsp;-&nbsp;{{ currentSBOM.details?.name }}</span>
        </div>
      </template>
    </v-select>
  </div>
</template>
