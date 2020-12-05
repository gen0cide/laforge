import { ID, Tag } from './common.model';

interface FileDownload {
  id: ID;
  sourceType: String;
  source: String;
  destination: String;
  template: Boolean;
  mode: String;
  disabled: Boolean;
  md5: String;
  absPath: String;
  tags: Tag[];
}

interface FileDelete {
  id: ID;
  path: String;
}

interface FileExtract {
  id: ID;
  source: String;
  destination: String;
  type: String;
}

export {
  FileDownload,
  FileDelete,
  FileExtract
}