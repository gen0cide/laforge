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
          }
        }
      }
    }
  }
`;

export { getEnvConfigQuery };
