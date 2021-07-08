import { ID, varsMap, Status, tagMap, Team, Plan } from './common.model';
import { Build, Environment } from './environment.model';
import { hostChildrenCompleted, ProvisionedHost } from './host.model';

export interface Network {
  id: ID;
  vdi_visible: boolean;
  name?: string;
  cidr?: string;
  vars?: varsMap[];
  tags?: tagMap[];
  NetworkToEnvironment?: Environment;
}

export interface ProvisionedNetwork {
  id: ID;
  name: string;
  cidr: string;
  ProvisionedNetworkToStatus?: Status;
  ProvisionedNetworkToNetwork?: Network;
  ProvisionedNetworkToBuild?: Build;
  ProvisionedNetworkToTeam?: Team;
  ProvisionedNetworkToProvisionedHost?: ProvisionedHost[];
  ProvisionedNetworkToPlan?: Plan;
}

export const networkChildrenCompleted = (network: ProvisionedNetwork): boolean => {
  let numCompleted = 0;
  let totalHosts = 0;
  for (const host of network.ProvisionedNetworkToProvisionedHost) {
    totalHosts++;
    if (hostChildrenCompleted(host)) numCompleted++;
  }
  if (numCompleted === totalHosts) return true;
  else return false;
};
