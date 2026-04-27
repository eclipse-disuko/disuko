// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {DataTableHeader, DataTableHeaderFilterItems} from '@shared/types/table';
import {useStorage} from '@vueuse/core';
import {computed, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {InternalDataTableHeader} from 'vuetify/lib/components/VDataTable/types';
import {useHeaderSettingsStore} from '@shared/stores/headerSettings.store';

const STORE_KEY = 'gridHeaderSettings';

type UseHeaderSettingsParams = {
  tableName: string;
  initiallyHiddenList?: string[];
  headers?: DataTableHeader[];
  settingsColumn?: InternalDataTableHeader;
};

export const useHeaderSettings = (props: UseHeaderSettingsParams) => {
  const headerSettingsStore = useHeaderSettingsStore();

  const {t} = useI18n();
  const localStorage = useStorage(STORE_KEY, {} as Record<string, number[]>);

  const tableName = ref<string>('');

  const initialSelectedHeaders = computed(
    (): number[] => headerSettingsStore.$state[tableName.value]?.initialSelectedHeaders ?? [],
  );
  const headers = computed((): DataTableHeader[] => headerSettingsStore.$state[tableName.value]?.headers ?? []);
  const settingsColumn = computed(
    (): InternalDataTableHeader | undefined => headerSettingsStore.$state[tableName.value]?.settingsColumn,
  );
  const selectedHeaders = computed((): number[] => localStorage.value[tableName.value] ?? []);
  const selectableHeaders = computed((): DataTableHeaderFilterItems[] =>
    headers.value.map(
      (header, headerIndex) =>
        ({
          ...header,
          text: headers.value[headerIndex].title.includes(',')
            ? getMultiTitle(headers.value[headerIndex].title)
            : t(headers.value[headerIndex].title),
          disabled: settingsColumn.value?.value === headers.value[headerIndex]?.value,
        }) as DataTableHeaderFilterItems,
    ),
  );
  const filteredHeaders = computed(() =>
    headerSettingsStore.$state[tableName.value]
      ? selectedHeaders.value
          .toSorted((a, b) => a - b)
          .map((columnNumber) => {
            // This is needed so the tooltip can be translated instantly
            const tooltipObject = headerSettingsStore.$state?.[tableName.value]?.headers[columnNumber]?.tooltipText
              ? {tooltipText: t(headerSettingsStore.$state[tableName.value].headers[columnNumber].tooltipText ?? '')}
              : {};

            const titleSet = headerSettingsStore.$state?.[tableName.value]?.headers[columnNumber]?.title.includes(',')
              ? headerSettingsStore.$state[tableName.value].headers[columnNumber].title.split(',')
              : [headerSettingsStore.$state?.[tableName.value]?.headers[columnNumber]?.title];

            const title = headerSettingsStore.$state?.[tableName.value]?.headers[columnNumber]?.title
              ? titleSet.map((titlePart) => t(titlePart.trim())).join(' ')
              : '';

            // This is needed so the title can be translated instantly
            const titleObject = {title: title};

            return {
              ...headerSettingsStore.$state[tableName.value].headers[columnNumber],
              ...titleObject,
              ...tooltipObject,
            } as DataTableHeader;
          })
      : [],
  );

  const getMultiTitle = (title: string) =>
    title
      .split(',')
      .map((part) => t(part.trim()))
      .join(' ');

  const resetSelectedHeaders = () => {
    updateSelectedHeaders(headerSettingsStore.$state[tableName.value].initialSelectedHeaders);
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

  const resetHeaderSettings = (newProps: UseHeaderSettingsParams) => {
    tableName.value = newProps.tableName;

    if (!headerSettingsStore.$state[tableName.value]) {
      headerSettingsStore.$patch({
        ...headerSettingsStore.$state,
        [tableName.value]: {
          headers: [],
          hideInitially: newProps.initiallyHiddenList || [],
          initialSelectedHeaders: [],
        },
      });
    }

    const headersBefore = [...(headerSettingsStore.$state[tableName.value].headers ?? [])];

    if (newProps.headers) {
      headerSettingsStore.$state[tableName.value].headers = newProps.headers;

      headerSettingsStore.$state[tableName.value].initialSelectedHeaders =
        headerSettingsStore.$state[tableName.value].hideInitially.length <= 0
          ? [...selectableHeaders.value.keys()]
          : selectableHeaders.value
              .map((header, index) => ({header, index}))
              .filter(({index}) =>
                newProps.headers?.[index]?.value
                  ? !headerSettingsStore.$state[tableName.value].hideInitially.includes(newProps.headers[index].value!)
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

    if (headerSettingsStore.$state[tableName.value] && newProps.settingsColumn) {
      headerSettingsStore.$state[tableName.value].settingsColumn = newProps.settingsColumn;
    }
  };

  resetHeaderSettings(props);

  return {
    headers,
    selectedHeaders,
    filteredHeaders,
    selectableHeaders,
    initialSelectedHeaders,
    settingsColumn,
    resetHeaderSettings,
    resetSelectedHeaders,
    updateSelectedHeaders,
    updateSelectedHeadersFromStringList,
  };
};
