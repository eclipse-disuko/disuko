<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts">
import {OverallReview} from '@disclosure-portal/model/VersionDetails';
import {
  formatDateAndTime,
  getOverallReviewColor,
  getOverallReviewIcon,
  getOverallReviewTranslationKey,
  sbomOutdated,
} from '@disclosure-portal/utils/Table';
import {formatDateTime} from '@disclosure-portal/utils/View';
import {useI18n} from 'vue-i18n';
import {Anchor} from 'vuetify/framework';

interface Props {
  review: OverallReview;
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
</script>

<template>
  <v-tooltip :location="tooltipPosition" max-width="500" content-class="dpTooltip">
    <template v-slot:activator="{props}">
      <v-icon
        :color="getOverallReviewColor(review.state)"
        :class="isGroup ? 'ml-7' : ''"
        style="cursor: pointer"
        v-bind="props"
        :size="isTable ? 'medium' : 'default'">
        {{ getOverallReviewIcon(review.state) }}
      </v-icon>
    </template>
    <span class="text-subtitle-1">{{ t(getOverallReviewTranslationKey(review.state)) }}</span>
    <br />
    <div v-if="review.state">
      <span>{{ review.comment }}</span>
      <br v-if="review.comment" />
      <span class="d-text d-secondary-text">
        {{ t('OVERALL_REVIEW_FOR_SBOM') }}
        {{ review.sbomName }} - {{ formatDateAndTime(review.sbomUploaded) }}
      </span>
      <span v-if="sbomOutdated(review.sbomUploaded)" class="d-text d-secondary-text">
        <br />
        <v-icon class="pr-2" color="red" small>mdi-exclamation</v-icon>
        <span>{{ t('SBOM_IS_OUTDATED') }}</span>
      </span>
      <br v-if="review.sbomName && review.sbomUploaded" />
      <span class="d-text d-secondary-text">
        {{ t('BY_CREATOR_OVERALL_REVIEW') }}
        {{ review.creatorFullName }} ({{ review.creator }})
      </span>
      <br v-if="review.creatorFullName && review.creator" />
      <span class="d-text d-secondary-text">
        {{ t('LAST_UPDATE_OVERALL_REVIEW') }} {{ formatDateTime(review.updated) }}
      </span>
    </div>
  </v-tooltip>
</template>
