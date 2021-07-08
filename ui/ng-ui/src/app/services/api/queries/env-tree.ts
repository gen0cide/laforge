import { gql } from 'apollo-angular';
import { DocumentNode } from 'graphql';
import { ID } from 'src/app/models/common.model';

const BAK_getEnvConfigQuery = (id: ID): DocumentNode => gql`
  {
    environment(envUUID: "${id}") {
      id
      CompetitionID
      Name
      Description
      Builder
      TeamCount
      AdminCIDRs
      ExposedVDIPorts
      tags {
        id
        name
        description
      }
      config {
        key
        value
      }
      maintainer {
        id
        name
        uuid
        email
      }
      build {
        id
        revision
        tags {
          id
          name
          description
        }
        config {
          key
          value
        }
        maintainer {
          id
          name
          uuid
          email
        }
        teams {
          id
          teamNumber
          provisionedNetworks {
            id
            name
            cidr
            network {
              id
              vdiVisible
            }
            provisionedHosts {
              id
              subnetIP
              status {
                state
                startedAt
                endedAt
                error
              }
              host {
                id
                hostname
                OS
                allowMacChanges
                exposedTCPPorts
                exposedUDPPorts
                userGroups
                overridePassword
                maintainer {
                  name
                  email
                }
                vars {
                  key
                  value
                }
                tags {
                  name
                  description
                }
              }
            }
            status {
              state
              startedAt
              endedAt
              error
            }
          }
        }
      }
    }
  }
`;

const getEnvConfigQuery = (id: ID): DocumentNode => gql`
  {
    environment(envUUID: "${id}") {
      id
      competition_id
      name
      description
      builder
      team_count
      admin_cidrs
      exposed_vdi_ports
      tags {
        key
        value
      }
      config {
        key
        value
      }
      EnvironmentToUser {
        id
        name
        uuid
        email
      }
      EnvironmentToBuild {
        id
        revision
        buildToTeam {
          id
          team_number
          TeamToProvisionedNetwork {
            id
            name
            cidr
            ProvisionedNetworkToNetwork {
              id
              vdi_visible
            }
            ProvisionedNetworkToProvisionedHost {
              id
              subnet_ip
              ProvisionedHostToHost {
                id
                hostname
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
            }
          }
        }
      }
    }
  }
`;

export { getEnvConfigQuery };
