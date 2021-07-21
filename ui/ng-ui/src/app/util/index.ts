import {
  LaForgeProvisionedHost,
  LaForgeProvisionedNetwork,
  LaForgeProvisionStatus,
  LaForgeSubscribeUpdatedStatusSubscription,
  LaForgeTeam
} from '@graphql';

export const hostChildrenCompleted = (
  host: LaForgeProvisionedHost,
  getStatus: (statusId: string) => LaForgeSubscribeUpdatedStatusSubscription['updatedStatus']
): boolean => {
  let numCompleted = 0;
  let totalSteps = 0;
  for (const step of host.ProvisionedHostToProvisioningStep) {
    totalSteps++;
    const stepStatus = getStatus(step.id);
    if (!stepStatus) numCompleted++;
    if (stepStatus.state === LaForgeProvisionStatus.Complete) numCompleted++;
    // if (step.ProvisioningStepToStatus.state === LaForgeProvisionStatus.Complete) numCompleted++;
  }
  if (numCompleted === totalSteps) return true;
  else return false;
};

export const networkChildrenCompleted = (
  network: LaForgeProvisionedNetwork,
  getStatus: (statusId: string) => LaForgeSubscribeUpdatedStatusSubscription['updatedStatus']
): boolean => {
  return false;
  // let numCompleted = 0;
  // let totalHosts = 0;
  // for (const host of network.ProvisionedNetworkToProvisionedHost) {
  //   totalHosts++;
  //   if (hostChildrenCompleted(host, getStatus)) numCompleted++;
  // }
  // if (numCompleted === totalHosts) return true;
  // else return false;
};

export const teamChildrenCompleted = (
  team: LaForgeTeam,
  getStatus: (statusId: string) => LaForgeSubscribeUpdatedStatusSubscription['updatedStatus']
): boolean => {
  return false;
  // let numCompleted = 0;
  // let totalNetworks = 0;
  // for (const network of team.TeamToProvisionedNetwork) {
  //   totalNetworks++;
  //   if (networkChildrenCompleted(network, getStatus)) numCompleted++;
  // }
  // if (numCompleted === totalNetworks) return true;
  // else return false;
};
