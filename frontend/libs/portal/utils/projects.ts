const getTextStatusColor = (itemStatus: string) => {
  const color = !itemStatus ? 'New' : itemStatus.substring(0, 1).toUpperCase() + itemStatus.substring(1);
  return `rgb(var(--v-theme-project${color}))`;
};

export const useProjectUtils = () => {
  return {
    getTextStatusColor,
  };
};
