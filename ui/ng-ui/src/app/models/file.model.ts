import { ID, Tag, varsMap } from './common.model';

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

export interface RemoteFile {
  id: ID;
  sourceType: string;
  sourec: string;
  destination: string;
  vars?: varsMap[];
  template: boolean;
  perms: string;
  disabled: boolean;
  md5: string;
  absPath: string;
  ext: string;
  tags?: Tag[];
}
