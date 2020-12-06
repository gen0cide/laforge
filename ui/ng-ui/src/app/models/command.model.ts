import { ID, Tag, User, varsMap } from './common.model';

interface Command {
  id: ID;
  name: string;
  description: string;
  program: string;
  args: string[];
  ignoreErrors: boolean;
  cooldown: number;
  timeout: number;
  disabled: boolean;
  vars: varsMap[];
  tags: Tag[];
  maintainer: User;
}

export { Command };
