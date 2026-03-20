export interface Decision {
  // common part
  key: string;
  created: string;
  updated: string;
  sbomId: string;
  sbomName: string;
  sbomUploaded: string;
  componentSpdxId: string;
  componentName: string;
  componentVersion: string;
  licenseExpression: string;
  comment: string;
  creator: string;
  active: string;

  type: string;

  // license decision part only
  licenseDecisionId: string;
  licenseDecisionName: string;

  // policy decision part only
  licenseMatchedId: string;
  policyId: string;
  policyEvaluated: string;
  policyDecision: string;
}
