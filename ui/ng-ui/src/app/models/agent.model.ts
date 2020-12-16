import { AgentStatusQueryResult, Team } from './common.model';
import { Environment } from './environment.model';
import { ProvisionedNetwork } from './network.model';
import { ProvisionedHost } from 'src/app/models/host.model';

interface AgentStatus {
  clientId: string;
  hostname: string;
  upTime: number;
  bootTime: number;
  numProcs: number;
  OS: string;
  hostID: string;
  load1?: number;
  load5?: number;
  load15?: number;
  totalMem: number;
  freeMem: number;
  usedMem: number;
}

/* eslint-disable @typescript-eslint/no-explicit-any */
function updateAgentStatuses(environment: Environment, statusQueryResult: AgentStatusQueryResult): Environment {
  return {
    ...environment,
    build: {
      ...environment.build,
      teams: environment.build.teams.map(
        (team: Team) =>
          ({
            ...team,
            provisionedNetworks: team.provisionedNetworks.map((provNetwork: ProvisionedNetwork) => ({
              ...provNetwork,
              provisionedHosts: provNetwork.provisionedHosts.map((provHost: ProvisionedHost) => {
                const matchedHost = statusQueryResult.environment.build.teams
                  .filter((t: any) => t.id === team.id)[0]
                  .provisionedNetworks.filter((pN: any) => pN.id === provNetwork.id)[0]
                  .provisionedHosts.filter((pH: any) => pH.id === provHost.id)[0];
                return {
                  ...provHost,
                  heartbeat: matchedHost ? { ...matchedHost.heartbeat } : null
                };
              })
            }))
          } as Team)
      )
    }
  };
}

export { AgentStatus, updateAgentStatuses };
