export type DataTableHeader = {
  key?: string;
  value?: SelectItemKey;
  /**
   * Title for the column.
   *
   * You can use one key or comma-separated keys to provide multiple keys for i18n.
   *
   * Example one key:
   * ``` ts
   * {
   *   title: 'header_name',
   *   // ...
   * }
   * ```
   *
   * Example multiple keys:
   * ``` ts
   * {
   *   title: 'header_name, header_surname',
   *   // ...
   * }
   * ```
   */
  title: string;
  tooltipText?: string;
  colspan?: number;
  rowspan?: number;
  fixed?: boolean;
  /**
   * @deprecated filterable is not supported anymore.
   */
  filterable?: boolean;
  selectable?: boolean;
  align?: 'start' | 'end' | 'center';
  width?: string | number;
  minWidth?: string | number;
  maxWidth?: string | number;
  sortable?: boolean;
  /**
   * @deprecated class is not supported anymore.
   */
  class?: string;
  sort?: DataTableCompareFunction;
  sortRaw?: DataTableCompareFunction;
};

export type DataTableHeaderFilterItems = {
  text?: string;
  value: string;
  icon?: string;
  iconColor?: string;
  chip?: string;
  chipColor?: string;
  disabled?: boolean;
};

export type DataTabelIndex = {
  searchIndex: string;
};

export interface SortItem {
  key: string;
  order: 'asc' | 'desc';
}

export interface DataTableItem<T> {
  item: T;
}
