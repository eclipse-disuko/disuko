// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

/**
 * Composable for creating custom filter functions for Vuetify data tables
 * Provides flexible searching across multiple fields in nested objects
 */

/**
 * Creates a custom filter function for data tables that searches across specified fields
 * @param searchableFields - Array of field paths to search (supports nested paths with dot notation)
 * @returns Filter function compatible with Vuetify v-data-table :custom-filter prop
 *
 * @example
 * // Simple usage with flat fields
 * const customFilter = useTableFilter(['name', 'email', 'status']);
 *
 *
 * @example
 * // In component
 * const customFilter = useTableFilter(['projectName', 'version.name', 'description']);
 * <v-data-table :custom-filter="customFilter" ... />
 */
export const useTableFilter = (searchableFields: string[]) => {
  return (value: unknown, query: string, item: unknown): boolean => {
    if (!query) return true;

    const searchText = query.toLowerCase();
    const raw = (item as any)?.raw || item;

    return searchableFields.some((fieldPath) => {
      const fieldValue = fieldPath.split('.').reduce((obj, key) => obj?.[key], raw);
      return String(fieldValue || '')
        .toLowerCase()
        .includes(searchText);
    });
  };
};
