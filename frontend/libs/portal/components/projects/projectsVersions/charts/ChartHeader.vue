// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script lang="ts">
import {TOOLTIP_OPEN_DELAY_IN_MS} from '@shared/utils/constant';
import {defineComponent, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import {useRouter} from 'vue-router';

export default defineComponent({
  name: 'ChartHeader',
  props: {
    headerText: {
      type: String,
      required: true,
    },
    helpText: {
      type: String,
      required: true,
    },
    navigationPath: {
      type: String,
      required: false,
    },
  },
  setup(props) {
    const showMenu = ref(false);
    const router = useRouter();
    const {t} = useI18n();

    const navigateTo = () => {
      if (props.navigationPath) {
        router.push(props.navigationPath);
      }
    };

    return {
      t,
      showMenu,
      navigateTo,
      TOOLTIP_OPEN_DELAY_IN_MS,
    };
  },
});
</script>

<template>
  <v-row>
    <v-col cols="12" xs="12" sm="10">
      <div>
        <v-tooltip :open-delay="TOOLTIP_OPEN_DELAY_IN_MS" location="bottom">
          <template v-slot:activator="{props}">
            <v-icon
              class="mt-n1"
              v-if="navigationPath"
              color="primary"
              @click="navigateTo()"
              size="small"
              v-bind="props"
            >
              mdi mdi-chevron-right
            </v-icon>
          </template>
          <span>{{ t('CHARTS_TT_DETAILS') }}</span>
        </v-tooltip>
        <span style="color: #33a4fd; cursor: pointer" class="pt-2" v-if="navigationPath" @click="navigateTo()">
          {{ t(headerText) }}
        </span>
      </div>
    </v-col>
    <v-col class="justify-end text-right" xs="12" sm="2">
      <v-menu class="ma-3" v-model="showMenu" :close-on-content-click="false" offset-y>
        <template v-slot:activator="{props}">
          <v-tooltip :open-delay="TOOLTIP_OPEN_DELAY_IN_MS" location="bottom" v-bind="props">
            <template v-slot:activator="{props}">
              <v-icon color="primary" v-bind="props" @click="showMenu = !showMenu" size="small" icon="mdi-help" />
            </template>
            <span>{{ t('CHARTS_TT_HELP') }}</span>
          </v-tooltip>
        </template>
        <v-card>
          <v-card-text class="pa-2 helptext-container">
            {{ t(helpText) }}
          </v-card-text>
        </v-card>
      </v-menu>
    </v-col>
  </v-row>
</template>
