// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

<template>
  <v-table fixed-header density="compact" class="borderTable tableNoHandCursor" v-if="stats">
    <template v-slot:default>
      <thead>
        <tr>
          <th class="text-left">
            {{ t('COL_NAME') }}
          </th>
          <th class="text-center" v-for="serverResponse in stats" :key="serverResponse._key">
            {{ formatDateAndTime(serverResponse.created) }}
          </th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td>Projects</td>
          <td class="text-center" v-for="serverResponse in stats" :key="serverResponse._key">
            {{ serverResponse.projectCount }}<span v-if="serverResponse.missingProjects" class="error--text"> !!!</span>
          </td>
        </tr>
        <tr>
          <td>Projects (active)</td>
          <td class="text-center" v-for="serverResponse in stats" :key="serverResponse._key">
            {{ serverResponse.projectActiveCount }}
          </td>
        </tr>
        <tr>
          <td>Projects (deleted)</td>
          <td class="text-center" v-for="serverResponse in stats" :key="serverResponse._key">
            {{ serverResponse.projectDeletedCount }}
          </td>
        </tr>
        <tr>
          <td>max. Versions in one Project</td>
          <td class="text-center" v-for="serverResponse in stats" :key="serverResponse._key">
            {{ serverResponse.maxVersionsInOneProject }}
          </td>
        </tr>
        <tr>
          <td>projects over or at version limit</td>
          <td class="text-center" v-for="serverResponse in stats" :key="serverResponse._key">
            {{ serverResponse.projectsOverOrAtVersionLimit }} (limit: {{ serverResponse.versionLimit }})
          </td>
        </tr>
        <tr>
          <td>Licenses</td>
          <td class="text-center" v-for="serverResponse in stats" :key="serverResponse._key">
            {{ serverResponse.licenseCount }}<span v-if="serverResponse.missingLicenses" class="error--text"> !!!</span>
          </td>
        </tr>
        <tr>
          <td>Licenses (active)</td>
          <td class="text-center" v-for="serverResponse in stats" :key="serverResponse._key">
            {{ serverResponse.licenseActiveCount }}
          </td>
        </tr>
        <tr>
          <td>Licenses (Chart)</td>
          <td class="text-center" v-for="serverResponse in stats" :key="serverResponse._key">
            {{ serverResponse.licenseChartCount }}
          </td>
        </tr>
        <tr>
          <td>Licenses (deleted)</td>
          <td class="text-center" v-for="serverResponse in stats" :key="serverResponse._key">
            {{ serverResponse.licenseDeletedCount }}
          </td>
        </tr>
        <tr>
          <td>Policy Rules</td>
          <td class="text-center" v-for="serverResponse in stats" :key="serverResponse._key">
            {{ serverResponse.policyRuleCount
            }}<span v-if="serverResponse.missingPolicyRules" class="error--text"> !!!</span>
          </td>
        </tr>
        <tr>
          <td>Policy Rules (active)</td>
          <td class="text-center" v-for="serverResponse in stats" :key="serverResponse._key">
            {{ serverResponse.policyRuleActiveCount }}
          </td>
        </tr>
        <tr>
          <td>Policy Rules (deleted)</td>
          <td class="text-center" v-for="serverResponse in stats" :key="serverResponse._key">
            {{ serverResponse.policyRuleDeletedCount }}
          </td>
        </tr>
        <tr>
          <td>Labels</td>
          <td class="text-center" v-for="serverResponse in stats" :key="serverResponse._key">
            {{ serverResponse.labelCount }}
          </td>
        </tr>
        <tr>
          <td>Schemas</td>
          <td class="text-center" v-for="serverResponse in stats" :key="serverResponse._key">
            {{ serverResponse.schemaCount }}
          </td>
        </tr>
        <tr>
          <td>Classifications</td>
          <td class="text-center" v-for="serverResponse in stats" :key="serverResponse._key">
            {{ serverResponse.obligationCount
            }}<span v-if="serverResponse.missingObligations" class="error--text"> !!!</span>
          </td>
        </tr>
        <tr>
          <td>Classifications (active)</td>
          <td class="text-center" v-for="serverResponse in stats" :key="serverResponse._key">
            {{ serverResponse.obligationActiveCount }}
          </td>
        </tr>
        <tr>
          <td>Classifications (deleted)</td>
          <td class="text-center" v-for="serverResponse in stats" :key="serverResponse._key">
            {{ serverResponse.obligationDeletedCount }}
          </td>
        </tr>
        <tr>
          <td>Users</td>
          <td class="text-center" v-for="serverResponse in stats" :key="serverResponse._key">
            {{ serverResponse.userCount }}<span v-if="serverResponse.missingUsers" class="error--text"> !!!</span>
          </td>
        </tr>
        <tr>
          <td>Users (active)</td>
          <td class="text-center" v-for="serverResponse in stats" :key="serverResponse._key">
            {{ serverResponse.userActiveCount }}
          </td>
        </tr>
        <tr>
          <td>Users (disabled)</td>
          <td class="text-center" v-for="serverResponse in stats" :key="serverResponse._key">
            {{ serverResponse.userDeactivateCount }}
          </td>
        </tr>
        <tr>
          <td>Users (TermsNotAccepted)</td>
          <td class="text-center" v-for="serverResponse in stats" :key="serverResponse._key">
            {{ serverResponse.userTermsNotAcceptedCount }}
          </td>
        </tr>
        <tr>
          <td>Users (Deprovisioned)</td>
          <td class="text-center" v-for="serverResponse in stats" :key="serverResponse._key">
            {{ serverResponse.userDeprovisionedCount }}
          </td>
        </tr>
        <tr>
          <td>Uploaded Files (all)</td>
          <td class="text-center" v-for="serverResponse in stats" :key="serverResponse._key">
            {{ serverResponse.uploadFileCnt
            }}<span v-if="serverResponse.missingUploadFiles" class="error--text"> !!!</span>
          </td>
        </tr>
        <tr>
          <td>Uploaded Files (sbom)</td>
          <td class="text-center" v-for="serverResponse in stats" :key="serverResponse._key">
            {{ serverResponse.uploadFileCntSBOM }}
          </td>
        </tr>
        <tr>
          <td>Uploaded Files (pdf)</td>
          <td class="text-center" v-for="serverResponse in stats" :key="serverResponse._key">
            {{ serverResponse.uploadFileCntPDF }}
          </td>
        </tr>
        <tr>
          <td>Uploaded Files (json)</td>
          <td class="text-center" v-for="serverResponse in stats" :key="serverResponse._key">
            {{ serverResponse.uploadFileCntJSON }}
          </td>
        </tr>
        <tr>
          <td>DB Backup Files</td>
          <td class="text-center" v-for="serverResponse in stats" :key="serverResponse._key">
            {{ serverResponse.dbBackupFileCnt }}
          </td>
        </tr>
      </tbody>
    </template>
  </v-table>
</template>

<script setup lang="ts">
import {SystemStats} from '@disclosure-portal/model/Statistic';
import {formatDateAndTime} from '@disclosure-portal/utils/Table';
import {useI18n} from 'vue-i18n';

defineProps<{stats: SystemStats[]}>();

const {t} = useI18n();
</script>

<style scoped></style>
