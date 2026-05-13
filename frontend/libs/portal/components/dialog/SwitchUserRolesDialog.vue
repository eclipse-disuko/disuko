<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script lang="ts">
import Icons from '@disclosure-portal/constants/icons';
import {UserDto} from '@disclosure-portal/model/Users';
import DCActionButton from '@shared/components/disco/DCActionButton.vue';
import DCloseButton from '@shared/components/disco/DCloseButton.vue';
import DLabel from '@disclosure-portal/components/disco/DLabel.vue';
import {defineComponent, nextTick, reactive, ref} from 'vue';
import {useI18n} from 'vue-i18n';

class RoleWrapper {
  value: string;
  text: string;

  constructor(value: string, text: string) {
    this.value = value;
    this.text = text;
  }
}

export default defineComponent({
  name: 'SwitchUserRolesDialog',
  components: {
    DCloseButton,
    DLabel,
    DCActionButton,
  },
  emits: ['applyNewUserRoles'],
  setup(_, {emit}) {
    const {t} = useI18n();

    const show = ref(false);
    const title = ref('');
    const item = reactive<UserDto>(new UserDto());
    const initialRoles = ref<RoleWrapper[]>([]);
    const forceNonInternal = ref(false);
    const icons = Icons;
    const showDialog = (model: UserDto): void => {
      Object.assign(item, JSON.parse(JSON.stringify(model)));
      show.value = true;
      initialRoles.value = model.roles.map((role) => new RoleWrapper(role, t(role)));
      title.value = 'RESTRICT_USER_ROLES';
    };

    const doDialogAction = (): void => {
      nextTick(() => {
        closeDialog();
        if (forceNonInternal.value) {
          item.roles.length = 0;
        }
        emit('applyNewUserRoles', item, forceNonInternal.value);
      });
    };

    const closeDialog = (): void => {
      show.value = false;
    };

    return {
      title,
      item,
      initialRoles,
      forceNonInternal,
      show,
      showDialog,
      icons,
      doDialogAction,
      closeDialog,
      t,
    };
  },
});
</script>

<template>
  <slot name="default" :showDialog="showDialog"> </slot>
  <v-form>
    <v-dialog v-model="show" content-class="medium" width="650" scrollable>
      <v-card class="pa-8 dDialog">
        <v-card-title>
          <v-row>
            <v-col cols="10">
              <span class="text-h5">{{ t(title) }}</span>
            </v-col>
            <v-col cols="2" align="right">
              <DCloseButton @click="closeDialog" />
            </v-col>
          </v-row>
        </v-card-title>
        <v-card-text>
          <v-col cols="12" xs="12">
            <v-select
              :items="initialRoles"
              v-model="item.roles"
              multiple
              clearable
              :label="t('ROLES')"
              class="pb-2"
              item-title="text"
              variant="outlined"
              density="compact"
              hide-details="auto">
              <template v-slot:chip="{item, props}">
                <DLabel closable :parentProps="props" :labelName="item.title" iconName="mdi-security" />
              </template>
            </v-select>
          </v-col>
          <v-col cols="12" xs="12">
            <v-checkbox
              v-model="forceNonInternal"
              hide-details
              color="primary"
              :label="t('FORCE_NON_INTERNAL')"
              class="mt-0 mr-2 shrink" />
          </v-col>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <DCActionButton
            isDialogButton
            size="small"
            variant="text"
            @click="closeDialog"
            class="mr-5"
            :text="t('BTN_CANCEL')" />
          <DCActionButton
            isDialogButton
            size="small"
            variant="flat"
            @click="doDialogAction"
            :text="t('BTN_APPLY_ROLES')" />
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-form>
</template>
