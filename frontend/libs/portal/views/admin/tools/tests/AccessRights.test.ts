// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {mountView} from '@disclosure-portal/test-utils/view-test-utils';
import {flushPromises} from '@vue/test-utils';
import {beforeEach, describe, expect, it, vi} from 'vitest';
import AccessRights from '../AccessRights.vue';

const {getAllProjectAccessRightsMock, getAllAccessRightsMock} = vi.hoisted(() => ({
  getAllProjectAccessRightsMock: vi.fn(),
  getAllAccessRightsMock: vi.fn(),
}));

vi.mock('@disclosure-portal/services/admin', () => ({
  default: {
    getAllProjectAccessRights: getAllProjectAccessRightsMock,
    getAllAccessRights: getAllAccessRightsMock,
  },
}));

const crud = {read: true, update: false, create: false, delete: false};

const projectAccessRights = {
  Owner: {
    allowProjectFoo: crud,
  },
  Supplier: {
    allowProjectFoo: crud,
  },
};

const accessRights = {
  ApplicationAdmin: {
    allowUsers: crud,
  },
  DomainAdmin: {
    allowUsers: crud,
  },
  Internal: {
    allowUsers: crud,
  },
};

describe('AccessRights', () => {
  beforeEach(() => {
    getAllProjectAccessRightsMock.mockReset();
    getAllAccessRightsMock.mockReset();
    getAllProjectAccessRightsMock.mockResolvedValue({data: projectAccessRights});
    getAllAccessRightsMock.mockResolvedValue({data: accessRights});
  });

  it('renders without throwing', async () => {
    const {wrapper} = mountView(AccessRights);
    await flushPromises();

    expect(wrapper.exists()).toBe(true);
  });

  it('loads project and further access rights on mount and builds the tables', async () => {
    const {wrapper} = mountView(AccessRights);
    await flushPromises();

    expect(getAllProjectAccessRightsMock).toHaveBeenCalledTimes(1);
    expect(getAllAccessRightsMock).toHaveBeenCalledTimes(1);

    const vm = wrapper.vm as unknown as {
      accessRightsLoaded: boolean;
      allProjectAccessRightsForTable: {accessRightName: string}[];
      allAccessRightsForTable: {accessRightName: string}[];
      allRoles: {value: string; text: string}[];
      selectedRoles: string[];
    };
    expect(vm.accessRightsLoaded).toBe(true);
    expect(vm.allProjectAccessRightsForTable.map((e) => e.accessRightName)).toContain('allowProjectFoo');
    expect(vm.allAccessRightsForTable.map((e) => e.accessRightName)).toContain('allowUsers');
    // The default role selection only covers these three groups per the component's constant.
    expect(vm.selectedRoles).toEqual(['ApplicationAdmin', 'DomainAdmin', 'Internal']);
  });

  it('re-filters the further-access-rights table when the role selection changes', async () => {
    const {wrapper} = mountView(AccessRights);
    await flushPromises();

    const vm = wrapper.vm as unknown as {
      selectedRoles: string[];
      dirtySelection: boolean;
      defaultSelection: boolean;
      prepareDataForSelection: () => void;
      accessRightsHeaderForSelection: {value: string}[];
    };
    // Column header per selected role (3 by default), plus the leading "access rights" name column.
    expect(vm.accessRightsHeaderForSelection).toHaveLength(4);

    vm.selectedRoles = ['ApplicationAdmin'];
    vm.prepareDataForSelection();
    await flushPromises();

    expect(vm.accessRightsHeaderForSelection.map((h) => h.value)).toEqual(['accessRightName', 'ApplicationAdmin']);
    expect(vm.dirtySelection).toBe(false);
    expect(vm.defaultSelection).toBe(false);
  });
});
