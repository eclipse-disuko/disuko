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

interface LicenseItemWithPolicyStatus {
  id: string;
  name?: string;
  policyType: string;
  icon: string;
  iconColor: string;
}

const {t} = useI18n();
const {info} = useSnackbar();
const rules = useRules();
const sbomStore = useSbomStore();
const userStore = useUserStore();
const emit = defineEmits(['reload']);
const projectStore = useProjectStore();

const form = ref<VForm | null>(null);
const isVisible = ref(false);

const selectedComponent = ref<ComponentInfoSlim | undefined>(undefined);

const selectedLicense = ref<LicenseItemWithPolicyStatus | undefined>(undefined);
const componentLicenses = ref<ComponentLicenses | undefined>(undefined);
const licensesLoading = ref(false);

const comment = ref<string | undefined>(undefined);
const selectedComponentStr = ref<string>('');
const licenseExpression = ref<string>('');
const verification = ref(false);
const errorDialog = ref<ErrorDialogInterface | null>(null);

const licenseDecisionRules = rules.required(t('LICENSE_DECISION'));
const commentRules = rules.minMax(t('LICENSE_RULE_COMMENT'), 0, 80, true);

const config = ref<DialogLicenseRuleConfig>({
  licenseId: '',
  component: new ComponentInfoSlim(),
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
    };
  });

  return [...known, ...unknown];
});

const open = async (
  newConfig: DialogLicenseRuleConfig = {
    licenseId: '',
    component: new ComponentInfoSlim(),
    policyStatus: [],
  },
) => {
  config.value = newConfig;
  await loadAndPrefillData();
  isVisible.value = true;
};

const loadAndPrefillData = async () => {
  if (!config.value.component?.spdxId) return;

  selectedComponent.value = config.value.component;
  selectedComponentStr.value = `${config.value.component.name} (${config.value.component.version})`;
  licenseExpression.value = config.value.component.licenseExpression;

  await loadLicenses();

  if (!config.value.licenseId) return;

  selectedLicense.value = licenses.value.find((license) => license.id === config.value.licenseId);
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

const dialogConfig = computed(() => ({
  title: t('LICENSE_RULE_CREATE'),
  primaryButton: {text: t('BTN_CREATE'), disabled: !verification.value},
  secondaryButton: {text: t('BTN_CANCEL')},
}));

const close = () => {
  form.value?.reset();
  isVisible.value = false;
};
defineExpose({open});
</script>

<template>
  <v-dialog v-model="isVisible" width="650" persistent>
    <DialogLayout :config="dialogConfig" @primary-action="doDialogAction" @secondary-action="close" @close="close">
      <v-form ref="form" @submit.prevent="doDialogAction">
        <Stack class="gap-4">
          <Stack direction="row" align="center">
            <v-icon class="mr-2">mdi-information-outline</v-icon>
            <span>{{ t('LICENSE_RULE_APPLIED_LATER_INFO') }}</span>
          </Stack>
          <v-text-field
            autocomplete="off"
            :model-value="selectedComponentStr"
            disabled
            variant="outlined"
            density="compact"
            hide-details
            :label="t('RELATED_COMPONENT')" />
          <v-textarea
            auto-grow
            rows="1"
            variant="outlined"
            density="compact"
            disabled
            :label="t('LICENSE_EXPRESSION')"
            :model-value="selectedComponent?.licenseExpression"
            hide-details />
          <v-select
            v-model="selectedLicense"
            clearable
            :label="t('LICENSE_DECISION')"
            :disabled="!selectedComponent"
            :items="licenses"
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
                <v-icon size="small" :color="item.raw.iconColor">
                  {{ item.raw.icon }}
                </v-icon>
                <span class="d-subtitle-2 ml-2">{{ item.raw.name }}</span>
                <span class="d-text d-secondary-text">&nbsp;({{ item.raw.id }})</span>
              </v-list-item>
            </template>
            <template #selection="{item}">
              <div class="d-inline">
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
          <v-checkbox v-model="verification" :label="t('LICENSE_RULE_VERIFICATION_NOTE_TEXT')" hide-details />
        </Stack>
      </v-form>
    </DialogLayout>
  </v-dialog>
  <ErrorDialog ref="errorDialog" />
</template>
