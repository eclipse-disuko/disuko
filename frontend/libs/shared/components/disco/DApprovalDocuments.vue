// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<script setup lang="ts">
import {DocumentDto} from '@disclosure-portal/model/Document';
import ProjectService from '@disclosure-portal/services/projects';
import {useProjectStore} from '@disclosure-portal/stores/project.store';
import {formatDateAndTime} from '@disclosure-portal/utils/Table';
import {formatDateTimeForFile, formatDateTimeShort} from '@disclosure-portal/utils/View';
import DCopyClipboardButton from '@shared/components/disco/DCopyClipboardButton.vue';
import DDateCellWithTooltip from '@shared/components/disco/DDateCellWithTooltip.vue';
import DIconButton from '@shared/components/disco/DIconButton.vue';
import dayjs from 'dayjs';
import {computed} from 'vue';
import {useI18n} from 'vue-i18n';

const {t} = useI18n();
const currentProject = computed(() => useProjectStore().currentProject!);

const props = defineProps<{
  documents: DocumentDto[];
}>();

const languageMap = {
  de: 'German',
  en: 'English',
};

const typeMap = {
  disclosure: 'Disclosure Document',
  policies: 'Policy Rules',
  policycheck: 'Policy Rules Check',
  archive: 'Archive',
};

const relevantDocuments = computed(() => {
  if (!props.documents) {
    return [];
  }
  const sortedDocuments = props.documents.sort((a, b) => a.updated.localeCompare(b.updated));
  sortedDocuments.reverse();
  const documentConditions = [
    (doc: DocumentDto) => doc.type === 'disclosure' && doc.lang === 'de',
    (doc: DocumentDto) => doc.type === 'disclosure' && doc.lang === 'en',
    (doc: DocumentDto) => doc.type === 'policies',
    (doc: DocumentDto) => doc.type === 'policycheck',
    (doc: DocumentDto) => doc.type === 'archive',
  ];

  const relevantDocuments = [];
  for (const condition of documentConditions) {
    const lastDoc = sortedDocuments.find(condition);
    if (lastDoc) {
      // Add chain for 'disclosure' documents of lang 'en' or 'de'
      if (lastDoc.type === 'disclosure' && ['en', 'de'].includes(lastDoc.lang)) {
        // Find all previous 'disclosure' documents for the same lang and older updated dates
        const prevDocs = sortedDocuments.filter(
          (doc) => doc.type === 'disclosure' && doc.lang === lastDoc.lang && doc.updated < lastDoc.updated,
        );

        // Map previous documents to ChainDetails and sort by 'updated' descending
        const chain = prevDocs
          .map((doc) => ({
            updated: doc.updated, // Directly use the Date object
            hash: doc.hash,
          }))
          .sort((a, b) => a.updated.localeCompare(b.updated));

        lastDoc.chain = chain;
      }

      relevantDocuments.push(lastDoc);
    }
  }

  return relevantDocuments;
});

const downloadDocument = (item: DocumentDto) => {
  const link = document.createElement('a');
  link.click();
  link.target = '_blank';
  link.rel = 'noopener noreferrer';
  ProjectService.downloadDocumentByTask(currentProject.value._key, item.approvalId, item.type, item.lang, item.version)
    .then((res) => {
      switch (item.type) {
        case 'disclosure': {
          link.download = 'Disclosure_' + item.approvalId + '_' + item.lang + '.pdf';
          break;
        }
        case 'policies': {
          link.download = 'PolicyRules_' + item.approvalId + '.json';
          break;
        }
        case 'policycheck': {
          link.download = 'PolicyCheck_' + item.approvalId + '.json';
          break;
        }
        case 'archive': {
          link.download =
            'Archive_' + item.created + item.approvalId + '_' + formatDateTimeForFile(item.created) + '.zip';
          break;
        }
        default: {
          link.download = 'UnknownType_' + item.approvalId + '.txt';
          break;
        }
      }
      link.href = URL.createObjectURL(new Blob([res.data as unknown as BlobPart]));
      link.click();
    })
    .catch((e) => {
      console.error('Error downloading document', e);
    });
};

const getReferenceInfoForClipboard = (item: DocumentDto): string => {
  let referenceInfo = `Disclosure Portal File Reference

Project Name: ${currentProject.value.name}
Project Identifier: ${currentProject.value._key}
Reference Timestamp: ${formatDateAndTime(dayjs().toISOString())} (UTC)
File Name: ${item.fileName}
File Identifier: ${item._key}
Created Date: ${formatDateTimeShort(item.created, true)} (UTC)
Updated Date: ${formatDateTimeShort(item.updated, true)} (UTC)
File SHA-256: ${item.hash}`;

  if (item.chain && item.chain.length > 0) {
    referenceInfo += `\nHash chain:\n`;
    referenceInfo += item.chain.map((chainItem) => `\tHash: ${chainItem.hash}`).join(',\n');
  }

  return referenceInfo;
};
</script>

<template>
  <div>
    <v-table>
      <thead>
        <tr>
          <th class="text-left">
            {{ t('COL_ACTIONS') }}
          </th>
          <th class="text-left">
            {{ t('COL_DOCUMENT') }}
          </th>
          <th class="text-left">
            {{ t('COL_UPDATED') }}
          </th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="doc in relevantDocuments" :key="doc._key">
          <td>
            <DCopyClipboardButton
              :tableButton="true"
              class="ml-2"
              :hint="t('TT_COPY_REFERENCE_INFO')"
              :content="getReferenceInfoForClipboard(doc)" />
            <DIconButton icon="mdi-download" :hint="t('TT_download')" @clicked="downloadDocument(doc)" />
          </td>
          <td>{{ typeMap[doc.type] }} {{ languageMap[doc.lang] ? `(${languageMap[doc.lang]})` : `` }}</td>
          <td>
            <DDateCellWithTooltip :value="doc.updated"></DDateCellWithTooltip>
          </td>
        </tr>
      </tbody>
    </v-table>
  </div>
</template>
