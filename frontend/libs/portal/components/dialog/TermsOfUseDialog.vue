<script setup lang="ts">
import {UserRequestDto} from '@disclosure-portal/model/Users';
import ProfileServer from '@disclosure-portal/services/profile';
import {useAppStore} from '@disclosure-portal/stores/app';
import {useUserStore} from '@disclosure-portal/stores/user';
import termsOfUseEn from '@shared/assets/documents/terms_of_use/TermsOfUseCurrent.md?raw';
import termsOfUseDe from '@shared/assets/documents/terms_of_use/TermsOfUseDe.md?raw';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import Markdown from '@shared/components/widgets/MarkdownWidget.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import {useClipboard} from '@shared/utils/clipboard';
import {storeToRefs} from 'pinia';
import {ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';

const showDialog = defineModel<boolean>({required: true, default: false});

const emit = defineEmits(['success']);

const {info: infoSnackbar, error: errorSnackbar} = useSnackbar();
const userStore = useUserStore();
const appStore = useAppStore();
const {t} = useI18n();
const {copyToClipboard} = useClipboard();

const {appLanguage} = storeToRefs(appStore);

const hasAgreedToU = ref(false);
const langTab = ref<'de' | 'en'>(appLanguage.value);
const hasScrolledToBottom = ref(false);

const close = async () => {
  if (!hasAgreedToU.value) {
    return;
  }
  const response = await ProfileServer.update(userStore.getProfile._key, {
    termsOfUse: true,
  } as UserRequestDto);
  userStore.updateTermsOfUse(true, response.data.termsOfUseDate);
  infoSnackbar(t('WELCOME_REGISTERED_USER_SNACKBAR_TEXT'));
  showDialog.value = false;
  emit('success', response);
};

const copyProviderStatementToClipboard = () => {
  const content = langTab.value === 'en' ? termsOfUseEn : termsOfUseDe;
  copyToClipboard(content);
};

const onIntersect = (isIntersecting: boolean, entries: IntersectionObserverEntry[], observer: IntersectionObserver) => {
  if (isIntersecting) {
    hasScrolledToBottom.value = true;
  }
};

watch(langTab, (newTab) => {
  appStore.setLanguage(newTab);
});
</script>

<template>
  <v-dialog v-model="showDialog" max-width="1045px" persistent scrollable>
    <v-card class="pa-8 dDialog">
      <v-card-title>
        <v-row>
          <v-col cols="10">
            <span class="text-h5">{{ t('TAB_TERMS_OF_USE') }}</span>
          </v-col>
          <v-col cols="2" align="right">
            <DCActionButton
              class="ml-2"
              :tableButton="false"
              :hint="t('TT_CopyText')"
              icon="mdi-content-copy"
              variant="text"
              @click="copyProviderStatementToClipboard" />
          </v-col>
        </v-row>
      </v-card-title>
      <v-card-text>
        <v-tabs v-model="langTab" slider-color="mbti" active-class="active" show-arrows bg-color="tabsHeader">
          <v-tab value="en">English</v-tab>
          <v-tab value="de">Deutsch</v-tab>
        </v-tabs>
        <v-tabs-window v-model="langTab" class="max-h-[400px] overflow-y-scroll">
          <v-tabs-window-item value="en">
            <Markdown :text="termsOfUseEn" :id="'providerStatementTOU'">
              <span id="end-of-tou" v-intersect="onIntersect">&nbsp;</span>
            </Markdown>
          </v-tabs-window-item>
          <v-tabs-window-item value="de">
            <Markdown :text="termsOfUseDe" :id="'providerStatementTOUDE'">
              <span id="end-of-tou" v-intersect="onIntersect">&nbsp;</span>
            </Markdown>
          </v-tabs-window-item>
        </v-tabs-window>
      </v-card-text>
      <v-card-actions class="pt-2">
        <v-checkbox
          v-model="hasAgreedToU"
          :hint="!hasScrolledToBottom ? t('LABEL_CHECKBOX_TERMS_OF_USE_HINT') : ''"
          :persistent-hint="!hasScrolledToBottom"
          color="primary"
          :disabled="!hasScrolledToBottom"
          :label="t('LABEL_CHECKBOX_TERMS_OF_USE')"
          class="shrink mr-2 mt-0"></v-checkbox>
        <v-spacer></v-spacer>
        <DCActionButton
          isDialogButton
          size="small"
          variant="text"
          @click="close"
          :disabled="!hasAgreedToU"
          :text="t('BTN_AGREE')" />
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
