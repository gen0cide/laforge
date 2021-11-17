import { gql } from 'apollo-angular';
import { DocumentNode } from 'graphql';
import { ID } from 'src/app/models/common.model';

const BAK_getEnvironmentQuery = (id: ID): DocumentNode => gql`
{
  environment(envUUID: "${id}") {
    id
    CompetitionID,
    Name,
    Description,
    Builder,
    TeamCount,
    AdminCIDRs,
    ExposedVDIPorts,
    tags {
      id,
      name,
      description
    },
    config {
      key,
      value
    },
    maintainer {
      id,
      name,
      uuid,
      email
    },
    build {
      id,
      revision,
      tags {
        id,
        name,
        description
      },
      config {
        key,
        value
      },
      maintainer {
        id,
        name,
        uuid,
        email
      },
      teams {
        id,
        teamNumber,
        provisionedNetworks {
          id,
          name,
          cidr,
          network {
            id,
            vdiVisible
          },
          vars {
            key,
            value
          },
          tags {
            id,
            name,
            description
          },
          provisionedHosts {
            id,
            subnetIP,
            status {
              state,
              startedAt,
              endedAt,
              error
            },
            host {
              id,
              hostname,
              OS,
              allowMacChanges,
              exposedTCPPorts,
              exposedUDPPorts,
              userGroups,
              overridePassword,
              maintainer {
                name,
                email
              },
              vars {
                key,
                value
              },
              tags {
                name,
                description
              },
            },
            provisionedSteps {
              id,
              provisionType,
              status {
                state,
                startedAt,
                endedAt,
                error
              },
              script {
                id,
                name,
                description,
                source,
                sourceType,
                disabled
              },
              command {
                id,
                name,
                description,
                args,
                disabled
              },
              DNSRecord {
                id,
                name,
                values,
                type,
                zone,
                disabled
              },
              fileDownload {
                id,
                source,
                sourceType,
                destination,
                mode,
                disabled,
              },
              fileDelete {
                id
                path,
              },
              fileExtract {
                id,
                source,
                destination,
                type
              }
            }
          },
          status {
            state,
            startedAt,
            endedAt,
            error
          },
        }
      }
    }
  }
}
`;

const getEnvironmentQuery = (id: ID): DocumentNode => gql`
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
        TeamToProvisionedNetwork {
          id
          name
          cidr
          ProvisionedNetworkToStatus {
            completed
            failed
            state
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
          }
        }
      }
    }
  }
}
`;

const getEnvironmentsQuery = (): DocumentNode => gql`
  {
    environments {
      id
      name
      competition_id
      EnvironmentToBuild {
        id
        revision
      }
    }
  }
`;

export { getEnvironmentQuery, getEnvironmentsQuery };
