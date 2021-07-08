import { ID, tagMap, varsMap } from './common.model';
import { Environment } from './environment.model';

export interface Command {
  id: ID;
  name: string;
  description: string;
  program: string;
  args: string[];
  disabled: boolean;
  hcl_id?: string;
  ignoreErrors?: boolean;
  cooldown?: number;
  timeout?: number;
  vars?: varsMap[];
  tags?: tagMap[];
  CommandToEnvironment?: Environment;
}
