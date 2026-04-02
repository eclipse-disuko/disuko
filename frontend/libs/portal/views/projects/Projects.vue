<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {onMounted} from 'vue';

const {dashboardCrumbs, projectsCrumb, ...breadcrumbs} = useBreadcrumbsStore();

onMounted(() => {
  breadcrumbs.setCurrentBreadcrumbs([...dashboardCrumbs, projectsCrumb]);
});
</script>

<template>
  <GridProjects></GridProjects>
</template>
