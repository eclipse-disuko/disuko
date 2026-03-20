// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {WizardCard} from '@disclosure-portal/model/Wizard';

export const useWizard = () => {
  const toI18n = (title: string) => title.toUpperCase().replace(/ /g, '_').replace(/-/g, '_');

  const flipCard = (card: WizardCard) => {
    card.isFlipped = !card.isFlipped;
  };

  return {
    toI18n,
    flipCard,
  };
};
