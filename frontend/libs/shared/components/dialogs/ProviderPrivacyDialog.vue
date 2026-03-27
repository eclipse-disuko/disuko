// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script lang="ts" setup>
import {useAppStore} from '@disclosure-portal/stores/app';
import LegalNoticeDe from '@shared/assets/documents/legal_notice/LegalNotice_de.md?raw';
import LegalNoticeEn from '@shared/assets/documents/legal_notice/LegalNotice_en.md?raw';
import PrivacyStatementDe from '@shared/assets/documents/privacy_statement/PrivacyStatement_de.md?raw';
import PrivacyStatementEn from '@shared/assets/documents/privacy_statement/PrivacyStatement_en.md?raw';
import ProviderDe from '@shared/assets/documents/provider/Provider_de.md?raw';
import ProviderEn from '@shared/assets/documents/provider/Provider_en.md?raw';
import NoticeDe from '@shared/assets/documents/provider_privacy_notice/ProviderPrivacyNotice_de.html?raw';
import NoticeEn from '@shared/assets/documents/provider_privacy_notice/ProviderPrivacyNotice_en.html?raw';
import TermsOfUseEn from '@shared/assets/documents/terms_of_use/TermsOfUseCurrent.md?raw';
import TermsOfUseDe from '@shared/assets/documents/terms_of_use/TermsOfUseDe.md?raw';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';

interface ITabs {
  key: string;
  name: string;
  content: string;
}

const isDialogVisible = ref(false);
const selectedTabIndex = ref(0);

const {t} = useI18n();
const appStore = useAppStore();

const tabs = computed<ITabs[]>(() => [
  {
    key: 'provider',
    name: 'TAB_PROVIDER',
    content: appStore.getAppLanguage === 'en' ? ProviderEn : ProviderDe,
  },
  {
    key: 'legalNotices',
    name: 'TAB_LEGAL_NOTICES',
    content: appStore.getAppLanguage === 'en' ? LegalNoticeEn : LegalNoticeDe,
  },
  {
    key: 'privacyStatement',
    name: 'TAB_PRIVACY_STATEMENT',
    content: appStore.getAppLanguage === 'en' ? PrivacyStatementEn : PrivacyStatementDe,
  },
  {
    key: 'termsOfUse',
    name: 'TAB_TERMS_OF_USE',
    content: appStore.getAppLanguage === 'en' ? TermsOfUseEn : TermsOfUseDe,
  },
  {
    key: 'notice',
    name: 'TAB_NOTICE',
    content: appStore.getAppLanguage === 'en' ? NoticeEn : NoticeDe,
  },
]);

const showDialog = () => {
  selectedTabIndex.value = 0;
  isDialogVisible.value = true;
};

const hideDialog = () => {
  isDialogVisible.value = false;
};

const clipboardContent = (key: string) => {
  if (key == 'termsOfUse') {
    const elements = document.getElementsByClassName('markdown');
    if (!elements) {
      return '';
    }
    const md = elements[0] as HTMLElement;
    return md ? md.innerText : '';
  }
  const element = document.getElementById('providerStatement' + key);
  return element ? element.innerText : '';
};
</script>

<template>
  <div>
    <slot :showDialog="showDialog"> </slot>

    <v-dialog v-model="isDialogVisible" content-class="large" max-width="1045px">
      <v-card flat class="dDialog pa-8">
        <v-card-title>
          <v-row>
            <v-col cols="10">
              <span class="text-h5">{{ t('PPS') }}</span>
            </v-col>
            <v-col cols="2" align="right">
              <DCloseButton @click="hideDialog" />
            </v-col>
          </v-row>
        </v-card-title>
        <v-card-text>
          <v-tabs
            v-model="selectedTabIndex"
            slider-color="mbti"
            active-class="active"
            show-arrows
            bg-color="tabsHeader"
          >
            <v-tab v-for="tab in tabs" :key="tab.key">
              {{ t(tab.name) }}
            </v-tab>
          </v-tabs>
          <v-tabs-window v-model="selectedTabIndex" class="max-h-[400px] overflow-y-scroll">
            <v-tabs-window-item v-for="tab in tabs" :key="tab.key">
              <v-row class="ml-1 pt-4 pr-2">
                <v-col>
                  <Markdown :text="tab.content" :id="'providerStatement' + tab.key"></Markdown>
                </v-col>
                <v-col class="d-flex justify-end ml-2" cols="1">
                  <DCopyClipboardButton :hint="t('TT_CopyText')" :content="clipboardContent(tab.key)" />
                </v-col>
              </v-row>
            </v-tabs-window-item>
          </v-tabs-window>
        </v-card-text>
      </v-card>
    </v-dialog>
  </div>
</template>
