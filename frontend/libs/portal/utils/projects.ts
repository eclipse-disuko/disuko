import {capitalizeFirstLetter} from '@disclosure-portal/utils/Tools';

const getTextStatusColor = (itemStatus: string) => {
  const color = !itemStatus ? 'New' : capitalizeFirstLetter(itemStatus);
  return `rgb(var(--v-theme-project${color}))`;
};

export const useProjectUtils = () => {
  return {
    getTextStatusColor,
  };
};
