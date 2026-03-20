<script setup lang="ts">
import {CustomId} from '@disclosure-portal/model/CustomId';
import adminService from '@disclosure-portal/services/admin';
import {useCustomIdStore} from '@disclosure-portal/stores/customid.store';
import {DiscoForm} from '@disclosure-portal/types/discobasics';
import useRules from '@disclosure-portal/utils/Rules';
import {isURL} from '@disclosure-portal/utils/Validation';
import useSnackbar from '@shared/composables/useSnackbar';
import {computed, nextTick, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const {info: snack} = useSnackbar();
const {minMax, longText} = useRules();
const customIdsStore = useCustomIdStore();

const emit = defineEmits(['reload']);

const isVisible = ref(false);
const isEdit = ref(false);
const tidManuallyChanged = ref(false);
const customId = ref(new CustomId());
const dialog = ref<DiscoForm | null>(null);
const title = ref('');
const confirmText = ref('');
const saving = ref(false);

const rules = {
  tid: [...minMax(t('COL_TID'), 4, 36, false), ...[(v: string) => !/[^A-Za-z0-9-]/.test(v) || t('MALFORMED_TID')]],
  name: minMax(t('NPV_DIALOG_TF_NAME'), 3, 80, false),
  description: longText(t('NP_DIALOG_TF_DESCRIPTION')),
  validURL: [(value: string) => !value || isURL(value) || t('VALIDATION_url')],
};

const dialogConfig = computed(() => ({
  title: t(title.value),
  loading: saving.value,
  primaryButton: t(confirmText.value),
  secondaryButton: t('BTN_CANCEL'),
}));

const open = (existing?: CustomId) => {
  tidManuallyChanged.value = false;
  if (existing) {
    isEdit.value = true;
    customId.value = {
      ...new CustomId(),
      _key: existing._key,
      name: existing.name,
      nameDE: existing.nameDE,
      description: existing.description,
      descriptionDE: existing.descriptionDE,
      linkTemplate: existing.linkTemplate,
    };
  } else {
    isEdit.value = false;
    customId.value = {} as CustomId;
  }
  dialog.value?.reset();
  title.value = existing ? 'CUSTOMID_DIALOG_EDIT_TITLE' : 'CUSTOMID_DIALOG_ADD_TITLE';
  confirmText.value = existing ? 'NP_DIALOG_BTN_EDIT' : 'NP_DIALOG_BTN_CREATE';
  isVisible.value = true;
};

const nameFocused = (focused: boolean) => {
  if (focused || isEdit.value || tidManuallyChanged.value) {
    return;
  }
  customId.value._key = customId.value.name.toLowerCase().replace(/[^a-z0-9-]/g, '');
};

const doDialogAction = async () => {
  await nextTick();
  const info = await dialog.value?.validate();
  if (!info?.valid) {
    return;
  }
  saving.value = true;
  try {
    if (isEdit.value) {
      await adminService.editCustomId(customId.value);
      snack(t('DIALOG_customid_edit_success'));
    } else {
      await adminService.createCustomId(customId.value);
      snack(t('DIALOG_customid_create_success'));
    }
    await customIdsStore.updateCustomIds();
    emit('reload');
    isVisible.value = false;
  } finally {
    saving.value = false;
  }
};

const close = () => {
  isVisible.value = false;
};

defineExpose({
  open,
});
</script>

<template>
  <v-dialog v-model="isVisible" content-class="large" scrollable width="500">
    <ReactiveDialogLayout
      :config="dialogConfig"
      @primary-action="doDialogAction"
      @secondary-action="close"
      @close="close">
      <v-form ref="dialog" @submit.prevent="doDialogAction">
        <Stack>
          <v-text-field
            class="required errorBorder"
            v-model="customId.name"
            :rules="rules.name"
            :label="t('AL_DIALOG_TF_NAME_EN')"
            autofocus
            hide-details="auto"
            variant="outlined"
            @update:focused="nameFocused" />
          <v-text-field
            class="required errorBorder"
            v-model="customId.nameDE"
            :rules="rules.name"
            :label="t('AL_DIALOG_TF_NAME_DE')"
            hide-details="auto"
            variant="outlined"
            @update:focused="nameFocused" />
          <v-text-field
            autocomplete="off"
            class="required errorBorder"
            v-model="customId._key"
            :rules="rules.tid"
            :label="t('COL_TID')"
            hide-details="auto"
            variant="outlined"
            :readonly="isEdit"
            @update:modelValue="tidManuallyChanged = true" />
          <v-text-field
            autocomplete="off"
            class="errorBorder"
            v-model="customId.linkTemplate"
            :label="t('LINK_TEMPLATE')"
            :rules="rules.validURL"
            hide-details="auto"
            variant="outlined"
            :hint="t('HINT_LINK_TEMPLATE')" />
          <v-textarea
            no-resize
            v-model="customId.description"
            :rules="rules.description"
            :label="t('AL_DIALOG_TF_DESCRIPTION_EN')"
            :counter="1000"
            variant="outlined" />
          <v-textarea
            no-resize
            v-model="customId.descriptionDE"
            :rules="rules.description"
            :label="t('AL_DIALOG_TF_DESCRIPTION_DE')"
            :counter="1000"
            variant="outlined" />
        </Stack>
      </v-form>
    </ReactiveDialogLayout>
  </v-dialog>
</template>
