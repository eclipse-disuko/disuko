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
    <div class="pa-2 border-data-table mt-2">
      <div class="mt-n5 text-caption bg-used-components pl-2">
        {{ t('TAD_COMPONENTS_IN_USE') }}
      </div>
      <v-table>
        <template v-slot:default>
          <thead>
            <tr>
              <v-tooltip v-bind:location="'bottom'" content-class="dpTooltip">
                <template v-slot:activator="{props: tooltipProps}">
                  <th
                    :class="stats.total > 0 ? 'text-center' : 'text-grey-lighten-1 text-center'"
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
                    :class="stats.denied > 0 ? 'text-center' : 'text-grey-lighten-1 text-center'"
                    v-bind="tooltipProps">
                    <div class="d-flex align-center justify-center">
                      <v-icon size="small" class="pr-1" :color="stats.denied > 0 ? 'policyStatusDeniedColor' : ''">
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
                    :class="stats.noAssertion > 0 ? 'text-center' : 'text-grey-lighten-1 text-center'"
                    v-bind="tooltipProps">
                    <div class="d-flex align-center justify-center">
                      <v-icon
                        size="small"
                        class="pr-1"
                        :color="stats.noAssertion > 0 ? 'policyStatusUnassertedColor' : ''">
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
                    :class="stats.warned > 0 ? 'text-center' : 'text-grey-lighten-1 text-center'"
                    v-bind="tooltipProps">
                    <div class="d-flex align-center justify-center">
                      <v-icon size="small" class="pr-1" :color="stats.warned > 0 ? 'policyStatusWarnedColor' : ''"
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
                    :class="stats.questioned > 0 ? 'text-center' : 'text-grey-lighten-1 text-center'"
                    v-bind="tooltipProps">
                    <div class="d-flex align-center justify-center">
                      <v-icon size="small" class="pr-1" :color="stats.questioned > 0 ? 'green' : ''">mdi-help </v-icon>
                      <span class="text-caption"> {{ t('TAD_QUESTIONED') }}</span>
                    </div>
                  </th>
                </template>
                <span>{{ t('TT_COMPONENTS_QUESTIONED') }}</span>
              </v-tooltip>
              <v-tooltip v-bind:location="'bottom'" content-class="dpTooltip">
                <template v-slot:activator="{props: tooltipProps}">
                  <th
                    :class="stats.allowed > 0 ? 'text-center' : 'text-grey-lighten-1 text-center'"
                    v-bind="tooltipProps">
                    <div class="d-flex align-center justify-center">
                      <v-icon size="small" class="pr-1" :color="stats.allowed > 0 ? 'green' : ''"
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
              <td :class="stats.total === 0 ? 'text-grey-lighten-1 text-center' : 'text-center'">
                {{ stats.total }}
              </td>
              <td :class="stats.denied === 0 ? 'text-grey-lighten-1 text-center' : 'text-center'">
                {{ stats.denied }}
              </td>
              <td :class="stats.noAssertion === 0 ? 'text-grey-lighten-1 text-center' : 'text-center'">
                {{ stats.noAssertion }}
              </td>
              <td :class="stats.warned === 0 ? 'text-grey-lighten-1 text-center' : 'text-center'">
                {{ stats.warned }}
              </td>
              <td :class="stats.questioned === 0 ? 'text-grey-lighten-1 text-center' : 'text-center'">
                {{ stats.questioned }}
              </td>
              <td :class="stats.allowed === 0 ? 'text-grey-lighten-1 text-center' : 'text-center'">
                {{ stats.allowed }}
              </td>
            </tr>
          </tbody>
        </template>
      </v-table>
    </div>
  </div>
</template>
