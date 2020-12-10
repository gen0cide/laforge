import { Injectable } from '@angular/core';
// import { ApolloQueryResult } from '@apollo/client/core';
import { Apollo, gql } from 'apollo-angular';
import { Environment, resolveStatuses } from 'src/app/models/environment.model';
import { bradley } from 'src/data/sample-config';
import { getEnvironmentQuery } from './queries/environment';

@Injectable({
  providedIn: 'root'
})
export class ApiService {
  constructor(private apollo: Apollo) {}

  getEnvironment(id: string): Promise<Environment> {
    // return new Promise((resolve) => resolve(bradley));

    // TODO: Resolve API CORS error

    return new Promise((resolve, reject) => {
      this.apollo
        .query({
          query: getEnvironmentQuery(id)
        })
        .toPromise()
        .then((result) => resolve(resolveStatuses((result as any).data.environment) as Environment), reject);
    });
  }
}
