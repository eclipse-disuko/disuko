<script setup lang="ts">
import {useThemeStore} from '@shared/stores/theme.store';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {useLanguageStore} from '@shared/stores/language.store';

const {toggleLanguage} = useLanguageStore();
const themeStore = useThemeStore();
const {t} = useI18n();

const themeLabel = computed(() => (themeStore.current === 'light' ? t('BTN_THEME_DARK') : t('BTN_THEME_LIGHT')));

const settingsMenu = ref(false);
</script>

<template>
  <v-menu
    v-model="settingsMenu"
    transition="slide-y-transition"
    :close-on-content-click="false"
    location="bottom"
    offset="4">
    <template #activator="{props}">
      <v-btn
        v-bind="props"
        variant="outlined"
        prepend-icon="mdi-cog-outline"
        class="border-sm !border-current/35 !text-current/70 hover:!text-current/90"
        :aria-label="t('BTN_SETTINGS')"
        :title="t('BTN_SETTINGS')">
        {{ t('BTN_SETTINGS') }}
      </v-btn>
    </template>
    <v-list class="pa-0">
      <v-list-item class="h-[56px] px-4" @click="themeStore.toggle" prepend-icon="mdi-theme-light-dark">
        <template #prepend>
          <v-icon color="primary"></v-icon>
        </template>
        {{ themeLabel }}
      </v-list-item>
      <v-list-item class="h-[56px] px-4" @click="toggleLanguage" prepend-icon="mdi-web">
        <template #prepend>
          <v-icon color="primary"></v-icon>
        </template>
        {{ t('BTN_LANGUAGE_SWITCH') }}
      </v-list-item>
    </v-list>
  </v-menu>
</template>
