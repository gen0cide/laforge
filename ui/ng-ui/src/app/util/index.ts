import { LaForgeProvisionedHost, LaForgeProvisionedNetwork, LaForgeProvisionStatus, LaForgeTeam } from '@graphql';

export const hostChildrenCompleted = (host: LaForgeProvisionedHost): boolean => {
  let numCompleted = 0;
  let totalSteps = 0;
  for (const step of host.ProvisionedHostToProvisioningStep) {
    totalSteps++;
    if (step.ProvisioningStepToStatus.state === LaForgeProvisionStatus.Complete) numCompleted++;
  }
  if (numCompleted === totalSteps) return true;
  else return false;
};

export const networkChildrenCompleted = (network: LaForgeProvisionedNetwork): boolean => {
  let numCompleted = 0;
  let totalHosts = 0;
  for (const host of network.ProvisionedNetworkToProvisionedHost) {
    totalHosts++;
    if (hostChildrenCompleted(host)) numCompleted++;
  }
  if (numCompleted === totalHosts) return true;
  else return false;
};

export const teamChildrenCompleted = (team: LaForgeTeam): boolean => {
  let numCompleted = 0;
  let totalNetworks = 0;
  for (const network of team.TeamToProvisionedNetwork) {
    totalNetworks++;
    if (networkChildrenCompleted(network)) numCompleted++;
  }
  if (numCompleted === totalNetworks) return true;
  else return false;
};
