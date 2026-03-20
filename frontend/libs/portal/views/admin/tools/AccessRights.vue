<script setup lang="ts">
import Icons from '@disclosure-portal/constants/icons';
import {AccessRights, ProjectAccessRights} from '@disclosure-portal/model/AccessRights';
import {ActionRights} from '@disclosure-portal/model/Credentials';
import {CRUDRights} from '@disclosure-portal/model/Rights';
import {IMap} from '@disclosure-portal/utils/View';
import {DataTableHeader} from '@shared/types/table';
import {TOOLTIP_OPEN_DELAY_IN_MS} from '@shared/utils/constant';
import _ from 'lodash';
import {nextTick, onMounted, ref} from 'vue';
import {useI18n} from 'vue-i18n';
import AdminService from '../../../services/admin';

class RoleWrapper {
  public value: string;
  public text: string;

  constructor(value: string, text: string) {
    this.value = value;
    this.text = text;
  }
}

class ProjectAccessRightsEntry {
  public accessRightName = '';
  public entry = new Map<string, CRUDRights | ActionRights>();
  public group = '';

  constructor(accessRightName: string, key: string, value: CRUDRights | ActionRights, group: string) {
    this.accessRightName = accessRightName;
    this.entry.set(key, value);
    this.group = group;
  }
}

class AccessRightsEntry {
  public accessRightName = '';
  public entry = new Map<string, CRUDRights>();
  public group = '';

  constructor(accessRightName: string, key: string, value: CRUDRights, group: string) {
    this.accessRightName = accessRightName;
    this.entry.set(key, value);
    this.group = group;
  }
}

const {t} = useI18n();
const allProjectAccessRights = ref<IMap<ProjectAccessRights>>({});
const projectAccessRightsHeaders = ref<DataTableHeader[]>([
  {
    title: t('ACCESS_RIGHTS'),
    width: 250,
    align: 'start',
    class: 'tableHeaderCell',
    value: 'accessRightName',
  },
]);
const allProjectAccessRightsForTable = ref<ProjectAccessRightsEntry[]>([]);

const allAccessRights = ref<IMap<AccessRights>>({});
const defaultAccessRightsHeaders = ref<DataTableHeader[]>([
  {
    title: t('ACCESS_RIGHTS'),
    width: 250,
    align: 'start',
    class: 'tableHeaderCell',
    value: 'accessRightName',
  },
]);
const allDynamicAccessRightsHeaders = ref<DataTableHeader[]>([]);
const allAccessRightsForTable = ref<AccessRightsEntry[]>([]);

const accessRightsLoaded = ref(false);

const allRoles = ref<RoleWrapper[]>([]);
const selectedRoles = ref<string[]>([]);
const selectedRolesDefault = ['ApplicationAdmin', 'DomainAdmin', 'Internal'];
const accessRightsHeaderForSelection = ref<DataTableHeader[]>([]);
const accessRightsForTableForSelection = ref<AccessRightsEntry[]>([]);
const dirtySelection = ref(false);
const defaultSelection = ref(true);

const icons = Icons;

const getProjectAccessRights = async () => {
  allProjectAccessRights.value = (await AdminService.getAllProjectAccessRights()).data;
  allAccessRights.value = (await AdminService.getAllAccessRights()).data;
  accessRightsLoaded.value = true;
  if (allProjectAccessRights.value) {
    const addedProjectAccessRightsHeaders: DataTableHeader[] = [];
    _.each(allProjectAccessRights.value, (accessRights, keyRole: string) => {
      addedProjectAccessRightsHeaders.push({
        title: t('USER_ROLE_' + keyRole),
        align: 'start',
        class: 'tableHeaderCell',
        value: keyRole,
      });
      _.each(accessRights, (crudRights, accessRightNameKey) => {
        const presentItem = _.find(allProjectAccessRightsForTable.value, {accessRightName: accessRightNameKey});
        if (presentItem) {
          const i = _.indexOf(allProjectAccessRightsForTable.value, presentItem);
          allProjectAccessRightsForTable.value[i].entry.set(keyRole, crudRights);
        } else {
          let group = '';
          if (_.size(crudRights) === 4) {
            group = 'CRUD';
          } else if (_.size(crudRights) === 3) {
            group = 'Upload Download Delete';
          }
          const item = new ProjectAccessRightsEntry(accessRightNameKey, keyRole, crudRights, group);
          allProjectAccessRightsForTable.value.push(item);
        }
      });
    });
    projectAccessRightsHeaders.value.push(
      ..._.orderBy(addedProjectAccessRightsHeaders, (el) => forProject(el.value), 'asc'),
    );
    allProjectAccessRightsForTable.value = _.orderBy(
      allProjectAccessRightsForTable.value,
      [(item) => (item.group === 'CRUD' ? 0 : 1)],
      ['asc']
    );
    allProjectAccessRightsForTable.value.forEach((item) => {
      item.entry = new Map(_.orderBy([...item.entry.entries()], (el) => forProject(el[0]), 'asc'));
    });
  }
  if (allAccessRights.value) {
    const addedAccessRightsHeaders: DataTableHeader[] = [];
    _.each(allAccessRights.value, (accessRights, keyRole: string) => {
      addedAccessRightsHeaders.push({
        title: t('USER_ROLE_' + keyRole),
        align: 'start',
        class: 'tableHeaderCell',
        value: keyRole,
      });
      _.each(accessRights, (crudRights, accessRightNameKey) => {
        const presentItem = _.find(allAccessRightsForTable.value, {accessRightName: accessRightNameKey});
        if (presentItem) {
          const i = _.indexOf(allAccessRightsForTable.value, presentItem);
          allAccessRightsForTable.value[i].entry.set(keyRole, crudRights);
        } else {
          const item = new AccessRightsEntry(accessRightNameKey, keyRole, crudRights, 'CRUD');
          allAccessRightsForTable.value.push(item);
        }
      });
    });
    allDynamicAccessRightsHeaders.value.push(
      ..._.orderBy(addedAccessRightsHeaders, (el) => forFurther(el.value), 'asc'),
    );
    allRoles.value = addedAccessRightsHeaders.map((header) => new RoleWrapper(header.value, header.title));
    selectedRoles.value = selectedRolesDefault;
    allAccessRightsForTable.value.forEach((item) => {
      item.entry = new Map(_.orderBy([...item.entry.entries()], (el) => forFurther(el[0]), 'asc'));
    });
    prepareDataForSelection();
  }
};

const prepareDataForDefaultSelection = () => {
  selectedRoles.value = selectedRolesDefault;
  prepareDataForSelection();
  defaultSelection.value = true;
};

const prepareDataForSelection = () => {
  accessRightsHeaderForSelection.value = [];
  accessRightsHeaderForSelection.value.push(...defaultAccessRightsHeaders.value);
  accessRightsHeaderForSelection.value.push(
    ...allDynamicAccessRightsHeaders.value.filter((h) => selectedRoles.value.some((r) => r === h.value)),
  );

  accessRightsForTableForSelection.value = [];
  const selected = _.cloneDeep(allAccessRightsForTable.value);
  selected.forEach((item) => {
    item.entry = new Map([...item.entry.entries()].filter((e) => selectedRoles.value.some((r) => r === e[0])));
  });
  accessRightsForTableForSelection.value.push(...selected);
  dirtySelection.value = false;
  defaultSelection.value = areArraysEqual(selectedRoles.value, selectedRolesDefault);
};

const areArraysEqual = (a1: string[], a2: string[]) => {
  return (
    a1.length === a2.length &&
    a1.every(function (e, i) {
      return e === a2[i];
    })
  );
};

const selectedRolesChanged = () => {
  nextTick(() => {
    dirtySelection.value = true;
    defaultSelection.value = areArraysEqual(selectedRoles.value, selectedRolesDefault);
  });
};

const forProject = (v: string): number => {
  switch (v) {
    case 'Owner':
      return 0;
    case 'Supplier':
      return 1;
    case 'Viewer':
      return 2;
    case 'PublicApi':
      return 3;
    default:
      return -1;
  }
};

const forFurther = (v: string): number => {
  switch (v) {
    case 'ApplicationAdmin':
      return 0;
    case 'DomainAdmin':
      return 1;
    case 'ProjectAnalyst':
      return 2;
    case 'PolicyManager':
      return 3;
    case 'LicenseManager':
      return 5;
    case 'Internal':
      return 6;
    case 'NonInternal':
      return 7;
    default:
      return -1;
  }
};

onMounted(() => {
  getProjectAccessRights();
});
</script>

<template>
  <v-container fluid>
    <v-card class="pa-4 mb-3" style="height: auto; border: none !important">
      <v-row>
        <v-col cols="12" xs="12">
          <h1 class="text-h5 pb-3">{{ t('PROJECT_ACCESS_RIGHTS') }}</h1>
          <v-row>
            <v-col cols="12" v-if="accessRightsLoaded">
              <v-data-table
                density="compact"
                class="striped-table"
                :headers="projectAccessRightsHeaders"
                fixed-header
                hide-default-footer
                :items="allProjectAccessRightsForTable"
                disable-sort
                :items-per-page="-1">
                <template v-slot:headers>
                  <tr>
                    <th>{{ t('ACCESS_RIGHTS') }}</th>
                    <th v-for="header in projectAccessRightsHeaders.slice(1)" :key="header.value">
                      {{ header.title }}
                    </th>
                  </tr>
                  <tr>
                    <th></th>
                    <th v-for="header in projectAccessRightsHeaders.slice(1)" :key="header.value">
                      <span class="tableCRUD">C</span>
                      <span class="borderCRUD px-1">|</span>
                      <span class="tableCRUD">R</span>
                      <span class="borderCRUD px-1">|</span>
                      <span class="tableCRUD">U</span>
                      <span class="borderCRUD px-1">|</span>
                      <span class="tableCRUD">D</span>
                    </th>
                  </tr>
                </template>

                <template v-slot:item="{item}">
                  <tr>
                    <td>{{ t('ACCESS_RIGHT_' + item.accessRightName) }}</td>
                    <td v-for="(rights, index) in item.entry" :key="index">
                      <v-tooltip
                        v-for="(v, k) in rights[1]"
                        :key="k"
                        :open-delay="TOOLTIP_OPEN_DELAY_IN_MS"
                        location="bottom">
                        <template v-slot:activator="{props}">
                          <span
                            v-if="k === 'upload' || k === 'download'"
                            v-bind="props"
                            class="tableCRUDIcon marginTableCRUD">
                            <v-icon v-if="v" color="primary">mdi-check</v-icon>
                            <v-icon v-else class="greyCheck">mdi-check</v-icon>
                            <span class="borderCRUD pl-6">|</span>
                          </span>
                          <span v-else v-bind="props" class="tableCRUDIcon">
                            <v-icon v-if="v" color="primary">mdi-check</v-icon>
                            <v-icon v-else class="greyCheck">mdi-check</v-icon>
                            <span class="borderCRUD px-2" v-if="k !== 'delete' && k !== 'download'">|</span>
                          </span>
                        </template>
                        <span>{{ k }}</span>
                      </v-tooltip>
                    </td>
                  </tr>
                  <tr v-if="item.accessRightName === 'allowRequestApproval'">
                    <th></th>
                    <th v-for="header in projectAccessRightsHeaders.length - 1">
                      <span>
                        Upload <span class="borderCRUD">|</span> Download <span class="borderCRUD">|</span> Delete
                      </span>
                    </th>
                  </tr>
                </template>
              </v-data-table>
            </v-col>
          </v-row>
        </v-col>
      </v-row>
    </v-card>

    <v-card class="pa-4 mb-3" style="height: auto; border: none !important">
      <v-row>
        <v-col cols="12" xs="12">
          <h1 class="text-h5 pb-3">{{ t('FURTHER_ACCESS_RIGHTS') }}</h1>
          <v-row style="height: auto">
            <v-col cols="12" md="4" xs="12" v-if="accessRightsLoaded">
              <v-select
                variant="outlined"
                density="compact"
                :items="allRoles"
                v-model="selectedRoles"
                v-bind:menu-props="{location: 'bottom'}"
                multiple
                :label="t('ROLES')"
                class="pb-2"
                item-title="text"
                hide-details="auto"
                @update:modelValue="selectedRolesChanged">
                <template v-slot:chip="{item, props}">
                  <DLabel :labelName="item.title" closable :parentProps="props" :iconName="icons.SECURITY" />
                </template>
              </v-select>
            </v-col>
            <v-col>
              <DCActionButton
                large
                :text="t('BTN_SHOW_ROLES')"
                :hint="t('TT_SHOW_ROLES')"
                @click="prepareDataForSelection"
                class="mx-2"
                :disabled="!dirtySelection" />
              <DCActionButton
                large
                :text="t('BTN_SHOW_DEFAULT_ROLES')"
                :hint="t('TT_SHOW_DEFAULT_ROLES')"
                @click="prepareDataForDefaultSelection"
                class="mx-2"
                :disabled="defaultSelection" />
            </v-col>
          </v-row>
          <v-row>
            <v-col cols="12" v-if="accessRightsLoaded">
              <v-data-table
                density="compact"
                class="striped-table"
                :headers="accessRightsHeaderForSelection"
                :items="accessRightsForTableForSelection"
                fixed-header
                hide-default-footer
                disable-sort
                :items-per-page="-1">
                <template v-slot:headers>
                  <tr>
                    <th>{{ t('ACCESS_RIGHTS') }}</th>
                    <th v-for="header in accessRightsHeaderForSelection.slice(1)" :key="header.value">
                      {{ header.title }}
                    </th>
                  </tr>
                  <tr>
                    <th></th>
                    <th v-for="header in accessRightsHeaderForSelection.slice(1)" :key="header.value">
                      <span class="tableCRUD">C</span>
                      <span class="borderCRUD px-1">|</span>
                      <span class="tableCRUD">R</span>
                      <span class="borderCRUD px-1">|</span>
                      <span class="tableCRUD">U</span>
                      <span class="borderCRUD px-1">|</span>
                      <span class="tableCRUD">D</span>
                    </th>
                  </tr>
                </template>

                <template v-slot:item="{item}">
                  <tr>
                    <td>{{ t('ACCESS_RIGHT_' + item.accessRightName) }}</td>
                    <td v-for="(rights, index) in item.entry" :key="index">
                      <v-tooltip
                        v-for="(v, k) in rights[1]"
                        :key="k"
                        :open-delay="TOOLTIP_OPEN_DELAY_IN_MS"
                        location="bottom">
                        <template v-slot:activator="{props}">
                          <span v-bind="props" class="tableCRUDIcon">
                            <v-icon v-if="v" color="primary">mdi-check</v-icon>
                            <v-icon v-else class="greyCheck">mdi-check</v-icon>
                            <span class="borderCRUD px-2" v-if="k !== 'delete'">|</span>
                          </span>
                        </template>
                        <span>{{ k }}</span>
                      </v-tooltip>
                    </td>
                  </tr>
                </template>
              </v-data-table>
            </v-col>
          </v-row>
        </v-col>
      </v-row>
    </v-card>
  </v-container>
</template>

<style scoped>
.greyCheck {
  color: rgb(var(--v-theme-tableBorderColor));
}

.marginTableCRUD {
  margin: 0 8px;
}

.tableCRUD {
  padding: 0 10px;
}

.tableCRUDIcon {
  padding: 0 2px;
}

.borderCRUD {
  color: rgb(var(--v-theme-tableBorderColor));
}
</style>
