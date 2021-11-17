import { NgModule } from '@angular/core';
import { ApolloClientOptions, InMemoryCache, split } from '@apollo/client/core';
import { WebSocketLink } from '@apollo/client/link/ws';
import { getMainDefinition } from '@apollo/client/utilities';
import { APOLLO_OPTIONS } from 'apollo-angular';
import { HttpLink } from 'apollo-angular/http';
import { environment } from 'src/environments/environment';

export function createApollo(httpLink: HttpLink): ApolloClientOptions<any> {
  const httpClient = httpLink.create({
    uri: environment.graphqlUrl,
    withCredentials: true
  });
  const wsClient = new WebSocketLink({
    uri: environment.wsUrl,
    options: {
      reconnect: true,
      timeout: 30000,
      minTimeout: 100000
    }
  });

  const link = split(
    ({ query }) => {
      const { kind, operation } = getMainDefinition(query) as any;
      return kind === 'OperationDefinition' && operation === 'subscription';
    },
    wsClient,
    httpClient
  );

  return {
    uri: environment.graphqlUrl,
    link,
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
