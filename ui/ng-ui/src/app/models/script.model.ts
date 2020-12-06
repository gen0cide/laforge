import { ID, Tag, User, varsMap } from './common.model';
import { Finding } from './finding.model';

interface Script {
  id: ID;
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
  tags: Tag[];
  absPath: string;
  maintainer: User;
  findings: Finding[];
}

export { Script };
