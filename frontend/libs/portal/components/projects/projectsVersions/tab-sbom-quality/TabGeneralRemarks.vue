<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
</script>

<template>
  <div class="p-4">
    <v-row>
      <v-col cols="12" xs="12" sm="12" md="10" lg="8">
        <span class="text-caption" v-html="t('QT_GENERAL_REMARKS')"></span>
      </v-col>
    </v-row>
  </div>
</template>
