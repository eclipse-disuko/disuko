<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script lang="ts" setup>
import adminService from '@disclosure-portal/services/admin';
import ProjectService from '@disclosure-portal/services/projects';
import {useAppStore} from '@disclosure-portal/stores/app';
import {formatDate} from '@disclosure-portal/utils/Table';
import {useClipboard} from '@shared/utils/clipboard';
import {TOOLTIP_OPEN_DELAY_IN_MS} from '@shared/utils/constant';
import {computed, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRouter} from 'vue-router';
import {useUrls} from '@shared/composables/useUrls';
import {ProjectModel} from '@disclosure-portal/model/Project';
import Icons from '@disclosure-portal/constants/icons';
import {ProjectSlimDto} from '@disclosure-portal/model/ProjectsResponse';

interface Props {
  item: ProjectModel;
  isAsync: boolean;
}

const props = defineProps<Props>();

const {t} = useI18n();
const appStore = useAppStore();
const {copyToClipboard} = useClipboard();
const router = useRouter();
const {openUrl} = useUrls();

const childProjects = ref<ProjectSlimDto[]>([]);
const labelTools = computed(() => appStore.getLabelsTools);
const childProjectAccessMap = ref(new Map<string, boolean>());
const icons = Icons;

const getAccessForProjectKey = (projectKey: string): boolean => {
  return childProjectAccessMap.value.get(projectKey) ?? false;
};

const getReferenceInfoForClipboard = async (item) => {
  const schemaLabelName = labelTools.value.schemaLabelsMap[item.schemaLabel]
    ? labelTools.value.schemaLabelsMap[item.schemaLabel].name
    : 'UNKNOWN_LABEL';
  const policyLabelNames = item.policyLabels
    .map((l) => (labelTools.value.policyLabelsMap[l] ? labelTools.value.policyLabelsMap[l].name : 'UNKNOWN_LABEL'))
    .join(', ');
  const projectLink = `${window.location.origin}/#/dashboard/projects/${encodeURIComponent(item._key)}`;

  const refInfo = `Disclosure Portal Project Reference

Project Name: ${item.name}
Project Identifier: ${item._key}
Project Schema Label: ${schemaLabelName}
Project Policy Labels: ${policyLabelNames}
Project Link: ${projectLink}
Application Name: ${item.applicationMeta.name}
Application Link: ${item.applicationMeta.externalLink}`;

  if (props.isAsync) {
    if (item.value.allowUsers.read || item.value.allowAllProjectUserManagement.read) {
      const userMail = await adminService.getUserMailById(item.responsible);
      return `${refInfo}
Project Responsible with Mail: ${item.responsible} (${userMail.email})`;
    } else {
      return `${refInfo}
Project Responsible: ${item.responsible}`;
    }
  }
  return refInfo;
};

const copyReferenceInfoToClipboard = async (item: ProjectSlimDto) => {
  const content = await getReferenceInfoForClipboard(item);
  copyToClipboard(content);
};

const onClickRow = (item: ProjectSlimDto) => {
  let url = '/dashboard/projects/';
  if (item.isGroup) {
    url = '/dashboard/groups/';
  }
  url += encodeURIComponent(item._key);
  openUrl(url, router);
};

onMounted(async () => {
  const response = await ProjectService.getChildren(props.item._key);
  childProjects.value = response.projects;
  childProjectAccessMap.value = new Map(
    response.list.map(({projectKey, hasProjectReadAccess}) => [projectKey, hasProjectReadAccess]),
  );
});
</script>

<template>
  <div class="ma-5">
    <div v-if="!item.isGroup">
      <v-table dark dense height="190px">
        <thead>
          <tr class="theme-background">
            <th
              :style="{
                backgroundColor: 'rgb(var(--v-theme-tableHeaderBackgroundColor))',
                borderBottom: 'solid 2px rgb(var(--v-theme-tableBorderColor))',
              }">
              {{ t('Labels') }}
            </th>
            <th
              :style="{
                backgroundColor: 'rgb(var(--v-theme-tableHeaderBackgroundColor))',
                borderBottom: 'solid 2px rgb(var(--v-theme-tableBorderColor))',
              }">
              {{ t('DESCRIPTION') }}
            </th>
          </tr>
        </thead>
        <tbody>
          <tr :style="{backgroundColor: 'rgb(var(--v-theme-tableHeaderBackgroundColor))'}">
            <td :style="{backgroundColor: 'rgb(var(--v-theme-tableHeaderBackgroundColor))', borderBottom: '0px'}">
              <v-row class="py-4" v-if="item.isGroup">
                <v-col cols="12" xs="12" class="py-1">
                  <DLabel :labelName="t('LBL_PROJECT_PARENT')" :iconName="icons.BACKUP" />
                </v-col>
              </v-row>
              <v-row class="py-4" v-if="!item.isGroup">
                <v-col cols="12" xs="12" class="py-1">
                  <v-tooltip :open-delay="TOOLTIP_OPEN_DELAY_IN_MS" bottom content-class="dpTooltip">
                    <template v-slot:activator="{props}">
                      <DLabel
                        :labelName="
                          labelTools.schemaLabelsMap[item.schemaLabel]
                            ? labelTools.schemaLabelsMap[item.schemaLabel].name
                            : 'UNKNOWN_LABEL'
                        "
                        :iconName="icons.SCHEMA"
                        v-bind="props" />
                    </template>
                    <span>{{ t('TT_schema_label_with_description') }}</span>
                    <span>
                      {{
                        labelTools.schemaLabelsMap[item.schemaLabel]
                          ? labelTools.schemaLabelsMap[item.schemaLabel].description
                          : ''
                      }}</span
                    >
                  </v-tooltip>
                  <v-tooltip
                    :open-delay="TOOLTIP_OPEN_DELAY_IN_MS"
                    bottom
                    v-for="(l, i) in item.freeLabels"
                    :key="'b' + i"
                    content-class="dpTooltip">
                    <template v-slot:activator="{props}">
                      <DLabel :labelName="l" :iconName="icons.TAG" v-bind="props" />
                    </template>
                    <span>{{ t('TT_free_label') }}</span>
                  </v-tooltip>
                  <DLabel :labelName="t('LBL_PROJECT_CHILD')" :iconName="icons.CHILD" v-if="item.parent" />
                </v-col>
                <v-col cols="12" xs="12" class="py-1">
                  <v-tooltip
                    :open-delay="TOOLTIP_OPEN_DELAY_IN_MS"
                    bottom
                    v-for="(l, i) in item.policyLabels"
                    :key="'a' + i"
                    content-class="dpTooltip">
                    <template v-slot:activator="{props}">
                      <DLabel
                        :labelName="
                          labelTools.policyLabelsMap[l] ? labelTools.policyLabelsMap[l].name : 'UNKNOWN_LABEL'
                        "
                        :iconName="icons.POLICY"
                        v-bind="props" />
                    </template>
                    <span>{{ t('TT_policy_label_with_description') }}</span>
                    <span>{{ labelTools.policyLabelsMap[l] ? labelTools.policyLabelsMap[l].description : '' }}</span>
                  </v-tooltip>
                </v-col>
                <v-col cols="12" xs="12" class="py-1">
                  <v-tooltip
                    :open-delay="TOOLTIP_OPEN_DELAY_IN_MS"
                    bottom
                    v-for="(l, i) in item.projectLabels"
                    :key="'a' + i"
                    content-class="dpTooltip">
                    <template v-slot:activator="{props}">
                      <DLabel
                        :labelName="
                          labelTools.projectLabelsMap[l] ? labelTools.projectLabelsMap[l].name : 'UNKNOWN_LABEL'
                        "
                        :iconName="icons.PROJECT_LABEL"
                        v-bind="props" />
                    </template>
                    <span>{{ t('TT_project_label_with_description') }}</span>
                    <span>{{ labelTools.projectLabelsMap[l] ? labelTools.projectLabelsMap[l].description : '' }}</span>
                  </v-tooltip>
                  <DLabel :labelName="t('LBL_PROJECT_NON_FOSS')" :iconName="icons.NON_FOSS" v-if="item.isNoFoss" />
                </v-col>
              </v-row>
            </td>
            <td :style="{backgroundColor: 'rgb(var(--v-theme-tableHeaderBackgroundColor))', borderBottom: '0px'}">
              {{ item.description }}
            </td>
          </tr>
        </tbody>
      </v-table>
    </div>

    <div v-if="item.isGroup">
      <v-card>
        <v-table dark dense>
          <colgroup>
            <col style="min-width: 6%" />
            <col style="min-width: 5%" />
            <col style="width: auto" />
            <col style="min-width: 5%" />
            <col style="min-width: 5%" />
          </colgroup>
          <thead>
            <tr style="cursor: default">
              <th class="w-24">{{ t('COL_ACTIONS') }}</th>
              <th class="w-40">{{ t('COL_STATUS') }}</th>
              <th>{{ t('COL_CHILD_NAME') }}</th>
              <th>{{ t('COL_DESCRIPTION') }}</th>
              <th>{{ t('COL_UPDATED') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr class="theme-background" v-if="childProjects.length == 0">
              <td>{{ t('NO_DATA_AVAILABLE') }}</td>
              <td></td>
              <td></td>
              <td></td>
              <td></td>
            </tr>
            <tr
              class="theme-background"
              @click="!child.isDeleted && getAccessForProjectKey(child._key) ? onClickRow(child) : undefined"
              :style="!child.isDeleted && getAccessForProjectKey(child._key) ? 'cursor: pointer' : 'cursor: default;'"
              v-for="child in childProjects"
              :key="child._key"
              v-else>
              <td>
                <DIconButton
                  v-if="item.accessRights.allowProject.read && getAccessForProjectKey(child._key)"
                  icon="mdi-content-copy"
                  :hint="t('TT_COPY_REFERENCE_INFO')"
                  @clicked="copyReferenceInfoToClipboard(child)" />
              </td>
              <td>
                <span v-if="!child.isDeleted" :class="'pStatus' + (!child.status ? 'new' : child.status)">{{
                  !child.status ? 'new' : child.status
                }}</span>
              </td>
              <td>
                <span v-if="child.isDeleted" class="disabledText">
                  {{ child.name }}
                  <span class="deleted">&nbsp;{{ t('PROJECT_DELETED') }}</span>
                </span>
                <span v-else-if="!getAccessForProjectKey(child._key)" class="disabledText">
                  {{ child.name }}
                  <span class="deleted">&nbsp;{{ t('INSUFFICIENT_PERMISSIONS') }}</span>
                </span>
                <span v-else>{{ child.name }}</span>
              </td>
              <td>{{ child.description.substring(0, 100) }}{{ child.description.length > 100 ? ' ...' : '' }}</td>
              <td>{{ formatDate(child.updated) }}</td>
            </tr>
          </tbody>
        </v-table>
      </v-card>
    </div>
  </div>
</template>

<style>
.theme-background {
  background-color: rgb(var(--v-theme-backgroundColor)) !important;
}

.deleted {
  color: rgb(var(--v-theme-error)) !important;
}
</style>
