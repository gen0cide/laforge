import { gql } from 'apollo-angular';
import { DocumentNode } from 'graphql';

const getAgentStatusesQuery = gql`
  query($id: String!) {
    environment(envUUID: $id) {
      build {
        teams {
          id
          provisionedNetworks {
            id
            provisionedHosts {
              id
              heartbeat {
                clientId
                hostname
                upTime
                bootTime
                numProcs
                OS
                hostID
                load1
                load5
                load15
                totalMem
                freeMem
                freeMem
                usedMem
              }
            }
          }
        }
      }
    }
  }
`;

export { getAgentStatusesQuery };
