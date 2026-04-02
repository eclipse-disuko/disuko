<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script setup lang="ts"></script>

<template>
  <div class="snowflake-container">
    <div class="snowflake" v-for="i in 15" :key="i" :style="`--delay: ${i * 0.5}s; --x: ${i * 6.5}%;`"></div>
  </div>
</template>

<style scoped lang="scss">
.snowflake-container {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  pointer-events: none;
  z-index: 10;
}

.snowflake {
  position: absolute;
  top: -20px;
  color: rgba(255, 255, 255, 0.9);
  font-size: 21px;
  pointer-events: none;
  animation: fall linear infinite;
  animation-delay: var(--delay, 0s);
  animation-duration: 8s;
  left: var(--x, 0%);
  will-change: transform;
}

.snowflake::before {
  content: '❄';
  text-shadow: 0 0 2px rgba(255, 255, 255, 1);
}

@keyframes fall {
  0% {
    transform: translateY(-20px) rotate(0deg);
    opacity: 1;
  }
  100% {
    transform: translateY(60px) rotate(360deg);
    opacity: 0.3;
  }
}

.snowflake:nth-child(2n) {
  font-size: 18px;
  animation-duration: 7s;
}

.snowflake:nth-child(3n) {
  font-size: 24px;
  animation-duration: 9s;
}

.snowflake:nth-child(4n) {
  font-size: 15px;
  animation-duration: 6s;
}

.snowflake:nth-child(5n) {
  font-size: 27px;
  animation-duration: 10s;
}

.snowflake:nth-child(2n) {
  animation-name: fall-drift-left;
}

.snowflake:nth-child(3n) {
  animation-name: fall-drift-right;
}

@keyframes fall-drift-left {
  0% {
    transform: translateY(-20px) translateX(0) rotate(0deg);
    opacity: 1;
  }
  100% {
    transform: translateY(60px) translateX(var(--drift, -25px)) rotate(360deg);
    opacity: 0.3;
  }
}

@keyframes fall-drift-right {
  0% {
    transform: translateY(-20px) translateX(0) rotate(0deg);
    opacity: 1;
  }
  100% {
    transform: translateY(60px) translateX(var(--drift, 25px)) rotate(360deg);
    opacity: 0.3;
  }
}
</style>
