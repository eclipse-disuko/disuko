// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {
  MatrixResponseDto,
  PolicyRuleClassificationDto,
  PolicyRuleClassificationRequestDto,
} from '@disclosure-portal/model/PolicyRuleClassification';
import {useApi} from '@shared/api/useApi';

const {api} = useApi();
const basePath = '/api/v1/admin/policyruleclassifications';

class PolicyRuleClassificationService {
  public getAll() {
    return api.get<PolicyRuleClassificationDto[]>(`${basePath}/`);
  }

  public getMatrix() {
    return api.get<MatrixResponseDto>(`${basePath}/matrix`);
  }

  public create(dto: PolicyRuleClassificationRequestDto) {
    return api.post<PolicyRuleClassificationDto>(`${basePath}/`, dto);
  }

  public update(key: string, dto: PolicyRuleClassificationRequestDto) {
    return api.put<PolicyRuleClassificationDto>(`${basePath}/${key}`, dto);
  }

  public delete(key: string) {
    return api.delete(`${basePath}/${key}`);
  }
}

const policyRuleClassificationService = new PolicyRuleClassificationService();
export default policyRuleClassificationService;
