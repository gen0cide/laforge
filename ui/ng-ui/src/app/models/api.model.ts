import { Build } from 'src/app/models/environment.model';

import { AgentStatus } from './agent.model';
import { ID } from './common.model';
import { Environment } from './environment.model';
import { ProvisioningStep } from './step.model';

/*
  query($id: String!) {
    environment(envUUID: $id) {
      EnvironmentToBuild {
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
  }
*/
export interface AgentStatusQueryResult {
  build: {
    id: ID;
    buildToTeam: [
      {
        id: ID;
        TeamToProvisionedNetwork: [
          {
            id: ID;
            ProvisionedNetworkToProvisionedHost: [
              {
                id: ID;
                ProvisionedHostToAgentStatus?: AgentStatus;
              }
            ];
          }
        ];
      }
    ];
  };
}

export interface EnvironmentQueryResult {
  environment: Environment;
}

export interface HostStepsQueryResult {
  provisionedHost: {
    id: ID;
    ProvisionedHostToProvisioningStep: ProvisioningStep[];
  };
}

export interface EnvironmentInfo {
  id: ID;
  name: string;
  competition_id: string;
  EnvironmentToBuild: {
    id: ID;
    revision: number;
  };
}

export interface EnvironmentInfoQueryResult {
  environments: EnvironmentInfo[];
}

export interface BuildQueryVars {
  buildId: ID;
}

export interface BuildQueryResult {
  build: Build;
}
