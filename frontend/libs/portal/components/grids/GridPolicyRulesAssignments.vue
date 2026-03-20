<template>
  <TableLayout has-title has-tab>
    <template #buttons>
      <v-spacer></v-spacer>
      <DCActionButton
        id="gridControlPolicyRules"
        icon="mdi-bank"
        size="large"
        @click="configurePoliciesForLicense"
        :text="t('BTN_MANAGE')"
        :hint="t('CONFIGURE_POLICIES_FOR_LICENSE_TOOLTIP')"
        v-if="
          rights &&
          rights.allowLicense &&
          rights.allowLicense.read &&
          rights.allowPolicy &&
          rights.allowPolicy.create &&
          rights.allowPolicy.read &&
          rights.allowPolicy.update &&
          rights.allowPolicy.delete
        "></DCActionButton>
      <v-text-field
        autocomplete="off"
        class="max-w-[500px]"
        v-model="search"
        variant="outlined"
        clearable
        density="compact"
        append-inner-icon="mdi-magnify"
        :label="t('labelSearch')"
        hide-details />
    </template>
    <template #table>
      <PolicyRulesTable class="fill-height" v-model="items" :loading="!dataLoaded"></PolicyRulesTable>
    </template>
  </TableLayout>
  <ConfigurePoliciesForLicenseDialog
    ref="configurePoliciesForLicenseDialogRef"
    @policyRulesSaved="reload"></ConfigurePoliciesForLicenseDialog>
</template>

<script setup lang="ts">
import License from '@disclosure-portal/model/License';
import {PolicyRulesAssignmentsDto} from '@disclosure-portal/model/PolicyRule';
import {Rights} from '@disclosure-portal/model/Rights';
import {default as LicenseService} from '@disclosure-portal/services/license';
import {useUserStore} from '@disclosure-portal/stores/user';
import {computed, onBeforeMount, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const props = defineProps<{
  licenseId: string;
  disabled?: boolean;
}>();

const {t} = useI18n();
const userStore = useUserStore();

const search = ref('');
const items = ref<PolicyRulesAssignmentsDto[]>([]);
const dataLoaded = ref(false);
const configurePoliciesForLicenseDialogRef = ref();

const rights = computed<Rights>(() => userStore.getRights);

const reload = async () => {
  dataLoaded.value = false;
  const response = await LicenseService.getPolicyRuleAssignmentsForThisLicence(props.licenseId);
  items.value = response.policyRulesAssignments;
  dataLoaded.value = true;
};

const configurePoliciesForLicense = async () => {
  const fullModel: License = (await LicenseService.get(props.licenseId)).data;
  configurePoliciesForLicenseDialogRef.value?.open(fullModel.licenseId, fullModel.name);
};

onBeforeMount(async () => {
  await reload();
});
</script>
