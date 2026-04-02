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
        :color="getOverallReviewColor(review.State)"
        :class="isGroup ? 'ml-7' : ''"
        style="cursor: pointer"
        v-bind="props"
        :size="isTable ? 'medium' : 'default'">
        {{ getOverallReviewIcon(review.State) }}
      </v-icon>
    </template>
    <span class="text-subtitle-1">{{ t(getOverallReviewTranslationKey(review.State)) }}</span>
    <br />
    <div v-if="review.State">
      <span>{{ review.Comment }}</span>
      <br v-if="review.Comment" />
      <span class="d-text d-secondary-text">
        {{ t('OVERALL_REVIEW_FOR_SBOM') }}
        {{ review.SBOMName }} - {{ formatDateAndTime(review.SBOMUploaded) }}
      </span>
      <span v-if="sbomOutdated(review.SBOMUploaded)" class="d-text d-secondary-text">
        <br />
        <v-icon class="pr-2" color="red" small>mdi-exclamation</v-icon>
        <span>{{ t('SBOM_IS_OUTDATED') }}</span>
      </span>
      <br v-if="review.SBOMName && review.SBOMUploaded" />
      <span class="d-text d-secondary-text">
        {{ t('BY_CREATOR_OVERALL_REVIEW') }}
        {{ review.CreatorFullName }} ({{ review.Creator }})
      </span>
      <br v-if="review.CreatorFullName && review.Creator" />
      <span class="d-text d-secondary-text">
        {{ t('LAST_UPDATE_OVERALL_REVIEW') }} {{ formatDateTime(review.Updated) }}
      </span>
    </div>
  </v-tooltip>
</template>
