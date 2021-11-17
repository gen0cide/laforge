import { ProvisionedHost } from 'src/app/models/host.model';

import { AgentStatusQueryResult } from './api.model';
import { Team } from './common.model';
import { Build, Environment } from './environment.model';
import { ProvisionedNetwork } from './network.model';

export enum AgentCommand {
  DEFAULT,
  DELETE,
  REBOOT,
  EXTRACT,
  DOWNLOAD,
  CREATEUSER,
  CREATEUSERPASS,
  ADDTOGROUP,
  EXECUTE,
  VALIDATE
}
export interface AgentStatus {
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
  timestamp: number;
}

/* eslint-disable @typescript-eslint/no-explicit-any */
export function updateEnvAgentStatuses(environment: Environment, statusQueryResult: AgentStatusQueryResult): Environment {
  return {
    ...environment,
    EnvironmentToBuild: environment.EnvironmentToBuild
      ? environment.EnvironmentToBuild.map(
          (build: any): Build => ({
            ...build,
            buildToTeam: build.buildToTeam
              ? build.buildToTeam.map(
                  (team: any): Team => ({
                    ...team,
                    TeamToProvisionedNetwork: team.TeamToProvisionedNetwork
                      ? team.TeamToProvisionedNetwork.map(
                          (provisionedNetwork: any): ProvisionedNetwork => ({
                            ...provisionedNetwork,
                            ProvisionedNetworkToProvisionedHost: provisionedNetwork.ProvisionedNetworkToProvisionedHost
                              ? provisionedNetwork.ProvisionedNetworkToProvisionedHost.map(
                                  (provisionedHost: any): ProvisionedHost => {
                                    const matchedHost = statusQueryResult.build.buildToTeam
                                      .filter((t: any) => t.id === team.id)[0]
                                      .TeamToProvisionedNetwork.filter((pN: any) => pN.id === provisionedNetwork.id)[0]
                                      .ProvisionedNetworkToProvisionedHost.filter((pH: any) => pH.id === provisionedHost.id)[0];
                                    // if (matchedHost != null) {
                                    //   console.log(matchedHost.id, matchedHost.ProvisionedHostToAgentStatus);
                                    // }
                                    return {
                                      ...provisionedHost,
                                      ProvisionedHostToAgentStatus: { ...matchedHost.ProvisionedHostToAgentStatus }
                                    };
                                  }
                                )
                              : null
                          })
                        )
                      : null
                  })
                )
              : null
          })
        )
      : null
  };
}

/* eslint-disable @typescript-eslint/no-explicit-any */
export function updateBuildAgentStatuses(build: Build, statusQueryResult: AgentStatusQueryResult): Build {
  return {
    ...build,
    buildToTeam: build.buildToTeam
      ? build.buildToTeam.map(
          (team: any): Team => ({
            ...team,
            TeamToProvisionedNetwork: team.TeamToProvisionedNetwork
              ? team.TeamToProvisionedNetwork.map(
                  (provisionedNetwork: any): ProvisionedNetwork => ({
                    ...provisionedNetwork,
                    ProvisionedNetworkToProvisionedHost: provisionedNetwork.ProvisionedNetworkToProvisionedHost
                      ? provisionedNetwork.ProvisionedNetworkToProvisionedHost.map(
                          (provisionedHost: any): ProvisionedHost => {
                            const matchedHost = statusQueryResult.build.buildToTeam
                              .filter((t: any) => t.id === team.id)[0]
                              .TeamToProvisionedNetwork.filter((pN: any) => pN.id === provisionedNetwork.id)[0]
                              .ProvisionedNetworkToProvisionedHost.filter((pH: any) => pH.id === provisionedHost.id)[0];
                            // if (matchedHost != null) {
                            //   console.log(matchedHost.id, matchedHost.ProvisionedHostToAgentStatus);
                            // }
                            return {
                              ...provisionedHost,
                              ProvisionedHostToAgentStatus: { ...matchedHost.ProvisionedHostToAgentStatus }
                            };
                          }
                        )
                      : null
                  })
                )
              : null
          })
        )
      : null
  };
}
