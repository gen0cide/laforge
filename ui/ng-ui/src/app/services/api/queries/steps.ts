import { gql } from 'apollo-angular';

const getProvisionedSteps = gql`
  query($hostId: String!) {
    provisionedHost(proHostUUID: $hostId) {
      id
      ProvisionedHostToProvisioningStep {
        id
        type
        ProvisioningStepToStatus {
          state
          started_at
          ended_at
          error
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
`;

export { getProvisionedSteps };
