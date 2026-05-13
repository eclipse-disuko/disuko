<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script lang="ts">
import DiffDialog from '@disclosure-portal/components/dialog/DiffDialog.vue';
import {LicenseWithSimilarity} from '@disclosure-portal/model/License';
import licenseService from '@disclosure-portal/services/license';
import {DiscoForm} from '@disclosure-portal/types/discobasics';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import DCloseButton from '@shared/components/disco/DCloseButton.vue';
import DLicenseChartIcon from '@disclosure-portal/components/disco/DLicenseChartIcon.vue';
import {DataTableHeader, SortItem} from '@shared/types/table';
import {CodeDiff} from 'v-code-diff';
import {defineComponent, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {openUrlInNewTab} from '@shared/utils/url';

class DiffDetails {
  public oldValue = '';
  public newValue = '';

  constructor(oldValue: string, newValue: string) {
    this.oldValue = oldValue;
    this.newValue = newValue;
  }
}

export default defineComponent({
  name: 'LicenseCompareDialog',
  components: {
    DCActionButton,
    DCloseButton,
    DLicenseChartIcon,
    CodeDiff,
    DiffDialog,
  },
  props: {
    class: {
      type: String,
      required: false,
      default: '',
    },
  },
  setup() {
    const {t} = useI18n();

    const isDialogVisible = ref(false);
    const title = ref('LIC_DIALOG_SEARCH_TEXT_TITLE');
    const form = ref<DiscoForm | null>(null);
    const similarLicenses = ref<LicenseWithSimilarity[]>([]);
    const licenseText = ref('');
    const sortBy = ref<SortItem[]>([]);
    sortBy.value = [{key: 'similarity', order: 'desc'}];

    const headers = ref<DataTableHeader[]>([
      {
        title: t('COL_LICENSE_TEXT_SIMILARITY'),
        value: t('COL_LICENSE_TEXT_SIMILARITY'),
        align: 'center',
        key: 'similarity',
        width: 110,
        sortable: true,
      },
      {
        title: t('COL_LICENSE_CHART_STATUS'),
        align: 'center',
        value: 'license.meta.isLicenseChart',
        width: 100,
      },
      {
        title: t('COL_LICENSE_NAME'),
        align: 'start',
        value: 'license.name',
        width: 300,
      },
      {
        title: t('COL_LICENSE_ID'),
        align: 'start',
        value: 'license.licenseId',
        width: 200,
      },
      {
        title: t('COL_APPROVAL_STATUS'),
        align: 'start',
        value: 'license.meta.approvalState',
        width: 150,
      },
      {
        title: t('COL_LICENSE_SOURCE'),
        align: 'start',
        value: 'license.source',
        width: 100,
      },
      {
        title: t('COL_ACTIONS'),
        align: 'start',
        value: 'actions',
        width: 120,
        sortable: false,
      },
    ]);
    const diffDetails = ref<DiffDetails[]>([]);
    const licenseId = ref('');
    const compareLoading = ref(false);

    onMounted(() => {
      form.value = null;
    });

    const showDialog = () => {
      isDialogVisible.value = true;
      title.value = 'LIC_DIALOG_SEARCH_TEXT_TITLE';
    };

    const search = (text: string) => {
      showDialog();
      licenseText.value = text;
      sendLicenseText();
    };

    const fillDiffDetails = (item: LicenseWithSimilarity) => {
      diffDetails.value = [];
      licenseId.value = item.license.licenseId;
      diffDetails.value.push(new DiffDetails(licenseText.value.trim(), item.license.text.trim()));
    };

    const goToSimilarLicenses = () => {
      diffDetails.value = [];
    };

    const hideDialog = () => {
      isDialogVisible.value = false;
    };

    const clearSearch = () => {
      diffDetails.value = [];
      licenseText.value = '';
      similarLicenses.value = [];
      hideDialog();
    };

    const doClose = () => {
      clearSearch();
      hideDialog();
    };

    const sendLicenseText = async () => {
      if (licenseText.value.length > 0) {
        compareLoading.value = true;
        try {
          const response = await licenseService.searchForSimilarLicenseText(licenseText.value);
          similarLicenses.value = response.data;
        } catch (e) {
          console.error(e);
        }
        compareLoading.value = false;
        localStorage.removeItem('licenseText');
      }
    };

    const doCompare = async () => {
      await sendLicenseText();
    };

    const formatPercentage = (license: LicenseWithSimilarity): string => {
      return (license.similarity * 100).toFixed(2) + ' %';
    };

    const onClickRow = (event: any, row: any) => {
      openUrlInNewTab('/dashboard/licenses/' + row.item.license.licenseId);
    };

    return {
      isDialogVisible,
      t,
      title,
      form,
      similarLicenses,
      licenseText,
      headers,
      sortBy,
      diffDetails,
      licenseId,
      compareLoading,
      showDialog,
      search,
      doCompare,
      fillDiffDetails,
      goToSimilarLicenses,
      clearSearch,
      formatPercentage,
      onClickRow,
      doClose,
    };
  },
});
</script>

<template>
  <div>
    <slot :showDialog="showDialog">
      <v-btn text="Replace me" size="small" color="primary" @click.stop="showDialog"></v-btn>
    </slot>

    <v-dialog v-model="isDialogVisible" width="75%" scrollable>
      <v-form ref="formDialog">
        <v-card class="pa-8 dDialog">
          <v-card-title>
            <v-row>
              <v-col cols="10">
                <span class="text-h5">{{ t(title) }}</span>
              </v-col>
              <v-col cols="2" align="right">
                <DCloseButton @click="doClose" />
              </v-col>
            </v-row>
          </v-card-title>
          <v-card-text class="px-0 pt-2">
            <v-col cols="12" xs="12">
              <span class="text-body-1">{{ t('LicenseTextSearchDialogHelperText') }}</span>
            </v-col>
            <v-col cols="12" xs="12" v-if="!(similarLicenses.length > 0)">
              <v-textarea
                no-resize
                variant="outlined"
                :label="t('CD_LICENSE_TEXT')"
                v-model="licenseText"
                :loading="compareLoading"
                :disabled="compareLoading"
                rows="10"
                flat />
            </v-col>

            <v-col cols="12" xs="12" v-if="similarLicenses.length > 0 && !(diffDetails.length > 0)">
              <v-data-table
                :headers="headers"
                :items="similarLicenses"
                :sort-by="sortBy"
                :item-class="(item: LicenseWithSimilarity) => (item.similarity < 0.8 ? 'irrelevant' : '')"
                @click:row="onClickRow"
                hide-default-footer
                class="license-compare-table striped-table"
                density="compact">
                <template v-slot:header.similarity="{column}">
                  {{ column.title }}
                  <tooltip :text="t('COL_LICENSE_TEXT_SIMILARITY_TOOLTIP')"></tooltip>
                </template>
                <template v-slot:header.actions="{column}">
                  {{ column.title }}
                  <tooltip :text="t('COL_LICENSE_TEXT_SIMILARITY_TOOLTIP')"></tooltip>
                </template>
                <template v-slot:item.similarity="{item}">
                  {{ formatPercentage(item) }}
                </template>
                <template v-slot:item.license.meta.isLicenseChart="{item}">
                  <DLicenseChartIcon :meta="item.license.meta" />
                </template>
                <template v-slot:item.actions="{item}">
                  <v-btn
                    class="bg-transparent"
                    variant="plain"
                    size="large"
                    @click.stop="fillDiffDetails(item)"
                    color="primary"
                    icon>
                    <v-icon>mdi-compare-horizontal</v-icon>
                  </v-btn>
                </template>
              </v-data-table>
            </v-col>
            <v-col v-if="diffDetails.length > 0">
              <h2 class="d-headline mb-3">{{ t('DIFF_DETAILS') }}: {{ licenseId }}</h2>
              <v-row class="mt-0 mr-0 mb-3 ml-0" v-for="(diff, i) in diffDetails" :key="i">
                <code-diff
                  :old-string="diff.oldValue"
                  :new-string="diff.newValue"
                  output-format="side-by-side"
                  diff-style="word" />
              </v-row>
            </v-col>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn
              isDialogButton
              @click="doCompare"
              class="mt-4"
              :disabled="licenseText.length === 0"
              :loading="compareLoading"
              v-if="!(similarLicenses.length > 0)"
              :variant="licenseText.length === 0 ? 'tonal' : 'flat'"
              :color="licenseText.length === 0 ? 'disabled' : 'primary'">
              {{ t('Btn_compare') }}
            </v-btn>
            <v-btn
              isDialogButton
              @click="doClose"
              class="primary mt-4"
              variant="flat"
              color="primary"
              v-if="similarLicenses.length > 0 && !(diffDetails.length > 0)">
              {{ t('BTN_CLOSE') }}
            </v-btn>
            <v-btn
              isDialogButton
              @click="goToSimilarLicenses"
              class="primary mt-4"
              v-if="diffDetails.length > 0"
              variant="flat"
              color="primary">
              {{ t('BTN_BACK') }}
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-form>
    </v-dialog>
  </div>
</template>

<style scoped>
.license-compare-table {
  font-size: 0.75em;

  .irrelevant {
    opacity: 0.7;
  }
}

.bg-transparent {
  background-color: transparent;
}
</style>
