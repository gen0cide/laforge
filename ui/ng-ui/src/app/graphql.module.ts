import { NgModule } from '@angular/core';
import { APOLLO_OPTIONS } from 'apollo-angular';
import { ApolloClientOptions, InMemoryCache } from '@apollo/client/core';
import { HttpLink } from 'apollo-angular/http';
import { environment } from 'src/environments/environment';

export function createApollo(httpLink: HttpLink): ApolloClientOptions<any> {
  return {
    uri: environment.graphqlUrl,
    link: httpLink.create({ uri: environment.graphqlUrl, withCredentials: true }),
    cache: new InMemoryCache({
      resultCaching: false
    }),
    credentials: 'include'
  };
}

@NgModule({
  providers: [
    {
      provide: APOLLO_OPTIONS,
      useFactory: createApollo,
      deps: [HttpLink]
    }
  ]
})
export class GraphQLModule {}
