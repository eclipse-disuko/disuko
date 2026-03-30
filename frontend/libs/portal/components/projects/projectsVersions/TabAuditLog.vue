// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import GridAuditLog from '@disclosure-portal/components/grids/GridAuditLog.vue';
import ProjectService from '@disclosure-portal/services/projects';
import VersionService from '@disclosure-portal/services/version';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useSbomStore} from '@disclosure-portal/stores/sbom.store';
import {computed} from 'vue';

const sbomStore = useSbomStore();

const currentProject = computed(() => useProjectStore().currentProject!);
const currentVersionId = computed(() => sbomStore.getCurrentVersion._key);
const currentProjectId = computed(() => currentProject.value._key!);

const fetchMethod = () => {
  if (currentProjectId.value && currentVersionId.value) {
    return VersionService.getAuditTrail(currentProjectId.value, currentVersionId.value);
  } else {
    return ProjectService.getAuditTrail(currentProjectId.value);
  }
};
</script>

<template>
  <GridAuditLog :fetch-method="fetchMethod" />
</template>
