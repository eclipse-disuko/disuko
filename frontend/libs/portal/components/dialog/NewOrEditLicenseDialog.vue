<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script lang="ts" setup>
import TabLinkObligations from '@disclosure-portal/components/dialog/classification/TabLinkObligations.vue';
import GridAliases from '@disclosure-portal/components/grids/GridAliases.vue';
import {useLicense} from '@disclosure-portal/composables/useLicense';
import License, {ITextValue} from '@disclosure-portal/model/License';
import licenseService from '@disclosure-portal/services/license';
import {DiscoForm} from '@disclosure-portal/types/discobasics';
import useRules from '@disclosure-portal/utils/Rules';
import {isSpdxIdentifier, isURLOrEmpty} from '@disclosure-portal/utils/Validation';
import DCloseButton from '@shared/components/disco/DCloseButton.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import dayjs from 'dayjs';
import {computed, nextTick, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';

const {getLicenseFamily, getLicenseTypes, getLicenseApprovalTypes, getLicenseReviewStates} = useLicense();
const licenseFamily: ITextValue[] = getLicenseFamily();
const licenseTypes: ITextValue[] = getLicenseTypes();
const licenseApproval: ITextValue[] = getLicenseApprovalTypes();
const reviewStates: ITextValue[] = getLicenseReviewStates();
const visible = defineModel<boolean>('visible');
const emits = defineEmits<{
  (e: 'update:modelValue'): void;
  (e: 'closed:successfully', id: string): void;
}>();
const props = defineProps({
  initialData: {
    type: Object as () => License,
    required: false,
    default: undefined,
    validator(value: unknown, rawProps: any): boolean {
      return (
        rawProps.mode === 'create' ||
        ((rawProps.mode === 'edit' || rawProps.mode === 'duplicate') && value !== undefined)
      );
    },
  },
  mode: {
    type: String,
    required: false,
    default: 'create',
    validator(value: any, rawProps: any): boolean {
      return ['create', 'edit', 'duplicate'].includes(value);
    },
  },
  testId: {
    type: String,
    required: false,
    default: 'licenses-editor',
  },
});

const isLoading = ref(false);
const isOpening = ref(false);
const hasOpenedObligationsTab = ref(false);
const hasOpenedAliasesTab = ref(false);

const {t} = useI18n();

const title = props.initialData ? 'LM_DIALOG_TITLE_EDIT_LICENSE' : 'LM_DIALOG_TITLE_CREATE_LICENSE';

const displayedTab = ref(0);
const selectedTab = ref(0);
watch(selectedTab, async () => {
  if (selectedTab.value === displayedTab.value) return;

  if (!formDialog.value) return;

  const isValid = (await formDialog.value.validate()).valid;

  await nextTick(() => {
    if (isValid) {
      displayedTab.value = selectedTab.value;
      hasOpenedObligationsTab.value ||= displayedTab.value === 2;
      hasOpenedAliasesTab.value ||= displayedTab.value === 5;
    } else {
      selectedTab.value = displayedTab.value;
    }
  });
});

const item = ref<License>({
  _key: '',
  active: false,
  aliases: [],
  created: '',
  isDeprecatedLicenseId: false,
  licenseId: '',
  meta: {
    approvalState: '',
    changelog: '',
    evaluation: '',
    family: '',
    isLicenseChart: false,
    legalComments: '',
    licenseType: '',
    licenseUrl: '',
    obligationsList: [],
    reviewDate: '',
    reviewState: '',
    sourceUrl: '',
    prevalentClassificationLevel: '',
    classifications: [],
  },
  name: '',
  source: '',
  text: '',
  type: 0,
  updated: '',
});

const isDatePickerVisible = ref(false);
const licenseTextArea = ref<any | null>(null);
// the date picker expects a date object, but this application does only work with strings, so we need a proxy object which does the conversion
const dateProxy = computed({
  get: () => new Date(item.value.meta.reviewDate),
  set: (value: Date) => (item.value.meta.reviewDate = dayjs(value).format().split('T')[0]),
});

const formError = '';
const duplicateIdMessage = '';

const rules = useRules();
const snackbar = useSnackbar();
const FLOW_BREADCRUMB_KEY = 'disuko.license.flow';

const writeFlowBreadcrumb = (step: string, details: Record<string, unknown> = {}) => {
  try {
    const entriesRaw = localStorage.getItem(FLOW_BREADCRUMB_KEY);
    const entries = entriesRaw ? JSON.parse(entriesRaw) : [];
    entries.push({
      ts: new Date().toISOString(),
      step,
      mode: props.mode,
      details,
    });
    localStorage.setItem(FLOW_BREADCRUMB_KEY, JSON.stringify(entries.slice(-80)));
  } catch (error) {
    console.error('Failed to write license flow breadcrumb', error);
  }
};

const toErrorDetails = (error: unknown) => {
  if (error instanceof Error) {
    return {message: error.message, stack: error.stack};
  }
  try {
    return {message: JSON.stringify(error)};
  } catch {
    return {message: String(error)};
  }
};

writeFlowBreadcrumb('module-loaded');

const licenseNameRules = [
  ...rules.required(t('COL_NAME')),
  (value: string) => (value && value.length >= 3) || t('LICENSE_NAME_VALIDATION_MIN_LENGTH', {min: 3}),
  (value: string) => (value && value.length <= 100) || t('LICENSE_NAME_VALIDATION_MAX_LENGTH', {max: 100}),
];
const licenseIdRules = rules.required(t('COL_LICENSE_ID'));
licenseIdRules.push((value) => isSpdxIdentifier(value, true) || t('ERROR_VALIDATION_LICENSE_ID_IS_NOT_VALID'));
const licenseFamilyRules = rules.requiredOrEmpty(t('COL_LICENSE_FAMILY'));
const licenseTypeRules = rules.requiredOrEmpty(t('COL_LICENSE_TYPE'));

const urlRules = [(value: string) => isURLOrEmpty(value) || t('VALIDATION_url')];

const reset = async () => {
  writeFlowBreadcrumb('reset:start');
  const licenseRequest =
    props.initialData && (props.mode === 'edit' || props.mode === 'duplicate')
      ? licenseService.get(props.initialData.licenseId)
      : null;
  selectedTab.value = 0;
  displayedTab.value = 0;

  item.value = new License({
    _key: '',
    active: false,
    aliases: [],
    created: '' + new Date().toISOString(),
    updated: '' + new Date().toISOString(),
    isDeprecatedLicenseId: false,
    licenseId: '',
    meta: {
      approvalState: '',
      changelog: '',
      evaluation: '',
      family: '',
      isLicenseChart: false,
      legalComments: '',
      licenseType: '',
      licenseUrl: '',
      obligationsList: [],
      reviewDate: '',
      reviewState: '',
      sourceUrl: '',
      classifications: [],
      prevalentClassificationLevel: '',
    },
    name: '',
    source: '',
    text: '',
    type: 0,
  });

  if (licenseRequest) {
    writeFlowBreadcrumb('reset:load-license-detail:start', {licenseId: props.initialData!.licenseId});
    const licenseDetail = (await licenseRequest).data;

    // duplicate content
    item.value._key = licenseDetail._key;
    item.value.name = licenseDetail.name;
    item.value.licenseId = licenseDetail.licenseId;
    item.value.active = licenseDetail.active;
    item.value.type = licenseDetail.type;
    item.value.isDeprecatedLicenseId = licenseDetail.isDeprecatedLicenseId;
    item.value.text = licenseDetail.text;
    item.value.source = licenseDetail.source;
    item.value.meta.family = licenseDetail.meta.family;
    item.value.meta.licenseType = licenseDetail.meta.licenseType;
    item.value.meta.obligationsList = licenseDetail.meta.obligationsList ?? [];
    item.value.meta.legalComments = licenseDetail.meta.legalComments;
    item.value.meta.evaluation = licenseDetail.meta.evaluation;
    item.value.meta.approvalState = licenseDetail.meta.approvalState;
    item.value.meta.reviewState = licenseDetail.meta.reviewState;
    item.value.meta.reviewDate = licenseDetail.meta.reviewDate;
    item.value.meta.licenseUrl = licenseDetail.meta.licenseUrl;
    item.value.meta.sourceUrl = licenseDetail.meta.sourceUrl;
    item.value.meta.isLicenseChart = licenseDetail.meta.isLicenseChart;
    item.value.meta.changelog = licenseDetail.meta.changelog;
    item.value.aliases = licenseDetail.aliases;

    // revert duplication for some properties if duplicate
    if (props.mode === 'duplicate') {
      item.value._key = '';
      item.value.name = '';
      item.value.licenseId = 'LicenseRef-MB-' + item.value.licenseId;
      item.value.aliases = [];
      item.value.source = 'custom';
      item.value.meta.approvalState = 'pending';
      item.value.meta.reviewState = '';
      item.value.meta.reviewDate = '';
      item.value.meta.changelog =
        'License derived from "' + licenseDetail.name + '" (' + licenseDetail.licenseId + ') \n';
    }
  }
  writeFlowBreadcrumb('reset:done');
};

/**
 * Users need to be able to paste encoded text into the textarea.
 * E.g. '\u003cCopyright Information\u003e'
 * should be pasted as '<Copyright Information>'
 * thus this method gets the text that is about to be pasted and decodes it via JSON.parse if necessary
 *
 * @param e ClipboardEvent
 */
const decodeIfNecessary = (e: ClipboardEvent) => {
  // @ts-expect-error window.clipboardData is not typed
  const paste = (e.clipboardData || window.clipboardData).getData('text/plain');
  if (paste === undefined) {
    console.error('Could not get clipboard data.');
    return;
  }
  try {
    const decodedPaste = JSON.parse(`{"text": "${paste}"}`).text;
    if (paste === decodedPaste) {
      // decoding was not necessary, thus manual event handling is not needed
      return;
    }
    e.stopPropagation();
    e.preventDefault();

    const textAreaRoot = licenseTextArea.value?.$el ?? licenseTextArea.value;
    const textAreaEl: HTMLTextAreaElement | null = textAreaRoot?.querySelector?.('textarea') ?? null;
    if (!textAreaEl) {
      console.error('Could not resolve license textarea element.');
      return;
    }
    const contentBefore = textAreaEl.value.substring(0, textAreaEl.selectionStart);
    const contentAfter = textAreaEl.value.substring(textAreaEl.selectionEnd);
    item.value.text = `${contentBefore}${decodedPaste}${contentAfter}`;

    // cursor position is lost when setting the text directly, we can restore it manually
    const cursorPositionAfterPaste =
      textAreaEl.selectionEnd - (textAreaEl.selectionEnd - textAreaEl.selectionStart) + decodedPaste.length;
    nextTick(() => {
      // this has to be done in $nextTick
      textAreaEl.setSelectionRange(cursorPositionAfterPaste, cursorPositionAfterPaste);
    });
  } catch (decodeError) {
    console.error('Pasted text could or did not have to be decoded', decodeError);
  }
};

const closeDialog = () => {
  visible.value = false;
};

const showDialog = async () => {
  if (isOpening.value) return;

  isOpening.value = true;
  writeFlowBreadcrumb('show-dialog:start');
  visible.value = true;
  hasOpenedObligationsTab.value = false;
  hasOpenedAliasesTab.value = false;

  try {
    await reset();
    writeFlowBreadcrumb('show-dialog:visible-true');
  } finally {
    isOpening.value = false;
  }
};

watch(
  () => item.value.meta.reviewState,
  (newValue) => {
    if (newValue === 'reviewed') {
      item.value.meta.reviewDate = new Date().toISOString().split('T')[0];
    } else {
      item.value.meta.reviewDate = '';
    }
  },
);

const formDialog = ref<DiscoForm | null>(null);
const doDialogAction = async () => {
  if (isLoading.value) return; // Verhindert mehrfaches Klicken
  isLoading.value = true;
  writeFlowBreadcrumb('submit:start');

  item.value.name = item.value.name.trim();
  item.value.licenseId = item.value.licenseId.trim();

  const result = await formDialog.value!.validate();
  writeFlowBreadcrumb('submit:validated', {valid: result.valid});
  if (result.valid) {
    // create a unique id for all new aliases
    item.value.aliases
      ?.filter((alias) => !alias._key)
      .forEach((alias) => (alias._key = alias.licenseId + '-' + new Date().getTime()));

    try {
      if (props.mode === 'create' || props.mode === 'duplicate') {
        writeFlowBreadcrumb('submit:create:request');
        await licenseService.create(item.value!);
        writeFlowBreadcrumb('submit:create:success', {licenseId: item.value.licenseId});
        snackbar.info(t('DIALOG_license_create_success'));
      } else if (props.mode === 'edit') {
        writeFlowBreadcrumb('submit:update:request', {key: item.value._key});
        await licenseService.update(item.value!, item.value!._key);
        writeFlowBreadcrumb('submit:update:success', {licenseId: item.value.licenseId});
        snackbar.info(t('DIALOG_license_edit_success'));
      }
      closeDialog();
      emits('closed:successfully', item.value.licenseId);
    } catch (error: any) {
      writeFlowBreadcrumb('submit:error', toErrorDetails(error));
      console.error(error);
      snackbar.error(t('DIALOG_license_create_error', {error: error.message}));
    } finally {
      isLoading.value = false;
      writeFlowBreadcrumb('submit:finally');
    }
  } else {
    isLoading.value = false;
    writeFlowBreadcrumb('submit:invalid');
  }
};
defineExpose({
  showDialog,
});
</script>

<template>
  <div>
    <slot :showDialog="showDialog"> </slot>

    <v-dialog v-model="visible" scrollable persistent width="80%">
      <v-form ref="formDialog">
        <v-card class="pa-8" ref="NElicenseDialog" :data-testid="testId">
          <v-card-title>
            <v-row>
              <v-col cols="10" class="d-flex align-center">
                <span class="text-h5">
                  {{ t(title) }}
                  <template v-if="item && item.name">
                    <q>{{ item.name }}</q>
                  </template>
                </span>
              </v-col>
              <v-col cols="2" align="right">
                <DCloseButton @click="closeDialog" />
              </v-col>
            </v-row>
          </v-card-title>
          <v-card-text class="expanding-container">
            <Stack v-if="item && item.meta && item.meta.obligationsList">
              <v-tabs
                v-model="selectedTab"
                slider-color="brand"
                active-class="active"
                show-arrows
                bg-color="tabsHeader">
                <v-tab :value="0" id="general">
                  {{ t('TAB_TITLE_GENERAL') }}
                </v-tab>
                <v-tab :value="1">
                  {{ t('TAB_TITLE_LICENSE') }}
                </v-tab>
                <v-tab :value="2">
                  {{ t('TAB_TITLE_CLASSIFICATIONS') }}
                </v-tab>
                <v-tab :value="3">
                  {{ t('TAB_TITLE_LEGALCOMMENTS') }}
                </v-tab>
                <v-tab :value="4">
                  {{ t('TAB_TITLE_CHANGELOG') }}
                </v-tab>
                <v-tab :value="5">
                  {{ t('TAB_TITLE_ALIASES') }}
                </v-tab>
                <v-tab :value="6">
                  {{ t('TAB_TITLE_EVALUATION') }}
                </v-tab>
              </v-tabs>
              <v-tabs-window class="pa-4" v-model="displayedTab">
                <v-tabs-window-item :key="0" :value="0">
                  <v-row>
                    <v-col cols="12">
                      <v-row justify="start" class="pt-2">
                        <v-col
                          v-if="mode === 'create' || mode === 'duplicate'"
                          sm="6"
                          md="6"
                          class="errorBorder">
                          <v-text-field
                            autocomplete="off"
                            variant="outlined"
                            :label="t('COL_NAME')"
                            v-model="item.name"
                            class="required"
                            hide-details="auto"
                            :rules="licenseNameRules"
                            autofocus
                            required />
                        </v-col>
                        <v-col
                          v-if="mode === 'edit' && item.source !== 'spdx'"
                          sm="6"
                          md="6"
                          class="errorBorder">
                          <v-text-field
                            autocomplete="off"
                            :readonly="item.source !== 'custom'"
                            variant="outlined"
                            :label="t('COL_NAME')"
                            v-model="item.name"
                            class="required"
                            hide-details="auto"
                            :rules="licenseNameRules"
                            autofocus
                            required />
                        </v-col>
                        <v-col
                          v-if="mode === 'edit' && item.source === 'spdx'"
                          sm="6"
                          md="6"
                          class="errorBorder">
                          <v-text-field
                            autocomplete="off"
                            readonly
                            variant="outlined"
                            :label="t('COL_NAME')"
                            v-model="item.name"
                            class="required"
                            hide-details="auto"
                            required />
                        </v-col>
                        <v-col
                          v-if="mode === 'create' || mode === 'duplicate'"
                          sm="6"
                          md="6"
                          class="errorBorder">
                          <v-text-field
                            autocomplete="off"
                            variant="outlined"
                            class="required"
                            v-model="item.licenseId"
                            :rules="licenseIdRules"
                            :error-messages="duplicateIdMessage"
                            hide-details="auto"
                            :label="t('COL_LICENSE_ID')"
                            required></v-text-field>
                        </v-col>
                        <v-col
                          v-if="mode === 'edit' && item.source !== 'spdx'"
                          sm="6"
                          md="6"
                          class="errorBorder">
                          <v-text-field
                            autocomplete="off"
                            variant="outlined"
                            class="required"
                            v-model="item.licenseId"
                            :rules="licenseIdRules"
                            :error-messages="duplicateIdMessage"
                            hide-details="auto"
                            :label="t('COL_LICENSE_ID')"
                            required></v-text-field>
                        </v-col>
                        <v-col v-if="mode === 'edit' && item.source === 'spdx'" sm="6" md="6">
                          <v-text-field
                            autocomplete="off"
                            variant="outlined"
                            readonly
                            v-model="item.licenseId"
                            hide-details="auto"
                            :label="t('COL_LICENSE_ID')"></v-text-field>
                        </v-col>
                      </v-row>
                      <v-row>
                        <v-col cols="12" sm="6" md="6" lg="6" class="errorBorder">
                          <v-select
                            :items="licenseTypes"
                            v-model="item.meta.licenseType"
                            :label="t('COL_LICENSE_TYPE')"
                            class="required"
                            :rules="licenseTypeRules"
                            variant="outlined"
                            item-value="value"
                            item-title="text"
                            hide-details="auto"
                            required></v-select>
                        </v-col>
                        <v-col cols="12" sm="6" md="6" lg="6" class="errorBorder">
                          <v-select
                            class="required"
                            :rules="licenseFamilyRules"
                            :items="licenseFamily"
                            v-model="item.meta.family"
                            :label="t('COL_LICENSE_FAMILY')"
                            variant="outlined"
                            item-value="value"
                            item-title="text"
                            hide-details="auto"
                            required></v-select>
                        </v-col>
                      </v-row>
                      <v-row justify="start" class="">
                        <v-col cols="12" sm="6" md="6" lg="6">
                          <v-select
                            :items="licenseApproval"
                            v-model="item.meta.approvalState"
                            :label="t('COL_APPROVAL_STATUS')"
                            variant="outlined"
                            item-value="value"
                            item-title="text"
                            hide-details="auto"></v-select>
                        </v-col>
                        <v-col cols="12" sm="6" md="6">
                          <v-select
                            :items="reviewStates"
                            v-model="item.meta.reviewState"
                            :label="t('COL_REVIEW_STATUS')"
                            variant="outlined"
                            item-value="value"
                            item-title="text"
                            hide-details="auto"></v-select>
                        </v-col>
                        <v-col
                          cols="12"
                          v-if="item.meta.reviewState === 'reviewed' && item.meta.reviewDate"
                          sm="6"
                          md="6">
                          <v-menu :close-on-content-click="false" location="bottom">
                            <template v-slot:activator="{props}">
                              <v-text-field
                                autocomplete="off"
                                class="cursor-pointer"
                                v-bind="props"
                                v-model="item.meta.reviewDate"
                                :label="t('COL_REVIEW_DATE')"
                                append-inner-icon="mdi-calendar"
                                readonly
                                @click:append="isDatePickerVisible = true"
                                hide-details="auto"
                                variant="outlined"></v-text-field>
                            </template>
                            <v-date-picker
                              v-model="dateProxy"
                              first-day-of-week="1"
                              color="primary"
                              hide-header
                              scrollable />
                          </v-menu>
                        </v-col>
                        <v-col sm="6" md="6" v-if="mode !== 'create'" lg="6">
                          <v-text-field
                            autocomplete="off"
                            variant="outlined"
                            readonly
                            disabled
                            v-model="item.source"
                            hide-details="auto"
                            :label="t('COL_LICENSE_SOURCE')"></v-text-field>
                        </v-col>
                      </v-row>
                      <v-row>
                        <v-col cols="12">
                          <v-text-field
                            autocomplete="off"
                            variant="outlined"
                            v-model="item.meta.licenseUrl"
                            hide-details="auto"
                            :rules="urlRules"
                            :label="t('COL_LICENSE_URL')"></v-text-field>
                        </v-col>
                        <v-col cols="12">
                          <v-text-field
                            autocomplete="off"
                            variant="outlined"
                            v-model="item.meta.sourceUrl"
                            hide-details="auto"
                            :rules="urlRules"
                            :label="t('COL_SOURCE_URL')"></v-text-field>
                        </v-col>
                        <v-col cols="12">
                          <v-switch
                            v-model="item.meta.isLicenseChart"
                            hide-details
                            color="primary"
                            :label="t('LICENSE_CHART_TITLE')"
                            class="mt-0 mr-2 shrink"></v-switch>
                        </v-col>
                        <v-col cols="12" v-if="formError">
                          <span class="text-[rgb(var(--v-theme-error))]">{{ t(formError) }}</span>
                        </v-col>
                      </v-row>
                    </v-col>
                  </v-row>
                </v-tabs-window-item>
                <v-tabs-window-item :key="1" :value="1" class="expanding-container no-scrollbar">
                  <v-textarea
                    ref="licenseTextArea"
                    class="expand pt-2"
                    rows="20"
                    no-resize
                    variant="outlined"
                    v-model="item.text"
                    hide-details="auto"
                    @paste="decodeIfNecessary"
                    :label="t('CD_LICENSE_TEXT')" />
                </v-tabs-window-item>
                <v-tabs-window-item
                  v-if="displayedTab === 2 || hasOpenedObligationsTab"
                  v-show="displayedTab === 2"
                  :key="2"
                  :value="2">
                  <TabLinkObligations ref="tabContent2" v-model:obligations="item.meta.obligationsList" />
                </v-tabs-window-item>
                <v-tabs-window-item :value="3" :key="3" class="expanding-container no-scrollbar" ref="tabItemWYSIWYG">
                  <v-textarea
                    class="expand pt-2"
                    rows="20"
                    no-resize
                    variant="outlined"
                    v-model="item.meta.legalComments"
                    hide-details="auto"
                    :label="t('TAB_TITLE_LEGALCOMMENTS')" />
                </v-tabs-window-item>
                <v-tabs-window-item :value="4" :key="4" class="expanding-container" style="overflow: hidden !important">
                  <v-textarea
                    class="expand pt-2"
                    no-resize
                    rows="20"
                    variant="outlined"
                    v-model="item.meta.changelog"
                    hide-details="auto"
                    :label="t('CD_LICENSE_CHANGELOG')" />
                </v-tabs-window-item>
                <v-tabs-window-item
                  v-if="displayedTab === 5 || hasOpenedAliasesTab"
                  v-show="displayedTab === 5"
                  :value="5"
                  :key="5">
                  <GridAliases ref="tabContentAliases" v-model:license="item" mode="edit" />
                </v-tabs-window-item>
                <v-tabs-window-item :value="6" :key="6" class="expanding-container">
                  <v-textarea
                    class="expand pt-2"
                    rows="20"
                    no-resize
                    variant="outlined"
                    v-model="item.meta.evaluation"
                    hide-details="auto"
                    :label="t('TAB_TITLE_EVALUATION')" />
                </v-tabs-window-item>
              </v-tabs-window>
            </Stack>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn @click="closeDialog" class="secondary mr-8" variant="plain" color="primary" size="small">
              <span>{{ t('BTN_CANCEL') }}</span>
            </v-btn>
            <v-btn
              @click="doDialogAction"
              class="primary confirmButton"
              variant="flat"
              color="primary"
              size="small"
              :disabled="isLoading">
              <v-progress-circular v-if="isLoading" indeterminate color="white" size="20"></v-progress-circular>
              <span v-else>
                <span v-if="mode === 'create' || mode === 'duplicate'">{{ t('NP_DIALOG_BTN_CREATE') }}</span>
                <span v-if="mode === 'edit'">{{ t('NP_DIALOG_BTN_EDIT') }}</span>
              </span>
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-form>
    </v-dialog>
  </div>
</template>

<style>
.cursor-pointer {
  input {
    cursor: pointer;
  }

  i {
    cursor: pointer;
  }
}
</style>
