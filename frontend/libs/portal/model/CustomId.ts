import {IMap} from '@disclosure-portal/utils/View';

export class CustomId {
  _key = '';
  updated?: string;
  created?: string;
  name = '';
  nameDE = '';
  description = '';
  descriptionDE = '';
  linkTemplate = '';
  technicalId = '';
  value = '';
}

export class CustomIds {
  public ids: CustomId[] = [] as CustomId[];
  public map: IMap<CustomId> = {} as IMap<CustomId>;

  constructor(customIds: CustomId[]) {
    this.ids = customIds;
    for (const id of this.ids) {
      this.map[id._key!] = id;
    }
  }
}

export interface CustomIdUsage {
  count: number;
}

export interface CustomIdName {
  name: string;
  nameDE: string;
  id: string;
}
