// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export const useTableSearch = (searchableFields: string[]) => {
  const fieldPaths = searchableFields.map((fieldPath) => fieldPath.split('.'));
  let lastRaw: unknown;
  let lastQuery = '';
  let lastResult = false;

  return (_value: unknown, query: string, item: unknown): boolean => {
    if (!query) return true;

    const raw = (item as {raw?: unknown})?.raw ?? item;
    if (raw !== undefined && raw === lastRaw && query === lastQuery) {
      return lastResult;
    }

    const searchText = query.toLowerCase();
    const result = fieldPaths.some((fieldPath) => {
      const fieldValue = fieldPath.reduce<unknown>((object, key) => {
        if (Array.isArray(object)) {
          return object.flatMap((entry) => (entry != null ? [(entry as Record<string, unknown>)[key]] : []));
        }
        return object != null ? (object as Record<string, unknown>)[key] : undefined;
      }, raw);
      return String(fieldValue ?? '')
        .toLowerCase()
        .includes(searchText);
    });

    lastRaw = raw;
    lastQuery = query;
    lastResult = result;
    return result;
  };
};
