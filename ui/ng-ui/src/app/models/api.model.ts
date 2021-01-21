import { AgentStatus } from './agent.model';
import { ID } from './common.model';
import { Environment } from './environment.model';
import { ProvisionedStep } from './host.model';

export interface AgentStatusQueryResult {
  environment: {
    build: {
      teams: [
        {
          id: ID;
          provisionedNetworks: [
            {
              id: ID;
              provisionedHosts: [
                {
                  id: ID;
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
