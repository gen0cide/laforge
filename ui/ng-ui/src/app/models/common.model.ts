import { Build, Environment } from './environment.model';
import { ProvisionedNetwork } from './network.model';

type ID = string;

interface varsMap {
  key?: string;
  value?: string;
}

interface configMap {
  key?: string;
  value?: string;
}

interface User {
  id: ID;
  name: string;
  uuid: string;
  email: string;
}

interface Tag {
  id: ID;
  name: string;
  description?: string;
}

interface Team {
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

enum ProvisionStatus {
  ProvStatusUndefined,
  ProvStatusAwaiting,
  ProvStatusInProgress,
  ProvStatusFailed,
  ProvStatusComplete,
  ProvStatusTainted
}

interface Status {
  state: ProvisionStatus;
  startedAt: string;
  endedAt: string;
  failed: boolean;
  completed: boolean;
  error: string;
}

export { ID, varsMap, configMap, User, Tag, Team, Status, ProvisionStatus };
