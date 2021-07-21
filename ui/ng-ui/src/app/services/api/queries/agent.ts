import { gql } from 'apollo-angular';

const GetAgentStatusesQuery = gql`
  query($buildId: String!) {
    build(buildUUID: $buildId) {
      id
      buildToTeam {
        id
        TeamToProvisionedNetwork {
          id
          ProvisionedNetworkToProvisionedHost {
            id
            ProvisionedHostToAgentStatus {
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
              usedMem
              timestamp
            }
          }
        }
      }
    }
  }
`;

// export { GetAgentStatusesQuery };
