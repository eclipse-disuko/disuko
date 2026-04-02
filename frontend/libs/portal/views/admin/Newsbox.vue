<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script lang="ts" setup>
import {ConfirmationType, IConfirmationDialogConfig} from '@disclosure-portal/components/dialog/ConfirmationDialog';
import ConfirmationDialog from '@disclosure-portal/components/dialog/ConfirmationDialog.vue';
import {NewsboxItem} from '@disclosure-portal/model/Newsbox';
import {useNewsboxStore} from '@disclosure-portal/stores/newsbox.store';
import {isURL} from '@disclosure-portal/utils/Validation';
import TextArea from '@shared/components/disco/TextArea.vue';
import TextField from '@shared/components/disco/TextField.vue';
import TableActionButtons, {TableActionButtonsProps} from '@shared/components/TableActionButtons.vue';
import {useBreadcrumbsStore} from '@shared/stores/breadcrumbs.store';
import {DataTableHeader} from '@shared/types/table';
import dayjs from 'dayjs';
import {computed, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {VForm} from 'vuetify/components';

const {t} = useI18n();
const breadcrumbs = useBreadcrumbsStore();
const newsboxStore = useNewsboxStore();

const form = ref({
  title: '',
  titleDE: '',
  description: '',
  descriptionDE: '',
  image: null as string | null,
  link: null as string | null,
  expiry: null as string | null,
} as NewsboxItem);
const submitting = ref(false);
const showDialog = ref(false);
const isEditMode = ref(false);
const search = ref('');
const formNewsboxDialog = ref<VForm | null>(null);
const editingItem = ref<NewsboxItem | null>(null);
const menuImageFilter = ref(false);
const menuLinkFilter = ref(false);
const menuStatusFilter = ref(false);
const selectedFilterImage = ref<string[]>([]);
const selectedFilterLink = ref<string[]>([]);
const selectedFilterStatus = ref<string[]>([]);
const sortBy = [{key: 'expiry', order: 'asc'}];
const confirmVisible = ref(false);
const confirmConfig = ref<IConfirmationDialogConfig>({} as IConfirmationDialogConfig);

const headers = computed<DataTableHeader[]>(() => [
  {title: t('COL_ACTIONS'), align: 'center', value: 'actions', width: 80, sortable: false, class: 'tableHeaderCell'},
  {title: t('TITLE'), align: 'start', value: 'title', width: 250, sortable: true, class: 'tableHeaderCell'},
  {title: t('DESCRIPTION'), align: 'start', value: 'description', width: 300, sortable: true, class: 'tableHeaderCell'},
  {title: t('TITLE_GERMAN'), align: 'start', value: 'titleDE', width: 200, sortable: true, class: 'tableHeaderCell'},
  {
    title: t('DESCRIPTION_GERMAN'),
    align: 'start',
    value: 'descriptionDE',
    width: 200,
    sortable: true,
    class: 'tableHeaderCell',
  },
  {
    title: t('HAS_IMAGE'),
    align: 'center',
    value: 'image',
    width: 120,
    sortable: false,
    filterable: true,
    class: 'tableHeaderCell',
  },
  {
    title: t('HAS_LINK'),
    align: 'center',
    value: 'link',
    width: 120,
    sortable: false,
    filterable: true,
    class: 'tableHeaderCell',
  },
  {
    title: t('STATUS'),
    align: 'center',
    value: 'expiry',
    width: 120,
    sortable: true,
    filterable: true,
    class: 'tableHeaderCell',
  },
]);

const initBreadcrumbs = () => {
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
      title: t('BC_NEWSBOX'),
      href: '/dashboard/admin/newsbox',
    },
  ]);
};

const resetForm = () => {
  form.value = {
    title: '',
    titleDE: '',
    description: '',
    descriptionDE: '',
    image: null,
    link: null,
    expiry: null,
  } as NewsboxItem;
  formNewsboxDialog.value?.resetValidation();
};

const openCreateDialog = () => {
  isEditMode.value = false;
  editingItem.value = null;
  resetForm();
  showDialog.value = true;
};

const closeDialog = () => {
  showDialog.value = false;
  isEditMode.value = false;
  editingItem.value = null;
  resetForm();
};

const showConfirmDelete = (item: NewsboxItem) => {
  confirmConfig.value = {
    type: ConfirmationType.DELETE,
    key: item._key,
    name: item.title,
    okButtonIsDisabled: false,
    okButton: 'BTN_DELETE',
    description: 'DLG_CONFIRMATION_DESCRIPTION',
  } as IConfirmationDialogConfig;
  confirmVisible.value = true;
};

const onConfirm = async () => {
  if (confirmConfig.value.type === ConfirmationType.DELETE && confirmConfig.value.key) {
    await newsboxStore.deleteItemsAdmin(confirmConfig.value.key);
    confirmVisible.value = false;
    reload();
  }
};

const openEditDialog = (item: NewsboxItem) => {
  isEditMode.value = true;
  editingItem.value = item;
  form.value = {
    title: '',
    titleDE: '',
    description: '',
    descriptionDE: '',
    image: null,
    link: null,
    expiry: null,
  } as NewsboxItem;
  Object.assign(form.value, item);

  if (item.expiry && item.expiry !== '' && item.expiry !== '0001-01-01T00:00:00Z') {
    form.value.expiry = dayjs(item.expiry).format('YYYY-MM-DD');
  } else {
    form.value.expiry = '';
  }
  showDialog.value = true;
};

const submit = async () => {
  submitting.value = true;
  if (isEditMode.value && editingItem.value?._key) {
    const updatedItem: NewsboxItem = {
      ...editingItem.value,
      title: form.value.title,
      titleDE: form.value.titleDE,
      description: form.value.description,
      descriptionDE: form.value.descriptionDE,
      image: form.value.image || null,
      link: form.value.link || null,
      expiry: form.value.expiry ? dayjs(form.value.expiry).toISOString() : '',
    };
    await newsboxStore.updateItemsAdmin(editingItem.value._key, updatedItem);
  } else {
    const newsboxItem = {
      title: form.value.title,
      titleDE: form.value.titleDE,
      description: form.value.description,
      descriptionDE: form.value.descriptionDE,
      image: form.value.image || null,
      link: form.value.link || null,
      expiry: form.value.expiry ? dayjs(form.value.expiry).toISOString() : null,
    };
    await newsboxStore.createItemsAdmin(newsboxItem);
  }
  closeDialog();
  reload();
  submitting.value = false;
};

const reload = async () => {
  await newsboxStore.fetchItemsAdmin();
};

const filteredItems = computed(() => {
  let items = newsboxStore.adminNewsItems?.items || [];

  if (selectedFilterImage.value.length > 0) {
    items = items.filter((item) => {
      const hasImage = !!item.image;
      return selectedFilterImage.value.includes(hasImage.toString());
    });
  }

  if (selectedFilterLink.value.length > 0) {
    items = items.filter((item) => {
      const hasLink = !!item.link;
      return selectedFilterLink.value.includes(hasLink.toString());
    });
  }

  if (selectedFilterStatus.value.length > 0) {
    items = items.filter((item) => {
      const expired = isExpired(item.expiry);
      return selectedFilterStatus.value.includes(expired.toString());
    });
  }

  return items;
});

const isExpired = (dateString: string | null | undefined) => {
  return !!(dateString && dateString !== '' && dateString !== '0001-01-01T00:00:00Z');
};

const handleImageError = (message: string) => {
  alert(message);
};

const dialogConfig = computed(() => ({
  title: isEditMode.value ? t('EDIT_NEWSBOX_ITEM') : t('TT_ADD_NEWSBOX_ITEM'),
  loading: submitting.value,
  primaryButton: isEditMode.value ? t('UPDATE') : t('CREATE'),
  secondaryButton: t('BTN_CANCEL'),
}));

const getActionButtons = (item: NewsboxItem): TableActionButtonsProps['buttons'] => {
  return [
    {
      icon: 'mdi-pencil',
      hint: t('TT_EDIT_NEWSBOX'),
      event: 'edit',
    },
    {
      icon: 'mdi-delete',
      hint: t('TT_DELETE_NEWSBOX'),
      event: 'delete',
    },
  ];
};

onMounted(() => {
  initBreadcrumbs();
  reload();
});
</script>

<template>
  <TableLayout data-testid="newsbox">
    <template #buttons>
      <h1 class="text-h5">{{ t('NEWSBOX') }}</h1>
      <DCActionButton
        large
        icon="mdi-plus"
        :hint="t('TT_ADD_NEWSBOX_ITEM')"
        :text="t('BTN_ADD')"
        class="mx-2"
        @click="openCreateDialog"></DCActionButton>
      <v-spacer></v-spacer>
      <v-text-field
        v-model="search"
        autocomplete="off"
        :max-width="500"
        append-inner-icon="mdi-magnify"
        variant="outlined"
        density="compact"
        :label="t('labelSearch')"
        single-line
        hide-details
        clearable></v-text-field>
    </template>
    <template #table>
      <div class="fill-height">
        <v-data-table
          v-model:search="search"
          density="comfortable"
          class="striped-table fill-height"
          fixed-header
          :headers="headers"
          :items="filteredItems"
          items-per-page="25"
          item-value="_key"
          :sort-by="sortBy"
          :cell-props="{class: 'py-3'}">
          <template v-slot:header.image="{column, getSortIcon, toggleSort}">
            <div class="v-data-table-header__content">
              <span>{{ column.title }}</span>
              <v-menu :close-on-content-click="false" v-model="menuImageFilter">
                <template v-slot:activator="{props}">
                  <DIconButton
                    :parentProps="props"
                    icon="mdi-filter-variant"
                    :hint="t('TT_SHOW_FILTER')"
                    :color="selectedFilterImage.length > 0 ? 'primary' : 'default'"
                    location="top" />
                </template>
                <div class="bg-background" style="width: 280px">
                  <v-row class="d-flex justify-end ma-1 mr-2">
                    <DIconButton icon="mdi-close" @clicked="menuImageFilter = false" color="default" />
                  </v-row>
                  <v-select
                    v-model="selectedFilterImage"
                    :items="[
                      {text: t('HAS_IMAGE'), value: 'true'},
                      {text: t('NO_IMAGE'), value: 'false'},
                    ]"
                    class="mx-2 pa-2 pb-4"
                    :label="t('Lbl_filter_on_type')"
                    clearable
                    multiple
                    item-title="text"
                    item-value="value"
                    variant="outlined"
                    density="compact"
                    menu
                    transition="scale-transition"
                    persistent-clear
                    :list-props="{class: 'striped-filter-dd py-0'}">
                    <template v-slot:item="{props}">
                      <v-list-item v-bind="props" class="py-0 px-2">
                        <template v-slot:prepend="{isSelected}">
                          <v-checkbox hide-details :model-value="isSelected" />
                        </template>
                        <template v-slot:title="{title}">
                          <span class="pFilterEntry">
                            {{ title }}
                          </span>
                        </template>
                      </v-list-item>
                    </template>
                    <template v-slot:selection="{item, index}">
                      <div v-if="index === 0" class="d-flex align-center">
                        <span class="pFilterEntry">{{ item.title }}</span>
                      </div>
                      <span v-if="index === 1" class="pAdditionalFilter">
                        +{{ selectedFilterImage.length - 1 }} others
                      </span>
                    </template>
                  </v-select>
                </div>
              </v-menu>
              <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
            </div>
          </template>

          <!-- Custom Header with Filter for HAS_LINK -->
          <template v-slot:header.link="{column, getSortIcon, toggleSort}">
            <div class="v-data-table-header__content">
              <span>{{ column.title }}</span>
              <v-menu :close-on-content-click="false" v-model="menuLinkFilter">
                <template v-slot:activator="{props}">
                  <DIconButton
                    :parentProps="props"
                    icon="mdi-filter-variant"
                    :hint="t('TT_SHOW_FILTER')"
                    :color="selectedFilterLink.length > 0 ? 'primary' : 'default'"
                    location="top" />
                </template>
                <div class="bg-background" style="width: 280px">
                  <v-row class="d-flex justify-end ma-1 mr-2">
                    <DIconButton icon="mdi-close" @clicked="menuLinkFilter = false" color="default" />
                  </v-row>
                  <v-select
                    v-model="selectedFilterLink"
                    :items="[
                      {text: t('HAS_LINK'), value: 'true'},
                      {text: t('NO_LINK'), value: 'false'},
                    ]"
                    class="mx-2 pa-2 pb-4"
                    :label="t('Lbl_filter_on_type')"
                    clearable
                    multiple
                    item-title="text"
                    item-value="value"
                    variant="outlined"
                    density="compact"
                    menu
                    transition="scale-transition"
                    persistent-clear
                    :list-props="{class: 'striped-filter-dd py-0'}">
                    <template v-slot:item="{props}">
                      <v-list-item v-bind="props" class="py-0 px-2">
                        <template v-slot:prepend="{isSelected}">
                          <v-checkbox hide-details :model-value="isSelected" />
                        </template>
                        <template v-slot:title="{title}">
                          <span class="pFilterEntry">
                            {{ title }}
                          </span>
                        </template>
                      </v-list-item>
                    </template>
                    <template v-slot:selection="{item, index}">
                      <div v-if="index === 0" class="d-flex align-center">
                        <span class="pFilterEntry">{{ item.title }}</span>
                      </div>
                      <span v-if="index === 1" class="pAdditionalFilter">
                        +{{ selectedFilterLink.length - 1 }} others
                      </span>
                    </template>
                  </v-select>
                </div>
              </v-menu>
              <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
            </div>
          </template>

          <!-- Custom Header with Filter for STATUS -->
          <template v-slot:header.expiry="{column, getSortIcon, toggleSort}">
            <div class="v-data-table-header__content">
              <span>{{ column.title }}</span>
              <v-menu :close-on-content-click="false" v-model="menuStatusFilter">
                <template v-slot:activator="{props}">
                  <DIconButton
                    :parentProps="props"
                    icon="mdi-filter-variant"
                    :hint="t('TT_SHOW_FILTER')"
                    :color="selectedFilterStatus.length > 0 ? 'primary' : 'default'"
                    location="top" />
                </template>
                <div class="bg-background" style="width: 280px">
                  <v-row class="d-flex justify-end ma-1 mr-2">
                    <DIconButton icon="mdi-close" @clicked="menuStatusFilter = false" color="default" />
                  </v-row>
                  <v-select
                    v-model="selectedFilterStatus"
                    :items="[
                      {text: t('ACTIVE'), value: 'false'},
                      {text: t('EXPIRED'), value: 'true'},
                    ]"
                    class="mx-2 pa-2 pb-4"
                    :label="t('Lbl_filter_on_type')"
                    clearable
                    multiple
                    item-title="text"
                    item-value="value"
                    variant="outlined"
                    density="compact"
                    menu
                    transition="scale-transition"
                    persistent-clear
                    :list-props="{class: 'striped-filter-dd py-0'}">
                    <template v-slot:item="{props}">
                      <v-list-item v-bind="props" class="py-0 px-2">
                        <template v-slot:prepend="{isSelected}">
                          <v-checkbox hide-details :model-value="isSelected" />
                        </template>
                        <template v-slot:title="{title}">
                          <span class="pFilterEntry">
                            {{ title }}
                          </span>
                        </template>
                      </v-list-item>
                    </template>
                    <template v-slot:selection="{item, index}">
                      <div v-if="index === 0" class="d-flex align-center">
                        <span class="pFilterEntry">{{ item.title }}</span>
                      </div>
                      <span v-if="index === 1" class="pAdditionalFilter">
                        +{{ selectedFilterStatus.length - 1 }} others
                      </span>
                    </template>
                  </v-select>
                </div>
              </v-menu>
              <v-icon class="v-data-table-header__sort-icon" :icon="getSortIcon(column)" @click="toggleSort(column)" />
            </div>
          </template>

          <template v-slot:[`item.description`]="{item}">
            <div class="text-truncate" style="max-width: 300px" :title="item.description">
              {{ item.description }}
            </div>
          </template>

          <template v-slot:[`item.titleDE`]="{item}">
            <div class="text-truncate" style="max-width: 200px" :title="item.titleDE">
              {{ item.titleDE || '-' }}
            </div>
          </template>

          <template v-slot:[`item.descriptionDE`]="{item}">
            <div class="text-truncate" style="max-width: 200px" :title="item.descriptionDE">
              {{ item.descriptionDE || '-' }}
            </div>
          </template>

          <template v-slot:[`item.image`]="{item}">
            <v-icon v-if="item.image" color="success">mdi-check</v-icon>
            <v-icon v-else color="grey">mdi-minus</v-icon>
          </template>

          <template v-slot:[`item.link`]="{item}">
            <v-icon v-if="item.link" color="success">mdi-check</v-icon>
            <v-icon v-else color="grey">mdi-minus</v-icon>
          </template>

          <template v-slot:[`item.expiry`]="{item}">
            <v-chip :color="isExpired(item.expiry) ? 'error' : 'success'" size="small">
              {{ isExpired(item.expiry) ? t('EXPIRED') : t('ACTIVE') }}
            </v-chip>
          </template>

          <template v-slot:[`item.actions`]="{item}">
            <TableActionButtons
              variant="compact"
              :buttons="getActionButtons(item)"
              @edit="openEditDialog(item)"
              @delete="showConfirmDelete(item)" />
          </template>
        </v-data-table>
      </div>
    </template>
  </TableLayout>

  <v-dialog v-model="showDialog" content-class="large" scrollable width="600">
    <ReactiveDialogLayout
      :config="dialogConfig"
      @primary-action="submit"
      @secondary-action="closeDialog"
      @close="closeDialog">
      <v-form ref="formNewsboxDialog" @submit.prevent="submit">
        <Stack>
          <TextField v-model="form.title" :label="t('TITLE')" required class="errorBorder" autofocus />
          <TextField v-model="form.titleDE" :label="t('TITLE_GERMAN')" />
          <TextArea v-model="form.description" :label="t('DESCRIPTION')" rows="3" required class="errorBorder" />
          <TextArea v-model="form.descriptionDE" :label="t('DESCRIPTION_GERMAN')" rows="3" />
          <div>
            <v-label class="text-subtitle-2 mb-2">{{ t('IMAGE') }}</v-label>
            <DImageUpload
              v-model="form.image"
              :upload-text="t('DRAG_DROP_IMAGE')"
              :button-text="t('SELECT_IMAGE')"
              @error="handleImageError" />
          </div>
          <TextField
            v-model="form.link"
            :label="t('LINK')"
            class="errorBorder"
            :rules="[(value: string) => !value || isURL(value) || t('VALIDATION_url')]"
            hint="Optional: External link URL" />
          <TextField v-if="isEditMode" v-model="form.expiry" :label="t('EXPIRY_DATE')" type="date">
            <template v-slot:append-inner>
              <v-btn
                v-if="form.expiry"
                icon="mdi-close"
                size="small"
                variant="text"
                @click="form.expiry = ''"
                :title="t('CLEAR_EXPIRY_DATE')" />
            </template>
          </TextField>
        </Stack>
      </v-form>
    </ReactiveDialogLayout>
  </v-dialog>
  <ConfirmationDialog v-model:showDialog="confirmVisible" :config="confirmConfig" @confirm="onConfirm" />
</template>
