import { GetBuildPlansData } from '../services/api/queries/build';

import { Command } from './command.model';
import { ID, configMap, User, Team, ProvisionStatus, Status, Plan, tagMap } from './common.model';
import { DNS, DNSRecord } from './dns.model';
import { FileDelete, FileDownload, FileExtract } from './file.model';
import { Host, ProvisionedHost } from './host.model';
import { Network, ProvisionedNetwork } from './network.model';
import { Script } from './script.model';
import { ProvisioningStep, ProvisioningStepType } from './step.model';
import { Identity } from './user.model';

export interface Build {
  id: ID;
  revision: number;
  completed_pla?: boolean;
  buildToStatus?: Status;
  buildToEnvironment?: Environment;
  buildToCompetition?: Competition;
  buildToProvisionedNetwork?: ProvisionedNetwork;
  buildToTeam?: Team[];
  buildToPlan?: Plan[];
}

export interface Competition {
  id: ID;
  root_password: string;
  hcl_id?: string;
  config?: configMap[];
  tags?: tagMap[];
  competitionToDns?: DNS;
  CompetitionToEnvironment?: Environment;
  CompetitionToBuild?: Build[];
}

export interface Environment {
  id: ID;
  hcl_id?: string;
  competition_id: string;
  name: string;
  description: string;
  builder: string;
  team_count: number;
  admin_cidrs: string[];
  exposed_vdi_ports: string[];
  config: configMap[];
  tags: tagMap[];
  EnvironmentToUser?: User[];
  EnvironmentToHost?: Host[];
  EnvironmentToCompetition?: Competition[];
  EnvironmentToIdentity?: Identity[];
  EnvironmentToCommand?: Command[];
  EnvironmentToScript?: Script[];
  EnvironmentToFileDownload?: FileDownload[];
  EnvironmentToFileDelete?: FileDelete[];
  EnvironmentToFileExtract?: FileExtract[];
  EnvironmentToDNSRecord?: DNSRecord[];
  EnvironmentToNetwork?: Network[];
  EnvironmentToBuild?: Build[];
}

/* eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types, @typescript-eslint/no-explicit-any */
export function resolveEnvEnums(environment: any): Environment {
  return {
    ...environment,
    EnvironmentToBuild: environment.EnvironmentToBuild
      ? environment.EnvironmentToBuild.map(
          (build: any): Build => ({
            ...build,
            buildToStatus: {
              ...build.buildToStatus,
              state: ProvisionStatus[build.buildToStatus.state]
            },
            buildToTeam: build.buildToTeam
              ? build.buildToTeam.map(
                  (team: any): Team => ({
                    ...team,
                    TeamToStatus: {
                      ...team.TeamToStatus,
                      state: ProvisionStatus[team.TeamToStatus.state]
                    },
                    TeamToProvisionedNetwork: team.TeamToProvisionedNetwork
                      ? team.TeamToProvisionedNetwork.map(
                          (provisionedNetwork: any): ProvisionedNetwork => ({
                            ...provisionedNetwork,
                            ProvisionedNetworkToStatus: {
                              ...provisionedNetwork.ProvisionedNetworkToStatus,
                              state: ProvisionStatus[provisionedNetwork.ProvisionedNetworkToStatus.state]
                            },
                            ProvisionedNetworkToProvisionedHost: provisionedNetwork.ProvisionedNetworkToProvisionedHost
                              ? provisionedNetwork.ProvisionedNetworkToProvisionedHost.map(
                                  (provisionedHost: any): ProvisionedHost => ({
                                    ...provisionedHost,
                                    ProvisionedHostToStatus: {
                                      ...provisionedHost.ProvisionedHostToStatus,
                                      state: ProvisionStatus[provisionedHost.ProvisionedHostToStatus.state]
                                    },
                                    ProvisionedHostToProvisioningStep: provisionedHost.ProvisionedHostToProvisioningStep
                                      ? provisionedHost.ProvisionedHostToProvisioningStep.map(
                                          (provisioningStep: any): ProvisioningStep => ({
                                            ...provisioningStep,
                                            type: ProvisioningStepType[provisioningStep.type],
                                            ProvisioningStepToStatus: {
                                              ...provisioningStep.ProvisioningStepToStatus,
                                              state: ProvisionStatus[provisioningStep.ProvisioningStepToStatus.state]
                                            }
                                          })
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
          })
        )
      : null
  };
}

/* eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types, @typescript-eslint/no-explicit-any */
export function resolveBuildEnums(build: any): Build {
  return {
    ...build,
    buildToStatus: build.buildToStatus
      ? {
          ...build.buildToStatus,
          state: typeof build.buildToStatus.state === 'string' ? ProvisionStatus[build.buildToStatus.state] : build.buildToStatus.state
        }
      : null,
    buildToTeam: build.buildToTeam
      ? build.buildToTeam.map(
          (team: any): Team => ({
            ...team,
            TeamToStatus: team.TeamToStatus
              ? {
                  ...team.TeamToStatus,
                  state: typeof team.TeamToStatus.state === 'string' ? ProvisionStatus[team.TeamToStatus.state] : team.TeamToStatus.state
                }
              : null,
            TeamToPlan: team.TeamToPlan
              ? {
                  ...team.TeamToPlan,
                  PlanToStatus: team.TeamToPlan.PlanToStatus
                    ? {
                        ...team.TeamToPlan.PlanToStatus,
                        state:
                          typeof team.TeamToPlan.PlanToStatus.state === 'string'
                            ? ProvisionStatus[team.TeamToPlan.PlanToStatus.state]
                            : team.TeamToPlan.PlanToStatus.state
                      }
                    : null
                }
              : null,
            TeamToProvisionedNetwork: team.TeamToProvisionedNetwork
              ? team.TeamToProvisionedNetwork.map(
                  (provisionedNetwork: any): ProvisionedNetwork => ({
                    ...provisionedNetwork,
                    ProvisionedNetworkToStatus: provisionedNetwork.ProvisionedNetworkToStatus
                      ? {
                          ...provisionedNetwork.ProvisionedNetworkToStatus,
                          state:
                            typeof provisionedNetwork.ProvisionedNetworkToStatus.state === 'string'
                              ? ProvisionStatus[provisionedNetwork.ProvisionedNetworkToStatus.state]
                              : provisionedNetwork.ProvisionedNetworkToStatus.state
                        }
                      : null,
                    ProvisionedNetworkToPlan: provisionedNetwork.ProvisionedNetworkToPlan
                      ? {
                          ...provisionedNetwork.ProvisionedNetworkToPlan,
                          PlanToStatus: provisionedNetwork.ProvisionedNetworkToPlan.PlanToStatus
                            ? {
                                ...provisionedNetwork.ProvisionedNetworkToPlan.PlanToStatus,
                                state:
                                  typeof provisionedNetwork.ProvisionedNetworkToPlan.PlanToStatus.state === 'string'
                                    ? ProvisionStatus[provisionedNetwork.ProvisionedNetworkToPlan.PlanToStatus.state]
                                    : provisionedNetwork.ProvisionedNetworkToPlan.PlanToStatus.state
                              }
                            : null
                        }
                      : null,
                    ProvisionedNetworkToProvisionedHost: provisionedNetwork.ProvisionedNetworkToProvisionedHost
                      ? provisionedNetwork.ProvisionedNetworkToProvisionedHost.map(
                          (provisionedHost: any): ProvisionedHost => ({
                            ...provisionedHost,
                            ProvisionedHostToStatus: provisionedHost.ProvisionedHostToStatus
                              ? {
                                  ...provisionedHost.ProvisionedHostToStatus,
                                  PlanToStatus: provisionedHost.ProvisionedHostToStatus.PlanToStatus
                                    ? {
                                        ...provisionedHost.ProvisionedHostToStatus.PlanToStatus,
                                        state:
                                          typeof provisionedHost.ProvisionedHostToStatus.PlanToStatus.state === 'string'
                                            ? ProvisionStatus[provisionedHost.ProvisionedHostToStatus.PlanToStatus.state]
                                            : provisionedHost.ProvisionedHostToStatus.PlanToStatus.state
                                      }
                                    : null
                                }
                              : null,
                            ProvisionedHostToPlan: provisionedHost.ProvisionedHostToPlan
                              ? {
                                  ...provisionedHost.ProvisionedHostToPlan,
                                  PlanToStatus: provisionedHost.ProvisionedHostToStatus.PlanToStatus
                                    ? {
                                        ...provisionedHost.ProvisionedHostToStatus.PlanToStatus,
                                        state:
                                          typeof provisionedHost.ProvisionedHostToStatus.PlanToStatus.state === 'string'
                                            ? ProvisionStatus[provisionedHost.ProvisionedHostToStatus.PlanToStatus.state]
                                            : provisionedHost.ProvisionedHostToStatus.PlanToStatus.state
                                      }
                                    : null
                                }
                              : null,
                            ProvisionedHostToProvisioningStep: provisionedHost.ProvisionedHostToProvisioningStep
                              ? provisionedHost.ProvisionedHostToProvisioningStep.map(
                                  (provisioningStep: any): ProvisioningStep => ({
                                    ...provisioningStep,
                                    type:
                                      typeof provisioningStep.type === 'string'
                                        ? ProvisioningStepType[provisioningStep.type]
                                        : provisioningStep.type,
                                    ProvisioningStepToStatus: provisioningStep.ProvisioningStepToStatus
                                      ? {
                                          ...provisioningStep.ProvisioningStepToStatus,
                                          state:
                                            typeof provisioningStep.ProvisioningStepToStatus.state === 'string'
                                              ? ProvisionStatus[provisioningStep.ProvisioningStepToStatus.state]
                                              : provisioningStep.ProvisioningStepToStatus.state
                                        }
                                      : null,
                                    ProvisioningStepToPlan: provisioningStep.ProvisioningStepToPlan
                                      ? {
                                          ...provisioningStep.ProvisioningStepToPlan,
                                          PlanToStatus: provisioningStep.ProvisioningStepToPlan.PlanToStatus
                                            ? {
                                                ...provisioningStep.ProvisioningStepToPlan.PlanToStatus,
                                                state:
                                                  typeof provisioningStep.ProvisioningStepToPlan.PlanToStatus.state === 'string'
                                                    ? ProvisionStatus[provisioningStep.ProvisioningStepToPlan.PlanToStatus.state]
                                                    : provisioningStep.ProvisioningStepToPlan.PlanToStatus.state
                                              }
                                            : null
                                        }
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
          })
        )
      : null
  };
}

export function updateBuildPlans(build: Build, buildPlansData: GetBuildPlansData): Build {
  return {
    ...build,
    buildToTeam: build.buildToTeam.map((team) => {
      const updatedTeam = buildPlansData.build.buildToTeam.find((t) => t.id === team.id);
      return {
        ...team,
        TeamToPlan: updatedTeam?.TeamToPlan ?? undefined,
        TeamToProvisionedNetwork: team.TeamToProvisionedNetwork.map((pnet) => {
          const updatedProvisionedNetwork = updatedTeam.TeamToProvisionedNetwork.find((n) => n.id === pnet.id);
          return {
            ...pnet,
            ProvisionedNetworkToPlan: updatedProvisionedNetwork?.ProvisionedNetworkToPlan ?? undefined,
            ProvisionedNetworkToProvisionedHost: pnet.ProvisionedNetworkToProvisionedHost.map((phost) => {
              const updatedProvisionedHost = updatedProvisionedNetwork.ProvisionedNetworkToProvisionedHost.find((h) => h.id === phost.id);
              return {
                ...phost,
                ProvisionedHostToPlan: updatedProvisionedHost?.ProvisionedHostToPlan ?? undefined,
                ProvisionedHostToProvisioningStep: phost.ProvisionedHostToProvisioningStep.map((pstep) => {
                  const updatedProvisioningStep = updatedProvisionedHost.ProvisionedHostToProvisioningStep.find((s) => s.id === pstep.id);
                  return {
                    ...pstep,
                    ProvisioningStepToPlan: updatedProvisioningStep?.ProvisioningStepToPlan ?? undefined
                  };
                })
              };
            })
          };
        })
      };
    })
  };
}
