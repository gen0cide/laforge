import { ID, Tag, configMap, User, Team, ProvisionStatus } from './common.model';
import { DNS } from './dns.model';
import { Host, ProvisionedHost, ProvisionedStep } from './host.model';
import { Network, ProvisionedNetwork } from './network.model';

export interface Build {
  id: ID;
  revision: number;
  tags: Tag[];
  config: configMap[];
  maintainer: User;
  teams: Team[];
}

export interface Competition {
  id: ID;
  rootPassword: string;
  config: configMap[];
  dns: DNS;
}

export interface Environment {
  id: ID;
  CompetitionID: string;
  Name: string;
  Description: string;
  Builder: string;
  TeamCount: number;
  AdminCIDRs: string[];
  ExposedVDIPorts: string[];
  tags: Tag[];
  config: configMap[];
  maintainer: User;
  networks: Network[];
  hosts: Host[];
  build: Build;
  competition: Competition;
}

/* eslint-disable-next-line @typescript-eslint/explicit-module-boundary-types, @typescript-eslint/no-explicit-any */
export function resolveStatuses(environment: any): any {
  return {
    ...environment,
    build: {
      ...environment.build,
      teams: [
        ...environment.build.teams.map((team: Team) => ({
          ...team,
          provisionedNetworks: team.provisionedNetworks.map((provisionedNetwork: ProvisionedNetwork) => ({
            ...provisionedNetwork,
            status: {
              ...provisionedNetwork.status,
              state: ProvisionStatus[provisionedNetwork.status.state]
            },
            provisionedHosts: provisionedNetwork.provisionedHosts.map((provisionedHost: ProvisionedHost) => ({
              ...provisionedHost,
              status: {
                ...provisionedHost.status,
                state: ProvisionStatus[provisionedHost.status.state]
              },
              provisionedSteps: provisionedHost.provisionedSteps.map((provisionedStep: ProvisionedStep) => ({
                ...provisionedStep,
                status: {
                  ...provisionedStep.status,
                  state: ProvisionStatus[provisionedStep.status.state]
                }
              }))
            }))
          }))
        }))
      ]
    }
  };
}
