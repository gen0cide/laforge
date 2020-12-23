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

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  constructor(private apollo: Apollo) {}

  /**
   * Sets up a subscription with the API to return an observable that updates as teh values change in the database
   * @param id The Environment ID of the environment
   */
  public getEnvironment(id: string): Observable<ApolloQueryResult<EnvironmentQueryResult>> {
    return this.apollo.watchQuery<EnvironmentQueryResult>({
      query: getEnvironmentQuery(id)
    }).valueChanges;
  }

  /**
   * Pulls an environment from the API once, without exposing a subscription or observable
   * @param id The Environment ID of the environment
   */
  public async pullEnvironments(): Promise<EnvironmentInfo[]> {
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
  public async pullEnvironment(id: string): Promise<Environment> {
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
  public async pullEnvTree(id: string): Promise<Environment> {
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
  public async pullHostSteps(hostId: string): Promise<ProvisionedStep[]> {
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
  public getAgentStatuses(envId: string): QueryRef<AgentStatusQueryResult, EmptyObject> {
    return this.apollo.watchQuery<AgentStatusQueryResult>({
      query: getAgentStatusesQuery,
      variables: {
        id: `${envId}`
      },
      // pollInterval: 0,
      // fetchPolicy: 'cache-and-network'
    });
  }

  /**
   * Pulls the statuses of all running agents from the API once, without exposing a subscription or observable
   * @param envId The Environment ID of the environment
   */
  public async pullAgentStatuses(envId: string): Promise<AgentStatusQueryResult> {
    console.log('test');
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
