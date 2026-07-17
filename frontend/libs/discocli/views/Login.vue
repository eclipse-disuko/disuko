<template>
  <section class="bg-background flex min-h-[calc(100vh-96px)] items-center justify-center px-4">
    <div class="w-full max-w-[460px]">
      <h1
        class="text-[3rem] leading-[1.05] font-semibold tracking-[0.01em] whitespace-nowrap text-[rgb(var(--v-theme-on-background))]">
        Supplier Portal
      </h1>
      <div class="mt-2 w-fit">
        <h2 class="text-primary font-serif text-[2.5rem] leading-[0.92] font-semibold tracking-[0.02em] uppercase">
          {{ t('LOGIN') }}
        </h2>
        <div class="mt-1 h-[2px] w-[60%] bg-[rgb(var(--v-theme-primary))]"></div>
      </div>

      <v-form ref="form" @submit.prevent="handleLogin" autocomplete="off" class="mt-8 w-[280px] space-y-3">
        <v-text-field
          class="required"
          v-model="credentials.projectUuid"
          :label="`${t('PROJECT_UUID') || 'Project UUID'}`"
          :rules="projectUuidRules"
          :error-messages="projectUuidError"
          :disabled="loading"
          persistent-hint
          hint=" "
          required
          hide-details="auto"
          variant="outlined"
          autocomplete="off"
          spellcheck="false" />

        <v-text-field
          class="required"
          v-model="credentials.token"
          :label="`${t('TOKEN') || 'Token'}`"
          :type="showToken ? 'text' : 'password'"
          :rules="tokenRules"
          :error-messages="tokenError"
          :disabled="loading"
          persistent-hint
          hint=" "
          required
          hide-details="auto"
          variant="outlined"
          autocomplete="off"
          spellcheck="false"
          :append-inner-icon="showToken ? 'mdi-eye-off' : 'mdi-eye'"
          @click:append-inner="showToken = !showToken" />

        <DCActionButton
          :text="(t('LOGIN') || 'Login').toUpperCase()"
          variant="flat"
          block
          size="default"
          :loading="loading"
          :disabled="loading"
          @clicked="handleLogin"
          class="!mt-8 h-[34px]" />
      </v-form>
    </div>
  </section>
</template>

<script setup lang="ts">
import {authService} from '@cli/services/authService';
import {useAppStore} from '@cli/stores/app';
import {getRule} from '@cli/utils/validationRules';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import {ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRouter} from 'vue-router';
import type {VForm} from 'vuetify/components';

const router = useRouter();
const appStore = useAppStore();
const {t} = useI18n();

const form = ref<InstanceType<typeof VForm> | null>(null);
const showToken = ref(false);
const loading = ref(false);
const error = ref<string | null>(null);
const projectUuidError = ref('');
const tokenError = ref('');

const credentials = ref({
  projectUuid: '',
  token: '',
});

const projectUuidRules = [getRule(t, 'required'), getRule(t, 'validUUID')];
const tokenRules = [getRule(t, 'required')];

const redirectAfterLogin = async (projectUuid?: string) => {
  const storedUuid = projectUuid || authService.getAuth() || appStore.authProjectUUID;
  if (storedUuid) {
    await router.push({name: 'ProjectDetails', params: {id: storedUuid}});
    return;
  }

  await router.push('/');
};

async function validateForm(): Promise<boolean> {
  const isValid = await form.value?.validate();
  return isValid?.valid ?? false;
}

function clearErrors(): void {
  error.value = null;
  projectUuidError.value = '';
  tokenError.value = '';
}

async function handleLogin() {
  const isFormValid = await validateForm();
  if (!isFormValid) return;

  loading.value = true;
  clearErrors();

  try {
    const authProjectUuid = await authService.login(
      credentials.value.projectUuid.trim(),
      credentials.value.token.trim(),
    );
    appStore.setAuth(authProjectUuid);
    await redirectAfterLogin(authProjectUuid);
  } catch (e: any) {
    const errorCode = e?.response?.data?.code || '';

    if (e?.response?.status === 401 || errorCode === 'DISCOTOKEN_UNAUTHORIZED') {
      tokenError.value = t('ERROR_INVALID_TOKEN');
      error.value = t('ERROR_INVALID_TOKEN');
    } else if (errorCode === 'ERROR_DB_NOT_FOUND') {
      projectUuidError.value = t('ERROR_INVALID_PROJECT_UUID');
      error.value = t('ERROR_INVALID_PROJECT_UUID');
    } else {
      error.value = e?.message || t('LOGIN_FAILED') || 'Login failed';
    }
  } finally {
    loading.value = false;
  }
}

if (authService.isAuthenticated()) {
  redirectAfterLogin();
}
</script>