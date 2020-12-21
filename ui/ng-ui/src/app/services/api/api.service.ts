import { Injectable } from '@angular/core';
import { ApolloQueryResult } from '@apollo/client/core';
import { Apollo } from 'apollo-angular';
import { getEnvironmentQuery } from './queries/environment';
import { Observable } from 'rxjs';
import { getAgentStatusesQuery } from './queries/agent';
import { AgentStatusQueryResult, EnvironmentQueryResult } from 'src/app/models/api.model';
import { Environment } from 'src/app/models/environment.model';

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
  public async pullEnvironment(id: string): Promise<Environment> {
    const res = await this.apollo
      .query<EnvironmentQueryResult>({
        query: getEnvironmentQuery(id)
      })
      .toPromise();
    return res.data.environment;
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
