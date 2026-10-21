// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {describe, expect, it} from 'vitest';
import {createGroupURL, createProjectURL, createSBOMURL, createVersionURL} from './apiUrls';

describe('dashboard URLs', () => {
  it.each([
    {key: 'item-id', encoded: 'item-id'},
    {key: 'id/with space?#%', encoded: 'id%2Fwith%20space%3F%23%25'},
  ])('builds routes with encoded identifiers: $key', ({key, encoded}) => {
    expect(createProjectURL(key)).toBe(`/dashboard/projects/${encoded}/overview`);
    expect(createGroupURL(key)).toBe(`/dashboard/groups/${encoded}/overview`);
    expect(createVersionURL(key, key)).toBe(`/dashboard/projects/${encoded}/versions/${encoded}`);
    expect(createSBOMURL(key, key, key)).toBe(
      `/dashboard/projects/${encoded}/versions/${encoded}/component/${encoded}`,
    );
  });

  it('keeps project, version, and SBOM identifiers in their respective segments', () => {
    expect(createVersionURL('project-id', 'version-id')).toBe('/dashboard/projects/project-id/versions/version-id');
    expect(createSBOMURL('project-id', 'version-id', 'sbom-id')).toBe(
      '/dashboard/projects/project-id/versions/version-id/component/sbom-id',
    );
  });
});
