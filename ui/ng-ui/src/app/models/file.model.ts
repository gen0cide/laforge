import { ID, Tag } from './common.model';

interface FileDownload {
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

interface FileDelete {
  id: ID;
  path: string;
}

interface FileExtract {
  id: ID;
  source: string;
  destination: string;
  type: string;
}

export { FileDownload, FileDelete, FileExtract };
