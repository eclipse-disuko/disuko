// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import * as help from '@disclosure-portal/assets/documents/help';
import releaseNotes from '@disclosure-portal/assets/documents/release_notes/ReleaseNotes.md?raw';
import HelpDialog from '@disclosure-portal/components/dialog/HelpDialog.vue';
import {INotificationMeta} from '@disclosure-portal/model/IdleInfo';
import {useAppStore} from '@disclosure-portal/stores/app';
import {createNavItemsGroup, useUserStore} from '@disclosure-portal/stores/user';
import eventbus from '@disclosure-portal/utils/eventbus';
import {logout} from '@disclosure-portal/utils/logout';
import {escapeHtml} from '@disclosure-portal/utils/Validation';
import {ThemeColor, useThemeStore} from '@shared/stores/theme.store';
import config from '@shared/utils/config';
import {storeToRefs} from 'pinia';
import {computed, onMounted, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRoute, useRouter} from 'vue-router';

const route = useRoute();
const router = useRouter();
const appStore = useAppStore();
const userStore = useUserStore();
const theme = useThemeStore();
const {t} = useI18n();

const {notificationClosed, notificationMessage, appLanguage, dummyDesignMode} = storeToRefs(appStore);

const username = ref('');
const navIsCollapsed = ref(true);
const helpOpen = ref(false);
const showSnow = ref(true);
const isHelpAnimating = ref(true);

const metaHelpText = computed(() => route.meta?.helpText);
const hasCustomHelp = computed(() => {
  // Only show + when there's specific context help (not just the default top section)
  return !!(metaHelpText.value && appLanguage.value in metaHelpText.value && metaHelpText.value[appLanguage.value]);
});
const helpText = computed(() => {
  const specificHelp = hasCustomHelp.value && metaHelpText.value ? '\n' + metaHelpText.value[appLanguage.value] : '';
  return help.top[appLanguage.value] + specificHelp;
});
const switchLabel = computed(() => {
  return theme.current === ThemeColor.dark ? t('SWITCH_LIGHT') : t('SWITCH_DARK');
});
const switchTooltip = computed(() => {
  return theme.current === ThemeColor.dark ? t('TT_SWITCH_LIGHT') : t('TT_SWITCH_DARK');
});
const version = computed(() => {
  return import.meta.env.VITE_VERSION;
});
const versionDate = computed(() => {
  return import.meta.env.VITE_BUILD_DATE + ' - INTERNAL DATA';
});

const collapseDrawer = () => {
  navIsCollapsed.value = true;
};

const openHelp = () => {
  helpOpen.value = true;
};

const disableNotification = () => {
  notificationClosed.value = true;
};

const onSetNotification = ({config}: {config: INotificationMeta}) => {
  if (notificationClosed.value) {
    return;
  }
  notificationClosed.value = !config.enabled;
  notificationMessage.value = config.text;
};

const escapeListener = (event: KeyboardEvent) => {
  if (event.key === 'Escape') {
    collapseDrawer();
  }
};

const toggleLanguage = () => {
  appStore.toggleLanguage();
  createNavItemsGroup();
};

const profile = () => {
  router.replace({path: '/dashboard/user'});
};

const helpDesk = () => {
  const url = t('URL_Helpdesk');
  window.open(url, '_blank');
};

const logoutUser = () => {
  userStore.clear();
  logout();
};

let animationTimeout: ReturnType<typeof setTimeout> | null = null;
const startHelpAnimation = () => {
  isHelpAnimating.value = true;
  if (animationTimeout) {
    clearTimeout(animationTimeout);
  }
  animationTimeout = setTimeout(() => {
    isHelpAnimating.value = false;
  }, 3700);
};

onMounted(() => {
  // initially set context help
  eventbus.on('set-notification', onSetNotification);
  window.addEventListener('keydown', escapeListener);

  username.value = userStore.getProfile.user;
  startHelpAnimation();
});
watch(
  () => route.path,
  () => {
    startHelpAnimation();
  },
);
</script>

<template>
  <v-system-bar
    class="pa-2 notification-bar"
    height="32"
    fixed
    app
    v-if="!appStore.notificationClosed && appStore.notificationMessage">
    <span class="font-bold" v-html="escapeHtml(appStore.notificationMessage)"></span>
    <v-spacer></v-spacer>
    <DCloseButton @click="disableNotification()"></DCloseButton>
  </v-system-bar>
  <v-navigation-drawer
    :rail="navIsCollapsed"
    permanent
    color="sidebar"
    :theme="config.useWinterTheme ? 'dark' : undefined">
    <div v-click-outside="collapseDrawer">
      <v-list-item
        @click="navIsCollapsed = !navIsCollapsed"
        :title="t('APP_NAME')"
        :prepend-icon="navIsCollapsed ? 'mdi-menu' : 'mdi-close'"
        class="pt-5 pb-3 navbar">
      </v-list-item>
      <v-list flat v-if="appStore.navItemGroup">
        <div v-for="item in appStore.navItemGroup.items" :key="item.path" class="mb-4">
          <NavItem :item="item" v-if="item.condition"></NavItem>
        </div>

        <template v-if="appStore.navItemGroup.adminItem && appStore.navItemGroup.adminItem.condition">
          <v-divider class="mx-2" />
          <NavItem :item="appStore.navItemGroup.adminItem"></NavItem>

          <div v-if="!navIsCollapsed">
            <div v-for="item in appStore.navItemGroup.adminItem.subItems" :key="item.path">
              <SubNavItem :item="item" v-if="item.condition"></SubNavItem>
            </div>
          </div>
        </template>
      </v-list>
    </div>
  </v-navigation-drawer>

  <v-app-bar
    color="sidebar"
    app
    class="navbar"
    :class="{
      'theme-dark': theme.current === ThemeColor.dark,
      'theme-light': theme.current === ThemeColor.light,
      'barrier-tape-background top': dummyDesignMode,
    }"
    height="56"
    id="disco-toolbar"
    :elevation="0"
    :theme="config.useWinterTheme ? 'dark' : undefined">
    <DBreadcrumb></DBreadcrumb>
    <v-spacer></v-spacer>
    <v-btn large color="button" variant="outlined" prepend-icon="mdi-help" class="mr-2 border-sm" @click="openHelp">
      <template v-slot:prepend>
        <span
          v-if="hasCustomHelp"
          class="rounded-full ml-2 size-2 bg-blue-45"
          :class="{'animate-pulse': isHelpAnimating}"></span>
        <v-icon v-else color="primary"></v-icon>
      </template>
      {{ t('HELP') }}
      <tooltip :text="t('HELP_SHOW')"></tooltip>
    </v-btn>
    <HelpDialog v-model="helpOpen" :title="t('HELP')" :text="helpText"></HelpDialog>

    <v-menu
      content-class="dropdown-menu dropdown-arrow"
      absolute
      transition="slide-y-transition"
      :position-y="56"
      offset-y
      :close-on-content-click="false">
      <template v-slot:activator="{props}">
        <v-btn
          large
          v-bind="props"
          class="mr-5 border-sm"
          color="button"
          variant="outlined"
          prepend-icon="mdi-account-circle-outline">
          <template v-slot:prepend>
            <v-icon color="primary"></v-icon>
          </template>

          {{ username }}
        </v-btn>
      </template>

      <v-list class="pa-0">
        <v-list-item class="h-[56px] px-4" @click="profile()" prepend-icon="mdi-account-circle-outline">
          <template v-slot:prepend>
            <v-icon color="primary"></v-icon>
          </template>
          {{ t('BTN_profile') }}
          <tooltip :text="t('TT_profile')"></tooltip>
        </v-list-item>
        <v-list-item class="h-[56px] px-4" @click="helpDesk()" prepend-icon="mdi-lifebuoy">
          <template v-slot:prepend>
            <v-icon color="primary"></v-icon>
          </template>
          {{ t('BTN_helpdesk') }}
          <tooltip :text="t('TT_helpdesk')"></tooltip>
        </v-list-item>
        <v-list-item class="h-[56px] px-4" @click="theme.toggle()" prepend-icon="mdi-theme-light-dark">
          <template v-slot:prepend>
            <v-icon color="primary"></v-icon>
          </template>
          {{ switchLabel }}
          <tooltip :text="switchTooltip"></tooltip>
        </v-list-item>
        <v-list-item class="h-[56px] px-4" @click="toggleLanguage()" prepend-icon="mdi-web">
          <template v-slot:prepend>
            <v-icon color="primary"></v-icon>
          </template>
          {{ t('BTN_LANGUAGE_SWITCH') }}
          <tooltip :text="t('TT_LANGUAGE_SWITCH')"></tooltip>
        </v-list-item>
        <v-list-item
          v-if="config.useWinterTheme"
          class="h-[56px] px-4"
          @click="showSnow = !showSnow"
          prepend-icon="mdi-snowflake">
          <template v-slot:prepend>
            <v-icon color="primary"></v-icon>
          </template>
          {{ t('TOGGLE_SNOW') }}
          <tooltip :text="t('TOGGLE_SNOW')"></tooltip>
        </v-list-item>
        <v-list-item class="h-[56px] px-4" @click="logoutUser()" prepend-icon="mdi-logout">
          <template v-slot:prepend>
            <v-icon color="primary"></v-icon>
          </template>
          {{ t('BTN_LOGOUT') }}
          <tooltip :text="t('TT_logout')"></tooltip>
        </v-list-item>
      </v-list>
    </v-menu>
    <Snowflake v-if="config.useWinterTheme && showSnow"></Snowflake>
  </v-app-bar>
  <v-main class="justify-start" id="disco-main">
    <div class="disco-full-height">
      <router-view></router-view>
    </div>
    <DSnackbar></DSnackbar>
  </v-main>
  <v-footer
    app
    id="disco-footer"
    class="gap-2"
    :class="{
      'theme-dark': theme.current === ThemeColor.dark,
      'theme-light': theme.current === ThemeColor.light,
      'barrier-tape-background bottom': dummyDesignMode,
    }">
    <v-spacer></v-spacer>
    <ProviderPrivacyDialog v-slot="{showDialog}">
      <span @click="showDialog" class="text-caption hover-underline cursor-pointer">
        {{ t('PPS') }}
      </span>
    </ProviderPrivacyDialog>
    <v-divider vertical></v-divider>
    <ReleaseNotesDialog :releaseNotes="releaseNotes" v-slot="{showDialog}">
      <span @click="showDialog" class="text-caption hover-underline cursor-pointer">
        <span v-if="version">{{ t('FT_VERSION') }}{{ `${version} ` }}</span>
        <span v-if="versionDate">{{ t('FT_VERSION_FROM') }} {{ versionDate }}</span>
      </span>
    </ReleaseNotesDialog>
  </v-footer>
</template>

<style scoped lang="scss">
.hover-underline:hover {
  text-decoration: underline;
}

.theme-light.barrier-tape-background.top {
  background-image:
    linear-gradient(to bottom, transparent 0%, rgb(var(--v-theme-surface)) 15px),
    repeating-linear-gradient(
      -45deg,
      transparent,
      transparent 7px,
      rgb(var(--v-theme-chartYellow)) 8px,
      rgb(var(--v-theme-chartYellow)) 12px
    );
}

.theme-light.barrier-tape-background.bottom {
  background-image:
    linear-gradient(to top, transparent 0%, rgb(var(--v-theme-surface)) 15px),
    repeating-linear-gradient(
      -45deg,
      transparent,
      transparent 7px,
      rgb(var(--v-theme-chartYellow)) 8px,
      rgb(var(--v-theme-chartYellow)) 12px
    );
}

.theme-dark.barrier-tape-background.top {
  background-image:
    linear-gradient(to bottom, transparent 0%, rgb(var(--v-theme-surface)) 10px),
    repeating-linear-gradient(
      -45deg,
      transparent,
      transparent 7px,
      rgb(var(--v-theme-chartYellow)) 8px,
      rgb(var(--v-theme-chartYellow)) 12px
    );
}

.theme-dark.barrier-tape-background.bottom {
  background-image:
    linear-gradient(to top, transparent 0%, rgb(var(--v-theme-surface)) 10px),
    repeating-linear-gradient(
      -45deg,
      transparent,
      transparent 7px,
      rgb(var(--v-theme-chartYellow)) 8px,
      rgb(var(--v-theme-chartYellow)) 12px
    );
}
</style>
