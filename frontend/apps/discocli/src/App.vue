<script setup lang="ts">
import AppLayout from '@cli/layouts/AppLayout.vue';
import AuthLayout from '@cli/layouts/AuthLayout.vue';
import {computed, onMounted} from 'vue';
import {useRoute} from 'vue-router';
import {useLanguageStore} from '@shared/stores/language.store';

const languageStore = useLanguageStore();
const route = useRoute();

const layoutComponent = computed(() => {
  return route.meta.layout === 'auth' ? AuthLayout : AppLayout;
});

onMounted(() => {
  languageStore.initializeLanguage();
});
</script>

<template>
  <component v-if="route.name" :is="layoutComponent">
    <router-view></router-view>
  </component>
</template>
