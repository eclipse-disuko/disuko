// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {OverallReviewState, VersionSlim} from '@disclosure-portal/model/VersionDetails';

/**
 * Composable for checking SBOM/SPDX approval and R&D confirmation status
 */
export function useApprovalCheck() {
  /**
   * Checks if a specific SBOM has been R&D confirmed (AUDITED status)
   * @param channel - The version/channel containing overall reviews
   * @param sbomKey - The SBOM key to check for R&D confirmation
   * @returns true if the SBOM has been R&D confirmed (AUDITED)
   */
  const isAudited = (channel: VersionSlim | null, sbomKey: string): boolean => {
    if (!channel || !sbomKey) {
      return false;
    }

    return (
      channel.overallReviews?.some(
        (review) => review.sbomId === sbomKey && review.state === OverallReviewState.AUDITED,
      ) ?? false
    );
  };

  return {
    isAudited,
  };
}
