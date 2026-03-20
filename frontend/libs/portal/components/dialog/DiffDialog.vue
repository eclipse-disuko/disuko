<script setup lang="ts">
import License, {LicenseDiff} from '@disclosure-portal/model/License';
import adminService from '@disclosure-portal/services/admin';
import {useIdleStore} from '@disclosure-portal/stores/idle.store';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import DCloseButton from '@shared/components/disco/DCloseButton.vue';
import {CodeDiff} from 'v-code-diff';
import {nextTick, reactive, ref} from 'vue';
import {useI18n} from 'vue-i18n';

class DiffDetails {
  public field = '';
  public oldValue = '';
  public newValue = '';

  constructor(field: string, oldValue: string, newValue: string) {
    this.field = field;
    this.oldValue = oldValue;
    this.newValue = newValue;
  }
}

const {t} = useI18n();
const idle = useIdleStore();

const show = ref(false);
const diffs = ref<LicenseDiff[]>([]);
const diffDetails = ref<DiffDetails[]>([]);
const currentPosition = ref(1);
const licenseId = ref('');
const oldLicense = reactive<License>({} as License);
const newLicense = reactive<License>({} as License);
const headerAdded = ref(false);

const open = async () => {
  idle.show();
  show.value = true;

  diffDetails.value = [];
  await reload();
  await nextTick();

  addCustomHeader();
  idle.hide();
};

const addCustomHeader = () => {
  if (headerAdded.value) {
    return;
  }
  const table: HTMLTableElement | null = document.querySelector('.file-diff-split.diff-table');
  if (!table) {
    console.error('Diff table could not be found!');
    return;
  }

  const row = table.insertRow(0);
  const beforeCell = row.insertCell(0);
  beforeCell.colSpan = 2;
  beforeCell.innerHTML = `<h2>${t('DIFF_BEFORE')}</h2>`;
  const afterCell = row.insertCell(1);
  afterCell.colSpan = 2;
  afterCell.innerHTML = `<h2>${t('DIFF_AFTER')}</h2>`;
  headerAdded.value = true;
};

const reload = async () => {
  diffs.value = await adminService.getLicensesDiffs();
  if (diffs.value.length === 0) {
    diffDetails.value = [];
    close();
  } else {
    currentPosition.value = 1;
    prepareDataForCurrentPosition();
  }
};
function toSpdxDeprecatedOrEmpty(value: boolean): string {
  return value ? t('SPDX_STATUS_VALUE') : '';
}
const prepareDataForCurrentPosition = () => {
  diffDetails.value = [];
  const diff = diffs.value[currentPosition.value - 1];
  licenseId.value = diff.licenseId;

  Object.assign(oldLicense, diff.oldLicense);
  Object.assign(newLicense, diff.newLicense);

  if (oldLicense.name !== newLicense.name) {
    diffDetails.value.push(new DiffDetails('' + t('COL_LICENSE_NAME'), oldLicense.name, newLicense.name));
  }
  if (oldLicense.text !== newLicense.text) {
    diffDetails.value.push(new DiffDetails('' + t('CD_LICENSE_TEXT'), oldLicense.text, newLicense.text));
  }
  if (oldLicense.meta.licenseUrl !== newLicense.meta.licenseUrl) {
    diffDetails.value.push(
      new DiffDetails('' + t('COL_LICENSE_URL'), oldLicense.meta.licenseUrl, newLicense.meta.licenseUrl),
    );
  }
  if (oldLicense.meta.sourceUrl !== newLicense.meta.sourceUrl) {
    diffDetails.value.push(
      new DiffDetails('' + t('COL_SOURCE_URL'), oldLicense.meta.sourceUrl, newLicense.meta.sourceUrl),
    );
  }
  if (oldLicense.isDeprecatedLicenseId !== newLicense.isDeprecatedLicenseId) {
    diffDetails.value.push(
      new DiffDetails(
        '' + t('COL_LICENSE_SPDX_STATUS'),
        toSpdxDeprecatedOrEmpty(oldLicense.isDeprecatedLicenseId),
        toSpdxDeprecatedOrEmpty(newLicense.isDeprecatedLicenseId),
      ),
    );
  }
};

const next = () => {
  currentPosition.value++;
  prepareDataForCurrentPosition();
};

const prev = () => {
  currentPosition.value--;
  prepareDataForCurrentPosition();
};

const accept = async () => {
  Object.assign(oldLicense, {
    name: newLicense.name,
    text: newLicense.text,
    meta: {
      ...oldLicense.meta,
      licenseUrl: newLicense.meta.licenseUrl,
      sourceUrl: newLicense.meta.sourceUrl,
    },
    isDeprecatedLicenseId: newLicense.isDeprecatedLicenseId,
  });
  await adminService.updateLicense(oldLicense, oldLicense._key);
  await adminService.deleteSpdxLicense(newLicense._key);
  await reload();
};

const getTheme = () => {
  const theme = localStorage.getItem('disco-theme');
  if (theme === 'dark' || theme === 'light') {
    return theme;
  }
  return 'dark'; // Default theme
};

const reject = async () => {
  await adminService.deleteSpdxLicense(newLicense._key);
  await reload();
};

const close = () => {
  show.value = false;
};
defineExpose({open, close});
</script>
<template>
  <v-dialog v-model="show" content-class="large fixed-height" scrollable>
    <v-card class="pa-8 dDialog" flat>
      <v-card-title>
        <v-row>
          <v-col cols="10">
            <h1 class="d-headline mb-3">{{ t('DIFF_DETAILS') }}: {{ licenseId }}</h1>
            <p class="d-headline" v-if="oldLicense && oldLicense.meta">
              {{ t('COL_APPROVAL_STATUS') }}:
              {{ t('LT_APP_' + oldLicense.meta.approvalState.toUpperCase()) }}
            </p>
            <p class="d-headline" v-if="oldLicense && oldLicense.meta">
              {{ t('LICENSE_CHART_TITLE') }}:
              {{ t('LICENSE_CHART_STATUS_' + oldLicense.meta.isLicenseChart) }}
            </p>
            <p class="d-headline" v-if="oldLicense && oldLicense.meta">
              {{ t('COL_LICENSE_FAMILY') }}:
              {{ t('LIC_FAMILY_' + oldLicense.meta.family?.toUpperCase().replace(' ', '_')) }}
            </p>
          </v-col>
          <v-col cols="2" class="text-right px-0">
            <DCloseButton @click="close()" />
          </v-col>
        </v-row>
      </v-card-title>
      <v-card-text>
        <v-row class="mt-0 mb-3 mr-0 ml-0" v-for="(diff, i) in diffDetails" :key="i">
          <h3 class="d-subtitle d-secondary-text d-inline">{{ diff.field }}</h3>
          <code-diff
            :theme="getTheme()"
            :old-string="diff.oldValue"
            :new-string="diff.newValue"
            output-format="side-by-side"
            diff-style="word" />
        </v-row>
      </v-card-text>
      <v-card-actions>
        <v-row class="ma-3">
          <DCActionButton
            large
            class="mr-4"
            text="<<"
            @click="prev"
            :disabled="currentPosition == 1"
            v-if="diffs.length > 1" />
          <v-spacer></v-spacer>
          <span v-if="diffs.length > 1">{{ currentPosition }}/{{ diffs.length }}</span>
          <v-spacer></v-spacer>
          <DCActionButton
            large
            class="mr-4"
            text=">>"
            @click="next"
            :disabled="currentPosition == diffs.length"
            v-if="diffs.length > 1" />
          <v-spacer></v-spacer>
          <DCActionButton isDialogButton large class="mr-4" :text="t('BTN_ACCEPT')" @click="accept" color="success" />
          <DCActionButton isDialogButton large :text="t('BTN_REJECT')" @click="reject" color="error" />
        </v-row>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
