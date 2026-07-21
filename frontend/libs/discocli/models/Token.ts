export interface Token {
  token: string;
  projectUuid: string;
  type: 'project' | 'group';
  projectName?: string;
}
