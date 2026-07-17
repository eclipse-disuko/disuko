<script setup lang="ts">
import {useAppStore} from '@cli/stores/app';
import {createReusableTemplate} from '@vueuse/core';
import dayjs from 'dayjs';
import {computed} from 'vue';
import {useI18n} from 'vue-i18n';

const [DefineTemplate, ReuseTemplate] = createReusableTemplate();
const {t} = useI18n();
const appStore = useAppStore();

const currentProject = computed(() => appStore.getCurrentProject());
const idLabel = computed(() => {
  return currentProject.value?.isGroup ? t('GROUP_IDENTIFIER') : t('PROJECT_IDENTIFIER');
});
const idCopyHint = computed(() => idLabel.value);
const idCopyContent = computed(() => `${idLabel.value}: ${currentProject.value?.uuid ?? ''}`);
const createdDate = computed(() => convertToShort(currentProject.value?.created || ''));
const updatedDate = computed(() => convertToShort(currentProject.value?.updated || ''));

const convertToShort = (str: string) => {
  if (!str) {
    return '';
  }
  return dayjs(str).format(t('DATETIME_FORMAT_SHORT'));
};
</script>

<template>
  <div v-if="currentProject" class="project-overview">
    <DefineTemplate>
      <Stack direction="row" justify="between">
        <Stack style="flex: 1.5">
          <div>
            <p class="text-caption text-grey-darken-1">{{ idLabel }}</p>
            <span class="text-body-2">{{ currentProject.uuid }}</span>
            <DCopyClipboardButton :hint="idCopyHint" :content="idCopyContent" />
          </div>
        </Stack>
        <Stack style="flex: 1">
          <p class="text-caption text-grey-darken-1">{{ t('PROJECT_CURRENT_SCHEMA') }}</p>
          <span class="text-body-2">{{ currentProject.schema || t('NO_SCHEMA_INFO') }}</span>
        </Stack>
        <Stack style="flex: 1">
          <p class="text-caption text-grey-darken-1">{{ t('CREATED') }}</p>
          <span class="text-body-2">{{ createdDate }}</span>
        </Stack>
        <Stack style="flex: 1">
          <p class="text-caption text-grey-darken-1">{{ t('UPDATED') }}</p>
          <span class="text-body-2">{{ updatedDate }}</span>
        </Stack>
      </Stack>

      <Stack>
        <p class="text-caption text-grey-darken-1">{{ t('DESCRIPTION') }}</p>
        <span class="text-body-2">{{ currentProject.description }}</span>
      </Stack>
    </DefineTemplate>

    <GridVersionsPublic v-if="!currentProject.isGroup && currentProject.versions?.length">
      <ReuseTemplate></ReuseTemplate>
    </GridVersionsPublic>

    <GridChildrenPublic v-if="currentProject.isGroup" :children="currentProject.children || []">
      <ReuseTemplate></ReuseTemplate>
    </GridChildrenPublic>
  </div>

  <Stack v-else justify="center" align="center" class="py-8">
    <v-alert type="info" variant="tonal" class="ma-4" width="100%" max-width="600" border="start">
      <template v-slot:prepend>
        <v-icon color="info">mdi-information</v-icon>
      </template>
      {{ t('SELECT_PROJECT_TOKEN') }}
    </v-alert>
  </Stack>
</template>
