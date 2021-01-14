import { Injectable } from '@angular/core';
import { ApolloQueryResult, FetchResult } from '@apollo/client/core';
import { Apollo, QueryRef } from 'apollo-angular';
import { getEnvironmentQuery, getEnvironmentsQuery } from './queries/environment';
import { Observable } from 'rxjs';
import { getAgentStatusesQuery } from './queries/agent';
import {
  AgentStatusQueryResult,
  EnvironmentInfo,
  EnvironmentInfoQueryResult,
  EnvironmentQueryResult,
  HostStepsQueryResult
} from 'src/app/models/api.model';
import { Environment } from 'src/app/models/environment.model';
import { getEnvConfigQuery } from './queries/env-tree';
import { getProvisionedSteps } from './queries/steps';
import { ProvisionedStep } from 'src/app/models/host.model';
import { EmptyObject } from 'apollo-angular/types';
import { environment } from 'src/environments/environment';
import { ID } from 'src/app/models/common.model';

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  private statusPollingInterval: number;

  constructor(private apollo: Apollo) {}

  /**
   * Set the interval used for agent status polling
   * @param interval the polling interval in milliseconds
   */
  public setStatusPollingInterval(interval: number): void {
    this.statusPollingInterval = interval;
  }
  /**
   * Get the interval used for agent status polling
   */
  public getStatusPollingInterval(): number {
    return this.statusPollingInterval;
  }

  /**
   * Sets up a subscription with the API to return an observable that updates as teh values change in the database
   * @param id The Environment ID of the environment
   */
  public getEnvironment(id: ID): Observable<ApolloQueryResult<EnvironmentQueryResult>> {
    return this.apollo.watchQuery<EnvironmentQueryResult>({
      query: getEnvironmentQuery(id)
    }).valueChanges;
  }

  /**
   * Pulls an environment from the API once, without exposing a subscription or observable
   * @param id The Environment ID of the environment
   */
  public async pullEnvironments(): Promise<EnvironmentInfo[]> {
    if (environment.isMockApi) {
      return new Promise<EnvironmentInfo[]>((resolve, reject) => {
        resolve((environment.mockEnvList as any).data.environments);
      });
    }
    const res = await this.apollo
      .query<EnvironmentInfoQueryResult>({
        query: getEnvironmentsQuery()
      })
      .toPromise();
    return res.data.environments;
  }

  /**
   * Pulls an environment from the API once, without exposing a subscription or observable
   * @param id The Environment ID of the environment
   */
  public async pullEnvironment(id: ID): Promise<Environment> {
    const res = await this.apollo
      .query<EnvironmentQueryResult>({
        query: getEnvironmentQuery(id)
      })
      .toPromise();
    return res.data.environment;
  }

  /**
   * Pulls an environment tree from the API once, without exposing a subscription or observable
   * @param id The Environment ID of the environment
   */
  public async pullEnvTree(id: ID): Promise<Environment> {
    if (environment.isMockApi) {
      return new Promise((resolve, reject) => {
        // const envTree = import('../../../data/CPTC_Env_Tree.json').then(
        //   (tree) => resolve((tree.default.data as EnvironmentQueryResult).environment),
        //   (err) => reject(err)
        // );
        resolve((environment.mockEnvTree as any).data.environment);
      });
    }
    const res = await this.apollo
      .query<EnvironmentQueryResult>({
        query: getEnvConfigQuery(id)
      })
      .toPromise();
    return res.data.environment;
  }

  /**
   * Pulls an environment tree from the API once, without exposing a subscription or observable
   * @param id The Environment ID of the environment
   */
  public async pullHostSteps(hostId: ID): Promise<ProvisionedStep[]> {
    const res = await this.apollo
      .query<HostStepsQueryResult>({
        query: getProvisionedSteps(hostId)
      })
      .toPromise();
    return res.data.provisionedHost.provisionedSteps;
  }

  /**
   * Pulls the statuses of all running agents from the API once, without exposing a subscription or observable
   * @param envId The Environment ID of the environment
   */
  public getAgentStatuses(envId: ID): QueryRef<AgentStatusQueryResult, EmptyObject> {
    return this.apollo.watchQuery<AgentStatusQueryResult>({
      query: getAgentStatusesQuery,
      variables: {
        id: `${envId}`
      }
      // pollInterval: 0,
      // fetchPolicy: 'cache-and-network'
    });
  }

  /**
   * Pulls the statuses of all running agents from the API once, without exposing a subscription or observable
   * @param envId The Environment ID of the environment
   */
  public async pullAgentStatuses(envId: ID): Promise<AgentStatusQueryResult> {
    if (environment.isMockApi) {
      return new Promise((resolve, reject) => {
        resolve((environment.mockAgentStatuses as any).data);
      });
    }
    const res = await this.apollo
      .query<AgentStatusQueryResult>({
        query: getAgentStatusesQuery,
        variables: {
          id: envId
        }
      })
      .toPromise()
      .then((result: ApolloQueryResult<AgentStatusQueryResult>) => result.data);
    console.log('test 2');
    return res;
  }
}
