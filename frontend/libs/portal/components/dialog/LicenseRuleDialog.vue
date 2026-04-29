<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {DialogLicenseRuleConfig} from '@disclosure-portal/components/dialog/DialogConfigs';
import {ErrorDialogInterface} from '@disclosure-portal/components/dialog/DialogInterfaces';
import ErrorDialog from '@disclosure-portal/components/dialog/ErrorDialog.vue';
import ErrorDialogConfig from '@disclosure-portal/model/ErrorDialogConfig';
import {LicenseRuleRequest} from '@disclosure-portal/model/LicenseRule';
import {ComponentLicenses} from '@disclosure-portal/model/Project';
import {ComponentInfoSlim} from '@disclosure-portal/model/VersionDetails';
import projectService from '@disclosure-portal/services/projects';
import versionService from '@disclosure-portal/services/version';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useSbomStore} from '@disclosure-portal/stores/sbom.store';
import {useUserStore} from '@disclosure-portal/stores/user';
import useRules from '@disclosure-portal/utils/Rules';
import {getIconColorForPolicyType, getIconForPolicyType} from '@disclosure-portal/utils/View';
import useSnackbar from '@shared/composables/useSnackbar';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {VForm} from 'vuetify/components';
import {escapeHtml} from '@disclosure-portal/utils/Validation';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';

interface LicenseItemWithPolicyStatus {
  id: string;
  name?: string;
  policyType: string;
  icon: string;
  iconColor: string;
  isRecommended: boolean;
  weight: number | null;
}

const {t} = useI18n();
const {info} = useSnackbar();
const rules = useRules();
const sbomStore = useSbomStore();
const userStore = useUserStore();
const emit = defineEmits<{
  reload: [];
  triggerComponentDetails: [spdxId: string];
}>();
const projectStore = useProjectStore();

const form = ref<VForm | null>(null);
const isVisible = ref(false);
const isDirectOpen = ref(false);

const selectedComponent = ref<ComponentInfoSlim | undefined>(undefined);

const selectedLicense = ref<LicenseItemWithPolicyStatus | undefined>(undefined);
const componentLicenses = ref<ComponentLicenses | undefined>(undefined);
const licensesLoading = ref(false);

const comment = ref<string | undefined>(undefined);
const selectedComponentStr = ref<string>('');
const licenseExpression = ref<string>('');
const licenseRecommended = ref<string>('');
const licenseRecommendedMsg = ref<string>('');
const verification = ref(false);
const errorDialog = ref<ErrorDialogInterface | null>(null);

const licenseDecisionRules = rules.required(t('LICENSE_DECISION'));
const commentRules = rules.minMax(t('LICENSE_RULE_COMMENT'), 0, 80, true);

const config = ref<DialogLicenseRuleConfig>({
  licenseId: '',
  component: new ComponentInfoSlim(),
  licenseRecommended: '',
  licenseRecommendedMsg: '',
});

const projectKey = computed(() => projectStore.currentProject!._key);
const currentVersionId = computed(() => sbomStore.getCurrentVersion._key);
const currentSbomId = computed(() => sbomStore.getSelectedSBOM?._key);
const currentSbomName = computed(() => sbomStore.getSelectedSBOM?.MetaInfo.Name);
const currentSbomUploaded = computed(() => sbomStore.getSelectedSBOM?.Uploaded);

const policyTypeMap = computed(() => {
  const statuses = config.value.policyStatus ?? [];
  return new Map<string, string>(statuses.map((p) => [p.licenseMatched, p.type]));
});

function getPolicyType(id: string): string {
  return policyTypeMap.value.get(id) ?? 'noassertion';
}

const licenseRecommendationWeightMap = computed(() => {
  const statuses = config.value.policyStatus ?? [];
  return new Map<string, number | null>(statuses.map((p) => [p.licenseMatched, p.licenseRecommendationWeight]));
});

function getLicenseRecommendationWeight(id: string): number | null {
  return licenseRecommendationWeightMap.value.get(id) ?? null;
}

const hasSelectableLicenses = computed(() =>
  licenses.value.some((l) => l.policyType === 'allow' || l.policyType === 'warn'),
);

const isLicenseDisabled = (policyType: string) =>
  hasSelectableLicenses.value && policyType !== 'allow' && policyType !== 'warn';

const licenseItemProps = (item: LicenseItemWithPolicyStatus) => ({
  disabled: isLicenseDisabled(item.policyType),
});

const licenses = computed((): LicenseItemWithPolicyStatus[] => {
  if (!componentLicenses.value) {
    return [];
  }

  const known = componentLicenses.value.KnownLicenses.map((l) => {
    const id = l.License.licenseId;
    const name = l.License.name;
    const type = getPolicyType(id);
    return {
      id,
      name,
      policyType: type,
      icon: getIconForPolicyType(type),
      iconColor: getIconColorForPolicyType(type),
      isRecommended: l.License.licenseId === licenseRecommended.value,
      weight: getLicenseRecommendationWeight(id),
    };
  });
  const unknown = componentLicenses.value.UnknownLicenses.map((id) => {
    const type = getPolicyType(id);
    return {
      id,
      name: id,
      policyType: type,
      icon: getIconForPolicyType(type),
      iconColor: getIconColorForPolicyType(type),
      isRecommended: false,
      weight: null,
    };
  });

  return [...sortByWeight(known), ...unknown];
});

function sortByWeight(items: LicenseItemWithPolicyStatus[]): LicenseItemWithPolicyStatus[] {
  return [...items].sort((a, b) => {
    const aw = a.weight;
    const bw = b.weight;

    if (aw === null && bw === null) return 0;
    if (aw === null) return 1;
    if (bw === null) return -1;

    return aw - bw;
  });
}

const open = async (
  newConfig: DialogLicenseRuleConfig = {
    licenseId: '',
    component: new ComponentInfoSlim(),
    policyStatus: [],
    licenseRecommended: '',
    licenseRecommendedMsg: '',
  },
  directOpen: boolean = false,
) => {
  config.value = newConfig;
  isDirectOpen.value = directOpen;
  await loadAndPrefillData();
  isVisible.value = true;
};

const loadAndPrefillData = async () => {
  if (!config.value.component?.spdxId) return;

  selectedComponent.value = config.value.component;
  selectedComponentStr.value = `${config.value.component.name} (${config.value.component.version})`;
  licenseExpression.value = config.value.component.licenseExpression;
  licenseRecommended.value = config.value.licenseRecommended;
  licenseRecommendedMsg.value = config.value.licenseRecommendedMsg;

  await loadLicenses();

  const licenseIdToSelect = config.value.licenseId || config.value.licenseRecommended;

  const candidate = licenses.value.find((license) => license.id === licenseIdToSelect);

  selectedLicense.value = candidate && isLicenseDisabled(candidate.policyType) ? undefined : candidate;
};

const loadLicenses = async () => {
  licensesLoading.value = true;
  return versionService
    .getVersionComponentsLicenses(
      projectKey.value,
      currentVersionId.value,
      currentSbomId.value,
      selectedComponent.value!.spdxId,
    )
    .then((res) => {
      componentLicenses.value = res;
      licensesLoading.value = false;
    });
};

const doDialogAction = async () => {
  if (!(await form.value?.validate())?.valid) {
    return;
  }

  const licenseRuleRequest: LicenseRuleRequest = {
    sbomId: currentSbomId.value,
    sbomName: currentSbomName.value,
    sbomUploaded: currentSbomUploaded.value,
    componentSpdxId: selectedComponent.value!.spdxId,
    componentName: selectedComponent.value!.name,
    componentVersion: selectedComponent.value!.version,
    licenseExpression: selectedComponent.value!.licenseExpression,
    licenseDecisionId: selectedLicense.value!.id,
    licenseDecisionName: selectedLicense.value?.name ?? '',
    comment: comment.value ?? '',
    creator: userStore.getProfile.user,
  };

  const response = (
    await projectService.createLicenseRule(projectKey.value, currentVersionId.value, licenseRuleRequest)
  ).data;
  if (!response.success) {
    const dialog = new ErrorDialogConfig();
    dialog.title = t('license_rule_create_error_title');
    dialog.description = t(response.message);
    errorDialog.value?.open(dialog);
    return;
  }
  form.value?.reset();
  emit('reload');
  close();
  info(t('LICENSE_RULE_CREATED'));
};

const translatedLicenseRecommendedMsg = computed(() =>
  licenseRecommendedMsg.value === 'LICENSE_RECOMMENDED_MSG'
    ? t(licenseRecommendedMsg.value, {license: licenseRecommended.value})
    : t(licenseRecommendedMsg.value),
);

const dialogConfig = computed(() => ({
  title: t('LICENSE_RULE_CREATE'),
  primaryButton: {text: t('BTN_CREATE'), disabled: !verification.value},
  secondaryButton: {text: t('BTN_CANCEL')},
}));

const close = () => {
  form.value?.reset();
  isVisible.value = false;
};

const formatText = (text: string): string => {
  text = escapeHtml(text);
  if (text.includes(' AND ') || text.includes(' OR ')) {
    return text
      .replace(/ AND /g, ' <strong class="db-highlight">AND</strong> ')
      .replace(/ OR /g, ' <strong class="db-highlight">OR</strong> ');
  }
  return text;
};

const closeAndTriggerComponentDetails = () => {
  close();
  emit('triggerComponentDetails', selectedComponent.value?.spdxId ?? '');
};

defineExpose({open});
</script>

<template>
  <v-dialog v-model="isVisible" width="1200" persistent>
    <DialogLayout :config="dialogConfig" @primary-action="doDialogAction" @secondary-action="close" @close="close">
      <template #left>
        <DCActionButton
          v-if="isDirectOpen"
          size="small"
          is-dialog-button
          @click="closeAndTriggerComponentDetails"
          :text="t('BTN_COMPONENT_DETAILS')" />
      </template>

      <v-form ref="form" @submit.prevent="doDialogAction">
        <Stack class="gap-4">
          <Stack direction="row" class="items-start gap-4">
            <Stack class="flex-1 gap-4 self-start">
              <v-field variant="outlined" density="compact" active :label="t('RELATED_COMPONENT')" hide-details>
                <span class="v-field__input text-title-1 py-2">{{ selectedComponentStr }}</span>
              </v-field>
              <v-field variant="outlined" density="compact" active :label="t('LICENSE_EXPRESSION')" hide-details>
                <span
                  class="v-field__input text-title-1 py-2"
                  v-html="formatText(selectedComponent?.licenseExpression ?? '')" />
              </v-field>
              <v-select
                v-model="selectedLicense"
                clearable
                :label="t('LICENSE_DECISION')"
                :disabled="!selectedComponent"
                :items="licenses"
                :item-props="licenseItemProps"
                return-object
                item-title="name"
                :loading="licensesLoading"
                variant="outlined"
                density="compact"
                hide-details
                required
                :rules="licenseDecisionRules">
                <template #item="{item, props}">
                  <v-list-item v-bind="props" title="">
                    <v-chip
                      v-if="item.raw.isRecommended"
                      variant="outlined"
                      label
                      size="x-small"
                      class="mr-1 font-bold">
                      {{ t('RECOMMENDED') }}
                    </v-chip>
                    <v-icon size="small" :color="item.raw.iconColor">
                      {{ item.raw.icon }}
                    </v-icon>
                    <span class="d-subtitle-2 ml-2">{{ item.raw.name }}</span>
                    <span class="d-text d-secondary-text">&nbsp;({{ item.raw.id }})</span>
                  </v-list-item>
                </template>
                <template #selection="{item}">
                  <div class="d-inline">
                    <v-chip
                      v-if="item.raw.isRecommended"
                      variant="outlined"
                      label
                      size="x-small"
                      class="mr-1 font-bold">
                      {{ t('RECOMMENDED') }}
                    </v-chip>
                    <v-icon size="small" :color="item.raw.iconColor">
                      {{ item.raw.icon }}
                    </v-icon>
                    <span class="d-subtitle-2 ml-2">{{ item.raw.name }}</span>
                    <span class="d-text d-secondary-text">&nbsp;({{ item.raw.id }})</span>
                  </div>
                </template>
              </v-select>
              <v-textarea
                v-model="comment"
                variant="outlined"
                density="compact"
                :label="t('LICENSE_RULE_COMMENT')"
                hide-details="auto"
                :rules="commentRules" />
            </Stack>
            <v-divider vertical class="mx-5" />
            <Stack class="flex-1 gap-4 self-start">
              <Stack direction="row" align="center">
                <span class="text-h6">{{ t('IMPORTANT_INFO_TEXT') }}</span>
              </Stack>
              <Stack direction="row" align="center">
                <v-icon class="mr-2" color="brand">mdi-arrow-right</v-icon>
                <span>{{ t('LICENSE_RULE_APPLIED_LATER_INFO') }}</span>
              </Stack>
              <Stack v-if="licenseRecommendedMsg" direction="row" align="center">
                <v-icon class="mr-2" color="brand">mdi-arrow-right</v-icon>
                <span>{{ translatedLicenseRecommendedMsg }}</span>
              </Stack>
              <Stack v-if="selectedLicense?.policyType === 'noassertion'" direction="row" align="center">
                <v-icon class="mr-2" color="error">mdi-arrow-right</v-icon>
                <span>{{ t('UNASSERTED_LICENSES_MSG', {license: selectedLicense?.id}) }}</span>
              </Stack>
              <v-checkbox
                v-model="verification"
                color="primary"
                :label="t('LICENSE_RULE_VERIFICATION_NOTE_TEXT')"
                hide-details />
            </Stack>
          </Stack>
        </Stack>
      </v-form>
    </DialogLayout>
  </v-dialog>
  <ErrorDialog ref="errorDialog" />
</template>
