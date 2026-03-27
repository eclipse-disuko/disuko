// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import Icons from '@disclosure-portal/constants/icons';
import type {Project} from '@disclosure-portal/model/Project';
import {useLabelStore} from '@disclosure-portal/stores/label.store';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {formatDate, getStrWithMaxLength} from '@disclosure-portal/utils/View';
import {createReusableTemplate} from '@vueuse/core';
import dayjs from 'dayjs';
import {computed} from 'vue';
import {useI18n} from 'vue-i18n';
import ProjectLabel from '../ProjectLabel.vue';

const projectStore = useProjectStore();
const labelStore = useLabelStore();
const {t} = useI18n();
const icons = Icons;
const [DefineTemplate, ReuseTemplate] = createReusableTemplate();

const currentProject = computed((): Project => projectStore.currentProject!);

const createdDate = computed(() => convertToShort(currentProject.value?.created || ''));
const updatedDate = computed(() => convertToShort(currentProject.value?.updated || ''));
const deleteDate = computed(() => formatDate(currentProject.value?.dummyDeletionDate || ''));

const convertToShort = (str: string) => {
  if (!str) {
    return '';
  }
  return dayjs(str).format(t('DATETIME_FORMAT_SHORT'));
};

const uuidLabel = computed(() => {
  return currentProject.value.isGroup ? t('GROUP_IDENTIFIER_UUID') : t('APPLICATION_UUID');
});
</script>

<template>
  <DefineTemplate>
    <v-row>
      <v-col cols="12" xs="12" md="3">
        <p class="text-caption text-grey-darken-1">{{ uuidLabel }}</p>
        <span class="text-body-2">{{ currentProject._key }}</span>
        <DCopyClipboardButton :hint="t('TT_UUIDCopyText')" :content="currentProject._key"></DCopyClipboardButton>
      </v-col>
      <v-col cols="12" xs="12" md="2" v-if="currentProject.applicationMeta.id">
        <p class="text-caption text-grey-darken-1">{{ t('APPLICATION_LINK') }}</p>
        <DExternalLink :text="currentProject.appLinkText" :url="currentProject.applicationMeta.externalLink">
        </DExternalLink>
      </v-col>
      <v-col cols="12" xs="12" md="3">
        <p class="text-caption text-grey-darken-1">{{ t('PROJECT_CURRENT_SCHEMA') }}</p>
        <span class="text-body-2" v-if="currentProject.correspondingSchema && currentProject.correspondingSchema.name">
          <v-icon color="primary" small>mdi mdi-chevron-right</v-icon>
          <a :href="'#/dashboard/schemas/' + currentProject.correspondingSchema._key" target="_blank">
            {{ currentProject.correspondingSchema.name }}-{{ currentProject.correspondingSchema.version }}
          </a>
        </span>
      </v-col>
      <v-col cols="12" xs="6" md="2">
        <p class="text-caption text-grey-darken-1">{{ t('CREATED') }}</p>
        <span class="text-body-2">{{ createdDate }}</span>
      </v-col>
      <v-col cols="12" xs="6" md="2">
        <p class="text-caption text-grey-darken-1">{{ t('UPDATED') }}</p>
        <span class="text-body-2">{{ updatedDate }}</span>
      </v-col>
      <v-col v-if="currentProject.isDummy" cols="12" xs="6" md="2">
        <p class="text-caption text-grey-darken-1">{{ t('DUMMY_DELETE_DATE') }}</p>
        <span class="text-body-2 text-yellow-500">{{ deleteDate }}</span>
      </v-col>
    </v-row>
    <v-row class="mb-4">
      <v-col cols="12" xs="12" :md="currentProject.applicationMeta.id ? '5' : '6'">
        <p class="text-caption text-grey-darken-1">{{ t('LABELS') }}</p>
        <div class="flex flex-wrap gap-1">
          <ProjectLabel :label="labelStore.getLabelByKey(currentProject.schemaLabel)" />

          <ProjectLabel
            v-for="(l, i) in currentProject.policyLabels"
            :key="'policy-' + i"
            :label="labelStore.getLabelByKey(l)" />

          <ProjectLabel
            v-for="(l, i) in currentProject.projectLabels"
            :key="'project-' + i"
            :label="labelStore.getLabelByKey(l)" />

          <span v-for="(l, i) in currentProject.freeLabels" :key="'tag-' + i">
            <DLabel :labelName="l" :iconName="icons.TAG"> </DLabel>
            <Tooltip>
              <span>{{ t('TT_free_label') }}</span>
            </Tooltip>
          </span>
          <DLabel :labelName="t('LBL_PROJECT_PARENT')" :iconName="icons.BACKUP" v-if="currentProject.isGroup"></DLabel>
          <DLabel :labelName="t('LBL_PROJECT_CHILD')" :iconName="icons.CHILD" v-if="currentProject.parent"></DLabel>
          <DLabel :labelName="t('LBL_PROJECT_NON_FOSS')" :iconName="icons.NON_FOSS" v-if="currentProject.isNoFoss">
          </DLabel>
        </div>
      </v-col>
      <v-col cols="12" xs="12" :md="currentProject.applicationMeta.id ? '7' : '6'">
        <p class="text-caption text-grey-darken-1">{{ t('DESCRIPTION') }}</p>

        <span class="text-body-2">
          {{ getStrWithMaxLength(250, currentProject.description) }}
          <Tooltip>{{ currentProject.description }}</Tooltip>
        </span>
      </v-col>
    </v-row>
  </DefineTemplate>
  <GridVersions v-if="currentProject && !currentProject.isGroup" ref="projectVersions">
    <ReuseTemplate></ReuseTemplate>
  </GridVersions>
  <GridChildren v-if="currentProject && currentProject.isGroup" ref="groupChildren">
    <ReuseTemplate></ReuseTemplate>
  </GridChildren>
</template>
