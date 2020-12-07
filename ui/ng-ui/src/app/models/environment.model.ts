import { ID, Tag, configMap, User, Team } from './common.model';
import { DNS } from './dns.model';
import { Host } from './host.model';
import { Network } from './network.model';

interface Build {
  id: ID;
  revision: number;
  tags: Tag[];
  config: configMap[];
  maintainer: User;
  teams: Team[];
}

interface Competition {
  id: ID;
  rootPassword: string;
  config: configMap[];
  dns: DNS;
}

interface Environment {
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

export { Environment, Build, Competition };
