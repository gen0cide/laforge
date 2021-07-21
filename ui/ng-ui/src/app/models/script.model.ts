import { ID, tagMap, varsMap } from './common.model';
import { Environment } from './environment.model';
import { Finding } from './finding.model';

export interface Script {
  id: ID;
  name: string;
  hcl_id?: string;
  language: string;
  description: string;
  source: string;
  source_type: string;
  disabled: boolean;
  args: string[];
  vars: varsMap[];
  tags: tagMap[];
  cooldown?: number;
  timeout?: number;
  ignore_errors?: boolean;
  absPath?: string;
  ScriptToFinding?: Finding[];
  ScriptToEnvironment?: Environment;
}
