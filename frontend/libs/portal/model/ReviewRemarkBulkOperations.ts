import {ReviewRemarkStatus} from '@disclosure-portal/model/Quality';

export interface BulkSetReviewRemarkStatusRequest {
  remarkKeys: string[];
  status: ReviewRemarkStatus;
}
