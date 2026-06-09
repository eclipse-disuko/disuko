<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import TableLayout from '@shared/layouts/TableLayout.vue';
import {computed, onMounted, ref} from 'vue';
import {ApprovableInfo} from '@disclosure-portal/model/Approval';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import projectService from '@disclosure-portal/services/projects';
import {useIdleStore} from '@shared/stores/idle.store';
const projectStore = useProjectStore();

const approvableInfo = ref<ApprovableInfo>({} as ApprovableInfo);
const useLatestSbom = ref(false);
// const childProjectChannels = ref<Map<string, VersionSlim>>(new Map());

const idle = useIdleStore();

const projectModel = computed(() => projectStore.currentProject!);

async function reload() {
  idle.showIdle = true;

  approvableInfo.value = await projectService.getApprovableInfo(projectModel.value._key, useLatestSbom.value);

  // childProjectChannels.value.clear();
  // const versionFetchPromises = approvableInfo.value.projects
  //   .filter((p) => p.approvablespdx.versionkey)
  //   .map(async (project) => {
  //     try {
  //       const versionDetails = await versionService.getVersion(project.projectKey, project.approvablespdx.versionkey);
  //       childProjectChannels.value.set(project.approvablespdx.versionkey, versionDetails.data);
  //     } catch (error) {
  //       console.error(`Failed to fetch version details for project ${project.projectKey}:`, error);
  //     }
  //   });
  // await Promise.all(versionFetchPromises);

  idle.showIdle = false;
}

onMounted(async () => {
  await reload();
});
</script>

<template>
  <TableLayout has-tab has-title>
    <template #description v-if="$slots.default">
      <slot></slot>
    </template>
    <template #table>
      <div ref="tableUserManagement" class="fill-height">
        <v-checkbox v-model="useLatestSbom" label="Use Latest SBOM" @change="reload" density="compact"> </v-checkbox>
        <!--        <GridSPDXList-->
        <!--          :projects="approvableInfo.projects"-->
        <!--          :channels="childProjectChannels"-->
        <!--          showSbomExtras-->
        <!--          showSupplier />-->
        <GridSPDXList
          :projects="approvableInfo.projects"
          :channels="projectModel.versions"
          showSbomExtras
          showSupplier />
      </div>
    </template>
  </TableLayout>
</template>
