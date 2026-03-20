// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

// src/composables/useDimensions.ts
import {onMounted, onUnmounted, reactive, UnwrapRef} from 'vue';

interface Dimensions {
  windowHeight: number;
  mainHeight: number;
  toolbarHeight: number;
  footerHeight: number;
}

const dimensions = reactive<Dimensions>({
  windowHeight: 0,
  mainHeight: 0,
  toolbarHeight: 0,
  footerHeight: 0,
});

const calculateTotalHeight = (parentElement: HTMLElement, padding: number) => {
  if (!parentElement) {
    return 0;
  }

  // Find all v-col elements within the parent element
  const vColElements = parentElement.querySelectorAll('.v-col');

  let totalHeight = 10;
  let found = false;

  vColElements.forEach((vColElement) => {
    const containsGrid = vColElement.querySelector('.v-data-table');
    if (!containsGrid) {
      if ((vColElement as HTMLElement).classList.contains('measure')) {
        let height = (vColElement as HTMLElement).offsetHeight + padding;
        totalHeight += height;
        found = true;
      }
    }
  });
  if (found) {
    totalHeight += padding;
  }
  return totalHeight;
};
const findMeasureHeightElement = (element: HTMLElement): HTMLElement | null => {
  let currentElement: HTMLElement | null = element;
  if (currentElement) {
    for (let i = 0; i < 3; i++) {
      if (currentElement && currentElement.classList && currentElement.classList.contains('measure-height')) {
        return currentElement;
      }
      currentElement = currentElement.parentElement;
      if (!currentElement) {
        return null;
      }
    }
  }

  return null;
};
const calculateHeight = (
  dataTableAsElement: UnwrapRef<HTMLElement | null>,
  withFooter: boolean = false,
  withInTab: boolean = false,
  heightAppender: string[] = [],
) => {
  updateDimensions();
  let tableHeight = dimensions.windowHeight;
  const margin = 20;
  tableHeight -= dimensions.toolbarHeight + margin;
  tableHeight -= dimensions.footerHeight + margin;

  if (dataTableAsElement) {
    let parentMeasureElement = findMeasureHeightElement(dataTableAsElement);

    if (parentMeasureElement) {
      const sumTotalHeight = calculateTotalHeight(parentMeasureElement, margin);
      if (sumTotalHeight) {
        tableHeight -= sumTotalHeight;
      }
    }
  }
  heightAppender.forEach((divName) => {
    const element = document.getElementById(divName); // Falls die Namen IDs sind
    // Oder: document.querySelector(`.${divName}`), falls es Klassen sind

    if (element) {
      tableHeight -= element.offsetHeight; // offsetHeight enthält die tatsächliche Höhe inkl. Padding
    } else {
      console.warn(`Element mit dem Namen ${divName} wurde nicht gefunden.`);
    }
  });
  if (withFooter) {
    tableHeight -= 56;
  }
  if (withInTab) {
    tableHeight -= 226;
  }
  return tableHeight - margin - 4;
};
const updateDimensions = () => {
  // const main = document.getElementById("disco-main");
  const toolbar = document.getElementById('disco-toolbar');
  const footer = document.getElementById('disco-footer');
  dimensions.windowHeight = window.innerHeight;
  dimensions.toolbarHeight = toolbar ? toolbar.offsetHeight : 0;
  dimensions.footerHeight = footer ? footer.offsetHeight : 0;
  dimensions.mainHeight = window.innerHeight;
};

const useDimensions = () => {
  onMounted(() => {});

  onUnmounted(() => {});

  return {
    calculateHeight,
    dimensions,
  };
};

export default useDimensions;
