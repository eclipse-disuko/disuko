<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {storeToRefs} from 'pinia';
import {computed} from 'vue';
import {InternalBreadcrumbItem} from 'vuetify/lib/components/VBreadcrumbs/VBreadcrumbs';

const breadcrumbsStore = useBreadcrumbsStore();
const {currentBreadcrumbs} = storeToRefs(breadcrumbsStore);

const breadcrumbs = computed<InternalBreadcrumbItem[]>(() =>
  currentBreadcrumbs.value.map((crumb: InternalBreadcrumbItem, index: number) => {
    const newCrumb: InternalBreadcrumbItem = {...crumb};
    newCrumb.to = crumb.href ?? '';
    if (index === currentBreadcrumbs.value.length - 1) {
      newCrumb.disabled = true;
    }
    return newCrumb;
  }),
);
</script>

<template>
  <v-breadcrumbs :items="breadcrumbs" color="grey-darken-1">
    <template v-slot:divider>
      <v-icon icon="mdi-chevron-right" size="x-small" color="grey-darken-1"></v-icon>
    </template>
  </v-breadcrumbs>
</template>

<style>
.v-breadcrumbs-item--disabled {
  color: #757575;
  opacity: 1;
}
</style>
