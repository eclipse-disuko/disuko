<script setup lang="ts">
import SBOMComponentLicenseTab from '@cli/components/dialogs/tabs/SBOMComponentLicenseTab.vue';
import SBOMComponentPolicyTab from '@cli/components/dialogs/tabs/SBOMComponentPolicyTab.vue';
import SBOMComponentScanTab from '@cli/components/dialogs/tabs/SBOMComponentScanTab.vue';
import {LicenseRemark, PolicyRule, ScanRemark, SpdxStatusComponent} from '@cli/models/Sbom';
import DialogLayout, {DialogLayoutConfig} from '@shared/layouts/DialogLayout.vue';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();

const show = ref(false);
const component = ref();

function open(c: SpdxStatusComponent) {
  component.value = c;
  activeTab.value = 'policy';
  show.value = true;
}
function close() {
  show.value = false;
  setTimeout(() => {
    component.value = null;
    activeTab.value = 'policy';
  }, 100);
}

const scanRemarks = computed<ScanRemark[]>(() => (component.value?.scanRemarks as ScanRemark[]) || []);
const licenseRemarks = computed<LicenseRemark[]>(() => (component.value?.licenseRemarks as LicenseRemark[]) || []);
const policyRules = computed<PolicyRule[]>(() => (component.value?.policyRuleStatus as PolicyRule[]) || []);
const activeTab = ref('policy');

defineExpose({open});

const dialogConfig = computed<DialogLayoutConfig>(() => ({
  title: `${component.value?.name || ''} ${component.value?.version || ''}`.trim(),
}));
</script>

<template>
  <v-dialog v-model="show" width="1200" max-width="90vw">
    <DialogLayout :config="dialogConfig" @close="close">
      <v-card height="500" flat class="pa-4 overflow-auto">
        <v-tabs
          v-model="activeTab"
          slider-color="mbti"
          active-class="active"
          show-arrows
          bg-color="tabsHeader"
          density="compact">
          <v-tab value="policy" class="px-8">{{ t('POLICY_RULES') }}</v-tab>
          <v-tab value="license" class="px-8">{{ t('LICENSES_REMARKS') }}</v-tab>
          <v-tab value="scan" class="px-8">{{ t('SCAN_REMARKS') }}</v-tab>
        </v-tabs>
        <v-tabs-window v-model="activeTab" class="flex-grow">
          <v-tabs-window-item value="policy">
            <SBOMComponentPolicyTab :policy-rules="policyRules" />
          </v-tabs-window-item>
          <v-tabs-window-item value="license">
            <SBOMComponentLicenseTab :license-remarks="licenseRemarks" />
          </v-tabs-window-item>
          <v-tabs-window-item value="scan">
            <SBOMComponentScanTab :scan-remarks="scanRemarks" />
          </v-tabs-window-item>
        </v-tabs-window>
      </v-card>
      <v-alert v-if="!component" type="warning" variant="tonal" density="compact" class="mt-4">{{
        t('NO_COMPONENT_SELECTED')
      }}</v-alert>
    </DialogLayout>
  </v-dialog>
</template>
