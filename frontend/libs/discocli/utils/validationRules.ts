import {FormRulesConfig} from './formRulesConfig';
import {isURL} from './validation';

type RuleFunction = (value: string) => boolean | string;

const ruleDefinitions = {
  required:
    (t: (key: string) => string): RuleFunction =>
    (value: string) =>
      !!value || t('VALIDATION_required'),

  validURL:
    (t: (key: string) => string): RuleFunction =>
    (value: string) =>
      (!!value && isURL(value)) || t('VALIDATION_url'),

  maxLength:
    (t: (key: string) => string, max: number = 2000): RuleFunction =>
    (value: string) =>
      (!!value && value.length <= max) ||
      `${t('VALIDATION_max_length')} [${value != null ? value.length.toString() : '0'}/${max}]`,

  minLength:
    (t: (key: string) => string, min: number = 1): RuleFunction =>
    (value: string) =>
      (!!value && value.length >= min) ||
      `${t('VALIDATION_min_length')} [${value != null ? value.length.toString() : '0'}/${min}]`,

  validUUID:
    (t: (key: string) => string): RuleFunction =>
    (value: string) => {
      const uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i;
      return (!!value && uuidRegex.test(value)) || t('VALIDATION_uuid');
    },
};

export function createFormRules(t: (key: string) => string, config?: FormRulesConfig) {
  const rules: Record<string, RuleFunction> = {};

  if (config?.required ?? true) {
    rules.required = ruleDefinitions.required(t);
  }

  if (config?.validURL ?? true) {
    rules.validURL = ruleDefinitions.validURL(t);
  }

  if (config?.maxLength ?? true) {
    const maxLength = config?.maxLength ?? 2000;
    rules.maxLength = ruleDefinitions.maxLength(t, maxLength);
  }

  return rules;
}

export function getRule(
  t: (key: string) => string,
  ruleName: keyof typeof ruleDefinitions,
  config?: {max?: number; min?: number},
) {
  if (ruleName === 'maxLength') {
    return ruleDefinitions[ruleName](t, config?.max ?? 2000);
  }
  if (ruleName === 'minLength') {
    return ruleDefinitions[ruleName](t, config?.min ?? 1);
  }
  return ruleDefinitions[ruleName](t);
}
