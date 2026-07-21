<script setup lang="ts">
import SettingsMenu from '@cli/components/SettingsMenu.vue';
import {useAppStore} from '@cli/stores/app';
import {computed, defineAsyncComponent} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRouter} from 'vue-router';

const {t} = useI18n();
const router = useRouter();
const appStore = useAppStore();

const isLoggedIn = computed(() => appStore.isAuthenticated);
const AsyncProviderPrivacyDialog = defineAsyncComponent(
  () => import('@shared/components/dialogs/ProviderPrivacyDialog.vue'),
);

const logout = async () => {
  await appStore.clearAuth();
  await router.push({name: 'Login'});
};
</script>

<template>
  <v-app>
    <v-app-bar color="primary" app class="navbar" height="56" id="disco-toolbar" :elevation="0">
      <div class="ml-4 flex items-center gap-2">
        <span class="text-subtitle-1 font-weight-bold text-[rgb(var(--v-theme-on-primary))]"> Supplier Portal</span>
        <v-divider vertical class="opacity-40"></v-divider>
        <div
          class="[&_.v-breadcrumbs]:min-h-0 [&_.v-breadcrumbs]:p-0 [&_.v-breadcrumbs-divider]:text-[rgba(var(--v-theme-on-primary),0.42)] [&_.v-breadcrumbs-item]:text-[rgba(var(--v-theme-on-primary),0.62)] [&_.v-breadcrumbs-item--disabled]:text-[rgba(var(--v-theme-on-primary),0.82)]">
          <DBreadcrumb />
        </div>
      </div>
      <v-spacer></v-spacer>
      <div class="ml-auto flex items-center">
        <v-btn
          v-if="isLoggedIn"
          variant="outlined"
          prepend-icon="mdi-logout"
          class="border-sm mr-2 !border-[rgba(var(--v-theme-on-primary),0.2)] !text-[rgba(var(--v-theme-on-primary),0.72)] hover:!text-[rgba(var(--v-theme-on-primary),0.92)]"
          :aria-label="t('BTN_LOGOUT')"
          :title="t('BTN_LOGOUT')"
          @click="logout">
          <template #prepend>
            <v-icon class="!text-[rgba(var(--v-theme-on-primary),0.72)]"></v-icon>
          </template>
          {{ t('BTN_LOGOUT') }}
        </v-btn>
        <SettingsMenu />
      </div>
    </v-app-bar>

    <v-main class="disco-full-height mt-[56px] justify-start" id="disco-main">
      <slot></slot>
    </v-main>

    <v-footer app height="40" class="gap-2 px-4" id="disco-footer">
      <v-spacer></v-spacer>
      <v-divider v-if="isLoggedIn" vertical></v-divider>
      <component :is="AsyncProviderPrivacyDialog" v-if="isLoggedIn" v-slot="{showDialog}">
        <span @click="showDialog" class="text-caption cursor-pointer hover:underline">
          {{ t('PPS') }}
        </span>
      </component>
    </v-footer>

    <DSnackbar />
  </v-app>
</template>

<style scoped>
.disco-full-height {
  min-height: calc(100vh - 96px);
}
</style>
