// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import LegalNoticeDe from '@shared/assets/documents/legal_notice/LegalNotice_de.md?raw';
import LegalNoticeEn from '@shared/assets/documents/legal_notice/LegalNotice_en.md?raw';
import PrivacyStatementDe from '@shared/assets/documents/privacy_statement/PrivacyStatement_de.md?raw';
import PrivacyStatementEn from '@shared/assets/documents/privacy_statement/PrivacyStatement_en.md?raw';
import ProviderDe from '@shared/assets/documents/provider/Provider_de.md?raw';
import ProviderEn from '@shared/assets/documents/provider/Provider_en.md?raw';
import NoticeDe from '@shared/assets/documents/provider_privacy_notice/ProviderPrivacyNotice_de.html?raw';
import NoticeEn from '@shared/assets/documents/provider_privacy_notice/ProviderPrivacyNotice_en.html?raw';
import TermsOfUseEn from '@shared/assets/documents/terms_of_use/TermsOfUseCurrent.md?raw';
import TermsOfUseDe from '@shared/assets/documents/terms_of_use/TermsOfUseDe.md?raw';
import {vuetifyStubs} from '@disclosure-portal/test-utils/vuetify-stubs';
import {config, mount} from '@vue/test-utils';
import {afterAll, beforeAll, beforeEach, describe, expect, it, vi} from 'vitest';
import {ref} from 'vue';
import ProviderPrivacyDialog from '../ProviderPrivacyDialog.vue';

const {useLanguageStoreMock} = vi.hoisted(() => ({useLanguageStoreMock: vi.fn()}));

// See PolicyRulesTable.test.ts: useLanguageStore() calls Vuetify's useLocale(), which needs a
// real Vuetify plugin instance. Mock the store module directly rather than pull one in.
vi.mock('@shared/stores/language.store', () => ({
  useLanguageStore: useLanguageStoreMock,
}));

const MarkdownStub = {
  props: ['text', 'id'],
  template: '<div class="markdown" :id="id">{{ text }}</div>',
};

const DCopyClipboardButtonStub = {
  props: ['hint', 'content'],
  template: '<button class="d-copy-clipboard-button" type="button" />',
};

const DCloseButtonStub = {
  template: '<button class="d-close-button" type="button" @click="$emit(\'click\')" />',
};

describe('ProviderPrivacyDialog', () => {
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
    vi.clearAllMocks();
    useLanguageStoreMock.mockReturnValue({appLanguage: ref('en')});
  });

  const createWrapper = () => {
    return mount(ProviderPrivacyDialog, {
      slots: {
        default: '<template #default="{ showDialog }"><button class="open-btn" @click="showDialog">open</button></template>',
      },
      global: {
        stubs: {
          ...vuetifyStubs,
          'v-tabs-window': {template: '<div><slot /></div>'},
          'v-tabs-window-item': {template: '<div><slot /></div>'},
          Markdown: MarkdownStub,
          DCopyClipboardButton: DCopyClipboardButtonStub,
          DCloseButton: DCloseButtonStub,
        },
      },
    });
  };

  it('does not render the dialog content before it is opened', () => {
    const wrapper = createWrapper();

    expect(wrapper.find('.v-dialog').exists()).toBe(false);
  });

  it('opens the dialog when the slot-exposed showDialog is invoked', async () => {
    const wrapper = createWrapper();

    await wrapper.find('.open-btn').trigger('click');

    expect(wrapper.find('.v-dialog').exists()).toBe(true);
  });

  it('renders a translated tab per document', async () => {
    const wrapper = createWrapper();
    await wrapper.find('.open-btn').trigger('click');

    const tabs = wrapper.findAllComponents(vuetifyStubs['v-tab']);
    expect(tabs).toHaveLength(5);
    expect(tabs.map((tab) => tab.text())).toEqual([
      'Provider',
      'Legal Notices',
      'Privacy Statement',
      'Terms of Use',
      'Third Party Notices',
    ]);
  });

  it('renders the English document content for each tab by default', async () => {
    const wrapper = createWrapper();
    await wrapper.find('.open-btn').trigger('click');

    const markdowns = wrapper.findAllComponents(MarkdownStub);
    expect(markdowns.map((md) => md.props('text'))).toEqual([
      ProviderEn,
      LegalNoticeEn,
      PrivacyStatementEn,
      TermsOfUseEn,
      NoticeEn,
    ]);
  });

  it('renders the German document content when appLanguage is de', async () => {
    useLanguageStoreMock.mockReturnValue({appLanguage: ref('de')});
    const wrapper = createWrapper();
    await wrapper.find('.open-btn').trigger('click');

    const markdowns = wrapper.findAllComponents(MarkdownStub);
    expect(markdowns.map((md) => md.props('text'))).toEqual([ProviderDe, LegalNoticeDe, PrivacyStatementDe, TermsOfUseDe, NoticeDe]);
  });

  it('closes the dialog when DCloseButton is clicked', async () => {
    const wrapper = createWrapper();
    await wrapper.find('.open-btn').trigger('click');
    expect(wrapper.find('.v-dialog').exists()).toBe(true);

    await wrapper.find('.d-close-button').trigger('click');

    expect(wrapper.find('.v-dialog').exists()).toBe(false);
  });
});
