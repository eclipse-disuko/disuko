<script lang="ts" setup>
import Stack from '@shared/layouts/Stack.vue';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import config from '@shared/utils/config';
import {computed, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const {dashboardCrumbs, ...breadcrumbs} = useBreadcrumbsStore();

interface FeatureFlag {
  name: string;
  enabled: boolean;
  originalValue: boolean;
}

const STORAGE_KEY = 'featureFlagOverrides';
const ORIGINAL_VALUES_KEY = 'featureFlagOriginals';
const EXCLUDED_KEYS = ['SERVER_URL', 'PUBLIC_API_ENDPOINT', 'OAUTH'];

const getOrInitializeOriginalValues = (): Record<string, boolean> => {
  const stored = sessionStorage.getItem(ORIGINAL_VALUES_KEY);
  if (stored) return JSON.parse(stored);

  const originalValues = Object.fromEntries(
    Object.entries(config).filter(([key, value]) => typeof value === 'boolean' && !EXCLUDED_KEYS.includes(key)),
  ) as Record<string, boolean>;

  sessionStorage.setItem(ORIGINAL_VALUES_KEY, JSON.stringify(originalValues));
  return originalValues;
};

const ORIGINAL_CONFIG_VALUES = getOrInitializeOriginalValues();

const loadOverrides = (): Record<string, boolean> => JSON.parse(sessionStorage.getItem(STORAGE_KEY) || '{}');

const saveOverrides = (overrides: Record<string, boolean>) =>
  sessionStorage.setItem(STORAGE_KEY, JSON.stringify(overrides));

const featureFlagOverrides = ref<Record<string, boolean>>(loadOverrides());

const initializeFlags = (): FeatureFlag[] => {
  const overrides = loadOverrides();

  return Object.entries(ORIGINAL_CONFIG_VALUES).map(([key, originalValue]) => {
    const enabled = overrides[key] ?? originalValue;
    (config as any)[key] = enabled;

    return {name: key, enabled, originalValue};
  });
};

const featureFlags = ref<FeatureFlag[]>(initializeFlags());

const updateConfig = (key: string, value: boolean) => {
  (config as any)[key] = value;
};

const toggleFeatureFlag = (flag: FeatureFlag) => {
  updateConfig(flag.name, flag.enabled);
  featureFlagOverrides.value[flag.name] = flag.enabled;
  saveOverrides(featureFlagOverrides.value);
};

const resetFlag = (flag: FeatureFlag) => {
  flag.enabled = flag.originalValue;
  updateConfig(flag.name, flag.originalValue);
  delete featureFlagOverrides.value[flag.name];
  saveOverrides(featureFlagOverrides.value);
};

const resetAll = () => {
  featureFlags.value.forEach((flag) => {
    flag.enabled = flag.originalValue;
    updateConfig(flag.name, flag.originalValue);
  });
  featureFlagOverrides.value = {};
  sessionStorage.removeItem(STORAGE_KEY);
  sessionStorage.removeItem(ORIGINAL_VALUES_KEY);
};

const hasOverrides = computed(() => featureFlags.value.some((flag) => flag.enabled !== flag.originalValue));

const initBreadcrumbs = () => {
  breadcrumbs.setCurrentBreadcrumbs([
    ...dashboardCrumbs,
    {
      title: t('FEATURE_FLAGS'),
    },
  ]);
};

onMounted(() => {
  featureFlagOverrides.value = loadOverrides();
  featureFlags.value = initializeFlags();
});

initBreadcrumbs();
</script>

<template>
  <div class="p-6">
    <Stack class="gap-6">
      <Stack direction="row" justify="between" align="center">
        <Stack>
          <h1 class="text-3xl font-semibold">{{ t('FEATURE_FLAGS') }}</h1>
          <p class="text-base text-gray-600">
            {{ t('FEATURE_FLAGS_DESCRIPTION') }}
          </p>
        </Stack>
        <v-btn v-if="hasOverrides" color="warning" variant="outlined" prepend-icon="mdi-restore" @click="resetAll">
          {{ t('RESET_ALL') }}
        </v-btn>
      </Stack>

      <v-alert type="info" variant="tonal">
        <Stack direction="row" align="center">
          <div>
            <strong>{{ t('DEBUG_MODE') }}</strong> - {{ t('FEATURE_FLAGS_SESSION_INFO') }}
          </div>
        </Stack>
      </v-alert>
    </Stack>

    <Stack class="mt-6">
      <v-card
        v-for="flag in featureFlags"
        :key="flag.name"
        :class="{
          'bg-green-500/70': flag.enabled,
          'border-2 border-yellow-500': flag.enabled !== flag.originalValue,
        }">
        <v-card-title>
          <Stack direction="row" justify="between">
            <Stack direction="row" align="center">
              <span class="text-xl font-medium">{{ flag.name }}</span>
              <v-chip v-if="flag.enabled !== flag.originalValue" color="warning" size="small" variant="outlined">
                {{ t('OVERRIDDEN') }}
              </v-chip>
            </Stack>
            <Stack direction="row" align="center">
              <v-btn
                v-if="flag.enabled !== flag.originalValue"
                size="small"
                variant="text"
                color="warning"
                @click="resetFlag(flag)">
                {{ t('RESET') }}
              </v-btn>
              <v-switch
                v-model="flag.enabled"
                color="primary"
                hide-details
                @update:model-value="toggleFeatureFlag(flag)">
              </v-switch>
            </Stack>
          </Stack>
        </v-card-title>
      </v-card>
    </Stack>
  </div>
</template>
