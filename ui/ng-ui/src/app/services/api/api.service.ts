import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { ApolloQueryResult, FetchResult } from '@apollo/client/core';
import { environment } from '@env';
import { AuthUser } from '@models/user.model';
import { Apollo, QueryRef } from 'apollo-angular';

import { EmptyObject } from 'apollo-angular/types';
import { Observable } from 'rxjs';

import {
  AgentStatusQueryResult,
  BuildQueryResult,
  BuildQueryVars,
  EnvironmentInfo,
  EnvironmentInfoQueryResult,
  EnvironmentQueryResult,
  HostStepsQueryResult
} from 'src/app/models/api.model';

import { ID } from 'src/app/models/common.model';
import { Environment, Build } from 'src/app/models/environment.model';
import { ProvisioningStep } from 'src/app/models/step.model';

import { GetAgentStatusesQuery } from './queries/agent';
import { GetBuildPlansData, GetBuildPlansQuery, GetBuildPlansVars, GetBuildQuery } from './queries/build';
import { getEnvConfigQuery } from './queries/env-tree';

import { getEnvironmentQuery, getEnvironmentsQuery } from './queries/environment';
import { getProvisionedSteps } from './queries/steps';

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  private statusPollingInterval: number;

  constructor(private apollo: Apollo, private http: HttpClient) {}

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
   * Pulls a build from the API once, without exposing a subscription or observable
   * @param id The ID of the build
   */
  public async pullBuild(id: ID): Promise<Build> {
    const res = await this.apollo
      .query<BuildQueryResult, BuildQueryVars>({
        query: GetBuildQuery,
        variables: {
          buildId: id
        }
      })
      .toPromise();
    return res.data.build;
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
  public async pullHostSteps(hostId: ID): Promise<ProvisioningStep[]> {
    const res = await this.apollo
      .query<HostStepsQueryResult>({
        query: getProvisionedSteps,
        variables: {
          hostId
        }
      })
      .toPromise();
    return res.data.provisionedHost.ProvisionedHostToProvisioningStep;
  }

  /**
   * Pulls the statuses of all running agents from the API once, without exposing a subscription or observable
   * @param buildId The ID of the build
   */
  public getAgentStatuses(buildId: ID): QueryRef<AgentStatusQueryResult, EmptyObject> {
    return this.apollo.watchQuery<AgentStatusQueryResult>({
      query: GetAgentStatusesQuery,
      variables: {
        buildId
      }
      // pollInterval: 0,
      // fetchPolicy: 'cache-and-network'
    });
  }

  /**
   * Pulls the statuses of all running agents from the API once, without exposing a subscription or observable
   * @param envId The Environment ID of the environment
   */
  public async pullAgentStatuses(buildId: ID): Promise<AgentStatusQueryResult> {
    if (environment.isMockApi) {
      return new Promise((resolve, reject) => {
        resolve((environment.mockAgentStatuses as any).data);
      });
    }
    const res = await this.apollo
      .query<AgentStatusQueryResult>({
        query: GetAgentStatusesQuery,
        variables: {
          buildId
        }
      })
      .toPromise()
      .then((result: ApolloQueryResult<AgentStatusQueryResult>) => result.data);
    return res;
  }

  public async pullBuildPlans(buildId: ID): Promise<GetBuildPlansData> {
    const res = await this.apollo
      .query<GetBuildPlansData, GetBuildPlansVars>({
        query: GetBuildPlansQuery,
        variables: {
          buildId
        }
      })
      .toPromise()
      .then((result) => result.data);
    return res;
  }

  public async fakeLogin(): Promise<AuthUser> {
    return new Promise((resolve, reject) => {
      this.http
        .get<AuthUser>(window.location.protocol + '//' + window.location.hostname + '/api/local/login', {
          withCredentials: true,
          observe: 'response'
        })
        .subscribe((resp) => {
          if (resp.status === 200) {
            resolve(resp.body);
          } else {
            reject(null);
          }
        });
    });
  }
}
