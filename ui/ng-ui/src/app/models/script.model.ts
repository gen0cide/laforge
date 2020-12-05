import { ID, Tag, User, varsMap } from './common.model';
import { Finding } from './finding.model';

interface Script {
  id: ID;
  name: String;
  language: String;
  description: String;
  source: String;
  sourceType: String;
  cooldown: Number;
  timeout: Number;
  ignoreErrors: Boolean;
  args: String[];
  disabled: Boolean;
  vars: varsMap[];
  tags: Tag[];
  absPath: String;
  maintainer: User;
  findings: Finding[];
}

export {
  Script
}