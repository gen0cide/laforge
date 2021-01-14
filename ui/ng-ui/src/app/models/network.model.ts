import { ID, varsMap, Tag, Status } from './common.model';
import { ProvisionedHost } from './host.model';

export interface Network {
  id: ID;
  name?: string;
  cidr?: string;
  vdiVisible: boolean;
  vars?: varsMap[];
  tags?: Tag[];
}

export interface ProvisionedNetwork {
  id: ID;
  name: string;
  cidr: string;
  // vars: varsMap[];
  // tags: Tag[];
  provisionedHosts: ProvisionedHost[];
  status?: Status;
  network: Network;
  // build: Build; Circular dependency
}
