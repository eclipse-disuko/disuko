<script setup lang="ts">
export interface WizardCardProps {
  key: string;
  image: string;
  title: string;
  subtitle: string;
  helptext: string;
  isFlipped: boolean;
  isActive?: boolean;
}
defineProps<{card: WizardCardProps}>();

const emit = defineEmits(['update']);

const flipCard = async (card: WizardCardProps) => {
  card.isFlipped = !card.isFlipped;
};
</script>

<template>
  <div class="w-full h-[360px] flex align-center justify-center">
    <div
      class="relative w-full h-full transition-transform duration-[0.6s] transform-3d"
      :class="{'transform-[rotateY(180deg)]': !card.isFlipped}">
      <v-card class="absolute w-full h-full" min-height="360" @click="emit('update', card.key)">
        <v-card-text class="overflow-auto max-h-[360px]">
          <Stack>
            <div class="self-end">
              <v-btn size="x-small" flat @click.stop="flipCard(card)" icon variant="text">
                <v-icon color="primary">mdi-close</v-icon>
              </v-btn>
            </div>
            <div class="flex justify-center font-[0.8rem]">
              {{ card.helptext }}
            </div>
          </Stack>
        </v-card-text>
      </v-card>
      <v-card
        class="absolute w-full h-full backface-hidden transform-[rotateY(180deg)]"
        :class="{active: card.isActive}"
        min-height="360"
        @click="emit('update', card.key)">
        <v-card-text>
          <Stack>
            <div class="self-end">
              <v-btn size="x-small" flat @click.stop="flipCard(card)" icon variant="text">
                <v-icon color="primary">mdi mdi-help</v-icon>
              </v-btn>
            </div>
            <div class="flex justify-center align-end">
              <v-img :src="card.image" height="80" max-width="130"></v-img>
            </div>
            <div class="flex justify-center">
              <h2 class="text-center py-3 text-h6 font-weight-light">{{ card.title }}</h2>
            </div>
            <div class="flex justify-center">
              <h4 class="text-center px-2 font-weight-light">{{ card.subtitle }}</h4>
            </div>
          </Stack>
        </v-card-text>
      </v-card>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.v-card {
  &.v-theme--dark {
    background-image: linear-gradient(180deg, rgba(51, 51, 51, 1) 0%, rgba(0, 0, 0, 1) 100%);
  }
  &.v-theme--light {
    background-image: linear-gradient(180deg, rgba(248, 248, 248, 1) 0%, rgba(255, 255, 255, 1) 100%);
  }
}
.active {
  border: 2px solid rgb(var(--v-theme-primary)) !important;
}
</style>
