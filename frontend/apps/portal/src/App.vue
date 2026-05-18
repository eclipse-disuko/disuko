<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {usePageTitle} from '@disclosure-portal/composables/usePageTitle';
import DHTTPError from '@shared/types/DHTTPError';
import ErrorDialogConfig from '@shared/types/ErrorDialogConfig';
import profileService from '@disclosure-portal/services/profile';
import {useAppStore} from '@disclosure-portal/stores/app';
import {useCustomIdStore} from '@disclosure-portal/stores/customid.store';
import {useLabelStore} from '@disclosure-portal/stores/label.store';
import {createNavItemsGroup, useUserStore} from '@disclosure-portal/stores/user';
import {useWizardStore} from '@disclosure-portal/stores/wizard.store';
import eventBus from '@shared/utils/eventbus';
import config from '@shared/utils/config';
import {onMounted, onUnmounted, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRoute, useRouter} from 'vue-router';
import {useLanguageStore} from '@shared/stores/language.store';
import {storeToRefs} from 'pinia';
import {useEventKeysStore} from '@shared/stores/eventKeys.store';

const {t} = useI18n();
const route = useRoute();
const userStore = useUserStore();
const appStore = useAppStore();
const languageStore = useLanguageStore();
const {appLanguage} = storeToRefs(languageStore);
const customIdsStore = useCustomIdStore();
const router = useRouter();
const wizardStore = useWizardStore();
const labelStore = useLabelStore();
const {useReactiveTitle} = usePageTitle();
const eventKeyStore = useEventKeysStore();

const hasAuthentication = ref(false);
const profileActive = ref(false);
const showDlgTos = ref(false);
const errorDialog = ref();
const dud = ref();

const backendCodesRedirectToProjectList: Set<string> = new Set([
  'ERROR_REPOSITORY_READ',
  'FIND_VERSION',
  'VERSION_DELETED',
  'PARAM_UUID_WRONG',
  'PARAM_VERSION_WRONG',
  'PARAM_VERSION_EMPTY',
  'AAR',
  'TASK_NOT_FOUND',
]);

const backendCodesRedirectToLicenseList: Set<string> = new Set(['LICENSE_DATA_MISSING']);

const login = () => {
  const url = config.SERVER_URL + config.OAUTH.LOGIN;
  window.location.replace(url);
};

const onCloseErrorDialog = (error: ErrorDialogConfig) => {
  if (backendCodesRedirectToProjectList.has(error.titleKeyOrCode)) {
    router.push('/dashboard/home');
  }

  if (backendCodesRedirectToLicenseList.has(error.titleKeyOrCode)) {
    router.push('/dashboard/licenses');
  }

  if (error.errorCode === 'ERROR_401' || error.errorCode === '401' || error.errorCode === 'UNAUTHORIZED') {
    login();
  }
};

const showError = ({error}: {error: ErrorDialogConfig}) => {
  if (!errorDialog.value) {
    return;
  }

  errorDialog.value.open(error);
};

const showAPIError = (error: DHTTPError) => {
  if (!errorDialog.value) {
    return;
  }

  const d = new ErrorDialogConfig();

  d.title = '' + t(error.title);
  d.titleKeyOrCode = error.title;
  d.description = '' + t(error.message);
  d.stackTrace = error.raw;
  d.errorCode = error.code;
  d.reqId = t(error.reqId);

  if (error.code === 'UNAUTHORIZED' || error.code === '401') {
    userStore.clear();
    login();
    return;
  }

  if (error.code === '403' && error.title === 'USER_DISABLED') {
    // TODO: Add Dialog here to inform user about status
    return;
  }

  errorDialog.value.open(d);
};

const onAcceptTOS = () => {
  window.location.reload();
};

const onResizeWindow = () => {
  eventBus.emit('window-resize', {});
};

// Set up reactive page title based on route meta.title only
watch(
  () => [route.meta?.title, appLanguage.value],
  () => {
    if (route.meta?.title) {
      let title = 'Disclosure Portal';
      const titleObj = route.meta.title as {[key: string]: string};
      title = titleObj[appLanguage.value];

      useReactiveTitle(title);
    }
  },
  {immediate: true},
);

onUnmounted(() => {
  window.removeEventListener('resize', onResizeWindow);
});

onMounted(async () => {
  eventBus.on('on-api-error', showAPIError);
  eventBus.on('on-error', showError);
  window.addEventListener('resize', onResizeWindow);

  languageStore.initializeLanguage();
  eventKeyStore.initEventKeyStore();

  const simpleProfileData = await profileService.getProfileData();

  await appStore.fetchLabelsTools();
  await labelStore.fetchAllLabels();
  await customIdsStore.updateCustomIds();

  hasAuthentication.value = true;

  userStore.setSimpleProfileData(simpleProfileData);

  if (!simpleProfileData.profile!.active) {
    dud.value?.open();
    return;
  }

  profileActive.value = true;

  if (!simpleProfileData.profile!.termsOfUse) {
    showDlgTos.value = true;
    return;
  }

  appStore.startTokenRefresher();
  createNavItemsGroup();
});
</script>

<template>
  <v-app>
    <Idle></Idle>
    <router-view v-if="hasAuthentication && profileActive"></router-view>
    <DisabledUserDialog ref="dud"></DisabledUserDialog>
    <ErrorDialog ref="errorDialog" @close="onCloseErrorDialog"></ErrorDialog>
    <TermsOfUseDialog v-model="showDlgTos" @success="onAcceptTOS"></TermsOfUseDialog>
    <NewWizardDialog v-if="wizardStore.isWizardOpen"></NewWizardDialog>
  </v-app>
</template>
