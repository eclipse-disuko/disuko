// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {describe, expect, it} from 'vitest';
import {useTableSearch} from './useTableSearch';

describe('useTableSearch', () => {
  it('preserves the CLI callback API for nested fields and Vuetify raw items', () => {
    const filter = useTableSearch(['projectName', 'version.name', 'description']);
    const item = {projectName: 'Engine', version: {name: 'Release'}, description: null};

    expect(filter(null, 'ENGINE', {raw: item})).toBe(true);
    expect(filter(null, 'RELEASE', item)).toBe(true);
    expect(filter(null, 'missing', {raw: item})).toBe(false);
    expect(filter(null, '', undefined)).toBe(true);
    expect(filter(null, 'missing', undefined)).toBe(false);
    expect(filter(null, 'missing', {version: null})).toBe(false);
  });

  it('matches substrings across configured fields without searching other values', () => {
    const filter = useTableSearch(['creator', 'delegatedToFullName']);
    const item = {creator: 'Alice', delegatedToFullName: 'Bob Smith', hidden: 'secret'};

    expect(filter(null, 'LIC', {raw: item})).toBe(true);
    expect(filter(null, 'SMITH', {raw: item})).toBe(true);
    expect(filter('secret', 'secret', {raw: item})).toBe(false);
  });

  it('matches all rows for empty queries and no rows without searchable fields', () => {
    const filter = useTableSearch([]);

    expect(filter(null, '', {name: 'Alpha'})).toBe(true);
    expect(filter(null, 'Alpha', {name: 'Alpha'})).toBe(false);
  });

  it('evaluates row fields only once across multiple column calls for the same row and query', () => {
    let reads = 0;
    const raw = {
      get projectName() {
        reads += 1;
        return 'Engine';
      },
    };
    const filter = useTableSearch(['projectName']);

    expect(filter('col1', 'ENGINE', {raw})).toBe(true);
    expect(filter('col2', 'ENGINE', {raw})).toBe(true);
    expect(filter('col3', 'ENGINE', {raw})).toBe(true);
    expect(reads).toBe(1);

    expect(filter('col1', 'OTHER', {raw})).toBe(false);
    expect(reads).toBe(2);
  });

  it('matches nested array fields such as components and licenses', () => {
    const filter = useTableSearch(['components.componentName', 'licenses.licenseId']);
    const item = {
      components: [{componentName: 'lodash', componentVersion: '4.17.21'}],
      licenses: [{licenseName: 'MIT License', licenseId: 'MIT'}],
    };

    expect(filter(null, 'LODASH', {raw: item})).toBe(true);
    expect(filter(null, 'mit', {raw: item})).toBe(true);
    expect(filter(null, 'apache', {raw: item})).toBe(false);
  });
});
