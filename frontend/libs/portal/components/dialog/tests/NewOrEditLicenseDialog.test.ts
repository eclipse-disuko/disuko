// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {mount} from '@vue/test-utils';
import {afterAll, beforeAll, beforeEach, describe, expect, it, vi} from 'vitest';
import {nextTick} from 'vue';
import {config} from '@vue/test-utils';
import {vuetifyStubs} from '@disclosure-portal/test-utils/vuetify-stubs';
import {hasLeadingBlank} from '@disclosure-portal/utils/Validation';
import NewOrEditLicenseDialog from '../NewOrEditLicenseDialog.vue';

const {createMock, updateMock, getMock, getAllObligationsMock, snackbarInfoMock, snackbarErrorMock} = vi.hoisted(
  () => ({
    createMock: vi.fn(),
    updateMock: vi.fn(),
    getMock: vi.fn(),
    getAllObligationsMock: vi.fn(),
    snackbarInfoMock: vi.fn(),
    snackbarErrorMock: vi.fn(),
  }),
);

vi.mock('@disclosure-portal/services/license', () => ({
  default: {
    create: createMock,
    update: updateMock,
    get: getMock,
  },
}));

vi.mock('@disclosure-portal/services/admin', () => ({
  default: {
    getAllObligations: getAllObligationsMock,
  },
}));

vi.mock('@shared/composables/useSnackbar', () => ({
  default: () => ({
    info: snackbarInfoMock,
    error: snackbarErrorMock,
  }),
}));

describe('NewOrEditLicenseDialog', () => {
  describe('hasLeadingBlank', () => {
    it('returns true for values starting with whitespace', () => {
      expect(hasLeadingBlank(' test')).toBe(true);
      expect(hasLeadingBlank('\ttest')).toBe(true);
    });

    it('returns false for values without leading whitespace', () => {
      expect(hasLeadingBlank('test')).toBe(false);
      expect(hasLeadingBlank('')).toBe(false);
    });
  });

  const originalTMock = config.global.mocks?.$t;

  beforeAll(() => {
    if (config.global.mocks && '$t' in config.global.mocks) {
      delete config.global.mocks.$t;
    }
  });

  afterAll(() => {
    if (originalTMock) {
      config.global.mocks = {
        ...config.global.mocks,
        $t: originalTMock,
      };
    }
  });

  beforeEach(() => {
    createMock.mockReset();
    updateMock.mockReset();
    getMock.mockReset();
    getAllObligationsMock.mockReset();
    snackbarInfoMock.mockReset();
    snackbarErrorMock.mockReset();
  });

  const createWrapper = (isValid: boolean) => {
    const formValidate = vi.fn(async () => ({valid: isValid}));

    const wrapper = mount(NewOrEditLicenseDialog, {
      props: {
        mode: 'create',
        visible: true,
      },
      global: {
        stubs: {
          ...vuetifyStubs,
          Stack: {template: '<div><slot /></div>'},
          TabLinkObligations: {template: '<div />'},
          GridAliases: {template: '<div />'},
          DCloseButton: {template: '<button type="button" />'},
          'v-tabs-window': {template: '<div><slot /></div>'},
          'v-tabs-window-item': {template: '<div><slot /></div>'},
          'v-date-picker': {template: '<div />'},
          'v-form': {
            template: '<form><slot /></form>',
            methods: {
              validate: formValidate,
              reset: vi.fn(),
              resetValidation: vi.fn(),
            },
          },
        },
      },
    });

    return {wrapper, formValidate};
  };

  it('calls create when form is valid', async () => {
    getAllObligationsMock.mockResolvedValue({data: {items: []}});
    createMock.mockResolvedValue({});

    const {wrapper, formValidate} = createWrapper(true);
    await (wrapper.vm as {showDialog: () => Promise<void>}).showDialog();
    await nextTick();

    const buttons = wrapper.findAll('button');
    await buttons[buttons.length - 1].trigger('click');
    await nextTick();

    expect(formValidate).toHaveBeenCalledTimes(1);
    expect(createMock).toHaveBeenCalledTimes(1);
    expect(snackbarInfoMock).toHaveBeenCalledTimes(1);
  });

  it('does not call create when form is invalid (validation fails)', async () => {
    getAllObligationsMock.mockResolvedValue({data: {items: []}});

    const {wrapper, formValidate} = createWrapper(false);
    await (wrapper.vm as {showDialog: () => Promise<void>}).showDialog();
    await nextTick();

    const buttons = wrapper.findAll('button');
    await buttons[buttons.length - 1].trigger('click');
    await nextTick();

    expect(formValidate).toHaveBeenCalledTimes(1);
    expect(createMock).not.toHaveBeenCalled();
  });
});