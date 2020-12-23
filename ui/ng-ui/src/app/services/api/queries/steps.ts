import { gql } from 'apollo-angular';
import { DocumentNode } from 'graphql';

const getProvisionedSteps = (hostId: string): DocumentNode => gql`
  {
    provisionedHost(proHostUUID: "${hostId}") {
      id
      provisionedSteps {
        id
        provisionType
        status {
          state
          startedAt
          endedAt
          error
        }
        script {
          id
          name
          description
          source
          sourceType
          disabled
        }
        command {
          id
          name
          description
          args
          disabled
        }
        DNSRecord {
          id
          name
          values
          type
          zone
          disabled
        }
      }
    }
  }
`;

export { getProvisionedSteps };
