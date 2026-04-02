<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script lang="ts" setup>
import Icons from '@disclosure-portal/constants/icons';
import Label from '@disclosure-portal/model/Label';
import {Rights} from '@disclosure-portal/model/Rights';
import SchemaModel from '@disclosure-portal/model/Schema';
import StatusDialogConfig from '@disclosure-portal/model/StatusDialogConfig';
import AdminService from '@disclosure-portal/services/admin';
import {useUserStore} from '@disclosure-portal/stores/user';
import {IMap} from '@disclosure-portal/utils/View';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {TOOLTIP_OPEN_DELAY_IN_MS} from '@shared/utils/constant';
import dayjs from 'dayjs';
import {onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import JsonViewer from 'vue-json-viewer';
import {useRoute, useRouter} from 'vue-router';

const route = useRoute();
const router = useRouter();
const {t} = useI18n();
const userStore = useUserStore();
const breadcrumbs = useBreadcrumbsStore();

const item = ref<SchemaModel>({} as SchemaModel);
const itemId = ref<string>(route.params.id as string);
const updated = ref<string>('');
const created = ref<string>('');
const schemaLabels = ref<Label[]>([]);
const labelsMap = ref<IMap<Label>>({});
const rights = ref<Rights>({} as Rights);
const icons = Icons;
const readOnly = ref<boolean>(true);

onMounted(async () => {
  rights.value = userStore.getRights;
  item.value = (await AdminService.getSchema(itemId.value)).data;
  updated.value = dayjs(item.value.updated.toString()).format(t('DATETIME_FORMAT_SHORT'));
  created.value = dayjs(item.value.created.toString()).format(t('DATETIME_FORMAT_SHORT'));
  if (route.path.includes('admin')) {
    readOnly.value = false;
  }
  initBreadcrumbs();
  await reloadLabels();
});

function close() {
  router.push('/dashboard/schemas');
}

async function activate() {
  const response = (await AdminService.activate(itemId.value)).data;
  item.value.active = true;

  if (response && response.success) {
    const d = new StatusDialogConfig();
    d.title = '' + t('SCHEMA_ACTIVATION');
    d.description = '' + t('schema_activation_successful') + item.value.name + ' / ' + item.value.version;
    Vue.prototype.$eventHub.$emit('onStatusMessage', d);
  }
}

async function download() {
  const link = document.createElement('a');
  link.target = '_blank';
  link.rel = 'noopener noreferrer';
  AdminService.downloadSchema(itemId.value)
    .then((res) => {
      link.download =
        'schema_' +
        item.value.name +
        '_' +
        item.value.version +
        (item.value.type === 0 ? '.json' : item.value.type === 1 ? '.xml' : '.dat');
      link.href = URL.createObjectURL(new Blob([res.data as unknown as BlobPart]));
      link.click();
    })
    .catch((e) => {
      console.error('cannot find schema ' + e);
    });
}

function initBreadcrumbs() {
  if (readOnly.value) {
    breadcrumbs.setCurrentBreadcrumbs([
      {
        title: t('BC_Dashboard'),
        disabled: false,
        href: '/dashboard/home',
      },
      {
        title: item.value.name,
        disabled: false,
        href: '/dashboard/admin/schemas/' + encodeURIComponent(item.value._key),
      },
    ]);
  } else {
    breadcrumbs.setCurrentBreadcrumbs([
      {
        title: t('BC_Dashboard'),
        href: '/dashboard/home',
      },
      {
        title: t('BC_ADMIN'),
        href: '/dashboard/admin',
      },
      {
        title: t('BC_SBOM_Schemes'),
        href: '/dashboard/admin/schemas/',
      },
      {
        title: item.value.name,
        href: '/dashboard/admin/schemas/' + encodeURIComponent(item.value._key),
      },
    ]);
  }
}

function getSchemaContent() {
  let content = '';
  try {
    content = JSON.parse(item.value.content);
  } catch (e) {
    content = item.value.content;
  }
  return content;
}

async function reloadLabels() {
  schemaLabels.value = (await AdminService.getSchemaLabels()).data;
  createLabelsMap();
}

function createLabelsMap() {
  labelsMap.value = {};
  for (const lbl of schemaLabels.value) {
    labelsMap.value[lbl._key] = lbl;
  }
}
</script>

<template>
  <v-container fluid>
    <v-row class="pt-1">
      <v-col>
        <h1 class="text-h5">
          {{ t('SCHEMA') }} <q>{{ item.name }}</q>
        </h1>
      </v-col>
    </v-row>
    <v-row>
      <v-col>
        <v-card>
          <v-row class="pa-4 pb-0">
            <v-col cols="5">
              <v-row>
                <v-col>
                  <span>{{ t('SCHEMA_VERSION') }}</span
                  ><br />
                  <span>{{ item.version }}</span>
                </v-col>
              </v-row>
              <v-row>
                <v-col>
                  <span>{{ t('DESCRIPTION') }}</span
                  ><br />
                  <span>{{ item.description }}</span>
                </v-col>
              </v-row>
            </v-col>
            <v-col cols="3">
              <v-row>
                <v-col>
                  <span>{{ t('CREATED') }}</span
                  ><br />
                  <span>{{ created }}</span>
                </v-col>
              </v-row>
              <v-row>
                <v-col>
                  <span>{{ t('UPDATED') }}</span
                  ><br />
                  <span>{{ updated }}</span>
                </v-col>
              </v-row>
            </v-col>
            <v-col cols="4">
              <v-row>
                <v-col>
                  <span>{{ t('LABELS') }}</span
                  ><br />
                  <v-tooltip :open-delay="TOOLTIP_OPEN_DELAY_IN_MS" location="bottom">
                    <template v-slot:activator="{props}">
                      <DLabel
                        :labelName="labelsMap[item.label] ? labelsMap[item.label].name : 'UNKNOWN_LABEL'"
                        :iconName="icons.SCHEMA"
                        v-bind="props" />
                    </template>
                    <span>{{ t('TT_schema_label_with_description') }}</span>
                    <span>{{ labelsMap[item.label] ? labelsMap[item.label].description : '' }}</span>
                  </v-tooltip>
                </v-col>
              </v-row>
              <v-row>
                <v-col>
                  <div>{{ t('SCHEMA_STATUS') }}</div>
                  <DCActionButton
                    large
                    icon="mdi-check"
                    :text="t('ACTIVATE_SCHEMA')"
                    :hint="t('TT_activate_schema')"
                    @click="activate"
                    class="mx-2"
                    v-if="!item.active && rights.allowSchema && rights.allowSchema.update && !readOnly" />
                  <div v-if="item.active">
                    <v-icon size="small" color="success">{{ icons.CIRCLE_FILLED }}</v-icon>
                    <span class="d-subtitle-2 pl-1">{{ 'active' }}</span>
                  </div>
                </v-col>
              </v-row>
            </v-col>
          </v-row>
          <v-row>
            <v-col col="12">
              <v-divider class="my-2" />
            </v-col>
          </v-row>
          <v-row class="px-4 pb-3">
            <v-col cols="12">
              <h1 class="text-h5">{{ t('CAPTION_SCHEMA_DETAILS') }}</h1>
              <span>{{ t('TXT_SBOM_JSON') }}</span>
              <DCActionButton
                large
                icon="mdi-download"
                :text="t('BTN_DOWNLOAD')"
                :hint="t('TT_download_schema')"
                @click="download"
                class="mx-2" />
            </v-col>
          </v-row>
          <v-row class="px-4 pb-3" v-if="item.content">
            <v-col>
              <h6 class="text-h6">{{ t('CD_RAW') }}:</h6>
              <json-viewer :value="getSchemaContent()" :expand-depth="3" aria-expanded="false" theme="jv-dark" sort />
            </v-col>
          </v-row>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>
