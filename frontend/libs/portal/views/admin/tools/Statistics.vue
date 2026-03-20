<script setup lang="ts">
import {MemStats} from '@disclosure-portal/model/Memstats';
import SystemStatsResponse from '@disclosure-portal/model/Statistic';
import AdminService from '@disclosure-portal/services/admin';
import {onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const systemProfile = ref<MemStats>({} as MemStats);
const systemProfiles = ref<MemStats[]>([]);
const stats = ref<SystemStatsResponse>(new SystemStatsResponse());
const tabs = ref<string | null>(null);

const getSystemProfile = async () => {
  const profile = await AdminService.getSystemProfile();
  systemProfile.value = profile;
  systemProfiles.value.push(profile);
};

const getStats = async () => {
  const serverResponses = await AdminService.getStats();
  stats.value = serverResponses.data;
};

const updateStats = async () => {
  const serverResponses = await AdminService.updateStats();
  stats.value = serverResponses.data;
};

onMounted(async () => {
  await getStats();
  await getSystemProfile();
});
</script>

<template>
  <TableLayout has-tab has-title>
    <template #buttons>
      <DCActionButton large @click="getStats" :hint="t('TT_GET_STATS')" :text="t('TT_GET_STATS')"></DCActionButton>
      <DCActionButton
        large
        @click="updateStats"
        :hint="t('TT_UPDATE_STATS')"
        :text="t('TT_UPDATE_STATS')"></DCActionButton>
      <DCActionButton
        large
        @click="getSystemProfile"
        :hint="t('TT_UPDATE_SYSTEM_PROFILE')"
        :text="t('TT_UPDATE_SYSTEM_PROFILE')"></DCActionButton>
    </template>
    <template #table>
      <v-tabs v-model="tabs" slider-color="mbti" active-class="active" show-arrows>
        <v-tab value="lastDays">
          {{ t('TITLE_STATISTIC_LAST_DAYS') }}
        </v-tab>
        <v-tab value="last12Months">
          {{ t('TITLE_STATISTIC_LAST_12_MONTHS') }}
        </v-tab>
      </v-tabs>
      <v-tabs-window v-model="tabs" style="min-height: 800px">
        <v-tabs-window-item value="lastDays">
          <SystemStatisticTable :stats="stats.dayStats" />
        </v-tabs-window-item>
        <v-tabs-window-item value="last12Months">
          <SystemStatisticTable :stats="stats.monthsStats" />
        </v-tabs-window-item>
      </v-tabs-window>
    </template>
  </TableLayout>
</template>

<style scoped></style>
