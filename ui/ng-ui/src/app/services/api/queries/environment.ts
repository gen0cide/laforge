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
    environmentToUser {
      id,
      name,
      uuid,
      email
    },
    environmentToBuild {
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
      buildToUser {
        id,
        name,
        uuid,
        email
      },
      buildToTeam {
        id,
        teamNumber,
        teamToProvisionedNetwork {
          id,
          name,
          cidr,
          provisionedNetworkToNetwork {
            id,
            vdiVisible
            vars {
              key,
              value
            },
            networkToTag {
              id,
              name,
              description
            },
          },
          provisionedNetworkToProvisionedHost {
            id,
            subnetIP,
            provisionedHostToHost {
              id,
              hostname,
              OS,
              allowMacChanges,
              exposedTCPPorts,
              exposedUDPPorts,
              userGroups,
              overridePassword,
              hostToUser {
                name,
                email
              },
              vars {
                key,
                value
              },
            	hostToTag {
                name,
                description
              },
            },
            provisionedHostToProvisioningStep {
              id,
              provisionType,
              provisioningStepToScript {
                id,
                name,
                description,
                source,
                sourceType,
                disabled
              },
              provisioningStepToCommand {
                id,
                name,
                description,
                args,
                disabled
              },
              provisioningStepToDNSRecord {
                id,
                name,
                values,
                type,
                zone,
                disabled
              },
              provisioningStepToFileDownload {
                id,
                source,
                sourceType,
                destination,
                disabled,
              },
              provisioningStepToFileDelete {
                id
                path,
              },
              provisioningStepToFileExtract {
                id,
                source,
                destination,
                type
              }
            }
          },
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
      Name
      CompetitionID
    }
  }
`;

export { getEnvironmentQuery, getEnvironmentsQuery };
