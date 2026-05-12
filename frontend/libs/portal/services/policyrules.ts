// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {useApi} from '@shared/api/useApi';
import PolicyRule from '@disclosure-portal/model/PolicyRule';

const {api} = useApi();

const modelName = 'policyrules';

class PolicyRuleService {
  public getAllPolicyRules() {
    return api.get<PolicyRule[]>(`/api/v1/${modelName}/`);
  }

  public getPolicyRule(id: string) {
    return api.get<PolicyRule>(`/api/v1/${modelName}/${id}`);
  }

  public async downloadSingleLPcsv(id: string) {
    return api.get(`/api/v1/${modelName}/${id}/csv`);
  }
}

const policyRuleService = new PolicyRuleService();
export default policyRuleService;
