export enum ChangeLogType {
  policyRuleChange = 'policy_rule_change',
}

export interface IChangeLogDto {
  _key: string;
  when: string;
  type: ChangeLogType;
  content: string;
}

export type PolicyRuleChangeLog = {
  licenseName: string;
  licenseId: string;
  policyStatus: 'Allow' | 'Deny' | 'Warn';
  change: 'Added' | 'Removed';
};

export type IPolicyRuleChangeLog = {
  _key: string;
  when: string;
} & {
  type: ChangeLogType.policyRuleChange;
  content: PolicyRuleChangeLog;
};

export class ChangeLogDto implements IChangeLogDto {
  public _key!: string;
  public when!: string;
  public type!: ChangeLogType;
  public content!: string;
}

export class ChangeLogResponse extends ChangeLogDto {
  constructor(dto: ChangeLogDto) {
    super();
    Object.assign(this, dto);
  }
}
