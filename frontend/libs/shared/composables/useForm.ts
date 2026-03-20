// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

export function useForm() {
  const rules = {
    required: (value: string) => !!value || 'Required.',
    requiredObject: (value: object) =>
      value && typeof value === 'object' ? Object.keys(value)?.length > 0 || 'Required.' : false,
  };

  /**
   * Converts a mapping object to an array of objects with name and value properties.
   * @example Input { "key1": "Value 1", "key2": "Value 2" }, Output [{ name: "Value 1", value: "key1" }, { name: "Value 2", value: "key2" }]
   * @param mapping
   */
  const getMappingNameAndValue = (mapping: Record<string, string>) => {
    return Object.entries(mapping).map(([key, value]) => ({
      name: value,
      value: key,
    }));
  };

  return {
    rules,
    getMappingNameAndValue,
  };
}
