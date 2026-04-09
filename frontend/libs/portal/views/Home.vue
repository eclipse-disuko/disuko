<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import ProfileService from '@disclosure-portal/services/profile';
import {useAppStore} from '@disclosure-portal/stores/app';
import {useNewsboxStore} from '@disclosure-portal/stores/newsbox.store';
import {useUserStore} from '@disclosure-portal/stores/user';
import {useWizardStore} from '@disclosure-portal/stores/wizard.store';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';

import {computed, onBeforeMount, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const userStore = useUserStore();
const appStore = useAppStore();
const {t} = useI18n();
const wizardStore = useWizardStore();
const newsboxStore = useNewsboxStore();
const {dashboard, ...breadcrumbs} = useBreadcrumbsStore();

const forename = ref(userStore.getProfile.forename);
const rights = ref(userStore.getRights);
const shouldShowNewsbox = computed(() => newsboxStore.showNewsbox);

const getDashboardCounts = async () => {
  try {
    const counts = await ProfileService.getDashboardCounts();
    appStore.updateTileCounts(counts);
    return counts;
  } catch (error) {
    console.error('Error getting dashboard counts: ', error);
    return {activeJobCount: 0, hasNewNewsboxItem: false};
  }
};

onBeforeMount(async () => {
  breadcrumbs.setCurrentBreadcrumbs([dashboard]);
  const counts = await getDashboardCounts();
  newsboxStore.hasNewNewsboxItem = counts.hasNewNewsboxItem;
});

const onDidYouKnowClick = () => {
  newsboxStore.showNewsbox = true;
};
</script>

<template>
  <v-main class="p-10">
    <div role="heading" class="text-h4 font-weight-light">
      {{ t('WELCOME', {forename: forename}) }}
    </div>
    <div class="flex flex-col md:flex-row gap-8 mt-10">
      <Stack v-if="rights.allowProject && rights.allowProject.create">
        <HomeTile
          type="action"
          :title="t('NP_DIALOG_TITLE')"
          icon="mdi-playlist-plus"
          :description="t('HOME_EXPLAIN_NEW_PROJECT')"
          @click="wizardStore.openWizard()" />
        <HomeTile
          type="action"
          :title="t('Title_New_Group')"
          icon="mdi-plus-box-multiple-outline"
          :description="t('HOME_EXPLAIN_NEW_GROUP')"
          @click="wizardStore.openWizard({isGroup: true})" />
      </Stack>
      <Stack>
        <HomeTile
          v-for="(module, index) in appStore.tiles"
          :key="index"
          type="navigation"
          :title="t(module.title)"
          icon="mdi-chevron-right"
          :url="module.url"
          :show-badge="module.cnt > -1"
          :badge-content="module.cnt" />

        <HomeTile
          type="navigation"
          :title="t('DID_YOU_KNOW')"
          icon="mdi-chevron-right"
          @click="onDidYouKnowClick"
          :show-badge="newsboxStore.hasNewNewsboxItem"
          badge-content="new" />
      </Stack>
    </div>
    <NewsboxDialog v-if="shouldShowNewsbox"></NewsboxDialog>
  </v-main>
</template>
