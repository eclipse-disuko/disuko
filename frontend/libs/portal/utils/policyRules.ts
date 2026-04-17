import {capitalizeFirstLetter} from '@disclosure-portal/utils/Tools';

const getTextStatusColor = (itemStatus: string) => {
  const color = capitalizeFirstLetter(itemStatus);
  return `rgb(var(--v-theme-pr${color}))`;
};

export const usePolicyRulesUtils = () => {
  return {
    getTextStatusColor,
  };
};
