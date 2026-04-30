import {capitalizeFirstLetter} from '@disclosure-portal/utils/Tools';

const getTextStatusColor = (itemStatus: string) => {
  const color = capitalizeFirstLetter(itemStatus);
  return `rgb(var(--v-theme-pr${color}))`;
};

export const getStringArray = (value: unknown): string[] => (Array.isArray(value) ? (value as string[]) : []);

export const getBoolArray = (value: unknown): boolean[] => {
  if (!Array.isArray(value)) {
    return [];
  }
  return value
    .filter((entry) => typeof entry === 'boolean' || entry === 'true' || entry === 'false')
    .map((entry) => (typeof entry === 'boolean' ? entry : entry === 'true'));
};

export const usePolicyRulesUtils = () => {
  return {
    getTextStatusColor,
  };
};
