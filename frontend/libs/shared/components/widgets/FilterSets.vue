<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script lang="ts">
import {IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import ConfirmationDialog from '@disclosure-portal/components/dialog/ConfirmationDialog.vue';
import {useLicense} from '@disclosure-portal/composables/useLicense';
import {Filter, FilterSetDto, FilterSetRequestDto} from '@disclosure-portal/model/FilterSet';
import {IObligation, ISelectItemWithCount} from '@disclosure-portal/model/IObligation';
import {ClassificationWithCount, compareFamily} from '@disclosure-portal/model/License';
import {Nullable} from '@disclosure-portal/model/VersionDetails';
import AdminService from '@disclosure-portal/services/admin';
import filterSetService from '@disclosure-portal/services/filtersets';
import licenseService from '@disclosure-portal/services/license';
import {DiscoForm} from '@disclosure-portal/types/discobasics';
import {LabelsTools} from '@disclosure-portal/utils/Labels';
import useRules from '@disclosure-portal/utils/Rules';
import {SearchOptions} from '@disclosure-portal/utils/Table';
import useViewTools, {getIconColorOfLevel, getIconOfLevel} from '@disclosure-portal/utils/View';
import Statistics from '@disclosure-portal/views/admin/tools/Statistics.vue';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import DCloseButton from '@shared/components/disco/DCloseButton.vue';
import {defineComponent, nextTick, onMounted, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRoute} from 'vue-router';

export default defineComponent({
  name: 'FilterSets',
  methods: {
    getIconOfLevel,
    getIconColorOfLevel,
  },
  components: {
    DCActionButton,
    DCloseButton,
    ConfirmationDialog,
    Statistics,
  },
  setup(props, ctx) {
    const {t} = useI18n();
    const {getI18NTextOfPrefixKey} = useLicense();

    const selectedFilters = ref<Record<string, string[]>>({});
    const excludeFilters = ref<Record<string, string[]>>({});

    const isValid = ref(false);
    const route = useRoute();
    const kIndex = ref(0);
    const options = ref<SearchOptions>({} as SearchOptions);
    const filterSetMenu = ref(false);
    const filterSetName = ref('');
    const filterMap = ref<Nullable<Record<string, ISelectItemWithCount[]>>>(null);
    const activeRules = ref({});
    const labelTools = new LabelsTools();
    const licensesLoading = ref(true);
    const filterSetLoading = ref(false);
    const selectedTab = ref(0);
    const selectedFilterSet = ref<FilterSetDto | null>(null);
    const filterSets = ref<FilterSetDto[]>([]);
    const isNew = ref(false);
    const tableName = 'licenses';
    const form = ref<DiscoForm | null>(null);
    const classifications = ref<IObligation[]>([]); // Define the type here
    const confirmVisible = ref(false);
    const confirmConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);

    const rules = useRules();
    activeRules.value = {
      name: rules.minMax(t('FS_NAME'), 3, 80, true),
    };

    watch(selectedFilterSet, (value) => {
      if (filterSetLoading.value) {
        return;
      }
      onFilterSetChange();
    });
    watch(filterSetMenu, (value) => {
      if (value) {
        selectedTab.value = 0;
        selectedFilterSet.value = null;
        isNew.value = false;
        loadFilter();
        loadSelectedFilter();
        setOptionsAndFilters();
      }
    });

    onMounted(async () => {
      await retrieveClassifications();
    });

    const doNew = () => {
      isNew.value = true;
      filterSetName.value = 'New Filter';
      selectedFilters.value = createDefaultFilterMap();
      excludeFilters.value = createDefaultFilterMap();
    };

    const doDelete = () => {
      if (selectedFilterSet.value === null) {
        return;
      }
      confirmConfig.value = {
        key: selectedFilterSet.value._key,
        name: selectedFilterSet.value.name,
        okButton: 'BTN_DELETE',
        description: 'DLG_CONFIRMATION_DESCRIPTION',
      } as IConfirmationDialogConfig;
      confirmVisible.value = true;
    };

    const doDeleteFilter = async () => {
      if (selectedFilterSet.value === null) {
        return;
      }
      await filterSetService.delete(selectedFilterSet.value._key);
      ctx.emit('reloadFilter', selectedFilterSet.value._key);
      close();
    };

    const requestFilterDataFromTable = () => {
      ctx.emit('requestFilterDataFromTable');
    };

    const onClickCreate = async () => {
      if (validate()) {
        const filterSet: FilterSetRequestDto = {
          name: filterSetName.value,
          includedFilters: convertRecordToFilterArray(selectedFilters.value),
          excludedFilters: convertRecordToFilterArray(excludeFilters.value),
          tableName: tableName,
        };
        await filterSetService.create(filterSet);
        ctx.emit('reloadFilter');
        close();
      }
    };

    const onClickCancel = () => {
      if (isNew.value) {
        isNew.value = false;
      } else {
        close();
      }
    };

    const onClickSave = async () => {
      if (selectedFilterSet.value !== null && validate()) {
        const filterSet: FilterSetRequestDto = {
          name: selectedFilterSet.value.name,
          includedFilters: convertRecordToFilterArray(selectedFilters.value),
          excludedFilters: convertRecordToFilterArray(excludeFilters.value),
          tableName: tableName,
        };
        await filterSetService.update(filterSet, selectedFilterSet.value._key);
        close();
      }
    };

    const close = () => {
      filterSetMenu.value = false;
    };

    const retrieveClassifications = async () => {
      const response = (await AdminService.getAllObligations()).data;
      classifications.value = response.items;
    };

    const convertRecordToFilterArray = (filtersToTransform: Record<string, string[]>): Filter[] => {
      const filters: Filter[] = [];
      for (const key in filtersToTransform) {
        if (Object.prototype.hasOwnProperty.call(filtersToTransform, key)) {
          filters.push(new Filter(key, filtersToTransform[key]));
        }
      }
      return filters;
    };

    const setOptionsAndFilters = async () => {
      const res = (await licenseService.getAllWithOptions(options.value)).data;
      licensesLoading.value = false;

      const possibleIsLicenseChart = Object.entries(res.meta.possibleCharts)
        .map(([k, count]) => ({
          text: k === 'true' ? t('TABLE_LICENSE_CHART_STATUS_IS') : t('TABLE_LICENSE_CHART_STATUS_IS_NOT'),
          value: k,
          count: count,
        }))
        .sort();

      const possibleSources = Object.entries(res.meta.possibleSources).map(([k, count]) => ({
        text: k,
        value: k,
        count: count,
      }));

      const possibleFamilies = Object.entries(res.meta.possibleFamilies)
        .sort((a, b) => compareFamily(a[0], b[0]))
        .map(([k, count]) => ({
          text: getI18NTextOfPrefixKey('LIC_FAMILY_', k),
          value: k.length === 0 ? 'not declared' : k,
          count: count,
        }));

      const possibleApproval = Object.entries(res.meta.possibleApproval).map(([k, count]) => ({
        text: getI18NTextOfPrefixKey('LT_APP_', k),
        value: k.length === 0 ? 'not set' : k,
        count: count,
      }));

      const possibleType = Object.entries(res.meta.possibleType).map(([k, count]) => ({
        text: getI18NTextOfPrefixKey('LT_', k),
        value: k.length === 0 ? 'not declared' : k,
        count: count,
      }));
      let vt = useViewTools();

      const possibleClassifications = res.meta.possibleClassifications.map((c: ClassificationWithCount) => ({
        text:
          vt.getNameForLanguage(c.classification) === ''
            ? t('NO_CLASSIFICATIONS')
            : vt.getNameForLanguage(c.classification),
        value: vt.getNameForLanguage(c.classification) === '' ? '' : c.classification.name,
        count: c.count,
      }));

      filterMap.value = {
        isLicenseChart: possibleIsLicenseChart,
        source: possibleSources,
        family: possibleFamilies,
        approvalState: possibleApproval,
        licenseType: possibleType,
        classifications: possibleClassifications,
      };
    };

    const loadFilter = async () => {
      filterSets.value = await filterSetService.getFilterSets(tableName);
    };

    const loadSelectedFilter = async () => {
      Object.assign(selectedFilters, createDefaultFilterMap());
      Object.assign(excludeFilters, createDefaultFilterMap());

      if (route.path.includes('filtersets') && route.params.id) {
        const filter = await filterSetService.getFilterSet(route.params.id as string);
        selectedFilterSet.value = filter;
        convertFilterToUiModel(filter);
      }
      kIndex.value++;
    };

    const onFilterSetChange = async () => {
      if (!selectedFilterSet.value) {
        return;
      }
      filterSetLoading.value = true;
      selectedFilterSet.value = await filterSetService.getFilterSet(selectedFilterSet.value._key);
      filterSetName.value = selectedFilterSet.value.name;
      convertFilterToUiModel(selectedFilterSet.value);
      await nextTick(() => {
        filterSetLoading.value = false;
      });
    };

    const validate = (): boolean => {
      return isValid.value;
    };

    const convertFilterToUiModel = (filter: FilterSetDto) => {
      filter.includedFilters.forEach((filter: Filter) => {
        selectedFilters.value[filter.name] = filter.values;
      });
      filter.excludedFilters.forEach((filter: Filter) => {
        excludeFilters.value[filter.name] = filter.values;
      });
    };

    const createDefaultFilterMap = () => {
      return {
        isLicenseChart: [],
        source: [],
        family: [],
        approvalState: [],
        licenseType: [],
      };
    };

    const setFilterData = (newSelectedFilters: Record<string, string[]>) => {
      selectedFilters.value = newSelectedFilters;
    };

    const classificationCheckboxCss = (state: string) => {
      return selectedFilters.value.classifications?.includes(state)
        ? 'v-icon notranslate mdi mdi-checkbox-marked primary--text'
        : 'v-icon notranslate mdi mdi-checkbox-blank-outline';
    };

    const getWarnLevel = (name: string) => {
      const classification = classifications.value.find((c) => c.name === name || c.nameDe === name);
      return classification ? classification.warnLevel : 'INFORMATION';
    };

    return {
      isValid,
      t,
      selectedFilters,
      excludeFilters,
      kIndex,
      options,
      filterSetMenu,
      filterSetName,
      filterMap,
      activeRules,
      labelTools,
      licensesLoading,
      selectedTab,
      selectedFilterSet,
      filterSets,
      isNew,
      tableName,
      form,
      classifications,
      doNew,
      doDelete,
      requestFilterDataFromTable,
      doDeleteFilter,
      onClickCreate,
      onClickCancel,
      onClickSave,
      close,
      retrieveClassifications,
      setOptionsAndFilters,
      validate,
      createDefaultFilterMap,
      onFilterSetChange,
      classificationCheckboxCss,
      getWarnLevel,
      setFilterData,
      confirmVisible,
      confirmConfig,
    };
  },
});
</script>
<template>
  <v-row tag="div" class="d-flex justify-end pr-3">
    <v-menu v-model="filterSetMenu" :close-on-content-click="false" max-width="1024">
      <template v-slot:activator="{props}">
        <span class="discoActionBtnHover">
          <v-btn variant="tonal" color="primary" prepend-icon="mdi mdi-filter-variant" class="text-none" v-bind="props">
            <span style="font-weight: 700 !important">
              {{ t('MANAGE_FILTER_BTN') }}
            </span>
          </v-btn>
        </span>
      </template>
      <v-form ref="filtersetsform" v-model="isValid">
        <v-card density="compact" class="pa-4" v-if="selectedFilters" style="min-width: 400px">
          <v-card-title class="d-flex justify-end">
            <v-col cols="12" xs="12" sm="10">
              <h1 class="d-subtitle">{{ t('MANAGE_FILTER_BTN') }}</h1>
            </v-col>
            <v-spacer></v-spacer>
            <DCloseButton @click="close()" />
          </v-card-title>
          <v-card-text>
            <v-row>
              <v-col cols="12" xs="12" md="6">
                <v-select
                  v-if="!isNew"
                  v-model="selectedFilterSet"
                  :items="filterSets"
                  :item-title="['name']"
                  :label="t('FILTER_SET_LABEL')"
                  variant="outlined"
                  density="compact"
                  hide-details="auto"
                  clearable
                  v-bind:menu-props="{location: 'bottom'}"
                  return-object></v-select>
                <v-text-field
                  autocomplete="off"
                  variant="outlined"
                  density="compact"
                  v-model="filterSetName"
                  :label="t('FS_NAME')"
                  v-if="isNew"
                  class="required"
                  :rules="activeRules.name"></v-text-field>
              </v-col>
              <v-col cols="12" xs="12" md="6">
                <DCActionButton
                  large
                  v-if="!isNew"
                  icon="mdi-plus"
                  :hint="t('TT_add_filter')"
                  :text="t('BTN_ADD')"
                  @click="doNew"
                  class="mx-2 py-3" />
                <DCActionButton
                  large
                  v-if="!isNew && selectedFilterSet?._key"
                  icon="mdi-delete"
                  :hint="t('TT_delete_filter')"
                  :text="t('BTN_DELETE')"
                  @click="doDelete"
                  class="mx-2 py-3" />
                <DCActionButton
                  v-if="isNew"
                  large
                  :hint="t('TT_apply_filter')"
                  :text="t('BTN_APPLY_FILTER')"
                  icon="mdi mdi-filter-variant-plus"
                  @click="requestFilterDataFromTable"
                  class="mx-2 py-3" />
              </v-col>
            </v-row>
            <v-divider class="my-5 mr-3"></v-divider>
            <v-row class="px-2">
              <v-col cols="12" xs="12" md="6">
                <v-text-field
                  autocomplete="off"
                  v-if="!isNew && selectedFilterSet?._key"
                  density="compact"
                  v-model="selectedFilterSet.name"
                  :label="t('FS_NAME')"
                  variant="outlined"
                  class="required"
                  :rules="activeRules.name"></v-text-field>
              </v-col>
              <template v-for="(elements, key) in filterMap">
                <v-col
                  cols="12"
                  xs="12"
                  md="6"
                  v-if="
                    key === 'classifications' && filterMap && ((selectedFilterSet && selectedFilterSet._key) || isNew)
                  "
                  :key="'fm' + key">
                  <v-select
                    density="compact"
                    v-model="selectedFilters[key]"
                    :items="elements"
                    :key="key"
                    :item-title="['text']"
                    multiple
                    clearable
                    variant="outlined"
                    :label="labelTools.camelCaseToLabel(key)"
                    hide-details="auto"
                    v-bind:menu-props="{location: 'bottom'}">
                    <template v-slot:item="{item, props}">
                      <v-list-item v-bind="props">
                        <template v-slot:prepend="{isSelected}">
                          <v-list-item-action start>
                            <v-checkbox-btn :model-value="isSelected"></v-checkbox-btn>
                          </v-list-item-action>
                        </template>
                        <template v-slot:title>
                          <v-icon
                            :color="getIconColorOfLevel(getWarnLevel(item.raw.value))"
                            :icon="getIconOfLevel(getWarnLevel(item.raw.value).toUpperCase())" />
                          {{ item.title }}
                          {{ key + ' ' + selectedFilters[key] }}
                        </template>
                      </v-list-item>
                    </template>

                    <template v-slot:selection="{item, index}">
                      <v-chip v-if="index < 2">
                        <span>{{ item.title }}</span>
                      </v-chip>
                      <span v-if="index === 2" class="text-grey text-caption align-self-center">
                        (+{{ selectedFilters[key].length - 2 }} others)
                      </span>
                    </template>
                  </v-select>
                </v-col>
              </template>
            </v-row>
            <v-row v-if="filterMap && ((selectedFilterSet && selectedFilterSet._key) || isNew)">
              <div v-for="(elements, key) in filterMap" :key="'fm' + key">
                <div v-if="key != 'classifications'">
                  <h3 class="d-subtitle-2 font-weight-bold">{{ labelTools.camelCaseToLabel(key) }}</h3>
                  <v-checkbox
                    v-for="element in elements"
                    :key="element.value"
                    v-model="selectedFilters[key]"
                    :label="element.text"
                    :value="element.value"
                    multiple
                    density="compact"></v-checkbox>
                </div>
              </div>
            </v-row>
          </v-card-text>
          <v-card-actions class="d-flex justify-end">
            <DCActionButton size="small" variant="text" @click="onClickCancel()" class="mr-5" :text="t('BTN_CANCEL')" />
            <DCActionButton
              size="small"
              variant="flat"
              @click="onClickSave()"
              :text="t('Btn_save')"
              v-if="!isNew && selectedFilterSet && selectedFilterSet._key" />
            <DCActionButton
              size="small"
              variant="flat"
              @click="onClickCreate"
              :text="t('NP_DIALOG_BTN_EDIT')"
              v-if="isNew" />
          </v-card-actions>
        </v-card>
      </v-form>
    </v-menu>
    <ConfirmationDialog v-model:showDialog="confirmVisible" :config="confirmConfig" @confirm="doDeleteFilter" />
  </v-row>
</template>
