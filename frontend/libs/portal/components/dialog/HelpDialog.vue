<template>
  <DSimpleDialog v-model="modelValue" :title="title">
    <Stack direction="row" align="start" class="gap-4">
      <div ref="contentRef" :class="tableOfContents.length > 0 ? 'flex-[3]' : 'flex-1'" class="pr-4 border-r border-grey-70/50">
        <Markdown :text="enhancedText"></Markdown>
      </div>
      <div v-if="tableOfContents.length > 0" class="flex-1 self-start sticky top-0 pl-4">
        <Stack class="max-h-[70vh]">
          <div class="text-sm font-semibold sticky top-0">Table of Contents</div>
          <nav class="flex flex-col">
            <a
              v-for="item in tableOfContents"
              :key="item.id"
              :href="`#${item.id}`"
              :class="[
                'block my-1 text-sm border-l-2 border-transparent text-font',
                'hover:text-yellow-45 hover:border-l-yellow-45',
                {
                  'text-yellow-45 border-l-yellow-45': activeId === item.id,
                  'pl-1 font-semibold': item.level === 1,
                  'pl-3': item.level === 2,
                  'pl-5 text-xs': item.level === 3,
                  'pl-7 text-xs': item.level === 4,
                  'pl-9 text-[11px]': item.level === 5,
                  'pl-[60px] text-[11px]': item.level === 6,
                },
              ]"
              @click.prevent="onTocItemClick(item.id)">
              {{ item.text }}
            </a>
          </nav>
        </Stack>
      </div>
    </Stack>
  </DSimpleDialog>
</template>

<script setup lang="ts">
import Markdown from '@shared/components/Markdown.vue';
import DSimpleDialog from '@shared/components/disco/DSimpleDialog.vue';
import Stack from '@shared/layouts/Stack.vue';
import {computed, ref, watchEffect} from 'vue';

interface Props {
  title: string;
  text: string;
}

interface TocItem {
  id: string;
  text: string;
  level: number;
}

const modelValue = defineModel<boolean>({required: true});
const props = defineProps<Props>();

const contentRef = ref<HTMLElement | null>(null);
const activeId = ref<string>('');

const cleanHeaderText = (text: string): string => {
  return text
    .replace(/<[^>]*>/g, '')
    .replace(/:\s*$/g, '')
    .trim();
};

const generateId = (text: string): string => {
  return text
    .toLowerCase()
    .replace(/[^\w\s-]/g, '')
    .replace(/\s+/g, '-');
};

const parseHeaderLine = (line: string) => {
  const match = line.match(/^(#{1,6})\s+(.+)$/);
  if (!match) return null;

  const level = match[1].length;
  const rawText = match[2].trim();
  const cleanedText = cleanHeaderText(rawText);
  const id = generateId(cleanedText);

  return {level, rawText, cleanedText, id};
};

const tableOfContents = computed<TocItem[]>(() => {
  const headings: TocItem[] = [];
  const lines = props.text.split('\n');

  lines.forEach((line) => {
    const parsed = parseHeaderLine(line);
    if (parsed) {
      headings.push({
        id: parsed.id,
        text: parsed.cleanedText,
        level: parsed.level,
      });
    }
  });

  return headings;
});

watchEffect(() => {
  if (tableOfContents.value.length > 0 && !activeId.value) {
    activeId.value = tableOfContents.value[0].id;
  }
});

const enhancedText = computed(() => {
  const lines = props.text.split('\n');

  const enhancedLines = lines.map((line) => {
    const parsed = parseHeaderLine(line);
    if (parsed) {
      return `${'#'.repeat(parsed.level)} <span id="${parsed.id}"></span>${parsed.rawText}`;
    }
    return line;
  });

  return enhancedLines.join('\n');
});

const scrollToHeading = (id: string) => {
  const element = document.getElementById(id);
  if (element) {
    element.scrollIntoView({behavior: 'smooth', block: 'start'});
  }
};

const onTocItemClick = (id: string) => {
  scrollToHeading(id);
  activeId.value = id;
};
</script>
