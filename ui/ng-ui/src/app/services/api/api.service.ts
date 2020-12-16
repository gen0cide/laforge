import { Injectable } from '@angular/core';
import { ApolloQueryResult } from '@apollo/client/core';
import { Apollo, QueryRef } from 'apollo-angular';
import { getEnvironmentQuery } from './queries/environment';
import { Observable } from 'rxjs';
import { getAgentStatusesQuery } from './queries/agent';
import { EmptyObject } from 'apollo-angular/types';
import { AgentStatusQueryResult } from 'src/app/models/common.model';

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  constructor(private apollo: Apollo) {}

  getEnvironment(id: string): Observable<ApolloQueryResult<unknown>> {
    // return new Promise((resolve) => resolve(bradley));

    return this.apollo.watchQuery({
      query: getEnvironmentQuery(id)
    }).valueChanges;
  }

  // pullEnvironment(id: string): Environment {
  //   // return new Promise((resolve) => resolve(bradley));

  //   return this.apollo.watchQuery({
  //     query: getEnvironmentQuery(id)
  //   }).valueChanges;
  // }

  getAgentStatuses(envId: string): Promise<AgentStatusQueryResult> {
    return this.apollo
      .query<AgentStatusQueryResult>({
        query: getAgentStatusesQuery(envId)
      })
      .toPromise()
      .then((res: ApolloQueryResult<AgentStatusQueryResult>) => res.data);
    // return this.apollo.watchQuery<any>({
    //   query: getAgentStatusesQuery(envId)
    // });
  }
}
