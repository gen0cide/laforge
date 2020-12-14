import { Injectable } from '@angular/core';
import { ApolloQueryResult } from '@apollo/client/core';
// import { ApolloQueryResult } from '@apollo/client/core';
import { Apollo, gql } from 'apollo-angular';
import { Environment, resolveStatuses } from 'src/app/models/environment.model';
// import { bradley } from 'src/data/sample-config';
import { getEnvironmentQuery } from './queries/environment';
import { Observable } from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  constructor(private apollo: Apollo) {}

  getEnvironment(id: string): Observable<ApolloQueryResult<unknown>> {
    // return new Promise((resolve) => resolve(bradley));

    // TODO: Resolve API CORS error

    return this.apollo.watchQuery({
      query: getEnvironmentQuery(id)
    }).valueChanges;
  }
}
