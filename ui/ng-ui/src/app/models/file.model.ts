import { ID, tagMap } from './common.model';
import { Environment } from './environment.model';

export interface FileDownload {
  id: ID;
  source: string;
  sourceType: string;
  destination: string;
  disabled: boolean;
  tags: tagMap[];
  hcl_id?: string;
  template?: boolean;
  perms?: string;
  md5?: string;
  absPath?: string;
  FileDownloadToEnvironment?: Environment;
}

export interface FileDelete {
  id: ID;
  hcl_id?: string;
  path: string;
  tags?: tagMap[];
  FileDeleteToEnvironment?: Environment;
}

export interface FileExtract {
  id: ID;
  source: string;
  destination: string;
  type: string;
  tags: tagMap[];
  hcl_id?: string;
  FileExtractToEnvironment?: Environment;
}
