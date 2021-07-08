import { ID, tagMap, varsMap } from './common.model';
import { Environment } from './environment.model';

export interface Role {
  id: ID;
  name: string;
  // permissions field if we want roles to act as "groups"
}

export interface User {
  id: ID;
  name: string;
  uuid: string;
  email: string;
}

export enum RoleLevel {
  ADMIN,
  USER,
  UNDEFINED
}

export enum ProviderType {
  LOCAL,
  GITHUB,
  OPENID,
  UNDEFINED
}
export interface AuthUser {
  id: ID;
  username: string;
  role: RoleLevel;
  provider: ProviderType;
}

export interface Identity {
  id: ID;
  first_name: string;
  last_name: string;
  email: string;
  password: string;
  hcl_id?: string;
  description?: string;
  avatar_file?: string;
  vars?: varsMap[];
  tags?: tagMap[];
  identityToEnvironment?: Environment;
}
