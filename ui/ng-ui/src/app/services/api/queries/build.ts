import { ID } from '@models/common.model';
import { gql } from 'apollo-angular';

import { Plan } from '../../../models/common.model';

export const GetBuildQuery = gql`
  query($buildId: String!) {
    build(buildUUID: $buildId) {
      id
      revision
      buildToStatus {
        completed
        failed
        state
      }
      buildToTeam {
        id
        team_number
        TeamToStatus {
          completed
          failed
          state
        }
        TeamToPlan {
          id
        }
        TeamToProvisionedNetwork {
          id
          name
          cidr
          ProvisionedNetworkToStatus {
            completed
            failed
            state
          }
          ProvisionedNetworkToPlan {
            id
          }
          ProvisionedNetworkToNetwork {
            id
            vdi_visible
            vars {
              key
              value
            }
            tags {
              key
              value
            }
          }
          ProvisionedNetworkToProvisionedHost {
            id
            subnet_ip
            ProvisionedHostToStatus {
              completed
              failed
              state
            }
            ProvisionedHostToPlan {
              id
            }
            ProvisionedHostToHost {
              id
              hostname
              description
              OS
              allow_mac_changes
              exposed_tcp_ports
              exposed_udp_ports
              user_groups
              override_password
              vars {
                key
                value
              }
              tags {
                key
                value
              }
            }
            ProvisionedHostToProvisioningStep {
              id
              type
              step_number
              ProvisioningStepToStatus {
                completed
                failed
                state
              }
              ProvisioningStepToScript {
                id
                name
                language
                description
                source
                source_type
                disabled
                args
                vars {
                  key
                  value
                }
                tags {
                  key
                  value
                }
              }
              ProvisioningStepToCommand {
                id
                name
                description
                program
                args
                disabled
                vars {
                  key
                  value
                }
                tags {
                  key
                  value
                }
              }
              ProvisioningStepToDNSRecord {
                id
                name
                values
                type
                zone
                disabled
                vars {
                  key
                  value
                }
                tags {
                  key
                  value
                }
              }
              ProvisioningStepToFileDownload {
                id
                source
                sourceType
                destination
                disabled
                tags {
                  key
                  value
                }
              }
              ProvisioningStepToFileDelete {
                id
                path
                tags {
                  key
                  value
                }
              }
              ProvisioningStepToFileExtract {
                id
                source
                destination
                type
                tags {
                  key
                  value
                }
              }
            }
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

export interface GetBuildPlansVars {
  buildId: ID;
}

export interface GetBuildPlansData {
  build: {
    id: ID;
    buildToTeam: {
      id: ID;
      TeamToPlan: Plan;
      TeamToProvisionedNetwork: {
        id: ID;
        ProvisionedNetworkToPlan: Plan;
        ProvisionedNetworkToProvisionedHost: {
          id: ID;
          ProvisionedHostToPlan: Plan;
          ProvisionedHostToProvisioningStep: {
            id: ID;
            ProvisioningStepToPlan: Plan;
          }[];
        }[];
      }[];
    }[];
  };
}

const GetBuildPlansQuery = gql`
  query($buildId: String!) {
    build(buildUUID: $buildId) {
      id
      buildToTeam {
        id
        TeamToPlan {
          ...PlanFields
        }
        TeamToProvisionedNetwork {
          id
          ProvisionedNetworkToPlan {
            ...PlanFields
          }
          ProvisionedNetworkToProvisionedHost {
            id
            ProvisionedHostToPlan {
              ...PlanFields
            }
            ProvisionedHostToProvisioningStep {
              id
              ProvisioningStepToPlan {
                ...PlanFields
              }
            }
          }
        }
      }
    }
  }
`;
