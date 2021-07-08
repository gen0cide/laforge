import { Injectable } from '@angular/core';
import { RebuildPlansData, RebuildPlansMutation, RebuildPlansVars } from '@services/api/queries/rebuild';
import { Apollo } from 'apollo-angular';
import { GraphQLError } from 'graphql';
import { ID, Team } from 'src/app/models/common.model';
import { ProvisionedHost } from 'src/app/models/host.model';
import { ProvisionedNetwork } from 'src/app/models/network.model';

@Injectable({
  providedIn: 'root'
})
export class RebuildService {
  rootPlans: ID[];

  constructor(private apollo: Apollo) {
    this.rootPlans = [];
  }

  /**
   * Executes the current rebuild based on the selected plans
   * @returns promise if the execution was a success (promise rejects with query errors)
   */
  executeRebuild = (): Promise<boolean> => {
    return new Promise<boolean>((resolve: (value: boolean) => void, reject: (reason?: readonly GraphQLError[]) => void) => {
      this.apollo
        .mutate<RebuildPlansData, RebuildPlansVars>({
          mutation: RebuildPlansMutation,
          variables: {
            rootPlans: this.rootPlans
          }
        })
        .subscribe((res) => {
          if (res.data?.rebuild) {
            resolve(true);
            // Clear the list of selections on success
            this.rootPlans = [];
          } else if (res.errors) {
            reject(res.errors);
          } else {
            resolve(false);
          }
        });
    });
  };

  addPlan = (planId: ID): void => {
    console.log(`add plan: ${planId}`);
    if (this.rootPlans.filter((id: ID) => id === planId).length === 0) this.rootPlans.push(planId);
  };

  removePlan = (planId: ID): void => {
    console.log(`rem plan: ${planId}`);
    if (this.rootPlans.indexOf(planId) >= 0) this.rootPlans.splice(this.rootPlans.indexOf(planId), 1);
  };

  /**
   * addTeam selects teams to rebuild
   * @param team team to rebuild
   * @returns successfully added to list to rebuild
   */
  addTeam = (team: Team): boolean => {
    const planId = team.TeamToPlan?.id ?? null;
    if (planId === null) return false;
    this.addPlan(planId);
    return true;
  };

  /**
   * removeTeam removes a team from the rebuild list
   * @param team team to remove from rebuild list
   * @returns successfully removed from the list to rebuild
   */
  removeTeam = (team: Team): boolean => {
    const planId = team.TeamToPlan?.id ?? null;
    if (planId === null) return false;
    this.removePlan(planId);
    return true;
  };

  /**
   * hasTeam checks if a team is in the list to rebuild
   * @param team team to check
   * @returns if team is in rebuild list
   */
  hasTeam = (team: Team): boolean => {
    const planId = team.TeamToPlan?.id ?? null;
    if (planId === null) return false;
    return this.rootPlans.indexOf(planId) >= 0;
  };

  /**
   * addNetwork selects networks to rebuild
   * @param network network to rebuild
   * @returns successfully added to list to rebuild
   */
  addNetwork = (network: ProvisionedNetwork): boolean => {
    const planId = network.ProvisionedNetworkToPlan?.id ?? null;
    if (planId === null) return false;
    this.addPlan(planId);
    return true;
  };

  /**
   * removeNetwork removes a network from the rebuild list
   * @param network network to remove from rebuild list
   * @returns successfully removed from the list to rebuild
   */
  removeNetwork = (network: ProvisionedNetwork): boolean => {
    const planId = network.ProvisionedNetworkToPlan?.id ?? null;
    if (planId === null) return false;
    this.removePlan(planId);
    return true;
  };

  /**
   * hasNetwork checks if a network is in the list to rebuild
   * @param network network to check
   * @returns if network is in rebuild list
   */
  hasNetwork = (network: ProvisionedNetwork): boolean => {
    const planId = network.ProvisionedNetworkToPlan?.id ?? null;
    if (planId === null) return false;
    return this.rootPlans.indexOf(planId) >= 0;
  };

  /**
   * addHost selects hosts to rebuild
   * @param host host to rebuild
   * @returns successfully added to list to rebuild
   */
  addHost = (host: ProvisionedHost): boolean => {
    const planId = host.ProvisionedHostToPlan?.id ?? null;
    if (planId === null) return false;
    this.addPlan(planId);
    return true;
  };

  /**
   * removeHost removes a host from the rebuild list
   * @param host host to remove from rebuild list
   * @returns successfully removed from the list to rebuild
   */
  removeHost = (host: ProvisionedHost): boolean => {
    const planId = host.ProvisionedHostToPlan?.id ?? null;
    if (planId === null) return false;
    this.removePlan(planId);
    return true;
  };

  /**
   * hasHost checks if a host is in the list to rebuild
   * @param host host to check
   * @returns if host is in rebuild list
   */
  hasHost = (host: ProvisionedHost): boolean => {
    const planId = host.ProvisionedHostToPlan?.id ?? null;
    if (planId === null) return false;
    return this.rootPlans.indexOf(planId) >= 0;
  };
}
