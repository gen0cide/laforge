import { AgentStatus } from './agent.model';
import { ID } from './common.model';
import { Environment } from './environment.model';
import { ProvisionedStep } from './host.model';

export interface AgentStatusQueryResult {
  environment: {
    build: {
      teams: [
        {
          id: string;
          provisionedNetworks: [
            {
              id: string;
              provisionedHosts: [
                {
                  id: string;
                  heartbeat?: AgentStatus;
                }
              ];
            }
          ];
        }
      ];
    };
  };
}

export interface EnvironmentQueryResult {
  environment: Environment;
}

export interface HostStepsQueryResult {
  provisionedHost: {
    id: ID;
    provisionedSteps: ProvisionedStep[];
  };
}

export interface EnvironmentInfo {
  id: ID;
  Name: string;
  CompetitionID: string;
}

export interface EnvironmentInfoQueryResult {
  environments: EnvironmentInfo[];
}
