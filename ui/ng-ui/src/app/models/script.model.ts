import { ID, Tag, User, varsMap } from './common.model';
import { Finding } from './finding.model';

export interface Script {
  id: ID;
  hcl_id?: string;
  name: string;
  language: string;
  description: string;
  source: string;
  sourceType: string;
  cooldown: number;
  timeout: number;
  ignoreErrors: boolean;
  args: string[];
  disabled: boolean;
  vars: varsMap[];
  scriptToTag: Tag[];
  absPath: string;
  scriptTouser: User;
  scriptToFinding: Finding[];
}
