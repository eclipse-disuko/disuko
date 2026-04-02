<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {ComponentStats} from '@disclosure-portal/model/VersionDetails';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const props = withDefaults(
  defineProps<{
    stats: ComponentStats;
    showRedWarnDeniedDecisionsMessage?: boolean;
  }>(),
  {
    showRedWarnDeniedDecisionsMessage: false,
  },
);
</script>

<template>
  <div :class="showRedWarnDeniedDecisionsMessage ? 'border-l-6 border-l-red-50' : ''">
    <div v-if="showRedWarnDeniedDecisionsMessage" class="pa-3">
      {{ t('PROJECT_HAS_DENIED_DECISIONS') }}
    </div>
    <div class="mt-2 pa-2 border-data-table">
      <div class="mt-n5 pl-2 text-caption bg-used-components">
        {{ t('TAD_COMPONENTS_IN_USE') }}
      </div>
      <v-table>
        <template v-slot:default>
          <thead>
            <tr>
              <v-tooltip v-bind:location="'bottom'" content-class="dpTooltip">
                <template v-slot:activator="{props: tooltipProps}">
                  <th
                    :class="stats.Total > 0 ? 'text-center' : 'text-center text-grey-lighten-1'"
                    v-bind="tooltipProps">
                    <div class="d-flex align-center justify-center">
                      <v-icon size="small" class="pr-1">mdi-layers</v-icon>
                      <span class="text-caption">{{ t('TAD_COMPONENTS') }}</span>
                    </div>
                  </th>
                </template>
                <span>{{ t('TT_COMPONENTS_TOTAL') }}</span>
              </v-tooltip>
              <v-tooltip v-bind:location="'bottom'" content-class="dpTooltip">
                <template v-slot:activator="{props: tooltipProps}">
                  <th
                    :class="stats.Denied > 0 ? 'text-center' : 'text-center text-grey-lighten-1'"
                    v-bind="tooltipProps">
                    <div class="d-flex align-center justify-center">
                      <v-icon size="small" class="pr-1" :color="stats.Denied > 0 ? 'policyStatusDeniedColor' : ''">
                        mdi-minus-circle
                      </v-icon>
                      <span class="text-caption">{{ t('TAD_DENIED') }}</span>
                    </div>
                  </th>
                </template>
                <span>{{ t('TT_COMPONENTS_DENIED') }}</span>
              </v-tooltip>
              <v-tooltip v-bind:location="'bottom'" content-class="dpTooltip">
                <template v-slot:activator="{props: tooltipProps}">
                  <th
                    :class="stats.NoAssertion > 0 ? 'text-center' : 'text-center text-grey-lighten-1'"
                    v-bind="tooltipProps">
                    <div class="d-flex align-center justify-center">
                      <v-icon
                        size="small"
                        class="pr-1"
                        :color="stats.NoAssertion > 0 ? 'policyStatusUnassertedColor' : ''">
                        mdi-lightning-bolt-circle
                      </v-icon>
                      <span class="text-caption">{{ t('TAD_UNASSERTED') }}</span>
                    </div>
                  </th>
                </template>
                <span>{{ t('TT_COMPONENTS_NOASSERTION') }}</span>
              </v-tooltip>
              <v-tooltip v-bind:location="'bottom'" content-class="dpTooltip">
                <template v-slot:activator="{props: tooltipProps}">
                  <th
                    :class="stats.Warned > 0 ? 'text-center' : 'text-center text-grey-lighten-1'"
                    v-bind="tooltipProps">
                    <div class="d-flex align-center justify-center">
                      <v-icon size="small" class="pr-1" :color="stats.Warned > 0 ? 'policyStatusWarnedColor' : ''"
                        >mdi-alert
                      </v-icon>
                      <span class="text-caption"> {{ t('TAD_WARNED') }}</span>
                    </div>
                  </th>
                </template>
                <span>{{ t('TT_COMPONENTS_WARNED') }}</span>
              </v-tooltip>
              <v-tooltip v-bind:location="'bottom'" content-class="dpTooltip">
                <template v-slot:activator="{props: tooltipProps}">
                  <th
                    :class="stats.Questioned > 0 ? 'text-center' : 'text-center text-grey-lighten-1'"
                    v-bind="tooltipProps">
                    <div class="d-flex align-center justify-center">
                      <v-icon size="small" class="pr-1" :color="stats.Questioned > 0 ? 'green' : ''">mdi-help </v-icon>
                      <span class="text-caption"> {{ t('TAD_QUESTIONED') }}</span>
                    </div>
                  </th>
                </template>
                <span>{{ t('TT_COMPONENTS_QUESTIONED') }}</span>
              </v-tooltip>
              <v-tooltip v-bind:location="'bottom'" content-class="dpTooltip">
                <template v-slot:activator="{props: tooltipProps}">
                  <th
                    :class="stats.Allowed > 0 ? 'text-center' : 'text-center text-grey-lighten-1'"
                    v-bind="tooltipProps">
                    <div class="d-flex align-center justify-center">
                      <v-icon size="small" class="pr-1" :color="stats.Allowed > 0 ? 'green' : ''"
                        >mdi-check-circle
                      </v-icon>
                      <span class="text-caption"> {{ t('TAD_ALLOWED') }}</span>
                    </div>
                  </th>
                </template>
                <span>{{ t('TT_COMPONENTS_ALLOWED') }}</span>
              </v-tooltip>
            </tr>
          </thead>
          <tbody>
            <tr>
              <td :class="stats.Total === 0 ? 'text-center text-grey-lighten-1' : 'text-center'">
                {{ stats.Total }}
              </td>
              <td :class="stats.Denied === 0 ? 'text-center text-grey-lighten-1' : 'text-center'">
                {{ stats.Denied }}
              </td>
              <td :class="stats.NoAssertion === 0 ? 'text-center text-grey-lighten-1' : 'text-center'">
                {{ stats.NoAssertion }}
              </td>
              <td :class="stats.Warned === 0 ? 'text-center text-grey-lighten-1' : 'text-center'">
                {{ stats.Warned }}
              </td>
              <td :class="stats.Questioned === 0 ? 'text-center text-grey-lighten-1' : 'text-center'">
                {{ stats.Questioned }}
              </td>
              <td :class="stats.Allowed === 0 ? 'text-center text-grey-lighten-1' : 'text-center'">
                {{ stats.Allowed }}
              </td>
            </tr>
          </tbody>
        </template>
      </v-table>
    </div>
  </div>
</template>
