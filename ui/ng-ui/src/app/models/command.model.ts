import { ID, Tag, User, varsMap } from './common.model';

interface Command {
  id: ID;
  name: String;
  description: String;
  program: String;
  args: String[];
  ignoreErrors: Boolean;
  cooldown: Number;
  timeout: Number;
  disabled: Boolean;
  vars: varsMap[];
  tags: Tag[];
  maintainer: User;
}

export {
  Command
}