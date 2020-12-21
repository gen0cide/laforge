import { ID, Tag } from './common.model';

export interface FileDownload {
  id: ID;
  sourceType: string;
  source: string;
  destination: string;
  template: boolean;
  mode: string;
  disabled: boolean;
  md5: string;
  absPath: string;
  tags: Tag[];
}

export interface FileDelete {
  id: ID;
  path: string;
}

export interface FileExtract {
  id: ID;
  source: string;
  destination: string;
  type: string;
}
