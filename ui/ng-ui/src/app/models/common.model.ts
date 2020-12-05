import { Build, Environment } from './environment.model'
import { ProvisionedNetwork } from './network.model'

type ID = String

interface varsMap {
  key?: String;
  value?: String;
}

interface configMap {
  key?: String;
  value?: String;
}

interface User {
  id: ID;
  name: String;
  uuid: String;
  email: String;
}

interface Tag {
  id: ID;
  name: String;
  description?: String;
}

interface Team {
  id: ID;
  teamNumber: Number;
  config: configMap[];
  revision: Number;
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
  startedAt: String;
  endedAt: String;
  failed: Boolean;
  completed: Boolean;
  error: String;
}

export {
  ID,
  varsMap,
  configMap,
  User,
  Tag,
  Team,
  Status
}