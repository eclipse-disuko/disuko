<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {DialogReviewRemarkConfig} from '@disclosure-portal/components/dialog/DialogConfigs';
import {type Project, SbomLicenses} from '@disclosure-portal/model/Project';
import {LicenseMeta, ReviewRemark, ReviewRemarkLevel, ReviewRemarkRequest} from '@disclosure-portal/model/Quality';
import {ReviewTemplate} from '@disclosure-portal/model/ReviewTemplate';
import {ComponentInfoSlim, SpdxFile} from '@disclosure-portal/model/VersionDetails';
import projectService from '@disclosure-portal/services/projects';
import versionService from '@disclosure-portal/services/version';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {useSbomStore} from '@disclosure-portal/stores/sbom.store';
import useRules from '@disclosure-portal/utils/Rules';
import {formatDateAndTime} from '@disclosure-portal/utils/Table';
import DialogLayout, {DialogLayoutConfig} from '@shared/layouts/DialogLayout.vue';
import useSnackbar from '@shared/composables/useSnackbar';
import _ from 'lodash';
import {computed, Ref, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {VForm} from 'vuetify/components';

const level: ReviewRemarkLevel[] = [ReviewRemarkLevel.GREEN, ReviewRemarkLevel.YELLOW, ReviewRemarkLevel.RED];

const {t} = useI18n();
const sbomStore = useSbomStore();
const projectStore = useProjectStore();
const {minMax} = useRules();
const {info: snack} = useSnackbar();
const emit = defineEmits(['reload']);

const rules = {
  name: minMax(t('NPV_DIALOG_TF_TITLE'), 5, 80, false),
  description: minMax(t('NP_DIALOG_TF_DESCRIPTION'), 10, 700, false),
  level: [(v: string) => v !== ReviewRemarkLevel.NOT_SET || t('DLG_LEVEL_REQUIRED')],
};

const isVisible = ref(false);
const config = ref<DialogReviewRemarkConfig>({} as DialogReviewRemarkConfig);
const item = ref<ReviewRemark>(new ReviewRemark());
const form = ref<VForm | null>(null);

const templates = ref<ReviewTemplate[]>([]);
const selectedTemplate = ref<ReviewTemplate | undefined>(undefined);
const templatesLoading = ref(false);

const selectedSbom = ref<SpdxFile | undefined>(undefined);
const sboms = ref<SpdxFile[]>([]);
const sbomsLoading = ref(false);

const comps = ref<ComponentInfoSlim[]>([]);
const compsLoading = ref(false);
const selectedComponents = ref<(ComponentInfoSlim | null)[]>([null]);

const sbomAllLicenses = ref<SbomLicenses | undefined>(undefined);
const licensesLoading = ref(false);
const selectedLicenses = ref<(LicenseMeta | null)[]>([null]);

const dialogActionDisabled = ref(false);

const licenses = computed((): LicenseMeta[] => {
  if (!sbomAllLicenses.value) {
    return [];
  }
  const known = sbomAllLicenses.value.known
    ? sbomAllLicenses.value.known.map((license) => ({
        licenseId: license.id,
        licenseName: license.name,
      }))
    : [];
  const unknown = sbomAllLicenses.value.unknown
    ? sbomAllLicenses.value.unknown.map((str) => ({licenseId: str, licenseName: str}))
    : [];

  return [...known, ...unknown];
});

const projectModel = computed((): Project => projectStore.currentProject!);
const version = computed(() => sbomStore.getCurrentVersion);
const versionID = computed(() => config.value.versionID || version.value._key);

const dialogConfig = computed((): DialogLayoutConfig => {
  const isEdit = !!config.value.presetItem;
  return {
    title: isEdit ? t('UM_DIALOG_TITLE_EDIT_REVIEW_REMARK') : t('UM_DIALOG_TITLE_NEW_REVIEW_REMARK'),
    primaryButton: {
      text: isEdit ? t('NP_DIALOG_BTN_EDIT') : t('NP_DIALOG_BTN_CREATE'),
      disabled: dialogActionDisabled.value,
      loading: dialogActionDisabled.value,
    },
    secondaryButton: {
      text: t('BTN_CANCEL'),
      disabled: dialogActionDisabled.value,
    },
  };
});

const open = (newConf: DialogReviewRemarkConfig) => {
  config.value = newConf;

  if (config.value.presetItem) {
    const preset = config.value.presetItem;
    item.value = _.cloneDeep(preset);

    config.value.spdxID = preset.sbomId;

    if (preset.components && preset.components.length) {
      config.value.components = preset.components.map((c) => {
        const slim = new ComponentInfoSlim();
        slim.spdxId = c.componentId;
        slim.name = c.componentName;
        slim.version = c.componentVersion;
        slim.licenseExpression = '';
        slim.componentInfo = [];
        return slim;
      });
    }

    if (preset.licenses && preset.licenses.length) {
      config.value.licenses = preset.licenses.map((l) => {
        const meta = new LicenseMeta();
        meta.licenseId = l.licenseId;
        meta.licenseName = l.licenseName;
        return meta;
      });
    }
  }

  loadTemplates();
  loadSboms();

  isVisible.value = true;
};

const setupAfterSbomsLoaded = () => {
  sbomsLoading.value = false;
  if (!config.value.spdxID) {
    return;
  }
  selectedSbom.value = sboms.value.find((sbom) => sbom._key === config.value.spdxID);
  if (config.value.components && config.value.components.length) {
    selectedComponents.value = [...config.value.components];
  }
  loadLicenses().then(() => {
    if (config.value.licenses && config.value.licenses.length) {
      selectedLicenses.value = config.value.licenses.map(
        (presetLicense) => licenses.value.find((l) => l.licenseId === presetLicense.licenseId) || presetLicense,
      );
    }
  });
};

const loadSboms = async () => {
  sbomsLoading.value = true;

  if (versionID.value === sbomStore.currentVersion._key) {
    // Use cached channel history from the store
    sboms.value = sbomStore.channelSpdxs;
    setupAfterSbomsLoaded();
  } else {
    // Dialog opened for a different version — fall back to direct fetch
    versionService.getSbomHistory(projectModel.value._key, versionID.value).then((res) => {
      const spdxFileHistory = res.data;
      if (spdxFileHistory[0]) {
        spdxFileHistory[0].isRecent = true;
      }
      sboms.value = spdxFileHistory;
      setupAfterSbomsLoaded();
    });
  }
};

const loadTemplates = () => {
  templatesLoading.value = true;
  projectService.getReviewTemplates(projectModel.value._key).then((res) => {
    templates.value = res.data;
    templates.value.sort((a, b) => {
      const titleA = a.title.toLowerCase();
      const titleB = b.title.toLowerCase();
      if (titleA < titleB) return -1;
      if (titleA > titleB) return 1;
      return 0;
    });
    templatesLoading.value = false;
  });
};

const templateChanged = () => {
  if (!selectedTemplate.value) {
    return;
  }
  item.value.level = selectedTemplate.value.level;
  item.value.title = selectedTemplate.value.title;
  item.value.description = selectedTemplate.value.description;
};

const sbomChanged = () => {
  selectedComponents.value = [null];
  selectedLicenses.value = [null];
  comps.value = [];
  loadLicenses();
};
const wait = 300;
const debouncedSearch = _.debounce((query: string) => {
  if (!query) {
    comps.value = [];
    return;
  }
  compsLoading.value = true;
  versionService
    .getVersionComponentsBySearch(projectModel.value._key, versionID.value, selectedSbom.value!._key, query)
    .then((res) => {
      comps.value = res;
      compsLoading.value = false;
    });
}, wait);
const compSearchChanged = async (query: string) => debouncedSearch(query);

const loadLicenses = async () => {
  if (!selectedSbom.value) {
    return;
  }
  licensesLoading.value = true;
  return versionService
    .getVersionSbomAllLicenses(projectModel.value._key, versionID.value, selectedSbom.value!._key)
    .then((res) => {
      sbomAllLicenses.value = res;
      licensesLoading.value = false;
    });
};

const close = () => {
  form.value?.reset();
  isVisible.value = false;
};

const doDialogAction = async () => {
  if (dialogActionDisabled.value) {
    return;
  }

  const validationResult = await form.value?.validate();
  if (!validationResult?.valid) {
    return;
  }

  dialogActionDisabled.value = true;
  try {
    item.value.sbomId = selectedSbom.value?._key || '';

    const reviewRemarkRequest = ReviewRemarkRequest.toRequest(item.value);

    reviewRemarkRequest.components = Array.from(
      new Set(
        selectedComponents.value
          .filter((c): c is ComponentInfoSlim => c !== null)
          .map((c) => c.spdxId)
          .filter((id) => !!id),
      ),
    );
    reviewRemarkRequest.licenses = Array.from(
      new Set(
        selectedLicenses.value
          .filter((c): c is LicenseMeta => c !== null)
          .map((l) => l.licenseId)
          .filter((id) => !!id),
      ),
    );

    if (config.value.presetItem) {
      await versionService.editReviewRemark(
        projectModel.value._key,
        versionID.value,
        config.value.presetItem.key,
        reviewRemarkRequest,
      );
      snack(t('DIALOG_remark_edit_success'));
    } else {
      await versionService.createReviewRemark(projectModel.value._key, versionID.value, reviewRemarkRequest);
      snack(t('DIALOG_remark_create_success'));
    }
    emit('reload');
    form.value?.reset();
    isVisible.value = false;
  } finally {
    dialogActionDisabled.value = false;
  }
};

const canAddComponent = computed(() => {
  const last = selectedComponents.value.at(-1);
  return !!last && !!last.spdxId;
});

const addComponent = () => {
  selectedComponents.value.push(null);
};

const canAddLicense = computed(() => {
  const last = selectedLicenses.value.at(-1);
  return !!last && !!last.licenseId;
});

const addLicense = () => {
  selectedLicenses.value.push(null);
};

function watchSingleEmpty<T>(listRef: Ref<(T | null)[]>) {
  watch(
    listRef,
    (newArr) => {
      // nichts tun, wenn gar kein null drinsteckt
      if (!newArr.some((c) => c === null)) return;

      const endsWithNull = newArr[newArr.length - 1] === null;
      const filtered = newArr.filter((c) => c !== null) as T[];
      const normalized = endsWithNull ? [...filtered, null] : filtered;

      // update nur, wenn sich wirklich etwas ändert
      const sameLength = normalized.length === newArr.length;
      const sameItems = sameLength && normalized.every((v, i) => v === newArr[i]);
      if (!sameItems) {
        listRef.value = normalized;
      }
    },
    {deep: true},
  );
}

watchSingleEmpty(selectedComponents);
watchSingleEmpty(selectedLicenses);
defineExpose({open});
</script>

<template>
  <v-dialog v-model="isVisible" content-class="large" width="1400" height="800">
    <DialogLayout :config="dialogConfig" @secondary-action="close" @primary-action="doDialogAction" @close="close">
      <v-form ref="form">
        <div class="grid grid-cols-1 md:grid-cols-[1fr_auto_1fr] gap-4">
          <Stack class="w-full h-min">
            <h3>{{ t('TAB_TITLE_GENERAL') }}</h3>
            <v-select
              item-text="title"
              return-object
              variant="outlined"
              clearable
              :label="t('REVIEW_REMARK_TEMPLATE_DIALOG')"
              v-bind:menu-props="{location: 'bottom'}"
              :items="templates"
              :loading="templatesLoading"
              v-model="selectedTemplate"
              @update:model-value="templateChanged"
              hide-details="auto" />
            <v-text-field
              autocomplete="off"
              required
              v-model="item.title"
              :rules="rules.name"
              :label="t('NPV_DIALOG_TF_TITLE')"
              autofocus
              hide-details="auto"
              variant="outlined" />
            <v-textarea
              v-model="item.description"
              no-resize
              variant="outlined"
              :label="t('NP_DIALOG_TF_DESCRIPTION') + '*'"
              :rules="rules.description"
              :counter="700"
              hide-details="auto" />
            <v-select
              :items="level"
              v-model="item.level"
              v-bind:menu-props="{location: 'bottom'}"
              :rules="rules.level"
              :label="t('COL_LEVEL')"
              variant="outlined"
              hide-details="auto"
              required>
              <template v-slot:item="{props}">
                <v-list-item v-bind="props">
                  <template v-slot:prepend="{isSelected}">
                    <v-checkbox hide-details :model-value="isSelected" />
                  </template>
                  <template v-slot:title="{title}">
                    <span>{{ t('REMARK_LEVEL_' + title) }}</span>
                  </template>
                </v-list-item>
              </template>
              <template v-slot:selection="{item}">
                <span>{{ item.raw ? t('REMARK_LEVEL_' + item.raw) : '' }}</span>
              </template>
            </v-select>
          </Stack>
          <v-divider vertical></v-divider>
          <Stack class="w-full h-min md:max-h-[580px] overflow-hidden md:overflow-auto">
            <h3>{{ t('COL_REFERENCES') }}</h3>
            <v-select
              v-model="selectedSbom"
              variant="outlined"
              :label="t('UM_DIALOG_REVIEW_REMARK_SBOM')"
              :items="sboms"
              return-object
              item-title="_key"
              clearable
              :loading="sbomsLoading"
              hide-details="auto"
              @update:model-value="sbomChanged">
              <template v-slot:item="{props, item}">
                <v-list-item v-bind="props">
                  <template v-slot:title>
                    <span class="d-subtitle-2">{{ formatDateAndTime(item.raw.Uploaded) }}&nbsp;</span>
                    <span class="d-text d-secondary-text">&nbsp;-&nbsp;{{ item.raw.MetaInfo.Name }}</span>
                    <span class="d-text d-secondary-text" v-if="item.raw.Tag">&nbsp;({{ item.raw.Tag }})</span>
                    <span class="d-text d-secondary-text" v-if="item.raw.isRecent"
                      >&nbsp;{{ '[' + t('SBOM_LATEST') + ']' }}</span
                    >
                    <span class="d-text d-secondary-text" v-else>&nbsp;{{ '[' + t('SBOM_FORMER') + ']' }}</span>
                  </template>
                </v-list-item>
              </template>
              <template v-slot:selection="{item}">
                <span class="d-subtitle-2">{{ formatDateAndTime(item.raw.Uploaded) }}&nbsp;</span>
                <span class="d-text d-secondary-text">&nbsp;-&nbsp;{{ item.raw.MetaInfo.Name }}</span>
                <span class="d-text d-secondary-text" v-if="item.raw.Tag">&nbsp;({{ item.raw.Tag }})</span>
                <span class="d-text d-secondary-text" v-if="item.raw.isRecent"
                  >&nbsp;{{ '[' + t('SBOM_LATEST') + ']' }}</span
                >
                <span class="d-text d-secondary-text" v-else>&nbsp;{{ '[' + t('SBOM_FORMER') + ']' }}</span>
              </template>
            </v-select>
            <v-autocomplete
              v-for="(_, index) in selectedComponents"
              :key="index"
              v-model="selectedComponents[index]"
              clearable
              :label="t('labelSearchComponent')"
              :disabled="!selectedSbom"
              :items="comps"
              @update:search="compSearchChanged"
              return-object
              item-title="name"
              :loading="compsLoading"
              variant="outlined"
              hide-details="auto">
              <template v-slot:item="{item, props}">
                <v-list-item v-bind="props" title="">
                  <span class="d-subtitle-2 ml-2">{{ item.raw.name }}</span>
                  <span class="d-text d-secondary-text">&nbsp;({{ item.raw.version }})</span>
                </v-list-item>
              </template>
              <template v-slot:selection="{item}">
                <div class="d-inline">
                  <span class="d-subtitle-2 ml-2">{{ item.raw.name }}</span>
                  <span class="d-text d-secondary-text">&nbsp;({{ item.raw.version }})</span>
                </div>
              </template>
            </v-autocomplete>
            <div
              v-if="canAddComponent"
              class="d-flex align-center border-md border-dashed border-opacity-25 p-3 mb-6"
              @click="addComponent">
              <v-icon color="primary">mdi-plus</v-icon>
              <span class="font-weight-light pl-1">{{ t('RR_DIALOG_MORE_COMPONENT') }}</span>
            </div>

            <v-autocomplete
              v-for="(_, index) in selectedLicenses"
              :key="index"
              v-model="selectedLicenses[index]"
              clearable
              :label="t('labelSearchLicense')"
              :disabled="!selectedSbom"
              :items="licenses"
              return-object
              item-title="name"
              :loading="licensesLoading"
              variant="outlined"
              hide-details="auto">
              <template v-slot:item="{item, props}">
                <v-list-item v-bind="props" title="">
                  <span class="d-subtitle-2 ml-2">{{ item.raw.licenseName }}</span>
                  <span class="d-text d-secondary-text">&nbsp;({{ item.raw.licenseId }})</span>
                </v-list-item>
              </template>
              <template v-slot:selection="{item}">
                <div class="d-inline">
                  <span class="d-subtitle-2 ml-2">{{ item.raw.licenseName }}</span>
                  <span class="d-text d-secondary-text">&nbsp;({{ item.raw.licenseId }})</span>
                </div>
              </template>
            </v-autocomplete>
            <div
              v-if="canAddLicense"
              class="d-flex align-center border-md border-dashed border-opacity-25 p-3 mb-6"
              @click="addLicense">
              <v-icon color="primary">mdi-plus</v-icon>
              <span class="font-weight-light pl-1">{{ t('RR_DIALOG_MORE_LICENSE') }}</span>
            </div>
          </Stack>
        </div>
      </v-form>
    </DialogLayout>
  </v-dialog>
</template>
