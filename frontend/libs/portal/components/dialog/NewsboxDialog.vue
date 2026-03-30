// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import {useAppStore} from '@disclosure-portal/stores/app';
import {useNewsboxStore} from '@disclosure-portal/stores/newsbox.store';
import {useUserStore} from '@disclosure-portal/stores/user';
import {storeToRefs} from 'pinia';
import {computed, ref, watch} from 'vue';

const appStore = useAppStore();
const newsboxStore = useNewsboxStore();
const userStore = useUserStore();

const {showNewsbox} = storeToRefs(newsboxStore);

const currentSlide = ref(0);
const lastSeenIndex = ref(-1);
const newsboxItems = computed(() => newsboxStore.newsItems?.items || []);

const currentTitle = computed(() => {
  const item = newsboxItems.value[currentSlide.value];
  if (!item) return '';
  return appStore.getAppLanguage === 'de' && item.titleDE ? item.titleDE : item.title;
});

const currentDescription = computed(() => {
  const item = newsboxItems.value[currentSlide.value];
  if (!item) return '';
  return appStore.getAppLanguage === 'de' && item.descriptionDE ? item.descriptionDE : item.description;
});

watch(
  () => newsboxStore.showNewsbox,
  (isOpen) => {
    if (isOpen) {
      newsboxStore.fetchItems();
    } else {
      currentSlide.value = 0;
    }
  },
  {immediate: true},
);

watch(
  () => newsboxStore.newsItems,
  (response) => {
    if (response && response.items.length > 0) {
      if (response?.toShow >= 0 && response?.toShow < response?.items?.length) {
        currentSlide.value = response.toShow;
        lastSeenIndex.value = response.toShow;
      } else {
        currentSlide.value = 0;
        lastSeenIndex.value = -1;
      }
    }
  },
  {immediate: true},
);

const close = () => {
  if (currentSlide.value >= lastSeenIndex.value) {
    lastSeenIndex.value = currentSlide.value;
    updateLastSeenForCurrentItem();
  }
  showNewsbox.value = false;
};

const goToNext = () => {
  if (currentSlide.value < newsboxItems.value.length - 1) {
    currentSlide.value++;
    if (currentSlide.value > lastSeenIndex.value) {
      lastSeenIndex.value = currentSlide.value;
      updateLastSeenForCurrentItem();
    }
  }
};

const goToPrevious = () => {
  if (currentSlide.value > 0) {
    currentSlide.value--;
  }
};

const updateLastSeenForCurrentItem = () => {
  const currentItem = newsboxItems.value[currentSlide.value];
  if (currentItem && currentItem._key) {
    const userLastSeenDto = {newsboxLastSeenId: currentItem._key};
    newsboxStore.updateLastSeen(userStore.getProfile._key, userLastSeenDto);
  }
};

const canGoNext = computed(() => newsboxItems.value.length > 1 && currentSlide.value < newsboxItems.value.length - 1);
const canGoPrevious = computed(() => newsboxItems.value.length > 1 && currentSlide.value > 0);
</script>

<template>
  <v-dialog v-model="newsboxStore.showNewsbox" content-class="large" scrollable width="900" height="800">
    <DialogLayout :config="{title: currentTitle}" @close="close" @primary-action="close">
      <Stack>
        <v-progress-circular
          v-if="newsboxStore.loading"
          indeterminate
          color="primary"
          class="mx-auto"></v-progress-circular>

        <Stack v-else-if="newsboxItems.length > 0" class="@container/newsbox h-full">
          <p class="text-body-1">
            {{ currentDescription }}
          </p>
          <DExternalLink
            v-if="newsboxItems[currentSlide]?.link"
            :url="newsboxItems[currentSlide]?.link || ''"
            text="Read More" />
          <div v-if="newsboxItems[currentSlide]?.image" class="flex justify-start items-start flex-1">
            <img
              :src="newsboxItems[currentSlide].image!"
              :alt="currentTitle"
              class="max-w-full object-contain h-auto @[400px]/newsbox:max-h-[500px] @[0px]/newsbox:max-h-[400px] rounded" />
          </div>
        </Stack>

        <div v-else-if="!newsboxStore.loading && newsboxItems.length === 0" class="text-center pa-4">
          <v-icon size="48" color="grey">mdi-newspaper-variant-outline</v-icon>
          <p class="text-body-1 mt-2">No news items available</p>
        </div>
      </Stack>
      <template #left>
        <div
          v-if="!newsboxStore.loading && newsboxItems.length > 0"
          class="newsbox-pagination flex-shrink-0 d-flex justify-center align-center py-3">
          <v-btn icon variant="text" size="medium" :disabled="!canGoPrevious" @click="goToPrevious" class="me-2">
            <v-icon>mdi-chevron-left</v-icon>
          </v-btn>

          <div class="d-flex align-center mx-1">
            <v-btn
              v-for="(_, index) in newsboxItems"
              :key="index"
              :variant="index === currentSlide ? 'flat' : 'flat'"
              :color="index === currentSlide ? 'primary' : 'grey-lighten-1'"
              size="x-small"
              disabled
              class="ma-1 rounded-full p-0 !rounded-full min-w-2.5 size-2.5">
            </v-btn>
          </div>

          <v-btn icon variant="text" size="medium" :disabled="!canGoNext" @click="goToNext" class="ms-2">
            <v-icon>mdi-chevron-right</v-icon>
          </v-btn>
        </div>
      </template>
    </DialogLayout>
  </v-dialog>
</template>
