import { AgentStatus } from './agent.model';
import { Environment } from './environment.model';

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
                  heartbeat: AgentStatus;
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
