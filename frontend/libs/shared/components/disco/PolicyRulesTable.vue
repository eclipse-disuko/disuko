<template>
  <v-data-table-virtual
    :loading="loading"
    :headers="filteredHeaders"
    :items="filteredList"
    :search="search"
    density="compact"
    class="striped-table custom-data-table v-data-table--dense fill-height"
    disable-pagination
    fixed-header
    hide-default-footer
    item-key="_key"
    :sort-by="sortItems">
    <template v-slot:[`header.status`]="{column, getSortIcon, toggleSort}">
      <div class="v-data-table-header__content">
        <span>{{ column.title }}</span>
        <v-menu offset-y :close-on-content-click="false" v-model="menu">
          <template v-slot:activator="{props}">
            <DIconButton
              :parentProps="props"
              icon="mdi-filter-variant"
              :hint="t('TT_SHOW_FILTER')"
              :color="selectedFilterStatus && selectedFilterStatus.length > 0 ? 'primary' : 'secondary'"></DIconButton>
          </template>
          <div style="width: 320px" class="bg-background">
            <v-card>
              <v-row class="d-flex justify-end ma-1 mr-2">
                <DCloseButton @click="menu = false"></DCloseButton>
              </v-row>

              <v-select
                variant="outlined"
                v-model="selectedFilterStatus"
                density="compact"
                class="mx-2 pa-2"
                autofocus
                clearable
                :items="possibleStatus"
                :label="t('lbl_filter_on_status')"
                hide-details
                item-title="text"
                item-value="value"
                multiple
                menu
                transition="scale-transition"
                persistent-clear
                :list-props="{class: 'striped-filter-dd py-0'}">
                <template v-slot:item="{props, item}">
                  <v-list-item v-bind="props" :title="undefined" class="py-0 px-2">
                    <template v-slot:prepend="{isSelected}">
                      <v-checkbox hide-details :model-value="isSelected" />
                    </template>
                    <span :class="'prStatus' + item.value + ' prStatusFilter'">{{ item.title }}</span>
                  </v-list-item>
                </template>
                <template v-slot:selection="{item, index}">
                  <div v-if="index === 0" class="d-flex align-center">
                    <span :class="'prStatus' + item.value + ' prStatusFilter'">{{ item.title }}</span>
                  </div>
                  <span v-if="index === 1" class="pAdditionalFilter">
                    +{{ selectedFilterStatus.length - 1 }} others
                  </span>
                </template>
              </v-select>
            </v-card>
          </div>
        </v-menu>
        <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
      </div>
    </template>
    <template v-slot:[`header.type`]="{}">
      <v-row>
        <v-col v-for="(rule, i) in rules" :key="i" :class="viewTools.gridPolicyRulesAssignmentsHeaderClassByLanguage()">
          <v-icon small :color="getIconColorForPolicyType(rule)">{{ getIconForPolicyType(rule) }}</v-icon>
          {{ t(policyStateToTranslationKey(rule)) }}
        </v-col>
        <v-col>
          <span :class="viewTools.gridPolicyRulesAssignmentsHeaderClassByLanguage()">
            {{ t(policyStateToTranslationKey(PolicyState.NOT_SET)) }}
          </span>
        </v-col>
      </v-row>
    </template>
    <template v-slot:[`item.status`]="{item}">
      <span :class="'prStatus' + item.status">{{ t('PR_STATUS_' + item.status.toUpperCase()) }}</span>
    </template>
    <template v-slot:[`item.description`]="{item}">
      <span>{{ getStrWithMaxLength(120, item.description) }}</span>
      <Tooltip :text="item.description"></Tooltip>
    </template>
    <template v-slot:[`item.type`]="{item}">
      <v-radio-group
        inline
        v-model="item.type"
        :disabled="!edit || item.status === 'deprecated'"
        class="v-data-table--dense reduced-radio-height">
        <v-col v-for="(rule, i) in rules">
          <v-radio :key="i" :value="rule" :class="viewTools.gridPolicyRulesAssignmentsRowClassByLanguage()" />
        </v-col>
        <v-col>
          <v-radio :value="PolicyState.NOT_SET" :class="viewTools.gridPolicyRulesAssignmentsRowClassByLanguage()" />
        </v-col>
      </v-radio-group>
    </template>
  </v-data-table-virtual>
</template>

<script setup lang="ts">
import {IDefaultSelectItem} from '@disclosure-portal/model/IObligation';
import {PolicyRules, PolicyRulesAssignmentsDto, PolicyState} from '@disclosure-portal/model/PolicyRule';
import {useAppStore} from '@disclosure-portal/stores/app';
import useViewTools, {
  getIconColorForPolicyType,
  getIconForPolicyType,
  getStrWithMaxLength,
  policyStateToTranslationKey,
} from '@disclosure-portal/utils/View';
import {DataTableHeader, SortItem} from '@shared/types/table';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';

interface Props {
  loading?: boolean;
  edit?: boolean;
  isDialog?: boolean;
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  edit: false,
  isDialog: false,
});

const item = defineModel<PolicyRulesAssignmentsDto[]>({required: true});

const {t} = useI18n();
const appStore = useAppStore();
const viewTools = useViewTools();

const search = ref('');
const sortItems = ref<SortItem[]>([{key: 'name', order: 'asc'}]);
const rules = ref(PolicyRules);
const menu = ref(false);
const selectedFilterStatus = ref<string[]>(['active', 'inactive']);
const possibleStatus = ref<IDefaultSelectItem[]>([
  {
    text: t('PR_STATUS_ACTIVE'),
    value: 'active',
  },
  {
    text: t('PR_STATUS_INACTIVE'),
    value: 'inactive',
  },
  {
    text: t('PR_STATUS_DEPRECATED'),
    value: 'deprecated',
  },
]);

const headers = computed<DataTableHeader[]>(() => {
  return [
    {
      title: t('COL_STATUS'),
      align: 'start',
      sortable: true,
      value: 'status',
      width: '120',
    },
    {
      title: t('COL_NAME'),
      align: 'start',
      sortable: true,
      value: 'name',
      width: '120',
    },
    {
      title: t('COL_DESCRIPTION'),
      align: 'start',
      filterable: false,
      value: 'description',
      sortable: false,
      width: '400',
    },
    {
      title: '',
      align: 'center',
      filterable: false,
      value: 'type',
      width: typeWidth.value,
      sortable: false,
    },
  ];
});

const filteredHeaders = computed(() => {
  return props.isDialog
    ? headers.value.filter((header: DataTableHeader) => header.value !== 'description')
    : headers.value;
});

const filteredList = computed<PolicyRulesAssignmentsDto[]>(() => {
  if (!Array.isArray(item.value)) {
    return [];
  }
  return item.value.filter((pr) => selectedFilterStatus.value.some((s) => s == pr.status));
});

const typeWidth = computed(() => (appStore.getAppLanguage === 'en' ? 450 : 470));
</script>

<style scoped>
.v-radio-group {
  font-size: 0.8rem !important;
  line-height: 0.5 !important;
  height: 45px !important;
  margin-top: -15px;
}
</style>
