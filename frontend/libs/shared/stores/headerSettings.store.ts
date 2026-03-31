// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {DataTableHeader, DataTableHeaderFilterItems} from '@shared/types/table';
import {useStorage} from '@vueuse/core';
import {defineStore} from 'pinia';
import {computed, reactive, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {InternalDataTableHeader} from 'vuetify/lib/components/VDataTable/types';

type HeaderSettings = {
  headers: DataTableHeader[];
  hideInitially: string[];
  initialSelectedHeaders: number[];
  settingsColumn?: InternalDataTableHeader;
};

const STORE_KEY = 'gridHeaderSettings';

export const useHeaderSettingsStore = defineStore(STORE_KEY, () => {
  const {t} = useI18n();
  const localStorage = useStorage(STORE_KEY, {} as Record<string, number[]>);

  const tableName = ref<string>('');

  const state = reactive<Record<string, HeaderSettings>>({});

  const initialSelectedHeaders = computed(() => state[tableName.value]?.initialSelectedHeaders ?? []);
  const headers = computed(() => state[tableName.value]?.headers ?? []);
  const settingsColumn = computed(() => state[tableName.value]?.settingsColumn);
  const selectedHeaders = computed(() => localStorage.value[tableName.value] ?? []);
  const selectableHeaders = computed((): DataTableHeaderFilterItems[] =>
    headers.value.map(
      (header) =>
        ({
          ...header,
          disabled: header?.value === settingsColumn.value?.value || header?.value === 'actionButtons',
          text: header.title.includes(',') ? getMultiTitle(header.title) : t(header.title),
        }) as DataTableHeaderFilterItems,
    ),
  );

  const setupStore = (newTableName: string, newHeaders?: DataTableHeader[], initiallyHiddenList: string[] = []) => {
    tableName.value = newTableName;

    if (!state[tableName.value]) {
      state[tableName.value] = {hideInitially: initiallyHiddenList} as HeaderSettings;
    }

    const headersBefore = [...(state[tableName.value].headers ?? [])];

    if (newHeaders) {
      state[tableName.value].headers = newHeaders;

      state[tableName.value].initialSelectedHeaders =
        state[tableName.value].hideInitially.length <= 0
          ? [...selectableHeaders.value.keys()]
          : selectableHeaders.value
              .map((header, index) => ({header, index}))
              .filter(({index}) =>
                newHeaders?.[index]?.value
                  ? !state[tableName.value].hideInitially.includes(newHeaders[index].value!)
                  : true,
              )
              .map(({index}) => index);

      if (headersBefore.length >= 1) {
        return;
      }

      if (localStorage.value?.[tableName.value]) {
        updateSelectedHeaders(localStorage.value[tableName.value]);
      } else {
        if (!localStorage.value?.[tableName.value]) {
          localStorage.value[tableName.value] = [] as number[];
        }
        resetSelectedHeaders();
      }
    }
  };

  const setSettingsColumn = (settingsColumn: InternalDataTableHeader) => {
    if (state[tableName.value]) {
      state[tableName.value].settingsColumn = settingsColumn;
    }
  };

  const getMultiTitle = (title: string) =>
    title
      .split(',')
      .map((part) => t(part.trim()))
      .join(' ');

  const resetSelectedHeaders = () => {
    updateSelectedHeaders(state[tableName.value].initialSelectedHeaders);
  };

  const updateSelectedHeadersFromStringList = (newHeaders: string[]) => {
    const newHeadersFromStringList = selectableHeaders.value
      .filter((header) =>
        Boolean(
          newHeaders.find((newHeader) => {
            // This should be impossible, but for security purposes we also check if it is the settings-column
            const isSettingsColumn = header.value === settingsColumn.value?.value;
            return newHeader === header.value || isSettingsColumn;
          }),
        ),
      )
      .map((header) => selectableHeaders.value.indexOf(header));

    updateSelectedHeaders(newHeadersFromStringList);
  };

  const updateSelectedHeaders = (newHeaders: number[]) => {
    localStorage.value[tableName.value] = newHeaders;
  };

  const filteredHeaders = computed(() =>
    state[tableName.value]
      ? [
          ...selectedHeaders.value
            .toSorted((a, b) => a - b)
            .map((columnNumber) => {
              // This is needed so the tooltip can be translated instantly
              const tooltipObject = state?.[tableName.value]?.headers[columnNumber]?.tooltipText
                ? {tooltipText: t(state[tableName.value].headers[columnNumber].tooltipText ?? '')}
                : {};

              const titleSet = state?.[tableName.value]?.headers[columnNumber]?.title.includes(',')
                ? state[tableName.value].headers[columnNumber].title.split(',')
                : [state?.[tableName.value]?.headers[columnNumber]?.title];

              const title = state?.[tableName.value]?.headers[columnNumber]?.title
                ? titleSet.map((titlePart) => t(titlePart.trim())).join(' ')
                : '';

              // This is needed so the title can be translated instantly
              const titleObject = {title: title};

              return {
                ...state[tableName.value].headers[columnNumber],
                ...titleObject,
                ...tooltipObject,
              } as DataTableHeader;
            }),
        ]
      : [],
  );

  return {
    headers,
    selectedHeaders,
    filteredHeaders,
    selectableHeaders,
    initialSelectedHeaders,
    settingsColumn,
    setupStore,
    setSettingsColumn,
    resetSelectedHeaders,
    updateSelectedHeaders,
    updateSelectedHeadersFromStringList,
  };
});
