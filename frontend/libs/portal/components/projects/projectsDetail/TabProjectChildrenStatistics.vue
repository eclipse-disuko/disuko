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
const search = ref<string | null>('');
const dataAreLoaded = ref(false);
// const childProjectChannels = ref<Map<string, VersionSlim>>(new Map());

const idle = useIdleStore();

const projectModel = computed(() => projectStore.currentProject!);
const filteredProjects = computed(() => {
  const projects = approvableInfo.value.projects ?? [];
  const normalizedSearch = (search.value ?? '').trim().toLowerCase();

  if (!normalizedSearch) {
    return projects;
  }

  return projects.filter((project) => {
    return [project.projectName, project.supplier].some((value) => value?.toLowerCase().includes(normalizedSearch));
  });
});

async function reload() {
  dataAreLoaded.value = false;
  idle.showIdle = true;

  approvableInfo.value = await projectService.getApprovableInfo(projectModel.value._key, true);

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
  dataAreLoaded.value = true;
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
      <div ref="tableUserManagement" class="flex h-full flex-col overflow-hidden">
        <div class="mb-2 flex flex-shrink-0 flex-wrap items-center gap-2">
          <v-spacer></v-spacer>
          <DSearchField v-model="search" class="ml-auto" />
        </div>
        <GridSPDXList
          class="flex-1 overflow-hidden"
          :projects="filteredProjects"
          :channels="projectModel.versions"
          showSbomExtras
          showSupplier
          showLoading
          :loading="!dataAreLoaded" />
      </div>
    </template>
  </TableLayout>
</template>
