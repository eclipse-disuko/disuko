<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script lang="ts">
import {useView} from '@disclosure-portal/composables/useView';
import {IObligation} from '@disclosure-portal/model/IObligation';
import {compareLevel} from '@disclosure-portal/model/Quality';
import {Rights} from '@disclosure-portal/model/Rights';
import {useUserStore} from '@disclosure-portal/stores/user';
import useViewTools, {getIconColorOfLevel, getIconOfLevel} from '@disclosure-portal/utils/View';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import DCloseButton from '@shared/components/disco/DCloseButton.vue';
import {DataTableHeader, SortItem} from '@shared/types/table';
import {TOOLTIP_OPEN_DELAY_IN_MS} from '@shared/utils/constant';
import {defineComponent, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRouter} from 'vue-router';

export default defineComponent({
  name: 'ClassificationsPerLicenseDialog',
  components: {
    DCActionButton,
    DCloseButton,
  },
  methods: {getIconOfLevel, getIconColorOfLevel},
  setup() {
    const sortItems = ref<SortItem[]>([{key: 'warnLevel', order: 'desc'}]);
    const viewTools = useViewTools();
    const {t} = useI18n();
    const classifications = ref<IObligation[]>([]);
    const headers = ref<DataTableHeader[]>([]);
    const licenseName = ref('');
    const show = ref(false);
    const licenseId = ref('');
    const rights = ref<Rights>({} as Rights);
    const userStore = useUserStore();
    const router = useRouter();
    const {getTextOfLevel} = useView();

    const open = (classificationsData: IObligation[], licenseNameData: string, licenseIdData: string) => {
      rights.value = userStore.getRights;
      classifications.value = classificationsData;
      licenseName.value = licenseNameData;
      licenseId.value = licenseIdData;
      headers.value = [
        {
          title: t('COL_WARN_LEVEL'),
          align: 'center',
          sortable: true,
          sort: compareLevel,
          width: 50,
          class: 'tableHeaderCell',
          value: 'warnLevel',
        },
        {
          title: t('COL_SHORT_NAME'),
          align: 'start',
          width: 280,
          class: 'tableHeaderCell',
          value: 'name',
        },
      ];
      show.value = true;
    };

    const close = () => {
      show.value = false;
    };

    const openLicenseClassificationTab = () => {
      if (rights.value.allowLicense && rights.value.allowLicense.read) {
        router.push({name: 'LicenseClassifications', params: {id: licenseId.value}});
      }
    };

    return {
      TOOLTIP_OPEN_DELAY_IN_MS,
      t,
      classifications,
      headers,
      licenseName,
      show,
      licenseId,
      rights,
      open,
      close,
      viewTools,
      openLicenseClassificationTab,
      sortItems,
      getTextOfLevel,
    };
  },
});
</script>
<template>
  <v-dialog v-model:model-value="show" scrollable max-width="800">
    <v-card class="pa-8">
      <v-card-title>
        <v-row>
          <v-col cols="10" class="d-flex align-center">
            <span class="text-h5">
              {{ licenseName }}
            </span>
          </v-col>
          <v-col cols="2" class="text-right px-0">
            <DCloseButton @click="close" />
          </v-col>
        </v-row>
      </v-card-title>
      <v-card-text class="px-2">
        <v-data-table
          v-if="classifications"
          :headers="headers"
          density="compact"
          class="striped-table custom-data-table"
          :items="classifications"
          :items-per-page="-1"
          :sort-by="sortItems"
          fixed-header
          height="400"
          hide-default-footer
          @click:row="openLicenseClassificationTab">
          <template v-slot:item.warnLevel="{item}">
            <span>
              <v-tooltip :open-delay="TOOLTIP_OPEN_DELAY_IN_MS" location="bottom" content-class="dpTooltip">
                <template v-slot:activator="{props}">
                  <v-icon v-bind="props" :color="getIconColorOfLevel(item.warnLevel)">
                    {{ getIconOfLevel(item.warnLevel) }}
                  </v-icon>
                </template>
                <span>{{ getTextOfLevel(item.warnLevel) }}</span>
              </v-tooltip>
            </span>
          </template>
          <template v-slot:item.name="{item}">
            {{ viewTools.getNameForLanguage(item) }}
          </template>
        </v-data-table>
        <span v-else>{{ t('NO_CLASSIFICATIONS') }}</span>
      </v-card-text>
      <v-card-actions class="justify-end">
        <DCActionButton isDialogButton size="small" variant="flat" @click="close" :text="t('BTN_OK')" />
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
<style scoped></style>
