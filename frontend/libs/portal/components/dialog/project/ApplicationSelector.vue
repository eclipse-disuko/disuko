<template>
  <v-form ref="appform">
    <v-row dense>
      <v-col cols="12" xs="12" v-if="!showSelect">
        <v-text-field
          autocomplete="off"
          v-if="model"
          hide-details="auto"
          v-model="model.name"
          readonly
          variant="outlined"
          :class="{required: isRequired}"
          :label="t('TITLE_APPLICATION_SEARCH')"
          clearable
          :required="isRequired"
          @click:clear="clear"
        >
        </v-text-field>
      </v-col>
      <v-col cols="12" xs="12" v-if="showSelect">
        <v-autocomplete
          v-model="selectedApp!"
          v-model:search="searchFieldInput"
          :rules="activeRules"
          @update:search="debouncedSearch"
          autocomplete="off"
          :item-text="getInputItemText"
          @click:clear="clear"
          append-icon=""
          :placeholder="t('TITLE_APPLICATION_PLACEHOLDER')"
          :label="t('WIZARD_title_application_search')"
          clearable
          :no-filter="true"
          :items="items"
          :required="isRequired"
          :loading="isLoading"
          @change="onSelect"
          return-object
          variant="outlined"
          hide-details="auto"
          :class="{required: isRequired}"
        >
          <template v-slot:item="{item, props}">
            <v-list-item v-bind="props" title="">
              <div class="d-flex align-center w-100">
                <div class="custom-cell">
                  <v-icon color="primary" size="large">mdi-view-dashboard-outline</v-icon>
                </div>
                <div class="pl-2">
                  <v-list-item-title>{{ getInputItemText(item.raw) }}</v-list-item-title>
                  <v-list-item-subtitle>{{ 'Application (' + item.raw.id + ')' }}</v-list-item-subtitle>
                </div>
              </div>
            </v-list-item>
          </template>
          <template v-slot:selection="{item}">
            <span> {{ getInputItemText(item.raw) }}</span>
          </template>
          <template v-slot:no-data>
            <v-list-item>
              <v-list-item-title>
                {{ t(message) }}
                <strong>Name</strong>
              </v-list-item-title>
            </v-list-item>
          </template>
        </v-autocomplete>
      </v-col>
    </v-row>
  </v-form>
</template>

<script lang="ts">
import {Application} from '@disclosure-portal/model/Application';
import ApplicationService from '@disclosure-portal/services/application';
import {DiscoForm} from '@disclosure-portal/types/discobasics';
import {RuleFunction} from '@disclosure-portal/types/rules';
import _ from 'lodash';
import {defineComponent, onMounted, PropType, ref, watch} from 'vue';
import {useI18n} from 'vue-i18n';

export default defineComponent({
  name: 'ApplicationSelector',
  props: {
    modelValue: {
      type: Object as PropType<Application>,
      required: true,
    },
    isRequired: {
      type: Boolean,
      default: false,
    },
  },
  emits: ['update:modelValue'],
  setup(props, {emit}) {
    const {t} = useI18n();
    const showSelect = ref(true);
    const model = ref<Application | null>(null);
    const selectedApp = ref<Application | null>(null);
    const searchFieldInput = ref('');
    const items = ref<Application[]>([]);
    const isLoading = ref(false);
    const error = ref(false);
    const message = ref('APPLICATION_MESSAGE_DEFAULT');
    const showAlert = ref(false);
    const activeRules = ref<RuleFunction[]>([]);
    const appform = ref<DiscoForm | null>(null);

    watch(selectedApp, (newValue) => {
      emit('update:modelValue', newValue);
    });
    onMounted(() => {
      model.value = props.modelValue;
      if (!model.value) {
        model.value = new Application();
      }
      showSelect.value = !props.modelValue || !props.modelValue.id;
    });

    const search = async () => {
      if (!searchFieldInput.value || searchFieldInput.value.length < 4) {
        return;
      }
      isLoading.value = true;
      try {
        const applicationResult = await ApplicationService.searchApplicationByQuery(
          encodeURIComponent(searchFieldInput.value),
        );
        if (applicationResult) {
          items.value = applicationResult;
          error.value = false;
        } else {
          message.value = 'APPLICATION_MESSAGE_EMPTY';
        }
        isLoading.value = false;
      } catch (e) {
        isLoading.value = false;
        error.value = true;
      }
    };

    const debouncedSearch = _.debounce(search, 200);

    const onSelect = () => {
      if (!selectedApp.value) {
        return;
      }
      showAlert.value = false;
    };

    const clear = () => {
      if (model.value) {
        model.value.id = '';
        model.value.externalLink = '';
        model.value.name = '';
      }
      showSelect.value = true;
    };

    const validate = async (): Promise<boolean> => {
      if (props.isRequired) {
        activeRules.value.push(() => {
          return selectedApp.value && selectedApp.value.id ? true : t('APPLICATION_REQUIRED');
        });
      }
      if (appform.value) {
        const result = await appform.value.validate();
        return result.valid;
      }
      return false;
    };

    const getInputItemText = (item: Application): string => {
      return item.name || '';
    };

    return {
      t,
      showSelect,
      model,
      selectedApp,
      searchFieldInput,
      items,
      isLoading,
      error,
      message,
      showAlert,
      activeRules,
      appform,
      debouncedSearch,
      onSelect,
      clear,
      validate,
      getInputItemText,
    };
  },
});
</script>
