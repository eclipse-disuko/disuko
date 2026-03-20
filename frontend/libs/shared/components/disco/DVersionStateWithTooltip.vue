<script setup lang="ts">
import {VersionSlim} from '@disclosure-portal/model/VersionDetails';
import {
  formatDateAndTime,
  getIconColor,
  getOverallReviewTranslationKey,
  getVersionStateIcon,
  sbomOutdated,
} from '@disclosure-portal/utils/Table';
import {formatDateTime} from '@disclosure-portal/utils/View';
import dayjs from 'dayjs';
import {computed} from 'vue';
import {useI18n} from 'vue-i18n';
import {Anchor} from 'vuetify/framework';

interface Props {
  version: VersionSlim;
  tooltipPosition?: Anchor;
  isGroup?: boolean;
  isTable?: boolean;
}
const props = withDefaults(defineProps<Props>(), {
  tooltipPosition: 'start',
  isGroup: false,
  isTable: true,
});

const {t} = useI18n();

const recentOverallReview = computed(() => {
  if (!props.version.overallReviews || props.version.overallReviews.length === 0) {
    return null;
  }
  return props.version.overallReviews.reduce((recent, current) => {
    return dayjs(current.updated).isAfter(dayjs(recent.updated)) ? current : recent;
  });
});
</script>

<template>
  <v-tooltip :location="tooltipPosition" max-width="500" content-class="dpTooltip">
    <template v-slot:activator="{props}">
      <v-icon
        v-if="version"
        :color="getIconColor(version.status)"
        :class="isGroup ? 'ml-7' : ''"
        style="cursor: pointer"
        v-bind="props"
        :size="isTable ? 'medium' : 'default'">
        {{ getVersionStateIcon(version.status) }}
      </v-icon>
    </template>
    <span class="text-subtitle-1">{{ t(getOverallReviewTranslationKey(version.status)) }}</span>
    <br />
    <div v-if="recentOverallReview?.state">
      <span>{{ recentOverallReview?.comment }}</span>
      <br v-if="recentOverallReview?.comment" />
      <span class="d-text d-secondary-text">
        {{ t('OVERALL_REVIEW_FOR_SBOM') }}
        {{ recentOverallReview?.sbomName }} - {{ formatDateAndTime(recentOverallReview?.sbomUploaded) }}
      </span>
      <span v-if="sbomOutdated(recentOverallReview?.sbomUploaded)" class="d-text d-secondary-text">
        <br />
        <v-icon class="pr-2" color="red" small>mdi-exclamation</v-icon>
        <span>{{ t('SBOM_IS_OUTDATED') }}</span>
      </span>
      <br v-if="recentOverallReview.sbomName && recentOverallReview?.sbomUploaded" />
      <span class="d-text d-secondary-text">
        {{ t('BY_CREATOR_OVERALL_REVIEW') }}
        {{ recentOverallReview?.creatorFullName }} ({{ recentOverallReview?.creator }})
      </span>
      <br v-if="recentOverallReview?.creatorFullName && recentOverallReview?.creator" />
      <span class="d-text d-secondary-text">
        {{ t('LAST_UPDATE_OVERALL_REVIEW') }} {{ formatDateTime(recentOverallReview?.updated) }}
      </span>
    </div>
  </v-tooltip>
</template>
