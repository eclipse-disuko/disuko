<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';

const {dashboardCrumbs, projectsCrumb, ...breadcrumbs} = useBreadcrumbsStore();

breadcrumbs.setCurrentBreadcrumbs([...dashboardCrumbs, projectsCrumb]);
</script>

<template>
  <AdminGridProjects></AdminGridProjects>
</template>
