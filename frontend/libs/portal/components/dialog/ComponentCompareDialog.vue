<!-- SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG -->
<!---->
<!-- SPDX-License-Identifier: Apache-2.0 -->

<script lang="ts">
import {SpdxIdentifier} from '@disclosure-portal/model/Spdx';
import {
  ComponentDiff,
  ComponentDiffType,
  ComponentDiffWrapper,
  ComponentInfo,
  ComponentMultiDiff,
} from '@disclosure-portal/model/VersionDetails';
import {getIconForDiffType} from '@disclosure-portal/utils/View';
import DCloseButton from '@shared/components/disco/DCloseButton.vue';
import {DataTableHeader} from '@shared/types/table';
import _ from 'lodash';
import {defineComponent, ref} from 'vue';
import {useI18n} from 'vue-i18n';

class ComponentCompareDetails {
  public key: string;
  public previous = '';
  public current = '';
  public isDiff = false;

  constructor(key: string, current: string, previous: string, isChanged: boolean, globalDiffType: ComponentDiffType) {
    this.key = key;
    switch (globalDiffType) {
      case ComponentDiffType.NEW:
        this.current = current;
        break;
      case ComponentDiffType.REMOVED:
        this.previous = current;
        break;
      case ComponentDiffType.CHANGED:
        this.current = current;
        this.previous = previous;
        this.isDiff = isChanged;
        break;
      default:
        this.current = current;
        this.previous = previous;
        break;
    }
  }
}

export default defineComponent({
  components: {
    DCloseButton,
  },
  setup() {
    const {t} = useI18n();
    const show = ref(false);
    const details = ref<ComponentMultiDiff>(new ComponentMultiDiff());
    const name = ref('');
    const selectedTab = ref(1);
    const componentDetails = ref<ComponentCompareDetails[]>([]);
    const sbomCurrent = ref('');
    const sbomPrevious = ref('');
    const headers = ref<DataTableHeader[]>([]);
    const DiffType = ref(ComponentDiffType.UNCHANGED);
    const currentVersionList = ref<string[]>([]);
    const previousVersionList = ref<string[]>([]);
    const selectedCurrentVersion = ref('');
    const selectedPreviousVersion = ref('');
    const openPanel = ref(1000000);

    function toYesOrNo(value: boolean): string {
      return t(value ? 'YES' : 'NO');
    }

    const setHeader = (withPrevious: boolean, withCurrent: boolean) => {
      headers.value = [
        {
          title: t('COL_COMP_COMPARE_DLG_KEY'),
          align: 'start',
          filterable: true,
          class: 'tableHeaderCell',
          value: 'key',
        },
      ];

      const widthWithOneColumn = 690;
      const widthWithTwoColumns = widthWithOneColumn / 2;

      if (withCurrent) {
        headers.value.push({
          title: t('COL_COMP_COMPARE_DLG_CURRENT'),
          align: 'start',
          filterable: true,
          class: 'tableHeaderCell',
          value: 'current',
          width: withPrevious ? widthWithTwoColumns : widthWithOneColumn,
        });
      }

      if (withPrevious) {
        headers.value.push({
          title: t('COL_COMP_COMPARE_DLG_PREVIOUS'),
          align: 'start',
          filterable: true,
          class: 'tableHeaderCell',
          value: 'previous',
          width: withCurrent ? widthWithTwoColumns : widthWithOneColumn,
        });
      }
    };

    const open = (
      detailsData: ComponentMultiDiff,
      nameData: string,
      spdxMetaPrevious: SpdxIdentifier,
      spdxMetaCurrent: SpdxIdentifier,
    ) => {
      details.value = detailsData;
      sbomCurrent.value = `[${spdxMetaCurrent.uploaded}] ${spdxMetaCurrent.label}`;
      sbomPrevious.value = `[${spdxMetaPrevious.uploaded}] ${spdxMetaPrevious.label}`;

      setHeader(details.value.DiffType !== ComponentDiffType.NEW, details.value.DiffType !== ComponentDiffType.REMOVED);

      show.value = true;
      name.value = nameData;
      componentDetails.value = [];
      selectedTab.value = 0;
      DiffType.value = details.value.DiffType;

      let foundCurrentVersion: ComponentInfo | undefined;
      let foundPreviousVersion: ComponentInfo | undefined;

      currentVersionList.value = _.orderBy(details.value.ComponentsNew.map((comp) => comp.version));
      if (currentVersionList.value.length > 0) {
        selectedCurrentVersion.value = currentVersionList.value[0];
        if (selectedCurrentVersion.value) {
          foundCurrentVersion = _.find(details.value.ComponentsNew, {version: selectedCurrentVersion.value});
        }
      }
      previousVersionList.value = _.orderBy(details.value.ComponentsOld.map((comp) => comp.version));
      if (previousVersionList.value.length > 0) {
        selectedPreviousVersion.value = previousVersionList.value[0];
        if (selectedPreviousVersion.value) {
          foundPreviousVersion = _.find(details.value.ComponentsOld, {version: selectedPreviousVersion.value});
        }
      }

      prepareDataToCompare(foundCurrentVersion, foundPreviousVersion);
    };

    const selectedVersionChanged = () => {
      componentDetails.value = [];
      let foundCurrentVersion: ComponentInfo | undefined;
      let foundPreviousVersion: ComponentInfo | undefined;

      if (selectedCurrentVersion.value) {
        foundCurrentVersion = _.find(details.value.ComponentsNew, {version: selectedCurrentVersion.value});
      }
      if (selectedPreviousVersion.value) {
        foundPreviousVersion = _.find(details.value.ComponentsOld, {version: selectedPreviousVersion.value});
      }

      prepareDataToCompare(foundCurrentVersion, foundPreviousVersion);
    };

    const prepareDataToCompare = (
      foundCurrentVersion: ComponentInfo | undefined,
      foundPreviousVersion: ComponentInfo | undefined,
    ) => {
      if (foundCurrentVersion || foundPreviousVersion) {
        const componentDiff = new ComponentDiff();
        componentDiff.DiffType = details.value.DiffType;
        componentDiff.ComponentOld = foundPreviousVersion ? foundPreviousVersion : new ComponentInfo();
        componentDiff.ComponentNew = foundCurrentVersion ? foundCurrentVersion : new ComponentInfo();
        if (details.value.Changes[`${selectedPreviousVersion.value}_${selectedCurrentVersion.value}`]) {
          const changes = details.value.Changes[`${selectedPreviousVersion.value}_${selectedCurrentVersion.value}`];
          componentDiff.SpdxId = changes.SpdxId;
          componentDiff.Name = changes.Name;
          componentDiff.Version = changes.Version;
          componentDiff.LicenseComments = changes.LicenseComments;
          componentDiff.LicenseDeclared = changes.LicenseDeclared;
          componentDiff.License = changes.License;
          componentDiff.LicenseEffective = changes.LicenseEffective;
          componentDiff.CopyrightText = changes.CopyrightText;
          componentDiff.Description = changes.Description;
          componentDiff.DownloadLocation = changes.DownloadLocation;
          componentDiff.prStatus = changes.prStatus;
          componentDiff.Type = changes.Type;
          componentDiff.Modified = changes.Modified;
          componentDiff.Questioned = changes.Questioned;
          componentDiff.Unasserted = changes.Unasserted;
          componentDiff.PURL = changes.PURL;
        }

        const diffWrapper = new ComponentDiffWrapper(componentDiff);

        componentDetails.value.push(
          new ComponentCompareDetails('SBOM', sbomCurrent.value, sbomPrevious.value, false, DiffType.value),
        );
        componentDetails.value.push(
          new ComponentCompareDetails(
            'Name',
            diffWrapper.name,
            diffWrapper.getOtherComponent().name,
            diffWrapper.diff.Name,
            DiffType.value,
          ),
        );
        componentDetails.value.push(
          new ComponentCompareDetails(
            'Status',
            diffWrapper.prStatus,
            diffWrapper.getOtherComponent().prStatus,
            diffWrapper.diff.prStatus,
            DiffType.value,
          ),
        );
        componentDetails.value.push(
          new ComponentCompareDetails(
            'Package-URL',
            diffWrapper.purl,
            diffWrapper.getOtherComponent().purl,
            diffWrapper.diff.PURL,
            DiffType.value,
          ),
        );
        componentDetails.value.push(
          new ComponentCompareDetails(
            'Version',
            diffWrapper.version,
            diffWrapper.getOtherComponent().version,
            diffWrapper.diff.Version,
            DiffType.value,
          ),
        );
        componentDetails.value.push(
          new ComponentCompareDetails(
            'License Declared',
            diffWrapper.licenseDeclared,
            diffWrapper.getOtherComponent().licenseDeclared,
            diffWrapper.diff.LicenseDeclared,
            DiffType.value,
          ),
        );
        componentDetails.value.push(
          new ComponentCompareDetails(
            'License Concluded',
            diffWrapper.license,
            diffWrapper.getOtherComponent().license,
            diffWrapper.diff.License,
            DiffType.value,
          ),
        );
        componentDetails.value.push(
          new ComponentCompareDetails(
            'License Effective',
            diffWrapper.licenseEffective,
            diffWrapper.getOtherComponent().licenseEffective,
            diffWrapper.diff.LicenseEffective,
            DiffType.value,
          ),
        );
        componentDetails.value.push(
          new ComponentCompareDetails(
            'Copyright Text',
            diffWrapper.copyrightText,
            diffWrapper.getOtherComponent().copyrightText,
            diffWrapper.diff.CopyrightText,
            DiffType.value,
          ),
        );
        componentDetails.value.push(
          new ComponentCompareDetails(
            'Download Location',
            diffWrapper.downloadLocation,
            diffWrapper.getOtherComponent().downloadLocation,
            diffWrapper.diff.DownloadLocation,
            DiffType.value,
          ),
        );
        componentDetails.value.push(
          new ComponentCompareDetails(
            'Description',
            diffWrapper.description,
            diffWrapper.getOtherComponent().description,
            diffWrapper.diff.Description,
            DiffType.value,
          ),
        );
        componentDetails.value.push(
          new ComponentCompareDetails(
            'License Comments',
            diffWrapper.licenseComments,
            diffWrapper.getOtherComponent().licenseComments,
            diffWrapper.diff.LicenseComments,
            DiffType.value,
          ),
        );
        componentDetails.value.push(
          new ComponentCompareDetails(
            'SPDX ID',
            diffWrapper.spdxId,
            diffWrapper.getOtherComponent().spdxId,
            diffWrapper.diff.SpdxId,
            DiffType.value,
          ),
        );
        componentDetails.value.push(
          new ComponentCompareDetails(
            'Type',
            diffWrapper.type,
            diffWrapper.getOtherComponent().type,
            diffWrapper.diff.Type,
            DiffType.value,
          ),
        );
        componentDetails.value.push(
          new ComponentCompareDetails(
            'Modified',
            toYesOrNo(diffWrapper.modified),
            toYesOrNo(diffWrapper.getOtherComponent().modified),
            diffWrapper.diff.Modified,
            DiffType.value,
          ),
        );
        componentDetails.value.push(
          new ComponentCompareDetails(
            'Questioned',
            toYesOrNo(diffWrapper.questioned),
            toYesOrNo(diffWrapper.getOtherComponent().questioned),
            diffWrapper.diff.Questioned,
            DiffType.value,
          ),
        );
        componentDetails.value.push(
          new ComponentCompareDetails(
            'Unasserted',
            toYesOrNo(diffWrapper.unasserted),
            toYesOrNo(diffWrapper.getOtherComponent().unasserted),
            diffWrapper.diff.Unasserted,
            DiffType.value,
          ),
        );
      }
    };

    const getCssClass = (data: {item: ComponentCompareDetails}): Record<string, any> => {
      if (data.item.isDiff) {
        return {
          class: {
            'diff-component': true,
          },
        };
      }
      return {};
    };

    const close = () => {
      show.value = false;
    };

    return {
      t,
      show,
      details,
      name,
      selectedTab,
      componentDetails,
      sbomCurrent,
      sbomPrevious,
      headers,
      DiffType,
      currentVersionList,
      previousVersionList,
      selectedCurrentVersion,
      selectedPreviousVersion,
      openPanel,
      setHeader,
      open,
      selectedVersionChanged,
      prepareDataToCompare,
      getCssClass,
      close,
      getIconForDiffType,
    };
  },
});
</script>

<template>
  <v-dialog v-model="show" content-class="large" scrollable style="max-width: 1000px">
    <v-card class="pa-8" flat>
      <v-card-title>
        <v-row>
          <v-col cols="10">
            <v-row>
              <span class="text-h5 mr-2"> {{ t('CD_COMPARE') }}: {{ name }} </span>
              <v-icon color="secondaryTextColor" class="pr-2" v-if="getIconForDiffType(DiffType) !== ''">{{
                getIconForDiffType(DiffType)
              }}</v-icon>
            </v-row>

            <v-row>
              <span class="text-h6"> {{ t(DiffType) }} </span>
            </v-row>
          </v-col>
          <v-col cols="2" align="right">
            <DCloseButton @click="close" />
          </v-col>
        </v-row>
      </v-card-title>
      <v-card-text>
        <v-row v-if="currentVersionList.length > 1 || previousVersionList.length > 1">
          <v-col v-if="currentVersionList.length > 1">
            <v-select
              variant="outlined"
              density="compact"
              :items="currentVersionList"
              v-model="selectedCurrentVersion"
              :label="t('SBOM_COMPARE_CURRENT_VERSION')"
              @update:modelValue="selectedVersionChanged">
            </v-select>
          </v-col>
          <v-col v-if="previousVersionList.length > 1">
            <v-select
              variant="outlined"
              density="compact"
              :items="previousVersionList"
              v-model="selectedPreviousVersion"
              :label="t('SBOM_COMPARE_PREVIOUS_VERSION')"
              @update:modelValue="selectedVersionChanged">
            </v-select>
          </v-col>
        </v-row>
        <v-tabs v-model="selectedTab" slider-color="brand" active-class="active" show-arrows bg-color="tabsHeader">
          <v-tab value="attributes">
            {{ t('TAB_TITLE_ATTRIBUTES') }}
          </v-tab>
        </v-tabs>
        <v-tabs-window style="max-height: 400px; overflow-y: scroll" v-model="selectedTab">
          <v-tabs-window-item value="attributes">
            <v-data-table
              class="striped-table"
              density="compact"
              hide-default-footer
              :items-per-page="100000"
              fixed-header
              :headers="headers"
              :items="componentDetails"
              :row-props="getCssClass">
              <template v-slot:item.current="{item}">
                <span class="tableCompareCell">{{ item.current }}</span>
              </template>
              <template v-slot:item.previous="{item}">
                <span class="tableCompareCell">{{ item.previous }}</span>
              </template>
            </v-data-table>
          </v-tabs-window-item>
        </v-tabs-window>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn @click="close" depressed color="primary">
          {{ t('BTN_CLOSE') }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>
