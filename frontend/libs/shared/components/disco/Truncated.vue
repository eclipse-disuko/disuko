<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<template>
  <span class="inline-block w-full">
    <span class="inline-block w-full truncate">
      <slot></slot>
    </span>
    <Tooltip>
      <slot></slot>
    </Tooltip>
  </span>
</template>
