<script lang="ts" setup>
import {Token} from '@disclosure-portal/model/Project';
import projectService from '@disclosure-portal/services/projects';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import useRules from '@disclosure-portal/utils/Rules';
import ReactiveDialogLayout, {type ReactiveDialogLayoutConfig} from '@shared/layouts/ReactiveDialogLayout.vue';
import dayjs from 'dayjs';
import {computed, nextTick, onMounted, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';
import {VForm, VSelect} from 'vuetify/components';

enum TokenDuration {
  month = 'month',
  threeMonths = 'three_months',
  sixMonths = 'six_months',
  year = 'year',
  max = 'max',
  custom = 'custom',
}

const allDurations: TokenDuration[] = [
  TokenDuration.month,
  TokenDuration.threeMonths,
  TokenDuration.sixMonths,
  TokenDuration.year,
  TokenDuration.max,
  TokenDuration.custom,
];

const nowPlusDuration = (duration: TokenDuration) => {
  switch (duration) {
    case TokenDuration.month:
      return dayjs().add(1, 'month');
    case TokenDuration.threeMonths:
      return dayjs().add(3, 'months');
    case TokenDuration.sixMonths:
      return dayjs().add(6, 'months');
    case TokenDuration.year:
      return dayjs().add(1, 'year');
    case TokenDuration.max:
      return dayjs().add(2, 'years');
    case TokenDuration.custom:
      return dayjs('');
  }
};

const emits = defineEmits(['onCreated']);
const {t} = useI18n();
const isDialogVisible = ref(false);
const isDatePickerVisible = ref(false);
const tokenForm = ref<VForm>();
onMounted(() => tokenForm.value?.reset());
const rules = useRules();
const loading = ref(false);
const projectUUID = computed(() => useProjectStore().currentProject!._key);

const nowDate = ref(dayjs().add(1, 'days').toISOString().slice(0, 10));
const maxDate = ref(dayjs().add(24, 'months').toISOString().slice(0, 10));

const name = ref('');
const nameRules = rules.required(t('TOKEN_NAME'));
const description = ref('');

const allTokenDurationOptions = allDurations;
const selectedDurationRef = ref<VSelect>();
const selectedDuration = ref<TokenDuration>(TokenDuration.year);
onMounted(() => (selectedDuration.value = TokenDuration.year));
const selectedDurationRules = rules.required(t('TKN_EXPIRATION'));
selectedDurationRules.push(() => selectedMoment.value !== undefined || t('TKN_date_required_hint'));

const selectedMoment = ref(dayjs().add(1, 'year'));
watch(selectedDuration, () => (selectedMoment.value = nowPlusDuration(selectedDuration.value)));
watch(selectedMoment, async () => await selectedDurationRef.value?.resetValidation());

const selectedDate = computed({
  get: () => selectedMoment.value?.toDate() ?? new Date(),
  set: (value: Date) => (selectedMoment.value = dayjs(value)),
});

const selectedDurationAsString = computed(() =>
  selectedMoment.value
    ? t('TKN_duration_' + selectedDuration.value) + ' - ' + selectedMoment.value.format(t('DATETIME_FORMAT_SHORT'))
    : t('TKN_duration_' + selectedDuration.value),
);

function generateToken() {
  nextTick(async () => {
    const validationResult = (await tokenForm.value?.validate())?.valid;
    if (validationResult) {
      loading.value = true;
      const token: Token = {
        _key: '',
        company: name.value,
        description: description.value,
        status: '',
        created: '',
        expiry: selectedMoment.value?.toISOString() ?? '',
        tokenSecret: '',
      };
      try {
        const response = await projectService.addProjectToken(projectUUID.value, token);
        closeDialog();
        emits('onCreated', response.data);
      } finally {
        loading.value = false;
      }
    }
  });
}

function showDialog() {
  tokenForm.value?.reset();
  isDialogVisible.value = true;
}

function closeDialog() {
  isDialogVisible.value = false;
}

const dialogConfig = computed(
  (): ReactiveDialogLayoutConfig => ({
    title: t('UM_DIALOG_TITLE_NEW_TOKEN'),
    loading: loading.value,
    secondaryButton: t('BTN_CANCEL'),
    primaryButton: t('NP_DIALOG_BTN_CREATE'),
  }),
);
</script>

<template>
  <slot :showDialog="showDialog" />

  <v-form ref="tokenForm">
    <v-dialog v-model="isDialogVisible" content-class="small" persistent width="600">
      <ReactiveDialogLayout
        :config="dialogConfig"
        @close="closeDialog"
        @secondary-action="closeDialog"
        @primary-action="generateToken">
        <Stack>
          <v-text-field
            autocomplete="off"
            v-model="name"
            :rules="nameRules"
            :label="t('TOKEN_NAME')"
            autofocus
            required
            variant="outlined" />
          <v-text-field autocomplete="off" v-model="description" :label="t('DESCRIPTION')" variant="outlined" />

          <v-menu v-model="isDatePickerVisible" :close-on-content-click="false" :target="selectedDurationRef?.$el">
            <template v-slot:activator>
              <v-select
                ref="selectedDurationRef"
                v-model="selectedDuration"
                :items="allTokenDurationOptions"
                :label="t('TKN_EXPIRATION')"
                :rules="selectedDurationRules"
                variant="outlined"
                return-object
                hide-details
                required>
                <template v-slot:item="{item, props}">
                  <v-list-item v-bind="props" :title="t('TKN_duration_' + item.value)"> </v-list-item>
                </template>
                <template v-if="selectedDuration === 'custom'" v-slot:append>
                  <v-btn variant="flat" large icon="mdi-calendar" @click="isDatePickerVisible = true">
                    <v-card>
                      <v-icon color="primary" />
                    </v-card>
                  </v-btn>
                </template>

                <template v-slot:selection>
                  {{ selectedDurationAsString }}
                </template>
              </v-select>
            </template>

            <v-date-picker
              show-current
              v-model="selectedDate"
              hide-header
              color="primary"
              :first-day-of-week="1"
              :min="nowDate"
              :max="maxDate"
              @update:modelValue="isDatePickerVisible = false">
            </v-date-picker>
          </v-menu>
        </Stack>
      </ReactiveDialogLayout>
    </v-dialog>
  </v-form>
</template>
