<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script lang="ts" setup>
import {CustomerMetaDTO, NoticeContactMetaDTO} from '@disclosure-portal/model/Project';
import {Group, Rights} from '@disclosure-portal/model/Rights';
import {RightsUtils} from '@disclosure-portal/utils/Rights';
import {useI18n} from 'vue-i18n';

interface Props {
  hasParent?: boolean;
  rights?: Rights;
  vehicleOnboard?: boolean;
  activeRules: Record<string, any>;
}

withDefaults(defineProps<Props>(), {
  hasParent: false,
  vehicleOnboard: false,
});

const {t} = useI18n();

const customerMeta = defineModel<CustomerMetaDTO>('customerMeta', {required: true});
const noticeMeta = defineModel<NoticeContactMetaDTO>('noticeMeta', {required: true});
</script>

<template>
  <Stack class="gap-3 pt-4">
    <Stack direction="row" v-if="hasParent">
      <v-icon color="warning" class="mr-2">mdi-alert</v-icon>
      <span>{{ t('OWNER_SETTINGS_FROM_PARENT') }}</span>
    </Stack>

    <DAutocompleteCompany
      id="owner-company"
      v-if="RightsUtils.rights().isInternal"
      v-model="customerMeta.dept"
      :readonly="hasParent || (rights && !rights.groups?.includes(Group.ProjectOwner))"
      :disabled="hasParent || (rights && !rights.groups?.includes(Group.ProjectOwner))"
      :label="t('COMPANY')"
      required
      ref="deptAutoComplete"
      aria="owner company"></DAutocompleteCompany>

    <v-textarea
      variant="outlined"
      no-resize
      rows="4"
      v-model="customerMeta.address"
      :label="t('PROJECT_SETTINGS_ADDRESS')"
      hide-details="auto"
      data-testid="OwnerSettings__Address"
      :rules="activeRules.address"
      :readonly="hasParent || (rights && !rights.groups?.includes(Group.ProjectOwner))"
      :disabled="hasParent || (rights && !rights.groups?.includes(Group.ProjectOwner))" />

    <v-textarea
      v-if="!vehicleOnboard"
      id="thirdparty-address"
      rows="5"
      autocomplete="off"
      :placeholder="t('PLACEHOLDER_NOTICE_CONTACT_ADDRESS')"
      persistent-placeholder
      variant="outlined"
      v-model="noticeMeta.address"
      :label="t('NOTICE_CONTACT_ADDRESS')"
      hide-details="auto"
      :rules="activeRules.address"
      :readonly="hasParent || (rights && !rights.groups?.includes(Group.ProjectOwner))"
      :disabled="hasParent || (rights && !rights.groups?.includes(Group.ProjectOwner))" />
  </Stack>
</template>
