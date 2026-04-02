// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {ReviewRemarkStatus} from '@disclosure-portal/model/Quality';

export interface BulkSetReviewRemarkStatusRequest {
  remarkKeys: string[];
  status: ReviewRemarkStatus;
}
