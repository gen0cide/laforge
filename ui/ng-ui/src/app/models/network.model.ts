import { ID, varsMap, Tag, Status } from './common.model';
import { Build } from './environment.model';
import { ProvisionedHost } from './host.model';

interface Network {
  id: ID;
  name: string;
  cidr: string;
  vdiVisible: boolean;
  vars: varsMap[];
  tags: Tag[];
}

interface ProvisionedNetwork {
  id: ID;
  name: string;
  cidr: string;
  vars: varsMap[];
  tags: Tag[];
  provisionedHosts: ProvisionedHost[];
  status: Status;
  network: Network;
  // build: Build; Circular dependency
}

export { Network, ProvisionedNetwork };
