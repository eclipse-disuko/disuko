<script lang="ts" setup>
import {CustomId} from '@disclosure-portal/model/CustomId';
import {useAppStore} from '@disclosure-portal/stores/app';
import {useCustomIdStore} from '@disclosure-portal/stores/customid.store';
import useRules from '@disclosure-portal/utils/Rules';
import {DataTableHeader} from '@shared/types/table';
import {storeToRefs} from 'pinia';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const props = defineProps({
  readonly: {
    type: Boolean,
    required: false,
    default: false,
  },
});

const currentCustomIds = defineModel<CustomId[]>({required: true});

const {t} = useI18n();
const {minMax} = useRules();
const appStore = useAppStore();
const customIdsStore = useCustomIdStore();

const {customIds} = storeToRefs(customIdsStore);

const lang = computed(() => appStore.getAppLanguage);

const headers = computed<DataTableHeader[]>(() => {
  const res: DataTableHeader[] = [
    {
      title: t('COL_NAME'),
      width: 220,
      align: 'start',
      class: 'tableHeaderCell',
      value: 'technicalId',
      sortable: true,
    },
    {
      title: t('COL_VALUES'),
      align: 'start',
      class: 'tableHeaderCell',
      value: 'value',
      sortable: true,
    },
  ];
  if (!props.readonly) {
    res.push({
      title: t('COL_ACTIONS'),
      align: 'center',
      width: 150,
      class: 'tableHeaderCell',
      value: 'actions',
      sortable: false,
    });
  }
  return res;
});

const deleteAliasByIndex = (index: number) => {
  currentCustomIds.value = currentCustomIds.value.filter((_, i) => i !== index);
};

const createEmptyId = async () => {
  const customId = new CustomId();
  customId.technicalId = customIds.value.ids[0]?._key || '';
  currentCustomIds.value.unshift(customId);
};

const buildLink = (item: CustomId) => {
  const tmpl = customIds.value.map[item.technicalId].linkTemplate;
  return tmpl.replace('$VALUE$', item.value);
};

const activeRules = ref({
  value: minMax(t('COL_VALUES'), 5, 80, false),
  required: [(value: string) => !!value || t('VALIDATION_required')],
});
</script>

<template>
  <div class="pa-4 pt-6">
    <v-row class="header d-flex flex-row align-center ga-4" v-if="customIds.ids.length >= 1">
      <DCActionButton
        large
        :hint="t('TT_add_customid')"
        :text="t('BTN_ADD')"
        icon="mdi-plus"
        @click="createEmptyId"></DCActionButton>
    </v-row>
    <v-row>
      <v-col cols="12" xs="12" class="pt-4 px-0">
        <v-data-table
          :items="currentCustomIds"
          :headers="headers"
          density="compact"
          class="striped-table"
          fixed-header
          item-key="_key"
          :hide-default-footer="true">
          <template v-slot:[`item.technicalId`]="{item}">
            <v-select
              v-if="!item._key"
              class="pt-5"
              required
              density="compact"
              variant="outlined"
              :items="customIds.ids"
              :item-title="lang === 'en' ? 'name' : 'nameDE'"
              item-value="_key"
              :rules="activeRules.required"
              v-bind:menu-props="{location: 'bottom'}"
              v-model="item.technicalId">
              <template v-slot:item="{props, item}">
                <v-list-item v-bind="props">
                  <template v-slot:title>
                    <v-tooltip
                      :text="
                        lang === 'en' ? customIds.map[item.value].description : customIds.map[item.value].descriptionDE
                      "
                      :disabled="
                        lang === 'en'
                          ? !customIds.map[item.value].description
                          : !customIds.map[item.value].descriptionDE
                      "
                      location="right"
                      content-class="dpTooltip">
                      <template v-slot:activator="{props}">
                        <div v-bind="props">
                          <div>
                            {{ lang === 'en' ? customIds.map[item.value].name : customIds.map[item.value].nameDE }}
                          </div>
                          <div class="d-text d-secondary-text">
                            {{ item.value }}
                          </div>
                        </div>
                      </template>
                    </v-tooltip>
                  </template>
                </v-list-item>
              </template>
            </v-select>
            <template v-else>
              <div>
                {{ lang === 'en' ? customIds.map[item.technicalId].name : customIds.map[item.technicalId].nameDE }}
              </div>
              <div class="d-text d-secondary-text">
                {{ item.technicalId }}
              </div>
            </template>
          </template>
          <template v-slot:[`item.value`]="{item}">
            <v-text-field
              autocomplete="off"
              class="my-1 pt-5"
              v-if="!item._key"
              density="compact"
              variant="outlined"
              v-model="item.value"
              :rules="activeRules.value"></v-text-field>
            <span v-else>{{ item.value }}</span>
          </template>
          <template v-slot:[`item.actions`]="{item, index}" v-if="!readonly">
            <v-row justify="center">
              <DExternalLinkIcon
                v-if="customIds.map[item.technicalId].linkTemplate"
                :url="buildLink(item)"
                :hint="t('BTN_OPEN_LINK')"></DExternalLinkIcon>
              <span v-else class="action-icon-placeholder"></span>
              <DCopyClipboardButton :hint="t('TT_COPY_REFERENCE_INFO')" :content="item.value" />
              <DeleteConfirmationDialog
                v-slot="{showDialog}"
                :title="item.technicalId"
                @confirmed="() => deleteAliasByIndex(index)">
                <DIconButton icon="mdi-close" :hint="t('TT_delete_customid')" @clicked="showDialog" />
              </DeleteConfirmationDialog>
            </v-row>
          </template>
        </v-data-table>
      </v-col>
    </v-row>
  </div>
</template>
