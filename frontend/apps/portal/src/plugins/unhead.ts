// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

import {createHead} from '@unhead/vue/client';

const head = createHead({
  init: [
    {
      titleTemplate: import.meta.env.DEV ? '%s | Local Disuko' : '%s | Disclosure Portal',
    },
  ],
});

export default head;
