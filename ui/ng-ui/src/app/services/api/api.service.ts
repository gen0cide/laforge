import { Injectable } from '@angular/core';
import { ApolloQueryResult } from '@apollo/client/core';
import { Apollo } from 'apollo-angular';
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
  public async getAgentStatuses(envId: string): Promise<AgentStatusQueryResult> {
    const res = await this.apollo
      .query<AgentStatusQueryResult>({
        query: getAgentStatusesQuery(envId)
      })
      .toPromise();
    return res.data;
  }
}
