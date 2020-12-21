import { Build, Environment } from './environment.model';
import { ProvisionedNetwork } from './network.model';

export type ID = string;

export interface varsMap {
  key?: string;
  value?: string;
}

export interface configMap {
  key?: string;
  value?: string;
}

export interface User {
  id: ID;
  name: string;
  uuid: string;
  email: string;
}

export interface Tag {
  id: ID;
  name: string;
  description?: string;
}

export interface Team {
  id: ID;
  teamNumber: number;
  config: configMap[];
  revision: number;
  maintainer: User;
  build: Build;
  environment: Environment;
  tags: Tag[];
  provisionedNetworks: ProvisionedNetwork[];
}

export enum ProvisionStatus {
  ProvStatusUndefined,
  ProvStatusAwaiting,
  ProvStatusInProgress,
  ProvStatusFailed,
  ProvStatusComplete,
  ProvStatusTainted
}

export interface Status {
  state: ProvisionStatus;
  startedAt: string;
  endedAt: string;
  failed: boolean;
  completed: boolean;
  error: string;
}
