export interface RemarkWithVersion {
  key?: string;
  version?: string;
  level?: string;
  status?: string; // Use ReviewRemarkStatus if available in shared types
  components?: {componentName: string; componentVersion: string}[] | string | null;
  sbomName?: string;
  licenses?: {licenseName: string; licenseId: string}[] | string | null;
  title?: string;
  origin?: string;
  author?: string;
  closed?: string | null;
  created?: string;
  updated?: string;
  description?: string;
  events?: Array<{
    author: string;
    content: string;
    created: string;
  }>;
}

export interface Remark {
  author: string;
}

export interface RRCommentExternDTO {
  content: string;
}

export enum ReviewRemarkLevel {
  GREEN = 'GREEN',
  YELLOW = 'YELLOW',
  RED = 'RED',
  NOT_SET = 'NOT_SET',
}

/**
 * Compare function for review remark levels to enable proper sorting
 * Order: GREEN (0) < YELLOW (1) < RED (2)
 * For descending sort (most severe first): RED > YELLOW > GREEN
 */
export function compareRRLevel(a: string, b: string): number {
  const levelWeight: Record<string, number> = {
    green: 0,
    yellow: 1,
    red: 2,
  };
  const aWeight = levelWeight[a?.toLowerCase()] ?? -1;
  const bWeight = levelWeight[b?.toLowerCase()] ?? -1;
  return aWeight - bWeight;
}
